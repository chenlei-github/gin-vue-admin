package brandtrekin

import (
	"fmt"
	"math"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
	"gorm.io/gorm"
)

type BtAggregateService struct{}

// AggregateProductToBrand 聚合商品月度销售数据到品牌月度趋势
// 目的：将每个商品的月度销售数据汇总到品牌级别
func (s *BtAggregateService) AggregateProductToBrand(marketID int64) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 1. 获取该市场下的所有品牌
		var brands []brandtrekin.BtBrand
		if err := tx.Where("market_id = ?", marketID).Find(&brands).Error; err != nil {
			return fmt.Errorf("查询品牌失败: %v", err)
		}

		// 2. 对每个品牌进行聚合
		for _, brand := range brands {
			brandID := int64(brand.ID)

			// 获取该品牌下所有商品的ASIN
			var products []brandtrekin.BtProduct
			if err := tx.Where("brand_id = ?", brandID).Find(&products).Error; err != nil {
				return fmt.Errorf("查询商品失败: %v", err)
			}

			if len(products) == 0 {
				continue
			}

			// 提取所有ASIN
			asins := make([]string, 0, len(products))
			for _, p := range products {
				if p.Asin != nil {
					asins = append(asins, *p.Asin)
				}
			}

			if len(asins) == 0 {
				continue
			}

			// 聚合：按日期分组，汇总所有商品的销售额
			type MonthlySalesSum struct {
				Date       time.Time
				TotalSales float64
				TotalUnits int
			}

			var monthlySums []MonthlySalesSum
			err := tx.Model(&brandtrekin.BtProductMonthlySales{}).
				Select("date, SUM(sales) as total_sales, SUM(units) as total_units").
				Where("asin IN ?", asins).
				Group("date").
				Order("date").
				Find(&monthlySums).Error

			if err != nil {
				return fmt.Errorf("聚合月度销售数据失败: %v", err)
			}

			// 3. 保存或更新品牌月度趋势
			for _, sum := range monthlySums {
				date := sum.Date
				revenue := sum.TotalSales

				trend := brandtrekin.BtBrandMonthlyTrend{
					BrandId: &brandID,
					Date:    &date,
					Revenue: &revenue,
				}

				// 查找或创建
				var existing brandtrekin.BtBrandMonthlyTrend
				err := tx.Where("brand_id = ? AND date = ?", brandID, date).
					First(&existing).Error

				if err == gorm.ErrRecordNotFound {
					if err := tx.Create(&trend).Error; err != nil {
						return fmt.Errorf("创建品牌月度趋势失败: %v", err)
					}
				} else if err != nil {
					return fmt.Errorf("查询品牌月度趋势失败: %v", err)
				} else {
					// 更新
					if err := tx.Model(&existing).Updates(&trend).Error; err != nil {
						return fmt.Errorf("更新品牌月度趋势失败: %v", err)
					}
				}
			}
		}

		return nil
	})
}

