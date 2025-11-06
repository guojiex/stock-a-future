package service

import (
	"fmt"
	"regexp"
	"sort"
	"stock-a-future/internal/indicators"
	"stock-a-future/internal/models"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// PredictionService 预测服务
type PredictionService struct {
	calculator        *indicators.Calculator
	patternRecognizer *indicators.PatternRecognizer
}

// NewPredictionService 创建预测服务
func NewPredictionService() *PredictionService {
	return &PredictionService{
		calculator:        indicators.NewCalculator(),
		patternRecognizer: indicators.NewPatternRecognizer(),
	}
}

// PredictTradingPoints 预测买卖点
func (s *PredictionService) PredictTradingPoints(data []models.StockDaily) (*models.PredictionResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("数据为空")
	}

	// 确保数据按日期排序（最新的在最后）
	latestData := data[len(data)-1]

	// 计算各种技术指标
	indicators := s.calculateAllIndicators(data)

	// 基于技术指标生成预测
	predictions := s.generatePredictions(data, indicators)

	// 对历史预测进行回测
	predictions = s.backTestPredictions(data, predictions)

	// 计算整体置信度
	confidence := s.calculateOverallConfidence(predictions)

	return &models.PredictionResult{
		TSCode:      latestData.TSCode,
		TradeDate:   latestData.TradeDate,
		Predictions: predictions,
		Confidence:  confidence,
		UpdatedAt:   time.Now(),
	}, nil
}

// calculateAllIndicators 计算所有技术指标
func (s *PredictionService) calculateAllIndicators(data []models.StockDaily) *models.TechnicalIndicators {
	if len(data) == 0 {
		return nil
	}

	latestData := data[len(data)-1]
	indicators := &models.TechnicalIndicators{
		TSCode:    latestData.TSCode,
		TradeDate: latestData.TradeDate,
	}

	// 计算MACD
	macdResults := s.calculator.CalculateMACD(data)
	if len(macdResults) > 0 {
		indicators.MACD = &macdResults[len(macdResults)-1]
	}

	// 计算RSI
	rsiResults := s.calculator.CalculateRSI(data, 14)
	if len(rsiResults) > 0 {
		indicators.RSI = &rsiResults[len(rsiResults)-1]
	}

	// 计算布林带
	bollResults := s.calculator.CalculateBollingerBands(data, 20, 2.0)
	if len(bollResults) > 0 {
		indicators.BOLL = &bollResults[len(bollResults)-1]
	}

	// 计算移动平均线
	ma5 := s.calculator.CalculateMA(data, 5)
	ma10 := s.calculator.CalculateMA(data, 10)
	ma20 := s.calculator.CalculateMA(data, 20)
	ma60 := s.calculator.CalculateMA(data, 60)
	ma120 := s.calculator.CalculateMA(data, 120)

	if len(ma5) > 0 && len(ma10) > 0 && len(ma20) > 0 {
		indicators.MA = &models.MovingAverageIndicator{
			MA5:  models.NewJSONDecimal(ma5[len(ma5)-1]),
			MA10: models.NewJSONDecimal(ma10[len(ma10)-1]),
			MA20: models.NewJSONDecimal(ma20[len(ma20)-1]),
		}
		if len(ma60) > 0 {
			indicators.MA.MA60 = models.NewJSONDecimal(ma60[len(ma60)-1])
		}
		if len(ma120) > 0 {
			indicators.MA.MA120 = models.NewJSONDecimal(ma120[len(ma120)-1])
		}
	}

	// 计算KDJ
	kdjResults := s.calculator.CalculateKDJ(data, 9)
	if len(kdjResults) > 0 {
		indicators.KDJ = &kdjResults[len(kdjResults)-1]
	}

	return indicators
}

