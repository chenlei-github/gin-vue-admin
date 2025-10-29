
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type BtImportLogSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      MarketId  *int `json:"marketId" form:"marketId"` 
      ImportMode  *string `json:"importMode" form:"importMode"` 
      Status  *string `json:"status" form:"status"` 
    request.PageInfo
}
