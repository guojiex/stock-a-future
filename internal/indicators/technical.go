package indicators

import (
	"math"
	"stock-a-future/internal/models"

	"github.com/shopspring/decimal"
)

// Calculator 技术指标计算器
type Calculator struct{}

// NewCalculator 创建新的技术指标计算器
func NewCalculator() *Calculator {
	return &Calculator{}
}

// CalculateMA 计算移动平均线
func (c *Calculator) CalculateMA(data []models.StockDaily, period int) []decimal.Decimal {
	if len(data) < period {
		return []decimal.Decimal{}
	}

	var mas []decimal.Decimal
	for i := period - 1; i < len(data); i++ {
		sum := decimal.Zero
		for j := i - period + 1; j <= i; j++ {
			sum = sum.Add(data[j].Close.Decimal)
		}
		ma := sum.Div(decimal.NewFromInt(int64(period)))
		mas = append(mas, ma)
	}
	return mas
}

// CalculateMACD 计算MACD指标
func (c *Calculator) CalculateMACD(data []models.StockDaily) []models.MACDIndicator {
	if len(data) < 26 {
		return []models.MACDIndicator{}
	}

	// 计算EMA12和EMA26
	ema12 := c.calculateEMA(data, 12)
	ema26 := c.calculateEMA(data, 26)

	var macdResults []models.MACDIndicator
	var difValues []decimal.Decimal

	// 计算DIF线
	for i := 0; i < len(ema12) && i < len(ema26); i++ {
		dif := ema12[i].Sub(ema26[i])
		difValues = append(difValues, dif)
	}

	// 计算DEA线(DIF的9日EMA)
	deaValues := c.calculateEMAFromValues(difValues, 9)

	// 计算MACD柱状图和信号
	for i := 0; i < len(difValues) && i < len(deaValues); i++ {
		macd := difValues[i].Sub(deaValues[i]).Mul(decimal.NewFromInt(2))

		signal := "HOLD"
		if i > 0 {
			prevMACD := difValues[i-1].Sub(deaValues[i-1]).Mul(decimal.NewFromInt(2))
			// 金叉：MACD由负转正
			if prevMACD.LessThan(decimal.Zero) && macd.GreaterThan(decimal.Zero) {
				signal = "BUY"
			}
			// 死叉：MACD由正转负
			if prevMACD.GreaterThan(decimal.Zero) && macd.LessThan(decimal.Zero) {
				signal = "SELL"
			}
		}

		macdResults = append(macdResults, models.MACDIndicator{
			DIF:       models.NewJSONDecimal(difValues[i]),
			DEA:       models.NewJSONDecimal(deaValues[i]),
			Histogram: models.NewJSONDecimal(macd),
			Signal:    signal,
		})
	}

	return macdResults
}

// CalculateRSI 计算RSI指标
func (c *Calculator) CalculateRSI(data []models.StockDaily, period int) []models.RSIIndicator {
	if len(data) < period+1 {
		return []models.RSIIndicator{}
	}

	var gains, losses []decimal.Decimal

	// 计算涨跌幅
	for i := 1; i < len(data); i++ {
		change := data[i].Close.Decimal.Sub(data[i-1].Close.Decimal)
		if change.GreaterThan(decimal.Zero) {
			gains = append(gains, change)
			losses = append(losses, decimal.Zero)
		} else {
			gains = append(gains, decimal.Zero)
			losses = append(losses, change.Abs())
		}
	}

	var rsiResults []models.RSIIndicator

	// 计算RSI
	for i := period - 1; i < len(gains); i++ {
		avgGain := c.averageOfSlice(gains[i-period+1 : i+1])
		avgLoss := c.averageOfSlice(losses[i-period+1 : i+1])

		var rsi decimal.Decimal
		if avgLoss.IsZero() {
			rsi = decimal.NewFromInt(100)
		} else {
			rs := avgGain.Div(avgLoss)
			rsi = decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).Div(decimal.NewFromInt(1).Add(rs)))
		}

		signal := "HOLD"
		if rsi.GreaterThan(decimal.NewFromInt(70)) {
			signal = "SELL" // 超买
		} else if rsi.LessThan(decimal.NewFromInt(30)) {
			signal = "BUY" // 超卖
		}

		// 这里简化处理，只计算一个周期的RSI
		rsiResults = append(rsiResults, models.RSIIndicator{
			RSI14:  models.NewJSONDecimal(rsi),
			Signal: signal,
		})
	}

	return rsiResults
}

