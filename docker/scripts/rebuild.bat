@echo off
REM Stock-A-Future Docker 重新构建脚本 (Windows)

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

echo 🔄 Stock-A-Future Docker 重新构建脚本
echo =====================================

REM 重新构建服务
echo %BLUE%[INFO]%NC% 重新构建 Docker 服务...

cd /d "%DOCKER_DIR%"

REM 检查使用哪个 compose 命令
docker-compose --version >nul 2>&1
if errorlevel 1 (
    set "COMPOSE_CMD=docker compose"
) else (
    set "COMPOSE_CMD=docker-compose"
)

REM 停止现有服务
echo %BLUE%[INFO]%NC% 停止现有服务...
%COMPOSE_CMD% down

REM 清理旧镜像
echo %BLUE%[INFO]%NC% 清理旧镜像...
%COMPOSE_CMD% down --rmi local

REM 重新构建并启动
echo %BLUE%[INFO]%NC% 重新构建并启动服务...
%COMPOSE_CMD% up --build -d

if errorlevel 1 (
    echo %RED%[ERROR]%NC% 服务重新构建失败
    pause
    exit /b 1
)

echo %GREEN%[SUCCESS]%NC% 服务重新构建成功

REM 等待服务就绪
echo %BLUE%[INFO]%NC% 等待服务就绪...

REM 等待 AKTools 服务
echo %BLUE%[INFO]%NC% 等待 AKTools 服务启动...
set /a "attempt=0"
set /a "max_attempts=30"

:wait_aktools
if !attempt! geq !max_attempts! (
    echo %RED%[ERROR]%NC% AKTools 服务启动超时
    pause
    exit /b 1
)

curl -s http://localhost:8080/health >nul 2>&1
if not errorlevel 1 (
    echo %GREEN%[SUCCESS]%NC% AKTools 服务已就绪
    goto wait_stock_future
)

set /a "attempt+=1"
echo|set /p="."
timeout /t 2 /nobreak >nul
goto wait_aktools

:wait_stock_future
REM 等待 Stock-A-Future 服务
echo %BLUE%[INFO]%NC% 等待 Stock-A-Future 服务启动...
set /a "attempt=0"

:wait_stock_future_loop
if !attempt! geq !max_attempts! (
    echo %RED%[ERROR]%NC% Stock-A-Future 服务启动超时
    pause
    exit /b 1
)

curl -s http://localhost:8081/api/v1/health >nul 2>&1
if not errorlevel 1 (
    echo %GREEN%[SUCCESS]%NC% Stock-A-Future 服务已就绪
    goto show_status
)

set /a "attempt+=1"
echo|set /p="."
timeout /t 2 /nobreak >nul
goto wait_stock_future_loop

:show_status
REM 显示服务状态
echo %BLUE%[INFO]%NC% 服务状态:
%COMPOSE_CMD% ps

echo.
echo %GREEN%[SUCCESS]%NC% 重新构建完成！
echo.
echo 📊 Stock-A-Future Web界面: http://localhost:8081
echo 🔗 Stock-A-Future API:    http://localhost:8081/api/v1/health
echo 📈 AKTools API:           http://localhost:8080

pause
