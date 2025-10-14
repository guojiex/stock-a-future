# 多策略回测展示功能

## 概述

回测页面现在完整支持多策略回测结果展示，可以同时查看多个策略的表现对比和组合整体效果。

## 功能特性

### 📊 1. 分层展示

回测结果按以下层次展示：

#### 第一层：组合整体表现
- 显示所有策略组合后的整体指标
- 包括总收益率、年化收益、最大回撤、夏普比率等
- 标题：**📊 组合整体表现**

#### 第二层：各策略表现对比
- 当选择多个策略（2个或以上）时显示
- 每个策略独立展示其性能指标
- 包含策略名称标题
- 标题：**🎯 各策略表现对比**

#### 第三层：单策略表现
- 当只选择1个策略时显示
- 直接展示该策略的性能指标
- 标题：**📊 策略表现**

### 📈 2. 组合权益曲线

- 展示所有策略组合后的整体权益曲线
- 包含组合价值、基准价值（如果有）、现金、持仓等
- 标题：**📈 组合权益曲线**

### 📋 3. 交易记录与策略关联

交易记录表格新增功能：
- **策略列**：当有多个策略时，显示策略列
- **策略标签**：每笔交易显示其所属策略的名称
- **颜色区分**：使用 Chip 组件标识不同策略
- 自动隐藏：只有1个策略时，策略列自动隐藏

## 数据结构

### 后端 API 响应

```json
{
  "success": true,
  "data": {
    "backtest_id": "uuid",
    
    // 组合整体指标（可选）
    "combined_metrics": {
      "total_return": 0.25,
      "annual_return": 0.30,
      "max_drawdown": -0.12,
      "sharpe_ratio": 2.1,
      "win_rate": 0.65,
      "total_trades": 150
    },
    
    // 各策略独立指标数组
    "performance": [
      {
        "strategy_id": "strategy-1-id",
        "strategy_name": "MACD金叉策略",
        "total_return": 0.20,
        "annual_return": 0.24,
        "max_drawdown": -0.10,
        "sharpe_ratio": 1.8,
        "win_rate": 0.60,
        "total_trades": 50
      },
      {
        "strategy_id": "strategy-2-id",
        "strategy_name": "双均线策略",
        "total_return": 0.18,
        "annual_return": 0.22,
        "max_drawdown": -0.08,
        "sharpe_ratio": 1.6,
        "win_rate": 0.62,
        "total_trades": 45
      },
      // ... 更多策略
    ],
    
    // 策略信息数组
    "strategies": [
      {
        "id": "strategy-1-id",
        "name": "MACD金叉策略",
        "description": "..."
      },
      {
        "id": "strategy-2-id",
        "name": "双均线策略",
        "description": "..."
      },
      // ... 更多策略
    ],
    
    // 组合权益曲线
    "equity_curve": [
      {
        "date": "2024-01-01",
        "portfolio_value": 1050000,
        "benchmark_value": 1020000,
        "cash": 200000,
        "holdings": 850000
      },
      // ... 更多数据点
    ],
    
    // 交易记录（包含策略ID）
    "trades": [
      {
        "id": "trade-1",
        "strategy_id": "strategy-1-id",  // 关联策略
        "symbol": "000001.SZ",
        "side": "buy",
        "quantity": 1000,
        "price": 15.50,
        "commission": 4.65,
        "timestamp": "2024-01-01T09:30:00Z"
      },
      // ... 更多交易
    ]
  }
}
```

## UI 展示逻辑

### 条件渲染规则

```typescript
// 1. 组合整体表现
{results.combined_metrics && (
  <Paper>组合整体表现</Paper>
)}

// 2. 多策略对比（2个或以上）
{results.performance?.length > 1 && (
  <Paper>
    {results.performance.map((perf, index) => (
      <Box key={index}>
        <Typography>{strategies[index].name}</Typography>
        {renderMetrics(perf)}
      </Box>
    ))}
  </Paper>
)}

// 3. 单策略表现（只有1个）
{results.performance?.length === 1 && (
  <Paper>
    {renderMetrics(results.performance[0])}
  </Paper>
)}

// 4. 组合权益曲线
{results.equity_curve?.length > 0 && (
  <Paper>
    <EquityCurveChart data={results.equity_curve} />
  </Paper>
)}

// 5. 交易记录（带策略信息）
{results.trades?.length > 0 && (
  <Paper>
    <TradesTable 
      trades={results.trades}
      strategies={results.strategies}
    />
  </Paper>
)}
```

### 交易表格策略列

```typescript
// 策略列只在多策略时显示
{strategies?.length > 1 && (
  <TableCell>
    <Chip label={strategyName} />
  </TableCell>
)}
```

## 示例场景

### 场景 1：选择 4 个策略