// CalculateBollingerBands 计算布林带
func (c *Calculator) CalculateBollingerBands(data []models.StockDaily, period int, multiplier float64) []models.BollingerBandsIndicator {
	if len(data) < period {
		return []models.BollingerBandsIndicator{}
	}

	mas := c.CalculateMA(data, period)
	var bollResults []models.BollingerBandsIndicator

	for i := period - 1; i < len(data); i++ {
		maIndex := i - period + 1
		if maIndex >= len(mas) {
			break
		}

		// 计算标准差
		sum := decimal.Zero
		for j := i - period + 1; j <= i; j++ {
			diff := data[j].Close.Decimal.Sub(mas[maIndex])
			sum = sum.Add(diff.Mul(diff))
		}
		variance := sum.Div(decimal.NewFromInt(int64(period)))
		stdDev := decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))

		middle := mas[maIndex]
		offset := stdDev.Mul(decimal.NewFromFloat(multiplier))
		upper := middle.Add(offset)
		lower := middle.Sub(offset)

		signal := "HOLD"
		currentPrice := data[i].Close.Decimal
		if currentPrice.GreaterThan(upper) {
			signal = "SELL" // 价格突破上轨，可能回调
		} else if currentPrice.LessThan(lower) {
			signal = "BUY" // 价格突破下轨，可能反弹
		}

		bollResults = append(bollResults, models.BollingerBandsIndicator{
			Upper:  models.NewJSONDecimal(upper),
			Middle: models.NewJSONDecimal(middle),
			Lower:  models.NewJSONDecimal(lower),
			Signal: signal,
		})
	}

	return bollResults
}

// CalculateKDJ 计算KDJ指标
func (c *Calculator) CalculateKDJ(data []models.StockDaily, period int) []models.KDJIndicator {
	if len(data) < period {
		return []models.KDJIndicator{}
	}

	var kdjResults []models.KDJIndicator
	var kValues, dValues []decimal.Decimal

	for i := period - 1; i < len(data); i++ {
		// 计算最高价和最低价
		high := data[i-period+1].High.Decimal
		low := data[i-period+1].Low.Decimal

		for j := i - period + 2; j <= i; j++ {
			if data[j].High.Decimal.GreaterThan(high) {
				high = data[j].High.Decimal
			}
			if data[j].Low.Decimal.LessThan(low) {
				low = data[j].Low.Decimal
			}
		}

		// 计算RSV
		var rsv decimal.Decimal
		if high.Equal(low) {
			rsv = decimal.NewFromInt(50) // 避免除零
		} else {
			rsv = data[i].Close.Decimal.Sub(low).Div(high.Sub(low)).Mul(decimal.NewFromInt(100))
		}

		// 计算K值
		var k decimal.Decimal
		if len(kValues) == 0 {
			k = rsv
		} else {
			k = kValues[len(kValues)-1].Mul(decimal.NewFromFloat(2.0 / 3.0)).Add(rsv.Mul(decimal.NewFromFloat(1.0 / 3.0)))
		}
		kValues = append(kValues, k)

		// 计算D值
		var d decimal.Decimal
		if len(dValues) == 0 {
			d = k
		} else {
			d = dValues[len(dValues)-1].Mul(decimal.NewFromFloat(2.0 / 3.0)).Add(k.Mul(decimal.NewFromFloat(1.0 / 3.0)))
		}
		dValues = append(dValues, d)

		// 计算J值
		j := k.Mul(decimal.NewFromInt(3)).Sub(d.Mul(decimal.NewFromInt(2)))

		signal := "HOLD"
		if k.GreaterThan(decimal.NewFromInt(80)) && d.GreaterThan(decimal.NewFromInt(80)) {
			signal = "SELL" // 超买
		} else if k.LessThan(decimal.NewFromInt(20)) && d.LessThan(decimal.NewFromInt(20)) {
			signal = "BUY" // 超卖
		}

		kdjResults = append(kdjResults, models.KDJIndicator{
			K:      models.NewJSONDecimal(k),
			D:      models.NewJSONDecimal(d),
			J:      models.NewJSONDecimal(j),
			Signal: signal,
		})
	}

	return kdjResults
}

