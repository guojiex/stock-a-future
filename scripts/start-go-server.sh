#!/bin/bash
# Unix/Linux/macOS脚本 - 启动Go API服务器
echo "启动Go API服务器在8081端口..."

# 设置环境变量并启动Go服务器
SERVER_PORT=8081 go run cmd/server/main.go
