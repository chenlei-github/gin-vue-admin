<template>
  <div class="market-detail-container">
    <!-- 返回按钮和标题 -->
    <div class="mb-4 flex items-center">
      <el-button @click="goBack">
        <el-icon class="mr-1"><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2 class="ml-4 text-2xl font-bold">{{ marketDetail?.name || '市场详情' }}</h2>
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
              <div class="metric-value">{{ formatCurrency(marketDetail?.metrics?.totalRevenue) }}</div>
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
                <el-tag v-if="marketDetail?.metrics?.cagr !== null" :type="getCAGRType(marketDetail?.metrics?.cagr)" size="large">
                  {{ formatPercent(marketDetail?.metrics?.cagr) }}
                </el-tag>
                <span v-else>-</span>
              </div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <div class="metric-card">
            <div class="metric-icon orange">
              <el-icon :size="40"><Search /></el-icon>
            </div>
            <div class="metric-content">
              <div class="metric-label">搜索量</div>
              <div class="metric-value">{{ formatNumber(marketDetail?.metrics?.searchVolume) }}</div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover">
          <div class="metric-card">
            <div class="metric-icon red">
              <el-icon :size="40"><ShoppingBag /></el-icon>
            </div>
            <div class="metric-content">
              <div class="metric-label">品牌数量</div>
              <div class="metric-value">{{ marketDetail?.metrics?.brandCount || 0 }}</div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 图表区域 -->
      <div class="grid grid-cols-1 xl:grid-cols-2 gap-4 mb-6">
        <!-- 营收和搜索量趋势图 (双Y轴) -->
        <el-card>
          <template #header>
            <span class="font-bold">营收与搜索量趋势</span>
          </template>
          <div ref="trendChartRef" class="chart-container"></div>
        </el-card>

        <!-- 品牌营收占比饼图 -->
        <el-card>
          <template #header>
            <span class="font-bold">品牌营收占比 (Top 8)</span>
          </template>
          <div ref="pieChartRef" class="chart-container"></div>
        </el-card>
      </div>

      <!-- 品牌排名柱状图 -->
      <el-card class="mb-6">
        <template #header>
          <span class="font-bold">品牌营收排名 (Top 8)</span>
        </template>
        <div ref="barChartRef" class="chart-container-large"></div>
      </el-card>

      <!-- 品牌列表表格 -->
      <el-card>
        <template #header>
          <span class="font-bold">品牌列表</span>
        </template>
        <el-table :data="marketDetail?.brands || []" style="width: 100%">
          <el-table-column prop="brandName" label="品牌名称" min-width="150">
            <template #default="{ row }">
              <el-button type="primary" link @click="viewBrandDetail(row.brandName)">
                {{ row.brandName }}
              </el-button>
            </template>
          </el-table-column>

          <el-table-column prop="revenue" label="营收" min-width="130" sortable>
            <template #default="{ row }">
              <span class="font-medium">{{ formatCurrency(row.revenue) }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="cagr" label="CAGR" width="100" sortable>
            <template #default="{ row }">
              <el-tag v-if="row.cagr !== null" :type="getCAGRType(row.cagr)">
                {{ formatPercent(row.cagr) }}
              </el-tag>
              <span v-else>-</span>
            </template>
          </el-table-column>

          <el-table-column prop="productCount" label="产品数" width="100" sortable>
            <template #default="{ row }">
              <el-tag type="warning">{{ row.productCount }}</el-tag>
            </template>
          </el-table-column>

          <el-table-column label="社交媒体" min-width="200">
            <template #default="{ row }">
              <div class="flex gap-2">
                <el-tooltip v-if="row.socialMedia?.youtubeChannel" content="YouTube" placement="top">
                  <a :href="row.socialMedia.youtubeChannel" target="_blank" class="social-icon youtube">
                    <el-icon :size="20"><VideoPlay /></el-icon>
                  </a>
                </el-tooltip>
                <el-tooltip v-if="row.socialMedia?.instagram" content="Instagram" placement="top">
                  <a :href="row.socialMedia.instagram" target="_blank" class="social-icon instagram">
                    <el-icon :size="20"><PictureFilled /></el-icon>
                  </a>
                </el-tooltip>
                <el-tooltip v-if="row.socialMedia?.facebook" content="Facebook" placement="top">
                  <a :href="row.socialMedia.facebook" target="_blank" class="social-icon facebook">
                    <el-icon :size="20"><User /></el-icon>
                  </a>
                </el-tooltip>
                <el-tooltip v-if="row.socialMedia?.reddit" content="Reddit" placement="top">
                  <a :href="row.socialMedia.reddit" target="_blank" class="social-icon reddit">
                    <el-icon :size="20"><ChatDotRound /></el-icon>
                  </a>
                </el-tooltip>
                <span v-if="!hasSocialMedia(row.socialMedia)" class="text-gray-400">-</span>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMarketDetail } from '@/api/brandtrekin/btDisplay'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft, Money, TrendCharts, Search, ShoppingBag,
  VideoPlay, PictureFilled, User, ChatDotRound
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'

defineOptions({
  name: 'MarketDetail'
})

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const marketDetail = ref(null)

