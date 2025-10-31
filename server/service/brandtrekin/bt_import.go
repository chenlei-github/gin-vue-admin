package brandtrekin

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type BtImportService struct{}

// ImportResult 导入结果
type ImportResult struct {
	Success  bool        `json:"success"`
	Total    int         `json:"total"`
	Errors   []RowError  `json:"errors"`
	Preview  interface{} `json:"preview"`  // 预览数据（前5条）
	FullData interface{} `json:"-"`        // 完整数据（不序列化到JSON）
}

// RowError 行错误
type RowError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

// BrandWithSocial 品牌及其社交媒体信息
type BrandWithSocial struct {
	BrandName         string
	Website           string
	YoutubeUrl        string
	YoutubeSubscribers int
	InstagramUrl      string
	InstagramFollowers int
	FacebookUrl       string
	FacebookFollowers int
	RedditUrl         string
	RedditPosts       int
}

// ProductInfo 商品信息
type ProductInfo struct {
	ASIN         string
	Title        string
	Brand        string
	Price        float64
	Rating       float64
	Reviews      int
	ImageUrl     string
	MonthlySales int
}

// KeywordInfo 关键词信息
type KeywordInfo struct {
	Keyword        string
	Source         string // google 或 amazon
	MonthlyVolumes []MonthlyVolume
}

// MonthlyVolume 月度数据
type MonthlyVolume struct {
	Date   time.Time
	Volume int
}

// ProductSalesInfo 商品销售信息
type ProductSalesInfo struct {
	ASIN         string
	MonthlySales []MonthlySales
}

// MonthlySales 月度销售
type MonthlySales struct {
	Date  time.Time
	Sales float64
	Units int
}

// ParseBrandSocial 解析品牌社交媒体数据
func (s *BtImportService) ParseBrandSocial(file *multipart.FileHeader) (*ImportResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("读取Excel文件失败: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}

	if len(rows) < 2 {
		return nil, errors.New("文件为空或缺少数据行")
	}

	header := rows[0]
	colMap := make(map[string]int)
	for i, col := range header {
		colMap[col] = i
	}

	if _, ok := colMap["Brand"]; !ok {
		return nil, errors.New("缺少必填列: Brand")
	}

	var brands []BrandWithSocial
	var rowErrors []RowError

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNumber := i + 1

		if len(row) == 0 || row[colMap["Brand"]] == "" {
			continue
		}

		brandName := strings.TrimSpace(row[colMap["Brand"]])
		if brandName == "" {
			continue
		}

		brand := BrandWithSocial{
			BrandName: brandName,
		}

		if idx, ok := colMap["Website"]; ok && idx < len(row) {
			brand.Website = strings.TrimSpace(row[idx])
		}

		if idx, ok := colMap["YouTube"]; ok && idx < len(row) {
			brand.YoutubeUrl = strings.TrimSpace(row[idx])
		}
		if idx, ok := colMap["YouTube Subscribers"]; ok && idx < len(row) && row[idx] != "" {
			if subs, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
				brand.YoutubeSubscribers = subs
			}
		}

		if idx, ok := colMap["Instagram"]; ok && idx < len(row) {
			brand.InstagramUrl = strings.TrimSpace(row[idx])
		}
		if idx, ok := colMap["Instagram Followers"]; ok && idx < len(row) && row[idx] != "" {
			if followers, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
				brand.InstagramFollowers = followers
			}
		}

		if idx, ok := colMap["Facebook"]; ok && idx < len(row) {
			brand.FacebookUrl = strings.TrimSpace(row[idx])
		}
		if idx, ok := colMap["Facebook Followers"]; ok && idx < len(row) && row[idx] != "" {
			if followers, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
				brand.FacebookFollowers = followers
			}
		}

		if idx, ok := colMap["Reddit"]; ok && idx < len(row) {
			brand.RedditUrl = strings.TrimSpace(row[idx])
		}
		if idx, ok := colMap["Reddit Posts"]; ok && idx < len(row) && row[idx] != "" {
			if posts, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
				brand.RedditPosts = posts
			}
		}

		brands = append(brands, brand)
	}

	preview := brands
	if len(preview) > 5 {
		preview = preview[:5]
	}

	return &ImportResult{
		Success:  true,
		Total:    len(brands),
		Errors:   rowErrors,
		Preview:  preview,
		FullData: brands,
	}, nil
}

