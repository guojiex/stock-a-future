# 策略权益曲线数据修复

## 问题描述

用户报告在使用后端生成的策略权益曲线数据后，权益曲线显示为没有更新数据的状态：
- 初始资金：¥1,000,000
- 最终价值：¥1,000,000.00
- 总收益率：+0.00%

这表明权益曲线只有起始和结束两个点，中间没有任何数据更新。

## 问题原因

经过分析代码，发现问题出在 `internal/service/backtest.go` 中的 `calculateEquityCurveFromTrades` 函数：

### 1. 只在卖出交易时记录权益点

```go
// 🐛 旧代码 - 问题版本
for _, trade := range sortedTrades {
    // 卖出交易才会产生盈亏
    if trade.Side == models.TradeSideSell && trade.PnL != 0 {
        currentEquity += trade.PnL
        equityCurve = append(equityCurve, models.EquityPoint{
            Date:           trade.Timestamp.Format("2006-01-02"),
            PortfolioValue: currentEquity,
            Cash:           trade.CashBalance,
            Holdings:       trade.HoldingAssets,
        })
    }
}
```

**问题**：
- 只在 `trade.Side == models.TradeSideSell` 时才添加权益点
- 如果策略只有买入交易而没有卖出，权益曲线就只有起点和终点
- 这导致图表显示为一条平线，收益率为0%

### 2. 使用PnL累加计算权益不准确

```go
// 🐛 错误的计算方式
currentEquity += trade.PnL
```

**问题**：
- `PnL`字段只在卖出交易时才有值
- 买入交易的`PnL`为0，无法反映资金的实际使用情况
- 应该使用`TotalAssets`字段，它包含了交易后的完整资产价值（现金+持仓）

## 解决方案

修改 `calculateEquityCurveFromTrades` 函数：

```go
// ✅ 修复后的代码
// 🔧 修复：遍历所有交易，使用TotalAssets字段（而不是累加PnL）
// TotalAssets = 现金 + 持仓市值，是每次交易后的完整资产价值
for _, trade := range sortedTrades {
    // 每次交易都记录权益点，无论买入还是卖出
    // 使用trade.TotalAssets字段，它包含了交易后的完整资产价值
    if trade.TotalAssets > 0 {
        equityCurve = append(equityCurve, models.EquityPoint{
            Date:           trade.Timestamp.Format("2006-01-02"),
            PortfolioValue: trade.TotalAssets,    // 使用TotalAssets而不是累加PnL
            Cash:           trade.CashBalance,
            Holdings:       trade.HoldingAssets,
        })
    }
}
```

### 修改要点

1. **记录所有交易的权益点**
   - 移除了 `trade.Side == models.TradeSideSell` 的条件
   - 买入和卖出交易都会记录权益变化

2. **使用TotalAssets字段**
   - 直接使用 `trade.TotalAssets`，它在每次交易执行时都会计算
   - `TotalAssets = Cash + HoldingAssets`（现金余额 + 持仓市值）
   - 准确反映了每次交易后的实际资产价值

3. **添加日志记录**
   - 记录策略ID、交易数量、权益点数量
   - 记录初始和最终资产价值
   - 便于调试和验证修复效果

## 验证要点

修复后，策略权益曲线应该：

1. ✅ **包含所有交易的权益点**
   - 买入交易后，现金减少，持仓增加，总资产基本不变（减去手续费）
   - 卖出交易后，现金增加，持仓减少，总资产根据盈亏变化

2. ✅ **准确反映资产变化**
   - 初始资金：`InitialCash`
   - 每次交易后：`TotalAssets = CashBalance + HoldingAssets`
   - 最终价值：最后一笔交易的`TotalAssets`

3. ✅ **支持多种交易模式**
   - 只买入未卖出：权益曲线显示买入后的持仓市值变化
   - 买入后卖出：权益曲线显示完整的买卖过程
   - 多次交易：每次交易都有对应的权益点

## 相关文件

- `internal/service/backtest.go` - 回测服务核心逻辑
  - `calculateEquityCurveFromTrades()` - 修复的函数
  - `generateStrategyPerformances()` - 调用上述函数生成策略性能数据

- `internal/models/backtest.go` - 数据模型
  - `EquityPoint` - 权益曲线点结构
  - `Trade` - 交易记录结构（包含TotalAssets字段）

- `web-react/src/pages/BacktestPage.tsx` - 前端回测页面
  - 第835-851行：使用后端提供的策略权益曲线数据

- `web-react/src/components/EquityCurveChart.tsx` - 权益曲线图表组件
  - 显示权益曲线和收益统计

## 测试建议

1. **单策略回测测试**
   - 测试只有买入交易的情况
   - 测试买入后卖出的情况
   - 测试多次买卖的情况

2. **多策略回测测试**
   - 测试每个策略都有独立的权益曲线
   - 验证不同策略的权益变化是否准确

3. **边界情况测试**
   - 测试没有交易的策略（应显示平线）
   - 测试只有一笔交易的策略
   - 测试交易频繁的策略

## 部署说明

1. 重新编译后端服务：
   ```bash
   go build -o bin/server.exe ./cmd/server
   ```

2. 重启服务器：
   ```bash
   ./bin/server.exe
   ```

3. 运行新的回测验证修复效果

## 相关文档

- [后端策略权益曲线实现指南](../features/BACKEND_STRATEGY_EQUITY_CURVES.md)
- [策略标签页权益曲线修复](./STRATEGY_TAB_EQUITY_CURVE_FIX.md)

---

**修复日期**: 2025-10-28
**修复人员**: AI Assistant
**问题影响**: 策略权益曲线显示不正确，影响回测结果分析
**修复版本**: 修复后版本使用`TotalAssets`字段记录所有交易的权益变化

