# 📊 前端展示需求文档

> **本文档详细描述所有前端展示页面的布局、数据指标和计算公式**

---

## 1️⃣ 首页 - 市场列表页 (`/`)

### 1.1 页面布局结构

```
┌─────────────────────────────────────────────────────────────┐
│  [页面标题区域]                                              │
│  标题：市场分析                                              │
│  副标题：探索和分析不同市场的竞争情况，了解市场趋势和品牌表现  │
└─────────────────────────────────────────────────────────────┘

┌──────────┬──────────┬──────────┬──────────┐
│ 卡片1    │ 卡片2    │ 卡片3    │ 卡片4    │
│ 追踪市场数│ 总市场规模│ 商品总数  │ 搜索量    │
│ [图标]   │ [图标]   │ [图标]   │ [图标]   │
│ 数值     │ 数值     │ 数值     │ 数值     │
│ 标签     │ 趋势     │ 标签     │ 标签     │
└──────────┴──────────┴──────────┴──────────┘

┌─────────────────────────────────────────────────────────────┐
│  [市场列表表格]                                              │
│  ┌────────┬────────┬────────┬────────┬────────┬────────┐   │
│  │市场名称│市场规模│市场增速│市场声量│品牌数量│操作    │   │
│  ├────────┼────────┼────────┼────────┼────────┼────────┤   │
│  │ ...    │ ...    │ ...    │ ...    │ ...    │[查看]  │   │
│  └────────┴────────┴────────┴────────┴────────┴────────┘   │
└─────────────────────────────────────────────────────────────┘

┌──────────────────────┬──────────────────────┐
│ [市场趋势预览1]      │ [市场趋势预览2]      │
│ 市场名称             │ 市场名称             │
│ 趋势百分比           │ 趋势百分比           │
│ [迷你柱状图]         │ [迷你柱状图]         │
└──────────────────────┴──────────────────────┘
```

### 1.2 核心指标卡片（4个）

#### 卡片1：追踪市场数
```typescript
{
  icon: "Building2",           // Lucide图标
  iconColor: "bg-blue-500",    // 蓝色背景
  label: "追踪市场数",
  value: number,               // 市场总数
  subLabel: "已激活市场"
}
```

**数据来源**：`markets` 表  
**计算公式**：`COUNT(id) WHERE status='active'`  
**数据类型**：Integer

---

#### 卡片2：总市场规模
```typescript
{
  icon: "TrendingUp",
  iconColor: "bg-green-500",   // 绿色背景
  label: "总市场规模",
  value: string,               // 格式化后的金额（$XXM/$XXB）
  trend: {
    value: number,             // 平均增长百分比
    direction: "up" | "down",
    icon: "ArrowUp" | "ArrowDown"
  },
  subLabel: "平均增长"
}
```

**数据来源**：`markets` 表  
**计算公式**：`SUM(total_revenue)`  
**数据类型**：Decimal(15,2)  
**格式化规则**：
- < 1,000,000: $XXK (千)
- < 1,000,000,000: $XXM (百万)
- >= 1,000,000,000: $XXB (十亿)

---

#### 卡片3：商品总数
```typescript
{
  icon: "Package",
  iconColor: "bg-pink-500",    // 粉色背景
  label: "商品总数",
  value: number,               // 商品数量
  subLabel: "跨所有市场"
}
```

**数据来源**：`markets` 表  
**计算公式**：`SUM(total_products)`  
**数据类型**：Integer

---

#### 卡片4：搜索量
```typescript
{
  icon: "Search",
  iconColor: "bg-yellow-500",  // 黄色背景
  label: "搜索量",
  value: string,               // 格式化后的数量（XXK/XXM）
  subLabel: "月度"
}
```

**数据来源**：`markets` 表  
**计算公式**：`SUM(search_volume)`  
**数据类型**：BigInt  
**格式化规则**：
- < 1,000: 直接显示数字
- < 1,000,000: XXK
- >= 1,000,000: XXM