// generatePredictions 基于技术指标和图形模式生成预测
func (s *PredictionService) generatePredictions(data []models.StockDaily, indicators *models.TechnicalIndicators) []models.TradingPointPrediction {
	var predictions []models.TradingPointPrediction
	currentPrice := data[len(data)-1].Close.Decimal
	signalDate := data[len(data)-1].TradeDate // 信号基于最新数据的日期

	// 基于MACD的预测
	if indicators.MACD != nil {
		// 查找MACD信号产生的日期对应的价格
		signalPrice := s.findPriceBySignalType(data, "MACD", indicators.MACD.Signal)
		if prediction := s.predictFromMACD(indicators.MACD, signalPrice, signalDate); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	// 基于RSI的预测
	if indicators.RSI != nil {
		// 查找RSI信号产生的日期对应的价格
		signalPrice := s.findPriceBySignalType(data, "RSI", indicators.RSI.Signal)
		if prediction := s.predictFromRSI(indicators.RSI, signalPrice, signalDate); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	// 基于布林带的预测
	if indicators.BOLL != nil {
		// 查找布林带信号产生的日期对应的价格
		signalPrice := s.findPriceBySignalType(data, "BOLL", indicators.BOLL.Signal)
		if prediction := s.predictFromBollingerBands(indicators.BOLL, signalPrice, signalDate); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	// 基于KDJ的预测
	if indicators.KDJ != nil {
		// 查找KDJ信号产生的日期对应的价格
		signalPrice := s.findPriceBySignalType(data, "KDJ", indicators.KDJ.Signal)
		if prediction := s.predictFromKDJ(indicators.KDJ, signalPrice, signalDate); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	// 基于移动平均线的预测
	if indicators.MA != nil {
		// 查找MA信号产生的日期对应的价格
		signalPrice := currentPrice // 默认使用当前价格
		if prediction := s.predictFromMA(indicators.MA, signalPrice, signalDate); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	// 基于图形模式的预测（新增）
	if patternPredictions := s.predictFromPatterns(data); len(patternPredictions) > 0 {
		predictions = append(predictions, patternPredictions...)
	}

	// 合并同一天的相同类型信号
	predictions = s.mergeSameDaySignals(predictions)

	// 对预测结果进行排序：置信度和强度最高的排在前面
	s.sortPredictionsByConfidenceAndStrength(predictions)

	return predictions
}

// predictFromMACD 基于MACD指标预测
func (s *PredictionService) predictFromMACD(macd *models.MACDIndicator, currentPrice decimal.Decimal, signalDate string) *models.TradingPointPrediction {
	if macd.Signal == "HOLD" {
		return nil
	}

	var predictType string
	var probability decimal.Decimal
	var reason string

	switch macd.Signal {
	case "BUY":
		predictType = "BUY"
		probability = decimal.NewFromFloat(0.65) // 基础概率65%
		reason = "MACD金叉信号，DIF线上穿DEA线"
	case "SELL":
		predictType = "SELL"
		probability = decimal.NewFromFloat(0.60) // 基础概率60%
		reason = "MACD死叉信号，DIF线下穿DEA线"
	default:
		return nil
	}

	// 根据MACD强度调整概率
	macdStrength := macd.Histogram.Decimal.Abs()
	if macdStrength.GreaterThan(decimal.NewFromFloat(0.1)) {
		probability = probability.Add(decimal.NewFromFloat(0.1))
	}

	// 使用信号产生时的价格
	signalPrice := currentPrice

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       models.NewJSONDecimal(signalPrice),
		Date:        signalDate, // 使用实际数据的交易日期
		Probability: models.NewJSONDecimal(probability),
		Reason:      reason,
		Indicators:  []string{"MACD"},
		SignalDate:  signalDate, // 信号基于传入的日期
	}
}

// predictFromRSI 基于RSI指标预测
func (s *PredictionService) predictFromRSI(rsi *models.RSIIndicator, currentPrice decimal.Decimal, signalDate string) *models.TradingPointPrediction {
	if rsi.Signal == "HOLD" {
		return nil
	}

	var predictType string
	var probability decimal.Decimal
	var reason string

	switch rsi.Signal {
	case "BUY":
		predictType = "BUY"
		probability = decimal.NewFromFloat(0.70) // RSI超卖信号相对可靠
		reason = fmt.Sprintf("RSI超卖信号，当前RSI值：%.2f", rsi.RSI14.Decimal.InexactFloat64())
	case "SELL":
		predictType = "SELL"
		probability = decimal.NewFromFloat(0.65)
		reason = fmt.Sprintf("RSI超买信号，当前RSI值：%.2f", rsi.RSI14.Decimal.InexactFloat64())
	default:
		return nil
	}

	// 使用信号产生时的价格
	signalPrice := currentPrice

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       models.NewJSONDecimal(signalPrice),
		Date:        signalDate, // 使用实际数据的交易日期
		Probability: models.NewJSONDecimal(probability),
		Reason:      reason,
		Indicators:  []string{"RSI"},
		SignalDate:  signalDate, // 信号基于传入的日期
	}
}

// predictFromBollingerBands 基于布林带预测
func (s *PredictionService) predictFromBollingerBands(boll *models.BollingerBandsIndicator, currentPrice decimal.Decimal, signalDate string) *models.TradingPointPrediction {
	if boll.Signal == "HOLD" {
		return nil
	}

	var predictType string
	var probability decimal.Decimal
	var reason string

	switch boll.Signal {
	case "BUY":
		predictType = "BUY"
		probability = decimal.NewFromFloat(0.60)
		reason = "价格触及布林带下轨，可能出现反弹"
	case "SELL":
		predictType = "SELL"
		probability = decimal.NewFromFloat(0.55)
		reason = "价格触及布林带上轨，可能出现回调"
	default:
		return nil
	}

	// 使用信号产生时的价格
	signalPrice := currentPrice

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       models.NewJSONDecimal(signalPrice),
		Date:        signalDate, // 使用实际数据的交易日期
		Probability: models.NewJSONDecimal(probability),
		Reason:      reason,
		Indicators:  []string{"BOLL"},
		SignalDate:  signalDate, // 信号基于传入的日期
	}
}

// predictFromKDJ 基于KDJ指标预测
func (s *PredictionService) predictFromKDJ(kdj *models.KDJIndicator, currentPrice decimal.Decimal, signalDate string) *models.TradingPointPrediction {
	if kdj.Signal == "HOLD" {
		return nil
	}

	var predictType string
	var probability decimal.Decimal
	var reason string

	switch kdj.Signal {
	case "BUY":
		predictType = "BUY"
		probability = decimal.NewFromFloat(0.58)
		reason = fmt.Sprintf("KDJ超卖信号，K值：%.2f，D值：%.2f", kdj.K.Decimal.InexactFloat64(), kdj.D.Decimal.InexactFloat64())
	case "SELL":
		predictType = "SELL"
		probability = decimal.NewFromFloat(0.55)
		reason = fmt.Sprintf("KDJ超买信号，K值：%.2f，D值：%.2f", kdj.K.Decimal.InexactFloat64(), kdj.D.Decimal.InexactFloat64())
	default:
		return nil
	}

	// 使用信号产生时的价格
	signalPrice := currentPrice

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       models.NewJSONDecimal(signalPrice),
		Date:        signalDate, // 使用实际数据的交易日期
		Probability: models.NewJSONDecimal(probability),
		Reason:      reason,
		Indicators:  []string{"KDJ"},
		SignalDate:  signalDate, // 信号基于传入的日期
	}
}

// predictFromMA 基于移动平均线预测
func (s *PredictionService) predictFromMA(ma *models.MovingAverageIndicator, currentPrice decimal.Decimal, signalDate string) *models.TradingPointPrediction {
	// 判断均线多头/空头排列
	var predictType string
	var probability decimal.Decimal
	var reason string

	// 多头排列：MA5 > MA10 > MA20
	if ma.MA5.Decimal.GreaterThan(ma.MA10.Decimal) && ma.MA10.Decimal.GreaterThan(ma.MA20.Decimal) {
		if currentPrice.GreaterThan(ma.MA5.Decimal) {
			predictType = "BUY"
			probability = decimal.NewFromFloat(0.62)
			reason = "均线多头排列，价格在5日均线之上"
		}
	}

	// 空头排列：MA5 < MA10 < MA20
	if ma.MA5.Decimal.LessThan(ma.MA10.Decimal) && ma.MA10.Decimal.LessThan(ma.MA20.Decimal) {
		if currentPrice.LessThan(ma.MA5.Decimal) {
			predictType = "SELL"
			probability = decimal.NewFromFloat(0.58)
			reason = "均线空头排列，价格在5日均线之下"
		}
	}

	if predictType == "" {
		return nil
	}

	// 使用信号产生时的价格
	signalPrice := currentPrice

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       models.NewJSONDecimal(signalPrice),
		Date:        signalDate, // 使用实际数据的交易日期
		Probability: models.NewJSONDecimal(probability),
		Reason:      reason,
		Indicators:  []string{"MA"},
		SignalDate:  signalDate, // 信号基于传入的日期
	}
}

