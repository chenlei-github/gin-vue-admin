
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="市场名称:" prop="marketName">
    <el-input v-model="formData.marketName" :clearable="true" placeholder="请输入市场名称" />
</el-form-item>
        <el-form-item label="市场标识符(用于URL):" prop="marketSlug">
    <el-input v-model="formData.marketSlug" :clearable="true" placeholder="请输入市场标识符(用于URL)" />
</el-form-item>
        <el-form-item label="市场描述:" prop="description">
    <RichEdit v-model="formData.description"/>
</el-form-item>
        <el-form-item label="状态:" prop="status">
    <el-tree-select v-model="formData.status" placeholder="请选择状态" :data="market_statusOptions" style="width:100%" filterable :clearable="false" check-strictly></el-tree-select>
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
  createBtMarket,
  updateBtMarket,
  findBtMarket
} from '@/api/brandtrekin/btMarket'

defineOptions({
    name: 'BtMarketForm'
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
               }],
               marketSlug : [{
                   required: true,
                   message: '请输入市场标识符',
                   trigger: ['input','blur'],
               }],
               status : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findBtMarket({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    market_statusOptions.value = await getDictFunc('market_status')
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