---

### 1.3 市场列表表格

#### 表头结构
```
| 市场名称 | 市场规模 | 市场增速 | 市场声量 | 品牌数量 | 商品数量 | 操作 |
```

#### 每行数据结构
```typescript
{
  // 列1：市场名称
  name: {
    primary: string,      // 主标题（如：CNC Router Machine）
    secondary: string     // 副标题（固定：亚马逊美国市场）
  },
  
  // 列2：市场规模
  revenue: {
    primary: string,      // 销售额（格式化：$13.8M）
    secondary: string     // 固定："年度预估"
  },
  
  // 列3：市场增速
  cagr: {
    primary: string,      // CAGR百分比（如：+15.5%）
    secondary: string,    // 固定："年复合增长率"
    trend: "up" | "down" | "neutral",
    icon: "TrendingUp" | "TrendingDown" | "Minus",
    color: "text-green-600" | "text-red-600" | "text-gray-600"
  },
  
  // 列4：市场声量
  searchVolume: {
    primary: string,      // 搜索量（格式化：1.5M）
    secondary: string     // 固定："月搜索量"
  },
  
  // 列5：品牌数量
  brandCount: {
    primary: number,      // 品牌数量
    secondary: string     // 固定："活跃品牌"
  },
  
  // 列6：商品数量
  productCount: {
    primary: number,      // 商品数量
    secondary: string     // 固定："在售商品"
  },
  
  // 列7：操作
  action: {
    label: "查看详情",
    icon: "ArrowRight",
    link: `/market/${market_slug}`
  }
}
```

#### 数据指标计算

| 指标名称 | 数据来源 | 计算公式 | 数据类型 |
|---------|---------|---------|---------|
| 市场规模 | `market_monthly_trends` | `SUM(revenue) WHERE market_id=X AND date >= DATE_SUB(NOW(), INTERVAL 12 MONTH)` | Decimal(15,2) |
| 市场增速(CAGR) | `market_monthly_trends` | 见下方CAGR计算公式 | Decimal(5,2) |
| 市场声量 | `keyword_monthly_volume` + `keywords` | `SUM(volume) WHERE market_id=X AND date=最近月份` | BigInt |
| 品牌数量 | `brands` | `COUNT(DISTINCT id) WHERE market_id=X` | Integer |
| 商品数量 | `products` | `COUNT(asin) WHERE market_id=X` | Integer |

**CAGR计算公式**：
```javascript
// 需要至少12个月的数据
function calculateCAGR(monthlyData) {
  if (monthlyData.length < 12) return null;
  
  const sortedData = monthlyData.sort((a, b) => new Date(a.date) - new Date(b.date));
  const beginningValue = sortedData[0].revenue;
  const endingValue = sortedData[sortedData.length - 1].revenue;
  const years = (sortedData.length - 1) / 12;
  
  if (beginningValue <= 0) return null;
  
  const cagr = (Math.pow(endingValue / beginningValue, 1 / years) - 1) * 100;
  
  // 限制在 -99% 到 +999% 之间
  return Math.max(-99, Math.min(999, cagr));
}
```

---

### 1.4 市场趋势预览区（2列网格）

#### 每个市场卡片结构
```typescript
{
  marketName: string,           // 市场名称
  title: string,                // "过去12个月趋势"
  trend: {
    value: number,              // 趋势百分比
    direction: "up" | "down",
    icon: "TrendingUp" | "TrendingDown"
  },
  chartData: Array<{
    month: string,              // 月份（简写：Jan, Feb...）
    revenue: number,            // 销售额
    opacity: number             // 透明度（0.3-1.0，渐变效果）
  }>,
  link: string                  // 跳转链接：/market/${market_slug}
}
```

