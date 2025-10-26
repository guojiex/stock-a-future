package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
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

	// 简单稳定排序：按ID字母顺序排序（保证每次顺序一致）
	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})

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

	// 移除基于策略状态的表现调整 - 策略表现应该只基于策略本身的逻辑

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

// GetStrategyTemplates 获取策略模板列表
func (s *StrategyService) GetStrategyTemplates() []models.StrategyTemplate {
	return []models.StrategyTemplate{
		{
			ID:          "macd_template",
			Name:        "MACD金叉策略模板",
			Description: "经典的MACD金叉死叉交易策略",
			Type:        models.StrategyTypeTechnical,
			Parameters: map[string]interface{}{
				"fast_period":    12,
				"slow_period":    26,
				"signal_period":  9,
				"buy_threshold":  0.0,
				"sell_threshold": 0.0,
			},
			Category:  "技术指标",
			Tags:      []string{"MACD", "趋势跟踪", "金叉"},
			CreatedAt: time.Now(),
		},
		{
			ID:          "ma_crossover_template",
			Name:        "双均线策略模板",
			Description: "短期均线突破长期均线的交易策略",
			Type:        models.StrategyTypeTechnical,
			Parameters: map[string]interface{}{
				"short_period": 5,
				"long_period":  20,
				"ma_type":      "sma",
				"threshold":    0.01,
			},
			Category:  "技术指标",
			Tags:      []string{"均线", "趋势跟踪", "突破"},
			CreatedAt: time.Now(),
		},
		{
			ID:          "rsi_template",
			Name:        "RSI超买超卖策略模板",
			Description: "基于RSI指标的超买超卖交易策略",
			Type:        models.StrategyTypeTechnical,
			Parameters: map[string]interface{}{
				"period":     14,
				"overbought": 70.0,
				"oversold":   30.0,
			},
			Category:  "技术指标",
			Tags:      []string{"RSI", "超买超卖", "震荡指标"},
			CreatedAt: time.Now(),
		},
		{
			ID:          "bollinger_template",
			Name:        "布林带策略模板",
			Description: "基于布林带的均值回归策略",
			Type:        models.StrategyTypeTechnical,
			Parameters: map[string]interface{}{
				"period":  20,
				"std_dev": 2.0,
			},
			Category:  "技术指标",
			Tags:      []string{"布林带", "均值回归", "波动率"},
			CreatedAt: time.Now(),
		},
	}
}

// ValidateParameters 验证策略参数
func (s *StrategyService) ValidateParameters(strategyType models.StrategyType, parameters map[string]interface{}) []map[string]string {
	var errors []map[string]string

	// 根据策略类型验证参数
	switch strategyType {
	case models.StrategyTypeTechnical:
		errors = s.validateTechnicalParams(parameters)
	case models.StrategyTypeFundamental:
		errors = s.validateFundamentalParams(parameters)
	case models.StrategyTypeML:
		errors = s.validateMLParams(parameters)
	case models.StrategyTypeComposite:
		errors = s.validateCompositeParams(parameters)
	}

	return errors
}

