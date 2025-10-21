# 测试策略API并显示创建时间
Write-Host "正在测试策略API..." -ForegroundColor Cyan
Write-Host ""

Start-Sleep -Seconds 2  # 等待服务器启动

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8081/api/v1/strategies" -Method Get
    
    if ($response.success) {
        Write-Host "✅ API调用成功" -ForegroundColor Green
        Write-Host ""
        Write-Host "策略列表（共 $($response.data.items.Count) 个）:" -ForegroundColor Yellow
        Write-Host ""
        
        foreach ($strategy in $response.data.items) {
            Write-Host "[$($strategy.id)]" -ForegroundColor Cyan
            Write-Host "  名称: $($strategy.name)"
            Write-Host "  状态: $($strategy.status)"
            Write-Host "  创建时间: $($strategy.created_at)" -ForegroundColor $(if ($strategy.created_at -match "2024-01-01") { "Green" } else { "Red" })
            Write-Host ""
        }
        
        # 检查创建时间是否不同
        $uniqueTimes = $response.data.items | Select-Object -ExpandProperty created_at | Get-Unique
        Write-Host "创建时间唯一值数量: $($uniqueTimes.Count)" -ForegroundColor $(if ($uniqueTimes.Count -eq 1) { "Red" } else { "Green" })
        
        if ($uniqueTimes.Count -eq 1) {
            Write-Host "⚠️  问题：所有策略的创建时间都相同！" -ForegroundColor Red
            Write-Host "   时间值: $($uniqueTimes[0])" -ForegroundColor Red
        } else {
            Write-Host "✅ 策略创建时间各不相同，排序应该稳定" -ForegroundColor Green
        }
        
    } else {
        Write-Host "❌ API返回失败: $($response.message)" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ API调用失败: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "请确保服务器正在运行: .\bin\server.exe" -ForegroundColor Yellow
}

