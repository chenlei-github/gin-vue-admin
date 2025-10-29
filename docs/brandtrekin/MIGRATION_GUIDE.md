# 🔄 BrandTrekin 迁移指南

> **从 Next.js 到 Gin-Vue-Admin 框架的完整迁移指南**

---

## 📋 迁移概述

### 迁移目标
将原有基于 Next.js + TypeScript + Prisma 的 BrandTrekin 系统迁移到 Gin-Vue-Admin 框架，实现：
- ✅ 更强大的权限管理系统
- ✅ 更完善的后台管理功能
- ✅ 更好的企业级特性支持
- ✅ 更高的开发效率

### 技术栈对比

| 组件 | 原版本 (Next.js) | 新版本 (GVA) |
|------|-----------------|--------------|
| 后端框架 | Next.js API Routes | Gin (Golang) |
| 前端框架 | Next.js 14 + React 18 | Vue 3 + Vite |
| 语言 | TypeScript | Golang + JavaScript |
| ORM | Prisma | GORM |
| UI库 | Tailwind CSS | Element Plus |
| 图表 | Recharts | ECharts |
| 状态管理 | React Hooks | Pinia |
| 路由 | Next.js Router | Vue Router 4 |

---

## ✅ 已完成的迁移工作

### 1. 数据库结构迁移 ✅

#### 原始表结构 → GVA表结构

| 原表名 | 新表名 | 状态 | 说明 |
|-------|--------|------|------|
| markets | bt_markets | ✅ | 添加GVA标准字段(ID, CreatedAt, UpdatedAt, DeletedAt) |
| brands | bt_brands | ✅ | 添加GVA标准字段 |
| brand_social_media | bt_brand_social_media | ✅ | 添加GVA标准字段 |
| products | bt_products | ✅ | 添加GVA标准字段 |
| product_monthly_sales | bt_product_monthly_sales | ✅ | 添加GVA标准字段 |
| keywords | bt_keywords | ✅ | 添加GVA标准字段 |
| keyword_monthly_volume | bt_keyword_monthly_volume | ✅ | 添加GVA标准字段 |
| brand_monthly_trends | bt_brand_monthly_trends | ✅ | 添加GVA标准字段 |
| market_monthly_trends | bt_market_monthly_trends | ✅ | 添加GVA标准字段 |
| import_logs | bt_import_logs | ✅ | 添加GVA标准字段 |

#### GVA标准字段说明
所有表都自动添加了以下字段：
```go
type GVA_MODEL struct {
    ID        uint           `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### 2. 后端模块迁移 ✅

#### 已创建的Go模块

**Model层** (`server/model/brandtrekin/`)
- ✅ btmarket.go - 市场数据模型
- ✅ btbrand.go - 品牌数据模型
- ✅ btbrandsocialmedia.go - 品牌社交媒体数据模型
- ✅ btproduct.go - 商品数据模型
- ✅ btproductmonthlysales.go - 商品月度销售数据模型
- ✅ btkeyword.go - 关键词数据模型
- ✅ btkeywordmonthlyvolume.go - 关键词月度搜索量数据模型
- ✅ btbrandmonthlytrend.go - 品牌月度趋势数据模型
- ✅ btmarketmonthlytrend.go - 市场月度趋势数据模型
- ✅ btimportlog.go - 导入日志数据模型

**Service层** (`server/service/brandtrekin/`)
- ✅ 所有模块的CRUD业务逻辑
- ✅ 分页查询
- ✅ 条件搜索
- ✅ 数据校验

**API层** (`server/api/v1/brandtrekin/`)
- ✅ 所有模块的RESTful API接口
- ✅ 统一响应格式
- ✅ 错误处理
- ✅ 权限验证

**Router层** (`server/router/brandtrekin/`)
- ✅ 所有模块的路由注册
- ✅ 中间件配置
- ✅ 权限控制

### 3. 前端模块迁移 ✅

#### 已创建的Vue组件

**管理页面** (`web/src/view/brandtrekin/`)
- ✅ btmarket.vue - 市场管理页面
- ✅ btbrand.vue - 品牌管理页面
- ✅ btbrandsocialmedia.vue - 品牌社交媒体管理页面
- ✅ btproduct.vue - 商品管理页面
- ✅ btproductmonthlysales.vue - 商品月度销售管理页面
- ✅ btkeyword.vue - 关键词管理页面
- ✅ btkeywordmonthlyvolume.vue - 关键词月度搜索量管理页面
- ✅ btbrandmonthlytrend.vue - 品牌月度趋势管理页面
- ✅ btmarketmonthlytrend.vue - 市场月度趋势管理页面
- ✅ btimportlog.vue - 导入日志管理页面

