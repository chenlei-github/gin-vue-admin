
package brandtrekin

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
)

type BtKeywordService struct {}
// CreateBtKeyword 创建关键词管理记录
// Author [yourname](https://github.com/yourname)
func (btKeywordService *BtKeywordService) CreateBtKeyword(ctx context.Context, btKeyword *brandtrekin.BtKeyword) (err error) {
	err = global.GVA_DB.Create(btKeyword).Error
	return err
}

// DeleteBtKeyword 删除关键词管理记录
// Author [yourname](https://github.com/yourname)
func (btKeywordService *BtKeywordService)DeleteBtKeyword(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&brandtrekin.BtKeyword{},"id = ?",ID).Error
	return err
}

// DeleteBtKeywordByIds 批量删除关键词管理记录
// Author [yourname](https://github.com/yourname)
func (btKeywordService *BtKeywordService)DeleteBtKeywordByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]brandtrekin.BtKeyword{},"id in ?",IDs).Error
	return err
}

// UpdateBtKeyword 更新关键词管理记录
// Author [yourname](https://github.com/yourname)
func (btKeywordService *BtKeywordService)UpdateBtKeyword(ctx context.Context, btKeyword brandtrekin.BtKeyword) (err error) {
	err = global.GVA_DB.Model(&brandtrekin.BtKeyword{}).Where("id = ?",btKeyword.ID).Updates(&btKeyword).Error
	return err
}

// GetBtKeyword 根据ID获取关键词管理记录
// Author [yourname](https://github.com/yourname)
func (btKeywordService *BtKeywordService)GetBtKeyword(ctx context.Context, ID string) (btKeyword brandtrekin.BtKeyword, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&btKeyword).Error
	return
}
// GetBtKeywordInfoList 分页获取关键词管理记录
// Author [yourname](https://github.com/yourname)
func (btKeywordService *BtKeywordService)GetBtKeywordInfoList(ctx context.Context, info brandtrekinReq.BtKeywordSearch) (list []brandtrekin.BtKeyword, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&brandtrekin.BtKeyword{})
    var btKeywords []brandtrekin.BtKeyword
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.MarketId != nil {
        db = db.Where("market_id = ?", *info.MarketId)
    }
    if info.Keyword != nil && *info.Keyword != "" {
        db = db.Where("keyword LIKE ?", "%"+ *info.Keyword+"%")
    }
    if info.Source != nil && *info.Source != "" {
        db = db.Where("source = ?", *info.Source)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&btKeywords).Error
	return  btKeywords, total, err
}
func (btKeywordService *BtKeywordService)GetBtKeywordDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   marketId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("bt_markets").Where("deleted_at IS NULL").Select("market_name as label,id as value").Scan(&marketId)
	   res["marketId"] = marketId
	return
}
func (btKeywordService *BtKeywordService)GetBtKeywordPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
