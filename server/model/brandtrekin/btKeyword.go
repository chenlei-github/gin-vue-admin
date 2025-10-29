
// 自动生成模板BtKeyword
package brandtrekin
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 关键词管理 结构体  BtKeyword
type BtKeyword struct {
    global.GVA_MODEL
  MarketId  *int64 `json:"marketId" form:"marketId" gorm:"comment:市场ID;column:market_id;" binding:"required"`  //所属市场
  Keyword  *string `json:"keyword" form:"keyword" gorm:"comment:关键词;column:keyword;size:200;" binding:"required"`  //关键词
  Source  *string `json:"source" form:"source" gorm:"comment:来源:google/amazon;column:source;size:20;" binding:"required"`  //来源
}


// TableName 关键词管理 BtKeyword自定义表名 bt_keywords
func (BtKeyword) TableName() string {
    return "bt_keywords"
}