// AggregateBrandToMarket 聚合品牌月度趋势到市场月度趋势
func (s *BtAggregateService) AggregateBrandToMarket(marketID int64) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 1. 获取该市场下所有品牌ID
		var brands []brandtrekin.BtBrand
		if err := tx.Where("market_id = ?", marketID).Find(&brands).Error; err != nil {
			return fmt.Errorf("查询品牌失败: %v", err)
		}

		if len(brands) == 0 {
			return nil
		}

		brandIDs := make([]int64, 0, len(brands))
		for _, b := range brands {
			brandIDs = append(brandIDs, int64(b.ID))
		}

		// 2. 按日期聚合所有品牌的销售数据
		type MonthlyMarketSum struct {
			Date       time.Time
			TotalSales float64
		}

		var monthlySums []MonthlyMarketSum
		err := tx.Model(&brandtrekin.BtBrandMonthlyTrend{}).
			Select("date, SUM(revenue) as total_sales").
			Where("brand_id IN ?", brandIDs).
			Group("date").
			Order("date").
			Find(&monthlySums).Error

		if err != nil {
			return fmt.Errorf("聚合市场月度数据失败: %v", err)
		}

		// 3. 获取该市场的月度搜索量数据
		// 先获取该市场的所有关键词ID
		var keywordIDs []int64
		err = tx.Model(&brandtrekin.BtKeyword{}).
			Where("market_id = ?", marketID).
			Pluck("id", &keywordIDs).Error

		if err != nil {
			return fmt.Errorf("获取关键词ID失败: %v", err)
		}

		// 创建搜索量映射：key为日期字符串，value为该月所有关键词的搜索量总和
		searchVolumeMap := make(map[string]int)

		if len(keywordIDs) > 0 {
			type MonthlySearchVolume struct {
				Date   time.Time
				Volume int
			}

			var searchVolumes []MonthlySearchVolume
			err = tx.Model(&brandtrekin.BtKeywordMonthlyVolume{}).
				Select("date, SUM(volume) as volume").
				Where("keyword_id IN ?", keywordIDs).
				Group("date").
				Order("date").
				Find(&searchVolumes).Error

			if err != nil {
				return fmt.Errorf("聚合搜索量数据失败: %v", err)
			}

			// 填充搜索量映射：每个月份对应该月的搜索量总和
			for _, sv := range searchVolumes {
				key := sv.Date.Format("2006-01-02")
				searchVolumeMap[key] = sv.Volume
			}
		}

		// 4. 保存或更新市场月度趋势
		for _, sum := range monthlySums {
			date := sum.Date
			revenue := sum.TotalSales

			// 获取该月份对应的搜索量（而不是最新月份的搜索量）
			dateKey := date.Format("2006-01-02")
			searchVolume := searchVolumeMap[dateKey]

			// 转换searchVolume类型
			searchVolumeInt64 := int64(searchVolume)
			trend := brandtrekin.BtMarketMonthlyTrend{
				MarketId:     &marketID,
				Date:         &date,
				Revenue:      &revenue,
				SearchVolume: &searchVolumeInt64,
			}

			// 查找或创建
			var existing brandtrekin.BtMarketMonthlyTrend
			err := tx.Where("market_id = ? AND date = ?", marketID, date).
				First(&existing).Error

			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&trend).Error; err != nil {
					return fmt.Errorf("创建市场月度趋势失败: %v", err)
				}
			} else if err != nil {
				return fmt.Errorf("查询市场月度趋势失败: %v", err)
			} else {
				// 更新
				if err := tx.Model(&existing).Updates(&trend).Error; err != nil {
					return fmt.Errorf("更新市场月度趋势失败: %v", err)
				}
			}
		}

		return nil
	})
}

// CalculateBrandCAGR 计算品牌的CAGR（年复合增长率）
// CAGR = (Ending Value / Beginning Value)^(1/years) - 1
// 使用前12个月和后12个月的平均值来计算，与trekin-main保持一致
func (s *BtAggregateService) CalculateBrandCAGR(marketID int64) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 获取该市场下的所有品牌
		var brands []brandtrekin.BtBrand
		if err := tx.Where("market_id = ?", marketID).Find(&brands).Error; err != nil {
			return fmt.Errorf("查询品牌失败: %v", err)
		}

		for _, brand := range brands {
			brandID := int64(brand.ID)

			// 获取品牌的月度趋势数据（按日期排序）
			var trends []brandtrekin.BtBrandMonthlyTrend
			err := tx.Where("brand_id = ?", brandID).
				Order("date ASC").
				Find(&trends).Error

			if err != nil {
				return fmt.Errorf("查询品牌月度趋势失败: %v", err)
			}

			// 至少需要24个月的数据才能计算CAGR（前12个月 + 后12个月）
			if len(trends) < 24 {
				continue
			}

			// 计算前12个月的平均收入
			var firstYearRevenue float64
			for i := 0; i < 12 && i < len(trends); i++ {
				if trends[i].Revenue != nil {
					firstYearRevenue += *trends[i].Revenue
				}
			}

			// 计算后12个月的平均收入
			var lastYearRevenue float64
			startIdx := len(trends) - 12
			for i := startIdx; i < len(trends); i++ {
				if trends[i].Revenue != nil {
					lastYearRevenue += *trends[i].Revenue
				}
			}

			// 起始值必须大于1000（设置最小阈值避免极端值）
			if firstYearRevenue <= 1000 {
				continue
			}

			// 计算年数：使用前12个月的中点和后12个月的中点之间的时间差
			firstDate := trends[5].Date            // 前12个月的中点
			lastDate := trends[len(trends)-6].Date // 后12个月的中点
			if firstDate == nil || lastDate == nil {
				continue
			}
			years := lastDate.Sub(*firstDate).Hours() / (365.25 * 24)

			// 至少需要半年的时间跨度
			if years < 0.5 {
				continue
			}

			// CAGR公式：(Ending Value / Beginning Value)^(1/years) - 1
			cagr := (math.Pow(lastYearRevenue/firstYearRevenue, 1.0/years) - 1.0) * 100.0

			// 限制在 -99% 到 +999% 之间
			cagr = math.Max(-99.0, math.Min(999.0, cagr))

			// 更新品牌的CAGR
			if err := tx.Model(&brand).Update("cagr", cagr).Error; err != nil {
				return fmt.Errorf("更新品牌CAGR失败: %v", err)
			}
		}

		return nil
	})
}