// calculateEMA 计算指数移动平均线
func (c *Calculator) calculateEMA(data []models.StockDaily, period int) []decimal.Decimal {
	if len(data) < period {
		return []decimal.Decimal{}
	}

	alpha := decimal.NewFromFloat(2.0 / float64(period+1))
	var emas []decimal.Decimal

	// 第一个EMA值使用简单移动平均
	sum := decimal.Zero
	for i := 0; i < period; i++ {
		sum = sum.Add(data[i].Close.Decimal)
	}
	firstEMA := sum.Div(decimal.NewFromInt(int64(period)))
	emas = append(emas, firstEMA)

	// 后续EMA值
	for i := period; i < len(data); i++ {
		ema := data[i].Close.Decimal.Mul(alpha).Add(emas[len(emas)-1].Mul(decimal.NewFromInt(1).Sub(alpha)))
		emas = append(emas, ema)
	}

	return emas
}

// calculateEMAFromValues 从数值数组计算EMA
func (c *Calculator) calculateEMAFromValues(values []decimal.Decimal, period int) []decimal.Decimal {
	if len(values) < period {
		return []decimal.Decimal{}
	}

	alpha := decimal.NewFromFloat(2.0 / float64(period+1))
	var emas []decimal.Decimal

	// 第一个EMA值使用简单移动平均
	sum := decimal.Zero
	for i := 0; i < period; i++ {
		sum = sum.Add(values[i])
	}
	firstEMA := sum.Div(decimal.NewFromInt(int64(period)))
	emas = append(emas, firstEMA)

	// 后续EMA值
	for i := period; i < len(values); i++ {
		ema := values[i].Mul(alpha).Add(emas[len(emas)-1].Mul(decimal.NewFromInt(1).Sub(alpha)))
		emas = append(emas, ema)
	}

	return emas
}

// averageOfSlice 计算切片的平均值
func (c *Calculator) averageOfSlice(values []decimal.Decimal) decimal.Decimal {
	if len(values) == 0 {
		return decimal.Zero
	}

	sum := decimal.Zero
	for _, v := range values {
		sum = sum.Add(v)
	}
	return sum.Div(decimal.NewFromInt(int64(len(values))))
}

// ===== 动量因子 =====

// CalculateWilliamsR 计算威廉指标 (%R)
func (c *Calculator) CalculateWilliamsR(data []models.StockDaily, period int) []models.WilliamsRIndicator {
	if len(data) < period {
		return []models.WilliamsRIndicator{}
	}

	var results []models.WilliamsRIndicator

	for i := period - 1; i < len(data); i++ {
		// 计算周期内最高价和最低价
		high := data[i-period+1].High.Decimal
		low := data[i-period+1].Low.Decimal

		for j := i - period + 2; j <= i; j++ {
			if data[j].High.Decimal.GreaterThan(high) {
				high = data[j].High.Decimal
			}
			if data[j].Low.Decimal.LessThan(low) {
				low = data[j].Low.Decimal
			}
		}

		// 计算威廉指标
		var wr decimal.Decimal
		if high.Equal(low) {
			wr = decimal.NewFromInt(-50) // 避免除零
		} else {
			wr = high.Sub(data[i].Close.Decimal).Div(high.Sub(low)).Mul(decimal.NewFromInt(-100))
		}

		signal := "HOLD"
		if wr.GreaterThan(decimal.NewFromInt(-20)) {
			signal = "SELL" // 超买
		} else if wr.LessThan(decimal.NewFromInt(-80)) {
			signal = "BUY" // 超卖
		}

		results = append(results, models.WilliamsRIndicator{
			WR14:   models.NewJSONDecimal(wr),
			Signal: signal,
		})
	}

	return results
}

// CalculateMomentum 计算动量指标
func (c *Calculator) CalculateMomentum(data []models.StockDaily) []models.MomentumIndicator {
	if len(data) < 20 {
		return []models.MomentumIndicator{}
	}

	var results []models.MomentumIndicator

	for i := 19; i < len(data); i++ {
		// 10日动量
		momentum10 := data[i].Close.Decimal.Sub(data[i-9].Close.Decimal)
		// 20日动量
		momentum20 := data[i].Close.Decimal.Sub(data[i-19].Close.Decimal)

		signal := "HOLD"
		if momentum10.GreaterThan(decimal.Zero) && momentum20.GreaterThan(decimal.Zero) {
			signal = "BUY" // 双重正动量
		} else if momentum10.LessThan(decimal.Zero) && momentum20.LessThan(decimal.Zero) {
			signal = "SELL" // 双重负动量
		}

		results = append(results, models.MomentumIndicator{
			Momentum10: models.NewJSONDecimal(momentum10),
			Momentum20: models.NewJSONDecimal(momentum20),
			Signal:     signal,
		})
	}

	return results
}

