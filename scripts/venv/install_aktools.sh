#!/bin/bash
# 在虚拟环境中安装AKTools的脚本

# 设置虚拟环境路径
VENV_PATH="$(pwd)/.venv"

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

echo "当前使用的Python: $(which python)"
echo "Python版本: $(python --version)"

# 安装AKTools
echo "开始安装AKTools..."
pip install aktools

# 检查安装结果
if [ $? -eq 0 ]; then
    echo "AKTools安装成功!"
    
    # 显示版本信息
    echo "已安装的AKTools版本:"
    pip show aktools
    
    echo ""
    echo "启动AKTools服务的命令:"
    echo "python -m aktools"
    echo ""
    echo "默认情况下，AKTools将在 http://127.0.0.1:8080 启动"
    echo "要指定端口，可以使用: python -m aktools --port 8080"
else
    echo "错误: AKTools安装失败!"
fi
