Write-Host "构建数据库工具..." -ForegroundColor Green

Write-Host "构建迁移工具..." -ForegroundColor Yellow
go build -o migrate.exe cmd/migrate/main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "构建迁移工具失败" -ForegroundColor Red
    exit 1
}

Write-Host "构建数据库工具..." -ForegroundColor Yellow
go build -o db-tools.exe cmd/db-tools/main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "构建数据库工具失败" -ForegroundColor Red
    exit 1
}

Write-Host "构建完成！" -ForegroundColor Green
Write-Host "可用工具:" -ForegroundColor Cyan
Write-Host "  migrate.exe - 数据迁移工具" -ForegroundColor White
Write-Host "  db-tools.exe - 数据库管理工具" -ForegroundColor White
