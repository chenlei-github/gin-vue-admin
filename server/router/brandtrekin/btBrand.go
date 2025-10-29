package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtBrandRouter struct {}

// InitBtBrandRouter 初始化 品牌管理 路由信息
func (s *BtBrandRouter) InitBtBrandRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	btBrandRouter := Router.Group("btBrand").Use(middleware.OperationRecord())
	btBrandRouterWithoutRecord := Router.Group("btBrand")
	btBrandRouterWithoutAuth := PublicRouter.Group("btBrand")
	{
		btBrandRouter.POST("createBtBrand", btBrandApi.CreateBtBrand)   // 新建品牌管理
		btBrandRouter.DELETE("deleteBtBrand", btBrandApi.DeleteBtBrand) // 删除品牌管理
		btBrandRouter.DELETE("deleteBtBrandByIds", btBrandApi.DeleteBtBrandByIds) // 批量删除品牌管理
		btBrandRouter.PUT("updateBtBrand", btBrandApi.UpdateBtBrand)    // 更新品牌管理
	}
	{
		btBrandRouterWithoutRecord.GET("findBtBrand", btBrandApi.FindBtBrand)        // 根据ID获取品牌管理
		btBrandRouterWithoutRecord.GET("getBtBrandList", btBrandApi.GetBtBrandList)  // 获取品牌管理列表
	}
	{
	    btBrandRouterWithoutAuth.GET("getBtBrandDataSource", btBrandApi.GetBtBrandDataSource)  // 获取品牌管理数据源
	    btBrandRouterWithoutAuth.GET("getBtBrandPublic", btBrandApi.GetBtBrandPublic)  // 品牌管理开放接口
	}
}
