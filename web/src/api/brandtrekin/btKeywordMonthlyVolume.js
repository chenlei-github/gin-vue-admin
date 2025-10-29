import service from '@/utils/request'
// @Tags BtKeywordMonthlyVolume
// @Summary 创建关键词月度搜索量
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtKeywordMonthlyVolume true "创建关键词月度搜索量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /btKeywordMonthlyVolume/createBtKeywordMonthlyVolume [post]
export const createBtKeywordMonthlyVolume = (data) => {
  return service({
    url: '/btKeywordMonthlyVolume/createBtKeywordMonthlyVolume',
    method: 'post',
    data
  })
}

// @Tags BtKeywordMonthlyVolume
// @Summary 删除关键词月度搜索量
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtKeywordMonthlyVolume true "删除关键词月度搜索量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btKeywordMonthlyVolume/deleteBtKeywordMonthlyVolume [delete]
export const deleteBtKeywordMonthlyVolume = (params) => {
  return service({
    url: '/btKeywordMonthlyVolume/deleteBtKeywordMonthlyVolume',
    method: 'delete',
    params
  })
}

// @Tags BtKeywordMonthlyVolume
// @Summary 批量删除关键词月度搜索量
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除关键词月度搜索量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btKeywordMonthlyVolume/deleteBtKeywordMonthlyVolume [delete]
export const deleteBtKeywordMonthlyVolumeByIds = (params) => {
  return service({
    url: '/btKeywordMonthlyVolume/deleteBtKeywordMonthlyVolumeByIds',
    method: 'delete',
    params
  })
}

// @Tags BtKeywordMonthlyVolume
// @Summary 更新关键词月度搜索量
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtKeywordMonthlyVolume true "更新关键词月度搜索量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /btKeywordMonthlyVolume/updateBtKeywordMonthlyVolume [put]
export const updateBtKeywordMonthlyVolume = (data) => {
  return service({
    url: '/btKeywordMonthlyVolume/updateBtKeywordMonthlyVolume',
    method: 'put',
    data
  })
}

// @Tags BtKeywordMonthlyVolume
// @Summary 用id查询关键词月度搜索量
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BtKeywordMonthlyVolume true "用id查询关键词月度搜索量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btKeywordMonthlyVolume/findBtKeywordMonthlyVolume [get]
export const findBtKeywordMonthlyVolume = (params) => {
  return service({
    url: '/btKeywordMonthlyVolume/findBtKeywordMonthlyVolume',
    method: 'get',
    params
  })
}

// @Tags BtKeywordMonthlyVolume
// @Summary 分页获取关键词月度搜索量列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取关键词月度搜索量列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /btKeywordMonthlyVolume/getBtKeywordMonthlyVolumeList [get]
export const getBtKeywordMonthlyVolumeList = (params) => {
  return service({
    url: '/btKeywordMonthlyVolume/getBtKeywordMonthlyVolumeList',
    method: 'get',
    params
  })
}
// @Tags BtKeywordMonthlyVolume
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btKeywordMonthlyVolume/findBtKeywordMonthlyVolumeDataSource [get]
export const getBtKeywordMonthlyVolumeDataSource = () => {
  return service({
    url: '/btKeywordMonthlyVolume/getBtKeywordMonthlyVolumeDataSource',
    method: 'get',
  })
}

// @Tags BtKeywordMonthlyVolume
// @Summary 不需要鉴权的关键词月度搜索量接口
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtKeywordMonthlyVolumeSearch true "分页获取关键词月度搜索量列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btKeywordMonthlyVolume/getBtKeywordMonthlyVolumePublic [get]
export const getBtKeywordMonthlyVolumePublic = () => {
  return service({
    url: '/btKeywordMonthlyVolume/getBtKeywordMonthlyVolumePublic',
    method: 'get',
  })
}
