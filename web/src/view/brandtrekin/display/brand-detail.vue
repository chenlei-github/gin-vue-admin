<template>
  <div class="brand-detail-container">
    <!-- 返回按钮和标题 -->
    <div class="mb-4 flex items-center justify-between">
      <div class="flex items-center">
        <el-button @click="goBack">
          <el-icon class="mr-1"><ArrowLeft /></el-icon>
          返回
        </el-button>
        <h2 class="ml-4 text-2xl font-bold">{{ brandDetail?.brand || '品牌详情' }}</h2>
      </div>
      <el-button v-if="brandDetail?.website" type="primary" @click="openWebsite">
        <el-icon class="mr-1"><Link /></el-icon>
        访问官网
      </el-button>
    </div>

    <div v-loading="loading">
      <!-- 顶部指标卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4 mb-6">
        <el-card shadow="hover">
          <div class="metric-card">
            <div class="metric-icon blue">
              <el-icon :size="40"><Money /></el-icon>
            </div>
            <div class="metric-content">
              <div class="metric-label">总营收 (近12月)</div>
              <div class="metric-value">{{ formatCurrency(brandDetail?.totalRevenue) }}</div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <div class="metric-card">
            <div class="metric-icon green">
              <el-icon :size="40"><TrendCharts /></el-icon>
            </div>
            <div class="metric-content">
              <div class="metric-label">复合增长率 (CAGR)</div>
              <div class="metric-value">
                <el-tag v-if="brandDetail?.cagr !== null" :type="getCAGRType(brandDetail?.cagr)" size="large">
                  {{ formatPercent(brandDetail?.cagr) }}
                </el-tag>
                <span v-else>-</span>
              </div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <div class="metric-card">
            <div class="metric-icon orange">
              <el-icon :size="40"><Goods /></el-icon>
            </div>
            <div class="metric-content">
              <div class="metric-label">产品数量</div>
              <div class="metric-value">{{ brandDetail?.productCount || 0 }}</div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <div class="metric-card">
            <div class="metric-icon purple">
              <el-icon :size="40"><Promotion /></el-icon>
            </div>
            <div class="metric-content">
              <div class="metric-label">官方网站</div>
              <div class="metric-value text-sm">
                <el-link v-if="brandDetail?.website" :href="brandDetail.website" target="_blank" type="primary">
                  {{ shortenUrl(brandDetail.website) }}
                </el-link>
                <span v-else class="text-gray-400">-</span>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 销售趋势图 -->
      <el-card class="mb-6">
        <template #header>
          <span class="font-bold">品牌营收趋势</span>
        </template>
        <div ref="salesTrendRef" class="chart-container"></div>
      </el-card>

      <!-- 社交媒体卡片 -->
      <div v-if="hasSocialMedia" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4 mb-6">
        <el-card v-if="brandDetail?.social?.youtube" shadow="hover" class="social-card youtube-card">
          <template #header>
            <div class="flex items-center gap-2">
              <el-icon :size="24" color="#ff0000"><VideoPlay /></el-icon>
              <span class="font-bold">YouTube</span>
            </div>
          </template>
          <div class="social-content">
            <div class="social-stat">
              <div class="stat-label">订阅数</div>
              <div class="stat-value">{{ formatNumber(brandDetail.social.youtube.subscribers) }}</div>
            </div>
            <el-button type="danger" class="w-full mt-3" @click="openLink(brandDetail.social.youtube.url)">
              访问频道
            </el-button>
          </div>
        </el-card>

        <el-card v-if="brandDetail?.social?.instagram" shadow="hover" class="social-card instagram-card">
          <template #header>
            <div class="flex items-center gap-2">
              <el-icon :size="24" color="#e4405f"><PictureFilled /></el-icon>
              <span class="font-bold">Instagram</span>
            </div>
          </template>
          <div class="social-content">
            <div class="social-stat">
              <div class="stat-label">粉丝数</div>
              <div class="stat-value">{{ formatNumber(brandDetail.social.instagram.followers) }}</div>
            </div>
            <div class="social-stat" v-if="brandDetail.social.instagram.posts">
              <div class="stat-label">帖子数</div>
              <div class="stat-value">{{ formatNumber(brandDetail.social.instagram.posts) }}</div>
            </div>
            <el-button type="danger" class="w-full mt-3" @click="openLink(brandDetail.social.instagram.url)">
              访问主页
            </el-button>
          </div>
        </el-card>

        <el-card v-if="brandDetail?.social?.facebook" shadow="hover" class="social-card facebook-card">
          <template #header>
            <div class="flex items-center gap-2">
              <el-icon :size="24" color="#1877f2"><User /></el-icon>
              <span class="font-bold">Facebook</span>
            </div>
          </template>
          <div class="social-content">
            <div class="social-stat">
              <div class="stat-label">粉丝数</div>
              <div class="stat-value">{{ formatNumber(brandDetail.social.facebook.followers) }}</div>
            </div>
            <el-button type="primary" class="w-full mt-3" @click="openLink(brandDetail.social.facebook.url)">
              访问主页
            </el-button>
          </div>
        </el-card>

        <el-card v-if="brandDetail?.social?.reddit" shadow="hover" class="social-card reddit-card">
          <template #header>
            <div class="flex items-center gap-2">
              <el-icon :size="24" color="#ff4500"><ChatDotRound /></el-icon>
              <span class="font-bold">Reddit</span>
            </div>
          </template>
          <div class="social-content">
            <div class="social-stat">
              <div class="stat-label">帖子数</div>
              <div class="stat-value">{{ formatNumber(brandDetail.social.reddit.posts) }}</div>
            </div>
            <el-button type="warning" class="w-full mt-3" @click="openLink(brandDetail.social.reddit.url)">
              访问社区
            </el-button>
          </div>
        </el-card>
      </div>

      <!-- 产品列表 -->
      <el-card>
        <template #header>
          <div class="flex justify-between items-center">
            <span class="font-bold">产品列表</span>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索产品..."
              style="width: 300px"
              clearable
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </template>

        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
          <div
            v-for="product in filteredProducts"
            :key="product.asin"
            class="product-card"
          >
            <div class="product-image-wrapper">
              <img
                :src="product.imageUrl || '/placeholder-product.png'"
                :alt="product.title"
                class="product-image"
                @error="handleImageError"
              />
            </div>
            <div class="product-content">
              <h4 class="product-title">{{ product.title }}</h4>
              <div class="product-meta">
                <div class="product-price">{{ formatCurrency(product.price) }}</div>
                <div class="product-rating">
                  <el-rate v-model="product.rating" disabled size="small" />
                  <span class="text-xs text-gray-500">({{ formatNumber(product.reviews) }})</span>
                </div>
              </div>
              <div class="product-revenue">
                <span class="text-xs text-gray-500">月销量:</span>
                <span class="font-bold text-green-600">{{ formatNumber(product.monthlySales) }}</span>
              </div>
              <!-- 产品销售趋势迷你图 -->
              <div v-if="product.salesTrend && product.salesTrend.length > 0" class="product-trend" :ref="(el) => setProductChartRef(el, product.asin)"></div>
              <el-button type="primary" link class="w-full mt-2" @click="openAmazonLink(product.asin)">
                在Amazon查看
              </el-button>
            </div>
          </div>
        </div>

        <el-empty v-if="filteredProducts.length === 0" description="暂无产品数据" />
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getBrandDetail } from '@/api/brandtrekin/btDisplay'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft, Link, Money, TrendCharts, Goods, Promotion, Search,
  VideoPlay, PictureFilled, User, ChatDotRound
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'

