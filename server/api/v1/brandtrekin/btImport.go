package brandtrekin

import (
	"mime/multipart"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/brandtrekin"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BtImportApi struct{}

// PreviewBrandSocial 预览品牌社交媒体文件
// @Tags BtImport
// @Summary 预览品牌社交媒体Excel文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "品牌社交媒体文件 (Brand-Social.xlsx)"
// @Success 200 {object} response.Response{data=object,msg=string} "预览成功"
// @Router /btImport/previewBrandSocial [post]
func (api *BtImportApi) PreviewBrandSocial(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("文件上传失败: "+err.Error(), c)
		return
	}

	result, err := btImportService.ParseBrandSocial(file)
	if err != nil {
		global.GVA_LOG.Error("解析文件失败!", zap.Error(err))
		response.FailWithMessage("解析文件失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "预览成功", c)
}

// PreviewGKW 预览Google关键词文件
// @Tags BtImport
// @Summary 预览Google关键词CSV文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "Google关键词文件 (GKW.csv)"
// @Success 200 {object} response.Response{data=object,msg=string} "预览成功"
// @Router /btImport/previewGKW [post]
func (api *BtImportApi) PreviewGKW(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("文件上传失败: "+err.Error(), c)
		return
	}

	result, err := btImportService.ParseGKW(file)
	if err != nil {
		global.GVA_LOG.Error("解析文件失败!", zap.Error(err))
		response.FailWithMessage("解析文件失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "预览成功", c)
}

// PreviewKeywordHistory 预览Amazon关键词历史文件
// @Tags BtImport
// @Summary 预览Amazon关键词历史Excel文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "关键词历史文件 (KeywordHistory.xlsx)"
// @Success 200 {object} response.Response{data=object,msg=string} "预览成功"
// @Router /btImport/previewKeywordHistory [post]
func (api *BtImportApi) PreviewKeywordHistory(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("文件上传失败: "+err.Error(), c)
		return
	}

	result, err := btImportService.ParseKeywordHistory(file)
	if err != nil {
		global.GVA_LOG.Error("解析文件失败!", zap.Error(err))
		response.FailWithMessage("解析文件失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "预览成功", c)
}

// PreviewProductUS 预览美国商品文件
// @Tags BtImport
// @Summary 预览美国商品Excel文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "商品文件 (Product-US.xlsx)"
// @Success 200 {object} response.Response{data=object,msg=string} "预览成功"
// @Router /btImport/previewProductUS [post]
func (api *BtImportApi) PreviewProductUS(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("文件上传失败: "+err.Error(), c)
		return
	}

	result, err := btImportService.ParseProductUS(file)
	if err != nil {
		global.GVA_LOG.Error("解析文件失败!", zap.Error(err))
		response.FailWithMessage("解析文件失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "预览成功", c)
}

// PreviewProductSales 预览商品销售数据文件
// @Tags BtImport
// @Summary 预览商品销售数据Excel文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "商品销售文件 (product-US-sales.xlsx)"
// @Success 200 {object} response.Response{data=object,msg=string} "预览成功"
// @Router /btImport/previewProductSales [post]
func (api *BtImportApi) PreviewProductSales(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("文件上传失败: "+err.Error(), c)
		return
	}

	result, err := btImportService.ParseProductSales(file)
	if err != nil {
		global.GVA_LOG.Error("解析文件失败!", zap.Error(err))
		response.FailWithMessage("解析文件失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "预览成功", c)
}

// BatchImport 批量导入所有数据
// @Tags BtImport
// @Summary 批量导入市场的所有数据文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param marketId formData int true "市场ID"
// @Param replaceMode formData bool false "是否全量替换模式（true=替换现有数据, false=增量导入）"
// @Param brandSocial formData file false "品牌社交媒体文件"
// @Param productUS formData file false "商品文件"
// @Param gkw formData file false "Google关键词文件"
// @Param keywordHistory formData file false "Amazon关键词历史文件"
// @Param productSales formData file false "商品销售文件"
// @Success 200 {object} response.Response{msg=string} "导入成功"
// @Router /btImport/batchImport [post]
func (api *BtImportApi) BatchImport(c *gin.Context) {
	ctx := c.Request.Context()

	// 获取market_id
	marketIDStr := c.PostForm("marketId")
	if marketIDStr == "" {
		response.FailWithMessage("缺少marketId参数", c)
		return
	}

	marketID, err := strconv.ParseInt(marketIDStr, 10, 64)
	if err != nil {
		response.FailWithMessage("marketId格式错误", c)
		return
	}

	// 获取导入模式
	replaceModeStr := c.DefaultPostForm("replaceMode", "false")
	replaceMode := replaceModeStr == "true"

	// 收集所有上传的文件
	fileKeys := []string{"brandSocial", "productUS", "gkw", "keywordHistory", "productSales"}

	form, err := c.MultipartForm()
	if err != nil {
		response.FailWithMessage("获取表单数据失败: "+err.Error(), c)
		return
	}

	filesMap := make(map[string]*multipart.FileHeader)
	for _, key := range fileKeys {
		if fileHeaders, ok := form.File[key]; ok && len(fileHeaders) > 0 {
			filesMap[key] = fileHeaders[0]
		}
	}

	if len(filesMap) == 0 {
		response.FailWithMessage("至少需要上传一个文件", c)
		return
	}

	// 记录导入开始
	importLog := brandtrekin.BtImportLog{
		MarketId: &marketID,
		Status:   ptrString("processing"),
	}
	err = btImportLogService.CreateBtImportLog(ctx, &importLog)
	if err != nil {
		global.GVA_LOG.Error("创建导入日志失败!", zap.Error(err))
	}

	// 执行批量导入（使用优化后的方法）
	err = btImportService.BatchImport(marketID, filesMap, replaceMode)
	if err != nil {
		global.GVA_LOG.Error("批量导入失败!", zap.Error(err))

		// 更新导入日志状态为失败
		if importLog.ID != 0 {
			status := "failed"
			errorMsg := err.Error()
			importLog.Status = &status
			importLog.ErrorMessage = &errorMsg
			_ = btImportLogService.UpdateBtImportLog(ctx, importLog)
		}

		response.FailWithMessage("导入失败: "+err.Error(), c)
		return
	}

	// 更新导入日志状态为成功
	if importLog.ID != 0 {
		status := "success"
		importLog.Status = &status
		_ = btImportLogService.UpdateBtImportLog(ctx, importLog)
	}

	response.OkWithMessage("导入成功", c)
}

// Helper function to create string pointer
func ptrString(s string) *string {
	return &s
}
