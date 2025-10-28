package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"

	"github.com/google/uuid"
)

// ParameterOptimizer 参数优化器
type ParameterOptimizer struct {
	backtestService *BacktestService
	strategyService *StrategyService
	logger          logger.Logger

	// 运行中的优化任务
	runningTasks map[string]*OptimizationTask
	tasksMutex   sync.RWMutex
}

// NewParameterOptimizer 创建参数优化器
func NewParameterOptimizer(backtestService *BacktestService, strategyService *StrategyService, log logger.Logger) *ParameterOptimizer {
	return &ParameterOptimizer{
		backtestService: backtestService,
		strategyService: strategyService,
		logger:          log,
		runningTasks:    make(map[string]*OptimizationTask),
	}
}

// OptimizationTask 优化任务
type OptimizationTask struct {
	ID                  string
	StrategyID          string
	Status              string // running, completed, failed, cancelled
	Progress            int    // 0-100
	CurrentCombo        int
	TotalCombos         int
	CurrentParams       map[string]interface{}
	BestParams          map[string]interface{}
	BestScore           float64
	BaselineParams      map[string]interface{} // 原始参数
	BaselinePerformance *models.BacktestResult // 原始参数的性能
	StartTime           time.Time
	EstimatedEndTime    time.Time
	CancelFunc          context.CancelFunc
	Results             []ParameterTestResult
}

// OptimizationConfig 优化配置
type OptimizationConfig struct {
	StrategyID         string                    `json:"strategy_id"`
	StrategyType       models.StrategyType       `json:"strategy_type"`
	ParameterRanges    map[string]ParameterRange `json:"parameter_ranges"`
	OptimizationTarget string                    `json:"optimization_target"` // total_return, sharpe_ratio, win_rate
	Symbols            []string                  `json:"symbols"`
	StartDate          string                    `json:"start_date"`
	EndDate            string                    `json:"end_date"`
	InitialCash        float64                   `json:"initial_cash"`
	Commission         float64                   `json:"commission"`
	Algorithm          string                    `json:"algorithm"` // grid_search, genetic
	MaxCombinations    int                       `json:"max_combinations"`
	GeneticConfig      *GeneticAlgorithmConfig   `json:"genetic_config,omitempty"`
}

// ParameterRange 参数范围
type ParameterRange struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Step float64 `json:"step"` // 用于网格搜索
}

// GeneticAlgorithmConfig 遗传算法配置
type GeneticAlgorithmConfig struct {
	PopulationSize int     `json:"population_size"`
	Generations    int     `json:"generations"`
	MutationRate   float64 `json:"mutation_rate"`
	CrossoverRate  float64 `json:"crossover_rate"`
	ElitismRate    float64 `json:"elitism_rate"`
}

// OptimizationResult 优化结果
type OptimizationResult struct {
	OptimizationID      string                 `json:"optimization_id"`
	StrategyID          string                 `json:"strategy_id"`
	BestParameters      map[string]interface{} `json:"best_parameters"`
	BestScore           float64                `json:"best_score"`
	Performance         *models.BacktestResult `json:"performance"`          // 优化后的性能（最佳参数）
	BaselinePerformance *models.BacktestResult `json:"baseline_performance"` // 优化前的性能（原始参数）
	BaselineParameters  map[string]interface{} `json:"baseline_parameters"`  // 原始参数
	AllResults          []ParameterTestResult  `json:"all_results"`
	TotalTested         int                    `json:"total_tested"`
	StartTime           time.Time              `json:"start_time"`
	EndTime             time.Time              `json:"end_time"`
	Duration            string                 `json:"duration"`
}

// ParameterTestResult 参数测试结果
type ParameterTestResult struct {
	Parameters  map[string]interface{} `json:"parameters"`
	Score       float64                `json:"score"`
	Performance *models.BacktestResult `json:"performance"`
}

