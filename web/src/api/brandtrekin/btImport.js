import service from '@/utils/request'

// @Tags BtImport
// @Summary 预览品牌社交媒体Excel文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "品牌社交媒体文件 (Brand-Social.xlsx)"
// @Success 200 {object} response.Response{data=object,msg=string} "预览成功"
// @Router /btImport/previewBrandSocial [post]
export const previewBrandSocial = (formData) => {
  return service({
    url: '/btImport/previewBrandSocial',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// @Tags BtImport
// @Summary 预览Google关键词CSV文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "Google关键词文件 (GKW.csv)"
// @Success 200 {object} response.Response{data=object,msg=string} "预览成功"
// @Router /btImport/previewGKW [post]
export const previewGKW = (formData) => {
  return service({
    url: '/btImport/previewGKW',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// @Tags BtImport
// @Summary 预览Amazon关键词历史Excel文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "关键词历史文件 (KeywordHistory.xlsx)"
// @Success 200 {object} response.Response{data=object,msg=string} "预览成功"
// @Router /btImport/previewKeywordHistory [post]
export const previewKeywordHistory = (formData) => {
  return service({
    url: '/btImport/previewKeywordHistory',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// @Tags BtImport
// @Summary 预览美国商品Excel文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "商品文件 (Product-US.xlsx)"
// @Success 200 {object} response.Response{data=object,msg=string} "预览成功"
// @Router /btImport/previewProductUS [post]
export const previewProductUS = (formData) => {
  return service({
    url: '/btImport/previewProductUS',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// @Tags BtImport
// @Summary 预览商品销售数据Excel文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "商品销售文件 (product-US-sales.xlsx)"
// @Success 200 {object} response.Response{data=object,msg=string} "预览成功"
// @Router /btImport/previewProductSales [post]
export const previewProductSales = (formData) => {
  return service({
    url: '/btImport/previewProductSales',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

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
export const batchImport = (formData) => {
  return service({
    url: '/btImport/batchImport',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    timeout: 300000 // 5分钟超时
  })
}

