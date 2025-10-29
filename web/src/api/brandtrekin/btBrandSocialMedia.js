import service from '@/utils/request'
// @Tags BtBrandSocialMedia
// @Summary 创建品牌社交媒体
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtBrandSocialMedia true "创建品牌社交媒体"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /btBrandSocialMedia/createBtBrandSocialMedia [post]
export const createBtBrandSocialMedia = (data) => {
  return service({
    url: '/btBrandSocialMedia/createBtBrandSocialMedia',
    method: 'post',
    data
  })
}

// @Tags BtBrandSocialMedia
// @Summary 删除品牌社交媒体
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtBrandSocialMedia true "删除品牌社交媒体"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btBrandSocialMedia/deleteBtBrandSocialMedia [delete]
export const deleteBtBrandSocialMedia = (params) => {
  return service({
    url: '/btBrandSocialMedia/deleteBtBrandSocialMedia',
    method: 'delete',
    params
  })
}

// @Tags BtBrandSocialMedia
// @Summary 批量删除品牌社交媒体
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除品牌社交媒体"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /btBrandSocialMedia/deleteBtBrandSocialMedia [delete]
export const deleteBtBrandSocialMediaByIds = (params) => {
  return service({
    url: '/btBrandSocialMedia/deleteBtBrandSocialMediaByIds',
    method: 'delete',
    params
  })
}

// @Tags BtBrandSocialMedia
// @Summary 更新品牌社交媒体
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BtBrandSocialMedia true "更新品牌社交媒体"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /btBrandSocialMedia/updateBtBrandSocialMedia [put]
export const updateBtBrandSocialMedia = (data) => {
  return service({
    url: '/btBrandSocialMedia/updateBtBrandSocialMedia',
    method: 'put',
    data
  })
}

// @Tags BtBrandSocialMedia
// @Summary 用id查询品牌社交媒体
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BtBrandSocialMedia true "用id查询品牌社交媒体"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btBrandSocialMedia/findBtBrandSocialMedia [get]
export const findBtBrandSocialMedia = (params) => {
  return service({
    url: '/btBrandSocialMedia/findBtBrandSocialMedia',
    method: 'get',
    params
  })
}

// @Tags BtBrandSocialMedia
// @Summary 分页获取品牌社交媒体列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取品牌社交媒体列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /btBrandSocialMedia/getBtBrandSocialMediaList [get]
export const getBtBrandSocialMediaList = (params) => {
  return service({
    url: '/btBrandSocialMedia/getBtBrandSocialMediaList',
    method: 'get',
    params
  })
}
// @Tags BtBrandSocialMedia
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /btBrandSocialMedia/findBtBrandSocialMediaDataSource [get]
export const getBtBrandSocialMediaDataSource = () => {
  return service({
    url: '/btBrandSocialMedia/getBtBrandSocialMediaDataSource',
    method: 'get',
  })
}

// @Tags BtBrandSocialMedia
// @Summary 不需要鉴权的品牌社交媒体接口
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtBrandSocialMediaSearch true "分页获取品牌社交媒体列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btBrandSocialMedia/getBtBrandSocialMediaPublic [get]
export const getBtBrandSocialMediaPublic = () => {
  return service({
    url: '/btBrandSocialMedia/getBtBrandSocialMediaPublic',
    method: 'get',
  })
}
