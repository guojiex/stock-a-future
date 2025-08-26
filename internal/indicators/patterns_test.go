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

	// 直接测试锤子线识别函数
	t.Logf("测试锤子线识别 - 输入数据: Open=%v, High=%v, Low=%v, Close=%v",
		data[0].Open.Decimal, data[0].High.Decimal, data[0].Low.Decimal, data[0].Close.Decimal)

	// 计算实体和影线
	body := data[0].Close.Decimal.Sub(data[0].Open.Decimal).Abs()
	higherPrice := data[0].Close.Decimal
	lowerPrice := data[0].Open.Decimal
	if data[0].Open.Decimal.GreaterThan(data[0].Close.Decimal) {
		higherPrice = data[0].Open.Decimal
		lowerPrice = data[0].Close.Decimal
	}
	upperShadow := data[0].High.Decimal.Sub(higherPrice)
	lowerShadow := lowerPrice.Sub(data[0].Low.Decimal)

	t.Logf("实体长度: %v, 上影线: %v, 下影线: %v", body, upperShadow, lowerShadow)
	t.Logf("下影线是否大于实体的2倍: %v", lowerShadow.GreaterThan(body.Mul(decimal.NewFromInt(2))))
	t.Logf("上影线是否小于实体的0.5倍: %v", upperShadow.LessThanOrEqual(body.Mul(decimal.NewFromFloat(0.5))))

	// 直接调用锤子线识别函数
	hammerPattern := recognizer.recognizeHammer(data[0])
	if hammerPattern != nil {
		t.Logf("✅ 直接调用recognizeHammer成功识别出锤子线，置信度: %v", hammerPattern.Confidence.Decimal)
	} else {
		t.Logf("❌ 直接调用recognizeHammer未识别出锤子线")
	}

	// 通过完整流程测试
	patterns := recognizer.RecognizeAllPatterns(data)
	if len(patterns) == 0 {
		t.Error("应该识别出锤子线模式")
		t.Logf("RecognizeAllPatterns返回空结果")
		return
	} else {
		t.Logf("RecognizeAllPatterns返回结果数量: %d", len(patterns))
	}

	found := false
	for _, pattern := range patterns {
		t.Logf("检查日期 %s 的模式:", pattern.TradeDate)
		t.Logf("  蜡烛图模式数量: %d", len(pattern.Candlestick))

		for _, candlestick := range pattern.Candlestick {
			t.Logf("  模式: %s, 信号: %s", candlestick.Pattern, candlestick.Signal)
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

	// 直接测试量价齐升识别函数
	t.Logf("测试量价齐升识别 - 当前数据: Close=%v, Vol=%v",
		data[1].Close.Decimal, data[1].Vol.Decimal)
	t.Logf("前一天数据: Close=%v, Vol=%v",
		data[0].Close.Decimal, data[0].Vol.Decimal)

	// 计算价格变化
	priceChange := data[1].Close.Decimal.Sub(data[0].Close.Decimal)
	priceChangePct := priceChange.Div(data[0].Close.Decimal).Mul(decimal.NewFromInt(100))

	// 计算成交量变化
	volumeChange := data[1].Vol.Decimal.Sub(data[0].Vol.Decimal)
	volumeChangePct := volumeChange.Div(data[0].Vol.Decimal).Mul(decimal.NewFromInt(100))

	t.Logf("价格变化: %v (%v%%)", priceChange, priceChangePct)
	t.Logf("成交量变化: %v (%v%%)", volumeChange, volumeChangePct)
	t.Logf("价格是否上涨: %v", priceChangePct.GreaterThanOrEqual(decimal.Zero))
	t.Logf("成交量是否增加: %v", volumeChangePct.GreaterThanOrEqual(decimal.Zero))

	// 直接调用量价齐升识别函数
	volumePricePattern := recognizer.recognizeVolumePriceRise(data[1], data[0], data[0])
	if volumePricePattern != nil {
		t.Logf("✅ 直接调用recognizeVolumePriceRise成功识别出量价齐升，置信度: %v",
			volumePricePattern.Confidence.Decimal)
	} else {
		t.Logf("❌ 直接调用recognizeVolumePriceRise未识别出量价齐升")
	}

	// 通过完整流程测试
	patterns := recognizer.RecognizeAllPatterns(data)
	if len(patterns) == 0 {
		t.Error("应该识别出量价齐升模式")
		t.Logf("RecognizeAllPatterns返回空结果")
		return
	} else {
		t.Logf("RecognizeAllPatterns返回结果数量: %d", len(patterns))
	}

	found := false
	for _, pattern := range patterns {
		t.Logf("检查日期 %s 的模式:", pattern.TradeDate)
		t.Logf("  量价模式数量: %d", len(pattern.VolumePrice))

		for _, volumePrice := range pattern.VolumePrice {
			t.Logf("  模式: %s, 信号: %s", volumePrice.Pattern, volumePrice.Signal)
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

// TestPatternRecognizer_RecognizeDoji 测试十字星模式识别
func TestPatternRecognizer_RecognizeDoji(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 创建测试数据：十字星模式
	data := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.5)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.02)), // 几乎等于开盘价
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
	}

	patterns := recognizer.RecognizeAllPatterns(data)
	found := false
	for _, pattern := range patterns {
		for _, candlestick := range pattern.Candlestick {
			if candlestick.Pattern == "十字星" {
				found = true
				if candlestick.Signal != "HOLD" {
					t.Errorf("十字星应该是观望信号，实际是: %s", candlestick.Signal)
				}
				break
			}
		}
	}

	if !found {
		t.Error("未找到十字星模式")
	}
}

