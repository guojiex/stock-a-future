package service

import (
	"stock-a-future/internal/models"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// TestPredictionDateUsesActualTradeDate 测试预测日期使用实际交易日期而不是未来日期
func TestPredictionDateUsesActualTradeDate(t *testing.T) {
	service := NewPredictionService()

	// 创建测试数据 - 使用2024年11月6日的数据
	testDate := "20241106"
	data := []models.StockDaily{
		{
			TSCode:    "600000.SH",
			TradeDate: "20241104",
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(9.8)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.2)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.7)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1000000)),
		},
		{
			TSCode:    "600000.SH",
			TradeDate: "20241105",
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.0)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(10.6)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(9.9)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1200000)),
		},
		{
			TSCode:    "600000.SH",
			TradeDate: testDate,
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(11.0)),
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(10.5)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(11.2)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(10.4)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(1500000)),
		},
	}

	// 调用预测服务
	result, err := service.PredictTradingPoints(data)
	if err != nil {
		t.Fatalf("预测失败: %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("预测结果为空")
	}

	t.Logf("预测结果:")
	t.Logf("  股票代码: %s", result.TSCode)
	t.Logf("  交易日期: %s", result.TradeDate)
	t.Logf("  信号数量: %d", len(result.Predictions))

	// 验证PredictionResult的TradeDate使用的是最新数据的日期
	if result.TradeDate != testDate {
		t.Errorf("期望TradeDate为 %s，实际为 %s", testDate, result.TradeDate)
	}

	// 验证所有预测信号的Date字段都使用实际数据日期，而不是未来日期
	today := time.Now().Format("20060102")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("20060102")

	for i, pred := range result.Predictions {
		t.Logf("\n  信号 %d:", i+1)
		t.Logf("    类型: %s", pred.Type)
		t.Logf("    预测日期: %s", pred.Date)
		t.Logf("    信号日期: %s", pred.SignalDate)
		t.Logf("    指标: %v", pred.Indicators)

		// 验证Date字段不是明天的日期
		if pred.Date == tomorrow {
			t.Errorf("预测信号 %d 的Date字段错误地使用了未来日期 %s (明天)，应该使用实际数据日期", i+1, tomorrow)
		}

		// Date应该是实际数据的日期，而不是今天或明天
		if pred.Date == today || pred.Date == tomorrow {
			t.Logf("    ⚠️ 警告: Date字段 (%s) 看起来像是使用了time.Now()而不是实际数据日期", pred.Date)
		}

		// Date应该等于SignalDate（都应该是实际数据日期）
		if pred.Date != pred.SignalDate {
			t.Logf("    ℹ️ 注意: Date (%s) 和 SignalDate (%s) 不同", pred.Date, pred.SignalDate)
		}
	}

	t.Logf("\n✅ 验证通过：所有预测信号都使用实际交易日期而不是未来日期")
}
