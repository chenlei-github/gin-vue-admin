<template>
  <div class="market-import-container">
    <!-- 页面标题 -->
<!--    <div class="mb-6">-->
<!--      <h2 class="text-2xl font-bold">导入市场数据 - {{ marketInfo.marketName || '未知市场' }}</h2>-->
<!--    </div>-->

    <!-- 5个文件上传组件 -->
    <div class="space-y-4 mb-6">
      <!-- 上传组件1: 品牌社交媒体数据 -->
      <el-card shadow="hover">
        <template #header>
          <div class="flex items-center justify-between">
            <div>
              <span class="font-semibold">📄 Brand-Social.xlsx</span>
              <span class="ml-2 text-sm text-gray-500">品牌社交媒体数据</span>
            </div>
            <el-tag :type="getStatusType(uploadStatus.brandSocial)">
              {{ getStatusText(uploadStatus.brandSocial) }}
            </el-tag>
          </div>
        </template>
        <div class="text-sm text-gray-600 mb-4">
          包含品牌名称、独立站、YouTube、Instagram、Facebook、Reddit数据
        </div>
        <div>
          <input
            ref="brandSocialFileInputRef"
            type="file"
            accept=".xlsx"
            style="display: none"
            @change="(e) => handleFileSelect('brandSocial', e, '.xlsx')"
          />
          <el-button
            type="primary"
            :loading="uploadStatus.brandSocial === 'uploading' || uploadStatus.brandSocial === 'parsing'"
            @click="$refs.brandSocialFileInputRef?.click()"
          >
            {{ uploadFiles.brandSocial ? '重新选择' : '选择文件' }}
          </el-button>
        </div>
        <div v-if="uploadFiles.brandSocial" class="mt-2 text-sm text-green-600">
          {{ uploadFiles.brandSocial.name }}
        </div>
        <div v-if="uploadErrors.brandSocial" class="mt-2 text-sm text-red-600">
          {{ uploadErrors.brandSocial }}
        </div>
        <div v-if="uploadPreviews.brandSocial && uploadStatus.brandSocial === 'success'" class="mt-4">
          <el-collapse>
            <el-collapse-item title="数据预览（前5行）" name="preview">
              <el-table :data="uploadPreviews.brandSocial" border size="small" max-height="300">
                <el-table-column v-for="(col, idx) in getPreviewColumns(uploadPreviews.brandSocial)" :key="idx" :prop="col" :label="col" />
              </el-table>
            </el-collapse-item>
          </el-collapse>
        </div>
      </el-card>

      <!-- 上传组件2: Google关键词数据 -->
      <el-card shadow="hover">
        <template #header>
          <div class="flex items-center justify-between">
            <div>
              <span class="font-semibold">📄 GKW.csv</span>
              <span class="ml-2 text-sm text-gray-500">Google关键词数据</span>
            </div>
            <el-tag :type="getStatusType(uploadStatus.gkw)">
              {{ getStatusText(uploadStatus.gkw) }}
            </el-tag>
          </div>
        </template>
        <div class="text-sm text-gray-600 mb-4">
          包含Google关键词及月度搜索量历史数据
        </div>
        <div>
          <input
            ref="gkwFileInputRef"
            type="file"
            accept=".csv"
            style="display: none"
            @change="(e) => handleFileSelect('gkw', e, '.csv')"
          />
          <el-button
            type="primary"
            :loading="uploadStatus.gkw === 'uploading' || uploadStatus.gkw === 'parsing'"
            @click="$refs.gkwFileInputRef?.click()"
          >
            {{ uploadFiles.gkw ? '重新选择' : '选择文件' }}
          </el-button>
        </div>
        <div v-if="uploadFiles.gkw" class="mt-2 text-sm text-green-600">
          {{ uploadFiles.gkw.name }}
        </div>
        <div v-if="uploadErrors.gkw" class="mt-2 text-sm text-red-600">
          {{ uploadErrors.gkw }}
        </div>
        <div v-if="uploadPreviews.gkw && uploadStatus.gkw === 'success'" class="mt-4">
          <el-collapse>
            <el-collapse-item title="数据预览（前5行）" name="preview">
              <el-table :data="uploadPreviews.gkw" border size="small" max-height="300">
                <el-table-column v-for="(col, idx) in getPreviewColumns(uploadPreviews.gkw)" :key="idx" :prop="col" :label="col" />
              </el-table>
            </el-collapse-item>
          </el-collapse>
        </div>
      </el-card>

      <!-- 上传组件3: Amazon关键词历史数据 -->
      <el-card shadow="hover">
        <template #header>
          <div class="flex items-center justify-between">
            <div>
              <span class="font-semibold">📄 KeywordHistory.xlsx</span>
              <span class="ml-2 text-sm text-gray-500">Amazon关键词历史数据</span>
            </div>
            <el-tag :type="getStatusType(uploadStatus.keywordHistory)">
              {{ getStatusText(uploadStatus.keywordHistory) }}
            </el-tag>
          </div>
        </template>
        <div class="text-sm text-gray-600 mb-4">
          包含Amazon关键词及月度搜索量历史数据
        </div>
        <div>
          <input
            ref="keywordHistoryFileInputRef"
            type="file"
            accept=".xlsx"
            style="display: none"
            @change="(e) => handleFileSelect('keywordHistory', e, '.xlsx')"
          />
          <el-button
            type="primary"
            :loading="uploadStatus.keywordHistory === 'uploading' || uploadStatus.keywordHistory === 'parsing'"
            @click="$refs.keywordHistoryFileInputRef?.click()"
          >
            {{ uploadFiles.keywordHistory ? '重新选择' : '选择文件' }}
          </el-button>
        </div>
        <div v-if="uploadFiles.keywordHistory" class="mt-2 text-sm text-green-600">
          {{ uploadFiles.keywordHistory.name }}
        </div>
        <div v-if="uploadErrors.keywordHistory" class="mt-2 text-sm text-red-600">
          {{ uploadErrors.keywordHistory }}
        </div>
        <div v-if="uploadPreviews.keywordHistory && uploadStatus.keywordHistory === 'success'" class="mt-4">
          <el-collapse>
            <el-collapse-item title="数据预览（前5行）" name="preview">
              <el-table :data="uploadPreviews.keywordHistory" border size="small" max-height="300">
                <el-table-column v-for="(col, idx) in getPreviewColumns(uploadPreviews.keywordHistory)" :key="idx" :prop="col" :label="col" />
              </el-table>
            </el-collapse-item>
          </el-collapse>
        </div>
      </el-card>

      <!-- 上传组件4: 商品基础信息 -->
      <el-card shadow="hover">
        <template #header>
          <div class="flex items-center justify-between">
            <div>
              <span class="font-semibold">📄 Product-US.xlsx</span>
              <span class="ml-2 text-sm text-gray-500">商品基础信息</span>
            </div>
            <el-tag :type="getStatusType(uploadStatus.productUS)">
              {{ getStatusText(uploadStatus.productUS) }}
            </el-tag>
          </div>
        </template>
        <div class="text-sm text-gray-600 mb-4">
          包含ASIN、标题、品牌、价格、评分、评论数、图片URL等
        </div>
        <div>
          <input
            ref="productUSFileInputRef"
            type="file"
            accept=".xlsx"
            style="display: none"
            @change="(e) => handleFileSelect('productUS', e, '.xlsx')"
          />
          <el-button
            type="primary"
            :loading="uploadStatus.productUS === 'uploading' || uploadStatus.productUS === 'parsing'"
            @click="$refs.productUSFileInputRef?.click()"
          >
            {{ uploadFiles.productUS ? '重新选择' : '选择文件' }}
          </el-button>
        </div>
        <div v-if="uploadFiles.productUS" class="mt-2 text-sm text-green-600">
          {{ uploadFiles.productUS.name }}
        </div>
        <div v-if="uploadErrors.productUS" class="mt-2 text-sm text-red-600">
          {{ uploadErrors.productUS }}
        </div>
        <div v-if="uploadPreviews.productUS && uploadStatus.productUS === 'success'" class="mt-4">
          <el-collapse>
            <el-collapse-item title="数据预览（前5行）" name="preview">
              <el-table :data="uploadPreviews.productUS" border size="small" max-height="300">
                <el-table-column v-for="(col, idx) in getPreviewColumns(uploadPreviews.productUS)" :key="idx" :prop="col" :label="col" />
              </el-table>
            </el-collapse-item>
          </el-collapse>
        </div>
      </el-card>

      <!-- 上传组件5: 商品月度销售数据 -->
      <el-card shadow="hover">
        <template #header>
          <div class="flex items-center justify-between">
            <div>
              <span class="font-semibold">📄 product-US-sales.xlsx</span>
              <span class="ml-2 text-sm text-gray-500">商品月度销售数据</span>
            </div>
            <el-tag :type="getStatusType(uploadStatus.productSales)">
              {{ getStatusText(uploadStatus.productSales) }}
            </el-tag>
          </div>
        </template>
        <div class="text-sm text-gray-600 mb-4">
          包含ASIN及每月销售额、销量数据
        </div>
        <div>
          <input
            ref="productSalesFileInputRef"
            type="file"
            accept=".xlsx"
            style="display: none"
            @change="(e) => handleFileSelect('productSales', e, '.xlsx')"
          />
          <el-button
            type="primary"
            :loading="uploadStatus.productSales === 'uploading' || uploadStatus.productSales === 'parsing'"
            @click="$refs.productSalesFileInputRef?.click()"
          >
            {{ uploadFiles.productSales ? '重新选择' : '选择文件' }}
          </el-button>
        </div>
        <div v-if="uploadFiles.productSales" class="mt-2 text-sm text-green-600">
          {{ uploadFiles.productSales.name }}
        </div>
        <div v-if="uploadErrors.productSales" class="mt-2 text-sm text-red-600">
          {{ uploadErrors.productSales }}
        </div>
        <div v-if="uploadPreviews.productSales && uploadStatus.productSales === 'success'" class="mt-4">
          <el-collapse>
            <el-collapse-item title="数据预览（前5行）" name="preview">
              <el-table :data="uploadPreviews.productSales" border size="small" max-height="300">
                <el-table-column v-for="(col, idx) in getPreviewColumns(uploadPreviews.productSales)" :key="idx" :prop="col" :label="col" />
              </el-table>
            </el-collapse-item>
          </el-collapse>
        </div>
      </el-card>
    </div>

    <!-- 导入选项 -->
    <el-card shadow="hover" class="mb-6">
      <template #header>
        <span class="font-semibold">导入选项</span>
      </template>
      <div class="space-y-4">
        <div>
          <div class="mb-2 font-medium">导入模式：</div>
          <el-radio-group v-model="importOptions.importMode">
            <el-radio label="incremental">增量导入（保留现有数据，仅添加或更新）</el-radio>
            <el-radio label="replace">全量替换（删除所有数据后重新导入）</el-radio>
          </el-radio-group>
        </div>
        <div>
          <div class="mb-2 font-medium">数据校验：</div>
          <el-checkbox v-model="importOptions.skipInvalid">跳过无效数据行</el-checkbox>
          <el-checkbox v-model="importOptions.autoCreateBrand">自动创建不存在的品牌</el-checkbox>
        </div>
      </div>
    </el-card>

    <!-- 导入进度区域 -->
    <el-card shadow="hover" class="mb-6" v-if="importProgress.show">
      <template #header>
        <span class="font-semibold">导入进度</span>
      </template>
      <div class="space-y-4">
        <div>
          <div class="mb-2 text-sm text-gray-600">当前步骤：{{ importProgress.currentStep }}</div>
          <el-progress :percentage="importProgress.percentage" :status="importProgress.status" />
        </div>
        <div v-if="importProgress.logs.length > 0">
          <div class="mb-2 font-medium">日志输出：</div>
          <div class="bg-gray-50 dark:bg-gray-800 p-4 rounded border max-h-64 overflow-y-auto">
            <div v-for="(log, idx) in importProgress.logs" :key="idx" class="text-sm font-mono mb-1">
              {{ log }}
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 操作按钮 -->
    <div class="flex justify-end space-x-4">
      <el-button @click="handleCancel">取消</el-button>
      <el-button @click="handleViewHistory">查看导入历史</el-button>
      <el-button type="primary" @click="handleStartImport" :loading="importProgress.show && importProgress.status !== 'success' && importProgress.status !== 'exception'">
        开始导入
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { findBtMarket } from '@/api/brandtrekin/btMarket'
import {
  previewBrandSocial,
  previewGKW,
  previewKeywordHistory,
  previewProductUS,
  previewProductSales,
  batchImport
} from '@/api/brandtrekin/btImport'
import { useUserStore } from '@/pinia'

