#!/bin/bash
# 启动AKTools服务的脚本

# 设置虚拟环境路径
VENV_PATH="$(pwd)/.venv"

# 默认端口
DEFAULT_PORT=8080
PORT=${1:-$DEFAULT_PORT}

# 检查虚拟环境是否存在
if [ ! -d "$VENV_PATH" ]; then
    echo "错误: 虚拟环境不存在! 请先运行 ./scripts/setup_venv.sh 创建虚拟环境。"
    exit 1
fi

# 检查是否在虚拟环境中
if [ -z "$VIRTUAL_ENV" ]; then
    echo "未检测到激活的虚拟环境，正在激活..."
    source "$VENV_PATH/bin/activate"
fi

# 确认已在虚拟环境中
if [ -z "$VIRTUAL_ENV" ]; then
    echo "错误: 无法激活虚拟环境!"
    exit 1
fi

# 检查AKTools是否已安装
if ! pip show aktools > /dev/null 2>&1; then
    echo "错误: AKTools未安装! 请先运行 ./scripts/install_aktools.sh 安装AKTools。"
    exit 1
fi

echo "当前使用的Python: $(which python)"
echo "Python版本: $(python --version)"
echo "启动AKTools服务，端口: $PORT..."

# 启动AKTools服务
python -m aktools --port $PORT

# 检查启动结果
if [ $? -ne 0 ]; then
    echo "错误: AKTools服务启动失败!"
    exit 1
fi
