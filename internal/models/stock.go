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
	// 新增动量因子
	WR       *WilliamsRIndicator `json:"wr,omitempty"`
	Momentum *MomentumIndicator  `json:"momentum,omitempty"`
	ROC      *ROCIndicator       `json:"roc,omitempty"`
	// 新增趋势因子
	ADX      *ADXIndicator      `json:"adx,omitempty"`
	SAR      *SARIndicator      `json:"sar,omitempty"`
	Ichimoku *IchimokuIndicator `json:"ichimoku,omitempty"`
	// 新增波动率因子
	ATR    *ATRIndicator                  `json:"atr,omitempty"`
	StdDev *StdDevIndicator               `json:"stddev,omitempty"`
	HV     *HistoricalVolatilityIndicator `json:"hv,omitempty"`
	// 新增成交量因子
	VWAP   *VWAPIndicator   `json:"vwap,omitempty"`
	ADLine *ADLineIndicator `json:"ad_line,omitempty"`
	EMV    *EMVIndicator    `json:"emv,omitempty"`
	VPT    *VPTIndicator    `json:"vpt,omitempty"`
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
	Type           string      `json:"type"`             // "BUY" 或 "SELL"
	Price          JSONDecimal `json:"price"`            // 预测价格
	Date           string      `json:"date"`             // 预测日期
	Probability    JSONDecimal `json:"probability"`      // 概率
	Reason         string      `json:"reason"`           // 预测理由
	Indicators     []string    `json:"indicators"`       // 相关指标
	SignalDate     string      `json:"signal_date"`      // 信号产生的日期（基于哪一天的数据）
	Backtested     bool        `json:"backtested"`       // 是否已回测
	IsCorrect      bool        `json:"is_correct"`       // 预测是否正确
	NextDayPrice   JSONDecimal `json:"next_day_price"`   // 第二天收盘价
	PriceDiff      JSONDecimal `json:"price_diff"`       // 价格差值
	PriceDiffRatio JSONDecimal `json:"price_diff_ratio"` // 价格差值百分比
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
	Total             int              `json:"total"`
	Signals           []FavoriteSignal `json:"signals"`
	Calculating       bool             `json:"calculating"`
	CalculationStatus interface{}      `json:"calculation_status,omitempty"`
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

// StockSignal 股票信号存储结构
type StockSignal struct {
	ID                  string                     `json:"id"`                             // 唯一标识
	TSCode              string                     `json:"ts_code"`                        // 股票代码
	Name                string                     `json:"name"`                           // 股票名称
	TradeDate           string                     `json:"trade_date"`                     // 信号基于的交易日期
	SignalDate          string                     `json:"signal_date"`                    // 信号计算日期
	SignalType          string                     `json:"signal_type"`                    // 信号类型: BUY, SELL, HOLD
	SignalStrength      string                     `json:"signal_strength"`                // 信号强度: STRONG, MEDIUM, WEAK
	Confidence          JSONDecimal                `json:"confidence"`                     // 置信度 0-1
	Patterns            []PatternRecognitionResult `json:"patterns,omitempty"`             // 识别到的图形模式
	TechnicalIndicators *TechnicalIndicators       `json:"technical_indicators,omitempty"` // 技术指标数据
	Predictions         *PredictionResult          `json:"predictions,omitempty"`          // 预测数据
	Description         string                     `json:"description"`                    // 信号描述
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
}

// SignalCalculationRequest 信号计算请求
type SignalCalculationRequest struct {
	TSCode    string `json:"ts_code"`    // 股票代码
	Name      string `json:"name"`       // 股票名称
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
	Force     bool   `json:"force"`      // 是否强制重新计算
}

// SignalCalculationResponse 信号计算响应
type SignalCalculationResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Signal  *StockSignal `json:"signal,omitempty"`
	Error   string       `json:"error,omitempty"`
}

// BatchSignalRequest 批量信号计算请求
type BatchSignalRequest struct {
	TSCodes []string `json:"ts_codes"` // 股票代码列表
	Force   bool     `json:"force"`    // 是否强制重新计算
}

// BatchSignalResponse 批量信号计算响应
type BatchSignalResponse struct {
	Total     int                         `json:"total"`   // 总数
	Success   int                         `json:"success"` // 成功数
	Failed    int                         `json:"failed"`  // 失败数
	Results   []SignalCalculationResponse `json:"results"` // 详细结果
	StartTime time.Time                   `json:"start_time"`
	EndTime   time.Time                   `json:"end_time"`
	Duration  string                      `json:"duration"`
}

