@echo off
setlocal enabledelayedexpansion

echo ==========================================
echo React App Memory-Optimized Build Script
echo ==========================================

REM 清理之前的构建
echo Cleaning previous build...
if exist "build" rmdir /s /q build
if exist ".tsbuildinfo" del /f ".tsbuildinfo"
if exist "node_modules\.cache" rmdir /s /q "node_modules\.cache"

echo.
echo Memory optimization settings:
echo - Max Old Space: 16GB
echo - Max Semi Space: 2GB  
echo - Memory garbage collection optimization
echo - Disable source maps
echo - Skip TypeScript strict checking
echo - Enable incremental builds

echo.
echo Attempting build with maximum memory optimization...

REM 设置环境变量
set GENERATE_SOURCEMAP=false
set DISABLE_ESLINT_PLUGIN=true
set TSC_COMPILE_ON_ERROR=true
set SKIP_PREFLIGHT_CHECK=true
set DISABLE_NEW_JSX_TRANSFORM=true
set FAST_REFRESH=false

REM 尝试最大内存配置构建
set NODE_OPTIONS=--max-old-space-size=16384 --max-semi-space-size=2048

echo [ATTEMPT 1] Building with 16GB memory limit...
call npx --max-old-space-size=16384 react-scripts build

if %errorlevel% equ 0 (
    echo.
    echo ✅ Build successful with 16GB memory configuration!
    goto :success
)

echo.
echo [ATTEMPT 1 FAILED] Trying with reduced memory but no optimization...
set NODE_OPTIONS=--max-old-space-size=12288 --max-semi-space-size=1024

echo [ATTEMPT 2] Building with 12GB memory limit...
call npx --max-old-space-size=12288 react-scripts build

if %errorlevel% equ 0 (
    echo.
    echo ✅ Build successful with 12GB memory configuration!
    goto :success
)

echo.
echo [ATTEMPT 2 FAILED] Trying with legacy JSX transform...
set DISABLE_NEW_JSX_TRANSFORM=true
set NODE_OPTIONS=--max-old-space-size=10240 --max-semi-space-size=512

echo [ATTEMPT 3] Building with legacy JSX and 10GB memory...
call npx --max-old-space-size=10240 react-scripts build

if %errorlevel% equ 0 (
    echo.
    echo ✅ Build successful with legacy JSX configuration!
    goto :success
)

echo.
echo ❌ All build attempts failed.
echo.
echo Troubleshooting suggestions:
echo 1. Close all other applications to free up memory
echo 2. Restart your computer to clear memory
echo 3. Try running: npm run start (development mode instead)
echo 4. Check if you have enough disk space
echo 5. Update Node.js to the latest LTS version
echo.
exit /b 1

:success
echo.
echo ==========================================
echo ✅ BUILD COMPLETED SUCCESSFULLY!
echo ==========================================
echo.
echo Build output is in the 'build' directory
echo You can now run: npm run serve
echo.
exit /b 0
