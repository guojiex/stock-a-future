# 在虚拟环境中安装AKTools的PowerShell脚本

# 设置虚拟环境路径
$VENV_PATH = Join-Path $PWD ".venv"

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

Write-Host "当前使用的Python: $(Get-Command python | Select-Object -ExpandProperty Source)"
python --version

# 安装AKTools
Write-Host "开始安装AKTools..." -ForegroundColor Cyan
pip install aktools

# 检查安装结果
if ($LASTEXITCODE -eq 0) {
    Write-Host "AKTools安装成功!" -ForegroundColor Green
    
    # 显示版本信息
    Write-Host "已安装的AKTools版本:" -ForegroundColor Cyan
    pip show aktools
    
    Write-Host ""
    Write-Host "启动AKTools服务的命令:" -ForegroundColor Cyan
    Write-Host "python -m aktools"
    Write-Host ""
    Write-Host "默认情况下，AKTools将在 http://127.0.0.1:8080 启动"
    Write-Host "要指定端口，可以使用: python -m aktools --port 8080"
}
else {
    Write-Host "错误: AKTools安装失败!" -ForegroundColor Red
}
