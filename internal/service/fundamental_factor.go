package service

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"stock-a-future/internal/client"
	"stock-a-future/internal/models"

	"github.com/shopspring/decimal"
)

// FundamentalFactorService 基本面因子计算服务
type FundamentalFactorService struct {
	client       client.DataSourceClient
	standardizer *FactorStandardizer
}

// NewFundamentalFactorService 创建基本面因子计算服务
func NewFundamentalFactorService(client client.DataSourceClient) *FundamentalFactorService {
	return &FundamentalFactorService{
		client:       client,
		standardizer: NewFactorStandardizer(),
	}
}

// CalculateFundamentalFactor 计算单个股票的基本面因子
func (s *FundamentalFactorService) CalculateFundamentalFactor(symbol, tradeDate string) (*models.FundamentalFactor, error) {
	startTime := time.Now()
	log.Printf("[FundamentalFactorService] 开始计算基本面因子: %s, 日期: %s", symbol, tradeDate)

	// 1. 获取每日基本面数据
	ctx := context.Background()
	dailyBasic, err := s.client.GetDailyBasic(ctx, symbol, tradeDate)
	if err != nil {
		log.Printf("[FundamentalFactorService] 获取每日基本面数据失败: %v", err)
		return nil, fmt.Errorf("获取每日基本面数据失败: %v", err)
	}

	// 2. 获取最新财务报表数据
	incomeStatement, err := s.getLatestIncomeStatement(symbol)
	if err != nil {
		log.Printf("[FundamentalFactorService] 获取利润表数据失败: %v", err)
		// 继续执行，使用空数据
	}

	balanceSheet, err := s.getLatestBalanceSheet(symbol)
	if err != nil {
		log.Printf("[FundamentalFactorService] 获取资产负债表数据失败: %v", err)
		// 继续执行，使用空数据
	}

	cashFlow, err := s.getLatestCashFlow(symbol)
	if err != nil {
		log.Printf("[FundamentalFactorService] 获取现金流量表数据失败: %v", err)
		// 继续执行，使用空数据
	}

	// 3. 计算基本面因子
	factor := &models.FundamentalFactor{
		TSCode:    s.determineTSCode(symbol),
		TradeDate: tradeDate,
		UpdatedAt: time.Now(),
	}

	// 计算价值因子
	s.calculateValueFactors(factor, dailyBasic)

	// 计算成长因子
	s.calculateGrowthFactors(factor, incomeStatement)

	// 计算质量因子
	s.calculateQualityFactors(factor, incomeStatement, balanceSheet)

	// 计算盈利因子
	s.calculateProfitabilityFactors(factor, incomeStatement, balanceSheet)

	// 计算运营效率因子
	s.calculateOperatingFactors(factor, incomeStatement, balanceSheet)

	// 计算现金流因子
	s.calculateCashFlowFactors(factor, incomeStatement, cashFlow)

	// 计算分红因子
	s.calculateDividendFactors(factor, dailyBasic)

	// 计算单个股票的因子得分（基于简化逻辑）
	s.calculateSingleStockScores(factor)

	totalDuration := time.Since(startTime)
	log.Printf("[FundamentalFactorService] 基本面因子计算完成: %s, 总耗时: %v", symbol, totalDuration)
	return factor, nil
}

// BatchCalculateFundamentalFactors 批量计算基本面因子并标准化
func (s *FundamentalFactorService) BatchCalculateFundamentalFactors(symbols []string, tradeDate string) ([]models.FundamentalFactor, error) {
	log.Printf("[FundamentalFactorService] 开始批量计算基本面因子，股票数量: %d", len(symbols))

	var factors []models.FundamentalFactor

	for i, symbol := range symbols {
		log.Printf("[FundamentalFactorService] 计算进度: %d/%d - %s", i+1, len(symbols), symbol)

		factor, err := s.CalculateFundamentalFactor(symbol, tradeDate)
		if err != nil {
			log.Printf("[FundamentalFactorService] 计算失败: %s - %v", symbol, err)
			continue
		}
		factors = append(factors, *factor)
	}

	// 对所有因子进行标准化处理
	if len(factors) > 1 {
		log.Printf("[FundamentalFactorService] 开始标准化处理，因子数量: %d", len(factors))
		standardizedFactors, err := s.standardizer.StandardizeFundamentalFactors(factors)
		if err != nil {
			log.Printf("[FundamentalFactorService] 标准化处理失败: %v", err)
			return factors, nil // 返回未标准化的数据
		}
		factors = standardizedFactors
	}

	log.Printf("[FundamentalFactorService] 批量计算完成，成功处理: %d 个股票", len(factors))
	return factors, nil
}

