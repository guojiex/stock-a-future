@echo off
REM 策略排序调试测试脚本
REM 用于快速启动服务器并打开浏览器进行测试

echo ========================================
echo 策略排序调试测试
echo ========================================
echo.

REM 检查服务器是否已编译
if not exist "bin\server.exe" (
    echo [步骤 1/3] 编译服务器...
    go build -o bin\server.exe .\cmd\server
    if errorlevel 1 (
        echo 编译失败！
        pause
        exit /b 1
    )
    echo 编译成功！
    echo.
) else (
    echo [步骤 1/3] 服务器已编译，跳过编译步骤
    echo.
)

echo [步骤 2/3] 启动后端服务器...
echo 请查看服务器控制台日志，观察"排序前顺序"和"排序后顺序"
echo.
start "Stock-A-Future Server - 策略排序调试" cmd /k ".\bin\server.exe"

REM 等待服务器启动
timeout /t 3 /nobreak >nul

echo [步骤 3/3] 打开浏览器...
echo.
echo ========================================
echo 调试指南：
echo ========================================
echo.
echo 1. 按 F12 打开浏览器开发者工具
echo 2. 切换到 Console (控制台) 标签
echo 3. 导航到策略管理页面
echo 4. 多次点击"刷新"按钮 (建议 5-10 次)
echo 5. 观察控制台日志：
echo    - 每次应该显示"策略顺序"
echo    - 如果顺序变化，会显示警告⚠️
echo    - 查看"原始API响应"确认数据
echo.
echo 6. 观察后端服务器日志：
echo    - "排序前顺序" - 可能每次不同(正常)
echo    - "排序后顺序" - 应该每次相同(关键)
echo    - "返回策略列表" - 最终返回的顺序
echo.
echo 详细调试指南: docs\debugging\STRATEGY_SORTING_DEBUG_GUIDE.md
echo.
echo ========================================

REM 等待用户确认
echo 按任意键打开浏览器并查看调试指南...
pause >nul

REM 打开浏览器
start http://localhost:3000/#/strategies

REM 打开调试指南
start notepad docs\debugging\STRATEGY_SORTING_DEBUG_GUIDE.md

echo.
echo 测试环境已准备就绪！
echo.
echo 提示：
echo - 如需停止服务器，关闭服务器窗口即可
echo - 前端日志在浏览器控制台
echo - 后端日志在服务器窗口
echo.
pause

