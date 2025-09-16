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
	backtests            map[string]*models.Backtest
	backtestResults      map[string]*models.BacktestResult  // 单策略结果（兼容性）
	backtestMultiResults map[string][]models.BacktestResult // 多策略结果
	backtestEquityCurves map[string][]models.EquityPoint    // 组合权益曲线
	backtestTrades       map[string][]models.Trade
	backtestProgress     map[string]*models.BacktestProgress
	runningBacktests     map[string]context.CancelFunc // 用于取消运行中的回测

	strategyService   *StrategyService
	dataSourceService *DataSourceService
	dailyCacheService *DailyCacheService // 使用现有的日线数据缓存服务
	logger            logger.Logger
	mutex             sync.RWMutex
}

// NewBacktestService 创建回测服务
func NewBacktestService(strategyService *StrategyService, dataSourceService *DataSourceService, dailyCacheService *DailyCacheService, log logger.Logger) *BacktestService {
	return &BacktestService{
		backtests:            make(map[string]*models.Backtest),
		backtestResults:      make(map[string]*models.BacktestResult),
		backtestMultiResults: make(map[string][]models.BacktestResult),
		backtestEquityCurves: make(map[string][]models.EquityPoint),
		backtestTrades:       make(map[string][]models.Trade),
		backtestProgress:     make(map[string]*models.BacktestProgress),
		runningBacktests:     make(map[string]context.CancelFunc),
		strategyService:      strategyService,
		dataSourceService:    dataSourceService,
		dailyCacheService:    dailyCacheService,
		logger:               log,
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

	// 自动处理重复名称，生成唯一名称
	backtest.Name = s.generateUniqueBacktestName(backtest.Name)

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
		// 生成唯一名称（排除当前回测）
		uniqueName := s.generateUniqueBacktestNameForUpdate(*req.Name, backtestID)
		backtest.Name = uniqueName
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

// generateUniqueBacktestName 生成唯一的回测名称
// 如果名称已存在，会自动在后面添加序号，如 "原名称 (2)"
func (s *BacktestService) generateUniqueBacktestName(originalName string) string {
	// 检查原始名称是否已存在
	if !s.isBacktestNameExists(originalName) {
		return originalName
	}

	// 如果存在，尝试添加序号
	counter := 2
	for {
		newName := fmt.Sprintf("%s (%d)", originalName, counter)
		if !s.isBacktestNameExists(newName) {
			return newName
		}
		counter++

		// 防止无限循环，最多尝试1000次
		if counter > 1000 {
			// 如果还是重复，使用时间戳确保唯一性
			return fmt.Sprintf("%s (%d)", originalName, time.Now().Unix())
		}
	}
}

// isBacktestNameExists 检查回测名称是否已存在
func (s *BacktestService) isBacktestNameExists(name string) bool {
	for _, existing := range s.backtests {
		if existing.Name == name {
			return true
		}
	}
	return false
}

// generateUniqueBacktestNameForUpdate 为更新操作生成唯一的回测名称
// 排除指定的回测ID，避免与自身冲突
func (s *BacktestService) generateUniqueBacktestNameForUpdate(originalName string, excludeID string) string {
	// 检查原始名称是否已存在（排除当前回测）
	if !s.isBacktestNameExistsExcluding(originalName, excludeID) {
		return originalName
	}

	// 如果存在，尝试添加序号
	counter := 2
	for {
		newName := fmt.Sprintf("%s (%d)", originalName, counter)
		if !s.isBacktestNameExistsExcluding(newName, excludeID) {
			return newName
		}
		counter++

		// 防止无限循环，最多尝试1000次
		if counter > 1000 {
			// 如果还是重复，使用时间戳确保唯一性
			return fmt.Sprintf("%s (%d)", originalName, time.Now().Unix())
		}
	}
}

// isBacktestNameExistsExcluding 检查回测名称是否已存在（排除指定ID）
func (s *BacktestService) isBacktestNameExistsExcluding(name string, excludeID string) bool {
	for id, existing := range s.backtests {
		if id != excludeID && existing.Name == name {
			return true
		}
	}
	return false
}

// StartBacktest 启动回测（支持多策略）
func (s *BacktestService) StartBacktest(ctx context.Context, backtest *models.Backtest, strategies []*models.Strategy) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 检查回测状态
	if backtest.Status != models.BacktestStatusPending {
		return fmt.Errorf("回测状态不允许启动: %s", backtest.Status)
	}

	// 验证策略数量
	if len(strategies) == 0 {
		return fmt.Errorf("至少需要一个策略")
	}

	if len(strategies) != len(backtest.StrategyIDs) {
		return fmt.Errorf("策略数量与配置不匹配")
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
		Message:    fmt.Sprintf("初始化多策略回测环境... (%d个策略)", len(strategies)),
	}

	// 创建可取消的上下文，根据回测时间范围动态设置超时
	totalDays := int(backtest.EndDate.Sub(backtest.StartDate).Hours() / 24)
	// 多策略需要更多时间，超时时间适当增加
	timeoutMinutes := maxInt(10, minInt(240, totalDays/3*len(strategies)))
	timeout := time.Duration(timeoutMinutes) * time.Minute

	s.logger.Info("设置多策略回测超时时间",
		logger.String("backtest_id", backtest.ID),
		logger.Int("strategies_count", len(strategies)),
		logger.Int("total_days", totalDays),
		logger.Int("timeout_minutes", timeoutMinutes),
	)

	backtestCtx, cancel := context.WithTimeout(ctx, timeout)
	s.runningBacktests[backtest.ID] = cancel

	s.logger.Info("准备启动多策略回测goroutine",
		logger.String("backtest_id", backtest.ID),
		logger.Any("strategy_ids", backtest.StrategyIDs),
		logger.Int("total_days", int(backtest.EndDate.Sub(backtest.StartDate).Hours()/24)),
		logger.String("start_date", backtest.StartDate.Format("2006-01-02")),
		logger.String("end_date", backtest.EndDate.Format("2006-01-02")),
		logger.Int("symbols_count", len(backtest.Symbols)),
	)

	// 启动后台多策略回测任务
	go s.runMultiStrategyBacktestTask(backtestCtx, backtest, strategies)

	s.logger.Info("多策略回测启动成功",
		logger.String("backtest_id", backtest.ID),
		logger.Any("strategy_ids", backtest.StrategyIDs),
		logger.String("status", "goroutine_launched"),
	)

	return nil
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

		// 存入缓存 - 使用多个时间范围键提高命中率
		if s.dailyCacheService != nil && len(data) > 0 {
			// 1. 存储原始回测期间范围
			s.dailyCacheService.Set(symbol, startDateStr, endDateStr, data)

			// 2. 按月份分别存储，提高月度查找命中率
			s.storeDataByMonths(symbol, data, startDate, endDate)

			s.logger.Info("预加载数据已缓存",
				logger.String("symbol", symbol),
				logger.String("range", fmt.Sprintf("%s至%s", startDateStr, endDateStr)),
				logger.Int("data_count", len(data)),
			)
		}
	}

	s.logger.Info("回测数据预加载完成")
	return nil
}

