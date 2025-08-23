package service

import (
	"stock-a-future/internal/indicators"
	"stock-a-future/internal/models"
	"time"
)

// StockServiceInterface 股票服务接口
type StockServiceInterface interface {
	GetDailyData(tsCode, startDate, endDate, adjust string) ([]models.StockDaily, error)
}

// PatternService 图形识别服务
type PatternService struct {
	stockService      StockServiceInterface
	patternRecognizer *indicators.PatternRecognizer
}

// NewPatternService 创建新的图形识别服务
func NewPatternService(stockService StockServiceInterface) *PatternService {
	return &PatternService{
		stockService:      stockService,
		patternRecognizer: indicators.NewPatternRecognizer(),
	}
}

// RecognizePatterns 识别指定股票的图形模式
func (p *PatternService) RecognizePatterns(tsCode, startDate, endDate string) ([]models.PatternRecognitionResult, error) {
	// 获取股票数据 - 使用前复权(qfq)而不是none，因为AKTools不支持none参数
	stockData, err := p.stockService.GetDailyData(tsCode, startDate, endDate, "qfq")
	if err != nil {
		return nil, err
	}

	// 识别图形模式
	patterns := p.patternRecognizer.RecognizeAllPatterns(stockData)
	return patterns, nil
}

// SearchPatterns 搜索指定图形模式
func (p *PatternService) SearchPatterns(request models.PatternSearchRequest) (*models.PatternSearchResponse, error) {
	// 获取股票数据 - 使用前复权(qfq)而不是none，因为AKTools不支持none参数
	stockData, err := p.stockService.GetDailyData(request.TSCode, request.StartDate, request.EndDate, "qfq")
	if err != nil {
		return nil, err
	}

	// 识别所有图形模式
	allPatterns := p.patternRecognizer.RecognizeAllPatterns(stockData)

	// 预分配切片容量，避免频繁扩容
	filteredPatterns := make([]models.PatternRecognitionResult, 0, len(allPatterns))
	for _, pattern := range allPatterns {
		// 检查置信度
		if request.MinConfidence > 0 {
			confidence, _ := pattern.OverallConfidence.Decimal.Float64()
			if confidence < request.MinConfidence {
				continue
			}
		}

		// 检查图形类型
		if len(request.Patterns) > 0 {
			hasPattern := false
			for _, targetPattern := range request.Patterns {
				// 检查蜡烛图模式
				for _, candlestick := range pattern.Candlestick {
					if candlestick.Pattern == targetPattern {
						hasPattern = true
						break
					}
				}
				// 检查量价图形
				for _, volumePrice := range pattern.VolumePrice {
					if volumePrice.Pattern == targetPattern {
						hasPattern = true
						break
					}
				}
				if hasPattern {
					break
				}
			}
			if !hasPattern {
				continue
			}
		}

		filteredPatterns = append(filteredPatterns, pattern)
	}

	return &models.PatternSearchResponse{
		Total:   len(filteredPatterns),
		Results: filteredPatterns,
	}, nil
}

// GetPatternSummary 获取图形模式摘要
func (p *PatternService) GetPatternSummary(tsCode string, days int) (*models.PatternSummary, error) {
	// 计算日期范围
	endDate := time.Now().Format("20060102")
	startDate := time.Now().AddDate(0, 0, -days).Format("20060102")

	// 获取股票数据 - 使用前复权(qfq)而不是none，因为AKTools不支持none参数
	stockData, err := p.stockService.GetDailyData(tsCode, startDate, endDate, "qfq")
	if err != nil {
		return nil, err
	}

	// 识别图形模式
	patterns := p.patternRecognizer.RecognizeAllPatterns(stockData)

	// 统计各种图形模式
	summary := &models.PatternSummary{
		TSCode:    tsCode,
		Period:    days,
		StartDate: startDate,
		EndDate:   endDate,
		Patterns:  make(map[string]int),
		Signals:   make(map[string]int),
		UpdatedAt: time.Now(),
	}

	// 统计图形模式
	for _, pattern := range patterns {
		// 统计蜡烛图模式
		for _, candlestick := range pattern.Candlestick {
			summary.Patterns[candlestick.Pattern]++
			summary.Signals[candlestick.Signal]++
		}

		// 统计量价图形
		for _, volumePrice := range pattern.VolumePrice {
			summary.Patterns[volumePrice.Pattern]++
			summary.Signals[volumePrice.Signal]++
		}
	}

	return summary, nil
}

// GetRecentSignals 获取最近的图形信号
func (p *PatternService) GetRecentSignals(tsCode string, limit int) ([]models.RecentSignal, error) {
	// 获取最近30天的数据
	endDate := time.Now().Format("20060102")
	startDate := time.Now().AddDate(0, 0, -30).Format("20060102")

	// 获取股票数据 - 使用前复权(qfq)而不是none，因为AKTools不支持none参数
	stockData, err := p.stockService.GetDailyData(tsCode, startDate, endDate, "qfq")
	if err != nil {
		return nil, err
	}

	// 识别图形模式
	patterns := p.patternRecognizer.RecognizeAllPatterns(stockData)

	// 转换为最近信号格式
	var recentSignals []models.RecentSignal
	for _, pattern := range patterns {
		// 处理蜡烛图模式
		for _, candlestick := range pattern.Candlestick {
			recentSignals = append(recentSignals, models.RecentSignal{
				TSCode:      candlestick.TSCode,
				TradeDate:   candlestick.TradeDate,
				Pattern:     candlestick.Pattern,
				Signal:      candlestick.Signal,
				Confidence:  candlestick.Confidence,
				Description: candlestick.Description,
				Strength:    candlestick.Strength,
				Type:        "CANDLESTICK",
			})
		}

		// 处理量价图形
		for _, volumePrice := range pattern.VolumePrice {
			recentSignals = append(recentSignals, models.RecentSignal{
				TSCode:      volumePrice.TSCode,
				TradeDate:   volumePrice.TradeDate,
				Pattern:     volumePrice.Pattern,
				Signal:      volumePrice.Signal,
				Confidence:  volumePrice.Confidence,
				Description: volumePrice.Description,
				Strength:    volumePrice.Strength,
				Type:        "VOLUME_PRICE",
			})
		}
	}

	// 按日期排序（最新的在前）
	// 这里简化处理，实际应该按日期排序
	if len(recentSignals) > limit {
		recentSignals = recentSignals[:limit]
	}

	return recentSignals, nil
}
