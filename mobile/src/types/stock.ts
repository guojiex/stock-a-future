/**
 * 股票相关数据类型定义
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

// 利润表
export interface IncomeStatement {
  ts_code: string;
  ann_date: string;       // 公告日期
  end_date: string;       // 报告期
  report_type: string;    // 报告类型
  comp_type: string;      // 公司类型
  basic_eps: number;      // 基本每股收益
  diluted_eps: number;    // 稀释每股收益
  total_revenue: number;  // 营业总收入
  revenue: number;        // 营业收入
  int_income?: number;    // 利息收入
  prem_earned?: number;   // 已赚保费
  comm_income?: number;   // 手续费及佣金收入
  n_commis_income?: number; // 手续费及佣金净收入
  n_oth_income?: number;  // 其他经营净收益
  n_oth_b_income?: number; // 加:其他业务净收益
  prem_income?: number;   // 保险业务收入
  out_prem?: number;      // 减:分出保费
  une_prem_reser?: number; // 提取未到期责任准备金
  reins_income?: number;  // 其中:分保费收入
  n_sec_tb_income?: number; // 代理买卖证券业务净收入
  n_sec_uw_income?: number; // 证券承销业务净收入
  n_asset_mg_income?: number; // 受托客户资产管理业务净收入
  oth_b_income?: number;  // 其他业务收入
  fv_value_chg_gain?: number; // 加:公允价值变动净收益
  invest_income?: number; // 加:投资净收益
  ass_invest_income?: number; // 其中:对联营企业和合营企业的投资收益
  forex_gain?: number;    // 加:汇兑净收益
  total_cogs: number;     // 营业总成本
  oper_cost: number;      // 减:营业成本
  int_exp?: number;       // 减:利息支出
  comm_exp?: number;      // 减:手续费及佣金支出
  biz_tax_surchg?: number; // 减:营业税金及附加
  sell_exp?: number;      // 减:销售费用
  admin_exp?: number;     // 减:管理费用
  fin_exp?: number;       // 减:财务费用
  assets_impair_loss?: number; // 减:资产减值损失
  prem_refund?: number;   // 退保金
  compens_payout?: number; // 赔付总支出
  reser_insur_liab?: number; // 提取保险责任准备金
  div_payt?: number;      // 保户红利支出
  reins_exp?: number;     // 分保费用
  oper_exp?: number;      // 营业支出
  compens_payout_refu?: number; // 减:摊回赔付支出
  insur_reser_refu?: number; // 减:摊回保险责任准备金
  reins_cost_refund?: number; // 减:摊回分保费用
  other_bus_cost?: number; // 其他业务成本
  operate_profit: number; // 营业利润
  non_oper_income?: number; // 加:营业外收入
  non_oper_exp?: number;  // 减:营业外支出
  nca_disploss?: number;  // 其中:减:非流动资产处置净损失
  total_profit: number;   // 利润总额
  income_tax: number;     // 减:所得税费用
  n_income: number;       // 净利润(含少数股东损益)
  n_income_attr_p: number; // 净利润(不含少数股东损益)
  minority_gain?: number; // 少数股东损益
  oth_compr_income?: number; // 其他综合收益
  t_compr_income: number; // 综合收益总额
  compr_inc_attr_p?: number; // 归属于母公司(或股东)的综合收益总额
  compr_inc_attr_m_s?: number; // 归属于少数股东的综合收益总额
  ebit?: number;          // 息税前利润
  ebitda?: number;        // 息税折旧摊销前利润
  insurance_exp?: number; // 保险业务支出
  undist_profit?: number; // 年初未分配利润
  distable_profit?: number; // 可分配利润
  rd_exp?: number;        // 研发费用
  fin_exp_int_exp?: number; // 财务费用:利息费用
  fin_exp_int_inc?: number; // 财务费用:利息收入
  transfer_surplus_rese?: number; // 盈余公积转增资本
  transfer_housing_imprest?: number; // 住房周转金转增资本
  transfer_oth?: number;  // 其他转增资本
  adj_lossgain?: number;  // 调整以前年度损益
  withdra_legal_surplus?: number; // 提取法定盈余公积金
  withdra_legal_pubfund?: number; // 提取法定公益金
  withdra_biz_devfund?: number; // 提取企业发展基金
  withdra_rese_fund?: number; // 提取储备基金
  withdra_oth_ersu?: number; // 提取任意盈余公积金
  workers_welfare?: number; // 职工奖金福利
  distr_profit_shrhder?: number; // 可供股东分配的利润
  prfshare_payable_dvd?: number; // 应付优先股股利
  comshare_payable_dvd?: number; // 应付普通股股利
  capit_comstock_div?: number; // 转作股本的普通股股利
  net_after_nr_lp_correct?: number;
  credit_impa_loss?: number; // 信用减值损失
  net_expo_hedging_benefits?: number;
  oth_impair_loss_assets?: number; // 其他资产减值损失
  total_opcost?: number;   // 营业总成本2
  amodcost_fin_assets?: number; // 以摊余成本计量的金融资产终止确认收益
  oth_income?: number;     // 其他收益
  asset_disp_income?: number; // 资产处置收益
  continued_net_profit?: number; // 持续经营净利润
  end_net_profit?: number; // 终止经营净利润
}

// 资产负债表
export interface BalanceSheet {
  ts_code: string;
  ann_date: string;       // 公告日期
  end_date: string;       // 报告期
  report_type: string;    // 报告类型
  comp_type: string;      // 公司类型
  total_share: number;    // 期末总股本
  cap_rese: number;       // 资本公积金
  undistr_porfit: number; // 未分配利润
  surplus_rese: number;   // 盈余公积金
  special_rese?: number;  // 专项储备
  money_cap: number;      // 货币资金
  trad_asset: number;     // 交易性金融资产
  notes_receiv: number;   // 应收票据
  accounts_receiv: number; // 应收账款
  oth_receiv: number;     // 其他应收款
  prepayment: number;     // 预付款项
  div_receiv?: number;    // 应收股利
  int_receiv?: number;    // 应收利息
  inventories: number;    // 存货
  amor_exp?: number;      // 待摊费用
  nca_within_1y?: number; // 一年内到期的非流动资产
  sett_rsrv?: number;     // 结算备付金
  loanto_oth_bank_fi?: number; // 拆出资金
  premium_receiv?: number; // 应收保费
  reinsur_receiv?: number; // 应收分保账款
  reinsur_res_receiv?: number; // 应收分保合同准备金
  pur_resale_fa?: number; // 买入返售金融资产
  oth_cur_assets: number; // 其他流动资产
  total_cur_assets: number; // 流动资产合计
  fa_avail_for_sale?: number; // 可供出售金融资产
  htm_invest?: number;    // 持有至到期投资
  lt_eqt_invest: number;  // 长期股权投资
  invest_real_estate: number; // 投资性房地产
  time_deposits?: number; // 定期存款
  oth_assets?: number;    // 其他资产
  lt_rec: number;         // 长期应收款
  fix_assets: number;     // 固定资产
  cip?: number;           // 在建工程
  const_materials?: number; // 工程物资
  fixed_assets_disp?: number; // 固定资产清理
  produc_bio_assets?: number; // 生产性生物资产
  oil_and_gas_assets?: number; // 油气资产
  intan_assets: number;   // 无形资产
  r_and_d?: number;       // 研发支出
  goodwill: number;       // 商誉
  lt_amor_exp?: number;   // 长期待摊费用
  defer_tax_assets: number; // 递延所得税资产
  decr_in_disbur?: number; // 发放贷款及垫款
  oth_nca: number;        // 其他非流动资产
  total_nca: number;      // 非流动资产合计
  cash_reser_cb?: number; // 现金及存放中央银行款项
  depos_in_oth_bfi?: number; // 存放同业和其金融机构款项
  prec_metals?: number;   // 贵金属
  deriv_assets?: number;  // 衍生金融资产
  rr_reinsur_une_prem?: number; // 应收分保未到期责任准备金
  rr_reinsur_outstd_cla?: number; // 应收分保未决赔款准备金
  rr_reinsur_lins_liab?: number; // 应收分保寿险责任准备金
  rr_reinsur_lthins_liab?: number; // 应收分保长期健康险责任准备金
  refund_depos?: number;  // 存出保证金
  ph_pledge_loans?: number; // 保户质押贷款
  refund_cap_depos?: number; // 存出资本保证金
  indep_acct_assets?: number; // 独立账户资产
  client_depos?: number;  // 其中：客户资金存款
  client_prov?: number;   // 其中：客户备付金
  transac_seat_fee?: number; // 其中:交易席位费
  invest_as_receiv?: number; // 应收款项类投资
  total_assets: number;   // 资产总计
  lt_borr: number;        // 长期借款
  st_borr: number;        // 短期借款
  cb_borr?: number;       // 向中央银行借款
  depos_ib_deposits?: number; // 吸收存款及同业存放
  loan_oth_bank?: number; // 拆入资金
  trading_fl: number;     // 交易性金融负债
  notes_payable: number;  // 应付票据
  acct_payable: number;   // 应付账款
  adv_receipts: number;   // 预收款项
  sold_for_repur_fa?: number; // 卖出回购金融资产款
  comm_payable?: number;  // 应付手续费及佣金
  payroll_payable: number; // 应付职工薪酬
  taxes_payable: number;  // 应交税费
  int_payable?: number;   // 应付利息
  div_payable?: number;   // 应付股利
  oth_payable: number;    // 其他应付款
  acc_exp?: number;       // 预提费用
  deferred_inc?: number;  // 递延收益
  st_bonds_payable?: number; // 应付短期债券
  payable_to_reinsurer?: number; // 应付分保账款
  rsrv_insur_cont?: number; // 保险合同准备金
  acting_trading_sec?: number; // 代理买卖证券款
  acting_uw_sec?: number; // 代理承销证券款
  non_cur_liab_due_1y: number; // 一年内到期的非流动负债
  oth_cur_liab: number;   // 其他流动负债
  total_cur_liab: number; // 流动负债合计
  bond_payable: number;   // 应付债券
  lt_payable: number;     // 长期应付款
  specific_payables?: number; // 专项应付款
  estimated_liab?: number; // 预计负债
  defer_tax_liab: number; // 递延所得税负债
  defer_inc_non_cur_liab?: number; // 递延收益-非流动负债
  oth_ncl: number;        // 其他非流动负债
  total_ncl: number;      // 非流动负债合计
  depos_oth_bfi?: number; // 同业和其他金融机构存放款项
  deriv_liab?: number;    // 衍生金融负债
  depos?: number;         // 吸收存款
  agency_bus_liab?: number; // 代理业务负债
  oth_liab?: number;      // 其他负债
  prem_receiv_adva?: number; // 预收保费
  depos_received?: number; // 存入保证金
  ph_invest?: number;     // 保户储金及投资款
  reser_une_prem?: number; // 未到期责任准备金
  reser_outstd_claims?: number; // 未决赔款准备金
  reser_lins_liab?: number; // 寿险责任准备金
  reser_lthins_liab?: number; // 长期健康险责任准备金
  indept_acc_liab?: number; // 独立账户负债
  pledge_borr?: number;   // 其中:质押借款
  indem_payable?: number; // 应付赔付款
  policy_div_payable?: number; // 应付保单红利
  total_liab: number;     // 负债合计
  treasury_share?: number; // 减:库存股
  ordin_risk_reser?: number; // 一般风险准备
  forex_differ?: number;  // 外币报表折算差额
  invest_loss_unconf?: number; // 未确认的投资损失
  minority_int?: number;  // 少数股东权益
  total_hldr_eqy_exc_min_int: number; // 股东权益合计(不含少数股东权益)
  total_hldr_eqy_inc_min_int: number; // 股东权益合计(含少数股东权益)
  total_liab_hldr_eqy: number; // 负债及股东权益总计
  lt_payroll_payable?: number; // 长期应付职工薪酬
  oth_comp_income?: number; // 其他综合收益
  oth_eqt_tools?: number; // 其他权益工具
  oth_eqt_tools_p_shr?: number; // 其他权益工具:优先股
  lending_funds?: number; // 融出资金
  acc_receivable?: number; // 应收款项
  st_fin_payable?: number; // 应付短期融资款
  payables?: number;      // 应付款项
}

// 现金流量表
export interface CashFlowStatement {
  ts_code: string;
  ann_date: string;       // 公告日期
  end_date: string;       // 报告期
  report_type: string;    // 报告类型
  comp_type: string;      // 公司类型
  net_profit: number;     // 净利润
  finan_exp?: number;     // 财务费用
  c_fr_sale_sg: number;   // 销售商品、提供劳务收到的现金
  recp_tax_rends?: number; // 收到的税费返还
  n_depos_incr_fi?: number; // 客户存款和同业存放款项净增加额
  n_incr_loans_cb?: number; // 向中央银行借款净增加额
  n_inc_borr_oth_fi?: number; // 向其他金融机构拆入资金净增加额
  prem_fr_orig_contr?: number; // 收到原保险合同保费取得的现金
  n_incr_insured_dep?: number; // 保户储金净增加额
  n_reinsur_prem?: number; // 收到再保业务现金净额
  n_incr_disp_tfa?: number; // 处置交易性金融资产净增加额
  ifc_cash_incr?: number; // 收取利息和手续费净增加额
  n_incr_disp_faas?: number; // 处置可供出售金融资产净增加额
  n_incr_loans_oth_bank?: number; // 拆入资金净增加额
  n_cap_incr_repur?: number; // 回购业务资金净增加额
  c_fr_oth_operate_a: number; // 收到其他与经营活动有关的现金
  c_inf_fr_operate_a: number; // 经营活动现金流入小计
  c_paid_goods_s: number; // 购买商品、接受劳务支付的现金
  c_paid_to_for_empl: number; // 支付给职工以及为职工支付的现金
  c_paid_for_taxes: number; // 支付的各项税费
  n_incr_clt_loan_adv?: number; // 客户贷款及垫款净增加额
  n_incr_dep_cbob?: number; // 存放央行和同业款项净增加额
  c_pay_claims_orig_inco?: number; // 支付原保险合同赔付款项的现金
  pay_handling_chrg?: number; // 支付手续费的现金
  pay_comm_insur_plcy?: number; // 支付保单红利的现金
  c_paid_oth_operate_a: number; // 支付其他与经营活动有关的现金
  c_outf_fr_operate_a: number; // 经营活动现金流出小计
  n_cashflow_operate_a: number; // 经营活动产生的现金流量净额
  c_recp_disp_withdrwl_invest?: number; // 收回投资收到的现金
  c_recp_return_invest: number; // 取得投资收益收到的现金
  n_recp_disp_fiolta?: number; // 处置固定资产、无形资产和其他长期资产收回的现金净额
  n_recp_disp_sobu?: number; // 处置子公司及其他营业单位收到的现金净额
  c_recp_oth_invest_a: number; // 收到其他与投资活动有关的现金
  c_inf_fr_invest_a: number; // 投资活动现金流入小计
  c_paid_acq_const_fiolta: number; // 购建固定资产、无形资产和其他长期资产支付的现金
  c_paid_invest: number;  // 投资支付的现金
  n_incr_pledge_loan?: number; // 质押贷款净增加额
  n_paid_acq_sobu?: number; // 取得子公司及其他营业单位支付的现金净额
  c_paid_oth_invest_a: number; // 支付其他与投资活动有关的现金
  c_outf_fr_invest_a: number; // 投资活动现金流出小计
  n_cashflow_invest_a: number; // 投资活动产生的现金流量净额
  c_recp_borrow: number;  // 取得借款收到的现金
  proc_issue_bonds?: number; // 发行债券收到的现金
  c_recp_oth_fin_a: number; // 收到其他与筹资活动有关的现金
  c_inf_fr_fin_a: number; // 筹资活动现金流入小计
  c_prepay_amt_borr: number; // 偿还债务支付的现金
  c_pay_dist_dpcp_int_exp: number; // 分配股利、利润或偿付利息支付的现金
  incl_dvd_profit_paid_sc_ms?: number; // 其中:子公司支付给少数股东的股利、利润
  c_paid_oth_fin_a: number; // 支付其他与筹资活动有关的现金
  c_outf_fr_fin_a: number; // 筹资活动现金流出小计
  n_cash_flows_fnc_act: number; // 筹资活动产生的现金流量净额
  eff_fx_flu_cash?: number; // 汇率变动对现金的影响
  n_incr_cash_cash_equ: number; // 现金及现金等价物净增加额
  c_cash_equ_beg_period: number; // 期初现金及现金等价物余额
  c_cash_equ_end_period: number; // 期末现金及现金等价物余额
  c_recp_cap_contrib?: number; // 吸收投资收到的现金
  incl_cash_rec_saims?: number; // 其中:子公司吸收少数股东投资收到的现金
  uncon_invest_loss?: number; // 未确认投资损失
  prov_depr_assets?: number; // 加:资产减值准备
  depr_fa_coga_dpba?: number; // 固定资产折旧、油气资产折耗、生产性生物资产折旧
  amort_intang_assets?: number; // 无形资产摊销
  lt_amort_deferred_exp?: number; // 长期待摊费用摊销
  decr_deferred_exp?: number; // 待摊费用减少
  incr_acc_exp?: number;  // 预提费用增加
  loss_disp_fiolta?: number; // 处置固定、无形资产和其他长期资产的损失
  loss_scr_fa?: number;   // 固定资产报废损失
  loss_fv_chg?: number;   // 公允价值变动损失
  invest_loss?: number;   // 投资损失
  decr_def_inc_tax_assets?: number; // 递延所得税资产减少
  incr_def_inc_tax_liab?: number; // 递延所得税负债增加
  decr_inventories?: number; // 存货的减少
  decr_oper_payable?: number; // 经营性应收项目的减少
  incr_oper_payable?: number; // 经营性应付项目的增加
  others?: number;        // 其他
  im_net_cashflow_oper_act?: number; // 经营活动产生的现金流量净额(间接法)
  conv_debt_into_cap?: number; // 债务转为资本
  conv_copbonds_due_within_1y?: number; // 一年内到期的可转换公司债券
  fa_fnc_leases?: number; // 融资租入固定资产
  end_bal_cash?: number;  // 现金的期末余额
  beg_bal_cash?: number;  // 减:现金的期初余额
  end_bal_cash_equ?: number; // 加:现金等价物的期末余额
  beg_bal_cash_equ?: number; // 减:现金等价物的期初余额
  im_n_incr_cash_equ?: number; // 现金及现金等价物净增加额(间接法)
}

// 每日基本面数据
export interface DailyBasic {
  ts_code: string;
  trade_date: string;
  close?: number;         // 当日收盘价
  turnover_rate?: number; // 换手率
  turnover_rate_f?: number; // 换手率(自由流通股)
  volume_ratio?: number;  // 量比
  pe?: number;           // 市盈率(总市值/净利润, 亏损的PE为空)
  pe_ttm?: number;       // 市盈率(TTM, 亏损的PE为空)
  pb?: number;           // 市净率(总市值/净资产)
  ps?: number;           // 市销率
  ps_ttm?: number;       // 市销率(TTM)
  dv_ratio?: number;     // 股息率(%)
  dv_ttm?: number;       // 股息率(TTM)(%)
  total_share?: number;  // 总股本(万股)
  float_share?: number;  // 流通股本(万股)
  free_share?: number;   // 自由流通股本(万)
  total_mv?: number;     // 总市值(万元)
  circ_mv?: number;      // 流通市值(万元)
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

// 预测数据
export interface PredictionData {
  ts_code: string;
  trade_date: string;
  predicted_direction?: string; // 预测方向: up/down/hold
  confidence?: number;          // 置信度 0-1
  buy_signal?: boolean;         // 买入信号
  sell_signal?: boolean;        // 卖出信号
  target_price?: number;        // 目标价格
  stop_loss?: number;          // 止损价格
  reasons?: string[];          // 预测理由
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
  x: number | string;     // X轴数据(通常是日期索引或日期字符串)
  y: number;              // Y轴数据
  date?: string;          // 日期字符串
  open?: number;          // 开盘价(K线用)
  high?: number;          // 最高价(K线用)
  low?: number;           // 最低价(K线用)
  close?: number;         // 收盘价(K线用)
  volume?: number;        // 成交量
}

// K线图数据
export interface CandlestickData extends ChartDataPoint {
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
}

// 技术指标图表数据
export interface IndicatorChartData {
  dates: string[];
  ma_lines?: {
    ma5?: number[];
    ma10?: number[];
    ma20?: number[];
    ma30?: number[];
    ma60?: number[];
  };
  macd?: {
    dif: number[];
    dea: number[];
    histogram: number[];
  };
  rsi?: {
    rsi6?: number[];
    rsi12?: number[];
    rsi24?: number[];
  };
  boll?: {
    upper: number[];
    mid: number[];
    lower: number[];
  };
  kdj?: {
    k: number[];
    d: number[];
    j: number[];
  };
}
