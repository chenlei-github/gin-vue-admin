package brandtrekin

import (
	"github.com/gin-gonic/gin"
)

type BtDisplayRouter struct{}

// InitBtDisplayRouter 初始化 前端展示 路由信息（公开接口，无需认证）
func (s *BtDisplayRouter) InitBtDisplayRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	// 使用PublicRouter，因为这些是公开的展示接口
	btDisplayRouter := PublicRouter.Group("api")

	{
		btDisplayRouter.GET("markets", btDisplayApi.GetMarketList)                            // 获取市场列表
		btDisplayRouter.GET("markets/:id", btDisplayApi.GetMarketDetail)                      // 获取市场详情
		btDisplayRouter.GET("markets/:marketId/brands/:brandName", btDisplayApi.GetBrandDetail) // 获取品牌详情
	}
}
