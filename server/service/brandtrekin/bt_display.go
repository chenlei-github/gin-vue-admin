package brandtrekin

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/response"
)

type BtDisplayService struct{}

// GetMarketList 获取市场列表
func (s *BtDisplayService) GetMarketList() ([]response.MarketListItem, error) {
	var markets []brandtrekin.BtMarket
	if err := global.GVA_DB.Where("status = ?", "active").Find(&markets).Error; err != nil {
		return nil, fmt.Errorf("查询市场失败: %v", err)
	}

	result := make([]response.MarketListItem, 0, len(markets))

	for _, market := range markets {
		marketID := int64(market.ID)

		// 获取市场的月度趋势数据（最近12个月）
		var trends []brandtrekin.BtMarketMonthlyTrend
		err := global.GVA_DB.Where("market_id = ?", marketID).
			Order("date DESC").
			Limit(12).
			Find(&trends).Error

		if err != nil {
			return nil, fmt.Errorf("查询市场月度趋势失败: %v", err)
		}

		// 反转顺序（从早到晚）
		monthlyTrends := make([]response.MonthlyTrend, len(trends))
		for i, trend := range trends {
			idx := len(trends) - 1 - i
			var searchVolume *int
			if trend.SearchVolume != nil {
				sv := int(*trend.SearchVolume)
				searchVolume = &sv
			}
			monthlyTrends[idx] = response.MonthlyTrend{
				Date:         trend.Date.Format("2006-01"),
				Revenue:      getFloat64(trend.Revenue),
				SearchVolume: searchVolume,
			}
		}

		metrics := response.MarketMetrics{
			TotalRevenue:  getFloat64(market.TotalRevenue),
			TotalProducts: int(getInt64(market.TotalProducts)),
			BrandCount:    int(getInt64(market.BrandCount)),
			SearchVolume:  int(getInt64(market.SearchVolume)),
			CAGR:          market.Cagr,
			MonthlyTrends: monthlyTrends,
		}

		item := response.MarketListItem{
			ID:      getString(market.MarketSlug),
			Name:    getString(market.MarketName),
			Metrics: metrics,
		}

		result = append(result, item)
	}

	return result, nil
}

// GetMarketDetail 获取市场详情
func (s *BtDisplayService) GetMarketDetail(marketSlug string) (*response.MarketDetail, error) {
	// 1. 查找市场
	var market brandtrekin.BtMarket
	if err := global.GVA_DB.Where("market_slug = ? AND status = ?", marketSlug, "active").
		First(&market).Error; err != nil {
		return nil, fmt.Errorf("市场不存在: %v", err)
	}

	marketID := int64(market.ID)

	// 2. 获取月度趋势
	var trends []brandtrekin.BtMarketMonthlyTrend
	err := global.GVA_DB.Where("market_id = ?", marketID).
		Order("date DESC").
		Limit(12).
		Find(&trends).Error

	if err != nil {
		return nil, fmt.Errorf("查询市场月度趋势失败: %v", err)
	}

	monthlyTrends := make([]response.MonthlyTrend, 0, len(trends))
	for i := len(trends) - 1; i >= 0; i-- {
		trend := trends[i]
		var searchVolume *int
		if trend.SearchVolume != nil {
			sv := int(*trend.SearchVolume)
			searchVolume = &sv
		}
		monthlyTrends = append(monthlyTrends, response.MonthlyTrend{
			Date:         trend.Date.Format("2006-01"),
			Revenue:      getFloat64(trend.Revenue),
			SearchVolume: searchVolume,
		})
	}

	// 3. 获取品牌列表
	var brands []brandtrekin.BtBrand
	if err := global.GVA_DB.Where("market_id = ?", marketID).
		Order("total_revenue DESC").
		Find(&brands).Error; err != nil {
		return nil, fmt.Errorf("查询品牌失败: %v", err)
	}

	brandSummaries := make([]response.BrandSummary, 0, len(brands))

	for _, brand := range brands {
		brandID := int64(brand.ID)

		// 获取社交媒体信息
		social, err := s.getBrandSocialMedia(brandID)
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("获取品牌社交媒体失败: %v", err))
		}

		// 获取品牌月度趋势
		var brandTrends []brandtrekin.BtBrandMonthlyTrend
		err = global.GVA_DB.Where("brand_id = ?", brandID).
			Order("date DESC").
			Limit(12).
			Find(&brandTrends).Error

		brandMonthlyTrends := make([]response.MonthlyTrend, 0, len(brandTrends))
		if err == nil {
			for i := len(brandTrends) - 1; i >= 0; i-- {
				trend := brandTrends[i]
				brandMonthlyTrends = append(brandMonthlyTrends, response.MonthlyTrend{
					Date:    trend.Date.Format("2006-01"),
					Revenue: getFloat64(trend.Revenue),
				})
			}
		} else {
			global.GVA_LOG.Warn(fmt.Sprintf("获取品牌月度趋势失败: %v", err))
		}

		summary := response.BrandSummary{
			Brand:         getString(brand.BrandName),
			TotalRevenue:  getFloat64(brand.TotalRevenue),
			ProductCount:  int(getInt64(brand.ProductCount)),
			CAGR:          brand.Cagr,
			Website:       brand.Website,
			Social:        social,
			MonthlyTrends: brandMonthlyTrends,
		}

		brandSummaries = append(brandSummaries, summary)
	}

	// 4. 组装响应
	detail := &response.MarketDetail{
		ID:   marketSlug,
		Name: getString(market.MarketName),
		Metrics: response.MarketMetrics{
			TotalRevenue:  getFloat64(market.TotalRevenue),
			TotalProducts: int(getInt64(market.TotalProducts)),
			BrandCount:    int(getInt64(market.BrandCount)),
			SearchVolume:  int(getInt64(market.SearchVolume)),
			CAGR:          market.Cagr,
			MonthlyTrends: monthlyTrends,
		},
		Brands: brandSummaries,
	}

	return detail, nil
}