#### 迷你柱状图规格
- **数据点数量**：12个月
- **柱状条样式**：圆角顶部，渐变透明度
- **透明度计算**：`opacity = 0.3 + (index / 11) * 0.7`
- **颜色**：蓝色 `#3b82f6`
- **高度**：80px
- **交互**：Hover显示具体数值

---

## 2️⃣ 市场详情页 (`/market/[id]`)

### 2.1 页面布局结构

```
┌─────────────────────────────────────────────────────────────┐
│  [面包屑导航]                                                │
│  市场分析 > {市场名称}                                       │
│                                                [返回]按钮    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [页面标题区域]                                              │
│  标题：{市场名称}                                            │
│  副标题：分析该细分市场的竞争情况、品牌表现和市场趋势        │
└─────────────────────────────────────────────────────────────┘

┌──────────┬──────────┬──────────┬──────────┐
│ 市场规模  │ 市场增速  │ 市场声量  │ 品牌数量  │
│ [图标]   │ [图标]   │ [图标]   │ [图标]   │
│ 数值     │ 数值     │ 数值     │ 数值     │
│ 说明     │ 说明     │ 说明     │ 说明     │
└──────────┴──────────┴──────────┴──────────┘

┌─────────────────────────────────────────────────────────────┐
│  [市场趋势图表]                                              │
│  标题：市场趋势分析                                          │
│  [双轴折线图：销售额(左轴) + 搜索量(右轴)]                   │
│  [可切换对比其他市场]                                        │
└─────────────────────────────────────────────────────────────┘

┌──────────────────────┬──────────────────────┐
│ [品牌销售额占比]      │ [品牌销售额排名]      │
│ 饼图（Top 8品牌）    │ 柱状图（Top 8品牌）  │
└──────────────────────┴──────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [品牌列表表格]                                              │
│  ┌────────┬────────┬────────┬────────┬────────┬────────┐   │
│  │品牌名称│独立站  │品牌规模│品牌增速│社交媒体│操作    │   │
│  ├────────┼────────┼────────┼────────┼────────┼────────┤   │
│  │ ...    │ ...    │ ...    │ ...    │ ...    │[查看]  │   │
│  └────────┴────────┴────────┴────────┴────────┴────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 核心指标卡片（4个）

#### 卡片1：市场规模
```typescript
{
  icon: "TrendingUp",
  iconColor: "bg-green-500",
  label: "市场规模",
  value: string,               // 格式化金额
  subLabel: "年度市场规模预估"
}
```

**数据来源**：`market_monthly_trends` 表  
**计算公式**：
```sql
SELECT SUM(revenue) 
FROM market_monthly_trends 
WHERE market_id = ? 
  AND date >= DATE_SUB(NOW(), INTERVAL 12 MONTH)
```

---

#### 卡片2：市场增速
```typescript
{
  icon: "TrendingUp" | "TrendingDown",
  iconColor: "bg-green-500" | "bg-red-500",
  label: "市场增速",
  value: string,               // "+15.5%" 或 "N/A"
  subLabel: "年复合增长率" | "数据不足"
}
```

**数据来源**：`market_monthly_trends` 表  
**计算公式**：同首页CAGR计算  
**显示规则**：
- 数据 >= 12个月：显示CAGR百分比
- 数据 < 12个月：显示 "N/A"，副标题显示 "数据不足"

---

#### 卡片3：市场声量
```typescript
{
  icon: "Search",
  iconColor: "bg-yellow-500",
  label: "市场声量",
  value: string,               // 格式化搜索量
  subLabel: "月度搜索量总和"
}
```

**数据来源**：`keyword_monthly_volume` + `keywords` 表  
**计算公式**：
```sql
SELECT SUM(kmv.volume)
FROM keyword_monthly_volume kmv
JOIN keywords k ON kmv.keyword_id = k.id
WHERE k.market_id = ?
  AND kmv.date = (SELECT MAX(date) FROM keyword_monthly_volume)
