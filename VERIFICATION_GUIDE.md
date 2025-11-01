# 修复验证指南

## 📋 修复内容

已修复 `bt_product_monthly_sales` 表中 `units > 0` 但 `sales = 0` 的问题。

**修改文件**：`server/service/brandtrekin/bt_import.go`

**修改内容**：
1. ✅ 允许处理销售额为 0 的数据
2. ✅ 允许处理销量为 0 的数据
3. ✅ 不跳过空单元格，将其视为 0
4. ✅ 修复数据更新逻辑

详细说明请查看：[FIX_SALES_ZERO_ISSUE.md](./FIX_SALES_ZERO_ISSUE.md)

---

## 🚀 快速验证步骤

### 步骤 1：重启后端服务

```bash
# 方式 1：使用重启脚本（推荐）
chmod +x restart_server.sh
./restart_server.sh

# 方式 2：手动重启
# 停止旧服务
kill -9 $(lsof -ti :8888)

# 启动新服务
cd server
nohup go run main.go > ../logs/server.log 2>&1 &
cd ..
```

### 步骤 2：导入测试数据

**方式 1：通过 API 导入（推荐）**

```bash
chmod +x test_import_api.sh
./test_import_api.sh
```

**方式 2：通过前端界面导入**

1. 访问 http://localhost:8888
2. 登录系统
3. 进入「市场管理」→「批量导入」
4. 选择参数：
   - 市场：CNC Router (marketId=1)
   - 导入模式：全量替换 (replaceMode=true)
   - 文件：上传 `trekin-main/data/CNCRouter/product-US-sales.xlsx`
5. 点击「开始导入」

**方式 3：使用 curl 命令**

```bash
curl -X POST http://localhost:8888/btImport/batchImport \
  -F "marketId=1" \
  -F "replaceMode=true" \
  -F "productSales=@trekin-main/data/CNCRouter/product-US-sales.xlsx"
```

### 步骤 3：验证修复效果

```bash
chmod +x verify_sales_fix.sh
./verify_sales_fix.sh
```

**预期结果**：
- ✅ 异常记录数应该为 0 或大幅减少
- ✅ 正常记录数应该增加
- ✅ 数据完整性提高

---

## 📊 验证检查项

### 1. 数据库检查

```sql
-- 检查异常数据数量
SELECT COUNT(*) as '异常记录数' 
FROM bt_product_monthly_sales 
WHERE units > 0 AND sales = 0;

-- 查看具体异常数据
SELECT asin, DATE_FORMAT(date, '%Y-%m') as month, units, sales 
FROM bt_product_monthly_sales 
WHERE units > 0 AND sales = 0 
ORDER BY date DESC 
LIMIT 10;

-- 统计数据分布
SELECT 
    COUNT(*) as '总记录数',
    SUM(CASE WHEN units > 0 AND sales = 0 THEN 1 ELSE 0 END) as '异常记录',
    SUM(CASE WHEN units > 0 AND sales > 0 THEN 1 ELSE 0 END) as '正常记录',
    SUM(CASE WHEN units = 0 AND sales = 0 THEN 1 ELSE 0 END) as '零值记录'
FROM bt_product_monthly_sales 
WHERE market_id = 1;
```

### 2. Excel 文件检查

打开 `trekin-main/data/CNCRouter/product-US-sales.xlsx`，检查：

- ✅ 是否有销量 > 0 但销售额为空的单元格
- ✅ 是否有销量 > 0 但销售额为 0 的单元格
- ✅ 数据格式是否正确

### 3. 日志检查

```bash
# 查看导入日志
tail -f logs/server.log

# 查看是否有错误
grep -i "error" logs/server.log | tail -20
```

---

## 🎯 测试场景

### 场景 1：销量有值，销售额为空

**Excel 数据**：
| ASIN | 2024-01 | 2024-02 |
|------|---------|---------|
| B09ZTYGGTJ | 100 | (空) |

**预期结果**：
- 2024-01: units=100, sales=0 ✅
- 2024-02: units=0, sales=0 ✅

### 场景 2：销量有值，销售额为 0

**Excel 数据**：
| ASIN | 2024-01($) |
|------|------------|
| B08XYZABC | 0 |

**预期结果**：
- 2024-01: units=0, sales=0 ✅

### 场景 3：销量和销售额都有值

**Excel 数据**：
| ASIN | 2024-01 | 2024-01($) |
|------|---------|------------|
| B07DEFGHI | 50 | 1234.56 |

**预期结果**：
- 2024-01: units=50, sales=1234.56 ✅

---

## 🔧 故障排查

### 问题 1：服务启动失败

**检查**：
```bash
# 查看端口占用
lsof -i :8888

# 查看日志
tail -50 logs/server.log
```

**解决**：
- 确保端口 8888 未被占用
- 检查配置文件 `server/config.yaml`
- 检查数据库连接

### 问题 2：导入失败

**检查**：
```bash
# 查看导入日志
grep "ParseProductSales" logs/server.log | tail -20

# 查看错误信息
grep -i "error" logs/server.log | tail -20
```

**解决**：
- 确认文件路径正确
- 确认文件格式正确（xlsx）
- 确认 sheet 名称正确
- 检查数据库连接

### 问题 3：数据仍然异常

**检查**：
```bash
# 确认代码已更新
grep -A 5 "允许处理空单元格" server/service/brandtrekin/bt_import.go

# 确认服务已重启
ps aux | grep "main" | grep -v grep
```

**解决**：
- 确认代码修改已保存
- 确认服务已重启（新代码生效）
- 清空旧数据后重新导入

---

## 📝 验证报告模板

```
验证时间：2025-11-01 22:34:22
验证人员：[你的名字]

1. 服务重启：✅ 成功 / ❌ 失败
   - 进程 ID：_______
   - 启动时间：_______

2. 数据导入：✅ 成功 / ❌ 失败
   - 导入文件：product-US-sales.xlsx
   - 导入模式：全量替换
   - 导入时间：_______
   - 导入记录数：_______

3. 数据验证：✅ 通过 / ❌ 未通过
   - 总记录数：_______
   - 异常记录数（units>0 且 sales=0）：_______
   - 正常记录数（units>0 且 sales>0）：_______
   - 零值记录数（units=0 且 sales=0）：_______

4. 修复效果：
   - 修复前异常记录数：_______
   - 修复后异常记录数：_______
   - 改善率：_______%

5. 备注：
   _______________________________________
```

---

## 🎉 验证成功标准

- ✅ 服务成功重启
- ✅ 数据成功导入
- ✅ 异常记录数为 0 或大幅减少（> 90%）
- ✅ 正常记录数增加
- ✅ 数据完整性提高
- ✅ 无新的错误日志

---

## 📞 联系支持

如果遇到问题，请提供：
1. 错误日志（`logs/server.log`）
2. 数据库查询结果
3. Excel 文件样本
4. 验证报告

---

**文档版本**：1.0  
**最后更新**：2025-11-01  
**相关文档**：[FIX_SALES_ZERO_ISSUE.md](./FIX_SALES_ZERO_ISSUE.md)
