package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(brandtrekin.BtMarket{}, brandtrekin.BtBrand{}, brandtrekin.BtBrandSocialMedia{}, brandtrekin.BtProduct{}, brandtrekin.BtProductMonthlySales{}, brandtrekin.BtKeyword{}, brandtrekin.BtKeywordMonthlyVolume{}, brandtrekin.BtBrandMonthlyTrend{}, brandtrekin.BtMarketMonthlyTrend{}, brandtrekin.BtImportLog{})
	if err != nil {
		return err
	}
	return nil
}
