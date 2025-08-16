@echo off
chcp 65001 >nul
echo ========================================
echo    Stock-A-Future API 服务器启动脚本
echo ========================================
echo.

:: 检查Go是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到Go环境，请先安装Go 1.22或更高版本
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

:: 检查.env文件是否存在
if not exist ".env" (
    echo 警告: 未找到.env配置文件
    echo 正在从cache.env.example创建.env文件...
    copy "cache.env.example" ".env" >nul
    if %errorlevel% neq 0 (
        echo 错误: 无法创建.env文件，请手动创建
        pause
        exit /b 1
    )
    echo 已创建.env文件，请编辑其中的TUSHARE_TOKEN配置
    echo.
)

:: 检查依赖
echo 正在检查Go模块依赖...
go mod tidy
if %errorlevel% neq 0 (
    echo 错误: 依赖检查失败
    pause
    exit /b 1
)

:: 启动服务器
echo.
echo 正在启动Stock-A-Future API服务器...
echo 服务器地址: http://localhost:8080
echo 按 Ctrl+C 停止服务器
echo.

:: 运行服务器
go run cmd/server/main.go

:: 如果服务器异常退出，暂停显示错误信息
if %errorlevel% neq 0 (
    echo.
    echo 服务器异常退出，错误代码: %errorlevel%
    pause
)
