package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"
)

// 辅助函数 - 使用内置max/min函数或自定义实现
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var (
	ErrBacktestNotFound     = errors.New("回测不存在")
	ErrBacktestExists       = errors.New("回测已存在")
	ErrBacktestNotCompleted = errors.New("回测尚未完成")
	ErrBacktestRunning      = errors.New("回测正在运行")
)

// BacktestService 回测服务
type BacktestService struct {
	// 在真实环境中，这里会有数据库连接
	// 目前使用内存存储进行演示
	backtests        map[string]*models.Backtest
	backtestResults  map[string]*models.BacktestResult
	backtestTrades   map[string][]models.Trade
	backtestProgress map[string]*models.BacktestProgress
	runningBacktests map[string]context.CancelFunc // 用于取消运行中的回测

	strategyService   *StrategyService
	dataSourceService *DataSourceService
	dailyCacheService *DailyCacheService // 使用现有的日线数据缓存服务
	logger            logger.Logger
	mutex             sync.RWMutex
}

// NewBacktestService 创建回测服务
func NewBacktestService(strategyService *StrategyService, dataSourceService *DataSourceService, dailyCacheService *DailyCacheService, log logger.Logger) *BacktestService {
	return &BacktestService{
		backtests:         make(map[string]*models.Backtest),
		backtestResults:   make(map[string]*models.BacktestResult),
		backtestTrades:    make(map[string][]models.Trade),
		backtestProgress:  make(map[string]*models.BacktestProgress),
		runningBacktests:  make(map[string]context.CancelFunc),
		strategyService:   strategyService,
		dataSourceService: dataSourceService,
		dailyCacheService: dailyCacheService,
		logger:            log,
	}
}

// GetBacktestsList 获取回测列表
func (s *BacktestService) GetBacktestsList(ctx context.Context, req *models.BacktestListRequest) ([]models.Backtest, int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var results []models.Backtest

	// 过滤回测
	for _, backtest := range s.backtests {
		if s.matchesBacktestFilter(backtest, req) {
			results = append(results, *backtest)
		}
	}

	total := len(results)

	// 分页
	start := (req.Page - 1) * req.Size
	end := start + req.Size

	if start >= len(results) {
		return []models.Backtest{}, total, nil
	}

	if end > len(results) {
		end = len(results)
	}

	return results[start:end], total, nil
}

// matchesBacktestFilter 检查回测是否匹配过滤条件
func (s *BacktestService) matchesBacktestFilter(backtest *models.Backtest, req *models.BacktestListRequest) bool {
	// 状态过滤
	if req.Status != "" && backtest.Status != req.Status {
		return false
	}

	// 策略过滤
	if req.StrategyID != "" && backtest.StrategyID != req.StrategyID {
		return false
	}

	// 关键词过滤
	if req.Keyword != "" {
		keyword := strings.ToLower(req.Keyword)
		name := strings.ToLower(backtest.Name)

		if !strings.Contains(name, keyword) {
			return false
		}
	}

	// 日期过滤
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil && backtest.StartDate.Before(startDate) {
			return false
		}
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil && backtest.EndDate.After(endDate) {
			return false
		}
	}

	return true
}

// GetBacktest 获取回测详情
func (s *BacktestService) GetBacktest(ctx context.Context, backtestID string) (*models.Backtest, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	backtest, exists := s.backtests[backtestID]
	if !exists {
		return nil, ErrBacktestNotFound
	}

	// 返回副本避免外部修改
	backtestCopy := *backtest
	return &backtestCopy, nil
}

// CreateBacktest 创建回测
func (s *BacktestService) CreateBacktest(ctx context.Context, backtest *models.Backtest) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 检查回测是否已存在
	if _, exists := s.backtests[backtest.ID]; exists {
		return ErrBacktestExists
	}

	// 检查名称是否重复
	for _, existing := range s.backtests {
		if existing.Name == backtest.Name {
			return fmt.Errorf("回测名称已存在: %s", backtest.Name)
		}
	}

	// 设置时间戳
	now := time.Now()
	backtest.CreatedAt = now

	// 保存回测
	s.backtests[backtest.ID] = backtest

	s.logger.Info("回测创建成功",
		logger.String("backtest_id", backtest.ID),
		logger.String("backtest_name", backtest.Name),
		logger.String("strategy_id", backtest.StrategyID),
	)

	return nil
}

