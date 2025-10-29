
package brandtrekin

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
)

type BtProductMonthlySalesService struct {}
// CreateBtProductMonthlySales 创建商品月度销售记录
// Author [yourname](https://github.com/yourname)
func (btProductMonthlySalesService *BtProductMonthlySalesService) CreateBtProductMonthlySales(ctx context.Context, btProductMonthlySales *brandtrekin.BtProductMonthlySales) (err error) {
	err = global.GVA_DB.Create(btProductMonthlySales).Error
	return err
}

// DeleteBtProductMonthlySales 删除商品月度销售记录
// Author [yourname](https://github.com/yourname)
func (btProductMonthlySalesService *BtProductMonthlySalesService)DeleteBtProductMonthlySales(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&brandtrekin.BtProductMonthlySales{},"id = ?",ID).Error
	return err
}

// DeleteBtProductMonthlySalesByIds 批量删除商品月度销售记录
// Author [yourname](https://github.com/yourname)
func (btProductMonthlySalesService *BtProductMonthlySalesService)DeleteBtProductMonthlySalesByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]brandtrekin.BtProductMonthlySales{},"id in ?",IDs).Error
	return err
}

// UpdateBtProductMonthlySales 更新商品月度销售记录
// Author [yourname](https://github.com/yourname)
func (btProductMonthlySalesService *BtProductMonthlySalesService)UpdateBtProductMonthlySales(ctx context.Context, btProductMonthlySales brandtrekin.BtProductMonthlySales) (err error) {
	err = global.GVA_DB.Model(&brandtrekin.BtProductMonthlySales{}).Where("id = ?",btProductMonthlySales.ID).Updates(&btProductMonthlySales).Error
	return err
}

// GetBtProductMonthlySales 根据ID获取商品月度销售记录
// Author [yourname](https://github.com/yourname)
func (btProductMonthlySalesService *BtProductMonthlySalesService)GetBtProductMonthlySales(ctx context.Context, ID string) (btProductMonthlySales brandtrekin.BtProductMonthlySales, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&btProductMonthlySales).Error
	return
}
// GetBtProductMonthlySalesInfoList 分页获取商品月度销售记录
// Author [yourname](https://github.com/yourname)
func (btProductMonthlySalesService *BtProductMonthlySalesService)GetBtProductMonthlySalesInfoList(ctx context.Context, info brandtrekinReq.BtProductMonthlySalesSearch) (list []brandtrekin.BtProductMonthlySales, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&brandtrekin.BtProductMonthlySales{})
    var btProductMonthlySaless []brandtrekin.BtProductMonthlySales
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.Asin != nil && *info.Asin != "" {
        db = db.Where("asin = ?", *info.Asin)
    }
			if len(info.DateRange) == 2 {
				db = db.Where("date BETWEEN ? AND ? ", info.DateRange[0], info.DateRange[1])
			}
	if info.StartSales != nil && info.EndSales != nil {
		db = db.Where("sales BETWEEN ? AND ? ", *info.StartSales, *info.EndSales)
	}
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
           orderMap["id"] = true
           orderMap["created_at"] = true
         	orderMap["date"] = true
         	orderMap["sales"] = true
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

	err = db.Find(&btProductMonthlySaless).Error
	return  btProductMonthlySaless, total, err
}
func (btProductMonthlySalesService *BtProductMonthlySalesService)GetBtProductMonthlySalesPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
