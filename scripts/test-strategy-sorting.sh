#!/bin/bash
# 策略排序调试测试脚本
# 用于快速启动服务器并打开浏览器进行测试

echo "========================================"
echo "策略排序调试测试"
echo "========================================"
echo ""

# 检查服务器是否已编译
if [ ! -f "bin/server" ]; then
    echo "[步骤 1/3] 编译服务器..."
    go build -o bin/server ./cmd/server
    if [ $? -ne 0 ]; then
        echo "编译失败！"
        exit 1
    fi
    echo "编译成功！"
    echo ""
else
    echo "[步骤 1/3] 服务器已编译，跳过编译步骤"
    echo ""
fi

echo "[步骤 2/3] 启动后端服务器..."
echo "请查看服务器控制台日志，观察'排序前顺序'和'排序后顺序'"
echo ""

# 在后台启动服务器
./bin/server > server-debug.log 2>&1 &
SERVER_PID=$!
echo "服务器 PID: $SERVER_PID"

# 等待服务器启动
sleep 3

echo "[步骤 3/3] 准备打开浏览器..."
echo ""
echo "========================================"
echo "调试指南："
echo "========================================"
echo ""
echo "1. 按 F12 打开浏览器开发者工具"
echo "2. 切换到 Console (控制台) 标签"
echo "3. 导航到策略管理页面"
echo "4. 多次点击'刷新'按钮 (建议 5-10 次)"
echo "5. 观察控制台日志："
echo "   - 每次应该显示'策略顺序'"
echo "   - 如果顺序变化，会显示警告⚠️"
echo "   - 查看'原始API响应'确认数据"
echo ""
echo "6. 观察后端服务器日志："
echo "   - '排序前顺序' - 可能每次不同(正常)"
echo "   - '排序后顺序' - 应该每次相同(关键)"
echo "   - '返回策略列表' - 最终返回的顺序"
echo ""
echo "   查看服务器日志: tail -f server-debug.log"
echo ""
echo "详细调试指南: docs/debugging/STRATEGY_SORTING_DEBUG_GUIDE.md"
echo ""
echo "========================================"
echo ""

# 打开浏览器
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    open http://localhost:3000/#/strategies
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    xdg-open http://localhost:3000/#/strategies 2>/dev/null || sensible-browser http://localhost:3000/#/strategies 2>/dev/null
fi

echo ""
echo "测试环境已准备就绪！"
echo ""
echo "提示："
echo "- 服务器运行在后台 (PID: $SERVER_PID)"
echo "- 查看服务器日志: tail -f server-debug.log"
echo "- 停止服务器: kill $SERVER_PID"
echo "- 前端日志在浏览器控制台"
echo ""
echo "按 Ctrl+C 停止服务器..."
echo ""

# 等待用户中断
wait $SERVER_PID