// 新增技术指标结构定义

// WilliamsRIndicator 威廉指标 (%R)
type WilliamsRIndicator struct {
	WR14   JSONDecimal `json:"wr14"`   // 14日威廉指标
	Signal string      `json:"signal"` // 超买/超卖信号
}

// MomentumIndicator 动量指标
type MomentumIndicator struct {
	Momentum10 JSONDecimal `json:"momentum10"` // 10日动量
	Momentum20 JSONDecimal `json:"momentum20"` // 20日动量
	Signal     string      `json:"signal"`     // 买卖信号
}

// ROCIndicator 变化率指标
type ROCIndicator struct {
	ROC10  JSONDecimal `json:"roc10"`  // 10日变化率
	ROC20  JSONDecimal `json:"roc20"`  // 20日变化率
	Signal string      `json:"signal"` // 买卖信号
}

// ADXIndicator 平均方向指数
type ADXIndicator struct {
	ADX    JSONDecimal `json:"adx"`    // 平均方向指数
	PDI    JSONDecimal `json:"pdi"`    // 正方向指数
	MDI    JSONDecimal `json:"mdi"`    // 负方向指数
	Signal string      `json:"signal"` // 趋势强度信号
}

// SARIndicator 抛物线转向指标
type SARIndicator struct {
	SAR    JSONDecimal `json:"sar"`    // SAR值
	Signal string      `json:"signal"` // 买卖信号
}

// IchimokuIndicator 一目均衡表
type IchimokuIndicator struct {
	TenkanSen   JSONDecimal `json:"tenkan_sen"`    // 转换线
	KijunSen    JSONDecimal `json:"kijun_sen"`     // 基准线
	SenkouSpanA JSONDecimal `json:"senkou_span_a"` // 先行带A
	SenkouSpanB JSONDecimal `json:"senkou_span_b"` // 先行带B
	ChikouSpan  JSONDecimal `json:"chikou_span"`   // 滞后线
	Signal      string      `json:"signal"`        // 买卖信号
}

// ATRIndicator 平均真实范围
type ATRIndicator struct {
	ATR14  JSONDecimal `json:"atr14"`  // 14日ATR
	Signal string      `json:"signal"` // 波动率信号
}

// StdDevIndicator 标准差指标
type StdDevIndicator struct {
	StdDev20 JSONDecimal `json:"stddev20"` // 20日标准差
	Signal   string      `json:"signal"`   // 波动率信号
}

// HistoricalVolatilityIndicator 历史波动率
type HistoricalVolatilityIndicator struct {
	HV20   JSONDecimal `json:"hv20"`   // 20日历史波动率
	HV60   JSONDecimal `json:"hv60"`   // 60日历史波动率
	Signal string      `json:"signal"` // 波动率信号
}

// VWAPIndicator 成交量加权平均价
type VWAPIndicator struct {
	VWAP   JSONDecimal `json:"vwap"`   // 成交量加权平均价
	Signal string      `json:"signal"` // 买卖信号
}

// ADLineIndicator 累积/派发线
type ADLineIndicator struct {
	ADLine JSONDecimal `json:"ad_line"` // 累积/派发线
	Signal string      `json:"signal"`  // 买卖信号
}

// EMVIndicator 简易波动指标
type EMVIndicator struct {
	EMV14  JSONDecimal `json:"emv14"`  // 14日EMV
	Signal string      `json:"signal"` // 买卖信号
}

// VPTIndicator 量价确认指标
type VPTIndicator struct {
	VPT    JSONDecimal `json:"vpt"`    // 量价确认指标
	Signal string      `json:"signal"` // 买卖信号
}

// ===== 基本面数据结构 =====

// FinancialStatement 财务报表基础结构
type FinancialStatement struct {
	TSCode     string `json:"ts_code"`     // 股票代码
	AnnDate    string `json:"ann_date"`    // 公告日期
	FDate      string `json:"f_date"`      // 报告期
	EndDate    string `json:"end_date"`    // 报告期结束日期
	ReportType string `json:"report_type"` // 报告类型：1-年报，2-中报，3-季报
	CompType   string `json:"comp_type"`   // 公司类型：1-一般工商业，2-银行，3-保险，4-证券
}

