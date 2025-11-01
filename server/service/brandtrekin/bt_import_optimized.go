package brandtrekin

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
	"gorm.io/gorm"
)

const (
	// 批量插入批次大小
	batchSize = 500
)

// BatchImport 优化后的批量导入方法
// 1. 使用协程并发解析文件
// 2. 拆分大事务为多个小事务
// 3. 使用批量插入优化数据库操作
// 4. 将聚合计算独立执行
func (s *BtImportService) BatchImport(marketID int64, files map[string]*multipart.FileHeader, replaceMode bool) error {
	// 如果是全量替换模式，先删除现有数据（独立事务）
	if replaceMode {
		if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			return s.deleteMarketData(tx, marketID)
		}); err != nil {
			return fmt.Errorf("删除现有数据失败: %v", err)
		}
	}

	// 1. 并发解析所有文件
	parseResults := make(map[string]interface{})
	parseErrors := make(map[string]error)
	var parseWg sync.WaitGroup
	var parseMu sync.Mutex

	for key, file := range files {
		parseWg.Add(1)
		go func(fileKey string, fileHeader *multipart.FileHeader) {
			defer parseWg.Done()

			var result *ImportResult
			var err error

			switch fileKey {
			case "brandSocial":
				result, err = s.ParseBrandSocial(fileHeader)
			case "productUS":
				result, err = s.ParseProductUS(fileHeader)
			case "gkw":
				result, err = s.ParseGKW(fileHeader)
			case "keywordHistory":
				result, err = s.ParseKeywordHistory(fileHeader)
			case "productSales":
				result, err = s.ParseProductSales(fileHeader)
			default:
				err = fmt.Errorf("未知的文件类型: %s", fileKey)
			}

			parseMu.Lock()
			if err != nil {
				parseErrors[fileKey] = err
			} else if result != nil {
				parseResults[fileKey] = result.FullData
			}
			parseMu.Unlock()
		}(key, file)
	}

	parseWg.Wait()

	// 检查解析错误
	for key, err := range parseErrors {
		return fmt.Errorf("解析%s文件失败: %v", key, err)
	}

	// 2. 按顺序保存数据（每个保存操作使用独立的小事务，避免长时间锁等待）
	// 2.1 保存品牌社交媒体数据
	if data, ok := parseResults["brandSocial"]; ok {
		if brands, ok := data.([]BrandWithSocial); ok {
			if err := s.SaveBrandSocialDataBatch(marketID, brands); err != nil {
				return fmt.Errorf("保存品牌社交媒体数据失败: %v", err)
			}
		}
	}

	// 2.2 保存商品数据
	if data, ok := parseResults["productUS"]; ok {
		if products, ok := data.([]ProductInfo); ok {
			if err := s.SaveProductDataBatch(marketID, products); err != nil {
				return fmt.Errorf("保存商品数据失败: %v", err)
			}
		}
	}

	// 2.3 保存Google关键词数据
	if data, ok := parseResults["gkw"]; ok {
		if keywords, ok := data.([]KeywordInfo); ok {
			if err := s.SaveKeywordDataBatch(marketID, keywords); err != nil {
				return fmt.Errorf("保存Google关键词数据失败: %v", err)
			}
		}
	}

	// 2.4 保存Amazon关键词数据
	if data, ok := parseResults["keywordHistory"]; ok {
		if keywords, ok := data.([]KeywordInfo); ok {
			if err := s.SaveKeywordDataBatch(marketID, keywords); err != nil {
				return fmt.Errorf("保存Amazon关键词数据失败: %v", err)
			}
		}
	}

	// 2.5 保存商品销售数据
	if data, ok := parseResults["productSales"]; ok {
		if salesData, ok := data.([]ProductSalesInfo); ok {
			if err := s.SaveProductSalesDataBatch(salesData); err != nil {
				return fmt.Errorf("保存商品销售数据失败: %v", err)
			}
		}
	}

	// 3. 导入完成后，独立执行数据聚合计算（不阻塞导入）
	// 使用异步方式执行，避免阻塞导入流程
	go func() {
		aggregateService := BtAggregateService{}
		if err := aggregateService.RunFullAggregation(marketID); err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("数据聚合计算失败: %v", err))
		}
	}()

	return nil
}

