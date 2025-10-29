package brandtrekin

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BtProductApi struct {}



// CreateBtProduct 创建商品管理
// @Tags BtProduct
// @Summary 创建商品管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtProduct true "创建商品管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /btProduct/createBtProduct [post]
func (btProductApi *BtProductApi) CreateBtProduct(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var btProduct brandtrekin.BtProduct
	err := c.ShouldBindJSON(&btProduct)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btProductService.CreateBtProduct(ctx,&btProduct)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBtProduct 删除商品管理
// @Tags BtProduct
// @Summary 删除商品管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtProduct true "删除商品管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /btProduct/deleteBtProduct [delete]
func (btProductApi *BtProductApi) DeleteBtProduct(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := btProductService.DeleteBtProduct(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBtProductByIds 批量删除商品管理
// @Tags BtProduct
// @Summary 批量删除商品管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /btProduct/deleteBtProductByIds [delete]
func (btProductApi *BtProductApi) DeleteBtProductByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := btProductService.DeleteBtProductByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBtProduct 更新商品管理
// @Tags BtProduct
// @Summary 更新商品管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtProduct true "更新商品管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /btProduct/updateBtProduct [put]
func (btProductApi *BtProductApi) UpdateBtProduct(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var btProduct brandtrekin.BtProduct
	err := c.ShouldBindJSON(&btProduct)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btProductService.UpdateBtProduct(ctx,btProduct)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBtProduct 用id查询商品管理
// @Tags BtProduct
// @Summary 用id查询商品管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询商品管理"
// @Success 200 {object} response.Response{data=brandtrekin.BtProduct,msg=string} "查询成功"
// @Router /btProduct/findBtProduct [get]
func (btProductApi *BtProductApi) FindBtProduct(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebtProduct, err := btProductService.GetBtProduct(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebtProduct, c)
}
// GetBtProductList 分页获取商品管理列表
// @Tags BtProduct
// @Summary 分页获取商品管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtProductSearch true "分页获取商品管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /btProduct/getBtProductList [get]
func (btProductApi *BtProductApi) GetBtProductList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo brandtrekinReq.BtProductSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := btProductService.GetBtProductInfoList(ctx,pageInfo)
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
// GetBtProductDataSource 获取BtProduct的数据源
// @Tags BtProduct
// @Summary 获取BtProduct的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /btProduct/getBtProductDataSource [get]
func (btProductApi *BtProductApi) GetBtProductDataSource(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口为获取数据源定义的数据
    dataSource, err := btProductService.GetBtProductDataSource(ctx)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
   		response.FailWithMessage("查询失败:" + err.Error(), c)
   		return
    }
   response.OkWithData(dataSource, c)
}

// GetBtProductPublic 不需要鉴权的商品管理接口
// @Tags BtProduct
// @Summary 不需要鉴权的商品管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btProduct/getBtProductPublic [get]
func (btProductApi *BtProductApi) GetBtProductPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    btProductService.GetBtProductPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的商品管理接口信息",
    }, "获取成功", c)
}