**API接口** (`web/src/api/brandtrekin/`)
- ✅ 所有模块的API接口定义
- ✅ Axios请求封装
- ✅ 统一错误处理

### 4. 权限系统迁移 ✅

- ✅ 所有API接口已自动注册到权限系统
- ✅ 所有菜单已自动注册到菜单系统
- ✅ 支持角色权限配置
- ✅ 支持按钮级权限控制

### 5. 字典系统迁移 ✅

已创建的业务字典：
- ✅ market_status - 市场状态（启用/禁用）
- ✅ social_platform - 社交媒体平台（YouTube/Instagram/Facebook/Reddit）
- ✅ keyword_source - 关键词来源（Google/Amazon）
- ✅ import_mode - 导入模式（增量导入/全量替换）
- ✅ import_status - 导入状态（成功/失败/部分成功）

---

## 🔄 待迁移功能

### 1. 数据导入功能 ⏳

#### 原版本实现
```typescript
// Next.js API Route
// /app/api/admin/markets/[id]/import/route.ts
export async function POST(request: Request) {
  // 解析5种文件
  // 导入数据到Prisma
  // 计算聚合指标
}
```

#### 需要在GVA中实现
```go
// server/service/brandtrekin/import_service.go
type ImportService struct{}

// 导入品牌社交媒体数据
func (s *ImportService) ImportBrandSocial(marketId uint, file *multipart.FileHeader) error {
    // 1. 解析Brand-Social.xlsx
    // 2. 校验数据
    // 3. 导入品牌数据
    // 4. 导入社交媒体数据
    // 5. 更新聚合指标
}

// 导入Google关键词数据
func (s *ImportService) ImportGKW(marketId uint, file *multipart.FileHeader) error {
    // 1. 解析GKW.csv
    // 2. 校验数据
    // 3. 导入关键词数据
    // 4. 导入月度搜索量数据
}

// 导入Amazon关键词历史数据
func (s *ImportService) ImportKeywordHistory(marketId uint, file *multipart.FileHeader) error {
    // 1. 解析KeywordHistory.xlsx
    // 2. 校验数据
    // 3. 导入关键词数据
    // 4. 导入月度搜索量数据
}

// 导入商品基础信息
func (s *ImportService) ImportProductUS(marketId uint, file *multipart.FileHeader) error {
    // 1. 解析Product-US.xlsx
    // 2. 校验数据
    // 3. 导入商品数据
}

// 导入商品月度销售数据
func (s *ImportService) ImportProductSales(marketId uint, file *multipart.FileHeader) error {
    // 1. 解析product-US-sales.xlsx
    // 2. 校验数据
    // 3. 导入月度销售数据
}

// 批量导入（5个文件）
func (s *ImportService) BatchImport(marketId uint, files map[string]*multipart.FileHeader, mode string) error {
    // 1. 开启事务
    // 2. 依次导入5个文件
    // 3. 计算聚合指标
    // 4. 提交事务或回滚
    // 5. 记录导入日志
}
```

#### 文件解析库推荐
```go
import (
    "github.com/xuri/excelize/v2"  // Excel文件解析
    "encoding/csv"                  // CSV文件解析
)
```

### 2. 数据聚合计算 ⏳

#### 原版本实现
```typescript
// /src/lib/calculators/cagr.ts
export function calculateCAGR(data: MonthlyData[]): number {
  if (data.length < 12) return null;
  const beginning = data[0].revenue;
  const ending = data[data.length - 1].revenue;
  const years = (data.length - 1) / 12;
  return Math.pow(ending / beginning, 1 / years) - 1;
}
```

