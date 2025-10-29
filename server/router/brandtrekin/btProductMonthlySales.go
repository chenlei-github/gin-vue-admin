package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtProductMonthlySalesRouter struct {}

// InitBtProductMonthlySalesRouter 初始化 商品月度销售 路由信息
func (s *BtProductMonthlySalesRouter) InitBtProductMonthlySalesRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	btProductMonthlySalesRouter := Router.Group("btProductMonthlySales").Use(middleware.OperationRecord())
	btProductMonthlySalesRouterWithoutRecord := Router.Group("btProductMonthlySales")
	btProductMonthlySalesRouterWithoutAuth := PublicRouter.Group("btProductMonthlySales")
	{
		btProductMonthlySalesRouter.POST("createBtProductMonthlySales", btProductMonthlySalesApi.CreateBtProductMonthlySales)   // 新建商品月度销售
		btProductMonthlySalesRouter.DELETE("deleteBtProductMonthlySales", btProductMonthlySalesApi.DeleteBtProductMonthlySales) // 删除商品月度销售
		btProductMonthlySalesRouter.DELETE("deleteBtProductMonthlySalesByIds", btProductMonthlySalesApi.DeleteBtProductMonthlySalesByIds) // 批量删除商品月度销售
		btProductMonthlySalesRouter.PUT("updateBtProductMonthlySales", btProductMonthlySalesApi.UpdateBtProductMonthlySales)    // 更新商品月度销售
	}
	{
		btProductMonthlySalesRouterWithoutRecord.GET("findBtProductMonthlySales", btProductMonthlySalesApi.FindBtProductMonthlySales)        // 根据ID获取商品月度销售
		btProductMonthlySalesRouterWithoutRecord.GET("getBtProductMonthlySalesList", btProductMonthlySalesApi.GetBtProductMonthlySalesList)  // 获取商品月度销售列表
	}
	{
	    btProductMonthlySalesRouterWithoutAuth.GET("getBtProductMonthlySalesPublic", btProductMonthlySalesApi.GetBtProductMonthlySalesPublic)  // 商品月度销售开放接口
	}
}
