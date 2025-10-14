# 回测图表实现文档

## 概述

为回测页面实现了完整的数据可视化功能，包括：
1. **权益曲线图表** - 展示回测期间的资产变化
2. **交易记录表格** - 详细列出所有交易信息

## 新增组件

### 1. EquityCurveChart（权益曲线图表）

**文件位置**: `web-react/src/components/EquityCurveChart.tsx`

**功能特性**：
- 📈 使用 Recharts 绘制专业的折线图
- 💰 显示组合价值、基准价值、现金、持仓等多条曲线
- 📊 实时统计初始资金、最终价值、总收益率
- 🎨 Material-UI 主题集成，深色模式友好
- 🖱️ 交互式 Tooltip，鼠标悬停显示详细数据
- 📱 响应式设计，自动适配不同屏幕尺寸

**数据格式**：
```typescript
interface EquityPoint {
  date: string;                // 日期
  portfolio_value: number;     // 组合价值（必需）
  benchmark_value?: number;    // 基准价值（可选）
  cash?: number;               // 现金（可选）
  holdings?: number;           // 持仓（可选）
}
```

**使用示例**：
```tsx
<EquityCurveChart
  data={results.equity_curve}
  initialCash={1000000}
/>
```

**视觉特性**：
- **组合价值线**: 蓝色实线，粗线条（主要曲线）
- **基准价值线**: 紫色虚线，用于对比
- **现金线**: 绿色细线，半透明
- **持仓线**: 橙色细线，半透明
- **Y轴**: 自动格式化为"万元"单位
- **图例**: 中文标签，易于理解

### 2. TradesTable（交易记录表格）

**文件位置**: `web-react/src/components/TradesTable.tsx`

**功能特性**：
- 📋 完整的交易记录列表
- 🔍 分页显示，支持每页 5/10/25/50 条
- 📊 统计信息：总交易次数、买卖比例、盈亏分布、累计盈亏
- 🎨 买入/卖出颜色区分（红色买入、绿色卖出）
- 💹 盈亏高亮显示（绿色盈利、红色亏损）
- 🏷️ 信号类型标签显示
- 💰 金额自动格式化，千分位分隔

**数据格式**：
```typescript
interface Trade {
  id: string;
  backtest_id: string;
  strategy_id?: string;
  symbol: string;              // 股票代码
  side: 'buy' | 'sell';       // 买入/卖出
  quantity: number;            // 数量
  price: number;               // 价格
  commission: number;          // 手续费
  pnl?: number;                // 盈亏（卖出时）
  signal_type?: string;        // 信号类型
  total_assets?: number;       // 总资产
  holding_assets?: number;     // 持仓资产
  cash_balance?: number;       // 现金余额
  timestamp: string;           // 交易时间
}
```

**使用示例**：
```tsx
<TradesTable trades={results.trades} />
```

**表格列**：
| 列名 | 说明 | 格式 |
|------|------|------|
| 时间 | 交易时间 | YYYY-MM-DD HH:MM |
| 股票代码 | 代码 | 文本 |
| 方向 | 买入/卖出 | Chip 标签 |
| 数量 | 股数 | 数字 |
| 价格 | 单价 | ¥X.XX |
| 金额 | 总金额 | ¥X,XXX.XX |
| 手续费 | 交易费用 | ¥X.XX |
| 盈亏 | 损益 | +/-¥X,XXX.XX |
| 总资产 | 交易后资产 | ¥X,XXX |
| 信号类型 | 策略信号 | 标签 |

## 集成到回测页面

### 修改的文件

**`web-react/src/pages/BacktestPage.tsx`**

```typescript
// 导入组件
import EquityCurveChart from '../components/EquityCurveChart';
import TradesTable from '../components/TradesTable';

// 在回测结果区域使用
{showResults && results && (
  <Paper sx={{ p: 3 }}>
    {/* 性能指标 */}
    {results.performance && renderMetrics(results.performance[0])}
    
    {/* 权益曲线 */}
    {results.equity_curve && results.equity_curve.length > 0 && (
      <EquityCurveChart
        data={results.equity_curve}
        initialCash={config.initial_cash}
      />
    )}
    
    {/* 交易记录 */}
    {results.trades && results.trades.length > 0 && (
      <TradesTable trades={results.trades} />
    )}
  </Paper>
)}
```

