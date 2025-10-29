# BrandTrekin Beta

BrandTrekin是一个基于线上商业数据分析从而给品牌卖家进行数据支持与建议的分析平台。

## 功能特点

### 1. 市场分析
- 追踪多个细分市场（CNC Router Machine、Laser Engraver）
- 展示市场规模、增速、搜索量等核心指标
- 可视化市场趋势和变化

### 2. 品牌分析
- 品牌销售数据追踪
- 社交媒体数据整合（YouTube、Instagram、Facebook、Reddit）
- 品牌产品组合分析
- 品牌竞争力对比

### 3. 产品分析
- 产品基础信息展示
- 销售趋势可视化
- 评分和评论数据
- 直接链接到亚马逊产品页

## 技术栈

- **前端框架**: Next.js 14 (React 18)
- **语言**: TypeScript
- **图表库**: Recharts
- **图标**: Lucide React
- **样式**: CSS Modules
- **数据解析**: XLSX, PapaParse

## 项目结构

```
/project
├── data/                    # 原始数据文件
│   ├── CNCRouter/          # CNC路由器市场数据
│   │   ├── Brand-Social.xlsx
│   │   ├── GKW.csv
│   │   ├── KeywordHistory.xlsx
│   │   ├── product-US-sales.xlsx
│   │   └── Product-US.xlsx
│   └── LaserEngraver/      # 激光雕刻机市场数据
│       └── ...
├── public/
│   └── data/               # 解析后的JSON数据
│       ├── markets.json
│       ├── cnc-router.json
│       └── laser-engraver.json
├── scripts/
│   └── parseData.js        # 数据解析脚本
├── src/
│   ├── app/
│   │   ├── layout.tsx      # 全局布局
│   │   ├── page.tsx        # 首页（市场列表）
│   │   ├── globals.css     # 全局样式
│   │   └── market/
│   │       └── [id]/
│   │           ├── page.tsx           # 市场详情页
│   │           └── brand/
│   │               └── [brand]/
│   │                   └── page.tsx   # 品牌详情页
│   └── types/
│       └── index.ts        # TypeScript类型定义
└── package.json
```

## 数据模型

### 1. 细分市场（Sub-market）
- 市场名称
- 市场规模（年度销售额）
- 市场增速（CAGR）
- 搜索量（Google + Amazon）
- 品牌数量
- 商品数量
- 月度趋势数据

### 2. 商品（Product）
- ASIN（唯一标识）
- 标题、品牌、价格
- 评分、评论数
- 月销量
- 图片URL
- 月度销售数据

### 3. 关键词（Keyword）
- Google关键词（GKW.csv）
- Amazon关键词（KeywordHistory.xlsx）
- 月度搜索量历史数据

### 4. 品牌（Brand）
- 品牌名称
- 年度销售规模
- 增长率
- 产品数量
- 独立站地址
- 社交媒体数据（YouTube、Instagram、Facebook、Reddit）

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 解析数据

首次运行或数据更新后，需要运行数据解析脚本：

```bash
node scripts/parseData.js
```

这将读取 `data/` 文件夹中的Excel和CSV文件，并生成JSON文件到 `public/data/` 目录。

### 3. 启动开发服务器

```bash
npm run dev
```

应用将在 http://localhost:3000 启动

### 4. 构建生产版本

```bash
npm run build
npm start
```

## 页面说明

### 首页 (`/`)
- 展示所有市场的概览
- 市场核心指标统计
- 市场列表表格
- 市场趋势预览图

### 市场详情页 (`/market/[id]`)
- 市场核心指标卡片
- 市场趋势图（双轴：销售额 + 搜索量）
- 品牌销售额占比饼图
- 品牌销售额排名柱状图
- 品牌列表表格（包含社交媒体链接）

### 品牌详情页 (`/market/[id]/brand/[brand]`)
- 品牌核心指标
- 品牌销售趋势图（24个月）
- 社交媒体数据卡片
- 品牌产品网格展示

## 数据来源

所有数据存储在 `data/` 文件夹中，按市场分类：

- `Product-US.xlsx`: 产品基础信息
- `product-US-sales.xlsx`: 产品月度销售数据
- `GKW.csv`: Google关键词搜索数据
- `KeywordHistory.xlsx`: Amazon关键词历史数据
- `Brand-Social.xlsx`: 品牌社交媒体数据

## 数据更新流程

1. 将新的Excel/CSV文件放入相应的 `data/[市场名称]/` 文件夹
2. 运行 `node scripts/parseData.js` 重新解析数据
3. 刷新浏览器查看更新后的数据

## 添加新市场

1. 在 `data/` 文件夹下创建新的市场文件夹
2. 添加所需的数据文件（参考现有市场的文件结构）
3. 在 `scripts/parseData.js` 中的 `markets` 数组添加新市场配置：

```javascript
const markets = [
  { id: 'cnc-router', name: 'CNC Router Machine', folder: 'CNCRouter' },
  { id: 'laser-engraver', name: 'Laser Engraver', folder: 'LaserEngraver' },
  { id: 'new-market', name: 'New Market Name', folder: 'NewMarketFolder' }
];
```

4. 运行数据解析脚本

## 设计风格

界面设计参考了 SimilarWeb 的现代化设计风格：
- 简洁的卡片式布局
- 清晰的数据可视化
- 响应式设计
- 直观的导航和面包屑
- 统一的配色方案

## 浏览器支持

- Chrome (推荐)
- Firefox
- Safari
- Edge

## 许可证

本项目为内部使用，未公开授权。

## 开发者

BrandTrekin Team
