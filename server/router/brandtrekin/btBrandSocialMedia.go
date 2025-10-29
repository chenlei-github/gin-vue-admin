package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtBrandSocialMediaRouter struct {}

// InitBtBrandSocialMediaRouter 初始化 品牌社交媒体 路由信息
func (s *BtBrandSocialMediaRouter) InitBtBrandSocialMediaRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	btBrandSocialMediaRouter := Router.Group("btBrandSocialMedia").Use(middleware.OperationRecord())
	btBrandSocialMediaRouterWithoutRecord := Router.Group("btBrandSocialMedia")
	btBrandSocialMediaRouterWithoutAuth := PublicRouter.Group("btBrandSocialMedia")
	{
		btBrandSocialMediaRouter.POST("createBtBrandSocialMedia", btBrandSocialMediaApi.CreateBtBrandSocialMedia)   // 新建品牌社交媒体
		btBrandSocialMediaRouter.DELETE("deleteBtBrandSocialMedia", btBrandSocialMediaApi.DeleteBtBrandSocialMedia) // 删除品牌社交媒体
		btBrandSocialMediaRouter.DELETE("deleteBtBrandSocialMediaByIds", btBrandSocialMediaApi.DeleteBtBrandSocialMediaByIds) // 批量删除品牌社交媒体
		btBrandSocialMediaRouter.PUT("updateBtBrandSocialMedia", btBrandSocialMediaApi.UpdateBtBrandSocialMedia)    // 更新品牌社交媒体
	}
	{
		btBrandSocialMediaRouterWithoutRecord.GET("findBtBrandSocialMedia", btBrandSocialMediaApi.FindBtBrandSocialMedia)        // 根据ID获取品牌社交媒体
		btBrandSocialMediaRouterWithoutRecord.GET("getBtBrandSocialMediaList", btBrandSocialMediaApi.GetBtBrandSocialMediaList)  // 获取品牌社交媒体列表
	}
	{
	    btBrandSocialMediaRouterWithoutAuth.GET("getBtBrandSocialMediaDataSource", btBrandSocialMediaApi.GetBtBrandSocialMediaDataSource)  // 获取品牌社交媒体数据源
	    btBrandSocialMediaRouterWithoutAuth.GET("getBtBrandSocialMediaPublic", btBrandSocialMediaApi.GetBtBrandSocialMediaPublic)  // 品牌社交媒体开放接口
	}
}