defineOptions({
  name: 'BtMarketImport'
})

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 支持通过 props 传递 marketId，或者从路由获取
const props = defineProps({
  marketId: {
    type: [String, Number],
    default: null
  }
})

const emit = defineEmits(['close', 'success'])

const marketId = computed(() => props.marketId || route.params.id || route.query.id)
const marketInfo = ref({})

// 文件输入refs
const brandSocialFileInputRef = ref()
const gkwFileInputRef = ref()
const keywordHistoryFileInputRef = ref()
const productUSFileInputRef = ref()
const productSalesFileInputRef = ref()

const uploadStatus = reactive({
  brandSocial: 'idle',
  gkw: 'idle',
  keywordHistory: 'idle',
  productUS: 'idle',
  productSales: 'idle'
})

const uploadFiles = reactive({
  brandSocial: null,
  gkw: null,
  keywordHistory: null,
  productUS: null,
  productSales: null
})

const uploadErrors = reactive({
  brandSocial: null,
  gkw: null,
  keywordHistory: null,
  productUS: null,
  productSales: null
})

const uploadPreviews = reactive({
  brandSocial: null,
  gkw: null,
  keywordHistory: null,
  productUS: null,
  productSales: null
})

// 导入选项
const importOptions = reactive({
  importMode: 'incremental', // incremental or replace
  skipInvalid: true,
  autoCreateBrand: true
})