defineOptions({
  name: 'BrandDetail'
})

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const brandDetail = ref(null)
const searchKeyword = ref('')

// 图表引用
const salesTrendRef = ref(null)
const productChartRefs = ref({})

let salesTrendChart = null
const productCharts = {}

// 设置产品图表引用
const setProductChartRef = (el, asin) => {
  if (el) {
    productChartRefs.value[asin] = el
  }
}

// 计算属性
const hasSocialMedia = computed(() => {
  const sm = brandDetail.value?.social
  return !!(sm?.youtube || sm?.instagram || sm?.facebook || sm?.reddit)
})

const filteredProducts = computed(() => {
  if (!brandDetail.value?.products) return []
  if (!searchKeyword.value) return brandDetail.value.products

  const keyword = searchKeyword.value.toLowerCase()
  return brandDetail.value.products.filter(p =>
    p.title?.toLowerCase().includes(keyword) ||
    p.asin?.toLowerCase().includes(keyword)
  )
})

// 加载品牌详情
const loadBrandDetail = async () => {
  loading.value = true
  try {
    const { marketId, brandName } = route.params
    const res = await getBrandDetail(marketId, brandName)
    if (res.code === 0) {
      brandDetail.value = res.data
      // 渲染图表
      setTimeout(() => {
        renderSalesTrendChart()
        renderProductTrendCharts()
      }, 100)
    } else {
      ElMessage.error(res.msg || '加载数据失败')
    }
  } catch (error) {
    console.error('加载品牌详情失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 渲染销售趋势图
const renderSalesTrendChart = () => {
  if (!salesTrendRef.value || !brandDetail.value) return

  salesTrendChart = echarts.init(salesTrendRef.value)
  const trends = brandDetail.value.monthlyTrends || []

  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: (params) => {
        const item = params[0]
        return `${item.name}<br/>营收: $${item.value.toLocaleString()}`
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: trends.map(t => t.date),
      boundaryGap: false
    },
    yAxis: {
      type: 'value',
      name: '营收 ($)',
      axisLabel: {
        formatter: (value) => '$' + (value / 1000).toFixed(0) + 'K'
      }
    },
    series: [{
      name: '营收',
      type: 'line',
      data: trends.map(t => t.revenue),
      smooth: true,
      lineStyle: { width: 3 },
      itemStyle: { color: '#409eff' },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
            { offset: 1, color: 'rgba(64, 158, 255, 0.05)' }
          ]
        }
      }
    }]
  }

  salesTrendChart.setOption(option)
}

