package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"
)

var (
	ErrStrategyNotFound = errors.New("策略不存在")
	ErrStrategyExists   = errors.New("策略已存在")
)

// StrategyService 策略服务
type StrategyService struct {
	// 在真实环境中，这里会有数据库连接
	// 目前使用内存存储进行演示
	strategies map[string]*models.Strategy
	logger     logger.Logger
}

// NewStrategyService 创建策略服务
func NewStrategyService(log logger.Logger) *StrategyService {
	service := &StrategyService{
		strategies: make(map[string]*models.Strategy),
		logger:     log,
	}

	// 初始化默认策略
	service.initDefaultStrategies()

	return service
}

// initDefaultStrategies 初始化默认策略
func (s *StrategyService) initDefaultStrategies() {
	for _, strategy := range models.DefaultStrategies {
		// 创建副本避免指针问题
		strategyCopy := strategy
		s.strategies[strategy.ID] = &strategyCopy
	}

	s.logger.Info("默认策略初始化完成", logger.Int("count", len(models.DefaultStrategies)))
}

// GetStrategiesList 获取策略列表
func (s *StrategyService) GetStrategiesList(ctx context.Context, req *models.StrategyListRequest) ([]models.Strategy, int, error) {
	var results []models.Strategy

	// 过滤策略
	for _, strategy := range s.strategies {
		if s.matchesFilter(strategy, req) {
			results = append(results, *strategy)
		}
	}

	total := len(results)

	// 分页
	start := (req.Page - 1) * req.Size
	end := start + req.Size

	if start >= len(results) {
		return []models.Strategy{}, total, nil
	}

	if end > len(results) {
		end = len(results)
	}

	return results[start:end], total, nil
}

// matchesFilter 检查策略是否匹配过滤条件
func (s *StrategyService) matchesFilter(strategy *models.Strategy, req *models.StrategyListRequest) bool {
	// 状态过滤
	if req.Status != "" && strategy.Status != req.Status {
		return false
	}

	// 类型过滤
	if req.Type != "" && strategy.Type != req.Type {
		return false
	}

	// 关键词过滤
	if req.Keyword != "" {
		keyword := strings.ToLower(req.Keyword)
		name := strings.ToLower(strategy.Name)
		desc := strings.ToLower(strategy.Description)

		if !strings.Contains(name, keyword) && !strings.Contains(desc, keyword) {
			return false
		}
	}

	return true
}

// GetStrategy 获取策略详情
func (s *StrategyService) GetStrategy(ctx context.Context, strategyID string) (*models.Strategy, error) {
	strategy, exists := s.strategies[strategyID]
	if !exists {
		return nil, ErrStrategyNotFound
	}

	// 返回副本避免外部修改
	strategyCopy := *strategy
	return &strategyCopy, nil
}

// CreateStrategy 创建策略
func (s *StrategyService) CreateStrategy(ctx context.Context, strategy *models.Strategy) error {
	// 检查策略是否已存在
	if _, exists := s.strategies[strategy.ID]; exists {
		return ErrStrategyExists
	}

	// 检查名称是否重复
	for _, existing := range s.strategies {
		if existing.Name == strategy.Name {
			return fmt.Errorf("策略名称已存在: %s", strategy.Name)
		}
	}

	// 设置时间戳
	now := time.Now()
	strategy.CreatedAt = now
	strategy.UpdatedAt = now

	// 验证策略参数
	if err := s.validateStrategyParameters(strategy); err != nil {
		return fmt.Errorf("策略参数验证失败: %w", err)
	}

	// 保存策略
	s.strategies[strategy.ID] = strategy

	s.logger.Info("策略创建成功",
		logger.String("strategy_id", strategy.ID),
		logger.String("strategy_name", strategy.Name),
		logger.String("strategy_type", string(strategy.Type)),
	)

	return nil
}

