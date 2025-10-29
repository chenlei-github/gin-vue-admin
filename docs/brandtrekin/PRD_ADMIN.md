# 🔧 后台管理需求文档

> **本文档详细描述后台管理系统的所有功能模块、数据导入流程和解析逻辑**

---

## 📋 目录

1. [市场管理模块](#1-市场管理模块)
2. [数据导入模块](#2-数据导入模块)
3. [数据解析逻辑](#3-数据解析逻辑)
4. [数据导入流程](#4-数据导入流程)
5. [导入历史模块](#5-导入历史模块)

---

## 1️⃣ 市场管理模块

### 1.1 市场列表页

#### 页面路径
```
/admin/markets
```

#### 页面布局
```
┌─────────────────────────────────────────────────────────────┐
│  [页面标题]  市场管理                    [+ 添加市场]按钮    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [市场列表表格]                                              │
│  ┌────┬──────────┬──────────┬────────┬──────────┬────────┐ │
│  │ ID │市场名称  │市场ID    │状态    │创建时间  │操作    │ │
│  ├────┼──────────┼──────────┼────────┼──────────┼────────┤ │
│  │ 1  │CNC...    │cnc-...   │[开关]  │2025-...  │[操作]  │ │
│  │ 2  │Laser...  │laser-... │[开关]  │2025-...  │[操作]  │ │
│  └────┴──────────┴──────────┴────────┴──────────┴────────┘ │
└─────────────────────────────────────────────────────────────┘
```

#### 表格列定义

| 列名 | 字段 | 类型 | 说明 |
|-----|------|------|------|
| ID | id | Integer | 自增主键 |
| 市场名称 | market_name | String | 中文名称，如"CNC Router Machine" |
| 市场ID | market_slug | String | 英文slug，用于URL，如"cnc-router-machine" |
| 状态 | status | Switch | 启用/禁用开关，实时切换 |
| 创建时间 | created_at | DateTime | 格式：YYYY-MM-DD HH:mm:ss |
| 操作 | - | Actions | [编辑] [删除] [导入数据] 三个按钮 |

#### 操作按钮

**编辑按钮**
```typescript
{
  label: "编辑",
  icon: "Edit",
  action: () => router.push(`/admin/markets/edit/${id}`)
}
```

**删除按钮**
```typescript
{
  label: "删除",
  icon: "Trash2",
  color: "danger",
  action: () => showDeleteConfirmModal(id)
}
```

**导入数据按钮**
```typescript
{
  label: "导入数据",
  icon: "Upload",
  color: "primary",
  action: () => router.push(`/admin/markets/${id}/import`)
}
```

---

### 1.2 添加市场页

#### 页面路径
```
/admin/markets/create
```

#### 页面布局
```
┌─────────────────────────────────────────────────────────────┐
│  [面包屑导航]  市场管理 > 添加市场                           │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [表单区域]                                                  │
│                                                              │
│  市场名称 *                                                  │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ CNC Router Machine                                     │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                              │
│  市场ID (slug) *                                             │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ cnc-router-machine                    [自动生成]按钮   │ │
│  └────────────────────────────────────────────────────────┘ │
│  提示：用于URL，只能包含小写字母、数字和连字符               │
│                                                              │
│  市场描述                                                    │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                                                        │ │
│  │                                                        │ │
│  │                                                        │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                              │
│  状态                                                        │
│  ○ 启用  ○ 禁用                                             │
│                                                              │
│  ┌────────────┬────────────────────┬────────────┐          │
│  │ [取消]     │ [保存并导入数据]   │ [保存]     │          │
│  └────────────┴────────────────────┴────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

#### 表单字段

**市场名称** (必填)
```typescript
{
  name: "market_name",
  type: "text",
  label: "市场名称",
  placeholder: "如：CNC Router Machine",
  required: true,
  maxLength: 100,
  validation: {
    required: "市场名称不能为空",
    maxLength: "市场名称不能超过100个字符"
  },
  onChange: (value) => {
    // 自动生成slug
    autoGenerateSlug(value);
  }
}
```

**市场ID (slug)** (必填)
```typescript
{
  name: "market_slug",
  type: "text",
  label: "市场ID (slug)",
  placeholder: "如：cnc-router-machine",
  required: true,
  maxLength: 100,
  pattern: /^[a-z0-9-]+$/,
  validation: {
    required: "市场ID不能为空",
    pattern: "只能包含小写字母、数字和连字符",
    unique: "该市场ID已存在"  // 需要异步校验
  },
  helperText: "用于URL，只能包含小写字母、数字和连字符",
  actions: [
    {
      label: "自动生成",
      action: () => generateSlugFromName()
    }
  ]
}
```

**自动生成Slug逻辑**
```javascript
function generateSlug(name) {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9\s-]/g, '')  // 移除特殊字符
    .replace(/\s+/g, '-')           // 空格替换为连字符
    .replace(/-+/g, '-')            // 多个连字符合并为一个
    .trim();
}

// 示例
generateSlug("CNC Router Machine") // => "cnc-router-machine"
generateSlug("Laser Engraver & Cutter") // => "laser-engraver-cutter"
```

**市场描述** (可选)
```typescript
{
  name: "description",
  type: "textarea",
  label: "市场描述",
  placeholder: "简要描述该市场...",
  rows: 4,
  maxLength: 500
}
```

**状态** (必填)
```typescript
{
  name: "status",
  type: "radio",
  label: "状态",
  options: [
    { label: "启用", value: "active" },
    { label: "禁用", value: "inactive" }
  ],
  defaultValue: "active"
}
```

#### 按钮操作

**取消按钮**
```typescript
{
  label: "取消",
  variant: "outline",
  action: () => router.back()
}
```

**保存按钮**
```typescript
{
  label: "保存",
  variant: "primary",
  action: async () => {
    const result = await createMarket(formData);
    if (result.success) {
      showSuccessMessage("市场创建成功");
      router.push("/admin/markets");
    }
  }
}
```

**保存并导入数据按钮**
```typescript
{
  label: "保存并导入数据",
  variant: "primary",
  action: async () => {
    const result = await createMarket(formData);
    if (result.success) {
      router.push(`/admin/markets/${result.data.id}/import`);
    }
  }
}
```

---

### 1.3 编辑市场页

#### 页面路径
```
/admin/markets/edit/:id
```

#### 页面布局
与添加市场页相同，但：
1. 表单字段预填充现有数据
2. 市场ID (slug) 字段禁用编辑（防止破坏现有链接）
3. 页面标题改为"编辑市场"

#### 特殊说明
```typescript
{
  name: "market_slug",
  type: "text",
  label: "市场ID (slug)",
  disabled: true,  // 禁用编辑
  helperText: "市场ID创建后不可修改，以保持URL稳定性"
}
```

---

### 1.4 删除市场功能

#### 触发方式
点击市场列表中的[删除]按钮

#### 确认弹窗
```
┌─────────────────────────────────────────────────────────────┐
│  ⚠️  确认删除                                                │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  删除市场将同时删除该市场下的所有数据，包括：                │
│  • 所有品牌信息                                              │
│  • 所有商品信息                                              │
│  • 所有关键词数据                                            │
│  • 所有销售数据                                              │
│  • 所有趋势数据                                              │
│                                                              │
│  ⚠️ 此操作不可恢复！                                         │
│                                                              │
│  请输入市场名称以确认删除：                                  │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                                                        │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                              │
│  ┌────────────┬────────────┐                               │
│  │ [取消]     │ [确认删除] │                               │
│  └────────────┴────────────┘                               │
└─────────────────────────────────────────────────────────────┘
```

#### 删除逻辑
```javascript
async function deleteMarket(marketId, confirmName) {
  // 1. 获取市场信息
  const market = await getMarketById(marketId);
  
  // 2. 校验输入的名称
  if (confirmName !== market.market_name) {
    throw new Error("市场名称不匹配");
  }
  
  // 3. 开启数据库事务
  const transaction = await db.beginTransaction();
  
  try {
    // 4. 按顺序删除关联数据（从子表到父表）
    await transaction.execute(`
      DELETE FROM product_monthly_sales 
      WHERE asin IN (
        SELECT asin FROM products WHERE market_id = ?
      )
    `, [marketId]);
    
    await transaction.execute(`
      DELETE FROM products WHERE market_id = ?
    `, [marketId]);
    
    await transaction.execute(`
      DELETE FROM brand_monthly_trends 
      WHERE brand_id IN (
        SELECT id FROM brands WHERE market_id = ?
      )
    `, [marketId]);
    
    await transaction.execute(`
      DELETE FROM brand_social_media 
      WHERE brand_id IN (
        SELECT id FROM brands WHERE market_id = ?
      )
    `, [marketId]);
    
    await transaction.execute(`
      DELETE FROM brands WHERE market_id = ?
    `, [marketId]);
    
    await transaction.execute(`
      DELETE FROM keyword_monthly_volume 
      WHERE keyword_id IN (
        SELECT id FROM keywords WHERE market_id = ?
      )
    `, [marketId]);
    
    await transaction.execute(`
      DELETE FROM keywords WHERE market_id = ?
    `, [marketId]);
    
    await transaction.execute(`
      DELETE FROM market_monthly_trends WHERE market_id = ?
    `, [marketId]);
    
    await transaction.execute(`
      DELETE FROM import_logs WHERE market_id = ?
    `, [marketId]);
    
    await transaction.execute(`
      DELETE FROM markets WHERE id = ?
    `, [marketId]);
    
    // 5. 提交事务
    await transaction.commit();
    
    return { success: true, message: "市场删除成功" };
    
  } catch (error) {
    // 6. 回滚事务
    await transaction.rollback();
    throw error;
  }
}
```

---

## 2️⃣ 数据导入模块

### 2.1 数据导入页

#### 页面路径
```
/admin/markets/:id/import
```

#### 页面布局
```
┌─────────────────────────────────────────────────────────────┐
│  [面包屑导航]  市场管理 > {市场名称} > 数据导入              │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [页面标题]  导入市场数据 - {市场名称}                       │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [上传组件1] 品牌社交媒体数据                                │
│  📄 Brand-Social.xlsx                                        │
│  包含品牌名称、独立站、YouTube、Instagram、Facebook、Reddit  │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ [选择文件]                              状态：未上传    │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [上传组件2] Google关键词数据                                │
│  📄 GKW.csv                                                  │
│  包含Google关键词及月度搜索量历史数据                        │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ [选择文件]                              状态：未上传    │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [上传组件3] Amazon关键词历史数据                            │
│  📄 KeywordHistory.xlsx                                      │
│  包含Amazon关键词及月度搜索量历史数据                        │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ [选择文件]                              状态：未上传    │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [上传组件4] 商品基础信息                                    │
│  📄 Product-US.xlsx                                          │
│  包含ASIN、标题、品牌、价格、评分、评论数、图片URL等         │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ [选择文件]                              状态：未上传    │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [上传组件5] 商品月度销售数据                                │
│  📄 product-US-sales.xlsx                                    │
│  包含ASIN及每月销售额、销量数据                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ [选择文件]                              状态：未上传    │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [导入选项]                                                  │
│                                                              │
│  导入模式：                                                  │
│  ● 增量导入（保留现有数据，仅添加或更新）                   │
│  ○ 全量替换（删除所有数据后重新导入）                       │
│                                                              │
│  数据校验：                                                  │
│  ☑ 跳过无效数据行                                           │
│  ☑ 自动创建不存在的品牌                                     │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [导入进度区域]                                              │
│  ████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░  65%          │
│  当前步骤：导入商品数据...                                   │
│                                                              │
│  [日志输出]                                                  │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ [2025-10-28 10:30:01] 开始解析文件...                  │ │
│  │ [2025-10-28 10:30:02] 解析Brand-Social.xlsx成功，25行  │ │
│  │ [2025-10-28 10:30:03] 解析GKW.csv成功，150行           │ │
│  │ [2025-10-28 10:30:05] 开始导入品牌数据...              │ │
│  │ [2025-10-28 10:30:06] 成功导入25个品牌                 │ │
│  │ ...                                                    │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  ┌────────────┬────────────────────┬────────────┐          │
│  │ [取消]     │ [查看导入历史]     │ [开始导入] │          │
│  └────────────┴────────────────────┴────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 文件上传组件

#### 组件状态
```typescript
type UploadStatus = 
  | "idle"           // 未上传
  | "uploading"      // 上传中
  | "parsing"        // 解析中
  | "success"        // 解析成功
  | "error";         // 解析失败

interface UploadComponent {
  fileType: "brand_social" | "gkw" | "keyword_history" | "product_us" | "product_sales";
  title: string;
  fileName: string;
  description: string;
  acceptedTypes: string[];
  status: UploadStatus;
  file: File | null;
  preview: any[] | null;      // 前5行数据预览
  error: string | null;
}
```

#### 5个上传组件配置

**组件1：品牌社交媒体数据**
```typescript
{
  fileType: "brand_social",
  title: "品牌社交媒体数据",
  fileName: "Brand-Social.xlsx",
  description: "包含品牌名称、独立站、YouTube、Instagram、Facebook、Reddit数据",
  acceptedTypes: [".xlsx"],
  expectedColumns: [
    "Brand",                    // 必填
    "Website",                  // 可选
    "YouTube",                  // 可选
    "YouTube Subscribers",      // 可选
    "Instagram",                // 可选
    "Instagram Followers",      // 可选
    "Facebook",                 // 可选
    "Facebook Followers",       // 可选
    "Reddit",                   // 可选
    "Reddit Posts"              // 可选
  ]
}
```

**组件2：Google关键词数据**
```typescript
{
  fileType: "gkw",
  title: "Google关键词数据",
  fileName: "GKW.csv",
  description: "包含Google关键词及月度搜索量历史数据",
  acceptedTypes: [".csv"],
  expectedColumns: [
    "Keyword",                  // 必填
    "YYYY-MM",                  // 动态列（月份）
    // ... 更多月份列
  ]
}
```

**组件3：Amazon关键词历史数据**
```typescript
{
  fileType: "keyword_history",
  title: "Amazon关键词历史数据",
  fileName: "KeywordHistory.xlsx",
  description: "包含Amazon关键词及月度搜索量历史数据",
  acceptedTypes: [".xlsx"],
  expectedColumns: [
    "Keyword",                  // 必填
    "YYYY-MM",                  // 动态列（月份）
    // ... 更多月份列
  ]
}
```

**组件4：商品基础信息**
```typescript
{
  fileType: "product_us",
  title: "商品基础信息",
  fileName: "Product-US.xlsx",
  description: "包含ASIN、标题、品牌、价格、评分、评论数、图片URL等",
  acceptedTypes: [".xlsx"],
  expectedColumns: [
    "ASIN",                     // 必填
    "Title",                    // 必填
    "Brand",                    // 必填
    "Price",                    // 可选
    "Rating",                   // 可选
    "Reviews",                  // 可选
    "Image",                    // 可选
    "Monthly Sales"             // 可选
  ]
}
```

**组件5：商品月度销售数据**
```typescript
{
  fileType: "product_sales",
  title: "商品月度销售数据",
  fileName: "product-US-sales.xlsx",
  description: "包含ASIN及每月销售额、销量数据",
  acceptedTypes: [".xlsx"],
  expectedColumns: [
    "ASIN",                     // 必填
    "YYYY-MM",                  // 动态列（月份）
    // ... 更多月份列
  ]
}
```

#### 上传流程
```javascript
async function handleFileUpload(fileType, file) {
  // 1. 更新状态为"上传中"
  setStatus(fileType, "uploading");
  
  // 2. 上传文件到服务器
  const uploadResult = await uploadFile(file);
  
  if (!uploadResult.success) {
    setStatus(fileType, "error");
    setError(fileType, uploadResult.error);
    return;
  }
  
  // 3. 更新状态为"解析中"
  setStatus(fileType, "parsing");
  
  // 4. 解析文件
  const parseResult = await parseFile(fileType, uploadResult.fileId);
  
  if (!parseResult.success) {
    setStatus(fileType, "error");
    setError(fileType, parseResult.error);
    return;
  }
  
  // 5. 更新状态为"成功"，显示预览
  setStatus(fileType, "success");
  setPreview(fileType, parseResult.preview);  // 前5行数据
}
```

#### 数据预览
```
┌─────────────────────────────────────────────────────────────┐
│  ✅ 解析成功！共 25 行数据                                   │
│                                                              │
│  [数据预览（前5行）]                                         │
│  ┌────────────┬────────────┬────────────┬────────────┐     │
│  │ Brand      │ Website    │ YouTube    │ YouTube... │     │
│  ├────────────┼────────────┼────────────┼────────────┤     │
│  │ Genmitsu   │ https://...│ https://...│ 50000      │     │
│  │ Ortur      │ https://...│ https://...│ 30000      │     │
│  │ ...        │ ...        │ ...        │ ...        │     │
│  └────────────┴────────────┴────────────┴────────────┘     │
│                                                              │
│  [重新上传]                                                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 3️⃣ 数据解析逻辑

### 3.1 解析Brand-Social.xlsx

#### 期望的Excel结构
```
| Brand      | Website         | YouTube         | YouTube Subscribers | Instagram       | Instagram Followers | Facebook        | Facebook Followers | Reddit          | Reddit Posts |
|------------|-----------------|-----------------|---------------------|-----------------|---------------------|-----------------|--------------------|-----------------|--------------| 
| Genmitsu   | https://...     | https://...     | 50000               | https://...     | 30000               | https://...     | 20000              | https://...     | 500          |
| Ortur      | https://...     | https://...     | 30000               | https://...     | 15000               | https://...     | 10000              | https://...     | 200          |
```

#### 解析代码
```javascript
import * as XLSX from 'xlsx';

async function parseBrandSocial(file, marketId) {
  // 1. 读取Excel文件
  const workbook = XLSX.read(await file.arrayBuffer());
  const sheetName = workbook.SheetNames[0];
  const sheet = workbook.Sheets[sheetName];
  
  // 2. 转换为JSON
  const data = XLSX.utils.sheet_to_json(sheet);
  
  // 3. 校验必填列
  if (data.length === 0) {
    throw new Error("文件为空");
  }
  
  const firstRow = data[0];
  if (!firstRow.Brand) {
    throw new Error("缺少必填列：Brand");
  }
  
  // 4. 解析每一行
  const brands = [];
  const errors = [];
  
  for (let i = 0; i < data.length; i++) {
    const row = data[i];
    const rowNumber = i + 2;  // Excel行号（从2开始，因为第1行是表头）
    
    // 跳过空行
    if (!row.Brand || row.Brand.trim() === '') {
      continue;
    }
    
    try {
      brands.push({
        market_id: marketId,
        brand_name: row.Brand.trim(),
        website: row.Website || null,
        
        // YouTube数据
        youtube_url: row.YouTube || null,
        youtube_subscribers: parseInt(row['YouTube Subscribers']) || 0,
        
        // Instagram数据
        instagram_url: row.Instagram || null,
        instagram_followers: parseInt(row['Instagram Followers']) || 0,
        
        // Facebook数据
        facebook_url: row.Facebook || null,
        facebook_followers: parseInt(row['Facebook Followers']) || 0,
        
        // Reddit数据
        reddit_url: row.Reddit || null,
        reddit_posts: parseInt(row['Reddit Posts']) || 0
      });
    } catch (error) {
      errors.push({
        row: rowNumber,
        error: error.message
      });
    }
  }
  
  return {
    success: true,
    data: brands,
    total: brands.length,
    errors: errors,
    preview: brands.slice(0, 5)  // 前5行预览
  };
}
```

---

### 3.2 解析GKW.csv

#### 期望的CSV结构
```
Keyword,2023-01,2023-02,2023-03,2023-04,...
cnc router,120000,125000,130000,135000,...
laser engraver,80000,82000,85000,88000,...
```

#### 解析代码
```javascript
import Papa from 'papaparse';

async function parseGKW(file, marketId) {
  return new Promise((resolve, reject) => {
    Papa.parse(file, {
      header: true,
      skipEmptyLines: true,
      complete: (results) => {
        try {
          const keywords = [];
          const errors = [];
          
          for (let i = 0; i < results.data.length; i++) {
            const row = results.data[i];
            const rowNumber = i + 2;
            
            // 跳过没有关键词的行
            if (!row.Keyword || row.Keyword.trim() === '') {
              continue;
            }
            
            const keyword = row.Keyword.trim();
            const monthlyVolumes = [];
            
            // 解析每个月份列
            for (const [key, value] of Object.entries(row)) {
              if (key === 'Keyword') continue;
              
              // 校验日期格式（YYYY-MM）
              if (!/^\d{4}-\d{2}$/.test(key)) {
                errors.push({
                  row: rowNumber,
                  error: `无效的日期格式：${key}`
                });
                continue;
              }
              
              // 解析搜索量
              const volume = parseInt(value);
              if (!isNaN(volume) && volume > 0) {
                monthlyVolumes.push({
                  date: `${key}-01`,  // 转换为完整日期格式
                  volume: volume
                });
              }
            }
            
            if (monthlyVolumes.length > 0) {
              keywords.push({
                market_id: marketId,
                keyword: keyword,
                source: 'google',
                monthly_volumes: monthlyVolumes
              });
            }
          }
          
          resolve({
            success: true,
            data: keywords,
            total: keywords.length,
            errors: errors,
            preview: keywords.slice(0, 5)
          });
          
        } catch (error) {
          reject(error);
        }
      },
      error: (error) => {
        reject(error);
      }
    });
  });
}
```

---

### 3.3 解析KeywordHistory.xlsx

#### 期望的Excel结构
```
| Keyword        | 2023-01 | 2023-02 | 2023-03 | ... |
|----------------|---------|---------|---------|-----|
| cnc router     | 150000  | 155000  | 160000  | ... |
| laser engraver | 100000  | 105000  | 110000  | ... |
```

#### 解析代码
```javascript
async function parseKeywordHistory(file, marketId) {
  // 与parseGKW类似，但source为'amazon'
  const workbook = XLSX.read(await file.arrayBuffer());
  const sheetName = workbook.SheetNames[0];
  const sheet = workbook.Sheets[sheetName];
  const data = XLSX.utils.sheet_to_json(sheet);
  
  const keywords = [];
  const errors = [];
  
  for (let i = 0; i < data.length; i++) {
    const row = data[i];
    const rowNumber = i + 2;
    
    if (!row.Keyword || row.Keyword.trim() === '') {
      continue;
    }
    
    const keyword = row.Keyword.trim();
    const monthlyVolumes = [];
    
    for (const [key, value] of Object.entries(row)) {
      if (key === 'Keyword') continue;
      
      if (!/^\d{4}-\d{2}$/.test(key)) {
        errors.push({
          row: rowNumber,
          error: `无效的日期格式：${key}`
        });
        continue;
      }
      
      const volume = parseInt(value);
      if (!isNaN(volume) && volume > 0) {
        monthlyVolumes.push({
          date: `${key}-01`,
          volume: volume
        });
      }
    }
    
    if (monthlyVolumes.length > 0) {
      keywords.push({
        market_id: marketId,
        keyword: keyword,
        source: 'amazon',  // 注意：这里是amazon
        monthly_volumes: monthlyVolumes
      });
    }
  }
  
  return {
    success: true,
    data: keywords,
    total: keywords.length,
    errors: errors,
    preview: keywords.slice(0, 5)
  };
}
```

---

### 3.4 解析Product-US.xlsx

#### 期望的Excel结构
```
| ASIN       | Title                  | Brand    | Price  | Rating | Reviews | Image                | Monthly Sales |
|------------|------------------------|----------|--------|--------|---------|----------------------|---------------|
| B08XYZ123  | CNC Router Machine...  | Genmitsu | 299.99 | 4.5    | 1234    | https://...          | 500           |
| B09ABC456  | Laser Engraver...      | Ortur    | 199.99 | 4.7    | 2345    | https://...          | 800           |
```

#### 解析代码
```javascript
async function parseProductUS(file, marketId) {
  const workbook = XLSX.read(await file.arrayBuffer());
  const sheetName = workbook.SheetNames[0];
  const sheet = workbook.Sheets[sheetName];
  const data = XLSX.utils.sheet_to_json(sheet);
  
  // 校验必填列
  if (data.length === 0) {
    throw new Error("文件为空");
  }
  
  const firstRow = data[0];
  const requiredColumns = ['ASIN', 'Title', 'Brand'];
  for (const col of requiredColumns) {
    if (!firstRow.hasOwnProperty(col)) {
      throw new Error(`缺少必填列：${col}`);
    }
  }
  
  const products = [];
  const errors = [];
  
  for (let i = 0; i < data.length; i++) {
    const row = data[i];
    const rowNumber = i + 2;
    
    // 校验必填字段
    if (!row.ASIN || !row.Title || !row.Brand) {
      errors.push({
        row: rowNumber,
        error: "缺少必填字段（ASIN、Title或Brand）"
      });
      continue;
    }
    
    try {
      products.push({
        market_id: marketId,
        asin: row.ASIN.trim(),
        title: row.Title.trim(),
        brand: row.Brand.trim(),
        price: parseFloat(row.Price) || 0,
        rating: parseFloat(row.Rating) || 0,
        reviews: parseInt(row.Reviews) || 0,
        image_url: row.Image || null,
        monthly_sales: parseInt(row['Monthly Sales']) || 0
      });
    } catch (error) {
      errors.push({
        row: rowNumber,
        error: error.message
      });
    }
  }
  
  return {
    success: true,
    data: products,
    total: products.length,
    errors: errors,
    preview: products.slice(0, 5)
  };
}
```

---

### 3.5 解析product-US-sales.xlsx

#### 期望的Excel结构
```
| ASIN       | 2023-01 | 2023-02 | 2023-03 | ... |
|------------|---------|---------|---------|-----|
| B08XYZ123  | 15000   | 16000   | 17000   | ... |
| B09ABC456  | 25000   | 26000   | 27000   | ... |
```

#### 解析代码
```javascript
async function parseProductSales(file, marketId) {
  const workbook = XLSX.read(await file.arrayBuffer());
  const sheetName = workbook.SheetNames[0];
  const sheet = workbook.Sheets[sheetName];
  const data = XLSX.utils.sheet_to_json(sheet);
  
  const productSales = [];
  const errors = [];
  
  for (let i = 0; i < data.length; i++) {
    const row = data[i];
    const rowNumber = i + 2;
    
    if (!row.ASIN || row.ASIN.trim() === '') {
      continue;
    }
    
    const asin = row.ASIN.trim();
    const monthlySales = [];
    
    for (const [key, value] of Object.entries(row)) {
      if (key === 'ASIN') continue;
      
      if (!/^\d{4}-\d{2}$/.test(key)) {
        errors.push({
          row: rowNumber,
          error: `无效的日期格式：${key}`
        });
        continue;
      }
      
      const sales = parseFloat(value);
      if (!isNaN(sales) && sales > 0) {
        monthlySales.push({
          date: `${key}-01`,
          sales: sales,
          units: 0  // 如果有销量数据，可以从另一个sheet读取
        });
      }
    }
    
    if (monthlySales.length > 0) {
      productSales.push({
        asin: asin,
        monthly_sales: monthlySales
      });
    }
  }
  
  return {
    success: true,
    data: productSales,
    total: productSales.length,
    errors: errors,
    preview: productSales.slice(0, 5)
  };
}
```

---

## 4️⃣ 数据导入流程

### 4.1 导入流程图

```
开始导入
    ↓
[Step 1] 校验所有文件已上传
    ↓
[Step 2] 开启数据库事务
    ↓
[Step 3] 如果是"全量替换"模式，删除现有数据
    ↓
[Step 4] 导入品牌基础信息
    ├─ 插入brands表
    └─ 插入brand_social_media表
    ↓
[Step 5] 导入商品基础信息
    ├─ 查找或创建brand_id
    └─ 插入products表
    ↓
[Step 6] 导入关键词数据
    ├─ 插入keywords表
    └─ 插入keyword_monthly_volume表
    ↓
[Step 7] 导入商品销售数据
    └─ 插入product_monthly_sales表
    ↓
[Step 8] 计算聚合指标
    ├─ 计算brand_monthly_trends
    ├─ 计算market_monthly_trends
    ├─ 计算CAGR
    └─ 更新markets表统计数据
    ↓
[Step 9] 提交事务
    ↓
[Step 10] 记录导入日志
    ↓
完成导入
```

### 4.2 导入代码实现

```javascript
async function importMarketData(marketId, files, options) {
  const {
    importMode,           // 'incremental' | 'replace'
    skipInvalid,          // boolean
    autoCreateBrand       // boolean
  } = options;
  
  const log = [];
  const stats = {
    brands: 0,
    products: 0,
    keywords: 0,
    salesRecords: 0,
    skipped: 0
  };
  
  try {
    log.push(`[${timestamp()}] 开始导入数据...`);
    
    // Step 1: 校验文件
    log.push(`[${timestamp()}] 校验文件...`);
    const requiredFiles = ['brand_social', 'gkw', 'keyword_history', 'product_us', 'product_sales'];
    for (const fileType of requiredFiles) {
      if (!files[fileType]) {
        throw new Error(`缺少必需文件：${fileType}`);
      }
    }
    
    // Step 2: 开启事务
    log.push(`[${timestamp()}] 开启数据库事务...`);
    const transaction = await db.beginTransaction();
    
    try {
      // Step 3: 全量替换模式 - 删除现有数据
      if (importMode === 'replace') {
        log.push(`[${timestamp()}] 删除现有数据...`);
        await deleteMarketData(transaction, marketId);
      }
      
      // Step 4: 导入品牌数据
      log.push(`[${timestamp()}] 导入品牌数据...`);
      const brandResult = await importBrands(
        transaction,
        marketId,
        files.brand_social,
        { skipInvalid }
      );
      stats.brands = brandResult.imported;
      stats.skipped += brandResult.skipped;
      log.push(`[${timestamp()}] 成功导入 ${brandResult.imported} 个品牌`);
      
      // Step 5: 导入商品数据
      log.push(`[${timestamp()}] 导入商品数据...`);
      const productResult = await importProducts(
        transaction,
        marketId,
        files.product_us,
        { skipInvalid, autoCreateBrand }
      );
      stats.products = productResult.imported;
      stats.skipped += productResult.skipped;
      log.push(`[${timestamp()}] 成功导入 ${productResult.imported} 个商品`);
      
      // Step 6: 导入关键词数据
      log.push(`[${timestamp()}] 导入关键词数据...`);
      const keywordResult = await importKeywords(
        transaction,
        marketId,
        files.gkw,
        files.keyword_history,
        { skipInvalid }
      );
      stats.keywords = keywordResult.imported;
      stats.skipped += keywordResult.skipped;
      log.push(`[${timestamp()}] 成功导入 ${keywordResult.imported} 个关键词`);
      
      // Step 7: 导入销售数据
      log.push(`[${timestamp()}] 导入销售数据...`);
      const salesResult = await importSalesData(
        transaction,
        files.product_sales,
        { skipInvalid }
      );
      stats.salesRecords = salesResult.imported;
      stats.skipped += salesResult.skipped;
      log.push(`[${timestamp()}] 成功导入 ${salesResult.imported} 条销售记录`);
      
      // Step 8: 计算聚合指标
      log.push(`[${timestamp()}] 计算聚合指标...`);
      await calculateAggregates(transaction, marketId);
      log.push(`[${timestamp()}] 聚合指标计算完成`);
      
      // Step 9: 提交事务
      log.push(`[${timestamp()}] 提交事务...`);
      await transaction.commit();
      log.push(`[${timestamp()}] 数据导入成功！`);
      
      // Step 10: 记录导入日志
      await saveImportLog(marketId, {
        importMode,
        status: 'success',
        stats,
        log: log.join('\n')
      });
      
      return {
        success: true,
        stats,
        log
      };
      
    } catch (error) {
      // 回滚事务
      log.push(`[${timestamp()}] 错误：${error.message}`);
      log.push(`[${timestamp()}] 回滚事务...`);
      await transaction.rollback();
      
      // 记录失败日志
      await saveImportLog(marketId, {
        importMode,
        status: 'failed',
        stats,
        error: error.message,
        log: log.join('\n')
      });
      
      throw error;
    }
    
  } catch (error) {
    return {
      success: false,
      error: error.message,
      log
    };
  }
}
```

### 4.3 计算聚合指标

```javascript
async function calculateAggregates(transaction, marketId) {
  // 1. 计算品牌月度趋势
  await transaction.execute(`
    INSERT INTO brand_monthly_trends (brand_id, date, revenue)
    SELECT 
      p.brand_id,
      pms.date,
      SUM(pms.sales) as revenue
    FROM product_monthly_sales pms
    JOIN products p ON pms.asin = p.asin
    WHERE p.market_id = ?
    GROUP BY p.brand_id, pms.date
    ON DUPLICATE KEY UPDATE revenue = VALUES(revenue)
  `, [marketId]);
  
  // 2. 计算市场月度趋势（销售额）
  await transaction.execute(`
    INSERT INTO market_monthly_trends (market_id, date, revenue)
    SELECT 
      b.market_id,
      bmt.date,
      SUM(bmt.revenue) as revenue
    FROM brand_monthly_trends bmt
    JOIN brands b ON bmt.brand_id = b.id
    WHERE b.market_id = ?
    GROUP BY b.market_id, bmt.date
    ON DUPLICATE KEY UPDATE revenue = VALUES(revenue)
  `, [marketId]);
  
  // 3. 计算市场月度趋势（搜索量）
  await transaction.execute(`
    UPDATE market_monthly_trends mmt
    SET search_volume = (
      SELECT COALESCE(SUM(kmv.volume), 0)
      FROM keyword_monthly_volume kmv
      JOIN keywords k ON kmv.keyword_id = k.id
      WHERE k.market_id = mmt.market_id
        AND kmv.date = mmt.date
    )
    WHERE mmt.market_id = ?
  `, [marketId]);
  
  // 4. 计算品牌CAGR
  await transaction.execute(`
    UPDATE brands b
    SET cagr = (
      SELECT 
        CASE 
          WHEN COUNT(*) >= 12 AND MIN(revenue) > 0 THEN
            LEAST(999, GREATEST(-99, 
              (POWER(MAX(revenue) / MIN(revenue), 1.0 / (COUNT(*) / 12.0)) - 1) * 100
            ))
          ELSE NULL
        END
      FROM brand_monthly_trends
      WHERE brand_id = b.id
    )
    WHERE b.market_id = ?
  `, [marketId]);
  
  // 5. 计算市场CAGR
  await transaction.execute(`
    UPDATE markets m
    SET cagr = (
      SELECT 
        CASE 
          WHEN COUNT(*) >= 12 AND MIN(revenue) > 0 THEN
            LEAST(999, GREATEST(-99, 
              (POWER(MAX(revenue) / MIN(revenue), 1.0 / (COUNT(*) / 12.0)) - 1) * 100
            ))
          ELSE NULL
        END
      FROM market_monthly_trends
      WHERE market_id = m.id
    )
    WHERE m.id = ?
  `, [marketId]);
  
  // 6. 更新品牌总销售额
  await transaction.execute(`
    UPDATE brands b
    SET total_revenue = (
      SELECT COALESCE(SUM(revenue), 0)
      FROM brand_monthly_trends
      WHERE brand_id = b.id
        AND date >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)
    )
    WHERE b.market_id = ?
  `, [marketId]);
  
  // 7. 更新品牌商品数量
  await transaction.execute(`
    UPDATE brands b
    SET product_count = (
      SELECT COUNT(*)
      FROM products
      WHERE brand_id = b.id
    )
    WHERE b.market_id = ?
  `, [marketId]);
  
  // 8. 更新市场总销售额
  await transaction.execute(`
    UPDATE markets m
    SET total_revenue = (
      SELECT COALESCE(SUM(revenue), 0)
      FROM market_monthly_trends
      WHERE market_id = m.id
        AND date >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)
    )
    WHERE m.id = ?
  `, [marketId]);
  
  // 9. 更新市场商品总数
  await transaction.execute(`
    UPDATE markets m
    SET total_products = (
      SELECT COUNT(*)
      FROM products
      WHERE market_id = m.id
    )
    WHERE m.id = ?
  `, [marketId]);
  
  // 10. 更新市场品牌数量
  await transaction.execute(`
    UPDATE markets m
    SET brand_count = (
      SELECT COUNT(*)
      FROM brands
      WHERE market_id = m.id
    )
    WHERE m.id = ?
  `, [marketId]);
  
  // 11. 更新市场搜索量
  await transaction.execute(`
    UPDATE markets m
    SET search_volume = (
      SELECT COALESCE(SUM(kmv.volume), 0)
      FROM keyword_monthly_volume kmv
      JOIN keywords k ON kmv.keyword_id = k.id
      WHERE k.market_id = m.id
        AND kmv.date = (
          SELECT MAX(date) FROM keyword_monthly_volume
        )
    )
    WHERE m.id = ?
  `, [marketId]);
}
```

---

## 5️⃣ 导入历史模块

### 5.1 导入历史列表页

#### 页面路径
```
/admin/markets/:id/import-history
```

#### 页面布局
```
┌─────────────────────────────────────────────────────────────┐
│  [面包屑导航]  市场管理 > {市场名称} > 导入历史              │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [页面标题]  导入历史 - {市场名称}                           │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  [导入历史列表表格]                                          │
│  ┌──────────┬────────┬────────┬────────┬────────┬────────┐ │
│  │导入时间  │操作人  │导入模式│状态    │统计信息│操作    │ │
│  ├──────────┼────────┼────────┼────────┼────────┼────────┤ │
│  │2025-...  │admin   │增量    │✅成功  │品牌25..│[详情]  │ │
│  │2025-...  │admin   │全量    │❌失败  │品牌0... │[日志]  │ │
│  └──────────┴────────┴────────┴────────┴────────┴────────┘ │
└─────────────────────────────────────────────────────────────┘
```

#### 表格列定义

| 列名 | 字段 | 类型 | 说明 |
|-----|------|------|------|
| 导入时间 | created_at | DateTime | YYYY-MM-DD HH:mm:ss |
| 操作人 | created_by | String | 管理员用户名 |
| 导入模式 | import_mode | Badge | "增量导入" 或 "全量替换" |
| 状态 | status | Badge | "成功"（绿色）/ "失败"（红色）/ "部分成功"（黄色）|
| 统计信息 | - | Text | "品牌X个，商品X个，关键词X个" |
| 操作 | - | Actions | [查看详情] [查看日志] |

---

## ✅ 后台管理开发检查清单

### 市场管理
- [ ] 市场列表页正确显示所有市场
- [ ] 状态开关可以实时切换
- [ ] 添加市场表单校验正确
- [ ] 市场ID自动生成功能正常
- [ ] 市场ID唯一性校验正常
- [ ] 编辑市场功能正常
- [ ] 删除市场需要输入名称确认
- [ ] 删除市场级联删除所有关联数据
- [ ] 删除操作使用事务，失败时回滚

### 数据导入
- [ ] 5个文件上传组件正常工作
- [ ] 文件类型校验正确
- [ ] 文件解析逻辑正确
- [ ] 解析成功后显示前5行预览
- [ ] 解析失败显示错误信息
- [ ] 导入模式选择正常
- [ ] 导入选项勾选正常
- [ ] 导入进度实时显示
- [ ] 导入日志实时输出
- [ ] 导入成功显示统计信息
- [ ] 导入失败回滚所有数据
- [ ] 聚合指标计算正确
- [ ] CAGR计算正确

### 导入历史
- [ ] 导入历史列表正确显示
- [ ] 可以查看导入详情
- [ ] 可以查看导入日志
- [ ] 状态显示正确（成功/失败/部分成功）

---

**后台管理需求文档完成！请继续阅读数据库设计文档。**
