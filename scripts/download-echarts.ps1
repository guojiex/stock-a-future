# 下载ECharts到本地的PowerShell脚本
# 使用方法: .\scripts\download-echarts.ps1

Write-Host "开始下载ECharts到本地..." -ForegroundColor Green

# 创建目标目录
$targetDir = "web/static/js/lib/echarts"
if (!(Test-Path $targetDir)) {
    New-Item -ItemType Directory -Path $targetDir -Force
    Write-Host "创建目录: $targetDir" -ForegroundColor Yellow
}

# ECharts CDN URL（使用最新稳定版本）
$echartsUrl = "https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js"
$targetFile = "$targetDir/echarts.min.js"

Write-Host "下载ECharts从: $echartsUrl" -ForegroundColor Cyan
Write-Host "保存到: $targetFile" -ForegroundColor Cyan

try {
    # 下载文件
    Invoke-WebRequest -Uri $echartsUrl -OutFile $targetFile
    
    # 检查文件大小
    $fileSize = (Get-Item $targetFile).Length
    $fileSizeKB = [math]::Round($fileSize / 1KB, 2)
    
    Write-Host "下载完成!" -ForegroundColor Green
    Write-Host "文件大小: $fileSizeKB KB" -ForegroundColor Green
    Write-Host "文件路径: $targetFile" -ForegroundColor Green
    
    # 验证文件内容
    $content = Get-Content $targetFile -Raw
    if ($content -match "echarts") {
        Write-Host "文件验证成功: 包含echarts关键字" -ForegroundColor Green
    } else {
        Write-Host "警告: 文件内容可能不完整" -ForegroundColor Yellow
    }
    
} catch {
    Write-Host "下载失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host "`n下一步操作:" -ForegroundColor Yellow
Write-Host "1. 修改 index.html 中的script标签" -ForegroundColor White
Write-Host "2. 将 src='https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js'" -ForegroundColor White
Write-Host "   改为 src='js/lib/echarts/echarts.min.js'" -ForegroundColor White
Write-Host "3. 测试图表功能是否正常" -ForegroundColor White
