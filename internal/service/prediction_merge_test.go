package service

import (
	"stock-a-future/internal/models"
	"testing"

	"github.com/shopspring/decimal"
)

// TestMergeSameDaySignals 测试同一天信号的合并
func TestMergeSameDaySignals(t *testing.T) {
	service := NewPredictionService()

	tests := []struct {
		name           string
		predictions    []models.TradingPointPrediction
		expectedCount  int
		expectedMerged bool // 是否期望发生合并
	}{
		{
			name: "同一天的多个买入信号应该被合并",
			predictions: []models.TradingPointPrediction{
				{
					Type:        "BUY",
					Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
					SignalDate:  "20240101",
					Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.65)),
					Reason:      "MACD金叉信号",
					Indicators:  []string{"MACD"},
				},
				{
					Type:        "BUY",
					Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.3)),
					SignalDate:  "20240101",
					Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.70)),
					Reason:      "RSI超卖信号",
					Indicators:  []string{"RSI"},
				},
				{
					Type:        "BUY",
					Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.4)),
					SignalDate:  "20240101",
					Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.60)),
					Reason:      "布林带下轨反弹",
					Indicators:  []string{"BOLL"},
				},
			},
			expectedCount:  1,
			expectedMerged: true,
		},
		{
			name: "同一天的买入和卖出信号不应该被合并",
			predictions: []models.TradingPointPrediction{
				{
					Type:        "BUY",
					Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
					SignalDate:  "20240101",
					Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.65)),
					Reason:      "MACD金叉信号",
					Indicators:  []string{"MACD"},
				},
				{
					Type:        "SELL",
					Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.8)),
					SignalDate:  "20240101",
					Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.60)),
					Reason:      "RSI超买信号",
					Indicators:  []string{"RSI"},
				},
			},
			expectedCount:  2,
			expectedMerged: false,
		},
		{
			name: "不同日期的信号不应该被合并",
			predictions: []models.TradingPointPrediction{
				{
					Type:        "BUY",
					Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
					SignalDate:  "20240101",
					Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.65)),
					Reason:      "MACD金叉信号",
					Indicators:  []string{"MACD"},
				},
				{
					Type:        "BUY",
					Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.6)),
					SignalDate:  "20240102",
					Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.70)),
					Reason:      "RSI超卖信号",
					Indicators:  []string{"RSI"},
				},
			},
			expectedCount:  2,
			expectedMerged: false,
		},
		{
			name: "单个信号不应该被修改",
			predictions: []models.TradingPointPrediction{
				{
					Type:        "BUY",
					Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
					SignalDate:  "20240101",
					Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.65)),
					Reason:      "MACD金叉信号",
					Indicators:  []string{"MACD"},
				},
			},
			expectedCount:  1,
			expectedMerged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.mergeSameDaySignals(tt.predictions)

			if len(result) != tt.expectedCount {
				t.Errorf("期望得到 %d 个信号，实际得到 %d 个", tt.expectedCount, len(result))
			}

			if tt.expectedMerged && len(tt.predictions) > 1 && len(result) == 1 {
				merged := result[0]

				// 验证合并后的信号包含所有原始指标
				if len(merged.Indicators) < len(tt.predictions) {
					t.Errorf("合并后的信号应该包含所有原始指标，期望至少 %d 个，实际 %d 个",
						len(tt.predictions), len(merged.Indicators))
				}

				// 验证原因已被合并
				if merged.Reason == "" {
					t.Error("合并后的信号应该包含合并的原因")
				}

				t.Logf("✅ 合并成功:")
				t.Logf("   - 类型: %s", merged.Type)
				t.Logf("   - 价格: %s", merged.Price.Decimal.String())
				t.Logf("   - 日期: %s", merged.SignalDate)
				t.Logf("   - 概率: %s", merged.Probability.Decimal.String())
				t.Logf("   - 指标数量: %d", len(merged.Indicators))
				t.Logf("   - 指标列表: %v", merged.Indicators)
				t.Logf("   - 合并原因: %s", merged.Reason)
			}
		})
	}
}