**显示内容**：
1. ✅ 组合整体表现卡片
2. ✅ 各策略表现对比卡片（4个策略分别展示）
3. ✅ 组合权益曲线图
4. ✅ 交易记录表（带策略列）

**布局**：
```
┌─────────────────────────────────┐
│ 📊 组合整体表现                 │
│ [6个性能指标卡片]               │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 🎯 各策略表现对比               │
│                                 │
│ MACD金叉策略                    │
│ [6个性能指标卡片]               │
│ ─────────────────────────────── │
│ 双均线策略                      │
│ [6个性能指标卡片]               │
│ ─────────────────────────────── │
│ 布林带策略                      │
│ [6个性能指标卡片]               │
│ ─────────────────────────────── │
│ RSI超买超卖策略                 │
│ [6个性能指标卡片]               │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 📈 组合权益曲线                 │
│ [统计信息 + 折线图]             │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 📋 交易记录                     │
│ 时间 | 代码 | 策略 | 方向 | ... │
└─────────────────────────────────┘
```

### 场景 2：选择 1 个策略

**显示内容**：
1. ✅ 策略表现卡片（单策略）
2. ✅ 组合权益曲线图
3. ✅ 交易记录表（无策略列）

### 场景 3：后端未返回 combined_metrics

- 自动跳过组合整体表现卡片
- 直接显示各策略表现对比

## 代码修改

### 文件：`web-react/src/pages/BacktestPage.tsx`

**修改内容**：
1. 将单个 Paper 拆分为多个独立的 Paper
2. 添加条件渲染逻辑
3. 支持 `combined_metrics` 和 `performance` 数组
4. 传递 `strategies` 给 TradesTable

### 文件：`web-react/src/components/TradesTable.tsx`

**修改内容**：
1. 添加 `strategies` 可选属性
2. 创建策略ID到名称的映射
3. 添加策略列（多策略时显示）
4. 使用 Chip 组件展示策略名称

## 性能优化

### 策略映射

使用 `useMemo` 缓存策略ID到名称的映射：

```typescript
const strategyMap = React.useMemo(() => {
  const map: Record<string, string> = {};
  if (strategies) {
    strategies.forEach(s => {
      map[s.id] = s.name;
    });
  }
  return map;
}, [strategies]);
```

### 条件渲染

- 只在有数据时渲染组件
- 避免空数据的不必要渲染
- 使用 `&&` 和 `?` 进行安全检查

## 用户体验

### 优势

1. **清晰的层次结构**：组合 → 各策略 → 详细数据
2. **灵活的展示**：自动适配单策略和多策略场景
3. **策略归属清晰**：每笔交易都能追溯到具体策略
4. **对比便捷**：多个策略性能一目了然

### 改进点

- ✅ 解决了"选择4个策略只显示1个"的问题
- ✅ 提供了完整的多策略对比视图
- ✅ 交易记录与策略关联清晰
- ✅ 支持组合整体和单策略表现

## 未来增强

### 可视化增强

1. **策略权益曲线对比**：
   - 在同一图表中展示多条策略曲线
   - 不同颜色区分不同策略
   - 支持显示/隐藏特定策略

2. **策略贡献分析**：
   - 饼图展示各策略对总收益的贡献
   - 柱状图对比各策略的夏普比率

3. **策略相关性分析**：
   - 热力图展示策略之间的相关性
   - 帮助优化策略组合

### 交互增强

1. **策略筛选**：
   - 交易记录支持按策略筛选
   - 权益曲线支持选择显示的策略

2. **策略权重调整**：
   - 模拟不同策略权重配置
   - 实时预览组合效果

## 测试建议

### 功能测试

1. **单策略场景**：
   - 选择1个策略
   - 验证显示"策略表现"而非"各策略表现对比"
   - 验证交易表格无策略列

2. **多策略场景**：
   - 选择2-5个策略
   - 验证显示所有策略的性能指标
   - 验证交易表格显示策略列
   - 验证策略名称正确映射

3. **边界条件**：
   - 无 combined_metrics
   - 无 strategies 数组
   - 策略ID与交易记录不匹配

### 数据验证

- 各策略指标总和应接近组合指标
- 交易记录的策略ID应存在于 strategies 数组中
- 策略数量应与 selectedStrategyIds 一致

## 相关文件

- `web-react/src/pages/BacktestPage.tsx` - 主页面逻辑
- `web-react/src/components/TradesTable.tsx` - 交易记录表格
- `web-react/src/components/EquityCurveChart.tsx` - 权益曲线图表

## 总结

✅ **完成的功能**:
- 组合整体表现展示
- 各策略独立表现对比
- 交易记录策略关联
- 自动适配单/多策略场景

🎉 **现在可以清晰地看到所有4个策略的表现了！**