// storeDataByMonths 按月份存储数据到缓存，提高查找命中率
func (s *BacktestService) storeDataByMonths(symbol string, data []models.StockDaily, startDate, endDate time.Time) {
	// 按月份组织数据
	monthlyData := make(map[string][]models.StockDaily)

	for _, daily := range data {
		// 解析交易日期
		tradeTime, err := time.Parse("2006-01-02T15:04:05.000", daily.TradeDate)
		if err != nil {
			continue
		}

		// 生成月份键 (YYYYMM)
		monthKey := tradeTime.Format("200601")
		if monthlyData[monthKey] == nil {
			monthlyData[monthKey] = make([]models.StockDaily, 0)
		}
		monthlyData[monthKey] = append(monthlyData[monthKey], daily)
	}

	// 为每个月份存储缓存
	for monthKey, monthData := range monthlyData {
		if len(monthData) == 0 {
			continue
		}

		// 解析月份，生成完整的月份范围
		year, _ := time.Parse("200601", monthKey)
		monthStart := time.Date(year.Year(), year.Month(), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, -1)

		monthStartStr := monthStart.Format("20060102")
		monthEndStr := monthEnd.Format("20060102")

		// 检查是否已经缓存
		if _, found := s.dailyCacheService.Get(symbol, monthStartStr, monthEndStr); !found {
			s.dailyCacheService.Set(symbol, monthStartStr, monthEndStr, monthData)
			s.logger.Debug("月度数据已缓存",
				logger.String("symbol", symbol),
				logger.String("month", monthKey),
				logger.Int("data_count", len(monthData)),
			)
		}
	}
}

