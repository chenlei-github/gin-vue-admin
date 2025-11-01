#!/bin/bash

echo "🔍 测试数据库连接"
echo "================================================"
echo ""

DB_HOST="rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com"
DB_USER="brandtrekin"
DB_PASS="cl@2025@!"
DB_NAME="brandtrekin"

echo "正在连接数据库..."
echo "Host: $DB_HOST"
echo "Database: $DB_NAME"
echo ""

python3 -c "
import pymysql
import sys

try:
    print('正在连接...')
    conn = pymysql.connect(
        host='$DB_HOST',
        user='$DB_USER',
        password='$DB_PASS',
        database='$DB_NAME',
        charset='utf8mb4',
        connect_timeout=10
    )
    print('✅ 数据库连接成功!')
    print()
    
    cursor = conn.cursor()
    
    # 查询 MySQL 版本
    cursor.execute('SELECT VERSION()')
    version = cursor.fetchone()
    print(f'MySQL 版本: {version[0]}')
    print()
    
    # 查询异常数据统计
    cursor.execute('SELECT COUNT(*) FROM bt_product_monthly_sales WHERE units > 0 AND sales = 0')
    abnormal_count = cursor.fetchone()[0]
    print(f'异常记录数 (units>0 但 sales=0): {abnormal_count}')
    
    # 查询总记录数
    cursor.execute('SELECT COUNT(*) FROM bt_product_monthly_sales WHERE market_id = 1')
    total_count = cursor.fetchone()[0]
    print(f'总记录数 (market_id=1): {total_count}')
    
    # 查询正常记录数
    cursor.execute('SELECT COUNT(*) FROM bt_product_monthly_sales WHERE market_id = 1 AND units > 0 AND sales > 0')
    normal_count = cursor.fetchone()[0]
    print(f'正常记录数 (units>0 且 sales>0): {normal_count}')
    
    print()
    print('=' * 60)
    print('数据分布:')
    print('=' * 60)
    
    cursor.execute('''
    SELECT 
        COUNT(*) as total,
        SUM(CASE WHEN units > 0 AND sales = 0 THEN 1 ELSE 0 END) as abnormal,
        SUM(CASE WHEN units > 0 AND sales > 0 THEN 1 ELSE 0 END) as normal,
        SUM(CASE WHEN units = 0 AND sales = 0 THEN 1 ELSE 0 END) as zero_both
    FROM bt_product_monthly_sales 
    WHERE market_id = 1
    ''')
    
    row = cursor.fetchone()
    print(f'总记录数: {row[0]}')
    print(f'异常记录 (units>0, sales=0): {row[1] or 0}')
    print(f'正常记录 (units>0, sales>0): {row[2] or 0}')
    print(f'零值记录 (units=0, sales=0): {row[3] or 0}')
    
    if row[1]:
        abnormal_ratio = (row[1] / row[0]) * 100
        print(f'异常比例: {abnormal_ratio:.2f}%')
    
    print()
    print('=' * 60)
    print('异常数据示例 (最新5条):')
    print('=' * 60)
    
    cursor.execute('''
    SELECT asin, date, units, sales 
    FROM bt_product_monthly_sales 
    WHERE market_id = 1 AND units > 0 AND sales = 0 
    ORDER BY date DESC 
    LIMIT 5
    ''')
    
    for row in cursor.fetchall():
        print(f'ASIN: {row[0]:<15} Date: {row[1]} Units: {row[2]:<5} Sales: ${row[3]}')
    
    cursor.close()
    conn.close()
    
    print()
    print('✅ 测试完成!')
    sys.exit(0)
    
except Exception as e:
    print(f'❌ 连接失败: {e}')
    import traceback
    traceback.print_exc()
    sys.exit(1)
"

echo ""
echo "================================================"
