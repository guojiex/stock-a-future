# 后端实现策略独立权益曲线

## 📌 概述

将策略权益曲线的计算从前端移到后端，确保数据的准确性和一致性。后端在回测执行时直接生成每个策略的独立权益曲线，并作为API响应的一部分返回给前端。

## 🎯 设计目标

1. **数据准确性**：后端在回测执行时已经处理所有交易数据，应该直接生成权益曲线
2. **前后端分离**：前端只负责展示，不应该重复计算业务逻辑
3. **性能优化**：避免前端重复计算，减少客户端负担
4. **向后兼容**：保持现有API字段不变，新增字段提供增强功能

## 🏗️ 架构设计

### 数据模型更新

#### 新增类型：`BacktestStrategyPerformance`

```go
// BacktestStrategyPerformance 单个策略的回测性能结果（包含独立权益曲线）
// 注意：与strategy.go中的StrategyPerformance区分，这里专门用于回测结果
type BacktestStrategyPerformance struct {
	StrategyID  string         `json:"strategy_id"`
	Metrics     BacktestResult `json:"metrics"`      // 策略性能指标
	EquityCurve []EquityPoint  `json:"equity_curve"` // 该策略的独立权益曲线
}
```

**命名说明**：
- 使用 `BacktestStrategyPerformance` 而不是 `StrategyPerformance`
- 避免与 `internal/models/strategy.go` 中的 `StrategyPerformance` 冲突
- 明确表示这是回测相关的性能数据

#### 更新：`BacktestResultsResponse`

```go
type BacktestResultsResponse struct {
	BacktestID           string                          `json:"backtest_id"`
	Performance          []BacktestResult                `json:"performance"`           // 多策略性能结果数组（向后兼容）
	StrategyPerformances []BacktestStrategyPerformance   `json:"strategy_performances"` // 🆕 每个策略的详细性能（含独立权益曲线）
	EquityCurve          []EquityPoint                   `json:"equity_curve"`          // 整体权益曲线
	Trades               []Trade                         `json:"trades"`
	Positions            []Position                      `json:"positions,omitempty"`
	Strategies           []*Strategy                     `json:"strategies"`
	BacktestConfig       BacktestConfig                  `json:"backtest_config"`
	CombinedMetrics      *BacktestResult                 `json:"combined_metrics,omitempty"`
}
```

**关键改进**：
- 新增 `StrategyPerformances` 字段，包含每个策略的完整性能数据和独立权益曲线
- 保留原有 `Performance` 字段用于向后兼容
- `EquityCurve` 继续保留为整体/组合权益曲线

### 服务层更新

#### BacktestService 新增字段

```go
type BacktestService struct {
	// ... 其他字段 ...
	
	// 🆕 每个策略的独立权益曲线
	// 结构: backtestID -> strategyID -> []EquityPoint
	backtestStrategyEquityCurves map[string]map[string][]models.EquityPoint
	
	// ... 其他字段 ...
}
```

**数据结构说明**：
```
backtestStrategyEquityCurves
└── "backtest_123"
    ├── "strategy_macd" → [{date, portfolio_value, ...}, ...]
    ├── "strategy_boll" → [{date, portfolio_value, ...}, ...]
    └── "strategy_kdj"  → [{date, portfolio_value, ...}, ...]
```

#### 新增方法：`generateStrategyPerformances`

```go
func (s *BacktestService) generateStrategyPerformances(
	backtestID string,
	backtest *models.Backtest,
	strategies []*models.Strategy,
	performanceResults []models.BacktestResult,
	allTrades []models.Trade,
) []models.BacktestStrategyPerformance
```

**功能**：
1. 检查是否已存储策略权益曲线（未来优化：回测执行时实时生成）
2. 如果没有存储的曲线，根据交易记录计算
3. 为每个策略组装完整的性能数据

**调用位置**：
- `GetBacktestResults` 方法中，在组装响应之前调用

#### 新增方法：`calculateEquityCurveFromTrades`

```go
func (s *BacktestService) calculateEquityCurveFromTrades(
	backtest *models.Backtest,
	strategyID string,
	allTrades []models.Trade,
) []models.EquityPoint
```