// StartOptimization 启动参数优化
func (s *ParameterOptimizer) StartOptimization(ctx context.Context, config *OptimizationConfig) (string, error) {
	optimizationID := uuid.New().String()

	s.logger.Info("启动参数优化",
		logger.String("optimization_id", optimizationID),
		logger.String("strategy_id", config.StrategyID),
		logger.String("algorithm", config.Algorithm),
	)

	// 创建独立的context，不依赖传入的ctx
	// 这样即使HTTP请求结束，优化任务仍然可以继续运行
	optimizationCtx, cancel := context.WithCancel(context.Background())

	// 创建任务
	task := &OptimizationTask{
		ID:         optimizationID,
		StrategyID: config.StrategyID,
		Status:     "running",
		Progress:   0,
		StartTime:  time.Now(),
		CancelFunc: cancel,
		Results:    make([]ParameterTestResult, 0),
	}

	s.tasksMutex.Lock()
	s.runningTasks[optimizationID] = task
	s.tasksMutex.Unlock()

	// 异步执行优化
	go func() {
		var result *OptimizationResult
		var err error

		switch config.Algorithm {
		case "grid_search", "":
			result, err = s.gridSearchOptimization(optimizationCtx, task, config)
		case "genetic":
			result, err = s.geneticAlgorithmOptimization(optimizationCtx, task, config)
		default:
			err = fmt.Errorf("不支持的优化算法: %s", config.Algorithm)
		}

		s.tasksMutex.Lock()
		if err != nil {
			task.Status = "failed"
			s.logger.Error("参数优化失败", logger.ErrorField(err))
		} else {
			task.Status = "completed"
			task.BestParams = result.BestParameters
			task.BestScore = result.BestScore
			task.Results = result.AllResults
			s.logger.Info("参数优化完成",
				logger.String("optimization_id", optimizationID),
				logger.Float64("best_score", result.BestScore),
			)
		}
		s.tasksMutex.Unlock()
	}()

	return optimizationID, nil
}

// gridSearchOptimization 网格搜索优化
func (s *ParameterOptimizer) gridSearchOptimization(ctx context.Context, task *OptimizationTask, config *OptimizationConfig) (*OptimizationResult, error) {
	startTime := time.Now()

	// 🔧 新增：获取原始策略并测试baseline性能
	originalStrategy, err := s.strategyService.GetStrategy(ctx, config.StrategyID)
	if err == nil && originalStrategy != nil {
		s.logger.Info("⏳ 测试原始参数性能作为baseline",
			logger.String("strategy_id", config.StrategyID),
		)

		baselineResult := s.testParameters(ctx, config, originalStrategy.Parameters)
		task.BaselineParams = originalStrategy.Parameters
		task.BaselinePerformance = baselineResult.Performance

		s.logger.Info("✅ Baseline测试完成",
			logger.String("strategy_id", config.StrategyID),
			logger.Float64("baseline_score", baselineResult.Score),
		)
	} else {
		s.logger.Warn("无法获取原始策略，跳过baseline测试",
			logger.String("strategy_id", config.StrategyID),
		)
	}

	// 生成参数组合
	parameterCombinations := s.generateParameterCombinations(config.ParameterRanges)

	// 限制组合数量
	if config.MaxCombinations > 0 && len(parameterCombinations) > config.MaxCombinations {
		parameterCombinations = parameterCombinations[:config.MaxCombinations]
	}

	task.TotalCombos = len(parameterCombinations)
	s.logger.Info("生成参数组合完成",
		logger.Int("total_combinations", len(parameterCombinations)),
	)

	// 并行测试每组参数
	results := make([]ParameterTestResult, 0, len(parameterCombinations))
	resultsChan := make(chan ParameterTestResult, len(parameterCombinations))
	semaphore := make(chan struct{}, runtime.NumCPU()) // 限制并发数

	var wg sync.WaitGroup
	for i, params := range parameterCombinations {
		select {
		case <-ctx.Done():
			s.logger.Info("优化任务被取消")
			return nil, errors.New("优化任务被取消")
		default:
		}

		wg.Add(1)
		go func(idx int, parameters map[string]interface{}) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 更新当前测试参数
			s.tasksMutex.Lock()
			task.CurrentCombo = idx + 1
			task.CurrentParams = parameters
			task.Progress = int(float64(idx+1) / float64(task.TotalCombos) * 100)
			s.tasksMutex.Unlock()

			// 测试这组参数
			result := s.testParameters(ctx, config, parameters)
			result.Parameters = parameters
			resultsChan <- result

			// 更新最佳结果
			s.tasksMutex.Lock()
			if result.Score > task.BestScore {
				task.BestScore = result.Score
				task.BestParams = parameters
			}
			s.tasksMutex.Unlock()

			s.logger.Debug("参数测试完成",
				logger.Int("index", idx),
				logger.Float64("score", result.Score),
			)
		}(i, params)
	}

	// 等待所有测试完成
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 收集结果
	for result := range resultsChan {
		results = append(results, result)
	}

	// 按得分排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	duration := time.Since(startTime)

	if len(results) == 0 {
		return nil, errors.New("没有获得任何测试结果")
	}

	return &OptimizationResult{
		OptimizationID: task.ID,
		StrategyID:     config.StrategyID,
		BestParameters: results[0].Parameters,
		BestScore:      results[0].Score,
		Performance:    results[0].Performance,
		AllResults:     results,
		TotalTested:    len(results),
		StartTime:      startTime,
		EndTime:        time.Now(),
		Duration:       duration.String(),
	}, nil
}

