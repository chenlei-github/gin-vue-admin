import service from '@/utils/request'
// @Tags BtBrandMonthlyTrend
// @Summary 创建品牌月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtBrandMonthlyTrend true "创建品牌月度趋势"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /btBrandMonthlyTrend/createBtBrandMonthlyTrend [post]
export const createBtBrandMonthlyTrend = (data) => {
  return service({
    url: '/btBrandMonthlyTrend/createBtBrandMonthlyTrend',
    method: 'post',
    data
  })
}

// @Tags BtBrandMonthlyTrend
// @Summary 删除品牌月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtBrandMonthlyTrend true "删除品牌月度趋势"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btBrandMonthlyTrend/deleteBtBrandMonthlyTrend [delete]
export const deleteBtBrandMonthlyTrend = (params) => {
  return service({
    url: '/btBrandMonthlyTrend/deleteBtBrandMonthlyTrend',
    method: 'delete',
    params
  })
}

// @Tags BtBrandMonthlyTrend
// @Summary 批量删除品牌月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除品牌月度趋势"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btBrandMonthlyTrend/deleteBtBrandMonthlyTrend [delete]
export const deleteBtBrandMonthlyTrendByIds = (params) => {
  return service({
    url: '/btBrandMonthlyTrend/deleteBtBrandMonthlyTrendByIds',
    method: 'delete',
    params
  })
}

// @Tags BtBrandMonthlyTrend
// @Summary 更新品牌月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtBrandMonthlyTrend true "更新品牌月度趋势"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /btBrandMonthlyTrend/updateBtBrandMonthlyTrend [put]
export const updateBtBrandMonthlyTrend = (data) => {
  return service({
    url: '/btBrandMonthlyTrend/updateBtBrandMonthlyTrend',
    method: 'put',
    data
  })
}

// @Tags BtBrandMonthlyTrend
// @Summary 用id查询品牌月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BtBrandMonthlyTrend true "用id查询品牌月度趋势"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btBrandMonthlyTrend/findBtBrandMonthlyTrend [get]
export const findBtBrandMonthlyTrend = (params) => {
  return service({
    url: '/btBrandMonthlyTrend/findBtBrandMonthlyTrend',
    method: 'get',
    params
  })
}

// @Tags BtBrandMonthlyTrend
// @Summary 分页获取品牌月度趋势列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取品牌月度趋势列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /btBrandMonthlyTrend/getBtBrandMonthlyTrendList [get]
export const getBtBrandMonthlyTrendList = (params) => {
  return service({
    url: '/btBrandMonthlyTrend/getBtBrandMonthlyTrendList',
    method: 'get',
    params
  })
}
// @Tags BtBrandMonthlyTrend
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btBrandMonthlyTrend/findBtBrandMonthlyTrendDataSource [get]
export const getBtBrandMonthlyTrendDataSource = () => {
  return service({
    url: '/btBrandMonthlyTrend/getBtBrandMonthlyTrendDataSource',
    method: 'get',
  })
}

// @Tags BtBrandMonthlyTrend
// @Summary 不需要鉴权的品牌月度趋势接口
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtBrandMonthlyTrendSearch true "分页获取品牌月度趋势列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btBrandMonthlyTrend/getBtBrandMonthlyTrendPublic [get]
export const getBtBrandMonthlyTrendPublic = () => {
  return service({
    url: '/btBrandMonthlyTrend/getBtBrandMonthlyTrendPublic',
    method: 'get',
  })
}
