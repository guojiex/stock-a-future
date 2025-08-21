@echo off
REM 测试AKTools连接的Windows批处理脚本

REM 默认AKTools服务URL
set DEFAULT_URL=http://127.0.0.1:8080
if "%~1"=="" (
    set URL=%DEFAULT_URL%
) else (
    set URL=%~1
)

echo 测试AKTools连接...
echo 服务URL: %URL%

REM 测试健康检查接口
echo 检查健康状态接口...
curl -s -o nul -w "%%{http_code}" "%URL%/health" > temp_status.txt
set /p HEALTH_RESPONSE=<temp_status.txt
del temp_status.txt

if "%HEALTH_RESPONSE%"=="200" (
    echo ✓ 健康检查接口正常
) else (
    echo ✗ 健康检查接口异常，HTTP状态码: %HEALTH_RESPONSE%
    echo 请确保AKTools服务已启动并在 %URL% 运行
    exit /b 1
)

REM 测试股票接口
echo 测试股票接口...
curl -s -o nul -w "%%{http_code}" "%URL%/api/public/stock_zh_a_info?symbol=000001" > temp_status.txt
set /p STOCK_RESPONSE=<temp_status.txt
del temp_status.txt

if "%STOCK_RESPONSE%"=="200" (
    echo ✓ 股票接口正常
) else (
    echo ✗ 股票接口异常，HTTP状态码: %STOCK_RESPONSE%
)

echo.
echo AKTools连接测试完成
echo 如果测试成功，您可以在config.env文件中设置:
echo DATA_SOURCE_TYPE=aktools
echo AKTOOLS_BASE_URL=%URL%
echo.
echo 然后启动Stock-A-Future服务器
