package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// StockDaily 股票日线数据
type StockDaily struct {
	TSCode    string          `json:"ts_code"`    // 股票代码
	TradeDate string          `json:"trade_date"` // 交易日期 YYYYMMDD
	Open      decimal.Decimal `json:"open"`       // 开盘价
	High      decimal.Decimal `json:"high"`       // 最高价
	Low       decimal.Decimal `json:"low"`        // 最低价
	Close     decimal.Decimal `json:"close"`      // 收盘价
	PreClose  decimal.Decimal `json:"pre_close"`  // 昨收价
	Change    decimal.Decimal `json:"change"`     // 涨跌额
	PctChg    decimal.Decimal `json:"pct_chg"`    // 涨跌幅
	Vol       decimal.Decimal `json:"vol"`        // 成交量(手)
	Amount    decimal.Decimal `json:"amount"`     // 成交额(千元)
}

// TechnicalIndicators 技术指标
type TechnicalIndicators struct {
	TSCode    string                   `json:"ts_code"`
	TradeDate string                   `json:"trade_date"`
	MACD      *MACDIndicator           `json:"macd,omitempty"`
	RSI       *RSIIndicator            `json:"rsi,omitempty"`
	BOLL      *BollingerBandsIndicator `json:"boll,omitempty"`
	MA        *MovingAverageIndicator  `json:"ma,omitempty"`
	KDJ       *KDJIndicator            `json:"kdj,omitempty"`
}

// MACDIndicator MACD指标
type MACDIndicator struct {
	DIF    decimal.Decimal `json:"dif"`    // DIF线
	DEA    decimal.Decimal `json:"dea"`    // DEA线(信号线)
	MACD   decimal.Decimal `json:"macd"`   // MACD柱状图
	Signal string          `json:"signal"` // 金叉/死叉信号
}

// RSIIndicator RSI指标
type RSIIndicator struct {
	RSI6   decimal.Decimal `json:"rsi6"`   // 6日RSI
	RSI12  decimal.Decimal `json:"rsi12"`  // 12日RSI
	RSI24  decimal.Decimal `json:"rsi24"`  // 24日RSI
	Signal string          `json:"signal"` // 超买/超卖信号
}

// BollingerBandsIndicator 布林带指标
type BollingerBandsIndicator struct {
	Upper  decimal.Decimal `json:"upper"`  // 上轨
	Middle decimal.Decimal `json:"middle"` // 中轨(MA20)
	Lower  decimal.Decimal `json:"lower"`  // 下轨
	Signal string          `json:"signal"` // 突破信号
}

// MovingAverageIndicator 移动平均线指标
type MovingAverageIndicator struct {
	MA5   decimal.Decimal `json:"ma5"`   // 5日均线
	MA10  decimal.Decimal `json:"ma10"`  // 10日均线
	MA20  decimal.Decimal `json:"ma20"`  // 20日均线
	MA60  decimal.Decimal `json:"ma60"`  // 60日均线
	MA120 decimal.Decimal `json:"ma120"` // 120日均线
}

// KDJIndicator KDJ指标
type KDJIndicator struct {
	K      decimal.Decimal `json:"k"`      // K值
	D      decimal.Decimal `json:"d"`      // D值
	J      decimal.Decimal `json:"j"`      // J值
	Signal string          `json:"signal"` // 买卖信号
}

// PredictionResult 预测结果
type PredictionResult struct {
	TSCode      string                   `json:"ts_code"`
	TradeDate   string                   `json:"trade_date"`
	Predictions []TradingPointPrediction `json:"predictions"`
	Confidence  decimal.Decimal          `json:"confidence"` // 预测置信度
	UpdatedAt   time.Time                `json:"updated_at"`
}

// TradingPointPrediction 买卖点预测
type TradingPointPrediction struct {
	Type        string          `json:"type"`        // "BUY" 或 "SELL"
	Price       decimal.Decimal `json:"price"`       // 预测价格
	Date        string          `json:"date"`        // 预测日期
	Probability decimal.Decimal `json:"probability"` // 概率
	Reason      string          `json:"reason"`      // 预测理由
	Indicators  []string        `json:"indicators"`  // 相关指标
}

// APIResponse 通用API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HealthStatus 健康检查状态
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}
