package brandtrekin

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BtDisplayApi struct{}

// GetMarketList 获取市场列表（公开接口）
// @Tags BtDisplay
// @Summary 获取市场列表
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]response.MarketListItem,msg=string} "获取成功"
// @Router /api/markets [get]
func (api *BtDisplayApi) GetMarketList(c *gin.Context) {
	marketList, err := btDisplayService.GetMarketList()
	if err != nil {
		global.GVA_LOG.Error("获取市场列表失败!", zap.Error(err))
		response.FailWithMessage("获取市场列表失败: "+err.Error(), c)
		return
	}

	response.OkWithData(marketList, c)
}

// GetMarketDetail 获取市场详情（公开接口）
// @Tags BtDisplay
// @Summary 获取市场详情
// @Accept application/json
// @Produce application/json
// @Param id path string true "市场ID (slug)"
// @Success 200 {object} response.Response{data=response.MarketDetail,msg=string} "获取成功"
// @Router /api/markets/{id} [get]
func (api *BtDisplayApi) GetMarketDetail(c *gin.Context) {
	marketSlug := c.Param("id")
	if marketSlug == "" {
		response.FailWithMessage("市场ID不能为空", c)
		return
	}

	detail, err := btDisplayService.GetMarketDetail(marketSlug)
	if err != nil {
		global.GVA_LOG.Error("获取市场详情失败!", zap.Error(err))
		response.FailWithMessage("获取市场详情失败: "+err.Error(), c)
		return
	}

	response.OkWithData(detail, c)
}

// GetBrandDetail 获取品牌详情（公开接口）
// @Tags BtDisplay
// @Summary 获取品牌详情
// @Accept application/json
// @Produce application/json
// @Param marketId path string true "市场ID (slug)"
// @Param brandName path string true "品牌名称"
// @Success 200 {object} response.Response{data=response.BrandDetail,msg=string} "获取成功"
// @Router /api/markets/{marketId}/brands/{brandName} [get]
func (api *BtDisplayApi) GetBrandDetail(c *gin.Context) {
	marketSlug := c.Param("marketId")
	brandName := c.Param("brandName")

	if marketSlug == "" || brandName == "" {
		response.FailWithMessage("市场ID和品牌名称不能为空", c)
		return
	}

	detail, err := btDisplayService.GetBrandDetail(marketSlug, brandName)
	if err != nil {
		global.GVA_LOG.Error("获取品牌详情失败!", zap.Error(err))
		response.FailWithMessage("获取品牌详情失败: "+err.Error(), c)
		return
	}

	response.OkWithData(detail, c)
}
