# 测试AKTools连接的PowerShell脚本

# 默认AKTools服务URL
$DEFAULT_URL = "http://127.0.0.1:8080"
$URL = if ($args[0]) { $args[0] } else { $DEFAULT_URL }

Write-Host "测试AKTools连接..." -ForegroundColor Cyan
Write-Host "服务URL: $URL"

# 测试健康检查接口
Write-Host "检查健康状态接口..." -ForegroundColor Cyan
try {
    $HEALTH_RESPONSE = Invoke-WebRequest -Uri "$URL/health" -Method GET -UseBasicParsing
    $HEALTH_STATUS = $HEALTH_RESPONSE.StatusCode
}
catch {
    $HEALTH_STATUS = $_.Exception.Response.StatusCode.value__
}

if ($HEALTH_STATUS -eq 200) {
    Write-Host "✓ 健康检查接口正常" -ForegroundColor Green
}
else {
    Write-Host "✗ 健康检查接口异常，HTTP状态码: $HEALTH_STATUS" -ForegroundColor Red
    Write-Host "请确保AKTools服务已启动并在 $URL 运行" -ForegroundColor Yellow
    exit 1
}

# 测试股票接口
Write-Host "测试股票接口..." -ForegroundColor Cyan
try {
    $STOCK_RESPONSE = Invoke-WebRequest -Uri "$URL/api/public/stock_zh_a_info?symbol=000001" -Method GET -UseBasicParsing
    $STOCK_STATUS = $STOCK_RESPONSE.StatusCode
}
catch {
    $STOCK_STATUS = $_.Exception.Response.StatusCode.value__
}

if ($STOCK_STATUS -eq 200) {
    Write-Host "✓ 股票接口正常" -ForegroundColor Green
}
else {
    Write-Host "✗ 股票接口异常，HTTP状态码: $STOCK_STATUS" -ForegroundColor Red
}

Write-Host ""
Write-Host "AKTools连接测试完成" -ForegroundColor Cyan
Write-Host "如果测试成功，您可以在config.env文件中设置:"
Write-Host "DATA_SOURCE_TYPE=aktools"
Write-Host "AKTOOLS_BASE_URL=$URL"
Write-Host ""
Write-Host "然后启动Stock-A-Future服务器"
