package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtImportLogRouter struct {}

// InitBtImportLogRouter 初始化 数据导入日志 路由信息
func (s *BtImportLogRouter) InitBtImportLogRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	btImportLogRouter := Router.Group("btImportLog").Use(middleware.OperationRecord())
	btImportLogRouterWithoutRecord := Router.Group("btImportLog")
	btImportLogRouterWithoutAuth := PublicRouter.Group("btImportLog")
	{
		btImportLogRouter.POST("createBtImportLog", btImportLogApi.CreateBtImportLog)   // 新建数据导入日志
		btImportLogRouter.DELETE("deleteBtImportLog", btImportLogApi.DeleteBtImportLog) // 删除数据导入日志
		btImportLogRouter.DELETE("deleteBtImportLogByIds", btImportLogApi.DeleteBtImportLogByIds) // 批量删除数据导入日志
		btImportLogRouter.PUT("updateBtImportLog", btImportLogApi.UpdateBtImportLog)    // 更新数据导入日志
	}
	{
		btImportLogRouterWithoutRecord.GET("findBtImportLog", btImportLogApi.FindBtImportLog)        // 根据ID获取数据导入日志
		btImportLogRouterWithoutRecord.GET("getBtImportLogList", btImportLogApi.GetBtImportLogList)  // 获取数据导入日志列表
	}
	{
	    btImportLogRouterWithoutAuth.GET("getBtImportLogDataSource", btImportLogApi.GetBtImportLogDataSource)  // 获取数据导入日志数据源
	    btImportLogRouterWithoutAuth.GET("getBtImportLogPublic", btImportLogApi.GetBtImportLogPublic)  // 数据导入日志开放接口
	}
}