```

---

#### 卡片4：品牌数量
```typescript
{
  icon: "Building2",
  iconColor: "bg-blue-500",
  label: "品牌数量",
  value: number,
  subLabel: "活跃品牌总数"
}
```

**数据来源**：`brands` 表  
**计算公式**：`COUNT(DISTINCT id) WHERE market_id = ?`

---

### 2.3 市场趋势图表

#### 图表配置
```typescript
{
  type: "ComposedChart",       // Recharts组件
  height: 400,
  data: Array<{
    date: string,              // 月份（YYYY-MM）
    revenue: number,           // 销售额（左Y轴）
    searchVolume: number       // 搜索量（右Y轴）
  }>,
  
  // 左Y轴：销售额
  leftAxis: {
    dataKey: "revenue",
    stroke: "#3b82f6",         // 蓝色
    strokeWidth: 2,
    name: "销售额"
  },
  
  // 右Y轴：搜索量
  rightAxis: {
    dataKey: "searchVolume",
    stroke: "#f97316",         // 橙色
    strokeWidth: 2,
    name: "搜索量"
  },
  
  // 交互
  tooltip: true,
  legend: true,
  grid: true
}
```

#### 数据来源
**销售额趋势**：
```sql
SELECT date, revenue 
FROM market_monthly_trends 
WHERE market_id = ? 
ORDER BY date ASC
LIMIT 24  -- 最多显示24个月
```

**搜索量趋势**：
```sql
SELECT kmv.date, SUM(kmv.volume) as volume
FROM keyword_monthly_volume kmv
JOIN keywords k ON kmv.keyword_id = k.id
WHERE k.market_id = ?
GROUP BY kmv.date
ORDER BY kmv.date ASC
LIMIT 24
```

---

### 2.4 品牌分布图表（2列）

#### 左侧：品牌销售额占比（饼图）
```typescript
{
  type: "PieChart",
  data: Array<{
    name: string,              // 品牌名称
    value: number,             // 销售额
    percentage: number,        // 占比百分比
    fill: string               // 颜色
  }>,
  colors: [
    "#3b82f6", "#10b981", "#f59e0b", "#ef4444",
    "#8b5cf6", "#ec4899", "#06b6d4", "#84cc16"
  ],
  showTop: 8,                  // 只显示Top 8品牌
  label: true,                 // 显示百分比标签
  legend: true
}
```

**数据来源**：
```sql
SELECT 
  b.brand_name,
  SUM(bmt.revenue) as total_revenue
FROM brand_monthly_trends bmt
JOIN brands b ON bmt.brand_id = b.id
WHERE b.market_id = ?
  AND bmt.date >= DATE_SUB(NOW(), INTERVAL 12 MONTH)
