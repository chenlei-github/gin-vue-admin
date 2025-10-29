
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type BtKeywordSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      MarketId  *int `json:"marketId" form:"marketId"` 
      Keyword  *string `json:"keyword" form:"keyword"` 
      Source  *string `json:"source" form:"source"` 
    request.PageInfo
}
