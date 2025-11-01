package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GKWRow 定义GKW CSV文件的行结构（使用map来处理动态列）
type GKWRow map[string]string

// MonthlyVolume 月度数据
type MonthlyVolume struct {
	Date   time.Time
	Volume int
}

// KeywordInfo 关键词信息
type KeywordInfo struct {
	Keyword        string
	Source         string // google 或 amazon
	MonthlyVolumes []MonthlyVolume
}

// RowError 行错误
type RowError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

// ImportResult 导入结果
type ImportResult struct {
	Success  bool        `json:"success"`
	Total    int         `json:"total"`
	Errors   []RowError  `json:"errors"`
	Preview  interface{} `json:"preview"` // 预览数据（前5条）
	FullData interface{} `json:"-"`       // 完整数据（不序列化到JSON）
}

// ParseGKW 解析Google关键词数据
func ParseGKW(file *multipart.FileHeader) (*ImportResult, error) {
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

		// 使用csv包解析CSV文件
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

			// 解析搜索量（支持带小数点的数字，如 "14800.0"）
			volumeFloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
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

// 使用反射设置 FileHeader 的 Open 方法
func setFileHeaderOpen(header *multipart.FileHeader, filePath string) {
	// 获取 FileHeader 的反射值
	v := reflect.ValueOf(header).Elem()
	
	// 查找 content 字段（私有字段）
	// multipart.FileHeader 有一个私有的 content 字段
	// 我们需要通过反射来访问它
	
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	
	// 使用 unsafe 包来设置私有字段
	contentField := v.FieldByName("content")
	if contentField.IsValid() {
		// 使用 unsafe 来修改私有字段
		reflect.NewAt(contentField.Type(), unsafe.Pointer(contentField.UnsafeAddr())).
			Elem().Set(reflect.ValueOf(content))
	}
}

func testFile(filePath string) {
	fmt.Printf("\n========== 测试文件: %s ==========\n", filePath)

	// 打开文件获取信息
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("❌ 无法打开文件: %v\n", err)
		return
	}

	info, err := file.Stat()
	if err != nil {
		fmt.Printf("❌ 无法获取文件信息: %v\n", err)
		file.Close()
		return
	}
	file.Close()

	// 创建 FileHeader
	header := &multipart.FileHeader{
		Filename: info.Name(),
		Size:     info.Size(),
	}

	// 设置文件内容
	setFileHeaderOpen(header, filePath)

	// 解析文件
	result, err := ParseGKW(header)
	if err != nil {
		fmt.Printf("❌ 解析文件失败: %v\n", err)
		return
	}

	// 验证结果
	if !result.Success {
		fmt.Printf("❌ 解析失败\n")
		return
	}

	if result.Total == 0 {
		fmt.Printf("❌ 未解析到任何关键词\n")
		return
	}

	fmt.Printf("✅ 成功解析 %d 个关键词\n", result.Total)

	// 打印错误信息
	if len(result.Errors) > 0 {
		fmt.Printf("\n⚠️  解析过程中的错误:\n")
		for _, err := range result.Errors {
			fmt.Printf("  行 %d: %s\n", err.Row, err.Error)
		}
	}

	// 打印预览数据
	if keywords, ok := result.Preview.([]KeywordInfo); ok {
		fmt.Printf("\n📊 预览数据:\n")
		for i, kw := range keywords {
			fmt.Printf("  %d. 关键词: %s\n", i+1, kw.Keyword)
			fmt.Printf("     来源: %s\n", kw.Source)
			fmt.Printf("     月度数据数量: %d\n", len(kw.MonthlyVolumes))
			if len(kw.MonthlyVolumes) > 0 {
				fmt.Printf("     第一个月度数据: %s - %d\n",
					kw.MonthlyVolumes[0].Date.Format("2006-01"),
					kw.MonthlyVolumes[0].Volume)
				if len(kw.MonthlyVolumes) > 1 {
					fmt.Printf("     最后一个月度数据: %s - %d\n",
						kw.MonthlyVolumes[len(kw.MonthlyVolumes)-1].Date.Format("2006-01"),
						kw.MonthlyVolumes[len(kw.MonthlyVolumes)-1].Volume)
				}
			}
			fmt.Println()
		}
	}
}

func main() {
	// 测试 LaserEngraver 文件
	testFile("/Users/leon/code/ai-code/gin-vue-admin/trekin-main/data/LaserEngraver/GKW.csv")

	// 测试 CNCRouter 文件
	testFile("/Users/leon/code/ai-code/gin-vue-admin/trekin-main/data/CNCRouter/GKW.csv")

	fmt.Println("\n========== 测试完成 ==========")
}
