package service

import (
	"context"
	"testing"
	"time"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"
)

// TestStrategyListSorting 测试策略列表排序的稳定性
func TestStrategyListSorting(t *testing.T) {
	// 创建策略服务
	log, _ := logger.NewLogger(nil)
	service := NewStrategyService(log)

	// 多次获取策略列表，验证排序的稳定性
	req := &models.StrategyListRequest{
		Page: 1,
		Size: 20,
	}

	// 获取10次策略列表
	var previousOrder []string
	for i := 0; i < 10; i++ {
		strategies, _, err := service.GetStrategiesList(context.Background(), req)
		if err != nil {
			t.Fatalf("获取策略列表失败: %v", err)
		}

		// 提取策略ID顺序
		currentOrder := make([]string, len(strategies))
		for j, strategy := range strategies {
			currentOrder[j] = strategy.ID
		}

		// 第一次保存顺序
		if i == 0 {
			previousOrder = currentOrder
			t.Logf("第 %d 次获取策略列表顺序: %v", i+1, currentOrder)
			continue
		}

		// 后续验证顺序是否一致
		if len(currentOrder) != len(previousOrder) {
			t.Errorf("第 %d 次获取的策略数量不一致: 期望 %d, 实际 %d",
				i+1, len(previousOrder), len(currentOrder))
			continue
		}

		for j := range currentOrder {
			if currentOrder[j] != previousOrder[j] {
				t.Errorf("第 %d 次获取的策略顺序不一致:\n期望: %v\n实际: %v",
					i+1, previousOrder, currentOrder)
				break
			}
		}

		t.Logf("第 %d 次获取策略列表顺序: %v (✓ 与第1次一致)", i+1, currentOrder)
	}
}

// TestStrategyListSortingRules 测试策略列表排序规则
func TestStrategyListSortingRules(t *testing.T) {
	log, _ := logger.NewLogger(nil)
	service := NewStrategyService(log)

	// 添加测试策略
	now := time.Now()
	testStrategies := []models.Strategy{
		{
			ID:        "test_active_1",
			Name:      "测试活跃策略1",
			Type:      models.StrategyTypeTechnical,
			Status:    models.StrategyStatusActive,
			CreatedAt: now.Add(-2 * time.Hour),
			UpdatedAt: now,
			CreatedBy: "test",
		},
		{
			ID:        "test_active_2",
			Name:      "测试活跃策略2",
			Type:      models.StrategyTypeTechnical,
			Status:    models.StrategyStatusActive,
			CreatedAt: now.Add(-1 * time.Hour), // 更新
			UpdatedAt: now,
			CreatedBy: "test",
		},
		{
			ID:        "test_testing_1",
			Name:      "测试测试中策略1",
			Type:      models.StrategyTypeTechnical,
			Status:    models.StrategyStatusTesting,
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: "test",
		},
		{
			ID:        "test_inactive_1",
			Name:      "测试非活跃策略1",
			Type:      models.StrategyTypeTechnical,
			Status:    models.StrategyStatusInactive,
			CreatedAt: now.Add(-3 * time.Hour),
			UpdatedAt: now,
			CreatedBy: "test",
		},
	}

	for _, strategy := range testStrategies {
		strategyCopy := strategy
		service.strategies[strategy.ID] = &strategyCopy
	}

	// 获取策略列表
	req := &models.StrategyListRequest{
		Page: 1,
		Size: 20,
	}

	strategies, _, err := service.GetStrategiesList(context.Background(), req)
	if err != nil {
		t.Fatalf("获取策略列表失败: %v", err)
	}

	t.Logf("策略排序结果:")
	for i, strategy := range strategies {
		t.Logf("%d. %s (状态: %s, 创建时间: %s)",
			i+1, strategy.ID, strategy.Status, strategy.CreatedAt.Format("15:04:05"))
	}

	// 验证排序规则
	// 1. 所有 Active 策略应该在最前面
	// 2. 然后是 Testing 策略
	// 3. 最后是 Inactive 策略
	// 4. 同状态下，创建时间新的在前
	// 5. 创建时间相同时，按ID字典序

	var lastStatus models.StrategyStatus
	var lastCreatedAt time.Time

	for i, strategy := range strategies {
		// 检查状态顺序
		if i > 0 {
			statusPriority := map[models.StrategyStatus]int{
				models.StrategyStatusActive:   1,
				models.StrategyStatusTesting:  2,
				models.StrategyStatusInactive: 3,
			}

			currentPriority := statusPriority[strategy.Status]
			lastPriority := statusPriority[lastStatus]

			if currentPriority < lastPriority {
				t.Errorf("策略状态顺序错误: %s (优先级 %d) 出现在 %s (优先级 %d) 之后",
					strategy.Status, currentPriority, lastStatus, lastPriority)
			}

			// 如果状态相同，检查创建时间
			if strategy.Status == lastStatus {
				if strategy.CreatedAt.After(lastCreatedAt) {
					t.Errorf("相同状态下创建时间顺序错误: %s (%s) 应该在 %s 之前",
						strategy.ID, strategy.CreatedAt.Format("15:04:05"),
						lastCreatedAt.Format("15:04:05"))
				}
			}
		}

		lastStatus = strategy.Status
		lastCreatedAt = strategy.CreatedAt
	}
}

// TestStrategyListSortingWithSameTimestamp 测试相同时间戳的稳定排序
func TestStrategyListSortingWithSameTimestamp(t *testing.T) {
	log, _ := logger.NewLogger(nil)
	service := NewStrategyService(log)

	// 清空默认策略
	service.strategies = make(map[string]*models.Strategy)

	// 创建多个具有相同状态和创建时间的策略
	now := time.Now()
	sameTimeStrategies := []models.Strategy{
		{
			ID:        "strategy_c",
			Name:      "策略C",
			Type:      models.StrategyTypeTechnical,
			Status:    models.StrategyStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: "test",
		},
		{
			ID:        "strategy_a",
			Name:      "策略A",
			Type:      models.StrategyTypeTechnical,
			Status:    models.StrategyStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: "test",
		},
		{
			ID:        "strategy_b",
			Name:      "策略B",
			Type:      models.StrategyTypeTechnical,
			Status:    models.StrategyStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: "test",
		},
	}

	for _, strategy := range sameTimeStrategies {
		strategyCopy := strategy
		service.strategies[strategy.ID] = &strategyCopy
	}

	req := &models.StrategyListRequest{
		Page: 1,
		Size: 20,
	}

	// 多次获取并验证顺序
	expectedOrder := []string{"strategy_a", "strategy_b", "strategy_c"} // 应该按ID字典序

	for i := 0; i < 5; i++ {
		strategies, _, err := service.GetStrategiesList(context.Background(), req)
		if err != nil {
			t.Fatalf("获取策略列表失败: %v", err)
		}

		if len(strategies) != len(expectedOrder) {
			t.Errorf("第 %d 次: 策略数量不正确，期望 %d，实际 %d",
				i+1, len(expectedOrder), len(strategies))
			continue
		}

		for j, strategy := range strategies {
			if strategy.ID != expectedOrder[j] {
				t.Errorf("第 %d 次: 位置 %d 的策略ID不正确，期望 %s，实际 %s",
					i+1, j, expectedOrder[j], strategy.ID)
			}
		}

		t.Logf("第 %d 次获取: %v (✓ 顺序正确)", i+1, extractIDs(strategies))
	}
}

// extractIDs 提取策略ID列表
func extractIDs(strategies []models.Strategy) []string {
	ids := make([]string, len(strategies))
	for i, strategy := range strategies {
		ids[i] = strategy.ID
	}
	return ids
}
