package service

import (
	"context"
	"testing"
	"time"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"
)

// TestParameterOptimizerCreation 测试参数优化器的创建
func TestParameterOptimizerCreation(t *testing.T) {
	log := logger.GetGlobalLogger()
	strategyService := NewStrategyService(log)
	backtestService := NewBacktestService(strategyService, nil, nil, log)

	optimizer := NewParameterOptimizer(backtestService, strategyService, log)

	if optimizer == nil {
		t.Fatal("参数优化器创建失败")
	}

	if optimizer.runningTasks == nil {
		t.Error("运行任务映射未初始化")
	}
}

// TestParameterCombinationsGeneration 测试参数组合生成
func TestParameterCombinationsGeneration(t *testing.T) {
	log := logger.GetGlobalLogger()
	strategyService := NewStrategyService(log)
	backtestService := NewBacktestService(strategyService, nil, nil, log)
	optimizer := NewParameterOptimizer(backtestService, strategyService, log)

	ranges := map[string]ParameterRange{
		"fast_period": {
			Min:  10,
			Max:  14,
			Step: 2,
		},
		"slow_period": {
			Min:  20,
			Max:  24,
			Step: 2,
		},
	}

	combinations := optimizer.generateParameterCombinations(ranges)

	// 应该生成 3 * 3 = 9 个组合
	// fast_period: 10, 12, 14
	// slow_period: 20, 22, 24
	expectedCount := 9
	if len(combinations) != expectedCount {
		t.Errorf("组合数量错误: 期望 %d, 实际 %d", expectedCount, len(combinations))
	}

	// 验证第一个组合
	if combinations[0]["fast_period"] != 10.0 || combinations[0]["slow_period"] != 20.0 {
		t.Errorf("第一个组合值错误: %+v", combinations[0])
	}

	// 验证最后一个组合
	lastIdx := len(combinations) - 1
	if combinations[lastIdx]["fast_period"] != 14.0 || combinations[lastIdx]["slow_period"] != 24.0 {
		t.Errorf("最后一个组合值错误: %+v", combinations[lastIdx])
	}
}

// TestScoreCalculation 测试得分计算
func TestScoreCalculation(t *testing.T) {
	log := logger.GetGlobalLogger()
	strategyService := NewStrategyService(log)
	backtestService := NewBacktestService(strategyService, nil, nil, log)
	optimizer := NewParameterOptimizer(backtestService, strategyService, log)

	result := &models.BacktestResult{
		TotalReturn: 0.25,
		SharpeRatio: 1.5,
		MaxDrawdown: -0.15,
		WinRate:     0.60,
	}

	tests := []struct {
		target        string
		expectedScore float64
	}{
		{"total_return", 0.25},
		{"sharpe_ratio", 1.5},
		{"win_rate", 0.60},
	}

	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			score := optimizer.calculateScore(result, tt.target)
			if score != tt.expectedScore {
				t.Errorf("目标 %s 得分错误: 期望 %.2f, 实际 %.2f", tt.target, tt.expectedScore, score)
			}
		})
	}

	// 测试综合得分（默认情况）
	compositeScore := optimizer.calculateScore(result, "composite")
	if compositeScore <= 0 {
		t.Error("综合得分应该大于0")
	}
}

