
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="所属市场:" prop="marketId">
    <el-select v-model="formData.marketId" placeholder="请选择所属市场" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.marketId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="关键词:" prop="keyword">
    <el-input v-model="formData.keyword" :clearable="true" placeholder="请输入关键词" />
</el-form-item>
        <el-form-item label="来源:" prop="source">
    <el-tree-select v-model="formData.source" placeholder="请选择来源" :data="keyword_sourceOptions" style="width:100%" filterable :clearable="false" check-strictly></el-tree-select>
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
    getBtKeywordDataSource,
  createBtKeyword,
  updateBtKeyword,
  findBtKeyword
} from '@/api/brandtrekin/btKeyword'

defineOptions({
    name: 'BtKeywordForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const keyword_sourceOptions = ref([])
const formData = ref({
            marketId: undefined,
            keyword: '',
            source: '',
        })
// 验证规则
const rule = reactive({
               marketId : [{
                   required: true,
                   message: '请选择所属市场',
                   trigger: ['input','blur'],
               }],
               keyword : [{
                   required: true,
                   message: '请输入关键词',
                   trigger: ['input','blur'],
               }],
               source : [{
                   required: true,
                   message: '请选择来源',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getBtKeywordDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findBtKeyword({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    keyword_sourceOptions.value = await getDictFunc('keyword_source')
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
               res = await createBtKeyword(formData.value)
               break
             case 'update':
               res = await updateBtKeyword(formData.value)
               break
             default:
               res = await createBtKeyword(formData.value)
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
