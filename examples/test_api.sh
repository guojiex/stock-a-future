#!/bin/bash

# Stock-A-Future API测试脚本
# 使用前请确保服务器已启动: make dev

BASE_URL="http://localhost:8080"

echo "=== Stock-A-Future API 测试 ==="
echo "基础URL: $BASE_URL"
echo ""

# 1. 健康检查
echo "1. 测试健康检查..."
curl -s "$BASE_URL/api/v1/health" | jq '.' || echo "请安装jq: brew install jq"
echo -e "\n"

# 2. 获取根路径信息
echo "2. 测试根路径..."
curl -s "$BASE_URL/" | jq '.' || echo "JSON格式化失败"
echo -e "\n"

# 3. 测试股票日线数据（需要有效的Tushare Token）
echo "3. 测试股票日线数据 (平安银行 000001.SZ)..."
echo "注意: 此测试需要有效的Tushare Token"
curl -s "$BASE_URL/api/v1/stocks/000001.SZ/daily?start_date=20240101&end_date=20240105" | jq '.' || echo "请求失败 - 检查Tushare Token配置"
echo -e "\n"

# 4. 测试技术指标
echo "4. 测试技术指标..."
curl -s "$BASE_URL/api/v1/stocks/000001.SZ/indicators" | jq '.' || echo "请求失败"
echo -e "\n"

# 5. 测试买卖点预测
echo "5. 测试买卖点预测..."
curl -s "$BASE_URL/api/v1/stocks/000001.SZ/predictions" | jq '.' || echo "请求失败"
echo -e "\n"

echo "=== 测试完成 ==="
echo "如果看到错误，请检查:"
echo "1. 服务器是否已启动 (make dev)"
echo "2. .env文件中的TUSHARE_TOKEN是否正确"
echo "3. 网络连接是否正常"
