#!/bin/bash

# 测试API响应头是否正确设置
# 使用方法: ./test_headers.sh [服务器地址]

SERVER=${1:-"http://localhost:8080"}
ENDPOINTS=(
    "/api/v1/health"
    "/api/v1/stocks"
    "/api/v1/stocks/search?q=test"
    "/api/v1/favorites"
    "/api/v1/groups"
)

echo "测试API响应头..."
echo "服务器地址: $SERVER"
echo

for endpoint in "${ENDPOINTS[@]}"; do
    echo "测试端点: $endpoint"
    curl -s -I "$SERVER$endpoint" | grep -E "Content-Type|Access-Control-Allow"
    echo
done

echo "测试完成"
