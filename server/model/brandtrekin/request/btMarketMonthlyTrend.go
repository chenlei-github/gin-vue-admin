
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type BtMarketMonthlyTrendSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      MarketId  *int `json:"marketId" form:"marketId"` 
      DateRange  []time.Time  `json:"dateRange" form:"dateRange[]"`
      StartRevenue  *float64  `json:"startRevenue" form:"startRevenue"`
EndRevenue  *float64  `json:"endRevenue" form:"endRevenue"`
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
