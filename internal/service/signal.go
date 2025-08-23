package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"stock-a-future/internal/models"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// SignalService 异步信号计算服务
type SignalService struct {
	db              *DatabaseService
	patternService  *PatternService
	stockService    StockServiceInterface
	favoriteService *FavoriteService
	mutex           sync.RWMutex
	running         bool
	stopChan        chan struct{}

	// 信号计算状态
	statusMutex       sync.RWMutex
	calculationStatus *CalculationStatus
}

// NewSignalService 创建信号计算服务
func NewSignalService(dataDir string, patternService *PatternService, stockService StockServiceInterface, favoriteService *FavoriteService) (*SignalService, error) {
	// 创建数据库服务
	dbService, err := NewDatabaseService(dataDir)
	if err != nil {
		return nil, fmt.Errorf("创建数据库服务失败: %v", err)
	}

	service := &SignalService{
		db:              dbService,
		patternService:  patternService,
		stockService:    stockService,
		favoriteService: favoriteService,
		running:         false,
		stopChan:        make(chan struct{}),
		calculationStatus: &CalculationStatus{
			IsCalculating: false,
			Total:         0,
			Completed:     0,
			Failed:        0,
			StartTime:     time.Time{},
			EndTime:       time.Time{},
		},
	}

	return service, nil
}

// Start 启动异步信号计算服务
func (s *SignalService) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		log.Printf("信号计算服务已在运行中")
		return
	}

	s.running = true
	log.Printf("启动异步信号计算服务...")

	// 启动后台goroutine进行信号计算
	go s.runSignalCalculation()
}

// Stop 停止信号计算服务
func (s *SignalService) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		return
	}

	s.running = false
	close(s.stopChan)
	log.Printf("信号计算服务已停止")
}

// Close 关闭服务
func (s *SignalService) Close() error {
	s.Stop()
	return s.db.Close()
}

// runSignalCalculation 运行信号计算的后台任务
func (s *SignalService) runSignalCalculation() {
	log.Printf("开始计算收藏股票信号...")

	// 立即启动一个goroutine执行计算，不阻塞服务启动
	go s.calculateFavoriteStocksSignals()

	// 设置定时器，每天凌晨2点执行一次
	// 计算下一个凌晨2点的时间
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if now.After(nextRun) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	// 计算第一次执行的等待时间
	initialDelay := nextRun.Sub(now)
	log.Printf("下次定时计算将在 %v 后执行（%s）", initialDelay, nextRun.Format("2006-01-02 15:04:05"))

	// 首次延迟后启动定时器
	time.AfterFunc(initialDelay, func() {
		// 首次定时执行
		if s.shouldCalculateToday() {
			go s.calculateFavoriteStocksSignals()
		}

		// 然后设置每24小时执行一次的定时器
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-s.stopChan:
				log.Printf("信号计算服务收到停止信号")
				return
			case <-ticker.C:
				// 检查是否需要计算（避免重复计算）
				if s.shouldCalculateToday() {
					go s.calculateFavoriteStocksSignals()
				}
			}
		}
	})
}

// shouldCalculateToday 检查今天是否需要计算信号
func (s *SignalService) shouldCalculateToday() bool {
	today := time.Now().Format("20060102")

	// 查询今天是否已经计算过
	query := `
		SELECT COUNT(*) FROM stock_signals 
		WHERE signal_date = ?
	`

	var count int
	err := s.db.GetDB().QueryRow(query, today).Scan(&count)
	if err != nil {
		log.Printf("检查今日信号计算状态失败: %v", err)
		return true // 出错时默认需要计算
	}

	return count == 0 // 如果今天没有计算过，则需要计算
}

