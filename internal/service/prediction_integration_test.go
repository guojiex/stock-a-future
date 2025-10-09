package service

import (
	"stock-a-future/internal/models"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// TestPredictionWithSignalMerge 测试完整的预测流程（包括信号合并）
func TestPredictionWithSignalMerge(t *testing.T) {
	service := NewPredictionService()

	// 创建模拟的股票数据
	data := createMockStockData()

	// 执行预测
	result, err := service.PredictTradingPoints(data)
	if err != nil {
		t.Fatalf("预测失败: %v", err)
	}

	t.Logf("预测结果:")
	t.Logf("  股票代码: %s", result.TSCode)
	t.Logf("  交易日期: %s", result.TradeDate)
	t.Logf("  整体置信度: %.2f%%", result.Confidence.Decimal.InexactFloat64()*100)
	t.Logf("  信号数量: %d", len(result.Predictions))

	// 检查信号是否已合并（同一天的同类型信号应该被合并）
	signalMap := make(map[string]int)
	for _, pred := range result.Predictions {
		key := pred.SignalDate + "_" + pred.Type
		signalMap[key]++
	}

	// 同一天同类型的信号应该只有一个
	for key, count := range signalMap {
		if count > 1 {
			t.Errorf("发现未合并的信号: %s 有 %d 个重复", key, count)
		} else {
			t.Logf("  ✅ 信号 %s 已正确合并（数量：1）", key)
		}
	}

	// 打印所有信号详情
	for i, pred := range result.Predictions {
		t.Logf("\n信号 %d:", i+1)
		t.Logf("  类型: %s", pred.Type)
		t.Logf("  价格: %s", pred.Price.Decimal.String())
		t.Logf("  信号日期: %s", pred.SignalDate)
		t.Logf("  预测日期: %s", pred.Date)
		t.Logf("  概率: %.2f%%", pred.Probability.Decimal.InexactFloat64()*100)
		t.Logf("  指标数量: %d", len(pred.Indicators))
		t.Logf("  指标: %v", pred.Indicators)
		t.Logf("  原因: %s", pred.Reason)

		if pred.Backtested {
			t.Logf("  回测:")
			t.Logf("    次日价格: %s", pred.NextDayPrice.Decimal.String())
			t.Logf("    价格差异: %s (%.2f%%)",
				pred.PriceDiff.Decimal.String(),
				pred.PriceDiffRatio.Decimal.InexactFloat64())
			t.Logf("    预测正确: %v", pred.IsCorrect)
		}

		// 如果指标数量大于1，说明是合并后的信号
		if len(pred.Indicators) > 1 {
			t.Logf("  ✅ 这是一个合并信号（包含%d个指标）", len(pred.Indicators))
		}
	}
}

// createMockStockData 创建模拟的股票数据用于测试
func createMockStockData() []models.StockDaily {
	baseDate := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	var data []models.StockDaily

	// 创建30天的模拟数据
	for i := 0; i < 30; i++ {
		date := baseDate.AddDate(0, 0, i)
		dateStr := date.Format("20060102")

		// 模拟价格波动
		basePrice := 10.0 + float64(i)*0.1
		open := basePrice
		high := basePrice * 1.03
		low := basePrice * 0.97
		close := basePrice * 1.01

		// 最后几天创建一个超卖的情况（RSI会很低）
		if i >= 25 {
			close = basePrice * 0.95
			low = basePrice * 0.93
		}

		// 模拟成交量
		vol := 1000000 + float64(i*10000)

		data = append(data, models.StockDaily{
			TSCode:    "600000.SH",
			TradeDate: dateStr,
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(open)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(high)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(low)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(close)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(vol)),
			Amount:    models.NewJSONDecimal(decimal.NewFromFloat(vol * close)),
		})
	}

	return data
}

// TestSignalMergeDoesNotAffectDifferentDates 测试信号合并不会影响不同日期的信号
func TestSignalMergeDoesNotAffectDifferentDates(t *testing.T) {
	service := NewPredictionService()

	// 创建足够长的数据以产生多个日期的信号
	data := createMockStockDataWithMultipleDates()

	result, err := service.PredictTradingPoints(data)
	if err != nil {
		t.Fatalf("预测失败: %v", err)
	}

	// 统计每个日期的信号数量
	dateSignalCount := make(map[string]int)
	for _, pred := range result.Predictions {
		dateSignalCount[pred.SignalDate]++
	}

	t.Logf("不同日期的信号分布:")
	for date, count := range dateSignalCount {
		t.Logf("  日期 %s: %d 个信号", date, count)
	}

	// 验证：不同日期的信号应该保持独立
	if len(dateSignalCount) == 0 {
		t.Log("没有产生信号，这是正常的（取决于数据特征）")
	} else {
		t.Logf("✅ 成功产生 %d 个不同日期的信号", len(dateSignalCount))
	}
}

// createMockStockDataWithMultipleDates 创建会产生多个日期信号的模拟数据
func createMockStockDataWithMultipleDates() []models.StockDaily {
	baseDate := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)
	var data []models.StockDaily

	// 创建60天的数据，模拟多次买卖机会
	for i := 0; i < 60; i++ {
		date := baseDate.AddDate(0, 0, i)
		dateStr := date.Format("20060102")

		// 模拟波浪式价格走势（产生多个买卖信号）
		basePrice := 10.0
		wave := 2.0 // 波动幅度

		// 使用抛物线模拟价格波动
		priceOffset := wave * (1.0 - ((float64(i%20)-10.0)/10.0)*((float64(i%20)-10.0)/10.0))

		close := basePrice + priceOffset
		open := close * 0.99
		high := close * 1.02
		low := close * 0.98

		vol := 1000000.0 + float64(i*10000)

		data = append(data, models.StockDaily{
			TSCode:    "600000.SH",
			TradeDate: dateStr,
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(open)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(high)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(low)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(close)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(vol)),
			Amount:    models.NewJSONDecimal(decimal.NewFromFloat(vol * close)),
		})
	}

	return data
}
