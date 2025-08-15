@echo off
title Stock-A-Future 快速启动
color 0A

echo.
echo 正在启动 Stock-A-Future API 服务器...
echo.

:: 直接启动服务器
go run cmd/server/main.go

pause
