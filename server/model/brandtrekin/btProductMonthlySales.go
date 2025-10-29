
// 自动生成模板BtProductMonthlySales
package brandtrekin
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// 商品月度销售 结构体  BtProductMonthlySales
type BtProductMonthlySales struct {
    global.GVA_MODEL
  Asin  *string `json:"asin" form:"asin" gorm:"comment:商品ASIN;column:asin;size:20;" binding:"required"`  //商品ASIN
  Date  *time.Time `json:"date" form:"date" gorm:"comment:月份;column:date;" binding:"required"`  //月份(YYYY-MM-01)
  Sales  *float64 `json:"sales" form:"sales" gorm:"default:0;comment:销售额;column:sales;size:12,2;"`  //销售额
  Units  *int64 `json:"units" form:"units" gorm:"default:0;comment:销量;column:units;"`  //销量
}


// TableName 商品月度销售 BtProductMonthlySales自定义表名 bt_product_monthly_sales
func (BtProductMonthlySales) TableName() string {
    return "bt_product_monthly_sales"
}





