package brandtrekin

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
)

type BtImportService struct{}

// ImportResult 导入结果
type ImportResult struct {
	Success  bool        `json:"success"`
	Total    int         `json:"total"`
	Errors   []RowError  `json:"errors"`
	Preview  interface{} `json:"preview"` // 预览数据（前5条）
	FullData interface{} `json:"-"`       // 完整数据（不序列化到JSON）
}

// RowError 行错误
type RowError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

// BrandWithSocial 品牌及其社交媒体信息
type BrandWithSocial struct {
	BrandName          string
	Website            string
	YoutubeUrl         string
	YoutubeSubscribers int
	InstagramUrl       string
	InstagramFollowers int
	FacebookUrl        string
	FacebookFollowers  int
	RedditUrl          string
	RedditPosts        int
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
	ExtendedData map[string]interface{} // 存储Excel中的额外字段
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

	// 识别Report sheet（跳过Method sheet）
	sheetNames := f.GetSheetList()
	var sheetName string
	for _, name := range sheetNames {
		if name == "Report" {
			sheetName = name
			break
		}
	}
	if sheetName == "" {
		return nil, errors.New("找不到Report sheet")
	}

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

		// 处理YouTube数据（支持"YouTube URL"和"YouTube"两种列名）
		if idx, ok := colMap["YouTube URL"]; ok && idx < len(row) {
			brand.YoutubeUrl = strings.TrimSpace(row[idx])
		} else if idx, ok := colMap["YouTube"]; ok && idx < len(row) {
			brand.YoutubeUrl = strings.TrimSpace(row[idx])
		}
		if idx, ok := colMap["YouTube Subscribers"]; ok && idx < len(row) && row[idx] != "" {
			if subs, valid := parseNumericValue(row[idx]); valid {
				brand.YoutubeSubscribers = int(subs)
			}
		}

		// 处理Instagram数据（支持"Instagram URL"和"Instagram"两种列名）
		if idx, ok := colMap["Instagram URL"]; ok && idx < len(row) {
			brand.InstagramUrl = strings.TrimSpace(row[idx])
		} else if idx, ok := colMap["Instagram"]; ok && idx < len(row) {
			brand.InstagramUrl = strings.TrimSpace(row[idx])
		}
		if idx, ok := colMap["Instagram Followers"]; ok && idx < len(row) && row[idx] != "" {
			if followers, valid := parseNumericValue(row[idx]); valid {
				brand.InstagramFollowers = int(followers)
			}
		}

		// 处理Facebook数据（支持"Facebook URL"和"Facebook"两种列名）
		if idx, ok := colMap["Facebook URL"]; ok && idx < len(row) {
			brand.FacebookUrl = strings.TrimSpace(row[idx])
		} else if idx, ok := colMap["Facebook"]; ok && idx < len(row) {
			brand.FacebookUrl = strings.TrimSpace(row[idx])
		}
		// 支持"Facebook Followers/Likes"和"Facebook Followers"两种列名
		if idx, ok := colMap["Facebook Followers/Likes"]; ok && idx < len(row) && row[idx] != "" {
			if followers, valid := parseNumericValue(row[idx]); valid {
				brand.FacebookFollowers = int(followers)
			}
		} else if idx, ok := colMap["Facebook Followers"]; ok && idx < len(row) && row[idx] != "" {
			if followers, valid := parseNumericValue(row[idx]); valid {
				brand.FacebookFollowers = int(followers)
			}
		}

		// 处理Reddit数据（支持"Reddit URL/Search"和"Reddit"两种列名）
		if idx, ok := colMap["Reddit URL/Search"]; ok && idx < len(row) {
			brand.RedditUrl = strings.TrimSpace(row[idx])
		} else if idx, ok := colMap["Reddit"]; ok && idx < len(row) {
			brand.RedditUrl = strings.TrimSpace(row[idx])
		}
		// 支持"Reddit Mentions (approx)"和"Reddit Posts"两种列名
		if idx, ok := colMap["Reddit Mentions (approx)"]; ok && idx < len(row) && row[idx] != "" {
			if posts, valid := parseNumericValue(row[idx]); valid {
				brand.RedditPosts = int(posts)
			}
		} else if idx, ok := colMap["Reddit Posts"]; ok && idx < len(row) && row[idx] != "" {
			if posts, valid := parseNumericValue(row[idx]); valid {
				brand.RedditPosts = int(posts)
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

// GKWRow 定义GKW CSV文件的行结构（使用map来处理动态列）
type GKWRow map[string]string

// ParseGKW 解析Google关键词数据（使用gocsv包）
func (s *BtImportService) ParseGKW(file *multipart.FileHeader) (*ImportResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 读取文件内容
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 尝试不同的编码：UTF-8 -> GBK -> GB18030
	encodings := []encoding.Encoding{
		nil, // UTF-8不需要转换
		simplifiedchinese.GBK,
		simplifiedchinese.GB18030,
	}

	var rows []GKWRow
	var parseErr error

	for _, enc := range encodings {
		var reader io.Reader
		if enc == nil {
			// UTF-8，直接读取
			reader = bytes.NewReader(fileBytes)
		} else {
			// 转换编码
			decoder := enc.NewDecoder()
			decoded, _, err := transform.Bytes(decoder, fileBytes)
			if err != nil {
				continue
			}
			reader = bytes.NewReader(decoded)
		}

		// 使用gocsv解析CSV文件
		// 先读取原始CSV数据
		csvReader := csv.NewReader(reader)
		csvReader.LazyQuotes = true
		csvReader.TrimLeadingSpace = true

		// 读取所有行
		records, err := csvReader.ReadAll()
		if err != nil {
			parseErr = err
			continue
		}

		if len(records) < 2 {
			parseErr = errors.New("文件为空或缺少数据行")
			continue
		}

		// 验证是否包含Keyword列
		header := records[0]
		hasKeyword := false
		for _, col := range header {
			if strings.TrimSpace(col) == "Keyword" {
				hasKeyword = true
				break
			}
		}

		if !hasKeyword {
			parseErr = errors.New("缺少必填列: Keyword")
			continue
		}

		// 手动将records转换为map结构
		rows = make([]GKWRow, 0, len(records)-1)
		for i := 1; i < len(records); i++ {
			row := make(GKWRow)
			for j, col := range header {
				if j < len(records[i]) {
					row[strings.TrimSpace(col)] = strings.TrimSpace(records[i][j])
				}
			}
			rows = append(rows, row)
		}

		// 解析成功，跳出循环
		parseErr = nil
		break
	}

	if parseErr != nil {
		return nil, fmt.Errorf("解析CSV文件失败，尝试了多种编码: %v", parseErr)
	}

	if len(rows) == 0 {
		return nil, errors.New("文件为空或缺少数据行")
	}

	// 验证数据格式：确保有Keyword列
	firstRow := rows[0]
	if _, ok := firstRow["Keyword"]; !ok {
		return nil, errors.New("缺少必填列: Keyword")
	}

	// 支持两种日期格式：
	// 1. "YYYY-MM" (如 "2021-09")
	// 2. "Searches: Mon YYYY" (如 "Searches: Sep 2021")
	datePattern1 := regexp.MustCompile(`^\d{4}-\d{2}$`)
	datePattern2 := regexp.MustCompile(`^Searches:\s+(\w+)\s+(\d{4})$`)

	// 月份名称映射
	monthMap := map[string]string{
		"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04",
		"May": "05", "Jun": "06", "Jul": "07", "Aug": "08",
		"Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12",
	}

	var keywords []KeywordInfo
	var rowErrors []RowError

	// 解析每一行数据
	for i, row := range rows {
		keyword := row["Keyword"]
		if keyword == "" {
			continue
		}

		var monthlyVolumes []MonthlyVolume

		// 遍历所有列，查找日期列
		for colName, value := range row {
			if colName == "Keyword" || colName == "Currency" || colName == "Segmentation" {
				continue
			}

			if value == "" {
				continue
			}

			// 尝试解析日期格式
			var date time.Time
			var dateStr string

			// 格式1: "YYYY-MM"
			if datePattern1.MatchString(colName) {
				dateStr = colName
				var err error
				date, err = time.Parse("2006-01", dateStr)
				if err != nil {
					continue
				}
			} else if matches := datePattern2.FindStringSubmatch(colName); matches != nil {
				// 格式2: "Searches: Mon YYYY"
				monthName := matches[1]
				year := matches[2]
				month, ok := monthMap[monthName]
				if !ok {
					continue
				}
				dateStr = fmt.Sprintf("%s-%s", year, month)
				var err error
				date, err = time.Parse("2006-01", dateStr)
				if err != nil {
					continue
				}
			} else {
				// 跳过非日期列
				continue
			}

			// 解析搜索量（支持带小数点的数字和千位分隔符，如 "14,800.0"）
			volumeFloat, valid := parseNumericValue(value)
			if !valid {
				// 如果解析失败，记录错误但继续处理
				rowErrors = append(rowErrors, RowError{
					Row:   i + 2, // +2 因为跳过了header行，且行号从1开始
					Error: fmt.Sprintf("列 %s 的值 '%s' 无法解析为数字", colName, value),
				})
				continue
			}

			volume := int(volumeFloat)
			if volume <= 0 {
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

	// 识别所有History-开头的sheet（排除Notes sheet）
	sheetNames := f.GetSheetList()
	var historySheets []string
	for _, name := range sheetNames {
		if strings.HasPrefix(name, "History-") && !strings.Contains(name, "Notes") {
			historySheets = append(historySheets, name)
		}
	}

	if len(historySheets) == 0 {
		return nil, errors.New("找不到History sheet")
	}

	var keywords []KeywordInfo
	var rowErrors []RowError
	keywordMap := make(map[string]*KeywordInfo) // 用于按关键词分组

	// 遍历所有历史sheet
	for _, sheetName := range historySheets {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			rowErrors = append(rowErrors, RowError{
				Row:   0,
				Error: fmt.Sprintf("读取sheet %s失败: %v", sheetName, err),
			})
			continue
		}

		if len(rows) < 3 {
			continue // 跳过空的sheet
		}

		// 跳过第1行标题行，第2行是表头（索引1）
		header := rows[1]
		colMap := make(map[string]int)
		for i, col := range header {
			if col != "" {
				colMap[col] = i
			}
		}

		// 校验必填列（中文列名）
		keywordCol, hasKeyword := colMap["关键词"]
		monthCol, hasMonth := colMap["月份"]
		volumeCol, hasVolume := colMap["月搜索量"]

		if !hasKeyword || !hasMonth || !hasVolume {
			rowErrors = append(rowErrors, RowError{
				Row:   0,
				Error: fmt.Sprintf("Sheet %s 缺少必填列（关键词、月份或月搜索量）", sheetName),
			})
			continue
		}

		// 从sheet名称提取关键词
		keywordFromSheet := strings.TrimPrefix(sheetName, "History-")
		keywordFromSheet = strings.TrimSuffix(keywordFromSheet, "-US")
		keywordFromSheet = strings.TrimSpace(keywordFromSheet)

		// 从第3行开始解析数据（索引2）
		for i := 2; i < len(rows); i++ {
			row := rows[i]

			if len(row) == 0 {
				continue
			}

			// 提取月份和搜索量
			if monthCol >= len(row) || volumeCol >= len(row) {
				continue
			}

			month := row[monthCol]
			volume := row[volumeCol]

			if month == "" || volume == "" {
				continue
			}

			// 解析月份格式（可能是 "2025-09" 或其他格式）
			var date time.Time
			monthStr := strings.TrimSpace(month)
			if matched, _ := regexp.MatchString(`^\d{4}-\d{2}$`, monthStr); matched {
				date, err = time.Parse("2006-01", monthStr)
				if err != nil {
					continue
				}
			} else {
				// 尝试其他日期格式
				continue
			}

			volumeFloat, valid := parseNumericValue(volume)
			if !valid || volumeFloat <= 0 {
				continue
			}
			volumeNum := int(volumeFloat)

			// 获取关键词（优先使用数据行中的关键词，否则使用sheet名称）
			keyword := keywordFromSheet
			if keywordCol < len(row) && row[keywordCol] != "" {
				keyword = strings.TrimSpace(row[keywordCol])
			}

			// 创建或更新关键词的月度数据
			key := keyword + "_amazon"
			if kw, exists := keywordMap[key]; exists {
				// 检查是否已存在该月份的数据
				exists := false
				for j, mv := range kw.MonthlyVolumes {
					if mv.Date.Year() == date.Year() && mv.Date.Month() == date.Month() {
						// 更新现有数据（如果新数据更大）
						if volumeNum > mv.Volume {
							kw.MonthlyVolumes[j].Volume = volumeNum
						}
						exists = true
						break
					}
				}
				if !exists {
					kw.MonthlyVolumes = append(kw.MonthlyVolumes, MonthlyVolume{
						Date:   date,
						Volume: volumeNum,
					})
				}
			} else {
				keywordMap[key] = &KeywordInfo{
					Keyword: keyword,
					Source:  "amazon",
					MonthlyVolumes: []MonthlyVolume{
						{
							Date:   date,
							Volume: volumeNum,
						},
					},
				}
			}
		}
	}

	// 转换为切片
	for _, kw := range keywordMap {
		if len(kw.MonthlyVolumes) > 0 {
			keywords = append(keywords, *kw)
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

// ParseProductUS 解析商品数据
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

	// 识别主数据sheet（排除Notes sheet）
	sheetNames := f.GetSheetList()
	var sheetName string
	for _, name := range sheetNames {
		if strings.HasPrefix(name, "US-") && !strings.Contains(name, "Notes") {
			sheetName = name
			break
		}
	}
	if sheetName == "" {
		return nil, errors.New("找不到主数据sheet")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}

	if len(rows) < 3 {
		return nil, errors.New("文件为空或缺少数据行（需要至少包含标题行、表头行和数据行）")
	}

	// 跳过第1行标题行，第2行是表头（索引1）
	header := rows[1]
	colMap := make(map[string]int)
	for i, col := range header {
		if col != "" {
			colMap[col] = i
		}
	}

	// 校验必填列（支持中英文列名）
	asinCol := -1
	brandCol := -1
	titleCol := -1

	// 查找ASIN列
	if idx, ok := colMap["ASIN"]; ok {
		asinCol = idx
	} else {
		for k, idx := range colMap {
			if strings.Contains(k, "ASIN") {
				asinCol = idx
				break
			}
		}
	}

	// 查找品牌列（支持"品牌"或"Brand"）
	if idx, ok := colMap["品牌"]; ok {
		brandCol = idx
	} else if idx, ok := colMap["Brand"]; ok {
		brandCol = idx
	} else {
		for k, idx := range colMap {
			if strings.Contains(k, "品牌") || strings.Contains(k, "Brand") {
				brandCol = idx
				break
			}
		}
	}

	// 查找商品标题列（支持"商品标题"或"Title"）
	if idx, ok := colMap["商品标题"]; ok {
		titleCol = idx
	} else if idx, ok := colMap["Title"]; ok {
		titleCol = idx
	} else {
		for k, idx := range colMap {
			if strings.Contains(k, "标题") || strings.Contains(k, "Title") {
				titleCol = idx
				break
			}
		}
	}

	if asinCol < 0 || brandCol < 0 || titleCol < 0 {
		return nil, errors.New("缺少必填列: ASIN、品牌/Brand或商品标题/Title")
	}

	var products []ProductInfo
	var rowErrors []RowError

	// 从第3行开始解析数据（索引2，跳过标题行和表头行）
	for i := 2; i < len(rows); i++ {
		row := rows[i]
		rowNumber := i + 1

		if len(row) == 0 {
			continue
		}

		if asinCol >= len(row) || brandCol >= len(row) || titleCol >= len(row) {
			rowErrors = append(rowErrors, RowError{
				Row:   rowNumber,
				Error: "缺少必填字段（ASIN、品牌或商品标题）",
			})
			continue
		}

		asin := strings.TrimSpace(row[asinCol])
		brand := strings.TrimSpace(row[brandCol])
		title := strings.TrimSpace(row[titleCol])

		if asin == "" || brand == "" || title == "" {
			rowErrors = append(rowErrors, RowError{
				Row:   rowNumber,
				Error: "ASIN、品牌或商品标题不能为空",
			})
			continue
		}

		product := ProductInfo{
			ASIN:  asin,
			Title: title,
			Brand: brand,
		}

	// 提取可选字段（支持中英文列名）
	// 价格（支持"价格"、"价格($)"、"Price"等列名）
	// 注意：需要移除千位分隔符（逗号）后再解析
	if idx, ok := colMap["价格($)"]; ok && idx < len(row) && row[idx] != "" {
		// 移除千位分隔符和货币符号
		priceStr := strings.TrimSpace(row[idx])
		priceStr = strings.ReplaceAll(priceStr, ",", "")
		priceStr = strings.TrimPrefix(priceStr, "$")
		priceStr = strings.TrimSpace(priceStr)
		
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			product.Price = price
		}
	} else if idx, ok := colMap["价格"]; ok && idx < len(row) && row[idx] != "" {
		priceStr := strings.TrimSpace(row[idx])
		priceStr = strings.ReplaceAll(priceStr, ",", "")
		priceStr = strings.TrimPrefix(priceStr, "$")
		priceStr = strings.TrimSpace(priceStr)
		
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			product.Price = price
		}
	} else if idx, ok := colMap["Price"]; ok && idx < len(row) && row[idx] != "" {
		priceStr := strings.TrimSpace(row[idx])
		priceStr = strings.ReplaceAll(priceStr, ",", "")
		priceStr = strings.TrimPrefix(priceStr, "$")
		priceStr = strings.TrimSpace(priceStr)
		
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			product.Price = price
		}
	}

		// 评分（支持"评分"、"Rating"等列名）
		if idx, ok := colMap["评分"]; ok && idx < len(row) && row[idx] != "" {
			if rating, valid := parseNumericValue(row[idx]); valid {
				product.Rating = rating
			}
		} else if idx, ok := colMap["Rating"]; ok && idx < len(row) && row[idx] != "" {
			if rating, valid := parseNumericValue(row[idx]); valid {
				product.Rating = rating
			}
		}

		// 评论数（支持"评分数"、"评论数"、"Reviews"等列名）
		if idx, ok := colMap["评分数"]; ok && idx < len(row) && row[idx] != "" {
			if reviews, valid := parseNumericValue(row[idx]); valid {
				product.Reviews = int(reviews)
			}
		} else if idx, ok := colMap["评论数"]; ok && idx < len(row) && row[idx] != "" {
			if reviews, valid := parseNumericValue(row[idx]); valid {
				product.Reviews = int(reviews)
			}
		} else if idx, ok := colMap["Reviews"]; ok && idx < len(row) && row[idx] != "" {
			if reviews, valid := parseNumericValue(row[idx]); valid {
				product.Reviews = int(reviews)
			}
		}

		// 图片URL
		if idx, ok := colMap["商品主图"]; ok && idx < len(row) {
			product.ImageUrl = strings.TrimSpace(row[idx])
		} else if idx, ok := colMap["Image"]; ok && idx < len(row) {
			product.ImageUrl = strings.TrimSpace(row[idx])
		}

		// 月销量
		if idx, ok := colMap["月销量"]; ok && idx < len(row) && row[idx] != "" {
			if sales, valid := parseNumericValue(row[idx]); valid {
				product.MonthlySales = int(sales)
			}
		} else if idx, ok := colMap["Monthly Sales"]; ok && idx < len(row) && row[idx] != "" {
			if sales, valid := parseNumericValue(row[idx]); valid {
				product.MonthlySales = int(sales)
			}
		}

		// 收集扩展数据（Excel中的其他字段）
		extendedData := make(map[string]interface{})
		
		// 定义已处理的核心字段，避免重复存储
		coreFields := map[string]bool{
			"ASIN": true, "品牌": true, "Brand": true, "商品标题": true, "Title": true,
			"价格($)": true, "价格": true, "Price": true, "评分": true, "Rating": true,
			"评分数": true, "评论数": true, "Reviews": true, "商品主图": true, "Image": true,
			"月销量": true, "Monthly Sales": true,
		}
		
		// 遍历所有列，收集未处理的字段
		for colName, colIdx := range colMap {
			if colName == "" || coreFields[colName] {
				continue // 跳过空列名和核心字段
			}
			
			if colIdx < len(row) {
				cellValue := strings.TrimSpace(row[colIdx])
				if cellValue != "" {
					// 尝试解析为数字，如果失败则保存为字符串
					if numValue, valid := parseNumericValue(cellValue); valid && numValue != 0 {
						extendedData[colName] = numValue
					} else {
						extendedData[colName] = cellValue
					}
				}
			}
		}
		
		product.ExtendedData = extendedData
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

// parseNumericValue 解析数值，处理各种格式：空格、逗号分隔符、货币符号等
// 返回 float64 值和是否解析成功
func parseNumericValue(value string) (float64, bool) {
	// 去除前后空格
	value = strings.TrimSpace(value)

	// 空值返回0
	if value == "" {
		return 0, true
	}

	// 去除常见的货币符号
	value = strings.TrimPrefix(value, "$")
	value = strings.TrimPrefix(value, "￥")
	value = strings.TrimPrefix(value, "€")
	value = strings.TrimPrefix(value, "£")
	value = strings.TrimSpace(value)

	// 去除千位分隔符（英文逗号）
	value = strings.ReplaceAll(value, ",", "")

	// 尝试解析为浮点数
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}

	// 确保非负数
	if result < 0 {
		return 0, false
	}

	return result, true
}

// processSheet 通用的sheet处理函数，用于处理销量或销售额sheet
// isVolumeSheet: true表示销量sheet，false表示销售额sheet
func (s *BtImportService) processSheet(f *excelize.File, sheetName string, isVolumeSheet bool,
	salesMap map[string]*ProductSalesInfo, datePattern *regexp.Regexp) []RowError {

	var rowErrors []RowError

	rows, err := f.GetRows(sheetName)
	if err != nil {
		rowErrors = append(rowErrors, RowError{
			Row:   0,
			Error: fmt.Sprintf("读取sheet %s失败: %v", sheetName, err),
		})
		return rowErrors
	}

	if len(rows) < 2 {
		return rowErrors
	}

	// 第一步：建立列索引映射（只遍历一次表头）
	header := rows[0]
	asinCol := -1
	dateColumns := make(map[int]time.Time) // 列索引 -> 日期

	for i, col := range header {
		col = strings.TrimSpace(col)

		if col == "ASIN" {
			asinCol = i
			continue
		}

		// 跳过非数据列
		if col == "图片" || col == "SKU" || col == "URL" ||
			col == "所属类目" || col == "商品标题" || col == "" {
			continue
		}

		// 提取日期（去除末尾的($)符号）
		dateStr := strings.TrimSuffix(col, "($)")
		dateStr = strings.TrimSpace(dateStr)

		// 匹配日期格式
		if datePattern.MatchString(dateStr) {
			if date, err := time.Parse("2006-01", dateStr); err == nil {
				dateColumns[i] = date
			}
		}
	}

	// 验证必填列
	if asinCol < 0 {
		rowErrors = append(rowErrors, RowError{
			Row:   0,
			Error: fmt.Sprintf("Sheet %s 缺少ASIN列", sheetName),
		})
		return rowErrors
	}

	// 第二步：遍历数据行（从第2行开始）
	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]

		// 跳过空行或ASIN列超出范围的行
		if len(row) == 0 || asinCol >= len(row) {
			continue
		}

		asin := strings.TrimSpace(row[asinCol])
		if asin == "" {
			continue
		}

		// 获取或创建该ASIN的销售数据记录
		salesInfo, exists := salesMap[asin]
		if !exists {
			salesInfo = &ProductSalesInfo{
				ASIN:         asin,
				MonthlySales: []MonthlySales{},
			}
			salesMap[asin] = salesInfo
		}

		// 第三步：遍历所有日期列，提取数据
		for colIdx, date := range dateColumns {
			// 获取单元格值
			var cellValue string
			if colIdx < len(row) {
				cellValue = row[colIdx]
			}

			// 解析数值
			numValue, valid := parseNumericValue(cellValue)
			if !valid {
				// 记录解析错误（可选）
				if cellValue != "" {
					rowErrors = append(rowErrors, RowError{
						Row: rowIdx + 1,
						Error: fmt.Sprintf("ASIN %s 的 %s 列值 '%s' 无法解析为数字",
							asin, date.Format("2006-01"), cellValue),
					})
				}
				continue
			}

			// 查找或创建该月份的销售记录
			found := false
			for j := range salesInfo.MonthlySales {
				ms := &salesInfo.MonthlySales[j]
				if ms.Date.Year() == date.Year() && ms.Date.Month() == date.Month() {
					// 更新对应字段
					if isVolumeSheet {
						ms.Units = int(numValue)
					} else {
						ms.Sales = numValue
					}
					found = true
					break
				}
			}

			// 如果不存在该月份的记录，创建新记录
			if !found {
				newRecord := MonthlySales{
					Date:  date,
					Sales: 0,
					Units: 0,
				}
				if isVolumeSheet {
					newRecord.Units = int(numValue)
				} else {
					newRecord.Sales = numValue
				}
				salesInfo.MonthlySales = append(salesInfo.MonthlySales, newRecord)
			}
		}
	}

	return rowErrors
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

	// 识别两个sheet：产品历史月销量和产品历史月销售额
	sheetNames := f.GetSheetList()
	var volumeSheet, revenueSheet string
	for _, name := range sheetNames {
		if name == "产品历史月销量" {
			volumeSheet = name
		} else if name == "产品历史月销售额" {
			revenueSheet = name
		}
	}

	if volumeSheet == "" && revenueSheet == "" {
		return nil, errors.New("找不到销售数据sheet（产品历史月销量或产品历史月销售额）")
	}

	// 编译日期正则表达式（只编译一次）
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}`)

	// 使用map按ASIN分组，避免重复查找
	salesMap := make(map[string]*ProductSalesInfo)
	var rowErrors []RowError

	// 处理产品历史月销量sheet
	if volumeSheet != "" {
		errors := s.processSheet(f, volumeSheet, true, salesMap, datePattern)
		rowErrors = append(rowErrors, errors...)
	}

	// 处理产品历史月销售额sheet
	if revenueSheet != "" {
		errors := s.processSheet(f, revenueSheet, false, salesMap, datePattern)
		rowErrors = append(rowErrors, errors...)
	}

	// 转换为切片并排序月度数据
	var productSalesData []ProductSalesInfo
	for _, salesInfo := range salesMap {
		if len(salesInfo.MonthlySales) > 0 {
			// 按日期排序月度数据
			sortMonthlySales(salesInfo.MonthlySales)
			productSalesData = append(productSalesData, *salesInfo)
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

// sortMonthlySales 按日期排序月度销售数据
func sortMonthlySales(sales []MonthlySales) {
	// 简单的冒泡排序（数据量不大时效率足够）
	for i := 0; i < len(sales)-1; i++ {
		for j := i + 1; j < len(sales); j++ {
			if sales[i].Date.After(sales[j].Date) {
				sales[i], sales[j] = sales[j], sales[i]
			}
		}
	}
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

		// 序列化扩展数据为JSON
		var extendedDataJSON string
		if len(productData.ExtendedData) > 0 {
			if jsonBytes, err := json.Marshal(productData.ExtendedData); err == nil {
				extendedDataJSON = string(jsonBytes)
			} else {
				// 如果JSON序列化失败，记录错误但继续处理
				fmt.Printf("警告: ASIN %s 的扩展数据JSON序列化失败: %v\n", productData.ASIN, err)
				extendedDataJSON = "{}"
			}
		} else {
			extendedDataJSON = "{}"
		}

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
			ExtendedData: extendedDataJSON,
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

			// 转换volume类型
			volumeInt64 := int64(volume)
			monthlyVol := brandtrekin.BtKeywordMonthlyVolume{
				KeywordId: &keywordID,
				Date:      &date,
				Volume:    &volumeInt64,
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

			// 转换units类型
			unitsInt64 := int64(units)
			monthlySales := brandtrekin.BtProductMonthlySales{
				Asin:  &asin,
				Date:  &date,
				Sales: &sales,
				Units: &unitsInt64,
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

// deleteMarketData 删除市场的所有关联数据
// 注意：使用 Unscoped() 进行物理删除，避免软删除导致的唯一索引冲突
func (s *BtImportService) deleteMarketData(tx *gorm.DB, marketID int64) error {
	// 按顺序删除关联数据（从子表到父表）

	// 1. 删除商品月度销售数据（物理删除）
	if err := tx.Unscoped().Where("asin IN (SELECT asin FROM bt_products WHERE market_id = ?)", marketID).
		Delete(&brandtrekin.BtProductMonthlySales{}).Error; err != nil {
		return err
	}

	// 2. 删除商品数据（物理删除，避免 ASIN 唯一索引冲突）
	if err := tx.Unscoped().Where("market_id = ?", marketID).Delete(&brandtrekin.BtProduct{}).Error; err != nil {
		return err
	}

	// 3. 删除品牌月度趋势数据（物理删除）
	if err := tx.Unscoped().Where("brand_id IN (SELECT id FROM bt_brands WHERE market_id = ? AND deleted_at IS NULL)", marketID).
		Delete(&brandtrekin.BtBrandMonthlyTrend{}).Error; err != nil {
		return err
	}

	// 4. 删除品牌社交媒体数据（物理删除）
	if err := tx.Unscoped().Where("brand_id IN (SELECT id FROM bt_brands WHERE market_id = ? AND deleted_at IS NULL)", marketID).
		Delete(&brandtrekin.BtBrandSocialMedia{}).Error; err != nil {
		return err
	}

	// 5. 删除品牌数据（物理删除，避免品牌名唯一索引冲突）
	if err := tx.Unscoped().Where("market_id = ?", marketID).Delete(&brandtrekin.BtBrand{}).Error; err != nil {
		return err
	}

	// 6. 删除关键词月度搜索量数据（物理删除）
	if err := tx.Unscoped().Where("keyword_id IN (SELECT id FROM bt_keywords WHERE market_id = ? AND deleted_at IS NULL)", marketID).
		Delete(&brandtrekin.BtKeywordMonthlyVolume{}).Error; err != nil {
		return err
	}

	// 7. 删除关键词数据（物理删除，避免关键词唯一索引冲突）
	if err := tx.Unscoped().Where("market_id = ?", marketID).Delete(&brandtrekin.BtKeyword{}).Error; err != nil {
		return err
	}

	// 8. 删除市场月度趋势数据（物理删除）
	if err := tx.Unscoped().Where("market_id = ?", marketID).Delete(&brandtrekin.BtMarketMonthlyTrend{}).Error; err != nil {
		return err
	}

	return nil
}
