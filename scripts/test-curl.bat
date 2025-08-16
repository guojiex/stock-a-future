@echo off
echo 测试Go语言版本的Curl工具
echo ================================
echo.

REM 检查curl.exe是否存在
if not exist "curl.exe" (
    echo 错误: 未找到curl.exe，请先运行build-curl.bat编译
    pause
    exit /b 1
)

echo 1. 测试帮助信息
echo --------------------------------
curl.exe
echo.

echo 2. 测试GET请求到httpbin.org
echo --------------------------------
curl.exe http://httpbin.org/get
echo.

echo 3. 测试POST请求到httpbin.org
echo --------------------------------
curl.exe -X POST -d "test=data" http://httpbin.org/post
echo.

echo 4. 测试详细输出模式
echo --------------------------------
curl.exe -v http://httpbin.org/json
echo.

echo 5. 测试JSON数据POST
echo --------------------------------
curl.exe -X POST -H "Content-Type: application/json" -d "{\"name\":\"test\",\"value\":123}" http://httpbin.org/post
echo.

echo 6. 测试基本认证
echo --------------------------------
curl.exe -u "user:pass" http://httpbin.org/basic-auth/user/pass
echo.

echo 7. 测试超时设置
echo --------------------------------
curl.exe --timeout 5 http://httpbin.org/delay/3
echo.

echo 测试完成！
echo.
echo 现在你可以使用curl.exe来测试你的本地API服务器了
echo 例如: curl.exe http://localhost:8080/api/stocks
echo.
pause
