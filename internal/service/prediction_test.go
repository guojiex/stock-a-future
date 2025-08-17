package service

import (
	"stock-a-future/internal/models"
	"testing"

	"github.com/shopspring/decimal"
)

func TestSortPredictionsByConfidenceAndStrength(t *testing.T) {
	service := NewPredictionService()

	// 创建测试数据
	predictions := []models.TradingPointPrediction{
		{
			Type:        "BUY",
			Reason:      "识别到锤子线模式，置信度：45.2，强度：WEAK",
			Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.78)),
		},
		{
			Type:        "BUY",
			Reason:      "识别到双响炮模式，置信度：85.5，强度：STRONG",
			Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.92)),
		},
		{
			Type:        "BUY",
			Reason:      "识别到红三兵模式，置信度：78.2，强度：STRONG",
			Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.90)),
		},
		{
			Type:        "BUY",
			Reason:      "识别到放量突破模式，置信度：72.1，强度：MEDIUM",
			Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.88)),
		},
		{
			Type:        "BUY",
			Reason:      "识别到量价齐升模式，置信度：65.8，强度：MEDIUM",
			Probability: models.NewJSONDecimal(decimal.NewFromFloat(0.84)),
		},
	}

	// 执行排序
	service.sortPredictionsByConfidenceAndStrength(predictions)

	// 验证排序结果
	expectedOrder := []string{
		"STRONG", // 双响炮
		"STRONG", // 红三兵
		"MEDIUM", // 放量突破
		"MEDIUM", // 量价齐升
		"WEAK",   // 锤子线
	}

	for i, prediction := range predictions {
		strength := service.extractStrengthFromReason(prediction.Reason)
		if strength != expectedOrder[i] {
			t.Errorf("位置 %d: 期望强度 %s, 实际强度 %s", i, expectedOrder[i], strength)
		}
	}

	// 验证STRONG组内的置信度排序
	strongPredictions := []models.TradingPointPrediction{}
	for _, pred := range predictions {
		if service.extractStrengthFromReason(pred.Reason) == "STRONG" {
			strongPredictions = append(strongPredictions, pred)
		}
	}

	if len(strongPredictions) >= 2 {
		confidence1 := service.extractConfidenceFromReason(strongPredictions[0].Reason)
		confidence2 := service.extractConfidenceFromReason(strongPredictions[1].Reason)
		if confidence1.LessThan(confidence2) {
			t.Errorf("STRONG组内置信度排序错误: 第一个置信度 %s 应该大于第二个 %s",
				confidence1.String(), confidence2.String())
		}
	}

	t.Logf("排序后的预测结果:")
	for i, pred := range predictions {
		strength := service.extractStrengthFromReason(pred.Reason)
		confidence := service.extractConfidenceFromReason(pred.Reason)
		t.Logf("%d. 强度: %s, 置信度: %s, 概率: %s",
			i+1, strength, confidence.String(), pred.Probability.Decimal.String())
	}
}

func TestExtractStrengthFromReason(t *testing.T) {
	service := NewPredictionService()

	testCases := []struct {
		reason   string
		expected string
	}{
		{"识别到双响炮模式，置信度：85.5，强度：STRONG", "STRONG"},
		{"识别到红三兵模式，置信度：78.2，强度：MEDIUM", "MEDIUM"},
		{"识别到锤子线模式，置信度：45.2，强度：WEAK", "WEAK"},
		{"MACD金叉信号，DIF线上穿DEA线", "WEAK"}, // 默认值
		{"", "WEAK"}, // 空字符串
	}

	for _, tc := range testCases {
		result := service.extractStrengthFromReason(tc.reason)
		if result != tc.expected {
			t.Errorf("提取强度失败: 输入 '%s', 期望 '%s', 实际 '%s'",
				tc.reason, tc.expected, result)
		}
	}
}

func TestExtractConfidenceFromReason(t *testing.T) {
	service := NewPredictionService()

	testCases := []struct {
		reason   string
		expected string
	}{
		{"识别到双响炮模式，置信度：85.5，强度：STRONG", "85.5"},
		{"识别到红三兵模式，置信度：78.2，强度：MEDIUM", "78.2"},
		{"识别到锤子线模式，置信度：45.2，强度：WEAK", "45.2"},
		{"MACD金叉信号，DIF线上穿DEA线", "0"}, // 默认值
		{"", "0"}, // 空字符串
	}

	for _, tc := range testCases {
		result := service.extractConfidenceFromReason(tc.reason)
		expected, _ := decimal.NewFromString(tc.expected)
		if !result.Equal(expected) {
			t.Errorf("提取置信度失败: 输入 '%s', 期望 '%s', 实际 '%s'",
				tc.reason, expected.String(), result.String())
		}
	}
}