**功能**：
1. 过滤出指定策略的所有交易
2. 按时间戳排序交易
3. 从初始资金开始，累计计算每笔交易后的权益变化
4. 生成时间序列的权益点位

**计算逻辑**：
```go
// 1. 初始点：回测开始日期 + 初始资金
equityCurve = [{date: start_date, portfolio_value: initial_cash}]

// 2. 遍历该策略的所有交易，按时间排序
for trade in sorted_trades {
	if trade.side == "SELL" && trade.pnl != 0 {
		// 卖出交易产生盈亏
		current_equity += trade.pnl
		equityCurve.append({
			date: trade.timestamp,
			portfolio_value: current_equity,
			cash: trade.cash_balance,
			holdings: trade.holding_assets,
		})
	}
}

// 3. 终点：回测结束日期 + 最终权益
if len(equityCurve) == 1 {
	equityCurve.append({date: end_date, portfolio_value: current_equity})
}
```

## 📊 数据流程

### 回测执行流程（当前实现）

```
1. StartBacktest
   └── 执行回测
       └── 生成交易记录
           └── 存储到 backtestTrades[backtestID]

2. GetBacktestResults（查询结果）
   ├── 获取回测配置和状态
   ├── 获取性能指标 (performanceResults)
   ├── 获取所有交易记录 (trades)
   ├── 🆕 generateStrategyPerformances
   │   └── 为每个策略调用 calculateEquityCurveFromTrades
   │       ├── 过滤该策略的交易
   │       ├── 按时间排序
   │       ├── 累计计算权益变化
   │       └── 生成 EquityPoint 序列
   └── 组装 BacktestResultsResponse
       ├── performance (兼容)
       ├── strategy_performances (新增) ← 包含独立权益曲线
       ├── equity_curve (整体)
       └── trades, strategies, config...
```

### 未来优化方向

```
1. StartBacktest
   └── 执行回测
       └── 实时生成权益曲线
           ├── 每次交易后更新整体权益曲线
           └── 🆕 同时更新每个策略的独立权益曲线
               └── 存储到 backtestStrategyEquityCurves[backtestID][strategyID]

2. GetBacktestResults
   └── 直接使用存储的权益曲线
       ├── backtestEquityCurves[backtestID] (整体)
       └── backtestStrategyEquityCurves[backtestID][strategyID] (各策略)
```

## 🔧 前端集成

### API响应示例

```json
{
	"backtest_id": "bt_123456",
	"performance": [
		{"strategy_id": "macd", "total_return": 0.15, ...},
		{"strategy_id": "boll", "total_return": 0.08, ...}
	],
	"strategy_performances": [
		{
			"strategy_id": "macd",
			"metrics": {"total_return": 0.15, "sharpe_ratio": 1.5, ...},
			"equity_curve": [
				{"date": "2024-01-01", "portfolio_value": 1000000, ...},
				{"date": "2024-01-15", "portfolio_value": 1050000, ...},
				...
			]
		},
		{
			"strategy_id": "boll",
			"metrics": {"total_return": 0.08, "sharpe_ratio": 1.2, ...},
			"equity_curve": [
				{"date": "2024-01-01", "portfolio_value": 1000000, ...},
				{"date": "2024-01-20", "portfolio_value": 1020000, ...},
				...
			]
		}
	],
	"equity_curve": [...],  // 整体权益曲线
	"trades": [...],
	"strategies": [...]
}
```

### 前端使用

```typescript
// 优先使用后端提供的策略权益曲线
const strategyPerf = results.strategy_performances?.[index];
const equityCurveData = strategyPerf?.equity_curve || strategyEquityCurves[index];

// 前端的 calculateStrategyEquityCurve 作为后备方案
// 如果后端没有提供 strategy_performances，前端仍然可以计算
```

**优势**：
1. **向后兼容**：旧版本后端不影响前端正常工作
2. **渐进式升级**：前端自动使用后端提供的数据
3. **降级策略**：如果后端数据缺失，前端有 fallback

## 🧪 测试验证

### 单元测试

创建测试文件：`internal/service/backtest_equity_curve_test.go`