// UpdateStrategy 更新策略
func (s *StrategyService) UpdateStrategy(ctx context.Context, strategyID string, req *models.UpdateStrategyRequest) error {
	strategy, exists := s.strategies[strategyID]
	if !exists {
		return ErrStrategyNotFound
	}

	// 更新字段
	if req.Name != nil {
		// 检查名称是否重复
		for id, existing := range s.strategies {
			if id != strategyID && existing.Name == *req.Name {
				return fmt.Errorf("策略名称已存在: %s", *req.Name)
			}
		}
		strategy.Name = *req.Name
	}

	if req.Description != nil {
		strategy.Description = *req.Description
	}

	if req.Status != nil {
		strategy.Status = *req.Status
	}

	if req.Parameters != nil {
		strategy.Parameters = *req.Parameters
		// 验证更新后的参数
		if err := s.validateStrategyParameters(strategy); err != nil {
			return fmt.Errorf("策略参数验证失败: %w", err)
		}
	}

	if req.Code != nil {
		strategy.Code = *req.Code
	}

	// 更新时间戳
	strategy.UpdatedAt = time.Now()

	s.logger.Info("策略更新成功",
		logger.String("strategy_id", strategyID),
		logger.String("strategy_name", strategy.Name),
	)

	return nil
}

// DeleteStrategy 删除策略
func (s *StrategyService) DeleteStrategy(ctx context.Context, strategyID string) error {
	if _, exists := s.strategies[strategyID]; !exists {
		return ErrStrategyNotFound
	}

	delete(s.strategies, strategyID)

	s.logger.Info("策略删除成功", logger.String("strategy_id", strategyID))

	return nil
}

// UpdateStrategyStatus 更新策略状态
func (s *StrategyService) UpdateStrategyStatus(ctx context.Context, strategyID string, status models.StrategyStatus) error {
	strategy, exists := s.strategies[strategyID]
	if !exists {
		return ErrStrategyNotFound
	}

	oldStatus := strategy.Status
	strategy.Status = status
	strategy.UpdatedAt = time.Now()

	s.logger.Info("策略状态更新成功",
		logger.String("strategy_id", strategyID),
		logger.String("old_status", string(oldStatus)),
		logger.String("new_status", string(status)),
	)

	return nil
}

// GetStrategyPerformance 获取策略表现
func (s *StrategyService) GetStrategyPerformance(ctx context.Context, strategyID string) (*models.StrategyPerformance, error) {
	strategy, exists := s.strategies[strategyID]
	if !exists {
		return nil, ErrStrategyNotFound
	}

	// 生成模拟性能数据
	performance := s.generateMockPerformance(strategy)

	return performance, nil
}

// generateMockPerformance 生成模拟性能数据
func (s *StrategyService) generateMockPerformance(strategy *models.Strategy) *models.StrategyPerformance {
	// 根据策略类型生成不同的模拟数据
	baseReturn := 0.15 // 基础收益率15%

	switch strategy.Type {
	case models.StrategyTypeTechnical:
		baseReturn = 0.12
	case models.StrategyTypeFundamental:
		baseReturn = 0.18
	case models.StrategyTypeML:
		baseReturn = 0.22
	case models.StrategyTypeComposite:
		baseReturn = 0.20
	}

	// 根据策略状态调整表现
	switch strategy.Status {
	case models.StrategyStatusActive:
		baseReturn *= 1.1 // 活跃策略表现更好
	case models.StrategyStatusTesting:
		baseReturn *= 0.9 // 测试中的策略表现稍差
	case models.StrategyStatusInactive:
		baseReturn *= 0.8 // 非活跃策略表现较差
	}

	// 添加一些随机性
	randomFactor := 0.8 + rand.Float64()*0.4 // 0.8 - 1.2
	baseReturn *= randomFactor

	return &models.StrategyPerformance{
		ID:              fmt.Sprintf("perf_%s", strategy.ID),
		StrategyID:      strategy.ID,
		TotalReturn:     baseReturn,
		AnnualReturn:    baseReturn * 0.85,          // 年化收益稍低
		MaxDrawdown:     -baseReturn * 0.3,          // 最大回撤为收益的30%
		SharpeRatio:     1.2 + rand.Float64()*0.8,   // 1.2-2.0
		SortinoRatio:    1.5 + rand.Float64()*1.0,   // 1.5-2.5
		WinRate:         0.55 + rand.Float64()*0.25, // 55%-80%
		ProfitFactor:    1.3 + rand.Float64()*0.7,   // 1.3-2.0
		TotalTrades:     100 + rand.Intn(400),       // 100-500次交易
		AvgTradeReturn:  baseReturn / 200,           // 平均每笔交易收益
		BenchmarkReturn: 0.08,                       // 基准收益8%
		Alpha:           baseReturn - 0.08,          // Alpha = 策略收益 - 基准收益
		Beta:            0.8 + rand.Float64()*0.4,   // 0.8-1.2
		LastUpdated:     time.Now(),
	}
}