// calculateOverallConfidence 计算整体置信度
func (s *PredictionService) calculateOverallConfidence(predictions []models.TradingPointPrediction) models.JSONDecimal {
	if len(predictions) == 0 {
		return models.NewJSONDecimal(decimal.Zero)
	}

	// 统计买卖信号
	buyCount := 0
	sellCount := 0
	totalProbability := decimal.Zero

	for _, pred := range predictions {
		totalProbability = totalProbability.Add(pred.Probability.Decimal)
		switch pred.Type {
		case "BUY":
			buyCount++
		case "SELL":
			sellCount++
		}
	}

	// 计算平均概率
	avgProbability := totalProbability.Div(decimal.NewFromInt(int64(len(predictions))))

	// 如果信号一致性高，提升置信度
	consistency := decimal.Zero
	if buyCount > 0 && sellCount == 0 {
		consistency = decimal.NewFromFloat(0.1) // 全部买入信号
	} else if sellCount > 0 && buyCount == 0 {
		consistency = decimal.NewFromFloat(0.1) // 全部卖出信号
	}

	confidence := avgProbability.Add(consistency)

	// 限制置信度在0-1之间
	if confidence.GreaterThan(decimal.NewFromInt(1)) {
		confidence = decimal.NewFromInt(1)
	}

	return models.NewJSONDecimal(confidence)
}