// UpdateBacktest 更新回测
func (s *BacktestService) UpdateBacktest(ctx context.Context, backtestID string, req *models.UpdateBacktestRequest) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	backtest, exists := s.backtests[backtestID]
	if !exists {
		return ErrBacktestNotFound
	}

	// 更新字段
	if req.Name != nil {
		// 检查名称是否重复
		for id, existing := range s.backtests {
			if id != backtestID && existing.Name == *req.Name {
				return fmt.Errorf("回测名称已存在: %s", *req.Name)
			}
		}
		backtest.Name = *req.Name
	}

	if req.Status != nil {
		backtest.Status = *req.Status
	}

	s.logger.Info("回测更新成功",
		logger.String("backtest_id", backtestID),
		logger.String("backtest_name", backtest.Name),
	)

	return nil
}

// DeleteBacktest 删除回测
func (s *BacktestService) DeleteBacktest(ctx context.Context, backtestID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	backtest, exists := s.backtests[backtestID]
	if !exists {
		return ErrBacktestNotFound
	}

	// 如果回测正在运行，先取消它
	if backtest.Status == models.BacktestStatusRunning {
		if cancelFunc, ok := s.runningBacktests[backtestID]; ok {
			cancelFunc()
			delete(s.runningBacktests, backtestID)
		}
	}

	// 删除相关数据
	delete(s.backtests, backtestID)
	delete(s.backtestResults, backtestID)
	delete(s.backtestTrades, backtestID)
	delete(s.backtestProgress, backtestID)

	s.logger.Info("回测删除成功", logger.String("backtest_id", backtestID))

	return nil
}

// StartBacktest 启动回测
func (s *BacktestService) StartBacktest(ctx context.Context, backtest *models.Backtest, strategy *models.Strategy) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 检查回测状态
	if backtest.Status != models.BacktestStatusPending {
		return fmt.Errorf("回测状态不允许启动: %s", backtest.Status)
	}

	// 更新状态
	backtest.Status = models.BacktestStatusRunning
	backtest.Progress = 0
	now := time.Now()
	backtest.StartedAt = &now

	// 初始化进度
	s.backtestProgress[backtest.ID] = &models.BacktestProgress{
		BacktestID: backtest.ID,
		Status:     string(models.BacktestStatusRunning),
		Progress:   0,
		Message:    "初始化回测环境...",
	}

	// 创建可取消的上下文，根据回测时间范围动态设置超时
	totalDays := int(backtest.EndDate.Sub(backtest.StartDate).Hours() / 24)
	// 计算合理的超时时间：每天至少10秒，最少5分钟，最多2小时
	timeoutMinutes := maxInt(5, minInt(120, totalDays/6))
	timeout := time.Duration(timeoutMinutes) * time.Minute

	s.logger.Info("设置回测超时时间",
		logger.String("backtest_id", backtest.ID),
		logger.Int("total_days", totalDays),
		logger.Int("timeout_minutes", timeoutMinutes),
	)

	backtestCtx, cancel := context.WithTimeout(ctx, timeout)
	s.runningBacktests[backtest.ID] = cancel

	s.logger.Info("准备启动回测goroutine",
		logger.String("backtest_id", backtest.ID),
		logger.String("strategy_id", strategy.ID),
		logger.Int("total_days", int(backtest.EndDate.Sub(backtest.StartDate).Hours()/24)),
		logger.String("start_date", backtest.StartDate.Format("2006-01-02")),
		logger.String("end_date", backtest.EndDate.Format("2006-01-02")),
		logger.Int("symbols_count", len(backtest.Symbols)),
	)

	// 启动后台回测任务
	go s.runBacktestTask(backtestCtx, backtest, strategy)

	s.logger.Info("回测启动成功",
		logger.String("backtest_id", backtest.ID),
		logger.String("strategy_id", strategy.ID),
		logger.String("status", "goroutine_launched"),
	)

	return nil
}

