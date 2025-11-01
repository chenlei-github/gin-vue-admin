#!/bin/bash

echo "🔄 重启 BrandTrekin 后端服务..."

# 查找并停止现有服务
PID=$(lsof -ti :8888)
if [ ! -z "$PID" ]; then
    echo "📌 找到运行中的服务 (PID: $PID)"
    kill -9 $PID
    echo "✅ 已停止旧服务"
    sleep 2
else
    echo "ℹ️  没有找到运行中的服务"
fi

# 进入服务器目录
cd /Users/leon/code/ai-code/gin-vue-admin/server

# 启动新服务
echo "🚀 启动新服务..."
nohup go run main.go > ../logs/server.log 2>&1 &

# 等待服务启动
sleep 3

# 检查服务状态
NEW_PID=$(lsof -ti :8888)
if [ ! -z "$NEW_PID" ]; then
    echo "✅ 服务启动成功 (PID: $NEW_PID)"
    echo "📝 日志文件: /Users/leon/code/ai-code/gin-vue-admin/logs/server.log"
    echo "🌐 服务地址: http://localhost:8888"
else
    echo "❌ 服务启动失败，请查看日志"
    tail -n 20 ../logs/server.log
    exit 1
fi

echo ""
echo "🎉 重启完成！"
echo ""
echo "📋 下一步操作："
echo "1. 访问 http://localhost:8888 确认服务正常"
echo "2. 通过 API 导入数据进行验证："
echo "   POST /btImport/batchImport"
echo "   参数："
echo "   - marketId: 1"
echo "   - replaceMode: true"
echo "   - productSales: trekin-main/data/CNCRouter/product-US-sales.xlsx"
echo ""
