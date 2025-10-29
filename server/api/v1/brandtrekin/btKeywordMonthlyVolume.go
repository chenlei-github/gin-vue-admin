package brandtrekin

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BtKeywordMonthlyVolumeApi struct {}



// CreateBtKeywordMonthlyVolume 创建关键词月度搜索量
// @Tags BtKeywordMonthlyVolume
// @Summary 创建关键词月度搜索量
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtKeywordMonthlyVolume true "创建关键词月度搜索量"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /btKeywordMonthlyVolume/createBtKeywordMonthlyVolume [post]
func (btKeywordMonthlyVolumeApi *BtKeywordMonthlyVolumeApi) CreateBtKeywordMonthlyVolume(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var btKeywordMonthlyVolume brandtrekin.BtKeywordMonthlyVolume
	err := c.ShouldBindJSON(&btKeywordMonthlyVolume)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btKeywordMonthlyVolumeService.CreateBtKeywordMonthlyVolume(ctx,&btKeywordMonthlyVolume)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBtKeywordMonthlyVolume 删除关键词月度搜索量
// @Tags BtKeywordMonthlyVolume
// @Summary 删除关键词月度搜索量
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtKeywordMonthlyVolume true "删除关键词月度搜索量"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /btKeywordMonthlyVolume/deleteBtKeywordMonthlyVolume [delete]
func (btKeywordMonthlyVolumeApi *BtKeywordMonthlyVolumeApi) DeleteBtKeywordMonthlyVolume(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := btKeywordMonthlyVolumeService.DeleteBtKeywordMonthlyVolume(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBtKeywordMonthlyVolumeByIds 批量删除关键词月度搜索量
// @Tags BtKeywordMonthlyVolume
// @Summary 批量删除关键词月度搜索量
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /btKeywordMonthlyVolume/deleteBtKeywordMonthlyVolumeByIds [delete]
func (btKeywordMonthlyVolumeApi *BtKeywordMonthlyVolumeApi) DeleteBtKeywordMonthlyVolumeByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := btKeywordMonthlyVolumeService.DeleteBtKeywordMonthlyVolumeByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBtKeywordMonthlyVolume 更新关键词月度搜索量
// @Tags BtKeywordMonthlyVolume
// @Summary 更新关键词月度搜索量
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtKeywordMonthlyVolume true "更新关键词月度搜索量"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /btKeywordMonthlyVolume/updateBtKeywordMonthlyVolume [put]
func (btKeywordMonthlyVolumeApi *BtKeywordMonthlyVolumeApi) UpdateBtKeywordMonthlyVolume(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var btKeywordMonthlyVolume brandtrekin.BtKeywordMonthlyVolume
	err := c.ShouldBindJSON(&btKeywordMonthlyVolume)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btKeywordMonthlyVolumeService.UpdateBtKeywordMonthlyVolume(ctx,btKeywordMonthlyVolume)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBtKeywordMonthlyVolume 用id查询关键词月度搜索量
// @Tags BtKeywordMonthlyVolume
// @Summary 用id查询关键词月度搜索量
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询关键词月度搜索量"
// @Success 200 {object} response.Response{data=brandtrekin.BtKeywordMonthlyVolume,msg=string} "查询成功"
// @Router /btKeywordMonthlyVolume/findBtKeywordMonthlyVolume [get]
func (btKeywordMonthlyVolumeApi *BtKeywordMonthlyVolumeApi) FindBtKeywordMonthlyVolume(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebtKeywordMonthlyVolume, err := btKeywordMonthlyVolumeService.GetBtKeywordMonthlyVolume(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebtKeywordMonthlyVolume, c)
}
// GetBtKeywordMonthlyVolumeList 分页获取关键词月度搜索量列表
// @Tags BtKeywordMonthlyVolume
// @Summary 分页获取关键词月度搜索量列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtKeywordMonthlyVolumeSearch true "分页获取关键词月度搜索量列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /btKeywordMonthlyVolume/getBtKeywordMonthlyVolumeList [get]
func (btKeywordMonthlyVolumeApi *BtKeywordMonthlyVolumeApi) GetBtKeywordMonthlyVolumeList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo brandtrekinReq.BtKeywordMonthlyVolumeSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := btKeywordMonthlyVolumeService.GetBtKeywordMonthlyVolumeInfoList(ctx,pageInfo)
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
// GetBtKeywordMonthlyVolumeDataSource 获取BtKeywordMonthlyVolume的数据源
// @Tags BtKeywordMonthlyVolume
// @Summary 获取BtKeywordMonthlyVolume的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /btKeywordMonthlyVolume/getBtKeywordMonthlyVolumeDataSource [get]
func (btKeywordMonthlyVolumeApi *BtKeywordMonthlyVolumeApi) GetBtKeywordMonthlyVolumeDataSource(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口为获取数据源定义的数据
    dataSource, err := btKeywordMonthlyVolumeService.GetBtKeywordMonthlyVolumeDataSource(ctx)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
   		response.FailWithMessage("查询失败:" + err.Error(), c)
   		return
    }
   response.OkWithData(dataSource, c)
}

// GetBtKeywordMonthlyVolumePublic 不需要鉴权的关键词月度搜索量接口
// @Tags BtKeywordMonthlyVolume
// @Summary 不需要鉴权的关键词月度搜索量接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btKeywordMonthlyVolume/getBtKeywordMonthlyVolumePublic [get]
func (btKeywordMonthlyVolumeApi *BtKeywordMonthlyVolumeApi) GetBtKeywordMonthlyVolumePublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    btKeywordMonthlyVolumeService.GetBtKeywordMonthlyVolumePublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的关键词月度搜索量接口信息",
    }, "获取成功", c)
}