// IncomeStatement 利润表数据
type IncomeStatement struct {
	FinancialStatement
	// 营业收入相关
	Revenue       JSONDecimal `json:"revenue"`        // 营业总收入
	OperRevenue   JSONDecimal `json:"oper_revenue"`   // 营业收入
	IntIncome     JSONDecimal `json:"int_income"`     // 利息收入
	CommissionInc JSONDecimal `json:"commission_inc"` // 手续费及佣金收入

	// 成本费用
	OperCost JSONDecimal `json:"oper_cost"` // 营业总成本
	OperExp  JSONDecimal `json:"oper_exp"`  // 营业费用
	AdminExp JSONDecimal `json:"admin_exp"` // 管理费用
	FinExp   JSONDecimal `json:"fin_exp"`   // 财务费用
	RdExp    JSONDecimal `json:"rd_exp"`    // 研发费用

	// 利润相关
	OperProfit    JSONDecimal `json:"oper_profit"`     // 营业利润
	TotalProfit   JSONDecimal `json:"total_profit"`    // 利润总额
	NetProfit     JSONDecimal `json:"net_profit"`      // 净利润
	NetProfitDedt JSONDecimal `json:"net_profit_dedt"` // 扣非净利润

	// 每股收益
	BasicEps   JSONDecimal `json:"basic_eps"`   // 基本每股收益
	DilutedEps JSONDecimal `json:"diluted_eps"` // 稀释每股收益
}

// BalanceSheet 资产负债表数据
type BalanceSheet struct {
	FinancialStatement
	// 资产
	TotalAssets     JSONDecimal `json:"total_assets"`     // 资产总计
	TotalCurAssets  JSONDecimal `json:"total_cur_assets"` // 流动资产合计
	Money           JSONDecimal `json:"money"`            // 货币资金
	TradAssets      JSONDecimal `json:"trad_assets"`      // 交易性金融资产
	NotesReceiv     JSONDecimal `json:"notes_receiv"`     // 应收票据
	AccountsReceiv  JSONDecimal `json:"accounts_receiv"`  // 应收账款
	InventoryAssets JSONDecimal `json:"inventory_assets"` // 存货
	TotalNcaAssets  JSONDecimal `json:"total_nca_assets"` // 非流动资产合计
	FixAssets       JSONDecimal `json:"fix_assets"`       // 固定资产
	CipAssets       JSONDecimal `json:"cip_assets"`       // 在建工程
	IntangAssets    JSONDecimal `json:"intang_assets"`    // 无形资产

	// 负债
	TotalLiab       JSONDecimal `json:"total_liab"`       // 负债合计
	TotalCurLiab    JSONDecimal `json:"total_cur_liab"`   // 流动负债合计
	ShortLoan       JSONDecimal `json:"short_loan"`       // 短期借款
	NotesPayable    JSONDecimal `json:"notes_payable"`    // 应付票据
	AccountsPayable JSONDecimal `json:"accounts_payable"` // 应付账款
	TotalNcaLiab    JSONDecimal `json:"total_nca_liab"`   // 非流动负债合计
	LongLoan        JSONDecimal `json:"long_loan"`        // 长期借款

	// 所有者权益
	TotalHldrEqy  JSONDecimal `json:"total_hldr_eqy"` // 所有者权益合计
	CapRese       JSONDecimal `json:"cap_rese"`       // 资本公积
	UndistrProfit JSONDecimal `json:"undistr_profit"` // 未分配利润
	TotalShare    JSONDecimal `json:"total_share"`    // 实收资本(或股本)
}

// CashFlowStatement 现金流量表数据
type CashFlowStatement struct {
	FinancialStatement
	// 经营活动现金流量
	NetCashOperAct    JSONDecimal `json:"net_cash_oper_act"`    // 经营活动产生的现金流量净额
	CashRecrSale      JSONDecimal `json:"cash_recr_sale"`       // 销售商品、提供劳务收到的现金
	CashPayGoods      JSONDecimal `json:"cash_pay_goods"`       // 购买商品、接受劳务支付的现金
	CashPayBehalfEmpl JSONDecimal `json:"cash_pay_behalf_empl"` // 支付给职工以及为职工支付的现金
	CashPayTax        JSONDecimal `json:"cash_pay_tax"`         // 支付的各项税费

	// 投资活动现金流量
	NetCashInvAct JSONDecimal `json:"net_cash_inv_act"` // 投资活动产生的现金流量净额
	CashRecvDisp  JSONDecimal `json:"cash_recv_disp"`   // 收回投资收到的现金
	CashPayAcq    JSONDecimal `json:"cash_pay_acq"`     // 投资支付的现金

	// 筹资活动现金流量
	NetCashFinAct  JSONDecimal `json:"net_cash_fin_act"` // 筹资活动产生的现金流量净额
	CashRecvInvest JSONDecimal `json:"cash_recv_invest"` // 吸收投资收到的现金
	CashPayDist    JSONDecimal `json:"cash_pay_dist"`    // 分配股利、利润或偿付利息支付的现金

	// 汇率变动影响
	FxEffectCash JSONDecimal `json:"fx_effect_cash"` // 汇率变动对现金及现金等价物的影响

	// 现金净增加额
	NetIncrCashCce JSONDecimal `json:"net_incr_cash_cce"` // 现金及现金等价物净增加额
	CashBegPeriod  JSONDecimal `json:"cash_beg_period"`   // 期初现金及现金等价物余额
	CashEndPeriod  JSONDecimal `json:"cash_end_period"`   // 期末现金及现金等价物余额
}