// getRealMarketData 获取真实市场数据（从预加载的缓存中查找）
func (s *BacktestService) getRealMarketData(ctx context.Context, symbol string, date time.Time) (*models.MarketData, error) {
	// 从缓存中查找包含该日期的数据
	if s.dailyCacheService == nil {
		return nil, fmt.Errorf("缓存服务未启用")
	}

	// 优化缓存查找逻辑：尝试从不同时间范围的缓存中查找

	// 1. 首先尝试从预加载的回测期间缓存中查找
	// 这是最可能命中的情况，因为preloadBacktestData已经加载了完整回测期间的数据
	if cachedData := s.findCachedDataContainingDate(symbol, date); cachedData != nil {
		return s.findDataForDate(cachedData, symbol, date)
	}

	// 2. 如果没有找到预加载缓存，尝试月度缓存
	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	monthEnd := monthStart.AddDate(0, 1, -1)
	monthStartStr := monthStart.Format("20060102")
	monthEndStr := monthEnd.Format("20060102")

	if cachedData, found := s.dailyCacheService.Get(symbol, monthStartStr, monthEndStr); found {
		return s.findDataForDate(cachedData, symbol, date)
	}

	// 3. 最后才回退到直接API调用，但优化为获取更大时间段以减少API调用
	return s.getRealMarketDataDirectOptimized(ctx, symbol, date)
}

// findCachedDataContainingDate 查找包含指定日期的缓存数据
func (s *BacktestService) findCachedDataContainingDate(symbol string, targetDate time.Time) []models.StockDaily {
	// 遍历所有缓存条目，查找包含目标日期的时间范围

	// 由于cache是sync.Map，我们需要通过反射或其他方式遍历
	// 这里采用一种更直接的方式：尝试常见的时间范围

	// 尝试不同的时间范围组合，从大到小
	timeRanges := []struct {
		start, end time.Time
	}{
		// 尝试季度范围
		{
			start: time.Date(targetDate.Year(), ((targetDate.Month()-1)/3)*3+1, 1, 0, 0, 0, 0, targetDate.Location()),
			end:   time.Date(targetDate.Year(), ((targetDate.Month()-1)/3+1)*3, 31, 0, 0, 0, 0, targetDate.Location()),
		},
		// 尝试半年范围
		{
			start: time.Date(targetDate.Year(), 1, 1, 0, 0, 0, 0, targetDate.Location()),
			end:   time.Date(targetDate.Year(), 6, 30, 0, 0, 0, 0, targetDate.Location()),
		},
		{
			start: time.Date(targetDate.Year(), 7, 1, 0, 0, 0, 0, targetDate.Location()),
			end:   time.Date(targetDate.Year(), 12, 31, 0, 0, 0, 0, targetDate.Location()),
		},
		// 尝试全年范围
		{
			start: time.Date(targetDate.Year(), 1, 1, 0, 0, 0, 0, targetDate.Location()),
			end:   time.Date(targetDate.Year(), 12, 31, 0, 0, 0, 0, targetDate.Location()),
		},
		// 尝试跨年范围（前一年下半年到当年上半年）
		{
			start: time.Date(targetDate.Year()-1, 7, 1, 0, 0, 0, 0, targetDate.Location()),
			end:   time.Date(targetDate.Year(), 6, 30, 0, 0, 0, 0, targetDate.Location()),
		},
	}

	for _, tr := range timeRanges {
		startStr := tr.start.Format("20060102")
		endStr := tr.end.Format("20060102")

		if cachedData, found := s.dailyCacheService.Get(symbol, startStr, endStr); found {
			// 验证缓存数据确实包含目标日期
			if s.dataContainsDate(cachedData, targetDate) {
				return cachedData
			}
		}
	}

	return nil
}

