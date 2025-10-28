# 策略标签页权益曲线显示问题修复

## 🐛 问题描述

在回测结果页面，当选择多个策略进行回测时，切换不同的策略标签页，所有策略显示的都是**相同的权益曲线**。

### 问题根因

1. **后端数据结构**：
   - `BacktestResultsResponse.EquityCurve` 存储的是所有策略合并后的整体权益曲线
   - `BacktestResultsResponse.Performance` 数组中每个策略只有性能指标，没有独立的权益曲线数据

2. **前端渲染逻辑**：
   - 所有策略标签页都使用相同的 `results.equity_curve` 数据源
   - 没有区分每个策略的独立权益曲线

```typescript
// ❌ 问题代码
{results.equity_curve && results.equity_curve.length > 0 && (
  <EquityCurveChart
    data={results.equity_curve}  // 所有策略都用同一份数据
    initialCash={config.initial_cash}
  />
)}
```

## ✅ 解决方案

### 方案选择

有两种解决方案：

1. **前端方案**（已实现）：根据每个策略的交易记录，在前端重新计算该策略的权益曲线
2. **后端方案**（未来优化）：修改后端，为每个策略单独存储权益曲线

我们选择了**前端方案**作为快速修复，原因：
- 不需要修改后端API和数据库结构
- 可以立即解决问题
- 基于现有的交易记录数据，计算逻辑简单准确

### 实现细节

#### 1. 添加状态管理

```typescript
const [strategyEquityCurves, setStrategyEquityCurves] = useState<Record<number, any[]>>({});
```

#### 2. 实现权益曲线计算函数

```typescript
const calculateStrategyEquityCurve = useCallback((strategyId: string, trades: any[], initialCash: number) => {
  if (!trades || trades.length === 0) {
    // 如果没有交易，返回初始资金的平线
    return [{ date: config.start_date, equity: initialCash }];
  }

  // 过滤出该策略的交易
  const strategyTrades = trades.filter((t: any) => t.strategy_id === strategyId);
  
  if (strategyTrades.length === 0) {
    return [{ date: config.start_date, equity: initialCash }];
  }

  // 按日期排序
  const sortedTrades = [...strategyTrades].sort((a, b) => 
    new Date(a.exit_date || a.entry_date).getTime() - new Date(b.exit_date || b.entry_date).getTime()
  );

  // 计算权益曲线
  const equityCurve: any[] = [{ date: config.start_date, equity: initialCash }];
  let currentEquity = initialCash;

  sortedTrades.forEach((trade) => {
    if (trade.exit_date && trade.pnl !== undefined) {
      // 已平仓交易，更新权益
      currentEquity += trade.pnl;
      equityCurve.push({
        date: trade.exit_date,
        equity: currentEquity,
      });
    }
  });

  // 如果没有完成的交易，至少返回初始状态
  if (equityCurve.length === 1) {
    equityCurve.push({ date: config.end_date, equity: initialCash });
  }

  return equityCurve;
}, [config.start_date, config.end_date]);
```

**计算逻辑**：
1. 过滤出该策略的所有交易
2. 按平仓日期排序
3. 从初始资金开始，累加每笔交易的盈亏（PnL）
4. 生成时间序列的权益点位

#### 3. 自动计算所有策略的权益曲线

```typescript
useEffect(() => {
  if (results?.performance && Array.isArray(results.performance) && results.trades) {
    const curves: Record<number, any[]> = {};
    
    results.performance.forEach((_, index) => {
      const strategy = results.strategies?.[index];
      if (strategy) {
        curves[index] = calculateStrategyEquityCurve(
          strategy.id,
          results.trades,
          config.initial_cash
        );
      }
    });
    
    setStrategyEquityCurves(curves);
  }
}, [results, config.initial_cash, calculateStrategyEquityCurve]);
```

#### 4. 使用策略特定的权益曲线

```typescript
// ✅ 修复后的代码
{strategyEquityCurves[index] && strategyEquityCurves[index].length > 0 && (
  <Box sx={{ mb: 3 }}>
    <Typography variant="subtitle2" gutterBottom>
      📈 权益曲线 ({strategy?.name || `策略 ${index + 1}`})
    </Typography>
    <EquityCurveChart
      data={strategyEquityCurves[index]}  // 使用该策略独立的权益曲线
      initialCash={config.initial_cash}
    />
  </Box>
)}
```

