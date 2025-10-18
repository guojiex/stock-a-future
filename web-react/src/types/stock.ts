/**
 * 股票相关数据类型定义 - Web版本
 * 对应Go后端的models包中的结构体
 */

// 股票基本信息
export interface StockBasic {
  ts_code: string;        // 股票代码
  symbol: string;         // 股票代码(不含后缀)
  name: string;           // 股票名称
  area?: string;          // 地域
  industry?: string;      // 行业
  market?: string;        // 市场类型
  list_date?: string;     // 上市日期
  is_hs?: string;         // 是否沪深港通
}

// 股票日线数据
export interface StockDaily {
  ts_code: string;        // 股票代码
  trade_date: string;     // 交易日期
  open: string;           // 开盘价
  high: string;           // 最高价
  low: string;            // 最低价
  close: string;          // 收盘价
  pre_close?: string;     // 前收盘价
  change?: string;        // 涨跌额
  pct_chg?: string;       // 涨跌幅
  vol: string;            // 成交量
  amount?: string;        // 成交额
}

// 技术指标数据 - 匹配后端Go结构
export interface TechnicalIndicators {
  ts_code: string;
  trade_date: string;
  
  // 传统指标
  macd?: MACDIndicator;
  rsi?: RSIIndicator;
  boll?: BollingerBandsIndicator;
  ma?: MovingAverageIndicator;
  kdj?: KDJIndicator;
  
  // 动量因子
  wr?: WilliamsRIndicator;
  momentum?: MomentumIndicator;
  roc?: ROCIndicator;
  
  // 趋势因子
  adx?: ADXIndicator;
  sar?: SARIndicator;
  ichimoku?: IchimokuIndicator;
  
  // 波动率因子
  atr?: ATRIndicator;
  stddev?: StdDevIndicator;
  hv?: HistoricalVolatilityIndicator;
  
  // 成交量因子
  vwap?: VWAPIndicator;
  ad_line?: ADLineIndicator;
  emv?: EMVIndicator;
  vpt?: VPTIndicator;
}

// MACD指标
export interface MACDIndicator {
  dif: number;
  dea: number;
  histogram: number;
  signal?: string;
}

// RSI指标
export interface RSIIndicator {
  rsi14: number;
  signal?: string;
}

// 布林带指标
export interface BollingerBandsIndicator {
  upper: number;
  middle: number;
  lower: number;
  signal?: string;
}

// 移动平均线指标
export interface MovingAverageIndicator {
  ma5: number;
  ma10: number;
  ma20: number;
  ma60: number;
  ma120: number;
}

// KDJ指标
export interface KDJIndicator {
  k: number;
  d: number;
  j: number;
  signal?: string;
}

// 威廉指标
export interface WilliamsRIndicator {
  wr14: number;
  signal?: string;
}

// 动量指标
export interface MomentumIndicator {
  momentum10: number;
  momentum20: number;
  signal?: string;
}

// 变化率指标
export interface ROCIndicator {
  roc10: number;
  roc20: number;
  signal?: string;
}

// ADX指标
export interface ADXIndicator {
  adx: number;
  pdi: number;
  mdi: number;
  signal?: string;
}

// SAR指标
export interface SARIndicator {
  sar: number;
  signal?: string;
}

// 一目均衡表
export interface IchimokuIndicator {
  tenkan_sen: number;
  kijun_sen: number;
  senkou_span_a: number;
  senkou_span_b: number;
  signal?: string;
}

// ATR指标
export interface ATRIndicator {
  atr14: number;
  signal?: string;
}

// 标准差指标
export interface StdDevIndicator {
  stddev20: number;
  signal?: string;
}

// 历史波动率指标
export interface HistoricalVolatilityIndicator {
  hv20: number;
  hv60: number;
  signal?: string;
}

// VWAP指标
export interface VWAPIndicator {
  vwap: number;
  signal?: string;
}

// A/D线指标
export interface ADLineIndicator {
  ad_line: number;
  signal?: string;
}

// EMV指标
export interface EMVIndicator {
  emv14: number;
  signal?: string;
}

// VPT指标
export interface VPTIndicator {
  vpt: number;
  signal?: string;
}

// 基本面数据
export interface FundamentalData {
  stock_basic?: StockBasic;
  income_statement?: IncomeStatement;
  balance_sheet?: BalanceSheet;
  cash_flow_statement?: CashFlowStatement;
  daily_basic?: DailyBasic;
}

// 利润表（简化版）
export interface IncomeStatement {
  ts_code: string;
  ann_date: string;       // 公告日期
  end_date: string;       // 报告期
  report_type: string;    // 报告类型
  basic_eps: number;      // 基本每股收益
  total_revenue: number;  // 营业总收入
  revenue: number;        // 营业收入
  total_profit: number;   // 利润总额
  n_income: number;       // 净利润
}

// 资产负债表（简化版）
export interface BalanceSheet {
  ts_code: string;
  ann_date: string;       // 公告日期
  end_date: string;       // 报告期
  total_share: number;    // 期末总股本
  total_assets: number;   // 资产总计
  total_liab: number;     // 负债合计
  total_hldr_eqy_inc_min_int: number; // 股东权益合计
}