// validateTechnicalParams 验证技术指标策略参数
func (s *StrategyService) validateTechnicalParams(parameters map[string]interface{}) []map[string]string {
	var errors []map[string]string

	// MACD参数验证
	if fastPeriod, ok := parameters["fast_period"].(float64); ok {
		if fastPeriod < 1 || fastPeriod > 50 {
			errors = append(errors, map[string]string{
				"field":   "fast_period",
				"message": "快线周期必须在1-50之间",
			})
		}
	}

	if slowPeriod, ok := parameters["slow_period"].(float64); ok {
		if slowPeriod < 1 || slowPeriod > 100 {
			errors = append(errors, map[string]string{
				"field":   "slow_period",
				"message": "慢线周期必须在1-100之间",
			})
		}
	}

	// 验证快线周期必须小于慢线周期
	if fastPeriod, ok1 := parameters["fast_period"].(float64); ok1 {
		if slowPeriod, ok2 := parameters["slow_period"].(float64); ok2 {
			if fastPeriod >= slowPeriod {
				errors = append(errors, map[string]string{
					"field":   "slow_period",
					"message": "慢线周期必须大于快线周期",
				})
			}
		}
	}

	if signalPeriod, ok := parameters["signal_period"].(float64); ok {
		if signalPeriod < 1 || signalPeriod > 50 {
			errors = append(errors, map[string]string{
				"field":   "signal_period",
				"message": "信号线周期必须在1-50之间",
			})
		}
	}

	// 双均线参数验证
	if shortPeriod, ok := parameters["short_period"].(float64); ok {
		if shortPeriod < 1 || shortPeriod > 50 {
			errors = append(errors, map[string]string{
				"field":   "short_period",
				"message": "短期均线周期必须在1-50之间",
			})
		}
	}

	if longPeriod, ok := parameters["long_period"].(float64); ok {
		if longPeriod < 1 || longPeriod > 200 {
			errors = append(errors, map[string]string{
				"field":   "long_period",
				"message": "长期均线周期必须在1-200之间",
			})
		}
	}

	// 验证短期均线必须小于长期均线
	if shortPeriod, ok1 := parameters["short_period"].(float64); ok1 {
		if longPeriod, ok2 := parameters["long_period"].(float64); ok2 {
			if shortPeriod >= longPeriod {
				errors = append(errors, map[string]string{
					"field":   "long_period",
					"message": "长期均线周期必须大于短期均线周期",
				})
			}
		}
	}

	// RSI参数验证
	if period, ok := parameters["period"].(float64); ok {
		if period < 1 || period > 50 {
			errors = append(errors, map[string]string{
				"field":   "period",
				"message": "RSI周期必须在1-50之间",
			})
		}
	}

	if overbought, ok := parameters["overbought"].(float64); ok {
		if overbought < 50 || overbought > 100 {
			errors = append(errors, map[string]string{
				"field":   "overbought",
				"message": "超买阈值必须在50-100之间",
			})
		}
	}

	if oversold, ok := parameters["oversold"].(float64); ok {
		if oversold < 0 || oversold > 50 {
			errors = append(errors, map[string]string{
				"field":   "oversold",
				"message": "超卖阈值必须在0-50之间",
			})
		}
	}

	// 布林带参数验证
	if stdDev, ok := parameters["std_dev"].(float64); ok {
		if stdDev < 0.5 || stdDev > 5 {
			errors = append(errors, map[string]string{
				"field":   "std_dev",
				"message": "标准差倍数必须在0.5-5之间",
			})
		}
	}

	return errors
}

// validateFundamentalParams 验证基本面策略参数
func (s *StrategyService) validateFundamentalParams(parameters map[string]interface{}) []map[string]string {
	var errors []map[string]string
	// TODO: 实现基本面参数验证
	return errors
}

// validateMLParams 验证机器学习策略参数
func (s *StrategyService) validateMLParams(parameters map[string]interface{}) []map[string]string {
	var errors []map[string]string
	// TODO: 实现机器学习参数验证
	return errors
}

// validateCompositeParams 验证复合策略参数
func (s *StrategyService) validateCompositeParams(parameters map[string]interface{}) []map[string]string {
	var errors []map[string]string
	// TODO: 实现复合策略参数验证
	return errors
}

