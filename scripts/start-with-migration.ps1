Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Stock-A-Future 数据库迁移启动脚本" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "正在启动服务器..." -ForegroundColor Yellow
Write-Host "注意：首次启动时会自动迁移收藏数据到数据库" -ForegroundColor Yellow
Write-Host ""

Write-Host "启动服务器..." -ForegroundColor Green
Start-Process powershell -ArgumentList "-NoExit", "-Command", "go run ./cmd/server/main.go"

Write-Host ""
Write-Host "服务器已启动！" -ForegroundColor Green
Write-Host "等待数据迁移完成..." -ForegroundColor Yellow
Write-Host ""
Write-Host "迁移完成后，您可以在浏览器中访问：" -ForegroundColor Cyan
Write-Host "  http://localhost:8080" -ForegroundColor White
Write-Host ""
Write-Host "按任意键退出..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
