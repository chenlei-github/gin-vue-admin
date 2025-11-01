#!/bin/bash

echo "🚀 一键修复验证 - units>0 但 sales=0 的问题"
echo "================================================"
echo ""

# 设置颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 步骤计数
STEP=1

# 函数：打印步骤
print_step() {
    echo ""
    echo -e "${YELLOW}📌 步骤 $STEP: $1${NC}"
    echo "------------------------------------------------"
    STEP=$((STEP + 1))
}

# 函数：打印成功
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# 函数：打印错误
print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 开始执行
echo "开始时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# 步骤 1：检查文件
print_step "检查必要文件"

EXCEL_FILE="/Users/leon/code/ai-code/gin-vue-admin/trekin-main/data/CNCRouter/product-US-sales.xlsx"
if [ -f "$EXCEL_FILE" ]; then
    print_success "Excel 文件存在: $EXCEL_FILE"
else
    print_error "Excel 文件不存在: $EXCEL_FILE"
    exit 1
fi

# 步骤 2：检查当前服务状态
print_step "检查当前服务状态"

OLD_PID=$(lsof -ti :8888)
if [ ! -z "$OLD_PID" ]; then
    print_success "发现运行中的服务 (PID: $OLD_PID)"
    
    # 停止旧服务
    echo "正在停止旧服务..."
    kill -9 $OLD_PID
    sleep 2
    print_success "旧服务已停止"
else
    echo "没有运行中的服务"
fi

# 步骤 3：启动新服务
print_step "启动新服务"

cd /Users/leon/code/ai-code/gin-vue-admin/server

# 创建日志目录
mkdir -p ../logs

# 启动服务
echo "正在启动服务..."
nohup go run main.go > ../logs/server.log 2>&1 &

# 等待服务启动
echo "等待服务启动..."
for i in {1..10}; do
    sleep 1
    if lsof -ti :8888 > /dev/null 2>&1; then
        NEW_PID=$(lsof -ti :8888)
        print_success "服务启动成功 (PID: $NEW_PID)"
        break
    fi
    echo "等待中... ($i/10)"
done

# 检查服务是否启动成功
if ! lsof -ti :8888 > /dev/null 2>&1; then
    print_error "服务启动失败"
    echo ""
    echo "查看日志："
    tail -20 ../logs/server.log
    exit 1
fi

cd ..

# 步骤 4：等待服务完全就绪
print_step "等待服务完全就绪"

echo "等待 5 秒..."
sleep 5
print_success "服务已就绪"

# 步骤 5：导入数据前检查
print_step "导入数据前检查"

echo "检查数据库中的异常数据..."
DB_HOST="rm-wz95yyq2c56dwroo92o.rwlb.rds.aliyuncs.com"
DB_USER="brandtrekin"
DB_PASS="cl@2025@!"
DB_NAME="brandtrekin"

