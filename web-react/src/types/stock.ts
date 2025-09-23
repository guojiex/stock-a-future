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

// 技术指标数据
export interface TechnicalIndicators {
  ts_code: string;
  trade_date: string;
  
  // 移动平均线
  ma5?: number;
  ma10?: number;
  ma20?: number;
  ma30?: number;
  ma60?: number;
  
  // MACD指标
  macd_dif?: number;
  macd_dea?: number;
  macd_histogram?: number;
  
  // RSI指标
  rsi6?: number;
  rsi12?: number;
  rsi24?: number;
  
  // 布林带
  boll_upper?: number;
  boll_mid?: number;
  boll_lower?: number;
  
  // KDJ指标
  kdj_k?: number;
  kdj_d?: number;
  kdj_j?: number;
  
  // 成交量指标
  volume_ratio?: number;
  turnover_rate?: number;
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
  id: number;
  ts_code: string;
  name: string;
  group_id?: number;
  notes?: string;
  created_at: string;
  updated_at: string;
  order_index?: number;
}

// 添加收藏请求
export interface AddFavoriteRequest {
  ts_code: string;
  name: string;
  group_id?: number;
  notes?: string;
}

// 股票组
export interface StockGroup {
  id: number;
  name: string;
  description?: string;
  created_at: string;
  updated_at: string;
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
