package models

import (
	"math"
	"time"
)

// BacktestStatus 回测状态
type BacktestStatus string

const (
	BacktestStatusPending   BacktestStatus = "pending"   // 待执行
	BacktestStatusRunning   BacktestStatus = "running"   // 运行中
	BacktestStatusCompleted BacktestStatus = "completed" // 已完成
	BacktestStatusFailed    BacktestStatus = "failed"    // 失败
	BacktestStatusCancelled BacktestStatus = "cancelled" // 已取消
)

// Backtest 回测模型
type Backtest struct {
	ID           string         `json:"id" db:"id"`
	Name         string         `json:"name" db:"name"`
	StrategyID   string         `json:"strategy_id" db:"strategy_id"`
	StrategyName string         `json:"strategy_name,omitempty"` // 关联查询时填充
	Symbols      []string       `json:"symbols" db:"symbols"`
	StartDate    time.Time      `json:"start_date" db:"start_date"`
	EndDate      time.Time      `json:"end_date" db:"end_date"`
	InitialCash  float64        `json:"initial_cash" db:"initial_cash"`
	Commission   float64        `json:"commission" db:"commission"` // 手续费率
	Slippage     float64        `json:"slippage" db:"slippage"`     // 滑点
	Benchmark    string         `json:"benchmark" db:"benchmark"`   // 基准指数
	Status       BacktestStatus `json:"status" db:"status"`
	Progress     int            `json:"progress" db:"progress"` // 进度百分比
	ErrorMessage string         `json:"error_message,omitempty" db:"error_message"`
	CreatedBy    string         `json:"created_by" db:"created_by"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	StartedAt    *time.Time     `json:"started_at,omitempty" db:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty" db:"completed_at"`
}

