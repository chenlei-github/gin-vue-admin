
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="createdAtRange">
      <template #label>
        <span>
          创建日期
          <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>

      <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
       </el-form-item>
      
            <el-form-item label="市场名称" prop="marketName">
  <el-input v-model="searchInfo.marketName" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="市场标识符(用于URL)" prop="marketSlug">
  <el-input v-model="searchInfo.marketSlug" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="状态" prop="status">
    <el-tree-select v-model="formData.status" placeholder="请选择状态" :data="market_statusOptions" style="width:100%" filterable :clearable="false" check-strictly ></el-tree-select>
</el-form-item>
            
            <el-form-item label="总销售额(最近12个月)" prop="totalRevenue">
  <el-input class="!w-40" v-model.number="searchInfo.startTotalRevenue" placeholder="最小值" />
  —
  <el-input class="!w-40" v-model.number="searchInfo.endTotalRevenue" placeholder="最大值" />
</el-form-item>
            

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button  type="primary" icon="plus" @click="openDialog()">新增</el-button>
            <el-button  icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
            <ExportTemplate  template-id="brandtrekin_BtMarket" />
            <ExportExcel  template-id="brandtrekin_BtMarket" filterDeleted/>
            <ImportExcel  template-id="brandtrekin_BtMarket" @on-success="getTableData" />
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange"
        >
        <el-table-column type="selection" width="55" />
        
        <el-table-column sortable align="left" label="ID" prop="ID" width="80" />
        
        <el-table-column sortable align="left" label="创建时间" prop="CreatedAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        
            <el-table-column align="left" label="市场名称" prop="marketName" width="120" />

            <el-table-column align="left" label="市场标识符(用于URL)" prop="marketSlug" width="150" />

            <el-table-column align="left" label="状态" prop="status" width="100">
    <template #default="scope">
      <el-switch
        v-model="scope.row.status"
        :active-value="'active'"
        :inactive-value="'inactive'"
        @change="handleStatusChange(scope.row)"
      />
    </template>
</el-table-column>
            <el-table-column sortable align="left" label="总销售额(最近12个月)" prop="totalRevenue" width="120" />

            <el-table-column sortable align="left" label="商品总数" prop="totalProducts" width="120" />

            <el-table-column sortable align="left" label="品牌数量" prop="brandCount" width="120" />

            <el-table-column sortable align="left" label="搜索量(最近月份)" prop="searchVolume" width="120" />

            <el-table-column sortable align="left" label="年复合增长率(%)" prop="cagr" width="120" />

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <el-button  type="primary" link icon="edit" class="table-button" @click="updateBtMarketFunc(scope.row)">编辑</el-button>
            <el-button  type="primary" link icon="Upload" class="table-button" @click="goToImport(scope.row)">导入数据</el-button>
            <el-button   type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'新增':'编辑'}}</span>
                <div>
                  <el-button @click="closeDialog">取 消</el-button>
                  <el-button v-if="type==='create'" :loading="btnLoading" type="primary" @click="enterDialogAndImport">保存并导入数据</el-button>
                  <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="市场名称:" prop="marketName">
    <el-input v-model="formData.marketName" :clearable="true" placeholder="请输入市场名称" @blur="autoGenerateSlug" />
</el-form-item>
            <el-form-item label="市场标识符(用于URL):" prop="marketSlug">
              <div style="display: flex; gap: 8px;">
                <el-input 
                  v-model="formData.marketSlug" 
                  :clearable="true" 
                  :disabled="type==='update'"
                  placeholder="请输入市场标识符(用于URL)" 
                  @blur="validateSlug"
                />
                <el-button v-if="type==='create'" @click="manualGenerateSlug">自动生成</el-button>
              </div>
              <div v-if="type==='update'" style="color: #909399; font-size: 12px; margin-top: 4px;">
                市场ID创建后不可修改，以保持URL稳定性
              </div>
              <div v-else style="color: #909399; font-size: 12px; margin-top: 4px;">
                用于URL，只能包含小写字母、数字和连字符
              </div>
</el-form-item>
            <el-form-item label="市场描述:" prop="description">
    <RichEdit v-model="formData.description"/>
</el-form-item>
            <el-form-item label="状态:" prop="status">
    <el-tree-select v-model="formData.status" placeholder="请选择状态" :data="market_statusOptions" style="width:100%" filterable :clearable="false" check-strictly></el-tree-select>