// GetFactorRanking 获取因子排名
func (s *FundamentalFactorService) GetFactorRanking(factors []models.FundamentalFactor, factorType string) []models.FundamentalFactor {
	// 复制切片避免修改原数据
	rankedFactors := make([]models.FundamentalFactor, len(factors))
	copy(rankedFactors, factors)

	// 根据因子类型排序
	switch factorType {
	case "value":
		sort.Slice(rankedFactors, func(i, j int) bool {
			return rankedFactors[i].ValueScore.Decimal.GreaterThan(rankedFactors[j].ValueScore.Decimal)
		})
	case "growth":
		sort.Slice(rankedFactors, func(i, j int) bool {
			return rankedFactors[i].GrowthScore.Decimal.GreaterThan(rankedFactors[j].GrowthScore.Decimal)
		})
	case "quality":
		sort.Slice(rankedFactors, func(i, j int) bool {
			return rankedFactors[i].QualityScore.Decimal.GreaterThan(rankedFactors[j].QualityScore.Decimal)
		})
	case "profitability":
		sort.Slice(rankedFactors, func(i, j int) bool {
			return rankedFactors[i].ProfitabilityScore.Decimal.GreaterThan(rankedFactors[j].ProfitabilityScore.Decimal)
		})
	case "composite":
		sort.Slice(rankedFactors, func(i, j int) bool {
			return rankedFactors[i].CompositeScore.Decimal.GreaterThan(rankedFactors[j].CompositeScore.Decimal)
		})
	}

	return rankedFactors
}

// 私有方法：获取最新财务报表数据

func (s *FundamentalFactorService) getLatestIncomeStatement(symbol string) (*models.IncomeStatement, error) {
	// 尝试获取最近几个报告期的数据
	periods := []string{"20241231", "20240930", "20240630", "20240331", "20231231"}

	log.Printf("[FundamentalFactorService] 开始获取利润表数据: %s", symbol)

	for i, period := range periods {
		log.Printf("[FundamentalFactorService] 尝试获取利润表数据 - 期间 %d/%d: %s", i+1, len(periods), period)
		statement, err := s.client.GetIncomeStatement(symbol, period, "1")
		if err != nil {
			log.Printf("[FundamentalFactorService] 期间 %s 获取失败: %v", period, err)
			continue
		}
		if statement == nil {
			log.Printf("[FundamentalFactorService] 期间 %s 返回空数据", period)
			continue
		}

		// 检查关键字段是否有效
		if statement.OperRevenue.Decimal.IsZero() && statement.NetProfit.Decimal.IsZero() {
			log.Printf("[FundamentalFactorService] 期间 %s 数据无效 - 营业收入和净利润都为零", period)
			continue
		}

		log.Printf("[FundamentalFactorService] 成功获取利润表数据 - 期间: %s, 营业收入: %v, 净利润: %v",
			period, statement.OperRevenue.Decimal, statement.NetProfit.Decimal)
		return statement, nil
	}

	// 尝试不指定期间，获取最新数据
	log.Printf("[FundamentalFactorService] 所有指定期间都失败，尝试获取最新利润表数据")
	statement, err := s.client.GetIncomeStatement(symbol, "", "1")
	if err != nil {
		log.Printf("[FundamentalFactorService] 获取最新利润表数据失败: %v", err)
		return nil, fmt.Errorf("未找到有效的利润表数据: 所有期间都失败，最新数据获取错误: %v", err)
	}
	if statement == nil {
		log.Printf("[FundamentalFactorService] 最新利润表数据为空")
		return nil, fmt.Errorf("未找到有效的利润表数据: 最新数据为空")
	}

	// 检查最新数据的关键字段
	if statement.OperRevenue.Decimal.IsZero() && statement.NetProfit.Decimal.IsZero() {
		log.Printf("[FundamentalFactorService] 最新利润表数据无效 - 营业收入和净利润都为零")
		return nil, fmt.Errorf("未找到有效的利润表数据: 最新数据营业收入和净利润都为零")
	}

	log.Printf("[FundamentalFactorService] 成功获取最新利润表数据 - 营业收入: %v, 净利润: %v",
		statement.OperRevenue.Decimal, statement.NetProfit.Decimal)
	return statement, nil
}

