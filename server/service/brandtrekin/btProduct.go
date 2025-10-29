
package brandtrekin

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
)

type BtProductService struct {}
// CreateBtProduct 创建商品管理记录
// Author [yourname](https://github.com/yourname)
func (btProductService *BtProductService) CreateBtProduct(ctx context.Context, btProduct *brandtrekin.BtProduct) (err error) {
	err = global.GVA_DB.Create(btProduct).Error
	return err
}

// DeleteBtProduct 删除商品管理记录
// Author [yourname](https://github.com/yourname)
func (btProductService *BtProductService)DeleteBtProduct(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&brandtrekin.BtProduct{},"id = ?",ID).Error
	return err
}

// DeleteBtProductByIds 批量删除商品管理记录
// Author [yourname](https://github.com/yourname)
func (btProductService *BtProductService)DeleteBtProductByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]brandtrekin.BtProduct{},"id in ?",IDs).Error
	return err
}

// UpdateBtProduct 更新商品管理记录
// Author [yourname](https://github.com/yourname)
func (btProductService *BtProductService)UpdateBtProduct(ctx context.Context, btProduct brandtrekin.BtProduct) (err error) {
	err = global.GVA_DB.Model(&brandtrekin.BtProduct{}).Where("id = ?",btProduct.ID).Updates(&btProduct).Error
	return err
}

// GetBtProduct 根据ID获取商品管理记录
// Author [yourname](https://github.com/yourname)
func (btProductService *BtProductService)GetBtProduct(ctx context.Context, ID string) (btProduct brandtrekin.BtProduct, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&btProduct).Error
	return
}
// GetBtProductInfoList 分页获取商品管理记录
// Author [yourname](https://github.com/yourname)
func (btProductService *BtProductService)GetBtProductInfoList(ctx context.Context, info brandtrekinReq.BtProductSearch) (list []brandtrekin.BtProduct, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&brandtrekin.BtProduct{})
    var btProducts []brandtrekin.BtProduct
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.MarketId != nil {
        db = db.Where("market_id = ?", *info.MarketId)
    }
    if info.BrandId != nil {
        db = db.Where("brand_id = ?", *info.BrandId)
    }
    if info.Asin != nil && *info.Asin != "" {
        db = db.Where("asin = ?", *info.Asin)
    }
    if info.Title != nil && *info.Title != "" {
        db = db.Where("title LIKE ?", "%"+ *info.Title+"%")
    }
	if info.StartPrice != nil && info.EndPrice != nil {
		db = db.Where("price BETWEEN ? AND ? ", *info.StartPrice, *info.EndPrice)
	}
	if info.StartRating != nil && info.EndRating != nil {
		db = db.Where("rating BETWEEN ? AND ? ", *info.StartRating, *info.EndRating)
	}
	if info.StartMonthlySales != nil && info.EndMonthlySales != nil {
		db = db.Where("monthly_sales BETWEEN ? AND ? ", *info.StartMonthlySales, *info.EndMonthlySales)
	}
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
           orderMap["id"] = true
           orderMap["created_at"] = true
         	orderMap["price"] = true
         	orderMap["rating"] = true
         	orderMap["reviews"] = true
         	orderMap["monthly_sales"] = true
       if orderMap[info.Sort] {
          OrderStr = info.Sort
          if info.Order == "descending" {
             OrderStr = OrderStr + " desc"
          }
          db = db.Order(OrderStr)
       }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&btProducts).Error
	return  btProducts, total, err
}
func (btProductService *BtProductService)GetBtProductDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   brandId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("bt_brands").Where("deleted_at IS NULL").Select("brand_name as label,id as value").Scan(&brandId)
	   res["brandId"] = brandId
	   marketId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("bt_markets").Where("deleted_at IS NULL").Select("market_name as label,id as value").Scan(&marketId)
	   res["marketId"] = marketId
	return
}
func (btProductService *BtProductService)GetBtProductPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
