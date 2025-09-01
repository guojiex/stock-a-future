package client

import (
	"context"
	"errors"
	"stock-a-future/internal/models"
)

// MockDataSourceClient 模拟数据源客户端，用于测试
type MockDataSourceClient struct {
	// 预设的测试数据
	DailyBasicData         *models.DailyBasic
	IncomeStatementData    *models.IncomeStatement
	BalanceSheetData       *models.BalanceSheet
	CashFlowData           *models.CashFlowStatement
	FinancialIndicatorData *models.FinancialIndicator
	StockBasicData         *models.StockBasic
	StockDailyData         []models.StockDaily
	FundamentalFactorData  *models.FundamentalFactor

	// 控制行为
	ShouldFail bool
	BaseURL    string
}

// NewMockDataSourceClient 创建新的模拟客户端
func NewMockDataSourceClient() *MockDataSourceClient {
	return &MockDataSourceClient{
		BaseURL: "http://127.0.0.1:8080",
	}
}

// ===== 基础数据接口 =====

func (m *MockDataSourceClient) GetDailyData(symbol, startDate, endDate, adjust string) ([]models.StockDaily, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	return m.StockDailyData, nil
}

func (m *MockDataSourceClient) GetDailyDataByDate(tradeDate string) ([]models.StockDaily, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	return m.StockDailyData, nil
}

func (m *MockDataSourceClient) GetStockBasic(symbol string) (*models.StockBasic, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	return m.StockBasicData, nil
}

func (m *MockDataSourceClient) GetStockList() ([]models.StockBasic, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	if m.StockBasicData != nil {
		return []models.StockBasic{*m.StockBasicData}, nil
	}
	return nil, nil
}

// ===== 基本面数据接口 =====

func (m *MockDataSourceClient) GetIncomeStatement(symbol, period, reportType string) (*models.IncomeStatement, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	return m.IncomeStatementData, nil
}

func (m *MockDataSourceClient) GetIncomeStatements(symbol, startPeriod, endPeriod, reportType string) ([]models.IncomeStatement, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	if m.IncomeStatementData != nil {
		return []models.IncomeStatement{*m.IncomeStatementData}, nil
	}
	return nil, nil
}

func (m *MockDataSourceClient) GetBalanceSheet(symbol, period, reportType string) (*models.BalanceSheet, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	return m.BalanceSheetData, nil
}

func (m *MockDataSourceClient) GetBalanceSheets(symbol, startPeriod, endPeriod, reportType string) ([]models.BalanceSheet, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	if m.BalanceSheetData != nil {
		return []models.BalanceSheet{*m.BalanceSheetData}, nil
	}
	return nil, nil
}

func (m *MockDataSourceClient) GetCashFlowStatement(symbol, period, reportType string) (*models.CashFlowStatement, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	return m.CashFlowData, nil
}

func (m *MockDataSourceClient) GetCashFlowStatements(symbol, startPeriod, endPeriod, reportType string) ([]models.CashFlowStatement, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	if m.CashFlowData != nil {
		return []models.CashFlowStatement{*m.CashFlowData}, nil
	}
	return nil, nil
}

func (m *MockDataSourceClient) GetFinancialIndicator(symbol, period, reportType string) (*models.FinancialIndicator, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	return m.FinancialIndicatorData, nil
}

func (m *MockDataSourceClient) GetFinancialIndicators(symbol, startPeriod, endPeriod, reportType string) ([]models.FinancialIndicator, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	if m.FinancialIndicatorData != nil {
		return []models.FinancialIndicator{*m.FinancialIndicatorData}, nil
	}
	return nil, nil
}

func (m *MockDataSourceClient) GetDailyBasic(ctx context.Context, symbol, tradeDate string) (*models.DailyBasic, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	return m.DailyBasicData, nil
}

func (m *MockDataSourceClient) GetDailyBasics(ctx context.Context, symbol, startDate, endDate string) ([]models.DailyBasic, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	if m.DailyBasicData != nil {
		return []models.DailyBasic{*m.DailyBasicData}, nil
	}
	return nil, nil
}

func (m *MockDataSourceClient) GetDailyBasicsByDate(ctx context.Context, tradeDate string) ([]models.DailyBasic, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	if m.DailyBasicData != nil {
		return []models.DailyBasic{*m.DailyBasicData}, nil
	}
	return nil, nil
}

// ===== 基本面因子接口 =====

func (m *MockDataSourceClient) GetFundamentalFactor(symbol, tradeDate string) (*models.FundamentalFactor, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	return m.FundamentalFactorData, nil
}

func (m *MockDataSourceClient) GetFundamentalFactors(symbol, startDate, endDate string) ([]models.FundamentalFactor, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	if m.FundamentalFactorData != nil {
		return []models.FundamentalFactor{*m.FundamentalFactorData}, nil
	}
	return nil, nil
}

func (m *MockDataSourceClient) GetFundamentalFactorsByDate(tradeDate string) ([]models.FundamentalFactor, error) {
	if m.ShouldFail {
		return nil, errors.New("模拟错误")
	}
	if m.FundamentalFactorData != nil {
		return []models.FundamentalFactor{*m.FundamentalFactorData}, nil
	}
	return nil, nil
}

// ===== 系统接口 =====

func (m *MockDataSourceClient) GetBaseURL() string {
	return m.BaseURL
}

func (m *MockDataSourceClient) TestConnection() error {
	if m.ShouldFail {
		return errors.New("连接失败")
	}
	return nil
}
