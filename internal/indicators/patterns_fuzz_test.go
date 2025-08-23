package indicators

import (
	"encoding/json"
	"stock-a-future/internal/models"
	"testing"

	"github.com/shopspring/decimal"
)

// FuzzPatternRecognizer_RecognizeAllPatterns 模糊测试图形识别
func FuzzPatternRecognizer_RecognizeAllPatterns(f *testing.F) {
	recognizer := NewPatternRecognizer()

	// 添加种子数据
	seedData := []models.StockDaily{
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

	f.Add(mustMarshalJSON(seedData))

	f.Fuzz(func(t *testing.T, data []byte) {
		var stockData []models.StockDaily
		if err := json.Unmarshal(data, &stockData); err != nil {
			// 如果数据无效，跳过测试
			return
		}

		// 验证数据有效性
		if len(stockData) == 0 {
			return
		}

		// 检查每个股票数据的基本字段
		for _, stock := range stockData {
			if stock.TSCode == "" || stock.TradeDate == "" {
				return
			}

			// 检查价格数据的有效性
			if stock.Open.Decimal.LessThan(decimal.Zero) ||
				stock.High.Decimal.LessThan(decimal.Zero) ||
				stock.Low.Decimal.LessThan(decimal.Zero) ||
				stock.Close.Decimal.LessThan(decimal.Zero) {
				return
			}

			// 检查价格逻辑关系
			if stock.High.Decimal.LessThan(stock.Low.Decimal) {
				return
			}
		}

		// 执行图形识别
		patterns := recognizer.RecognizeAllPatterns(stockData)

		// 验证结果不为nil
		if patterns == nil {
			t.Error("RecognizeAllPatterns应该返回非nil结果")
		}

		// 验证结果类型
		for _, pattern := range patterns {
			if pattern.TSCode == "" {
				t.Error("模式结果应该包含TSCode")
			}

			if pattern.TradeDate == "" {
				t.Error("模式结果应该包含TradeDate")
			}
		}
	})
}

// FuzzPatternRecognizer_CalculateConfidence 模糊测试置信度计算
func FuzzPatternRecognizer_CalculateConfidence(f *testing.F) {
	recognizer := NewPatternRecognizer()

	// 添加种子数据
	f.Add(5.0, 30.0, 2.0)
	f.Add(0.0, 0.0, 1.0)
	f.Add(-5.0, -30.0, 0.5)
	f.Add(100.0, 500.0, 10.0)

	f.Fuzz(func(t *testing.T, priceChange, volumeChange, volumeRatio float64) {
		// 转换为decimal
		priceChangeDec := decimal.NewFromFloat(priceChange)
		volumeChangeDec := decimal.NewFromFloat(volumeChange)
		volumeRatioDec := decimal.NewFromFloat(volumeRatio)

		// 执行置信度计算
		confidence := recognizer.calculateConfidence(priceChangeDec, volumeChangeDec, volumeRatioDec)

		// 验证置信度在合理范围内
		if confidence.LessThan(decimal.Zero) {
			t.Errorf("置信度不应该为负数: %v", confidence)
		}

		if confidence.GreaterThan(decimal.NewFromInt(100)) {
			t.Errorf("置信度不应该超过100: %v", confidence)
		}
	})
}

// FuzzPatternRecognizer_CalculateStrength 模糊测试强度计算
func FuzzPatternRecognizer_CalculateStrength(f *testing.F) {
	recognizer := NewPatternRecognizer()

	// 添加种子数据
	f.Add(0.0)
	f.Add(50.0)
	f.Add(100.0)
	f.Add(-10.0)
	f.Add(150.0)

	f.Fuzz(func(t *testing.T, confidence float64) {
		// 转换为decimal
		confidenceDec := decimal.NewFromFloat(confidence)

		// 执行强度计算
		strength := recognizer.calculateStrength(confidenceDec)

		// 验证强度值在有效范围内
		validStrengths := map[string]bool{
			"WEAK":   true,
			"MEDIUM": true,
			"STRONG": true,
		}

		if !validStrengths[strength] {
			t.Errorf("强度值无效: %s", strength)
		}
	})
}

// FuzzJSONDecimal_MarshalUnmarshal 模糊测试JSONDecimal的序列化和反序列化
func FuzzJSONDecimal_MarshalUnmarshal(f *testing.F) {
	// 添加种子数据
	f.Add(0.0)
	f.Add(123.45)
	f.Add(-123.45)
	f.Add(999999.99)
	f.Add(0.001)

	f.Fuzz(func(t *testing.T, value float64) {
		// 创建JSONDecimal
		original := models.NewJSONDecimal(decimal.NewFromFloat(value))

		// 序列化
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}

		// 反序列化
		var restored models.JSONDecimal
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatalf("反序列化失败: %v", err)
		}

		// 验证值相等
		if !original.Decimal.Equal(restored.Decimal) {
			t.Errorf("序列化前后值不相等: 原始值 %v, 恢复值 %v",
				original.Decimal, restored.Decimal)
		}
	})
}

// mustMarshalJSON 辅助函数，将数据序列化为JSON
func mustMarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// FuzzStockDaily_Validation 模糊测试股票日线数据验证
func FuzzStockDaily_Validation(f *testing.F) {
	// 添加种子数据
	seedData := models.StockDaily{
		TSCode:    "000001.SZ",
		TradeDate: "20240101",
		Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
		High:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
		Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.8)),
		Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.2)),
		Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
	}

	f.Add(mustMarshalJSON(seedData))

	f.Fuzz(func(t *testing.T, data []byte) {
		var stock models.StockDaily
		if err := json.Unmarshal(data, &stock); err != nil {
			// 如果数据无效，跳过测试
			return
		}

		// 验证基本字段
		if stock.TSCode == "" {
			t.Error("TSCode不应该为空")
		}

		if stock.TradeDate == "" {
			t.Error("TradeDate不应该为空")
		}

		// 验证价格数据
		if stock.Open.Decimal.LessThan(decimal.Zero) {
			t.Error("开盘价不应该为负数")
		}

		if stock.High.Decimal.LessThan(decimal.Zero) {
			t.Error("最高价不应该为负数")
		}

		if stock.Low.Decimal.LessThan(decimal.Zero) {
			t.Error("最低价不应该为负数")
		}

		if stock.Close.Decimal.LessThan(decimal.Zero) {
			t.Error("收盘价不应该为负数")
		}

		// 验证价格逻辑关系
		if stock.High.Decimal.LessThan(stock.Low.Decimal) {
			t.Error("最高价应该大于等于最低价")
		}

		if stock.Vol.Decimal.LessThan(decimal.Zero) {
			t.Error("成交量不应该为负数")
		}
	})
}
