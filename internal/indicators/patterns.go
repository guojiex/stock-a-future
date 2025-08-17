package indicators

import (
	"log"
	"stock-a-future/internal/models"

	"github.com/shopspring/decimal"
)

// PatternRecognizer 图形识别器
type PatternRecognizer struct{}

// NewPatternRecognizer 创建新的图形识别器
func NewPatternRecognizer() *PatternRecognizer {
	return &PatternRecognizer{}
}

// RecognizeAllPatterns 识别所有图形模式
func (p *PatternRecognizer) RecognizeAllPatterns(data []models.StockDaily) []models.PatternRecognitionResult {
	if len(data) < 3 {
		log.Printf("[模式识别] 数据不足，需要至少3天数据，当前只有 %d 天", len(data))
		return []models.PatternRecognitionResult{}
	}

	log.Printf("[模式识别] 开始识别模式，共 %d 天数据", len(data))
	var results []models.PatternRecognitionResult

	// 从第3个交易日开始识别（需要至少3天数据）
	for i := 2; i < len(data); i++ {
		// 获取当前日期和前两天的数据
		current := data[i]
		prev1 := data[i-1]
		prev2 := data[i-2]

		log.Printf("📅 [模式识别] 分析日期: %s (索引: %d)", current.TradeDate, i)
		log.Printf("🔍 [模式识别] 开始识别各种技术形态...")

		// 识别蜡烛图模式
		candlestickPatterns := p.recognizeCandlestickPatterns(current, prev1, prev2, i, data)

		// 识别量价图形
		volumePricePatterns := p.recognizeVolumePricePatterns(current, prev1, prev2, i, data)

		// 如果有识别到图形，创建结果
		if len(candlestickPatterns) > 0 || len(volumePricePatterns) > 0 {
			log.Printf("🎯 [模式识别] 在日期 %s 成功识别到技术形态:", current.TradeDate)
			log.Printf("   📊 蜡烛图模式: %d 个", len(candlestickPatterns))
			log.Printf("   📈 量价模式: %d 个", len(volumePricePatterns))

			// 计算综合信号和置信度
			combinedSignal, overallConfidence, riskLevel := p.calculateCombinedSignal(
				candlestickPatterns, volumePricePatterns)

			result := models.PatternRecognitionResult{
				TSCode:            current.TSCode,
				TradeDate:         current.TradeDate,
				Candlestick:       candlestickPatterns,
				VolumePrice:       volumePricePatterns,
				CombinedSignal:    combinedSignal,
				OverallConfidence: models.NewJSONDecimal(overallConfidence),
				RiskLevel:         riskLevel,
			}
			results = append(results, result)
		} else {
			log.Printf("[模式识别] 在日期 %s 未识别到任何模式", current.TradeDate)
		}
	}

	log.Printf("[模式识别] 识别完成，共找到 %d 个结果", len(results))
	return results
}

// recognizeCandlestickPatterns 识别蜡烛图模式
func (p *PatternRecognizer) recognizeCandlestickPatterns(current, prev1, prev2 models.StockDaily, index int, data []models.StockDaily) []models.CandlestickPattern {
	var patterns []models.CandlestickPattern

	log.Printf("[蜡烛图识别] 开始识别蜡烛图模式，日期: %s", current.TradeDate)

	// 双响炮模式 - 连续两根大阳线，成交量放大
	if pattern := p.recognizeDoubleCannon(current, prev1, prev2); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	// 红三兵模式 - 连续三根上涨K线
	if pattern := p.recognizeRedThreeSoldiers(current, prev1, prev2); pattern != nil {
		log.Printf("[蜡烛图识别] ✅ 识别到红三兵模式")
		patterns = append(patterns, *pattern)
	}

	// 乌云盖顶模式 - 大阳线后跟大阴线
	if pattern := p.recognizeDarkCloudCover(current, prev1, prev2); pattern != nil {
		log.Printf("[蜡烛图识别] ✅ 识别到乌云盖顶模式")
		patterns = append(patterns, *pattern)
	}

	// 锤子线模式 - 下影线很长的K线
	if pattern := p.recognizeHammer(current); pattern != nil {
		log.Printf("[蜡烛图识别] ✅ 识别到锤子线模式")
		patterns = append(patterns, *pattern)
	}

	// 启明星模式 - 下跌趋势中的反转信号
	if pattern := p.recognizeMorningStar(current, prev1, prev2, index, data); pattern != nil {
		log.Printf("[蜡烛图识别] ✅ 识别到启明星模式")
		patterns = append(patterns, *pattern)
	}

	// 黄昏星模式 - 上涨趋势中的反转信号
	if pattern := p.recognizeEveningStar(current, prev1, prev2, index, data); pattern != nil {
		log.Printf("[蜡烛图识别] ✅ 识别到黄昏星模式")
		patterns = append(patterns, *pattern)
	}

	log.Printf("[蜡烛图识别] 识别完成，共找到 %d 个蜡烛图模式", len(patterns))
	return patterns
}

