#!/bin/bash

# Stock-A-Future API测试脚本
# 使用前请确保服务器已启动: make dev

# 从.env文件读取配置
ENV_FILE="../.env"
DEFAULT_PORT="8080"
DEFAULT_HOST="localhost"

# 读取.env文件中的配置
if [ -f "$ENV_FILE" ]; then
    echo "🔧 从.env文件读取配置..."
    # 使用source读取.env文件，但只提取需要的变量
    SERVER_PORT=$(grep "^SERVER_PORT=" "$ENV_FILE" 2>/dev/null | cut -d'=' -f2 | tr -d '"' | tr -d "'" | head -1)
    SERVER_HOST=$(grep "^SERVER_HOST=" "$ENV_FILE" 2>/dev/null | cut -d'=' -f2 | tr -d '"' | tr -d "'" | head -1)
else
    echo "⚠️  未找到.env文件，使用默认配置"
    echo "💡 提示: 可以创建.env文件来配置服务器设置："
    echo "   echo 'SERVER_PORT=8081' > .env"
    echo "   echo 'SERVER_HOST=localhost' >> .env"
    echo "   echo 'TUSHARE_TOKEN=your_token_here' >> .env"
    echo ""
fi

# 设置默认值
SERVER_PORT=${SERVER_PORT:-$DEFAULT_PORT}
SERVER_HOST=${SERVER_HOST:-$DEFAULT_HOST}

# 支持命令行参数覆盖配置
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--port)
            SERVER_PORT="$2"
            shift 2
            ;;
        -h|--host)
            SERVER_HOST="$2"
            shift 2
            ;;
        --help)
            echo "用法: $0 [选项]"
            echo "选项:"
            echo "  -p, --port PORT    指定服务器端口 (默认: $DEFAULT_PORT)"
            echo "  -h, --host HOST    指定服务器地址 (默认: $DEFAULT_HOST)"
            echo "  --help             显示此帮助信息"
            echo ""
            echo "示例:"
            echo "  $0                 # 使用.env文件或默认配置"
            echo "  $0 -p 8081         # 使用端口8081"
            echo "  $0 -h 192.168.1.100 -p 8080  # 使用自定义地址和端口"
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            echo "使用 --help 查看帮助信息"
            exit 1
            ;;
    esac
done

BASE_URL="http://${SERVER_HOST}:${SERVER_PORT}"

echo "=== Stock-A-Future API 测试 ==="
echo "配置信息:"
echo "  服务器地址: $SERVER_HOST"
echo "  服务器端口: $SERVER_PORT"
echo "  基础URL: $BASE_URL"
if [ -f "$ENV_FILE" ]; then
    echo "  配置来源: .env文件"
else
    echo "  配置来源: 默认值"
fi
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
echo "1. 确保API服务器正在运行在端口 $SERVER_PORT"
echo "2. 在浏览器中打开 examples/index.html"
echo "3. 或使用本地服务器: cd examples && python3 -m http.server 8000"
echo "4. 网页客户端会自动检测服务器端口"
echo ""
echo "如果看到错误，请检查:"
echo "1. 服务器是否已启动 (make dev 或 SERVER_PORT=$SERVER_PORT go run cmd/server/main.go)"
echo "2. .env文件中的TUSHARE_TOKEN是否正确"
echo "3. .env文件中的SERVER_PORT配置是否正确"
echo "4. 网络连接是否正常"
echo "5. 防火墙设置是否阻止了请求"
