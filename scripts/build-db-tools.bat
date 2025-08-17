@echo off
echo 构建数据库工具...

echo 构建迁移工具...
go build -o migrate.exe cmd/migrate/main.go
if %errorlevel% neq 0 (
    echo 构建迁移工具失败
    exit /b 1
)

echo 构建数据库工具...
go build -o db-tools.exe cmd/db-tools/main.go
if %errorlevel% neq 0 (
    echo 构建数据库工具失败
    exit /b 1
)

echo 构建完成！
echo 可用工具:
echo   migrate.exe - 数据迁移工具
echo   db-tools.exe - 数据库管理工具
