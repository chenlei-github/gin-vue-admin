
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="所属品牌:" prop="brandId">
    <el-select v-model="formData.brandId" placeholder="请选择所属品牌" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.brandId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="社交平台:" prop="platform">
    <el-tree-select v-model="formData.platform" placeholder="请选择社交平台" :data="social_platformOptions" style="width:100%" filterable :clearable="false" check-strictly></el-tree-select>
</el-form-item>
        <el-form-item label="平台链接:" prop="url">
    <el-input v-model="formData.url" :clearable="true" placeholder="请输入平台链接" />
</el-form-item>
        <el-form-item label="订阅数(YouTube):" prop="subscribers">
    <el-input v-model.number="formData.subscribers" :clearable="true" placeholder="请输入订阅数(YouTube)" />
</el-form-item>
        <el-form-item label="粉丝数(Instagram/Facebook):" prop="followers">
    <el-input v-model.number="formData.followers" :clearable="true" placeholder="请输入粉丝数(Instagram/Facebook)" />
</el-form-item>
        <el-form-item label="帖子数(Reddit):" prop="posts">
    <el-input v-model.number="formData.posts" :clearable="true" placeholder="请输入帖子数(Reddit)" />
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
    getBtBrandSocialMediaDataSource,
  createBtBrandSocialMedia,
  updateBtBrandSocialMedia,
  findBtBrandSocialMedia
} from '@/api/brandtrekin/btBrandSocialMedia'

defineOptions({
    name: 'BtBrandSocialMediaForm'
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
const social_platformOptions = ref([])
const formData = ref({
            brandId: undefined,
            platform: '',
            url: '',
            subscribers: 0,
            followers: 0,
            posts: 0,
        })
// 验证规则
const rule = reactive({
               brandId : [{
                   required: true,
                   message: '请选择所属品牌',
                   trigger: ['input','blur'],
               }],
               platform : [{
                   required: true,
                   message: '请选择社交平台',
                   trigger: ['input','blur'],
               }],
               url : [{
                   required: true,
                   message: '请输入平台链接',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getBtBrandSocialMediaDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findBtBrandSocialMedia({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    social_platformOptions.value = await getDictFunc('social_platform')
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
               res = await createBtBrandSocialMedia(formData.value)
               break
             case 'update':
               res = await updateBtBrandSocialMedia(formData.value)
               break
             default:
               res = await createBtBrandSocialMedia(formData.value)
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