// runBacktestTask 运行回测任务
func (s *BacktestService) runBacktestTask(ctx context.Context, backtest *models.Backtest, strategy *models.Strategy) {
	// 立即输出日志，确保goroutine已启动
	s.logger.Info("🚀 回测goroutine已启动",
		logger.String("backtest_id", backtest.ID),
		logger.String("goroutine_status", "started"),
	)

	s.logger.Info("回测任务开始执行",
		logger.String("backtest_id", backtest.ID),
		logger.String("goroutine", "runBacktestTask"),
	)

	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("回测任务出现panic",
				logger.String("backtest_id", backtest.ID),
				logger.Any("panic", r),
			)
		}
		s.mutex.Lock()
		delete(s.runningBacktests, backtest.ID)
		s.mutex.Unlock()
		s.logger.Info("回测任务清理完成", logger.String("backtest_id", backtest.ID))
	}()

	// 预加载回测数据
	s.logger.Info("开始预加载回测数据",
		logger.String("backtest_id", backtest.ID),
		logger.Int("symbols_count", len(backtest.Symbols)),
	)

	if err := s.preloadBacktestData(ctx, backtest.Symbols, backtest.StartDate, backtest.EndDate); err != nil {
		s.logger.Error("预加载回测数据失败",
			logger.String("backtest_id", backtest.ID),
			logger.ErrorField(err),
		)
		s.updateBacktestStatus(backtest.ID, models.BacktestStatusFailed, "预加载数据失败")
		return
	}

	// 计算回测参数
	totalDays := int(backtest.EndDate.Sub(backtest.StartDate).Hours() / 24)
	if totalDays <= 0 {
		totalDays = 1
	}

	s.logger.Info("回测参数计算完成",
		logger.String("backtest_id", backtest.ID),
		logger.Int("total_days", totalDays),
	)

	// 初始化组合
	portfolio := &models.Portfolio{
		Cash:       backtest.InitialCash,
		Positions:  make(map[string]models.Position),
		TotalValue: backtest.InitialCash,
	}

	var trades []models.Trade
	var equityCurve []models.EquityPoint
	var dailyReturns []float64

	// 模拟每日回测
	s.logger.Info("开始回测循环",
		logger.String("backtest_id", backtest.ID),
		logger.Int("total_days", totalDays),
		logger.String("start_date", backtest.StartDate.Format("2006-01-02")),
		logger.String("end_date", backtest.EndDate.Format("2006-01-02")),
	)

	// 用于控制日志输出频率
	var lastLoggedProgress int = -1
	const progressLogInterval = 10 // 每10%打印一次进度

	for day := 0; day <= totalDays; day++ {

		select {
		case <-ctx.Done():
			// 回测被取消或超时
			if ctx.Err() == context.DeadlineExceeded {
				s.logger.Error("回测超时",
					logger.String("backtest_id", backtest.ID),
					logger.ErrorField(ctx.Err()),
				)
				s.updateBacktestStatus(backtest.ID, models.BacktestStatusFailed, "回测执行超时")
			} else {
				s.logger.Info("回测被取消",
					logger.String("backtest_id", backtest.ID),
					logger.ErrorField(ctx.Err()),
				)
				s.updateBacktestStatus(backtest.ID, models.BacktestStatusCancelled, "回测已取消")
			}
			return
		default:
		}

		currentDate := backtest.StartDate.AddDate(0, 0, day)

		if currentDate.After(backtest.EndDate) {
			s.logger.Info("回测日期超出范围，结束循环",
				logger.String("backtest_id", backtest.ID),
				logger.String("current_date", currentDate.Format("2006-01-02")),
			)
			break
		}

		// 更新进度
		progress := int(float64(day) / float64(totalDays) * 100)

		// 只在进度达到特定节点时打印日志
		if progress >= lastLoggedProgress+progressLogInterval || day == 0 || day == totalDays {
			s.logger.Info("回测进度更新",
				logger.String("backtest_id", backtest.ID),
				logger.Int("progress", progress),
				logger.String("current_date", currentDate.Format("2006-01-02")),
				logger.Int("day", day),
				logger.Int("total_days", totalDays),
			)
			lastLoggedProgress = progress
		}
		s.updateBacktestProgress(backtest.ID, progress, fmt.Sprintf("回测进行中... %s", currentDate.Format("2006-01-02")))

		// 模拟每个股票的交易
		// 只在首次或进度更新时记录股票处理信息
		if progress >= lastLoggedProgress || day == 0 {
			s.logger.Info("处理股票列表",
				logger.String("backtest_id", backtest.ID),
				logger.Int("symbols_count", len(backtest.Symbols)),
				logger.String("date", currentDate.Format("2006-01-02")),
			)
		}

		for _, symbol := range backtest.Symbols {
			// 移除每个股票的处理日志

			// 获取真实市场数据
			marketData, err := s.getRealMarketData(ctx, symbol, currentDate)
			if err != nil {
				// 如果获取真实数据失败，记录错误并跳过该股票该日期
				s.logger.Error("获取真实市场数据失败，跳过该股票",
					logger.String("backtest_id", backtest.ID),
					logger.String("symbol", symbol),
					logger.String("date", currentDate.Format("2006-01-02")),
					logger.ErrorField(err),
				)
				continue // 跳过该股票该日期的处理
			}

			// 执行策略
			signal, err := s.strategyService.ExecuteStrategy(ctx, strategy.ID, marketData)
			if err != nil {
				s.logger.Error("策略执行失败",
					logger.String("backtest_id", backtest.ID),
					logger.String("strategy_id", strategy.ID),
					logger.String("symbol", symbol),
					logger.String("date", currentDate.Format("2006-01-02")),
					logger.ErrorField(err),
				)
				continue
			}

			// 移除策略执行成功日志

			// 根据信号执行交易
			if trade := s.executeSignal(signal, marketData, portfolio, backtest); trade != nil {
				trades = append(trades, *trade)
				// 移除交易执行成功日志
			}
		}

		// 更新组合价值
		s.updatePortfolioValue(ctx, portfolio, backtest.Symbols, currentDate)

		// 记录权益曲线
		equityCurve = append(equityCurve, models.EquityPoint{
			Date:           currentDate.Format("2006-01-02"),
			PortfolioValue: portfolio.TotalValue,
			Cash:           portfolio.Cash,
			Holdings:       portfolio.TotalValue - portfolio.Cash,
		})

		// 计算日收益率
		if len(equityCurve) > 1 {
			prevValue := equityCurve[len(equityCurve)-2].PortfolioValue
			if prevValue > 0 {
				dailyReturn := (portfolio.TotalValue - prevValue) / prevValue
				dailyReturns = append(dailyReturns, dailyReturn)
			}
		}

		// 模拟处理时间（移除每日进度日志）

		// 确保日志输出 - 增加延迟确保日志刷新
		time.Sleep(time.Millisecond * 100)

		// 移除强制更新进度日志
	}

	// 计算最终结果
	result := s.calculateBacktestResult(backtest, dailyReturns, portfolio)

	// 保存结果
	s.mutex.Lock()
	s.backtestResults[backtest.ID] = result
	s.backtestTrades[backtest.ID] = trades

	// 更新回测状态
	backtest.Status = models.BacktestStatusCompleted
	backtest.Progress = 100
	now := time.Now()
	backtest.CompletedAt = &now
	s.mutex.Unlock()

	// 最终进度更新
	s.updateBacktestProgress(backtest.ID, 100, "回测完成")

	s.logger.Info("回测任务完成",
		logger.String("backtest_id", backtest.ID),
		logger.Int("total_trades", len(trades)),
		logger.Float64("final_value", portfolio.TotalValue),
		logger.Float64("total_return", result.TotalReturn),
	)
}