```go
func TestCalculateEquityCurveFromTrades(t *testing.T) {
	// 测试用例：
	// 1. 没有交易 → 平坦曲线
	// 2. 单笔盈利交易 → 上升曲线
	// 3. 单笔亏损交易 → 下降曲线
	// 4. 多笔交易 → 复杂曲线
	// 5. 过滤策略ID → 只包含指定策略的交易
}

func TestGenerateStrategyPerformances(t *testing.T) {
	// 测试多策略场景
	// 验证每个策略的权益曲线独立计算
}
```

### 集成测试

```bash
# 1. 启动回测服务
go run cmd/server/main.go

# 2. 创建多策略回测
curl -X POST http://localhost:8081/api/v1/backtests \
  -H "Content-Type: application/json" \
  -d '{
    "name": "多策略测试",
    "strategy_ids": ["macd_strategy", "boll_strategy"],
    "symbols": ["000001.SZ", "600976.SH"],
    "start_date": "2024-01-01",
    "end_date": "2024-10-01",
    "initial_cash": 1000000
  }'

# 3. 启动回测
curl -X POST http://localhost:8081/api/v1/backtests/{id}/start

# 4. 获取结果并验证
curl http://localhost:8081/api/v1/backtests/{id}/results | jq '.data.strategy_performances'

# 预期：
# - strategy_performances 数组长度为2
# - 每个元素包含 strategy_id, metrics, equity_curve
# - equity_curve 是非空数组
# - 两个策略的曲线不同
```

### 前端验证

1. **查看浏览器控制台**：
   ```javascript
   console.log(results.strategy_performances);
   // 应该看到每个策略的equity_curve数组
   ```

2. **检查网络请求**：
   - DevTools → Network → 找到 `/backtests/{id}/results`
   - 查看 Response 中是否包含 `strategy_performances`

3. **视觉验证**：
   - 切换不同的策略标签页
   - 观察权益曲线是否不同
   - 曲线走势应与该策略的盈亏一致

## 📈 性能考虑

### 当前实现（计算时）

**时间复杂度**：
- `calculateEquityCurveFromTrades`: O(N log N)，N为交易数量（排序主导）
- `generateStrategyPerformances`: O(S * N log N)，S为策略数量

**空间复杂度**：
- 每个策略的权益曲线：O(T)，T为该策略的交易笔数
- 总空间：O(S * T)

**优化建议**：
1. **回测执行时实时生成**：避免查询时计算，时间复杂度变为 O(1)
2. **持久化存储**：将权益曲线存入数据库，减少内存占用
3. **增量更新**：新交易产生时只更新曲线，不重新计算全部

### 未来优化（实时生成）

```go
// 在回测执行时，每次交易后更新权益曲线
func (s *BacktestService) recordTrade(trade *models.Trade) {
	// 更新整体权益曲线
	s.updateOverallEquityCurve(trade)
	
	// 🆕 更新策略独立权益曲线
	s.updateStrategyEquityCurve(trade.StrategyID, trade)
}

func (s *BacktestService) updateStrategyEquityCurve(strategyID string, trade *models.Trade) {
	curves := s.backtestStrategyEquityCurves[trade.BacktestID]
	if curves == nil {
		curves = make(map[string][]models.EquityPoint)
		s.backtestStrategyEquityCurves[trade.BacktestID] = curves
	}
	
	// 获取当前权益
	currentCurve := curves[strategyID]
	lastEquity := s.getLastEquity(currentCurve)
	
	// 卖出交易产生盈亏，更新权益
	if trade.Side == models.TradeSideSell {
		newEquity := lastEquity + trade.PnL
		curves[strategyID] = append(currentCurve, models.EquityPoint{
			Date:           trade.Timestamp.Format("2006-01-02"),
			PortfolioValue: newEquity,
			Cash:           trade.CashBalance,
			Holdings:       trade.HoldingAssets,
		})
	}
}
```

## 🔍 故障排查

### 问题：strategy_performances 为空

**可能原因**：
1. 后端版本未更新，不包含新逻辑
2. 回测还未完成
3. strategies 数组为空