func (s *FundamentalFactorService) getLatestBalanceSheet(symbol string) (*models.BalanceSheet, error) {
	// 尝试获取最近几个报告期的数据
	periods := []string{"20241231", "20240930", "20240630", "20240331", "20231231"}

	log.Printf("[FundamentalFactorService] 开始获取资产负债表数据: %s", symbol)

	for i, period := range periods {
		log.Printf("[FundamentalFactorService] 尝试获取资产负债表数据 - 期间 %d/%d: %s", i+1, len(periods), period)
		sheet, err := s.client.GetBalanceSheet(symbol, period, "1")
		if err != nil {
			log.Printf("[FundamentalFactorService] 期间 %s 获取失败: %v", period, err)
			continue
		}
		if sheet == nil {
			log.Printf("[FundamentalFactorService] 期间 %s 返回空数据", period)
			continue
		}

		// 检查关键字段是否有效
		if sheet.TotalAssets.Decimal.IsZero() && sheet.TotalHldrEqy.Decimal.IsZero() {
			log.Printf("[FundamentalFactorService] 期间 %s 数据无效 - 总资产和所有者权益都为零", period)
			continue
		}

		log.Printf("[FundamentalFactorService] 成功获取资产负债表数据 - 期间: %s, 总资产: %v, 所有者权益: %v",
			period, sheet.TotalAssets.Decimal, sheet.TotalHldrEqy.Decimal)
		return sheet, nil
	}

	// 尝试不指定期间，获取最新数据
	log.Printf("[FundamentalFactorService] 所有指定期间都失败，尝试获取最新资产负债表数据")
	sheet, err := s.client.GetBalanceSheet(symbol, "", "1")
	if err != nil {
		log.Printf("[FundamentalFactorService] 获取最新资产负债表数据失败: %v", err)
		return nil, fmt.Errorf("未找到有效的资产负债表数据: 所有期间都失败，最新数据获取错误: %v", err)
	}
	if sheet == nil {
		log.Printf("[FundamentalFactorService] 最新资产负债表数据为空")
		return nil, fmt.Errorf("未找到有效的资产负债表数据: 最新数据为空")
	}

	// 检查最新数据的关键字段
	if sheet.TotalAssets.Decimal.IsZero() && sheet.TotalHldrEqy.Decimal.IsZero() {
		log.Printf("[FundamentalFactorService] 最新资产负债表数据无效 - 总资产和所有者权益都为零")
		return nil, fmt.Errorf("未找到有效的资产负债表数据: 最新数据总资产和所有者权益都为零")
	}

	log.Printf("[FundamentalFactorService] 成功获取最新资产负债表数据 - 总资产: %v, 所有者权益: %v",
		sheet.TotalAssets.Decimal, sheet.TotalHldrEqy.Decimal)
	return sheet, nil
}

func (s *FundamentalFactorService) getLatestCashFlow(symbol string) (*models.CashFlowStatement, error) {
	// 尝试获取最近几个报告期的数据
	periods := []string{"20241231", "20240930", "20240630", "20240331", "20231231"}

	log.Printf("[FundamentalFactorService] 开始获取现金流量表数据: %s", symbol)

	for i, period := range periods {
		log.Printf("[FundamentalFactorService] 尝试获取现金流量表数据 - 期间 %d/%d: %s", i+1, len(periods), period)
		cashFlow, err := s.client.GetCashFlowStatement(symbol, period, "1")
		if err != nil {
			log.Printf("[FundamentalFactorService] 期间 %s 获取失败: %v", period, err)
			continue
		}
		if cashFlow == nil {
			log.Printf("[FundamentalFactorService] 期间 %s 返回空数据", period)
			continue
		}

		// 检查关键字段是否有效
		if cashFlow.NetCashOperAct.Decimal.IsZero() &&
			cashFlow.NetCashInvAct.Decimal.IsZero() &&
			cashFlow.NetCashFinAct.Decimal.IsZero() {
			log.Printf("[FundamentalFactorService] 期间 %s 数据无效 - 所有现金流净额为零", period)
			continue
		}

		log.Printf("[FundamentalFactorService] 成功获取现金流量表数据 - 期间: %s, 经营现金流: %v, 投资现金流: %v, 筹资现金流: %v",
			period, cashFlow.NetCashOperAct.Decimal, cashFlow.NetCashInvAct.Decimal, cashFlow.NetCashFinAct.Decimal)
		return cashFlow, nil
	}

	// 尝试不指定期间，获取最新数据
	log.Printf("[FundamentalFactorService] 所有指定期间都失败，尝试获取最新现金流量表数据")
	cashFlow, err := s.client.GetCashFlowStatement(symbol, "", "1")
	if err != nil {
		log.Printf("[FundamentalFactorService] 获取最新现金流量表数据失败: %v", err)
		return nil, fmt.Errorf("未找到有效的现金流量表数据: 所有期间都失败，最新数据获取错误: %v", err)
	}
	if cashFlow == nil {
		log.Printf("[FundamentalFactorService] 最新现金流量表数据为空")
		return nil, fmt.Errorf("未找到有效的现金流量表数据: 最新数据为空")
	}

	// 检查最新数据的关键字段
	if cashFlow.NetCashOperAct.Decimal.IsZero() &&
		cashFlow.NetCashInvAct.Decimal.IsZero() &&
		cashFlow.NetCashFinAct.Decimal.IsZero() {
		log.Printf("[FundamentalFactorService] 最新现金流量表数据无效 - 所有现金流净额为零")
		return nil, fmt.Errorf("未找到有效的现金流量表数据: 最新数据所有现金流净额为零")
	}

	log.Printf("[FundamentalFactorService] 成功获取最新现金流量表数据 - 经营现金流: %v, 投资现金流: %v, 筹资现金流: %v",
		cashFlow.NetCashOperAct.Decimal, cashFlow.NetCashInvAct.Decimal, cashFlow.NetCashFinAct.Decimal)
	return cashFlow, nil
}

// 私有方法：计算各类因子

