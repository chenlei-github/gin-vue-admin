
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="所属市场:" prop="marketId">
    <el-select v-model="formData.marketId" placeholder="请选择所属市场" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.marketId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="所属品牌:" prop="brandId">
    <el-select v-model="formData.brandId" placeholder="请选择所属品牌" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.brandId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="亚马逊ASIN:" prop="asin">
    <el-input v-model="formData.asin" :clearable="true" placeholder="请输入亚马逊ASIN" />
</el-form-item>
        <el-form-item label="商品标题:" prop="title">
    <el-input v-model="formData.title" :clearable="true" placeholder="请输入商品标题" />
</el-form-item>
        <el-form-item label="价格:" prop="price">
    <el-input-number v-model="formData.price" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="评分(0-5):" prop="rating">
    <el-input-number v-model="formData.rating" style="width:100%" :precision="2" :clearable="true" />
</el-form-item>
        <el-form-item label="评论数:" prop="reviews">
    <el-input v-model.number="formData.reviews" :clearable="true" placeholder="请输入评论数" />
</el-form-item>
        <el-form-item label="月销量:" prop="monthlySales">
    <el-input v-model.number="formData.monthlySales" :clearable="true" placeholder="请输入月销量" />
</el-form-item>
        <el-form-item label="商品图片:" prop="imageUrl">
    <SelectImage
     v-model="formData.imageUrl"
     file-type="image"
    />
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
    getBtProductDataSource,
  createBtProduct,
  updateBtProduct,
  findBtProduct
} from '@/api/brandtrekin/btProduct'

defineOptions({
    name: 'BtProductForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
// 图片选择组件
import SelectImage from '@/components/selectImage/selectImage.vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const formData = ref({
            marketId: undefined,
            brandId: undefined,
            asin: '',
            title: '',
            price: 0,
            rating: 0,
            reviews: 0,
            monthlySales: 0,
            imageUrl: "",
        })
// 验证规则
const rule = reactive({
               marketId : [{
                   required: true,
                   message: '请选择所属市场',
                   trigger: ['input','blur'],
               }],
               brandId : [{
                   required: true,
                   message: '请选择所属品牌',
                   trigger: ['input','blur'],
               }],
               asin : [{
                   required: true,
                   message: '请输入ASIN',
                   trigger: ['input','blur'],
               }],
               title : [{
                   required: true,
                   message: '请输入商品标题',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getBtProductDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findBtProduct({ ID: route.query.id })
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
               res = await createBtProduct(formData.value)
               break
             case 'update':
               res = await updateBtProduct(formData.value)
               break
             default:
               res = await createBtProduct(formData.value)
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
