package brandtrekin

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BtImportLogApi struct {}



// CreateBtImportLog 创建数据导入日志
// @Tags BtImportLog
// @Summary 创建数据导入日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtImportLog true "创建数据导入日志"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /btImportLog/createBtImportLog [post]
func (btImportLogApi *BtImportLogApi) CreateBtImportLog(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var btImportLog brandtrekin.BtImportLog
	err := c.ShouldBindJSON(&btImportLog)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btImportLogService.CreateBtImportLog(ctx,&btImportLog)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBtImportLog 删除数据导入日志
// @Tags BtImportLog
// @Summary 删除数据导入日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtImportLog true "删除数据导入日志"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /btImportLog/deleteBtImportLog [delete]
func (btImportLogApi *BtImportLogApi) DeleteBtImportLog(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := btImportLogService.DeleteBtImportLog(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBtImportLogByIds 批量删除数据导入日志
// @Tags BtImportLog
// @Summary 批量删除数据导入日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /btImportLog/deleteBtImportLogByIds [delete]
func (btImportLogApi *BtImportLogApi) DeleteBtImportLogByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := btImportLogService.DeleteBtImportLogByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBtImportLog 更新数据导入日志
// @Tags BtImportLog
// @Summary 更新数据导入日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtImportLog true "更新数据导入日志"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /btImportLog/updateBtImportLog [put]
func (btImportLogApi *BtImportLogApi) UpdateBtImportLog(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var btImportLog brandtrekin.BtImportLog
	err := c.ShouldBindJSON(&btImportLog)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btImportLogService.UpdateBtImportLog(ctx,btImportLog)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBtImportLog 用id查询数据导入日志
// @Tags BtImportLog
// @Summary 用id查询数据导入日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询数据导入日志"
// @Success 200 {object} response.Response{data=brandtrekin.BtImportLog,msg=string} "查询成功"
// @Router /btImportLog/findBtImportLog [get]
func (btImportLogApi *BtImportLogApi) FindBtImportLog(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebtImportLog, err := btImportLogService.GetBtImportLog(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebtImportLog, c)
}
// GetBtImportLogList 分页获取数据导入日志列表
// @Tags BtImportLog
// @Summary 分页获取数据导入日志列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtImportLogSearch true "分页获取数据导入日志列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /btImportLog/getBtImportLogList [get]
func (btImportLogApi *BtImportLogApi) GetBtImportLogList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo brandtrekinReq.BtImportLogSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := btImportLogService.GetBtImportLogInfoList(ctx,pageInfo)
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
// GetBtImportLogDataSource 获取BtImportLog的数据源
// @Tags BtImportLog
// @Summary 获取BtImportLog的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /btImportLog/getBtImportLogDataSource [get]
func (btImportLogApi *BtImportLogApi) GetBtImportLogDataSource(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口为获取数据源定义的数据
    dataSource, err := btImportLogService.GetBtImportLogDataSource(ctx)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
   		response.FailWithMessage("查询失败:" + err.Error(), c)
   		return
    }
   response.OkWithData(dataSource, c)
}

// GetBtImportLogPublic 不需要鉴权的数据导入日志接口
// @Tags BtImportLog
// @Summary 不需要鉴权的数据导入日志接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btImportLog/getBtImportLogPublic [get]
func (btImportLogApi *BtImportLogApi) GetBtImportLogPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    btImportLogService.GetBtImportLogPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的数据导入日志接口信息",
    }, "获取成功", c)
}