// GetStrategyTypeDefinitions 获取策略类型定义
func (s *StrategyService) GetStrategyTypeDefinitions() []models.StrategyTypeDefinition {
	return []models.StrategyTypeDefinition{
		{
			Type:        models.StrategyTypeTechnical,
			Name:        "技术指标策略",
			Description: "基于技术指标的交易策略",
			Parameters: []models.ParameterDefinition{
				{
					Name:         "fast_period",
					DisplayName:  "快线周期",
					Type:         "int",
					DefaultValue: 12,
					MinValue:     1,
					MaxValue:     50,
					Required:     true,
					Description:  "MACD快线的计算周期，通常为12天",
				},
				{
					Name:         "slow_period",
					DisplayName:  "慢线周期",
					Type:         "int",
					DefaultValue: 26,
					MinValue:     1,
					MaxValue:     100,
					Required:     true,
					Description:  "MACD慢线的计算周期，通常为26天",
				},
				{
					Name:         "signal_period",
					DisplayName:  "信号线周期",
					Type:         "int",
					DefaultValue: 9,
					MinValue:     1,
					MaxValue:     50,
					Required:     true,
					Description:  "信号线的计算周期，通常为9天",
				},
			},
		},
		{
			Type:        models.StrategyTypeFundamental,
			Name:        "基本面策略",
			Description: "基于公司基本面数据的投资策略",
			Parameters:  []models.ParameterDefinition{},
		},
		{
			Type:        models.StrategyTypeML,
			Name:        "机器学习策略",
			Description: "基于机器学习模型的预测策略",
			Parameters:  []models.ParameterDefinition{},
		},
		{
			Type:        models.StrategyTypeComposite,
			Name:        "复合策略",
			Description: "结合多种策略类型的综合策略",
			Parameters:  []models.ParameterDefinition{},
		},
	}
}

// ExecuteStrategy 执行策略（生成交易信号）
func (s *StrategyService) ExecuteStrategy(ctx context.Context, strategyID string, marketData *models.MarketData) (*models.Signal, error) {
	strategy, exists := s.strategies[strategyID]
	if !exists {
		return nil, ErrStrategyNotFound
	}

	// 移除策略状态检查 - 任何定义好的策略都应该可以执行

	// 根据策略类型执行不同的逻辑
	switch strategy.ID {
	case "macd_strategy":
		return s.executeMACDStrategyImproved(strategy, marketData)
	case "ma_crossover":
		return s.executeMAStrategyImproved(strategy, marketData)
	case "rsi_strategy":
		return s.executeRSIStrategyImproved(strategy, marketData)
	case "bollinger_strategy":
		return s.executeBollingerStrategyImproved(strategy, marketData)
	}

	s.logger.Error("未知的策略类型",
		logger.String("strategy_id", strategy.ID),
		logger.String("strategy_name", strategy.Name),
	)
	return nil, fmt.Errorf("未知的策略类型: %s", strategy.ID)
}

// ==================== 改进的策略实现 ====================

// executeMACDStrategyImproved 改进的MACD策略
func (s *StrategyService) executeMACDStrategyImproved(strategy *models.Strategy, marketData *models.MarketData) (*models.Signal, error) {
	signal := &models.Signal{
		ID:         fmt.Sprintf("signal_%d", time.Now().Unix()),
		StrategyID: strategy.ID,
		Symbol:     marketData.Symbol,
		Price:      marketData.Close,
		Timestamp:  marketData.Date,
		CreatedAt:  time.Now(),
	}

	// 基于价格趋势的简化MACD逻辑
	// 使用价格相对于某个基准的变化来模拟MACD信号
	priceBase := 5.0 // 假设基准价格
	priceChange := (marketData.Close - priceBase) / priceBase

	// 添加一些随机性，但更有逻辑
	randomFactor := (rand.Float64() - 0.5) * 0.2 // -0.1 到 0.1 的随机因子
	adjustedChange := priceChange + randomFactor

	// 更严格的买入/卖出条件，减少频繁交易
	if adjustedChange > 0.05 { // 价格上涨超过5%
		signal.SignalType = models.SignalTypeBuy
		signal.Side = models.TradeSideBuy
		signal.Strength = math.Min(adjustedChange*2, 1.0)
		signal.Confidence = 0.75
		signal.Reason = fmt.Sprintf("MACD金叉信号 (价格变化: %.2f%%)", adjustedChange*100)
	} else if adjustedChange < -0.05 { // 价格下跌超过5%
		signal.SignalType = models.SignalTypeSell
		signal.Side = models.TradeSideSell
		signal.Strength = math.Min(-adjustedChange*2, 1.0)
		signal.Confidence = 0.75
		signal.Reason = fmt.Sprintf("MACD死叉信号 (价格变化: %.2f%%)", adjustedChange*100)
	} else {
		signal.SignalType = models.SignalTypeHold
		signal.Strength = 0.5
		signal.Confidence = 0.5
		signal.Reason = "MACD震荡区间，等待明确信号"
	}

	return signal, nil
}

