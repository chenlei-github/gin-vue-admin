
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type BtMarketSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      MarketName  *string `json:"marketName" form:"marketName"` 
      MarketSlug  *string `json:"marketSlug" form:"marketSlug"` 
      Status  *string `json:"status" form:"status"` 
      StartTotalRevenue  *float64  `json:"startTotalRevenue" form:"startTotalRevenue"`
EndTotalRevenue  *float64  `json:"endTotalRevenue" form:"endTotalRevenue"`
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
