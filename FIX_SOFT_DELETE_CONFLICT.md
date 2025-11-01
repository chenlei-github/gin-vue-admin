# 修复：软删除导致的全量替换重复键冲突问题

## 🐛 问题描述

在使用 **全量替换模式** 导入市场数据时，出现以下错误：

```
批量导入失败! {"error": "保存商品数据失败: 批量创建商品失败: Error 1062 (23000): Duplicate entry 'B09ZTYGGTJ' for key 'bt_products.uni_bt_products_asin'"}
```

## 🔍 问题根源分析

### 1. 模型使用了软删除

所有 BrandTrekin 模型都继承自 `global.GVA_MODEL`，包含软删除字段：

```go
// global/model.go
type GVA_MODEL struct {
    ID        uint           `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // ✅ 软删除字段
}
```

### 2. 唯一索引不包含 DeletedAt

商品模型的 ASIN 字段有唯一索引，但不包含 `deleted_at`：

```go
// model/brandtrekin/btProduct.go
type BtProduct struct {
    global.GVA_MODEL
    Asin *string `gorm:"unique;comment:亚马逊ASIN;column:asin;size:20;"`
    // ...
}
```

数据库约束：`uni_bt_products_asin` 只对 `asin` 字段建立唯一索引

### 3. 软删除与唯一索引的冲突

#### 问题场景：

```sql
-- 第一次导入：插入数据
INSERT INTO bt_products (id, asin, deleted_at) 
VALUES (1, 'B09ZTYGGTJ', NULL);

-- 全量替换：执行软删除（只更新 deleted_at）
UPDATE bt_products 
SET deleted_at = '2025-11-01 21:00:00' 
WHERE asin = 'B09ZTYGGTJ';

-- 数据库中的实际状态：
-- id=1, asin='B09ZTYGGTJ', deleted_at='2025-11-01 21:00:00' ✅ 数据仍然存在

-- 第二次导入：尝试插入新数据
INSERT INTO bt_products (id, asin, deleted_at) 
VALUES (2, 'B09ZTYGGTJ', NULL);

-- ❌ 错误！唯一索引冲突
-- Error 1062: Duplicate entry 'B09ZTYGGTJ' for key 'uni_bt_products_asin'
```

#### 为什么会冲突？

1. **软删除不是真正删除**：
   - GORM 的 `Delete()` 操作只是设置 `deleted_at` 字段
   - 数据仍然存在于数据库中

2. **唯一索引不区分软删除状态**：
   - 唯一索引 `uni_bt_products_asin` 只对 `asin` 字段生效
   - 不管 `deleted_at` 是否为 NULL，都会检查唯一性
   - 因此无法插入相同的 ASIN

### 4. 代码流程分析

```go
// bt_import.go

// 全量替换模式：删除现有数据
if replaceMode {
    if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
        return s.deleteMarketData(tx, marketID)  // ❌ 软删除
    }); err != nil {
        return fmt.Errorf("删除现有数据失败: %v", err)
    }
}

// deleteMarketData 方法（修复前）
func (s *BtImportService) deleteMarketData(tx *gorm.DB, marketID int64) error {
    // ❌ 软删除：只设置 deleted_at，数据仍在数据库中
    if err := tx.Model(&brandtrekin.BtProduct{}).
        Where("market_id = ?", marketID).
        Delete(&brandtrekin.BtProduct{}).Error; err != nil {
        return err
    }
    // ...
}

// 保存新数据
func (s *BtImportService) SaveProductDataBatch(...) error {
    // ❌ 插入失败：ASIN 唯一索引冲突
    if err := tx.Create(&toCreate).Error; err != nil {
        return fmt.Errorf("批量创建商品失败: %v", err)
    }
}
```

## ✅ 解决方案

### 核心思路

在全量替换模式下，使用 **物理删除** 而非软删除，确保旧数据真正从数据库中移除。

### 修复方法：使用 Unscoped()

GORM 提供了 `Unscoped()` 方法来执行物理删除：

```go
// 软删除（默认）
db.Delete(&model)  // UPDATE models SET deleted_at = NOW()

