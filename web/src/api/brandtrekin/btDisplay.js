import service from '@/utils/request'

// @Tags BtDisplay
// @Summary 获取市场列表
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]response.MarketListItem,msg=string} "获取成功"
// @Router /api/markets [get]
export const getMarketList = () => {
  return service({
    url: '/api/markets',
    method: 'get'
  })
}

// @Tags BtDisplay
// @Summary 获取市场详情
// @Accept application/json
// @Produce application/json
// @Param id path string true "市场ID"
// @Success 200 {object} response.Response{data=response.MarketDetail,msg=string} "获取成功"
// @Router /api/markets/:id [get]
export const getMarketDetail = (id) => {
  return service({
    url: `/api/markets/${id}`,
    method: 'get'
  })
}

// @Tags BtDisplay
// @Summary 获取品牌详情
// @Accept application/json
// @Produce application/json
// @Param marketId path string true "市场ID"
// @Param brandName path string true "品牌名称"
// @Success 200 {object} response.Response{data=response.BrandDetail,msg=string} "获取成功"
// @Router /api/markets/:marketId/brands/:brandName [get]
export const getBrandDetail = (marketId, brandName) => {
  return service({
    url: `/api/markets/brands/${marketId}/${encodeURIComponent(brandName)}`,
    method: 'get'
  })
}
