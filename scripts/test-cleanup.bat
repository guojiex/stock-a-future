@echo off
echo 测试Stock-A-Future数据清理功能
echo ================================

set SERVER_URL=http://localhost:8080

echo.
echo 1. 检查清理服务状态...
curl -s "%SERVER_URL%/api/v1/cleanup/status" | python -m json.tool

echo.
echo 2. 手动触发数据清理...
curl -s -X POST "%SERVER_URL%/api/v1/cleanup/manual" | python -m json.tool

echo.
echo 3. 等待5秒后再次检查状态...
timeout /t 5 /nobreak >nul
curl -s "%SERVER_URL%/api/v1/cleanup/status" | python -m json.tool

echo.
echo 4. 更新清理配置（将股票信号保留天数改为60天）...
curl -s -X PUT "%SERVER_URL%/api/v1/cleanup/config" ^
  -H "Content-Type: application/json" ^
  -d "{\"retention_days\": 60}" | python -m json.tool

echo.
echo 5. 验证配置更新...
curl -s "%SERVER_URL%/api/v1/cleanup/status" | python -m json.tool

echo.
echo 测试完成！
echo 注意：收藏数据不会被清理，只清理过期的股票信号数据
pause
