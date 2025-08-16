# Go语言版本Curl工具编译脚本
# 适用于Windows PowerShell

Write-Host "正在编译Go语言版本的Curl工具..." -ForegroundColor Green
Write-Host ""

# 检查Go是否安装
try {
    $goVersion = go version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "Go未安装"
    }
    Write-Host "Go版本信息:" -ForegroundColor Yellow
    Write-Host $goVersion -ForegroundColor Cyan
    Write-Host ""
} catch {
    Write-Host "错误: 未找到Go语言环境，请先安装Go 1.22或更高版本" -ForegroundColor Red
    Write-Host "下载地址: https://golang.org/dl/" -ForegroundColor Yellow
    Read-Host "按回车键退出"
    exit 1
}

# 检查项目结构
if (-not (Test-Path "cmd\curl\main.go")) {
    Write-Host "错误: 未找到cmd\curl\main.go文件" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

# 编译curl工具
Write-Host "编译中..." -ForegroundColor Yellow
try {
    $buildResult = go build -o curl.exe .\cmd\curl\ 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "编译失败: $buildResult"
    }
} catch {
    Write-Host "编译失败！" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

Write-Host ""
Write-Host "编译成功！curl.exe已生成" -ForegroundColor Green
Write-Host ""

# 显示使用说明
Write-Host "使用方法示例:" -ForegroundColor Yellow
Write-Host "  .\curl.exe http://localhost:8080/api/stocks" -ForegroundColor Cyan
Write-Host "  .\curl.exe -X POST -d '{\"name\":\"AAPL\"}' http://localhost:8080/api/stocks" -ForegroundColor Cyan
Write-Host "  .\curl.exe -v http://localhost:8080/api/stocks" -ForegroundColor Cyan
Write-Host ""

# 检查是否要测试编译结果
$testCompile = Read-Host "是否要测试编译结果？(y/n)"
if ($testCompile -eq "y" -or $testCompile -eq "Y") {
    Write-Host ""
    Write-Host "测试编译结果..." -ForegroundColor Yellow
    
    # 显示帮助信息
    Write-Host "显示帮助信息:" -ForegroundColor Cyan
    & .\curl.exe | Out-String
    
    Write-Host ""
    Write-Host "测试完成！" -ForegroundColor Green
}

Write-Host ""
Write-Host "查看帮助: .\curl.exe" -ForegroundColor Yellow
Write-Host ""

Read-Host "按回车键退出"