// 物理删除
db.Unscoped().Delete(&model)  // DELETE FROM models
```

### 修复后的代码

```go
// deleteMarketData 删除市场的所有关联数据
// 注意：使用 Unscoped() 进行物理删除，避免软删除导致的唯一索引冲突
func (s *BtImportService) deleteMarketData(tx *gorm.DB, marketID int64) error {
    // 按顺序删除关联数据（从子表到父表）

    // 1. 删除商品月度销售数据（物理删除）
    if err := tx.Unscoped().Where("asin IN (SELECT asin FROM bt_products WHERE market_id = ?)", marketID).
        Delete(&brandtrekin.BtProductMonthlySales{}).Error; err != nil {
        return err
    }

    // 2. 删除商品数据（物理删除，避免 ASIN 唯一索引冲突）
    if err := tx.Unscoped().Where("market_id = ?", marketID).Delete(&brandtrekin.BtProduct{}).Error; err != nil {
        return err
    }

    // 3. 删除品牌月度趋势数据（物理删除）
    if err := tx.Unscoped().Where("brand_id IN (SELECT id FROM bt_brands WHERE market_id = ? AND deleted_at IS NULL)", marketID).
        Delete(&brandtrekin.BtBrandMonthlyTrend{}).Error; err != nil {
        return err
    }

    // 4. 删除品牌社交媒体数据（物理删除）
    if err := tx.Unscoped().Where("brand_id IN (SELECT id FROM bt_brands WHERE market_id = ? AND deleted_at IS NULL)", marketID).
        Delete(&brandtrekin.BtBrandSocialMedia{}).Error; err != nil {
        return err
    }

    // 5. 删除品牌数据（物理删除，避免品牌名唯一索引冲突）
    if err := tx.Unscoped().Where("market_id = ?", marketID).Delete(&brandtrekin.BtBrand{}).Error; err != nil {
        return err
    }

    // 6. 删除关键词月度搜索量数据（物理删除）
    if err := tx.Unscoped().Where("keyword_id IN (SELECT id FROM bt_keywords WHERE market_id = ? AND deleted_at IS NULL)", marketID).
        Delete(&brandtrekin.BtKeywordMonthlyVolume{}).Error; err != nil {
        return err
    }

    // 7. 删除关键词数据（物理删除，避免关键词唯一索引冲突）
    if err := tx.Unscoped().Where("market_id = ?", marketID).Delete(&brandtrekin.BtKeyword{}).Error; err != nil {
        return err
    }

    // 8. 删除市场月度趋势数据（物理删除）
    if err := tx.Unscoped().Where("market_id = ?", marketID).Delete(&brandtrekin.BtMarketMonthlyTrend{}).Error; err != nil {
        return err
    }

    return nil
}
```

## 📊 修复效果对比

### 修复前（软删除）

```sql
-- 执行 Delete() 操作
UPDATE bt_products SET deleted_at = '2025-11-01 21:00:00' WHERE market_id = 1;

-- 数据仍然存在
SELECT * FROM bt_products WHERE market_id = 1;
-- 结果：id=1, asin='B09ZTYGGTJ', deleted_at='2025-11-01 21:00:00'

-- 插入新数据
INSERT INTO bt_products (asin, deleted_at) VALUES ('B09ZTYGGTJ', NULL);
-- ❌ Error 1062: Duplicate entry 'B09ZTYGGTJ'
```

### 修复后（物理删除）

```sql
-- 执行 Unscoped().Delete() 操作
DELETE FROM bt_products WHERE market_id = 1;

-- 数据已被真正删除
SELECT * FROM bt_products WHERE market_id = 1;
-- 结果：空（0 rows）

