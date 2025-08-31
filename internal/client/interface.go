package client

import (
	"context"
	"fmt"
	"stock-a-future/internal/models"
)

// DataSourceClient 数据源客户端接口
type DataSourceClient interface {
	// ===== 基础数据接口 =====

	// 获取股票日线数据
	GetDailyData(symbol, startDate, endDate, adjust string) ([]models.StockDaily, error)

	// 根据交易日期获取所有股票数据
	GetDailyDataByDate(tradeDate string) ([]models.StockDaily, error)

	// 获取股票基本信息
	GetStockBasic(symbol string) (*models.StockBasic, error)

	// 获取股票列表
	GetStockList() ([]models.StockBasic, error)

	// ===== 基本面数据接口 =====

	// 获取利润表数据
	// symbol: 股票代码, period: 报告期(YYYYMMDD), reportType: 报告类型(1-年报,2-中报,3-季报)
	GetIncomeStatement(symbol, period, reportType string) (*models.IncomeStatement, error)

	// 批量获取利润表数据
	GetIncomeStatements(symbol, startPeriod, endPeriod, reportType string) ([]models.IncomeStatement, error)

	// 获取资产负债表数据
	GetBalanceSheet(symbol, period, reportType string) (*models.BalanceSheet, error)

	// 批量获取资产负债表数据
	GetBalanceSheets(symbol, startPeriod, endPeriod, reportType string) ([]models.BalanceSheet, error)

	// 获取现金流量表数据
	GetCashFlowStatement(symbol, period, reportType string) (*models.CashFlowStatement, error)

	// 批量获取现金流量表数据
	GetCashFlowStatements(symbol, startPeriod, endPeriod, reportType string) ([]models.CashFlowStatement, error)

	// 获取财务指标数据
	GetFinancialIndicator(symbol, period, reportType string) (*models.FinancialIndicator, error)

	// 批量获取财务指标数据
	GetFinancialIndicators(symbol, startPeriod, endPeriod, reportType string) ([]models.FinancialIndicator, error)

	// 获取每日基本面指标
	GetDailyBasic(ctx context.Context, symbol, tradeDate string) (*models.DailyBasic, error)

	// 批量获取每日基本面指标
	GetDailyBasics(ctx context.Context, symbol, startDate, endDate string) ([]models.DailyBasic, error)

	// 根据交易日期获取所有股票的每日基本面指标
	GetDailyBasicsByDate(ctx context.Context, tradeDate string) ([]models.DailyBasic, error)

	// ===== 基本面因子接口 =====

	// 获取基本面因子数据
	GetFundamentalFactor(symbol, tradeDate string) (*models.FundamentalFactor, error)

	// 批量获取基本面因子数据
	GetFundamentalFactors(symbol, startDate, endDate string) ([]models.FundamentalFactor, error)

	// 根据交易日期获取所有股票的基本面因子
	GetFundamentalFactorsByDate(tradeDate string) ([]models.FundamentalFactor, error)

	// ===== 系统接口 =====

	// 获取基础URL
	GetBaseURL() string

	// 测试连接
	TestConnection() error
}

// FundamentalFactorCalculator 基本面因子计算器接口
type FundamentalFactorCalculator interface {
	// 计算单个股票的基本面因子
	CalculateFundamentalFactor(symbol, tradeDate string) (*models.FundamentalFactor, error)

	// 批量计算基本面因子
	CalculateFundamentalFactors(symbols []string, tradeDate string) ([]models.FundamentalFactor, error)

	// 计算价值因子
	CalculateValueFactors(symbol, tradeDate string, dailyBasic *models.DailyBasic, financialIndicator *models.FinancialIndicator) (map[string]models.JSONDecimal, error)

	// 计算成长因子
	CalculateGrowthFactors(symbol, tradeDate string, currentIndicator, previousIndicator *models.FinancialIndicator) (map[string]models.JSONDecimal, error)

	// 计算质量因子
	CalculateQualityFactors(symbol, tradeDate string, balanceSheet *models.BalanceSheet, financialIndicator *models.FinancialIndicator) (map[string]models.JSONDecimal, error)

	// 计算盈利因子
	CalculateProfitabilityFactors(symbol, tradeDate string, incomeStatement *models.IncomeStatement, balanceSheet *models.BalanceSheet) (map[string]models.JSONDecimal, error)

	// 计算运营效率因子
	CalculateOperatingFactors(symbol, tradeDate string, incomeStatement *models.IncomeStatement, balanceSheet *models.BalanceSheet) (map[string]models.JSONDecimal, error)

	// 计算现金流因子
	CalculateCashFlowFactors(symbol, tradeDate string, incomeStatement *models.IncomeStatement, cashFlowStatement *models.CashFlowStatement) (map[string]models.JSONDecimal, error)

	// 标准化因子得分
	NormalizeFactorScores(factors []models.FundamentalFactor) ([]models.FundamentalFactor, error)

	// 计算行业排名
	CalculateIndustryRankings(factors []models.FundamentalFactor, industry string) ([]models.FundamentalFactor, error)

	// 计算市场排名
	CalculateMarketRankings(factors []models.FundamentalFactor) ([]models.FundamentalFactor, error)
}

// DataSourceType 数据源类型
type DataSourceType string

const (
	DataSourceTushare DataSourceType = "tushare"
	DataSourceAKTools DataSourceType = "aktools"
)

// DataSourceFactory 数据源工厂
type DataSourceFactory struct{}

// NewDataSourceFactory 创建数据源工厂
func NewDataSourceFactory() *DataSourceFactory {
	return &DataSourceFactory{}
}

// CreateClient 根据类型创建数据源客户端
func (f *DataSourceFactory) CreateClient(sourceType DataSourceType, config map[string]string) (DataSourceClient, error) {
	switch sourceType {
	case DataSourceTushare:
		token := config["token"]
		baseURL := config["base_url"]
		if baseURL == "" {
			baseURL = "http://api.tushare.pro"
		}
		return NewTushareClient(token, baseURL), nil

	case DataSourceAKTools:
		baseURL := config["base_url"]
		if baseURL == "" {
			baseURL = "http://127.0.0.1:8080"
		}
		return NewAKToolsClient(baseURL), nil

	default:
		return nil, fmt.Errorf("不支持的数据源类型: %s", sourceType)
	}
}
