# Stock-A-Future 服务器启动脚本 (AKTools)
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Stock-A-Future 服务器启动脚本 (AKTools)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "正在检查配置..." -ForegroundColor Yellow
if (-not (Test-Path "config.env")) {
    Write-Host "错误: 未找到 config.env 配置文件" -ForegroundColor Red
    Write-Host "请确保已创建 config.env 文件" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

Write-Host "配置文件检查通过" -ForegroundColor Green
Write-Host ""

Write-Host "正在启动 Stock-A-Future 服务器..." -ForegroundColor Yellow
Write-Host "数据源: AKTools" -ForegroundColor White
Write-Host "服务器端口: 8081" -ForegroundColor White
Write-Host ""

try {
    go run cmd/server/main.go
} catch {
    Write-Host "启动失败: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "服务器已停止" -ForegroundColor Yellow
Read-Host "按回车键退出"
