@echo off
REM Stock-A-Future Docker 启动脚本 (Windows)
REM 用于一键启动 AKTools 和 Golang 程序

setlocal EnableDelayedExpansion

REM 脚本配置
set "SCRIPT_DIR=%~dp0"
set "DOCKER_DIR=%SCRIPT_DIR%.."
set "PROJECT_ROOT=%DOCKER_DIR%\.."

REM 颜色定义 (Windows 10+ 支持 ANSI 颜色)
set "RED=[31m"
set "GREEN=[32m"
set "YELLOW=[33m"
set "BLUE=[34m"
set "NC=[0m"

echo 🚀 Stock-A-Future Docker 启动脚本
echo ==================================

REM 检查 Docker 环境
echo %BLUE%[INFO]%NC% 检查 Docker 环境...

docker --version >nul 2>&1
if errorlevel 1 (
    echo %RED%[ERROR]%NC% Docker 未安装或不在 PATH 中
    pause
    exit /b 1
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    docker compose version >nul 2>&1
    if errorlevel 1 (
        echo %RED%[ERROR]%NC% Docker Compose 未安装或不在 PATH 中
        pause
        exit /b 1
    )
    set "COMPOSE_CMD=docker compose"
) else (
    set "COMPOSE_CMD=docker-compose"
)

docker info >nul 2>&1
if errorlevel 1 (
    echo %RED%[ERROR]%NC% Docker 守护进程未运行
    pause
    exit /b 1
)

echo %GREEN%[SUCCESS]%NC% Docker 环境检查通过

REM 创建必要的目录
echo %BLUE%[INFO]%NC% 创建必要的目录...

if not exist "%DOCKER_DIR%\volumes\data" mkdir "%DOCKER_DIR%\volumes\data"
if not exist "%DOCKER_DIR%\volumes\logs" mkdir "%DOCKER_DIR%\volumes\logs"
if not exist "%DOCKER_DIR%\volumes\aktools-data" mkdir "%DOCKER_DIR%\volumes\aktools-data"
if not exist "%DOCKER_DIR%\volumes\aktools-logs" mkdir "%DOCKER_DIR%\volumes\aktools-logs"

echo %GREEN%[SUCCESS]%NC% 目录创建完成

REM 检查端口占用
echo %BLUE%[INFO]%NC% 检查端口占用情况...

netstat -an | findstr ":8080" >nul 2>&1
if not errorlevel 1 (
    echo %YELLOW%[WARNING]%NC% 端口 8080 已被占用
    set "PORT_WARNING=1"
)

netstat -an | findstr ":8081" >nul 2>&1
if not errorlevel 1 (
    echo %YELLOW%[WARNING]%NC% 端口 8081 已被占用
    set "PORT_WARNING=1"
)

if defined PORT_WARNING (
    echo %YELLOW%[WARNING]%NC% 请确保端口可用，或修改 docker-compose.yml 中的端口映射
    set /p "CONTINUE=是否继续启动? (y/N): "
    if /i not "!CONTINUE!"=="y" (
        echo %BLUE%[INFO]%NC% 启动已取消
        pause
        exit /b 0
    )
) else (
    echo %GREEN%[SUCCESS]%NC% 端口检查通过
)

REM 启动服务
echo %BLUE%[INFO]%NC% 启动 Docker 服务...

cd /d "%DOCKER_DIR%"

%COMPOSE_CMD% up --build -d

if errorlevel 1 (
    echo %RED%[ERROR]%NC% 服务启动失败
    pause
    exit /b 1
)

echo %GREEN%[SUCCESS]%NC% 服务启动成功

REM 等待服务就绪
echo %BLUE%[INFO]%NC% 等待服务就绪...

REM 等待 AKTools 服务
echo %BLUE%[INFO]%NC% 等待 AKTools 服务启动...
set /a "attempt=0"
set /a "max_attempts=30"

:wait_aktools
if !attempt! geq !max_attempts! (
    echo %RED%[ERROR]%NC% AKTools 服务启动超时
    goto show_logs
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
    goto show_logs
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

REM 显示访问信息
echo.
echo %GREEN%[SUCCESS]%NC% === 服务启动完成 ===
echo.
echo 📊 Stock-A-Future Web界面: http://localhost:8081
echo 🔗 Stock-A-Future API:    http://localhost:8081/api/v1/health
echo 📈 AKTools API:           http://localhost:8080
echo.
echo 📋 常用 API 端点:
echo    健康检查: curl http://localhost:8081/api/v1/health
echo    股票信息: curl http://localhost:8081/api/v1/stocks/000001/basic
echo    日线数据: curl http://localhost:8081/api/v1/stocks/000001/daily
echo.
echo 📝 查看日志: %SCRIPT_DIR%logs.bat
echo ⏹️  停止服务: %SCRIPT_DIR%stop.bat
echo 🔄 重新构建: %SCRIPT_DIR%rebuild.bat
echo.

echo %GREEN%[SUCCESS]%NC% 启动完成！
pause
exit /b 0

:show_logs
echo %BLUE%[INFO]%NC% 显示最近的服务日志:
%COMPOSE_CMD% logs --tail=20
pause
exit /b 1
