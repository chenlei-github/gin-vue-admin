# 📋 BrandTrekin 品牌市场分析平台 - GVA版本

> **基于 Gin-Vue-Admin 框架的品牌市场分析平台**  
> **版本**: v2.0 (GVA迁移版)  
> **更新时间**: 2025-10-29

---

## 🎯 项目概述

**项目名称**：BrandTrekin - 品牌市场分析平台  
**项目类型**：全栈Web应用（基于Gin-Vue-Admin框架）  
**开发目标**：将原有基于Next.js的系统迁移到Gin-Vue-Admin框架，实现企业级的品牌市场数据管理和分析

### 核心功能模块

1. **市场管理模块**
   - 市场基础信息管理（市场名称、市场slug、描述、状态）
   - 市场聚合统计（总销售额、商品总数、品牌数量、搜索量、CAGR年复合增长率）
   - 支持市场的增删改查

2. **品牌管理模块**
   - 品牌基础信息（品牌名称、品牌官网）
   - 品牌聚合统计（总销售额、商品数量、CAGR）
   - 品牌社交媒体数据（YouTube、Instagram、Facebook、Reddit）

3. **商品管理模块**
   - 商品基础信息（ASIN、标题、价格、评分、评论数、月销量、图片）
   - 商品月度销售数据（每月的销售额和销量）

4. **关键词管理模块**
   - 关键词基础信息（关键词、来源：google/amazon）
   - 关键词月度搜索量数据

5. **趋势数据模块**
   - 品牌月度趋势（每月销售额）
   - 市场月度趋势（每月销售额和搜索量）
   - 用于计算CAGR和展示趋势图表

6. **数据导入模块**
   - 支持导入5种类型的Excel/CSV文件
   - 支持增量导入和全量替换两种模式
   - 导入过程记录日志

7. **导入历史模块**
   - 记录每次导入的详细信息
   - 支持查看历史导入记录和日志

---

## 🛠 技术栈

### 后端技术
```json
{
  "framework": "Gin",
  "language": "Golang 1.22+",
  "orm": "GORM",
  "database": "MySQL 8.0+",
  "auth": "JWT + Casbin",
  "file_parser": {
    "excel": "excelize",
    "csv": "encoding/csv"
  }
}
```

### 前端技术
```json
{
  "framework": "Vue 3",
  "ui_library": "Element Plus",
  "build_tool": "Vite",
  "state_management": "Pinia",
  "router": "Vue Router 4",
  "http_client": "Axios",
  "charts": "ECharts"
}
```

---

## 📁 项目结构

### 后端目录结构
```
server/
├── api/v1/brandtrekin/          # API控制器
│   ├── btmarket.go              # 市场管理API
│   ├── btbrand.go               # 品牌管理API
│   ├── btbrandsocialmedia.go    # 品牌社交媒体API
│   ├── btproduct.go             # 商品管理API
│   ├── btproductmonthlysales.go # 商品月度销售API
│   ├── btkeyword.go             # 关键词管理API
│   ├── btkeywordmonthlyvolume.go# 关键词月度搜索量API
│   ├── btbrandmonthlytrend.go   # 品牌月度趋势API
│   ├── btmarketmonthlytrend.go  # 市场月度趋势API
│   └── btimportlog.go           # 导入日志API
├── model/brandtrekin/           # 数据模型
│   ├── btmarket.go
│   ├── btbrand.go
│   ├── btbrandsocialmedia.go
│   ├── btproduct.go
│   ├── btproductmonthlysales.go
│   ├── btkeyword.go
│   ├── btkeywordmonthlyvolume.go
│   ├── btbrandmonthlytrend.go
│   ├── btmarketmonthlytrend.go
│   ├── btimportlog.go
│   ├── request/                 # 请求参数结构体
│   └── response/                # 响应结构体
├── service/brandtrekin/         # 业务逻辑层
│   ├── btmarket.go
│   ├── btbrand.go
│   ├── btbrandsocialmedia.go
│   ├── btproduct.go
│   ├── btproductmonthlysales.go
│   ├── btkeyword.go
│   ├── btkeywordmonthlyvolume.go
│   ├── btbrandmonthlytrend.go
│   ├── btmarketmonthlytrend.go
│   └── btimportlog.go
└── router/brandtrekin/          # 路由定义
    ├── btmarket.go
    ├── btbrand.go
    ├── btbrandsocialmedia.go
    ├── btproduct.go
    ├── btproductmonthlysales.go
    ├── btkeyword.go
    ├── btkeywordmonthlyvolume.go
    ├── btbrandmonthlytrend.go
    ├── btmarketmonthlytrend.go
    └── btimportlog.go
```