## 后端数据对接

### API 响应格式

后端 `GET /api/v1/backtests/{id}/results` 应返回：

```json
{
  "success": true,
  "data": {
    "backtest_id": "uuid",
    "performance": [
      {
        "total_return": 0.15,
        "annual_return": 0.18,
        "max_drawdown": -0.08,
        "sharpe_ratio": 1.5,
        "win_rate": 0.6,
        "total_trades": 50
      }
    ],
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
    "trades": [
      {
        "id": "uuid",
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

### 数据映射

| 后端字段 | 前端字段 | 说明 |
|---------|---------|------|
| `portfolio_value` | `portfolio_value` | 组合总价值 |
| `benchmark_value` | `benchmark_value` | 基准价值（可选） |
| `cash` | `cash` | 现金（可选） |
| `holdings` | `holdings` | 持仓价值（可选） |
| `pnl` | `pnl` | 交易盈亏 |
| `total_assets` | `total_assets` | 总资产 |

## 样式和主题

### 颜色方案

- **主色调**: Material-UI primary color（蓝色）
- **次要色**: Material-UI secondary color（紫色）
- **成功色**: 绿色（盈利、买入）
- **错误色**: 红色（亏损、卖出）
- **警告色**: 橙色（持仓）

### 响应式设计

- **桌面**: 完整显示所有列和图表
- **平板**: 自动调整图表高度和表格列宽
- **手机**: 保持可读性，支持横向滚动

## 性能优化

1. **图表渲染**:
   - 使用 Recharts 的虚拟化渲染
   - 数据点过多时自动抽样显示

2. **表格分页**:
   - 默认每页 10 条记录
   - 大数据集时避免一次渲染过多行

3. **条件渲染**:
   - 只在有数据时渲染组件
   - 避免空数据的不必要渲染

## 测试建议

### 功能测试

1. **权益曲线测试**:
   - 测试只有组合价值的情况
   - 测试包含基准的情况
   - 测试包含现金和持仓的情况
   - 测试空数据显示

2. **交易记录测试**:
   - 测试不同交易方向的显示
   - 测试盈亏计算和颜色
   - 测试分页功能
   - 测试排序功能

### 边界测试

- 空数据集
- 单条数据
- 超大数据集（1000+ 条）
- 极端数值（负数、零、超大数）

### UI/UX 测试

- 不同屏幕尺寸
- 深色/浅色模式切换
- 交互响应速度
- Tooltip 显示准确性

## 未来改进

### 功能增强

1. **图表交互**:
   - 缩放功能
   - 时间范围选择
   - 多策略对比

2. **数据导出**:
   - 导出为 CSV
   - 导出为 Excel
   - 导出图表为图片

3. **高级分析**:
   - 回撤分析图
   - 收益分布图
   - 持仓热力图

4. **性能优化**:
   - 虚拟滚动
   - 数据懒加载
   - 图表缓存

## 相关文件

- `web-react/src/components/EquityCurveChart.tsx` - 权益曲线组件
- `web-react/src/components/TradesTable.tsx` - 交易记录组件
- `web-react/src/pages/BacktestPage.tsx` - 回测页面
- `web-react/package.json` - 依赖配置（Recharts）

## 依赖项

```json
{
  "recharts": "^3.2.1",
  "@mui/material": "^7.3.2",
  "@mui/icons-material": "^7.3.2"
}
```

已包含在项目中，无需额外安装。

## 总结

✅ **完成的功能**:
- 权益曲线可视化
- 交易记录表格
- 统计信息展示
- 响应式设计
- Material-UI 集成

✅ **代码质量**:
- TypeScript 类型安全
- 组件化设计
- 性能优化
- 无 linter 错误

🎉 **回测页面现在具备完整的数据可视化能力！**