// ParseGKW 解析Google关键词数据
func (s *BtImportService) ParseGKW(file *multipart.FileHeader) (*ImportResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	reader := csv.NewReader(src)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取CSV文件失败: %v", err)
	}

	if len(rows) < 2 {
		return nil, errors.New("文件为空或缺少数据行")
	}

	header := rows[0]
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}$`)
	var keywords []KeywordInfo
	var rowErrors []RowError

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNumber := i + 1

		if len(row) == 0 || row[0] == "" {
			continue
		}

		keyword := strings.TrimSpace(row[0])
		if keyword == "" {
			continue
		}

		var monthlyVolumes []MonthlyVolume

		for j := 1; j < len(header) && j < len(row); j++ {
			dateStr := strings.TrimSpace(header[j])

			if !datePattern.MatchString(dateStr) {
				rowErrors = append(rowErrors, RowError{
					Row:   rowNumber,
					Error: fmt.Sprintf("无效的日期格式: %s", dateStr),
				})
				continue
			}

			if row[j] == "" {
				continue
			}

			volume, err := strconv.Atoi(strings.TrimSpace(row[j]))
			if err != nil || volume <= 0 {
				continue
			}

			date, err := time.Parse("2006-01", dateStr)
			if err != nil {
				continue
			}

			monthlyVolumes = append(monthlyVolumes, MonthlyVolume{
				Date:   date,
				Volume: volume,
			})
		}

		if len(monthlyVolumes) > 0 {
			keywords = append(keywords, KeywordInfo{
				Keyword:        keyword,
				Source:         "google",
				MonthlyVolumes: monthlyVolumes,
			})
		}
	}

	preview := keywords
	if len(preview) > 5 {
		preview = preview[:5]
	}

	return &ImportResult{
		Success:  true,
		Total:    len(keywords),
		Errors:   rowErrors,
		Preview:  preview,
		FullData: keywords,
	}, nil
}

// ParseKeywordHistory 解析Amazon关键词历史数据
func (s *BtImportService) ParseKeywordHistory(file *multipart.FileHeader) (*ImportResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("读取Excel文件失败: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}

	if len(rows) < 2 {
		return nil, errors.New("文件为空或缺少数据行")
	}

	header := rows[0]
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}$`)
	var keywords []KeywordInfo
	var rowErrors []RowError

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNumber := i + 1

		if len(row) == 0 || row[0] == "" {
			continue
		}

		keyword := strings.TrimSpace(row[0])
		if keyword == "" {
			continue
		}

		var monthlyVolumes []MonthlyVolume

		for j := 1; j < len(header) && j < len(row); j++ {
			dateStr := strings.TrimSpace(header[j])

			if !datePattern.MatchString(dateStr) {
				rowErrors = append(rowErrors, RowError{
					Row:   rowNumber,
					Error: fmt.Sprintf("无效的日期格式: %s", dateStr),
				})
				continue
			}

			if row[j] == "" {
				continue
			}

			volume, err := strconv.Atoi(strings.TrimSpace(row[j]))
			if err != nil || volume <= 0 {
				continue
			}

			date, err := time.Parse("2006-01", dateStr)
			if err != nil {
				continue
			}

			monthlyVolumes = append(monthlyVolumes, MonthlyVolume{
				Date:   date,
				Volume: volume,
			})
		}

		if len(monthlyVolumes) > 0 {
			keywords = append(keywords, KeywordInfo{
				Keyword:        keyword,
				Source:         "amazon",
				MonthlyVolumes: monthlyVolumes,
			})
		}
	}

	preview := keywords
	if len(preview) > 5 {
		preview = preview[:5]
	}

	return &ImportResult{
		Success:  true,
		Total:    len(keywords),
		Errors:   rowErrors,
		Preview:  preview,
		FullData: keywords,
	}, nil
}