GROUP BY b.id
ORDER BY total_revenue DESC
LIMIT 8
```

---

#### 右侧：品牌销售额排名（柱状图）
```typescript
{
  type: "BarChart",
  data: Array<{
    name: string,              // 品牌名称
    revenue: number            // 销售额
  }>,
  xAxis: {
    dataKey: "name",
    angle: -45,                // 斜45度显示
    textAnchor: "end"
  },
  yAxis: {
    label: "销售额"
  },
  bar: {
    fill: "#3b82f6",
    radius: [8, 8, 0, 0]       // 圆角顶部
  }
}
```

**数据来源**：同饼图

---

### 2.5 品牌列表表格

#### 表头结构
```
| 品牌名称 | 品牌独立站 | 品牌规模 | 品牌增速 | 社交媒体热度 | 操作 |
```

#### 每行数据结构
```typescript
{
  // 列1：品牌名称
  brand: {
    primary: string,           // 品牌名
    secondary: string          // "X个商品"
  },
  
  // 列2：品牌独立站
  website: {
    url: string | null,
    icon: "Globe",
    display: string | "-"
  },
  
  // 列3：品牌规模
  revenue: {
    primary: string,           // 格式化金额
    secondary: "年度预估"
  },
  
  // 列4：品牌增速
  cagr: {
    primary: string,           // "+20.3%" 或 "N/A"
    secondary: string,
    trend: "up" | "down" | "neutral",
    icon: string
  },
  
  // 列5：社交媒体热度
  social: {
    youtube: {
      icon: "Youtube",
      color: "text-red-600",
      value: string,           // 格式化订阅数
      url: string | null
    },
    instagram: {
      icon: "Instagram",
      color: "text-pink-600",
      value: string,
      url: string | null
    },
    facebook: {
      icon: "Facebook",
      color: "text-blue-600",
      value: string,
      url: string | null
    },
    reddit: {
      icon: "MessageCircle",
      color: "text-orange-600",
      value: string,           // 帖子数
      url: string | null
    }
  },
  
  // 列6：操作
  action: {
    label: "查看详情",
    link: `/market/${market_slug}/brand/${brand_name}`
  }
}
```

#### 数据来源
```sql
SELECT 
  b.id,
  b.brand_name,
  b.website,
  b.total_revenue,
  b.cagr,
  b.product_count,
  bsm_yt.url as youtube_url,
  bsm_yt.subscribers as youtube_subscribers,
  bsm_ig.url as instagram_url,
  bsm_ig.followers as instagram_followers,
  bsm_fb.url as facebook_url,
  bsm_fb.followers as facebook_followers,
  bsm_rd.url as reddit_url,
  bsm_rd.posts as reddit_posts
FROM brands b
LEFT JOIN brand_social_media bsm_yt ON b.id = bsm_yt.brand_id AND bsm_yt.platform = 'youtube'
LEFT JOIN brand_social_media bsm_ig ON b.id = bsm_ig.brand_id AND bsm_ig.platform = 'instagram'
LEFT JOIN brand_social_media bsm_fb ON b.id = bsm_fb.brand_id AND bsm_fb.platform = 'facebook'
LEFT JOIN brand_social_media bsm_rd ON b.id = bsm_rd.brand_id AND bsm_rd.platform = 'reddit'
WHERE b.market_id = ?
ORDER BY b.total_revenue DESC
```

---

## 3️⃣ 品牌详情页 (`/market/[id]/brand/[brand]`)

### 3.1 页面布局结构

```
┌─────────────────────────────────────────────────────────────┐
│  [面包屑导航]                                                │
│  市场分析 > {市场名称} > {品牌名称}                          │
│                                                [返回市场]    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [页面标题区域]                                              │
│  标题：{品牌名称}                                            │
│  副标题：品牌销售趋势、社交媒体数据和产品信息                │
└─────────────────────────────────────────────────────────────┘

┌──────────┬──────────┬──────────┬──────────┐
│ 品牌规模  │ 品牌增速  │ 商品数量  │ 品牌独立站│
│ [图标]   │ [图标]   │ [图标]   │ [按钮]   │
│ 数值     │ 数值     │ 数值     │          │
│ 说明     │ 说明     │ 说明     │          │
└──────────┴──────────┴──────────┴──────────┘

┌─────────────────────────────────────────────────────────────┐
│  [品牌销售趋势图表]                                          │
│  标题：品牌销售趋势                                          │
│  副标题：品牌销售额历史趋势（X个月数据）                     │
│  [折线图]                                                    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [社交媒体数据卡片]                                          │
│  标题：社交媒体数据                                          │
│  副标题：品牌在各社交媒体平台的表现                          │
│  ┌──────────┬──────────┬──────────┬──────────┐            │
│  │ YouTube  │Instagram │ Facebook │ Reddit   │            │
│  │ [图标]   │ [图标]   │ [图标]   │ [图标]   │            │
│  │ 订阅数   │ 粉丝数   │ 关注数   │ 帖子数   │            │
│  │ [访问]   │ [访问]   │ [访问]   │ [访问]   │            │
│  └──────────┴──────────┴──────────┴──────────┘            │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [品牌商品网格]                                              │
│  标题：品牌商品                                              │
│  副标题：该品牌在亚马逊的所有商品列表                        │
│  ┌────────┬────────┬────────┬────────┐                    │
│  │商品卡片│商品卡片│商品卡片│商品卡片│                    │
│  │ [图片] │ [图片] │ [图片] │ [图片] │                    │
│  │ 标题   │ 标题   │ 标题   │ 标题   │                    │
│  │ 评分   │ 评分   │ 评分   │ 评分   │                    │
│  │ 趋势图 │ 趋势图 │ 趋势图 │ 趋势图 │                    │
│  │ 价格   │ 价格   │ 价格   │ 价格   │                    │
│  │ [查看] │ [查看] │ [查看] │ [查看] │                    │
│  └────────┴────────┴────────┴────────┘                    │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 核心指标卡片（4个）

