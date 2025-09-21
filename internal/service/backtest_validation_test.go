package service

import (
	"context"
	"testing"
	"time"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"
)

func TestValidateTradesData(t *testing.T) {
	// 创建测试用的日志配置
	logConfig := &logger.Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	}
	log, err := logger.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("创建logger失败: %v", err)
	}

	// 创建回测服务
	dailyCacheService := NewDailyCacheService(nil)
	backtestService := NewBacktestService(nil, nil, dailyCacheService, log)

	tests := []struct {
		name        string
		trades      []models.Trade
		expectError bool
		description string
	}{
		{
			name:        "空交易记录",
			trades:      []models.Trade{},
			expectError: false,
			description: "空交易记录应该通过验证",
		},
		{
			name: "正常的买入卖出序列",
			trades: []models.Trade{
				{
					ID:            "buy1",
					Symbol:        "600976",
					Side:          models.TradeSideBuy,
					HoldingAssets: 100000.0, // 买入后持仓10万
					Timestamp:     time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
				},
				{
					ID:            "sell1",
					Symbol:        "600976",
					Side:          models.TradeSideSell,
					HoldingAssets: 80000.0, // 卖出后持仓8万（正常减少）
					Timestamp:     time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC),
				},
			},
			expectError: false,
			description: "正常的买入卖出序列应该通过验证",
		},
		{
			name: "异常：单股票买卖序列中卖出后持仓资产增加",
			trades: []models.Trade{
				{
					ID:            "buy1",
					Symbol:        "600976",
					Side:          models.TradeSideBuy,
					HoldingAssets: 100000.0, // 买入后持仓10万
					Timestamp:     time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
				},
				{
					ID:            "sell1",
					Symbol:        "600976",
					Side:          models.TradeSideSell,
					HoldingAssets: 120000.0, // 🚨 卖出后持仓12万（异常增加）
					Timestamp:     time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC),
				},
			},
			expectError: true,
			description: "单股票买卖序列中卖出后持仓资产增加应该触发错误",
		},
		{
			name: "多只股票混合交易",
			trades: []models.Trade{
				{
					ID:            "buy1",
					Symbol:        "600976",
					Side:          models.TradeSideBuy,
					HoldingAssets: 100000.0,
					Timestamp:     time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
				},
				{
					ID:            "buy2",
					Symbol:        "000001",
					Side:          models.TradeSideBuy,
					HoldingAssets: 200000.0,
					Timestamp:     time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
				},
				{
					ID:            "sell1",
					Symbol:        "600976",
					Side:          models.TradeSideSell,
					HoldingAssets: 180000.0, // 正常：只剩000001的持仓
					Timestamp:     time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC),
				},
				{
					ID:            "sell2",
					Symbol:        "000001",
					Side:          models.TradeSideSell,
					HoldingAssets: 0.0, // 正常：全部卖出
					Timestamp:     time.Date(2024, 1, 2, 10, 30, 0, 0, time.UTC),
				},
			},
			expectError: false,
			description: "多只股票的正常交易序列应该通过验证",
		},
		{
			name: "时间顺序混乱的交易记录",
			trades: []models.Trade{
				{
					ID:            "sell1",
					Symbol:        "600976",
					Side:          models.TradeSideSell,
					HoldingAssets: 80000.0,
					Timestamp:     time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC),
				},
				{
					ID:            "buy1",
					Symbol:        "600976",
					Side:          models.TradeSideBuy,
					HoldingAssets: 100000.0,
					Timestamp:     time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC), // 时间更早
				},
			},
			expectError: false,
			description: "函数应该自动按时间排序后再验证",
		},
		{
			name: "复杂场景：多次买入后卖出（不触发错误）",
			trades: []models.Trade{
				{
					ID:            "buy1",
					Symbol:        "600976",
					Side:          models.TradeSideBuy,
					HoldingAssets: 50000.0,
					Timestamp:     time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
				},
				{
					ID:            "buy2",
					Symbol:        "600976",
					Side:          models.TradeSideBuy,
					HoldingAssets: 100000.0, // 加仓后持仓10万
					Timestamp:     time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
				},
				{
					ID:            "sell1",
					Symbol:        "600976",
					Side:          models.TradeSideSell,
					HoldingAssets: 150000.0, // 卖出后持仓15万（多次交易不触发简单序列检查）
					Timestamp:     time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC),
				},
			},
			expectError: false,
			description: "多次买入后的卖出不会触发简单序列检查，但可能记录警告",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := backtestService.validateTradesData(tt.trades, "test-backtest-id")

			if tt.expectError && err == nil {
				t.Errorf("期望出现错误，但没有错误返回。%s", tt.description)
			}

			if !tt.expectError && err != nil {
				t.Errorf("不期望出现错误，但返回了错误: %v。%s", err, tt.description)
			}

			if err != nil {
				t.Logf("验证错误信息: %v", err)
			}
		})
	}
}