func (s *FundamentalFactorService) calculateValueFactors(factor *models.FundamentalFactor, dailyBasic *models.DailyBasic) {
	if dailyBasic == nil {
		return
	}

	// PE 市盈率
	factor.PE = dailyBasic.Pe

	// PB 市净率
	factor.PB = dailyBasic.Pb

	// PS 市销率
	factor.PS = dailyBasic.Ps

	// PCF 市现率 - 暂时使用 PE 作为替代
	factor.PCF = dailyBasic.Pe

	// EV/EBITDA - 暂时设为0，需要更多数据计算
	factor.EVToEBITDA = models.NewJSONDecimal(decimal.Zero)
}

func (s *FundamentalFactorService) calculateGrowthFactors(factor *models.FundamentalFactor, incomeStatement *models.IncomeStatement) {
	if incomeStatement == nil {
		// 设置默认值
		factor.RevenueGrowth = models.NewJSONDecimal(decimal.Zero)
		factor.NetProfitGrowth = models.NewJSONDecimal(decimal.Zero)
		factor.EPSGrowth = models.NewJSONDecimal(decimal.Zero)
		factor.ROEGrowth = models.NewJSONDecimal(decimal.Zero)
		return
	}

	// 尝试获取历史数据计算真实的增长率
	historicalData := s.getHistoricalIncomeStatement(factor.TSCode, incomeStatement.EndDate)

	if historicalData != nil {
		// 计算营收增长率
		factor.RevenueGrowth = s.calculateGrowthRate(incomeStatement.OperRevenue.Decimal, historicalData.OperRevenue.Decimal)

		// 计算净利润增长率
		factor.NetProfitGrowth = s.calculateGrowthRate(incomeStatement.NetProfit.Decimal, historicalData.NetProfit.Decimal)

		// EPS增长率和ROE增长率需要更多数据，暂时使用营收增长率作为估算
		factor.EPSGrowth = factor.NetProfitGrowth
		factor.ROEGrowth = factor.NetProfitGrowth

		log.Printf("[FundamentalFactorService] 成功计算成长因子 - 营收增长率: %v%%, 净利润增长率: %v%%",
			factor.RevenueGrowth.Decimal.Mul(decimal.NewFromInt(100)),
			factor.NetProfitGrowth.Decimal.Mul(decimal.NewFromInt(100)))
	} else {
		// 无法获取历史数据，使用当前数据与行业平均值比较（简化处理）
		factor.RevenueGrowth = s.estimateGrowthFromCurrentData(incomeStatement.OperRevenue.Decimal)
		factor.NetProfitGrowth = s.estimateGrowthFromCurrentData(incomeStatement.NetProfit.Decimal)
		factor.EPSGrowth = factor.NetProfitGrowth
		factor.ROEGrowth = factor.NetProfitGrowth

		log.Printf("[FundamentalFactorService] 使用估算方法计算成长因子 - 营收增长率: %v%%, 净利润增长率: %v%%",
			factor.RevenueGrowth.Decimal.Mul(decimal.NewFromInt(100)),
			factor.NetProfitGrowth.Decimal.Mul(decimal.NewFromInt(100)))
	}
}

func (s *FundamentalFactorService) calculateQualityFactors(factor *models.FundamentalFactor, incomeStatement *models.IncomeStatement, balanceSheet *models.BalanceSheet) {
	if incomeStatement == nil || balanceSheet == nil {
		// 设置默认值
		factor.ROE = models.NewJSONDecimal(decimal.Zero)
		factor.ROA = models.NewJSONDecimal(decimal.Zero)
		factor.DebtToAssets = models.NewJSONDecimal(decimal.Zero)
		factor.CurrentRatio = models.NewJSONDecimal(decimal.Zero)
		factor.QuickRatio = models.NewJSONDecimal(decimal.Zero)
		return
	}

	// ROE = 净利润 / 所有者权益
	if !balanceSheet.TotalHldrEqy.Decimal.IsZero() {
		roe := incomeStatement.NetProfit.Decimal.Div(balanceSheet.TotalHldrEqy.Decimal).Mul(decimal.NewFromInt(100))
		factor.ROE = models.NewJSONDecimal(roe)
	} else {
		factor.ROE = models.NewJSONDecimal(decimal.Zero)
	}

	// ROA = 净利润 / 总资产
	if !balanceSheet.TotalAssets.Decimal.IsZero() {
		roa := incomeStatement.NetProfit.Decimal.Div(balanceSheet.TotalAssets.Decimal).Mul(decimal.NewFromInt(100))
		factor.ROA = models.NewJSONDecimal(roa)
	} else {
		factor.ROA = models.NewJSONDecimal(decimal.Zero)
	}

	// 资产负债率 = 总负债 / 总资产
	if !balanceSheet.TotalAssets.Decimal.IsZero() {
		debtRatio := balanceSheet.TotalLiab.Decimal.Div(balanceSheet.TotalAssets.Decimal).Mul(decimal.NewFromInt(100))
		factor.DebtToAssets = models.NewJSONDecimal(debtRatio)
	} else {
		factor.DebtToAssets = models.NewJSONDecimal(decimal.Zero)
	}

	// 流动比率 = 流动资产 / 流动负债
	if !balanceSheet.TotalCurLiab.Decimal.IsZero() {
		currentRatio := balanceSheet.TotalCurAssets.Decimal.Div(balanceSheet.TotalCurLiab.Decimal)
		factor.CurrentRatio = models.NewJSONDecimal(currentRatio)
	} else {
		factor.CurrentRatio = models.NewJSONDecimal(decimal.Zero)
	}

	// 速动比率 = (流动资产 - 存货) / 流动负债
	if !balanceSheet.TotalCurLiab.Decimal.IsZero() {
		quickAssets := balanceSheet.TotalCurAssets.Decimal.Sub(balanceSheet.InventoryAssets.Decimal)
		quickRatio := quickAssets.Div(balanceSheet.TotalCurLiab.Decimal)
		factor.QuickRatio = models.NewJSONDecimal(quickRatio)
	} else {
		factor.QuickRatio = models.NewJSONDecimal(decimal.Zero)
	}
}

