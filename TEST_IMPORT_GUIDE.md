# 全量替换导入修复效果测试指南

## 📋 测试目标

验证全量替换模式下的两个修复：
1. ✅ **软删除冲突修复** - 使用 `Unscoped()` 进行物理删除
2. ✅ **CSV重复数据修复** - 批量保存前自动去重

## 🔧 测试环境

- **服务器状态**: ✅ 已启动 (端口 8888)
- **数据库**: brandtrekin @ rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com
- **测试市场**: cnc-router
- **测试数据**: 
  - `/Users/leon/code/ai-code/gin-vue-admin/trekin-main/data/CNCRouter/GKW.csv`
  - `/Users/leon/code/ai-code/gin-vue-admin/trekin-main/data/CNCRouter/` 下的其他 CSV 文件

## 📝 测试步骤

### 步骤 1: 访问前端界面

1. 打开浏览器访问: `http://localhost:8080` (或前端运行的端口)
2. 登录系统
3. 进入 **BrandTrekin** > **市场数据导入** 页面

### 步骤 2: 第一次全量导入（基准测试）

1. **选择市场**: `cnc-router`
2. **选择导入模式**: `全量替换` ✅
3. **上传文件**: 
   - 选择 `/Users/leon/code/ai-code/gin-vue-admin/trekin-main/data/CNCRouter/` 目录下的所有 CSV 文件
   - 包括: GKW.csv, Products.csv, Brands.csv 等
4. **点击导入**
5. **预期结果**: 
   - ✅ 导入成功
   - ✅ 显示导入的商品、品牌、关键词数量
   - ✅ 无错误信息

### 步骤 3: 记录第一次导入的数据

记录导入后的数据统计：
- 商品数量: ______
- 品牌数量: ______
- 关键词数量: ______
- 销售数据: ______

### 步骤 4: 第二次全量导入（关键测试）

**这是测试修复效果的关键步骤！**

1. **不要关闭浏览器，不要退出登录**
2. **再次进入市场数据导入页面**
3. **选择市场**: `cnc-router` (相同市场)
4. **选择导入模式**: `全量替换` ✅ (相同模式)
5. **上传文件**: 
   - **上传相同的 CSV 文件**（模拟重复导入场景）
6. **点击导入**
7. **预期结果**: 
   - ✅ 导入成功（修复前会失败）
   - ✅ 无重复键错误
   - ✅ 数据量与第一次相同或相近

### 步骤 5: 验证修复效果

#### 5.1 检查导入结果

- [ ] 第二次导入是否成功？
- [ ] 是否出现 `Duplicate entry` 错误？
- [ ] 数据量是否正确？

#### 5.2 检查数据完整性

在前端查看市场数据：
1. 进入 **市场分析** 页面
2. 选择 `cnc-router` 市场
3. 检查：
   - [ ] 商品列表是否正常显示
   - [ ] 品牌列表是否正常显示
   - [ ] 关键词列表是否正常显示
   - [ ] 市场规模、增速等指标是否正确

#### 5.3 检查数据库（可选）

如果有数据库访问权限，可以执行以下 SQL 验证：

```sql
-- 1. 检查软删除数据（应该为 0）
SELECT COUNT(*) as soft_deleted_count
FROM bt_products 
WHERE market_id = 1 AND deleted_at IS NOT NULL;

-- 2. 检查重复ASIN（应该为空）
SELECT asin, COUNT(*) as count
FROM bt_products
WHERE market_id = 1 AND deleted_at IS NULL
GROUP BY asin
HAVING count > 1;

-- 3. 检查总数据量
SELECT 
    (SELECT COUNT(*) FROM bt_products WHERE market_id = 1) as products,
    (SELECT COUNT(*) FROM bt_brands WHERE market_id = 1) as brands,
    (SELECT COUNT(*) FROM bt_keywords WHERE market_id = 1) as keywords;
```

## ✅ 测试结果判定

### 修复成功的标志