// SaveBrandSocialDataBatch 批量保存品牌社交媒体数据（使用批量插入优化）
func (s *BtImportService) SaveBrandSocialDataBatch(marketID int64, brands []BrandWithSocial) error {
	if len(brands) == 0 {
		return nil
	}

	// 分批处理
	for i := 0; i < len(brands); i += batchSize {
		end := i + batchSize
		if end > len(brands) {
			end = len(brands)
		}
		batch := brands[i:end]

		// 每个批次使用独立事务
		if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			// 1. 批量查询现有品牌（减少数据库查询）
			brandNames := make([]string, len(batch))
			for j, b := range batch {
				brandNames[j] = b.BrandName
			}

			var existingBrands []brandtrekin.BtBrand
			tx.Where("market_id = ? AND brand_name IN ?", marketID, brandNames).
				Find(&existingBrands)

			// 创建品牌名到ID的映射
			brandMap := make(map[string]int64)
			for _, b := range existingBrands {
				if b.BrandName != nil {
					brandMap[*b.BrandName] = int64(b.ID)
				}
			}

			// 2. 准备批量插入和更新的品牌数据
			toCreate := []brandtrekin.BtBrand{}
			toUpdate := []brandtrekin.BtBrand{}

			for _, brandData := range batch {
				brand := brandtrekin.BtBrand{
					MarketId:  &marketID,
					BrandName: &brandData.BrandName,
				}
				if brandData.Website != "" {
					brand.Website = &brandData.Website
				}

				if brandID, exists := brandMap[brandData.BrandName]; exists {
					brand.ID = uint(brandID)
					toUpdate = append(toUpdate, brand)
				} else {
					toCreate = append(toCreate, brand)
				}
			}

			// 3. 批量创建新品牌
			if len(toCreate) > 0 {
				if err := tx.CreateInBatches(toCreate, batchSize).Error; err != nil {
					return fmt.Errorf("批量创建品牌失败: %v", err)
				}

				// 重新查询新创建的品牌ID
				var newBrands []brandtrekin.BtBrand
				if err := tx.Where("market_id = ? AND brand_name IN ?", marketID, brandNames).
					Find(&newBrands).Error; err != nil {
					return fmt.Errorf("查询新创建品牌失败: %v", err)
				}

				for _, b := range newBrands {
					if b.BrandName != nil {
						brandMap[*b.BrandName] = int64(b.ID)
					}
				}
			}

			// 4. 批量更新现有品牌
			if len(toUpdate) > 0 {
				for _, brand := range toUpdate {
					if err := tx.Model(&brandtrekin.BtBrand{}).
						Where("id = ?", brand.ID).
						Updates(map[string]interface{}{
							"website": brand.Website,
						}).Error; err != nil {
						return fmt.Errorf("更新品牌失败: %v", err)
					}
				}
			}

			// 5. 批量处理社交媒体数据
			toCreateSocial := []brandtrekin.BtBrandSocialMedia{}

			for _, brandData := range batch {
				brandID := brandMap[brandData.BrandName]
				socialMedia := []struct {
					platform    string
					url         string
					subscribers int
					followers   int
					posts       int
				}{
					{"youtube", brandData.YoutubeUrl, brandData.YoutubeSubscribers, 0, 0},
					{"instagram", brandData.InstagramUrl, 0, brandData.InstagramFollowers, 0},
					{"facebook", brandData.FacebookUrl, 0, brandData.FacebookFollowers, 0},
					{"reddit", brandData.RedditUrl, 0, 0, brandData.RedditPosts},
				}

				for _, sm := range socialMedia {
					if sm.url == "" {
						continue
					}

					platform := sm.platform
					url := sm.url
					social := brandtrekin.BtBrandSocialMedia{
						BrandId:  &brandID,
						Platform: &platform,
						Url:      &url,
					}

					if sm.subscribers > 0 {
						subs := int64(sm.subscribers)
						social.Subscribers = &subs
					}
					if sm.followers > 0 {
						fol := int64(sm.followers)
						social.Followers = &fol
					}
					if sm.posts > 0 {
						p := int64(sm.posts)
						social.Posts = &p
					}

					toCreateSocial = append(toCreateSocial, social)
				}
			}

			// 使用 ON DUPLICATE KEY UPDATE 或批量插入/更新
			// 由于GORM不支持复杂的upsert，我们使用批量删除+插入的方式
			if len(toCreateSocial) > 0 {
				// 先批量删除可能存在的记录
				var brandIDs []int64
				for _, social := range toCreateSocial {
					brandIDs = append(brandIDs, *social.BrandId)
				}
				tx.Where("brand_id IN ?", brandIDs).
					Delete(&brandtrekin.BtBrandSocialMedia{})

				// 批量插入
				if err := tx.CreateInBatches(toCreateSocial, batchSize).Error; err != nil {
					return fmt.Errorf("批量创建社交媒体记录失败: %v", err)
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// SaveProductDataBatch 批量保存商品数据
func (s *BtImportService) SaveProductDataBatch(marketID int64, products []ProductInfo) error {
	if len(products) == 0 {
		return nil
	}

	// 分批处理
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}
		batch := products[i:end]

		if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			// 1. 批量查询品牌ID
			brandNames := make([]string, 0, len(batch))
			for _, p := range batch {
				brandNames = append(brandNames, p.Brand)
			}

			var brands []brandtrekin.BtBrand
			tx.Where("market_id = ? AND brand_name IN ?", marketID, brandNames).
				Find(&brands)

			brandMap := make(map[string]int64)
			for _, b := range brands {
				if b.BrandName != nil {
					brandMap[*b.BrandName] = int64(b.ID)
				}
			}

			// 2. 创建不存在的品牌
			for _, p := range batch {
				if _, exists := brandMap[p.Brand]; !exists {
					brand := brandtrekin.BtBrand{
						MarketId:  &marketID,
						BrandName: &p.Brand,
					}
					if err := tx.Create(&brand).Error; err != nil {
						return fmt.Errorf("创建品牌失败: %v", err)
					}
					brandMap[p.Brand] = int64(brand.ID)
				}
			}

			// 3. 批量准备商品数据
			toCreate := []brandtrekin.BtProduct{}
			toUpdate := []brandtrekin.BtProduct{}

			asins := make([]string, len(batch))
			for j, p := range batch {
				asins[j] = p.ASIN
			}

			var existingProducts []brandtrekin.BtProduct
			tx.Where("asin IN ?", asins).Find(&existingProducts)

			existingMap := make(map[string]bool)
			for _, ep := range existingProducts {
				if ep.Asin != nil {
					existingMap[*ep.Asin] = true
				}
			}

			for _, p := range batch {
				brandID := brandMap[p.Brand]
				reviews := int64(p.Reviews)
				monthlySales := int64(p.MonthlySales)
				
				// 处理扩展数据JSON
				var extendedDataJSON string
				if len(p.ExtendedData) > 0 {
					if jsonBytes, err := json.Marshal(p.ExtendedData); err == nil {
						extendedDataJSON = string(jsonBytes)
					} else {
						// 如果JSON序列化失败，使用空JSON对象
						extendedDataJSON = "{}"
					}
				} else {
					extendedDataJSON = "{}"
				}
				
				product := brandtrekin.BtProduct{
					MarketId:     &marketID,
					BrandId:      &brandID,
					Asin:         &p.ASIN,
					Title:        &p.Title,
					Price:        &p.Price,
					Rating:       &p.Rating,
					Reviews:      &reviews,
					MonthlySales: &monthlySales,
					ImageUrl:     p.ImageUrl,
					ExtendedData: extendedDataJSON,
				}

				if existingMap[p.ASIN] {
					toUpdate = append(toUpdate, product)
				} else {
					toCreate = append(toCreate, product)
				}
			}

			// 4. 批量创建
			if len(toCreate) > 0 {
				if err := tx.CreateInBatches(toCreate, batchSize).Error; err != nil {
					return fmt.Errorf("批量创建商品失败: %v", err)
				}
			}

			// 5. 批量更新
			if len(toUpdate) > 0 {
				for _, product := range toUpdate {
					if err := tx.Model(&brandtrekin.BtProduct{}).
						Where("asin = ?", product.Asin).
						Updates(&product).Error; err != nil {
						return fmt.Errorf("更新商品失败: %v", err)
					}
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// SaveKeywordDataBatch 批量保存关键词数据
func (s *BtImportService) SaveKeywordDataBatch(marketID int64, keywords []KeywordInfo) error {
	if len(keywords) == 0 {
		return nil
	}

	// 分批处理
	for i := 0; i < len(keywords); i += batchSize {
		end := i + batchSize
		if end > len(keywords) {
			end = len(keywords)
		}
		batch := keywords[i:end]

		if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			// 1. 批量查询现有关键词
			keywordKeys := make([]struct {
				keyword string
				source  string
			}, len(batch))
			for j, k := range batch {
				keywordKeys[j] = struct {
					keyword string
					source  string
				}{k.Keyword, k.Source}
			}

			// 使用 OR 条件批量查询现有关键词
			var existingKeywords []brandtrekin.BtKeyword
			if len(batch) > 0 {
				// 构建查询条件
				var conditions []string
				var args []interface{}
				args = append(args, marketID)

				for _, k := range batch {
					conditions = append(conditions, "(keyword = ? AND source = ?)")
					args = append(args, k.Keyword, k.Source)
				}

				query := fmt.Sprintf("market_id = ? AND (%s)", strings.Join(conditions, " OR "))
				tx.Where(query, args...).Find(&existingKeywords)
			}

			keywordMap := make(map[string]int64)
			for _, k := range existingKeywords {
				key := fmt.Sprintf("%s|%s", *k.Keyword, *k.Source)
				keywordMap[key] = int64(k.ID)
			}

			// 2. 批量创建关键词
			toCreate := []brandtrekin.BtKeyword{}
			for _, k := range batch {
				key := fmt.Sprintf("%s|%s", k.Keyword, k.Source)
				if _, exists := keywordMap[key]; !exists {
					source := k.Source
					keyword := k.Keyword
					kw := brandtrekin.BtKeyword{
						MarketId: &marketID,
						Keyword:  &keyword,
						Source:   &source,
					}
					toCreate = append(toCreate, kw)
				}
			}

			if len(toCreate) > 0 {
				if err := tx.CreateInBatches(toCreate, batchSize).Error; err != nil {
					return fmt.Errorf("批量创建关键词失败: %v", err)
				}

				// 重新查询关键词ID
				for _, k := range batch {
					key := fmt.Sprintf("%s|%s", k.Keyword, k.Source)
					if _, exists := keywordMap[key]; !exists {
						var kw brandtrekin.BtKeyword
						if err := tx.Where("market_id = ? AND keyword = ? AND source = ?", marketID, k.Keyword, k.Source).
							First(&kw).Error; err == nil {
							keywordMap[key] = int64(kw.ID)
						}
					}
				}
			}

			// 3. 批量保存月度搜索量
			allVolumes := []brandtrekin.BtKeywordMonthlyVolume{}
			for _, k := range batch {
				key := fmt.Sprintf("%s|%s", k.Keyword, k.Source)
				keywordID := keywordMap[key]

				for _, mv := range k.MonthlyVolumes {
					volumeInt64 := int64(mv.Volume)
					allVolumes = append(allVolumes, brandtrekin.BtKeywordMonthlyVolume{
						KeywordId: &keywordID,
						Date:      &mv.Date,
						Volume:    &volumeInt64,
					})
				}
			}

			// 使用批量删除+插入的方式
			if len(allVolumes) > 0 {
				var keywordIDs []int64
				for _, v := range allVolumes {
					keywordIDs = append(keywordIDs, *v.KeywordId)
				}
				tx.Where("keyword_id IN ?", keywordIDs).
					Delete(&brandtrekin.BtKeywordMonthlyVolume{})

				if err := tx.CreateInBatches(allVolumes, batchSize).Error; err != nil {
					return fmt.Errorf("批量创建月度搜索量失败: %v", err)
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// SaveProductSalesDataBatch 批量保存商品销售数据
func (s *BtImportService) SaveProductSalesDataBatch(salesData []ProductSalesInfo) error {
	if len(salesData) == 0 {
		return nil
	}

	// 分批处理
	for i := 0; i < len(salesData); i += batchSize {
		end := i + batchSize
		if end > len(salesData) {
			end = len(salesData)
		}
		batch := salesData[i:end]

		if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			// 收集所有ASIN
			asins := make([]string, len(batch))
			for j, s := range batch {
				asins[j] = s.ASIN
			}

			// 批量准备销售数据
			allSales := []brandtrekin.BtProductMonthlySales{}
			for _, s := range batch {
				for _, ms := range s.MonthlySales {
					unitsInt64 := int64(ms.Units)
					asin := s.ASIN
					allSales = append(allSales, brandtrekin.BtProductMonthlySales{
						Asin:  &asin,
						Date:  &ms.Date,
						Sales: &ms.Sales,
						Units: &unitsInt64,
					})
				}
			}

			// 使用批量删除+插入的方式
			if len(allSales) > 0 {
				tx.Where("asin IN ?", asins).
					Delete(&brandtrekin.BtProductMonthlySales{})

				if err := tx.CreateInBatches(allSales, batchSize).Error; err != nil {
					return fmt.Errorf("批量创建月度销售记录失败: %v", err)
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
