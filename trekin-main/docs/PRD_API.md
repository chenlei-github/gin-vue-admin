# 🔌 API接口设计文档

> **本文档详细描述前端展示API和后台管理API的接口规范**

---

## 📋 API规范

### 基础信息
- **Base URL**: `/api`
- **Content-Type**: `application/json`
- **字符编码**: `UTF-8`

### 统一响应格式
```typescript
interface ApiResponse<T> {
  code: number;        // 0表示成功，非0表示失败
  message: string;     // 响应消息
  data?: T;            // 响应数据
}
```

### 错误码
| 错误码 | 说明 |
|-------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 1️⃣ 前端展示API

### API 1.1：获取市场列表

**接口地址**：`GET /api/markets`

**请求参数**：无

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "cnc-router-machine",
      "name": "CNC Router Machine",
      "metrics": {
        "totalRevenue": 13800000,
        "totalProducts": 114,
        "brandCount": 25,
        "searchVolume": 1500000,
        "cagr": 15.5,
        "monthlyTrends": [
          {
            "date": "2023-01",
            "revenue": 1100000
          },
          {
            "date": "2023-02",
            "revenue": 1150000
          }
        ]
      }
    }
  ]
}
```

---

### API 1.2：获取市场详情

**接口地址**：`GET /api/markets/:id`

**路径参数**：
- `id`: 市场slug（如：cnc-router-machine）

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "cnc-router-machine",
    "name": "CNC Router Machine",
    "metrics": {
      "totalRevenue": 13800000,
      "totalProducts": 114,
      "brandCount": 25,
      "searchVolume": 1500000,
      "cagr": 15.5,
      "monthlyTrends": [
        {
          "date": "2023-01",
          "revenue": 1100000,
          "searchVolume": 1400000
        }
      ]
    },
    "brands": [
      {
        "brand": "Genmitsu",
        "totalRevenue": 5200000,
        "productCount": 15,
        "cagr": 20.3,
        "website": "https://genmitsu.com",
        "social": {
          "youtube": {
            "url": "https://youtube.com/@genmitsu",
            "subscribers": 50000
          },
          "instagram": {
            "url": "https://instagram.com/genmitsu",
            "followers": 30000
          }
        }
      }
    ]
  }
}
```

---

### API 1.3：获取品牌详情

**接口地址**：`GET /api/markets/:marketId/brands/:brandName`

**路径参数**：
- `marketId`: 市场slug
- `brandName`: 品牌名称

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "brand": "Genmitsu",
    "totalRevenue": 5200000,
    "productCount": 15,
    "cagr": 20.3,
    "website": "https://genmitsu.com",
    "social": {
      "youtube": {
        "url": "https://youtube.com/@genmitsu",
        "subscribers": 50000
      }
    },
    "monthlyTrends": [
      {
        "date": "2023-01",
        "revenue": 430000
      }
    ],
    "products": [
      {
        "asin": "B08XYZ123",
        "title": "CNC Router Machine...",
        "price": 299.99,
        "rating": 4.5,
        "reviews": 1234,
        "imageUrl": "https://...",
        "monthlySales": 500,
        "salesTrend": [
          {
            "date": "2023-01",
            "sales": 14000
          }
        ]
      }
    ]
  }
}
```

---

## 2️⃣ 后台管理API

### API 2.1：市场管理

#### 获取市场列表
**接口地址**：`GET /api/admin/markets`

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "marketName": "CNC Router Machine",
      "marketSlug": "cnc-router-machine",
      "status": "active",
      "createdAt": "2025-10-28 10:00:00"
    }
  ]
}
```

#### 创建市场
**接口地址**：`POST /api/admin/markets`

**请求体**：
```json
{
  "marketName": "CNC Router Machine",
  "marketSlug": "cnc-router-machine",
  "description": "...",
  "status": "active"
}
```

#### 更新市场
**接口地址**：`PUT /api/admin/markets/:id`

#### 删除市场
**接口地址**：`DELETE /api/admin/markets/:id`

**请求体**：
```json
{
  "confirmName": "CNC Router Machine"
}
```

---

### API 2.2：数据导入

#### 上传文件
**接口地址**：`POST /api/admin/markets/:id/upload`

**Content-Type**：`multipart/form-data`

**请求参数**：
- `fileType`: brand_social | gkw | keyword_history | product_us | product_sales
- `file`: File

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "fileId": "xxx",
    "preview": [
      {
        "Brand": "Genmitsu",
        "Website": "https://..."
      }
    ]
  }
}
```

#### 开始导入
**接口地址**：`POST /api/admin/markets/:id/import`

**请求体**：
```json
{
  "fileIds": {
    "brand_social": "xxx",
    "gkw": "xxx",
    "keyword_history": "xxx",
    "product_us": "xxx",
    "product_sales": "xxx"
  },
  "importMode": "incremental",
  "skipInvalid": true,
  "autoCreateBrand": true
}
```

#### 获取导入进度
**接口地址**：`GET /api/admin/markets/:id/import/:importId/progress`

**响应示例**：
```json
{
  "code": 0,
  "data": {
    "status": "processing",
    "progress": 65,
    "currentStep": "导入商品数据",
    "logs": [
      "[2025-10-28 10:30:01] 开始解析文件...",
      "[2025-10-28 10:30:02] 解析成功"
    ]
  }
}
```

---

## ✅ API开发检查清单

### 前端展示API
- [ ] 获取市场列表API正常
- [ ] 获取市场详情API正常
- [ ] 获取品牌详情API正常
- [ ] 所有数据格式化正确
- [ ] 错误处理完善

### 后台管理API
- [ ] 市场CRUD API正常
- [ ] 文件上传API正常
- [ ] 数据导入API正常
- [ ] 导入进度API正常
- [ ] 导入历史API正常

---

**API接口设计文档完成！**
