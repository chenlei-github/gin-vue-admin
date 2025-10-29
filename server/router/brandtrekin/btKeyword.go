package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtKeywordRouter struct {}

// InitBtKeywordRouter 初始化 关键词管理 路由信息
func (s *BtKeywordRouter) InitBtKeywordRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	btKeywordRouter := Router.Group("btKeyword").Use(middleware.OperationRecord())
	btKeywordRouterWithoutRecord := Router.Group("btKeyword")
	btKeywordRouterWithoutAuth := PublicRouter.Group("btKeyword")
	{
		btKeywordRouter.POST("createBtKeyword", btKeywordApi.CreateBtKeyword)   // 新建关键词管理
		btKeywordRouter.DELETE("deleteBtKeyword", btKeywordApi.DeleteBtKeyword) // 删除关键词管理
		btKeywordRouter.DELETE("deleteBtKeywordByIds", btKeywordApi.DeleteBtKeywordByIds) // 批量删除关键词管理
		btKeywordRouter.PUT("updateBtKeyword", btKeywordApi.UpdateBtKeyword)    // 更新关键词管理
	}
	{
		btKeywordRouterWithoutRecord.GET("findBtKeyword", btKeywordApi.FindBtKeyword)        // 根据ID获取关键词管理
		btKeywordRouterWithoutRecord.GET("getBtKeywordList", btKeywordApi.GetBtKeywordList)  // 获取关键词管理列表
	}
	{
	    btKeywordRouterWithoutAuth.GET("getBtKeywordDataSource", btKeywordApi.GetBtKeywordDataSource)  // 获取关键词管理数据源
	    btKeywordRouterWithoutAuth.GET("getBtKeywordPublic", btKeywordApi.GetBtKeywordPublic)  // 关键词管理开放接口
	}
}