// TestPatternRecognizer_RecognizeEngulfing 测试吞没模式识别
func TestPatternRecognizer_RecognizeEngulfing(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 创建测试数据：看涨吞没模式
	data := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.6)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.1)), // 小阴线
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240102",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(9.9)),  // 低开
			High:      models.NewJSONDecimal(decimal.NewFromFloat(11.0)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.8)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.8)), // 大阳线，吞没前一根
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1500)),
		},
	}

	patterns := recognizer.RecognizeAllPatterns(data)
	found := false
	for _, pattern := range patterns {
		for _, candlestick := range pattern.Candlestick {
			if candlestick.Pattern == "吞没模式" {
				found = true
				if candlestick.Signal != "BUY" {
					t.Errorf("看涨吞没应该是买入信号，实际是: %s", candlestick.Signal)
				}
				break
			}
		}
	}

	if !found {
		t.Error("未找到吞没模式")
	}
}

// TestPatternRecognizer_RecognizeShootingStar 测试射击之星模式识别
func TestPatternRecognizer_RecognizeShootingStar(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 创建测试数据：射击之星模式
	data := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(11.0)), // 很长的上影线
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.95)), // 很短的下影线
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.1)), // 小实体
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
	}

	patterns := recognizer.RecognizeAllPatterns(data)
	found := false
	for _, pattern := range patterns {
		for _, candlestick := range pattern.Candlestick {
			if candlestick.Pattern == "射击之星" {
				found = true
				if candlestick.Signal != "SELL" {
					t.Errorf("射击之星应该是卖出信号，实际是: %s", candlestick.Signal)
				}
				break
			}
		}
	}

	if !found {
		t.Error("未找到射击之星模式")
	}
}

// TestPatternRecognizer_RecognizeThreeBlackCrows 测试三只乌鸦模式识别
func TestPatternRecognizer_RecognizeThreeBlackCrows(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 创建测试数据：三只乌鸦模式
	data := []models.StockDaily{
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240101",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(11.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(11.1)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.6)), // 第一根阴线
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)),
		},
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240102",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.8)), // 在前一根实体内开盘
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.9)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.1)), // 第二根阴线，收盘更低
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1200)),
		},
		{
			TSCode:    "000001.SZ",
			TradeDate: "20240103",
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.3)), // 在前一根实体内开盘
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.4)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.5)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(9.6)), // 第三根阴线，收盘再低
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1500)),
		},
	}

	patterns := recognizer.RecognizeAllPatterns(data)
	found := false
	for _, pattern := range patterns {
		for _, candlestick := range pattern.Candlestick {
			if candlestick.Pattern == "三只乌鸦" {
				found = true
				if candlestick.Signal != "SELL" {
					t.Errorf("三只乌鸦应该是卖出信号，实际是: %s", candlestick.Signal)
				}
				break
			}
		}
	}

	if !found {
		t.Error("未找到三只乌鸦模式")
	}
}