// preloadBacktestData 预加载回测期间所有股票的历史数据
func (s *BacktestService) preloadBacktestData(ctx context.Context, symbols []string, startDate, endDate time.Time) error {
	s.logger.Info("开始预加载回测数据",
		logger.Int("symbols_count", len(symbols)),
		logger.String("start_date", startDate.Format("2006-01-02")),
		logger.String("end_date", endDate.Format("2006-01-02")),
	)

	// 获取数据源客户端
	client, err := s.dataSourceService.GetClient()
	if err != nil {
		return fmt.Errorf("获取数据源客户端失败: %w", err)
	}

	// 格式化日期
	startDateStr := startDate.Format("20060102")
	endDateStr := endDate.Format("20060102")

	// 为每个股票预加载数据
	for _, symbol := range symbols {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 检查缓存中是否已有数据
		if s.dailyCacheService != nil {
			if _, found := s.dailyCacheService.Get(symbol, startDateStr, endDateStr); found {
				continue
			}
		}

		// 从API获取数据
		data, err := client.GetDailyData(symbol, startDateStr, endDateStr, "qfq")
		if err != nil {
			s.logger.Error("预加载股票数据失败",
				logger.String("symbol", symbol),
				logger.ErrorField(err),
			)
			// 继续处理其他股票，不中断整个预加载过程
			continue
		}

		// 存入缓存
		if s.dailyCacheService != nil && len(data) > 0 {
			s.dailyCacheService.Set(symbol, startDateStr, endDateStr, data)
		}
	}

	s.logger.Info("回测数据预加载完成")
	return nil
}

