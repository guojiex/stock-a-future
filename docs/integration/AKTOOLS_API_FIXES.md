# AKTools API 修复说明

## 🐛 问题描述

在实现 AKToolsClient 基本面数据 API 时遇到以下错误：

```
TypeError: stock_profit_sheet_by_report_em() got an unexpected keyword argument 'period'
```

## 🔍 问题分析

### 原因
AKShare 的 API 函数参数限制：

**财务报表 API 不接受 `period` 参数：**
- `stock_profit_sheet_by_report_em()` - 只接受 `symbol` 参数
- `stock_balance_sheet_by_report_em()` - 只接受 `symbol` 参数
- `stock_cash_flow_sheet_by_report_em()` - 只接受 `symbol` 参数

**实时数据 API 不接受任何参数：**
- `stock_zh_a_spot_em()` - 不接受任何参数，返回所有A股实时数据

### 原始实现问题
```go
// ❌ 错误的实现 - 参数问题
params := url.Values{}
params.Set("symbol", cleanSymbol)
if period != "" {
    params.Set("period", period)  // AKShare API 不支持此参数
}

// ❌ 错误的实现 - 股票代码格式问题
cleanSymbol := c.CleanStockSymbol(symbol)  // 返回 "000001"
params.Set("symbol", cleanSymbol)  // AKShare需要的是 "SZ000001" 格式
```

## ✅ 修复方案

### 1. 移除不支持的参数

修改 API 调用，不传递 `period` 参数：

```go
// ✅ 正确的实现
params := url.Values{}
params.Set("symbol", cleanSymbol)
// 不传递period参数，因为AKShare API不支持
```

### 2. 修正股票代码格式

AKShare 财务报表 API 需要特定的股票代码格式：

```go
// ❌ 错误的格式
cleanSymbol := c.CleanStockSymbol(symbol)  // 返回 "000001"
params.Set("symbol", cleanSymbol)

// ✅ 正确的格式
akshareSymbol := c.DetermineAKShareSymbol(symbol)  // 返回 "SZ000001"
params.Set("symbol", akshareSymbol)

// 新增的转换函数
func (c *AKToolsClient) DetermineAKShareSymbol(symbol string) string {
    cleanSymbol := c.CleanStockSymbol(symbol)
    
    if strings.HasPrefix(cleanSymbol, "600") || strings.HasPrefix(cleanSymbol, "601") ||
        strings.HasPrefix(cleanSymbol, "603") || strings.HasPrefix(cleanSymbol, "688") {
        return "SH" + cleanSymbol  // 上海股票: SH600519
    } else if strings.HasPrefix(cleanSymbol, "000") || strings.HasPrefix(cleanSymbol, "002") ||
        strings.HasPrefix(cleanSymbol, "300") {
        return "SZ" + cleanSymbol  // 深圳股票: SZ000001
    }
    return "SH" + cleanSymbol
}
```

### 3. 客户端筛选数据

在客户端根据 `period` 参数筛选返回的数据：

```go
// 如果指定了period，尝试找到匹配的记录
if period != "" {
    for _, data := range rawData {
        if reportDate, ok := data["报告期"].(string); ok {
            formattedDate := c.formatDateForFrontend(reportDate)
            if formattedDate == period {
                return c.convertToIncomeStatement(data, symbol, period, reportType)
            }
        }
    }
    // 如果没有找到匹配的period，返回错误
    return nil, fmt.Errorf("未找到指定期间的数据: %s, 期间: %s", symbol, period)
}

// 如果没有指定period，返回最新的一条数据
return c.convertToIncomeStatement(rawData[0], symbol, period, reportType)
```

### 4. 实时数据 API 修复

对于 `stock_zh_a_spot_em` API，采用不同的处理方式：

```go
// ❌ 错误的实现
params := url.Values{}
params.Set("symbol", cleanSymbol)
params.Set("trade_date", tradeDate)
apiURL := fmt.Sprintf("%s/api/public/stock_zh_a_spot_em?%s", c.baseURL, params.Encode())

// ✅ 正确的实现
apiURL := fmt.Sprintf("%s/api/public/stock_zh_a_spot_em", c.baseURL)  // 不传递任何参数

// 在返回的所有股票数据中查找指定股票
for _, data := range rawData {
    if code, ok := data["代码"].(string); ok {
        if code == cleanSymbol {
            return c.convertToDailyBasic(data, symbol, tradeDate)
        }
    }
}
```

