@echo off
REM 创建和激活Python虚拟环境的Windows批处理脚本

REM 设置虚拟环境名称和路径
set VENV_NAME=stock_env
set VENV_PATH=%CD%\.venv

REM 检查Python版本
echo 检查Python版本...
python --version 2>NUL
if %ERRORLEVEL% NEQ 0 (
    echo 错误: 未找到Python，请确保已安装Python 3.7+
    exit /b 1
)

REM 检查是否已存在虚拟环境
if exist "%VENV_PATH%" (
    echo 发现已存在的虚拟环境: %VENV_PATH%
    set /p RECREATE=是否重新创建虚拟环境? (y/n): 
    if /i "%RECREATE%"=="y" (
        echo 删除现有虚拟环境...
        rmdir /s /q "%VENV_PATH%"
    ) else (
        echo 使用现有虚拟环境...
        call "%VENV_PATH%\Scripts\activate.bat"
        echo 虚拟环境已激活: %VENV_NAME%
        echo 可以使用 'deactivate' 命令退出虚拟环境
        exit /b 0
    )
)

REM 创建虚拟环境
echo 创建Python虚拟环境: %VENV_NAME%...
python -m venv "%VENV_PATH%"

REM 检查虚拟环境是否创建成功
if not exist "%VENV_PATH%" (
    echo 错误: 创建虚拟环境失败!
    exit /b 1
)

REM 激活虚拟环境
echo 激活虚拟环境...
call "%VENV_PATH%\Scripts\activate.bat"

REM 更新pip
echo 更新pip到最新版本...
python -m pip install --upgrade pip

echo 虚拟环境设置完成!
echo 虚拟环境名称: %VENV_NAME%
echo 虚拟环境路径: %VENV_PATH%
python --version
pip --version
echo.
echo 可以使用 'deactivate' 命令退出虚拟环境
echo 使用 '%VENV_PATH%\Scripts\activate.bat' 重新激活虚拟环境
