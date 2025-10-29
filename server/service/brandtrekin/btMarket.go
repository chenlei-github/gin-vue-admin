
package brandtrekin

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
)

type BtMarketService struct {}
// CreateBtMarket 创建市场管理记录
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService) CreateBtMarket(ctx context.Context, btMarket *brandtrekin.BtMarket) (err error) {
	err = global.GVA_DB.Create(btMarket).Error
	return err
}

// DeleteBtMarket 删除市场管理记录
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService)DeleteBtMarket(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&brandtrekin.BtMarket{},"id = ?",ID).Error
	return err
}

// DeleteBtMarketByIds 批量删除市场管理记录
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService)DeleteBtMarketByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]brandtrekin.BtMarket{},"id in ?",IDs).Error
	return err
}

// UpdateBtMarket 更新市场管理记录
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService)UpdateBtMarket(ctx context.Context, btMarket brandtrekin.BtMarket) (err error) {
	err = global.GVA_DB.Model(&brandtrekin.BtMarket{}).Where("id = ?",btMarket.ID).Updates(&btMarket).Error
	return err
}

// GetBtMarket 根据ID获取市场管理记录
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService)GetBtMarket(ctx context.Context, ID string) (btMarket brandtrekin.BtMarket, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&btMarket).Error
	return
}
// GetBtMarketInfoList 分页获取市场管理记录
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService)GetBtMarketInfoList(ctx context.Context, info brandtrekinReq.BtMarketSearch) (list []brandtrekin.BtMarket, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&brandtrekin.BtMarket{})
    var btMarkets []brandtrekin.BtMarket
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.MarketName != nil && *info.MarketName != "" {
        db = db.Where("market_name LIKE ?", "%"+ *info.MarketName+"%")
    }
    if info.MarketSlug != nil && *info.MarketSlug != "" {
        db = db.Where("market_slug = ?", *info.MarketSlug)
    }
    if info.Status != nil && *info.Status != "" {
        db = db.Where("status = ?", *info.Status)
    }
	if info.StartTotalRevenue != nil && info.EndTotalRevenue != nil {
		db = db.Where("total_revenue BETWEEN ? AND ? ", *info.StartTotalRevenue, *info.EndTotalRevenue)
	}
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
           orderMap["id"] = true
           orderMap["created_at"] = true
         	orderMap["total_revenue"] = true
         	orderMap["total_products"] = true
         	orderMap["brand_count"] = true
         	orderMap["search_volume"] = true
         	orderMap["cagr"] = true
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

	err = db.Find(&btMarkets).Error
	return  btMarkets, total, err
}
func (btMarketService *BtMarketService)GetBtMarketPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
