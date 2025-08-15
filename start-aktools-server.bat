@echo off
echo ========================================
echo Stock-A-Future 服务器启动脚本 (AKTools)
echo ========================================
echo.

echo 正在检查配置...
if not exist "config.env" (
    echo 错误: 未找到 config.env 配置文件
    echo 请确保已创建 config.env 文件
    pause
    exit /b 1
)

echo 配置文件检查通过
echo.

echo 正在启动 Stock-A-Future 服务器...
echo 数据源: AKTools
echo 服务器端口: 8081
echo.

go run cmd/server/main.go

echo.
echo 服务器已停止
pause