// geneticAlgorithmOptimization 遗传算法优化（简化版）
func (s *ParameterOptimizer) geneticAlgorithmOptimization(ctx context.Context, task *OptimizationTask, config *OptimizationConfig) (*OptimizationResult, error) {
	startTime := time.Now()

	if config.GeneticConfig == nil {
		return nil, errors.New("遗传算法配置为空")
	}

	// 🔧 新增：获取原始策略并测试baseline性能
	originalStrategy, err := s.strategyService.GetStrategy(ctx, config.StrategyID)
	if err == nil && originalStrategy != nil {
		s.logger.Info("⏳ 测试原始参数性能作为baseline",
			logger.String("strategy_id", config.StrategyID),
		)

		baselineResult := s.testParameters(ctx, config, originalStrategy.Parameters)
		task.BaselineParams = originalStrategy.Parameters
		task.BaselinePerformance = baselineResult.Performance

		s.logger.Info("✅ Baseline测试完成",
			logger.String("strategy_id", config.StrategyID),
			logger.Float64("baseline_score", baselineResult.Score),
		)
	} else {
		s.logger.Warn("无法获取原始策略，跳过baseline测试",
			logger.String("strategy_id", config.StrategyID),
		)
	}

	ga := config.GeneticConfig
	population := s.initializePopulation(config.ParameterRanges, ga.PopulationSize)
	task.TotalCombos = ga.PopulationSize * ga.Generations

	var bestIndividual map[string]interface{}
	var bestScore float64
	allResults := make([]ParameterTestResult, 0)

	for gen := 0; gen < ga.Generations; gen++ {
		select {
		case <-ctx.Done():
			return nil, errors.New("优化任务被取消")
		default:
		}

		// 评估当前种群
		fitness := make([]float64, len(population))
		for i, individual := range population {
			result := s.testParameters(ctx, config, individual)
			fitness[i] = result.Score
			allResults = append(allResults, result)

			if result.Score > bestScore {
				bestScore = result.Score
				bestIndividual = individual
			}

			// 更新进度
			s.tasksMutex.Lock()
			task.CurrentCombo = gen*ga.PopulationSize + i + 1
			task.Progress = int(float64(task.CurrentCombo) / float64(task.TotalCombos) * 100)
			task.BestScore = bestScore
			task.BestParams = bestIndividual
			s.tasksMutex.Unlock()
		}

		// 选择、交叉、变异
		newPopulation := s.evolvePopulation(population, fitness, config.ParameterRanges, ga)
		population = newPopulation

		s.logger.Info("遗传算法代数完成",
			logger.Int("generation", gen+1),
			logger.Float64("best_score", bestScore),
		)
	}

	// 按得分排序结果
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].Score > allResults[j].Score
	})

	duration := time.Since(startTime)

	return &OptimizationResult{
		OptimizationID: task.ID,
		StrategyID:     config.StrategyID,
		BestParameters: bestIndividual,
		BestScore:      bestScore,
		Performance:    allResults[0].Performance,
		AllResults:     allResults[:min(len(allResults), 100)], // 只保留前100个结果
		TotalTested:    len(allResults),
		StartTime:      startTime,
		EndTime:        time.Now(),
		Duration:       duration.String(),
	}, nil
}