// executeMAStrategyImproved 改进的移动平均策略
func (s *StrategyService) executeMAStrategyImproved(strategy *models.Strategy, marketData *models.MarketData) (*models.Signal, error) {
	signal := &models.Signal{
		ID:         fmt.Sprintf("signal_%d", time.Now().Unix()),
		StrategyID: strategy.ID,
		Symbol:     marketData.Symbol,
		Price:      marketData.Close,
		Timestamp:  marketData.Date,
		CreatedAt:  time.Now(),
	}

	// 模拟短期和长期均线
	shortMA := marketData.Close * (0.95 + rand.Float64()*0.1) // 短期均线在价格附近波动
	longMA := marketData.Close * (0.90 + rand.Float64()*0.2)  // 长期均线波动更大

	crossover := (marketData.Close - shortMA) / shortMA
	trend := (shortMA - longMA) / longMA

	// 更智能的信号生成：需要价格突破且趋势确认
	if crossover > 0.02 && trend > 0.03 { // 价格突破短期均线且短期均线在长期均线之上
		signal.SignalType = models.SignalTypeBuy
		signal.Side = models.TradeSideBuy
		signal.Strength = math.Min((crossover+trend)*5, 1.0)
		signal.Confidence = 0.8
		signal.Reason = "短期均线上穿长期均线，趋势向上"
	} else if crossover < -0.02 && trend < -0.03 { // 价格跌破短期均线且短期均线在长期均线之下
		signal.SignalType = models.SignalTypeSell
		signal.Side = models.TradeSideSell
		signal.Strength = math.Min(-(crossover+trend)*5, 1.0)
		signal.Confidence = 0.8
		signal.Reason = "短期均线下穿长期均线，趋势向下"
	} else {
		signal.SignalType = models.SignalTypeHold
		signal.Strength = 0.5
		signal.Confidence = 0.5
		signal.Reason = "均线纠缠状态，观望"
	}

	return signal, nil
}

// executeRSIStrategyImproved 改进的RSI策略
func (s *StrategyService) executeRSIStrategyImproved(strategy *models.Strategy, marketData *models.MarketData) (*models.Signal, error) {
	signal := &models.Signal{
		ID:         fmt.Sprintf("signal_%d", time.Now().Unix()),
		StrategyID: strategy.ID,
		Symbol:     marketData.Symbol,
		Price:      marketData.Close,
		Timestamp:  marketData.Date,
		CreatedAt:  time.Now(),
	}

	// 基于价格位置模拟RSI值
	priceRange := marketData.High - marketData.Low
	if priceRange == 0 {
		priceRange = marketData.Close * 0.02 // 假设2%的波动
	}

	pricePosition := (marketData.Close - marketData.Low) / priceRange
	rsiValue := 30 + pricePosition*40 + (rand.Float64()-0.5)*20 // 30-70区间，加随机波动

	overbought := 75.0
	oversold := 25.0

	if params, ok := strategy.Parameters["overbought"].(float64); ok {
		overbought = params
	}
	if params, ok := strategy.Parameters["oversold"].(float64); ok {
		oversold = params
	}

	// 更严格的RSI条件，避免频繁交易
	if rsiValue < oversold && pricePosition < 0.3 { // RSI超卖且价格在低位
		signal.SignalType = models.SignalTypeBuy
		signal.Side = models.TradeSideBuy
		signal.Strength = (oversold - rsiValue) / oversold
		signal.Confidence = 0.85
		signal.Reason = fmt.Sprintf("RSI超卖信号 (RSI: %.1f, 价格低位)", rsiValue)
	} else if rsiValue > overbought && pricePosition > 0.7 { // RSI超买且价格在高位
		signal.SignalType = models.SignalTypeSell
		signal.Side = models.TradeSideSell
		signal.Strength = (rsiValue - overbought) / (100 - overbought)
		signal.Confidence = 0.85
		signal.Reason = fmt.Sprintf("RSI超买信号 (RSI: %.1f, 价格高位)", rsiValue)
	} else {
		signal.SignalType = models.SignalTypeHold
		signal.Strength = 0.5
		signal.Confidence = 0.5
		signal.Reason = fmt.Sprintf("RSI正常区间 (RSI: %.1f)", rsiValue)
	}

	return signal, nil
}

