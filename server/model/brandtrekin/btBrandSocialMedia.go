
// 自动生成模板BtBrandSocialMedia
package brandtrekin
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 品牌社交媒体 结构体  BtBrandSocialMedia
type BtBrandSocialMedia struct {
    global.GVA_MODEL
  BrandId  *int64 `json:"brandId" form:"brandId" gorm:"comment:品牌ID;column:brand_id;" binding:"required"`  //所属品牌
  Platform  *string `json:"platform" form:"platform" gorm:"comment:平台:youtube/instagram/facebook/reddit;column:platform;size:20;" binding:"required"`  //社交平台
  Url  *string `json:"url" form:"url" gorm:"comment:链接;column:url;size:500;" binding:"required"`  //平台链接
  Subscribers  *int64 `json:"subscribers" form:"subscribers" gorm:"default:0;comment:订阅数;column:subscribers;"`  //订阅数(YouTube)
  Followers  *int64 `json:"followers" form:"followers" gorm:"default:0;comment:粉丝数;column:followers;"`  //粉丝数(Instagram/Facebook)
  Posts  *int64 `json:"posts" form:"posts" gorm:"default:0;comment:帖子数;column:posts;"`  //帖子数(Reddit)
}


// TableName 品牌社交媒体 BtBrandSocialMedia自定义表名 bt_brand_social_media
func (BtBrandSocialMedia) TableName() string {
    return "bt_brand_social_media"
}