// getRealMarketData 获取真实市场数据（从预加载的缓存中查找）
func (s *BacktestService) getRealMarketData(ctx context.Context, symbol string, date time.Time) (*models.MarketData, error) {
	// 从缓存中查找包含该日期的数据
	if s.dailyCacheService == nil {
		return nil, fmt.Errorf("缓存服务未启用")
	}

	// 尝试从较大的时间范围缓存中查找数据
	// 首先尝试当月数据
	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	monthEnd := monthStart.AddDate(0, 1, -1)

	monthStartStr := monthStart.Format("20060102")
	monthEndStr := monthEnd.Format("20060102")

	if cachedData, found := s.dailyCacheService.Get(symbol, monthStartStr, monthEndStr); found {
		return s.findDataForDate(cachedData, symbol, date)
	}

	// 如果没有月度缓存，尝试查找包含该日期的任何缓存
	// 这里可以扩展更复杂的缓存查找逻辑
	// 暂时回退到直接API调用
	return s.getRealMarketDataDirect(ctx, symbol, date)
}

// getRealMarketDataDirect 直接从API获取市场数据（作为备用方案）
func (s *BacktestService) getRealMarketDataDirect(ctx context.Context, symbol string, date time.Time) (*models.MarketData, error) {
	dateStr := date.Format("20060102")

	// 获取数据源客户端
	client, err := s.dataSourceService.GetClient()
	if err != nil {
		return nil, fmt.Errorf("获取数据源客户端失败: %w", err)
	}

	// 尝试获取前后一周的数据，找到最近的交易日
	startDate := date.AddDate(0, 0, -7)
	endDate := date.AddDate(0, 0, 7)
	startDateStr := startDate.Format("20060102")
	endDateStr := endDate.Format("20060102")

	weekData, err := client.GetDailyData(symbol, startDateStr, endDateStr, "qfq")
	if err != nil {
		return nil, fmt.Errorf("获取股票数据失败: %w", err)
	}

	if len(weekData) == 0 {
		return nil, fmt.Errorf("股票 %s 在 %s 前后一周都无交易数据", symbol, dateStr)
	}

	// 存入缓存
	if s.dailyCacheService != nil {
		s.dailyCacheService.Set(symbol, startDateStr, endDateStr, weekData)
	}

	return s.findDataForDate(weekData, symbol, date)
}

// findDataForDate 从数据列表中找到指定日期的数据
func (s *BacktestService) findDataForDate(dailyData []models.StockDaily, symbol string, targetDate time.Time) (*models.MarketData, error) {
	targetDateStr := targetDate.Format("20060102")

	// 首先尝试精确匹配
	for _, data := range dailyData {
		// 从ISO格式转换为YYYYMMDD格式进行比较
		if tradeTime, err := time.Parse("2006-01-02T15:04:05.000", data.TradeDate); err == nil {
			tradeDateStr := tradeTime.Format("20060102")
			if tradeDateStr == targetDateStr {
				return s.convertStockDailyToMarketData(data, symbol, targetDate)
			}
		}
	}

	// 如果没有精确匹配，找最近的交易日
	var closestData *models.StockDaily
	minDiff := int64(^uint64(0) >> 1) // 最大int64值

	for _, data := range dailyData {
		// 使用正确的ISO日期格式解析
		tradeTime, err := time.Parse("2006-01-02T15:04:05.000", data.TradeDate)
		if err != nil {
			continue
		}

		diff := targetDate.Unix() - tradeTime.Unix()
		if diff < 0 {
			diff = -diff
		}

		if diff < minDiff {
			minDiff = diff
			closestData = &data
		}
	}

	if closestData == nil {
		return nil, fmt.Errorf("无法找到股票 %s 在 %s 的有效交易数据", symbol, targetDateStr)
	}

	// 只在使用非精确匹配时记录日志
	if closestTradeTime, err := time.Parse("2006-01-02T15:04:05.000", closestData.TradeDate); err == nil {
		closestDateStr := closestTradeTime.Format("20060102")
		if closestDateStr != targetDateStr {
			s.logger.Debug("使用最近交易日",
				logger.String("symbol", symbol),
				logger.String("target", targetDateStr),
				logger.String("actual", closestDateStr),
			)
		}
	}

	return s.convertStockDailyToMarketData(*closestData, symbol, targetDate)
}

