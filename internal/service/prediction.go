package service

import (
	"fmt"
	"stock-a-future/internal/indicators"
	"stock-a-future/internal/models"
	"time"

	"github.com/shopspring/decimal"
)

// PredictionService 预测服务
type PredictionService struct {
	calculator *indicators.Calculator
}

// NewPredictionService 创建预测服务
func NewPredictionService() *PredictionService {
	return &PredictionService{
		calculator: indicators.NewCalculator(),
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
			MA5:  ma5[len(ma5)-1],
			MA10: ma10[len(ma10)-1],
			MA20: ma20[len(ma20)-1],
		}
		if len(ma60) > 0 {
			indicators.MA.MA60 = ma60[len(ma60)-1]
		}
		if len(ma120) > 0 {
			indicators.MA.MA120 = ma120[len(ma120)-1]
		}
	}

	// 计算KDJ
	kdjResults := s.calculator.CalculateKDJ(data, 9)
	if len(kdjResults) > 0 {
		indicators.KDJ = &kdjResults[len(kdjResults)-1]
	}

	return indicators
}

// generatePredictions 基于技术指标生成预测
func (s *PredictionService) generatePredictions(data []models.StockDaily, indicators *models.TechnicalIndicators) []models.TradingPointPrediction {
	var predictions []models.TradingPointPrediction
	currentPrice := data[len(data)-1].Close

	// 基于MACD的预测
	if indicators.MACD != nil {
		if prediction := s.predictFromMACD(indicators.MACD, currentPrice); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	// 基于RSI的预测
	if indicators.RSI != nil {
		if prediction := s.predictFromRSI(indicators.RSI, currentPrice); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	// 基于布林带的预测
	if indicators.BOLL != nil {
		if prediction := s.predictFromBollingerBands(indicators.BOLL, currentPrice); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	// 基于KDJ的预测
	if indicators.KDJ != nil {
		if prediction := s.predictFromKDJ(indicators.KDJ, currentPrice); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	// 基于移动平均线的预测
	if indicators.MA != nil {
		if prediction := s.predictFromMA(indicators.MA, currentPrice); prediction != nil {
			predictions = append(predictions, *prediction)
		}
	}

	return predictions
}

// predictFromMACD 基于MACD指标预测
func (s *PredictionService) predictFromMACD(macd *models.MACDIndicator, currentPrice decimal.Decimal) *models.TradingPointPrediction {
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
	macdStrength := macd.MACD.Abs()
	if macdStrength.GreaterThan(decimal.NewFromFloat(0.1)) {
		probability = probability.Add(decimal.NewFromFloat(0.1))
	}

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       currentPrice,
		Date:        time.Now().AddDate(0, 0, 1).Format("20060102"), // 预测明天
		Probability: probability,
		Reason:      reason,
		Indicators:  []string{"MACD"},
	}
}

// predictFromRSI 基于RSI指标预测
func (s *PredictionService) predictFromRSI(rsi *models.RSIIndicator, currentPrice decimal.Decimal) *models.TradingPointPrediction {
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
		reason = fmt.Sprintf("RSI超卖信号，当前RSI值：%.2f", rsi.RSI12.InexactFloat64())
	case "SELL":
		predictType = "SELL"
		probability = decimal.NewFromFloat(0.65)
		reason = fmt.Sprintf("RSI超买信号，当前RSI值：%.2f", rsi.RSI12.InexactFloat64())
	default:
		return nil
	}

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       currentPrice,
		Date:        time.Now().AddDate(0, 0, 1).Format("20060102"),
		Probability: probability,
		Reason:      reason,
		Indicators:  []string{"RSI"},
	}
}

// predictFromBollingerBands 基于布林带预测
func (s *PredictionService) predictFromBollingerBands(boll *models.BollingerBandsIndicator, currentPrice decimal.Decimal) *models.TradingPointPrediction {
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

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       currentPrice,
		Date:        time.Now().AddDate(0, 0, 2).Format("20060102"), // 布林带信号预测后天
		Probability: probability,
		Reason:      reason,
		Indicators:  []string{"BOLL"},
	}
}

// predictFromKDJ 基于KDJ指标预测
func (s *PredictionService) predictFromKDJ(kdj *models.KDJIndicator, currentPrice decimal.Decimal) *models.TradingPointPrediction {
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
		reason = fmt.Sprintf("KDJ超卖信号，K值：%.2f，D值：%.2f", kdj.K.InexactFloat64(), kdj.D.InexactFloat64())
	case "SELL":
		predictType = "SELL"
		probability = decimal.NewFromFloat(0.55)
		reason = fmt.Sprintf("KDJ超买信号，K值：%.2f，D值：%.2f", kdj.K.InexactFloat64(), kdj.D.InexactFloat64())
	default:
		return nil
	}

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       currentPrice,
		Date:        time.Now().AddDate(0, 0, 1).Format("20060102"),
		Probability: probability,
		Reason:      reason,
		Indicators:  []string{"KDJ"},
	}
}

// predictFromMA 基于移动平均线预测
func (s *PredictionService) predictFromMA(ma *models.MovingAverageIndicator, currentPrice decimal.Decimal) *models.TradingPointPrediction {
	// 判断均线多头/空头排列
	var predictType string
	var probability decimal.Decimal
	var reason string

	// 多头排列：MA5 > MA10 > MA20
	if ma.MA5.GreaterThan(ma.MA10) && ma.MA10.GreaterThan(ma.MA20) {
		if currentPrice.GreaterThan(ma.MA5) {
			predictType = "BUY"
			probability = decimal.NewFromFloat(0.62)
			reason = "均线多头排列，价格在5日均线之上"
		}
	}

	// 空头排列：MA5 < MA10 < MA20
	if ma.MA5.LessThan(ma.MA10) && ma.MA10.LessThan(ma.MA20) {
		if currentPrice.LessThan(ma.MA5) {
			predictType = "SELL"
			probability = decimal.NewFromFloat(0.58)
			reason = "均线空头排列，价格在5日均线之下"
		}
	}

	if predictType == "" {
		return nil
	}

	return &models.TradingPointPrediction{
		Type:        predictType,
		Price:       currentPrice,
		Date:        time.Now().AddDate(0, 0, 3).Format("20060102"), // 均线信号预测3天后
		Probability: probability,
		Reason:      reason,
		Indicators:  []string{"MA"},
	}
}

// calculateOverallConfidence 计算整体置信度
func (s *PredictionService) calculateOverallConfidence(predictions []models.TradingPointPrediction) decimal.Decimal {
	if len(predictions) == 0 {
		return decimal.Zero
	}

	// 统计买卖信号
	buyCount := 0
	sellCount := 0
	totalProbability := decimal.Zero

	for _, pred := range predictions {
		totalProbability = totalProbability.Add(pred.Probability)
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

	return confidence
}
