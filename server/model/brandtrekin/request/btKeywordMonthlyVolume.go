
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type BtKeywordMonthlyVolumeSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      KeywordId  *int `json:"keywordId" form:"keywordId"` 
      DateRange  []time.Time  `json:"dateRange" form:"dateRange[]"`
      StartVolume  *int  `json:"startVolume" form:"startVolume"`
EndVolume  *int  `json:"endVolume" form:"endVolume"`
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
