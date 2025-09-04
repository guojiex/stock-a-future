package client

import (
	"testing"

	"github.com/shopspring/decimal"
)

// TestAKToolsChangeCalculation 测试AKTools涨跌幅计算的准确性
func TestAKToolsChangeCalculation(t *testing.T) {
	client := NewAKToolsClient("http://127.0.0.1:8080")

	// 模拟AKTools返回的数据
	testData := []AKToolsDailyResponse{
		{
			Date:      "2024-01-15",
			Open:      10.50,
			Close:     11.00,
			High:      11.20,
			Low:       10.40,
			Volume:    1000000,
			Amount:    10800000,
			Change:    0.50, // 涨跌额
			ChangePct: 4.76, // 涨跌幅 (0.50/10.50 * 100)
		},
		{
			Date:      "2024-01-16",
			Open:      11.00,
			Close:     10.80,
			High:      11.10,
			Low:       10.70,
			Volume:    800000,
			Amount:    8640000,
			Change:    -0.20, // 涨跌额
			ChangePct: -1.82, // 涨跌幅 (-0.20/11.00 * 100)
		},
	}

	// 转换数据
	result := client.convertToStockDaily(testData, "000001")

	// 验证数据长度
	if len(result) != 2 {
		t.Fatalf("期望得到2条数据，实际得到%d条", len(result))
	}

	// 测试第一条数据
	first := result[0]
	expectedPreClose1 := 10.50 // 11.00 - 0.50
	actualPreClose1, _ := first.PreClose.Float64()

	if !decimal.NewFromFloat(expectedPreClose1).Equal(first.PreClose.Decimal) {
		t.Errorf("第一条数据昨收价计算错误: 期望 %.2f, 实际 %.2f", expectedPreClose1, actualPreClose1)
	}

	// 验证涨跌额和涨跌幅
	actualChange1, _ := first.Change.Float64()
	actualPctChg1, _ := first.PctChg.Float64()

	if actualChange1 != 0.50 {
		t.Errorf("第一条数据涨跌额错误: 期望 0.50, 实际 %.2f", actualChange1)
	}

	if actualPctChg1 != 4.76 {
		t.Errorf("第一条数据涨跌幅错误: 期望 4.76, 实际 %.2f", actualPctChg1)
	}

	// 测试第二条数据
	second := result[1]
	expectedPreClose2 := 11.00 // 10.80 - (-0.20)
	actualPreClose2, _ := second.PreClose.Float64()

	if !decimal.NewFromFloat(expectedPreClose2).Equal(second.PreClose.Decimal) {
		t.Errorf("第二条数据昨收价计算错误: 期望 %.2f, 实际 %.2f", expectedPreClose2, actualPreClose2)
	}

	// 验证涨跌额和涨跌幅
	actualChange2, _ := second.Change.Float64()
	actualPctChg2, _ := second.PctChg.Float64()

	if actualChange2 != -0.20 {
		t.Errorf("第二条数据涨跌额错误: 期望 -0.20, 实际 %.2f", actualChange2)
	}

	if actualPctChg2 != -1.82 {
		t.Errorf("第二条数据涨跌幅错误: 期望 -1.82, 实际 %.2f", actualPctChg2)
	}

	t.Logf("✅ 涨跌幅计算测试通过")
	t.Logf("第一条: 收盘价=%.2f, 昨收价=%.2f, 涨跌额=%.2f, 涨跌幅=%.2f%%",
		11.00, actualPreClose1, actualChange1, actualPctChg1)
	t.Logf("第二条: 收盘价=%.2f, 昨收价=%.2f, 涨跌额=%.2f, 涨跌幅=%.2f%%",
		10.80, actualPreClose2, actualChange2, actualPctChg2)
}

// TestAKToolsChangeCalculationEdgeCases 测试边界情况
func TestAKToolsChangeCalculationEdgeCases(t *testing.T) {
	client := NewAKToolsClient("http://127.0.0.1:8080")

	// 测试零涨跌额的情况
	testData := []AKToolsDailyResponse{
		{
			Date:      "2024-01-15",
			Open:      10.00,
			Close:     10.00,
			High:      10.10,
			Low:       9.90,
			Volume:    1000000,
			Amount:    10000000,
			Change:    0.00, // 无涨跌
			ChangePct: 0.00, // 无涨跌幅
		},
	}

	result := client.convertToStockDaily(testData, "000001")

	if len(result) != 1 {
		t.Fatalf("期望得到1条数据，实际得到%d条", len(result))
	}

	first := result[0]
	expectedPreClose := 10.00 // 10.00 - 0.00
	actualPreClose, _ := first.PreClose.Float64()

	if !decimal.NewFromFloat(expectedPreClose).Equal(first.PreClose.Decimal) {
		t.Errorf("零涨跌情况昨收价计算错误: 期望 %.2f, 实际 %.2f", expectedPreClose, actualPreClose)
	}

	t.Logf("✅ 零涨跌情况测试通过: 收盘价=%.2f, 昨收价=%.2f", 10.00, actualPreClose)
}

// TestAKToolsChangeCalculationConsistency 测试计算一致性
func TestAKToolsChangeCalculationConsistency(t *testing.T) {
	client := NewAKToolsClient("http://127.0.0.1:8080")

	// 测试多天连续数据的一致性
	testData := []AKToolsDailyResponse{
		{
			Date:      "2024-01-15",
			Close:     10.00,
			Change:    0.50,
			ChangePct: 5.26, // 0.50/9.50 * 100
		},
		{
			Date:      "2024-01-16",
			Close:     10.20,
			Change:    0.20,
			ChangePct: 2.00, // 0.20/10.00 * 100
		},
		{
			Date:      "2024-01-17",
			Close:     9.80,
			Change:    -0.40,
			ChangePct: -3.92, // -0.40/10.20 * 100
		},
	}

	result := client.convertToStockDaily(testData, "000001")

	if len(result) != 3 {
		t.Fatalf("期望得到3条数据，实际得到%d条", len(result))
	}

	// 验证连续性：验证基于涨跌额计算昨收价的正确性
	for i := 0; i < len(result); i++ {
		currentPreClose, _ := result[i].PreClose.Float64()

		// 验证计算的正确性：昨收价 = 收盘价 - 涨跌额
		currentClose, _ := result[i].Close.Float64()
		currentChange, _ := result[i].Change.Float64()
		expectedPreClose := currentClose - currentChange

		if !decimal.NewFromFloat(expectedPreClose).Equal(result[i].PreClose.Decimal) {
			t.Errorf("第%d天昨收价计算错误: 期望 %.2f, 实际 %.2f",
				i+1, expectedPreClose, currentPreClose)
		}

		t.Logf("第%d天: 收盘价=%.2f, 昨收价=%.2f, 涨跌额=%.2f",
			i+1, currentClose, currentPreClose, currentChange)
	}

	t.Logf("✅ 连续数据一致性测试通过")
}