// TestOptimizationTaskLifecycle 测试优化任务生命周期
func TestOptimizationTaskLifecycle(t *testing.T) {
	log := logger.GetGlobalLogger()
	strategyService := NewStrategyService(log)
	backtestService := NewBacktestService(strategyService, nil, nil, log)
	optimizer := NewParameterOptimizer(backtestService, strategyService, log)

	config := &OptimizationConfig{
		StrategyID:   "test_strategy",
		StrategyType: models.StrategyTypeTechnical,
		ParameterRanges: map[string]ParameterRange{
			"period": {Min: 5, Max: 15, Step: 5},
		},
		OptimizationTarget: "sharpe_ratio",
		Symbols:            []string{"000001.SZ"},
		StartDate:          "2024-01-01",
		EndDate:            "2024-12-31",
		InitialCash:        1000000,
		Commission:         0.0003,
		Algorithm:          "grid_search",
		MaxCombinations:    10,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动优化
	optimizationID, err := optimizer.StartOptimization(ctx, config)
	if err != nil {
		t.Fatalf("启动优化失败: %v", err)
	}

	if optimizationID == "" {
		t.Error("优化ID为空")
	}

	// 等待一小段时间让优化开始
	time.Sleep(100 * time.Millisecond)

	// 检查任务状态
	task, err := optimizer.GetOptimizationProgress(optimizationID)
	if err != nil {
		t.Fatalf("获取优化进度失败: %v", err)
	}

	if task.Status != "running" && task.Status != "completed" {
		t.Errorf("任务状态错误: %s", task.Status)
	}

	// 测试取消优化
	err = optimizer.CancelOptimization(optimizationID)
	if err != nil && task.Status == "running" {
		t.Errorf("取消优化失败: %v", err)
	}

	// 等待任务状态更新
	time.Sleep(100 * time.Millisecond)

	// 再次检查状态
	task, _ = optimizer.GetOptimizationProgress(optimizationID)
	if task.Status != "cancelled" && task.Status != "completed" {
		t.Logf("任务状态: %s (可能已完成或被取消)", task.Status)
	}
}

// TestGeneticAlgorithmComponents 测试遗传算法组件
func TestGeneticAlgorithmComponents(t *testing.T) {
	log := logger.GetGlobalLogger()
	strategyService := NewStrategyService(log)
	backtestService := NewBacktestService(strategyService, nil, nil, log)
	optimizer := NewParameterOptimizer(backtestService, strategyService, log)

	ranges := map[string]ParameterRange{
		"param1": {Min: 1, Max: 10, Step: 1},
		"param2": {Min: 0, Max: 1, Step: 0.1},
	}

	// 测试种群初始化
	population := optimizer.initializePopulation(ranges, 10)
	if len(population) != 10 {
		t.Errorf("种群大小错误: 期望 10, 实际 %d", len(population))
	}

	// 验证每个个体的参数在范围内
	for i, individual := range population {
		for paramName, paramValue := range individual {
			r := ranges[paramName]
			val := paramValue.(float64)
			if val < r.Min || val > r.Max {
				t.Errorf("个体 %d 的参数 %s 超出范围: %.2f (范围: %.2f-%.2f)",
					i, paramName, val, r.Min, r.Max)
			}
		}
	}

	// 测试交叉操作
	parent1 := population[0]
	parent2 := population[1]
	child := optimizer.crossover(parent1, parent2, 0.5)

	if len(child) != len(parent1) {
		t.Error("交叉后的子代参数数量不正确")
	}

	// 测试变异操作
	mutated := optimizer.mutate(child, ranges, 1.0) // 100%变异率
	if len(mutated) != len(child) {
		t.Error("变异后的个体参数数量不正确")
	}

	// 验证变异后的参数在范围内
	for paramName, paramValue := range mutated {
		r := ranges[paramName]
		val := paramValue.(float64)
		if val < r.Min || val > r.Max {
			t.Errorf("变异后的参数 %s 超出范围: %.2f (范围: %.2f-%.2f)",
				paramName, val, r.Min, r.Max)
		}
	}
}

// BenchmarkParameterCombinationsGeneration 基准测试参数组合生成
func BenchmarkParameterCombinationsGeneration(b *testing.B) {
	log := logger.GetGlobalLogger()
	strategyService := NewStrategyService(log)
	backtestService := NewBacktestService(strategyService, nil, nil, log)
	optimizer := NewParameterOptimizer(backtestService, strategyService, log)

	ranges := map[string]ParameterRange{
		"fast_period":   {Min: 5, Max: 20, Step: 1},
		"slow_period":   {Min: 20, Max: 50, Step: 1},
		"signal_period": {Min: 5, Max: 15, Step: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = optimizer.generateParameterCombinations(ranges)
	}
}