#### 卡片1：品牌规模
```typescript
{
  icon: "TrendingUp",
  iconColor: "bg-green-500",
  label: "品牌规模",
  value: string,               // 格式化金额
  subLabel: "年度销售额预估"
}
```

**数据来源**：`brand_monthly_trends` 表  
**计算公式**：
```sql
SELECT SUM(revenue)
FROM brand_monthly_trends
WHERE brand_id = ?
  AND date >= DATE_SUB(NOW(), INTERVAL 12 MONTH)
```

---

#### 卡片2：品牌增速
```typescript
{
  icon: "TrendingUp" | "TrendingDown",
  iconColor: "bg-green-500" | "bg-red-500",
  label: "品牌增速",
  value: string,               // "+20.3%" 或 "N/A"
  subLabel: "年复合增长率" | "数据不足"
}
```

**数据来源**：`brand_monthly_trends` 表  
**计算公式**：同CAGR计算

---

#### 卡片3：商品数量
```typescript
{
  icon: "Package",
  iconColor: "bg-pink-500",
  label: "商品数量",
  value: number,
  subLabel: "在售商品总数"
}
```

**数据来源**：`products` 表  
**计算公式**：`COUNT(asin) WHERE brand_id = ?`

---

#### 卡片4：品牌独立站
```typescript
{
  icon: "Globe",
  iconColor: "bg-blue-500",
  label: "品牌独立站",
  button: {
    label: "访问网站",
    icon: "ExternalLink",
    url: string | null,
    disabled: boolean          // 无网站时禁用
  } | {
    label: "暂无数据",
    disabled: true
  }
}
```

**数据来源**：`brands` 表的 `website` 字段

---

### 3.3 品牌销售趋势图表

#### 图表配置
```typescript
{
  type: "LineChart",
  height: 300,
  data: Array<{
    date: string,              // 月份
    revenue: number            // 销售额
  }>,
  line: {
    dataKey: "revenue",
    stroke: "#3b82f6",
    strokeWidth: 3,
    dot: true,                 // 显示数据点
    activeDot: { r: 6 }
  },
  xAxis: {
    dataKey: "date",
    label: "月份"
  },
  yAxis: {
    label: "销售额"
  },
  tooltip: true,
  grid: true
}
```

**数据来源**：
```sql
SELECT date, revenue
FROM brand_monthly_trends
WHERE brand_id = ?
ORDER BY date ASC
```

---

### 3.4 社交媒体数据卡片（4列）

#### YouTube卡片
```typescript
{
  platform: "youtube",
  icon: "Youtube",
  iconColor: "text-red-600",
  bgColor: "bg-red-50",
  label: "YouTube",
  value: string,               // 格式化订阅数（如：50K）
  subLabel: "订阅用户",
  button: {
    label: "访问频道",
    icon: "ExternalLink",
    url: string | null
  },
  visible: boolean             // 有数据时才显示
}
```

---

