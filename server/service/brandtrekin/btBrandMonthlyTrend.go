
package brandtrekin

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
)

type BtBrandMonthlyTrendService struct {}
// CreateBtBrandMonthlyTrend 创建品牌月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btBrandMonthlyTrendService *BtBrandMonthlyTrendService) CreateBtBrandMonthlyTrend(ctx context.Context, btBrandMonthlyTrend *brandtrekin.BtBrandMonthlyTrend) (err error) {
	err = global.GVA_DB.Create(btBrandMonthlyTrend).Error
	return err
}

// DeleteBtBrandMonthlyTrend 删除品牌月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btBrandMonthlyTrendService *BtBrandMonthlyTrendService)DeleteBtBrandMonthlyTrend(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&brandtrekin.BtBrandMonthlyTrend{},"id = ?",ID).Error
	return err
}

// DeleteBtBrandMonthlyTrendByIds 批量删除品牌月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btBrandMonthlyTrendService *BtBrandMonthlyTrendService)DeleteBtBrandMonthlyTrendByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]brandtrekin.BtBrandMonthlyTrend{},"id in ?",IDs).Error
	return err
}

// UpdateBtBrandMonthlyTrend 更新品牌月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btBrandMonthlyTrendService *BtBrandMonthlyTrendService)UpdateBtBrandMonthlyTrend(ctx context.Context, btBrandMonthlyTrend brandtrekin.BtBrandMonthlyTrend) (err error) {
	err = global.GVA_DB.Model(&brandtrekin.BtBrandMonthlyTrend{}).Where("id = ?",btBrandMonthlyTrend.ID).Updates(&btBrandMonthlyTrend).Error
	return err
}

// GetBtBrandMonthlyTrend 根据ID获取品牌月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btBrandMonthlyTrendService *BtBrandMonthlyTrendService)GetBtBrandMonthlyTrend(ctx context.Context, ID string) (btBrandMonthlyTrend brandtrekin.BtBrandMonthlyTrend, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&btBrandMonthlyTrend).Error
	return
}
// GetBtBrandMonthlyTrendInfoList 分页获取品牌月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btBrandMonthlyTrendService *BtBrandMonthlyTrendService)GetBtBrandMonthlyTrendInfoList(ctx context.Context, info brandtrekinReq.BtBrandMonthlyTrendSearch) (list []brandtrekin.BtBrandMonthlyTrend, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&brandtrekin.BtBrandMonthlyTrend{})
    var btBrandMonthlyTrends []brandtrekin.BtBrandMonthlyTrend
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.BrandId != nil {
        db = db.Where("brand_id = ?", *info.BrandId)
    }
			if len(info.DateRange) == 2 {
				db = db.Where("date BETWEEN ? AND ? ", info.DateRange[0], info.DateRange[1])
			}
	if info.StartRevenue != nil && info.EndRevenue != nil {
		db = db.Where("revenue BETWEEN ? AND ? ", *info.StartRevenue, *info.EndRevenue)
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
         	orderMap["revenue"] = true
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

	err = db.Find(&btBrandMonthlyTrends).Error
	return  btBrandMonthlyTrends, total, err
}
func (btBrandMonthlyTrendService *BtBrandMonthlyTrendService)GetBtBrandMonthlyTrendDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   brandId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("bt_brands").Where("deleted_at IS NULL").Select("brand_name as label,id as value").Scan(&brandId)
	   res["brandId"] = brandId
	return
}
func (btBrandMonthlyTrendService *BtBrandMonthlyTrendService)GetBtBrandMonthlyTrendPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
