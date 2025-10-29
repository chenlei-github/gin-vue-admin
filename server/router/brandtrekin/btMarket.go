package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtMarketRouter struct {}

// InitBtMarketRouter 初始化 市场管理 路由信息
func (s *BtMarketRouter) InitBtMarketRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	btMarketRouter := Router.Group("btMarket").Use(middleware.OperationRecord())
	btMarketRouterWithoutRecord := Router.Group("btMarket")
	btMarketRouterWithoutAuth := PublicRouter.Group("btMarket")
	{
		btMarketRouter.POST("createBtMarket", btMarketApi.CreateBtMarket)   // 新建市场管理
		btMarketRouter.DELETE("deleteBtMarket", btMarketApi.DeleteBtMarket) // 删除市场管理
		btMarketRouter.DELETE("deleteBtMarketByIds", btMarketApi.DeleteBtMarketByIds) // 批量删除市场管理
		btMarketRouter.PUT("updateBtMarket", btMarketApi.UpdateBtMarket)    // 更新市场管理
	}
	{
		btMarketRouterWithoutRecord.GET("findBtMarket", btMarketApi.FindBtMarket)        // 根据ID获取市场管理
		btMarketRouterWithoutRecord.GET("getBtMarketList", btMarketApi.GetBtMarketList)  // 获取市场管理列表
	}
	{
	    btMarketRouterWithoutAuth.GET("getBtMarketPublic", btMarketApi.GetBtMarketPublic)  // 市场管理开放接口
	}
}
