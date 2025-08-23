package service

import (
	"stock-a-future/internal/models"
	"testing"

	"github.com/shopspring/decimal"
)

// MockStockServiceForPattern 模拟股票服务接口
type MockStockServiceForPattern struct {
	shouldFail bool
	stockData  []models.StockDaily
}

func (m *MockStockServiceForPattern) GetDailyData(tsCode, startDate, endDate, adjust string) ([]models.StockDaily, error) {
	if m.shouldFail {
		return nil, &MockErrorForPattern{message: "模拟错误"}
	}
	return m.stockData, nil
}

// MockErrorForPattern 模拟错误类型
type MockErrorForPattern struct {
	message string
}

func (e *MockErrorForPattern) Error() string {
	return e.message
}

func TestNewPatternService(t *testing.T) {
	stockService := &MockStockServiceForPattern{}

	service := NewPatternService(stockService)

	if service == nil {
		t.Fatal("NewPatternService应该返回非空的PatternService实例")
	}

	if service.stockService != stockService {
		t.Error("StockService应该正确设置")
	}

	if service.patternRecognizer == nil {
		t.Error("PatternRecognizer应该被初始化")
	}
}

func TestPatternService_RecognizePatterns(t *testing.T) {
	// 创建模拟股票数据
	stockData := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.8)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.2)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
	}

	stockService := &MockStockServiceForPattern{
		shouldFail: false,
		stockData:  stockData,
	}

	service := NewPatternService(stockService)

	// 测试正常情况
	patterns, err := service.RecognizePatterns("000001.SZ", "20240101", "20240101")

	if err != nil {
		t.Fatalf("RecognizePatterns不应该返回错误: %v", err)
	}

	if patterns == nil {
		t.Error("RecognizePatterns应该返回非空结果")
	}
}

func TestPatternService_RecognizePatterns_Error(t *testing.T) {
	stockService := &MockStockServiceForPattern{
		shouldFail: true,
	}

	service := NewPatternService(stockService)

	// 测试错误情况
	patterns, err := service.RecognizePatterns("000001.SZ", "20240101", "20240101")

	if err == nil {
		t.Error("当底层服务失败时，RecognizePatterns应该返回错误")
	}

	if patterns != nil {
		t.Error("当发生错误时，patterns应该为nil")
	}
}

func TestPatternService_SearchPatterns(t *testing.T) {
	// 创建模拟股票数据
	stockData := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.8)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.2)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
	}

	stockService := &MockStockServiceForPattern{
		shouldFail: false,
		stockData:  stockData,
	}

	service := NewPatternService(stockService)

	// 创建搜索请求
	request := models.PatternSearchRequest{
		TSCode:        "000001.SZ",
		StartDate:     "20240101",
		EndDate:       "20240101",
		Patterns:      []string{"双响炮"},
		MinConfidence: 0.5,
	}

	// 测试搜索
	response, err := service.SearchPatterns(request)

	if err != nil {
		t.Fatalf("SearchPatterns不应该返回错误: %v", err)
	}

	if response == nil {
		t.Error("SearchPatterns应该返回非空响应")
	}

	if response.Total < 0 {
		t.Error("Total字段应该大于等于0")
	}
}

func TestPatternService_SearchPatterns_Error(t *testing.T) {
	stockService := &MockStockServiceForPattern{
		shouldFail: true,
	}

	service := NewPatternService(stockService)

	request := models.PatternSearchRequest{
		TSCode:    "000001.SZ",
		StartDate: "20240101",
		EndDate:   "20240101",
	}

	// 测试错误情况
	response, err := service.SearchPatterns(request)

	if err == nil {
		t.Error("当底层服务失败时，SearchPatterns应该返回错误")
	}

	if response != nil {
		t.Error("当发生错误时，response应该为nil")
	}
}

func TestPatternService_GetPatternSummary(t *testing.T) {
	// 创建模拟股票数据
	stockData := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.8)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.2)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
	}

	stockService := &MockStockServiceForPattern{
		shouldFail: false,
		stockData:  stockData,
	}

	service := NewPatternService(stockService)

	// 测试获取摘要
	summary, err := service.GetPatternSummary("000001.SZ", 30)

	if err != nil {
		t.Fatalf("GetPatternSummary不应该返回错误: %v", err)
	}

	if summary == nil {
		t.Error("GetPatternSummary应该返回非空摘要")
	}
}

func TestPatternService_GetPatternSummary_Error(t *testing.T) {
	stockService := &MockStockServiceForPattern{
		shouldFail: true,
	}

	service := NewPatternService(stockService)

	// 测试错误情况
	summary, err := service.GetPatternSummary("000001.SZ", 30)

	if err == nil {
		t.Error("当底层服务失败时，GetPatternSummary应该返回错误")
	}

	if summary != nil {
		t.Error("当发生错误时，summary应该为nil")
	}
}
