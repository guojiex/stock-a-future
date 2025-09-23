@echo off
setlocal enabledelayedexpansion
REM Stock-A-Future Full Stack Startup Script (English Version)
REM Starts Go backend and React Native/Web frontend

echo Starting Stock-A-Future Full Stack Application...
echo.
echo Service startup order:
echo 1. AKTools Service (Data Provider) - Port 8080
echo 2. Go Backend API - Port 8081  
echo 3. Frontend Application - Port 3000
echo.

REM Check Node.js
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Node.js not found. Please install Node.js 18+
    echo Download: https://nodejs.org/
    pause
    exit /b 1
) else (
    for /f "tokens=*" %%i in ('node --version') do set NODE_VERSION=%%i
    echo OK: Node.js version: %NODE_VERSION%
)

REM Check and start AKTools service
echo.
echo ==========================================
echo STEP 1: Checking AKTools service...
echo ==========================================
curl -s http://127.0.0.1:8080/api/public/stock_zh_a_hist >nul 2>&1
if %errorlevel% equ 0 (
    echo [SUCCESS] AKTools service is already running (http://127.0.0.1:8080)
) else (
    echo [INFO] AKTools service is not running. Starting it now...
    
    REM Check if Python is available
    echo [DEBUG] Checking Python installation...
    python --version
    if %errorlevel% neq 0 (
        echo [ERROR] Python not found. Please install Python 3.8+
        echo [INFO] Continuing without AKTools (will use fallback data source)
        goto check_go_backend
    )
    echo [SUCCESS] Python is available
    
    REM Check if AKTools is installed
    echo [DEBUG] Checking AKTools installation...
    python -c "import aktools; print('AKTools version:', aktools.__version__)"
    if %errorlevel% neq 0 (
        echo [WARNING] AKTools not installed. Installing now...
        pip install aktools
        if %errorlevel% neq 0 (
            echo [ERROR] Failed to install AKTools
            echo [INFO] Continuing without AKTools (will use fallback data source)
            goto check_go_backend
        )
        echo [SUCCESS] AKTools installed successfully
    )
    echo [SUCCESS] AKTools is available
    
    echo [INFO] Starting AKTools service in background...
    echo [DEBUG] Command: python -m aktools --port 8080
    start "AKTools Service" cmd /k "python -m aktools --port 8080"
    
    REM Wait for AKTools to start
    echo [INFO] Waiting 10 seconds for AKTools to start...
    timeout /t 10 /nobreak >nul
    
    REM Check if port is listening
    echo [DEBUG] Checking if port 8080 is listening...
    netstat -an | findstr :8080
    
    REM Test AKTools API endpoint
    echo [DEBUG] Testing AKTools API endpoint...
    curl -s "http://127.0.0.1:8080/api/public/stock_zh_a_hist" >nul 2>&1
    if %errorlevel% equ 0 (
        echo [SUCCESS] AKTools service is fully ready (http://127.0.0.1:8080)
    ) else (
        echo [WARNING] AKTools API test failed, but continuing...
        echo [INFO] Go backend will test the connection during startup
    )
)

:check_go_backend
echo.
echo ==========================================
echo STEP 2: Checking Go backend service...
echo ==========================================
curl -s http://localhost:8081/api/v1/health >nul 2>&1
if %errorlevel% equ 0 (
    echo [SUCCESS] Go backend is already running (http://localhost:8081)
) else (
    echo [INFO] Go backend is not running. Starting it now...
    
    REM Check if Go is installed
    echo [DEBUG] Checking Go installation...
    go version
    if %errorlevel% neq 0 (
        echo [ERROR] Go not found. Please install Go 1.22+
        echo [INFO] Download: https://golang.org/dl/
        pause
        exit /b 1
    ) else (
        echo [SUCCESS] Go is available
    )
    
    echo [INFO] Starting Go backend server in background...
    echo [DEBUG] Command: go run cmd/server/main.go
    echo [INFO] This will test AKTools connection during startup...
    start "Stock-A-Future Backend" cmd /k "go run cmd/server/main.go"
    
    REM Wait for server to start
    echo [INFO] Waiting 5 seconds for server to start...
    timeout /t 5 /nobreak >nul
    
    REM Check if server started successfully
    echo [DEBUG] Testing Go backend health endpoint...
    curl -s http://localhost:8081/api/v1/health >nul 2>&1
    if %errorlevel% equ 0 (
        echo [SUCCESS] Go backend started successfully (http://localhost:8081)
    ) else (
        echo [WARNING] Go backend may still be starting or failed...
        echo [INFO] If it failed due to AKTools connection, check the backend terminal window
        echo [INFO] If the web app fails to load data, please wait a moment and refresh.
    )
)

echo.
echo Choose frontend to start:
echo 1. React Web App (browser-based)
echo 2. React Native Mobile App (mobile simulator)
echo 3. Both Web and Mobile
echo.
set /p choice="Enter your choice (1-3): "

if "%choice%"=="1" goto start_web
if "%choice%"=="2" goto start_mobile
if "%choice%"=="3" goto start_both
goto start_web

:start_web
echo.
echo Starting React Web application...
cd web-react
if not exist "node_modules" (
    echo Installing web dependencies...
    call npm install
    if %errorlevel% neq 0 (
        echo ERROR: Failed to install web dependencies
        pause
        exit /b 1
    )
)
echo Web app will open in browser: http://localhost:3000
call npm start
goto end

:start_mobile
echo.
echo Starting React Native mobile application...
cd mobile
if not exist "node_modules" (
    echo Installing mobile dependencies...
    call npm install
    if %errorlevel% neq 0 (
        echo ERROR: Failed to install mobile dependencies
        pause
        exit /b 1
    )
)
echo Starting React Native Metro bundler...
echo Use 'npm run android' or 'npm run ios' in another terminal to run on device/simulator
call npm start
goto end

:start_both
echo.
echo Starting both Web and Mobile applications...

REM Start Web App
echo Starting React Web App...
cd web-react
if not exist "node_modules" (
    echo Installing web dependencies...
    call npm install
    if %errorlevel% neq 0 (
        echo ERROR: Failed to install web dependencies
        pause
        exit /b 1
    )
)
start "React Web App" cmd /k "npm start"
cd ..

REM Start Mobile App
echo Starting React Native Mobile App...
cd mobile
if not exist "node_modules" (
    echo Installing mobile dependencies...
    call npm install
    if %errorlevel% neq 0 (
        echo ERROR: Failed to install mobile dependencies
        pause
        exit /b 1
    )
)
echo.
echo All services are starting:
echo - Web App: http://localhost:3000
echo - Mobile Metro: http://localhost:8081 (Metro bundler)
echo - Go Backend: http://localhost:8081 (API)
echo - AKTools Service: http://127.0.0.1:8080 (Data provider)
echo.
echo To run on mobile device/simulator, use:
echo   npm run android  (Android)
echo   npm run ios      (iOS)
echo.
call npm start

:end
