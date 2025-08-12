#!/bin/bash

# Stock-A-Future API测试脚本
# 使用前请确保服务器已启动: make dev

BASE_URL="http://localhost:8080"

echo "=== Stock-A-Future API 测试 ==="
echo "基础URL: $BASE_URL"
echo ""

# 检查jq是否安装
if ! command -v jq &> /dev/null; then
    echo "⚠️  建议安装jq以获得更好的JSON格式化显示: brew install jq"
    echo ""
fi

# 1. 健康检查
echo "1. 测试健康检查..."
if curl -s "$BASE_URL/api/v1/health" | jq '.' 2>/dev/null; then
    echo "✅ 健康检查通过"
else
    echo "❌ 健康检查失败"
fi
echo -e "\n"

# 2. 获取根路径信息
echo "2. 测试根路径..."
if curl -s "$BASE_URL/" | jq '.' 2>/dev/null; then
    echo "✅ 根路径访问正常"
else
    echo "❌ 根路径访问失败"
fi
echo -e "\n"

# 3. 测试CORS头（用于网页客户端）
echo "3. 测试CORS设置..."
CORS_HEADERS=$(curl -s -I "$BASE_URL/api/v1/health" | grep -i "access-control")
if [ -n "$CORS_HEADERS" ]; then
    echo "✅ CORS设置正常："
    echo "$CORS_HEADERS"
else
    echo "❌ CORS设置可能有问题"
fi
echo -e "\n"

# 4. 测试股票日线数据（需要有效的Tushare Token）
echo "4. 测试股票日线数据 (平安银行 000001.SZ)..."
echo "注意: 此测试需要有效的Tushare Token"
RESPONSE=$(curl -s "$BASE_URL/api/v1/stocks/000001.SZ/daily?start_date=20240101&end_date=20240105")
if echo "$RESPONSE" | jq -e '.success' >/dev/null 2>&1; then
    echo "✅ 股票数据获取成功"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
else
    echo "❌ 股票数据获取失败"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
fi
echo -e "\n"

# 5. 测试技术指标
echo "5. 测试技术指标..."
RESPONSE=$(curl -s "$BASE_URL/api/v1/stocks/000001.SZ/indicators")
if echo "$RESPONSE" | jq -e '.success' >/dev/null 2>&1; then
    echo "✅ 技术指标获取成功"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
else
    echo "❌ 技术指标获取失败"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
fi
echo -e "\n"

# 6. 测试买卖点预测
echo "6. 测试买卖点预测..."
RESPONSE=$(curl -s "$BASE_URL/api/v1/stocks/000001.SZ/predictions")
if echo "$RESPONSE" | jq -e '.success' >/dev/null 2>&1; then
    echo "✅ 买卖点预测获取成功"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
else
    echo "❌ 买卖点预测获取失败"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
fi
echo -e "\n"

# 7. 测试不同股票代码
echo "7. 测试其他股票代码..."
TEST_STOCKS=("600000.SH" "000002.SZ")
for stock in "${TEST_STOCKS[@]}"; do
    echo "测试 $stock..."
    RESPONSE=$(curl -s "$BASE_URL/api/v1/stocks/$stock/daily?start_date=20240101&end_date=20240105")
    if echo "$RESPONSE" | jq -e '.success' >/dev/null 2>&1; then
        echo "✅ $stock 数据获取成功"
    else
        echo "❌ $stock 数据获取失败"
    fi
done
echo -e "\n"

echo "=== 测试完成 ==="
echo ""
echo "🌐 网页示例使用说明："
echo "1. 确保API服务器正在运行"
echo "2. 在浏览器中打开 examples/index.html"
echo "3. 或使用本地服务器: cd examples && python3 -m http.server 8000"
echo ""
echo "如果看到错误，请检查:"
echo "1. 服务器是否已启动 (make dev)"
echo "2. .env文件中的TUSHARE_TOKEN是否正确"
echo "3. 网络连接是否正常"
echo "4. 防火墙设置是否阻止了请求"
