<template>
  <div class="market-list-container">
    <!-- 顶部指标卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4 mb-6">
      <el-card shadow="hover">
        <div class="metric-card">
          <div class="metric-icon">
            <el-icon :size="40" color="#409eff"><TrendCharts /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">市场总数</div>
            <div class="metric-value">{{ totalMarkets }}</div>
          </div>
        </div>
      </el-card>

      <el-card shadow="hover">
        <div class="metric-card">
          <div class="metric-icon">
            <el-icon :size="40" color="#67c23a"><Money /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">总营收</div>
            <div class="metric-value">{{ formatCurrency(totalRevenue) }}</div>
          </div>
        </div>
      </el-card>

      <el-card shadow="hover">
        <div class="metric-card">
          <div class="metric-icon">
            <el-icon :size="40" color="#e6a23c"><Goods /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">产品总数</div>
            <div class="metric-value">{{ totalProducts }}</div>
          </div>
        </div>
      </el-card>

      <el-card shadow="hover">
        <div class="metric-card">
          <div class="metric-icon">
            <el-icon :size="40" color="#f56c6c"><Search /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">搜索量</div>
            <div class="metric-value">{{ formatNumber(totalSearchVolume) }}</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 市场列表表格 -->
    <el-card>
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg font-bold">市场列表</span>
          <el-button type="primary" @click="refreshData">
            <el-icon class="mr-1"><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="marketList"
        style="width: 100%"
        @row-click="handleRowClick"
        :row-style="{ cursor: 'pointer' }"
      >
        <el-table-column prop="name" label="市场名称" min-width="150">
          <template #default="{ row }">
            <div class="font-semibold text-blue-600">{{ row.name }}</div>
          </template>
        </el-table-column>

        <el-table-column prop="metrics.totalRevenue" label="总营收" min-width="130" sortable>
          <template #default="{ row }">
            <span class="font-medium">{{ formatCurrency(row.metrics.totalRevenue) }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="metrics.cagr" label="CAGR" width="100" sortable>
          <template #default="{ row }">
            <el-tag
              v-if="row.metrics.cagr !== null && row.metrics.cagr !== undefined"
              :type="getCAGRType(row.metrics.cagr)"
            >
              {{ formatPercent(row.metrics.cagr) }}
            </el-tag>
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>

        <el-table-column prop="metrics.brandCount" label="品牌数" width="100" sortable>
          <template #default="{ row }">
            <el-tag type="info">{{ row.metrics.brandCount }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="metrics.totalProducts" label="产品数" width="100" sortable>
          <template #default="{ row }">
            <el-tag type="warning">{{ row.metrics.totalProducts }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="metrics.searchVolume" label="搜索量" min-width="120" sortable>
          <template #default="{ row }">
            {{ formatNumber(row.metrics.searchVolume) }}
          </template>
        </el-table-column>

        <el-table-column label="12月趋势" min-width="200">
          <template #default="{ row }">
            <div class="trend-chart" :ref="(el) => setChartRef(el, row.id)"></div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click.stop="viewDetail(row.id)">
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { getMarketList } from '@/api/brandtrekin/btDisplay'
import { ElMessage } from 'element-plus'
import { TrendCharts, Money, Goods, Search, Refresh } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

defineOptions({
  name: 'MarketList'
})

const router = useRouter()
const loading = ref(false)
const marketList = ref([])
const chartRefs = ref({})

// 汇总指标
const totalMarkets = ref(0)
const totalRevenue = ref(0)
const totalProducts = ref(0)
const totalSearchVolume = ref(0)

// 设置图表引用
const setChartRef = (el, id) => {
  if (el) {
    chartRefs.value[id] = el
  }
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await getMarketList()
    if (res.code === 0) {
      marketList.value = res.data || []

      // 计算汇总指标
      totalMarkets.value = marketList.value.length
      totalRevenue.value = marketList.value.reduce((sum, m) => sum + (m.metrics.totalRevenue || 0), 0)
      totalProducts.value = marketList.value.reduce((sum, m) => sum + (m.metrics.totalProducts || 0), 0)
      totalSearchVolume.value = marketList.value.reduce((sum, m) => sum + (m.metrics.searchVolume || 0), 0)

      // 渲染趋势图表
      await nextTick()
      renderTrendCharts()
    } else {
      ElMessage.error(res.msg || '加载数据失败')
    }
  } catch (error) {
    console.error('加载市场列表失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 渲染趋势图表
const renderTrendCharts = () => {
  marketList.value.forEach((market) => {
    const chartEl = chartRefs.value[market.id]
    if (!chartEl) return

    const chart = echarts.init(chartEl)
    const trends = market.metrics.monthlyTrends || []

    const option = {
      grid: {
        left: 0,
        right: 0,
        top: 5,
        bottom: 5
      },
      xAxis: {
        type: 'category',
        show: false,
        data: trends.map(t => t.month)
      },
      yAxis: {
        type: 'value',
        show: false
      },
      series: [{
        type: 'line',
        data: trends.map(t => t.revenue),
        smooth: true,
        showSymbol: false,
        lineStyle: {
          width: 2,
          color: '#409eff'
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [{
              offset: 0, color: 'rgba(64, 158, 255, 0.3)'
            }, {
              offset: 1, color: 'rgba(64, 158, 255, 0.05)'
            }]
          }
        }
      }]
    }

    chart.setOption(option)

    // 响应式处理
    window.addEventListener('resize', () => {
      chart.resize()
    })
  })
}

// 格式化货币
const formatCurrency = (value) => {
  if (!value && value !== 0) return '$0'
  return '$' + value.toLocaleString('en-US', { maximumFractionDigits: 0 })
}

// 格式化数字
const formatNumber = (value) => {
  if (!value && value !== 0) return '0'
  return value.toLocaleString('en-US')
}

// 格式化百分比
const formatPercent = (value) => {
  if (value === null || value === undefined) return '-'
  return (value >= 0 ? '+' : '') + value.toFixed(1) + '%'
}

// 获取 CAGR 标签类型
const getCAGRType = (cagr) => {
  if (cagr >= 20) return 'success'
  if (cagr >= 10) return 'primary'
  if (cagr >= 0) return 'warning'
  return 'danger'
}

// 查看详情
const viewDetail = (marketId) => {
  router.push({ name: 'MarketDetail', params: { id: marketId } })
}

// 行点击事件
const handleRowClick = (row) => {
  viewDetail(row.id)
}

// 刷新数据
const refreshData = () => {
  loadData()
}

// 页面加载时获取数据
onMounted(() => {
  loadData()
})
</script>

<style scoped>
.market-list-container {
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
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.1) 0%, rgba(64, 158, 255, 0.05) 100%);
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

.trend-chart {
  width: 100%;
  height: 50px;
}

:deep(.el-table__row) {
  transition: all 0.3s;
}

:deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}
</style>
