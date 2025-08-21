@echo off
REM 启动AKTools服务的Windows批处理脚本

REM 设置虚拟环境路径
set VENV_PATH=%CD%\.venv

REM 默认端口
set DEFAULT_PORT=8080
if "%~1"=="" (
    set PORT=%DEFAULT_PORT%
) else (
    set PORT=%~1
)

REM 检查虚拟环境是否存在
if not exist "%VENV_PATH%" (
    echo 错误: 虚拟环境不存在! 请先运行 scripts\setup_venv.bat 创建虚拟环境。
    exit /b 1
)

REM 检查是否在虚拟环境中
if "%VIRTUAL_ENV%"=="" (
    echo 未检测到激活的虚拟环境，正在激活...
    call "%VENV_PATH%\Scripts\activate.bat"
)

REM 确认已在虚拟环境中
if "%VIRTUAL_ENV%"=="" (
    echo 错误: 无法激活虚拟环境!
    exit /b 1
)

REM 检查AKTools是否已安装
pip show aktools > nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo 错误: AKTools未安装! 请先运行 scripts\install_aktools.bat 安装AKTools。
    exit /b 1
)

echo 当前使用的Python: %VIRTUAL_ENV%\Scripts\python.exe
python --version
echo 启动AKTools服务，端口: %PORT%...

REM 启动AKTools服务
python -m aktools --port %PORT%

REM 检查启动结果
if %ERRORLEVEL% NEQ 0 (
    echo 错误: AKTools服务启动失败!
    exit /b 1
)
