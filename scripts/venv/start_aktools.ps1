# 启动AKTools服务的PowerShell脚本

# 设置虚拟环境路径
$VENV_PATH = Join-Path $PWD ".venv"

# 默认端口
$DEFAULT_PORT = 8080
$PORT = if ($args[0]) { $args[0] } else { $DEFAULT_PORT }

# 检查虚拟环境是否存在
if (-not (Test-Path $VENV_PATH)) {
    Write-Host "错误: 虚拟环境不存在! 请先运行 .\scripts\setup_venv.ps1 创建虚拟环境。" -ForegroundColor Red
    exit 1
}

# 检查是否在虚拟环境中
if (-not $env:VIRTUAL_ENV) {
    Write-Host "未检测到激活的虚拟环境，正在激活..." -ForegroundColor Yellow
    & "$VENV_PATH\Scripts\Activate.ps1"
}

# 确认已在虚拟环境中
if (-not $env:VIRTUAL_ENV) {
    Write-Host "错误: 无法激活虚拟环境!" -ForegroundColor Red
    exit 1
}

# 检查AKTools是否已安装
try {
    $null = pip show aktools
}
catch {
    Write-Host "错误: AKTools未安装! 请先运行 .\scripts\install_aktools.ps1 安装AKTools。" -ForegroundColor Red
    exit 1
}

Write-Host "当前使用的Python: $(Get-Command python | Select-Object -ExpandProperty Source)"
python --version
Write-Host "启动AKTools服务，端口: $PORT..." -ForegroundColor Cyan

# 启动AKTools服务
python -m aktools --port $PORT

# 检查启动结果
if ($LASTEXITCODE -ne 0) {
    Write-Host "错误: AKTools服务启动失败!" -ForegroundColor Red
    exit 1
}