// BacktestResult 回测结果
type BacktestResult struct {
	ID              string    `json:"id" db:"id"`
	BacktestID      string    `json:"backtest_id" db:"backtest_id"`
	TotalReturn     float64   `json:"total_return" db:"total_return"`         // 总收益率
	AnnualReturn    float64   `json:"annual_return" db:"annual_return"`       // 年化收益率
	MaxDrawdown     float64   `json:"max_drawdown" db:"max_drawdown"`         // 最大回撤
	SharpeRatio     float64   `json:"sharpe_ratio" db:"sharpe_ratio"`         // 夏普比率
	SortinoRatio    float64   `json:"sortino_ratio" db:"sortino_ratio"`       // 索提诺比率
	WinRate         float64   `json:"win_rate" db:"win_rate"`                 // 胜率
	ProfitFactor    float64   `json:"profit_factor" db:"profit_factor"`       // 盈亏比
	TotalTrades     int       `json:"total_trades" db:"total_trades"`         // 总交易次数
	AvgTradeReturn  float64   `json:"avg_trade_return" db:"avg_trade_return"` // 平均交易收益
	BenchmarkReturn float64   `json:"benchmark_return" db:"benchmark_return"` // 基准收益
	Alpha           float64   `json:"alpha" db:"alpha"`                       // Alpha
	Beta            float64   `json:"beta" db:"beta"`                         // Beta
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// Trade 交易记录
type Trade struct {
	ID         string    `json:"id" db:"id"`
	BacktestID string    `json:"backtest_id" db:"backtest_id"`
	Symbol     string    `json:"symbol" db:"symbol"`
	Side       TradeSide `json:"side" db:"side"`                         // 买入/卖出
	Quantity   int       `json:"quantity" db:"quantity"`                 // 数量
	Price      float64   `json:"price" db:"price"`                       // 价格
	Commission float64   `json:"commission" db:"commission"`             // 手续费
	PnL        float64   `json:"pnl,omitempty" db:"pnl"`                 // 盈亏（卖出时计算）
	SignalType string    `json:"signal_type,omitempty" db:"signal_type"` // 触发信号类型
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Position 持仓记录
type Position struct {
	ID           string    `json:"id" db:"id"`
	BacktestID   string    `json:"backtest_id" db:"backtest_id"`
	Symbol       string    `json:"symbol" db:"symbol"`
	Quantity     int       `json:"quantity" db:"quantity"`           // 持仓数量
	AvgPrice     float64   `json:"avg_price" db:"avg_price"`         // 平均成本
	MarketValue  float64   `json:"market_value" db:"market_value"`   // 市值
	UnrealizedPL float64   `json:"unrealized_pl" db:"unrealized_pl"` // 未实现盈亏
	Timestamp    time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// EquityPoint 权益曲线点
type EquityPoint struct {
	Date           string  `json:"date"`
	PortfolioValue float64 `json:"portfolio_value"`
	BenchmarkValue float64 `json:"benchmark_value,omitempty"`
	Cash           float64 `json:"cash"`
	Holdings       float64 `json:"holdings"`
}

// BacktestProgress 回测进度
type BacktestProgress struct {
	BacktestID          string     `json:"backtest_id"`
	Status              string     `json:"status"`
	Progress            int        `json:"progress"` // 0-100
	Message             string     `json:"message"`
	CurrentDate         string     `json:"current_date"`
	EstimatedCompletion *time.Time `json:"estimated_completion,omitempty"`
	Error               string     `json:"error,omitempty"`
}

// BacktestResultsResponse 回测结果响应
type BacktestResultsResponse struct {
	BacktestID     string         `json:"backtest_id"`
	Performance    BacktestResult `json:"performance"`
	EquityCurve    []EquityPoint  `json:"equity_curve"`
	Trades         []Trade        `json:"trades"`
	Positions      []Position     `json:"positions,omitempty"`
	Strategy       *Strategy      `json:"strategy"`        // 策略信息
	BacktestConfig BacktestConfig `json:"backtest_config"` // 回测配置
}

// BacktestConfig 回测配置信息
type BacktestConfig struct {
	Name        string   `json:"name"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	InitialCash float64  `json:"initial_cash"`
	Symbols     []string `json:"symbols"`
	Commission  float64  `json:"commission"`
	CreatedAt   string   `json:"created_at"`
}

// CreateBacktestRequest 创建回测请求
type CreateBacktestRequest struct {
	Name        string   `json:"name" validate:"required,max=100"`
	StrategyID  string   `json:"strategy_id" validate:"required"`
	Symbols     []string `json:"symbols" validate:"required,min=1"`
	StartDate   string   `json:"start_date" validate:"required"`
	EndDate     string   `json:"end_date" validate:"required"`
	InitialCash float64  `json:"initial_cash" validate:"min=10000"`
	Commission  float64  `json:"commission" validate:"min=0,max=0.01"` // 0-1%
	Slippage    float64  `json:"slippage" validate:"min=0,max=0.01"`   // 0-1%
	Benchmark   string   `json:"benchmark"`
}

// UpdateBacktestRequest 更新回测请求
type UpdateBacktestRequest struct {
	Name   *string         `json:"name,omitempty" validate:"omitempty,max=100"`
	Status *BacktestStatus `json:"status,omitempty"`
}

// BacktestListRequest 回测列表请求
type BacktestListRequest struct {
	Page       int            `json:"page" form:"page"`
	Size       int            `json:"size" form:"size"`
	Status     BacktestStatus `json:"status" form:"status"`
	StrategyID string         `json:"strategy_id" form:"strategy_id"`
	Keyword    string         `json:"keyword" form:"keyword"`
	StartDate  string         `json:"start_date" form:"start_date"`
	EndDate    string         `json:"end_date" form:"end_date"`
}

// BacktestListResponse 回测列表响应
type BacktestListResponse struct {
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
	Items []Backtest `json:"items"`
}

// BacktestEngine 回测引擎接口
type BacktestEngine interface {
	// StartBacktest 启动回测
	StartBacktest(backtest *Backtest, strategy *Strategy) error

	// GetProgress 获取回测进度
	GetProgress(backtestID string) (*BacktestProgress, error)

	// CancelBacktest 取消回测
	CancelBacktest(backtestID string) error

	// GetResults 获取回测结果
	GetResults(backtestID string) (*BacktestResultsResponse, error)
}

// Portfolio 投资组合
type Portfolio struct {
	Cash        float64             `json:"cash"`
	Positions   map[string]Position `json:"positions"`
	TotalValue  float64             `json:"total_value"`
	DailyReturn float64             `json:"daily_return"`
}

// OrderType 订单类型
type OrderType string

const (
	OrderTypeMarket OrderType = "market" // 市价单
	OrderTypeLimit  OrderType = "limit"  // 限价单
)

// Order 订单
type Order struct {
	ID         string     `json:"id"`
	Symbol     string     `json:"symbol"`
	Side       TradeSide  `json:"side"`
	Type       OrderType  `json:"type"`
	Quantity   int        `json:"quantity"`
	Price      float64    `json:"price"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	ExecutedAt *time.Time `json:"executed_at,omitempty"`
}

// MarketData 市场数据
type MarketData struct {
	Symbol   string    `json:"symbol"`
	Date     time.Time `json:"date"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   int64     `json:"volume"`
	Amount   float64   `json:"amount"`
	AdjClose float64   `json:"adj_close,omitempty"` // 复权收盘价
}

// PerformanceMetrics 性能指标计算器
type PerformanceMetrics struct {
	Returns          []float64 // 日收益率序列
	BenchmarkReturns []float64 // 基准收益率序列
	RiskFreeRate     float64   // 无风险利率
}

// CalculateMetrics 计算性能指标
func (pm *PerformanceMetrics) CalculateMetrics() *BacktestResult {
	if len(pm.Returns) == 0 {
		return &BacktestResult{}
	}

	result := &BacktestResult{
		TotalReturn:    pm.calculateTotalReturn(),
		AnnualReturn:   pm.calculateAnnualReturn(),
		MaxDrawdown:    pm.calculateMaxDrawdown(),
		SharpeRatio:    pm.calculateSharpeRatio(),
		SortinoRatio:   pm.calculateSortinoRatio(),
		WinRate:        pm.calculateWinRate(),
		TotalTrades:    len(pm.Returns),
		AvgTradeReturn: pm.calculateAvgReturn(),
	}

	if len(pm.BenchmarkReturns) > 0 {
		result.BenchmarkReturn = pm.calculateBenchmarkReturn()
		result.Alpha = pm.calculateAlpha()
		result.Beta = pm.calculateBeta()
	}

	return result
}

// calculateTotalReturn 计算总收益率
func (pm *PerformanceMetrics) calculateTotalReturn() float64 {
	totalReturn := 1.0
	for _, ret := range pm.Returns {
		totalReturn *= (1 + ret)
	}
	return totalReturn - 1
}

// calculateAnnualReturn 计算年化收益率
func (pm *PerformanceMetrics) calculateAnnualReturn() float64 {
	if len(pm.Returns) == 0 {
		return 0
	}

	totalReturn := pm.calculateTotalReturn()
	days := float64(len(pm.Returns))
	years := days / 252.0 // 假设一年252个交易日

	if years <= 0 {
		return 0
	}

	return math.Pow(1+totalReturn, 1/years) - 1
}

// calculateMaxDrawdown 计算最大回撤
func (pm *PerformanceMetrics) calculateMaxDrawdown() float64 {
	if len(pm.Returns) == 0 {
		return 0
	}

	peak := 1.0
	maxDrawdown := 0.0
	currentValue := 1.0

	for _, ret := range pm.Returns {
		currentValue *= (1 + ret)
		if currentValue > peak {
			peak = currentValue
		}

		drawdown := (peak - currentValue) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return -maxDrawdown // 返回负值表示回撤
}

// calculateSharpeRatio 计算夏普比率
func (pm *PerformanceMetrics) calculateSharpeRatio() float64 {
	if len(pm.Returns) == 0 {
		return 0
	}

	avgReturn := pm.calculateAvgReturn()
	stdDev := pm.calculateStdDev()

	if stdDev == 0 {
		return 0
	}

	return (avgReturn - pm.RiskFreeRate) / stdDev * math.Sqrt(252) // 年化
}

// calculateSortinoRatio 计算索提诺比率
func (pm *PerformanceMetrics) calculateSortinoRatio() float64 {
	if len(pm.Returns) == 0 {
		return 0
	}

	avgReturn := pm.calculateAvgReturn()
	downstdDev := pm.calculateDownsideStdDev()

	if downstdDev == 0 {
		return 0
	}

	return (avgReturn - pm.RiskFreeRate) / downstdDev * math.Sqrt(252)
}

// calculateWinRate 计算胜率
func (pm *PerformanceMetrics) calculateWinRate() float64 {
	if len(pm.Returns) == 0 {
		return 0
	}

	wins := 0
	for _, ret := range pm.Returns {
		if ret > 0 {
			wins++
		}
	}

	return float64(wins) / float64(len(pm.Returns))
}

// calculateAvgReturn 计算平均收益率
func (pm *PerformanceMetrics) calculateAvgReturn() float64 {
	if len(pm.Returns) == 0 {
		return 0
	}

	sum := 0.0
	for _, ret := range pm.Returns {
		sum += ret
	}

	return sum / float64(len(pm.Returns))
}

// calculateStdDev 计算标准差
func (pm *PerformanceMetrics) calculateStdDev() float64 {
	if len(pm.Returns) == 0 {
		return 0
	}

	avgReturn := pm.calculateAvgReturn()
	sumSquaredDiff := 0.0

	for _, ret := range pm.Returns {
		diff := ret - avgReturn
		sumSquaredDiff += diff * diff
	}

	variance := sumSquaredDiff / float64(len(pm.Returns))
	return math.Sqrt(variance)
}

// calculateDownsideStdDev 计算下行标准差
func (pm *PerformanceMetrics) calculateDownsideStdDev() float64 {
	if len(pm.Returns) == 0 {
		return 0
	}

	sumSquaredDownside := 0.0
	count := 0

	for _, ret := range pm.Returns {
		if ret < pm.RiskFreeRate {
			diff := ret - pm.RiskFreeRate
			sumSquaredDownside += diff * diff
			count++
		}
	}

	if count == 0 {
		return 0
	}

	variance := sumSquaredDownside / float64(count)
	return math.Sqrt(variance)
}

// calculateBenchmarkReturn 计算基准收益率
func (pm *PerformanceMetrics) calculateBenchmarkReturn() float64 {
	if len(pm.BenchmarkReturns) == 0 {
		return 0
	}

	totalReturn := 1.0
	for _, ret := range pm.BenchmarkReturns {
		totalReturn *= (1 + ret)
	}
	return totalReturn - 1
}

// calculateAlpha 计算Alpha
func (pm *PerformanceMetrics) calculateAlpha() float64 {
	if len(pm.Returns) == 0 || len(pm.BenchmarkReturns) == 0 {
		return 0
	}

	portfolioReturn := pm.calculateTotalReturn()
	benchmarkReturn := pm.calculateBenchmarkReturn()
	beta := pm.calculateBeta()

	return portfolioReturn - (pm.RiskFreeRate + beta*(benchmarkReturn-pm.RiskFreeRate))
}

// calculateBeta 计算Beta
func (pm *PerformanceMetrics) calculateBeta() float64 {
	if len(pm.Returns) == 0 || len(pm.BenchmarkReturns) == 0 {
		return 0
	}

	minLen := len(pm.Returns)
	if len(pm.BenchmarkReturns) < minLen {
		minLen = len(pm.BenchmarkReturns)
	}

	portfolioReturns := pm.Returns[:minLen]
	benchmarkReturns := pm.BenchmarkReturns[:minLen]

	// 计算协方差和方差
	portfolioAvg := 0.0
	benchmarkAvg := 0.0

	for i := 0; i < minLen; i++ {
		portfolioAvg += portfolioReturns[i]
		benchmarkAvg += benchmarkReturns[i]
	}

	portfolioAvg /= float64(minLen)
	benchmarkAvg /= float64(minLen)

	covariance := 0.0
	benchmarkVariance := 0.0

	for i := 0; i < minLen; i++ {
		portfolioDiff := portfolioReturns[i] - portfolioAvg
		benchmarkDiff := benchmarkReturns[i] - benchmarkAvg

		covariance += portfolioDiff * benchmarkDiff
		benchmarkVariance += benchmarkDiff * benchmarkDiff
	}

	if benchmarkVariance == 0 {
		return 0
	}

	return covariance / benchmarkVariance
}
