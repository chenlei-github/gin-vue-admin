package brandtrekin

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	BtMarketRouter
	BtBrandRouter
	BtBrandSocialMediaRouter
	BtProductRouter
	BtProductMonthlySalesRouter
	BtKeywordRouter
	BtKeywordMonthlyVolumeRouter
	BtBrandMonthlyTrendRouter
	BtMarketMonthlyTrendRouter
	BtImportLogRouter
	BtImportRouter
	BtDisplayRouter
}

var (
	btMarketApi               = api.ApiGroupApp.BrandtrekinApiGroup.BtMarketApi
	btBrandApi                = api.ApiGroupApp.BrandtrekinApiGroup.BtBrandApi
	btBrandSocialMediaApi     = api.ApiGroupApp.BrandtrekinApiGroup.BtBrandSocialMediaApi
	btProductApi              = api.ApiGroupApp.BrandtrekinApiGroup.BtProductApi
	btProductMonthlySalesApi  = api.ApiGroupApp.BrandtrekinApiGroup.BtProductMonthlySalesApi
	btKeywordApi              = api.ApiGroupApp.BrandtrekinApiGroup.BtKeywordApi
	btKeywordMonthlyVolumeApi = api.ApiGroupApp.BrandtrekinApiGroup.BtKeywordMonthlyVolumeApi
	btBrandMonthlyTrendApi    = api.ApiGroupApp.BrandtrekinApiGroup.BtBrandMonthlyTrendApi
	btMarketMonthlyTrendApi   = api.ApiGroupApp.BrandtrekinApiGroup.BtMarketMonthlyTrendApi
	btImportLogApi            = api.ApiGroupApp.BrandtrekinApiGroup.BtImportLogApi
	btImportApi               = api.ApiGroupApp.BrandtrekinApiGroup.BtImportApi
	btDisplayApi              = api.ApiGroupApp.BrandtrekinApiGroup.BtDisplayApi
)