// CalculateROC 计算变化率指标
func (c *Calculator) CalculateROC(data []models.StockDaily) []models.ROCIndicator {
	if len(data) < 20 {
		return []models.ROCIndicator{}
	}

	var results []models.ROCIndicator

	for i := 19; i < len(data); i++ {
		// 10日变化率
		roc10 := data[i].Close.Decimal.Sub(data[i-9].Close.Decimal).Div(data[i-9].Close.Decimal).Mul(decimal.NewFromInt(100))
		// 20日变化率
		roc20 := data[i].Close.Decimal.Sub(data[i-19].Close.Decimal).Div(data[i-19].Close.Decimal).Mul(decimal.NewFromInt(100))

		signal := "HOLD"
		if roc10.GreaterThan(decimal.NewFromInt(5)) && roc20.GreaterThan(decimal.NewFromInt(10)) {
			signal = "BUY" // 强势上涨
		} else if roc10.LessThan(decimal.NewFromInt(-5)) && roc20.LessThan(decimal.NewFromInt(-10)) {
			signal = "SELL" // 强势下跌
		}

		results = append(results, models.ROCIndicator{
			ROC10:  models.NewJSONDecimal(roc10),
			ROC20:  models.NewJSONDecimal(roc20),
			Signal: signal,
		})
	}

	return results
}

// ===== 趋势因子 =====

// CalculateADX 计算平均方向指数
func (c *Calculator) CalculateADX(data []models.StockDaily, period int) []models.ADXIndicator {
	if len(data) < period+1 {
		return []models.ADXIndicator{}
	}

	var results []models.ADXIndicator
	var trueRanges, plusDMs, minusDMs []decimal.Decimal

	// 计算TR、+DM、-DM
	for i := 1; i < len(data); i++ {
		// 真实范围 TR
		tr1 := data[i].High.Decimal.Sub(data[i].Low.Decimal)
		tr2 := data[i].High.Decimal.Sub(data[i-1].Close.Decimal).Abs()
		tr3 := data[i].Low.Decimal.Sub(data[i-1].Close.Decimal).Abs()
		tr := decimal.Max(tr1, decimal.Max(tr2, tr3))
		trueRanges = append(trueRanges, tr)

		// +DM 和 -DM
		upMove := data[i].High.Decimal.Sub(data[i-1].High.Decimal)
		downMove := data[i-1].Low.Decimal.Sub(data[i].Low.Decimal)

		var plusDM, minusDM decimal.Decimal
		if upMove.GreaterThan(downMove) && upMove.GreaterThan(decimal.Zero) {
			plusDM = upMove
		}
		if downMove.GreaterThan(upMove) && downMove.GreaterThan(decimal.Zero) {
			minusDM = downMove
		}

		plusDMs = append(plusDMs, plusDM)
		minusDMs = append(minusDMs, minusDM)
	}

	// 计算平滑的TR、+DM、-DM
	for i := period - 1; i < len(trueRanges); i++ {
		avgTR := c.averageOfSlice(trueRanges[i-period+1 : i+1])
		avgPlusDM := c.averageOfSlice(plusDMs[i-period+1 : i+1])
		avgMinusDM := c.averageOfSlice(minusDMs[i-period+1 : i+1])

		// 计算+DI和-DI
		var pdi, mdi decimal.Decimal
		if !avgTR.IsZero() {
			pdi = avgPlusDM.Div(avgTR).Mul(decimal.NewFromInt(100))
			mdi = avgMinusDM.Div(avgTR).Mul(decimal.NewFromInt(100))
		}

		// 计算ADX
		var adx decimal.Decimal
		diSum := pdi.Add(mdi)
		if !diSum.IsZero() {
			dx := pdi.Sub(mdi).Abs().Div(diSum).Mul(decimal.NewFromInt(100))
			adx = dx // 简化处理，实际应该是DX的移动平均
		}

		signal := "HOLD"
		if adx.GreaterThan(decimal.NewFromInt(25)) {
			if pdi.GreaterThan(mdi) {
				signal = "BUY" // 强势上涨趋势
			} else {
				signal = "SELL" // 强势下跌趋势
			}
		}

		results = append(results, models.ADXIndicator{
			ADX:    models.NewJSONDecimal(adx),
			PDI:    models.NewJSONDecimal(pdi),
			MDI:    models.NewJSONDecimal(mdi),
			Signal: signal,
		})
	}

	return results
}

