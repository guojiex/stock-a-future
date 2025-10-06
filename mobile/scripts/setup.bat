@echo off
REM Stock-A-Future Mobile App 安装脚本 (Windows版本)

echo 🚀 开始设置 Stock-A-Future Mobile App...

REM 检查Node.js
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 未检测到Node.js，请先安装Node.js 18+
    pause
    exit /b 1
) else (
    echo ✅ Node.js检查通过
)

REM 安装依赖
echo 📦 安装依赖包...
call npm install

if %errorlevel% neq 0 (
    echo ❌ 依赖安装失败
    pause
    exit /b 1
) else (
    echo ✅ 依赖安装成功
)

REM 检查Android环境
if defined ANDROID_HOME (
    echo ✅ 检测到Android SDK: %ANDROID_HOME%
) else (
    echo ⚠️  未检测到ANDROID_HOME环境变量
    echo    请确保已安装Android Studio并设置环境变量
)

echo.
echo 🎉 设置完成！
echo.
echo 📱 运行应用：
echo    npm run android  # 运行Android版本
echo    npm start        # 启动Metro bundler
echo.
echo 🔧 开发工具：
echo    npm run type-check  # TypeScript类型检查
echo    npm run lint        # ESLint代码检查
echo.
echo 📖 更多信息请查看 README.md
pause