// testParameters 测试一组参数
func (s *ParameterOptimizer) testParameters(ctx context.Context, config *OptimizationConfig, parameters map[string]interface{}) ParameterTestResult {
	// 创建临时策略
	strategy := &models.Strategy{
		ID:         fmt.Sprintf("temp_opt_strategy_%d", time.Now().UnixNano()),
		Name:       "临时优化策略",
		Type:       config.StrategyType,
		Parameters: parameters,
		Status:     models.StrategyStatusInactive,
	}

	// 运行快速回测
	performance := s.runQuickBacktest(ctx, config, strategy)

	// 根据优化目标计算得分
	score := s.calculateScore(performance, config.OptimizationTarget)

	return ParameterTestResult{
		Parameters:  parameters,
		Score:       score,
		Performance: performance,
	}
}

// runQuickBacktest 运行快速回测
func (s *ParameterOptimizer) runQuickBacktest(ctx context.Context, config *OptimizationConfig, strategy *models.Strategy) *models.BacktestResult {
	// 简化的回测执行（这里需要调用实际的回测服务）
	// 为了简化，这里返回模拟结果
	// 实际实现时应该调用 backtestService.RunBacktest

	// TODO: 实际调用回测服务
	// 可以参考 BacktestService.RunBacktest 的实现

	// 临时返回模拟数据
	return &models.BacktestResult{
		ID:             uuid.New().String(),
		BacktestID:     uuid.New().String(),
		StrategyID:     strategy.ID,
		TotalReturn:    rand.Float64()*0.4 - 0.2, // -20% to +20%
		AnnualReturn:   rand.Float64()*0.4 - 0.2,
		MaxDrawdown:    -rand.Float64() * 0.3, // 0 to -30%
		SharpeRatio:    rand.Float64() * 3,
		WinRate:        0.4 + rand.Float64()*0.3, // 40% to 70%
		TotalTrades:    50 + rand.Intn(150),
		ProfitFactor:   1 + rand.Float64()*2,
		AvgTradeReturn: rand.Float64()*0.1 - 0.05,
		CreatedAt:      time.Now(),
	}
}

// calculateScore 根据优化目标计算得分
func (s *ParameterOptimizer) calculateScore(result *models.BacktestResult, target string) float64 {
	if result == nil {
		return -math.MaxFloat64
	}

	switch target {
	case "total_return":
		return result.TotalReturn
	case "sharpe_ratio":
		return result.SharpeRatio
	case "win_rate":
		return result.WinRate
	case "profit_factor":
		return result.ProfitFactor
	default:
		// 综合得分：夏普比率为主，收益率和最大回撤为辅
		return result.SharpeRatio*0.5 + result.TotalReturn*0.3 + (1+result.MaxDrawdown)*0.2
	}
}

