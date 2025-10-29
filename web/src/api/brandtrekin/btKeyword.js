import service from '@/utils/request'
// @Tags BtKeyword
// @Summary 创建关键词管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtKeyword true "创建关键词管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /btKeyword/createBtKeyword [post]
export const createBtKeyword = (data) => {
  return service({
    url: '/btKeyword/createBtKeyword',
    method: 'post',
    data
  })
}

// @Tags BtKeyword
// @Summary 删除关键词管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtKeyword true "删除关键词管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btKeyword/deleteBtKeyword [delete]
export const deleteBtKeyword = (params) => {
  return service({
    url: '/btKeyword/deleteBtKeyword',
    method: 'delete',
    params
  })
}

// @Tags BtKeyword
// @Summary 批量删除关键词管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除关键词管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btKeyword/deleteBtKeyword [delete]
export const deleteBtKeywordByIds = (params) => {
  return service({
    url: '/btKeyword/deleteBtKeywordByIds',
    method: 'delete',
    params
  })
}

// @Tags BtKeyword
// @Summary 更新关键词管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtKeyword true "更新关键词管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /btKeyword/updateBtKeyword [put]
export const updateBtKeyword = (data) => {
  return service({
    url: '/btKeyword/updateBtKeyword',
    method: 'put',
    data
  })
}

// @Tags BtKeyword
// @Summary 用id查询关键词管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BtKeyword true "用id查询关键词管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btKeyword/findBtKeyword [get]
export const findBtKeyword = (params) => {
  return service({
    url: '/btKeyword/findBtKeyword',
    method: 'get',
    params
  })
}

// @Tags BtKeyword
// @Summary 分页获取关键词管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取关键词管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /btKeyword/getBtKeywordList [get]
export const getBtKeywordList = (params) => {
  return service({
    url: '/btKeyword/getBtKeywordList',
    method: 'get',
    params
  })
}
// @Tags BtKeyword
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btKeyword/findBtKeywordDataSource [get]
export const getBtKeywordDataSource = () => {
  return service({
    url: '/btKeyword/getBtKeywordDataSource',
    method: 'get',
  })
}

// @Tags BtKeyword
// @Summary 不需要鉴权的关键词管理接口
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtKeywordSearch true "分页获取关键词管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btKeyword/getBtKeywordPublic [get]
export const getBtKeywordPublic = () => {
  return service({
    url: '/btKeyword/getBtKeywordPublic',
    method: 'get',
  })
}