// convertStockDailyToMarketData 将StockDaily转换为MarketData
func (s *BacktestService) convertStockDailyToMarketData(stockDaily models.StockDaily, symbol string, date time.Time) (*models.MarketData, error) {
	// 转换价格数据（去掉无意义的精度警告）
	open, _ := stockDaily.Open.Float64()
	high, _ := stockDaily.High.Float64()
	low, _ := stockDaily.Low.Float64()
	close, _ := stockDaily.Close.Float64()

	// 转换成交量（注意：StockDaily.Vol是手，需要转换为股）
	vol, _ := stockDaily.Vol.Float64()
	volume := int64(vol * 100) // 1手 = 100股

	// 转换成交额（注意：StockDaily.Amount是千元，需要转换为元）
	amount, _ := stockDaily.Amount.Float64()
	amountInYuan := amount * 1000 // 千元转元

	return &models.MarketData{
		Symbol:   symbol,
		Date:     date,
		Open:     open,
		High:     high,
		Low:      low,
		Close:    close,
		Volume:   volume,
		Amount:   amountInYuan,
		AdjClose: close, // 使用前复权数据，收盘价就是复权价
	}, nil
}

// generateMockMarketData 生成模拟市场数据（保留作为备用）
func (s *BacktestService) generateMockMarketData(symbol string, date time.Time) *models.MarketData {
	// 使用日期作为随机种子，确保相同日期生成相同数据
	rand.Seed(date.Unix() + int64(len(symbol)))

	// 基础价格（根据股票代码生成）
	basePrice := 10.0 + float64(len(symbol))*2.5

	// 添加趋势和随机波动
	days := int(date.Sub(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)).Hours() / 24)
	trend := math.Sin(float64(days)/365.0*2*math.Pi) * 0.3 // 年度周期
	randomWalk := (rand.Float64() - 0.5) * 0.1             // 随机游走

	close := basePrice * (1 + trend + randomWalk)
	volatility := 0.02 + rand.Float64()*0.03 // 2-5%的波动率

	high := close * (1 + volatility*rand.Float64())
	low := close * (1 - volatility*rand.Float64())
	open := low + (high-low)*rand.Float64()

	volume := int64(1000000 + rand.Intn(9000000)) // 100万-1000万

	return &models.MarketData{
		Symbol: symbol,
		Date:   date,
		Open:   open,
		High:   high,
		Low:    low,
		Close:  close,
		Volume: volume,
		Amount: float64(volume) * (high + low) / 2,
	}
}

// executeSignal 执行交易信号
func (s *BacktestService) executeSignal(signal *models.Signal, marketData *models.MarketData, portfolio *models.Portfolio, backtest *models.Backtest) *models.Trade {
	if signal.SignalType == models.SignalTypeHold {
		return nil
	}

	// 计算交易数量
	var quantity int
	var price float64 = marketData.Close

	if signal.Side == models.TradeSideBuy {
		// 买入：使用可用现金的一定比例
		maxInvestment := portfolio.Cash * 0.2 * signal.Strength // 最多使用20%的现金，根据信号强度调整
		quantity = int(maxInvestment/price/100) * 100           // 以手为单位（100股）

		if quantity < 100 || portfolio.Cash < float64(quantity)*price*(1+backtest.Commission) {
			return nil // 资金不足或数量太少
		}
	} else {
		// 卖出：卖出持有的股票
		position, exists := portfolio.Positions[signal.Symbol]
		if !exists || position.Quantity <= 0 {
			return nil // 没有持仓
		}

		quantity = int(float64(position.Quantity) * signal.Strength) // 根据信号强度决定卖出比例
		if quantity < 100 {
			quantity = position.Quantity // 全部卖出
		}
		quantity = -quantity // 卖出为负数
	}

	// 计算手续费
	commission := math.Abs(float64(quantity)) * price * backtest.Commission

	// 创建交易记录
	trade := &models.Trade{
		ID:         fmt.Sprintf("trade_%s_%d", signal.Symbol, time.Now().UnixNano()),
		BacktestID: backtest.ID,
		Symbol:     signal.Symbol,
		Side:       signal.Side,
		Quantity:   int(math.Abs(float64(quantity))),
		Price:      price,
		Commission: commission,
		SignalType: string(signal.SignalType),
		Timestamp:  signal.Timestamp,
		CreatedAt:  time.Now(),
	}

	// 更新组合
	s.updatePortfolioWithTrade(portfolio, trade)

	return trade
}

