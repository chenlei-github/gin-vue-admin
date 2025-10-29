import service from '@/utils/request'
// @Tags BtProductMonthlySales
// @Summary 创建商品月度销售
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtProductMonthlySales true "创建商品月度销售"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /btProductMonthlySales/createBtProductMonthlySales [post]
export const createBtProductMonthlySales = (data) => {
  return service({
    url: '/btProductMonthlySales/createBtProductMonthlySales',
    method: 'post',
    data
  })
}

// @Tags BtProductMonthlySales
// @Summary 删除商品月度销售
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtProductMonthlySales true "删除商品月度销售"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btProductMonthlySales/deleteBtProductMonthlySales [delete]
export const deleteBtProductMonthlySales = (params) => {
  return service({
    url: '/btProductMonthlySales/deleteBtProductMonthlySales',
    method: 'delete',
    params
  })
}

// @Tags BtProductMonthlySales
// @Summary 批量删除商品月度销售
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除商品月度销售"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btProductMonthlySales/deleteBtProductMonthlySales [delete]
export const deleteBtProductMonthlySalesByIds = (params) => {
  return service({
    url: '/btProductMonthlySales/deleteBtProductMonthlySalesByIds',
    method: 'delete',
    params
  })
}

// @Tags BtProductMonthlySales
// @Summary 更新商品月度销售
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtProductMonthlySales true "更新商品月度销售"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /btProductMonthlySales/updateBtProductMonthlySales [put]
export const updateBtProductMonthlySales = (data) => {
  return service({
    url: '/btProductMonthlySales/updateBtProductMonthlySales',
    method: 'put',
    data
  })
}

// @Tags BtProductMonthlySales
// @Summary 用id查询商品月度销售
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BtProductMonthlySales true "用id查询商品月度销售"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btProductMonthlySales/findBtProductMonthlySales [get]
export const findBtProductMonthlySales = (params) => {
  return service({
    url: '/btProductMonthlySales/findBtProductMonthlySales',
    method: 'get',
    params
  })
}

// @Tags BtProductMonthlySales
// @Summary 分页获取商品月度销售列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取商品月度销售列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /btProductMonthlySales/getBtProductMonthlySalesList [get]
export const getBtProductMonthlySalesList = (params) => {
  return service({
    url: '/btProductMonthlySales/getBtProductMonthlySalesList',
    method: 'get',
    params
  })
}

// @Tags BtProductMonthlySales
// @Summary 不需要鉴权的商品月度销售接口
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtProductMonthlySalesSearch true "分页获取商品月度销售列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btProductMonthlySales/getBtProductMonthlySalesPublic [get]
export const getBtProductMonthlySalesPublic = () => {
  return service({
    url: '/btProductMonthlySales/getBtProductMonthlySalesPublic',
    method: 'get',
  })
}
