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
				units := sum.TotalUnits

				trend := brandtrekin.BtBrandMonthlyTrend{
					BrandId: &brandID,
					Date:    &date,
					Revenue: &revenue,
					Units:   &units,
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
			TotalUnits int
		}

		var monthlySums []MonthlyMarketSum
		err := tx.Model(&brandtrekin.BtBrandMonthlyTrend{}).
			Select("date, SUM(revenue) as total_sales, SUM(units) as total_units").
			Where("brand_id IN ?", brandIDs).
			Group("date").
			Order("date").
			Find(&monthlySums).Error

		if err != nil {
			return fmt.Errorf("聚合市场月度数据失败: %v", err)
		}

		// 3. 获取该市场的月度搜索量数据
		type MonthlySearchVolume struct {
			Date   time.Time
			Volume int
		}

		var searchVolumes []MonthlySearchVolume
		err = tx.Model(&brandtrekin.BtKeywordMonthlyVolume{}).
			Select("date, SUM(volume) as volume").
			Where("market_id = ?", marketID).
			Group("date").
			Order("date").
			Find(&searchVolumes).Error

		if err != nil {
			return fmt.Errorf("聚合搜索量数据失败: %v", err)
		}

		// 创建搜索量映射
		searchVolumeMap := make(map[string]int)
		for _, sv := range searchVolumes {
			key := sv.Date.Format("2006-01-02")
			searchVolumeMap[key] = sv.Volume
		}

		// 4. 保存或更新市场月度趋势
		for _, sum := range monthlySums {
			date := sum.Date
			revenue := sum.TotalSales
			units := sum.TotalUnits

			// 获取该月的搜索量
			dateKey := date.Format("2006-01-02")
			searchVolume := searchVolumeMap[dateKey]

			trend := brandtrekin.BtMarketMonthlyTrend{
				MarketId:     &marketID,
				Date:         &date,
				Revenue:      &revenue,
				Units:        &units,
				SearchVolume: &searchVolume,
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

			// 至少需要12个月的数据
			if len(trends) < 12 {
				continue
			}

			// 获取第一个月和最后一个月的收入
			var beginningValue, endingValue float64
			if trends[0].Revenue != nil {
				beginningValue = *trends[0].Revenue
			}
			if trends[len(trends)-1].Revenue != nil {
				endingValue = *trends[len(trends)-1].Revenue
			}

			// 起始值必须大于0
			if beginningValue <= 0 {
				continue
			}

			// 计算年数（月数 / 12）
			years := float64(len(trends)) / 12.0

			// CAGR公式：(Ending Value / Beginning Value)^(1/years) - 1
			cagr := (math.Pow(endingValue/beginningValue, 1.0/years) - 1.0) * 100.0

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

		// 至少需要12个月的数据
		if len(trends) < 12 {
			return nil
		}

		// 获取第一个月和最后一个月的收入
		var beginningValue, endingValue float64
		if trends[0].Revenue != nil {
			beginningValue = *trends[0].Revenue
		}
		if trends[len(trends)-1].Revenue != nil {
			endingValue = *trends[len(trends)-1].Revenue
		}

		// 起始值必须大于0
		if beginningValue <= 0 {
			return nil
		}

		// 计算年数
		years := float64(len(trends)) / 12.0

		// CAGR公式
		cagr := (math.Pow(endingValue/beginningValue, 1.0/years) - 1.0) * 100.0

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
func (s *BtAggregateService) UpdateBrandMetrics(marketID int64) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 获取该市场下的所有品牌
		var brands []brandtrekin.BtBrand
		if err := tx.Where("market_id = ?", marketID).Find(&brands).Error; err != nil {
			return fmt.Errorf("查询品牌失败: %v", err)
		}

		for _, brand := range brands {
			brandID := int64(brand.ID)

			// 1. 更新总销售额（最近12个月）
			var totalRevenue float64
			twelveMonthsAgo := time.Now().AddDate(0, -12, 0)

			err := tx.Model(&brandtrekin.BtBrandMonthlyTrend{}).
				Select("COALESCE(SUM(revenue), 0)").
				Where("brand_id = ? AND date >= ?", brandID, twelveMonthsAgo).
				Scan(&totalRevenue).Error

			if err != nil {
				return fmt.Errorf("计算品牌总销售额失败: %v", err)
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
func (s *BtAggregateService) UpdateMarketMetrics(marketID int64) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 获取市场
		var market brandtrekin.BtMarket
		if err := tx.First(&market, marketID).Error; err != nil {
			return fmt.Errorf("查询市场失败: %v", err)
		}

		// 1. 更新总销售额（最近12个月）
		var totalRevenue float64
		twelveMonthsAgo := time.Now().AddDate(0, -12, 0)

		err := tx.Model(&brandtrekin.BtMarketMonthlyTrend{}).
			Select("COALESCE(SUM(revenue), 0)").
			Where("market_id = ? AND date >= ?", marketID, twelveMonthsAgo).
			Scan(&totalRevenue).Error

		if err != nil {
			return fmt.Errorf("计算市场总销售额失败: %v", err)
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

		// 4. 获取最近月份的搜索量
		var searchVolume int64
		err = tx.Model(&brandtrekin.BtMarketMonthlyTrend{}).
			Select("COALESCE(search_volume, 0)").
			Where("market_id = ?", marketID).
			Order("date DESC").
			Limit(1).
			Scan(&searchVolume).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("获取搜索量失败: %v", err)
		}

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
