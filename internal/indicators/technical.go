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
			sum = sum.Add(data[j].Close)
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
			DIF:    difValues[i],
			DEA:    deaValues[i],
			MACD:   macd,
			Signal: signal,
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
		change := data[i].Close.Sub(data[i-1].Close)
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
			RSI6:   rsi, // 简化为单一RSI值
			RSI12:  rsi,
			RSI24:  rsi,
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
			diff := data[j].Close.Sub(mas[maIndex])
			sum = sum.Add(diff.Mul(diff))
		}
		variance := sum.Div(decimal.NewFromInt(int64(period)))
		stdDev := decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))

		middle := mas[maIndex]
		offset := stdDev.Mul(decimal.NewFromFloat(multiplier))
		upper := middle.Add(offset)
		lower := middle.Sub(offset)

		signal := "HOLD"
		currentPrice := data[i].Close
		if currentPrice.GreaterThan(upper) {
			signal = "SELL" // 价格突破上轨，可能回调
		} else if currentPrice.LessThan(lower) {
			signal = "BUY" // 价格突破下轨，可能反弹
		}

		bollResults = append(bollResults, models.BollingerBandsIndicator{
			Upper:  upper,
			Middle: middle,
			Lower:  lower,
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
		high := data[i-period+1].High
		low := data[i-period+1].Low

		for j := i - period + 2; j <= i; j++ {
			if data[j].High.GreaterThan(high) {
				high = data[j].High
			}
			if data[j].Low.LessThan(low) {
				low = data[j].Low
			}
		}

		// 计算RSV
		var rsv decimal.Decimal
		if high.Equal(low) {
			rsv = decimal.NewFromInt(50) // 避免除零
		} else {
			rsv = data[i].Close.Sub(low).Div(high.Sub(low)).Mul(decimal.NewFromInt(100))
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
			K:      k,
			D:      d,
			J:      j,
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
		sum = sum.Add(data[i].Close)
	}
	firstEMA := sum.Div(decimal.NewFromInt(int64(period)))
	emas = append(emas, firstEMA)

	// 后续EMA值
	for i := period; i < len(data); i++ {
		ema := data[i].Close.Mul(alpha).Add(emas[len(emas)-1].Mul(decimal.NewFromInt(1).Sub(alpha)))
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
