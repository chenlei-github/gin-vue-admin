# BrandTrekin 聚合计算逻辑修复总结

## 问题描述

同一市场 cnc-router 的数据，trekin-main 和 brandtrekin 计算结果存在巨大差异：

| 指标 | trekin-main（正确） | brandtrekin（修复前） | 差异 |
|------|-------------------|---------------------|------|
| 市场规模（年度） | $13.87M | $32,222 | 相差 430 倍 |
| 市场增速（CAGR） | 28.9% | -99.0% | 完全相反 |
| 月搜索量总和 | 241.5K | 0 | 无数据 |

## 根本原因分析

### 1. 市场规模计算错误

**trekin-main 的正确逻辑**：
```javascript
// 计算最近12个月的销售额（基于数据中的最新日期）
if (sortedMonths.length >= 12) {
  const last12Months = sortedMonths.slice(-12);
  totalRevenue = last12Months.reduce((sum, month) => sum + monthlySales[month], 0);
}
```

**brandtrekin 修复前的错误逻辑**：
```go
// 错误：使用当前系统时间往前推12个月
twelveMonthsAgo := time.Now().AddDate(0, -12, 0)
err := tx.Model(&brandtrekin.BtMarketMonthlyTrend{}).
    Select("COALESCE(SUM(revenue), 0)").
    Where("market_id = ? AND date >= ?", marketID, twelveMonthsAgo).
    Scan(&totalRevenue).Error
```

**问题**：数据的最新日期是 2025-06，而当前时间是 2025-11，导致查询不到任何数据。

### 2. CAGR（年复合增长率）计算错误

**trekin-main 的正确逻辑**：
```javascript
// 使用前12个月和后12个月的总收入来计算
const firstYearRevenue = sortedMonths.slice(0, 12).reduce(...);
const lastYearRevenue = sortedMonths.slice(-12).reduce(...);

// 计算两个时期中点之间的年数
const firstDate = new Date(sortedMonths[5]); // 前12个月的中点
const lastDate = new Date(sortedMonths[sortedMonths.length - 6]); // 后12个月的中点
const years = (lastDate - firstDate) / (365.25 * 24 * 60 * 60 * 1000);

// CAGR = (Ending Value / Beginning Value)^(1/years) - 1
cagr = (Math.pow(lastYearRevenue / firstYearRevenue, 1 / years) - 1) * 100;
```

**brandtrekin 修复前的错误逻辑**：
```go
// 错误：使用第一个月和最后一个月的单月数据
beginningValue = *trends[0].Revenue
endingValue = *trends[len(trends)-1].Revenue

// 错误：使用总月数除以12作为年数
years := float64(len(trends)) / 12.0

cagr := (math.Pow(endingValue/beginningValue, 1.0/years) - 1.0) * 100.0
```

**问题**：
1. 单月数据波动大，不能代表整体趋势
2. 年数计算不准确，应该使用实际的时间差

### 3. 搜索量获取错误

**trekin-main 的正确逻辑**：
```javascript
// 获取数据中最新月份的搜索量
if (sortedKeywordMonths.length > 0) {
  const latestMonth = sortedKeywordMonths[sortedKeywordMonths.length - 1];
  totalSearchVolume = keywordTrends[latestMonth] || 0;
}
```

**brandtrekin 修复前的错误逻辑**：
```go
// 错误：按日期倒序查询，但可能查不到数据
err = tx.Model(&brandtrekin.BtMarketMonthlyTrend{}).
    Select("COALESCE(search_volume, 0)").
    Where("market_id = ?", marketID).
    Order("date DESC").
    Limit(1).
    Scan(&searchVolume).Error
```

**问题**：查询逻辑本身没问题，但因为前面的聚合步骤有问题，导致没有数据。

## 修复方案

### 1. 修复市场规模计算

```go
// 修复后：获取数据中最新的12个月
var trends []brandtrekin.BtMarketMonthlyTrend
err := tx.Where("market_id = ?", marketID).
    Order("date DESC").
    Limit(12).
    Find(&trends).Error

// 计算最近12个月的总销售额
var totalRevenue float64
for _, trend := range trends {
    if trend.Revenue != nil {
        totalRevenue += *trend.Revenue
    }
}
```

### 2. 修复 CAGR 计算

```go
// 修复后：使用前12个月和后12个月的总收入
if len(trends) < 24 {
    return nil  // 至少需要24个月数据
}

// 计算前12个月的总收入
var firstYearRevenue float64
for i := 0; i < 12 && i < len(trends); i++ {
    if trends[i].Revenue != nil {
        firstYearRevenue += *trends[i].Revenue
    }
}

// 计算后12个月的总收入
var lastYearRevenue float64
startIdx := len(trends) - 12
for i := startIdx; i < len(trends); i++ {
    if trends[i].Revenue != nil {
        lastYearRevenue += *trends[i].Revenue
    }
}

// 计算实际年数（使用两个时期中点之间的时间差）
firstDate := trends[5].Date  // 前12个月的中点
lastDate := trends[len(trends)-6].Date  // 后12个月的中点
years := lastDate.Sub(*firstDate).Hours() / (365.25 * 24)

// CAGR公式
cagr := (math.Pow(lastYearRevenue/firstYearRevenue, 1.0/years) - 1.0) * 100.0
```

### 3. 修复搜索量获取

```go
// 修复后：从已查询的trends中获取最新月份的搜索量
var searchVolume int64
if len(trends) > 0 && trends[0].SearchVolume != nil {
    searchVolume = *trends[0].SearchVolume
}
```

## 修复后的预期结果

修复后，brandtrekin 的计算结果应该与 trekin-main 保持一致：

| 指标 | 预期值 |
|------|--------|
| 市场规模（年度） | ~$13.87M |
| 市场增速（CAGR） | ~28.9% |
| 月搜索量总和 | ~241.5K |

## 关键改进点

1. **时间范围修正**：从"基于当前系统时间"改为"基于数据中的实际日期"
2. **CAGR 计算优化**：从"单月对比"改为"年度总收入对比"
3. **年数计算精确化**：从"简单除法"改为"实际时间差计算"
4. **数据阈值保护**：添加最小值检查（如 firstYearRevenue > 1000），避免极端值

## 测试建议

1. 重新导入 cnc-router 市场的数据
2. 运行聚合计算
3. 对比结果与 trekin-main 的 cnc-router.json
4. 验证其他市场（laser-engraver）的计算结果

## 相关文件

- 修复文件：`/server/service/brandtrekin/bt_aggregate.go`
- 参考实现：`/trekin-main/scripts/parseData.js`
- 测试数据：`/trekin-main/public/data/cnc-router.json`
