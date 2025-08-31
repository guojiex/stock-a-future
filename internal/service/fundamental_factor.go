package service

import (
	"context"
	"fmt"
	"log"
	"sort"
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

	for _, period := range periods {
		statement, err := s.client.GetIncomeStatement(symbol, period, "1")
		if err == nil && statement != nil {
			return statement, nil
		}
	}

	return nil, fmt.Errorf("未找到有效的利润表数据")
}

func (s *FundamentalFactorService) getLatestBalanceSheet(symbol string) (*models.BalanceSheet, error) {
	// 尝试获取最近几个报告期的数据
	periods := []string{"20241231", "20240930", "20240630", "20240331", "20231231"}

	for _, period := range periods {
		sheet, err := s.client.GetBalanceSheet(symbol, period, "1")
		if err == nil && sheet != nil {
			return sheet, nil
		}
	}

	return nil, fmt.Errorf("未找到有效的资产负债表数据")
}

func (s *FundamentalFactorService) getLatestCashFlow(symbol string) (*models.CashFlowStatement, error) {
	// 尝试获取最近几个报告期的数据
	periods := []string{"20241231", "20240930", "20240630", "20240331", "20231231"}

	for _, period := range periods {
		cashFlow, err := s.client.GetCashFlowStatement(symbol, period, "1")
		if err == nil && cashFlow != nil {
			return cashFlow, nil
		}
	}

	return nil, fmt.Errorf("未找到有效的现金流量表数据")
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

	// 营收增长率 - 需要历史数据计算，暂时设为0
	factor.RevenueGrowth = models.NewJSONDecimal(decimal.Zero)

	// 净利润增长率 - 需要历史数据计算，暂时设为0
	factor.NetProfitGrowth = models.NewJSONDecimal(decimal.Zero)

	// EPS增长率 - 需要历史数据计算，暂时设为0
	factor.EPSGrowth = models.NewJSONDecimal(decimal.Zero)

	// ROE增长率 - 需要历史数据计算，暂时设为0
	factor.ROEGrowth = models.NewJSONDecimal(decimal.Zero)
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