#### 需要在GVA中实现
```go
// server/service/brandtrekin/aggregate_service.go
type AggregateService struct{}

// 计算CAGR
func (s *AggregateService) CalculateCAGR(trendData []TrendData) (float64, error) {
    if len(trendData) < 12 {
        return 0, errors.New("需要至少12个月的数据")
    }
    
    beginning := trendData[0].Revenue
    ending := trendData[len(trendData)-1].Revenue
    years := float64(len(trendData)-1) / 12.0
    
    if beginning <= 0 {
        return 0, errors.New("起始值必须大于0")
    }
    
    cagr := math.Pow(ending/beginning, 1/years) - 1
    
    // 限制在 -99% 到 +999% 之间
    if cagr < -0.99 {
        cagr = -0.99
    } else if cagr > 9.99 {
        cagr = 9.99
    }
    
    return cagr * 100, nil
}

// 聚合品牌月度趋势
func (s *AggregateService) AggregateBrandMonthlyTrends(marketId uint) error {
    // 1. 从商品月度销售数据聚合到品牌月度趋势
    // 2. 按品牌和月份分组求和
    // 3. 更新或插入品牌月度趋势数据
}

// 聚合市场月度趋势
func (s *AggregateService) AggregateMarketMonthlyTrends(marketId uint) error {
    // 1. 从品牌月度趋势聚合到市场月度趋势（销售额）
    // 2. 从关键词月度搜索量聚合到市场月度趋势（搜索量）
    // 3. 更新或插入市场月度趋势数据
}

// 更新市场聚合指标
func (s *AggregateService) UpdateMarketMetrics(marketId uint) error {
    // 1. 计算总销售额（最近12个月）
    // 2. 计算商品总数
    // 3. 计算品牌数量
    // 4. 计算搜索量（最近月份）
    // 5. 计算CAGR
    // 6. 更新市场表
}

// 更新品牌聚合指标
func (s *AggregateService) UpdateBrandMetrics(brandId uint) error {
    // 1. 计算总销售额（最近12个月）
    // 2. 计算商品数量
    // 3. 计算CAGR
    // 4. 更新品牌表
}
```

### 3. 前端展示页面 ⏳

#### 原版本页面
- `/` - 市场列表页
- `/market/[id]` - 市场详情页
- `/market/[id]/brand/[brand]` - 品牌详情页

#### 需要在GVA中实现

**市场列表页** (`web/src/view/brandtrekin/display/market-list.vue`)
```vue
<template>
  <div class="market-list">
    <!-- 4个核心指标卡片 -->
    <el-row :gutter="20">
      <el-col :span="6">
        <metric-card icon="Building2" label="追踪市场数" :value="metrics.marketCount" />
      </el-col>
      <el-col :span="6">
        <metric-card icon="TrendingUp" label="总市场规模" :value="formatCurrency(metrics.totalRevenue)" />
      </el-col>
      <el-col :span="6">
        <metric-card icon="Package" label="商品总数" :value="metrics.totalProducts" />
      </el-col>
      <el-col :span="6">
        <metric-card icon="Search" label="搜索量" :value="formatNumber(metrics.searchVolume)" />
      </el-col>
    </el-row>

    <!-- 市场列表表格 -->
    <el-table :data="markets" style="margin-top: 20px">
      <el-table-column prop="marketName" label="市场名称" />
      <el-table-column prop="totalRevenue" label="市场规模" :formatter="formatRevenue" />
      <el-table-column prop="cagr" label="市场增速" :formatter="formatCAGR" />
      <el-table-column prop="searchVolume" label="市场声量" :formatter="formatVolume" />
      <el-table-column prop="brandCount" label="品牌数量" />
      <el-table-column prop="totalProducts" label="商品数量" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button type="primary" @click="viewDetail(row.id)">查看详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 市场趋势预览 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12" v-for="market in topMarkets" :key="market.id">
        <market-trend-card :market="market" />
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getMarketList } from '@/api/brandtrekin/display'

const markets = ref([])
const metrics = ref({})

onMounted(async () => {
  const res = await getMarketList()
  markets.value = res.data.markets
  metrics.value = res.data.metrics
})
</script>
```

