package service

import (
	"context"
	"testing"
	"time"

	"stock-a-future/config"
	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"

	"github.com/google/uuid"
)

// TestBacktestAutoRename 测试回测自动重命名功能
func TestBacktestAutoRename(t *testing.T) {
	// 创建测试服务
	logConfig := &logger.Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	}
	log, err := logger.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("创建日志器失败: %v", err)
	}

	strategyService := NewStrategyService(log)

	cfg := &config.Config{} // 使用空配置进行测试
	dataSourceService := NewDataSourceService(cfg)

	dailyCacheService := NewDailyCacheService(nil) // 使用默认配置

	backtestService := NewBacktestService(strategyService, dataSourceService, dailyCacheService, log)

	ctx := context.Background()

	// 创建第一个回测
	backtest1 := &models.Backtest{
		ID:          uuid.New().String(),
		Name:        "中国东航 - 回测分析 2025-09-14",
		StrategyID:  "test-strategy-1",
		Symbols:     []string{"600115"},
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		InitialCash: 100000,
		Commission:  0.001,
		Slippage:    0.001,
		Benchmark:   "000300",
		Status:      models.BacktestStatusPending,
		Progress:    0,
		CreatedBy:   "test-user",
	}

	err = backtestService.CreateBacktest(ctx, backtest1)
	if err != nil {
		t.Fatalf("创建第一个回测失败: %v", err)
	}

	t.Logf("✅ 第一个回测创建成功: %s", backtest1.Name)

	// 创建第二个回测，使用相同名称
	backtest2 := &models.Backtest{
		ID:          uuid.New().String(),
		Name:        "中国东航 - 回测分析 2025-09-14", // 相同名称
		StrategyID:  "test-strategy-2",
		Symbols:     []string{"600115"},
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		InitialCash: 100000,
		Commission:  0.001,
		Slippage:    0.001,
		Benchmark:   "000300",
		Status:      models.BacktestStatusPending,
		Progress:    0,
		CreatedBy:   "test-user",
	}

	originalName := backtest2.Name
	err = backtestService.CreateBacktest(ctx, backtest2)
	if err != nil {
		t.Fatalf("创建第二个回测失败: %v", err)
	}

	// 检查名称是否被自动重命名
	if backtest2.Name == originalName {
		t.Errorf("期望名称被自动重命名，但仍然是原名称: %s", backtest2.Name)
	} else {
		t.Logf("✅ 第二个回测自动重命名成功: %s -> %s", originalName, backtest2.Name)
	}

	// 验证重命名的格式是否正确
	expectedName := "中国东航 - 回测分析 2025-09-14 (2)"
	if backtest2.Name != expectedName {
		t.Errorf("重命名格式不正确。期望: %s, 实际: %s", expectedName, backtest2.Name)
	} else {
		t.Logf("✅ 重命名格式正确: %s", backtest2.Name)
	}

	// 创建第三个回测，再次使用相同的原始名称
	backtest3 := &models.Backtest{
		ID:          uuid.New().String(),
		Name:        "中国东航 - 回测分析 2025-09-14", // 再次使用相同名称
		StrategyID:  "test-strategy-3",
		Symbols:     []string{"600115"},
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		InitialCash: 100000,
		Commission:  0.001,
		Slippage:    0.001,
		Benchmark:   "000300",
		Status:      models.BacktestStatusPending,
		Progress:    0,
		CreatedBy:   "test-user",
	}

	err = backtestService.CreateBacktest(ctx, backtest3)
	if err != nil {
		t.Fatalf("创建第三个回测失败: %v", err)
	}

	// 验证第三个回测的名称
	expectedName3 := "中国东航 - 回测分析 2025-09-14 (3)"
	if backtest3.Name != expectedName3 {
		t.Errorf("第三个回测重命名不正确。期望: %s, 实际: %s", expectedName3, backtest3.Name)
	} else {
		t.Logf("✅ 第三个回测重命名正确: %s", backtest3.Name)
	}

	// 验证所有回测都存在且名称不同
	list, total, err := backtestService.GetBacktestsList(ctx, &models.BacktestListRequest{
		Page: 1,
		Size: 10,
	})
	if err != nil {
		t.Fatalf("获取回测列表失败: %v", err)
	}

	if total != 3 {
		t.Errorf("期望3个回测，实际: %d", total)
	}

	t.Logf("✅ 总共创建了 %d 个回测", total)
	for i, bt := range list {
		t.Logf("  回测 %d: %s", i+1, bt.Name)
	}
}

// TestBacktestUpdateAutoRename 测试回测更新时的自动重命名功能
func TestBacktestUpdateAutoRename(t *testing.T) {
	// 创建测试服务
	logConfig := &logger.Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	}
	log, err := logger.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("创建日志器失败: %v", err)
	}

	strategyService := NewStrategyService(log)

	cfg := &config.Config{} // 使用空配置进行测试
	dataSourceService := NewDataSourceService(cfg)

	dailyCacheService := NewDailyCacheService(nil) // 使用默认配置

	backtestService := NewBacktestService(strategyService, dataSourceService, dailyCacheService, log)

	ctx := context.Background()

	// 创建两个不同名称的回测
	backtest1 := &models.Backtest{
		ID:          uuid.New().String(),
		Name:        "回测A",
		StrategyID:  "test-strategy-1",
		Symbols:     []string{"600115"},
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		InitialCash: 100000,
		Commission:  0.001,
		Slippage:    0.001,
		Benchmark:   "000300",
		Status:      models.BacktestStatusPending,
		Progress:    0,
		CreatedBy:   "test-user",
	}

	backtest2 := &models.Backtest{
		ID:          uuid.New().String(),
		Name:        "回测B",
		StrategyID:  "test-strategy-2",
		Symbols:     []string{"600115"},
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		InitialCash: 100000,
		Commission:  0.001,
		Slippage:    0.001,
		Benchmark:   "000300",
		Status:      models.BacktestStatusPending,
		Progress:    0,
		CreatedBy:   "test-user",
	}

	// 创建两个回测
	err = backtestService.CreateBacktest(ctx, backtest1)
	if err != nil {
		t.Fatalf("创建第一个回测失败: %v", err)
	}

	err = backtestService.CreateBacktest(ctx, backtest2)
	if err != nil {
		t.Fatalf("创建第二个回测失败: %v", err)
	}

	t.Logf("✅ 创建了两个回测: %s, %s", backtest1.Name, backtest2.Name)

	// 尝试将回测2的名称更新为与回测1相同的名称
	newName := "回测A"
	updateReq := &models.UpdateBacktestRequest{
		Name: &newName,
	}

	err = backtestService.UpdateBacktest(ctx, backtest2.ID, updateReq)
	if err != nil {
		t.Fatalf("更新回测失败: %v", err)
	}

	// 获取更新后的回测
	updatedBacktest, err := backtestService.GetBacktest(ctx, backtest2.ID)
	if err != nil {
		t.Fatalf("获取更新后的回测失败: %v", err)
	}

	// 检查名称是否被自动重命名
	expectedName := "回测A (2)"
	if updatedBacktest.Name != expectedName {
		t.Errorf("更新时重命名不正确。期望: %s, 实际: %s", expectedName, updatedBacktest.Name)
	} else {
		t.Logf("✅ 更新时自动重命名成功: %s -> %s", newName, updatedBacktest.Name)
	}
}
