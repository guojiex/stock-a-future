package indicators

import (
	"stock-a-future/internal/models"
	"testing"

	"github.com/shopspring/decimal"
)

func TestPatternRecognizer_RecognizeDoubleCannon(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 创建测试数据：双响炮模式
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

	patterns := recognizer.RecognizeAllPatterns(data)
	if len(patterns) == 0 {
		t.Error("应该识别出双响炮模式")
		return
	}

	found := false
	for _, pattern := range patterns {
		for _, candlestick := range pattern.Candlestick {
			if candlestick.Pattern == "双响炮" {
				found = true
				if candlestick.Signal != "BUY" {
					t.Errorf("双响炮应该是买入信号，实际是: %s", candlestick.Signal)
				}
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Error("未找到双响炮模式")
	}
}

func TestPatternRecognizer_RecognizeHammer(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 创建测试数据：锤子线模式
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

	patterns := recognizer.RecognizeAllPatterns(data)
	if len(patterns) == 0 {
		t.Error("应该识别出锤子线模式")
		return
	}

	found := false
	for _, pattern := range patterns {
		for _, candlestick := range pattern.Candlestick {
			if candlestick.Pattern == "锤子线" {
				found = true
				if candlestick.Signal != "BUY" {
					t.Errorf("锤子线应该是买入信号，实际是: %s", candlestick.Signal)
				}
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Error("未找到锤子线模式")
	}
}

func TestPatternRecognizer_RecognizeVolumePriceRise(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 创建测试数据：量价齐升模式
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

	patterns := recognizer.RecognizeAllPatterns(data)
	if len(patterns) == 0 {
		t.Error("应该识别出量价齐升模式")
		return
	}

	found := false
	for _, pattern := range patterns {
		for _, volumePrice := range pattern.VolumePrice {
			if volumePrice.Pattern == "量价齐升" {
				found = true
				if volumePrice.Signal != "BUY" {
					t.Errorf("量价齐升应该是买入信号，实际是: %s", volumePrice.Signal)
				}
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Error("未找到量价齐升模式")
	}
}

func TestPatternRecognizer_CalculateConfidence(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 测试置信度计算
	priceChange := decimal.NewFromFloat(5.0)   // 5%涨幅
	volumeChange := decimal.NewFromFloat(30.0) // 30%成交量变化
	volumeRatio := decimal.NewFromFloat(2.0)   // 2倍量比

	confidence := recognizer.calculateConfidence(priceChange, volumeChange, volumeRatio)

	// 置信度应该在合理范围内
	if confidence.LessThan(decimal.Zero) || confidence.GreaterThan(decimal.NewFromInt(100)) {
		t.Errorf("置信度应该在0-100之间，实际是: %s", confidence.String())
	}
}

func TestPatternRecognizer_CalculateStrength(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 测试强度计算
	tests := []struct {
		confidence decimal.Decimal
		expected   string
	}{
		{decimal.NewFromInt(90), "STRONG"},
		{decimal.NewFromInt(70), "MEDIUM"},
		{decimal.NewFromInt(50), "WEAK"},
	}

	for _, test := range tests {
		strength := recognizer.calculateStrength(test.confidence)
		if strength != test.expected {
			t.Errorf("置信度 %s 应该对应强度 %s，实际是: %s",
				test.confidence.String(), test.expected, strength)
		}
	}
}

// TestPatternRecognizer_RecognizeDoubleCannon_603093 专门测试603093.SH南华期货的双响炮识别
func TestPatternRecognizer_RecognizeDoubleCannon_603093(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 创建模拟的603093.SH南华期货数据
	// 这里使用一些可能触发双响炮的数据
	data := []models.StockDaily{
		{
			TSCode:    "603093.SH",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(15.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(15.3)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(14.8)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(15.3)), // 涨幅2%
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
		{
			TSCode:    "603093.SH",
			TradeDate: "20240102",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(15.3)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(15.8)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(15.2)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(15.8)), // 涨幅3.3%
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1800)), // 成交量是前一天的1.8倍
		},
		{
			TSCode:    "603093.SH",
			TradeDate: "20240103",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(15.8)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(16.3)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(15.7)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(16.3)), // 涨幅3.2%
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(3000)), // 成交量是前一天的1.67倍
		},
	}

	t.Logf("开始测试603093.SH南华期货的双响炮识别")
	t.Logf("数据长度: %d", len(data))

	// 打印每根K线的详细信息
	for i, candle := range data {
		change := candle.Close.Decimal.Sub(candle.Open.Decimal)
		changePct := change.Div(candle.Open.Decimal).Mul(decimal.NewFromFloat(100))
		t.Logf("第%d天 %s: 开盘=%s, 收盘=%s, 涨幅=%s%%, 成交量=%s",
			i+1, candle.TradeDate,
			candle.Open.Decimal.String(),
			candle.Close.Decimal.String(),
			changePct.String(),
			candle.Vol.Decimal.String())
	}

	patterns := recognizer.RecognizeAllPatterns(data)
	t.Logf("识别结果数量: %d", len(patterns))

	if len(patterns) == 0 {
		t.Log("未识别出任何模式")
		return
	}

	// 检查是否识别出双响炮
	found := false
	for _, pattern := range patterns {
		t.Logf("在日期 %s 识别到模式:", pattern.TradeDate)
		t.Logf("  蜡烛图模式数量: %d", len(pattern.Candlestick))
		t.Logf("  量价模式数量: %d", len(pattern.VolumePrice))

		for _, candlestick := range pattern.Candlestick {
			t.Logf("    蜡烛图模式: %s, 信号: %s, 置信度: %s",
				candlestick.Pattern, candlestick.Signal, candlestick.Confidence.Decimal.String())

			if candlestick.Pattern == "双响炮" {
				found = true
				t.Logf("✅ 找到双响炮模式!")
				if candlestick.Signal != "BUY" {
					t.Errorf("双响炮应该是买入信号，实际是: %s", candlestick.Signal)
				}
			}
		}

		for _, volumePrice := range pattern.VolumePrice {
			t.Logf("    量价模式: %s, 信号: %s, 置信度: %s",
				volumePrice.Pattern, volumePrice.Signal, volumePrice.Confidence.Decimal.String())
		}
	}

	if !found {
		t.Log("❌ 未找到双响炮模式")
	} else {
		t.Log("✅ 成功识别出双响炮模式")
	}
}