// 图表引用
const trendChartRef = ref(null)
const pieChartRef = ref(null)
const barChartRef = ref(null)

let trendChart = null
let pieChart = null
let barChart = null

// 加载市场详情
const loadMarketDetail = async () => {
  loading.value = true
  try {
    const marketId = route.params.id
    const res = await getMarketDetail(marketId)
    if (res.code === 0) {
      marketDetail.value = res.data
      // 渲染图表
      setTimeout(() => {
        renderTrendChart()
        renderPieChart()
        renderBarChart()
      }, 100)
    } else {
      ElMessage.error(res.msg || '加载数据失败')
    }
  } catch (error) {
    console.error('加载市场详情失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 渲染趋势图 (双Y轴)
const renderTrendChart = () => {
  if (!trendChartRef.value || !marketDetail.value) return

  trendChart = echarts.init(trendChartRef.value)
  const trends = marketDetail.value.metrics?.monthlyTrends || []

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      }
    },
    legend: {
      data: ['营收', '搜索量']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: trends.map(t => t.month),
      boundaryGap: false
    },
    yAxis: [
      {
        type: 'value',
        name: '营收 ($)',
        position: 'left',
        axisLabel: {
          formatter: (value) => '$' + (value / 1000).toFixed(0) + 'K'
        }
      },
      {
        type: 'value',
        name: '搜索量',
        position: 'right',
        axisLabel: {
          formatter: (value) => (value / 1000).toFixed(0) + 'K'
        }
      }
    ],
    series: [
      {
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
      },
      {
        name: '搜索量',
        type: 'line',
        yAxisIndex: 1,
        data: trends.map(t => t.searchVolume),
        smooth: true,
        lineStyle: { width: 3 },
        itemStyle: { color: '#67c23a' }
      }
    ]
  }

  trendChart.setOption(option)
}

// 渲染饼图
const renderPieChart = () => {
  if (!pieChartRef.value || !marketDetail.value) return

  pieChart = echarts.init(pieChartRef.value)
  const brands = marketDetail.value.brands || []
  const topBrands = brands.slice(0, 8)

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: ${c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'right',
      top: 'center'
    },
    series: [
      {
        name: '品牌营收',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 20,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: topBrands.map(brand => ({
          name: brand.brandName,
          value: brand.revenue || 0
        }))
      }
    ]
  }

  pieChart.setOption(option)
}

// 渲染柱状图
const renderBarChart = () => {
  if (!barChartRef.value || !marketDetail.value) return

  barChart = echarts.init(barChartRef.value)
  const brands = marketDetail.value.brands || []
  const topBrands = brands.slice(0, 8)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
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
      data: topBrands.map(b => b.brandName),
      axisLabel: {
        interval: 0,
        rotate: 30
      }
    },
    yAxis: {
      type: 'value',
      name: '营收 ($)',
      axisLabel: {
        formatter: (value) => '$' + (value / 1000).toFixed(0) + 'K'
      }
    },
    series: [
      {
        name: '营收',
        type: 'bar',
        data: topBrands.map(b => b.revenue || 0),
        itemStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: '#409eff' },
              { offset: 1, color: '#67c23a' }
            ]
          },
          borderRadius: [5, 5, 0, 0]
        },
        label: {
          show: true,
          position: 'top',
          formatter: (params) => '$' + (params.value / 1000).toFixed(0) + 'K'
        }
      }
    ]
  }

  barChart.setOption(option)
}

// 工具函数
const formatCurrency = (value) => {
  if (!value && value !== 0) return '$0'
  return '$' + value.toLocaleString('en-US', { maximumFractionDigits: 0 })
}

const formatNumber = (value) => {
  if (!value && value !== 0) return '0'
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

const hasSocialMedia = (socialMedia) => {
  if (!socialMedia) return false
  return !!(socialMedia.youtubeChannel || socialMedia.instagram || socialMedia.facebook || socialMedia.reddit)
}

// 导航
const goBack = () => {
  router.back()
}

const viewBrandDetail = (brandName) => {
  router.push({
    name: 'BrandDetail',
    params: {
      marketId: route.params.id,
      brandName: brandName
    }
  })
}

// 响应式处理
const handleResize = () => {
  trendChart?.resize()
  pieChart?.resize()
  barChart?.resize()
}

// 生命周期
onMounted(() => {
  loadMarketDetail()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  pieChart?.dispose()
  barChart?.dispose()
})
</script>

<style scoped>
.market-detail-container {
  padding: 20px;
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

.metric-icon.red {
  background: linear-gradient(135deg, rgba(245, 108, 108, 0.1) 0%, rgba(245, 108, 108, 0.05) 100%);
  color: #f56c6c;
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
  height: 300px;
}

.chart-container-large {
  width: 100%;
  height: 400px;
}

.social-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  color: white;
  transition: all 0.3s;
}

.social-icon:hover {
  transform: scale(1.1);
}

.social-icon.youtube {
  background-color: #ff0000;
}

.social-icon.instagram {
  background: linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%);
}

.social-icon.facebook {
  background-color: #1877f2;
}

.social-icon.reddit {
  background-color: #ff4500;
}
</style>