### 前端目录结构
```
web/src/
├── api/brandtrekin/             # API接口定义
│   ├── btmarket.js
│   ├── btbrand.js
│   ├── btbrandsocialmedia.js
│   ├── btproduct.js
│   ├── btproductmonthlysales.js
│   ├── btkeyword.js
│   ├── btkeywordmonthlyvolume.js
│   ├── btbrandmonthlytrend.js
│   ├── btmarketmonthlytrend.js
│   └── btimportlog.js
└── view/brandtrekin/            # 页面视图
    ├── btmarket.vue
    ├── btbrand.vue
    ├── btbrandsocialmedia.vue
    ├── btproduct.vue
    ├── btproductmonthlysales.vue
    ├── btkeyword.vue
    ├── btkeywordmonthlyvolume.vue
    ├── btbrandmonthlytrend.vue
    ├── btmarketmonthlytrend.vue
    └── btimportlog.vue
```

---

## 💾 数据库设计

### 数据表清单

| 序号 | 表名 | 说明 | 预估数据量 |
|-----|------|------|-----------| 
| 1 | bt_markets | 市场表 | 10-100 |
| 2 | bt_brands | 品牌表 | 100-1000 |
| 3 | bt_brand_social_media | 品牌社交媒体表 | 400-4000 |
| 4 | bt_products | 商品表 | 1000-10000 |
| 5 | bt_product_monthly_sales | 商品月度销售表 | 10万-100万 |
| 6 | bt_keywords | 关键词表 | 1000-10000 |
| 7 | bt_keyword_monthly_volume | 关键词月度搜索量表 | 10万-100万 |
| 8 | bt_brand_monthly_trends | 品牌月度趋势表 | 1万-10万 |
| 9 | bt_market_monthly_trends | 市场月度趋势表 | 1000-10000 |
| 10 | bt_import_logs | 导入日志表 | 100-1000 |

### 数据关系图

```
bt_markets (市场)
    ├── 1:N → bt_brands (品牌)
    │         ├── 1:N → bt_brand_social_media (品牌社交媒体)
    │         ├── 1:N → bt_products (商品)
    │         └── 1:N → bt_brand_monthly_trends (品牌月度趋势)
    ├── 1:N → bt_products (商品)
    │         └── 1:N → bt_product_monthly_sales (商品月度销售)
    ├── 1:N → bt_keywords (关键词)
    │         └── 1:N → bt_keyword_monthly_volume (关键词月度搜索量)
    ├── 1:N → bt_market_monthly_trends (市场月度趋势)
    └── 1:N → bt_import_logs (导入日志)
```

---

## 📊 核心数据指标

### 市场级别指标

| 指标名称 | 计算公式 | 数据来源 |
|---------|---------| ---------|
| 市场规模 | 最近12个月销售额总和 | bt_market_monthly_trends |
| 市场增速(CAGR) | (Ending/Beginning)^(1/years)-1 | bt_market_monthly_trends |
| 市场声量 | 最近月份搜索量总和 | bt_keyword_monthly_volume |
| 品牌数量 | COUNT(DISTINCT brand_id) | bt_brands |
| 商品数量 | COUNT(asin) | bt_products |

### 品牌级别指标

| 指标名称 | 计算公式 | 数据来源 |
|---------|---------| ---------|
| 品牌规模 | 最近12个月销售额总和 | bt_brand_monthly_trends |
| 品牌增速(CAGR) | (Ending/Beginning)^(1/years)-1 | bt_brand_monthly_trends |
| 商品数量 | COUNT(asin) | bt_products |
| 社交媒体数据 | 直接读取 | bt_brand_social_media |

### CAGR计算公式

```
CAGR = (Ending Value / Beginning Value)^(1/years) - 1

其中：
- Ending Value: 最近月份的销售额
- Beginning Value: 最早月份的销售额
- years: 时间跨度（月数 / 12）
- 需要至少12个月的数据才能计算

示例：
- 12个月前销售额：$100,000
- 最近月份销售额：$150,000
- CAGR = (150000 / 100000)^(1/1) - 1 = 0.5 = 50%
```

---

## 🎨 字典配置

系统已自动创建以下字典：

| 字典类型 | 字典名称 | 说明 | 选项 |
|---------|---------|------|------|
| market_status | 市场状态 | 市场的启用/禁用状态 | 启用(active)、禁用(inactive) |
| social_platform | 社交媒体平台 | 品牌社交媒体平台类型 | YouTube、Instagram、Facebook、Reddit |
| keyword_source | 关键词来源 | 关键词数据来源平台 | Google、Amazon |
| import_mode | 导入模式 | 数据导入模式类型 | 增量导入、全量替换 |
| import_status | 导入状态 | 数据导入执行状态 | 成功、失败、部分成功 |