// 渲染产品趋势迷你图
const renderProductTrendCharts = () => {
  if (!brandDetail.value?.products) return

  brandDetail.value.products.forEach(product => {
    const chartEl = productChartRefs.value[product.asin]
    if (!chartEl || !product.salesTrend || product.salesTrend.length === 0) return

    const chart = echarts.init(chartEl)
    productCharts[product.asin] = chart

    const option = {
      grid: { left: 0, right: 0, top: 5, bottom: 5 },
      xAxis: {
        type: 'category',
        show: false,
        data: product.salesTrend.map(t => t.date)
      },
      yAxis: {
        type: 'value',
        show: false
      },
      series: [{
        type: 'line',
        data: product.salesTrend.map(t => t.sales),
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 2, color: '#67c23a' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(103, 194, 58, 0.3)' },
              { offset: 1, color: 'rgba(103, 194, 58, 0.05)' }
            ]
          }
        }
      }]
    }

    chart.setOption(option)
  })
}

// 工具函数
const formatCurrency = (value) => {
  if (!value && value !== 0) return '-'
  if (value === 0) return '-'
  return '$' + value.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatNumber = (value) => {
  if (!value && value !== 0) return '0'
  if (value >= 1000000) return (value / 1000000).toFixed(1) + 'M'
  if (value >= 1000) return (value / 1000).toFixed(1) + 'K'
  return value.toLocaleString('en-US')
}

const formatPercent = (value) => {
  if (value === null || value === undefined) return '-'
  return (value >= 0 ? '+' : '') + value.toFixed(1) + '%'
}

const getCAGRType = (cagr) => {
  if (cagr >= 20) return 'success'
  if (cagr >= 10) return 'primary'
  if (cagr >= 0) return 'warning'
  return 'danger'
}

const shortenUrl = (url) => {
  if (!url) return '-'
  try {
    const urlObj = new URL(url)
    return urlObj.hostname.replace('www.', '')
  } catch {
    return url.substring(0, 30) + '...'
  }
}

const handleImageError = (e) => {
  e.target.src = '/placeholder-product.png'
}

// 导航
const goBack = () => {
  router.back()
}

const openWebsite = () => {
  if (brandDetail.value?.website) {
    window.open(brandDetail.value.website, '_blank')
  }
}

const openLink = (url) => {
  if (url) window.open(url, '_blank')
}

const openAmazonLink = (asin) => {
  window.open(`https://www.amazon.com/dp/${asin}`, '_blank')
}

// 响应式处理
const handleResize = () => {
  salesTrendChart?.resize()
  Object.values(productCharts).forEach(chart => chart?.resize())
}

// 生命周期
onMounted(() => {
  loadBrandDetail()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  salesTrendChart?.dispose()
  Object.values(productCharts).forEach(chart => chart?.dispose())
})
</script>

<style scoped>
.brand-detail-container {
  padding: 20px;
  height: 100%;
  overflow-y: auto;
}

.metric-card {
  display: flex;
  align-items: center;
  gap: 16px;
}

.metric-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 60px;
  height: 60px;
  border-radius: 12px;
}

.metric-icon.blue {
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.1) 0%, rgba(64, 158, 255, 0.05) 100%);
  color: #409eff;
}

.metric-icon.green {
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.1) 0%, rgba(103, 194, 58, 0.05) 100%);
  color: #67c23a;
}

.metric-icon.orange {
  background: linear-gradient(135deg, rgba(230, 162, 60, 0.1) 0%, rgba(230, 162, 60, 0.05) 100%);
  color: #e6a23c;
}

.metric-icon.purple {
  background: linear-gradient(135deg, rgba(138, 43, 226, 0.1) 0%, rgba(138, 43, 226, 0.05) 100%);
  color: #8a2be2;
}

.metric-content {
  flex: 1;
}

.metric-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.metric-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.chart-container {
  width: 100%;
  height: 350px;
}

.social-card {
  border-radius: 12px;
  overflow: hidden;
}

.social-content {
  padding: 10px 0;
}

.social-stat {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.stat-value {
  font-size: 20px;
  font-weight: bold;
  color: #303133;
}

.product-card {
  border: 1px solid #ebeef5;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s;
  background: white;
}

.product-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-4px);
}

.product-image-wrapper {
  position: relative;
  width: 100%;
  padding-top: 100%;
  background: #f5f7fa;
  overflow: hidden;
}

.product-image {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-rank {
  position: absolute;
  top: 10px;
  right: 10px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: bold;
}

.product-content {
  padding: 16px;
}

.product-title {
  font-size: 14px;
  font-weight: 600;
  line-height: 1.4;
  height: 2.8em;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  margin-bottom: 12px;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.product-price {
  font-size: 18px;
  font-weight: bold;
  color: #f56c6c;
}

.product-rating {
  display: flex;
  align-items: center;
  gap: 4px;
}

.product-revenue {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 12px;
}

.product-trend {
  width: 100%;
  height: 50px;
  margin-bottom: 8px;
}
</style>