#### Instagram卡片
```typescript
{
  platform: "instagram",
  icon: "Instagram",
  iconColor: "text-pink-600",
  bgColor: "bg-pink-50",
  label: "Instagram",
  value: string,               // 格式化粉丝数
  subLabel: "粉丝数量",
  button: {
    label: "访问主页",
    icon: "ExternalLink",
    url: string | null
  },
  visible: boolean
}
```

---

#### Facebook卡片
```typescript
{
  platform: "facebook",
  icon: "Facebook",
  iconColor: "text-blue-600",
  bgColor: "bg-blue-50",
  label: "Facebook",
  value: string,               // 格式化关注数
  subLabel: "关注数量",
  button: {
    label: "访问主页",
    icon: "ExternalLink",
    url: string | null
  },
  visible: boolean
}
```

---

#### Reddit卡片
```typescript
{
  platform: "reddit",
  icon: "MessageCircle",
  iconColor: "text-orange-600",
  bgColor: "bg-orange-50",
  label: "Reddit",
  value: string,               // 帖子数
  subLabel: "讨论数量",
  button: {
    label: "访问社区",
    icon: "ExternalLink",
    url: string | null
  },
  visible: boolean
}
```

**数据来源**：
```sql
SELECT 
  platform,
  url,
  subscribers,
  followers,
  posts
FROM brand_social_media
WHERE brand_id = ?
```

---

### 3.5 品牌商品网格

#### 响应式布局
- **桌面端**：每行4个商品
- **平板端**：每行3个商品
- **移动端**：每行1个商品

#### 商品卡片结构
```typescript
{
  // 商品图片
  image: {
    url: string,
    alt: string,
    placeholder: "/placeholder-product.png"  // 无图时显示
  },
  
  // 商品标题
  title: {
    text: string,
    maxLines: 2,               // 最多显示2行
    overflow: "ellipsis"
  },
  
  // 品牌名称
  brand: string,
  
  // 评分和评论
  rating: {
    score: number,             // 0-5
    stars: JSX.Element,        // 星级组件
    reviews: number,           // 评论数
    display: "4.5 (1,234)"
  },
  
  // 销售趋势（迷你图）
  salesTrend: {
    title: "销售趋势(6个月)",
    data: Array<{
      month: string,
      sales: number,
      opacity: number          // 渐变透明度
    }>,
    trend: {
      value: number,           // 趋势百分比
      direction: "up" | "down",
      icon: "ArrowUp" | "ArrowDown",
      color: "text-green-600" | "text-red-600"
    },
    chart: {
      type: "bar",
      height: 40,
      barWidth: 8,
      radius: [4, 4, 0, 0]
    }
  },
  
  // 价格和月销量
  price: {
    value: number,
    display: string,           // "$299.99"
    currency: "USD"
  },
  monthlySales: {
    value: number,
    display: string            // "1.2K/月"
  },
  
  // 操作按钮
  action: {
    label: "在亚马逊查看",
    icon: "ExternalLink",
    url: string,               // Amazon产品链接
    target: "_blank"
  }
}
```

#### 数据来源

**商品基础信息**：
```sql
SELECT 
  asin,
  title,
  brand,
  price,
  rating,
  reviews,
  image_url,
  monthly_sales
FROM products
WHERE brand_id = ?
ORDER BY monthly_sales DESC
```

**商品销售趋势（最近6个月）**：
```sql
SELECT 
  date,
  sales,
  units
FROM product_monthly_sales
WHERE asin = ?
ORDER BY date DESC
LIMIT 6
```

**趋势计算**：
```javascript
function calculateProductTrend(salesData) {
  if (salesData.length < 2) return null;
  
  const latest = salesData[0].sales;
  const oldest = salesData[salesData.length - 1].sales;
  
  if (oldest === 0) return null;
  
  return ((latest - oldest) / oldest) * 100;
}
```

---

## 🎨 UI/UX设计规范

