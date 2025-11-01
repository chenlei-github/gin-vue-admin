#!/usr/bin/env python3
import pymysql

# 连接数据库
conn = pymysql.connect(
    host='rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com',
    user='brandtrekin',
    password='cl@2025@!',
    database='brandtrekin',
    charset='utf8mb4'
)

cursor = conn.cursor()

# 查询统计信息
cursor.execute("SELECT COUNT(*) FROM bt_product_monthly_sales WHERE market_id = 1")
total = cursor.fetchone()[0]

cursor.execute("SELECT COUNT(*) FROM bt_product_monthly_sales WHERE market_id = 1 AND units > 0 AND sales = 0")
abnormal = cursor.fetchone()[0]

cursor.execute("SELECT COUNT(*) FROM bt_product_monthly_sales WHERE market_id = 1 AND units > 0 AND sales > 0")
normal = cursor.fetchone()[0]

print("=" * 60)
print("数据统计")
print("=" * 60)
print(f"总记录数: {total}")
print(f"异常记录 (units>0, sales=0): {abnormal}")
print(f"正常记录 (units>0, sales>0): {normal}")
print(f"异常比例: {abnormal/total*100:.2f}%")
print()

# 查询异常数据示例
cursor.execute("SELECT asin, date, units, sales FROM bt_product_monthly_sales WHERE market_id = 1 AND units > 0 AND sales = 0 ORDER BY date DESC LIMIT 10")
print("=" * 60)
print("异常数据示例 (最新10条)")
print("=" * 60)
for row in cursor.fetchall():
    print(f"ASIN: {row[0]:<15} Date: {row[1]} Units: {row[2]:<5} Sales: ${row[3]}")

cursor.close()
conn.close()
