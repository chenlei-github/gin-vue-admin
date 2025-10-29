
// 自动生成模板BtProduct
package brandtrekin
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 商品管理 结构体  BtProduct
type BtProduct struct {
    global.GVA_MODEL
  MarketId  *int64 `json:"marketId" form:"marketId" gorm:"comment:市场ID;column:market_id;" binding:"required"`  //所属市场
  BrandId  *int64 `json:"brandId" form:"brandId" gorm:"comment:品牌ID;column:brand_id;" binding:"required"`  //所属品牌
  Asin  *string `json:"asin" form:"asin" gorm:"unique;comment:亚马逊ASIN;column:asin;size:20;" binding:"required"`  //亚马逊ASIN
  Title  *string `json:"title" form:"title" gorm:"comment:商品标题;column:title;size:500;" binding:"required"`  //商品标题
  Price  *float64 `json:"price" form:"price" gorm:"default:0;comment:价格;column:price;size:10,2;"`  //价格
  Rating  *float64 `json:"rating" form:"rating" gorm:"default:0;comment:评分;column:rating;size:3,2;"`  //评分(0-5)
  Reviews  *int64 `json:"reviews" form:"reviews" gorm:"default:0;comment:评论数;column:reviews;"`  //评论数
  MonthlySales  *int64 `json:"monthlySales" form:"monthlySales" gorm:"default:0;comment:月销量;column:monthly_sales;"`  //月销量
  ImageUrl  string `json:"imageUrl" form:"imageUrl" gorm:"comment:图片URL;column:image_url;size:500;"`  //商品图片
}


// TableName 商品管理 BtProduct自定义表名 bt_products
func (BtProduct) TableName() string {
    return "bt_products"
}