// recognizeVolumePricePatterns 识别量价图形
func (p *PatternRecognizer) recognizeVolumePricePatterns(current, prev1, prev2 models.StockDaily, index int, data []models.StockDaily) []models.VolumePricePattern {
	var patterns []models.VolumePricePattern

	// 量价齐升模式
	if pattern := p.recognizeVolumePriceRise(current, prev1, prev2); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	// 量价背离模式
	if pattern := p.recognizeVolumePriceDivergence(current, prev1, prev2); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	// 放量突破模式
	if pattern := p.recognizeVolumeBreakout(current, prev1, prev2, index, data); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	return patterns
}

// recognizeDoubleCannon 识别双响炮模式
func (p *PatternRecognizer) recognizeDoubleCannon(current, prev1, _ models.StockDaily) *models.CandlestickPattern {
	// 双响炮：连续两根大阳线，成交量放大
	// 第一根：中阳线或大阳线
	// 第二根：大阳线，成交量明显放大

	// 计算第一根K线的涨幅
	prev1Change := prev1.Close.Decimal.Sub(prev1.Open.Decimal)
	prev1ChangePct := prev1Change.Div(prev1.Open.Decimal).Mul(decimal.NewFromFloat(100))

	// 计算第二根K线的涨幅
	currentChange := current.Close.Decimal.Sub(current.Open.Decimal)
	currentChangePct := currentChange.Div(current.Open.Decimal).Mul(decimal.NewFromFloat(100))

	// 计算成交量比率
	volumeRatio := current.Vol.Decimal.Div(prev1.Vol.Decimal)

	// 逐个检查判断条件
	condition1 := prev1ChangePct.GreaterThan(decimal.NewFromFloat(2.0))
	condition2 := currentChangePct.GreaterThan(decimal.NewFromFloat(3.0))
	condition3 := volumeRatio.GreaterThan(decimal.NewFromFloat(1.5))
	condition4 := prev1Change.GreaterThan(decimal.Zero)
	condition5 := currentChange.GreaterThan(decimal.Zero)

	// 判断条件：
	// 1. 第一根K线涨幅 > 2%
	// 2. 第二根K线涨幅 > 3%
	// 3. 第二根成交量是第一根的1.5倍以上
	// 4. 两根K线都是阳线
	if condition1 && condition2 && condition3 && condition4 && condition5 {
		// 计算置信度
		confidence := p.calculateConfidence(prev1ChangePct, currentChangePct, volumeRatio)
		strength := p.calculateStrength(confidence)

		return &models.CandlestickPattern{
			TSCode:      current.TSCode,
			TradeDate:   current.TradeDate,
			Pattern:     "双响炮",
			Signal:      "BUY",
			Confidence:  models.NewJSONDecimal(confidence),
			Description: "连续两根大阳线，成交量放大，强势上涨信号",
			Strength:    strength,
			Volume:      current.Vol,
			PriceChange: models.NewJSONDecimal(currentChange),
		}
	}

	return nil
}

