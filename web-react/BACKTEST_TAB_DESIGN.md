# 回测结果 Tab 式设计

## 设计理念

参考 web 版的实现，采用 **Tab 切换** 的方式展示多策略回测结果，实现：
1. ✅ **紧凑的横向布局** - 性能指标横向排列，节省垂直空间
2. ✅ **策略独立展示** - 每个策略有自己的 Tab，包含性能、曲线、交易记录
3. ✅ **分策略数据** - 权益曲线和交易记录按策略分离展示

## 界面布局

### 整体结构

```
┌─────────────────────────────────────────────┐
│ 📊 回测结果                                  │
├─────────────────────────────────────────────┤
│ 组合整体表现（可选）                         │
│ [6个指标横向排列]                            │
├─────────────────────────────────────────────┤
│ ┌Tab1: MACD金叉┐ Tab2 │ Tab3 │ Tab4 │      │
│ └───────────────┘      │      │      │      │
├─────────────────────────────────────────────┤
│ 当前策略名称和描述                           │
│ [6个指标横向排列]                            │
│                                             │
│ 📈 权益曲线                                  │
│ [图表]                                       │
│                                             │
│ 📋 交易记录 (该策略的50笔)                   │
│ [表格 - 仅显示该策略的交易]                  │
└─────────────────────────────────────────────┘
```

## 关键特性

### 1. 紧凑的性能指标展示

**横向排列**，每个指标占用最小宽度：

```tsx
<Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
  <Box sx={{ minWidth: 140 }}>
    <Typography variant="caption">总收益率</Typography>
    <Typography variant="h6">+15.25%</Typography>
  </Box>
  <Box sx={{ minWidth: 140 }}>
    <Typography variant="caption">年化收益</Typography>
    <Typography variant="h6">+18.50%</Typography>
  </Box>
  // ... 更多指标
</Box>
```

**优势**：
- 节省垂直空间，一行显示所有指标
- 易于对比不同策略的同一指标
- 响应式设计，小屏幕自动换行

### 2. Tab 切换策略

使用 Material-UI `Tabs` 组件：

```tsx
<Tabs 
  value={selectedStrategyTab}
  onChange={(_, newValue) => setSelectedStrategyTab(newValue)}
  variant="scrollable"
>
  <Tab label="MACD金叉策略" icon={<Chip label="1" />} />
  <Tab label="双均线策略" icon={<Chip label="2" />} />
  <Tab label="布林带策略" icon={<Chip label="3" />} />
  <Tab label="RSI超买超卖策略" icon={<Chip label="4" />} />
</Tabs>
```

**优势**：
- 清晰的策略分离
- 一次只展示一个策略的详细信息
- 避免页面过长难以浏览
- 支持滚动，可容纳更多策略

### 3. 策略数据分离

#### 权益曲线
- 每个 Tab 显示该策略的权益曲线
- 如果后端提供分策略曲线，直接显示
- 如果只有组合曲线，所有 Tab 显示相同曲线

#### 交易记录
- 根据 `trade.strategy_id` 过滤交易记录
- 每个 Tab 只显示属于该策略的交易
- 交易数量显示在标题中

```tsx
const strategyTrades = results.trades?.filter(
  (t: any) => t.strategy_id === strategy?.id
) || [];

<Typography>📋 交易记录 ({strategyTrades.length} 笔)</Typography>
<TradesTable trades={strategyTrades} strategies={[strategy]} />
```

## 数据流

### 1. 组合整体表现（可选）

```json
{
  "combined_metrics": {
    "total_return": 0.25,
    "annual_return": 0.30,
    "max_drawdown": -0.12,
    "sharpe_ratio": 2.1,
    "win_rate": 0.65,
    "total_trades": 150
  }
}
```

在顶部显示，突出显示（蓝色背景）。

### 2. 各策略独立表现

```json
{
  "performance": [
    {
      "strategy_id": "strategy-1",
      "total_return": 0.20,
      ...
    },
    {
      "strategy_id": "strategy-2",
      "total_return": 0.18,
      ...
    }
  ],
  "strategies": [
    {
      "id": "strategy-1",
      "name": "MACD金叉策略",
      "description": "..."
    },
    {
      "id": "strategy-2",
      "name": "双均线策略",
      "description": "..."
    }
  ]
}
```

每个策略一个 Tab，切换 Tab 时显示对应策略的数据。

### 3. 交易记录过滤

