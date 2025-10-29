package brandtrekin

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BtMarketMonthlyTrendApi struct {}



// CreateBtMarketMonthlyTrend 创建市场月度趋势
// @Tags BtMarketMonthlyTrend
// @Summary 创建市场月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtMarketMonthlyTrend true "创建市场月度趋势"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /btMarketMonthlyTrend/createBtMarketMonthlyTrend [post]
func (btMarketMonthlyTrendApi *BtMarketMonthlyTrendApi) CreateBtMarketMonthlyTrend(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var btMarketMonthlyTrend brandtrekin.BtMarketMonthlyTrend
	err := c.ShouldBindJSON(&btMarketMonthlyTrend)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btMarketMonthlyTrendService.CreateBtMarketMonthlyTrend(ctx,&btMarketMonthlyTrend)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBtMarketMonthlyTrend 删除市场月度趋势
// @Tags BtMarketMonthlyTrend
// @Summary 删除市场月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtMarketMonthlyTrend true "删除市场月度趋势"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /btMarketMonthlyTrend/deleteBtMarketMonthlyTrend [delete]
func (btMarketMonthlyTrendApi *BtMarketMonthlyTrendApi) DeleteBtMarketMonthlyTrend(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := btMarketMonthlyTrendService.DeleteBtMarketMonthlyTrend(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBtMarketMonthlyTrendByIds 批量删除市场月度趋势
// @Tags BtMarketMonthlyTrend
// @Summary 批量删除市场月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /btMarketMonthlyTrend/deleteBtMarketMonthlyTrendByIds [delete]
func (btMarketMonthlyTrendApi *BtMarketMonthlyTrendApi) DeleteBtMarketMonthlyTrendByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := btMarketMonthlyTrendService.DeleteBtMarketMonthlyTrendByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBtMarketMonthlyTrend 更新市场月度趋势
// @Tags BtMarketMonthlyTrend
// @Summary 更新市场月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtMarketMonthlyTrend true "更新市场月度趋势"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /btMarketMonthlyTrend/updateBtMarketMonthlyTrend [put]
func (btMarketMonthlyTrendApi *BtMarketMonthlyTrendApi) UpdateBtMarketMonthlyTrend(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var btMarketMonthlyTrend brandtrekin.BtMarketMonthlyTrend
	err := c.ShouldBindJSON(&btMarketMonthlyTrend)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btMarketMonthlyTrendService.UpdateBtMarketMonthlyTrend(ctx,btMarketMonthlyTrend)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBtMarketMonthlyTrend 用id查询市场月度趋势
// @Tags BtMarketMonthlyTrend
// @Summary 用id查询市场月度趋势
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询市场月度趋势"
// @Success 200 {object} response.Response{data=brandtrekin.BtMarketMonthlyTrend,msg=string} "查询成功"
// @Router /btMarketMonthlyTrend/findBtMarketMonthlyTrend [get]
func (btMarketMonthlyTrendApi *BtMarketMonthlyTrendApi) FindBtMarketMonthlyTrend(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebtMarketMonthlyTrend, err := btMarketMonthlyTrendService.GetBtMarketMonthlyTrend(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebtMarketMonthlyTrend, c)
}
// GetBtMarketMonthlyTrendList 分页获取市场月度趋势列表
// @Tags BtMarketMonthlyTrend
// @Summary 分页获取市场月度趋势列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtMarketMonthlyTrendSearch true "分页获取市场月度趋势列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /btMarketMonthlyTrend/getBtMarketMonthlyTrendList [get]
func (btMarketMonthlyTrendApi *BtMarketMonthlyTrendApi) GetBtMarketMonthlyTrendList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo brandtrekinReq.BtMarketMonthlyTrendSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := btMarketMonthlyTrendService.GetBtMarketMonthlyTrendInfoList(ctx,pageInfo)
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
// GetBtMarketMonthlyTrendDataSource 获取BtMarketMonthlyTrend的数据源
// @Tags BtMarketMonthlyTrend
// @Summary 获取BtMarketMonthlyTrend的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /btMarketMonthlyTrend/getBtMarketMonthlyTrendDataSource [get]
func (btMarketMonthlyTrendApi *BtMarketMonthlyTrendApi) GetBtMarketMonthlyTrendDataSource(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口为获取数据源定义的数据
    dataSource, err := btMarketMonthlyTrendService.GetBtMarketMonthlyTrendDataSource(ctx)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
   		response.FailWithMessage("查询失败:" + err.Error(), c)
   		return
    }
   response.OkWithData(dataSource, c)
}

// GetBtMarketMonthlyTrendPublic 不需要鉴权的市场月度趋势接口
// @Tags BtMarketMonthlyTrend
// @Summary 不需要鉴权的市场月度趋势接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btMarketMonthlyTrend/getBtMarketMonthlyTrendPublic [get]
func (btMarketMonthlyTrendApi *BtMarketMonthlyTrendApi) GetBtMarketMonthlyTrendPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    btMarketMonthlyTrendService.GetBtMarketMonthlyTrendPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的市场月度趋势接口信息",
    }, "获取成功", c)
}