// recognizeRedThreeSoldiers 识别红三兵模式
func (p *PatternRecognizer) recognizeRedThreeSoldiers(current, prev1, prev2 models.StockDaily) *models.CandlestickPattern {
	// 红三兵：连续三根上涨K线，每根K线都在前一根K线的实体范围内开盘

	// 检查三根K线是否都是阳线
	if current.Close.Decimal.LessThanOrEqual(current.Open.Decimal) ||
		prev1.Close.Decimal.LessThanOrEqual(prev1.Open.Decimal) ||
		prev2.Close.Decimal.LessThanOrEqual(prev2.Open.Decimal) {
		return nil
	}

	// 检查开盘价是否在前一根K线的实体范围内
	if current.Open.Decimal.LessThan(prev1.Open.Decimal) ||
		current.Open.Decimal.GreaterThan(prev1.Close.Decimal) ||
		prev1.Open.Decimal.LessThan(prev2.Open.Decimal) ||
		prev1.Open.Decimal.GreaterThan(prev2.Close.Decimal) {
		return nil
	}

	// 计算三根K线的平均涨幅
	change1 := prev2.Close.Decimal.Sub(prev2.Open.Decimal)
	change2 := prev1.Close.Decimal.Sub(prev1.Open.Decimal)
	change3 := current.Close.Decimal.Sub(current.Open.Decimal)

	avgChange := change1.Add(change2).Add(change3).Div(decimal.NewFromInt(3))
	avgChangePct := avgChange.Div(prev2.Open.Decimal).Mul(decimal.NewFromInt(100))

	// 判断条件：平均涨幅 > 1.5%
	if avgChangePct.GreaterThan(decimal.NewFromFloat(1.5)) {
		confidence := p.calculateConfidence(avgChangePct, decimal.Zero, decimal.Zero)
		strength := p.calculateStrength(confidence)

		return &models.CandlestickPattern{
			TSCode:      current.TSCode,
			TradeDate:   current.TradeDate,
			Pattern:     "红三兵",
			Signal:      "BUY",
			Confidence:  models.NewJSONDecimal(confidence),
			Description: "连续三根上涨K线，稳步上涨信号",
			Strength:    strength,
			Volume:      current.Vol,
			PriceChange: models.NewJSONDecimal(change3),
		}
	}

	return nil
}

// recognizeDarkCloudCover 识别乌云盖顶模式
func (p *PatternRecognizer) recognizeDarkCloudCover(current, prev1, _ models.StockDaily) *models.CandlestickPattern {
	// 乌云盖顶：大阳线后跟大阴线，阴线开盘价高于前一根阳线的最高价，收盘价低于前一根阳线实体的中点

	// 检查前一根是否为大阳线
	prev1Change := prev1.Close.Decimal.Sub(prev1.Open.Decimal)
	prev1ChangePct := prev1Change.Div(prev1.Open.Decimal).Mul(decimal.NewFromInt(100))

	// 检查当前是否为大阴线
	currentChange := current.Close.Decimal.Sub(current.Open.Decimal)
	currentChangePct := currentChange.Div(current.Open.Decimal).Mul(decimal.NewFromInt(100))

	// 判断条件：
	// 1. 前一根K线涨幅 > 2%
	// 2. 当前K线跌幅 > 2%
	// 3. 当前开盘价高于前一根最高价
	// 4. 当前收盘价低于前一根实体中点
	if prev1ChangePct.GreaterThan(decimal.NewFromFloat(2.0)) &&
		currentChangePct.LessThan(decimal.NewFromFloat(-2.0)) &&
		current.Open.Decimal.GreaterThan(prev1.High.Decimal) &&
		current.Close.Decimal.LessThan(prev1.Open.Decimal.Add(prev1Change.Div(decimal.NewFromInt(2)))) {

		confidence := p.calculateConfidence(prev1ChangePct.Abs(), currentChangePct.Abs(), decimal.Zero)
		strength := p.calculateStrength(confidence)

		return &models.CandlestickPattern{
			TSCode:      current.TSCode,
			TradeDate:   current.TradeDate,
			Pattern:     "乌云盖顶",
			Signal:      "SELL",
			Confidence:  models.NewJSONDecimal(confidence),
			Description: "大阳线后跟大阴线，可能见顶信号",
			Strength:    strength,
			Volume:      current.Vol,
			PriceChange: models.NewJSONDecimal(currentChange),
		}
	}

	return nil
}