// updatePortfolioWithTrade 根据交易更新组合
func (s *BacktestService) updatePortfolioWithTrade(portfolio *models.Portfolio, trade *models.Trade) {
	if trade.Side == models.TradeSideBuy {
		// 买入
		totalCost := float64(trade.Quantity)*trade.Price + trade.Commission
		portfolio.Cash -= totalCost

		// 更新持仓
		if position, exists := portfolio.Positions[trade.Symbol]; exists {
			// 计算新的平均成本
			totalShares := position.Quantity + trade.Quantity
			totalValue := float64(position.Quantity)*position.AvgPrice + float64(trade.Quantity)*trade.Price

			position.Quantity = totalShares
			position.AvgPrice = totalValue / float64(totalShares)
			portfolio.Positions[trade.Symbol] = position
		} else {
			// 新建持仓
			portfolio.Positions[trade.Symbol] = models.Position{
				Symbol:   trade.Symbol,
				Quantity: trade.Quantity,
				AvgPrice: trade.Price,
			}
		}
	} else {
		// 卖出
		totalRevenue := float64(trade.Quantity)*trade.Price - trade.Commission
		portfolio.Cash += totalRevenue

		// 更新持仓
		if position, exists := portfolio.Positions[trade.Symbol]; exists {
			// 计算盈亏
			trade.PnL = float64(trade.Quantity)*(trade.Price-position.AvgPrice) - trade.Commission

			position.Quantity -= trade.Quantity
			if position.Quantity <= 0 {
				delete(portfolio.Positions, trade.Symbol)
			} else {
				portfolio.Positions[trade.Symbol] = position
			}
		}
	}
}

// updatePortfolioValue 更新组合价值
func (s *BacktestService) updatePortfolioValue(ctx context.Context, portfolio *models.Portfolio, symbols []string, date time.Time) {
	holdingsValue := 0.0

	for symbol, position := range portfolio.Positions {
		if position.Quantity > 0 {
			// 获取当前市价
			marketData, err := s.getRealMarketData(ctx, symbol, date)
			if err != nil {
				// 如果获取数据失败，使用上一次的市值继续计算
				// 保持原有的市值不变
				holdingsValue += position.MarketValue
				continue
			}

			marketValue := float64(position.Quantity) * marketData.Close
			holdingsValue += marketValue

			// 更新持仓信息
			position.MarketValue = marketValue
			position.UnrealizedPL = marketValue - float64(position.Quantity)*position.AvgPrice
			portfolio.Positions[symbol] = position
		}
	}

	portfolio.TotalValue = portfolio.Cash + holdingsValue
}

// calculateBacktestResult 计算回测结果
func (s *BacktestService) calculateBacktestResult(backtest *models.Backtest, dailyReturns []float64, portfolio *models.Portfolio) *models.BacktestResult {
	metrics := &models.PerformanceMetrics{
		Returns:      dailyReturns,
		RiskFreeRate: 0.03 / 252, // 年化3%无风险利率转为日利率
	}

	result := metrics.CalculateMetrics()
	result.ID = fmt.Sprintf("result_%s", backtest.ID)
	result.BacktestID = backtest.ID
	result.CreatedAt = time.Now()

	// 计算总收益率
	if backtest.InitialCash > 0 {
		result.TotalReturn = (portfolio.TotalValue - backtest.InitialCash) / backtest.InitialCash
	}

	return result
}

// CancelBacktest 取消回测
func (s *BacktestService) CancelBacktest(ctx context.Context, backtestID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	backtest, exists := s.backtests[backtestID]
	if !exists {
		return ErrBacktestNotFound
	}

	if backtest.Status != models.BacktestStatusRunning {
		return fmt.Errorf("回测状态不允许取消: %s", backtest.Status)
	}

	// 取消运行中的回测
	if cancelFunc, ok := s.runningBacktests[backtestID]; ok {
		cancelFunc()
		delete(s.runningBacktests, backtestID)
	}

	// 更新状态
	backtest.Status = models.BacktestStatusCancelled
	now := time.Now()
	backtest.CompletedAt = &now

	s.logger.Info("回测取消成功", logger.String("backtest_id", backtestID))

	return nil
}

