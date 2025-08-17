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
	SignalDate  string      `json:"signal_date"` // 信号产生的日期（基于哪一天的数据）
}

// APIResponse 通用API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// FavoriteStock 收藏股票结构
type FavoriteStock struct {
	ID        string    `json:"id"`         // 唯一标识
	TSCode    string    `json:"ts_code"`    // 股票代码
	Name      string    `json:"name"`       // 股票名称
	StartDate string    `json:"start_date"` // 收藏时的开始日期
	EndDate   string    `json:"end_date"`   // 收藏时的结束日期
	GroupID   string    `json:"group_id"`   // 所属分组ID
	SortOrder int       `json:"sort_order"` // 排序顺序
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// FavoriteGroup 收藏分组结构
type FavoriteGroup struct {
	ID        string    `json:"id"`         // 唯一标识
	Name      string    `json:"name"`       // 分组名称
	Color     string    `json:"color"`      // 分组颜色
	SortOrder int       `json:"sort_order"` // 排序顺序
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// FavoriteStockRequest 添加收藏股票请求
type FavoriteStockRequest struct {
	TSCode    string `json:"ts_code"`    // 股票代码
	Name      string `json:"name"`       // 股票名称
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
	GroupID   string `json:"group_id"`   // 所属分组ID
}

// UpdateFavoriteRequest 更新收藏股票请求
type UpdateFavoriteRequest struct {
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
	GroupID   string `json:"group_id"`   // 所属分组ID
	SortOrder int    `json:"sort_order"` // 排序顺序
}

// CreateGroupRequest 创建分组请求
type CreateGroupRequest struct {
	Name  string `json:"name"`  // 分组名称
	Color string `json:"color"` // 分组颜色
}

// UpdateGroupRequest 更新分组请求
type UpdateGroupRequest struct {
	Name      string `json:"name"`       // 分组名称
	Color     string `json:"color"`      // 分组颜色
	SortOrder int    `json:"sort_order"` // 排序顺序
}

// UpdateFavoritesOrderRequest 更新收藏排序请求
type UpdateFavoritesOrderRequest struct {
	FavoriteOrders []FavoriteOrderItem `json:"favorite_orders"` // 收藏排序列表
}

// FavoriteOrderItem 收藏排序项
type FavoriteOrderItem struct {
	ID        string `json:"id"`         // 收藏ID
	GroupID   string `json:"group_id"`   // 所属分组ID
	SortOrder int    `json:"sort_order"` // 排序顺序
}

// HealthStatus 健康检查状态
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
}

// SimpleIndicator 简单技术指标
type SimpleIndicator struct {
	MA5            string `json:"ma5"`
	MA10           string `json:"ma10"`
	MA20           string `json:"ma20"`
	CurrentPrice   string `json:"current_price"`
	PriceChange    string `json:"price_change"`
	PriceChangePct string `json:"price_change_pct"`
	Trend          string `json:"trend"`
	LastUpdate     string `json:"last_update"`
}

// FavoriteSignal 收藏股票信号
type FavoriteSignal struct {
	ID           string          `json:"id"`
	TSCode       string          `json:"ts_code"`
	Name         string          `json:"name"`
	GroupID      string          `json:"group_id"`
	CurrentPrice string          `json:"current_price"`
	TradeDate    string          `json:"trade_date"`
	Indicators   SimpleIndicator `json:"indicators"`
	Predictions  interface{}     `json:"predictions"` // 保持interface{}因为预测结构可能变化
	UpdatedAt    string          `json:"updated_at"`
}

// FavoritesSignalsResponse 收藏股票信号响应
type FavoritesSignalsResponse struct {
	Total   int              `json:"total"`
	Signals []FavoriteSignal `json:"signals"`
}

// CandlestickPattern 蜡烛图模式识别结果
type CandlestickPattern struct {
	TSCode      string      `json:"ts_code"`      // 股票代码
	TradeDate   string      `json:"trade_date"`   // 交易日期
	Pattern     string      `json:"pattern"`      // 图形名称
	Signal      string      `json:"signal"`       // 买卖信号
	Confidence  JSONDecimal `json:"confidence"`   // 置信度
	Description string      `json:"description"`  // 图形描述
	Strength    string      `json:"strength"`     // 信号强度 (STRONG, MEDIUM, WEAK)
	Volume      JSONDecimal `json:"volume"`       // 成交量
	PriceChange JSONDecimal `json:"price_change"` // 价格变化
}

// VolumePricePattern 量价图形识别结果
type VolumePricePattern struct {
	TSCode      string      `json:"ts_code"`      // 股票代码
	TradeDate   string      `json:"trade_date"`   // 交易日期
	Pattern     string      `json:"pattern"`      // 图形名称
	Signal      string      `json:"signal"`       // 买卖信号
	Confidence  JSONDecimal `json:"confidence"`   // 置信度
	Description string      `json:"description"`  // 图形描述
	Strength    string      `json:"strength"`     // 信号强度
	Volume      JSONDecimal `json:"volume"`       // 成交量
	PriceChange JSONDecimal `json:"price_change"` // 价格变化
	VolumeRatio JSONDecimal `json:"volume_ratio"` // 量比
}

// PatternRecognitionResult 图形识别结果
type PatternRecognitionResult struct {
	TSCode            string               `json:"ts_code"`            // 股票代码
	TradeDate         string               `json:"trade_date"`         // 交易日期
	Candlestick       []CandlestickPattern `json:"candlestick"`        // 蜡烛图模式
	VolumePrice       []VolumePricePattern `json:"volume_price"`       // 量价图形
	CombinedSignal    string               `json:"combined_signal"`    // 综合信号
	OverallConfidence JSONDecimal          `json:"overall_confidence"` // 综合置信度
	RiskLevel         string               `json:"risk_level"`         // 风险等级
}

// PatternSearchRequest 图形搜索请求
type PatternSearchRequest struct {
	TSCode        string   `json:"ts_code"`        // 股票代码
	StartDate     string   `json:"start_date"`     // 开始日期
	EndDate       string   `json:"end_date"`       // 结束日期
	Patterns      []string `json:"patterns"`       // 要搜索的图形类型
	MinConfidence float64  `json:"min_confidence"` // 最小置信度
}

// PatternSearchResponse 图形搜索响应
type PatternSearchResponse struct {
	Total   int                        `json:"total"`
	Results []PatternRecognitionResult `json:"results"`
}

// PatternSummary 图形模式摘要
type PatternSummary struct {
	TSCode    string         `json:"ts_code"`    // 股票代码
	Period    int            `json:"period"`     // 统计周期（天数）
	StartDate string         `json:"start_date"` // 开始日期
	EndDate   string         `json:"end_date"`   // 结束日期
	Patterns  map[string]int `json:"patterns"`   // 各种图形模式统计
	Signals   map[string]int `json:"signals"`    // 各种信号统计
	UpdatedAt time.Time      `json:"updated_at"` // 更新时间
}

// RecentSignal 最近的图形信号
type RecentSignal struct {
	TSCode      string      `json:"ts_code"`     // 股票代码
	TradeDate   string      `json:"trade_date"`  // 交易日期
	Pattern     string      `json:"pattern"`     // 图形名称
	Signal      string      `json:"signal"`      // 买卖信号
	Confidence  JSONDecimal `json:"confidence"`  // 置信度
	Description string      `json:"description"` // 图形描述
	Strength    string      `json:"strength"`    // 信号强度
	Type        string      `json:"type"`        // 信号类型
}