// TestPatternRecognizer_RecognizeLowVolumePrice 测试地量地价模式识别
func TestPatternRecognizer_RecognizeLowVolumePrice(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 创建测试数据：地量地价模式
	// 需要至少20天的历史数据
	data := make([]models.StockDaily, 21)
	
	// 填充前20天的数据（作为基准）
	for i := 0; i < 20; i++ {
		data[i] = models.StockDaily{
			TSCode:    "000001.SZ",
			TradeDate: "2024010" + string(rune('1'+i%10)),
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.5)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(2000)), // 平均成交量
		}
	}
	
	// 最后一天：地量地价
	data[20] = models.StockDaily{
		TSCode:    "000001.SZ",
		TradeDate: "20240121",
		Open:      models.NewJSONDecimal(decimal.NewFromFloat(9.0)),
		High:      models.NewJSONDecimal(decimal.NewFromFloat(9.2)),
		Low:       models.NewJSONDecimal(decimal.NewFromFloat(8.8)),
		Close:     models.NewJSONDecimal(decimal.NewFromFloat(9.0)), // 价格低于平均价格的95%
		Vol:       models.NewJSONDecimal(decimal.NewFromFloat(800)), // 成交量低于平均的50%
	}

	patterns := recognizer.RecognizeAllPatterns(data)
	found := false
	for _, pattern := range patterns {
		for _, volumePrice := range pattern.VolumePrice {
			if volumePrice.Pattern == "地量地价" {
				found = true
				if volumePrice.Signal != "BUY" {
					t.Errorf("地量地价应该是买入信号，实际是: %s", volumePrice.Signal)
				}
				break
			}
		}
	}

	if !found {
		t.Error("未找到地量地价模式")
	}
}

// TestPatternRecognizer_RecognizeHighVolumePrice 测试天量天价模式识别
func TestPatternRecognizer_RecognizeHighVolumePrice(t *testing.T) {
	recognizer := NewPatternRecognizer()

	// 创建测试数据：天量天价模式
	// 需要至少20天的历史数据
	data := make([]models.StockDaily, 21)
	
	// 填充前20天的数据（作为基准）
	for i := 0; i < 20; i++ {
		data[i] = models.StockDaily{
			TSCode:    "000001.SZ",
			TradeDate: "2024010" + string(rune('1'+i%10)),
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.5)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000)), // 平均成交量
		}
	}
	
	// 最后一天：天量天价
	data[20] = models.StockDaily{
		TSCode:    "000001.SZ",
		TradeDate: "20240121",
		Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.4)),
		High:      models.NewJSONDecimal(decimal.NewFromFloat(11.0)),
		Low:       models.NewJSONDecimal(decimal.NewFromFloat(10.3)),
		Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.8)), // 价格接近最高价的95%以上
		Vol:       models.NewJSONDecimal(decimal.NewFromFloat(3000)), // 成交量是平均的3倍
	}

	patterns := recognizer.RecognizeAllPatterns(data)
	found := false
	for _, pattern := range patterns {
		for _, volumePrice := range pattern.VolumePrice {
			if volumePrice.Pattern == "天量天价" {
				found = true
				if volumePrice.Signal != "SELL" {
					t.Errorf("天量天价应该是卖出信号，实际是: %s", volumePrice.Signal)
				}
				break
			}
		}
	}

	if !found {
		t.Error("未找到天量天价模式")
	}
}
