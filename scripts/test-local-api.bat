@echo off
echo 测试本地API服务器
echo ==================
echo.

echo 1. 健康检查
.\curl.exe http://localhost:8081/api/v1/health

echo.
echo 2. 搜索股票（平安）
.\curl.exe -X POST -d "{\"query\":\"平安\"}" http://localhost:8081/api/v1/stocks/search

echo.
echo 3. 获取股票列表
.\curl.exe http://localhost:8081/api/v1/stocks

echo.
echo 4. 获取特定股票信息（平安银行）
.\curl.exe http://localhost:8081/api/v1/stocks/000001.SZ/basic

echo.
echo 测试完成！
pause
