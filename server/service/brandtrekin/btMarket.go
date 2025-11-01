
package brandtrekin

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
	brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
	"gorm.io/gorm"
)

type BtMarketService struct {}
// CreateBtMarket 创建市场管理记录
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService) CreateBtMarket(ctx context.Context, btMarket *brandtrekin.BtMarket) (err error) {
	err = global.GVA_DB.Create(btMarket).Error
	return err
}

// DeleteBtMarket 删除市场管理记录（级联删除所有关联数据）
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService) DeleteBtMarket(ctx context.Context, ID string) (err error) {
	marketID, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return fmt.Errorf("无效的市场ID: %v", err)
	}

	// 使用事务确保数据一致性
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 调用级联删除方法
		if err := btMarketService.deleteMarketData(tx, marketID); err != nil {
			return err
		}

		// 最后删除市场记录本身
		if err := tx.Model(&brandtrekin.BtMarket{}).Where("id = ?", marketID).Delete(&brandtrekin.BtMarket{}).Error; err != nil {
			return fmt.Errorf("删除市场记录失败: %v", err)
		}

		return nil
	})
}

// deleteMarketData 删除市场的所有关联数据（从子表到父表）
func (btMarketService *BtMarketService) deleteMarketData(tx *gorm.DB, marketID int64) error {
	// 按顺序删除关联数据（从子表到父表）

	// 1. 删除商品月度销售数据
	if err := tx.Model(&brandtrekin.BtProductMonthlySales{}).Where("asin IN (SELECT asin FROM bt_products WHERE market_id = ?)", marketID).
		Delete(&brandtrekin.BtProductMonthlySales{}).Error; err != nil {
		return fmt.Errorf("删除商品月度销售数据失败: %v", err)
	}

	// 2. 删除商品数据
	if err := tx.Model(&brandtrekin.BtProduct{}).Where("market_id = ?", marketID).Delete(&brandtrekin.BtProduct{}).Error; err != nil {
		return fmt.Errorf("删除商品数据失败: %v", err)
	}

	// 3. 删除品牌月度趋势数据
	if err := tx.Model(&brandtrekin.BtBrandMonthlyTrend{}).Where("brand_id IN (SELECT id FROM bt_brands WHERE market_id = ?)", marketID).
		Delete(&brandtrekin.BtBrandMonthlyTrend{}).Error; err != nil {
		return fmt.Errorf("删除品牌月度趋势数据失败: %v", err)
	}

	// 4. 删除品牌社交媒体数据
	if err := tx.Model(&brandtrekin.BtBrandSocialMedia{}).Where("brand_id IN (SELECT id FROM bt_brands WHERE market_id = ?)", marketID).
		Delete(&brandtrekin.BtBrandSocialMedia{}).Error; err != nil {
		return fmt.Errorf("删除品牌社交媒体数据失败: %v", err)
	}

	// 5. 删除品牌数据
	if err := tx.Model(&brandtrekin.BtBrand{}).Where("market_id = ?", marketID).Delete(&brandtrekin.BtBrand{}).Error; err != nil {
		return fmt.Errorf("删除品牌数据失败: %v", err)
	}

	// 6. 删除关键词月度搜索量数据
	if err := tx.Model(&brandtrekin.BtKeywordMonthlyVolume{}).Where("keyword_id IN (SELECT id FROM bt_keywords WHERE market_id = ?)", marketID).
		Delete(&brandtrekin.BtKeywordMonthlyVolume{}).Error; err != nil {
		return fmt.Errorf("删除关键词月度搜索量数据失败: %v", err)
	}

	// 7. 删除关键词数据
	if err := tx.Model(&brandtrekin.BtKeyword{}).Where("market_id = ?", marketID).Delete(&brandtrekin.BtKeyword{}).Error; err != nil {
		return fmt.Errorf("删除关键词数据失败: %v", err)
	}

	// 8. 删除市场月度趋势数据
	if err := tx.Model(&brandtrekin.BtMarketMonthlyTrend{}).Where("market_id = ?", marketID).Delete(&brandtrekin.BtMarketMonthlyTrend{}).Error; err != nil {
		return fmt.Errorf("删除市场月度趋势数据失败: %v", err)
	}

	// 9. 删除导入日志（如果存在该表）
	// 注意：这里假设有导入日志表，如果没有可以注释掉
	// if err := tx.Where("market_id = ?", marketID).Delete(&brandtrekin.BtImportLog{}).Error; err != nil {
	// 	return fmt.Errorf("删除导入日志失败: %v", err)
	// }

	return nil
}

// DeleteBtMarketByIds 批量删除市场管理记录（级联删除所有关联数据）
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService) DeleteBtMarketByIds(ctx context.Context, IDs []string) (err error) {
	// 使用事务确保数据一致性
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 遍历每个市场ID，进行级联删除
		for _, ID := range IDs {
			marketID, err := strconv.ParseInt(ID, 10, 64)
			if err != nil {
				return fmt.Errorf("无效的市场ID: %s, %v", ID, err)
			}

			// 调用级联删除方法
			if err := btMarketService.deleteMarketData(tx, marketID); err != nil {
				return err
			}

			// 删除市场记录本身
			if err := tx.Model(&brandtrekin.BtMarket{}).Where("id = ?", marketID).Delete(&brandtrekin.BtMarket{}).Error; err != nil {
				return fmt.Errorf("删除市场记录失败 (ID: %s): %v", ID, err)
			}
		}

		return nil
	})
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
func (btMarketService *BtMarketService) GetBtMarketPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

// GenerateSlugFromName 根据市场名称自动生成slug
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService) GenerateSlugFromName(ctx context.Context, name string) (slug string, err error) {
	if name == "" {
		return "", fmt.Errorf("市场名称不能为空")
	}

	// 转换为小写
	slug = strings.ToLower(name)

	// 移除特殊字符，只保留字母、数字、空格和连字符
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = reg.ReplaceAllString(slug, "")

	// 将空格替换为连字符
	slug = strings.ReplaceAll(slug, " ", "-")

	// 将多个连字符合并为一个
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// 去除首尾的连字符
	slug = strings.Trim(slug, "-")

	if slug == "" {
		return "", fmt.Errorf("生成的市场ID无效")
	}

	return slug, nil
}

// ValidateSlugUnique 校验slug唯一性
// Author [yourname](https://github.com/yourname)
func (btMarketService *BtMarketService) ValidateSlugUnique(ctx context.Context, slug string, excludeID uint) (isUnique bool, err error) {
	if slug == "" {
		return false, fmt.Errorf("市场ID不能为空")
	}

	var count int64
	query := global.GVA_DB.Model(&brandtrekin.BtMarket{}).Where("market_slug = ?", slug)

	// 如果提供了excludeID，则排除该ID（用于编辑时检查）
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("校验失败: %v", err)
	}

	return count == 0, nil
}
