package models

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// JSONDecimal 用于JSON序列化的decimal类型，确保序列化为数字而不是字符串
type JSONDecimal struct {
	decimal.Decimal
}

// MarshalJSON 实现JSON序列化，将decimal转换为float64
func (d JSONDecimal) MarshalJSON() ([]byte, error) {
	f, _ := d.Float64()
	return json.Marshal(f)
}

// UnmarshalJSON 实现JSON反序列化
func (d *JSONDecimal) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		// 如果解析float64失败，尝试解析字符串
		var s string
		if err2 := json.Unmarshal(data, &s); err2 != nil {
			return err
		}
		dec, err3 := decimal.NewFromString(s)
		if err3 != nil {
			return err3
		}
		d.Decimal = dec
		return nil
	}
	d.Decimal = decimal.NewFromFloat(f)
	return nil
}

// NewJSONDecimal 创建新的JSONDecimal
func NewJSONDecimal(d decimal.Decimal) JSONDecimal {
	return JSONDecimal{d}
}

// StockBasic 股票基本信息
type StockBasic struct {
	TSCode   string `json:"ts_code"`   // 股票代码
	Symbol   string `json:"symbol"`    // 股票简称
	Name     string `json:"name"`      // 股票名称
	Area     string `json:"area"`      // 所在地域
	Industry string `json:"industry"`  // 所属行业
	Market   string `json:"market"`    // 市场类型
	ListDate string `json:"list_date"` // 上市日期
}

// StockDaily 股票日线数据
type StockDaily struct {
	TSCode    string      `json:"ts_code"`    // 股票代码
	TradeDate string      `json:"trade_date"` // 交易日期 YYYYMMDD
	Open      JSONDecimal `json:"open"`       // 开盘价
	High      JSONDecimal `json:"high"`       // 最高价
	Low       JSONDecimal `json:"low"`        // 最低价
	Close     JSONDecimal `json:"close"`      // 收盘价
	PreClose  JSONDecimal `json:"pre_close"`  // 昨收价
	Change    JSONDecimal `json:"change"`     // 涨跌额
	PctChg    JSONDecimal `json:"pct_chg"`    // 涨跌幅
	Vol       JSONDecimal `json:"vol"`        // 成交量(手)
	Amount    JSONDecimal `json:"amount"`     // 成交额(千元)
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
	DIF       JSONDecimal `json:"dif"`       // DIF线
	DEA       JSONDecimal `json:"dea"`       // DEA线(信号线)
	Histogram JSONDecimal `json:"histogram"` // MACD柱状图
	Signal    string      `json:"signal"`    // 金叉/死叉信号
}

// RSIIndicator RSI指标
type RSIIndicator struct {
	RSI14  JSONDecimal `json:"rsi14"`  // 14日RSI
	Signal string      `json:"signal"` // 超买/超卖信号
}

// BollingerBandsIndicator 布林带指标
type BollingerBandsIndicator struct {
	Upper  JSONDecimal `json:"upper"`  // 上轨
	Middle JSONDecimal `json:"middle"` // 中轨(MA20)
	Lower  JSONDecimal `json:"lower"`  // 下轨
	Signal string      `json:"signal"` // 突破信号
}

// MovingAverageIndicator 移动平均线指标
type MovingAverageIndicator struct {
	MA5   JSONDecimal `json:"ma5"`   // 5日均线
	MA10  JSONDecimal `json:"ma10"`  // 10日均线
	MA20  JSONDecimal `json:"ma20"`  // 20日均线
	MA60  JSONDecimal `json:"ma60"`  // 60日均线
	MA120 JSONDecimal `json:"ma120"` // 120日均线
}

// KDJIndicator KDJ指标
type KDJIndicator struct {
	K      JSONDecimal `json:"k"`      // K值
	D      JSONDecimal `json:"d"`      // D值
	J      JSONDecimal `json:"j"`      // J值
	Signal string      `json:"signal"` // 买卖信号
}

// PredictionResult 预测结果
type PredictionResult struct {
	TSCode      string                   `json:"ts_code"`
	TradeDate   string                   `json:"trade_date"`
	Predictions []TradingPointPrediction `json:"predictions"`
	Confidence  JSONDecimal              `json:"confidence"` // 预测置信度
	UpdatedAt   time.Time                `json:"updated_at"`
}

// TradingPointPrediction 买卖点预测
type TradingPointPrediction struct {
	Type        string      `json:"type"`        // "BUY" 或 "SELL"
	Price       JSONDecimal `json:"price"`       // 预测价格
	Date        string      `json:"date"`        // 预测日期
	Probability JSONDecimal `json:"probability"` // 概率
	Reason      string      `json:"reason"`      // 预测理由
	Indicators  []string    `json:"indicators"`  // 相关指标
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