## 📊 效果对比

### 修复前

```
Tab: 布林带策略     → 显示整体权益曲线
Tab: MACD策略      → 显示整体权益曲线（相同）
Tab: KDJ策略       → 显示整体权益曲线（相同）
```

### 修复后

```
Tab: 布林带策略     → 显示布林带策略的独立权益曲线
Tab: MACD策略      → 显示MACD策略的独立权益曲线
Tab: KDJ策略       → 显示KDJ策略的独立权益曲线
```

## 🔍 数据流程

```
1. 后端返回回测结果
   └─ performance: [策略1指标, 策略2指标, ...]
   └─ trades: [所有交易记录]
   └─ strategies: [策略1信息, 策略2信息, ...]
   └─ equity_curve: [整体权益曲线]

2. 前端接收结果
   └─ useEffect 监听 results 变化

3. 计算每个策略的权益曲线
   └─ 遍历 results.performance
   └─ 对每个策略：
      ├─ 过滤该策略的交易记录
      ├─ 按日期排序
      └─ 累计计算权益变化

4. 存储到 strategyEquityCurves
   └─ { 0: [曲线数据], 1: [曲线数据], ... }

5. 渲染策略标签页
   └─ Tab 0 → 使用 strategyEquityCurves[0]
   └─ Tab 1 → 使用 strategyEquityCurves[1]
   └─ ...
```

## 🎯 测试验证

### 测试步骤

1. **选择多个策略**：
   - 选择3个不同的策略（如MACD、布林带、KDJ）

2. **配置回测参数**：
   - 选择多只股票
   - 设置合理的时间范围（如最近1年）
   - 启动回测

3. **查看结果**：
   - 等待回测完成
   - 查看策略标签页

4. **切换标签验证**：
   - 依次点击每个策略标签
   - 观察权益曲线是否不同
   - 确认曲线与该策略的交易记录和盈亏一致

### 预期结果

- ✅ 每个策略的权益曲线独立显示
- ✅ 权益曲线的走势与该策略的盈亏表现一致
- ✅ 标题显示正确的策略名称
- ✅ 没有策略之间的数据混淆

## 🚀 未来优化

### 后端优化方案

如果后端需要提供更精确的权益曲线，可以考虑：

1. **修改数据模型**：
   ```go
   type StrategyPerformance struct {
       StrategyID   string        `json:"strategy_id"`
       Metrics      BacktestResult `json:"metrics"`
       EquityCurve  []EquityPoint  `json:"equity_curve"`  // 新增
       Trades       []Trade        `json:"trades"`
   }
   
   type BacktestResultsResponse struct {
       BacktestID         string                 `json:"backtest_id"`
       StrategyPerformance []StrategyPerformance `json:"strategy_performance"`
       OverallEquityCurve []EquityPoint          `json:"overall_equity_curve"`
       CombinedMetrics    *BacktestResult        `json:"combined_metrics"`
   }
   ```

2. **在回测执行时记录每个策略的权益变化**：
   - 为每个策略维护独立的权益序列
   - 实时记录每笔交易后的权益变化
   - 生成更精细的权益曲线（包含未平仓持仓的浮动盈亏）

3. **优点**：
   - 更精确的权益曲线（包含持仓浮盈浮亏）
   - 前端不需要计算，直接展示
   - 可以记录更多细节（如日内波动）

### 前端增强

1. **添加对比视图**：
   - 在同一图表中叠加显示多个策略的权益曲线
   - 便于直观比较策略表现

2. **交互式标注**：
   - 在权益曲线上标注买入/卖出点
   - 鼠标悬停显示交易详情

3. **性能优化**：
   - 大量交易记录时使用虚拟化
   - 权益曲线计算结果缓存

## 📝 相关文件

- `web-react/src/pages/BacktestPage.tsx` - 回测页面主文件
- `internal/models/backtest.go` - 回测数据模型
- `internal/service/backtest.go` - 回测服务逻辑

## 🔖 版本信息

- **修复日期**: 2025-10-28
- **问题发现**: 用户反馈策略标签切换时权益曲线相同
- **修复方式**: 前端计算每个策略的独立权益曲线
- **影响范围**: 仅前端展示，不影响后端逻辑和数据存储

