@echo off
REM 停用Python虚拟环境的Windows批处理脚本

REM 检查是否在虚拟环境中
if "%VIRTUAL_ENV%"=="" (
    echo 错误: 当前未激活任何虚拟环境!
    exit /b 1
)

REM 显示当前虚拟环境信息
echo 当前虚拟环境: %VIRTUAL_ENV%

REM 停用虚拟环境
echo 停用虚拟环境...
call deactivate

REM 检查是否成功停用
if "%VIRTUAL_ENV%"=="" (
    echo 虚拟环境已成功停用!
) else (
    echo 警告: 虚拟环境可能未正确停用。
)

echo.
echo 如需重新激活虚拟环境，请运行:
echo .venv\Scripts\activate.bat 或 scripts\setup_venv.bat