// predictFromPatterns 基于图形模式生成预测
func (s *PredictionService) predictFromPatterns(data []models.StockDaily) []models.TradingPointPrediction {
	var predictions []models.TradingPointPrediction

	// 识别所有图形模式
	patterns := s.patternRecognizer.RecognizeAllPatterns(data)

	// 用于去重的映射：模式类型 -> 最佳预测
	patternMap := make(map[string]*models.TradingPointPrediction)

	// 遍历识别到的模式，生成预测
	for _, pattern := range patterns {
		// 处理蜡烛图模式
		for _, candlestick := range pattern.Candlestick {
			// 查找该蜡烛图模式对应的价格
			patternPrice := s.findPriceForPattern(data, candlestick.TradeDate, candlestick.Signal)

			if prediction := s.predictFromCandlestickPattern(candlestick, patternPrice); prediction != nil {
				// 生成模式标识符（类型+指标）
				patternKey := fmt.Sprintf("candlestick:%s", candlestick.Pattern)

				// 如果已存在相同模式，比较置信度，保留最高的
				if existing, exists := patternMap[patternKey]; exists {
					if prediction.Probability.Decimal.GreaterThan(existing.Probability.Decimal) {
						patternMap[patternKey] = prediction
					}
				} else {
					patternMap[patternKey] = prediction
				}
			}
		}

		// 处理量价模式
		for _, volumePrice := range pattern.VolumePrice {
			// 查找该量价模式对应的价格
			patternPrice := s.findPriceForPattern(data, volumePrice.TradeDate, volumePrice.Signal)

			if prediction := s.predictFromVolumePricePattern(volumePrice, patternPrice); prediction != nil {
				// 生成模式标识符（类型+指标）
				patternKey := fmt.Sprintf("volume_price:%s", volumePrice.Pattern)

				// 如果已存在相同模式，比较置信度，保留最高的
				if existing, exists := patternMap[patternKey]; exists {
					if prediction.Probability.Decimal.GreaterThan(existing.Probability.Decimal) {
						patternMap[patternKey] = prediction
					}
				} else {
					patternMap[patternKey] = prediction
				}
			}
		}
	}

	// 将去重后的预测添加到结果中
	for _, prediction := range patternMap {
		predictions = append(predictions, *prediction)
	}

	return predictions
}

// predictFromCandlestickPattern 基于蜡烛图模式生成预测
func (s *PredictionService) predictFromCandlestickPattern(pattern models.CandlestickPattern, currentPrice decimal.Decimal) *models.TradingPointPrediction {
	var predictType string
	var probability decimal.Decimal
	var reason string

	switch pattern.Pattern {
	case "双响炮":
		predictType = "BUY"
		probability = decimal.NewFromFloat(0.75) // 双响炮是强势买入信号
		reason = fmt.Sprintf("识别到双响炮模式，置信度：%s，强度：%s",
			pattern.Confidence.Decimal.String(), pattern.Strength)
	case "红三兵":
		predictType = "BUY"
		probability = decimal.NewFromFloat(0.70)
		reason = fmt.Sprintf("识别到红三兵模式，置信度：%s，强度：%s",
			pattern.Confidence.Decimal.String(), pattern.Strength)
	case "乌云盖顶":
		predictType = "SELL"
		probability = decimal.NewFromFloat(0.65)
		reason = fmt.Sprintf("识别到乌云盖顶模式，置信度：%s，强度：%s",
			pattern.Confidence.Decimal.String(), pattern.Strength)
	case "锤子线":
		predictType = "BUY"
		probability = decimal.NewFromFloat(0.60)
		reason = fmt.Sprintf("识别到锤子线模式，置信度：%s，强度：%s",
			pattern.Confidence.Decimal.String(), pattern.Strength)
	case "启明星":
		predictType = "BUY"
		probability = decimal.NewFromFloat(0.70)
		reason = fmt.Sprintf("识别到启明星模式，置信度：%s，强度：%s",
			pattern.Confidence.Decimal.String(), pattern.Strength)
	case "黄昏星":
		predictType = "SELL"
		probability = decimal.NewFromFloat(0.65)
		reason = fmt.Sprintf("识别到黄昏星模式，置信度：%s，强度：%s",
			pattern.Confidence.Decimal.String(), pattern.Strength)
	default:
		return nil
	}

	// 根据置信度调整概率
	if pattern.Confidence.Decimal.GreaterThan(decimal.Zero) {
		confidenceFactor := pattern.Confidence.Decimal.Mul(decimal.NewFromFloat(0.2))
		probability = probability.Add(confidenceFactor)
		// 限制概率不超过0.95
		if probability.GreaterThan(decimal.NewFromFloat(0.95)) {
			probability = decimal.NewFromFloat(0.95)
		}
	}

	// 使用模式识别时的价格，而不是当前价格
	// 这样可以确保每个信号显示的是信号产生时的价格
	signalPrice := currentPrice

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       models.NewJSONDecimal(signalPrice),
		Date:        pattern.TradeDate, // 使用实际数据的交易日期
		Probability: models.NewJSONDecimal(probability),
		Reason:      reason,
		Indicators:  []string{fmt.Sprintf("图形模式:%s", pattern.Pattern)},
		SignalDate:  pattern.TradeDate, // 信号基于识别到模式的日期
	}
}

