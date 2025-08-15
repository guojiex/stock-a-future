@echo off
echo ========================================
echo AKTools 连接测试工具
echo ========================================
echo.

echo 正在测试 AKTools API 连接...
echo.

go run cmd/aktools-test/main.go

echo.
echo 测试完成
pause