// CalculateSAR 计算抛物线转向指标
func (c *Calculator) CalculateSAR(data []models.StockDaily) []models.SARIndicator {
	if len(data) < 5 {
		return []models.SARIndicator{}
	}

	var results []models.SARIndicator
	af := decimal.NewFromFloat(0.02) // 加速因子
	maxAF := decimal.NewFromFloat(0.2)

	// 初始化
	isUpTrend := data[1].Close.Decimal.GreaterThan(data[0].Close.Decimal)
	var sar, ep decimal.Decimal

	if isUpTrend {
		sar = data[0].Low.Decimal
		ep = data[1].High.Decimal
	} else {
		sar = data[0].High.Decimal
		ep = data[1].Low.Decimal
	}

	for i := 1; i < len(data); i++ {
		// 计算新的SAR
		newSAR := sar.Add(af.Mul(ep.Sub(sar)))

		// 检查趋势反转
		var trendReversed bool
		if isUpTrend {
			if data[i].Low.Decimal.LessThan(newSAR) {
				trendReversed = true
				isUpTrend = false
				newSAR = ep
				ep = data[i].Low.Decimal
				af = decimal.NewFromFloat(0.02)
			} else {
				if data[i].High.Decimal.GreaterThan(ep) {
					ep = data[i].High.Decimal
					af = decimal.Min(af.Add(decimal.NewFromFloat(0.02)), maxAF)
				}
			}
		} else {
			if data[i].High.Decimal.GreaterThan(newSAR) {
				trendReversed = true
				isUpTrend = true
				newSAR = ep
				ep = data[i].High.Decimal
				af = decimal.NewFromFloat(0.02)
			} else {
				if data[i].Low.Decimal.LessThan(ep) {
					ep = data[i].Low.Decimal
					af = decimal.Min(af.Add(decimal.NewFromFloat(0.02)), maxAF)
				}
			}
		}

		sar = newSAR

		signal := "HOLD"
		if trendReversed {
			if isUpTrend {
				signal = "BUY"
			} else {
				signal = "SELL"
			}
		}

		results = append(results, models.SARIndicator{
			SAR:    models.NewJSONDecimal(sar),
			Signal: signal,
		})
	}

	return results
}

// CalculateIchimoku 计算一目均衡表
func (c *Calculator) CalculateIchimoku(data []models.StockDaily) []models.IchimokuIndicator {
	if len(data) < 52 {
		return []models.IchimokuIndicator{}
	}

	var results []models.IchimokuIndicator

	for i := 51; i < len(data); i++ {
		// 转换线 (9日最高价+最低价)/2
		tenkanHigh := data[i-8].High.Decimal
		tenkanLow := data[i-8].Low.Decimal
		for j := i - 7; j <= i; j++ {
			if data[j].High.Decimal.GreaterThan(tenkanHigh) {
				tenkanHigh = data[j].High.Decimal
			}
			if data[j].Low.Decimal.LessThan(tenkanLow) {
				tenkanLow = data[j].Low.Decimal
			}
		}
		tenkanSen := tenkanHigh.Add(tenkanLow).Div(decimal.NewFromInt(2))

		// 基准线 (26日最高价+最低价)/2
		kijunHigh := data[i-25].High.Decimal
		kijunLow := data[i-25].Low.Decimal
		for j := i - 24; j <= i; j++ {
			if data[j].High.Decimal.GreaterThan(kijunHigh) {
				kijunHigh = data[j].High.Decimal
			}
			if data[j].Low.Decimal.LessThan(kijunLow) {
				kijunLow = data[j].Low.Decimal
			}
		}
		kijunSen := kijunHigh.Add(kijunLow).Div(decimal.NewFromInt(2))

		// 先行带A (转换线+基准线)/2
		senkouSpanA := tenkanSen.Add(kijunSen).Div(decimal.NewFromInt(2))

		// 先行带B (52日最高价+最低价)/2
		senkouHigh := data[i-51].High.Decimal
		senkouLow := data[i-51].Low.Decimal
		for j := i - 50; j <= i; j++ {
			if data[j].High.Decimal.GreaterThan(senkouHigh) {
				senkouHigh = data[j].High.Decimal
			}
			if data[j].Low.Decimal.LessThan(senkouLow) {
				senkouLow = data[j].Low.Decimal
			}
		}
		senkouSpanB := senkouHigh.Add(senkouLow).Div(decimal.NewFromInt(2))

		// 滞后线 (当前收盘价)
		chikouSpan := data[i].Close.Decimal

		signal := "HOLD"
		currentPrice := data[i].Close.Decimal
		if tenkanSen.GreaterThan(kijunSen) && currentPrice.GreaterThan(senkouSpanA) && currentPrice.GreaterThan(senkouSpanB) {
			signal = "BUY" // 多头排列
		} else if tenkanSen.LessThan(kijunSen) && currentPrice.LessThan(senkouSpanA) && currentPrice.LessThan(senkouSpanB) {
			signal = "SELL" // 空头排列
		}

		results = append(results, models.IchimokuIndicator{
			TenkanSen:   models.NewJSONDecimal(tenkanSen),
			KijunSen:    models.NewJSONDecimal(kijunSen),
			SenkouSpanA: models.NewJSONDecimal(senkouSpanA),
			SenkouSpanB: models.NewJSONDecimal(senkouSpanB),
			ChikouSpan:  models.NewJSONDecimal(chikouSpan),
			Signal:      signal,
		})
	}

	return results
}