// CalculationStatus 表示信号计算的状态
type CalculationStatus struct {
	IsCalculating bool      `json:"is_calculating"`
	Total         int       `json:"total"`
	Completed     int       `json:"completed"`
	Failed        int       `json:"failed"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Duration      string    `json:"duration,omitempty"`
}

// GetCalculationStatus 获取当前计算状态
func (s *SignalService) GetCalculationStatus() *CalculationStatus {
	s.statusMutex.RLock()
	defer s.statusMutex.RUnlock()

	status := &CalculationStatus{
		IsCalculating: s.calculationStatus.IsCalculating,
		Total:         s.calculationStatus.Total,
		Completed:     s.calculationStatus.Completed,
		Failed:        s.calculationStatus.Failed,
		StartTime:     s.calculationStatus.StartTime,
		EndTime:       s.calculationStatus.EndTime,
	}

	if !s.calculationStatus.StartTime.IsZero() {
		if s.calculationStatus.IsCalculating {
			status.Duration = time.Since(s.calculationStatus.StartTime).String()
		} else if !s.calculationStatus.EndTime.IsZero() {
			status.Duration = s.calculationStatus.EndTime.Sub(s.calculationStatus.StartTime).String()
		}
	}

	return status
}

// calculateFavoriteStocksSignals 计算收藏股票的信号
func (s *SignalService) calculateFavoriteStocksSignals() {
	log.Printf("开始计算收藏股票信号...")
	startTime := time.Now()

	// 更新计算状态
	s.statusMutex.Lock()
	s.calculationStatus.IsCalculating = true
	s.calculationStatus.StartTime = startTime
	s.calculationStatus.Completed = 0
	s.calculationStatus.Failed = 0
	s.statusMutex.Unlock()

	// 获取所有收藏股票
	favorites := s.favoriteService.GetFavorites()
	if len(favorites) == 0 {
		log.Printf("没有收藏股票，跳过信号计算")

		// 更新计算状态
		s.statusMutex.Lock()
		s.calculationStatus.IsCalculating = false
		s.calculationStatus.EndTime = time.Now()
		s.calculationStatus.Total = 0
		s.statusMutex.Unlock()

		return
	}

	total := len(favorites)

	// 更新总数
	s.statusMutex.Lock()
	s.calculationStatus.Total = total
	s.statusMutex.Unlock()

	// 并发计算信号
	semaphore := make(chan struct{}, 5) // 限制并发数为5
	var wg sync.WaitGroup

	for _, favorite := range favorites {
		wg.Add(1)
		go func(fav *models.FavoriteStock) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			if err := s.calculateStockSignal(fav.TSCode, fav.Name, false); err != nil {
				log.Printf("计算股票 %s 信号失败: %v", fav.TSCode, err)

				// 更新失败计数
				s.statusMutex.Lock()
				s.calculationStatus.Failed++
				s.statusMutex.Unlock()
			} else {
				// 更新成功计数
				s.statusMutex.Lock()
				s.calculationStatus.Completed++
				s.statusMutex.Unlock()
			}
		}(favorite)
	}

	wg.Wait()

	// 计算完成，更新状态
	s.statusMutex.Lock()
	s.calculationStatus.IsCalculating = false
	s.calculationStatus.EndTime = time.Now()
	duration := s.calculationStatus.EndTime.Sub(s.calculationStatus.StartTime)
	s.statusMutex.Unlock()

	log.Printf("收藏股票信号计算完成: 总数=%d, 成功=%d, 失败=%d, 耗时=%v",
		total, s.calculationStatus.Completed, s.calculationStatus.Failed, duration)
}

// CalculateStockSignal 计算单个股票信号（公开方法）
func (s *SignalService) CalculateStockSignal(tsCode, name string, force bool) error {
	return s.calculateStockSignal(tsCode, name, force)
}

// calculateStockSignal 计算单个股票的信号
func (s *SignalService) calculateStockSignal(tsCode, name string, force bool) error {
	today := time.Now().Format("20060102")

	// 检查是否已经计算过今天的信号
	if !force && s.hasSignalForToday(tsCode, today) {
		log.Printf("股票 %s 今日信号已存在，跳过计算", tsCode)
		return nil
	}

	// 获取股票名称
	name = s.getStockName(tsCode, name)

	// 获取股票数据
	stockData, err := s.getStockDataForAnalysis(tsCode)
	if err != nil {
		return err
	}

	// 计算信号
	signal, err := s.computeStockSignal(tsCode, name, stockData, today)
	if err != nil {
		return err
	}

	// 保存信号到数据库
	if err := s.saveSignalToDB(signal); err != nil {
		return fmt.Errorf("保存信号到数据库失败: %v", err)
	}

	log.Printf("股票 %s 信号计算完成: 类型=%s, 强度=%s, 置信度=%.2f",
		tsCode, signal.SignalType, signal.SignalStrength, signal.Confidence.Decimal.InexactFloat64())

	return nil
}

// getStockName 获取股票名称
func (s *SignalService) getStockName(tsCode, name string) string {
	if name != "" {
		return name
	}

	// 尝试从本地数据库获取股票名称
	stockInfo, err := s.getStockInfoFromDB(tsCode)
	if err == nil && stockInfo.Name != "" {
		log.Printf("从数据库获取到股票名称: %s -> %s", tsCode, stockInfo.Name)
		return stockInfo.Name
	}

	// 使用代码作为名称
	log.Printf("未能获取股票名称，使用代码作为名称: %s", tsCode)
	return tsCode
}

// getStockDataForAnalysis 获取用于分析的股票数据
func (s *SignalService) getStockDataForAnalysis(tsCode string) ([]models.StockDaily, error) {
	// 获取最近30天的数据进行分析
	endDate := time.Now().Format("20060102")
	startDate := time.Now().AddDate(0, 0, -30).Format("20060102")

	// 获取股票数据 - 使用前复权(qfq)而不是none，因为AKTools不支持none参数
	stockData, err := s.stockService.GetDailyData(tsCode, startDate, endDate, "qfq")
	if err != nil {
		log.Printf("获取股票 %s 数据详细错误: %v", tsCode, err)
		return nil, fmt.Errorf("获取股票数据失败: %v", err)
	}

	if len(stockData) == 0 {
		return nil, fmt.Errorf("没有获取到股票数据")
	}

	return stockData, nil
}

// computeStockSignal 计算股票信号
func (s *SignalService) computeStockSignal(tsCode, name string, stockData []models.StockDaily, today string) (*models.StockSignal, error) {
	// 获取最新的交易日期
	latestData := stockData[len(stockData)-1]
	tradeDate := latestData.TradeDate

	// 计算图形模式
	patterns, err := s.patternService.RecognizePatterns(tsCode, tradeDate, tradeDate)
	if err != nil {
		log.Printf("识别图形模式失败: %v", err)
		patterns = []models.PatternRecognitionResult{} // 继续处理，但无图形模式
	}

	// 计算技术指标（简化版本）
	indicators := s.calculateTechnicalIndicators(stockData)

	// 生成预测（简化版本）
	predictions := s.generatePredictions(stockData, patterns)

	// 综合分析生成信号
	signal := s.generateSignal(stockData, patterns, indicators, predictions)
	signal.ID = s.generateID()
	signal.TSCode = tsCode
	signal.Name = name
	signal.TradeDate = tradeDate
	signal.SignalDate = today
	signal.CreatedAt = time.Now()
	signal.UpdatedAt = time.Now()

	return signal, nil
}

// getStockInfoFromDB 从数据库获取股票信息
func (s *SignalService) getStockInfoFromDB(tsCode string) (*models.StockBasic, error) {
	query := `
		SELECT ts_code, symbol, name, area, industry, market, list_date
		FROM stock_basic 
		WHERE ts_code = ?
	`

	var stock models.StockBasic
	err := s.db.GetDB().QueryRow(query, tsCode).Scan(
		&stock.TSCode,
		&stock.Symbol,
		&stock.Name,
		&stock.Area,
		&stock.Industry,
		&stock.Market,
		&stock.ListDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("未找到股票信息")
		}
		return nil, err
	}

	return &stock, nil
}

// hasSignalForToday 检查今天是否已有信号
func (s *SignalService) hasSignalForToday(tsCode, signalDate string) bool {
	query := `
		SELECT COUNT(*) FROM stock_signals 
		WHERE ts_code = ? AND signal_date = ?
	`

	var count int
	err := s.db.GetDB().QueryRow(query, tsCode, signalDate).Scan(&count)
	if err != nil {
		log.Printf("检查信号存在性失败: %v", err)
		return false
	}

	return count > 0
}

// generateID 生成唯一ID
func (s *SignalService) generateID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// calculateTechnicalIndicators 计算技术指标（简化版本）
func (s *SignalService) calculateTechnicalIndicators(stockData []models.StockDaily) *models.TechnicalIndicators {
	if len(stockData) == 0 {
		return nil
	}

	latest := stockData[len(stockData)-1]

	// 计算简单移动平均线
	ma5 := s.calculateMA(stockData, 5)
	ma10 := s.calculateMA(stockData, 10)
	ma20 := s.calculateMA(stockData, 20)

	return &models.TechnicalIndicators{
		TSCode:    latest.TSCode,
		TradeDate: latest.TradeDate,
		MA: &models.MovingAverageIndicator{
			MA5:  models.NewJSONDecimal(ma5),
			MA10: models.NewJSONDecimal(ma10),
			MA20: models.NewJSONDecimal(ma20),
		},
	}
}

// calculateMA 计算移动平均线
func (s *SignalService) calculateMA(data []models.StockDaily, period int) decimal.Decimal {
	if len(data) < period {
		return decimal.Zero
	}

	sum := decimal.Zero
	for i := len(data) - period; i < len(data); i++ {
		sum = sum.Add(data[i].Close.Decimal)
	}

	return sum.Div(decimal.NewFromInt(int64(period)))
}

// generatePredictions 生成预测（简化版本）
func (s *SignalService) generatePredictions(stockData []models.StockDaily, _ []models.PatternRecognitionResult) *models.PredictionResult {
	if len(stockData) == 0 {
		return nil
	}

	latest := stockData[len(stockData)-1]

	return &models.PredictionResult{
		TSCode:     latest.TSCode,
		TradeDate:  latest.TradeDate,
		Confidence: models.NewJSONDecimal(decimal.NewFromFloat(0.6)), // 简化的置信度
		UpdatedAt:  time.Now(),
	}
}

// generateSignal 生成综合信号
func (s *SignalService) generateSignal(stockData []models.StockDaily, patterns []models.PatternRecognitionResult, indicators *models.TechnicalIndicators, predictions *models.PredictionResult) *models.StockSignal {
	if len(stockData) == 0 {
		return &models.StockSignal{
			SignalType:     "HOLD",
			SignalStrength: "WEAK",
			Confidence:     models.NewJSONDecimal(decimal.NewFromFloat(0.3)),
			Description:    "数据不足，建议观望",
		}
	}

	latest := stockData[len(stockData)-1]

	// 简化的信号生成逻辑
	signalType := "HOLD"
	signalStrength := "MEDIUM"
	confidence := decimal.NewFromFloat(0.5)
	description := "基于技术分析的综合判断"

	// 基于图形模式的信号判断
	buySignals := 0
	sellSignals := 0
	totalConfidence := decimal.Zero

	for _, pattern := range patterns {
		for _, candlestick := range pattern.Candlestick {
			switch candlestick.Signal {
			case "BUY":
				buySignals++
			case "SELL":
				sellSignals++
			}
			totalConfidence = totalConfidence.Add(candlestick.Confidence.Decimal)
		}

		for _, volumePrice := range pattern.VolumePrice {
			switch volumePrice.Signal {
			case "BUY":
				buySignals++
			case "SELL":
				sellSignals++
			}
			totalConfidence = totalConfidence.Add(volumePrice.Confidence.Decimal)
		}
	}

	// 基于移动平均线的趋势判断
	if indicators != nil && indicators.MA != nil {
		currentPrice := latest.Close.Decimal
		ma5 := indicators.MA.MA5.Decimal
		ma10 := indicators.MA.MA10.Decimal

		if currentPrice.GreaterThan(ma5) && ma5.GreaterThan(ma10) {
			buySignals++
		} else if currentPrice.LessThan(ma5) && ma5.LessThan(ma10) {
			sellSignals++
		}
	}

	// 综合判断
	if buySignals > sellSignals {
		signalType = "BUY"
		if buySignals >= 3 {
			signalStrength = "STRONG"
			confidence = decimal.NewFromFloat(0.8)
		} else {
			signalStrength = "MEDIUM"
			confidence = decimal.NewFromFloat(0.6)
		}
		description = fmt.Sprintf("检测到 %d 个买入信号，建议关注", buySignals)
	} else if sellSignals > buySignals {
		signalType = "SELL"
		if sellSignals >= 3 {
			signalStrength = "STRONG"
			confidence = decimal.NewFromFloat(0.8)
		} else {
			signalStrength = "MEDIUM"
			confidence = decimal.NewFromFloat(0.6)
		}
		description = fmt.Sprintf("检测到 %d 个卖出信号，建议谨慎", sellSignals)
	}

	return &models.StockSignal{
		SignalType:          signalType,
		SignalStrength:      signalStrength,
		Confidence:          models.NewJSONDecimal(confidence),
		Patterns:            patterns,
		TechnicalIndicators: indicators,
		Predictions:         predictions,
		Description:         description,
	}
}

// saveSignalToDB 保存信号到数据库
func (s *SignalService) saveSignalToDB(signal *models.StockSignal) error {
	// 序列化复杂字段
	patternsJSON, _ := json.Marshal(signal.Patterns)
	indicatorsJSON, _ := json.Marshal(signal.TechnicalIndicators)
	predictionsJSON, _ := json.Marshal(signal.Predictions)

	stmt := `
		INSERT OR REPLACE INTO stock_signals 
		(id, ts_code, name, trade_date, signal_date, signal_type, signal_strength, 
		 confidence, patterns, technical_indicators, predictions, description, 
		 created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.GetDB().Exec(stmt,
		signal.ID,
		signal.TSCode,
		signal.Name,
		signal.TradeDate,
		signal.SignalDate,
		signal.SignalType,
		signal.SignalStrength,
		signal.Confidence.Decimal.InexactFloat64(),
		string(patternsJSON),
		string(indicatorsJSON),
		string(predictionsJSON),
		signal.Description,
		signal.CreatedAt,
		signal.UpdatedAt,
	)

	return err
}

// GetSignal 获取股票信号
func (s *SignalService) GetSignal(tsCode, signalDate string) (*models.StockSignal, error) {
	query := `
		SELECT id, ts_code, name, trade_date, signal_date, signal_type, signal_strength,
		       confidence, patterns, technical_indicators, predictions, description,
		       created_at, updated_at
		FROM stock_signals 
		WHERE ts_code = ? AND signal_date = ?
	`

	var signal models.StockSignal
	var patternsJSON, indicatorsJSON, predictionsJSON string
	var confidence float64

	err := s.db.GetDB().QueryRow(query, tsCode, signalDate).Scan(
		&signal.ID,
		&signal.TSCode,
		&signal.Name,
		&signal.TradeDate,
		&signal.SignalDate,
		&signal.SignalType,
		&signal.SignalStrength,
		&confidence,
		&patternsJSON,
		&indicatorsJSON,
		&predictionsJSON,
		&signal.Description,
		&signal.CreatedAt,
		&signal.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("未找到信号记录")
		}
		return nil, err
	}

	signal.Confidence = models.NewJSONDecimal(decimal.NewFromFloat(confidence))

	// 反序列化复杂字段
	if patternsJSON != "" {
		json.Unmarshal([]byte(patternsJSON), &signal.Patterns)
	}
	if indicatorsJSON != "" {
		json.Unmarshal([]byte(indicatorsJSON), &signal.TechnicalIndicators)
	}
	if predictionsJSON != "" {
		json.Unmarshal([]byte(predictionsJSON), &signal.Predictions)
	}

	return &signal, nil
}

