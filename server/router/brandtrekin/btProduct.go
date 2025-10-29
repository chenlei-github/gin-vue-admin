package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtProductRouter struct {}

// InitBtProductRouter 初始化 商品管理 路由信息
func (s *BtProductRouter) InitBtProductRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	btProductRouter := Router.Group("btProduct").Use(middleware.OperationRecord())
	btProductRouterWithoutRecord := Router.Group("btProduct")
	btProductRouterWithoutAuth := PublicRouter.Group("btProduct")
	{
		btProductRouter.POST("createBtProduct", btProductApi.CreateBtProduct)   // 新建商品管理
		btProductRouter.DELETE("deleteBtProduct", btProductApi.DeleteBtProduct) // 删除商品管理
		btProductRouter.DELETE("deleteBtProductByIds", btProductApi.DeleteBtProductByIds) // 批量删除商品管理
		btProductRouter.PUT("updateBtProduct", btProductApi.UpdateBtProduct)    // 更新商品管理
	}
	{
		btProductRouterWithoutRecord.GET("findBtProduct", btProductApi.FindBtProduct)        // 根据ID获取商品管理
		btProductRouterWithoutRecord.GET("getBtProductList", btProductApi.GetBtProductList)  // 获取商品管理列表
	}
	{
	    btProductRouterWithoutAuth.GET("getBtProductDataSource", btProductApi.GetBtProductDataSource)  // 获取商品管理数据源
	    btProductRouterWithoutAuth.GET("getBtProductPublic", btProductApi.GetBtProductPublic)  // 商品管理开放接口
	}
}