// recognizeHammer 识别锤子线模式
func (p *PatternRecognizer) recognizeHammer(current models.StockDaily) *models.CandlestickPattern {
	// 锤子线：下影线很长，实体较小，上影线很短或没有

	// 计算各部分长度
	body := current.Close.Decimal.Sub(current.Open.Decimal).Abs()

	// 确定最高和最低价格
	var higherPrice, lowerPrice decimal.Decimal
	if current.Close.Decimal.GreaterThan(current.Open.Decimal) {
		higherPrice = current.Close.Decimal
		lowerPrice = current.Open.Decimal
	} else {
		higherPrice = current.Open.Decimal
		lowerPrice = current.Close.Decimal
	}

	upperShadow := current.High.Decimal.Sub(higherPrice)
	lowerShadow := lowerPrice.Sub(current.Low.Decimal)

	// 判断条件：
	// 1. 下影线长度 > 实体长度的2倍
	// 2. 上影线长度 < 实体长度的0.5倍
	// 3. 实体长度 > 0（避免十字星）
	if lowerShadow.GreaterThan(body.Mul(decimal.NewFromInt(2))) &&
		upperShadow.LessThan(body.Mul(decimal.NewFromFloat(0.5))) &&
		body.GreaterThan(decimal.Zero) {

		// 计算置信度（下影线越长，置信度越高）
		confidence := lowerShadow.Div(body)
		if confidence.GreaterThan(decimal.NewFromInt(10)) {
			confidence = decimal.NewFromInt(10)
		}
		confidence = confidence.Div(decimal.NewFromInt(10)).Mul(decimal.NewFromInt(100))
		strength := p.calculateStrength(confidence)

		return &models.CandlestickPattern{
			TSCode:      current.TSCode,
			TradeDate:   current.TradeDate,
			Pattern:     "锤子线",
			Signal:      "BUY",
			Confidence:  models.NewJSONDecimal(confidence),
			Description: "下影线很长，可能见底信号",
			Strength:    strength,
			Volume:      current.Vol,
			PriceChange: models.NewJSONDecimal(current.Close.Decimal.Sub(current.Open.Decimal)),
		}
	}

	return nil
}

// recognizeMorningStar 识别启明星模式
func (p *PatternRecognizer) recognizeMorningStar(current, prev1, prev2 models.StockDaily, index int, data []models.StockDaily) *models.CandlestickPattern {
	// 启明星：下跌趋势中的反转信号
	// 第一根：大阴线
	// 第二根：小实体K线（十字星或小阴小阳）
	// 第三根：大阳线

	// 检查是否有足够的历史数据判断趋势
	if index < 5 {
		return nil
	}

	// 判断前5天是否处于下跌趋势
	trendDown := true
	for i := index - 4; i < index-1; i++ {
		if data[i].Close.Decimal.GreaterThanOrEqual(data[i-1].Close.Decimal) {
			trendDown = false
			break
		}
	}

	if !trendDown {
		return nil
	}

	// 第一根：大阴线
	prev2Change := prev2.Close.Decimal.Sub(prev2.Open.Decimal)
	prev2ChangePct := prev2Change.Div(prev2.Open.Decimal).Mul(decimal.NewFromInt(100))

	// 第二根：小实体K线
	prev1Body := prev1.Close.Decimal.Sub(prev1.Open.Decimal).Abs()
	prev1BodyPct := prev1Body.Div(prev1.Open.Decimal).Mul(decimal.NewFromInt(100))

	// 第三根：大阳线
	currentChange := current.Close.Decimal.Sub(current.Open.Decimal)
	currentChangePct := currentChange.Div(current.Open.Decimal).Mul(decimal.NewFromInt(100))

	// 判断条件：
	// 1. 第一根跌幅 > 2%
	// 2. 第二根实体 < 1%
	// 3. 第三根涨幅 > 2%
	if prev2ChangePct.LessThan(decimal.NewFromFloat(-2.0)) &&
		prev1BodyPct.LessThan(decimal.NewFromFloat(1.0)) &&
		currentChangePct.GreaterThan(decimal.NewFromFloat(2.0)) {

		confidence := p.calculateConfidence(prev2ChangePct.Abs(), currentChangePct, decimal.Zero)
		strength := p.calculateStrength(confidence)

		return &models.CandlestickPattern{
			TSCode:      current.TSCode,
			TradeDate:   current.TradeDate,
			Pattern:     "启明星",
			Signal:      "BUY",
			Confidence:  models.NewJSONDecimal(confidence),
			Description: "下跌趋势中的反转信号",
			Strength:    strength,
			Volume:      current.Vol,
			PriceChange: models.NewJSONDecimal(currentChange),
		}
	}

	return nil
}