// GetLatestSignals 获取最新的信号列表
func (s *SignalService) GetLatestSignals(limit int) ([]*models.StockSignal, error) {
	query := `
		SELECT id, ts_code, name, trade_date, signal_date, signal_type, signal_strength,
		       confidence, patterns, technical_indicators, predictions, description,
		       created_at, updated_at
		FROM stock_signals 
		ORDER BY signal_date DESC, updated_at DESC
		LIMIT ?
	`

	rows, err := s.db.GetDB().Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signals []*models.StockSignal
	for rows.Next() {
		var signal models.StockSignal
		var patternsJSON, indicatorsJSON, predictionsJSON string
		var confidence float64

		err := rows.Scan(
			&signal.ID,
			&signal.TSCode,
			&signal.Name,
			&signal.TradeDate,
			&signal.SignalDate,
			&signal.SignalType,
			&signal.SignalStrength,
			&confidence,
			&patternsJSON,
			&indicatorsJSON,
			&predictionsJSON,
			&signal.Description,
			&signal.CreatedAt,
			&signal.UpdatedAt,
		)
		if err != nil {
			continue
		}

		signal.Confidence = models.NewJSONDecimal(decimal.NewFromFloat(confidence))

		// 反序列化复杂字段
		if patternsJSON != "" {
			json.Unmarshal([]byte(patternsJSON), &signal.Patterns)
		}
		if indicatorsJSON != "" {
			json.Unmarshal([]byte(indicatorsJSON), &signal.TechnicalIndicators)
		}
		if predictionsJSON != "" {
			json.Unmarshal([]byte(predictionsJSON), &signal.Predictions)
		}

		signals = append(signals, &signal)
	}

	return signals, nil
}
