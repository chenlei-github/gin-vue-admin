package brandtrekin

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	BtMarketApi
	BtBrandApi
	BtBrandSocialMediaApi
	BtProductApi
	BtProductMonthlySalesApi
	BtKeywordApi
	BtKeywordMonthlyVolumeApi
	BtBrandMonthlyTrendApi
	BtMarketMonthlyTrendApi
	BtImportLogApi
}

var (
	btMarketService               = service.ServiceGroupApp.BrandtrekinServiceGroup.BtMarketService
	btBrandService                = service.ServiceGroupApp.BrandtrekinServiceGroup.BtBrandService
	btBrandSocialMediaService     = service.ServiceGroupApp.BrandtrekinServiceGroup.BtBrandSocialMediaService
	btProductService              = service.ServiceGroupApp.BrandtrekinServiceGroup.BtProductService
	btProductMonthlySalesService  = service.ServiceGroupApp.BrandtrekinServiceGroup.BtProductMonthlySalesService
	btKeywordService              = service.ServiceGroupApp.BrandtrekinServiceGroup.BtKeywordService
	btKeywordMonthlyVolumeService = service.ServiceGroupApp.BrandtrekinServiceGroup.BtKeywordMonthlyVolumeService
	btBrandMonthlyTrendService    = service.ServiceGroupApp.BrandtrekinServiceGroup.BtBrandMonthlyTrendService
	btMarketMonthlyTrendService   = service.ServiceGroupApp.BrandtrekinServiceGroup.BtMarketMonthlyTrendService
	btImportLogService            = service.ServiceGroupApp.BrandtrekinServiceGroup.BtImportLogService
)