// predictFromVolumePricePattern 基于量价模式生成预测
func (s *PredictionService) predictFromVolumePricePattern(pattern models.VolumePricePattern, currentPrice decimal.Decimal) *models.TradingPointPrediction {
	var predictType string
	var probability decimal.Decimal
	var reason string

	switch pattern.Pattern {
	case "量价齐升":
		predictType = "BUY"
		probability = decimal.NewFromFloat(0.65)
		reason = fmt.Sprintf("识别到量价齐升模式，置信度：%s，强度：%s",
			pattern.Confidence.Decimal.String(), pattern.Strength)
	case "量价背离":
		predictType = "SELL"
		probability = decimal.NewFromFloat(0.60)
		reason = fmt.Sprintf("识别到量价背离模式，置信度：%s，强度：%s",
			pattern.Confidence.Decimal.String(), pattern.Strength)
	case "放量突破":
		predictType = "BUY"
		probability = decimal.NewFromFloat(0.70)
		reason = fmt.Sprintf("识别到放量突破模式，置信度：%s，强度：%s",
			pattern.Confidence.Decimal.String(), pattern.Strength)
	default:
		return nil
	}

	// 根据置信度调整概率
	if pattern.Confidence.Decimal.GreaterThan(decimal.Zero) {
		confidenceFactor := pattern.Confidence.Decimal.Mul(decimal.NewFromFloat(0.15))
		probability = probability.Add(confidenceFactor)
		// 限制概率不超过0.90
		if probability.GreaterThan(decimal.NewFromFloat(0.90)) {
			probability = decimal.NewFromFloat(0.90)
		}
	}

	// 使用模式识别时的价格，而不是当前价格
	// 这样可以确保每个信号显示的是信号产生时的价格
	signalPrice := currentPrice

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       models.NewJSONDecimal(signalPrice),
		Date:        pattern.TradeDate, // 使用实际数据的交易日期
		Probability: models.NewJSONDecimal(probability),
		Reason:      reason,
		Indicators:  []string{fmt.Sprintf("量价模式:%s", pattern.Pattern)},
		SignalDate:  pattern.TradeDate, // 信号基于识别到模式的日期
	}
}

// mergeSameDaySignals 合并同一天的相同类型信号
func (s *PredictionService) mergeSameDaySignals(predictions []models.TradingPointPrediction) []models.TradingPointPrediction {
	if len(predictions) == 0 {
		return predictions
	}

	// 创建映射：SignalDate + Type -> []TradingPointPrediction
	signalMap := make(map[string][]models.TradingPointPrediction)

	for _, pred := range predictions {
		key := fmt.Sprintf("%s_%s", pred.SignalDate, pred.Type)
		signalMap[key] = append(signalMap[key], pred)
	}

	// 合并信号
	var mergedPredictions []models.TradingPointPrediction

	for _, signals := range signalMap {
		if len(signals) == 1 {
			// 只有一个信号，直接添加
			mergedPredictions = append(mergedPredictions, signals[0])
		} else {
			// 多个信号，需要合并
			merged := s.mergeMultipleSignals(signals)
			mergedPredictions = append(mergedPredictions, merged)
		}
	}

	return mergedPredictions
}

