@echo off
echo ========================================
echo Stock-A-Future 数据库迁移启动脚本
echo ========================================
echo.

echo 正在启动服务器...
echo 注意：首次启动时会自动迁移收藏数据到数据库
echo.

echo 启动服务器...
start "Stock-A-Future Server" cmd /k "go run ./cmd/server/main.go"

echo.
echo 服务器已启动！
echo 等待数据迁移完成...
echo.
echo 迁移完成后，您可以在浏览器中访问：
echo   http://localhost:8080
echo.
echo 按任意键退出...
pause > nul
