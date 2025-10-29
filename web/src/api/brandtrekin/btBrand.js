import service from '@/utils/request'
// @Tags BtBrand
// @Summary 创建品牌管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtBrand true "创建品牌管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /btBrand/createBtBrand [post]
export const createBtBrand = (data) => {
  return service({
    url: '/btBrand/createBtBrand',
    method: 'post',
    data
  })
}

// @Tags BtBrand
// @Summary 删除品牌管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtBrand true "删除品牌管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btBrand/deleteBtBrand [delete]
export const deleteBtBrand = (params) => {
  return service({
    url: '/btBrand/deleteBtBrand',
    method: 'delete',
    params
  })
}

// @Tags BtBrand
// @Summary 批量删除品牌管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除品牌管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btBrand/deleteBtBrand [delete]
export const deleteBtBrandByIds = (params) => {
  return service({
    url: '/btBrand/deleteBtBrandByIds',
    method: 'delete',
    params
  })
}

// @Tags BtBrand
// @Summary 更新品牌管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtBrand true "更新品牌管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /btBrand/updateBtBrand [put]
export const updateBtBrand = (data) => {
  return service({
    url: '/btBrand/updateBtBrand',
    method: 'put',
    data
  })
}

// @Tags BtBrand
// @Summary 用id查询品牌管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BtBrand true "用id查询品牌管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btBrand/findBtBrand [get]
export const findBtBrand = (params) => {
  return service({
    url: '/btBrand/findBtBrand',
    method: 'get',
    params
  })
}

// @Tags BtBrand
// @Summary 分页获取品牌管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取品牌管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /btBrand/getBtBrandList [get]
export const getBtBrandList = (params) => {
  return service({
    url: '/btBrand/getBtBrandList',
    method: 'get',
    params
  })
}
// @Tags BtBrand
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btBrand/findBtBrandDataSource [get]
export const getBtBrandDataSource = () => {
  return service({
    url: '/btBrand/getBtBrandDataSource',
    method: 'get',
  })
}

// @Tags BtBrand
// @Summary 不需要鉴权的品牌管理接口
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtBrandSearch true "分页获取品牌管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btBrand/getBtBrandPublic [get]
export const getBtBrandPublic = () => {
  return service({
    url: '/btBrand/getBtBrandPublic',
    method: 'get',
  })
}
