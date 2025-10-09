#!/bin/bash

# 打开拖拽分组功能测试指南

echo ""
echo "========================================"
echo "  拖拽分组功能测试指南"
echo "========================================"
echo ""
echo "正在打开测试指南..."
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TEST_FILE="$SCRIPT_DIR/test-drag-groups.html"

# 根据操作系统打开HTML文件
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    open "$TEST_FILE"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    xdg-open "$TEST_FILE" 2>/dev/null || \
    sensible-browser "$TEST_FILE" 2>/dev/null || \
    firefox "$TEST_FILE" 2>/dev/null || \
    chromium-browser "$TEST_FILE" 2>/dev/null
else
    # Windows (Git Bash)
    start "" "$TEST_FILE"
fi

echo ""
echo "测试指南已在浏览器中打开！"
echo ""
echo "如需查看详细技术文档，请查看："
echo "  - docs/features/DRAG_AND_DROP_GROUPS.md"
echo "  - DRAG_DROP_ENHANCEMENT_SUMMARY.md"
echo ""

