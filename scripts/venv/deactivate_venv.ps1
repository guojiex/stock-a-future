# 停用Python虚拟环境的PowerShell脚本

# 检查是否在虚拟环境中
if (-not $env:VIRTUAL_ENV) {
    Write-Host "错误: 当前未激活任何虚拟环境!" -ForegroundColor Red
    exit 1
}

# 显示当前虚拟环境信息
Write-Host "当前虚拟环境: $env:VIRTUAL_ENV"

# 停用虚拟环境
Write-Host "停用虚拟环境..."
deactivate

# 检查是否成功停用
if (-not $env:VIRTUAL_ENV) {
    Write-Host "虚拟环境已成功停用!" -ForegroundColor Green
}
else {
    Write-Host "警告: 虚拟环境可能未正确停用。" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "如需重新激活虚拟环境，请运行:"
Write-Host "& .venv\Scripts\Activate.ps1 或 .\scripts\setup_venv.ps1"