// generateParameterCombinations 生成所有参数组合（网格搜索）
func (s *ParameterOptimizer) generateParameterCombinations(ranges map[string]ParameterRange) []map[string]interface{} {
	// 获取参数名称列表
	paramNames := make([]string, 0, len(ranges))
	for name := range ranges {
		paramNames = append(paramNames, name)
	}
	sort.Strings(paramNames) // 确保顺序一致

	// 生成每个参数的可能值
	paramValues := make([][]float64, len(paramNames))
	for i, name := range paramNames {
		r := ranges[name]
		values := make([]float64, 0)
		for v := r.Min; v <= r.Max; v += r.Step {
			values = append(values, v)
		}
		if len(values) == 0 {
			values = append(values, r.Min)
		}
		paramValues[i] = values
	}

	// 生成笛卡尔积
	combinations := s.cartesianProduct(paramNames, paramValues)
	return combinations
}

// cartesianProduct 生成笛卡尔积
func (s *ParameterOptimizer) cartesianProduct(names []string, values [][]float64) []map[string]interface{} {
	if len(names) == 0 {
		return []map[string]interface{}{{}}
	}

	result := []map[string]interface{}{}
	subResult := s.cartesianProduct(names[1:], values[1:])

	for _, val := range values[0] {
		for _, sub := range subResult {
			combo := make(map[string]interface{})
			combo[names[0]] = val
			for k, v := range sub {
				combo[k] = v
			}
			result = append(result, combo)
		}
	}

	return result
}

// initializePopulation 初始化遗传算法种群
func (s *ParameterOptimizer) initializePopulation(ranges map[string]ParameterRange, size int) []map[string]interface{} {
	population := make([]map[string]interface{}, size)
	for i := 0; i < size; i++ {
		individual := make(map[string]interface{})
		for name, r := range ranges {
			individual[name] = r.Min + rand.Float64()*(r.Max-r.Min)
		}
		population[i] = individual
	}
	return population
}

// evolvePopulation 进化种群
func (s *ParameterOptimizer) evolvePopulation(population []map[string]interface{}, fitness []float64, ranges map[string]ParameterRange, ga *GeneticAlgorithmConfig) []map[string]interface{} {
	newPopulation := make([]map[string]interface{}, len(population))

	// 精英保留
	eliteCount := int(float64(len(population)) * ga.ElitismRate)
	eliteIndices := s.selectTopIndices(fitness, eliteCount)
	for i, idx := range eliteIndices {
		newPopulation[i] = population[idx]
	}

	// 生成新个体
	for i := eliteCount; i < len(population); i++ {
		// 选择父代
		parent1 := population[s.tournamentSelection(fitness)]
		parent2 := population[s.tournamentSelection(fitness)]

		// 交叉
		child := s.crossover(parent1, parent2, ga.CrossoverRate)

		// 变异
		child = s.mutate(child, ranges, ga.MutationRate)

		newPopulation[i] = child
	}

	return newPopulation
}

// selectTopIndices 选择得分最高的个体索引
func (s *ParameterOptimizer) selectTopIndices(fitness []float64, count int) []int {
	type indexedFitness struct {
		index   int
		fitness float64
	}

	indexed := make([]indexedFitness, len(fitness))
	for i, f := range fitness {
		indexed[i] = indexedFitness{index: i, fitness: f}
	}

	sort.Slice(indexed, func(i, j int) bool {
		return indexed[i].fitness > indexed[j].fitness
	})

	indices := make([]int, count)
	for i := 0; i < count; i++ {
		indices[i] = indexed[i].index
	}

	return indices
}

// tournamentSelection 锦标赛选择
func (s *ParameterOptimizer) tournamentSelection(fitness []float64) int {
	tournamentSize := 3
	best := rand.Intn(len(fitness))
	for i := 1; i < tournamentSize; i++ {
		competitor := rand.Intn(len(fitness))
		if fitness[competitor] > fitness[best] {
			best = competitor
		}
	}
	return best
}

