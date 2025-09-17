package service

import (
	"context"
	"testing"
	"time"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"
)

func TestMultiStrategyBacktestResultsDisplay(t *testing.T) {
	// 创建测试用的日志配置
	logConfig := &logger.Config{
		Level:  "info",
		Format: "console",
	}
	log, _ := logger.NewLogger(logConfig)

	// 创建测试用的服务依赖（简化版本，避免依赖问题）
	dailyCacheService := NewDailyCacheService(nil) // 使用默认配置

	// 创建回测服务（传入nil的服务依赖进行测试）
	backtestService := NewBacktestService(nil, nil, dailyCacheService, log)

	// 创建测试回测
	backtest := &models.Backtest{
		ID:          "test-multi-strategy",
		Name:        "多策略回测测试",
		StrategyIDs: []string{"strategy1", "strategy2", "strategy3"},
		Symbols:     []string{"600976", "000001"},
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
		InitialCash: 1000000,
		Commission:  0.001,
		Status:      models.BacktestStatusCompleted,
		CreatedAt:   time.Now(),
	}

	// 创建测试用的多策略结果
	multiResults := []models.BacktestResult{
		{
			ID:             "result1",
			BacktestID:     backtest.ID,
			StrategyID:     "strategy1",
			StrategyName:   "MA交叉策略",
			TotalReturn:    0.15,  // 15%
			AnnualReturn:   0.20,  // 20%
			MaxDrawdown:    -0.08, // -8%
			SharpeRatio:    1.5,
			WinRate:        0.65, // 65%
			TotalTrades:    45,
			AvgTradeReturn: 0.003, // 0.3%
			ProfitFactor:   2.1,
			CreatedAt:      time.Now(),
		},
		{
			ID:             "result2",
			BacktestID:     backtest.ID,
			StrategyID:     "strategy2",
			StrategyName:   "RSI策略",
			TotalReturn:    0.12,  // 12%
			AnnualReturn:   0.16,  // 16%
			MaxDrawdown:    -0.06, // -6%
			SharpeRatio:    1.3,
			WinRate:        0.70, // 70%
			TotalTrades:    38,
			AvgTradeReturn: 0.0032, // 0.32%
			ProfitFactor:   2.3,
			CreatedAt:      time.Now(),
		},
		{
			ID:             "result3",
			BacktestID:     backtest.ID,
			StrategyID:     "strategy3",
			StrategyName:   "MACD策略",
			TotalReturn:    0.18,  // 18%
			AnnualReturn:   0.24,  // 24%
			MaxDrawdown:    -0.10, // -10%
			SharpeRatio:    1.7,
			WinRate:        0.60, // 60%
			TotalTrades:    52,
			AvgTradeReturn: 0.0035, // 0.35%
			ProfitFactor:   1.9,
			CreatedAt:      time.Now(),
		},
	}

	// 保存测试数据到服务
	backtestService.backtests[backtest.ID] = backtest
	backtestService.backtestMultiResults[backtest.ID] = multiResults

	// 计算组合指标
	combinedMetrics := backtestService.calculateCombinedMetrics(multiResults)

	// 保存组合指标
	if combinedMetrics != nil {
		backtestService.backtestResults[backtest.ID] = combinedMetrics
	}

	// 创建测试交易记录
	testTrades := []models.Trade{
		{
			ID:         "trade1",
			BacktestID: backtest.ID,
			StrategyID: "strategy1",
			Symbol:     "600976",
			Side:       models.TradeSideBuy,
			Quantity:   1000,
			Price:      10.50,
			Commission: 10.50,
			Timestamp:  time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC),
			CreatedAt:  time.Now(),
		},
		{
			ID:         "trade2",
			BacktestID: backtest.ID,
			StrategyID: "strategy2",
			Symbol:     "000001",
			Side:       models.TradeSideBuy,
			Quantity:   500,
			Price:      15.20,
			Commission: 7.60,
			Timestamp:  time.Date(2024, 1, 20, 10, 0, 0, 0, time.UTC),
			CreatedAt:  time.Now(),
		},
	}
	backtestService.backtestTrades[backtest.ID] = testTrades

	// 创建测试权益曲线
	equityCurve := []models.EquityPoint{
		{
			Date:           "2024-01-01",
			PortfolioValue: 1000000,
			Cash:           1000000,
			Holdings:       0,
		},
		{
			Date:           "2024-02-01",
			PortfolioValue: 1080000,
			Cash:           800000,
			Holdings:       280000,
		},
		{
			Date:           "2024-03-31",
			PortfolioValue: 1150000,
			Cash:           750000,
			Holdings:       400000,
		},
	}
	backtestService.backtestEquityCurves[backtest.ID] = equityCurve

	// 测试获取多策略回测结果
	ctx := context.Background()
	results, err := backtestService.GetBacktestResults(ctx, backtest.ID)
	if err != nil {
		t.Fatalf("获取多策略回测结果失败: %v", err)
	}

	// 验证结果结构
	t.Logf("✅ 多策略回测结果测试开始")
	t.Logf("回测ID: %s", results.BacktestID)
	t.Logf("策略数量: %d", len(results.Strategies))
	t.Logf("性能结果数量: %d", len(results.Performance))

	// 验证多策略结果
	if len(results.Performance) != 3 {
		t.Errorf("期望3个策略结果，实际获得%d个", len(results.Performance))
	}

	// 验证组合指标
	if results.CombinedMetrics == nil {
		t.Errorf("组合指标不应为空")
	} else {
		t.Logf("✅ 组合指标存在")
		t.Logf("组合总收益率: %.2f%%", results.CombinedMetrics.TotalReturn*100)
		t.Logf("组合年化收益率: %.2f%%", results.CombinedMetrics.AnnualReturn*100)
		t.Logf("组合最大回撤: %.2f%%", results.CombinedMetrics.MaxDrawdown*100)
		t.Logf("组合夏普比率: %.2f", results.CombinedMetrics.SharpeRatio)
		t.Logf("组合胜率: %.2f%%", results.CombinedMetrics.WinRate*100)
		t.Logf("组合总交易次数: %d", results.CombinedMetrics.TotalTrades)
		t.Logf("组合平均交易收益: %.3f%%", results.CombinedMetrics.AvgTradeReturn*100)
		t.Logf("组合盈亏比: %.2f", results.CombinedMetrics.ProfitFactor)

		// 验证组合指标的合理性
		if results.CombinedMetrics.TotalReturn == 0 {
			t.Errorf("组合总收益率不应为0")
		}
		if results.CombinedMetrics.SharpeRatio == 0 {
			t.Errorf("组合夏普比率不应为0")
		}
	}

	// 验证各策略的详细结果
	for i, performance := range results.Performance {
		t.Logf("策略%d (%s):", i+1, performance.StrategyName)
		t.Logf("  总收益率: %.2f%%", performance.TotalReturn*100)
		t.Logf("  年化收益率: %.2f%%", performance.AnnualReturn*100)
		t.Logf("  最大回撤: %.2f%%", performance.MaxDrawdown*100)
		t.Logf("  夏普比率: %.2f", performance.SharpeRatio)
		t.Logf("  胜率: %.2f%%", performance.WinRate*100)
		t.Logf("  交易次数: %d", performance.TotalTrades)

		// 验证数据有效性
		if performance.StrategyID == "" {
			t.Errorf("策略ID不应为空")
		}
		if performance.StrategyName == "" {
			t.Errorf("策略名称不应为空")
		}
	}

	// 验证权益曲线
	if len(results.EquityCurve) == 0 {
		t.Errorf("权益曲线不应为空")
	} else {
		t.Logf("✅ 权益曲线包含%d个数据点", len(results.EquityCurve))
	}

	// 验证交易记录
	if len(results.Trades) == 0 {
		t.Errorf("交易记录不应为空")
	} else {
		t.Logf("✅ 交易记录包含%d条记录", len(results.Trades))
	}

	// 验证策略信息
	if len(results.Strategies) != 3 {
		t.Errorf("期望3个策略信息，实际获得%d个", len(results.Strategies))
	} else {
		t.Logf("✅ 策略信息完整")
		for i, strategy := range results.Strategies {
			t.Logf("策略%d: %s (%s)", i+1, strategy.Name, strategy.ID)
		}
	}

	t.Logf("🎉 多策略回测结果测试完成")
}

