# PowerShell脚本 - 启动Go API服务器
Write-Host "启动Go API服务器在8081端口..." -ForegroundColor Green

# 设置环境变量
$env:SERVER_PORT = "8081"

# 启动Go服务器
go run cmd/server/main.go