// dataContainsDate 检查数据是否包含指定日期
func (s *BacktestService) dataContainsDate(data []models.StockDaily, targetDate time.Time) bool {
	targetDateStr := targetDate.Format("20060102")

	for _, daily := range data {
		// 从ISO格式转换为YYYYMMDD格式进行比较
		if tradeTime, err := time.Parse("2006-01-02T15:04:05.000", daily.TradeDate); err == nil {
			tradeDateStr := tradeTime.Format("20060102")
			if tradeDateStr == targetDateStr {
				return true
			}
		}
	}
	return false
}

// getRealMarketDataDirectOptimized 优化的直接API获取（获取更大时间段减少API调用）
func (s *BacktestService) getRealMarketDataDirectOptimized(ctx context.Context, symbol string, date time.Time) (*models.MarketData, error) {
	// 获取数据源客户端
	client, err := s.dataSourceService.GetClient()
	if err != nil {
		return nil, fmt.Errorf("获取数据源客户端失败: %w", err)
	}

	// 优化：获取更大的时间范围（一个月）而不是一周，减少API调用次数
	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	monthEnd := monthStart.AddDate(0, 1, -1)

	startDateStr := monthStart.Format("20060102")
	endDateStr := monthEnd.Format("20060102")

	s.logger.Info("直接API获取月度数据",
		logger.String("symbol", symbol),
		logger.String("target_date", date.Format("2006-01-02")),
		logger.String("fetch_range", fmt.Sprintf("%s至%s", startDateStr, endDateStr)),
	)

	monthData, err := client.GetDailyData(symbol, startDateStr, endDateStr, "qfq")
	if err != nil {
		// 如果月度获取失败，回退到原来的周度获取
		s.logger.Warn("月度数据获取失败，回退到周度获取",
			logger.String("symbol", symbol),
			logger.ErrorField(err),
		)
		return s.getRealMarketDataDirect(ctx, symbol, date)
	}

	if len(monthData) == 0 {
		return nil, fmt.Errorf("股票 %s 在 %s 月份无交易数据", symbol, date.Format("2006-01"))
	}

	// 存入缓存（使用月度范围）
	if s.dailyCacheService != nil {
		s.dailyCacheService.Set(symbol, startDateStr, endDateStr, monthData)
		s.logger.Info("月度数据已缓存",
			logger.String("symbol", symbol),
			logger.String("range", fmt.Sprintf("%s至%s", startDateStr, endDateStr)),
			logger.Int("data_count", len(monthData)),
		)
	}

	return s.findDataForDate(monthData, symbol, date)
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
	trades := s.backtestTrades[backtestID]

	// 获取策略信息 - 修复兼容性处理逻辑
	var strategy *models.Strategy

	// 优先使用多策略ID，兼容旧的单策略ID
	var strategyID string
	if len(backtest.StrategyIDs) > 0 {
		// 多策略情况下，使用第一个策略作为主策略（兼容旧接口）
		strategyID = backtest.StrategyIDs[0]
	} else if backtest.StrategyID != "" {
		// 兼容旧的单策略
		strategyID = backtest.StrategyID
	}

	if strategyID != "" {
		var err error
		strategy, err = s.strategyService.GetStrategy(ctx, strategyID)
		if err != nil {
			s.logger.Warn("获取策略信息失败",
				logger.String("strategy_id", strategyID),
				logger.ErrorField(err),
			)
			// 即使获取策略失败也不应该影响回测结果返回
			strategy = nil
		}
	}

	// 生成权益曲线（简化版）
	equityCurve := s.generateEquityCurve(backtest, result)

	// 构建回测配置信息
	backtestConfig := models.BacktestConfig{
		Name:        backtest.Name,
		StartDate:   backtest.StartDate.Format("2006-01-02"),
		EndDate:     backtest.EndDate.Format("2006-01-02"),
		InitialCash: backtest.InitialCash,
		Symbols:     backtest.Symbols,
		Commission:  backtest.Commission,
		CreatedAt:   backtest.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 检查是否有多策略结果
	multiResults, hasMultiResults := s.backtestMultiResults[backtestID]
	combinedEquityCurve, hasEquityCurve := s.backtestEquityCurves[backtestID]

	var strategies []*models.Strategy
	var performanceResults []models.BacktestResult
	var finalEquityCurve []models.EquityPoint

	if hasMultiResults && len(multiResults) > 0 {
		// 多策略结果
		performanceResults = multiResults

		// 获取所有策略信息
		strategyIDs := backtest.StrategyIDs
		if len(strategyIDs) == 0 && backtest.StrategyID != "" {
			// 兼容性处理
			strategyIDs = []string{backtest.StrategyID}
		}

		for _, strategyID := range strategyIDs {
			// 跳过空的策略ID
			if strings.TrimSpace(strategyID) == "" {
				s.logger.Warn("跳过空的策略ID",
					logger.String("backtest_id", backtestID),
				)
				continue
			}

			strategy, err := s.strategyService.GetStrategy(ctx, strategyID)
			if err != nil {
				s.logger.Warn("获取策略信息失败",
					logger.String("strategy_id", strategyID),
					logger.ErrorField(err),
				)
				// 创建默认策略信息
				strategy = &models.Strategy{
					ID:   strategyID,
					Name: fmt.Sprintf("策略-%s", strategyID),
				}
			}
			strategies = append(strategies, strategy)
		}

		// 使用组合权益曲线
		if hasEquityCurve && len(combinedEquityCurve) > 0 {
			finalEquityCurve = combinedEquityCurve
		} else {
			finalEquityCurve = equityCurve
		}
	} else {
		// 单策略结果（兼容性）
		performanceResults = []models.BacktestResult{*result}
		strategies = []*models.Strategy{strategy}
		finalEquityCurve = equityCurve
	}

	// 计算组合整体指标（如果是多策略）
	var combinedMetrics *models.BacktestResult
	if len(performanceResults) > 1 {
		combinedMetrics = s.calculateCombinedMetrics(performanceResults)
	}

	response := &models.BacktestResultsResponse{
		BacktestID:      backtestID,
		Performance:     performanceResults,
		EquityCurve:     finalEquityCurve,
		Trades:          trades,
		Strategies:      strategies,
		BacktestConfig:  backtestConfig,
		CombinedMetrics: combinedMetrics,
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

// runMultiStrategyBacktestTask 运行多策略回测任务
func (s *BacktestService) runMultiStrategyBacktestTask(ctx context.Context, backtest *models.Backtest, strategies []*models.Strategy) {
	// 立即输出日志，确保goroutine已启动
	s.logger.Info("🚀 多策略回测goroutine已启动",
		logger.String("backtest_id", backtest.ID),
		logger.String("goroutine_status", "started"),
		logger.Int("strategies_count", len(strategies)),
	)

	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("多策略回测任务出现panic",
				logger.String("backtest_id", backtest.ID),
				logger.Any("panic", r),
			)
		}
		s.mutex.Lock()
		delete(s.runningBacktests, backtest.ID)
		s.mutex.Unlock()
		s.logger.Info("多策略回测任务清理完成", logger.String("backtest_id", backtest.ID))
	}()

	// 预加载回测数据
	s.logger.Info("开始预加载多策略回测数据",
		logger.String("backtest_id", backtest.ID),
		logger.Int("symbols_count", len(backtest.Symbols)),
		logger.Int("strategies_count", len(strategies)),
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

	s.logger.Info("多策略回测参数计算完成",
		logger.String("backtest_id", backtest.ID),
		logger.Int("total_days", totalDays),
		logger.Int("strategies_count", len(strategies)),
	)

	// 为每个策略创建独立的投资组合
	strategyPortfolios := make(map[string]*models.Portfolio)
	strategyTrades := make(map[string][]models.Trade)
	strategyEquityCurves := make(map[string][]models.EquityPoint)
	strategyDailyReturns := make(map[string][]float64)

	// 计算每个策略的初始资金（平均分配）
	initialCashPerStrategy := backtest.InitialCash / float64(len(strategies))

	for _, strategy := range strategies {
		strategyPortfolios[strategy.ID] = &models.Portfolio{
			Cash:       initialCashPerStrategy,
			Positions:  make(map[string]models.Position),
			TotalValue: initialCashPerStrategy,
		}
		strategyTrades[strategy.ID] = []models.Trade{}
		strategyEquityCurves[strategy.ID] = []models.EquityPoint{}
		strategyDailyReturns[strategy.ID] = []float64{}
	}

	// 模拟每日回测
	s.logger.Info("开始多策略回测循环",
		logger.String("backtest_id", backtest.ID),
		logger.Int("total_days", totalDays),
		logger.String("start_date", backtest.StartDate.Format("2006-01-02")),
		logger.String("end_date", backtest.EndDate.Format("2006-01-02")),
		logger.Int("strategies_count", len(strategies)),
	)

	// 用于控制日志输出频率
	var lastLoggedProgress int = -1
	const progressLogInterval = 10 // 每10%打印一次进度

	for day := 0; day <= totalDays; day++ {
		select {
		case <-ctx.Done():
			// 回测被取消或超时
			if ctx.Err() == context.DeadlineExceeded {
				s.logger.Error("多策略回测超时",
					logger.String("backtest_id", backtest.ID),
					logger.ErrorField(ctx.Err()),
				)
				s.updateBacktestStatus(backtest.ID, models.BacktestStatusFailed, "回测执行超时")
			} else {
				s.logger.Info("多策略回测被取消",
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
			s.logger.Info("多策略回测日期超出范围，结束循环",
				logger.String("backtest_id", backtest.ID),
				logger.String("current_date", currentDate.Format("2006-01-02")),
			)
			break
		}

		// 更新进度
		progress := int(float64(day) / float64(totalDays) * 100)

		// 只在进度达到特定节点时打印日志
		if progress >= lastLoggedProgress+progressLogInterval || day == 0 || day == totalDays {
			s.logger.Info("多策略回测进度更新",
				logger.String("backtest_id", backtest.ID),
				logger.Int("progress", progress),
				logger.String("current_date", currentDate.Format("2006-01-02")),
				logger.Int("day", day),
				logger.Int("total_days", totalDays),
			)
			lastLoggedProgress = progress
		}
		s.updateBacktestProgress(backtest.ID, progress, fmt.Sprintf("多策略回测进行中... %s", currentDate.Format("2006-01-02")))

		// 对每个股票执行所有策略
		for _, symbol := range backtest.Symbols {
			// 获取真实市场数据
			marketData, err := s.getRealMarketData(ctx, symbol, currentDate)
			if err != nil {
				s.logger.Error("获取真实市场数据失败，跳过该股票",
					logger.String("backtest_id", backtest.ID),
					logger.String("symbol", symbol),
					logger.String("date", currentDate.Format("2006-01-02")),
					logger.ErrorField(err),
				)
				continue
			}

			// 为每个策略执行交易逻辑
			for _, strategy := range strategies {
				portfolio := strategyPortfolios[strategy.ID]

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

				// 根据信号执行交易
				if trade := s.executeSignalForStrategy(signal, marketData, portfolio, backtest, strategy.ID); trade != nil {
					strategyTrades[strategy.ID] = append(strategyTrades[strategy.ID], *trade)
				}
			}
		}

		// 更新每个策略的组合价值
		for _, strategy := range strategies {
			portfolio := strategyPortfolios[strategy.ID]
			s.updatePortfolioValue(ctx, portfolio, backtest.Symbols, currentDate)

			// 记录权益曲线
			strategyEquityCurves[strategy.ID] = append(strategyEquityCurves[strategy.ID], models.EquityPoint{
				Date:           currentDate.Format("2006-01-02"),
				PortfolioValue: portfolio.TotalValue,
				Cash:           portfolio.Cash,
				Holdings:       portfolio.TotalValue - portfolio.Cash,
			})

			// 计算日收益率
			if len(strategyEquityCurves[strategy.ID]) > 1 {
				prevValue := strategyEquityCurves[strategy.ID][len(strategyEquityCurves[strategy.ID])-2].PortfolioValue
				if prevValue > 0 {
					dailyReturn := (portfolio.TotalValue - prevValue) / prevValue
					strategyDailyReturns[strategy.ID] = append(strategyDailyReturns[strategy.ID], dailyReturn)
				}
			}
		}

		// 添加小延迟以避免过于频繁的操作
		time.Sleep(1 * time.Millisecond)
	}

	// 回测完成，计算和保存结果
	s.logger.Info("多策略回测循环完成，开始计算结果",
		logger.String("backtest_id", backtest.ID),
		logger.Int("strategies_count", len(strategies)),
	)

	// 为每个策略计算性能指标
	var allResults []models.BacktestResult
	var allTrades []models.Trade
	var combinedEquityCurve []models.EquityPoint

	for _, strategy := range strategies {
		// 计算该策略的性能指标
		performanceMetrics := &models.PerformanceMetrics{
			Returns:      strategyDailyReturns[strategy.ID],
			RiskFreeRate: 0.03 / 252, // 假设年化无风险利率3%
		}

		result := performanceMetrics.CalculateMetrics()
		result.ID = fmt.Sprintf("%s_%s", backtest.ID, strategy.ID)
		result.BacktestID = backtest.ID
		result.StrategyID = strategy.ID
		result.StrategyName = strategy.Name
		result.CreatedAt = time.Now()

		allResults = append(allResults, *result)

		// 合并交易记录
		allTrades = append(allTrades, strategyTrades[strategy.ID]...)

		s.logger.Info("策略性能计算完成",
			logger.String("backtest_id", backtest.ID),
			logger.String("strategy_id", strategy.ID),
			logger.String("strategy_name", strategy.Name),
			logger.Float64("total_return", result.TotalReturn),
			logger.Float64("sharpe_ratio", result.SharpeRatio),
			logger.Int("total_trades", result.TotalTrades),
		)
	}

	// 计算组合整体权益曲线（所有策略的平均或加权组合）
	if len(strategies) > 0 {
		maxLen := 0
		for _, curve := range strategyEquityCurves {
			if len(curve) > maxLen {
				maxLen = len(curve)
			}
		}

		for i := 0; i < maxLen; i++ {
			var totalValue, totalCash, totalHoldings float64
			var date string
			count := 0

			for _, curve := range strategyEquityCurves {
				if i < len(curve) {
					totalValue += curve[i].PortfolioValue
					totalCash += curve[i].Cash
					totalHoldings += curve[i].Holdings
					date = curve[i].Date
					count++
				}
			}

			if count > 0 {
				combinedEquityCurve = append(combinedEquityCurve, models.EquityPoint{
					Date:           date,
					PortfolioValue: totalValue / float64(count), // 平均值
					Cash:           totalCash / float64(count),
					Holdings:       totalHoldings / float64(count),
				})
			}
		}
	}

	// 保存结果到内存（在真实环境中应该保存到数据库）
	s.backtestTrades[backtest.ID] = allTrades

	// 保存多策略结果到新的存储结构
	s.backtestMultiResults[backtest.ID] = allResults
	s.backtestEquityCurves[backtest.ID] = combinedEquityCurve

	// 保存第一个策略的结果作为主结果（兼容性）
	if len(allResults) > 0 {
		s.backtestResults[backtest.ID] = &allResults[0]
	}

	// 更新回测状态
	backtest.Status = models.BacktestStatusCompleted
	backtest.Progress = 100
	now := time.Now()
	backtest.CompletedAt = &now

	s.updateBacktestStatus(backtest.ID, models.BacktestStatusCompleted, "多策略回测完成")

	s.logger.Info("🎉 多策略回测任务完成",
		logger.String("backtest_id", backtest.ID),
		logger.Int("strategies_count", len(strategies)),
		logger.Int("total_trades", len(allTrades)),
		logger.Int("equity_points", len(combinedEquityCurve)),
	)
}

// executeSignalForStrategy 为特定策略执行交易信号
func (s *BacktestService) executeSignalForStrategy(signal *models.Signal, marketData *models.MarketData, portfolio *models.Portfolio, backtest *models.Backtest, strategyID string) *models.Trade {
	if signal == nil || signal.SignalType == models.SignalTypeHold {
		return nil
	}

	symbol := marketData.Symbol
	price := marketData.Close

	switch signal.SignalType {
	case models.SignalTypeBuy:
		// 买入逻辑
		maxInvestment := portfolio.Cash * 0.2 // 每次最多投入20%的现金
		if maxInvestment < 1000 {             // 最小投资金额
			return nil
		}

		quantity := int(maxInvestment / price)
		if quantity <= 0 {
			return nil
		}

		cost := float64(quantity) * price
		commission := cost * backtest.Commission
		totalCost := cost + commission

		if totalCost > portfolio.Cash {
			return nil
		}

		// 执行买入
		portfolio.Cash -= totalCost
		if position, exists := portfolio.Positions[symbol]; exists {
			// 更新现有持仓
			totalShares := position.Quantity + quantity
			totalCost := position.AvgPrice*float64(position.Quantity) + cost
			position.AvgPrice = totalCost / float64(totalShares)
			position.Quantity = totalShares
			position.MarketValue = float64(totalShares) * price
			position.UnrealizedPL = position.MarketValue - totalCost
			portfolio.Positions[symbol] = position
		} else {
			// 创建新持仓
			portfolio.Positions[symbol] = models.Position{
				Symbol:       symbol,
				Quantity:     quantity,
				AvgPrice:     price,
				MarketValue:  float64(quantity) * price,
				UnrealizedPL: 0,
				Timestamp:    marketData.Date,
			}
		}

		return &models.Trade{
			ID:         fmt.Sprintf("%s_%s_%d", backtest.ID, symbol, time.Now().UnixNano()),
			BacktestID: backtest.ID,
			StrategyID: strategyID,
			Symbol:     symbol,
			Side:       models.TradeSideBuy,
			Quantity:   quantity,
			Price:      price,
			Commission: commission,
			SignalType: string(signal.SignalType),
			Timestamp:  marketData.Date,
			CreatedAt:  time.Now(),
		}

	case models.SignalTypeSell:
		// 卖出逻辑
		position, exists := portfolio.Positions[symbol]
		if !exists || position.Quantity <= 0 {
			return nil
		}

		quantity := position.Quantity
		revenue := float64(quantity) * price
		commission := revenue * backtest.Commission
		netRevenue := revenue - commission

		// 计算盈亏
		pnl := netRevenue - (position.AvgPrice * float64(quantity))

		// 执行卖出
		portfolio.Cash += netRevenue
		delete(portfolio.Positions, symbol)

		return &models.Trade{
			ID:         fmt.Sprintf("%s_%s_%d", backtest.ID, symbol, time.Now().UnixNano()),
			BacktestID: backtest.ID,
			StrategyID: strategyID,
			Symbol:     symbol,
			Side:       models.TradeSideSell,
			Quantity:   quantity,
			Price:      price,
			Commission: commission,
			PnL:        pnl,
			SignalType: string(signal.SignalType),
			Timestamp:  marketData.Date,
			CreatedAt:  time.Now(),
		}
	}

	return nil
}

// calculateCombinedMetrics 计算多策略组合的整体指标
func (s *BacktestService) calculateCombinedMetrics(results []models.BacktestResult) *models.BacktestResult {
	if len(results) == 0 {
		return nil
	}

	// 计算平均指标
	combined := &models.BacktestResult{
		ID:           "combined",
		BacktestID:   results[0].BacktestID,
		StrategyID:   "combined",
		StrategyName: "组合策略",
		CreatedAt:    time.Now(),
	}

	var totalReturn, annualReturn, maxDrawdown, sharpeRatio, sortinoRatio, winRate, profitFactor, avgTradeReturn, benchmarkReturn, alpha, beta float64
	var totalTrades int

	for _, result := range results {
		totalReturn += result.TotalReturn
		annualReturn += result.AnnualReturn
		maxDrawdown += result.MaxDrawdown
		sharpeRatio += result.SharpeRatio
		sortinoRatio += result.SortinoRatio
		winRate += result.WinRate
		profitFactor += result.ProfitFactor
		avgTradeReturn += result.AvgTradeReturn
		benchmarkReturn += result.BenchmarkReturn
		alpha += result.Alpha
		beta += result.Beta
		totalTrades += result.TotalTrades
	}

	count := float64(len(results))
	combined.TotalReturn = totalReturn / count
	combined.AnnualReturn = annualReturn / count
	combined.MaxDrawdown = maxDrawdown / count
	combined.SharpeRatio = sharpeRatio / count
	combined.SortinoRatio = sortinoRatio / count
	combined.WinRate = winRate / count
	combined.ProfitFactor = profitFactor / count
	combined.AvgTradeReturn = avgTradeReturn / count
	combined.BenchmarkReturn = benchmarkReturn / count
	combined.Alpha = alpha / count
	combined.Beta = beta / count
	combined.TotalTrades = totalTrades

	return combined
}
