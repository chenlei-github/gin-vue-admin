package main

import (
	"fmt"
	"time"
)

// 模拟 trekin-main 的逻辑
func trekinMainLogic() {
	// 模拟关键词月度搜索量数据
	keywordData := map[string][]struct {
		Date   string
		Volume int
	}{
		"cnc machine": {
			{"2024-06", 150000},
			{"2024-07", 160000},
			{"2024-08", 165000},
		},
		"cnc router": {
			{"2024-06", 20000},
			{"2024-07", 21000},
			{"2024-08", 22200},
		},
		"desktop cnc": {
			{"2024-06", 4000},
			{"2024-07", 4200},
			{"2024-08", 4400},
		},
	}

	// 1. 按日期汇总所有关键词的搜索量
	keywordTrends := make(map[string]int)
	for keyword, volumes := range keywordData {
		for _, v := range volumes {
			keywordTrends[v.Date] += v.Volume
		}
		fmt.Printf("关键词: %s\n", keyword)
		for _, v := range volumes {
			fmt.Printf("  %s: %d\n", v.Date, v.Volume)
		}
	}

	fmt.Println("\n=== 按月份汇总的搜索量 ===")
	dates := []string{"2024-06", "2024-07", "2024-08"}
	for _, date := range dates {
		fmt.Printf("%s: %d\n", date, keywordTrends[date])
	}

	// 2. 取最新月份的搜索量
	latestMonth := dates[len(dates)-1]
	totalSearchVolume := keywordTrends[latestMonth]

	fmt.Printf("\n=== trekin-main 逻辑 ===\n")
	fmt.Printf("最新月份: %s\n", latestMonth)
	fmt.Printf("月度搜索量总和: %d (%.1fK)\n", totalSearchVolume, float64(totalSearchVolume)/1000)
}

// 模拟修正后的 Go 代码逻辑
func goCodeLogic() {
	// 模拟数据库查询结果
	type MonthlySearchVolume struct {
		Date   time.Time
		Volume int
	}

	searchVolumes := []MonthlySearchVolume{
		{time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), 174000},
		{time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), 185200},
		{time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC), 191600},
	}

	// 创建搜索量映射
	searchVolumeMap := make(map[string]int)
	for _, sv := range searchVolumes {
		key := sv.Date.Format("2006-01-02")
		searchVolumeMap[key] = sv.Volume
	}

	fmt.Println("\n=== 修正后的 Go 代码逻辑 ===")
	fmt.Println("每个月份的搜索量映射:")
	for _, sv := range searchVolumes {
		dateKey := sv.Date.Format("2006-01-02")
		fmt.Printf("  %s: %d\n", dateKey, searchVolumeMap[dateKey])
	}

	// 模拟市场月度趋势数据
	type MonthlyMarketSum struct {
		Date       time.Time
		TotalSales float64
	}

	monthlySums := []MonthlyMarketSum{
		{time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), 1000000},
		{time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), 1100000},
		{time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC), 1200000},
	}

	fmt.Println("\n保存到市场月度趋势:")
	for _, sum := range monthlySums {
		dateKey := sum.Date.Format("2006-01-02")
		searchVolume := searchVolumeMap[dateKey]
		fmt.Printf("  日期: %s, 销售额: $%.0f, 搜索量: %d\n", 
			dateKey, sum.TotalSales, searchVolume)
	}

	// 获取最新月份的搜索量（用于市场汇总指标）
	latestSearchVolume := searchVolumes[len(searchVolumes)-1].Volume
	fmt.Printf("\n市场汇总指标中的搜索量（最新月份）: %d (%.1fK)\n", 
		latestSearchVolume, float64(latestSearchVolume)/1000)
}

func main() {
	fmt.Println("========================================")
	fmt.Println("月度搜索量计算逻辑对比")
	fmt.Println("========================================\n")

	trekinMainLogic()
	goCodeLogic()

	fmt.Println("\n========================================")
	fmt.Println("✅ 关键点:")
	fmt.Println("1. trekin-main: 取最新月份的搜索量总和")
	fmt.Println("2. Go 代码修正前: 所有月份都用最新月份的搜索量（错误）")
	fmt.Println("3. Go 代码修正后: 每个月份用对应月份的搜索量（正确）")
	fmt.Println("4. 市场汇总指标: 使用最新月份的搜索量（与 trekin-main 一致）")
	fmt.Println("========================================")
}
