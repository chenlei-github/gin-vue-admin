
package brandtrekin

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
)

type BtImportLogService struct {}
// CreateBtImportLog 创建数据导入日志记录
// Author [yourname](https://github.com/yourname)
func (btImportLogService *BtImportLogService) CreateBtImportLog(ctx context.Context, btImportLog *brandtrekin.BtImportLog) (err error) {
	err = global.GVA_DB.Create(btImportLog).Error
	return err
}

// DeleteBtImportLog 删除数据导入日志记录
// Author [yourname](https://github.com/yourname)
func (btImportLogService *BtImportLogService)DeleteBtImportLog(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&brandtrekin.BtImportLog{},"id = ?",ID).Error
	return err
}

// DeleteBtImportLogByIds 批量删除数据导入日志记录
// Author [yourname](https://github.com/yourname)
func (btImportLogService *BtImportLogService)DeleteBtImportLogByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]brandtrekin.BtImportLog{},"id in ?",IDs).Error
	return err
}

// UpdateBtImportLog 更新数据导入日志记录
// Author [yourname](https://github.com/yourname)
func (btImportLogService *BtImportLogService)UpdateBtImportLog(ctx context.Context, btImportLog brandtrekin.BtImportLog) (err error) {
	err = global.GVA_DB.Model(&brandtrekin.BtImportLog{}).Where("id = ?",btImportLog.ID).Updates(&btImportLog).Error
	return err
}

// GetBtImportLog 根据ID获取数据导入日志记录
// Author [yourname](https://github.com/yourname)
func (btImportLogService *BtImportLogService)GetBtImportLog(ctx context.Context, ID string) (btImportLog brandtrekin.BtImportLog, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&btImportLog).Error
	return
}
// GetBtImportLogInfoList 分页获取数据导入日志记录
// Author [yourname](https://github.com/yourname)
func (btImportLogService *BtImportLogService)GetBtImportLogInfoList(ctx context.Context, info brandtrekinReq.BtImportLogSearch) (list []brandtrekin.BtImportLog, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&brandtrekin.BtImportLog{})
    var btImportLogs []brandtrekin.BtImportLog
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.MarketId != nil {
        db = db.Where("market_id = ?", *info.MarketId)
    }
    if info.ImportMode != nil && *info.ImportMode != "" {
        db = db.Where("import_mode = ?", *info.ImportMode)
    }
    if info.Status != nil && *info.Status != "" {
        db = db.Where("status = ?", *info.Status)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&btImportLogs).Error
	return  btImportLogs, total, err
}
func (btImportLogService *BtImportLogService)GetBtImportLogDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   marketId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("bt_markets").Where("deleted_at IS NULL").Select("market_name as label,id as value").Scan(&marketId)
	   res["marketId"] = marketId
	return
}
func (btImportLogService *BtImportLogService)GetBtImportLogPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