---

## 🚀 快速开始

### 第一步：启动后端服务

```bash
# 1. 进入后端目录
cd server

# 2. 配置数据库连接
# 编辑 config.yaml，配置MySQL连接信息

# 3. 初始化数据库
# 系统会自动创建所有表结构

# 4. 启动服务
go run main.go
```

### 第二步：启动前端服务

```bash
# 1. 进入前端目录
cd web

# 2. 安装依赖
npm install

# 3. 启动开发服务器
npm run serve
```

### 第三步：访问系统

- **前端地址**：http://localhost:8080
- **后端地址**：http://localhost:8888
- **默认账号**：admin / 123456

---

## 📝 功能清单

### ✅ 已完成功能

#### 基础模块（10个）
- [x] 市场管理模块
- [x] 品牌管理模块
- [x] 品牌社交媒体模块
- [x] 商品管理模块
- [x] 商品月度销售模块
- [x] 关键词管理模块
- [x] 关键词月度搜索量模块
- [x] 品牌月度趋势模块
- [x] 市场月度趋势模块
- [x] 导入日志模块

#### 自动生成内容
- [x] 所有模块的CRUD API接口
- [x] 所有模块的前端管理页面
- [x] 所有模块的路由配置
- [x] 所有模块的权限配置
- [x] 5个业务字典配置

### 🔄 待开发功能

#### 数据导入功能
- [ ] Excel/CSV文件上传接口
- [ ] Brand-Social.xlsx解析逻辑
- [ ] GKW.csv解析逻辑
- [ ] KeywordHistory.xlsx解析逻辑
- [ ] Product-US.xlsx解析逻辑
- [ ] product-US-sales.xlsx解析逻辑
- [ ] 数据导入主流程
- [ ] 导入进度跟踪
- [ ] 导入日志记录

#### 数据聚合计算
- [ ] 品牌月度趋势聚合
- [ ] 市场月度趋势聚合
- [ ] CAGR自动计算
- [ ] 市场规模自动计算
- [ ] 品牌规模自动计算
- [ ] 搜索量自动汇总

#### 前端展示页面
- [ ] 市场列表页（展示核心指标和趋势）
- [ ] 市场详情页（详细数据、品牌分布、趋势图表）
- [ ] 品牌详情页（销售趋势、社交媒体、商品列表）
- [ ] 数据可视化图表（ECharts集成）

---

## 📚 文档索引

本项目文档分为以下几个部分：

### 1️⃣ [数据库设计文档](./PRD_DATABASE.md)
- 10张核心数据表的完整设计
- 表关系和外键约束
- 索引设计
- 核心计算逻辑（CAGR、聚合指标）

### 2️⃣ [API接口设计文档](./PRD_API.md)
- 前端展示API（市场列表、市场详情、品牌详情）
- 后台管理API（CRUD接口）
- 数据导入API
- 统一响应格式和错误码

### 3️⃣ [前端展示需求文档](./PRD_FRONTEND.md)
- 市场列表页完整布局和数据指标
- 市场详情页完整布局和数据指标
- 品牌详情页完整布局和数据指标
- 所有数据指标的计算公式
- UI/UX设计规范

### 4️⃣ [后台管理需求文档](./PRD_ADMIN.md)
- 市场管理模块（CRUD操作）
- 数据导入模块（5种文件类型）
- 数据解析逻辑（详细的解析规则）
- 导入流程和进度显示
- 导入历史记录

### 5️⃣ [开发任务清单](./PRD_TASKS.md)
- Phase 1: 数据库设计与搭建 ✅
- Phase 2: 基础模块开发 ✅
- Phase 3: 数据导入功能开发 🔄
- Phase 4: 数据聚合计算开发 🔄
- Phase 5: 前端展示页面开发 🔄
- Phase 6: 测试与优化 ⏳
- Phase 7: 部署上线 ⏳

---

## ⚠️ 重要注意事项

### 数据完整性
1. **不要遗漏任何指标**：前端展示的每一个数据指标都必须在数据库中有对应字段
2. **计算公式必须准确**：特别是CAGR、市场规模、搜索量等核心指标
3. **数据关联必须正确**：品牌、商品、关键词等数据必须正确关联到市场

### 数据导入
1. **5个文件缺一不可**：Brand-Social.xlsx、GKW.csv、KeywordHistory.xlsx、Product-US.xlsx、product-US-sales.xlsx
2. **解析逻辑必须严格**：按照文档中的解析规则，不能随意修改
3. **数据校验必须完善**：导入前必须校验数据格式和必填字段
4. **事务处理必须完整**：导入失败时必须回滚所有数据

