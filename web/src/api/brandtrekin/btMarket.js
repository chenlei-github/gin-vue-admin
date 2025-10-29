import service from '@/utils/request'
// @Tags BtMarket
// @Summary 创建市场管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtMarket true "创建市场管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /btMarket/createBtMarket [post]
export const createBtMarket = (data) => {
  return service({
    url: '/btMarket/createBtMarket',
    method: 'post',
    data
  })
}

// @Tags BtMarket
// @Summary 删除市场管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtMarket true "删除市场管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btMarket/deleteBtMarket [delete]
export const deleteBtMarket = (params) => {
  return service({
    url: '/btMarket/deleteBtMarket',
    method: 'delete',
    params
  })
}

// @Tags BtMarket
// @Summary 批量删除市场管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除市场管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btMarket/deleteBtMarket [delete]
export const deleteBtMarketByIds = (params) => {
  return service({
    url: '/btMarket/deleteBtMarketByIds',
    method: 'delete',
    params
  })
}

// @Tags BtMarket
// @Summary 更新市场管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtMarket true "更新市场管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /btMarket/updateBtMarket [put]
export const updateBtMarket = (data) => {
  return service({
    url: '/btMarket/updateBtMarket',
    method: 'put',
    data
  })
}

// @Tags BtMarket
// @Summary 用id查询市场管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BtMarket true "用id查询市场管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btMarket/findBtMarket [get]
export const findBtMarket = (params) => {
  return service({
    url: '/btMarket/findBtMarket',
    method: 'get',
    params
  })
}

// @Tags BtMarket
// @Summary 分页获取市场管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取市场管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /btMarket/getBtMarketList [get]
export const getBtMarketList = (params) => {
  return service({
    url: '/btMarket/getBtMarketList',
    method: 'get',
    params
  })
}

// @Tags BtMarket
// @Summary 不需要鉴权的市场管理接口
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtMarketSearch true "分页获取市场管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btMarket/getBtMarketPublic [get]
export const getBtMarketPublic = () => {
  return service({
    url: '/btMarket/getBtMarketPublic',
    method: 'get',
  })
}
