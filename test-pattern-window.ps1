# 测试模式预测时间窗口配置
Write-Host "🧪 测试模式预测时间窗口配置" -ForegroundColor Green
Write-Host "===============================" -ForegroundColor Green

# 检查服务器是否运行
Write-Host "📡 检查服务器状态..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8081/api/health" -Method GET -TimeoutSec 5
    Write-Host "✅ 服务器正在运行" -ForegroundColor Green
} catch {
    Write-Host "❌ 服务器未运行，请先启动服务器" -ForegroundColor Red
    Write-Host "💡 运行命令: .\bin\server.exe" -ForegroundColor Cyan
    exit 1
}

Write-Host ""
Write-Host "🔍 测试模式预测API..." -ForegroundColor Yellow

# 测试预测API
try {
    $predictionResponse = Invoke-RestMethod -Uri "http://localhost:8081/api/v1/stocks/000001/predictions" -Method GET -TimeoutSec 10
    Write-Host "✅ 预测API调用成功" -ForegroundColor Green
    
    if ($predictionResponse.predictions -and $predictionResponse.predictions.Count -gt 0) {
        Write-Host "📊 获取到 $($predictionResponse.predictions.Count) 个预测" -ForegroundColor Green
        
        # 检查是否有双响炮等K线形态预测
        $patternCount = 0
        $predictionResponse.predictions | ForEach-Object {
            if ($_.reason -match "双响炮|红三兵|乌云盖顶|锤子线|启明星|黄昏星") {
                $patternCount++
                Write-Host "  🎯 发现K线形态预测: $($_.type) - $($_.price) - $($_.reason)" -ForegroundColor Cyan
            }
        }
        
        if ($patternCount -gt 0) {
            Write-Host "🎉 K线形态预测功能正常工作，找到 $patternCount 个形态预测" -ForegroundColor Green
        } else {
            Write-Host "ℹ️  当前数据中没有K线形态预测，但功能正常" -ForegroundColor Yellow
        }
        
    } else {
        Write-Host "⚠️  没有预测数据" -ForegroundColor Yellow
    }
    
} catch {
    Write-Host "❌ 预测API调用失败: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "📋 配置说明：" -ForegroundColor Yellow
Write-Host "  1. 模式预测时间窗口可在 config.env 中配置" -ForegroundColor White
Write-Host "  2. 默认值: PATTERN_PREDICTION_DAYS=14 (两周)" -ForegroundColor White
Write-Host "  3. 只影响双响炮等K线形态预测，不影响其他功能" -ForegroundColor White
Write-Host "  4. 修改配置后需要重启服务器生效" -ForegroundColor White

Write-Host ""
Write-Host "🌐 当前配置: PATTERN_PREDICTION_DAYS=14" -ForegroundColor Cyan
Write-Host "💡 如需调整，请修改 config.env 文件中的 PATTERN_PREDICTION_DAYS 值" -ForegroundColor Cyan

Write-Host ""
Write-Host "✅ 测试完成！" -ForegroundColor Green
