import service from '@/utils/request'
// @Tags BtProduct
// @Summary 创建商品管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtProduct true "创建商品管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /btProduct/createBtProduct [post]
export const createBtProduct = (data) => {
  return service({
    url: '/btProduct/createBtProduct',
    method: 'post',
    data
  })
}

// @Tags BtProduct
// @Summary 删除商品管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtProduct true "删除商品管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btProduct/deleteBtProduct [delete]
export const deleteBtProduct = (params) => {
  return service({
    url: '/btProduct/deleteBtProduct',
    method: 'delete',
    params
  })
}

// @Tags BtProduct
// @Summary 批量删除商品管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除商品管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btProduct/deleteBtProduct [delete]
export const deleteBtProductByIds = (params) => {
  return service({
    url: '/btProduct/deleteBtProductByIds',
    method: 'delete',
    params
  })
}

// @Tags BtProduct
// @Summary 更新商品管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtProduct true "更新商品管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /btProduct/updateBtProduct [put]
export const updateBtProduct = (data) => {
  return service({
    url: '/btProduct/updateBtProduct',
    method: 'put',
    data
  })
}

// @Tags BtProduct
// @Summary 用id查询商品管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BtProduct true "用id查询商品管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btProduct/findBtProduct [get]
export const findBtProduct = (params) => {
  return service({
    url: '/btProduct/findBtProduct',
    method: 'get',
    params
  })
}

// @Tags BtProduct
// @Summary 分页获取商品管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取商品管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /btProduct/getBtProductList [get]
export const getBtProductList = (params) => {
  return service({
    url: '/btProduct/getBtProductList',
    method: 'get',
    params
  })
}
// @Tags BtProduct
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btProduct/findBtProductDataSource [get]
export const getBtProductDataSource = () => {
  return service({
    url: '/btProduct/getBtProductDataSource',
    method: 'get',
  })
}

// @Tags BtProduct
// @Summary 不需要鉴权的商品管理接口
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtProductSearch true "分页获取商品管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btProduct/getBtProductPublic [get]
export const getBtProductPublic = () => {
  return service({
    url: '/btProduct/getBtProductPublic',
    method: 'get',
  })
}
