# 月度搜索量计算逻辑修正说明

## 问题描述

在 `bt_aggregate.go` 的 `AggregateBrandToMarket` 方法中，第 140-220 行的月度搜索量计算逻辑存在问题。

### 修正前的错误逻辑

```go
// 错误：所有月份都使用最新月份的搜索量
var searchVolumes []MonthlySearchVolume
err = tx.Model(&brandtrekin.BtKeywordMonthlyVolume{}).
    Select("date, SUM(volume) as volume").
    Where("keyword_id IN ?", keywordIDs).
    Group("date").
    Order("date").
    Find(&searchVolumes).Error

// 取最新月份的搜索量
latestSearchVolume := searchVolumes[len(searchVolumes)-1]

// 错误：所有月份的市场趋势都使用同一个搜索量值
for _, sum := range monthlySums {
    searchVolumeInt64 := int64(latestSearchVolume.Volume)  // ❌ 错误
    trend := brandtrekin.BtMarketMonthlyTrend{
        SearchVolume: &searchVolumeInt64,
    }
}
```

**问题**：所有月份的市场趋势记录都被赋予了最新月份的搜索量，导致历史月份的搜索量数据不准确。

## trekin-main 的正确逻辑

```javascript
// 1. 按日期汇总所有关键词的搜索量
const keywordTrends = {};
keywords.forEach(keyword => {
  keyword.monthlyVolume.forEach(month => {
    if (!keywordTrends[month.date]) {
      keywordTrends[month.date] = 0;
    }
    keywordTrends[month.date] += month.volume || 0;  // 每个月份独立汇总
  });
});

// 2. 取最新月份的搜索量（仅用于市场汇总指标）
const sortedKeywordMonths = Object.keys(keywordTrends).sort();
let totalSearchVolume = 0;
if (sortedKeywordMonths.length > 0) {
  const latestMonth = sortedKeywordMonths[sortedKeywordMonths.length - 1];
  totalSearchVolume = keywordTrends[latestMonth] || 0;
}
```

**关键点**：
1. 每个月份都有自己独立的搜索量总和
2. 最新月份的搜索量仅用于市场的汇总指标（`searchVolume` 字段）
3. 市场月度趋势表中，每条记录应该使用对应月份的搜索量

## 修正后的正确逻辑

```go
// 创建搜索量映射：key为日期字符串，value为该月所有关键词的搜索量总和
searchVolumeMap := make(map[string]int)

if len(keywordIDs) > 0 {
    var searchVolumes []MonthlySearchVolume
    err = tx.Model(&brandtrekin.BtKeywordMonthlyVolume{}).
        Select("date, SUM(volume) as volume").
        Where("keyword_id IN ?", keywordIDs).
        Group("date").
        Order("date").
        Find(&searchVolumes).Error

    // 填充搜索量映射：每个月份对应该月的搜索量总和
    for _, sv := range searchVolumes {
        key := sv.Date.Format("2006-01-02")
        searchVolumeMap[key] = sv.Volume  // ✅ 每个月份独立存储
    }
}

// 保存或更新市场月度趋势
for _, sum := range monthlySums {
    date := sum.Date
    revenue := sum.TotalSales

    // 获取该月份对应的搜索量（而不是最新月份的搜索量）
    dateKey := date.Format("2006-01-02")
    searchVolume := searchVolumeMap[dateKey]  // ✅ 使用对应月份的搜索量

    searchVolumeInt64 := int64(searchVolume)
    trend := brandtrekin.BtMarketMonthlyTrend{
        MarketId:     &marketID,
        Date:         &date,
        Revenue:      &revenue,
        SearchVolume: &searchVolumeInt64,  // ✅ 正确
    }
    // ... 保存逻辑
}
```

## 数据示例对比

### 修正前（错误）

| 日期 | 销售额 | 搜索量 | 说明 |
|------|--------|--------|------|
| 2024-06 | $1,000,000 | 191,600 | ❌ 错误：使用了最新月份的搜索量 |
| 2024-07 | $1,100,000 | 191,600 | ❌ 错误：使用了最新月份的搜索量 |
| 2024-08 | $1,200,000 | 191,600 | ✅ 正确：最新月份 |

### 修正后（正确）

| 日期 | 销售额 | 搜索量 | 说明 |
|------|--------|--------|------|
| 2024-06 | $1,000,000 | 174,000 | ✅ 正确：使用该月的搜索量 |
| 2024-07 | $1,100,000 | 185,200 | ✅ 正确：使用该月的搜索量 |
| 2024-08 | $1,200,000 | 191,600 | ✅ 正确：使用该月的搜索量 |

## 市场汇总指标

在 `UpdateMarketMetrics` 方法中，市场的 `search_volume` 字段应该使用**最新月份**的搜索量，这与 trekin-main 的逻辑一致：

```go
// 4. 获取最新月份的搜索量（数据中的最新日期，而不是当前时间）
// 注意：trends 已经按 date DESC 排序，所以 trends[0] 是最新月份
var searchVolume int64
if len(trends) > 0 && trends[0].SearchVolume != nil {
    searchVolume = *trends[0].SearchVolume  // ✅ 最新月份的搜索量
}
```

## 总结

### 两个不同的概念

1. **市场月度趋势表** (`bt_market_monthly_trend`)
   - 每条记录代表一个月份
   - `search_volume` 字段应该是**该月份**所有关键词的搜索量总和
   - 用于展示历史趋势图

2. **市场汇总指标** (`bt_market` 表的 `search_volume` 字段)
   - 代表市场的当前搜索热度
   - 应该使用**最新月份**的搜索量
   - 用于市场列表展示

### 与 trekin-main 的一致性

修正后的逻辑与 trekin-main 完全一致：
- ✅ 每个月份使用对应月份的搜索量
- ✅ 市场汇总指标使用最新月份的搜索量
- ✅ 搜索量计算方式：SUM(该月所有关键词的搜索量)

## 测试验证

运行测试脚本验证修正效果：
```bash
go run test_search_volume.go
```

测试结果显示修正后的逻辑与 trekin-main 完全一致。
