# 创建和激活Python虚拟环境的PowerShell脚本

# 设置虚拟环境名称和路径
$VENV_NAME = "stock_env"
$VENV_PATH = Join-Path $PWD ".venv"

# 检查Python版本
Write-Host "检查Python版本..."
try {
    $pythonVersion = python --version
    Write-Host $pythonVersion
}
catch {
    Write-Host "错误: 未找到Python，请确保已安装Python 3.7+" -ForegroundColor Red
    exit 1
}

# 检查是否已存在虚拟环境
if (Test-Path $VENV_PATH) {
    Write-Host "发现已存在的虚拟环境: $VENV_PATH"
    $RECREATE = Read-Host "是否重新创建虚拟环境? (y/n)"
    if ($RECREATE -eq "y" -or $RECREATE -eq "Y") {
        Write-Host "删除现有虚拟环境..."
        Remove-Item -Recurse -Force $VENV_PATH
    }
    else {
        Write-Host "使用现有虚拟环境..."
        & "$VENV_PATH\Scripts\Activate.ps1"
        Write-Host "虚拟环境已激活: $VENV_NAME" -ForegroundColor Green
        Write-Host "可以使用 'deactivate' 命令退出虚拟环境"
        exit 0
    }
}

# 创建虚拟环境
Write-Host "创建Python虚拟环境: $VENV_NAME..."
python -m venv $VENV_PATH

# 检查虚拟环境是否创建成功
if (-not (Test-Path $VENV_PATH)) {
    Write-Host "错误: 创建虚拟环境失败!" -ForegroundColor Red
    exit 1
}

# 激活虚拟环境
Write-Host "激活虚拟环境..."
& "$VENV_PATH\Scripts\Activate.ps1"

# 更新pip
Write-Host "更新pip到最新版本..."
python -m pip install --upgrade pip

Write-Host "虚拟环境设置完成!" -ForegroundColor Green
Write-Host "虚拟环境名称: $VENV_NAME"
Write-Host "虚拟环境路径: $VENV_PATH"
python --version
pip --version
Write-Host ""
Write-Host "可以使用 'deactivate' 命令退出虚拟环境"
Write-Host "使用 '& $VENV_PATH\Scripts\Activate.ps1' 重新激活虚拟环境"