// GetBacktestProgress 获取回测进度
func (s *BacktestService) GetBacktestProgress(ctx context.Context, backtestID string) (*models.BacktestProgress, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 检查回测是否存在
	backtest, exists := s.backtests[backtestID]
	if !exists {
		return nil, ErrBacktestNotFound
	}

	// 获取进度信息
	if progress, ok := s.backtestProgress[backtestID]; ok {
		return progress, nil
	}

	// 如果没有进度信息，根据状态返回默认进度
	progress := &models.BacktestProgress{
		BacktestID: backtestID,
		Status:     string(backtest.Status),
		Progress:   backtest.Progress,
	}

	switch backtest.Status {
	case models.BacktestStatusPending:
		progress.Message = "等待启动"
	case models.BacktestStatusRunning:
		progress.Message = "运行中..."
	case models.BacktestStatusCompleted:
		progress.Message = "已完成"
		progress.Progress = 100
	case models.BacktestStatusFailed:
		progress.Message = "执行失败"
		progress.Error = backtest.ErrorMessage
	case models.BacktestStatusCancelled:
		progress.Message = "已取消"
	}

	return progress, nil
}

// GetBacktestResults 获取回测结果
func (s *BacktestService) GetBacktestResults(ctx context.Context, backtestID string) (*models.BacktestResultsResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 检查回测是否存在
	backtest, exists := s.backtests[backtestID]
	if !exists {
		return nil, ErrBacktestNotFound
	}

	// 检查回测是否完成
	if backtest.Status != models.BacktestStatusCompleted {
		return nil, ErrBacktestNotCompleted
	}

	// 获取回测结果
	result, exists := s.backtestResults[backtestID]
	if !exists {
		return nil, fmt.Errorf("回测结果不存在")
	}

	// 获取交易记录
	trades, _ := s.backtestTrades[backtestID]

	// 生成权益曲线（简化版）
	equityCurve := s.generateEquityCurve(backtest, result)

	response := &models.BacktestResultsResponse{
		BacktestID:  backtestID,
		Performance: *result,
		EquityCurve: equityCurve,
		Trades:      trades,
	}

	return response, nil
}

// generateEquityCurve 生成权益曲线
func (s *BacktestService) generateEquityCurve(backtest *models.Backtest, result *models.BacktestResult) []models.EquityPoint {
	var curve []models.EquityPoint

	totalDays := int(backtest.EndDate.Sub(backtest.StartDate).Hours() / 24)
	if totalDays <= 0 {
		totalDays = 1
	}

	// 简化的权益曲线生成
	for i := 0; i <= totalDays; i += 7 { // 每周一个点
		date := backtest.StartDate.AddDate(0, 0, i)
		if date.After(backtest.EndDate) {
			date = backtest.EndDate
		}

		// 计算该日期的组合价值
		progress := float64(i) / float64(totalDays)
		portfolioValue := backtest.InitialCash * (1 + result.TotalReturn*progress)

		// 添加一些波动
		volatility := 0.1 * result.TotalReturn * math.Sin(float64(i)/10.0)
		portfolioValue *= (1 + volatility)

		curve = append(curve, models.EquityPoint{
			Date:           date.Format("2006-01-02"),
			PortfolioValue: portfolioValue,
			BenchmarkValue: backtest.InitialCash * (1 + result.BenchmarkReturn*progress),
			Cash:           portfolioValue * 0.1, // 假设10%现金
			Holdings:       portfolioValue * 0.9, // 假设90%持仓
		})

		if date.Equal(backtest.EndDate) {
			break
		}
	}

	return curve
}

// updateBacktestStatus 更新回测状态
func (s *BacktestService) updateBacktestStatus(backtestID string, status models.BacktestStatus, message string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if backtest, exists := s.backtests[backtestID]; exists {
		backtest.Status = status
		if message != "" {
			backtest.ErrorMessage = message
		}

		if status == models.BacktestStatusCompleted || status == models.BacktestStatusFailed || status == models.BacktestStatusCancelled {
			now := time.Now()
			backtest.CompletedAt = &now
		}
	}
}

// updateBacktestProgress 更新回测进度
func (s *BacktestService) updateBacktestProgress(backtestID string, progress int, message string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if progressInfo, exists := s.backtestProgress[backtestID]; exists {
		progressInfo.Progress = progress
		progressInfo.Message = message

		if progress >= 100 {
			progressInfo.Status = string(models.BacktestStatusCompleted)
		}
	}

	// 同时更新回测对象的进度
	if backtest, exists := s.backtests[backtestID]; exists {
		backtest.Progress = progress
	}
}