### 配色方案
```css
/* 主色调 */
--primary: #2563eb;          /* 蓝色 */
--success: #10b981;          /* 绿色 */
--warning: #f59e0b;          /* 黄色 */
--danger: #ef4444;           /* 红色 */

/* 背景色 */
--bg-primary: #ffffff;       /* 白色 */
--bg-secondary: #f8fafc;     /* 浅灰 */
--bg-tertiary: #f1f5f9;      /* 更浅灰 */

/* 边框色 */
--border: #e2e8f0;

/* 文字色 */
--text-primary: #1e293b;     /* 深灰 */
--text-secondary: #64748b;   /* 中灰 */
--text-tertiary: #94a3b8;    /* 浅灰 */

/* 社交媒体色 */
--youtube: #ef4444;          /* 红色 */
--instagram: #ec4899;        /* 粉色 */
--facebook: #3b82f6;         /* 蓝色 */
--reddit: #f97316;           /* 橙色 */
```

### 组件规范

#### 卡片
```css
.card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 24px;
  transition: box-shadow 0.2s;
}

.card:hover {
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
```

#### 按钮
```css
.button-primary {
  background: #2563eb;
  color: white;
  border-radius: 6px;
  height: 40px;
  padding: 0 24px;
  font-weight: 500;
  transition: background 0.2s;
}

.button-primary:hover {
  background: #1d4ed8;
}
```

#### 表格
```css
.table {
  width: 100%;
  border-collapse: collapse;
}

.table thead {
  background: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
}

.table tbody tr:nth-child(even) {
  background: #f8fafc;
}

.table tbody tr:hover {
  background: #f1f5f9;
}

.table td, .table th {
  padding: 12px 16px;
  text-align: left;
}
```

### 图标规范
- **库**：Lucide React
- **大小**：16px（小）、20px（中）、24px（大）
- **颜色**：继承父元素或使用主题色

### 图表规范
- **库**：Recharts
- **配色**：使用主题色
- **交互**：Tooltip、Legend、Grid
- **响应式**：自适应容器宽度

### 响应式断点
```css
/* 移动端 */
@media (max-width: 767px) {
  /* 1列布局 */
}

/* 平板端 */
@media (min-width: 768px) and (max-width: 1279px) {
  /* 2列布局 */
}

/* 桌面端 */
@media (min-width: 1280px) {
  /* 4列布局 */
}
```

---

## ✅ 前端开发检查清单

### 首页
- [ ] 4个核心指标卡片正确显示
- [ ] 市场列表表格所有列正确显示
- [ ] 市场趋势预览图正确渲染
- [ ] 点击"查看详情"跳转到市场详情页
- [ ] 所有数据格式化正确（金额、数量）
- [ ] 响应式布局正常

### 市场详情页
- [ ] 面包屑导航正确
- [ ] 4个核心指标卡片正确显示
- [ ] 市场趋势双轴图表正确渲染
- [ ] 品牌销售额占比饼图正确显示Top 8
- [ ] 品牌销售额排名柱状图正确显示Top 8
- [ ] 品牌列表表格所有列正确显示
- [ ] 社交媒体图标和数据正确显示
- [ ] 点击"查看详情"跳转到品牌详情页
- [ ] CAGR为null时显示"N/A"
- [ ] 响应式布局正常

### 品牌详情页
- [ ] 面包屑导航正确
- [ ] 4个核心指标卡片正确显示
- [ ] 品牌独立站按钮正确（有/无网站）
- [ ] 品牌销售趋势折线图正确渲染
- [ ] 社交媒体4个卡片正确显示（有数据时）
- [ ] 社交媒体链接可点击跳转
- [ ] 商品网格响应式布局正常
- [ ] 商品卡片所有信息正确显示
- [ ] 商品销售趋势迷你图正确渲染
- [ ] 商品图片加载失败时显示占位图
- [ ] 点击"在亚马逊查看"跳转到Amazon
- [ ] 响应式布局正常

---

**前端展示需求文档完成！请继续阅读后台管理需求文档。**
