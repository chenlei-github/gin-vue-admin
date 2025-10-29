
package brandtrekin

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
)

type BtBrandService struct {}
// CreateBtBrand 创建品牌管理记录
// Author [yourname](https://github.com/yourname)
func (btBrandService *BtBrandService) CreateBtBrand(ctx context.Context, btBrand *brandtrekin.BtBrand) (err error) {
	err = global.GVA_DB.Create(btBrand).Error
	return err
}

// DeleteBtBrand 删除品牌管理记录
// Author [yourname](https://github.com/yourname)
func (btBrandService *BtBrandService)DeleteBtBrand(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&brandtrekin.BtBrand{},"id = ?",ID).Error
	return err
}

// DeleteBtBrandByIds 批量删除品牌管理记录
// Author [yourname](https://github.com/yourname)
func (btBrandService *BtBrandService)DeleteBtBrandByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]brandtrekin.BtBrand{},"id in ?",IDs).Error
	return err
}

// UpdateBtBrand 更新品牌管理记录
// Author [yourname](https://github.com/yourname)
func (btBrandService *BtBrandService)UpdateBtBrand(ctx context.Context, btBrand brandtrekin.BtBrand) (err error) {
	err = global.GVA_DB.Model(&brandtrekin.BtBrand{}).Where("id = ?",btBrand.ID).Updates(&btBrand).Error
	return err
}

// GetBtBrand 根据ID获取品牌管理记录
// Author [yourname](https://github.com/yourname)
func (btBrandService *BtBrandService)GetBtBrand(ctx context.Context, ID string) (btBrand brandtrekin.BtBrand, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&btBrand).Error
	return
}
// GetBtBrandInfoList 分页获取品牌管理记录
// Author [yourname](https://github.com/yourname)
func (btBrandService *BtBrandService)GetBtBrandInfoList(ctx context.Context, info brandtrekinReq.BtBrandSearch) (list []brandtrekin.BtBrand, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&brandtrekin.BtBrand{})
    var btBrands []brandtrekin.BtBrand
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.MarketId != nil {
        db = db.Where("market_id = ?", *info.MarketId)
    }
    if info.BrandName != nil && *info.BrandName != "" {
        db = db.Where("brand_name LIKE ?", "%"+ *info.BrandName+"%")
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
         	orderMap["product_count"] = true
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

	err = db.Find(&btBrands).Error
	return  btBrands, total, err
}
func (btBrandService *BtBrandService)GetBtBrandDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   marketId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("bt_markets").Where("deleted_at IS NULL").Select("market_name as label,id as value").Scan(&marketId)
	   res["marketId"] = marketId
	return
}
func (btBrandService *BtBrandService)GetBtBrandPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