</el-form-item>
          </el-form>
    </el-drawer>

    <!-- 数据导入弹框 -->
    <el-drawer 
      destroy-on-close 
      :size="appStore.drawerSize" 
      v-model="importDrawerVisible" 
      :show-close="false" 
      :before-close="closeImportDrawer"
      direction="rtl"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">导入市场数据 - {{ currentMarket?.marketName || '' }}</span>
          <div>
            <el-button @click="closeImportDrawer">关闭</el-button>
          </div>
        </div>
      </template>
      <MarketImportComponent 
        v-if="currentMarket"
        :market-id="currentMarket.ID" 
        @close="closeImportDrawer"
        @success="handleImportSuccess"
      />
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="市场名称">
    {{ detailForm.marketName }}
</el-descriptions-item>
                    <el-descriptions-item label="市场标识符(用于URL)">
    {{ detailForm.marketSlug }}
</el-descriptions-item>
                    <el-descriptions-item label="市场描述">
    <RichView v-model="detailForm.description" />
</el-descriptions-item>
                    <el-descriptions-item label="状态">
    {{ detailForm.status }}
</el-descriptions-item>
                    <el-descriptions-item label="总销售额(最近12个月)">
    {{ detailForm.totalRevenue }}
</el-descriptions-item>
                    <el-descriptions-item label="商品总数">
    {{ detailForm.totalProducts }}
</el-descriptions-item>
                    <el-descriptions-item label="品牌数量">
    {{ detailForm.brandCount }}
</el-descriptions-item>
                    <el-descriptions-item label="搜索量(最近月份)">
    {{ detailForm.searchVolume }}
</el-descriptions-item>
                    <el-descriptions-item label="年复合增长率(%)">
    {{ detailForm.cagr }}
</el-descriptions-item>
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
  createBtMarket,
  deleteBtMarket,
  deleteBtMarketByIds,
  updateBtMarket,
  findBtMarket,
  getBtMarketList,
  toggleMarketStatus,
  validateDeleteMarket,
  generateSlugFromName,
  validateSlugUnique
} from '@/api/brandtrekin/btMarket'
import MarketImportComponent from './btMarketImport.vue'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'
import RichView from '@/components/richtext/rich-view.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"
import { useRouter } from 'vue-router'

// 导出组件
import ExportExcel from '@/components/exportExcel/exportExcel.vue'
// 导入组件
import ImportExcel from '@/components/exportExcel/importExcel.vue'
// 导出模板组件
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'


