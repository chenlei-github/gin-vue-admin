
// 自动生成模板BtBrandMonthlyTrend
package brandtrekin
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// 品牌月度趋势 结构体  BtBrandMonthlyTrend
type BtBrandMonthlyTrend struct {
    global.GVA_MODEL
  BrandId  *int64 `json:"brandId" form:"brandId" gorm:"comment:品牌ID;column:brand_id;" binding:"required"`  //所属品牌
  Date  *time.Time `json:"date" form:"date" gorm:"comment:月份;column:date;" binding:"required"`  //月份(YYYY-MM-01)
  Revenue  *float64 `json:"revenue" form:"revenue" gorm:"default:0;comment:销售额;column:revenue;size:12,2;"`  //销售额
}


// TableName 品牌月度趋势 BtBrandMonthlyTrend自定义表名 bt_brand_monthly_trends
func (BtBrandMonthlyTrend) TableName() string {
    return "bt_brand_monthly_trends"
}





