@echo off
setlocal enabledelayedexpansion
REM Stock-A-Future Full Stack Startup Script (English Version)
REM Starts Go backend and React Native/Web frontend

echo Starting Stock-A-Future Full Stack Application...
echo.
echo Service startup order:
echo 1. AKTools Service (Data Provider) - Port 8080
echo 2. Build Go Backend API (compile to bin/server.exe)
echo 3. Start Go Backend API - Port 8081  
echo 4. Frontend Application - Port 3000
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

REM Start AKTools service
echo.
echo ==========================================
echo STEP 1: Starting AKTools service...
echo ==========================================
echo Starting AKTools on port 8080...
start "AKTools Service" cmd /k "python -m aktools --port 8080"

echo Waiting 8 seconds for AKTools to initialize...
timeout /t 8 /nobreak >nul
echo AKTools should now be running on http://127.0.0.1:8080

:build_go_backend
echo.
echo ==========================================
echo STEP 2: Building Go backend service...
echo ==========================================

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

REM Build the Go backend
echo [INFO] Building Go backend server...
echo [DEBUG] Command: go build -o bin/server.exe cmd/server/main.go
go build -o bin/server.exe cmd/server/main.go
if %errorlevel% neq 0 (
    echo [ERROR] Failed to build Go backend
    echo [INFO] Please check the Go code for compilation errors
    pause
    exit /b 1
) else (
    echo [SUCCESS] Go backend built successfully: bin/server.exe
)

:check_go_backend
echo.
echo ==========================================
echo STEP 3: Checking Go backend service...
echo ==========================================
curl -s http://localhost:8081/api/v1/health >nul 2>&1
if %errorlevel% equ 0 (
    echo [SUCCESS] Go backend is already running (http://localhost:8081)
) else (
    echo [INFO] Go backend is not running. Starting it now...
    
    echo [INFO] Starting Go backend server in background...
    echo [DEBUG] Command: bin\server.exe
    echo [INFO] This will test AKTools connection during startup...
    start "Stock-A-Future Backend" cmd /k "bin\server.exe"
    
    REM Wait for server to start
    echo [INFO] Waiting 8 seconds for server to start...
    timeout /t 8 /nobreak >nul
    
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
echo 1. React Web App (Development Mode - Hot Reload)
echo 2. React Web App (Production Mode - Faster Startup)
echo 3. React Native Mobile App (Development Mode)
echo 4. React Native Mobile App (Production Mode)
echo 5. Both Web and Mobile (Development Mode)
echo 6. Both Web and Mobile (Production Mode)
echo.
set /p choice="Enter your choice (1-6): "

if "%choice%"=="1" goto start_web_dev
if "%choice%"=="2" goto start_web_prod
if "%choice%"=="3" goto start_mobile_dev
if "%choice%"=="4" goto start_mobile_prod
if "%choice%"=="5" goto start_both_dev
if "%choice%"=="6" goto start_both_prod
goto start_web_dev

:start_web_dev
echo.
echo Starting React Web application (Development Mode)...
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
echo Development mode: Hot reload enabled, slower startup
call npm start
goto end

:start_web_prod
echo.
echo Starting React Web application (Production Mode)...
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
echo Building production version with memory optimization...
call npm run build
if %errorlevel% neq 0 (
    echo ERROR: Failed to build web application
    echo The build process ran out of memory. Try closing other applications.
    pause
    exit /b 1
)
echo Web app will open in browser: http://localhost:3000
echo Production mode: Optimized build, faster startup, no hot reload
call npm run serve
goto end

:start_mobile_dev
echo.
echo Starting React Native mobile application (Development Mode)...
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
echo Development mode: Hot reload enabled, slower startup
echo Use 'npm run android' or 'npm run ios' in another terminal to run on device/simulator
call npm start
goto end

:start_mobile_prod
echo.
echo Starting React Native mobile application (Production Mode)...
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
echo Building production bundles...
echo Creating Android bundle...
call npm run bundle:android
echo Creating iOS bundle...
call npm run bundle:ios
echo Production mode: Optimized bundles, faster app performance
echo Use 'npm run android:release' or 'npm run ios:release' to run production builds
echo Starting Metro bundler for additional development...
call npm start
goto end

:start_both_dev
echo.
echo Starting both Web and Mobile applications (Development Mode)...

REM Start Web App
echo Starting React Web App (Development)...
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
start "React Web App (Dev)" cmd /k "npm start"
cd ..

REM Start Mobile App
echo Starting React Native Mobile App (Development)...
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
echo - Web App: http://localhost:3000 (Development Mode)
echo - Mobile Metro: http://localhost:8081 (Metro bundler - Development)
echo - Go Backend: http://localhost:8081 (API)
echo - AKTools Service: http://127.0.0.1:8080 (Data provider)
echo.
echo To run on mobile device/simulator, use:
echo   npm run android  (Android - Development)
echo   npm run ios      (iOS - Development)
echo.
call npm start
goto end

:start_both_prod
echo.
echo Starting both Web and Mobile applications (Production Mode)...

REM Start Web App
echo Building and starting React Web App (Production)...
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
echo Building production web app with memory optimization...
call npm run build
if %errorlevel% neq 0 (
    echo ERROR: Failed to build web application
    echo The build process ran out of memory. Try closing other applications.
    pause
    exit /b 1
)
start "React Web App (Prod)" cmd /k "npm run serve"
cd ..

REM Start Mobile App
echo Building React Native Mobile App (Production)...
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
echo Building production mobile bundles...
call npm run bundle:android
call npm run bundle:ios
echo.
echo All services are starting:
echo - Web App: http://localhost:3000 (Production Mode)
echo - Mobile Metro: http://localhost:8081 (Metro bundler)
echo - Go Backend: http://localhost:8081 (API)
echo - AKTools Service: http://127.0.0.1:8080 (Data provider)
echo.
echo To run on mobile device/simulator, use:
echo   npm run android:release  (Android - Production)
echo   npm run ios:release      (iOS - Production)
echo.
call npm start

:end
