
// 自动生成模板BtBrand
package brandtrekin
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 品牌管理 结构体  BtBrand
type BtBrand struct {
    global.GVA_MODEL
  MarketId  *int64 `json:"marketId" form:"marketId" gorm:"comment:市场ID;column:market_id;" binding:"required"`  //所属市场
  BrandName  *string `json:"brandName" form:"brandName" gorm:"comment:品牌名称;column:brand_name;size:100;" binding:"required"`  //品牌名称
  Website  *string `json:"website" form:"website" gorm:"comment:品牌独立站;column:website;size:500;"`  //品牌官网
  TotalRevenue  *float64 `json:"totalRevenue" form:"totalRevenue" gorm:"default:0;comment:总销售额;column:total_revenue;size:15,2;"`  //总销售额(最近12个月)
  ProductCount  *int64 `json:"productCount" form:"productCount" gorm:"default:0;comment:商品数量;column:product_count;"`  //商品数量
  Cagr  *float64 `json:"cagr" form:"cagr" gorm:"comment:CAGR;column:cagr;size:5,2;"`  //年复合增长率(%)
}


// TableName 品牌管理 BtBrand自定义表名 bt_brands
func (BtBrand) TableName() string {
    return "bt_brands"
}





