# 测试信号计算服务的PowerShell脚本

$baseUrl = "http://localhost:8080"

Write-Host "=== 信号计算服务测试 ===" -ForegroundColor Green

# 1. 测试健康检查
Write-Host "`n1. 测试健康检查..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/health" -Method GET
    Write-Host "健康检查成功: $($response.status)" -ForegroundColor Green
} catch {
    Write-Host "健康检查失败: $_" -ForegroundColor Red
    exit 1
}

# 2. 测试计算单个股票信号
Write-Host "`n2. 测试计算单个股票信号..." -ForegroundColor Yellow
$signalRequest = @{
    ts_code = "000001.SZ"
    name = "平安银行"
    force = $true
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/signals/calculate" -Method POST -Body $signalRequest -ContentType "application/json"
    Write-Host "信号计算成功:" -ForegroundColor Green
    Write-Host "  股票代码: $($response.signal.ts_code)"
    Write-Host "  信号类型: $($response.signal.signal_type)"
    Write-Host "  信号强度: $($response.signal.signal_strength)"
    Write-Host "  置信度: $($response.signal.confidence)"
    Write-Host "  描述: $($response.signal.description)"
} catch {
    Write-Host "信号计算失败: $_" -ForegroundColor Red
}

# 3. 测试获取股票信号
Write-Host "`n3. 测试获取股票信号..." -ForegroundColor Yellow
try {
    $today = Get-Date -Format "yyyyMMdd"
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/signals/000001.SZ?signal_date=$today" -Method GET
    Write-Host "获取信号成功:" -ForegroundColor Green
    Write-Host "  信号日期: $($response.data.signal_date)"
    Write-Host "  交易日期: $($response.data.trade_date)"
    Write-Host "  信号类型: $($response.data.signal_type)"
} catch {
    Write-Host "获取信号失败: $_" -ForegroundColor Red
}

# 4. 测试获取最新信号列表
Write-Host "`n4. 测试获取最新信号列表..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/signals?limit=5" -Method GET
    Write-Host "获取最新信号成功:" -ForegroundColor Green
    Write-Host "  信号总数: $($response.data.total)"
    foreach ($signal in $response.data.signals) {
        Write-Host "    $($signal.ts_code) - $($signal.signal_type) - $($signal.signal_strength)"
    }
} catch {
    Write-Host "获取最新信号失败: $_" -ForegroundColor Red
}

# 5. 测试批量计算信号
Write-Host "`n5. 测试批量计算信号..." -ForegroundColor Yellow
$batchRequest = @{
    ts_codes = @("000002.SZ", "000858.SZ")
    force = $true
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/signals/batch" -Method POST -Body $batchRequest -ContentType "application/json"
    Write-Host "批量计算成功:" -ForegroundColor Green
    Write-Host "  总数: $($response.total)"
    Write-Host "  成功: $($response.success)"
    Write-Host "  失败: $($response.failed)"
    Write-Host "  耗时: $($response.duration)"
} catch {
    Write-Host "批量计算失败: $_" -ForegroundColor Red
}

Write-Host "`n=== 测试完成 ===" -ForegroundColor Green