// mergeMultipleSignals 合并多个相同日期和类型的信号
func (s *PredictionService) mergeMultipleSignals(signals []models.TradingPointPrediction) models.TradingPointPrediction {
	if len(signals) == 0 {
		return models.TradingPointPrediction{}
	}

	// 使用第一个信号作为基础
	merged := signals[0]

	// 收集所有指标
	indicatorSet := make(map[string]bool)
	for _, sig := range signals {
		for _, ind := range sig.Indicators {
			indicatorSet[ind] = true
		}
	}

	// 转换为切片
	var indicators []string
	for ind := range indicatorSet {
		indicators = append(indicators, ind)
	}
	merged.Indicators = indicators

	// 合并原因：选择最有说服力的几个原因
	var reasons []string
	reasonSet := make(map[string]bool)
	for _, sig := range signals {
		if sig.Reason != "" && !reasonSet[sig.Reason] {
			reasons = append(reasons, sig.Reason)
			reasonSet[sig.Reason] = true
		}
	}

	// 将原因合并，但限制长度
	if len(reasons) > 3 {
		// 只保留前3个最重要的原因
		merged.Reason = fmt.Sprintf("%s 等%d个信号", strings.Join(reasons[:3], "；"), len(reasons))
	} else {
		merged.Reason = strings.Join(reasons, "；")
	}

	// 计算综合概率：多个信号共识时应该提升置信度
	totalProbability := decimal.Zero
	maxProbability := decimal.Zero

	for _, sig := range signals {
		totalProbability = totalProbability.Add(sig.Probability.Decimal)
		if sig.Probability.Decimal.GreaterThan(maxProbability) {
			maxProbability = sig.Probability.Decimal
		}
	}

	// 使用加权平均：基础概率是平均值，然后根据信号数量提升
	avgProbability := totalProbability.Div(decimal.NewFromInt(int64(len(signals))))

	// 信号共识度提升：每多一个信号，提升5%，但总提升不超过15%
	consensusBonus := decimal.NewFromFloat(0.05).Mul(decimal.NewFromInt(int64(len(signals) - 1)))
	maxBonus := decimal.NewFromFloat(0.15)
	if consensusBonus.GreaterThan(maxBonus) {
		consensusBonus = maxBonus
	}

	// 最终概率 = 平均概率 + 共识提升，但不超过最大概率的1.2倍，且不超过0.95
	finalProbability := avgProbability.Add(consensusBonus)
	maxAllowed := maxProbability.Mul(decimal.NewFromFloat(1.2))
	if finalProbability.GreaterThan(maxAllowed) {
		finalProbability = maxAllowed
	}
	if finalProbability.GreaterThan(decimal.NewFromFloat(0.95)) {
		finalProbability = decimal.NewFromFloat(0.95)
	}

	merged.Probability = models.NewJSONDecimal(finalProbability)

	// 价格选择：买入信号用最低价，卖出信号用最高价
	if merged.Type == "BUY" {
		lowestPrice := signals[0].Price.Decimal
		for _, sig := range signals[1:] {
			if sig.Price.Decimal.LessThan(lowestPrice) {
				lowestPrice = sig.Price.Decimal
			}
		}
		merged.Price = models.NewJSONDecimal(lowestPrice)
	} else if merged.Type == "SELL" {
		highestPrice := signals[0].Price.Decimal
		for _, sig := range signals[1:] {
			if sig.Price.Decimal.GreaterThan(highestPrice) {
				highestPrice = sig.Price.Decimal
			}
		}
		merged.Price = models.NewJSONDecimal(highestPrice)
	}

	// 回测信息：如果所有信号都有回测数据，使用平均值
	allBacktested := true
	backtestCount := 0
	totalNextDayPrice := decimal.Zero
	totalPriceDiff := decimal.Zero
	totalPriceDiffRatio := decimal.Zero
	correctCount := 0

	for _, sig := range signals {
		if sig.Backtested {
			backtestCount++
			totalNextDayPrice = totalNextDayPrice.Add(sig.NextDayPrice.Decimal)
			totalPriceDiff = totalPriceDiff.Add(sig.PriceDiff.Decimal)
			totalPriceDiffRatio = totalPriceDiffRatio.Add(sig.PriceDiffRatio.Decimal)
			if sig.IsCorrect {
				correctCount++
			}
		} else {
			allBacktested = false
		}
	}

	if allBacktested && backtestCount > 0 {
		merged.Backtested = true
		merged.NextDayPrice = models.NewJSONDecimal(totalNextDayPrice.Div(decimal.NewFromInt(int64(backtestCount))))
		merged.PriceDiff = models.NewJSONDecimal(totalPriceDiff.Div(decimal.NewFromInt(int64(backtestCount))))
		merged.PriceDiffRatio = models.NewJSONDecimal(totalPriceDiffRatio.Div(decimal.NewFromInt(int64(backtestCount))))
		// 如果超过一半的信号预测正确，则认为合并后的信号正确
		merged.IsCorrect = correctCount > len(signals)/2
	}

	return merged
}

// sortPredictionsByConfidenceAndStrength 对预测结果进行排序
// 排序规则：1. 信号日期（最新的在前） 2. 强度等级（STRONG > MEDIUM > WEAK） 3. 置信度（高到低） 4. 概率（高到低）
func (s *PredictionService) sortPredictionsByConfidenceAndStrength(predictions []models.TradingPointPrediction) {
	// 定义强度等级权重
	strengthWeight := map[string]int{
		"STRONG": 3,
		"MEDIUM": 2,
		"WEAK":   1,
	}

	// 使用sort.Slice进行排序
	sort.Slice(predictions, func(i, j int) bool {
		pred1 := predictions[i]
		pred2 := predictions[j]

		// 首先按信号日期排序（最新的在前）
		if pred1.SignalDate != pred2.SignalDate {
			return pred1.SignalDate > pred2.SignalDate
		}

		// 日期相同时，按强度等级排序
		strength1 := s.extractStrengthFromReason(pred1.Reason)
		strength2 := s.extractStrengthFromReason(pred2.Reason)

		if strengthWeight[strength1] != strengthWeight[strength2] {
			return strengthWeight[strength1] > strengthWeight[strength2]
		}

		// 强度等级相同时，按置信度排序
		confidence1 := s.extractConfidenceFromReason(pred1.Reason)
		confidence2 := s.extractConfidenceFromReason(pred2.Reason)

		if !confidence1.Equal(confidence2) {
			return confidence1.GreaterThan(confidence2)
		}

		// 置信度相同时，按概率排序
		return pred1.Probability.Decimal.GreaterThan(pred2.Probability.Decimal)
	})
}

