# AKTools 基本面数据 API 实现

## 📋 概述

本文档描述了 AKToolsClient 中基本面数据 API 的实现，包括财务报表数据、每日基本面指标和基本面因子的获取方法。

## 🎯 已实现功能

### 1. 财务报表数据

#### 利润表 (Income Statement)
- `GetIncomeStatement(symbol, period, reportType)` - 获取单个利润表
- `GetIncomeStatements(symbol, startPeriod, endPeriod, reportType)` - 批量获取利润表

**支持的数据字段：**
- 营业总收入、营业收入
- 营业总成本、管理费用、财务费用、研发费用
- 营业利润、利润总额、净利润、扣非净利润
- 基本每股收益、稀释每股收益

#### 资产负债表 (Balance Sheet)
- `GetBalanceSheet(symbol, period, reportType)` - 获取单个资产负债表
- `GetBalanceSheets(symbol, startPeriod, endPeriod, reportType)` - 批量获取资产负债表

**支持的数据字段：**
- 资产：资产总计、流动资产合计、货币资金、应收账款、存货、固定资产
- 负债：负债合计、流动负债合计、短期借款、应付账款、长期借款
- 权益：所有者权益合计、资本公积、未分配利润、实收资本

#### 现金流量表 (Cash Flow Statement)
- `GetCashFlowStatement(symbol, period, reportType)` - 获取单个现金流量表
- `GetCashFlowStatements(symbol, startPeriod, endPeriod, reportType)` - 批量获取现金流量表

**支持的数据字段：**
- 经营活动：经营活动现金流净额、销售收到现金、购买支付现金、税费支付
- 投资活动：投资活动现金流净额、收回投资现金、投资支付现金
- 筹资活动：筹资活动现金流净额、吸收投资现金、分配股利现金
- 现金净增加额、期初期末现金余额

### 2. 每日基本面指标

#### 每日基本面数据 (Daily Basic)
- `GetDailyBasic(symbol, tradeDate)` - 获取单个每日基本面指标
- `GetDailyBasics(symbol, startDate, endDate)` - 批量获取（当前返回空）
- `GetDailyBasicsByDate(tradeDate)` - 按日期获取所有股票（当前返回空）

**支持的数据字段：**
- 基本数据：收盘价、换手率、量比
- 估值指标：市盈率、市盈率TTM、市净率、市销率、市销率TTM
- 股本市值：总股本、流通股本、总市值、流通市值
- 分红指标：股息率、股息率TTM

### 3. 基本面因子

#### 基本面因子数据 (Fundamental Factor)
- `GetFundamentalFactor(symbol, tradeDate)` - 获取基本面因子（基础结构）
- `GetFundamentalFactors(symbol, startDate, endDate)` - 批量获取（当前返回空）
- `GetFundamentalFactorsByDate(tradeDate)` - 按日期获取所有股票（当前返回空）

**注意：** 基本面因子需要通过计算服务生成，当前返回基础数据结构。

## 🔧 API 端点映射

### AKTools HTTP API 端点

| 功能 | AKTools API 端点 | 参数 | 备注 |
|-----|-----------------|------|------|
| 利润表 | `/api/public/stock_profit_sheet_by_report_em` | symbol | 返回所有期间数据 |
| 资产负债表 | `/api/public/stock_balance_sheet_by_report_em` | symbol | 返回所有期间数据 |
| 现金流量表 | `/api/public/stock_cash_flow_sheet_by_report_em` | symbol | 返回所有期间数据 |
| 每日基本面 | `/api/public/stock_zh_a_spot_em` | 无参数 | 返回所有A股实时数据 |

## 📝 使用示例

### 基本用法

```go
// 创建AKTools客户端
client := NewAKToolsClient("http://127.0.0.1:8080")

// 获取利润表数据
incomeStatement, err := client.GetIncomeStatement("000001", "20231231", "1")
if err != nil {
    log.Printf("获取利润表失败: %v", err)
    return
}

fmt.Printf("营业收入: %s\n", incomeStatement.Revenue.String())
fmt.Printf("净利润: %s\n", incomeStatement.NetProfit.String())

// 获取资产负债表数据
balanceSheet, err := client.GetBalanceSheet("000001", "20231231", "1")
if err != nil {
    log.Printf("获取资产负债表失败: %v", err)
    return
}

fmt.Printf("资产总计: %s\n", balanceSheet.TotalAssets.String())
fmt.Printf("负债合计: %s\n", balanceSheet.TotalLiab.String())

// 获取现金流量表数据
cashFlowStatement, err := client.GetCashFlowStatement("000001", "20231231", "1")
if err != nil {
    log.Printf("获取现金流量表失败: %v", err)
    return
}

fmt.Printf("经营现金流: %s\n", cashFlowStatement.NetCashOperAct.String())

// 获取每日基本面指标
dailyBasic, err := client.GetDailyBasic("000001", "20240115")
if err != nil {
    log.Printf("获取每日基本面失败: %v", err)
    return
}

fmt.Printf("市盈率: %s\n", dailyBasic.Pe.String())
fmt.Printf("市净率: %s\n", dailyBasic.Pb.String())
```