// FinancialIndicator 财务指标数据
type FinancialIndicator struct {
	FinancialStatement
	// 每股指标
	Eps         JSONDecimal `json:"eps"`          // 每股收益
	DtEps       JSONDecimal `json:"dt_eps"`       // 稀释每股收益
	TotalRevPs  JSONDecimal `json:"total_rev_ps"` // 每股营业总收入
	BvPs        JSONDecimal `json:"bv_ps"`        // 每股净资产
	OcfPs       JSONDecimal `json:"ocf_ps"`       // 每股经营现金流
	NetprofitPs JSONDecimal `json:"netprofit_ps"` // 每股净利润

	// 成长能力指标
	OrLastYear  JSONDecimal `json:"or_last_year"`  // 营业收入同比增长率(%)
	OpLastYear  JSONDecimal `json:"op_last_year"`  // 营业利润同比增长率(%)
	EpsLastYear JSONDecimal `json:"eps_last_year"` // 每股收益同比增长率(%)
	NetprofitGr JSONDecimal `json:"netprofit_gr"`  // 净利润同比增长率(%)

	// 盈利能力指标
	Roe         JSONDecimal `json:"roe"`          // 净资产收益率
	Roa         JSONDecimal `json:"roa"`          // 资产收益率
	GrossMargin JSONDecimal `json:"gross_margin"` // 销售毛利率
	NetMargin   JSONDecimal `json:"net_margin"`   // 销售净利率
	OperMargin  JSONDecimal `json:"oper_margin"`  // 营业利润率
	EbitMargin  JSONDecimal `json:"ebit_margin"`  // EBIT利润率

	// 运营能力指标
	InvTurn    JSONDecimal `json:"inv_turn"`    // 存货周转率
	ArTurn     JSONDecimal `json:"ar_turn"`     // 应收账款周转率
	CaTurn     JSONDecimal `json:"ca_turn"`     // 流动资产周转率
	FaTurn     JSONDecimal `json:"fa_turn"`     // 固定资产周转率
	AssetsTurn JSONDecimal `json:"assets_turn"` // 总资产周转率

	// 偿债能力指标
	CurrentRatio JSONDecimal `json:"current_ratio"`  // 流动比率
	QuickRatio   JSONDecimal `json:"quick_ratio"`    // 速动比率
	CashRatio    JSONDecimal `json:"cash_ratio"`     // 现金比率
	LrRatio      JSONDecimal `json:"lr_ratio"`       // 资产负债率
	DebtToAssets JSONDecimal `json:"debt_to_assets"` // 负债合计/资产总计
	DebtToEqt    JSONDecimal `json:"debt_to_eqt"`    // 产权比率

	// 现金流指标
	CfPs          JSONDecimal `json:"cf_ps"`           // 每股现金流量净额
	CfNetprofitPs JSONDecimal `json:"cf_netprofit_ps"` // 每股经营现金流净额
	CfLiab        JSONDecimal `json:"cf_liab"`         // 现金流量负债比
}

