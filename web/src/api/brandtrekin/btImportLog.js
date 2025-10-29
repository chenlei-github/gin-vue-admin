import service from '@/utils/request'
// @Tags BtImportLog
// @Summary 创建数据导入日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtImportLog true "创建数据导入日志"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /btImportLog/createBtImportLog [post]
export const createBtImportLog = (data) => {
  return service({
    url: '/btImportLog/createBtImportLog',
    method: 'post',
    data
  })
}

// @Tags BtImportLog
// @Summary 删除数据导入日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtImportLog true "删除数据导入日志"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btImportLog/deleteBtImportLog [delete]
export const deleteBtImportLog = (params) => {
  return service({
    url: '/btImportLog/deleteBtImportLog',
    method: 'delete',
    params
  })
}

// @Tags BtImportLog
// @Summary 批量删除数据导入日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除数据导入日志"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btImportLog/deleteBtImportLog [delete]
export const deleteBtImportLogByIds = (params) => {
  return service({
    url: '/btImportLog/deleteBtImportLogByIds',
    method: 'delete',
    params
  })
}

// @Tags BtImportLog
// @Summary 更新数据导入日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtImportLog true "更新数据导入日志"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /btImportLog/updateBtImportLog [put]
export const updateBtImportLog = (data) => {
  return service({
    url: '/btImportLog/updateBtImportLog',
    method: 'put',
    data
  })
}

// @Tags BtImportLog
// @Summary 用id查询数据导入日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BtImportLog true "用id查询数据导入日志"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btImportLog/findBtImportLog [get]
export const findBtImportLog = (params) => {
  return service({
    url: '/btImportLog/findBtImportLog',
    method: 'get',
    params
  })
}

// @Tags BtImportLog
// @Summary 分页获取数据导入日志列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取数据导入日志列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /btImportLog/getBtImportLogList [get]
export const getBtImportLogList = (params) => {
  return service({
    url: '/btImportLog/getBtImportLogList',
    method: 'get',
    params
  })
}
// @Tags BtImportLog
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btImportLog/findBtImportLogDataSource [get]
export const getBtImportLogDataSource = () => {
  return service({
    url: '/btImportLog/getBtImportLogDataSource',
    method: 'get',
  })
}

// @Tags BtImportLog
// @Summary 不需要鉴权的数据导入日志接口
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtImportLogSearch true "分页获取数据导入日志列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btImportLog/getBtImportLogPublic [get]
export const getBtImportLogPublic = () => {
  return service({
    url: '/btImportLog/getBtImportLogPublic',
    method: 'get',
  })
}