func TestCalculateCombinedMetrics(t *testing.T) {
	// 创建测试用的日志配置
	logConfig := &logger.Config{
		Level:  "info",
		Format: "console",
	}
	log, _ := logger.NewLogger(logConfig)

	// 创建回测服务
	backtestService := NewBacktestService(nil, nil, nil, log)

	// 创建测试数据
	results := []models.BacktestResult{
		{
			TotalReturn:    0.15,
			AnnualReturn:   0.20,
			MaxDrawdown:    -0.08,
			SharpeRatio:    1.5,
			WinRate:        0.65,
			TotalTrades:    45,
			AvgTradeReturn: 0.003,
			ProfitFactor:   2.1,
		},
		{
			TotalReturn:    0.12,
			AnnualReturn:   0.16,
			MaxDrawdown:    -0.06,
			SharpeRatio:    1.3,
			WinRate:        0.70,
			TotalTrades:    38,
			AvgTradeReturn: 0.0032,
			ProfitFactor:   2.3,
		},
		{
			TotalReturn:    0.18,
			AnnualReturn:   0.24,
			MaxDrawdown:    -0.10,
			SharpeRatio:    1.7,
			WinRate:        0.60,
			TotalTrades:    52,
			AvgTradeReturn: 0.0035,
			ProfitFactor:   1.9,
		},
	}

	// 计算组合指标
	combined := backtestService.calculateCombinedMetrics(results)

	if combined == nil {
		t.Fatalf("组合指标计算结果不应为空")
	}

	// 验证组合指标计算的正确性
	expectedTotalReturn := (0.15 + 0.12 + 0.18) / 3
	if abs(combined.TotalReturn-expectedTotalReturn) > 0.001 {
		t.Errorf("组合总收益率计算错误，期望%.3f，实际%.3f", expectedTotalReturn, combined.TotalReturn)
	}

	expectedAnnualReturn := (0.20 + 0.16 + 0.24) / 3
	if abs(combined.AnnualReturn-expectedAnnualReturn) > 0.001 {
		t.Errorf("组合年化收益率计算错误，期望%.3f，实际%.3f", expectedAnnualReturn, combined.AnnualReturn)
	}

	expectedMaxDrawdown := (-0.08 + -0.06 + -0.10) / 3
	if abs(combined.MaxDrawdown-expectedMaxDrawdown) > 0.001 {
		t.Errorf("组合最大回撤计算错误，期望%.3f，实际%.3f", expectedMaxDrawdown, combined.MaxDrawdown)
	}

	expectedSharpeRatio := (1.5 + 1.3 + 1.7) / 3
	if abs(combined.SharpeRatio-expectedSharpeRatio) > 0.001 {
		t.Errorf("组合夏普比率计算错误，期望%.3f，实际%.3f", expectedSharpeRatio, combined.SharpeRatio)
	}

	expectedTotalTrades := 45 + 38 + 52
	if combined.TotalTrades != expectedTotalTrades {
		t.Errorf("组合总交易次数计算错误，期望%d，实际%d", expectedTotalTrades, combined.TotalTrades)
	}

	t.Logf("✅ 组合指标计算验证通过")
	t.Logf("组合总收益率: %.2f%%", combined.TotalReturn*100)
	t.Logf("组合年化收益率: %.2f%%", combined.AnnualReturn*100)
	t.Logf("组合最大回撤: %.2f%%", combined.MaxDrawdown*100)
	t.Logf("组合夏普比率: %.2f", combined.SharpeRatio)
	t.Logf("组合胜率: %.2f%%", combined.WinRate*100)
	t.Logf("组合总交易次数: %d", combined.TotalTrades)
	t.Logf("组合盈亏比: %.2f", combined.ProfitFactor)
}

// 辅助函数：计算绝对值
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