// recognizeEveningStar 识别黄昏星模式
func (p *PatternRecognizer) recognizeEveningStar(current, prev1, prev2 models.StockDaily, index int, data []models.StockDaily) *models.CandlestickPattern {
	// 黄昏星：上涨趋势中的反转信号
	// 第一根：大阳线
	// 第二根：小实体K线（十字星或小阴小阳）
	// 第三根：大阴线

	// 检查是否有足够的历史数据判断趋势
	if index < 5 {
		return nil
	}

	// 判断前5天是否处于上涨趋势
	trendUp := true
	for i := index - 4; i < index-1; i++ {
		if data[i].Close.Decimal.LessThanOrEqual(data[i-1].Close.Decimal) {
			trendUp = false
			break
		}
	}

	if !trendUp {
		return nil
	}

	// 第一根：大阳线
	prev2Change := prev2.Close.Decimal.Sub(prev2.Open.Decimal)
	prev2ChangePct := prev2Change.Div(prev2.Open.Decimal).Mul(decimal.NewFromInt(100))

	// 第二根：小实体K线
	prev1Body := prev1.Close.Decimal.Sub(prev1.Open.Decimal).Abs()
	prev1BodyPct := prev1Body.Div(prev1.Open.Decimal).Mul(decimal.NewFromInt(100))

	// 第三根：大阴线
	currentChange := current.Close.Decimal.Sub(current.Open.Decimal)
	currentChangePct := currentChange.Div(current.Open.Decimal).Mul(decimal.NewFromInt(100))

	// 判断条件：
	// 1. 第一根涨幅 > 2%
	// 2. 第二根实体 < 1%
	// 3. 第三根跌幅 > 2%
	if prev2ChangePct.GreaterThan(decimal.NewFromFloat(2.0)) &&
		prev1BodyPct.LessThan(decimal.NewFromFloat(1.0)) &&
		currentChangePct.LessThan(decimal.NewFromFloat(-2.0)) {

		confidence := p.calculateConfidence(prev2ChangePct, currentChangePct.Abs(), decimal.Zero)
		strength := p.calculateStrength(confidence)

		return &models.CandlestickPattern{
			TSCode:      current.TSCode,
			TradeDate:   current.TradeDate,
			Pattern:     "黄昏星",
			Signal:      "SELL",
			Confidence:  models.NewJSONDecimal(confidence),
			Description: "上涨趋势中的反转信号",
			Strength:    strength,
			Volume:      current.Vol,
			PriceChange: models.NewJSONDecimal(currentChange),
		}
	}

	return nil
}

// recognizeVolumePriceRise 识别量价齐升模式
func (p *PatternRecognizer) recognizeVolumePriceRise(current, prev1, _ models.StockDaily) *models.VolumePricePattern {
	// 量价齐升：价格和成交量同时上涨

	// 计算价格变化
	priceChange := current.Close.Decimal.Sub(prev1.Close.Decimal)
	priceChangePct := priceChange.Div(prev1.Close.Decimal).Mul(decimal.NewFromInt(100))

	// 计算成交量变化
	volumeChange := current.Vol.Decimal.Sub(prev1.Vol.Decimal)
	volumeChangePct := volumeChange.Div(prev1.Vol.Decimal).Mul(decimal.NewFromInt(100))

	// 判断条件：
	// 1. 价格上涨 > 1%
	// 2. 成交量增加 > 20%
	if priceChangePct.GreaterThan(decimal.NewFromFloat(1.0)) &&
		volumeChangePct.GreaterThan(decimal.NewFromFloat(20.0)) {

		confidence := p.calculateConfidence(priceChangePct, volumeChangePct, decimal.Zero)
		strength := p.calculateStrength(confidence)

		return &models.VolumePricePattern{
			TSCode:      current.TSCode,
			TradeDate:   current.TradeDate,
			Pattern:     "量价齐升",
			Signal:      "BUY",
			Confidence:  models.NewJSONDecimal(confidence),
			Description: "价格和成交量同时上涨，强势信号",
			Strength:    strength,
			Volume:      current.Vol,
			PriceChange: models.NewJSONDecimal(priceChange),
			VolumeRatio: models.NewJSONDecimal(volumeChangePct),
		}
	}

	return nil
}

