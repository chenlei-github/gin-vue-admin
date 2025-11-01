package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BtImportRouter struct{}

// InitBtImportRouter 初始化 数据导入 路由信息
func (s *BtImportRouter) InitBtImportRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	btImportRouter := Router.Group("btImport").Use(middleware.OperationRecord())

	{
		// 文件预览接口（用于上传前预览文件内容）
		btImportRouter.POST("previewBrandSocial", btImportApi.PreviewBrandSocial)           // 预览品牌社交媒体文件
		btImportRouter.POST("previewGKW", btImportApi.PreviewGKW)                           // 预览Google关键词文件
		btImportRouter.POST("previewKeywordHistory", btImportApi.PreviewKeywordHistory)     // 预览Amazon关键词历史文件
		btImportRouter.POST("previewProductUS", btImportApi.PreviewProductUS)               // 预览商品文件
		btImportRouter.POST("previewProductSales", btImportApi.PreviewProductSales)         // 预览商品销售文件

		// 批量导入接口
		btImportRouter.POST("batchImport", btImportApi.BatchImport)                         // 批量导入所有数据
	}
}
