# Stock-A-Future API 服务器启动脚本
# 使用PowerShell启动Go服务器

Write-Host "=== Stock-A-Future API 服务器启动脚本 ===" -ForegroundColor Green
Write-Host ""

# 检查Go是否安装
try {
    $goVersion = go version
    Write-Host "✓ Go已安装: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go未安装或不在PATH中" -ForegroundColor Red
    Write-Host "请先安装Go 1.22或更高版本" -ForegroundColor Yellow
    exit 1
}

# 检查环境变量
Write-Host "检查环境变量配置..." -ForegroundColor Yellow

$envFile = ".env"
if (Test-Path $envFile) {
    Write-Host "✓ 找到.env配置文件" -ForegroundColor Green
    Get-Content $envFile | ForEach-Object {
        if ($_ -match "^TUSHARE_TOKEN=") {
            $token = $_.Split("=", 2)[1]
            if ($token -and $token -ne "") {
                Write-Host "✓ TUSHARE_TOKEN已配置" -ForegroundColor Green
            } else {
                Write-Host "❌ TUSHARE_TOKEN为空" -ForegroundColor Red
            }
        }
    }
} else {
    Write-Host "⚠ 未找到.env文件，将使用系统环境变量" -ForegroundColor Yellow
}

# 检查TUSHARE_TOKEN环境变量
$tushareToken = $env:TUSHARE_TOKEN
if (-not $tushareToken) {
    Write-Host "❌ TUSHARE_TOKEN环境变量未设置" -ForegroundColor Red
    Write-Host "请设置TUSHARE_TOKEN环境变量或创建.env文件" -ForegroundColor Yellow
    Write-Host "示例: .env文件内容:" -ForegroundColor Cyan
    Write-Host "TUSHARE_TOKEN=your_tushare_token_here" -ForegroundColor Cyan
    Write-Host "TUSHARE_BASE_URL=http://api.tushare.pro" -ForegroundColor Cyan
    Write-Host "SERVER_PORT=8080" -ForegroundColor Cyan
    exit 1
} else {
    Write-Host "✓ TUSHARE_TOKEN已设置" -ForegroundColor Green
}

# 检查项目依赖
Write-Host "检查Go模块依赖..." -ForegroundColor Yellow
try {
    go mod tidy
    Write-Host "✓ Go模块依赖已更新" -ForegroundColor Green
} catch {
    Write-Host "❌ 更新Go模块依赖失败" -ForegroundColor Red
    exit 1
}

# 构建项目
Write-Host "构建项目..." -ForegroundColor Yellow
try {
    go build -o bin/server.exe ./cmd/server
    Write-Host "✓ 项目构建成功" -ForegroundColor Green
} catch {
    Write-Host "❌ 项目构建失败" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=== 启动服务器 ===" -ForegroundColor Green
Write-Host "服务器启动后将自动测试Tushare API连接..." -ForegroundColor Cyan
Write-Host "如果连接失败，服务器将无法启动" -ForegroundColor Yellow
Write-Host ""

# 启动服务器
try {
    & .\bin\server.exe
} catch {
    Write-Host "❌ 服务器启动失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