// recognizeVolumePriceDivergence 识别量价背离模式
func (p *PatternRecognizer) recognizeVolumePriceDivergence(current, prev1, _ models.StockDaily) *models.VolumePricePattern {
	// 量价背离：价格上涨但成交量下降，或价格下跌但成交量上升

	// 计算价格变化
	priceChange := current.Close.Decimal.Sub(prev1.Close.Decimal)
	priceChangePct := priceChange.Div(prev1.Close.Decimal).Mul(decimal.NewFromInt(100))

	// 计算成交量变化
	volumeChange := current.Vol.Decimal.Sub(prev1.Vol.Decimal)
	volumeChangePct := volumeChange.Div(prev1.Vol.Decimal).Mul(decimal.NewFromInt(100))

	// 判断条件：
	// 1. 价格上涨 > 1% 且成交量下降 > 20%（顶背离）
	// 2. 价格下跌 > 1% 且成交量上升 > 20%（底背离）
	if (priceChangePct.GreaterThan(decimal.NewFromFloat(1.0)) &&
		volumeChangePct.LessThan(decimal.NewFromFloat(-20.0))) ||
		(priceChangePct.LessThan(decimal.NewFromFloat(-1.0)) &&
			volumeChangePct.GreaterThan(decimal.NewFromFloat(20.0))) {

		var signal, description string
		if priceChangePct.GreaterThan(decimal.Zero) {
			signal = "SELL" // 顶背离
			description = "价格上涨但成交量下降，可能见顶"
		} else {
			signal = "BUY" // 底背离
			description = "价格下跌但成交量上升，可能见底"
		}

		confidence := p.calculateConfidence(priceChangePct.Abs(), volumeChangePct.Abs(), decimal.Zero)
		strength := p.calculateStrength(confidence)

		return &models.VolumePricePattern{
			TSCode:      current.TSCode,
			TradeDate:   current.TradeDate,
			Pattern:     "量价背离",
			Signal:      signal,
			Confidence:  models.NewJSONDecimal(confidence),
			Description: description,
			Strength:    strength,
			Volume:      current.Vol,
			PriceChange: models.NewJSONDecimal(priceChange),
			VolumeRatio: models.NewJSONDecimal(volumeChangePct),
		}
	}

	return nil
}

// recognizeVolumeBreakout 识别放量突破模式
func (p *PatternRecognizer) recognizeVolumeBreakout(current, prev1, _ models.StockDaily, index int, data []models.StockDaily) *models.VolumePricePattern {
	// 放量突破：价格突破重要阻力位，成交量明显放大

	// 检查是否有足够的历史数据
	if index < 20 {
		return nil
	}

	// 计算20日均线作为阻力位
	var sum decimal.Decimal
	for i := index - 19; i <= index; i++ {
		sum = sum.Add(data[i].Close.Decimal)
	}
	ma20 := sum.Div(decimal.NewFromInt(20))

	// 计算成交量比率（当前成交量与20日平均成交量比较）
	var volumeSum decimal.Decimal
	for i := index - 19; i <= index; i++ {
		volumeSum = volumeSum.Add(data[i].Vol.Decimal)
	}
	avgVolume := volumeSum.Div(decimal.NewFromInt(20))
	volumeRatio := current.Vol.Decimal.Div(avgVolume)

	// 判断条件：
	// 1. 价格突破20日均线
	// 2. 成交量是20日平均成交量的2倍以上
	// 3. 突破幅度 > 1%
	if current.Close.Decimal.GreaterThan(ma20) &&
		prev1.Close.Decimal.LessThanOrEqual(ma20) &&
		volumeRatio.GreaterThan(decimal.NewFromFloat(2.0)) {

		breakoutPct := current.Close.Decimal.Sub(ma20).Div(ma20).Mul(decimal.NewFromInt(100))
		confidence := p.calculateConfidence(breakoutPct, volumeRatio, decimal.Zero)
		strength := p.calculateStrength(confidence)

		return &models.VolumePricePattern{
			TSCode:      current.TSCode,
			TradeDate:   current.TradeDate,
			Pattern:     "放量突破",
			Signal:      "BUY",
			Confidence:  models.NewJSONDecimal(confidence),
			Description: "价格突破重要阻力位，成交量放大",
			Strength:    strength,
			Volume:      current.Vol,
			PriceChange: models.NewJSONDecimal(current.Close.Decimal.Sub(prev1.Close.Decimal)),
			VolumeRatio: models.NewJSONDecimal(volumeRatio),
		}
	}

	return nil
}

