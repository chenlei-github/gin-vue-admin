
// 自动生成模板BtMarketMonthlyTrend
package brandtrekin
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// 市场月度趋势 结构体  BtMarketMonthlyTrend
type BtMarketMonthlyTrend struct {
    global.GVA_MODEL
  MarketId  *int64 `json:"marketId" form:"marketId" gorm:"comment:市场ID;column:market_id;" binding:"required"`  //所属市场
  Date  *time.Time `json:"date" form:"date" gorm:"comment:月份;column:date;" binding:"required"`  //月份(YYYY-MM-01)
  Revenue  *float64 `json:"revenue" form:"revenue" gorm:"default:0;comment:销售额;column:revenue;size:15,2;"`  //销售额
  SearchVolume  *int64 `json:"searchVolume" form:"searchVolume" gorm:"default:0;comment:搜索量;column:search_volume;"`  //搜索量
}


// TableName 市场月度趋势 BtMarketMonthlyTrend自定义表名 bt_market_monthly_trends
func (BtMarketMonthlyTrend) TableName() string {
    return "bt_market_monthly_trends"
}