func (s *FundamentalFactorService) calculateProfitabilityFactors(factor *models.FundamentalFactor, incomeStatement *models.IncomeStatement, balanceSheet *models.BalanceSheet) {
	if incomeStatement == nil {
		// 设置默认值
		factor.GrossMargin = models.NewJSONDecimal(decimal.Zero)
		factor.NetMargin = models.NewJSONDecimal(decimal.Zero)
		factor.OperatingMargin = models.NewJSONDecimal(decimal.Zero)
		factor.EBITDAMargin = models.NewJSONDecimal(decimal.Zero)
		factor.ROIC = models.NewJSONDecimal(decimal.Zero)
		return
	}

	// 毛利率 = (营业收入 - 营业成本) / 营业收入
	if !incomeStatement.OperRevenue.Decimal.IsZero() {
		grossProfit := incomeStatement.OperRevenue.Decimal.Sub(incomeStatement.OperCost.Decimal)
		grossMargin := grossProfit.Div(incomeStatement.OperRevenue.Decimal).Mul(decimal.NewFromInt(100))
		factor.GrossMargin = models.NewJSONDecimal(grossMargin)
	} else {
		factor.GrossMargin = models.NewJSONDecimal(decimal.Zero)
	}

	// 净利率 = 净利润 / 营业收入
	if !incomeStatement.OperRevenue.Decimal.IsZero() {
		netMargin := incomeStatement.NetProfit.Decimal.Div(incomeStatement.OperRevenue.Decimal).Mul(decimal.NewFromInt(100))
		factor.NetMargin = models.NewJSONDecimal(netMargin)
	} else {
		factor.NetMargin = models.NewJSONDecimal(decimal.Zero)
	}

	// 营业利润率 = 营业利润 / 营业收入
	if !incomeStatement.OperRevenue.Decimal.IsZero() {
		operMargin := incomeStatement.OperProfit.Decimal.Div(incomeStatement.OperRevenue.Decimal).Mul(decimal.NewFromInt(100))
		factor.OperatingMargin = models.NewJSONDecimal(operMargin)
	} else {
		factor.OperatingMargin = models.NewJSONDecimal(decimal.Zero)
	}

	// EBITDA利润率 - 暂时使用营业利润率
	factor.EBITDAMargin = factor.OperatingMargin

	// ROIC - 投入资本回报率，暂时使用ROE
	factor.ROIC = factor.ROE
}

func (s *FundamentalFactorService) calculateOperatingFactors(factor *models.FundamentalFactor, incomeStatement *models.IncomeStatement, balanceSheet *models.BalanceSheet) {
	if incomeStatement == nil || balanceSheet == nil {
		// 设置默认值
		factor.AssetTurnover = models.NewJSONDecimal(decimal.Zero)
		factor.InventoryTurnover = models.NewJSONDecimal(decimal.Zero)
		factor.ReceivableTurnover = models.NewJSONDecimal(decimal.Zero)
		return
	}

	// 总资产周转率 = 营业收入 / 平均总资产
	if !balanceSheet.TotalAssets.Decimal.IsZero() {
		assetTurnover := incomeStatement.OperRevenue.Decimal.Div(balanceSheet.TotalAssets.Decimal)
		factor.AssetTurnover = models.NewJSONDecimal(assetTurnover)
	} else {
		factor.AssetTurnover = models.NewJSONDecimal(decimal.Zero)
	}

	// 存货周转率 = 营业成本 / 平均存货
	if !balanceSheet.InventoryAssets.Decimal.IsZero() {
		invTurnover := incomeStatement.OperCost.Decimal.Div(balanceSheet.InventoryAssets.Decimal)
		factor.InventoryTurnover = models.NewJSONDecimal(invTurnover)
	} else {
		factor.InventoryTurnover = models.NewJSONDecimal(decimal.Zero)
	}

	// 应收账款周转率 = 营业收入 / 平均应收账款
	if !balanceSheet.AccountsReceiv.Decimal.IsZero() {
		arTurnover := incomeStatement.OperRevenue.Decimal.Div(balanceSheet.AccountsReceiv.Decimal)
		factor.ReceivableTurnover = models.NewJSONDecimal(arTurnover)
	} else {
		factor.ReceivableTurnover = models.NewJSONDecimal(decimal.Zero)
	}
}

