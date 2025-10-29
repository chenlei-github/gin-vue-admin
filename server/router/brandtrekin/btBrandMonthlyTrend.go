package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtBrandMonthlyTrendRouter struct {}

// InitBtBrandMonthlyTrendRouter 初始化 品牌月度趋势 路由信息
func (s *BtBrandMonthlyTrendRouter) InitBtBrandMonthlyTrendRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	btBrandMonthlyTrendRouter := Router.Group("btBrandMonthlyTrend").Use(middleware.OperationRecord())
	btBrandMonthlyTrendRouterWithoutRecord := Router.Group("btBrandMonthlyTrend")
	btBrandMonthlyTrendRouterWithoutAuth := PublicRouter.Group("btBrandMonthlyTrend")
	{
		btBrandMonthlyTrendRouter.POST("createBtBrandMonthlyTrend", btBrandMonthlyTrendApi.CreateBtBrandMonthlyTrend)   // 新建品牌月度趋势
		btBrandMonthlyTrendRouter.DELETE("deleteBtBrandMonthlyTrend", btBrandMonthlyTrendApi.DeleteBtBrandMonthlyTrend) // 删除品牌月度趋势
		btBrandMonthlyTrendRouter.DELETE("deleteBtBrandMonthlyTrendByIds", btBrandMonthlyTrendApi.DeleteBtBrandMonthlyTrendByIds) // 批量删除品牌月度趋势
		btBrandMonthlyTrendRouter.PUT("updateBtBrandMonthlyTrend", btBrandMonthlyTrendApi.UpdateBtBrandMonthlyTrend)    // 更新品牌月度趋势
	}
	{
		btBrandMonthlyTrendRouterWithoutRecord.GET("findBtBrandMonthlyTrend", btBrandMonthlyTrendApi.FindBtBrandMonthlyTrend)        // 根据ID获取品牌月度趋势
		btBrandMonthlyTrendRouterWithoutRecord.GET("getBtBrandMonthlyTrendList", btBrandMonthlyTrendApi.GetBtBrandMonthlyTrendList)  // 获取品牌月度趋势列表
	}
	{
	    btBrandMonthlyTrendRouterWithoutAuth.GET("getBtBrandMonthlyTrendDataSource", btBrandMonthlyTrendApi.GetBtBrandMonthlyTrendDataSource)  // 获取品牌月度趋势数据源
	    btBrandMonthlyTrendRouterWithoutAuth.GET("getBtBrandMonthlyTrendPublic", btBrandMonthlyTrendApi.GetBtBrandMonthlyTrendPublic)  // 品牌月度趋势开放接口
	}
}
