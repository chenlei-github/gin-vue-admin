
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="商品ASIN:" prop="asin">
    <el-input v-model="formData.asin" :clearable="true" placeholder="请输入商品ASIN" />
</el-form-item>
        <el-form-item label="月份(YYYY-MM-01):" prop="date">
    <el-date-picker v-model="formData.date" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
        <el-form-item label="销售额:" prop="sales">
    <el-input-number v-model="formData.sales" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="销量:" prop="units">
    <el-input v-model.number="formData.units" :clearable="true" placeholder="请输入销量" />
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
  createBtProductMonthlySales,
  updateBtProductMonthlySales,
  findBtProductMonthlySales
} from '@/api/brandtrekin/btProductMonthlySales'

defineOptions({
    name: 'BtProductMonthlySalesForm'
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
const formData = ref({
            asin: '',
            date: new Date(),
            sales: 0,
            units: 0,
        })
// 验证规则
const rule = reactive({
               asin : [{
                   required: true,
                   message: '请输入ASIN',
                   trigger: ['input','blur'],
               }],
               date : [{
                   required: true,
                   message: '请选择月份',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findBtProductMonthlySales({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
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
               res = await createBtProductMonthlySales(formData.value)
               break
             case 'update':
               res = await updateBtProductMonthlySales(formData.value)
               break
             default:
               res = await createBtProductMonthlySales(formData.value)
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