func TestValidateSingleStockTrades(t *testing.T) {
	// 创建测试用的日志配置
	logConfig := &logger.Config{
		Level:  "info",
		Format: "console",
	}
	log, _ := logger.NewLogger(logConfig)

	// 创建回测服务
	dailyCacheService := NewDailyCacheService(nil)
	backtestService := NewBacktestService(nil, nil, dailyCacheService, log)

	tests := []struct {
		name        string
		symbol      string
		trades      []models.Trade
		expectError bool
		description string
	}{
		{
			name:   "单只股票正常交易",
			symbol: "600976",
			trades: []models.Trade{
				{
					ID:            "buy1",
					Symbol:        "600976",
					Side:          models.TradeSideBuy,
					HoldingAssets: 100000.0,
					Timestamp:     time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
				},
				{
					ID:            "sell1",
					Symbol:        "600976",
					Side:          models.TradeSideSell,
					HoldingAssets: 0.0, // 全部卖出
					Timestamp:     time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC),
				},
			},
			expectError: false,
			description: "单只股票的正常买入卖出应该通过验证",
		},
		{
			name:   "单只股票异常交易",
			symbol: "600976",
			trades: []models.Trade{
				{
					ID:            "buy1",
					Symbol:        "600976",
					Side:          models.TradeSideBuy,
					HoldingAssets: 100000.0,
					Timestamp:     time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
				},
				{
					ID:            "sell1",
					Symbol:        "600976",
					Side:          models.TradeSideSell,
					HoldingAssets: 150000.0, // 🚨 异常增加
					Timestamp:     time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC),
				},
			},
			expectError: true,
			description: "单只股票的异常交易应该被检测到",
		},
		{
			name:   "只有买入没有卖出",
			symbol: "600976",
			trades: []models.Trade{
				{
					ID:            "buy1",
					Symbol:        "600976",
					Side:          models.TradeSideBuy,
					HoldingAssets: 100000.0,
					Timestamp:     time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
				},
			},
			expectError: false,
			description: "只有买入操作应该通过验证",
		},
		{
			name:   "只有卖出没有买入",
			symbol: "600976",
			trades: []models.Trade{
				{
					ID:            "sell1",
					Symbol:        "600976",
					Side:          models.TradeSideSell,
					HoldingAssets: 50000.0,
					Timestamp:     time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC),
				},
			},
			expectError: false,
			description: "只有卖出操作应该通过验证（可能是已有持仓的卖出）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := backtestService.validateSingleStockTrades(tt.symbol, tt.trades, "test-backtest-id", false)

			if tt.expectError && err == nil {
				t.Errorf("期望出现错误，但没有错误返回。%s", tt.description)
			}

			if !tt.expectError && err != nil {
				t.Errorf("不期望出现错误，但返回了错误: %v。%s", err, tt.description)
			}

			if err != nil {
				t.Logf("验证错误信息: %v", err)
			}
		})
	}
}

// TestValidateTradesDataIntegration 集成测试：测试在GetBacktestResults中的实际使用
func TestValidateTradesDataIntegration(t *testing.T) {
	// 创建测试用的日志配置
	logConfig := &logger.Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	}
	log, err := logger.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("创建logger失败: %v", err)
	}

	// 创建回测服务
	dailyCacheService := NewDailyCacheService(nil)
	backtestService := NewBacktestService(nil, nil, dailyCacheService, log)

	// 创建测试回测
	backtestID := "test-validation-integration"
	backtest := &models.Backtest{
		ID:          backtestID,
		Name:        "数据验证集成测试",
		StrategyID:  "test-strategy",
		Symbols:     []string{"600976"},
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
		InitialCash: 1000000,
		Commission:  0.001,
		Status:      models.BacktestStatusCompleted,
		CreatedAt:   time.Now(),
	}

	// 创建测试结果
	result := &models.BacktestResult{
		ID:           "test-result",
		BacktestID:   backtestID,
		StrategyID:   "test-strategy",
		TotalReturn:  0.10,
		AnnualReturn: 0.15,
		CreatedAt:    time.Now(),
	}

	// 创建包含异常数据的交易记录
	trades := []models.Trade{
		{
			ID:            "buy1",
			BacktestID:    backtestID,
			Symbol:        "600976",
			Side:          models.TradeSideBuy,
			HoldingAssets: 100000.0,
			Timestamp:     time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
		},
		{
			ID:            "sell1",
			BacktestID:    backtestID,
			Symbol:        "600976",
			Side:          models.TradeSideSell,
			HoldingAssets: 120000.0, // 🚨 异常：卖出后持仓增加
			Timestamp:     time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC),
		},
	}

	// 设置测试数据
	backtestService.backtests[backtestID] = backtest
	backtestService.backtestResults[backtestID] = result
	backtestService.backtestTrades[backtestID] = trades

	// 调用GetBacktestResults，应该记录错误日志但不影响结果返回
	response, err := backtestService.GetBacktestResults(context.Background(), backtestID)

	// 验证结果
	if err != nil {
		t.Errorf("GetBacktestResults不应该返回错误，但返回了: %v", err)
	}

	if response == nil {
		t.Fatal("GetBacktestResults应该返回结果，但返回了nil")
	}

	if len(response.Trades) != 2 {
		t.Errorf("期望返回2条交易记录，但返回了%d条", len(response.Trades))
	}

	t.Logf("✅ 集成测试通过：即使存在数据异常，GetBacktestResults仍能正常返回结果，同时记录错误日志")
}