// CalculateMarketCAGR 计算市场的CAGR
// 使用前12个月和后12个月的总收入来计算，与trekin-main保持一致
func (s *BtAggregateService) CalculateMarketCAGR(marketID int64) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 获取市场
		var market brandtrekin.BtMarket
		if err := tx.First(&market, marketID).Error; err != nil {
			return fmt.Errorf("查询市场失败: %v", err)
		}

		// 获取市场的月度趋势数据
		var trends []brandtrekin.BtMarketMonthlyTrend
		err := tx.Where("market_id = ?", marketID).
			Order("date ASC").
			Find(&trends).Error

		if err != nil {
			return fmt.Errorf("查询市场月度趋势失败: %v", err)
		}

		// 至少需要24个月的数据才能计算CAGR（前12个月 + 后12个月）
		if len(trends) < 24 {
			return nil
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

		// 起始值必须大于1000（设置最小阈值避免极端值）
		if firstYearRevenue <= 1000 {
			return nil
		}

		// 计算年数：使用前12个月的中点和后12个月的中点之间的时间差
		firstDate := trends[5].Date            // 前12个月的中点
		lastDate := trends[len(trends)-6].Date // 后12个月的中点
		if firstDate == nil || lastDate == nil {
			return nil
		}
		years := lastDate.Sub(*firstDate).Hours() / (365.25 * 24)

		// 至少需要半年的时间跨度
		if years < 0.5 {
			return nil
		}

		// CAGR公式：(Ending Value / Beginning Value)^(1/years) - 1
		cagr := (math.Pow(lastYearRevenue/firstYearRevenue, 1.0/years) - 1.0) * 100.0

		// 限制范围
		cagr = math.Max(-99.0, math.Min(999.0, cagr))

		// 更新市场的CAGR
		if err := tx.Model(&market).Update("cagr", cagr).Error; err != nil {
			return fmt.Errorf("更新市场CAGR失败: %v", err)
		}

		return nil
	})
}

// UpdateBrandMetrics 更新品牌的聚合指标
// 使用数据中最新的12个月，而不是从当前时间往前推12个月
func (s *BtAggregateService) UpdateBrandMetrics(marketID int64) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 获取该市场下的所有品牌
		var brands []brandtrekin.BtBrand
		if err := tx.Where("market_id = ?", marketID).Find(&brands).Error; err != nil {
			return fmt.Errorf("查询品牌失败: %v", err)
		}

		for _, brand := range brands {
			brandID := int64(brand.ID)

			// 1. 获取该品牌的所有月度趋势数据（按日期倒序）
			var trends []brandtrekin.BtBrandMonthlyTrend
			err := tx.Where("brand_id = ?", brandID).
				Order("date DESC").
				Limit(12).
				Find(&trends).Error

			if err != nil {
				return fmt.Errorf("查询品牌月度趋势失败: %v", err)
			}

			// 计算最近12个月的总销售额
			var totalRevenue float64
			for _, trend := range trends {
				if trend.Revenue != nil {
					totalRevenue += *trend.Revenue
				}
			}

			// 2. 更新商品数量
			var productCount int64
			err = tx.Model(&brandtrekin.BtProduct{}).
				Where("brand_id = ?", brandID).
				Count(&productCount).Error

			if err != nil {
				return fmt.Errorf("计算品牌商品数量失败: %v", err)
			}

			// 3. 更新品牌记录
			updates := map[string]interface{}{
				"total_revenue": totalRevenue,
				"product_count": productCount,
			}

			if err := tx.Model(&brand).Updates(updates).Error; err != nil {
				return fmt.Errorf("更新品牌指标失败: %v", err)
			}
		}

		return nil
	})
}

