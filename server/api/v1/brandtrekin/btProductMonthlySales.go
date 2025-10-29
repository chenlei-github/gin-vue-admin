package brandtrekin

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BtProductMonthlySalesApi struct {}



// CreateBtProductMonthlySales 创建商品月度销售
// @Tags BtProductMonthlySales
// @Summary 创建商品月度销售
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtProductMonthlySales true "创建商品月度销售"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /btProductMonthlySales/createBtProductMonthlySales [post]
func (btProductMonthlySalesApi *BtProductMonthlySalesApi) CreateBtProductMonthlySales(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var btProductMonthlySales brandtrekin.BtProductMonthlySales
	err := c.ShouldBindJSON(&btProductMonthlySales)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btProductMonthlySalesService.CreateBtProductMonthlySales(ctx,&btProductMonthlySales)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBtProductMonthlySales 删除商品月度销售
// @Tags BtProductMonthlySales
// @Summary 删除商品月度销售
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtProductMonthlySales true "删除商品月度销售"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /btProductMonthlySales/deleteBtProductMonthlySales [delete]
func (btProductMonthlySalesApi *BtProductMonthlySalesApi) DeleteBtProductMonthlySales(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := btProductMonthlySalesService.DeleteBtProductMonthlySales(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBtProductMonthlySalesByIds 批量删除商品月度销售
// @Tags BtProductMonthlySales
// @Summary 批量删除商品月度销售
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /btProductMonthlySales/deleteBtProductMonthlySalesByIds [delete]
func (btProductMonthlySalesApi *BtProductMonthlySalesApi) DeleteBtProductMonthlySalesByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := btProductMonthlySalesService.DeleteBtProductMonthlySalesByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBtProductMonthlySales 更新商品月度销售
// @Tags BtProductMonthlySales
// @Summary 更新商品月度销售
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtProductMonthlySales true "更新商品月度销售"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /btProductMonthlySales/updateBtProductMonthlySales [put]
func (btProductMonthlySalesApi *BtProductMonthlySalesApi) UpdateBtProductMonthlySales(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var btProductMonthlySales brandtrekin.BtProductMonthlySales
	err := c.ShouldBindJSON(&btProductMonthlySales)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btProductMonthlySalesService.UpdateBtProductMonthlySales(ctx,btProductMonthlySales)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBtProductMonthlySales 用id查询商品月度销售
// @Tags BtProductMonthlySales
// @Summary 用id查询商品月度销售
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询商品月度销售"
// @Success 200 {object} response.Response{data=brandtrekin.BtProductMonthlySales,msg=string} "查询成功"
// @Router /btProductMonthlySales/findBtProductMonthlySales [get]
func (btProductMonthlySalesApi *BtProductMonthlySalesApi) FindBtProductMonthlySales(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebtProductMonthlySales, err := btProductMonthlySalesService.GetBtProductMonthlySales(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebtProductMonthlySales, c)
}
// GetBtProductMonthlySalesList 分页获取商品月度销售列表
// @Tags BtProductMonthlySales
// @Summary 分页获取商品月度销售列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtProductMonthlySalesSearch true "分页获取商品月度销售列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /btProductMonthlySales/getBtProductMonthlySalesList [get]
func (btProductMonthlySalesApi *BtProductMonthlySalesApi) GetBtProductMonthlySalesList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo brandtrekinReq.BtProductMonthlySalesSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := btProductMonthlySalesService.GetBtProductMonthlySalesInfoList(ctx,pageInfo)
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

// GetBtProductMonthlySalesPublic 不需要鉴权的商品月度销售接口
// @Tags BtProductMonthlySales
// @Summary 不需要鉴权的商品月度销售接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btProductMonthlySales/getBtProductMonthlySalesPublic [get]
func (btProductMonthlySalesApi *BtProductMonthlySalesApi) GetBtProductMonthlySalesPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    btProductMonthlySalesService.GetBtProductMonthlySalesPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的商品月度销售接口信息",
    }, "获取成功", c)
}
