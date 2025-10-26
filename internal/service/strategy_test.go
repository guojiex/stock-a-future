package service

import (
	"context"
	"testing"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"
)

func TestGetStrategyTemplates(t *testing.T) {
	log, err := logger.NewLogger(&logger.Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	service := NewStrategyService(log)

	templates := service.GetStrategyTemplates()

	if len(templates) == 0 {
		t.Fatal("Expected at least one template")
	}

	// 验证模板包含必要字段
	for _, template := range templates {
		if template.ID == "" {
			t.Error("Template ID should not be empty")
		}
		if template.Name == "" {
			t.Error("Template Name should not be empty")
		}
		if template.Type == "" {
			t.Error("Template Type should not be empty")
		}
		if template.Parameters == nil {
			t.Error("Template Parameters should not be nil")
		}
	}

	t.Logf("✅ 成功获取 %d 个策略模板", len(templates))
}

func TestGetStrategyTypeDefinitions(t *testing.T) {
	log, err := logger.NewLogger(&logger.Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	service := NewStrategyService(log)

	typeDefs := service.GetStrategyTypeDefinitions()

	if len(typeDefs) == 0 {
		t.Fatal("Expected at least one type definition")
	}

	// 验证类型定义包含必要字段
	for _, typeDef := range typeDefs {
		if typeDef.Type == "" {
			t.Error("Type should not be empty")
		}
		if typeDef.Name == "" {
			t.Error("Name should not be empty")
		}
		t.Logf("  类型: %s - %s", typeDef.Type, typeDef.Name)
	}

	t.Logf("✅ 成功获取 %d 个策略类型定义", len(typeDefs))
}

func TestValidateParameters(t *testing.T) {
	log, err := logger.NewLogger(&logger.Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	service := NewStrategyService(log)

	tests := []struct {
		name          string
		strategyType  models.StrategyType
		parameters    map[string]interface{}
		expectErrors  bool
		errorContains string
	}{
		{
			name:         "有效的MACD参数",
			strategyType: models.StrategyTypeTechnical,
			parameters: map[string]interface{}{
				"fast_period": 12.0,
				"slow_period": 26.0,
			},
			expectErrors: false,
		},
		{
			name:         "快线周期超出范围",
			strategyType: models.StrategyTypeTechnical,
			parameters: map[string]interface{}{
				"fast_period": 100.0,
				"slow_period": 26.0,
			},
			expectErrors:  true,
			errorContains: "快线周期",
		},
		{
			name:         "慢线小于快线",
			strategyType: models.StrategyTypeTechnical,
			parameters: map[string]interface{}{
				"fast_period": 26.0,
				"slow_period": 12.0,
			},
			expectErrors:  true,
			errorContains: "慢线周期必须大于快线周期",
		},
		{
			name:         "RSI参数超出范围",
			strategyType: models.StrategyTypeTechnical,
			parameters: map[string]interface{}{
				"period":     14.0,
				"overbought": 120.0, // 无效值
				"oversold":   30.0,
			},
			expectErrors:  true,
			errorContains: "超买阈值",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := service.ValidateParameters(tt.strategyType, tt.parameters)

			if tt.expectErrors {
				if len(errors) == 0 {
					t.Errorf("Expected validation errors but got none")
				} else {
					t.Logf("✅ 验证失败(预期): %v", errors)
					// 检查错误信息是否包含预期内容
					if tt.errorContains != "" {
						found := false
						for _, err := range errors {
							if msg, ok := err["message"]; ok {
								if contains(msg, tt.errorContains) {
									found = true
									break
								}
							}
						}
						if !found {
							t.Errorf("Expected error to contain '%s', got %v", tt.errorContains, errors)
						}
					}
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("Expected no validation errors but got: %v", errors)
				} else {
					t.Logf("✅ 参数验证通过")
				}
			}
		})
	}
}

func TestCreateStrategyWithValidation(t *testing.T) {
	log, err := logger.NewLogger(&logger.Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	service := NewStrategyService(log)
	ctx := context.Background()

	// 创建一个有效的策略
	strategy := &models.Strategy{
		ID:          "test_strategy_1",
		Name:        "测试MACD策略",
		Description: "用于测试的MACD策略",
		Type:        models.StrategyTypeTechnical,
		Status:      models.StrategyStatusInactive,
		Parameters: map[string]interface{}{
			"fast_period": 12.0,
			"slow_period": 26.0,
		},
	}

	// 验证参数
	errors := service.ValidateParameters(strategy.Type, strategy.Parameters)
	if len(errors) > 0 {
		t.Fatalf("参数验证失败: %v", errors)
	}

	// 创建策略
	err = service.CreateStrategy(ctx, strategy)
	if err != nil {
		t.Fatalf("创建策略失败: %v", err)
	}

	t.Logf("✅ 策略创建成功: %s", strategy.Name)

	// 验证策略已创建
	req := &models.StrategyListRequest{
		Page: 1,
		Size: 100,
	}
	strategies, total, err := service.GetStrategiesList(ctx, req)
	if err != nil {
		t.Fatalf("获取策略列表失败: %v", err)
	}

	found := false
	for _, s := range strategies {
		if s.ID == strategy.ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("创建的策略未在列表中找到")
	} else {
		t.Logf("✅ 策略列表中共有 %d 个策略，包含新创建的策略", total)
	}
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

// indexOf 查找子串位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
