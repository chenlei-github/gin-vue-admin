# 重复键冲突问题修复总结

## 📋 问题概述

在使用 **全量替换模式** 导入市场数据时，可能会遇到重复键冲突错误：

```
Error 1062 (23000): Duplicate entry 'XXX' for key 'bt_products.uni_bt_products_asin'
```

这个错误有 **两个不同的根本原因**，需要分别处理。

---

## 🔴 问题 1：软删除导致的重复键冲突（已修复）

### 问题描述

全量替换模式下，即使选择了"替换现有数据"，仍然报重复键错误。

### 根本原因

- 模型使用了 **软删除**（`DeletedAt` 字段）
- 删除操作只是设置 `deleted_at` 字段，**数据仍在数据库中**
- 唯一索引 **不包含** `deleted_at` 字段
- 插入新数据时，与软删除的旧数据冲突

### 解决方案

在全量替换模式下，使用 **物理删除** 而非软删除：

```go
// 修复前：软删除
tx.Delete(&brandtrekin.BtProduct{})

// 修复后：物理删除
tx.Unscoped().Delete(&brandtrekin.BtProduct{})
```

### 详细文档

👉 [FIX_SOFT_DELETE_CONFLICT.md](./FIX_SOFT_DELETE_CONFLICT.md)

### 修改的文件

- [bt_import.go](/Users/leon/code/ai-code/gin-vue-admin/server/service/brandtrekin/bt_import.go) - `deleteMarketData()` 方法

---

## 🟡 问题 2：CSV 文件重复数据导致的冲突（已修复）

### 问题描述

CSV 文件中存在重复的记录（如同一个 ASIN 出现多次），导致批量插入时冲突。

### 根本原因

- CSV 文件中有重复数据
- 分批处理时，第一批创建了某个 ASIN
- 第二批又遇到相同的 ASIN，尝试再次创建，导致冲突

### 解决方案

在批量保存前，先对数据进行 **去重处理**：

```go
// 使用 map 自动去重
productMap := make(map[string]ProductInfo)
for _, p := range products {
    productMap[p.ASIN] = p  // 相同 ASIN 会被覆盖
}

// 转换回切片
uniqueProducts := make([]ProductInfo, 0, len(productMap))
for _, p := range productMap {
    uniqueProducts = append(uniqueProducts, p)
}
```

### 修改的文件

- [bt_import_optimized.go](/Users/leon/code/ai-code/gin-vue-admin/server/service/brandtrekin/bt_import_optimized.go)
  - `SaveProductDataBatch()` - 商品数据去重
  - `SaveBrandSocialDataBatch()` - 品牌数据去重
  - `SaveKeywordDataBatch()` - 关键词数据去重
  - `SaveProductSalesDataBatch()` - 销售数据去重

---

## 🎯 如何判断是哪个问题？

### 判断方法

1. **查看错误日志**：
   - 如果是第一次导入就报错 → 可能是 **问题 2**（CSV 重复数据）
   - 如果是第二次导入才报错 → 可能是 **问题 1**（软删除冲突）

2. **检查数据库**：
   ```sql
   -- 查看是否有软删除的数据
   SELECT COUNT(*) FROM bt_products WHERE deleted_at IS NOT NULL;
   
   -- 如果有数据，说明是问题 1
   ```

3. **检查 CSV 文件**：
   ```bash
   # 检查是否有重复的 ASIN
   cut -d',' -f1 Products.csv | sort | uniq -d
   
   # 如果有输出，说明是问题 2
   ```

---

## ✅ 修复效果

### 修复前

```
❌ 全量替换失败
Error 1062: Duplicate entry 'B09ZTYGGTJ' for key 'bt_products.uni_bt_products_asin'
```

### 修复后

```
✅ 全量替换成功
- 旧数据被物理删除
- CSV 重复数据自动去重
- 新数据成功导入
```

---

## 📊 涉及的唯一索引

| 表名 | 唯一索引 | 索引字段 |
|------|---------|---------|
| bt_products | uni_bt_products_asin | asin |
| bt_brands | uni_bt_brands_market_brand | market_id + brand_name |
| bt_keywords | uni_bt_keywords_market_keyword_source | market_id + keyword + source |
| bt_product_monthly_sales | uni_bt_product_monthly_sales_asin_date | asin + date |

---

## 🧪 测试验证

### 1. 测试全量替换（验证问题 1 已修复）

```bash
# 第一次导入
curl -X POST http://localhost:8888/api/v1/brandtrekin/import \
  -F "marketId=1" \
  -F "replaceMode=true" \
  -F "files=@Products.csv"

# 第二次导入（全量替换）
curl -X POST http://localhost:8888/api/v1/brandtrekin/import \
  -F "marketId=1" \
  -F "replaceMode=true" \
  -F "files=@Products.csv"

# 预期：✅ 成功，不报错
```

### 2. 测试重复数据（验证问题 2 已修复）

```bash
# 创建包含重复 ASIN 的 CSV 文件
cat > test_duplicate.csv << EOF
ASIN,Title,Price
B09ZTYGGTJ,Product 1,29.99
B09ZTYGGTJ,Product 1 Duplicate,29.99
B0CMTJ6CZC,Product 2,39.99
EOF

# 导入
curl -X POST http://localhost:8888/api/v1/brandtrekin/import \
  -F "marketId=1" \
  -F "replaceMode=true" \
  -F "files=@test_duplicate.csv"

# 预期：✅ 成功，自动去重，只保留最后一个
```

### 3. 验证数据完整性

```sql
-- 检查是否有软删除的数据
SELECT COUNT(*) FROM bt_products WHERE deleted_at IS NOT NULL;
-- 预期：0（全量替换后，旧数据被物理删除）

-- 检查是否有重复的 ASIN
SELECT asin, COUNT(*) as count 
FROM bt_products 
WHERE deleted_at IS NULL 
GROUP BY asin 
HAVING count > 1;
-- 预期：空结果（没有重复）
```

---

## 🎉 总结

| 问题 | 原因 | 解决方案 | 状态 |
|------|------|---------|------|
| 软删除冲突 | 唯一索引不包含 deleted_at | 使用 Unscoped() 物理删除 | ✅ 已修复 |
| CSV 重复数据 | 文件中有重复记录 | 批量保存前去重 | ✅ 已修复 |

两个问题都已修复，现在可以正常使用全量替换模式导入数据了！🎉