func (s *FundamentalFactorService) calculateCashFlowFactors(factor *models.FundamentalFactor, incomeStatement *models.IncomeStatement, cashFlow *models.CashFlowStatement) {
	if incomeStatement == nil || cashFlow == nil {
		// 设置默认值
		factor.OCFToRevenue = models.NewJSONDecimal(decimal.Zero)
		factor.OCFToNetProfit = models.NewJSONDecimal(decimal.Zero)
		factor.FCFYield = models.NewJSONDecimal(decimal.Zero)
		return
	}

	// 经营现金流/营收比 = 经营活动现金流净额 / 营业收入
	if !incomeStatement.OperRevenue.Decimal.IsZero() {
		ocfRevRatio := cashFlow.NetCashOperAct.Decimal.Div(incomeStatement.OperRevenue.Decimal).Mul(decimal.NewFromInt(100))
		factor.OCFToRevenue = models.NewJSONDecimal(ocfRevRatio)
	} else {
		factor.OCFToRevenue = models.NewJSONDecimal(decimal.Zero)
	}

	// 经营现金流/净利润比 = 经营活动现金流净额 / 净利润
	if !incomeStatement.NetProfit.Decimal.IsZero() {
		ocfNetRatio := cashFlow.NetCashOperAct.Decimal.Div(incomeStatement.NetProfit.Decimal)
		factor.OCFToNetProfit = models.NewJSONDecimal(ocfNetRatio)
	} else {
		factor.OCFToNetProfit = models.NewJSONDecimal(decimal.Zero)
	}

	// 自由现金流收益率 - 暂时使用经营现金流/营收比
	factor.FCFYield = factor.OCFToRevenue
}

func (s *FundamentalFactorService) calculateDividendFactors(factor *models.FundamentalFactor, dailyBasic *models.DailyBasic) {
	if dailyBasic == nil {
		// 设置默认值
		factor.DividendYield = models.NewJSONDecimal(decimal.Zero)
		factor.PayoutRatio = models.NewJSONDecimal(decimal.Zero)
		return
	}

	// 股息率
	factor.DividendYield = dailyBasic.DvRatio

	// 分红率 - 暂时设为0，需要更多数据计算
	factor.PayoutRatio = models.NewJSONDecimal(decimal.Zero)
}

// determineTSCode 智能判断股票代码的市场后缀
func (s *FundamentalFactorService) determineTSCode(symbol string) string {
	// 如果symbol已经包含市场后缀，直接返回
	if strings.Contains(symbol, ".") {
		return symbol
	}

	// 根据股票代码规则判断市场
	// 600xxx, 601xxx, 603xxx, 688xxx -> 上海
	// 000xxx, 002xxx, 300xxx -> 深圳
	// 430xxx, 830xxx, 870xxx -> 北京
	if strings.HasPrefix(symbol, "600") || strings.HasPrefix(symbol, "601") ||
		strings.HasPrefix(symbol, "603") || strings.HasPrefix(symbol, "688") {
		return symbol + ".SH"
	} else if strings.HasPrefix(symbol, "000") || strings.HasPrefix(symbol, "002") ||
		strings.HasPrefix(symbol, "300") {
		return symbol + ".SZ"
	} else if strings.HasPrefix(symbol, "430") || strings.HasPrefix(symbol, "830") ||
		strings.HasPrefix(symbol, "870") {
		return symbol + ".BJ"
	}

	// 默认返回上海市场
	return symbol + ".SH"
}

// calculateSingleStockScores 为单个股票计算因子得分（简化版本）
// 由于缺乏市场对比数据，使用基于阈值的评分逻辑
func (s *FundamentalFactorService) calculateSingleStockScores(factor *models.FundamentalFactor) {
	log.Printf("[FundamentalFactorService] 开始计算单股票因子得分: %s", factor.TSCode)

	// 计算价值因子得分
	factor.ValueScore = s.calculateValueScore(factor)

	// 计算成长因子得分
	factor.GrowthScore = s.calculateGrowthScore(factor)

	// 计算质量因子得分
	factor.QualityScore = s.calculateQualityScore(factor)

	// 计算盈利因子得分
	factor.ProfitabilityScore = s.calculateProfitabilityScore(factor)

	// 计算综合得分
	s.calculateCompositeScore(factor)

	// 设置默认排名和分位数（单股票无法计算真实排名）
	factor.MarketRank = 1
	factor.IndustryRank = 1
	factor.MarketPercentile = models.NewJSONDecimal(decimal.NewFromFloat(50.0))   // 默认50%分位
	factor.IndustryPercentile = models.NewJSONDecimal(decimal.NewFromFloat(50.0)) // 默认50%分位

	log.Printf("[FundamentalFactorService] 单股票因子得分计算完成 - 价值: %v, 成长: %v, 质量: %v, 盈利: %v, 综合: %v",
		factor.ValueScore.Decimal, factor.GrowthScore.Decimal, factor.QualityScore.Decimal,
		factor.ProfitabilityScore.Decimal, factor.CompositeScore.Decimal)
}

