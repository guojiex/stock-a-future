#!/bin/bash
# 创建和激活Python虚拟环境的脚本

# 设置虚拟环境名称和路径
VENV_NAME="stock_env"
VENV_PATH="$(pwd)/.venv"

# 检查Python版本
echo "检查Python版本..."
python3 --version || { echo "错误: 未找到Python3，请确保已安装Python 3.7+"; exit 1; }

# 检查是否已存在虚拟环境
if [ -d "$VENV_PATH" ]; then
    echo "发现已存在的虚拟环境: $VENV_PATH"
    read -p "是否重新创建虚拟环境? (y/n): " RECREATE
    if [ "$RECREATE" = "y" ] || [ "$RECREATE" = "Y" ]; then
        echo "删除现有虚拟环境..."
        rm -rf "$VENV_PATH"
    else
        echo "使用现有虚拟环境..."
        source "$VENV_PATH/bin/activate"
        echo "虚拟环境已激活: $VENV_NAME"
        echo "可以使用 'deactivate' 命令退出虚拟环境"
        exit 0
    fi
fi

# 创建虚拟环境
echo "创建Python虚拟环境: $VENV_NAME..."
python3 -m venv "$VENV_PATH"

# 检查虚拟环境是否创建成功
if [ ! -d "$VENV_PATH" ]; then
    echo "错误: 创建虚拟环境失败!"
    exit 1
fi

# 激活虚拟环境
echo "激活虚拟环境..."
source "$VENV_PATH/bin/activate"

# 更新pip
echo "更新pip到最新版本..."
pip install --upgrade pip

echo "虚拟环境设置完成!"
echo "虚拟环境名称: $VENV_NAME"
echo "虚拟环境路径: $VENV_PATH"
echo "Python版本: $(python --version)"
echo "Pip版本: $(pip --version)"
echo ""
echo "可以使用 'deactivate' 命令退出虚拟环境"
echo "使用 'source $VENV_PATH/bin/activate' 重新激活虚拟环境"