```json
{
  "trades": [
    {
      "id": "trade-1",
      "strategy_id": "strategy-1",  // 关键字段
      "symbol": "000001.SZ",
      "side": "buy",
      ...
    },
    {
      "id": "trade-2",
      "strategy_id": "strategy-2",
      "symbol": "600000.SH",
      "side": "sell",
      ...
    }
  ]
}
```

根据 `strategy_id` 过滤，每个 Tab 只显示对应策略的交易。

## 对比：旧设计 vs 新设计

### 旧设计（纵向展示）

```
❌ 组合整体表现
   [6个大卡片]

❌ MACD金叉策略
   [6个大卡片]

❌ 双均线策略
   [6个大卡片]

❌ 布林带策略
   [6个大卡片]

❌ RSI超买超卖策略
   [6个大卡片]

❌ 组合权益曲线
   [图表]

❌ 所有交易记录（混在一起）
   [大表格]
```

**问题**：
- 页面极长，需要大量滚动
- 策略指标重复展示，效率低
- 交易记录混在一起，难以区分策略
- 视觉疲劳，信息密度过高

### 新设计（Tab 切换）

```
✅ 组合整体表现（可选）
   [6个指标横向]

✅ [Tab1] [Tab2] [Tab3] [Tab4]
   
   当前策略：MACD金叉策略
   [6个指标横向]
   
   📈 权益曲线
   [该策略的曲线]
   
   📋 交易记录 (50笔)
   [该策略的交易]
```

**优势**：
- 页面紧凑，一屏显示更多信息
- Tab 切换快速，策略独立清晰
- 交易记录自动过滤，精准定位
- 视觉舒适，信息组织合理

## 响应式设计

### 桌面端（宽度 > 900px）
- 指标横向排列，一行6个
- Tab 全部展示
- 图表和表格完整显示

### 平板端（600px - 900px）
- 指标横向排列，一行3个，自动换行
- Tab 滚动显示
- 图表自适应宽度

### 移动端（宽度 < 600px）
- 指标横向排列，一行2个
- Tab 滚动显示
- 表格横向滚动

## 实现细节

### 状态管理

```typescript
const [selectedStrategyTab, setSelectedStrategyTab] = useState(0);
```

### Tab 切换处理

```typescript
<Tabs 
  value={selectedStrategyTab}
  onChange={(_, newValue) => setSelectedStrategyTab(newValue)}
>
```

### Tab 内容渲染

```typescript
{results.performance.map((perfMetrics, index) => (
  <Box 
    key={index}
    role="tabpanel"
    hidden={selectedStrategyTab !== index}
  >
    {selectedStrategyTab === index && (
      <Box>
        {/* 只渲染当前激活的 Tab */}
      </Box>
    )}
  </Box>
))}
```

### 交易记录过滤

```typescript
const strategyTrades = results.trades?.filter(
  (t: any) => t.strategy_id === strategy?.id
) || [];
```

## 性能优化

### 1. 条件渲染
- 使用 `hidden` 属性隐藏非激活 Tab
- 只有当前 Tab 的内容被渲染

### 2. 数据过滤
- 交易记录在切换 Tab 时即时过滤
- 避免渲染所有交易记录

### 3. 懒加载
- 图表只在 Tab 激活时初始化
- 避免一次性加载多个图表

## 用户体验

### 优势

1. **信息密度优化**
   - 横向布局，一屏显示更多内容
   - 减少滚动次数

2. **导航清晰**
   - Tab 标签明确标识策略
   - 一键切换，快速对比

3. **数据精准**
   - 每个策略的数据独立展示
   - 交易记录自动过滤

4. **视觉舒适**
   - 避免长页面滚动疲劳
   - 信息分组清晰

### 改进点

- ✅ 解决了"页面太长"的问题
- ✅ 实现了策略独立展示
- ✅ 权益曲线和交易记录按策略分离
- ✅ 横向紧凑布局，节省空间

## 未来增强

### 1. 策略对比模式
- 添加"对比视图"按钮
- 并排显示多个策略的指标
- 方便直接对比

### 2. 自定义 Tab 顺序
- 拖拽排序
- 根据收益率自动排序

### 3. 图表叠加模式
- 在同一图表中显示多条策略曲线
- 不同颜色区分
- 可选择显示/隐藏特定策略

## 总结

新设计完美解决了你提出的问题：

1. ✅ **紧凑布局** - 性能指标横向排列，不再占用大量垂直空间
2. ✅ **Tab 切换** - 参考 web 版设计，每个策略独立展示
3. ✅ **数据分离** - 权益曲线和交易记录按策略过滤展示

现在 4 个策略可以清晰、紧凑地展示，页面不再冗长！🎉

