Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Stock-A-Future 数据迁移工具" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "正在构建迁移工具..." -ForegroundColor Yellow
go build -o migrate.exe cmd/migrate/main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "构建失败，请检查Go环境" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

Write-Host ""
Write-Host "开始数据迁移..." -ForegroundColor Green
Write-Host "注意：这将把JSON文件数据迁移到SQLite数据库" -ForegroundColor Yellow
Write-Host ""

Write-Host "执行迁移..." -ForegroundColor Green
./migrate.exe -data data

Write-Host ""
Write-Host "迁移完成！" -ForegroundColor Green
Write-Host "数据库文件位置: data/favorites.db" -ForegroundColor Cyan
Write-Host ""
Write-Host "按回车键退出..." -ForegroundColor Gray
Read-Host