// 现金流量表（简化版）
export interface CashFlowStatement {
  ts_code: string;
  ann_date: string;       // 公告日期
  end_date: string;       // 报告期
  n_cashflow_operate_a: number; // 经营活动现金流量净额
  n_cashflow_invest_a: number;  // 投资活动现金流量净额
  n_cash_flows_fnc_act: number; // 筹资活动现金流量净额
}

// 每日基本面数据
export interface DailyBasic {
  ts_code: string;
  trade_date: string;
  close?: number;         // 当日收盘价
  turnover_rate?: number; // 换手率
  pe?: number;           // 市盈率
  pb?: number;           // 市净率
  ps?: number;           // 市销率
  total_mv?: number;     // 总市值
  circ_mv?: number;      // 流通市值
}

// 收藏股票
export interface Favorite {
  id: string;           // 改为string以匹配后端UUID
  ts_code: string;
  name: string;
  start_date?: string;  // 添加开始日期
  end_date?: string;    // 添加结束日期
  group_id?: string;    // 改为string
  notes?: string;
  sort_order?: number;  // 改为sort_order
  created_at: string;
  updated_at: string;
}

// 添加收藏请求
export interface AddFavoriteRequest {
  ts_code: string;
  name: string;
  start_date?: string;
  end_date?: string;
  group_id?: string;
  notes?: string;
}

// 收藏分组
export interface FavoriteGroup {
  id: string;
  name: string;
  color?: string;       // 添加颜色
  sort_order?: number;
  created_at: string;
  updated_at: string;
}

// 创建分组请求
export interface CreateGroupRequest {
  name: string;
  color?: string;
}

// 更新分组请求
export interface UpdateGroupRequest {
  name?: string;
  color?: string;
  sort_order?: number;
}

// 更新收藏请求
export interface UpdateFavoriteRequest {
  start_date?: string;
  end_date?: string;
  group_id?: string;
  sort_order?: number;
}

// API响应格式
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

// 分页响应
export interface PaginatedResponse<T = any> {
  total: number;
  page: number;
  page_size: number;
  data: T[];
}

// 搜索请求参数
export interface SearchParams {
  q: string;              // 搜索关键词
  limit?: number;         // 限制数量
  offset?: number;        // 偏移量
}

// 日期范围
export interface DateRange {
  start_date: string;     // 开始日期 YYYYMMDD
  end_date: string;       // 结束日期 YYYYMMDD
}

// 图表数据点
export interface ChartDataPoint {
  x: number | string;     // X轴数据
  y: number;              // Y轴数据
  date?: string;          // 日期字符串
}

// K线图数据
export interface CandlestickData extends ChartDataPoint {
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
}

// 买卖点预测
export interface TradingPointPrediction {
  type: 'BUY' | 'SELL';           // 买入或卖出
  price: number;                   // 预测价格
  date: string;                    // 预测日期 YYYYMMDD
  probability: number;             // 概率 (0-1)
  reason: string;                  // 预测理由
  indicators: string[];            // 相关指标
  signal_date: string;             // 信号产生的日期（基于哪一天的数据）
  backtested: boolean;             // 是否已回测
  is_correct?: boolean;            // 预测是否正确
  next_day_price?: number;         // 第二天收盘价
  price_diff?: number;             // 价格差值
  price_diff_ratio?: number;       // 价格差值百分比
}

// 预测结果
export interface PredictionResult {
  ts_code: string;
  trade_date: string;
  predictions: TradingPointPrediction[];
  confidence: number;              // 预测置信度 (0-1)
  updated_at: string;
}

// 形态模式摘要
export interface PatternSummary {
  ts_code: string;                 // 股票代码
  period: number;                  // 统计周期（天数）
  start_date: string;              // 开始日期 YYYYMMDD
  end_date: string;                // 结束日期 YYYYMMDD
  patterns: Record<string, number>; // 各种图形模式统计（形态名称 -> 出现次数）
  signals: Record<string, number>;  // 各种信号统计（信号类型 -> 出现次数）
  updated_at: string;              // 更新时间
}

// 收藏股票的信号数据
export interface FavoriteSignal {
  id: string;                      // 收藏ID
  ts_code: string;                 // 股票代码
  name: string;                    // 股票名称
  group_id?: string;               // 分组ID
  current_price: string;           // 当前价格
  trade_date: string;              // 交易日期
  indicators: TechnicalIndicators; // 技术指标
  predictions: TradingPointPrediction[]; // 买卖预测
  updated_at: string;              // 更新时间
}

// 收藏股票信号汇总响应
export interface FavoritesSignalsResponse {
  total: number;                   // 总数
  signals: FavoriteSignal[];       // 信号列表
  calculating?: boolean;           // 是否正在计算
  calculation_status?: {           // 计算状态
    status: string;
    message: {
      message: string;
    };
    detail: {
      is_calculating: boolean;
      completed: number;
      total: number;
      last_updated: string;
    };
  };
}