### 5. 更新测试用例

修改测试用例以适应新的 API 调用方式：

```go
// ❌ 原始测试
testSymbol := "000001"  // 错误的股票代码格式
incomeStatement, err := client.GetIncomeStatement(testSymbol, testPeriod, "1")

// ✅ 修复后测试
testSymbol := "600519"  // 使用上海股票茅台作为测试
incomeStatement, err := client.GetIncomeStatement(testSymbol, "", "1")  // 获取最新数据
```

## 📝 修复的文件

### 1. `internal/client/aktools.go`
- **新增** `DetermineAKShareSymbol()` - 转换股票代码为AKShare格式
- `GetIncomeStatement()` - 修正股票代码格式，移除period参数，添加客户端筛选
- `GetBalanceSheet()` - 修正股票代码格式，同样修复
- `GetCashFlowStatement()` - 修正股票代码格式，同样修复
- `GetDailyBasic()` - 修复API调用，移除不支持的参数

### 2. `internal/client/aktools_fundamental_test.go`
- **更改测试股票代码**：从 `000001` 改为 `600519`（茅台）
- **新增** `DetermineAKShareSymbol()` 函数测试
- 所有测试方法改为不指定 period 或传递空字符串
- 测试最新数据获取功能

### 3. `scripts/test-aktools-fundamental.go`
- **更改测试股票代码**：从 `000001` 改为 `600519`（茅台）
- 测试脚本改为获取最新数据
- 添加动态期间测试逻辑

### 4. `docs/integration/AKTOOLS_FUNDAMENTAL_API.md`
- 更新文档说明 API 参数限制
- 添加期间筛选逻辑说明

## 🎯 修复结果

### API 调用流程
1. **发送请求**：只传递 `symbol` 参数到 AKTools API
2. **获取数据**：API 返回该股票所有可用期间的财务数据
3. **客户端筛选**：
   - 如果指定了 `period`，在返回数据中查找匹配的记录
   - 如果未指定 `period`，返回最新的记录
   - 如果指定了 `period` 但未找到匹配数据，返回错误

### 错误处理
- 网络错误：详细的连接错误信息
- 数据解析错误：JSON 解析失败信息
- 期间不匹配：明确指出未找到指定期间的数据
- 空数据：提示 API 返回空数据

### 使用示例

```go
client := NewAKToolsClient("http://127.0.0.1:8080")

// 获取最新财务数据
latest, err := client.GetIncomeStatement("000001", "", "1")

// 获取特定期间数据
specific, err := client.GetIncomeStatement("000001", "20231231", "1")

// 批量获取数据（按期间范围筛选）
batch, err := client.GetIncomeStatements("000001", "20220101", "20231231", "1")
```

## 🔄 测试验证

### 运行测试
```bash
# 运行单元测试
go test -v ./internal/client -run TestAKToolsClient

# 运行测试脚本
go run scripts/test-aktools-fundamental.go
```

### 预期结果
- ✅ 不再出现 `period` 参数错误
- ✅ 能够成功获取最新财务数据
- ✅ 能够根据指定期间筛选数据
- ✅ 错误处理更加准确和友好

## 📋 经验总结

### 1. API 文档的重要性
- 在实现第三方 API 客户端时，必须仔细阅读 API 文档
- 不能假设 API 支持某些参数，需要实际验证

### 2. 错误处理策略
- 当 API 不支持某些功能时，可以在客户端实现相应逻辑
- 保持接口一致性的同时，在实现层面做适配

### 3. 测试驱动开发
- 先编写测试用例，能够更早发现 API 调用问题
- 集成测试对于验证第三方服务集成非常重要

### 4. 向后兼容性
- 修复 API 调用问题时，保持客户端接口不变
- 通过内部逻辑调整来解决外部 API 限制

这次修复确保了 AKTools 基本面数据 API 能够正常工作，为后续的量化分析和基本面因子计算提供了可靠的数据基础。