### 性能优化
1. **大数据量处理**：商品销售数据可能有数万条，需要批量插入
2. **索引优化**：为常用查询字段添加索引
3. **缓存策略**：前端展示数据可以考虑缓存

---

## 🎓 开发指南

### GVA框架特性

#### 1. 自动生成的CRUD功能
每个模块都自动生成了完整的CRUD功能：
- 创建（Create）
- 查询（Read）- 支持分页、搜索、排序
- 更新（Update）
- 删除（Delete）

#### 2. 权限管理
- 所有API接口已自动注册到权限系统
- 可在系统管理界面配置角色权限
- 支持按钮级权限控制

#### 3. 数据关联
- 品牌关联市场（下拉选择）
- 商品关联市场和品牌（下拉选择）
- 关键词关联市场（下拉选择）
- 所有趋势数据关联对应实体

#### 4. 字典管理
- 5个业务字典已自动创建
- 可在系统管理界面维护字典选项
- 前端表单自动使用字典数据

### 扩展开发建议

#### 1. 数据导入功能开发
建议在 `service/brandtrekin` 目录下创建专门的导入服务：
```go
// service/brandtrekin/import_service.go
type ImportService struct{}

func (s *ImportService) ImportBrandSocial(file *multipart.FileHeader) error {
    // 解析Brand-Social.xlsx
    // 导入品牌和社交媒体数据
}

func (s *ImportService) ImportGKW(file *multipart.FileHeader) error {
    // 解析GKW.csv
    // 导入Google关键词数据
}

// ... 其他导入方法
```

#### 2. 数据聚合计算
建议创建定时任务或触发器：
```go
// service/brandtrekin/aggregate_service.go
type AggregateService struct{}

func (s *AggregateService) CalculateMarketMetrics(marketId uint) error {
    // 计算市场聚合指标
    // 更新 bt_markets 表的统计字段
}

func (s *AggregateService) CalculateCAGR(entityType string, entityId uint) (float64, error) {
    // 计算CAGR
}
```

#### 3. 前端展示页面
建议在 `web/src/view/brandtrekin` 目录下创建展示页面：
```
web/src/view/brandtrekin/
├── display/                    # 前端展示页面
│   ├── market-list.vue        # 市场列表页
│   ├── market-detail.vue      # 市场详情页
│   └── brand-detail.vue       # 品牌详情页
└── manage/                     # 后台管理页面（已自动生成）
    ├── btmarket.vue
    ├── btbrand.vue
    └── ...
```

---

## 📞 技术支持

### 参考资料
- **Gin-Vue-Admin官方文档**：https://www.gin-vue-admin.com
- **原始需求文档**：`trekin-main/docs/` 目录
- **GVA视频教程**：https://www.bilibili.com/video/BV1Rg411u7xH

### 常见问题

1. **Q: 如何添加自定义API接口？**  
   A: 在对应的 `api/v1/brandtrekin/*.go` 文件中添加方法，然后在 `router/brandtrekin/*.go` 中注册路由

2. **Q: 如何修改表单字段？**  
   A: 修改 `model/brandtrekin/*.go` 中的结构体定义，然后重新生成前端页面

3. **Q: 如何实现数据导入？**  
   A: 参考 `docs/brandtrekin/PRD_ADMIN.md` 中的数据导入模块设计

---

## 📝 版本历史

- **v2.0** (2025-10-29): 迁移到Gin-Vue-Admin框架，创建10个核心模块
- **v1.0** (2025-10-28): 原始Next.js版本

---

## ✅ 迁移完成检查清单

### 数据库
- [x] 所有10张表创建成功
- [x] 所有字段定义正确
- [x] 所有关联关系配置正确
- [x] 所有字典创建成功

### 后端
- [x] 所有模块的Model创建成功
- [x] 所有模块的API创建成功
- [x] 所有模块的Service创建成功
- [x] 所有模块的Router创建成功
- [x] 所有API权限自动注册

### 前端
- [x] 所有模块的管理页面创建成功
- [x] 所有模块的API接口定义创建成功
- [x] 所有菜单自动注册
- [ ] 前端展示页面（待开发）

### 功能
- [ ] 数据导入功能（待开发）
- [ ] 数据聚合计算（待开发）
- [ ] CAGR自动计算（待开发）
- [ ] 前端数据可视化（待开发）

---

**迁移基础框架已完成！接下来可以开始开发数据导入和展示功能。🚀**
