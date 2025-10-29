
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type BtBrandSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      MarketId  *int `json:"marketId" form:"marketId"` 
      BrandName  *string `json:"brandName" form:"brandName"` 
      StartTotalRevenue  *float64  `json:"startTotalRevenue" form:"startTotalRevenue"`
EndTotalRevenue  *float64  `json:"endTotalRevenue" form:"endTotalRevenue"`
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
