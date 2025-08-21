# 测试Stock-A-Future数据清理功能
Write-Host "测试Stock-A-Future数据清理功能" -ForegroundColor Green
Write-Host "===============================" -ForegroundColor Green

$SERVER_URL = "http://localhost:8080"

Write-Host "`n1. 检查清理服务状态..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$SERVER_URL/api/v1/cleanup/status" -Method Get
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "错误: $_" -ForegroundColor Red
}

Write-Host "`n2. 手动触发数据清理..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$SERVER_URL/api/v1/cleanup/manual" -Method Post
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "错误: $_" -ForegroundColor Red
}

Write-Host "`n3. 等待5秒后再次检查状态..." -ForegroundColor Yellow
Start-Sleep -Seconds 5
try {
    $response = Invoke-RestMethod -Uri "$SERVER_URL/api/v1/cleanup/status" -Method Get
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "错误: $_" -ForegroundColor Red
}

Write-Host "`n4. 更新清理配置（将股票信号保留天数改为60天）..." -ForegroundColor Yellow
$config = @{
    retention_days = 60
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$SERVER_URL/api/v1/cleanup/config" -Method Put -Body $config -ContentType "application/json"
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "错误: $_" -ForegroundColor Red
}

Write-Host "`n5. 验证配置更新..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$SERVER_URL/api/v1/cleanup/status" -Method Get
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "错误: $_" -ForegroundColor Red
}

Write-Host "`n测试完成！" -ForegroundColor Green
Write-Host "注意：收藏数据不会被清理，只清理过期的股票信号数据" -ForegroundColor Cyan
Read-Host "按回车键继续..."
