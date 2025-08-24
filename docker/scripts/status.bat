@echo off
REM Stock-A-Future Docker 状态检查脚本 (Windows)

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

echo 📊 Stock-A-Future Docker 状态检查
echo =================================

REM 检查服务状态
echo %BLUE%[INFO]%NC% 检查 Docker 服务状态...

cd /d "%DOCKER_DIR%"

REM 检查使用哪个 compose 命令
docker-compose --version >nul 2>&1
if errorlevel 1 (
    set "COMPOSE_CMD=docker compose"
) else (
    set "COMPOSE_CMD=docker-compose"
)

REM 显示服务状态
echo.
%COMPOSE_CMD% ps
echo.

REM 检查服务健康状态
echo %BLUE%[INFO]%NC% 检查服务健康状态...

REM 检查 AKTools
curl -s http://localhost:8080/health >nul 2>&1
if not errorlevel 1 (
    echo %GREEN%[SUCCESS]%NC% AKTools 服务正常
) else (
    echo %RED%[ERROR]%NC% AKTools 服务不可用
)

REM 检查 Stock-A-Future
curl -s http://localhost:8081/api/v1/health >nul 2>&1
if not errorlevel 1 (
    echo %GREEN%[SUCCESS]%NC% Stock-A-Future 服务正常
) else (
    echo %RED%[ERROR]%NC% Stock-A-Future 服务不可用
)

REM 显示访问信息
echo.
echo %BLUE%[INFO]%NC% 服务访问地址:
echo 📊 Stock-A-Future Web界面: http://localhost:8081
echo 🔗 Stock-A-Future API:    http://localhost:8081/api/v1/health
echo 📈 AKTools API:           http://localhost:8080
echo.

pause
