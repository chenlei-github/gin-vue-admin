package brandtrekin

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
	brandtrekinReq "github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BtMarketApi struct {}



// CreateBtMarket 创建市场管理
// @Tags BtMarket
// @Summary 创建市场管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtMarket true "创建市场管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /btMarket/createBtMarket [post]
func (btMarketApi *BtMarketApi) CreateBtMarket(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var btMarket brandtrekin.BtMarket
	err := c.ShouldBindJSON(&btMarket)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btMarketService.CreateBtMarket(ctx,&btMarket)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBtMarket 删除市场管理
// @Tags BtMarket
// @Summary 删除市场管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtMarket true "删除市场管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /btMarket/deleteBtMarket [delete]
func (btMarketApi *BtMarketApi) DeleteBtMarket(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := btMarketService.DeleteBtMarket(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBtMarketByIds 批量删除市场管理
// @Tags BtMarket
// @Summary 批量删除市场管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /btMarket/deleteBtMarketByIds [delete]
func (btMarketApi *BtMarketApi) DeleteBtMarketByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := btMarketService.DeleteBtMarketByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBtMarket 更新市场管理
// @Tags BtMarket
// @Summary 更新市场管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body brandtrekin.BtMarket true "更新市场管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /btMarket/updateBtMarket [put]
func (btMarketApi *BtMarketApi) UpdateBtMarket(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var btMarket brandtrekin.BtMarket
	err := c.ShouldBindJSON(&btMarket)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = btMarketService.UpdateBtMarket(ctx,btMarket)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBtMarket 用id查询市场管理
// @Tags BtMarket
// @Summary 用id查询市场管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询市场管理"
// @Success 200 {object} response.Response{data=brandtrekin.BtMarket,msg=string} "查询成功"
// @Router /btMarket/findBtMarket [get]
func (btMarketApi *BtMarketApi) FindBtMarket(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebtMarket, err := btMarketService.GetBtMarket(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebtMarket, c)
}
// GetBtMarketList 分页获取市场管理列表
// @Tags BtMarket
// @Summary 分页获取市场管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query brandtrekinReq.BtMarketSearch true "分页获取市场管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /btMarket/getBtMarketList [get]
func (btMarketApi *BtMarketApi) GetBtMarketList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo brandtrekinReq.BtMarketSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := btMarketService.GetBtMarketInfoList(ctx,pageInfo)
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

// ToggleMarketStatus 切换市场状态
// @Tags BtMarket
// @Summary 切换市场状态（用于状态开关实时切换）
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "市场ID"
// @Param status body object{status=string} true "状态 (active/inactive)"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /btMarket/toggleMarketStatus [put]
func (btMarketApi *BtMarketApi) ToggleMarketStatus(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("市场ID不能为空", c)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证状态值
	if req.Status != "active" && req.Status != "inactive" {
		response.FailWithMessage("状态值必须为active或inactive", c)
		return
	}

	// 获取市场信息
	btMarket, err := btMarketService.GetBtMarket(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}

	// 更新状态
	status := req.Status
	btMarket.Status = &status
	err = btMarketService.UpdateBtMarket(ctx, btMarket)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("状态更新成功", c)
}

// GenerateSlugFromName 根据市场名称自动生成slug
// @Tags BtMarket
// @Summary 根据市场名称自动生成slug
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param name query string true "市场名称"
// @Success 200 {object} response.Response{data=object{slug=string},msg=string} "生成成功"
// @Router /btMarket/generateSlugFromName [get]
func (btMarketApi *BtMarketApi) GenerateSlugFromName(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	name := c.Query("name")
	if name == "" {
		response.FailWithMessage("市场名称不能为空", c)
		return
	}

	slug, err := btMarketService.GenerateSlugFromName(ctx, name)
	if err != nil {
		global.GVA_LOG.Error("生成失败!", zap.Error(err))
		response.FailWithMessage("生成失败:"+err.Error(), c)
		return
	}

	response.OkWithData(gin.H{"slug": slug}, c)
}

// ValidateSlugUnique 校验slug唯一性
// @Tags BtMarket
// @Summary 校验slug唯一性
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param slug query string true "市场slug"
// @Param excludeID query uint false "排除的市场ID（用于编辑时检查）"
// @Success 200 {object} response.Response{data=object{isUnique=bool},msg=string} "校验成功"
// @Router /btMarket/validateSlugUnique [get]
func (btMarketApi *BtMarketApi) ValidateSlugUnique(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	slug := c.Query("slug")
	if slug == "" {
		response.FailWithMessage("市场ID不能为空", c)
		return
	}

	excludeIDStr := c.Query("excludeID")
	var excludeID uint
	if excludeIDStr != "" {
		if id, err := strconv.ParseUint(excludeIDStr, 10, 32); err == nil {
			excludeID = uint(id)
		}
	}

	isUnique, err := btMarketService.ValidateSlugUnique(ctx, slug, excludeID)
	if err != nil {
		global.GVA_LOG.Error("校验失败!", zap.Error(err))
		response.FailWithMessage("校验失败:"+err.Error(), c)
		return
	}

	response.OkWithData(gin.H{"isUnique": isUnique}, c)
}

// ValidateDeleteMarket 删除前市场名称校验
// @Tags BtMarket
// @Summary 删除前市场名称校验（返回市场信息用于前端确认）
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "市场ID"
// @Success 200 {object} response.Response{data=brandtrekin.BtMarket,msg=string} "查询成功"
// @Router /btMarket/validateDeleteMarket [get]
func (btMarketApi *BtMarketApi) ValidateDeleteMarket(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("市场ID不能为空", c)
		return
	}

	btMarket, err := btMarketService.GetBtMarket(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(btMarket, c)
}

// GetBtMarketPublic 不需要鉴权的市场管理接口
// @Tags BtMarket
// @Summary 不需要鉴权的市场管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /btMarket/getBtMarketPublic [get]
func (btMarketApi *BtMarketApi) GetBtMarketPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	btMarketService.GetBtMarketPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的市场管理接口信息",
	}, "获取成功", c)
}
