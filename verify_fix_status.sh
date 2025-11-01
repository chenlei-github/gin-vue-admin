#!/bin/bash

echo "=== 验证全量替换导入修复状态 ==="
echo ""

echo "📋 检查项 1: bt_import.go 中的 Unscoped() 物理删除"
echo "----------------------------------------"
UNSCOPED_COUNT=$(grep -c "Unscoped()" /Users/leon/code/ai-code/gin-vue-admin/server/service/brandtrekin/bt_import.go)
if [ "$UNSCOPED_COUNT" -gt 0 ]; then
    echo "✅ 找到 $UNSCOPED_COUNT 处 Unscoped() 调用"
    echo "   修复状态: 已应用"
else
    echo "❌ 未找到 Unscoped() 调用"
    echo "   修复状态: 未应用"
fi
echo ""

echo "📋 检查项 2: bt_import_optimized.go 中的去重逻辑"
echo "----------------------------------------"
DEDUP_COUNT=$(grep -c "productMap\|brandMap\|keywordMap" /Users/leon/code/ai-code/gin-vue-admin/server/service/brandtrekin/bt_import_optimized.go)
if [ "$DEDUP_COUNT" -gt 10 ]; then
    echo "✅ 找到去重相关代码"
    echo "   修复状态: 已应用"
else
    echo "⚠️  去重逻辑可能未完全应用"
    echo "   修复状态: 需要检查"
fi
echo ""

echo "📋 检查项 3: 服务器运行状态"
echo "----------------------------------------"
if lsof -i :8888 > /dev/null 2>&1; then
    echo "✅ 服务器正在运行 (端口 8888)"
    PID=$(lsof -ti :8888)
    echo "   进程 PID: $PID"
else
    echo "❌ 服务器未运行"
fi
echo ""

echo "📋 检查项 4: 修复文档"
echo "----------------------------------------"
if [ -f "/Users/leon/code/ai-code/gin-vue-admin/FIX_SOFT_DELETE_CONFLICT.md" ]; then
    echo "✅ FIX_SOFT_DELETE_CONFLICT.md 存在"
else
    echo "⚠️  FIX_SOFT_DELETE_CONFLICT.md 不存在"
fi

if [ -f "/Users/leon/code/ai-code/gin-vue-admin/FIX_DUPLICATE_KEY_SUMMARY.md" ]; then
    echo "✅ FIX_DUPLICATE_KEY_SUMMARY.md 存在"
else
    echo "⚠️  FIX_DUPLICATE_KEY_SUMMARY.md 不存在"
fi

if [ -f "/Users/leon/code/ai-code/gin-vue-admin/TEST_IMPORT_GUIDE.md" ]; then
    echo "✅ TEST_IMPORT_GUIDE.md 存在"
else
    echo "⚠️  TEST_IMPORT_GUIDE.md 不存在"
fi
echo ""

echo "📋 总结"
echo "----------------------------------------"
echo "修复 1: 软删除冲突 (Unscoped 物理删除)"
if [ "$UNSCOPED_COUNT" -gt 0 ]; then
    echo "  状态: ✅ 已修复"
else
    echo "  状态: ❌ 未修复"
fi

echo ""
echo "修复 2: CSV 重复数据去重"
echo "  状态: ⚠️  需要添加到 bt_import_optimized.go"
echo ""

echo "🎯 下一步操作:"
echo "1. 添加去重逻辑到 bt_import_optimized.go"
echo "2. 重启服务器（如果代码已修改）"
echo "3. 按照 TEST_IMPORT_GUIDE.md 进行测试"
echo ""
