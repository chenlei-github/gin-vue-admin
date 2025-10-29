
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type BtProductSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      MarketId  *int `json:"marketId" form:"marketId"` 
      BrandId  *int `json:"brandId" form:"brandId"` 
      Asin  *string `json:"asin" form:"asin"` 
      Title  *string `json:"title" form:"title"` 
      StartPrice  *float64  `json:"startPrice" form:"startPrice"`
EndPrice  *float64  `json:"endPrice" form:"endPrice"`
      StartRating  *float64  `json:"startRating" form:"startRating"`
EndRating  *float64  `json:"endRating" form:"endRating"`
      StartMonthlySales  *int  `json:"startMonthlySales" form:"startMonthlySales"`
EndMonthlySales  *int  `json:"endMonthlySales" form:"endMonthlySales"`
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
