package brandtrekin

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BtKeywordApi struct {}



// CreateBtKeyword 创建关键词管理
// @Tags BtKeyword
// @Summary 创建关键词管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtKeyword true "创建关键词管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /btKeyword/createBtKeyword [post]
func (btKeywordApi *BtKeywordApi) CreateBtKeyword(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var btKeyword brandtrekin.BtKeyword
	err := c.ShouldBindJSON(&btKeyword)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btKeywordService.CreateBtKeyword(ctx,&btKeyword)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBtKeyword 删除关键词管理
// @Tags BtKeyword
// @Summary 删除关键词管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtKeyword true "删除关键词管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /btKeyword/deleteBtKeyword [delete]
func (btKeywordApi *BtKeywordApi) DeleteBtKeyword(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := btKeywordService.DeleteBtKeyword(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBtKeywordByIds 批量删除关键词管理
// @Tags BtKeyword
// @Summary 批量删除关键词管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /btKeyword/deleteBtKeywordByIds [delete]
func (btKeywordApi *BtKeywordApi) DeleteBtKeywordByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := btKeywordService.DeleteBtKeywordByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBtKeyword 更新关键词管理
// @Tags BtKeyword
// @Summary 更新关键词管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtKeyword true "更新关键词管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /btKeyword/updateBtKeyword [put]
func (btKeywordApi *BtKeywordApi) UpdateBtKeyword(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var btKeyword brandtrekin.BtKeyword
	err := c.ShouldBindJSON(&btKeyword)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btKeywordService.UpdateBtKeyword(ctx,btKeyword)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBtKeyword 用id查询关键词管理
// @Tags BtKeyword
// @Summary 用id查询关键词管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询关键词管理"
// @Success 200 {object} response.Response{data=brandtrekin.BtKeyword,msg=string} "查询成功"
// @Router /btKeyword/findBtKeyword [get]
func (btKeywordApi *BtKeywordApi) FindBtKeyword(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebtKeyword, err := btKeywordService.GetBtKeyword(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebtKeyword, c)
}
// GetBtKeywordList 分页获取关键词管理列表
// @Tags BtKeyword
// @Summary 分页获取关键词管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtKeywordSearch true "分页获取关键词管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /btKeyword/getBtKeywordList [get]
func (btKeywordApi *BtKeywordApi) GetBtKeywordList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo brandtrekinReq.BtKeywordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := btKeywordService.GetBtKeywordInfoList(ctx,pageInfo)
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
// GetBtKeywordDataSource 获取BtKeyword的数据源
// @Tags BtKeyword
// @Summary 获取BtKeyword的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /btKeyword/getBtKeywordDataSource [get]
func (btKeywordApi *BtKeywordApi) GetBtKeywordDataSource(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口为获取数据源定义的数据
    dataSource, err := btKeywordService.GetBtKeywordDataSource(ctx)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
   		response.FailWithMessage("查询失败:" + err.Error(), c)
   		return
    }
   response.OkWithData(dataSource, c)
}

// GetBtKeywordPublic 不需要鉴权的关键词管理接口
// @Tags BtKeyword
// @Summary 不需要鉴权的关键词管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btKeyword/getBtKeywordPublic [get]
func (btKeywordApi *BtKeywordApi) GetBtKeywordPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    btKeywordService.GetBtKeywordPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的关键词管理接口信息",
    }, "获取成功", c)
}