// validateStrategyParameters 验证策略参数
func (s *StrategyService) validateStrategyParameters(strategy *models.Strategy) error {
	if strategy.Parameters == nil {
		return nil
	}

	switch strategy.Type {
	case models.StrategyTypeTechnical:
		return s.validateTechnicalStrategyParameters(strategy)
	case models.StrategyTypeFundamental:
		return s.validateFundamentalStrategyParameters(strategy)
	case models.StrategyTypeML:
		return s.validateMLStrategyParameters(strategy)
	case models.StrategyTypeComposite:
		return s.validateCompositeStrategyParameters(strategy)
	}

	return nil
}

// validateTechnicalStrategyParameters 验证技术指标策略参数
func (s *StrategyService) validateTechnicalStrategyParameters(strategy *models.Strategy) error {
	params := strategy.Parameters

	// 根据策略ID判断具体的技术指标类型
	switch strategy.ID {
	case "macd_strategy":
		return s.validateMACDParameters(params)
	case "ma_crossover":
		return s.validateMAParameters(params)
	case "rsi_strategy":
		return s.validateRSIParameters(params)
	case "bollinger_strategy":
		return s.validateBollingerParameters(params)
	}

	return nil
}

// validateMACDParameters 验证MACD参数
func (s *StrategyService) validateMACDParameters(params map[string]interface{}) error {
	fastPeriod, ok := params["fast_period"].(float64)
	if !ok || fastPeriod < 1 || fastPeriod > 50 {
		return errors.New("MACD快线周期必须在1-50之间")
	}

	slowPeriod, ok := params["slow_period"].(float64)
	if !ok || slowPeriod < 1 || slowPeriod > 100 {
		return errors.New("MACD慢线周期必须在1-100之间")
	}

	if fastPeriod >= slowPeriod {
		return errors.New("MACD快线周期必须小于慢线周期")
	}

	signalPeriod, ok := params["signal_period"].(float64)
	if !ok || signalPeriod < 1 || signalPeriod > 50 {
		return errors.New("MACD信号线周期必须在1-50之间")
	}

	return nil
}

// validateMAParameters 验证移动平均参数
func (s *StrategyService) validateMAParameters(params map[string]interface{}) error {
	shortPeriod, ok := params["short_period"].(float64)
	if !ok || shortPeriod < 1 || shortPeriod > 50 {
		return errors.New("短期均线周期必须在1-50之间")
	}

	longPeriod, ok := params["long_period"].(float64)
	if !ok || longPeriod < 1 || longPeriod > 200 {
		return errors.New("长期均线周期必须在1-200之间")
	}

	if shortPeriod >= longPeriod {
		return errors.New("短期均线周期必须小于长期均线周期")
	}

	maType, ok := params["ma_type"].(string)
	if !ok {
		return errors.New("均线类型不能为空")
	}

	validMATypes := []string{"sma", "ema", "wma"}
	isValid := false
	for _, validType := range validMATypes {
		if maType == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		return errors.New("均线类型必须是sma、ema或wma")
	}

	return nil
}

// validateRSIParameters 验证RSI参数
func (s *StrategyService) validateRSIParameters(params map[string]interface{}) error {
	period, ok := params["period"].(float64)
	if !ok || period < 1 || period > 50 {
		return errors.New("RSI周期必须在1-50之间")
	}

	overbought, ok := params["overbought"].(float64)
	if !ok || overbought < 50 || overbought > 100 {
		return errors.New("RSI超买阈值必须在50-100之间")
	}

	oversold, ok := params["oversold"].(float64)
	if !ok || oversold < 0 || oversold > 50 {
		return errors.New("RSI超卖阈值必须在0-50之间")
	}

	if oversold >= overbought {
		return errors.New("RSI超卖阈值必须小于超买阈值")
	}

	return nil
}