// executeBollingerStrategyImproved 改进的布林带策略
func (s *StrategyService) executeBollingerStrategyImproved(strategy *models.Strategy, marketData *models.MarketData) (*models.Signal, error) {
	signal := &models.Signal{
		ID:         fmt.Sprintf("signal_%d", time.Now().Unix()),
		StrategyID: strategy.ID,
		Symbol:     marketData.Symbol,
		Price:      marketData.Close,
		Timestamp:  marketData.Date,
		CreatedAt:  time.Now(),
	}

	// 获取策略参数
	stdDevMultiplier := 2.0
	if params, ok := strategy.Parameters["std_dev"].(float64); ok {
		stdDevMultiplier = params
	}
	// period 参数暂时未使用，因为我们使用简化的布林带计算

	// 简化的布林带计算：使用当前价格的移动平均估算
	// 在实际应用中，这里应该获取历史数据来计算真正的移动平均线
	// 为了演示，我们使用一个简化的方法：基于当日的开高低收来估算趋势

	// 估算移动平均线（中轨）：使用当日均价作为近似
	middleBand := (marketData.High + marketData.Low + marketData.Close*2) / 4

	// 计算当日波动率作为标准差的估算
	dailyRange := marketData.High - marketData.Low
	estimatedStdDev := dailyRange / 4 // 简化估算：日内波动的1/4作为标准差

	// 计算布林带上下轨
	upperBand := middleBand + (estimatedStdDev * stdDevMultiplier)
	lowerBand := middleBand - (estimatedStdDev * stdDevMultiplier)

	// 计算价格相对于布林带的位置
	bandWidth := upperBand - lowerBand
	if bandWidth <= 0 {
		// 避免除零错误，使用默认持有信号
		signal.SignalType = models.SignalTypeHold
		signal.Strength = 0.5
		signal.Confidence = 0.5
		signal.Reason = "布林带宽度过小，无法判断"
		return signal, nil
	}

	pricePosition := (marketData.Close - lowerBand) / bandWidth

	// 计算相对波动率
	relativeVolatility := dailyRange / marketData.Close

	// 布林带策略逻辑：更宽松的条件以便产生交易信号
	if pricePosition < 0.3 && relativeVolatility > 0.01 { // 价格在下30%且有波动
		signal.SignalType = models.SignalTypeBuy
		signal.Side = models.TradeSideBuy
		signal.Strength = math.Min((0.3-pricePosition)*3, 1.0)
		signal.Confidence = 0.75
		signal.Reason = fmt.Sprintf("价格接近布林带下轨 (位置: %.1f%%, 中轨: %.2f, 下轨: %.2f)",
			pricePosition*100, middleBand, lowerBand)
	} else if pricePosition > 0.7 && relativeVolatility > 0.01 { // 价格在上30%且有波动
		signal.SignalType = models.SignalTypeSell
		signal.Side = models.TradeSideSell
		signal.Strength = math.Min((pricePosition-0.7)*3, 1.0)
		signal.Confidence = 0.75
		signal.Reason = fmt.Sprintf("价格接近布林带上轨 (位置: %.1f%%, 中轨: %.2f, 上轨: %.2f)",
			pricePosition*100, middleBand, upperBand)
	} else {
		signal.SignalType = models.SignalTypeHold
		signal.Strength = 0.5
		signal.Confidence = 0.5
		signal.Reason = fmt.Sprintf("价格在布林带中部 (位置: %.1f%%, 中轨: %.2f)",
			pricePosition*100, middleBand)
	}

	return signal, nil
}
