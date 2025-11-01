#!/bin/bash

echo "🧪 验证修复效果 - units>0 但 sales=0 的问题"
echo "================================================"
echo ""

# 数据库连接信息
DB_HOST="rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com"
DB_USER="brandtrekin"
DB_PASS="cl@2025@!"
DB_NAME="brandtrekin"

echo "📊 检查数据库中的异常数据..."
echo ""

# 检查 units>0 但 sales=0 的记录数量
echo "1️⃣ 统计 units>0 但 sales=0 的记录数量："
mysql -h $DB_HOST -u $DB_USER -p"$DB_PASS" $DB_NAME -e "
SELECT COUNT(*) as '异常记录数' 
FROM bt_product_monthly_sales 
WHERE units > 0 AND sales = 0;
"

echo ""
echo "2️⃣ 查看具体的异常数据（前10条）："
mysql -h $DB_HOST -u $DB_USER -p"$DB_PASS" $DB_NAME -e "
SELECT asin, DATE_FORMAT(date, '%Y-%m') as month, units, sales 
FROM bt_product_monthly_sales 
WHERE units > 0 AND sales = 0 
ORDER BY date DESC, asin 
LIMIT 10;
"

echo ""
echo "3️⃣ 查看 CNC Router 市场的数据统计："
mysql -h $DB_HOST -u $DB_USER -p"$DB_PASS" $DB_NAME -e "
SELECT 
    market_id,
    COUNT(*) as '总记录数',
    SUM(CASE WHEN units > 0 AND sales = 0 THEN 1 ELSE 0 END) as '异常记录数',
    SUM(CASE WHEN units > 0 AND sales > 0 THEN 1 ELSE 0 END) as '正常记录数',
    SUM(CASE WHEN units = 0 AND sales = 0 THEN 1 ELSE 0 END) as '零值记录数'
FROM bt_product_monthly_sales 
WHERE market_id = 1
GROUP BY market_id;
"

echo ""
echo "4️⃣ 查看最近的数据样本（前20条）："
mysql -h $DB_HOST -u $DB_USER -p"$DB_PASS" $DB_NAME -e "
SELECT 
    asin, 
    DATE_FORMAT(date, '%Y-%m') as month, 
    units, 
    ROUND(sales, 2) as sales,
    CASE 
        WHEN units > 0 AND sales = 0 THEN '❌ 异常'
        WHEN units > 0 AND sales > 0 THEN '✅ 正常'
        WHEN units = 0 AND sales = 0 THEN '⚪ 零值'
        ELSE '❓ 其他'
    END as status
FROM bt_product_monthly_sales 
WHERE market_id = 1
ORDER BY date DESC, asin 
LIMIT 20;
"

echo ""
echo "================================================"
echo "✅ 验证完成！"
echo ""
echo "📝 说明："
echo "- 异常记录：units>0 但 sales=0（需要修复）"
echo "- 正常记录：units>0 且 sales>0"
echo "- 零值记录：units=0 且 sales=0（正常情况）"
echo ""