// 导入进度
const importProgress = reactive({
  show: false,
  percentage: 0,
  status: '', // success, exception, ''
  currentStep: '',
  logs: []
})


// 获取状态类型
const getStatusType = (status) => {
  const statusMap = {
    idle: 'info',
    uploading: 'warning',
    parsing: 'warning',
    success: 'success',
    error: 'danger'
  }
  return statusMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    idle: '未上传',
    uploading: '上传中',
    parsing: '解析中',
    success: '解析成功',
    error: '解析失败'
  }
  return statusMap[status] || '未上传'
}

// 获取预览列
const getPreviewColumns = (data) => {
  if (!data || data.length === 0) return []
  return Object.keys(data[0])
}

// 文件选择处理
const handleFileSelect = async (fileType, event, expectedExt) => {
  const file = event.target.files?.[0]
  if (!file) return
  
  const ext = file.name.substring(file.name.lastIndexOf('.'))
  if (ext.toLowerCase() !== expectedExt.toLowerCase()) {
    ElMessage.error(`文件格式不正确，请上传${expectedExt}文件`)
    event.target.value = '' // 清空选择
    return
  }
  
  // 更新状态和文件
  uploadStatus[fileType] = 'uploading'
  uploadFiles[fileType] = file
  uploadErrors[fileType] = null
  uploadPreviews[fileType] = null
  
  // 根据文件类型选择预览API
  const previewApis = {
    brandSocial: previewBrandSocial,
    gkw: previewGKW,
    keywordHistory: previewKeywordHistory,
    productUS: previewProductUS,
    productSales: previewProductSales
  }
  
  const previewApi = previewApis[fileType]
  if (!previewApi) {
    uploadStatus[fileType] = 'error'
    uploadErrors[fileType] = '未知的文件类型'
    return
  }
  
  // 创建FormData并调用预览API
  const formData = new FormData()
  formData.append('file', file)
  
  uploadStatus[fileType] = 'parsing'
  
  try {
    const res = await previewApi(formData)
    if (res.code === 0) {
      uploadStatus[fileType] = 'success'
      // 提取预览数据（根据实际API返回结构调整）
      // API返回格式：{ code: 0, data: { success: true, total: 25, preview: [...], errors: [] }, msg: '预览成功' }
      const data = res.data || {}
      uploadPreviews[fileType] = data.preview || []
      ElMessage.success(`文件解析成功，共 ${data.total || 0} 条数据`)
    } else {
      uploadStatus[fileType] = 'error'
      uploadErrors[fileType] = res.msg || '解析失败'
      ElMessage.error(res.msg || '解析失败')
      event.target.value = '' // 清空选择
    }
  } catch (err) {
    uploadStatus[fileType] = 'error'
    uploadErrors[fileType] = err.message || '解析失败'
    ElMessage.error(err.message || '解析失败')
    event.target.value = '' // 清空选择
  }
}