// calculateConfidence 计算置信度
func (p *PatternRecognizer) calculateConfidence(priceChange, volumeChange, volumeRatio decimal.Decimal) decimal.Decimal {
	// 基于价格变化、成交量变化和量比计算置信度
	var confidence decimal.Decimal

	// 价格变化权重：40%
	if !priceChange.IsZero() {
		priceConfidence := priceChange.Abs()
		if priceConfidence.GreaterThan(decimal.NewFromInt(10)) {
			priceConfidence = decimal.NewFromInt(10)
		}
		priceConfidence = priceConfidence.Div(decimal.NewFromInt(10)).Mul(decimal.NewFromInt(40))
		confidence = confidence.Add(priceConfidence)
	}

	// 成交量变化权重：30%
	if !volumeChange.IsZero() {
		volumeConfidence := volumeChange.Abs()
		if volumeConfidence.GreaterThan(decimal.NewFromInt(50)) {
			volumeConfidence = decimal.NewFromInt(50)
		}
		volumeConfidence = volumeConfidence.Div(decimal.NewFromInt(50)).Mul(decimal.NewFromInt(30))
		confidence = confidence.Add(volumeConfidence)
	}

	// 量比权重：30%
	if !volumeRatio.IsZero() {
		ratioConfidence := volumeRatio
		if ratioConfidence.GreaterThan(decimal.NewFromInt(5)) {
			ratioConfidence = decimal.NewFromInt(5)
		}
		ratioConfidence = ratioConfidence.Div(decimal.NewFromInt(5)).Mul(decimal.NewFromInt(30))
		confidence = confidence.Add(ratioConfidence)
	}

	// 限制置信度最大值为100
	if confidence.GreaterThan(decimal.NewFromInt(100)) {
		confidence = decimal.NewFromInt(100)
	}

	return confidence
}

// calculateStrength 计算信号强度
func (p *PatternRecognizer) calculateStrength(confidence decimal.Decimal) string {
	if confidence.GreaterThanOrEqual(decimal.NewFromInt(80)) {
		return "STRONG"
	} else if confidence.GreaterThanOrEqual(decimal.NewFromInt(60)) {
		return "MEDIUM"
	} else {
		return "WEAK"
	}
}

// calculateCombinedSignal 计算综合信号
func (p *PatternRecognizer) calculateCombinedSignal(candlestick []models.CandlestickPattern, volumePrice []models.VolumePricePattern) (string, decimal.Decimal, string) {
	if len(candlestick) == 0 && len(volumePrice) == 0 {
		return "HOLD", decimal.Zero, "LOW"
	}

	var totalConfidence decimal.Decimal
	var buySignals, sellSignals int
	var totalSignals int

	// 统计蜡烛图信号
	for _, pattern := range candlestick {
		totalConfidence = totalConfidence.Add(pattern.Confidence.Decimal)
		totalSignals++
		switch pattern.Signal {
		case "BUY":
			buySignals++
		case "SELL":
			sellSignals++
		}
	}

	// 统计量价图形信号
	for _, pattern := range volumePrice {
		totalConfidence = totalConfidence.Add(pattern.Confidence.Decimal)
		totalSignals++
		switch pattern.Signal {
		case "BUY":
			buySignals++
		case "SELL":
			sellSignals++
		}
	}

	// 计算平均置信度
	avgConfidence := totalConfidence.Div(decimal.NewFromInt(int64(totalSignals)))

	// 确定综合信号
	var combinedSignal string
	if buySignals > sellSignals {
		combinedSignal = "BUY"
	} else if sellSignals > buySignals {
		combinedSignal = "SELL"
	} else {
		combinedSignal = "HOLD"
	}

	// 确定风险等级
	var riskLevel string
	if avgConfidence.GreaterThanOrEqual(decimal.NewFromInt(80)) {
		riskLevel = "LOW"
	} else if avgConfidence.GreaterThanOrEqual(decimal.NewFromInt(60)) {
		riskLevel = "MEDIUM"
	} else {
		riskLevel = "HIGH"
	}

	return combinedSignal, avgConfidence, riskLevel
}
