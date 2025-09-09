package models

import (
	"time"
)

// StrategyType 策略类型
type StrategyType string

const (
	StrategyTypeTechnical   StrategyType = "technical"   // 技术指标策略
	StrategyTypeFundamental StrategyType = "fundamental" // 基本面策略
	StrategyTypeML          StrategyType = "ml"          // 机器学习策略
	StrategyTypeComposite   StrategyType = "composite"   // 复合策略
)

// StrategyStatus 策略状态
type StrategyStatus string

const (
	StrategyStatusActive   StrategyStatus = "active"   // 活跃
	StrategyStatusInactive StrategyStatus = "inactive" // 非活跃
	StrategyStatusTesting  StrategyStatus = "testing"  // 测试中
)

// Strategy 策略模型
type Strategy struct {
	ID          string                 `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	Type        StrategyType           `json:"strategy_type" db:"strategy_type"`
	Status      StrategyStatus         `json:"status" db:"status"`
	Parameters  map[string]interface{} `json:"parameters" db:"parameters"`
	Code        string                 `json:"code,omitempty" db:"code"` // 策略代码（敏感信息，通常不返回给前端）
	CreatedBy   string                 `json:"created_by" db:"created_by"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// StrategyVersion 策略版本
type StrategyVersion struct {
	ID         string                 `json:"id" db:"id"`
	StrategyID string                 `json:"strategy_id" db:"strategy_id"`
	Version    string                 `json:"version" db:"version"`
	Code       string                 `json:"code" db:"code"`
	Parameters map[string]interface{} `json:"parameters" db:"parameters"`
	Changelog  string                 `json:"changelog" db:"changelog"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at"`
}

// StrategyPerformance 策略表现
type StrategyPerformance struct {
	ID              string    `json:"id" db:"id"`
	StrategyID      string    `json:"strategy_id" db:"strategy_id"`
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
	LastUpdated     time.Time `json:"last_updated" db:"last_updated"`
}

// Signal 交易信号
type Signal struct {
	ID         string     `json:"id" db:"id"`
	StrategyID string     `json:"strategy_id" db:"strategy_id"`
	Symbol     string     `json:"symbol" db:"symbol"`
	SignalType SignalType `json:"signal_type" db:"signal_type"`
	Side       TradeSide  `json:"side" db:"side"`             // 买入/卖出
	Strength   float64    `json:"strength" db:"strength"`     // 信号强度 (0-1)
	Confidence float64    `json:"confidence" db:"confidence"` // 置信度 (0-1)
	Price      float64    `json:"price" db:"price"`           // 触发价格
	Reason     string     `json:"reason" db:"reason"`         // 信号原因
	Timestamp  time.Time  `json:"timestamp" db:"timestamp"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// SignalType 信号类型
type SignalType string

const (
	SignalTypeBuy  SignalType = "buy"  // 买入信号
	SignalTypeSell SignalType = "sell" // 卖出信号
	SignalTypeHold SignalType = "hold" // 持有信号
	SignalTypeExit SignalType = "exit" // 退出信号
)

// TradeSide 交易方向
type TradeSide string

const (
	TradeSideBuy  TradeSide = "buy"  // 买入
	TradeSideSell TradeSide = "sell" // 卖出
)

// StrategyListRequest 策略列表请求
type StrategyListRequest struct {
	Page    int            `json:"page" form:"page"`
	Size    int            `json:"size" form:"size"`
	Status  StrategyStatus `json:"status" form:"status"`
	Type    StrategyType   `json:"type" form:"type"`
	Keyword string         `json:"keyword" form:"keyword"`
}

// StrategyListResponse 策略列表响应
type StrategyListResponse struct {
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
	Items []Strategy `json:"items"`
}

// CreateStrategyRequest 创建策略请求
type CreateStrategyRequest struct {
	Name        string                 `json:"name" validate:"required,max=100"`
	Description string                 `json:"description" validate:"max=1000"`
	Type        StrategyType           `json:"strategy_type" validate:"required"`
	Parameters  map[string]interface{} `json:"parameters"`
	Code        string                 `json:"code" validate:"required"`
}

// UpdateStrategyRequest 更新策略请求
type UpdateStrategyRequest struct {
	Name        *string                 `json:"name,omitempty" validate:"omitempty,max=100"`
	Description *string                 `json:"description,omitempty" validate:"omitempty,max=1000"`
	Status      *StrategyStatus         `json:"status,omitempty"`
	Parameters  *map[string]interface{} `json:"parameters,omitempty"`
	Code        *string                 `json:"code,omitempty"`
}

// 预定义策略参数结构

// MACDStrategyParams MACD策略参数
type MACDStrategyParams struct {
	FastPeriod    int     `json:"fast_period" validate:"min=1,max=50"`    // 快线周期，默认12
	SlowPeriod    int     `json:"slow_period" validate:"min=1,max=100"`   // 慢线周期，默认26
	SignalPeriod  int     `json:"signal_period" validate:"min=1,max=50"`  // 信号线周期，默认9
	BuyThreshold  float64 `json:"buy_threshold" validate:"min=-1,max=1"`  // 买入阈值，默认0
	SellThreshold float64 `json:"sell_threshold" validate:"min=-1,max=1"` // 卖出阈值，默认0
}

// MAStrategyParams 移动平均策略参数
type MAStrategyParams struct {
	ShortPeriod int     `json:"short_period" validate:"min=1,max=50"` // 短期均线周期，默认5
	LongPeriod  int     `json:"long_period" validate:"min=1,max=200"` // 长期均线周期，默认20
	MAType      string  `json:"ma_type" validate:"oneof=sma ema wma"` // 均线类型：sma/ema/wma
	Threshold   float64 `json:"threshold" validate:"min=0,max=0.1"`   // 突破阈值，默认0.01
}

// RSIStrategyParams RSI策略参数
type RSIStrategyParams struct {
	Period     int     `json:"period" validate:"min=1,max=50"`       // RSI周期，默认14
	Overbought float64 `json:"overbought" validate:"min=50,max=100"` // 超买阈值，默认70
	Oversold   float64 `json:"oversold" validate:"min=0,max=50"`     // 超卖阈值，默认30
}

// BollingerStrategyParams 布林带策略参数
type BollingerStrategyParams struct {
	Period int     `json:"period" validate:"min=1,max=50"`   // 周期，默认20
	StdDev float64 `json:"std_dev" validate:"min=0.5,max=5"` // 标准差倍数，默认2
}

// DefaultStrategies 默认策略配置
var DefaultStrategies = []Strategy{
	{
		ID:          "macd_strategy",
		Name:        "MACD金叉策略",
		Description: "基于MACD指标的金叉死叉交易策略，当MACD线上穿信号线时买入，下穿时卖出",
		Type:        StrategyTypeTechnical,
		Status:      StrategyStatusActive,
		Parameters: map[string]interface{}{
			"fast_period":    12,
			"slow_period":    26,
			"signal_period":  9,
			"buy_threshold":  0.0,
			"sell_threshold": 0.0,
		},
		CreatedBy: "system",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:          "ma_crossover",
		Name:        "双均线策略",
		Description: "短期均线突破长期均线的交易策略，短线上穿长线时买入，下穿时卖出",
		Type:        StrategyTypeTechnical,
		Status:      StrategyStatusActive,
		Parameters: map[string]interface{}{
			"short_period": 5,
			"long_period":  20,
			"ma_type":      "sma",
			"threshold":    0.01,
		},
		CreatedBy: "system",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:          "rsi_strategy",
		Name:        "RSI超买超卖策略",
		Description: "基于RSI指标的超买超卖交易策略，RSI低于超卖线时买入，高于超买线时卖出",
		Type:        StrategyTypeTechnical,
		Status:      StrategyStatusInactive,
		Parameters: map[string]interface{}{
			"period":     14,
			"overbought": 70.0,
			"oversold":   30.0,
		},
		CreatedBy: "system",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:          "bollinger_strategy",
		Name:        "布林带策略",
		Description: "基于布林带的均值回归策略，价格触及下轨时买入，触及上轨时卖出",
		Type:        StrategyTypeTechnical,
		Status:      StrategyStatusActive,
		Parameters: map[string]interface{}{
			"period":  20,
			"std_dev": 2.0,
		},
		CreatedBy: "system",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

// 注意: 移除了自定义的JSON序列化方法，让Go使用默认的JSON序列化
// 这样Parameters字段将正确地序列化为JSON对象而不是字符串
