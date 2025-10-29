
// 自动生成模板BtMarket
package brandtrekin
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 市场管理 结构体  BtMarket
type BtMarket struct {
    global.GVA_MODEL
  MarketName  *string `json:"marketName" form:"marketName" gorm:"comment:市场名称;column:market_name;size:100;" binding:"required"`  //市场名称
  MarketSlug  *string `json:"marketSlug" form:"marketSlug" gorm:"unique;comment:市场slug;column:market_slug;size:100;" binding:"required"`  //市场标识符(用于URL)
  Description  *string `json:"description" form:"description" gorm:"comment:市场描述;column:description;size:500;type:text;"`  //市场描述
  Status  *string `json:"status" form:"status" gorm:"default:active;comment:状态:active/inactive;column:status;size:20;" binding:"required"`  //状态
  TotalRevenue  *float64 `json:"totalRevenue" form:"totalRevenue" gorm:"default:0;comment:总销售额;column:total_revenue;size:15,2;"`  //总销售额(最近12个月)
  TotalProducts  *int64 `json:"totalProducts" form:"totalProducts" gorm:"default:0;comment:商品总数;column:total_products;"`  //商品总数
  BrandCount  *int64 `json:"brandCount" form:"brandCount" gorm:"default:0;comment:品牌数量;column:brand_count;"`  //品牌数量
  SearchVolume  *int64 `json:"searchVolume" form:"searchVolume" gorm:"default:0;comment:搜索量;column:search_volume;"`  //搜索量(最近月份)
  Cagr  *float64 `json:"cagr" form:"cagr" gorm:"comment:CAGR年复合增长率;column:cagr;size:5,2;"`  //年复合增长率(%)
}


// TableName 市场管理 BtMarket自定义表名 bt_markets
func (BtMarket) TableName() string {
    return "bt_markets"
}