// ParseProductUS 解析美国商品数据
func (s *BtImportService) ParseProductUS(file *multipart.FileHeader) (*ImportResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("读取Excel文件失败: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}

	if len(rows) < 2 {
		return nil, errors.New("文件为空或缺少数据行")
	}

	header := rows[0]
	colMap := make(map[string]int)
	for i, col := range header {
		colMap[col] = i
	}

	requiredCols := []string{"ASIN", "Title", "Brand"}
	for _, col := range requiredCols {
		if _, ok := colMap[col]; !ok {
			return nil, fmt.Errorf("缺少必填列: %s", col)
		}
	}

	var products []ProductInfo
	var rowErrors []RowError

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNumber := i + 1

		if colMap["ASIN"] >= len(row) || colMap["Title"] >= len(row) || colMap["Brand"] >= len(row) {
			rowErrors = append(rowErrors, RowError{
				Row:   rowNumber,
				Error: "缺少必填字段（ASIN、Title或Brand）",
			})
			continue
		}

		asin := strings.TrimSpace(row[colMap["ASIN"]])
		title := strings.TrimSpace(row[colMap["Title"]])
		brand := strings.TrimSpace(row[colMap["Brand"]])

		if asin == "" || title == "" || brand == "" {
			rowErrors = append(rowErrors, RowError{
				Row:   rowNumber,
				Error: "ASIN、Title或Brand不能为空",
			})
			continue
		}

		product := ProductInfo{
			ASIN:  asin,
			Title: title,
			Brand: brand,
		}

		if idx, ok := colMap["Price"]; ok && idx < len(row) && row[idx] != "" {
			if price, err := strconv.ParseFloat(strings.TrimSpace(row[idx]), 64); err == nil {
				product.Price = price
			}
		}

		if idx, ok := colMap["Rating"]; ok && idx < len(row) && row[idx] != "" {
			if rating, err := strconv.ParseFloat(strings.TrimSpace(row[idx]), 64); err == nil {
				product.Rating = rating
			}
		}

		if idx, ok := colMap["Reviews"]; ok && idx < len(row) && row[idx] != "" {
			if reviews, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
				product.Reviews = reviews
			}
		}

		if idx, ok := colMap["Image"]; ok && idx < len(row) {
			product.ImageUrl = strings.TrimSpace(row[idx])
		}

		if idx, ok := colMap["Monthly Sales"]; ok && idx < len(row) && row[idx] != "" {
			if sales, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
				product.MonthlySales = sales
			}
		}

		products = append(products, product)
	}

	preview := products
	if len(preview) > 5 {
		preview = preview[:5]
	}

	return &ImportResult{
		Success:  true,
		Total:    len(products),
		Errors:   rowErrors,
		Preview:  preview,
		FullData: products,
	}, nil
}