// 开始导入
const handleStartImport = async () => {
  // 检查是否有文件已上传
  const hasFile = Object.values(uploadFiles).some(file => file !== null)
  if (!hasFile) {
    ElMessage.warning('请至少上传一个文件')
    return
  }
  
  // 确认导入
  await ElMessageBox.confirm(
    `确定要开始导入吗？\n导入模式：${importOptions.importMode === 'incremental' ? '增量导入' : '全量替换'}`,
    '确认导入',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  )
  
  // 创建FormData
  const formData = new FormData()
  formData.append('marketId', marketId.value)
  formData.append('replaceMode', importOptions.importMode === 'replace')
  
  // 添加文件
  const fileMap = {
    brandSocial: 'brandSocial',
    gkw: 'gkw',
    keywordHistory: 'keywordHistory',
    productUS: 'productUS',
    productSales: 'productSales'
  }
  
  for (const [key, formKey] of Object.entries(fileMap)) {
    if (uploadFiles[key]) {
      formData.append(formKey, uploadFiles[key])
    }
  }
  
  // 显示进度
  importProgress.show = true
  importProgress.percentage = 0
  importProgress.status = ''
  importProgress.currentStep = '准备导入...'
  importProgress.logs = []
  
  // 添加日志
  const addLog = (message) => {
    const timestamp = new Date().toLocaleString('zh-CN')
    importProgress.logs.push(`[${timestamp}] ${message}`)
  }
  
  addLog('开始导入数据...')
  
  try {
    // 调用批量导入API
    const res = await batchImport(formData)
    
    if (res.code === 0) {
      importProgress.percentage = 100
      importProgress.status = 'success'
      importProgress.currentStep = '导入完成'
      addLog('数据导入成功！')
      ElMessage.success('数据导入成功')
      
      // 3秒后关闭弹框并刷新列表
      setTimeout(() => {
        emit('success')
        emit('close')
      }, 3000)
    } else {
      importProgress.status = 'exception'
      importProgress.currentStep = '导入失败'
      addLog(`导入失败: ${res.msg || '未知错误'}`)
      ElMessage.error(res.msg || '导入失败')
    }
  } catch (error) {
    importProgress.status = 'exception'
    importProgress.currentStep = '导入失败'
    addLog(`导入失败: ${error.message || '未知错误'}`)
    ElMessage.error(error.message || '导入失败')
  }
}

// 取消
const handleCancel = () => {
  emit('close')
}

// 查看导入历史
const handleViewHistory = () => {
  emit('close')
  router.push({ path: '/layout/btImportLog', query: { marketId: marketId.value } })
}

// 获取市场信息的函数
const loadMarketInfo = async (id) => {
  if (!id) return
  try {
    const res = await findBtMarket({ ID: id })
    if (res.code === 0) {
      marketInfo.value = res.data
    } else {
      ElMessage.error('获取市场信息失败')
    }
  } catch (error) {
    ElMessage.error('获取市场信息失败')
  }
}

// 初始化
onMounted(async () => {
  // 获取市场信息
  await loadMarketInfo(marketId.value)
})

// 监听 marketId 变化（当作为组件使用时）
watch(() => props.marketId, async (newMarketId) => {
  await loadMarketInfo(newMarketId)
}, { immediate: true })
</script>

<style scoped>
.market-import-container {
  padding: 20px;
}

.space-y-4 > * + * {
  margin-top: 1rem;
}

.space-x-4 > * + * {
  margin-left: 1rem;
}
</style>