// DailyBasic 每日基本面指标
type DailyBasic struct {
	TSCode      string      `json:"ts_code"`      // 股票代码
	TradeDate   string      `json:"trade_date"`   // 交易日期
	Close       JSONDecimal `json:"close"`        // 当日收盘价
	Turnover    JSONDecimal `json:"turnover"`     // 换手率（%）
	VolumeRatio JSONDecimal `json:"volume_ratio"` // 量比
	Pe          JSONDecimal `json:"pe"`           // 市盈率（总市值/净利润，亏损的PE为空）
	PeTtm       JSONDecimal `json:"pe_ttm"`       // 市盈率（TTM，亏损的PE为空）
	Pb          JSONDecimal `json:"pb"`           // 市净率（总市值/净资产）
	Ps          JSONDecimal `json:"ps"`           // 市销率
	PsTtm       JSONDecimal `json:"ps_ttm"`       // 市销率（TTM）
	DvRatio     JSONDecimal `json:"dv_ratio"`     // 股息率（%）
	DvTtm       JSONDecimal `json:"dv_ttm"`       // 股息率（TTM）（%）
	TotalShare  JSONDecimal `json:"total_share"`  // 总股本（万股）
	FloatShare  JSONDecimal `json:"float_share"`  // 流通股本（万股）
	FreeShare   JSONDecimal `json:"free_share"`   // 自由流通股本（万）
	TotalMv     JSONDecimal `json:"total_mv"`     // 总市值（万元）
	CircMv      JSONDecimal `json:"circ_mv"`      // 流通市值（万元）
}

// FundamentalFactor 基本面因子
type FundamentalFactor struct {
	TSCode    string `json:"ts_code"`    // 股票代码
	TradeDate string `json:"trade_date"` // 计算日期

	// 价值因子
	PE         JSONDecimal `json:"pe"`           // 市盈率
	PB         JSONDecimal `json:"pb"`           // 市净率
	PS         JSONDecimal `json:"ps"`           // 市销率
	PCF        JSONDecimal `json:"pcf"`          // 市现率
	EVToEBITDA JSONDecimal `json:"ev_to_ebitda"` // EV/EBITDA

	// 成长因子
	RevenueGrowth   JSONDecimal `json:"revenue_growth"`    // 营收增长率
	NetProfitGrowth JSONDecimal `json:"net_profit_growth"` // 净利润增长率
	EPSGrowth       JSONDecimal `json:"eps_growth"`        // EPS增长率
	ROEGrowth       JSONDecimal `json:"roe_growth"`        // ROE增长率

	// 质量因子
	ROE          JSONDecimal `json:"roe"`            // 净资产收益率
	ROA          JSONDecimal `json:"roa"`            // 资产收益率
	GrossMargin  JSONDecimal `json:"gross_margin"`   // 毛利率
	NetMargin    JSONDecimal `json:"net_margin"`     // 净利率
	DebtToAssets JSONDecimal `json:"debt_to_assets"` // 资产负债率
	CurrentRatio JSONDecimal `json:"current_ratio"`  // 流动比率
	QuickRatio   JSONDecimal `json:"quick_ratio"`    // 速动比率

	// 盈利因子
	ROIC            JSONDecimal `json:"roic"`             // 投入资本回报率
	OperatingMargin JSONDecimal `json:"operating_margin"` // 营业利润率
	EBITDAMargin    JSONDecimal `json:"ebitda_margin"`    // EBITDA利润率

	// 运营效率因子
	AssetTurnover      JSONDecimal `json:"asset_turnover"`      // 总资产周转率
	InventoryTurnover  JSONDecimal `json:"inventory_turnover"`  // 存货周转率
	ReceivableTurnover JSONDecimal `json:"receivable_turnover"` // 应收账款周转率

	// 现金流因子
	OCFToRevenue   JSONDecimal `json:"ocf_to_revenue"`    // 经营现金流/营收比
	OCFToNetProfit JSONDecimal `json:"ocf_to_net_profit"` // 经营现金流/净利润比
	FCFYield       JSONDecimal `json:"fcf_yield"`         // 自由现金流收益率

	// 分红因子
	DividendYield JSONDecimal `json:"dividend_yield"` // 股息率
	PayoutRatio   JSONDecimal `json:"payout_ratio"`   // 分红率

	// 因子标准化得分
	ValueScore         JSONDecimal `json:"value_score"`         // 价值因子得分
	GrowthScore        JSONDecimal `json:"growth_score"`        // 成长因子得分
	QualityScore       JSONDecimal `json:"quality_score"`       // 质量因子得分
	ProfitabilityScore JSONDecimal `json:"profitability_score"` // 盈利因子得分
	CompositeScore     JSONDecimal `json:"composite_score"`     // 综合得分

	// 行业和市场相对位置
	IndustryRank       int         `json:"industry_rank"`       // 行业排名
	MarketRank         int         `json:"market_rank"`         // 全市场排名
	IndustryPercentile JSONDecimal `json:"industry_percentile"` // 行业分位数
	MarketPercentile   JSONDecimal `json:"market_percentile"`   // 市场分位数

	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}
