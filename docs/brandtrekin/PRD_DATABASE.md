# 💾 数据库设计文档

> **本文档详细描述MySQL数据库的表结构、索引、关系和计算逻辑**

---

## 📋 目录

1. [数据库概览](#1-数据库概览)
2. [核心数据表](#2-核心数据表)
3. [表关系图](#3-表关系图)
4. [核心计算逻辑](#4-核心计算逻辑)
5. [索引优化](#5-索引优化)

---

## 1️⃣ 数据库概览

### 数据库配置
```sql
CREATE DATABASE brandtrekin 
  CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;
```

### 表清单
| 序号 | 表名 | 说明 | 预估数据量 |
|-----|------|------|-----------|
| 1 | markets | 市场表 | 10-100 |
| 2 | brands | 品牌表 | 100-1000 |
| 3 | brand_social_media | 品牌社交媒体表 | 400-4000 |
| 4 | products | 商品表 | 1000-10000 |
| 5 | product_monthly_sales | 商品月度销售表 | 10万-100万 |
| 6 | keywords | 关键词表 | 1000-10000 |
| 7 | keyword_monthly_volume | 关键词月度搜索量表 | 10万-100万 |
| 8 | brand_monthly_trends | 品牌月度趋势表 | 1万-10万 |
| 9 | market_monthly_trends | 市场月度趋势表 | 1000-10000 |
| 10 | import_logs | 导入日志表 | 100-1000 |

---

## 2️⃣ 核心数据表

### 表1：markets（市场表）

#### 表结构
```sql
CREATE TABLE `markets` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '市场ID',
  `market_name` VARCHAR(100) NOT NULL COMMENT '市场名称',
  `market_slug` VARCHAR(100) NOT NULL COMMENT '市场slug（用于URL）',
  `description` TEXT COMMENT '市场描述',
  `status` ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
  
  -- 聚合统计字段（由计算任务更新）
  `total_revenue` DECIMAL(15,2) DEFAULT 0 COMMENT '总销售额（最近12个月）',
  `total_products` INT DEFAULT 0 COMMENT '商品总数',
  `brand_count` INT DEFAULT 0 COMMENT '品牌数量',
  `search_volume` BIGINT DEFAULT 0 COMMENT '搜索量（最近月份）',
  `cagr` DECIMAL(5,2) COMMENT 'CAGR（年复合增长率，-99.00 到 999.00）',
  
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_market_slug` (`market_slug`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='市场表';
```

#### 字段说明
| 字段名 | 类型 | 说明 | 示例值 |
|-------|------|------|--------|
| id | INT UNSIGNED | 自增主键 | 1 |
| market_name | VARCHAR(100) | 市场名称 | "CNC Router Machine" |
| market_slug | VARCHAR(100) | URL友好的标识符 | "cnc-router-machine" |
| description | TEXT | 市场描述 | "..." |
| status | ENUM | 状态：active/inactive | "active" |
| total_revenue | DECIMAL(15,2) | 最近12个月总销售额 | 13800000.00 |
| total_products | INT | 商品总数 | 114 |
| brand_count | INT | 品牌数量 | 25 |
| search_volume | BIGINT | 最近月份搜索量 | 1500000 |
| cagr | DECIMAL(5,2) | 年复合增长率 | 15.50 |

---

### 表2：brands（品牌表）

#### 表结构
```sql
CREATE TABLE `brands` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '品牌ID',
  `market_id` INT UNSIGNED NOT NULL COMMENT '市场ID',
  `brand_name` VARCHAR(100) NOT NULL COMMENT '品牌名称',
  `website` VARCHAR(500) COMMENT '品牌独立站',
  
  -- 聚合统计字段
  `total_revenue` DECIMAL(15,2) DEFAULT 0 COMMENT '总销售额（最近12个月）',
  `product_count` INT DEFAULT 0 COMMENT '商品数量',
  `cagr` DECIMAL(5,2) COMMENT 'CAGR',
  
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_market_brand` (`market_id`, `brand_name`),
  KEY `idx_market_id` (`market_id`),
  KEY `idx_brand_name` (`brand_name`),
  CONSTRAINT `fk_brands_market` FOREIGN KEY (`market_id`) 
    REFERENCES `markets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌表';
```

#### 字段说明
| 字段名 | 类型 | 说明 | 示例值 |
|-------|------|------|--------|
| id | INT UNSIGNED | 自增主键 | 1 |
| market_id | INT UNSIGNED | 所属市场ID | 1 |
| brand_name | VARCHAR(100) | 品牌名称 | "Genmitsu" |
| website | VARCHAR(500) | 品牌官网 | "https://genmitsu.com" |
| total_revenue | DECIMAL(15,2) | 最近12个月总销售额 | 5200000.00 |
| product_count | INT | 商品数量 | 15 |
| cagr | DECIMAL(5,2) | 年复合增长率 | 20.30 |

---

### 表3：brand_social_media（品牌社交媒体表）

#### 表结构
```sql
CREATE TABLE `brand_social_media` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `brand_id` INT UNSIGNED NOT NULL COMMENT '品牌ID',
  `platform` ENUM('youtube', 'instagram', 'facebook', 'reddit') NOT NULL COMMENT '平台',
  `url` VARCHAR(500) NOT NULL COMMENT '链接',
  `subscribers` INT DEFAULT 0 COMMENT '订阅数（YouTube）',
  `followers` INT DEFAULT 0 COMMENT '粉丝数（Instagram/Facebook）',
  `posts` INT DEFAULT 0 COMMENT '帖子数（Reddit）',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_brand_platform` (`brand_id`, `platform`),
  KEY `idx_brand_id` (`brand_id`),
  KEY `idx_platform` (`platform`),
  CONSTRAINT `fk_social_brand` FOREIGN KEY (`brand_id`) 
    REFERENCES `brands` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌社交媒体表';
```

#### 字段说明
| 字段名 | 类型 | 说明 | 示例值 |
|-------|------|------|--------|
| id | INT UNSIGNED | 自增主键 | 1 |
| brand_id | INT UNSIGNED | 品牌ID | 1 |
| platform | ENUM | 平台类型 | "youtube" |
| url | VARCHAR(500) | 平台链接 | "https://youtube.com/@genmitsu" |
| subscribers | INT | YouTube订阅数 | 50000 |
| followers | INT | Instagram/Facebook粉丝数 | 30000 |
| posts | INT | Reddit帖子数 | 500 |

---

### 表4：products（商品表）

#### 表结构
```sql
CREATE TABLE `products` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `market_id` INT UNSIGNED NOT NULL COMMENT '市场ID',
  `brand_id` INT UNSIGNED NOT NULL COMMENT '品牌ID',
  `asin` VARCHAR(20) NOT NULL COMMENT '亚马逊ASIN',
  `title` VARCHAR(500) NOT NULL COMMENT '商品标题',
  `price` DECIMAL(10,2) DEFAULT 0 COMMENT '价格',
  `rating` DECIMAL(3,2) DEFAULT 0 COMMENT '评分（0-5）',
  `reviews` INT DEFAULT 0 COMMENT '评论数',
  `monthly_sales` INT DEFAULT 0 COMMENT '月销量',
  `image_url` VARCHAR(500) COMMENT '图片URL',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_asin` (`asin`),
  KEY `idx_market_id` (`market_id`),
  KEY `idx_brand_id` (`brand_id`),
  KEY `idx_monthly_sales` (`monthly_sales`),
  CONSTRAINT `fk_products_market` FOREIGN KEY (`market_id`) 
    REFERENCES `markets` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_products_brand` FOREIGN KEY (`brand_id`) 
    REFERENCES `brands` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';
```

#### 字段说明
| 字段名 | 类型 | 说明 | 示例值 |
|-------|------|------|--------|
| id | INT UNSIGNED | 自增主键 | 1 |
| market_id | INT UNSIGNED | 市场ID | 1 |
| brand_id | INT UNSIGNED | 品牌ID | 1 |
| asin | VARCHAR(20) | 亚马逊ASIN | "B08XYZ123" |
| title | VARCHAR(500) | 商品标题 | "CNC Router Machine..." |
| price | DECIMAL(10,2) | 价格 | 299.99 |
| rating | DECIMAL(3,2) | 评分 | 4.50 |
| reviews | INT | 评论数 | 1234 |
| monthly_sales | INT | 月销量 | 500 |
| image_url | VARCHAR(500) | 图片URL | "https://..." |

---

### 表5：product_monthly_sales（商品月度销售表）

#### 表结构
```sql
CREATE TABLE `product_monthly_sales` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `asin` VARCHAR(20) NOT NULL COMMENT '商品ASIN',
  `date` DATE NOT NULL COMMENT '月份（YYYY-MM-01）',
  `sales` DECIMAL(12,2) DEFAULT 0 COMMENT '销售额',
  `units` INT DEFAULT 0 COMMENT '销量',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_asin_date` (`asin`, `date`),
  KEY `idx_asin` (`asin`),
  KEY `idx_date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品月度销售表';
```

#### 字段说明
| 字段名 | 类型 | 说明 | 示例值 |
|-------|------|------|--------|
| id | BIGINT UNSIGNED | 自增主键 | 1 |
| asin | VARCHAR(20) | 商品ASIN | "B08XYZ123" |
| date | DATE | 月份（每月1号） | "2023-01-01" |
| sales | DECIMAL(12,2) | 销售额 | 15000.00 |
| units | INT | 销量 | 50 |

---

### 表6：keywords（关键词表）

#### 表结构
```sql
CREATE TABLE `keywords` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `market_id` INT UNSIGNED NOT NULL COMMENT '市场ID',
  `keyword` VARCHAR(200) NOT NULL COMMENT '关键词',
  `source` ENUM('google', 'amazon') NOT NULL COMMENT '来源',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_market_keyword_source` (`market_id`, `keyword`, `source`),
  KEY `idx_market_id` (`market_id`),
  KEY `idx_keyword` (`keyword`),
  KEY `idx_source` (`source`),
  CONSTRAINT `fk_keywords_market` FOREIGN KEY (`market_id`) 
    REFERENCES `markets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='关键词表';
```

#### 字段说明
| 字段名 | 类型 | 说明 | 示例值 |
|-------|------|------|--------|
| id | INT UNSIGNED | 自增主键 | 1 |
| market_id | INT UNSIGNED | 市场ID | 1 |
| keyword | VARCHAR(200) | 关键词 | "cnc router" |
| source | ENUM | 来源：google/amazon | "google" |

---

### 表7：keyword_monthly_volume（关键词月度搜索量表）

#### 表结构
```sql
CREATE TABLE `keyword_monthly_volume` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `keyword_id` INT UNSIGNED NOT NULL COMMENT '关键词ID',
  `date` DATE NOT NULL COMMENT '月份（YYYY-MM-01）',
  `volume` BIGINT DEFAULT 0 COMMENT '搜索量',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_keyword_date` (`keyword_id`, `date`),
  KEY `idx_keyword_id` (`keyword_id`),
  KEY `idx_date` (`date`),
  CONSTRAINT `fk_volume_keyword` FOREIGN KEY (`keyword_id`) 
    REFERENCES `keywords` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='关键词月度搜索量表';
```

#### 字段说明
| 字段名 | 类型 | 说明 | 示例值 |
|-------|------|------|--------|
| id | BIGINT UNSIGNED | 自增主键 | 1 |
| keyword_id | INT UNSIGNED | 关键词ID | 1 |
| date | DATE | 月份（每月1号） | "2023-01-01" |
| volume | BIGINT | 搜索量 | 120000 |

---

### 表8：brand_monthly_trends（品牌月度趋势表）

#### 表结构
```sql
CREATE TABLE `brand_monthly_trends` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `brand_id` INT UNSIGNED NOT NULL COMMENT '品牌ID',
  `date` DATE NOT NULL COMMENT '月份（YYYY-MM-01）',
  `revenue` DECIMAL(12,2) DEFAULT 0 COMMENT '销售额',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_brand_date` (`brand_id`, `date`),
  KEY `idx_brand_id` (`brand_id`),
  KEY `idx_date` (`date`),
  CONSTRAINT `fk_trends_brand` FOREIGN KEY (`brand_id`) 
    REFERENCES `brands` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌月度趋势表';
```

#### 字段说明
| 字段名 | 类型 | 说明 | 示例值 |
|-------|------|------|--------|
| id | BIGINT UNSIGNED | 自增主键 | 1 |
| brand_id | INT UNSIGNED | 品牌ID | 1 |
| date | DATE | 月份（每月1号） | "2023-01-01" |
| revenue | DECIMAL(12,2) | 销售额 | 450000.00 |

---

### 表9：market_monthly_trends（市场月度趋势表）

#### 表结构
```sql
CREATE TABLE `market_monthly_trends` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `market_id` INT UNSIGNED NOT NULL COMMENT '市场ID',
  `date` DATE NOT NULL COMMENT '月份（YYYY-MM-01）',
  `revenue` DECIMAL(15,2) DEFAULT 0 COMMENT '销售额',
  `search_volume` BIGINT DEFAULT 0 COMMENT '搜索量',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_market_date` (`market_id`, `date`),
  KEY `idx_market_id` (`market_id`),
  KEY `idx_date` (`date`),
  CONSTRAINT `fk_trends_market` FOREIGN KEY (`market_id`) 
    REFERENCES `markets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='市场月度趋势表';
```

#### 字段说明
| 字段名 | 类型 | 说明 | 示例值 |
|-------|------|------|--------|
| id | BIGINT UNSIGNED | 自增主键 | 1 |
| market_id | INT UNSIGNED | 市场ID | 1 |
| date | DATE | 月份（每月1号） | "2023-01-01" |
| revenue | DECIMAL(15,2) | 销售额 | 1150000.00 |
| search_volume | BIGINT | 搜索量 | 1500000 |

---

### 表10：import_logs（导入日志表）

#### 表结构
```sql
CREATE TABLE `import_logs` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `market_id` INT UNSIGNED NOT NULL COMMENT '市场ID',
  `import_mode` ENUM('incremental', 'replace') NOT NULL COMMENT '导入模式',
  `status` ENUM('success', 'failed', 'partial') NOT NULL COMMENT '状态',
  `brands_count` INT DEFAULT 0 COMMENT '导入品牌数',
  `products_count` INT DEFAULT 0 COMMENT '导入商品数',
  `keywords_count` INT DEFAULT 0 COMMENT '导入关键词数',
  `sales_records_count` INT DEFAULT 0 COMMENT '导入销售记录数',
  `skipped_count` INT DEFAULT 0 COMMENT '跳过记录数',
  `error_message` TEXT COMMENT '错误信息',
  `log_content` LONGTEXT COMMENT '日志内容',
  `created_by` VARCHAR(50) COMMENT '操作人',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  KEY `idx_market_id` (`market_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_logs_market` FOREIGN KEY (`market_id`) 
    REFERENCES `markets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='导入日志表';
```

#### 字段说明
| 字段名 | 类型 | 说明 | 示例值 |
|-------|------|------|--------|
| id | INT UNSIGNED | 自增主键 | 1 |
| market_id | INT UNSIGNED | 市场ID | 1 |
| import_mode | ENUM | 导入模式 | "incremental" |
| status | ENUM | 状态 | "success" |
| brands_count | INT | 导入品牌数 | 25 |
| products_count | INT | 导入商品数 | 114 |
| keywords_count | INT | 导入关键词数 | 150 |
| sales_records_count | INT | 导入销售记录数 | 2736 |
| skipped_count | INT | 跳过记录数 | 5 |
| error_message | TEXT | 错误信息 | null |
| log_content | LONGTEXT | 日志内容 | "..." |
| created_by | VARCHAR(50) | 操作人 | "admin" |

---

## 3️⃣ 表关系图

```
┌─────────────────┐
│    markets      │
│  (市场表)       │
└────────┬────────┘
         │ 1
         │
         │ N
    ┌────┴────┬──────────────────────────┐
    │         │                          │
    │ N       │ N                        │ N
┌───┴────┐ ┌─┴──────────┐         ┌─────┴────────┐
│ brands │ │  products  │         │   keywords   │
│(品牌表)│ │  (商品表)  │         │  (关键词表)  │
└───┬────┘ └─────┬──────┘         └──────┬───────┘
    │ 1          │ 1                     │ 1
    │            │                       │
    │ N          │ N                     │ N
┌───┴──────────┐ │              ┌────────┴────────────┐
│brand_social_ │ │              │keyword_monthly_     │
│media         │ │              │volume               │
│(品牌社交媒体)│ │              │(关键词月度搜索量)   │
└──────────────┘ │              └─────────────────────┘
                 │
    ┌────────────┴────────────┐
    │                         │
    │ N                       │ N
┌───┴──────────────┐  ┌───────┴─────────────┐
│product_monthly_  │  │brand_monthly_       │
│sales             │  │trends               │
│(商品月度销售)    │  │(品牌月度趋势)       │
└──────────────────┘  └─────────────────────┘
                              │
                              │ (聚合)
                              │
                      ┌───────┴─────────────┐
                      │market_monthly_      │
                      │trends               │
                      │(市场月度趋势)       │
                      └─────────────────────┘
```

---

## 4️⃣ 核心计算逻辑

### 4.1 品牌月度趋势计算

**目的**：从商品月度销售数据聚合到品牌月度趋势

```sql
-- 计算品牌月度销售额
INSERT INTO brand_monthly_trends (brand_id, date, revenue)
SELECT 
  p.brand_id,
  pms.date,
  SUM(pms.sales) as revenue
FROM product_monthly_sales pms
JOIN products p ON pms.asin = p.asin
WHERE p.market_id = ?
GROUP BY p.brand_id, pms.date
ON DUPLICATE KEY UPDATE revenue = VALUES(revenue);
```

**说明**：
- 从 `product_monthly_sales` 表读取每个商品的月度销售额
- 通过 `products` 表关联到品牌
- 按品牌和月份分组求和
- 使用 `ON DUPLICATE KEY UPDATE` 实现增量更新

---

### 4.2 市场月度趋势计算（销售额）

**目的**：从品牌月度趋势聚合到市场月度趋势

```sql
-- 计算市场月度销售额
INSERT INTO market_monthly_trends (market_id, date, revenue)
SELECT 
  b.market_id,
  bmt.date,
  SUM(bmt.revenue) as revenue
FROM brand_monthly_trends bmt
JOIN brands b ON bmt.brand_id = b.id
WHERE b.market_id = ?
GROUP BY b.market_id, bmt.date
ON DUPLICATE KEY UPDATE revenue = VALUES(revenue);
```

---

### 4.3 市场月度趋势计算（搜索量）

**目的**：计算每个月的关键词搜索量总和

```sql
-- 计算市场月度搜索量
UPDATE market_monthly_trends mmt
SET search_volume = (
  SELECT COALESCE(SUM(kmv.volume), 0)
  FROM keyword_monthly_volume kmv
  JOIN keywords k ON kmv.keyword_id = k.id
  WHERE k.market_id = mmt.market_id
    AND kmv.date = mmt.date
)
WHERE mmt.market_id = ?;
```

---

### 4.4 CAGR计算（品牌）

**目的**：计算品牌的年复合增长率

```sql
-- 计算品牌CAGR
UPDATE brands b
SET cagr = (
  SELECT 
    CASE 
      -- 至少需要12个月的数据
      WHEN COUNT(*) >= 12 AND MIN(revenue) > 0 THEN
        -- CAGR公式：(Ending Value / Beginning Value)^(1/years) - 1
        -- 限制在 -99% 到 +999% 之间
        LEAST(999, GREATEST(-99, 
          (POWER(MAX(revenue) / MIN(revenue), 1.0 / (COUNT(*) / 12.0)) - 1) * 100
        ))
      ELSE NULL
    END
  FROM brand_monthly_trends
  WHERE brand_id = b.id
)
WHERE b.market_id = ?;
```

**CAGR公式详解**：
```
CAGR = (Ending Value / Beginning Value)^(1/years) - 1

其中：
- Ending Value: 最近月份的销售额
- Beginning Value: 最早月份的销售额
- years: 时间跨度（月数 / 12）

示例：
- 12个月前销售额：$100,000
- 最近月份销售额：$150,000
- CAGR = (150000 / 100000)^(1/1) - 1 = 0.5 = 50%
```

---

### 4.5 CAGR计算（市场）

**目的**：计算市场的年复合增长率

```sql
-- 计算市场CAGR
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
WHERE m.id = ?;
```

---

### 4.6 更新品牌总销售额

**目的**：计算品牌最近12个月的总销售额

```sql
-- 更新品牌总销售额（最近12个月）
UPDATE brands b
SET total_revenue = (
  SELECT COALESCE(SUM(revenue), 0)
  FROM brand_monthly_trends
  WHERE brand_id = b.id
    AND date >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)
)
WHERE b.market_id = ?;
```

---

### 4.7 更新品牌商品数量

```sql
-- 更新品牌商品数量
UPDATE brands b
SET product_count = (
  SELECT COUNT(*)
  FROM products
  WHERE brand_id = b.id
)
WHERE b.market_id = ?;
```

---

### 4.8 更新市场总销售额

```sql
-- 更新市场总销售额（最近12个月）
UPDATE markets m
SET total_revenue = (
  SELECT COALESCE(SUM(revenue), 0)
  FROM market_monthly_trends
  WHERE market_id = m.id
    AND date >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)
)
WHERE m.id = ?;
```

---

### 4.9 更新市场商品总数

```sql
-- 更新市场商品总数
UPDATE markets m
SET total_products = (
  SELECT COUNT(*)
  FROM products
  WHERE market_id = m.id
)
WHERE m.id = ?;
```

---

### 4.10 更新市场品牌数量

```sql
-- 更新市场品牌数量
UPDATE markets m
SET brand_count = (
  SELECT COUNT(*)
  FROM brands
  WHERE market_id = m.id
)
WHERE m.id = ?;
```

---

### 4.11 更新市场搜索量

```sql
-- 更新市场搜索量（最近月份）
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
WHERE m.id = ?;
```

---

## 5️⃣ 索引优化

### 5.1 索引策略

| 表名 | 索引类型 | 索引字段 | 目的 |
|-----|---------|---------|------|
| markets | UNIQUE | market_slug | 确保slug唯一性，加速URL查询 |
| markets | INDEX | status | 加速状态筛选 |
| brands | UNIQUE | (market_id, brand_name) | 确保同一市场内品牌名唯一 |
| brands | INDEX | market_id | 加速市场关联查询 |
| brands | INDEX | brand_name | 加速品牌名搜索 |
| brand_social_media | UNIQUE | (brand_id, platform) | 确保同一品牌每个平台只有一条记录 |
| products | UNIQUE | asin | 确保ASIN唯一性 |
| products | INDEX | market_id | 加速市场关联查询 |
| products | INDEX | brand_id | 加速品牌关联查询 |
| products | INDEX | monthly_sales | 加速销量排序 |
| product_monthly_sales | UNIQUE | (asin, date) | 确保同一商品每月只有一条记录 |
| product_monthly_sales | INDEX | asin | 加速商品销售查询 |
| product_monthly_sales | INDEX | date | 加速日期范围查询 |
| keywords | UNIQUE | (market_id, keyword, source) | 确保同一市场同一来源的关键词唯一 |
| keywords | INDEX | market_id | 加速市场关联查询 |
| keyword_monthly_volume | UNIQUE | (keyword_id, date) | 确保同一关键词每月只有一条记录 |
| keyword_monthly_volume | INDEX | keyword_id | 加速关键词搜索量查询 |
| keyword_monthly_volume | INDEX | date | 加速日期范围查询 |
| brand_monthly_trends | UNIQUE | (brand_id, date) | 确保同一品牌每月只有一条记录 |
| market_monthly_trends | UNIQUE | (market_id, date) | 确保同一市场每月只有一条记录 |

### 5.2 复合索引建议

```sql
-- 商品表：按市场和销量排序
CREATE INDEX idx_market_sales ON products(market_id, monthly_sales DESC);

-- 商品表：按品牌和销量排序
CREATE INDEX idx_brand_sales ON products(brand_id, monthly_sales DESC);

-- 商品月度销售：按日期范围查询
CREATE INDEX idx_date_range ON product_monthly_sales(date, asin);

-- 关键词月度搜索量：按日期范围查询
CREATE INDEX idx_date_range ON keyword_monthly_volume(date, keyword_id);
```

---

## 6️⃣ 数据库初始化脚本

### 完整建表脚本

```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS brandtrekin 
  CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;

USE brandtrekin;

-- 1. 市场表
CREATE TABLE `markets` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '市场ID',
  `market_name` VARCHAR(100) NOT NULL COMMENT '市场名称',
  `market_slug` VARCHAR(100) NOT NULL COMMENT '市场slug（用于URL）',
  `description` TEXT COMMENT '市场描述',
  `status` ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
  `total_revenue` DECIMAL(15,2) DEFAULT 0 COMMENT '总销售额（最近12个月）',
  `total_products` INT DEFAULT 0 COMMENT '商品总数',
  `brand_count` INT DEFAULT 0 COMMENT '品牌数量',
  `search_volume` BIGINT DEFAULT 0 COMMENT '搜索量（最近月份）',
  `cagr` DECIMAL(5,2) COMMENT 'CAGR（年复合增长率）',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_market_slug` (`market_slug`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='市场表';

-- 2. 品牌表
CREATE TABLE `brands` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '品牌ID',
  `market_id` INT UNSIGNED NOT NULL COMMENT '市场ID',
  `brand_name` VARCHAR(100) NOT NULL COMMENT '品牌名称',
  `website` VARCHAR(500) COMMENT '品牌独立站',
  `total_revenue` DECIMAL(15,2) DEFAULT 0 COMMENT '总销售额（最近12个月）',
  `product_count` INT DEFAULT 0 COMMENT '商品数量',
  `cagr` DECIMAL(5,2) COMMENT 'CAGR',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_market_brand` (`market_id`, `brand_name`),
  KEY `idx_market_id` (`market_id`),
  KEY `idx_brand_name` (`brand_name`),
  CONSTRAINT `fk_brands_market` FOREIGN KEY (`market_id`) 
    REFERENCES `markets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌表';

-- 3. 品牌社交媒体表
CREATE TABLE `brand_social_media` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `brand_id` INT UNSIGNED NOT NULL COMMENT '品牌ID',
  `platform` ENUM('youtube', 'instagram', 'facebook', 'reddit') NOT NULL COMMENT '平台',
  `url` VARCHAR(500) NOT NULL COMMENT '链接',
  `subscribers` INT DEFAULT 0 COMMENT '订阅数（YouTube）',
  `followers` INT DEFAULT 0 COMMENT '粉丝数（Instagram/Facebook）',
  `posts` INT DEFAULT 0 COMMENT '帖子数（Reddit）',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_brand_platform` (`brand_id`, `platform`),
  KEY `idx_brand_id` (`brand_id`),
  KEY `idx_platform` (`platform`),
  CONSTRAINT `fk_social_brand` FOREIGN KEY (`brand_id`) 
    REFERENCES `brands` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌社交媒体表';

-- 4. 商品表
CREATE TABLE `products` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `market_id` INT UNSIGNED NOT NULL COMMENT '市场ID',
  `brand_id` INT UNSIGNED NOT NULL COMMENT '品牌ID',
  `asin` VARCHAR(20) NOT NULL COMMENT '亚马逊ASIN',
  `title` VARCHAR(500) NOT NULL COMMENT '商品标题',
  `price` DECIMAL(10,2) DEFAULT 0 COMMENT '价格',
  `rating` DECIMAL(3,2) DEFAULT 0 COMMENT '评分（0-5）',
  `reviews` INT DEFAULT 0 COMMENT '评论数',
  `monthly_sales` INT DEFAULT 0 COMMENT '月销量',
  `image_url` VARCHAR(500) COMMENT '图片URL',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_asin` (`asin`),
  KEY `idx_market_id` (`market_id`),
  KEY `idx_brand_id` (`brand_id`),
  KEY `idx_monthly_sales` (`monthly_sales`),
  CONSTRAINT `fk_products_market` FOREIGN KEY (`market_id`) 
    REFERENCES `markets` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_products_brand` FOREIGN KEY (`brand_id`) 
    REFERENCES `brands` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- 5. 商品月度销售表
CREATE TABLE `product_monthly_sales` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `asin` VARCHAR(20) NOT NULL COMMENT '商品ASIN',
  `date` DATE NOT NULL COMMENT '月份（YYYY-MM-01）',
  `sales` DECIMAL(12,2) DEFAULT 0 COMMENT '销售额',
  `units` INT DEFAULT 0 COMMENT '销量',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_asin_date` (`asin`, `date`),
  KEY `idx_asin` (`asin`),
  KEY `idx_date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品月度销售表';

-- 6. 关键词表
CREATE TABLE `keywords` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `market_id` INT UNSIGNED NOT NULL COMMENT '市场ID',
  `keyword` VARCHAR(200) NOT NULL COMMENT '关键词',
  `source` ENUM('google', 'amazon') NOT NULL COMMENT '来源',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_market_keyword_source` (`market_id`, `keyword`, `source`),
  KEY `idx_market_id` (`market_id`),
  KEY `idx_keyword` (`keyword`),
  KEY `idx_source` (`source`),
  CONSTRAINT `fk_keywords_market` FOREIGN KEY (`market_id`) 
    REFERENCES `markets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='关键词表';

-- 7. 关键词月度搜索量表
CREATE TABLE `keyword_monthly_volume` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `keyword_id` INT UNSIGNED NOT NULL COMMENT '关键词ID',
  `date` DATE NOT NULL COMMENT '月份（YYYY-MM-01）',
  `volume` BIGINT DEFAULT 0 COMMENT '搜索量',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_keyword_date` (`keyword_id`, `date`),
  KEY `idx_keyword_id` (`keyword_id`),
  KEY `idx_date` (`date`),
  CONSTRAINT `fk_volume_keyword` FOREIGN KEY (`keyword_id`) 
    REFERENCES `keywords` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='关键词月度搜索量表';

-- 8. 品牌月度趋势表
CREATE TABLE `brand_monthly_trends` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `brand_id` INT UNSIGNED NOT NULL COMMENT '品牌ID',
  `date` DATE NOT NULL COMMENT '月份（YYYY-MM-01）',
  `revenue` DECIMAL(12,2) DEFAULT 0 COMMENT '销售额',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_brand_date` (`brand_id`, `date`),
  KEY `idx_brand_id` (`brand_id`),
  KEY `idx_date` (`date`),
  CONSTRAINT `fk_trends_brand` FOREIGN KEY (`brand_id`) 
    REFERENCES `brands` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌月度趋势表';

-- 9. 市场月度趋势表
CREATE TABLE `market_monthly_trends` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `market_id` INT UNSIGNED NOT NULL COMMENT '市场ID',
  `date` DATE NOT NULL COMMENT '月份（YYYY-MM-01）',
  `revenue` DECIMAL(15,2) DEFAULT 0 COMMENT '销售额',
  `search_volume` BIGINT DEFAULT 0 COMMENT '搜索量',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_market_date` (`market_id`, `date`),
  KEY `idx_market_id` (`market_id`),
  KEY `idx_date` (`date`),
  CONSTRAINT `fk_trends_market` FOREIGN KEY (`market_id`) 
    REFERENCES `markets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='市场月度趋势表';

-- 10. 导入日志表
CREATE TABLE `import_logs` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `market_id` INT UNSIGNED NOT NULL COMMENT '市场ID',
  `import_mode` ENUM('incremental', 'replace') NOT NULL COMMENT '导入模式',
  `status` ENUM('success', 'failed', 'partial') NOT NULL COMMENT '状态',
  `brands_count` INT DEFAULT 0 COMMENT '导入品牌数',
  `products_count` INT DEFAULT 0 COMMENT '导入商品数',
  `keywords_count` INT DEFAULT 0 COMMENT '导入关键词数',
  `sales_records_count` INT DEFAULT 0 COMMENT '导入销售记录数',
  `skipped_count` INT DEFAULT 0 COMMENT '跳过记录数',
  `error_message` TEXT COMMENT '错误信息',
  `log_content` LONGTEXT COMMENT '日志内容',
  `created_by` VARCHAR(50) COMMENT '操作人',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_market_id` (`market_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_logs_market` FOREIGN KEY (`market_id`) 
    REFERENCES `markets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='导入日志表';
```

---

## ✅ 数据库开发检查清单

### 表结构
- [ ] 所有10张表创建成功
- [ ] 所有字段类型正确
- [ ] 所有注释完整
- [ ] 字符集为utf8mb4

### 索引
- [ ] 所有主键索引创建成功
- [ ] 所有唯一索引创建成功
- [ ] 所有普通索引创建成功
- [ ] 外键约束创建成功

### 计算逻辑
- [ ] 品牌月度趋势计算正确
- [ ] 市场月度趋势计算正确
- [ ] CAGR计算正确
- [ ] 所有聚合指标更新正确

### 性能
- [ ] 查询性能测试通过
- [ ] 大数据量导入测试通过
- [ ] 索引使用率检查通过

---

**数据库设计文档完成！请继续阅读API接口设计文档。**