// ===== 波动率因子 =====

// CalculateATR 计算平均真实范围
func (c *Calculator) CalculateATR(data []models.StockDaily, period int) []models.ATRIndicator {
	if len(data) < period+1 {
		return []models.ATRIndicator{}
	}

	var results []models.ATRIndicator
	var trueRanges []decimal.Decimal

	// 计算真实范围
	for i := 1; i < len(data); i++ {
		tr1 := data[i].High.Decimal.Sub(data[i].Low.Decimal)
		tr2 := data[i].High.Decimal.Sub(data[i-1].Close.Decimal).Abs()
		tr3 := data[i].Low.Decimal.Sub(data[i-1].Close.Decimal).Abs()
		tr := decimal.Max(tr1, decimal.Max(tr2, tr3))
		trueRanges = append(trueRanges, tr)
	}

	// 计算ATR
	for i := period - 1; i < len(trueRanges); i++ {
		atr := c.averageOfSlice(trueRanges[i-period+1 : i+1])

		signal := "HOLD"
		// ATR主要用于衡量波动率，这里简化处理
		if i >= period {
			prevATR := c.averageOfSlice(trueRanges[i-period : i])
			if atr.GreaterThan(prevATR.Mul(decimal.NewFromFloat(1.2))) {
				signal = "HIGH_VOLATILITY" // 高波动率
			} else if atr.LessThan(prevATR.Mul(decimal.NewFromFloat(0.8))) {
				signal = "LOW_VOLATILITY" // 低波动率
			}
		}

		results = append(results, models.ATRIndicator{
			ATR14:  models.NewJSONDecimal(atr),
			Signal: signal,
		})
	}

	return results
}

// CalculateStdDev 计算标准差
func (c *Calculator) CalculateStdDev(data []models.StockDaily, period int) []models.StdDevIndicator {
	if len(data) < period {
		return []models.StdDevIndicator{}
	}

	var results []models.StdDevIndicator

	for i := period - 1; i < len(data); i++ {
		// 计算平均值
		sum := decimal.Zero
		for j := i - period + 1; j <= i; j++ {
			sum = sum.Add(data[j].Close.Decimal)
		}
		mean := sum.Div(decimal.NewFromInt(int64(period)))

		// 计算方差
		varianceSum := decimal.Zero
		for j := i - period + 1; j <= i; j++ {
			diff := data[j].Close.Decimal.Sub(mean)
			varianceSum = varianceSum.Add(diff.Mul(diff))
		}
		variance := varianceSum.Div(decimal.NewFromInt(int64(period)))
		stdDev := decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))

		signal := "HOLD"
		// 标准差反映价格波动程度
		if stdDev.GreaterThan(mean.Mul(decimal.NewFromFloat(0.05))) {
			signal = "HIGH_VOLATILITY"
		} else if stdDev.LessThan(mean.Mul(decimal.NewFromFloat(0.02))) {
			signal = "LOW_VOLATILITY"
		}

		results = append(results, models.StdDevIndicator{
			StdDev20: models.NewJSONDecimal(stdDev),
			Signal:   signal,
		})
	}

	return results
}