// validateBollingerParameters 验证布林带参数
func (s *StrategyService) validateBollingerParameters(params map[string]interface{}) error {
	period, ok := params["period"].(float64)
	if !ok || period < 1 || period > 50 {
		return errors.New("布林带周期必须在1-50之间")
	}

	stdDev, ok := params["std_dev"].(float64)
	if !ok || stdDev < 0.5 || stdDev > 5 {
		return errors.New("布林带标准差倍数必须在0.5-5之间")
	}

	return nil
}

// validateFundamentalStrategyParameters 验证基本面策略参数
func (s *StrategyService) validateFundamentalStrategyParameters(strategy *models.Strategy) error {
	// 基本面策略参数验证
	// 这里可以添加具体的基本面指标验证逻辑
	return nil
}

// validateMLStrategyParameters 验证机器学习策略参数
func (s *StrategyService) validateMLStrategyParameters(strategy *models.Strategy) error {
	// 机器学习策略参数验证
	// 这里可以添加具体的ML模型参数验证逻辑
	return nil
}

// validateCompositeStrategyParameters 验证复合策略参数
func (s *StrategyService) validateCompositeStrategyParameters(strategy *models.Strategy) error {
	// 复合策略参数验证
	// 这里可以添加具体的复合策略参数验证逻辑
	return nil
}

// ExecuteStrategy 执行策略（生成交易信号）
func (s *StrategyService) ExecuteStrategy(ctx context.Context, strategyID string, marketData *models.MarketData) (*models.Signal, error) {
	strategy, exists := s.strategies[strategyID]
	if !exists {
		return nil, ErrStrategyNotFound
	}

	if strategy.Status != models.StrategyStatusActive {
		return nil, fmt.Errorf("策略未激活: %s", strategy.Status)
	}

	// 根据策略类型执行不同的逻辑
	switch strategy.ID {
	case "macd_strategy":
		return s.executeMACDStrategy(strategy, marketData)
	case "ma_crossover":
		return s.executeMAStrategy(strategy, marketData)
	case "rsi_strategy":
		return s.executeRSIStrategy(strategy, marketData)
	case "bollinger_strategy":
		return s.executeBollingerStrategy(strategy, marketData)
	}

	return nil, fmt.Errorf("未知的策略类型: %s", strategy.ID)
}

// executeMACDStrategy 执行MACD策略
func (s *StrategyService) executeMACDStrategy(strategy *models.Strategy, marketData *models.MarketData) (*models.Signal, error) {
	// 这里是MACD策略的简化实现
	// 在实际应用中，需要计算MACD指标并根据金叉死叉生成信号

	// 模拟MACD计算结果
	// 实际实现需要历史数据来计算MACD线、信号线和柱状图

	signal := &models.Signal{
		ID:         fmt.Sprintf("signal_%d", time.Now().Unix()),
		StrategyID: strategy.ID,
		Symbol:     marketData.Symbol,
		Price:      marketData.Close,
		Timestamp:  marketData.Date,
		CreatedAt:  time.Now(),
	}

	// 简化的信号生成逻辑（实际需要技术指标计算）
	// 这里使用随机数模拟

	rand.Seed(time.Now().UnixNano())
	signalValue := rand.Float64()

	if signalValue > 0.7 {
		signal.SignalType = models.SignalTypeBuy
		signal.Side = models.TradeSideBuy
		signal.Strength = signalValue
		signal.Confidence = 0.8
		signal.Reason = "MACD金叉信号"
	} else if signalValue < 0.3 {
		signal.SignalType = models.SignalTypeSell
		signal.Side = models.TradeSideSell
		signal.Strength = 1 - signalValue
		signal.Confidence = 0.8
		signal.Reason = "MACD死叉信号"
	} else {
		signal.SignalType = models.SignalTypeHold
		signal.Strength = 0.5
		signal.Confidence = 0.5
		signal.Reason = "MACD震荡区间"
	}

	return signal, nil
}

