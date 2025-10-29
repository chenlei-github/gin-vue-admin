package brandtrekin

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BtBrandSocialMediaApi struct {}



// CreateBtBrandSocialMedia 创建品牌社交媒体
// @Tags BtBrandSocialMedia
// @Summary 创建品牌社交媒体
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtBrandSocialMedia true "创建品牌社交媒体"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /btBrandSocialMedia/createBtBrandSocialMedia [post]
func (btBrandSocialMediaApi *BtBrandSocialMediaApi) CreateBtBrandSocialMedia(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var btBrandSocialMedia brandtrekin.BtBrandSocialMedia
	err := c.ShouldBindJSON(&btBrandSocialMedia)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btBrandSocialMediaService.CreateBtBrandSocialMedia(ctx,&btBrandSocialMedia)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBtBrandSocialMedia 删除品牌社交媒体
// @Tags BtBrandSocialMedia
// @Summary 删除品牌社交媒体
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtBrandSocialMedia true "删除品牌社交媒体"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /btBrandSocialMedia/deleteBtBrandSocialMedia [delete]
func (btBrandSocialMediaApi *BtBrandSocialMediaApi) DeleteBtBrandSocialMedia(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := btBrandSocialMediaService.DeleteBtBrandSocialMedia(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBtBrandSocialMediaByIds 批量删除品牌社交媒体
// @Tags BtBrandSocialMedia
// @Summary 批量删除品牌社交媒体
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /btBrandSocialMedia/deleteBtBrandSocialMediaByIds [delete]
func (btBrandSocialMediaApi *BtBrandSocialMediaApi) DeleteBtBrandSocialMediaByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := btBrandSocialMediaService.DeleteBtBrandSocialMediaByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBtBrandSocialMedia 更新品牌社交媒体
// @Tags BtBrandSocialMedia
// @Summary 更新品牌社交媒体
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtBrandSocialMedia true "更新品牌社交媒体"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /btBrandSocialMedia/updateBtBrandSocialMedia [put]
func (btBrandSocialMediaApi *BtBrandSocialMediaApi) UpdateBtBrandSocialMedia(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var btBrandSocialMedia brandtrekin.BtBrandSocialMedia
	err := c.ShouldBindJSON(&btBrandSocialMedia)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btBrandSocialMediaService.UpdateBtBrandSocialMedia(ctx,btBrandSocialMedia)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBtBrandSocialMedia 用id查询品牌社交媒体
// @Tags BtBrandSocialMedia
// @Summary 用id查询品牌社交媒体
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询品牌社交媒体"
// @Success 200 {object} response.Response{data=brandtrekin.BtBrandSocialMedia,msg=string} "查询成功"
// @Router /btBrandSocialMedia/findBtBrandSocialMedia [get]
func (btBrandSocialMediaApi *BtBrandSocialMediaApi) FindBtBrandSocialMedia(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebtBrandSocialMedia, err := btBrandSocialMediaService.GetBtBrandSocialMedia(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebtBrandSocialMedia, c)
}
// GetBtBrandSocialMediaList 分页获取品牌社交媒体列表
// @Tags BtBrandSocialMedia
// @Summary 分页获取品牌社交媒体列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtBrandSocialMediaSearch true "分页获取品牌社交媒体列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /btBrandSocialMedia/getBtBrandSocialMediaList [get]
func (btBrandSocialMediaApi *BtBrandSocialMediaApi) GetBtBrandSocialMediaList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo brandtrekinReq.BtBrandSocialMediaSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := btBrandSocialMediaService.GetBtBrandSocialMediaInfoList(ctx,pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}
// GetBtBrandSocialMediaDataSource 获取BtBrandSocialMedia的数据源
// @Tags BtBrandSocialMedia
// @Summary 获取BtBrandSocialMedia的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /btBrandSocialMedia/getBtBrandSocialMediaDataSource [get]
func (btBrandSocialMediaApi *BtBrandSocialMediaApi) GetBtBrandSocialMediaDataSource(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口为获取数据源定义的数据
    dataSource, err := btBrandSocialMediaService.GetBtBrandSocialMediaDataSource(ctx)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
   		response.FailWithMessage("查询失败:" + err.Error(), c)
   		return
    }
   response.OkWithData(dataSource, c)
}

// GetBtBrandSocialMediaPublic 不需要鉴权的品牌社交媒体接口
// @Tags BtBrandSocialMedia
// @Summary 不需要鉴权的品牌社交媒体接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btBrandSocialMedia/getBtBrandSocialMediaPublic [get]
func (btBrandSocialMediaApi *BtBrandSocialMediaApi) GetBtBrandSocialMediaPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    btBrandSocialMediaService.GetBtBrandSocialMediaPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的品牌社交媒体接口信息",
    }, "获取成功", c)
}
