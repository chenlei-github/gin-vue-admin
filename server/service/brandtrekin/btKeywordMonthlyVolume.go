
package brandtrekin

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
)

type BtKeywordMonthlyVolumeService struct {}
// CreateBtKeywordMonthlyVolume 创建关键词月度搜索量记录
// Author [yourname](https://github.com/yourname)
func (btKeywordMonthlyVolumeService *BtKeywordMonthlyVolumeService) CreateBtKeywordMonthlyVolume(ctx context.Context, btKeywordMonthlyVolume *brandtrekin.BtKeywordMonthlyVolume) (err error) {
	err = global.GVA_DB.Create(btKeywordMonthlyVolume).Error
	return err
}

// DeleteBtKeywordMonthlyVolume 删除关键词月度搜索量记录
// Author [yourname](https://github.com/yourname)
func (btKeywordMonthlyVolumeService *BtKeywordMonthlyVolumeService)DeleteBtKeywordMonthlyVolume(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&brandtrekin.BtKeywordMonthlyVolume{},"id = ?",ID).Error
	return err
}

// DeleteBtKeywordMonthlyVolumeByIds 批量删除关键词月度搜索量记录
// Author [yourname](https://github.com/yourname)
func (btKeywordMonthlyVolumeService *BtKeywordMonthlyVolumeService)DeleteBtKeywordMonthlyVolumeByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]brandtrekin.BtKeywordMonthlyVolume{},"id in ?",IDs).Error
	return err
}

// UpdateBtKeywordMonthlyVolume 更新关键词月度搜索量记录
// Author [yourname](https://github.com/yourname)
func (btKeywordMonthlyVolumeService *BtKeywordMonthlyVolumeService)UpdateBtKeywordMonthlyVolume(ctx context.Context, btKeywordMonthlyVolume brandtrekin.BtKeywordMonthlyVolume) (err error) {
	err = global.GVA_DB.Model(&brandtrekin.BtKeywordMonthlyVolume{}).Where("id = ?",btKeywordMonthlyVolume.ID).Updates(&btKeywordMonthlyVolume).Error
	return err
}

// GetBtKeywordMonthlyVolume 根据ID获取关键词月度搜索量记录
// Author [yourname](https://github.com/yourname)
func (btKeywordMonthlyVolumeService *BtKeywordMonthlyVolumeService)GetBtKeywordMonthlyVolume(ctx context.Context, ID string) (btKeywordMonthlyVolume brandtrekin.BtKeywordMonthlyVolume, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&btKeywordMonthlyVolume).Error
	return
}
// GetBtKeywordMonthlyVolumeInfoList 分页获取关键词月度搜索量记录
// Author [yourname](https://github.com/yourname)
func (btKeywordMonthlyVolumeService *BtKeywordMonthlyVolumeService)GetBtKeywordMonthlyVolumeInfoList(ctx context.Context, info brandtrekinReq.BtKeywordMonthlyVolumeSearch) (list []brandtrekin.BtKeywordMonthlyVolume, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&brandtrekin.BtKeywordMonthlyVolume{})
    var btKeywordMonthlyVolumes []brandtrekin.BtKeywordMonthlyVolume
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.KeywordId != nil {
        db = db.Where("keyword_id = ?", *info.KeywordId)
    }
			if len(info.DateRange) == 2 {
				db = db.Where("date BETWEEN ? AND ? ", info.DateRange[0], info.DateRange[1])
			}
	if info.StartVolume != nil && info.EndVolume != nil {
		db = db.Where("volume BETWEEN ? AND ? ", *info.StartVolume, *info.EndVolume)
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
         	orderMap["volume"] = true
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

	err = db.Find(&btKeywordMonthlyVolumes).Error
	return  btKeywordMonthlyVolumes, total, err
}
func (btKeywordMonthlyVolumeService *BtKeywordMonthlyVolumeService)GetBtKeywordMonthlyVolumeDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   keywordId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("bt_keywords").Where("deleted_at IS NULL").Select("keyword as label,id as value").Scan(&keywordId)
	   res["keywordId"] = keywordId
	return
}
func (btKeywordMonthlyVolumeService *BtKeywordMonthlyVolumeService)GetBtKeywordMonthlyVolumePublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
