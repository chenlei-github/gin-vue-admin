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
              <div class="metric-label">月度搜索量总和</div>
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

      <!-- 图表区域 - 柱状图和饼图并排 -->
      <div class="grid grid-cols-1 xl:grid-cols-2 gap-4 mb-6">
        <!-- 品牌排名柱状图 -->
        <el-card>
          <template #header>
            <span class="font-bold">品牌营收排名 (Top 8)</span>
          </template>
          <div ref="barChartRef" class="chart-container"></div>
        </el-card>

        <!-- 品牌营收占比饼图 -->
        <el-card>
          <template #header>
            <span class="font-bold">品牌营收占比 (Top 8)</span>
          </template>
          <div ref="pieChartRef" class="chart-container"></div>
        </el-card>
      </div>

      <!-- 营收和搜索量趋势图 (双Y轴) - 单独一行 -->
      <el-card class="mb-6">
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-bold">营收与搜索量趋势</span>
            <el-select 
              v-model="selectedBrands" 
              multiple
              collapse-tags
              collapse-tags-tooltip
              placeholder="选择品牌对比" 
              style="width: 300px"
              @change="handleBrandChange"
            >
              <el-option label="全部品牌" value="all" />
              <el-option 
                v-for="brand in marketDetail?.brands || []" 
                :key="brand.brand" 
                :label="brand.brand" 
                :value="brand.brand" 
              />
            </el-select>
          </div>
        </template>
        <div ref="trendChartRef" class="chart-container-large"></div>
      </el-card>

      <!-- 品牌列表表格 - 最下方 -->
      <el-card class="mb-6">
        <template #header>
          <span class="font-bold">品牌列表</span>
        </template>
        <el-table 
          :data="marketDetail?.brands || []" 
          style="width: 100%"
          :max-height="400"
          stripe
        >
          <el-table-column prop="brand" label="品牌名称" min-width="150"  align="center" fixed>
            <template #default="{ row }">
              <el-button type="primary" link @click="viewBrandDetail(row.brand)">
                {{ row.brand }}
              </el-button>
            </template>
          </el-table-column>

          <el-table-column prop="totalRevenue" label="营收"  sortable>
            <template #default="{ row }">
              <span class="font-medium">{{ formatCurrency(row.totalRevenue) }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="cagr" label="品牌增速CAGR" sortable>
            <template #default="{ row }">
              <el-tag v-if="row.cagr !== null" :type="getCAGRType(row.cagr)">
                {{ formatPercent(row.cagr) }}
              </el-tag>
              <span v-else>-</span>
            </template>
          </el-table-column>

          <el-table-column prop="productCount" label="产品数" sortable>
            <template #default="{ row }">
              <el-tag type="warning">{{ row.productCount }}</el-tag>
            </template>
          </el-table-column>

          <el-table-column label="独立站">
            <template #default="{ row }">
              <el-tooltip v-if="row.website" placement="top" content="访问品牌官网">
                <a :href="row.website" target="_blank" class="website-link">
                  <el-icon :size="20" color="#409eff"><Link /></el-icon>
                </a>
              </el-tooltip>
              <span v-else class="text-gray-400">-</span>
            </template>
          </el-table-column>

          <el-table-column label="社交媒体" min-width="250">
            <template #default="{ row }">
              <div class="flex gap-2 items-center">
                <el-tooltip v-if="row.social?.youtube?.url" placement="top">
                  <template #content>
                    YouTube<br/>
                    订阅者: {{ formatNumber(row.social.youtube.subscribers) }}
                  </template>
                  <a :href="row.social.youtube.url" target="_blank" class="social-icon youtube">
                    <el-icon :size="16"><VideoPlay /></el-icon>
                  </a>
                </el-tooltip>
                <el-tooltip v-if="row.social?.instagram?.url" placement="top">
                  <template #content>
                    Instagram<br/>
                    粉丝: {{ formatNumber(row.social.instagram.followers) }}
                  </template>
                  <a :href="row.social.instagram.url" target="_blank" class="social-icon instagram">
                    <el-icon :size="16"><PictureFilled /></el-icon>
                  </a>
                </el-tooltip>
                <el-tooltip v-if="row.social?.facebook?.url" placement="top">
                  <template #content>
                    Facebook<br/>
                    粉丝: {{ formatNumber(row.social.facebook.followers) }}
                  </template>
                  <a :href="row.social.facebook.url" target="_blank" class="social-icon facebook">
                    <el-icon :size="16"><User /></el-icon>
                  </a>
                </el-tooltip>
                <el-tooltip v-if="row.social?.reddit?.url" placement="top">
                  <template #content>
                    Reddit<br/>
                    帖子数: {{ formatNumber(row.social.reddit.posts) }}
                  </template>
                  <a :href="row.social.reddit.url" target="_blank" class="social-icon reddit">
                    <el-icon :size="16"><ChatDotRound /></el-icon>
                  </a>
                </el-tooltip>
                <span v-if="!hasSocialMedia(row.social)" class="text-gray-400">-</span>
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
  VideoPlay, PictureFilled, User, ChatDotRound, Link
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'

defineOptions({
  name: 'MarketDetail'
})

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const marketDetail = ref(null)
const selectedBrands = ref(['all'])

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

// 品牌颜色配置
const brandColors = [
  '#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399',
  '#00d4ff', '#ff6b9d', '#c990ff', '#ffb800', '#00c9a7'
]

// 渲染趋势图 (支持多品牌对比)
const renderTrendChart = () => {
  if (!trendChartRef.value || !marketDetail.value) return

  if (!trendChart) {
    trendChart = echarts.init(trendChartRef.value)
  }

  // 获取所有日期（使用市场总趋势的日期作为基准）
  const allDates = marketDetail.value.metrics?.monthlyTrends?.map(t => t.date) || []
  
  const legendData = []
  const series = []
  
  // 如果选择了"全部品牌"或没有选择任何品牌
  if (selectedBrands.value.includes('all') || selectedBrands.value.length === 0) {
    const trends = marketDetail.value.metrics?.monthlyTrends || []
    legendData.push('总营收', '搜索量')
    
    series.push({
      name: '总营收',
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
    })
    
    series.push({
      name: '搜索量',
      type: 'line',
      yAxisIndex: 1,
      data: trends.map(t => t.searchVolume),
      smooth: true,
      lineStyle: { width: 3 },
      itemStyle: { color: '#67c23a' }
    })
  } else {
    // 显示选中的多个品牌
    const validBrands = selectedBrands.value.filter(b => b !== 'all')
    
    validBrands.forEach((brandName, index) => {
      const brand = marketDetail.value.brands?.find(b => b.brand === brandName)
      if (brand && brand.monthlyTrends) {
        const color = brandColors[index % brandColors.length]
        legendData.push(brandName)
        
        series.push({
          name: brandName,
          type: 'line',
          data: brand.monthlyTrends.map(t => t.revenue),
          smooth: true,
          lineStyle: { width: 3 },
          itemStyle: { color: color },
          emphasis: {
            focus: 'series'
          }
        })
      }
    })
    
    // 如果选择了具体品牌，也显示总搜索量作为参考
    if (validBrands.length > 0) {
      const trends = marketDetail.value.metrics?.monthlyTrends || []
      legendData.push('市场搜索量')
      series.push({
        name: '市场搜索量',
        type: 'line',
        yAxisIndex: 1,
        data: trends.map(t => t.searchVolume),
        smooth: true,
        lineStyle: { width: 2, type: 'dashed' },
        itemStyle: { color: '#67c23a' },
        opacity: 0.6
      })
    }
  }

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      formatter: (params) => {
        let result = params[0].axisValue + '<br/>'
        params.forEach(param => {
          if (param.seriesName.includes('搜索量')) {
            result += `${param.marker}${param.seriesName}: ${formatNumber(param.value)}<br/>`
          } else {
            result += `${param.marker}${param.seriesName}: ${formatCurrency(param.value)}<br/>`
          }
        })
        return result
      }
    },
    legend: {
      data: legendData,
      top: 10,
      type: 'scroll'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '60px',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: allDates,
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
    series: series
  }

  trendChart.setOption(option, true)
}

// 处理品牌切换
const handleBrandChange = () => {
  renderTrendChart()
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
          name: brand.brand,
          value: brand.totalRevenue || 0
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
      bottom: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: topBrands.map(b => b.brand),
      axisLabel: {
        interval: 0,
        rotate: 45,
        fontSize: 11
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
        data: topBrands.map(b => b.totalRevenue || 0),
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
          fontSize: 10,
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

const hasSocialMedia = (social) => {
  if (!social) return false
  return !!(social.youtube?.url || social.instagram?.url || social.facebook?.url || social.reddit?.url)
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
  width: 28px;
  height: 28px;
  border-radius: 50%;
  color: white;
  transition: all 0.3s;
}

.social-icon:hover {
  transform: scale(1.15);
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

.website-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  cursor: pointer;
}

.website-link:hover {
  transform: scale(1.2);
  opacity: 0.8;
}
</style>
