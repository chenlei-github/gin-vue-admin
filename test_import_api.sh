#!/bin/bash

echo "🚀 通过 API 导入 product-US-sales.xlsx"
echo "================================================"
echo ""

# API 配置
API_URL="http://localhost:8888/btImport/batchImport"
EXCEL_FILE="/Users/leon/code/ai-code/gin-vue-admin/trekin-main/data/CNCRouter/product-US-sales.xlsx"

# 检查文件是否存在
if [ ! -f "$EXCEL_FILE" ]; then
    echo "❌ 错误：找不到文件 $EXCEL_FILE"
    exit 1
fi

echo "📁 文件路径: $EXCEL_FILE"
echo "🌐 API 地址: $API_URL"
echo ""

# 检查服务是否运行
if ! lsof -i :8888 > /dev/null 2>&1; then
    echo "❌ 错误：后端服务未运行（端口 8888）"
    echo "请先运行: ./restart_server.sh"
    exit 1
fi

echo "✅ 后端服务正在运行"
echo ""

# 获取 token（如果需要）
# 这里假设你已经登录，或者 API 不需要认证
# 如果需要 token，请先调用登录接口获取

echo "📤 开始上传文件..."
echo ""

# 调用 API 导入数据
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL" \
  -F "marketId=1" \
  -F "replaceMode=true" \
  -F "productSales=@$EXCEL_FILE")

# 分离响应体和状态码
HTTP_BODY=$(echo "$RESPONSE" | head -n -1)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

echo "📊 HTTP 状态码: $HTTP_CODE"
echo ""

if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ 导入成功！"
    echo ""
    echo "📋 响应内容:"
    echo "$HTTP_BODY" | jq '.' 2>/dev/null || echo "$HTTP_BODY"
    echo ""
    echo "🎉 下一步："
    echo "运行验证脚本检查数据: ./verify_sales_fix.sh"
else
    echo "❌ 导入失败！"
    echo ""
    echo "📋 错误信息:"
    echo "$HTTP_BODY" | jq '.' 2>/dev/null || echo "$HTTP_BODY"
    echo ""
    echo "💡 可能的原因："
    echo "1. 需要登录认证（token）"
    echo "2. API 路径不正确"
    echo "3. 参数格式不正确"
    echo "4. 服务器内部错误"
    echo ""
    echo "请检查服务器日志: tail -f logs/server.log"
fi

echo ""
echo "================================================"