// executeMAStrategy 执行移动平均策略
func (s *StrategyService) executeMAStrategy(strategy *models.Strategy, marketData *models.MarketData) (*models.Signal, error) {
	// 移动平均策略的简化实现
	// 实际需要计算短期和长期移动平均线

	signal := &models.Signal{
		ID:         fmt.Sprintf("signal_%d", time.Now().Unix()),
		StrategyID: strategy.ID,
		Symbol:     marketData.Symbol,
		Price:      marketData.Close,
		Timestamp:  marketData.Date,
		CreatedAt:  time.Now(),
	}

	// 简化的信号生成逻辑

	signalValue := rand.Float64()

	if signalValue > 0.6 {
		signal.SignalType = models.SignalTypeBuy
		signal.Side = models.TradeSideBuy
		signal.Strength = signalValue
		signal.Confidence = 0.75
		signal.Reason = "短期均线上穿长期均线"
	} else if signalValue < 0.4 {
		signal.SignalType = models.SignalTypeSell
		signal.Side = models.TradeSideSell
		signal.Strength = 1 - signalValue
		signal.Confidence = 0.75
		signal.Reason = "短期均线下穿长期均线"
	} else {
		signal.SignalType = models.SignalTypeHold
		signal.Strength = 0.5
		signal.Confidence = 0.5
		signal.Reason = "均线纠缠状态"
	}

	return signal, nil
}

// executeRSIStrategy 执行RSI策略
func (s *StrategyService) executeRSIStrategy(strategy *models.Strategy, marketData *models.MarketData) (*models.Signal, error) {
	// RSI策略的简化实现

	signal := &models.Signal{
		ID:         fmt.Sprintf("signal_%d", time.Now().Unix()),
		StrategyID: strategy.ID,
		Symbol:     marketData.Symbol,
		Price:      marketData.Close,
		Timestamp:  marketData.Date,
		CreatedAt:  time.Now(),
	}

	// 模拟RSI值

	rsiValue := rand.Float64() * 100

	overbought := 70.0
	oversold := 30.0

	if params, ok := strategy.Parameters["overbought"].(float64); ok {
		overbought = params
	}

	if params, ok := strategy.Parameters["oversold"].(float64); ok {
		oversold = params
	}

	if rsiValue < oversold {
		signal.SignalType = models.SignalTypeBuy
		signal.Side = models.TradeSideBuy
		signal.Strength = (oversold - rsiValue) / oversold
		signal.Confidence = 0.8
		signal.Reason = fmt.Sprintf("RSI超卖信号 (RSI: %.2f)", rsiValue)
	} else if rsiValue > overbought {
		signal.SignalType = models.SignalTypeSell
		signal.Side = models.TradeSideSell
		signal.Strength = (rsiValue - overbought) / (100 - overbought)
		signal.Confidence = 0.8
		signal.Reason = fmt.Sprintf("RSI超买信号 (RSI: %.2f)", rsiValue)
	} else {
		signal.SignalType = models.SignalTypeHold
		signal.Strength = 0.5
		signal.Confidence = 0.5
		signal.Reason = fmt.Sprintf("RSI正常区间 (RSI: %.2f)", rsiValue)
	}

	return signal, nil
}

// executeBollingerStrategy 执行布林带策略
func (s *StrategyService) executeBollingerStrategy(strategy *models.Strategy, marketData *models.MarketData) (*models.Signal, error) {
	// 布林带策略的简化实现

	signal := &models.Signal{
		ID:         fmt.Sprintf("signal_%d", time.Now().Unix()),
		StrategyID: strategy.ID,
		Symbol:     marketData.Symbol,
		Price:      marketData.Close,
		Timestamp:  marketData.Date,
		CreatedAt:  time.Now(),
	}

	// 模拟布林带计算

	// 假设当前价格相对于布林带的位置
	position := rand.Float64() // 0表示下轨，0.5表示中轨，1表示上轨

	if position < 0.1 {
		signal.SignalType = models.SignalTypeBuy
		signal.Side = models.TradeSideBuy
		signal.Strength = 1 - position*10 // 越接近下轨强度越大
		signal.Confidence = 0.75
		signal.Reason = "价格触及布林带下轨，超卖信号"
	} else if position > 0.9 {
		signal.SignalType = models.SignalTypeSell
		signal.Side = models.TradeSideSell
		signal.Strength = (position - 0.9) * 10 // 越接近上轨强度越大
		signal.Confidence = 0.75
		signal.Reason = "价格触及布林带上轨，超买信号"
	} else {
		signal.SignalType = models.SignalTypeHold
		signal.Strength = 0.5
		signal.Confidence = 0.5
		signal.Reason = "价格在布林带中轨附近"
	}

	return signal, nil
}
