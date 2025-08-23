package indicators

import (
	"stock-a-future/internal/models"
	"testing"

	"github.com/shopspring/decimal"
)

// BenchmarkPatternRecognizer_RecognizeAllPatterns 基准测试图形识别性能
func BenchmarkPatternRecognizer_RecognizeAllPatterns(b *testing.B) {
	recognizer := NewPatternRecognizer()

	// 创建测试数据
	data := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.8)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.2)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240102",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.2)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.8)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(10.1)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.6)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1500)),
		},
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240103",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.6)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(11.2)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(11.0)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(2500)),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		recognizer.RecognizeAllPatterns(data)
	}
}

// BenchmarkPatternRecognizer_RecognizeDoubleCannon 基准测试双响炮识别性能
func BenchmarkPatternRecognizer_RecognizeDoubleCannon(b *testing.B) {
	recognizer := NewPatternRecognizer()

	// 创建双响炮测试数据
	data := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.3)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.8)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.3)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240102",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.3)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.8)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(10.2)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.8)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1800)),
		},
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240103",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.8)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(11.3)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(10.7)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(11.3)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(3000)),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		recognizer.RecognizeAllPatterns(data)
	}
}

// BenchmarkPatternRecognizer_RecognizeHammer 基准测试锤子线识别性能
func BenchmarkPatternRecognizer_RecognizeHammer(b *testing.B) {
	recognizer := NewPatternRecognizer()

	// 创建锤子线测试数据
	data := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.1)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.0)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		recognizer.RecognizeAllPatterns(data)
	}
}

// BenchmarkPatternRecognizer_RecognizeVolumePriceRise 基准测试量价齐升识别性能
func BenchmarkPatternRecognizer_RecognizeVolumePriceRise(b *testing.B) {
	recognizer := NewPatternRecognizer()

	// 创建量价齐升测试数据
	data := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.8)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240102",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.8)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.9)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.2)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1500)),
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		recognizer.RecognizeAllPatterns(data)
	}
}

// BenchmarkPatternRecognizer_CalculateConfidence 基准测试置信度计算性能
func BenchmarkPatternRecognizer_CalculateConfidence(b *testing.B) {
	recognizer := NewPatternRecognizer()

	priceChange := decimal.NewFromFloat(5.0)
	volumeChange := decimal.NewFromFloat(30.0)
	volumeRatio := decimal.NewFromFloat(2.0)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		recognizer.calculateConfidence(priceChange, volumeChange, volumeRatio)
	}
}

// BenchmarkPatternRecognizer_CalculateStrength 基准测试强度计算性能
func BenchmarkPatternRecognizer_CalculateStrength(b *testing.B) {
	recognizer := NewPatternRecognizer()

	confidence := decimal.NewFromInt(85)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		recognizer.calculateStrength(confidence)
	}
}
