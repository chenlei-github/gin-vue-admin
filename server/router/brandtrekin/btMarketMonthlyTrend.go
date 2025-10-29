package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtMarketMonthlyTrendRouter struct {}

// InitBtMarketMonthlyTrendRouter 初始化 市场月度趋势 路由信息
func (s *BtMarketMonthlyTrendRouter) InitBtMarketMonthlyTrendRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	btMarketMonthlyTrendRouter := Router.Group("btMarketMonthlyTrend").Use(middleware.OperationRecord())
	btMarketMonthlyTrendRouterWithoutRecord := Router.Group("btMarketMonthlyTrend")
	btMarketMonthlyTrendRouterWithoutAuth := PublicRouter.Group("btMarketMonthlyTrend")
	{
		btMarketMonthlyTrendRouter.POST("createBtMarketMonthlyTrend", btMarketMonthlyTrendApi.CreateBtMarketMonthlyTrend)   // 新建市场月度趋势
		btMarketMonthlyTrendRouter.DELETE("deleteBtMarketMonthlyTrend", btMarketMonthlyTrendApi.DeleteBtMarketMonthlyTrend) // 删除市场月度趋势
		btMarketMonthlyTrendRouter.DELETE("deleteBtMarketMonthlyTrendByIds", btMarketMonthlyTrendApi.DeleteBtMarketMonthlyTrendByIds) // 批量删除市场月度趋势
		btMarketMonthlyTrendRouter.PUT("updateBtMarketMonthlyTrend", btMarketMonthlyTrendApi.UpdateBtMarketMonthlyTrend)    // 更新市场月度趋势
	}
	{
		btMarketMonthlyTrendRouterWithoutRecord.GET("findBtMarketMonthlyTrend", btMarketMonthlyTrendApi.FindBtMarketMonthlyTrend)        // 根据ID获取市场月度趋势
		btMarketMonthlyTrendRouterWithoutRecord.GET("getBtMarketMonthlyTrendList", btMarketMonthlyTrendApi.GetBtMarketMonthlyTrendList)  // 获取市场月度趋势列表
	}
	{
	    btMarketMonthlyTrendRouterWithoutAuth.GET("getBtMarketMonthlyTrendDataSource", btMarketMonthlyTrendApi.GetBtMarketMonthlyTrendDataSource)  // 获取市场月度趋势数据源
	    btMarketMonthlyTrendRouterWithoutAuth.GET("getBtMarketMonthlyTrendPublic", btMarketMonthlyTrendApi.GetBtMarketMonthlyTrendPublic)  // 市场月度趋势开放接口
	}
}
