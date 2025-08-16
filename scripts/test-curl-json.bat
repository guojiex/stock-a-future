@echo off
echo 测试JSON POST请求
echo ===================
echo.

echo 1. 测试JSON数据POST请求
.\curl.exe -X POST -H "Content-Type: application/json" -d "{\"name\":\"test\",\"value\":123}" http://httpbin.org/post

echo.
echo 2. 测试表单数据POST请求
.\curl.exe -X POST -d "name=test&value=123" http://httpbin.org/post

echo.
echo 3. 测试基本认证
.\curl.exe -u "user:pass" http://httpbin.org/basic-auth/user/pass

echo.
echo 测试完成！
pause