**排查步骤**：
```bash
# 检查后端是否包含新字段
curl http://localhost:8081/api/v1/backtests/{id}/results | jq '.data | keys'
# 应该看到 "strategy_performances"

# 检查策略信息
curl http://localhost:8081/api/v1/backtests/{id}/results | jq '.data.strategies'
# 应该非空
```

### 问题：equity_curve 全为同一值

**可能原因**：
1. 所有交易都是买入（未平仓）
2. 交易的 PnL 字段为0或未设置
3. 交易记录中 StrategyID 不匹配

**排查步骤**：
```bash
# 检查交易记录
curl http://localhost:8081/api/v1/backtests/{id}/results | jq '.data.trades[] | {strategy_id, side, pnl}'

# 预期：
# - 应该有 "SELL" 交易
# - pnl 不为0
# - strategy_id 与策略列表匹配
```

### 问题：不同策略的曲线相同

**可能原因**：
1. 前端仍在使用老的计算逻辑
2. 后端返回的 strategy_performances 有误
3. 交易记录中所有交易的 strategy_id 相同

**排查步骤**：
```javascript
// 前端控制台
console.log(results.strategy_performances?.[0]?.equity_curve);
console.log(results.strategy_performances?.[1]?.equity_curve);
// 比较两个曲线是否不同

// 如果相同，检查后端交易记录
console.log(results.trades.filter(t => t.strategy_id === 'macd').length);
console.log(results.trades.filter(t => t.strategy_id === 'boll').length);
```

## 📝 相关文件

### 后端
- `internal/models/backtest.go` - 数据模型定义
- `internal/service/backtest.go` - 回测服务逻辑
- `internal/handler/backtest.go` - API处理器

### 前端
- `web-react/src/pages/BacktestPage.tsx` - 回测页面
- `web-react/src/services/api.ts` - API服务

### 文档
- `docs/fixes/STRATEGY_TAB_EQUITY_CURVE_FIX.md` - 前端修复文档
- `docs/features/BACKEND_STRATEGY_EQUITY_CURVES.md` - 本文档

## 🚀 部署清单

### 后端更新
- [x] 数据模型：新增 `BacktestStrategyPerformance`
- [x] 数据模型：更新 `BacktestResultsResponse`
- [x] 服务层：新增字段 `backtestStrategyEquityCurves`
- [x] 服务层：新增方法 `generateStrategyPerformances`
- [x] 服务层：新增方法 `calculateEquityCurveFromTrades`
- [x] 服务层：更新 `GetBacktestResults` 调用新方法
- [ ] 单元测试：`backtest_equity_curve_test.go`
- [ ] 集成测试：多策略回测场景

### 前端更新
- [x] 回测页面：优先使用 `strategy_performances[].equity_curve`
- [x] 回测页面：保留前端计算作为 fallback
- [ ] TypeScript类型：更新API响应类型定义

### 文档更新
- [x] 功能文档：后端实现说明
- [x] 故障排查指南
- [ ] API文档：更新响应格式说明

## 📊 效果对比

### 修复前（前端计算）

**问题**：
- 前端重复计算业务逻辑
- 与后端计算可能存在差异
- 增加客户端负担
- 无法利用后端的优化

### 修复后（后端提供）

**优势**：
- ✅ 数据源统一，保证一致性
- ✅ 计算逻辑集中在后端，易于维护
- ✅ 前端专注于展示，职责清晰
- ✅ 向后兼容，渐进式升级
- ✅ 为未来实时生成奠定基础

## 🎯 未来规划

### 第一阶段（已完成）
- ✅ 后端在查询时根据交易记录计算权益曲线
- ✅ 前端优先使用后端数据，保留 fallback

### 第二阶段（计划中）
- [ ] 后端在回测执行时实时生成权益曲线
- [ ] 将权益曲线持久化到数据库
- [ ] 优化大数据量时的性能

### 第三阶段（未来）
- [ ] 支持自定义权益曲线粒度（日/周/月）
- [ ] 支持权益曲线的统计分析（波动率、回撤区间等）
- [ ] 支持多策略权益曲线对比视图

---

**版本信息**：
- 创建日期：2025-10-28
- 最后更新：2025-10-28
- 作者：Stock-A-Future Team

