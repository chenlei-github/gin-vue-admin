
// 自动生成模板BtKeywordMonthlyVolume
package brandtrekin
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// 关键词月度搜索量 结构体  BtKeywordMonthlyVolume
type BtKeywordMonthlyVolume struct {
    global.GVA_MODEL
  KeywordId  *int64 `json:"keywordId" form:"keywordId" gorm:"comment:关键词ID;column:keyword_id;" binding:"required"`  //所属关键词
  Date  *time.Time `json:"date" form:"date" gorm:"comment:月份;column:date;" binding:"required"`  //月份(YYYY-MM-01)
  Volume  *int64 `json:"volume" form:"volume" gorm:"default:0;comment:搜索量;column:volume;"`  //搜索量
}


// TableName 关键词月度搜索量 BtKeywordMonthlyVolume自定义表名 bt_keyword_monthly_volume
func (BtKeywordMonthlyVolume) TableName() string {
    return "bt_keyword_monthly_volume"
}





