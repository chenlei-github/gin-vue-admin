package brandtrekin

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
    brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BtBrandApi struct {}



// CreateBtBrand 创建品牌管理
// @Tags BtBrand
// @Summary 创建品牌管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtBrand true "创建品牌管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /btBrand/createBtBrand [post]
func (btBrandApi *BtBrandApi) CreateBtBrand(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var btBrand brandtrekin.BtBrand
	err := c.ShouldBindJSON(&btBrand)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btBrandService.CreateBtBrand(ctx,&btBrand)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBtBrand 删除品牌管理
// @Tags BtBrand
// @Summary 删除品牌管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtBrand true "删除品牌管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /btBrand/deleteBtBrand [delete]
func (btBrandApi *BtBrandApi) DeleteBtBrand(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := btBrandService.DeleteBtBrand(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBtBrandByIds 批量删除品牌管理
// @Tags BtBrand
// @Summary 批量删除品牌管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /btBrand/deleteBtBrandByIds [delete]
func (btBrandApi *BtBrandApi) DeleteBtBrandByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := btBrandService.DeleteBtBrandByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBtBrand 更新品牌管理
// @Tags BtBrand
// @Summary 更新品牌管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtBrand true "更新品牌管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /btBrand/updateBtBrand [put]
func (btBrandApi *BtBrandApi) UpdateBtBrand(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var btBrand brandtrekin.BtBrand
	err := c.ShouldBindJSON(&btBrand)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btBrandService.UpdateBtBrand(ctx,btBrand)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBtBrand 用id查询品牌管理
// @Tags BtBrand
// @Summary 用id查询品牌管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询品牌管理"
// @Success 200 {object} response.Response{data=brandtrekin.BtBrand,msg=string} "查询成功"
// @Router /btBrand/findBtBrand [get]
func (btBrandApi *BtBrandApi) FindBtBrand(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebtBrand, err := btBrandService.GetBtBrand(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebtBrand, c)
}
// GetBtBrandList 分页获取品牌管理列表
// @Tags BtBrand
// @Summary 分页获取品牌管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtBrandSearch true "分页获取品牌管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /btBrand/getBtBrandList [get]
func (btBrandApi *BtBrandApi) GetBtBrandList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo brandtrekinReq.BtBrandSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := btBrandService.GetBtBrandInfoList(ctx,pageInfo)
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
// GetBtBrandDataSource 获取BtBrand的数据源
// @Tags BtBrand
// @Summary 获取BtBrand的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /btBrand/getBtBrandDataSource [get]
func (btBrandApi *BtBrandApi) GetBtBrandDataSource(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口为获取数据源定义的数据
    dataSource, err := btBrandService.GetBtBrandDataSource(ctx)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
   		response.FailWithMessage("查询失败:" + err.Error(), c)
   		return
    }
   response.OkWithData(dataSource, c)
}

// GetBtBrandPublic 不需要鉴权的品牌管理接口
// @Tags BtBrand
// @Summary 不需要鉴权的品牌管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btBrand/getBtBrandPublic [get]
func (btBrandApi *BtBrandApi) GetBtBrandPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    btBrandService.GetBtBrandPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的品牌管理接口信息",
    }, "获取成功", c)
}