-- 插入新数据
INSERT INTO bt_products (asin, deleted_at) VALUES ('B09ZTYGGTJ', NULL);
-- ✅ 成功！
```

## 🔧 修改的文件

- **[bt_import.go](/Users/leon/code/ai-code/gin-vue-admin/server/service/brandtrekin/bt_import.go)**
  - `deleteMarketData()` - 所有删除操作都改为使用 `Unscoped()` 进行物理删除

## 📝 涉及的唯一索引

| 表名 | 唯一索引名 | 索引字段 | 是否包含 deleted_at |
|------|-----------|---------|-------------------|
| bt_products | uni_bt_products_asin | asin | ❌ 否 |
| bt_brands | uni_bt_brands_market_brand | market_id + brand_name | ❌ 否 |
| bt_keywords | uni_bt_keywords_market_keyword_source | market_id + keyword + source | ❌ 否 |
| bt_product_monthly_sales | uni_bt_product_monthly_sales_asin_date | asin + date | ❌ 否 |

**结论**：所有唯一索引都不包含 `deleted_at` 字段，因此在全量替换模式下必须使用物理删除。

## 🎯 为什么不修改唯一索引？

### 方案对比

#### 方案 1：修改唯一索引（不推荐）

```sql
-- 删除旧索引
ALTER TABLE bt_products DROP INDEX uni_bt_products_asin;

-- 创建包含 deleted_at 的唯一索引
ALTER TABLE bt_products ADD UNIQUE INDEX uni_bt_products_asin (asin, deleted_at);
```

**缺点**：
- ❌ 允许相同 ASIN 的多条软删除记录存在
- ❌ 数据库会积累大量软删除的历史数据
- ❌ 查询性能下降
- ❌ 违反业务逻辑（ASIN 应该全局唯一）

#### 方案 2：使用物理删除（推荐）✅

```go
// 全量替换时使用物理删除
tx.Unscoped().Delete(&model)
```

**优点**：
- ✅ 保持唯一索引的语义正确性
- ✅ 避免历史数据积累
- ✅ 查询性能更好
- ✅ 符合全量替换的业务逻辑

## ⚠️ 注意事项

### 1. 物理删除的影响

- **数据无法恢复**：物理删除后，数据将永久丢失
- **适用场景**：全量替换模式（旧数据不需要保留）
- **不适用场景**：增量导入模式（应该保留历史数据）

### 2. 其他删除场景

对于非全量替换的删除操作（如用户手动删除），仍然应该使用软删除：

```go
// 用户手动删除商品（软删除）
func (s *BtProductService) DeleteBtProduct(id uint) error {
    return global.GVA_DB.Delete(&brandtrekin.BtProduct{}, id).Error
}

// 全量替换删除（物理删除）
func (s *BtImportService) deleteMarketData(tx *gorm.DB, marketID int64) error {
    return tx.Unscoped().Where("market_id = ?", marketID).Delete(&brandtrekin.BtProduct{}).Error
}
```

## ✅ 测试验证

### 1. 测试全量替换

```bash
# 第一次导入
curl -X POST http://localhost:8888/api/v1/brandtrekin/import \
  -F "marketId=1" \
  -F "replaceMode=true" \
  -F "files=@Products.csv"

# 检查数据
SELECT COUNT(*) FROM bt_products WHERE market_id = 1 AND deleted_at IS NULL;
# 预期：有数据

# 第二次导入（全量替换）
curl -X POST http://localhost:8888/api/v1/brandtrekin/import \
  -F "marketId=1" \
  -F "replaceMode=true" \
  -F "files=@Products.csv"

# 检查数据
SELECT COUNT(*) FROM bt_products WHERE market_id = 1 AND deleted_at IS NULL;
# 预期：有数据（新数据）

SELECT COUNT(*) FROM bt_products WHERE market_id = 1 AND deleted_at IS NOT NULL;
# 预期：0（旧数据已被物理删除）
```

### 2. 验证唯一索引

```sql
-- 查询是否有重复的 ASIN
SELECT asin, COUNT(*) as count 
FROM bt_products 
WHERE deleted_at IS NULL 
GROUP BY asin 
HAVING count > 1;

-- 预期：空结果（没有重复）
```

## 🎉 总结

### 问题根源
- 软删除 + 不包含 deleted_at 的唯一索引 = 重复键冲突

### 解决方案
- 全量替换模式下使用 `Unscoped()` 进行物理删除

### 修复效果
- ✅ 全量替换不再报重复键错误
- ✅ 数据库不会积累软删除的历史数据
- ✅ 唯一索引语义正确
- ✅ 查询性能更好

修复完成！🎉