// GetBrandDetail 获取品牌详情
func (s *BtDisplayService) GetBrandDetail(marketSlug string, brandName string) (*response.BrandDetail, error) {
	// 1. 查找市场
	var market brandtrekin.BtMarket
	if err := global.GVA_DB.Where("market_slug = ? AND status = ?", marketSlug, "active").
		First(&market).Error; err != nil {
		return nil, fmt.Errorf("市场不存在: %v", err)
	}

	marketID := int64(market.ID)

	// 2. 查找品牌
	var brand brandtrekin.BtBrand
	if err := global.GVA_DB.Where("market_id = ? AND brand_name = ?", marketID, brandName).
		First(&brand).Error; err != nil {
		return nil, fmt.Errorf("品牌不存在: %v", err)
	}

	brandID := int64(brand.ID)

	// 3. 获取品牌月度趋势
	var trends []brandtrekin.BtBrandMonthlyTrend
	err := global.GVA_DB.Where("brand_id = ?", brandID).
		Order("date DESC").
		Limit(12).
		Find(&trends).Error

	if err != nil {
		return nil, fmt.Errorf("查询品牌月度趋势失败: %v", err)
	}

	monthlyTrends := make([]response.MonthlyTrend, 0, len(trends))
	for i := len(trends) - 1; i >= 0; i-- {
		trend := trends[i]
		monthlyTrends = append(monthlyTrends, response.MonthlyTrend{
			Date:    trend.Date.Format("2006-01"),
			Revenue: getFloat64(trend.Revenue),
		})
	}

	// 4. 获取社交媒体信息
	social, err := s.getBrandSocialMedia(brandID)
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("获取品牌社交媒体失败: %v", err))
	}

	// 5. 获取商品列表
	var products []brandtrekin.BtProduct
	if err := global.GVA_DB.Where("brand_id = ?", brandID).
		Order("monthly_sales DESC").
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("查询商品失败: %v", err)
	}

	productDetails := make([]response.ProductDetail, 0, len(products))

	for _, product := range products {
		asin := getString(product.Asin)

		// 获取商品的月度销售趋势
		var salesTrends []brandtrekin.BtProductMonthlySales
		err := global.GVA_DB.Where("asin = ?", asin).
			Order("date DESC").
			Limit(12).
			Find(&salesTrends).Error

		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("获取商品销售趋势失败: %v", err))
		}

		salesTrendList := make([]response.ProductSalesTrend, 0, len(salesTrends))
		for i := len(salesTrends) - 1; i >= 0; i-- {
			trend := salesTrends[i]
			salesTrendList = append(salesTrendList, response.ProductSalesTrend{
				Date:  trend.Date.Format("2006-01"),
				Sales: getFloat64(trend.Sales),
			})
		}

		detail := response.ProductDetail{
			ASIN:         asin,
			Title:        getString(product.Title),
			Price:        getFloat64(product.Price),
			Rating:       getFloat64(product.Rating),
			Reviews:      int(getInt64(product.Reviews)),
			ImageUrl:     product.ImageUrl,
			MonthlySales: int(getInt64(product.MonthlySales)),
			SalesTrend:   salesTrendList,
		}

		productDetails = append(productDetails, detail)
	}

	// 6. 组装响应
	brandDetail := &response.BrandDetail{
		Brand:         brandName,
		TotalRevenue:  getFloat64(brand.TotalRevenue),
		ProductCount:  int(getInt64(brand.ProductCount)),
		CAGR:          brand.Cagr,
		Website:       brand.Website,
		Social:        social,
		MonthlyTrends: monthlyTrends,
		Products:      productDetails,
	}

	return brandDetail, nil
}

// getBrandSocialMedia 获取品牌的社交媒体信息
func (s *BtDisplayService) getBrandSocialMedia(brandID int64) (*response.SocialMedia, error) {
	var socials []brandtrekin.BtBrandSocialMedia
	if err := global.GVA_DB.Where("brand_id = ?", brandID).Find(&socials).Error; err != nil {
		return nil, err
	}

	if len(socials) == 0 {
		return nil, nil
	}

	social := &response.SocialMedia{}

	for _, s := range socials {
		platform := getString(s.Platform)
		url := getString(s.Url)

		switch platform {
		case "youtube":
			var subscribers *int
			if s.Subscribers != nil {
				subs := int(*s.Subscribers)
				subscribers = &subs
			}
			social.Youtube = &response.SocialPlatform{
				URL:         url,
				Subscribers: subscribers,
			}
		case "instagram":
			var followers *int
			if s.Followers != nil {
				fol := int(*s.Followers)
				followers = &fol
			}
			social.Instagram = &response.SocialPlatform{
				URL:       url,
				Followers: followers,
			}
		case "facebook":
			var followers *int
			if s.Followers != nil {
				fol := int(*s.Followers)
				followers = &fol
			}
			social.Facebook = &response.SocialPlatform{
				URL:       url,
				Followers: followers,
			}
		case "reddit":
			var posts *int
			if s.Posts != nil {
				p := int(*s.Posts)
				posts = &p
			}
			social.Reddit = &response.SocialPlatform{
				URL:   url,
				Posts: posts,
			}
		}
	}

	return social, nil
}

// Helper functions
func getString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func getFloat64(ptr *float64) float64 {
	if ptr == nil {
		return 0.0
	}
	return *ptr
}

func getInt64(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
