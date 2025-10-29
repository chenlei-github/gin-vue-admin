import service from '@/utils/request'
// @Tags BtMarketMonthlyTrend
// @Summary 创建市场月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtMarketMonthlyTrend true "创建市场月度趋势"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /btMarketMonthlyTrend/createBtMarketMonthlyTrend [post]
export const createBtMarketMonthlyTrend = (data) => {
  return service({
    url: '/btMarketMonthlyTrend/createBtMarketMonthlyTrend',
    method: 'post',
    data
  })
}

// @Tags BtMarketMonthlyTrend
// @Summary 删除市场月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtMarketMonthlyTrend true "删除市场月度趋势"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btMarketMonthlyTrend/deleteBtMarketMonthlyTrend [delete]
export const deleteBtMarketMonthlyTrend = (params) => {
  return service({
    url: '/btMarketMonthlyTrend/deleteBtMarketMonthlyTrend',
    method: 'delete',
    params
  })
}

// @Tags BtMarketMonthlyTrend
// @Summary 批量删除市场月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除市场月度趋势"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btMarketMonthlyTrend/deleteBtMarketMonthlyTrend [delete]
export const deleteBtMarketMonthlyTrendByIds = (params) => {
  return service({
    url: '/btMarketMonthlyTrend/deleteBtMarketMonthlyTrendByIds',
    method: 'delete',
    params
  })
}

// @Tags BtMarketMonthlyTrend
// @Summary 更新市场月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtMarketMonthlyTrend true "更新市场月度趋势"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /btMarketMonthlyTrend/updateBtMarketMonthlyTrend [put]
export const updateBtMarketMonthlyTrend = (data) => {
  return service({
    url: '/btMarketMonthlyTrend/updateBtMarketMonthlyTrend',
    method: 'put',
    data
  })
}

// @Tags BtMarketMonthlyTrend
// @Summary 用id查询市场月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BtMarketMonthlyTrend true "用id查询市场月度趋势"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btMarketMonthlyTrend/findBtMarketMonthlyTrend [get]
export const findBtMarketMonthlyTrend = (params) => {
  return service({
    url: '/btMarketMonthlyTrend/findBtMarketMonthlyTrend',
    method: 'get',
    params
  })
}

// @Tags BtMarketMonthlyTrend
// @Summary 分页获取市场月度趋势列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取市场月度趋势列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /btMarketMonthlyTrend/getBtMarketMonthlyTrendList [get]
export const getBtMarketMonthlyTrendList = (params) => {
  return service({
    url: '/btMarketMonthlyTrend/getBtMarketMonthlyTrendList',
    method: 'get',
    params
  })
}
// @Tags BtMarketMonthlyTrend
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btMarketMonthlyTrend/findBtMarketMonthlyTrendDataSource [get]
export const getBtMarketMonthlyTrendDataSource = () => {
  return service({
    url: '/btMarketMonthlyTrend/getBtMarketMonthlyTrendDataSource',
    method: 'get',
  })
}

// @Tags BtMarketMonthlyTrend
// @Summary 不需要鉴权的市场月度趋势接口
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtMarketMonthlyTrendSearch true "分页获取市场月度趋势列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btMarketMonthlyTrend/getBtMarketMonthlyTrendPublic [get]
export const getBtMarketMonthlyTrendPublic = () => {
  return service({
    url: '/btMarketMonthlyTrend/getBtMarketMonthlyTrendPublic',
    method: 'get',
  })
}