### 批量数据获取

```go
// 批量获取利润表数据
incomeStatements, err := client.GetIncomeStatements("000001", "20220101", "20231231", "1")
if err != nil {
    log.Printf("批量获取利润表失败: %v", err)
    return
}

fmt.Printf("获取到 %d 条利润表数据\n", len(incomeStatements))

for _, stmt := range incomeStatements {
    fmt.Printf("期间: %s, 营业收入: %s, 净利润: %s\n", 
        stmt.FDate, stmt.Revenue.String(), stmt.NetProfit.String())
}
```

## 🧪 单元测试

### 运行测试

```bash
# 运行所有基本面数据测试
go test -v ./internal/client -run TestAKToolsClient

# 运行特定测试
go test -v ./internal/client -run TestAKToolsClient_GetIncomeStatement
go test -v ./internal/client -run TestAKToolsClient_GetBalanceSheet
go test -v ./internal/client -run TestAKToolsClient_GetCashFlowStatement
go test -v ./internal/client -run TestAKToolsClient_GetDailyBasic

# 运行不依赖外部服务的测试
go test -v ./internal/client -run TestAKToolsClient_SymbolHandling
go test -v ./internal/client -run TestAKToolsClient_DataConversion
```

### 测试覆盖的功能

1. **财务报表数据获取测试**
   - 利润表单个和批量获取
   - 资产负债表单个和批量获取
   - 现金流量表单个和批量获取

2. **每日基本面数据测试**
   - 每日基本面指标获取
   - 数据结构验证

3. **辅助功能测试**
   - 股票代码处理 (`CleanStockSymbol`, `DetermineTSCode`)
   - 数据类型转换 (`parseDecimalFromInterface`)
   - API连通性测试

4. **集成测试**
   - 多个API的综合调用测试
   - 数据一致性验证

## ⚠️ 注意事项

### 1. AKTools 服务依赖

所有测试都需要 AKTools 服务在 `http://127.0.0.1:8080` 运行。

**启动 AKTools 服务：**
```bash
# 使用 Docker 启动
docker run -p 8080:8080 akfamily/aktools

# 或者本地安装启动
pip install aktools
aktools --host 0.0.0.0 --port 8080
```

### 2. API 参数限制

**重要修正：** AKShare 的 API 参数限制和数据处理方式。

**财务报表 API 限制：**
- `stock_profit_sheet_by_report_em` - 仅接受 `symbol` 参数
- `stock_balance_sheet_by_report_em` - 仅接受 `symbol` 参数  
- `stock_cash_flow_sheet_by_report_em` - 仅接受 `symbol` 参数

**实时数据 API 限制：**
- `stock_zh_a_spot_em` - 不接受任何参数，返回所有A股实时数据

**数据处理逻辑：**
- **财务报表**：API 返回该股票的所有可用期间数据，客户端根据 `period` 参数筛选
- **每日基本面**：API 返回所有A股数据，客户端根据 `symbol` 参数筛选指定股票
- 如果不指定筛选条件，返回最新或第一条数据
- 如果指定了筛选条件但未找到匹配数据，返回错误

### 3. 数据可用性

- **财务报表数据**：依赖东方财富等数据源，可能存在延迟
- **每日基本面数据**：仅支持实时或最近交易日数据
- **历史数据**：部分API不支持历史数据批量获取

### 4. 错误处理

- 网络连接失败会返回详细错误信息
- API返回非200状态码会记录具体错误
- 数据解析失败会跳过错误数据继续处理
- 期间不匹配会返回具体的错误信息

### 5. 数据格式

- 所有数值字段使用 `decimal.Decimal` 类型确保精度
- 日期字段统一为 `YYYYMMDD` 格式
- 股票代码自动添加市场后缀（.SH/.SZ/.BJ）

## 🔄 下一步计划

1. **实现 TushareClient 基本面数据方法**
2. **创建基本面因子计算服务**
3. **设计数据库表结构存储基本面数据**
4. **添加 REST API 端点暴露基本面数据**
5. **集成到现有的股票分析系统**

## 📊 数据结构示例

### 利润表数据示例
```json
{
  "ts_code": "000001.SZ",
  "f_date": "20231231",
  "end_date": "20231231",
  "report_type": "1",
  "revenue": 123456789.00,
  "oper_revenue": 123456789.00,
  "net_profit": 12345678.90,
  "basic_eps": 1.23
}
```

### 资产负债表数据示例
```json
{
  "ts_code": "000001.SZ",
  "f_date": "20231231",
  "total_assets": 987654321.00,
  "total_liab": 654321098.00,
  "total_hldr_eqy": 333333223.00
}
```

### 每日基本面数据示例
```json
{
  "ts_code": "000001.SZ",
  "trade_date": "20240115",
  "close": 12.34,
  "pe": 15.67,
  "pb": 1.23,
  "total_mv": 123456789
}
```

这个实现为Stock-A-Future项目提供了完整的AKTools基本面数据获取能力，支持财务报表分析和基本面投资决策。