**市场详情页** (`web/src/view/brandtrekin/display/market-detail.vue`)
```vue
<template>
  <div class="market-detail">
    <!-- 面包屑导航 -->
    <el-breadcrumb>
      <el-breadcrumb-item :to="{ path: '/brandtrekin/display/market-list' }">市场分析</el-breadcrumb-item>
      <el-breadcrumb-item>{{ market.marketName }}</el-breadcrumb-item>
    </el-breadcrumb>

    <!-- 4个核心指标卡片 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6">
        <metric-card icon="TrendingUp" label="市场规模" :value="formatCurrency(market.totalRevenue)" />
      </el-col>
      <el-col :span="6">
        <metric-card icon="TrendingUp" label="市场增速" :value="formatPercent(market.cagr)" />
      </el-col>
      <el-col :span="6">
        <metric-card icon="Search" label="市场声量" :value="formatNumber(market.searchVolume)" />
      </el-col>
      <el-col :span="6">
        <metric-card icon="Building2" label="品牌数量" :value="market.brandCount" />
      </el-col>
    </el-row>

    <!-- 市场趋势图表 -->
    <el-card style="margin-top: 20px">
      <template #header>
        <span>市场趋势分析</span>
      </template>
      <div ref="trendChart" style="height: 400px"></div>
    </el-card>

    <!-- 品牌分布图表 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>品牌销售额占比</span>
          </template>
          <div ref="pieChart" style="height: 300px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>品牌销售额排名</span>
          </template>
          <div ref="barChart" style="height: 300px"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 品牌列表表格 -->
    <el-card style="margin-top: 20px">
      <template #header>
        <span>品牌列表</span>
      </template>
      <el-table :data="brands">
        <el-table-column prop="brandName" label="品牌名称" />
        <el-table-column prop="website" label="品牌独立站" />
        <el-table-column prop="totalRevenue" label="品牌规模" :formatter="formatRevenue" />
        <el-table-column prop="cagr" label="品牌增速" :formatter="formatCAGR" />
        <el-table-column label="社交媒体热度">
          <template #default="{ row }">
            <social-media-icons :social="row.social" />
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="{ row }">
            <el-button type="primary" @click="viewBrand(row.brandName)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import * as echarts from 'echarts'
import { getMarketDetail } from '@/api/brandtrekin/display'

const route = useRoute()
const market = ref({})
const brands = ref([])
const trendChart = ref(null)
const pieChart = ref(null)
const barChart = ref(null)

onMounted(async () => {
  const res = await getMarketDetail(route.params.id)
  market.value = res.data.market
  brands.value = res.data.brands
  
  // 初始化图表
  initTrendChart()
  initPieChart()
  initBarChart()
})

function initTrendChart() {
  const chart = echarts.init(trendChart.value)
  // ECharts配置...
}

function initPieChart() {
  const chart = echarts.init(pieChart.value)
  // ECharts配置...
}

function initBarChart() {
  const chart = echarts.init(barChart.value)
  // ECharts配置...
}
</script>
```

**品牌详情页** (`web/src/view/brandtrekin/display/brand-detail.vue`)
```vue
<template>
  <div class="brand-detail">
    <!-- 面包屑导航 -->
    <el-breadcrumb>
      <el-breadcrumb-item :to="{ path: '/brandtrekin/display/market-list' }">市场分析</el-breadcrumb-item>
      <el-breadcrumb-item :to="{ path: `/brandtrekin/display/market-detail/${marketId}` }">
        {{ market.marketName }}
      </el-breadcrumb-item>
      <el-breadcrumb-item>{{ brand.brandName }}</el-breadcrumb-item>
    </el-breadcrumb>

    <!-- 4个核心指标卡片 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6">
        <metric-card icon="TrendingUp" label="品牌规模" :value="formatCurrency(brand.totalRevenue)" />
      </el-col>
      <el-col :span="6">
        <metric-card icon="TrendingUp" label="品牌增速" :value="formatPercent(brand.cagr)" />
      </el-col>
      <el-col :span="6">
        <metric-card icon="Package" label="商品数量" :value="brand.productCount" />
      </el-col>
      <el-col :span="6">
        <metric-card icon="Globe" label="品牌独立站">
          <el-button v-if="brand.website" type="primary" @click="openWebsite">访问网站</el-button>
          <span v-else>暂无数据</span>
        </metric-card>
      </el-col>
    </el-row>

    <!-- 品牌销售趋势图表 -->
    <el-card style="margin-top: 20px">
      <template #header>
        <span>品牌销售趋势</span>
      </template>
      <div ref="trendChart" style="height: 300px"></div>
    </el-card>

    <!-- 社交媒体数据卡片 -->
    <el-card style="margin-top: 20px">
      <template #header>
        <span>社交媒体数据</span>
      </template>
      <el-row :gutter="20">
        <el-col :span="6" v-for="platform in socialPlatforms" :key="platform.name">
          <social-media-card :platform="platform" />
        </el-col>
      </el-row>
    </el-card>

    <!-- 品牌商品网格 -->
    <el-card style="margin-top: 20px">
      <template #header>
        <span>品牌商品</span>
      </template>
      <el-row :gutter="20">
        <el-col :span="6" v-for="product in products" :key="product.asin">
          <product-card :product="product" />
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import * as echarts from 'echarts'
import { getBrandDetail } from '@/api/brandtrekin/display'

const route = useRoute()
const market = ref({})
const brand = ref({})
const products = ref([])
const socialPlatforms = ref([])
const trendChart = ref(null)

onMounted(async () => {
  const res = await getBrandDetail(route.params.marketId, route.params.brandName)
  market.value = res.data.market
  brand.value = res.data.brand
  products.value = res.data.products
  socialPlatforms.value = res.data.socialPlatforms
  
  // 初始化图表
  initTrendChart()
})

function initTrendChart() {
  const chart = echarts.init(trendChart.value)
  // ECharts配置...
}
</script>
```