// TestMergeMultipleSignals 测试多个信号合并的详细逻辑
func TestMergeMultipleSignals(t *testing.T) {
	service := NewPredictionService()

	t.Run("买入信号应使用最低价", func(t *testing.T) {
		signals := []models.TradingPointPrediction{
			{
				Type:        "BUY",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.65)),
				Indicators:  []string{"MACD"},
			},
			{
				Type:        "BUY",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.2)), // 最低价
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.70)),
				Indicators:  []string{"RSI"},
			},
			{
				Type:        "BUY",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.4)),
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.60)),
				Indicators:  []string{"BOLL"},
			},
		}

		merged := service.mergeMultipleSignals(signals)

		expectedPrice := decimal.NewFromFloat(10.2)
		if !merged.Price.Decimal.Equal(expectedPrice) {
			t.Errorf("买入信号应使用最低价 %s，实际为 %s",
				expectedPrice.String(), merged.Price.Decimal.String())
		}

		t.Logf("✅ 买入信号使用最低价: %s", merged.Price.Decimal.String())
	})

	t.Run("卖出信号应使用最高价", func(t *testing.T) {
		signals := []models.TradingPointPrediction{
			{
				Type:        "SELL",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.65)),
				Indicators:  []string{"MACD"},
			},
			{
				Type:        "SELL",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.9)), // 最高价
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.70)),
				Indicators:  []string{"RSI"},
			},
			{
				Type:        "SELL",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.7)),
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.60)),
				Indicators:  []string{"BOLL"},
			},
		}

		merged := service.mergeMultipleSignals(signals)

		expectedPrice := decimal.NewFromFloat(10.9)
		if !merged.Price.Decimal.Equal(expectedPrice) {
			t.Errorf("卖出信号应使用最高价 %s，实际为 %s",
				expectedPrice.String(), merged.Price.Decimal.String())
		}

		t.Logf("✅ 卖出信号使用最高价: %s", merged.Price.Decimal.String())
	})

	t.Run("多个信号共识应提升概率", func(t *testing.T) {
		signals := []models.TradingPointPrediction{
			{
				Type:        "BUY",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.65)),
				Indicators:  []string{"MACD"},
			},
			{
				Type:        "BUY",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.3)),
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.70)),
				Indicators:  []string{"RSI"},
			},
			{
				Type:        "BUY",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.4)),
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.60)),
				Indicators:  []string{"BOLL"},
			},
		}

		merged := service.mergeMultipleSignals(signals)

		// 计算平均概率
		avgProb := decimal.NewFromFloat(0.65).
			Add(decimal.NewFromFloat(0.70)).
			Add(decimal.NewFromFloat(0.60)).
			Div(decimal.NewFromFloat(3.0))

		// 合并后的概率应该高于平均值（因为共识提升）
		if !merged.Probability.Decimal.GreaterThan(avgProb) {
			t.Errorf("合并后的概率 %s 应该高于平均概率 %s（因为多个信号共识）",
				merged.Probability.Decimal.String(), avgProb.String())
		}

		t.Logf("✅ 平均概率: %s", avgProb.String())
		t.Logf("✅ 合并后概率: %s (提升了 %s)",
			merged.Probability.Decimal.String(),
			merged.Probability.Decimal.Sub(avgProb).String())
	})

	t.Run("合并后应包含所有指标", func(t *testing.T) {
		signals := []models.TradingPointPrediction{
			{
				Type:        "BUY",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.65)),
				Indicators:  []string{"MACD"},
			},
			{
				Type:        "BUY",
				Price:       models.NewJSONDecimal(decimal.NewFromFloat(10.3)),
				SignalDate:  "20240101",
				Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.70)),
				Indicators:  []string{"RSI", "KDJ"},
			},
		}

		merged := service.mergeMultipleSignals(signals)

		// 应该包含所有唯一的指标
		indicatorMap := make(map[string]bool)
		for _, ind := range merged.Indicators {
			indicatorMap[ind] = true
		}

		expectedIndicators := []string{"MACD", "RSI", "KDJ"}
		for _, expected := range expectedIndicators {
			if !indicatorMap[expected] {
				t.Errorf("合并后的信号应包含指标 %s", expected)
			}
		}

		t.Logf("✅ 合并后的指标: %v", merged.Indicators)
	})
}

// TestMergeSameDaySignalsWithBacktest 测试带回测数据的信号合并
func TestMergeSameDaySignalsWithBacktest(t *testing.T) {
	service := NewPredictionService()

	t.Run("回测数据应该被平均", func(t *testing.T) {
		signals := []models.TradingPointPrediction{
			{
				Type:           "BUY",
				Price:          models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
				SignalDate:     "20240101",
				Probability:    models.NewJSONDecimal(decimal.NewFromFloat(0.65)),
				Indicators:     []string{"MACD"},
				Backtested:     true,
				NextDayPrice:   models.NewJSONDecimal(decimal.NewFromFloat(10.8)),
				PriceDiff:      models.NewJSONDecimal(decimal.NewFromFloat(0.3)),
				PriceDiffRatio: models.NewJSONDecimal(decimal.NewFromFloat(2.86)),
				IsCorrect:      true,
			},
			{
				Type:           "BUY",
				Price:          models.NewJSONDecimal(decimal.NewFromFloat(10.3)),
				SignalDate:     "20240101",
				Probability:    models.NewJSONDecimal(decimal.NewFromFloat(0.70)),
				Indicators:     []string{"RSI"},
				Backtested:     true,
				NextDayPrice:   models.NewJSONDecimal(decimal.NewFromFloat(10.6)),
				PriceDiff:      models.NewJSONDecimal(decimal.NewFromFloat(0.3)),
				PriceDiffRatio: models.NewJSONDecimal(decimal.NewFromFloat(2.91)),
				IsCorrect:      true,
			},
		}

		merged := service.mergeMultipleSignals(signals)

		if !merged.Backtested {
			t.Error("合并后的信号应该标记为已回测")
		}

		// 验证回测数据是平均值
		expectedNextDayPrice := decimal.NewFromFloat(10.8).
			Add(decimal.NewFromFloat(10.6)).
			Div(decimal.NewFromFloat(2.0))

		if !merged.NextDayPrice.Decimal.Equal(expectedNextDayPrice) {
			t.Errorf("回测的次日价格应该是平均值 %s，实际为 %s",
				expectedNextDayPrice.String(), merged.NextDayPrice.Decimal.String())
		}

		if !merged.IsCorrect {
			t.Error("如果所有信号都正确，合并后的信号也应该正确")
		}

		t.Logf("✅ 回测数据已正确平均")
		t.Logf("   - 次日价格: %s", merged.NextDayPrice.Decimal.String())
		t.Logf("   - 价格差异: %s", merged.PriceDiff.Decimal.String())
		t.Logf("   - 预测正确: %v", merged.IsCorrect)
	})
}
