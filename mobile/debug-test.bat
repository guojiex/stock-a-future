@echo off
REM 陕西煤业数据调试测试脚本 (Windows)
echo ================================================================================
echo 陕西煤业(601225)数据调试测试
echo ================================================================================
echo.

REM 检查node是否安装
where node >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo 错误: 未找到Node.js，请先安装Node.js
    echo 下载地址: https://nodejs.org/
    pause
    exit /b 1
)

echo [1/3] 检查依赖...
if not exist "node_modules\node-fetch" (
    echo 安装node-fetch依赖...
    call npm install node-fetch
    if %ERRORLEVEL% NEQ 0 (
        echo 错误: 依赖安装失败
        pause
        exit /b 1
    )
)

echo.
echo [2/3] 测试后端API数据...
echo ================================================================================
node test-api.js
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo 错误: API测试失败
    echo 请确保后端服务运行在 http://127.0.0.1:8081
    pause
    exit /b 1
)

echo.
echo ================================================================================
echo [3/3] 测试完成
echo ================================================================================
echo.
echo 接下来的步骤:
echo 1. 启动React Native应用: npm start
echo 2. 打开浏览器开发者工具 (F12)
echo 3. 在应用中搜索 "陕西煤业" 或 "601225"
echo 4. 查看控制台日志输出
echo.
echo 详细调试指南请查看:
echo - mobile\DEBUG_SHANXI_COAL_ISSUE.md
echo - mobile\DEBUG_LOGS_GUIDE.md
echo.

pause

