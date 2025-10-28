package service

import (
	"testing"
	"time"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"
)

// TestCalculateEquityCurveFromTrades 测试权益曲线计算
func TestCalculateEquityCurveFromTrades(t *testing.T) {
	// 创建测试服务
	log, err := logger.NewLogger(logger.DefaultConfig())
	if err != nil {
		t.Fatalf("创建logger失败: %v", err)
	}
	service := &BacktestService{
		logger: log,
	}

	// 创建测试回测配置
	backtest := &models.Backtest{
		ID:          "test-backtest-1",
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		InitialCash: 1000000.0,
	}

	strategyID := "strategy-1"

	// 测试用例1: 没有交易
	t.Run("无交易的策略", func(t *testing.T) {
		trades := []models.Trade{}
		curve := service.calculateEquityCurveFromTrades(backtest, strategyID, trades)

		if len(curve) != 2 {
			t.Errorf("预期权益曲线有2个点(起点和终点)，实际有%d个点", len(curve))
		}

		if curve[0].PortfolioValue != backtest.InitialCash {
			t.Errorf("起点权益值错误: 预期%.2f, 实际%.2f", backtest.InitialCash, curve[0].PortfolioValue)
		}

		if curve[1].PortfolioValue != backtest.InitialCash {
			t.Errorf("终点权益值错误: 预期%.2f, 实际%.2f", backtest.InitialCash, curve[1].PortfolioValue)
		}
	})

	// 测试用例2: 只有买入交易
	t.Run("只有买入交易", func(t *testing.T) {
		trades := []models.Trade{
			{
				ID:            "trade-1",
				StrategyID:    strategyID,
				Symbol:        "000001.SZ",
				Side:          models.TradeSideBuy,
				Quantity:      1000,
				Price:         10.0,
				Commission:    30.0,
				TotalAssets:   999970.0, // 1000000 - 10*1000 - 30 = 999970
				CashBalance:   989970.0, // 1000000 - 10*1000 - 30
				HoldingAssets: 10000.0,  // 10 * 1000
				Timestamp:     time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		curve := service.calculateEquityCurveFromTrades(backtest, strategyID, trades)

		// 应该有3个点：起点 + 1个买入交易
		if len(curve) != 2 {
			t.Errorf("预期权益曲线有2个点(起点+买入)，实际有%d个点", len(curve))
		}

		// 检查买入后的权益值
		if curve[1].PortfolioValue != 999970.0 {
			t.Errorf("买入后权益值错误: 预期999970.0, 实际%.2f", curve[1].PortfolioValue)
		}

		// 检查现金和持仓
		if curve[1].Cash != 989970.0 {
			t.Errorf("买入后现金错误: 预期989970.0, 实际%.2f", curve[1].Cash)
		}

		if curve[1].Holdings != 10000.0 {
			t.Errorf("买入后持仓错误: 预期10000.0, 实际%.2f", curve[1].Holdings)
		}
	})

	// 测试用例3: 买入后卖出（盈利）
	t.Run("买入后卖出-盈利", func(t *testing.T) {
		trades := []models.Trade{
			{
				ID:            "trade-1",
				StrategyID:    strategyID,
				Symbol:        "000001.SZ",
				Side:          models.TradeSideBuy,
				Quantity:      1000,
				Price:         10.0,
				Commission:    30.0,
				TotalAssets:   999970.0,
				CashBalance:   989970.0,
				HoldingAssets: 10000.0,
				Timestamp:     time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:            "trade-2",
				StrategyID:    strategyID,
				Symbol:        "000001.SZ",
				Side:          models.TradeSideSell,
				Quantity:      1000,
				Price:         12.0,
				Commission:    36.0,
				PnL:           1934.0, // (12-10)*1000 - 30 - 36
				TotalAssets:   1001904.0,
				CashBalance:   1001904.0,
				HoldingAssets: 0.0,
				Timestamp:     time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		curve := service.calculateEquityCurveFromTrades(backtest, strategyID, trades)

		// 应该有3个点：起点 + 买入 + 卖出
		if len(curve) != 3 {
			t.Errorf("预期权益曲线有3个点(起点+买入+卖出)，实际有%d个点", len(curve))
		}

		// 检查卖出后的权益值
		lastPoint := curve[len(curve)-1]
		if lastPoint.PortfolioValue != 1001904.0 {
			t.Errorf("卖出后权益值错误: 预期1001904.0, 实际%.2f", lastPoint.PortfolioValue)
		}

		// 验证盈利
		profit := lastPoint.PortfolioValue - backtest.InitialCash
		expectedProfit := 1904.0
		if profit < expectedProfit-1 || profit > expectedProfit+1 {
			t.Errorf("盈利金额错误: 预期约%.2f, 实际%.2f", expectedProfit, profit)
		}
	})

	// 测试用例4: 多次交易
	t.Run("多次交易", func(t *testing.T) {
		trades := []models.Trade{
			{
				ID:            "trade-1",
				StrategyID:    strategyID,
				Symbol:        "000001.SZ",
				Side:          models.TradeSideBuy,
				Quantity:      1000,
				Price:         10.0,
				Commission:    30.0,
				TotalAssets:   999970.0,
				CashBalance:   989970.0,
				HoldingAssets: 10000.0,
				Timestamp:     time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:            "trade-2",
				StrategyID:    strategyID,
				Symbol:        "000001.SZ",
				Side:          models.TradeSideSell,
				Quantity:      1000,
				Price:         11.0,
				Commission:    33.0,
				PnL:           937.0,
				TotalAssets:   1000907.0,
				CashBalance:   1000907.0,
				HoldingAssets: 0.0,
				Timestamp:     time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:            "trade-3",
				StrategyID:    strategyID,
				Symbol:        "000002.SZ",
				Side:          models.TradeSideBuy,
				Quantity:      2000,
				Price:         5.0,
				Commission:    30.0,
				TotalAssets:   1000877.0,
				CashBalance:   990877.0,
				HoldingAssets: 10000.0,
				Timestamp:     time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		curve := service.calculateEquityCurveFromTrades(backtest, strategyID, trades)

		// 应该有4个点：起点 + 3个交易
		if len(curve) != 4 {
			t.Errorf("预期权益曲线有4个点(起点+3个交易)，实际有%d个点", len(curve))
		}

		// 验证最终权益
		lastPoint := curve[len(curve)-1]
		if lastPoint.PortfolioValue != 1000877.0 {
			t.Errorf("最终权益值错误: 预期1000877.0, 实际%.2f", lastPoint.PortfolioValue)
		}
	})

	// 测试用例5: 其他策略的交易不应被包含
	t.Run("过滤其他策略的交易", func(t *testing.T) {
		trades := []models.Trade{
			{
				ID:            "trade-1",
				StrategyID:    strategyID,
				Symbol:        "000001.SZ",
				Side:          models.TradeSideBuy,
				Quantity:      1000,
				Price:         10.0,
				Commission:    30.0,
				TotalAssets:   999970.0,
				CashBalance:   989970.0,
				HoldingAssets: 10000.0,
				Timestamp:     time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:            "trade-2",
				StrategyID:    "strategy-2", // 不同的策略
				Symbol:        "000002.SZ",
				Side:          models.TradeSideBuy,
				Quantity:      2000,
				Price:         5.0,
				Commission:    30.0,
				TotalAssets:   999970.0,
				CashBalance:   989970.0,
				HoldingAssets: 10000.0,
				Timestamp:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		curve := service.calculateEquityCurveFromTrades(backtest, strategyID, trades)

		// 应该只有2个点：起点 + strategy-1的1个交易
		if len(curve) != 2 {
			t.Errorf("预期权益曲线有2个点(起点+1个交易)，实际有%d个点", len(curve))
		}
	})
}