// ParseProductSales 解析商品销售数据
func (s *BtImportService) ParseProductSales(file *multipart.FileHeader) (*ImportResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("读取Excel文件失败: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}

	if len(rows) < 2 {
		return nil, errors.New("文件为空或缺少数据行")
	}

	header := rows[0]
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}$`)
	var productSalesData []ProductSalesInfo
	var rowErrors []RowError

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNumber := i + 1

		if len(row) == 0 || row[0] == "" {
			continue
		}

		asin := strings.TrimSpace(row[0])
		if asin == "" {
			continue
		}

		var monthlySales []MonthlySales

		for j := 1; j < len(header) && j < len(row); j++ {
			dateStr := strings.TrimSpace(header[j])

			if !datePattern.MatchString(dateStr) {
				rowErrors = append(rowErrors, RowError{
					Row:   rowNumber,
					Error: fmt.Sprintf("无效的日期格式: %s", dateStr),
				})
				continue
			}

			if row[j] == "" {
				continue
			}

			sales, err := strconv.ParseFloat(strings.TrimSpace(row[j]), 64)
			if err != nil || sales <= 0 {
				continue
			}

			date, err := time.Parse("2006-01", dateStr)
			if err != nil {
				continue
			}

			monthlySales = append(monthlySales, MonthlySales{
				Date:  date,
				Sales: sales,
				Units: 0,
			})
		}

		if len(monthlySales) > 0 {
			productSalesData = append(productSalesData, ProductSalesInfo{
				ASIN:         asin,
				MonthlySales: monthlySales,
			})
		}
	}

	preview := productSalesData
	if len(preview) > 5 {
		preview = preview[:5]
	}

	return &ImportResult{
		Success:  true,
		Total:    len(productSalesData),
		Errors:   rowErrors,
		Preview:  preview,
		FullData: productSalesData,
	}, nil
}

// SaveBrandSocialData 保存品牌和社交媒体数据
func (s *BtImportService) SaveBrandSocialData(tx *gorm.DB, marketID int64, brands []BrandWithSocial) error {
	for _, brandData := range brands {
		// 创建或更新品牌
		brand := brandtrekin.BtBrand{
			MarketId:  &marketID,
			BrandName: &brandData.BrandName,
		}
		if brandData.Website != "" {
			brand.Website = &brandData.Website
		}

		// 查找或创建品牌
		var existingBrand brandtrekin.BtBrand
		err := tx.Where("market_id = ? AND brand_name = ?", marketID, brandData.BrandName).
			First(&existingBrand).Error

		var brandID int64
		if err == gorm.ErrRecordNotFound {
			if err := tx.Create(&brand).Error; err != nil {
				return fmt.Errorf("创建品牌失败: %v", err)
			}
			brandID = int64(brand.ID)
		} else if err != nil {
			return fmt.Errorf("查询品牌失败: %v", err)
		} else {
			brandID = int64(existingBrand.ID)
			// 更新品牌信息
			if err := tx.Model(&existingBrand).Updates(&brand).Error; err != nil {
				return fmt.Errorf("更新品牌失败: %v", err)
			}
		}

		// 创建社交媒体记录
		socialMedia := []struct {
			platform   string
			url        string
			subscribers int
			followers  int
			posts      int
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
				social.Subscribers = &sm.subscribers
			}
			if sm.followers > 0 {
				social.Followers = &sm.followers
			}
			if sm.posts > 0 {
				social.Posts = &sm.posts
			}

			// 查找或创建社交媒体记录
			var existingSocial brandtrekin.BtBrandSocialMedia
			err := tx.Where("brand_id = ? AND platform = ?", brandID, platform).
				First(&existingSocial).Error

			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&social).Error; err != nil {
					return fmt.Errorf("创建社交媒体记录失败: %v", err)
				}
			} else if err != nil {
				return fmt.Errorf("查询社交媒体记录失败: %v", err)
			} else {
				// 更新现有记录
				if err := tx.Model(&existingSocial).Updates(&social).Error; err != nil {
					return fmt.Errorf("更新社交媒体记录失败: %v", err)
				}
			}
		}
	}

	return nil
}

// SaveProductData 保存商品数据
func (s *BtImportService) SaveProductData(tx *gorm.DB, marketID int64, products []ProductInfo) error {
	for _, productData := range products {
		// 查找品牌ID
		var brand brandtrekin.BtBrand
		err := tx.Where("market_id = ? AND brand_name = ?", marketID, productData.Brand).
			First(&brand).Error

		if err == gorm.ErrRecordNotFound {
			// 品牌不存在，创建新品牌
			brand = brandtrekin.BtBrand{
				MarketId:  &marketID,
				BrandName: &productData.Brand,
			}
			if err := tx.Create(&brand).Error; err != nil {
				return fmt.Errorf("创建品牌失败: %v", err)
			}
		} else if err != nil {
			return fmt.Errorf("查询品牌失败: %v", err)
		}

		brandID := int64(brand.ID)

		// 创建或更新商品
		reviews := int64(productData.Reviews)
		monthlySales := int64(productData.MonthlySales)
		product := brandtrekin.BtProduct{
			MarketId:     &marketID,
			BrandId:      &brandID,
			Asin:         &productData.ASIN,
			Title:        &productData.Title,
			Price:        &productData.Price,
			Rating:       &productData.Rating,
			Reviews:      &reviews,
			MonthlySales: &monthlySales,
			ImageUrl:     productData.ImageUrl,
		}

		// 查找或创建商品
		var existingProduct brandtrekin.BtProduct
		err = tx.Where("asin = ?", productData.ASIN).First(&existingProduct).Error

		if err == gorm.ErrRecordNotFound {
			if err := tx.Create(&product).Error; err != nil {
				return fmt.Errorf("创建商品失败: %v", err)
			}
		} else if err != nil {
			return fmt.Errorf("查询商品失败: %v", err)
		} else {
			// 更新商品
			if err := tx.Model(&existingProduct).Updates(&product).Error; err != nil {
				return fmt.Errorf("更新商品失败: %v", err)
			}
		}
	}

	return nil
}

// SaveKeywordData 保存关键词数据
func (s *BtImportService) SaveKeywordData(tx *gorm.DB, marketID int64, keywords []KeywordInfo) error {
	for _, keywordData := range keywords {
		source := keywordData.Source
		keyword := keywordData.Keyword

		// 创建或更新关键词
		kw := brandtrekin.BtKeyword{
			MarketId: &marketID,
			Keyword:  &keyword,
			Source:   &source,
		}

		var existingKeyword brandtrekin.BtKeyword
		err := tx.Where("market_id = ? AND keyword = ? AND source = ?", marketID, keyword, source).
			First(&existingKeyword).Error

		var keywordID int64
		if err == gorm.ErrRecordNotFound {
			if err := tx.Create(&kw).Error; err != nil {
				return fmt.Errorf("创建关键词失败: %v", err)
			}
			keywordID = int64(kw.ID)
		} else if err != nil {
			return fmt.Errorf("查询关键词失败: %v", err)
		} else {
			keywordID = int64(existingKeyword.ID)
		}

		// 保存月度搜索量
		for _, mv := range keywordData.MonthlyVolumes {
			volume := mv.Volume
			date := mv.Date

			monthlyVol := brandtrekin.BtKeywordMonthlyVolume{
				KeywordId: &keywordID,
				MarketId:  &marketID,
				Keyword:   &keyword,
				Date:      &date,
				Volume:    &volume,
			}

			// 查找或创建月度记录
			var existing brandtrekin.BtKeywordMonthlyVolume
			err := tx.Where("keyword_id = ? AND date = ?", keywordID, date).
				First(&existing).Error

			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&monthlyVol).Error; err != nil {
					return fmt.Errorf("创建月度搜索量失败: %v", err)
				}
			} else if err != nil {
				return fmt.Errorf("查询月度搜索量失败: %v", err)
			} else {
				// 更新
				if err := tx.Model(&existing).Updates(&monthlyVol).Error; err != nil {
					return fmt.Errorf("更新月度搜索量失败: %v", err)
				}
			}
		}
	}

	return nil
}

// SaveProductSalesData 保存商品销售数据
func (s *BtImportService) SaveProductSalesData(tx *gorm.DB, productSalesData []ProductSalesInfo) error {
	for _, salesData := range productSalesData {
		asin := salesData.ASIN

		for _, ms := range salesData.MonthlySales {
			sales := ms.Sales
			units := ms.Units
			date := ms.Date

			monthlySales := brandtrekin.BtProductMonthlySales{
				Asin:  &asin,
				Date:  &date,
				Sales: &sales,
				Units: &units,
			}

			// 查找或创建月度销售记录
			var existing brandtrekin.BtProductMonthlySales
			err := tx.Where("asin = ? AND date = ?", asin, date).
				First(&existing).Error

			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&monthlySales).Error; err != nil {
					return fmt.Errorf("创建月度销售记录失败: %v", err)
				}
			} else if err != nil {
				return fmt.Errorf("查询月度销售记录失败: %v", err)
			} else {
				// 更新
				if err := tx.Model(&existing).Updates(&monthlySales).Error; err != nil {
					return fmt.Errorf("更新月度销售记录失败: %v", err)
				}
			}
		}
	}

	return nil
}

// BatchImport 批量导入所有数据（在事务中执行）
func (s *BtImportService) BatchImport(marketID int64, files map[string]*multipart.FileHeader, replaceMode bool) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 如果是全量替换模式，删除现有数据
		if replaceMode {
			if err := s.deleteMarketData(tx, marketID); err != nil {
				return fmt.Errorf("删除现有数据失败: %v", err)
			}
		}

		// 1. 导入品牌社交媒体数据
		if file, ok := files["brandSocial"]; ok {
			result, err := s.ParseBrandSocial(file)
			if err != nil {
				return fmt.Errorf("解析品牌社交媒体数据失败: %v", err)
			}

			if brands, ok := result.FullData.([]BrandWithSocial); ok {
				if err := s.SaveBrandSocialData(tx, marketID, brands); err != nil {
					return err
				}
			}
		}

		// 2. 导入商品数据
		if file, ok := files["productUS"]; ok {
			result, err := s.ParseProductUS(file)
			if err != nil {
				return fmt.Errorf("解析商品数据失败: %v", err)
			}

			if products, ok := result.FullData.([]ProductInfo); ok {
				if err := s.SaveProductData(tx, marketID, products); err != nil {
					return err
				}
			}
		}

		// 3. 导入Google关键词数据
		if file, ok := files["gkw"]; ok {
			result, err := s.ParseGKW(file)
			if err != nil {
				return fmt.Errorf("解析Google关键词数据失败: %v", err)
			}

			if keywords, ok := result.FullData.([]KeywordInfo); ok {
				if err := s.SaveKeywordData(tx, marketID, keywords); err != nil {
					return err
				}
			}
		}

		// 4. 导入Amazon关键词数据
		if file, ok := files["keywordHistory"]; ok {
			result, err := s.ParseKeywordHistory(file)
			if err != nil {
				return fmt.Errorf("解析Amazon关键词数据失败: %v", err)
			}

			if keywords, ok := result.FullData.([]KeywordInfo); ok {
				if err := s.SaveKeywordData(tx, marketID, keywords); err != nil {
					return err
				}
			}
		}

		// 5. 导入商品销售数据
		if file, ok := files["productSales"]; ok {
			result, err := s.ParseProductSales(file)
			if err != nil {
				return fmt.Errorf("解析商品销售数据失败: %v", err)
			}

			if salesData, ok := result.FullData.([]ProductSalesInfo); ok {
				if err := s.SaveProductSalesData(tx, salesData); err != nil {
					return err
				}
			}
		}

		// 6. 导入完成后，自动运行数据聚合计算
		// 注意：这里使用全局DB而不是事务tx，因为聚合服务自己管理事务
		aggregateService := BtAggregateService{}
		if err := aggregateService.RunFullAggregation(marketID); err != nil {
			return fmt.Errorf("数据聚合计算失败: %v", err)
		}

		return nil
	})
}

// deleteMarketData 删除市场的所有关联数据
func (s *BtImportService) deleteMarketData(tx *gorm.DB, marketID int64) error {
	// 按顺序删除关联数据（从子表到父表）

	// 1. 删除商品月度销售数据
	if err := tx.Where("asin IN (SELECT asin FROM bt_products WHERE market_id = ?)", marketID).
		Delete(&brandtrekin.BtProductMonthlySales{}).Error; err != nil {
		return err
	}

	// 2. 删除商品数据
	if err := tx.Where("market_id = ?", marketID).Delete(&brandtrekin.BtProduct{}).Error; err != nil {
		return err
	}

	// 3. 删除品牌月度趋势数据
	if err := tx.Where("brand_id IN (SELECT id FROM bt_brands WHERE market_id = ?)", marketID).
		Delete(&brandtrekin.BtBrandMonthlyTrend{}).Error; err != nil {
		return err
	}

	// 4. 删除品牌社交媒体数据
	if err := tx.Where("brand_id IN (SELECT id FROM bt_brands WHERE market_id = ?)", marketID).
		Delete(&brandtrekin.BtBrandSocialMedia{}).Error; err != nil {
		return err
	}

	// 5. 删除品牌数据
	if err := tx.Where("market_id = ?", marketID).Delete(&brandtrekin.BtBrand{}).Error; err != nil {
		return err
	}

	// 6. 删除关键词月度搜索量数据
	if err := tx.Where("market_id = ?", marketID).Delete(&brandtrekin.BtKeywordMonthlyVolume{}).Error; err != nil {
		return err
	}

	// 7. 删除关键词数据
	if err := tx.Where("market_id = ?", marketID).Delete(&brandtrekin.BtKeyword{}).Error; err != nil {
		return err
	}

	// 8. 删除市场月度趋势数据
	if err := tx.Where("market_id = ?", marketID).Delete(&brandtrekin.BtMarketMonthlyTrend{}).Error; err != nil {
		return err
	}

	return nil
}