// extractStrengthFromReason 从reason字段中提取强度等级
func (s *PredictionService) extractStrengthFromReason(reason string) string {
	if reason == "" {
		return "WEAK"
	}

	// 查找强度等级
	if strings.Contains(reason, "强度：STRONG") {
		return "STRONG"
	} else if strings.Contains(reason, "强度：MEDIUM") {
		return "MEDIUM"
	} else if strings.Contains(reason, "强度：WEAK") {
		return "WEAK"
	}

	return "WEAK" // 默认值
}

// extractConfidenceFromReason 从reason字段中提取置信度
func (s *PredictionService) extractConfidenceFromReason(reason string) decimal.Decimal {
	if reason == "" {
		return decimal.Zero
	}

	// 查找置信度值
	confidencePattern := regexp.MustCompile(`置信度：([\d.]+)`)
	matches := confidencePattern.FindStringSubmatch(reason)
	if len(matches) > 1 {
		if confidence, err := decimal.NewFromString(matches[1]); err == nil {
			return confidence
		}
	}

	return decimal.Zero
}

// findPriceBySignalType 根据信号类型查找对应的价格
func (s *PredictionService) findPriceBySignalType(data []models.StockDaily, indicatorType string, signalType string) decimal.Decimal {
	// 如果没有数据或信号类型为空，则返回最新价格
	if len(data) == 0 || signalType == "" || signalType == "HOLD" {
		return data[len(data)-1].Close.Decimal
	}

	// 默认使用最新价格
	currentPrice := data[len(data)-1].Close.Decimal

	// 根据不同的指标类型和信号类型，查找对应的价格
	// 这里简化处理，实际上应该根据指标计算逻辑找到信号产生的具体日期
	switch indicatorType {
	case "MACD":
		// MACD金叉/死叉通常在最近几天产生
		// 简化处理：使用最近3天内的价格
		if len(data) > 3 {
			switch signalType {
			case "BUY":
				// 金叉通常在低点附近，使用近期最低价
				lowestPrice := data[len(data)-3].Close.Decimal
				for i := len(data) - 3; i < len(data); i++ {
					if data[i].Close.Decimal.LessThan(lowestPrice) {
						lowestPrice = data[i].Close.Decimal
					}
				}
				return lowestPrice
			case "SELL":
				// 死叉通常在高点附近，使用近期最高价
				highestPrice := data[len(data)-3].Close.Decimal
				for i := len(data) - 3; i < len(data); i++ {
					if data[i].Close.Decimal.GreaterThan(highestPrice) {
						highestPrice = data[i].Close.Decimal
					}
				}
				return highestPrice
			}
		}
	case "RSI":
		// RSI超买/超卖通常在价格极值点产生
		switch signalType {
		case "BUY":
			// 超卖信号通常在低点，使用近期最低价
			lowestPrice := currentPrice
			for i := max(0, len(data)-5); i < len(data); i++ {
				if data[i].Close.Decimal.LessThan(lowestPrice) {
					lowestPrice = data[i].Close.Decimal
				}
			}
			return lowestPrice
		case "SELL":
			// 超买信号通常在高点，使用近期最高价
			highestPrice := currentPrice
			for i := max(0, len(data)-5); i < len(data); i++ {
				if data[i].Close.Decimal.GreaterThan(highestPrice) {
					highestPrice = data[i].Close.Decimal
				}
			}
			return highestPrice
		}
	case "BOLL":
		// 布林带突破信号
		switch signalType {
		case "BUY":
			// 价格触及下轨时的价格
			for i := max(0, len(data)-5); i < len(data); i++ {
				// 简化处理：使用近期最低价
				if i == max(0, len(data)-5) || data[i].Low.Decimal.LessThan(currentPrice) {
					currentPrice = data[i].Low.Decimal
				}
			}
			return currentPrice
		case "SELL":
			// 价格触及上轨时的价格
			for i := max(0, len(data)-5); i < len(data); i++ {
				// 简化处理：使用近期最高价
				if i == max(0, len(data)-5) || data[i].High.Decimal.GreaterThan(currentPrice) {
					currentPrice = data[i].High.Decimal
				}
			}
			return currentPrice
		}
	case "KDJ":
		// KDJ超买/超卖信号
		// 简化处理：与RSI类似
		switch signalType {
		case "BUY":
			// 超卖信号通常在低点
			lowestPrice := currentPrice
			for i := max(0, len(data)-5); i < len(data); i++ {
				if data[i].Close.Decimal.LessThan(lowestPrice) {
					lowestPrice = data[i].Close.Decimal
				}
			}
			return lowestPrice
		case "SELL":
			// 超买信号通常在高点
			highestPrice := currentPrice
			for i := max(0, len(data)-5); i < len(data); i++ {
				if data[i].Close.Decimal.GreaterThan(highestPrice) {
					highestPrice = data[i].Close.Decimal
				}
			}
			return highestPrice
		}
	}

	// 默认返回当前价格
	return currentPrice
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// findPriceForPattern 根据模式的日期和类型查找对应的价格
func (s *PredictionService) findPriceForPattern(data []models.StockDaily, patternDate, signalType string) decimal.Decimal {
	// 如果没有数据，返回零值
	if len(data) == 0 {
		return decimal.Zero
	}

	// 默认使用最新价格
	currentPrice := data[len(data)-1].Close.Decimal

	// 查找模式日期对应的价格
	for _, daily := range data {
		// 如果找到了完全匹配的日期
		if daily.TradeDate == patternDate {
			// 根据信号类型决定使用哪个价格
			switch signalType {
			case "BUY":
				// 买入信号通常在低点形成，使用当天的最低价
				return daily.Low.Decimal
			case "SELL":
				// 卖出信号通常在高点形成，使用当天的最高价
				return daily.High.Decimal
			default:
				// 其他情况使用收盘价
				return daily.Close.Decimal
			}
		}
	}

	// 如果没有找到完全匹配的日期，尝试找到最接近的日期
	// 这里简化处理，直接使用当前价格
	return currentPrice
}

// backTestPredictions 对预测进行回测
func (s *PredictionService) backTestPredictions(data []models.StockDaily, predictions []models.TradingPointPrediction) []models.TradingPointPrediction {
	// 如果数据不足，无法回测
	if len(data) <= 1 {
		return predictions
	}

	// 创建日期到价格的映射，方便快速查找
	dateToPrice := make(map[string]models.StockDaily)
	for _, daily := range data {
		dateToPrice[daily.TradeDate] = daily
	}

	// 对每个预测进行回测
	for i := range predictions {
		// 获取信号日期
		signalDate := predictions[i].SignalDate

		// 尝试找到信号日期后的第一个交易日
		nextDayData, nextDay := s.findNextTradingDay(data, signalDate)

		// 如果找到了下一个交易日，进行回测
		if nextDay != "" {
			predictions[i].Backtested = true
			predictions[i].NextDayPrice = models.NewJSONDecimal(nextDayData.Close.Decimal)

			// 计算价格差值
			priceDiff := nextDayData.Close.Decimal.Sub(predictions[i].Price.Decimal)
			predictions[i].PriceDiff = models.NewJSONDecimal(priceDiff)

			// 计算价格差值百分比
			if !predictions[i].Price.Decimal.IsZero() {
				priceDiffRatio := priceDiff.Div(predictions[i].Price.Decimal).Mul(decimal.NewFromInt(100))
				predictions[i].PriceDiffRatio = models.NewJSONDecimal(priceDiffRatio)
			} else {
				predictions[i].PriceDiffRatio = models.NewJSONDecimal(decimal.Zero)
			}

			// 判断预测是否正确
			switch predictions[i].Type {
			case "BUY":
				// 买入信号：如果第二天收盘价高于信号价格，则预测正确
				predictions[i].IsCorrect = nextDayData.Close.Decimal.GreaterThan(predictions[i].Price.Decimal)
			case "SELL":
				// 卖出信号：如果第二天收盘价低于信号价格，则预测正确
				predictions[i].IsCorrect = nextDayData.Close.Decimal.LessThan(predictions[i].Price.Decimal)
			}
		}
	}

	return predictions
}

// findNextTradingDay 查找指定日期之后的第一个交易日
func (s *PredictionService) findNextTradingDay(data []models.StockDaily, currentDate string) (models.StockDaily, string) {
	// 如果数据不足，无法查找
	if len(data) <= 1 {
		return models.StockDaily{}, ""
	}

	// 查找当前日期的索引
	currentIndex := -1
	for i, daily := range data {
		if daily.TradeDate == currentDate {
			currentIndex = i
			break
		}
	}

	// 如果找到了当前日期，且不是最后一天，返回下一个交易日
	if currentIndex >= 0 && currentIndex < len(data)-1 {
		return data[currentIndex+1], data[currentIndex+1].TradeDate
	}

	// 如果没有找到当前日期或者是最后一天，返回空
	return models.StockDaily{}, ""
}