// calculateValueScore 计算价值因子得分
func (s *FundamentalFactorService) calculateValueScore(factor *models.FundamentalFactor) models.JSONDecimal {
	score := decimal.Zero
	count := 0

	// PE评分：越低越好，参考范围 0-50
	if !factor.PE.Decimal.IsZero() && factor.PE.Decimal.GreaterThan(decimal.Zero) {
		pe, _ := factor.PE.Decimal.Float64()
		if pe <= 15 {
			score = score.Add(decimal.NewFromFloat(2.0)) // 优秀
		} else if pe <= 25 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if pe <= 50 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// PB评分：越低越好，参考范围 0-10
	if !factor.PB.Decimal.IsZero() && factor.PB.Decimal.GreaterThan(decimal.Zero) {
		pb, _ := factor.PB.Decimal.Float64()
		if pb <= 1.5 {
			score = score.Add(decimal.NewFromFloat(2.0)) // 优秀
		} else if pb <= 3.0 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if pb <= 5.0 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// PS评分：越低越好，参考范围 0-10
	if !factor.PS.Decimal.IsZero() && factor.PS.Decimal.GreaterThan(decimal.Zero) {
		ps, _ := factor.PS.Decimal.Float64()
		if ps <= 2.0 {
			score = score.Add(decimal.NewFromFloat(2.0)) // 优秀
		} else if ps <= 5.0 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if ps <= 10.0 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// 计算平均得分
	if count > 0 {
		return models.NewJSONDecimal(score.Div(decimal.NewFromInt(int64(count))))
	}
	return models.NewJSONDecimal(decimal.Zero)
}

// calculateGrowthScore 计算成长因子得分
func (s *FundamentalFactorService) calculateGrowthScore(factor *models.FundamentalFactor) models.JSONDecimal {
	score := decimal.Zero
	count := 0

	// 营收增长率评分
	if !factor.RevenueGrowth.Decimal.IsZero() {
		growth, _ := factor.RevenueGrowth.Decimal.Float64()
		if growth >= 20 {
			score = score.Add(decimal.NewFromFloat(2.0)) // 优秀
		} else if growth >= 10 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if growth >= 0 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// 净利润增长率评分
	if !factor.NetProfitGrowth.Decimal.IsZero() {
		growth, _ := factor.NetProfitGrowth.Decimal.Float64()
		if growth >= 20 {
			score = score.Add(decimal.NewFromFloat(2.0)) // 优秀
		} else if growth >= 10 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if growth >= 0 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// 如果没有有效的成长数据，返回中性得分
	if count == 0 {
		return models.NewJSONDecimal(decimal.Zero)
	}

	// 计算平均得分
	return models.NewJSONDecimal(score.Div(decimal.NewFromInt(int64(count))))
}

// getHistoricalIncomeStatement 获取历史利润表数据（同比去年同期）
func (s *FundamentalFactorService) getHistoricalIncomeStatement(tsCode, currentPeriod string) *models.IncomeStatement {
	if currentPeriod == "" {
		return nil
	}

	// 计算去年同期
	if len(currentPeriod) != 8 {
		return nil
	}

	year := currentPeriod[:4]
	monthDay := currentPeriod[4:]

	currentYear, err := strconv.Atoi(year)
	if err != nil {
		return nil
	}

	lastYear := currentYear - 1
	historicalPeriod := fmt.Sprintf("%d%s", lastYear, monthDay)

	// 尝试获取历史数据
	statement, err := s.client.GetIncomeStatement(tsCode, historicalPeriod, "1")
	if err != nil {
		log.Printf("[FundamentalFactorService] 获取历史利润表数据失败 - 期间: %s, 错误: %v", historicalPeriod, err)
		return nil
	}

	if statement == nil || (statement.OperRevenue.Decimal.IsZero() && statement.NetProfit.Decimal.IsZero()) {
		log.Printf("[FundamentalFactorService] 历史利润表数据无效 - 期间: %s", historicalPeriod)
		return nil
	}

	log.Printf("[FundamentalFactorService] 成功获取历史利润表数据 - 期间: %s", historicalPeriod)
	return statement
}

// calculateGrowthRate 计算增长率 (current - historical) / historical
func (s *FundamentalFactorService) calculateGrowthRate(current, historical decimal.Decimal) models.JSONDecimal {
	if historical.IsZero() {
		return models.NewJSONDecimal(decimal.Zero)
	}

	growthRate := current.Sub(historical).Div(historical)
	return models.NewJSONDecimal(growthRate)
}

// estimateGrowthFromCurrentData 基于当前数据估算增长率（简化方法）
func (s *FundamentalFactorService) estimateGrowthFromCurrentData(value decimal.Decimal) models.JSONDecimal {
	// 这是一个简化的估算方法，实际应该基于行业数据或其他基准
	// 这里我们基于数据的规模给出一个合理的估算

	if value.IsZero() {
		return models.NewJSONDecimal(decimal.Zero)
	}

	// 将值转换为亿元单位
	valueInYi := value.Div(decimal.NewFromInt(100000000))
	valueFloat, _ := valueInYi.Float64()

	var estimatedGrowth float64
	switch {
	case valueFloat > 100: // 大于100亿
		estimatedGrowth = 0.05 // 5%
	case valueFloat > 50: // 50-100亿
		estimatedGrowth = 0.08 // 8%
	case valueFloat > 10: // 10-50亿
		estimatedGrowth = 0.12 // 12%
	case valueFloat > 1: // 1-10亿
		estimatedGrowth = 0.15 // 15%
	default: // 小于1亿
		estimatedGrowth = 0.10 // 10%
	}

	return models.NewJSONDecimal(decimal.NewFromFloat(estimatedGrowth))
}

// calculateQualityScore 计算质量因子得分
func (s *FundamentalFactorService) calculateQualityScore(factor *models.FundamentalFactor) models.JSONDecimal {
	score := decimal.Zero
	count := 0

	// ROE评分：越高越好
	if !factor.ROE.Decimal.IsZero() {
		roe, _ := factor.ROE.Decimal.Float64()
		if roe >= 20 {
			score = score.Add(decimal.NewFromFloat(2.0)) // 优秀
		} else if roe >= 15 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if roe >= 10 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// ROA评分：越高越好
	if !factor.ROA.Decimal.IsZero() {
		roa, _ := factor.ROA.Decimal.Float64()
		if roa >= 10 {
			score = score.Add(decimal.NewFromFloat(2.0)) // 优秀
		} else if roa >= 5 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if roa >= 2 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// 资产负债率评分：适中为好 (30%-60%)
	if !factor.DebtToAssets.Decimal.IsZero() {
		debt, _ := factor.DebtToAssets.Decimal.Float64()
		if debt >= 30 && debt <= 60 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if debt < 30 || (debt > 60 && debt <= 80) {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// 流动比率评分：适中为好 (1.2-2.0)
	if !factor.CurrentRatio.Decimal.IsZero() {
		current, _ := factor.CurrentRatio.Decimal.Float64()
		if current >= 1.2 && current <= 2.0 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if current >= 1.0 && current < 1.2 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// 计算平均得分
	if count > 0 {
		return models.NewJSONDecimal(score.Div(decimal.NewFromInt(int64(count))))
	}
	return models.NewJSONDecimal(decimal.Zero)
}

// calculateProfitabilityScore 计算盈利因子得分
func (s *FundamentalFactorService) calculateProfitabilityScore(factor *models.FundamentalFactor) models.JSONDecimal {
	score := decimal.Zero
	count := 0

	// 毛利率评分：越高越好
	if !factor.GrossMargin.Decimal.IsZero() {
		margin, _ := factor.GrossMargin.Decimal.Float64()
		if margin >= 40 {
			score = score.Add(decimal.NewFromFloat(2.0)) // 优秀
		} else if margin >= 25 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if margin >= 15 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// 净利率评分：越高越好
	if !factor.NetMargin.Decimal.IsZero() {
		margin, _ := factor.NetMargin.Decimal.Float64()
		if margin >= 15 {
			score = score.Add(decimal.NewFromFloat(2.0)) // 优秀
		} else if margin >= 10 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if margin >= 5 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// 营业利润率评分：越高越好
	if !factor.OperatingMargin.Decimal.IsZero() {
		margin, _ := factor.OperatingMargin.Decimal.Float64()
		if margin >= 20 {
			score = score.Add(decimal.NewFromFloat(2.0)) // 优秀
		} else if margin >= 10 {
			score = score.Add(decimal.NewFromFloat(1.0)) // 良好
		} else if margin >= 5 {
			score = score.Add(decimal.NewFromFloat(0.0)) // 一般
		} else {
			score = score.Add(decimal.NewFromFloat(-1.0)) // 较差
		}
		count++
	}

	// 计算平均得分
	if count > 0 {
		return models.NewJSONDecimal(score.Div(decimal.NewFromInt(int64(count))))
	}
	return models.NewJSONDecimal(decimal.Zero)
}

// calculateCompositeScore 计算综合得分
func (s *FundamentalFactorService) calculateCompositeScore(factor *models.FundamentalFactor) {
	// 默认权重配置
	weights := map[string]decimal.Decimal{
		"value":         decimal.NewFromFloat(0.25), // 价值因子权重
		"growth":        decimal.NewFromFloat(0.25), // 成长因子权重
		"quality":       decimal.NewFromFloat(0.25), // 质量因子权重
		"profitability": decimal.NewFromFloat(0.25), // 盈利因子权重
	}

	// 计算加权综合得分
	compositeScore := factor.ValueScore.Decimal.Mul(weights["value"]).
		Add(factor.GrowthScore.Decimal.Mul(weights["growth"])).
		Add(factor.QualityScore.Decimal.Mul(weights["quality"])).
		Add(factor.ProfitabilityScore.Decimal.Mul(weights["profitability"]))

	factor.CompositeScore = models.NewJSONDecimal(compositeScore)
}
