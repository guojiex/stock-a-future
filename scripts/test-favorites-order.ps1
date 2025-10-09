# 测试收藏列表排序
# 验证收藏列表是否按创建时间从新到旧排序

Write-Host "=== 测试收藏列表排序 ===" -ForegroundColor Cyan

# 假设服务器运行在本地8080端口
$baseUrl = "http://127.0.0.1:8080"

# 获取收藏列表
Write-Host "`n1. 获取收藏列表..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/favorites" -Method GET
    
    if ($response.success) {
        $favorites = $response.data.favorites
        Write-Host "   ✅ 成功获取 $($favorites.Count) 个收藏" -ForegroundColor Green
        
        if ($favorites.Count -gt 1) {
            Write-Host "`n2. 检查排序（按创建时间）..." -ForegroundColor Yellow
            
            $isCorrectOrder = $true
            $previousTime = [DateTime]::MaxValue
            
            foreach ($fav in $favorites) {
                $createdAt = [DateTime]::Parse($fav.created_at)
                Write-Host "   - $($fav.name) ($($fav.ts_code)): $($fav.created_at)" -ForegroundColor Gray
                
                if ($createdAt -gt $previousTime) {
                    Write-Host "     ❌ 排序错误：此项比前一项更晚" -ForegroundColor Red
                    $isCorrectOrder = $false
                }
                
                $previousTime = $createdAt
            }
            
            if ($isCorrectOrder) {
                Write-Host "`n   ✅ 排序正确：按创建时间从新到旧" -ForegroundColor Green
            } else {
                Write-Host "`n   ❌ 排序错误：未按创建时间从新到旧" -ForegroundColor Red
            }
        } else {
            Write-Host "   ⚠️  收藏数量不足，无法测试排序" -ForegroundColor Yellow
        }
    } else {
        Write-Host "   ❌ 获取失败：$($response.error)" -ForegroundColor Red
    }
} catch {
    Write-Host "   ❌ 请求失败：$($_.Exception.Message)" -ForegroundColor Red
    Write-Host "   提示：请确保服务器正在运行（运行 .\bin\server.exe）" -ForegroundColor Yellow
}

Write-Host "`n=== 测试完成 ===" -ForegroundColor Cyan

