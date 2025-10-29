
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type BtBrandSocialMediaSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      BrandId  *int `json:"brandId" form:"brandId"` 
      Platform  *string `json:"platform" form:"platform"` 
    request.PageInfo
}