---

## 📝 迁移步骤建议

### Step 1: 数据导入功能开发（优先级：高）

1. **创建导入服务**
   ```bash
   # 创建文件
   touch server/service/brandtrekin/import_service.go
   ```

2. **实现文件解析**
   - 安装依赖：`go get github.com/xuri/excelize/v2`
   - 实现5个文件的解析方法
   - 添加数据校验逻辑

3. **实现导入API**
   ```bash
   # 在 server/api/v1/brandtrekin/ 中添加导入接口
   ```

4. **前端上传组件**
   ```bash
   # 创建文件上传页面
   touch web/src/view/brandtrekin/import/data-import.vue
   ```

### Step 2: 数据聚合计算开发（优先级：高）

1. **创建聚合服务**
   ```bash
   touch server/service/brandtrekin/aggregate_service.go
   ```

2. **实现CAGR计算**
   - 实现CAGR计算函数
   - 添加单元测试

3. **实现聚合逻辑**
   - 品牌月度趋势聚合
   - 市场月度趋势聚合
   - 市场指标更新
   - 品牌指标更新

4. **创建定时任务**
   - 定期更新聚合指标
   - 或在数据导入后触发

### Step 3: 前端展示页面开发（优先级：中）

1. **创建展示页面目录**
   ```bash
   mkdir -p web/src/view/brandtrekin/display
   ```

2. **开发市场列表页**
   - 实现4个指标卡片
   - 实现市场列表表格
   - 实现趋势预览图

3. **开发市场详情页**
   - 实现4个指标卡片
   - 实现趋势图表（ECharts）
   - 实现品牌分布图表
   - 实现品牌列表表格

4. **开发品牌详情页**
   - 实现4个指标卡片
   - 实现销售趋势图
   - 实现社交媒体卡片
   - 实现商品网格

### Step 4: 测试与优化（优先级：中）

1. **功能测试**
   - 测试数据导入功能
   - 测试数据聚合计算
   - 测试前端展示页面

2. **性能测试**
   - 测试大数据量导入
   - 测试图表渲染性能
   - 优化慢查询

3. **兼容性测试**
   - 测试不同浏览器
   - 测试响应式布局

---

## 🔧 开发工具和资源

### Go开发工具
- **IDE**: GoLand 或 VS Code + Go插件
- **调试**: Delve
- **测试**: go test
- **文档**: godoc

### Vue开发工具
- **IDE**: VS Code + Volar插件
- **调试**: Vue DevTools
- **测试**: Vitest
- **文档**: VitePress

### 推荐的Go包
```go
// Excel解析
github.com/xuri/excelize/v2

// CSV解析
encoding/csv

// 数学计算
math

// 时间处理
time

// 数据库操作
gorm.io/gorm
```

### 推荐的npm包
```json
{
  "echarts": "^5.5.1",           // 图表库
  "element-plus": "^2.10.2",     // UI组件库
  "axios": "^1.8.2",             // HTTP客户端
  "pinia": "^2.2.2",             // 状态管理
  "vue-router": "^4.4.3"         // 路由管理
}
```

---

## 📞 技术支持

### 遇到问题？

1. **查看GVA官方文档**
   - https://www.gin-vue-admin.com

2. **查看原始需求文档**
   - `docs/brandtrekin/PRD_*.md`

3. **参考示例代码**
   - `server/model/example/`
   - `web/src/view/example/`

4. **加入GVA社区**
   - QQ群：971857775

---

## ✅ 迁移检查清单

### 基础框架 ✅
- [x] 数据库表结构迁移
- [x] 后端模块创建
- [x] 前端管理页面创建
- [x] 权限系统配置
- [x] 字典系统配置

### 核心功能 ⏳
- [ ] 数据导入功能
- [ ] 数据聚合计算
- [ ] CAGR自动计算
- [ ] 前端展示页面
- [ ] 数据可视化图表

### 测试与优化 ⏳
- [ ] 功能测试
- [ ] 性能测试
- [ ] 兼容性测试
- [ ] 代码优化
- [ ] 文档完善

---

**迁移基础框架已完成！接下来按照上述步骤继续开发即可。🚀**