- ✅ 第二次全量导入成功，无错误
- ✅ 无 `Duplicate entry` 错误
- ✅ 数据量正确
- ✅ 软删除数据为 0
- ✅ 无重复ASIN数据

### 修复失败的标志

- ❌ 第二次导入报错: `Error 1062 (23000): Duplicate entry 'xxx' for key 'bt_products.uni_bt_products_asin'`
- ❌ 数据量异常
- ❌ 存在软删除数据残留
- ❌ 存在重复ASIN数据

## 🔄 额外测试场景

### 场景 1: 测试 CSV 重复数据去重

1. 手动编辑 CSV 文件，添加重复的 ASIN
2. 使用全量替换模式导入
3. 验证：
   - [ ] 导入成功
   - [ ] 重复数据被自动去重
   - [ ] 保留最后一条记录

### 场景 2: 测试增量导入模式

1. 第一次使用全量替换导入
2. 第二次使用增量导入模式
3. 验证：
   - [ ] 增量导入成功
   - [ ] 现有数据被更新
   - [ ] 新数据被添加

### 场景 3: 测试不同市场

1. 导入 `cnc-router` 市场
2. 导入 `laser-engraver` 市场
3. 验证：
   - [ ] 两个市场数据互不影响
   - [ ] 都能成功导入

## 📊 测试记录表

| 测试项 | 预期结果 | 实际结果 | 状态 |
|--------|---------|---------|------|
| 第一次全量导入 | 成功 | | ⬜ |
| 第二次全量导入 | 成功 | | ⬜ |
| 无重复键错误 | 是 | | ⬜ |
| 软删除数据 | 0 | | ⬜ |
| 重复ASIN数据 | 0 | | ⬜ |
| 数据完整性 | 正确 | | ⬜ |
| CSV重复数据去重 | 成功 | | ⬜ |
| 增量导入 | 成功 | | ⬜ |

## 🐛 问题排查

### 如果第二次导入仍然失败

1. **检查代码是否已更新**:
   ```bash
   cd /Users/leon/code/ai-code/gin-vue-admin/server
   grep -n "Unscoped()" service/brandtrekin/bt_import.go
   ```
   应该看到多处 `Unscoped()` 调用

2. **检查服务器是否重启**:
   - 修改代码后需要重启服务器
   - 如果使用热重载，确认已生效

3. **检查数据库连接**:
   - 确认数据库配置正确
   - 确认有足够的权限

4. **查看服务器日志**:
   - 检查 `/Users/leon/code/ai-code/gin-vue-admin/server/log/` 目录
   - 查找错误信息

### 如果数据量异常

1. **检查 CSV 文件**:
   - 确认文件格式正确
   - 确认文件编码正确（UTF-8 或 GBK）

2. **检查去重逻辑**:
   - 查看 `bt_import_optimized.go` 中的去重代码
   - 确认 map 去重逻辑正确

## 📞 联系支持

如果测试过程中遇到问题：
1. 记录详细的错误信息
2. 截图保存错误页面
3. 导出服务器日志
4. 提供测试步骤和数据

## 🎉 测试完成

完成所有测试后，请填写：

- **测试日期**: _______________
- **测试人员**: _______________
- **测试结果**: ⬜ 通过 / ⬜ 失败
- **备注**: _______________

---

**修复文件**:
- [bt_import.go](/Users/leon/code/ai-code/gin-vue-admin/server/service/brandtrekin/bt_import.go) - 添加 `Unscoped()` 物理删除
- [bt_import_optimized.go](/Users/leon/code/ai-code/gin-vue-admin/server/service/brandtrekin/bt_import_optimized.go) - 添加 CSV 去重逻辑

**文档**:
- [FIX_SOFT_DELETE_CONFLICT.md](/Users/leon/code/ai-code/gin-vue-admin/FIX_SOFT_DELETE_CONFLICT.md) - 软删除冲突详细说明
- [FIX_DUPLICATE_KEY_SUMMARY.md](/Users/leon/code/ai-code/gin-vue-admin/FIX_DUPLICATE_KEY_SUMMARY.md) - 问题总结对比
