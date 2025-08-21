@echo off
REM 在虚拟环境中安装AKTools的Windows批处理脚本

REM 设置虚拟环境路径
set VENV_PATH=%CD%\.venv

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

echo 当前使用的Python: %VIRTUAL_ENV%\Scripts\python.exe
python --version

REM 安装AKTools
echo 开始安装AKTools...
pip install aktools

REM 检查安装结果
if %ERRORLEVEL% EQU 0 (
    echo AKTools安装成功!
    
    REM 显示版本信息
    echo 已安装的AKTools版本:
    pip show aktools
    
    echo.
    echo 启动AKTools服务的命令:
    echo python -m aktools
    echo.
    echo 默认情况下，AKTools将在 http://127.0.0.1:8080 启动
    echo 要指定端口，可以使用: python -m aktools --port 8080
) else (
    echo 错误: AKTools安装失败!
)
