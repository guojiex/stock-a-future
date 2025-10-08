@echo off
REM Windows批处理脚本 - 启动Go API服务器
echo 启动Go API服务器在8081端口...

REM 设置环境变量
set SERVER_PORT=8081

REM 启动Go服务器
go run cmd/server/main.go

pause