// CalculateHistoricalVolatility 计算历史波动率
func (c *Calculator) CalculateHistoricalVolatility(data []models.StockDaily) []models.HistoricalVolatilityIndicator {
	if len(data) < 60 {
		return []models.HistoricalVolatilityIndicator{}
	}

	var results []models.HistoricalVolatilityIndicator

	for i := 59; i < len(data); i++ {
		// 计算20日历史波动率
		var returns20 []decimal.Decimal
		for j := i - 19; j <= i; j++ {
			if j > 0 {
				ret := data[j].Close.Decimal.Div(data[j-1].Close.Decimal).Sub(decimal.NewFromInt(1))
				returns20 = append(returns20, ret)
			}
		}

		// 计算60日历史波动率
		var returns60 []decimal.Decimal
		for j := i - 59; j <= i; j++ {
			if j > 0 {
				ret := data[j].Close.Decimal.Div(data[j-1].Close.Decimal).Sub(decimal.NewFromInt(1))
				returns60 = append(returns60, ret)
			}
		}

		// 计算标准差并年化
		hv20 := c.calculateVolatility(returns20).Mul(decimal.NewFromFloat(math.Sqrt(252)))
		hv60 := c.calculateVolatility(returns60).Mul(decimal.NewFromFloat(math.Sqrt(252)))

		signal := "HOLD"
		if hv20.GreaterThan(decimal.NewFromFloat(0.3)) {
			signal = "HIGH_VOLATILITY"
		} else if hv20.LessThan(decimal.NewFromFloat(0.15)) {
			signal = "LOW_VOLATILITY"
		}

		results = append(results, models.HistoricalVolatilityIndicator{
			HV20:   models.NewJSONDecimal(hv20),
			HV60:   models.NewJSONDecimal(hv60),
			Signal: signal,
		})
	}

	return results
}

// calculateVolatility 计算收益率序列的波动率
func (c *Calculator) calculateVolatility(returns []decimal.Decimal) decimal.Decimal {
	if len(returns) == 0 {
		return decimal.Zero
	}

	// 计算平均收益率
	mean := c.averageOfSlice(returns)

	// 计算方差
	varianceSum := decimal.Zero
	for _, ret := range returns {
		diff := ret.Sub(mean)
		varianceSum = varianceSum.Add(diff.Mul(diff))
	}

	variance := varianceSum.Div(decimal.NewFromInt(int64(len(returns))))
	return decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))
}

// ===== 成交量因子 =====

// CalculateVWAP 计算成交量加权平均价
func (c *Calculator) CalculateVWAP(data []models.StockDaily) []models.VWAPIndicator {
	if len(data) < 1 {
		return []models.VWAPIndicator{}
	}

	var results []models.VWAPIndicator
	var cumulativeVolumePrice, cumulativeVolume decimal.Decimal

	for i := 0; i < len(data); i++ {
		// 典型价格 (最高价+最低价+收盘价)/3
		typicalPrice := data[i].High.Decimal.Add(data[i].Low.Decimal).Add(data[i].Close.Decimal).Div(decimal.NewFromInt(3))
		volumePrice := typicalPrice.Mul(data[i].Vol.Decimal)

		cumulativeVolumePrice = cumulativeVolumePrice.Add(volumePrice)
		cumulativeVolume = cumulativeVolume.Add(data[i].Vol.Decimal)

		var vwap decimal.Decimal
		if !cumulativeVolume.IsZero() {
			vwap = cumulativeVolumePrice.Div(cumulativeVolume)
		}

		signal := "HOLD"
		currentPrice := data[i].Close.Decimal
		if currentPrice.GreaterThan(vwap) {
			signal = "BUY" // 价格高于VWAP
		} else if currentPrice.LessThan(vwap) {
			signal = "SELL" // 价格低于VWAP
		}

		results = append(results, models.VWAPIndicator{
			VWAP:   models.NewJSONDecimal(vwap),
			Signal: signal,
		})

		// 重置累积值（这里简化为每日重置，实际应该根据需要调整）
		if i > 0 && i%20 == 0 { // 每20天重置一次
			cumulativeVolumePrice = decimal.Zero
			cumulativeVolume = decimal.Zero
		}
	}

	return results
}

