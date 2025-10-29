
package brandtrekin

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
)

type BtBrandSocialMediaService struct {}
// CreateBtBrandSocialMedia 创建品牌社交媒体记录
// Author [yourname](https://github.com/yourname)
func (btBrandSocialMediaService *BtBrandSocialMediaService) CreateBtBrandSocialMedia(ctx context.Context, btBrandSocialMedia *brandtrekin.BtBrandSocialMedia) (err error) {
	err = global.GVA_DB.Create(btBrandSocialMedia).Error
	return err
}

// DeleteBtBrandSocialMedia 删除品牌社交媒体记录
// Author [yourname](https://github.com/yourname)
func (btBrandSocialMediaService *BtBrandSocialMediaService)DeleteBtBrandSocialMedia(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&brandtrekin.BtBrandSocialMedia{},"id = ?",ID).Error
	return err
}

// DeleteBtBrandSocialMediaByIds 批量删除品牌社交媒体记录
// Author [yourname](https://github.com/yourname)
func (btBrandSocialMediaService *BtBrandSocialMediaService)DeleteBtBrandSocialMediaByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]brandtrekin.BtBrandSocialMedia{},"id in ?",IDs).Error
	return err
}

// UpdateBtBrandSocialMedia 更新品牌社交媒体记录
// Author [yourname](https://github.com/yourname)
func (btBrandSocialMediaService *BtBrandSocialMediaService)UpdateBtBrandSocialMedia(ctx context.Context, btBrandSocialMedia brandtrekin.BtBrandSocialMedia) (err error) {
	err = global.GVA_DB.Model(&brandtrekin.BtBrandSocialMedia{}).Where("id = ?",btBrandSocialMedia.ID).Updates(&btBrandSocialMedia).Error
	return err
}

// GetBtBrandSocialMedia 根据ID获取品牌社交媒体记录
// Author [yourname](https://github.com/yourname)
func (btBrandSocialMediaService *BtBrandSocialMediaService)GetBtBrandSocialMedia(ctx context.Context, ID string) (btBrandSocialMedia brandtrekin.BtBrandSocialMedia, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&btBrandSocialMedia).Error
	return
}
// GetBtBrandSocialMediaInfoList 分页获取品牌社交媒体记录
// Author [yourname](https://github.com/yourname)
func (btBrandSocialMediaService *BtBrandSocialMediaService)GetBtBrandSocialMediaInfoList(ctx context.Context, info brandtrekinReq.BtBrandSocialMediaSearch) (list []brandtrekin.BtBrandSocialMedia, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&brandtrekin.BtBrandSocialMedia{})
    var btBrandSocialMedias []brandtrekin.BtBrandSocialMedia
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.BrandId != nil {
        db = db.Where("brand_id = ?", *info.BrandId)
    }
    if info.Platform != nil && *info.Platform != "" {
        db = db.Where("platform = ?", *info.Platform)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&btBrandSocialMedias).Error
	return  btBrandSocialMedias, total, err
}
func (btBrandSocialMediaService *BtBrandSocialMediaService)GetBtBrandSocialMediaDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   brandId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("bt_brands").Where("deleted_at IS NULL").Select("brand_name as label,id as value").Scan(&brandId)
	   res["brandId"] = brandId
	return
}
func (btBrandSocialMediaService *BtBrandSocialMediaService)GetBtBrandSocialMediaPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
