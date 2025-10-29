
package brandtrekin

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
)

type BtMarketMonthlyTrendService struct {}
// CreateBtMarketMonthlyTrend 创建市场月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btMarketMonthlyTrendService *BtMarketMonthlyTrendService) CreateBtMarketMonthlyTrend(ctx context.Context, btMarketMonthlyTrend *brandtrekin.BtMarketMonthlyTrend) (err error) {
	err = global.GVA_DB.Create(btMarketMonthlyTrend).Error
	return err
}

// DeleteBtMarketMonthlyTrend 删除市场月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btMarketMonthlyTrendService *BtMarketMonthlyTrendService)DeleteBtMarketMonthlyTrend(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&brandtrekin.BtMarketMonthlyTrend{},"id = ?",ID).Error
	return err
}

// DeleteBtMarketMonthlyTrendByIds 批量删除市场月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btMarketMonthlyTrendService *BtMarketMonthlyTrendService)DeleteBtMarketMonthlyTrendByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]brandtrekin.BtMarketMonthlyTrend{},"id in ?",IDs).Error
	return err
}

// UpdateBtMarketMonthlyTrend 更新市场月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btMarketMonthlyTrendService *BtMarketMonthlyTrendService)UpdateBtMarketMonthlyTrend(ctx context.Context, btMarketMonthlyTrend brandtrekin.BtMarketMonthlyTrend) (err error) {
	err = global.GVA_DB.Model(&brandtrekin.BtMarketMonthlyTrend{}).Where("id = ?",btMarketMonthlyTrend.ID).Updates(&btMarketMonthlyTrend).Error
	return err
}

// GetBtMarketMonthlyTrend 根据ID获取市场月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btMarketMonthlyTrendService *BtMarketMonthlyTrendService)GetBtMarketMonthlyTrend(ctx context.Context, ID string) (btMarketMonthlyTrend brandtrekin.BtMarketMonthlyTrend, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&btMarketMonthlyTrend).Error
	return
}
// GetBtMarketMonthlyTrendInfoList 分页获取市场月度趋势记录
// Author [yourname](https://github.com/yourname)
func (btMarketMonthlyTrendService *BtMarketMonthlyTrendService)GetBtMarketMonthlyTrendInfoList(ctx context.Context, info brandtrekinReq.BtMarketMonthlyTrendSearch) (list []brandtrekin.BtMarketMonthlyTrend, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&brandtrekin.BtMarketMonthlyTrend{})
    var btMarketMonthlyTrends []brandtrekin.BtMarketMonthlyTrend
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.MarketId != nil {
        db = db.Where("market_id = ?", *info.MarketId)
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

	err = db.Find(&btMarketMonthlyTrends).Error
	return  btMarketMonthlyTrends, total, err
}
func (btMarketMonthlyTrendService *BtMarketMonthlyTrendService)GetBtMarketMonthlyTrendDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   marketId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("bt_markets").Where("deleted_at IS NULL").Select("market_name as label,id as value").Scan(&marketId)
	   res["marketId"] = marketId
	return
}
func (btMarketMonthlyTrendService *BtMarketMonthlyTrendService)GetBtMarketMonthlyTrendPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
