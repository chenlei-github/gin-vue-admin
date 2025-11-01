#!/bin/bash

# 测试脚本：验证重复键冲突问题已修复
# 使用方法：./test_duplicate_fix.sh

echo "=========================================="
echo "测试：全量替换模式重复键冲突修复"
echo "=========================================="
echo ""

# 配置
API_URL="http://localhost:8888/api/v1/brandtrekin/import"
MARKET_ID=1
DATA_DIR="/Users/leon/code/ai-code/gin-vue-admin/trekin-main/data/CNCRouter"

echo "📋 测试配置："
echo "  - API URL: $API_URL"
echo "  - Market ID: $MARKET_ID"
echo "  - Data Dir: $DATA_DIR"
echo ""

# 测试 1：第一次导入
echo "🧪 测试 1：第一次全量导入"
echo "----------------------------------------"
curl -X POST "$API_URL" \
  -F "marketId=$MARKET_ID" \
  -F "replaceMode=true" \
  -F "files=@$DATA_DIR/Products.csv" \
  -F "files=@$DATA_DIR/Brands.csv" \
  -F "files=@$DATA_DIR/GKW.csv"

echo ""
echo ""

# 等待 2 秒
sleep 2

# 测试 2：第二次导入（验证软删除冲突已修复）
echo "🧪 测试 2：第二次全量导入（验证软删除冲突已修复）"
echo "----------------------------------------"
curl -X POST "$API_URL" \
  -F "marketId=$MARKET_ID" \
  -F "replaceMode=true" \
  -F "files=@$DATA_DIR/Products.csv" \
  -F "files=@$DATA_DIR/Brands.csv" \
  -F "files=@$DATA_DIR/GKW.csv"

echo ""
echo ""

# 测试 3：验证数据完整性
echo "🧪 测试 3：验证数据完整性"
echo "----------------------------------------"

# 需要 MySQL 连接信息
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="root"
DB_PASS="your_password"
DB_NAME="gva"

echo "检查软删除数据..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
SELECT 
    '商品' as 表名,
    COUNT(*) as 软删除数量
FROM bt_products 
WHERE deleted_at IS NOT NULL
UNION ALL
SELECT 
    '品牌' as 表名,
    COUNT(*) as 软删除数量
FROM bt_brands 
WHERE deleted_at IS NOT NULL
UNION ALL
SELECT 
    '关键词' as 表名,
    COUNT(*) as 软删除数量
FROM bt_keywords 
WHERE deleted_at IS NOT NULL;
"

echo ""
echo "检查重复数据..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
SELECT 
    asin, 
    COUNT(*) as count 
FROM bt_products 
WHERE deleted_at IS NULL 
GROUP BY asin 
HAVING count > 1;
"

echo ""
echo "=========================================="
echo "✅ 测试完成！"
echo "=========================================="
echo ""
echo "预期结果："
echo "  1. 两次导入都应该成功，不报错"
echo "  2. 软删除数量应该为 0"
echo "  3. 不应该有重复的 ASIN"
echo ""
