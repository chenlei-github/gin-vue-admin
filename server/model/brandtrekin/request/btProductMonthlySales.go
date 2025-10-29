
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type BtProductMonthlySalesSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      Asin  *string `json:"asin" form:"asin"` 
      DateRange  []time.Time  `json:"dateRange" form:"dateRange[]"`
      StartSales  *float64  `json:"startSales" form:"startSales"`
EndSales  *float64  `json:"endSales" form:"endSales"`
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