// crossover 交叉操作
func (s *ParameterOptimizer) crossover(parent1, parent2 map[string]interface{}, rate float64) map[string]interface{} {
	child := make(map[string]interface{})
	for key := range parent1 {
		if rand.Float64() < rate {
			child[key] = parent1[key]
		} else {
			child[key] = parent2[key]
		}
	}
	return child
}

// mutate 变异操作
func (s *ParameterOptimizer) mutate(individual map[string]interface{}, ranges map[string]ParameterRange, rate float64) map[string]interface{} {
	mutated := make(map[string]interface{})
	for key, val := range individual {
		if rand.Float64() < rate {
			r := ranges[key]
			mutated[key] = r.Min + rand.Float64()*(r.Max-r.Min)
		} else {
			mutated[key] = val
		}
	}
	return mutated
}

// GetOptimizationProgress 获取优化进度
func (s *ParameterOptimizer) GetOptimizationProgress(optimizationID string) (*OptimizationTask, error) {
	s.tasksMutex.RLock()
	defer s.tasksMutex.RUnlock()

	task, exists := s.runningTasks[optimizationID]
	if !exists {
		return nil, fmt.Errorf("优化任务不存在: %s", optimizationID)
	}

	return task, nil
}

// GetOptimizationResult 获取优化结果
func (s *ParameterOptimizer) GetOptimizationResult(optimizationID string) (*OptimizationResult, error) {
	s.tasksMutex.RLock()
	task, exists := s.runningTasks[optimizationID]
	s.tasksMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("优化任务不存在: %s", optimizationID)
	}

	if task.Status != "completed" {
		return nil, fmt.Errorf("优化任务尚未完成，当前状态: %s", task.Status)
	}

	// 🔧 修复：从最佳结果中获取Performance数据
	var bestPerformance *models.BacktestResult
	for _, result := range task.Results {
		if result.Score == task.BestScore {
			bestPerformance = result.Performance
			break
		}
	}

	// 如果没找到匹配的，使用第一个结果的Performance（兜底）
	if bestPerformance == nil && len(task.Results) > 0 {
		bestPerformance = task.Results[0].Performance
		s.logger.Warn("未找到最佳得分对应的Performance，使用第一个结果",
			logger.String("optimization_id", optimizationID),
			logger.Float64("best_score", task.BestScore),
		)
	}

	s.logger.Info("📊 优化结果准备完成",
		logger.String("optimization_id", optimizationID),
		logger.String("strategy_id", task.StrategyID),
		logger.Int("total_tested", len(task.Results)),
		logger.Float64("best_score", task.BestScore),
		logger.Bool("has_performance", bestPerformance != nil),
		logger.Bool("has_baseline", task.BaselinePerformance != nil),
	)

	return &OptimizationResult{
		OptimizationID:      task.ID,
		StrategyID:          task.StrategyID,
		BestParameters:      task.BestParams,
		BestScore:           task.BestScore,
		Performance:         bestPerformance,          // ✅ 优化后的性能
		BaselinePerformance: task.BaselinePerformance, // ✅ 优化前的性能
		BaselineParameters:  task.BaselineParams,      // ✅ 原始参数
		AllResults:          task.Results,
		TotalTested:         len(task.Results),
		StartTime:           task.StartTime,
		EndTime:             time.Now(),
	}, nil
}

// CancelOptimization 取消优化任务
func (s *ParameterOptimizer) CancelOptimization(optimizationID string) error {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()

	task, exists := s.runningTasks[optimizationID]
	if !exists {
		return fmt.Errorf("优化任务不存在: %s", optimizationID)
	}

	if task.Status != "running" {
		return fmt.Errorf("任务不在运行状态: %s", task.Status)
	}

	task.CancelFunc()
	task.Status = "cancelled"

	s.logger.Info("优化任务已取消", logger.String("optimization_id", optimizationID))
	return nil
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