defineOptions({
    name: 'BtMarket'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()
const router = useRouter()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const market_statusOptions = ref([])
const formData = ref({
            marketName: '',
            marketSlug: '',
            description: '',
            status: '',
        })



// 验证规则
const rule = reactive({
               marketName : [{
                   required: true,
                   message: '请输入市场名称',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               marketSlug : [{
                   required: true,
                   message: '请输入市场标识符',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
               },
               {
                   pattern: /^[a-z0-9-]+$/,
                   message: '只能包含小写字母、数字和连字符',
                   trigger: ['input', 'blur'],
               }
              ],
               status : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
// 排序
const sortChange = ({ prop, order }) => {
  const sortMap = {
    CreatedAt:"created_at",
    ID:"id",
            totalRevenue: 'total_revenue',
            totalProducts: 'total_products',
            brandCount: 'brand_count',
            searchVolume: 'search_volume',
            cagr: 'cagr',
  }

  let sort = sortMap[prop]
  if(!sort){
   sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  }

  searchInfo.value.sort = sort
  searchInfo.value.order = order
  getTableData()
}
// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getBtMarketList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
    market_statusOptions.value = await getDictFunc('market_status')
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = async (row) => {
  // 获取市场信息用于确认
  const res = await validateDeleteMarket({ ID: row.ID })
  if (res.code !== 0) {
    ElMessage.error('获取市场信息失败')
    return
  }

  const marketInfo = res.data
  
  // 显示自定义删除确认弹窗
  ElMessageBox.prompt(
    `删除市场将同时删除该市场下的所有数据，包括：
• 所有品牌信息
• 所有商品信息
• 所有关键词数据
• 所有销售数据
• 所有趋势数据

⚠️ 此操作不可恢复！

请输入市场名称以确认删除：`,
    '确认删除',
    {
      confirmButtonText: '确认删除',
      cancelButtonText: '取消',
      inputPlaceholder: `请输入：${marketInfo.marketName || ''}`,
      inputValidator: (value) => {
        if (!value || value.trim() === '') {
          return '请输入市场名称'
        }
        if (value.trim() !== marketInfo.marketName) {
          return '市场名称不匹配，请重新输入'
        }
        return true
      },
      type: 'warning',
      dangerouslyUseHTMLString: false
    }
  ).then(() => {
    deleteBtMarketFunc(row)
  }).catch(() => {
    // 取消删除
  })
}

// 状态切换
const handleStatusChange = async (row) => {
  try {
    const res = await toggleMarketStatus({ ID: row.ID }, { status: row.status })
    if (res.code === 0) {
      ElMessage.success('状态更新成功')
    } else {
      // 恢复原状态
      row.status = row.status === 'active' ? 'inactive' : 'active'
      ElMessage.error(res.msg || '状态更新失败')
    }
  } catch (error) {
    // 恢复原状态
    row.status = row.status === 'active' ? 'inactive' : 'active'
    ElMessage.error('状态更新失败')
  }
}

// 打开导入数据弹框
const goToImport = async (row) => {
  // 获取完整的市场信息
  try {
    const res = await findBtMarket({ ID: row.ID })
    if (res.code === 0) {
      currentMarket.value = res.data
      importDrawerVisible.value = true
    } else {
      ElMessage.error('获取市场信息失败')
    }
  } catch (error) {
    ElMessage.error('获取市场信息失败')
  }
}

// 关闭导入弹框
const closeImportDrawer = () => {
  importDrawerVisible.value = false
  currentMarket.value = null
}

// 导入成功回调
const handleImportSuccess = () => {
  // 刷新列表
  getTableData()
}

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          IDs.push(item.ID)
        })
      const res = await deleteBtMarketByIds({ IDs })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === IDs.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
      })
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateBtMarketFunc = async(row) => {
    const res = await findBtMarket({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteBtMarketFunc = async (row) => {
    const res = await deleteBtMarket({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 导入弹框控制
const importDrawerVisible = ref(false)
const currentMarket = ref(null)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        marketName: '',
        marketSlug: '',
        description: '',
        status: '',
        }
}

// 自动生成slug（从市场名称）
const autoGenerateSlug = async () => {
  if (type.value === 'create' && formData.value.marketName && !formData.value.marketSlug) {
    try {
      const res = await generateSlugFromName({ name: formData.value.marketName })
      if (res.code === 0 && res.data.slug) {
        formData.value.marketSlug = res.data.slug
      }
    } catch (error) {
      // 生成失败时不处理
    }
  }
}

// 手动生成slug
const manualGenerateSlug = async () => {
  if (!formData.value.marketName) {
    ElMessage.warning('请先输入市场名称')
    return
  }
  try {
    const res = await generateSlugFromName({ name: formData.value.marketName })
    if (res.code === 0 && res.data.slug) {
      formData.value.marketSlug = res.data.slug
      ElMessage.success('已自动生成市场ID')
      // 立即校验唯一性
      await validateSlug()
    } else {
      ElMessage.error('生成失败')
    }
  } catch (error) {
    ElMessage.error('生成失败')
  }
}

// 校验slug唯一性
const validateSlug = async () => {
  if (!formData.value.marketSlug || type.value === 'update') {
    return
  }
  try {
    const excludeID = type.value === 'update' ? formData.value.ID : undefined
    const params = { slug: formData.value.marketSlug }
    if (excludeID) {
      params.excludeID = excludeID
    }
    const res = await validateSlugUnique(params)
    if (res.code === 0) {
      if (!res.data.isUnique) {
        ElMessage.warning('该市场ID已存在，请修改')
      }
    }
  } catch (error) {
    // 校验失败时不处理
  }
}

// 弹窗确定
const enterDialog = async () => {
     btnLoading.value = true
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return btnLoading.value = false
              let res
              switch (type.value) {
                case 'create':
                  res = await createBtMarket(formData.value)
                  break
                case 'update':
                  res = await updateBtMarket(formData.value)
                  break
                default:
                  res = await createBtMarket(formData.value)
                  break
              }
              btnLoading.value = false
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
                })
                closeDialog()
                getTableData()
              }
      })
}

// 保存并导入数据
const enterDialogAndImport = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    if (type.value !== 'create') {
      btnLoading.value = false
      return
    }
    const res = await createBtMarket(formData.value)
    btnLoading.value = false
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '市场创建成功，即将打开导入页面'
      })
      // 通过slug查询新创建的市场信息
      const listRes = await getBtMarketList({ page: 1, pageSize: 10, marketSlug: formData.value.marketSlug })
      if (listRes.code === 0 && listRes.data.list && listRes.data.list.length > 0) {
        currentMarket.value = listRes.data.list[0]
        closeDialog()
        getTableData()
        // 打开导入弹框
        importDrawerVisible.value = true
      } else {
        ElMessage.error('无法获取新创建的市场信息')
        closeDialog()
        getTableData()
      }
    }
  })
}

const detailForm = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findBtMarket({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}


</script>

<style>

</style>