// UpdateMarketMetrics 更新市场的聚合指标
// 使用数据中最新的12个月，而不是从当前时间往前推12个月
func (s *BtAggregateService) UpdateMarketMetrics(marketID int64) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 获取市场
		var market brandtrekin.BtMarket
		if err := tx.First(&market, marketID).Error; err != nil {
			return fmt.Errorf("查询市场失败: %v", err)
		}

		// 1. 获取市场的所有月度趋势数据（按日期倒序，取最近12个月）
		var trends []brandtrekin.BtMarketMonthlyTrend
		err := tx.Where("market_id = ?", marketID).
			Order("date DESC").
			Limit(12).
			Find(&trends).Error

		if err != nil {
			return fmt.Errorf("查询市场月度趋势失败: %v", err)
		}

		// 计算最近12个月的总销售额和总搜索量
		var totalRevenue float64
		var searchVolume int64
		for _, trend := range trends {
			if trend.Revenue != nil {
				totalRevenue += *trend.Revenue
			}
			if trend.SearchVolume != nil {
				searchVolume += *trend.SearchVolume
			}
		}

		// 2. 更新品牌数量
		var brandCount int64
		err = tx.Model(&brandtrekin.BtBrand{}).
			Where("market_id = ?", marketID).
			Count(&brandCount).Error

		if err != nil {
			return fmt.Errorf("计算品牌数量失败: %v", err)
		}

		// 3. 更新商品总数
		var productCount int64
		err = tx.Model(&brandtrekin.BtProduct{}).
			Where("market_id = ?", marketID).
			Count(&productCount).Error

		if err != nil {
			return fmt.Errorf("计算商品数量失败: %v", err)
		}

		//// 4. 获取最新月份的搜索量（数据中的最新日期，而不是当前时间）
		//// 注意：trends 已经按 date DESC 排序，所以 trends[0] 是最新月份
		//var searchVolume int64
		//if len(trends) > 0 && trends[0].SearchVolume != nil {
		//	searchVolume = *trends[0].SearchVolume
		//}

		// 5. 更新市场记录
		updates := map[string]interface{}{
			"total_revenue":  totalRevenue,
			"brand_count":    brandCount,
			"total_products": productCount,
			"search_volume":  searchVolume,
		}

		if err := tx.Model(&market).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新市场指标失败: %v", err)
		}

		return nil
	})
}

// RunFullAggregation 运行完整的数据聚合流程
// 此函数会按顺序执行所有聚合操作
func (s *BtAggregateService) RunFullAggregation(marketID int64) error {
	// 1. 聚合商品销售数据到品牌月度趋势
	if err := s.AggregateProductToBrand(marketID); err != nil {
		return fmt.Errorf("商品到品牌聚合失败: %v", err)
	}

	// 2. 聚合品牌月度趋势到市场月度趋势
	if err := s.AggregateBrandToMarket(marketID); err != nil {
		return fmt.Errorf("品牌到市场聚合失败: %v", err)
	}

	// 3. 计算品牌CAGR
	if err := s.CalculateBrandCAGR(marketID); err != nil {
		return fmt.Errorf("品牌CAGR计算失败: %v", err)
	}

	// 4. 计算市场CAGR
	if err := s.CalculateMarketCAGR(marketID); err != nil {
		return fmt.Errorf("市场CAGR计算失败: %v", err)
	}

	// 5. 更新品牌指标
	if err := s.UpdateBrandMetrics(marketID); err != nil {
		return fmt.Errorf("品牌指标更新失败: %v", err)
	}

	// 6. 更新市场指标
	if err := s.UpdateMarketMetrics(marketID); err != nil {
		return fmt.Errorf("市场指标更新失败: %v", err)
	}

	return nil
}
