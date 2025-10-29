
// 自动生成模板BtImportLog
package brandtrekin
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 数据导入日志 结构体  BtImportLog
type BtImportLog struct {
    global.GVA_MODEL
  MarketId  *int64 `json:"marketId" form:"marketId" gorm:"comment:市场ID;column:market_id;" binding:"required"`  //所属市场
  ImportMode  *string `json:"importMode" form:"importMode" gorm:"comment:导入模式:incremental/replace;column:import_mode;size:20;" binding:"required"`  //导入模式
  Status  *string `json:"status" form:"status" gorm:"comment:状态:success/failed/partial;column:status;size:20;" binding:"required"`  //导入状态
  BrandsCount  *int64 `json:"brandsCount" form:"brandsCount" gorm:"default:0;comment:导入品牌数;column:brands_count;"`  //导入品牌数
  ProductsCount  *int64 `json:"productsCount" form:"productsCount" gorm:"default:0;comment:导入商品数;column:products_count;"`  //导入商品数
  KeywordsCount  *int64 `json:"keywordsCount" form:"keywordsCount" gorm:"default:0;comment:导入关键词数;column:keywords_count;"`  //导入关键词数
  SalesRecordsCount  *int64 `json:"salesRecordsCount" form:"salesRecordsCount" gorm:"default:0;comment:导入销售记录数;column:sales_records_count;"`  //导入销售记录数
  SkippedCount  *int64 `json:"skippedCount" form:"skippedCount" gorm:"default:0;comment:跳过记录数;column:skipped_count;"`  //跳过记录数
  ErrorMessage  *string `json:"errorMessage" form:"errorMessage" gorm:"comment:错误信息;column:error_message;size:1000;type:text;"`  //错误信息
  LogContent  *string `json:"logContent" form:"logContent" gorm:"comment:日志内容;column:log_content;size:10000;type:text;"`  //日志内容
  CreatedBy  *string `json:"createdBy" form:"createdBy" gorm:"comment:操作人;column:created_by;size:50;"`  //操作人
}


// TableName 数据导入日志 BtImportLog自定义表名 bt_import_logs
func (BtImportLog) TableName() string {
    return "bt_import_logs"
}





