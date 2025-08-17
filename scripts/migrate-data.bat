@echo off
echo ========================================
echo Stock-A-Future 数据迁移工具
echo ========================================
echo.

echo 正在构建迁移工具...
go build -o migrate.exe cmd/migrate/main.go
if %errorlevel% neq 0 (
    echo 构建失败，请检查Go环境
    pause
    exit /b 1
)

echo.
echo 开始数据迁移...
echo 注意：这将把JSON文件数据迁移到SQLite数据库
echo.

echo 执行迁移...
migrate.exe -data data

echo.
echo 迁移完成！
echo 数据库文件位置: data/favorites.db
echo.
echo 按任意键退出...
pause > nul
