
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="所属市场:" prop="marketId">
    <el-select v-model="formData.marketId" placeholder="请选择所属市场" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.marketId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="导入模式:" prop="importMode">
    <el-tree-select v-model="formData.importMode" placeholder="请选择导入模式" :data="import_modeOptions" style="width:100%" filterable :clearable="false" check-strictly></el-tree-select>
</el-form-item>
        <el-form-item label="导入状态:" prop="status">
    <el-tree-select v-model="formData.status" placeholder="请选择导入状态" :data="import_statusOptions" style="width:100%" filterable :clearable="false" check-strictly></el-tree-select>
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
    getBtImportLogDataSource,
  createBtImportLog,
  updateBtImportLog,
  findBtImportLog
} from '@/api/brandtrekin/btImportLog'

defineOptions({
    name: 'BtImportLogForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const import_modeOptions = ref([])
const import_statusOptions = ref([])
const formData = ref({
            marketId: undefined,
            importMode: '',
            status: '',
        })
// 验证规则
const rule = reactive({
               marketId : [{
                   required: true,
                   message: '请选择市场',
                   trigger: ['input','blur'],
               }],
               importMode : [{
                   required: true,
                   message: '请选择导入模式',
                   trigger: ['input','blur'],
               }],
               status : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getBtImportLogDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findBtImportLog({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    import_modeOptions.value = await getDictFunc('import_mode')
    import_statusOptions.value = await getDictFunc('import_status')
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createBtImportLog(formData.value)
               break
             case 'update':
               res = await updateBtImportLog(formData.value)
               break
             default:
               res = await createBtImportLog(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
