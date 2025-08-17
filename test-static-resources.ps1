# 测试静态资源请求不再打印日志
Write-Host "🧪 测试静态资源请求日志过滤功能" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green

# 检查服务器是否运行
Write-Host "📡 检查服务器状态..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/health" -Method GET -TimeoutSec 5
    Write-Host "✅ 服务器正在运行" -ForegroundColor Green
} catch {
    Write-Host "❌ 服务器未运行，请先启动服务器" -ForegroundColor Red
    Write-Host "💡 运行命令: .\bin\server.exe" -ForegroundColor Cyan
    exit 1
}

Write-Host ""
Write-Host "🔍 测试静态资源请求..." -ForegroundColor Yellow

# 测试JavaScript文件请求
Write-Host "📜 测试JavaScript文件请求..." -ForegroundColor Cyan
try {
    $jsResponse = Invoke-WebRequest -Uri "http://localhost:8080/js/services/favorites.js" -Method GET -TimeoutSec 5
    Write-Host "  ✅ JS文件请求成功，状态码: $($jsResponse.StatusCode)" -ForegroundColor Green
} catch {
    Write-Host "  ❌ JS文件请求失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 测试CSS文件请求
Write-Host "🎨 测试CSS文件请求..." -ForegroundColor Cyan
try {
    $cssResponse = Invoke-WebRequest -Uri "http://localhost:8080/styles.css" -Method GET -TimeoutSec 5
    Write-Host "  ✅ CSS文件请求成功，状态码: $($cssResponse.StatusCode)" -ForegroundColor Green
} catch {
    Write-Host "  ❌ CSS文件请求失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 测试图标文件请求
Write-Host "🖼️  测试图标文件请求..." -ForegroundColor Cyan
try {
    $iconResponse = Invoke-WebRequest -Uri "http://localhost:8080/favicon.png" -Method GET -TimeoutSec 5
    Write-Host "  ✅ 图标文件请求成功，状态码: $($iconResponse.StatusCode)" -ForegroundColor Green
} catch {
    Write-Host "  ❌ 图标文件请求失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 测试API请求（应该仍然有日志）
Write-Host ""
Write-Host "🔌 测试API请求（应该仍然有日志）..." -ForegroundColor Yellow
try {
    $apiResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/health" -Method GET -TimeoutSec 5
    Write-Host "  ✅ API请求成功，状态码: 200" -ForegroundColor Green
} catch {
    Write-Host "  ❌ API请求失败: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "📋 测试说明：" -ForegroundColor Yellow
Write-Host "  1. 静态资源请求（JS、CSS、图片）不再打印日志" -ForegroundColor White
Write-Host "  2. API请求仍然会打印日志" -ForegroundColor White
Write-Host "  3. 请查看服务器控制台，确认静态资源请求没有日志输出" -ForegroundColor White

Write-Host ""
Write-Host "✅ 测试完成！" -ForegroundColor Green
Write-Host "💡 现在静态资源请求不会产生日志噪音了" -ForegroundColor Cyan
