package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BtProductMonthlySales struct {
	Asin  *string
	Date  *string
	Sales *float64
	Units *int64
}

func main() {
	// 连接数据库
	dsn := "brandtrekin:cl@2025@!@tcp(rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com:3306)/brandtrekin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 查询问题记录总数
	var problemCount int64
	db.Table("bt_product_monthly_sales").
		Where("units > 0 AND sales = 0").
		Count(&problemCount)
	fmt.Printf("问题记录总数: %d\n\n", problemCount)

	// 查询总记录数和分类统计
	type Stats struct {
		Total       int64
		Problem     int64
		Normal      int64
		OnlyUnits   int64
		OnlySales   int64
		BothZero    int64
	}
	var stats Stats
	db.Table("bt_product_monthly_sales").Count(&stats.Total)
	db.Table("bt_product_monthly_sales").Where("units > 0 AND sales = 0").Count(&stats.Problem)
	db.Table("bt_product_monthly_sales").Where("units > 0 AND sales > 0").Count(&stats.Normal)
	db.Table("bt_product_monthly_sales").Where("units > 0 AND sales IS NULL").Count(&stats.OnlyUnits)
	db.Table("bt_product_monthly_sales").Where("(units IS NULL OR units = 0) AND sales > 0").Count(&stats.OnlySales)
	db.Table("bt_product_monthly_sales").Where("(units IS NULL OR units = 0) AND (sales IS NULL OR sales = 0)").Count(&stats.BothZero)

	fmt.Println("数据统计:")
	fmt.Printf("  总记录数: %d\n", stats.Total)
	fmt.Printf("  正常记录 (units>0 且 sales>0): %d\n", stats.Normal)
	fmt.Printf("  问题记录 (units>0 但 sales=0): %d\n", stats.Problem)
	fmt.Printf("  只有销量 (units>0 但 sales=NULL): %d\n", stats.OnlyUnits)
	fmt.Printf("  只有销售额 (units=0 但 sales>0): %d\n", stats.OnlySales)
	fmt.Printf("  两者都为空: %d\n\n", stats.BothZero)

	// 查询前10条问题记录的详细信息
	fmt.Println("前10条问题记录示例:")
	var records []BtProductMonthlySales
	db.Table("bt_product_monthly_sales").
		Select("asin, DATE_FORMAT(date, '%Y-%m') as date, sales, units").
		Where("units > 0 AND sales = 0").
		Limit(10).
		Scan(&records)

	for i, record := range records {
		asin := "NULL"
		if record.Asin != nil {
			asin = *record.Asin
		}
		date := "NULL"
		if record.Date != nil {
			date = *record.Date
		}
		sales := 0.0
		if record.Sales != nil {
			sales = *record.Sales
		}
		units := int64(0)
		if record.Units != nil {
			units = *record.Units
		}
		fmt.Printf("%d. ASIN: %s, Date: %s, Units: %d, Sales: %.2f\n", i+1, asin, date, units, sales)
	}

	// 按ASIN分组统计问题记录
	fmt.Println("\n按ASIN统计问题记录数量 (前10个):")
	type AsinStats struct {
		Asin   string
		Count  int64
	}
	var asinStats []AsinStats
	db.Table("bt_product_monthly_sales").
		Select("asin, COUNT(*) as count").
		Where("units > 0 AND sales = 0").
		Group("asin").
		Order("count DESC").
		Limit(10).
		Scan(&asinStats)

	for i, stat := range asinStats {
		fmt.Printf("%d. ASIN: %s, 问题记录数: %d\n", i+1, stat.Asin, stat.Count)
	}
}

