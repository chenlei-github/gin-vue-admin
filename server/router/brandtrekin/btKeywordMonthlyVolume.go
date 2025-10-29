package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtKeywordMonthlyVolumeRouter struct {}

// InitBtKeywordMonthlyVolumeRouter 初始化 关键词月度搜索量 路由信息
func (s *BtKeywordMonthlyVolumeRouter) InitBtKeywordMonthlyVolumeRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	btKeywordMonthlyVolumeRouter := Router.Group("btKeywordMonthlyVolume").Use(middleware.OperationRecord())
	btKeywordMonthlyVolumeRouterWithoutRecord := Router.Group("btKeywordMonthlyVolume")
	btKeywordMonthlyVolumeRouterWithoutAuth := PublicRouter.Group("btKeywordMonthlyVolume")
	{
		btKeywordMonthlyVolumeRouter.POST("createBtKeywordMonthlyVolume", btKeywordMonthlyVolumeApi.CreateBtKeywordMonthlyVolume)   // 新建关键词月度搜索量
		btKeywordMonthlyVolumeRouter.DELETE("deleteBtKeywordMonthlyVolume", btKeywordMonthlyVolumeApi.DeleteBtKeywordMonthlyVolume) // 删除关键词月度搜索量
		btKeywordMonthlyVolumeRouter.DELETE("deleteBtKeywordMonthlyVolumeByIds", btKeywordMonthlyVolumeApi.DeleteBtKeywordMonthlyVolumeByIds) // 批量删除关键词月度搜索量
		btKeywordMonthlyVolumeRouter.PUT("updateBtKeywordMonthlyVolume", btKeywordMonthlyVolumeApi.UpdateBtKeywordMonthlyVolume)    // 更新关键词月度搜索量
	}
	{
		btKeywordMonthlyVolumeRouterWithoutRecord.GET("findBtKeywordMonthlyVolume", btKeywordMonthlyVolumeApi.FindBtKeywordMonthlyVolume)        // 根据ID获取关键词月度搜索量
		btKeywordMonthlyVolumeRouterWithoutRecord.GET("getBtKeywordMonthlyVolumeList", btKeywordMonthlyVolumeApi.GetBtKeywordMonthlyVolumeList)  // 获取关键词月度搜索量列表
	}
	{
	    btKeywordMonthlyVolumeRouterWithoutAuth.GET("getBtKeywordMonthlyVolumeDataSource", btKeywordMonthlyVolumeApi.GetBtKeywordMonthlyVolumeDataSource)  // 获取关键词月度搜索量数据源
	    btKeywordMonthlyVolumeRouterWithoutAuth.GET("getBtKeywordMonthlyVolumePublic", btKeywordMonthlyVolumeApi.GetBtKeywordMonthlyVolumePublic)  // 关键词月度搜索量开放接口
	}
}