// CalculateADLine 计算累积/派发线
func (c *Calculator) CalculateADLine(data []models.StockDaily) []models.ADLineIndicator {
	if len(data) < 1 {
		return []models.ADLineIndicator{}
	}

	var results []models.ADLineIndicator
	var adLine decimal.Decimal

	for i := 0; i < len(data); i++ {
		// 计算资金流量乘数
		var moneyFlowMultiplier decimal.Decimal
		highLowDiff := data[i].High.Decimal.Sub(data[i].Low.Decimal)
		if !highLowDiff.IsZero() {
			moneyFlowMultiplier = data[i].Close.Decimal.Sub(data[i].Low.Decimal).Sub(data[i].High.Decimal.Sub(data[i].Close.Decimal)).Div(highLowDiff)
		}

		// 计算资金流量
		moneyFlowVolume := moneyFlowMultiplier.Mul(data[i].Vol.Decimal)
		adLine = adLine.Add(moneyFlowVolume)

		signal := "HOLD"
		if len(results) > 0 && i > 0 {
			prevADLine := results[len(results)-1].ADLine.Decimal
			if adLine.GreaterThan(prevADLine) && data[i].Close.Decimal.GreaterThan(data[i-1].Close.Decimal) {
				signal = "BUY" // 价格和A/D线同时上涨
			} else if adLine.LessThan(prevADLine) && data[i].Close.Decimal.LessThan(data[i-1].Close.Decimal) {
				signal = "SELL" // 价格和A/D线同时下跌
			}
		}

		results = append(results, models.ADLineIndicator{
			ADLine: models.NewJSONDecimal(adLine),
			Signal: signal,
		})
	}

	return results
}

// CalculateEMV 计算简易波动指标
func (c *Calculator) CalculateEMV(data []models.StockDaily, period int) []models.EMVIndicator {
	if len(data) < period+1 {
		return []models.EMVIndicator{}
	}

	var results []models.EMVIndicator
	var emvValues []decimal.Decimal

	for i := 1; i < len(data); i++ {
		// 计算距离移动
		distanceMove := data[i].High.Decimal.Add(data[i].Low.Decimal).Div(decimal.NewFromInt(2)).Sub(
			data[i-1].High.Decimal.Add(data[i-1].Low.Decimal).Div(decimal.NewFromInt(2)))

		// 计算EMV
		var emv decimal.Decimal

		// 检查最高价和最低价是否相等（避免除以零）
		priceRange := data[i].High.Decimal.Sub(data[i].Low.Decimal)
		if !priceRange.IsZero() && !data[i].Vol.Decimal.IsZero() {
			// 计算盒子高度
			boxHeight := data[i].Vol.Decimal.Div(priceRange)
			if !boxHeight.IsZero() {
				emv = distanceMove.Div(boxHeight)
			}
		}
		// 如果价格区间为零或成交量为零，EMV保持为零

		emvValues = append(emvValues, emv)
	}

	// 计算EMV的移动平均
	for i := period - 1; i < len(emvValues); i++ {
		emvMA := c.averageOfSlice(emvValues[i-period+1 : i+1])

		signal := "HOLD"
		if emvMA.GreaterThan(decimal.Zero) {
			signal = "BUY" // 正EMV表示上涨容易
		} else if emvMA.LessThan(decimal.Zero) {
			signal = "SELL" // 负EMV表示下跌容易
		}

		results = append(results, models.EMVIndicator{
			EMV14:  models.NewJSONDecimal(emvMA),
			Signal: signal,
		})
	}

	return results
}

// CalculateVPT 计算量价确认指标
func (c *Calculator) CalculateVPT(data []models.StockDaily) []models.VPTIndicator {
	if len(data) < 2 {
		return []models.VPTIndicator{}
	}

	var results []models.VPTIndicator
	var vpt decimal.Decimal

	for i := 1; i < len(data); i++ {
		// 计算价格变化率
		priceChange := data[i].Close.Decimal.Sub(data[i-1].Close.Decimal).Div(data[i-1].Close.Decimal)

		// 计算VPT
		vpt = vpt.Add(data[i].Vol.Decimal.Mul(priceChange))

		signal := "HOLD"
		if len(results) > 0 {
			prevVPT := results[len(results)-1].VPT.Decimal
			if vpt.GreaterThan(prevVPT) && data[i].Close.Decimal.GreaterThan(data[i-1].Close.Decimal) {
				signal = "BUY" // VPT和价格同时上涨
			} else if vpt.LessThan(prevVPT) && data[i].Close.Decimal.LessThan(data[i-1].Close.Decimal) {
				signal = "SELL" // VPT和价格同时下跌
			}
		}

		results = append(results, models.VPTIndicator{
			VPT:    models.NewJSONDecimal(vpt),
			Signal: signal,
		})
	}

	return results
}
