@echo off
chcp 65001 >nul
echo 开始下载ECharts到本地...

REM 创建目标目录
if not exist "web\static\js\lib\echarts" (
    mkdir "web\static\js\lib\echarts"
    echo 创建目录: web\static\js\lib\echarts
)

REM 下载ECharts文件
echo 正在下载ECharts...
powershell -Command "& {Invoke-WebRequest -Uri 'https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js' -OutFile 'web\static\js\lib\echarts\echarts.min.js'}"

if exist "web\static\js\lib\echarts\echarts.min.js" (
    echo 下载完成!
    echo 文件路径: web\static\js\lib\echarts\echarts.min.js
    
    REM 显示文件大小
    for %%A in ("web\static\js\lib\echarts\echarts.min.js") do echo 文件大小: %%~zA 字节
    
    echo.
    echo 下一步操作:
    echo 1. 修改 index.html 中的script标签
    echo 2. 将 src='https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js'
    echo    改为 src='js/lib/echarts/echarts.min.js'
    echo 3. 测试图表功能是否正常
) else (
    echo 下载失败!
)

pause
