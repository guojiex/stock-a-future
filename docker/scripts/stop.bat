@echo off
REM Stock-A-Future Docker 停止脚本 (Windows)

setlocal EnableDelayedExpansion

REM 脚本配置
set "SCRIPT_DIR=%~dp0"
set "DOCKER_DIR=%SCRIPT_DIR%.."

REM 颜色定义
set "RED=[31m"
set "GREEN=[32m"
set "YELLOW=[33m"
set "BLUE=[34m"
set "NC=[0m"

echo ⏹️ Stock-A-Future Docker 停止脚本
echo =================================

REM 停止服务
echo %BLUE%[INFO]%NC% 停止 Docker 服务...

cd /d "%DOCKER_DIR%"

REM 检查使用哪个 compose 命令
docker-compose --version >nul 2>&1
if errorlevel 1 (
    set "COMPOSE_CMD=docker compose"
) else (
    set "COMPOSE_CMD=docker-compose"
)

%COMPOSE_CMD% down

if errorlevel 1 (
    echo %RED%[ERROR]%NC% 服务停止失败
    pause
    exit /b 1
)

echo %GREEN%[SUCCESS]%NC% 服务停止成功
echo %GREEN%[SUCCESS]%NC% 所有服务已停止

pause
