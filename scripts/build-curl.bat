@echo off
echo 正在编译Go语言版本的Curl工具...
echo.

REM 检查Go是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到Go语言环境，请先安装Go 1.22或更高版本
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

echo Go版本信息:
go version
echo.

REM 编译curl工具
echo 编译中...
go build -o curl.exe ./cmd/curl

if errorlevel 1 (
    echo 编译失败！
    pause
    exit /b 1
)

echo.
echo 编译成功！curl.exe已生成
echo.
echo 使用方法示例:
echo   curl.exe http://localhost:8080/api/stocks
echo   curl.exe -X POST -d "{\"name\":\"AAPL\"}" http://localhost:8080/api/stocks
echo   curl.exe -v http://localhost:8080/api/stocks
echo.
echo 查看帮助: curl.exe
echo.
pause
