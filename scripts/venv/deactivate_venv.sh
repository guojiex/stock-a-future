#!/bin/bash
# 停用Python虚拟环境的脚本

# 检查是否在虚拟环境中
if [ -z "$VIRTUAL_ENV" ]; then
    echo "错误: 当前未激活任何虚拟环境!"
    exit 1
fi

# 显示当前虚拟环境信息
echo "当前虚拟环境: $VIRTUAL_ENV"

# 停用虚拟环境
echo "停用虚拟环境..."
deactivate

# 检查是否成功停用
if [ -z "$VIRTUAL_ENV" ]; then
    echo "虚拟环境已成功停用!"
else
    echo "警告: 虚拟环境可能未正确停用。"
fi

echo ""
echo "如需重新激活虚拟环境，请运行:"
echo "source .venv/bin/activate 或 ./scripts/setup_venv.sh"