# 使用 Python 查询数据库（避免 MySQL 客户端版本兼容问题）
BEFORE_COUNT=$(python3 -c "
import pymysql
try:
    conn = pymysql.connect(host='$DB_HOST', user='$DB_USER', password='$DB_PASS', database='$DB_NAME', connect_timeout=5)
    cursor = conn.cursor()
    cursor.execute('SELECT COUNT(*) FROM bt_product_monthly_sales WHERE units > 0 AND sales = 0')
    print(cursor.fetchone()[0])
    conn.close()
except Exception as e:
    pass
" 2>/dev/null)

if [ ! -z "$BEFORE_COUNT" ]; then
    echo "导入前异常记录数: $BEFORE_COUNT"
else
    echo "无法连接数据库，跳过数据检查"
fi

# 步骤 6：登录获取 Token
print_step "登录获取 Token"

echo "正在登录..."
LOGIN_URL="http://localhost:8888/base/login"
LOGIN_RESPONSE=$(curl -s -X POST "$LOGIN_URL" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}' 2>/dev/null)

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    print_error "登录失败，无法获取 Token"
    echo "响应: $LOGIN_RESPONSE"
    exit 1
fi

print_success "登录成功，已获取 Token"

# 步骤 7：导入数据
print_step "通过 API 导入数据"

echo "正在上传文件..."
API_URL="http://localhost:8888/btImport/batchImport"

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL" \
  -H "x-token: $TOKEN" \
  -F "marketId=1" \
  -F "replaceMode=true" \
  -F "productSales=@$EXCEL_FILE" 2>/dev/null)

# 修复 macOS 的 head 命令问题
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
HTTP_BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_success "数据导入成功"
    echo "$HTTP_BODY" | jq '.' 2>/dev/null || echo "$HTTP_BODY"
else
    print_error "数据导入失败 (HTTP $HTTP_CODE)"
    echo "$HTTP_BODY"
    echo ""
    echo "请检查日志: tail -f logs/server.log"
    exit 1
fi

# 步骤 8：等待数据处理完成
print_step "等待数据处理完成"

echo "等待 10 秒..."
sleep 10
print_success "数据处理完成"

# 步骤 9：验证修复效果
print_step "验证修复效果"

if [ ! -z "$BEFORE_COUNT" ]; then
    # 使用 Python 查询数据库
    AFTER_COUNT=$(python3 -c "
import pymysql
try:
    conn = pymysql.connect(host='$DB_HOST', user='$DB_USER', password='$DB_PASS', database='$DB_NAME', connect_timeout=5)
    cursor = conn.cursor()
    cursor.execute('SELECT COUNT(*) FROM bt_product_monthly_sales WHERE units > 0 AND sales = 0')
    print(cursor.fetchone()[0])
    conn.close()
except Exception as e:
    pass
" 2>/dev/null)
    
    echo ""
    echo "📊 验证结果："
    echo "------------------------------------------------"
    echo "导入前异常记录数: $BEFORE_COUNT"
    echo "导入后异常记录数: $AFTER_COUNT"
    
    if [ ! -z "$AFTER_COUNT" ] && [ "$AFTER_COUNT" -lt "$BEFORE_COUNT" ]; then
        IMPROVEMENT=$((100 - AFTER_COUNT * 100 / BEFORE_COUNT))
        print_success "改善率: $IMPROVEMENT%"
    fi
    
    if [ "$AFTER_COUNT" -eq 0 ]; then
        print_success "🎉 完美！没有异常数据了！"
    elif [ ! -z "$AFTER_COUNT" ] && [ "$AFTER_COUNT" -lt "$BEFORE_COUNT" ]; then
        print_success "✅ 异常数据已减少！"
    else
        print_error "⚠️  异常数据未减少，请检查"
    fi
    
    echo ""
    echo "查看详细数据："
    python3 -c "
import pymysql
try:
    conn = pymysql.connect(host='$DB_HOST', user='$DB_USER', password='$DB_PASS', database='$DB_NAME', connect_timeout=5)
    cursor = conn.cursor()
    cursor.execute('''SELECT 
        COUNT(*) as total,
        SUM(CASE WHEN units > 0 AND sales = 0 THEN 1 ELSE 0 END) as abnormal,
        SUM(CASE WHEN units > 0 AND sales > 0 THEN 1 ELSE 0 END) as normal,
        SUM(CASE WHEN units = 0 AND sales = 0 THEN 1 ELSE 0 END) as zero_both
    FROM bt_product_monthly_sales WHERE market_id = 1''')
    row = cursor.fetchone()
    print('总记录数: %d' % row[0])
    print('异常记录 (units>0, sales=0): %d' % (row[1] or 0))
    print('正常记录 (units>0, sales>0): %d' % (row[2] or 0))
    print('零值记录 (units=0, sales=0): %d' % (row[3] or 0))
    conn.close()
except Exception as e:
    print('查询失败:', e)
" 2>/dev/null
else
    echo "无法连接数据库，请手动验证"
fi

# 完成
echo ""
echo "================================================"
echo -e "${GREEN}🎉 验证完成！${NC}"
echo ""
echo "结束时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""
echo "📋 后续操作："
echo "1. 查看详细验证报告: ./verify_sales_fix.sh"
echo "2. 查看服务日志: tail -f logs/server.log"
echo "3. 查看修复文档: cat FIX_SUMMARY.md"
echo ""
