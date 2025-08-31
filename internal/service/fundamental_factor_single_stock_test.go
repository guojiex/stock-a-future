package service

import (
	"testing"

	"stock-a-future/internal/client"
	"stock-a-future/internal/models"

	"github.com/shopspring/decimal"
)

// TestCalculateFundamentalFactorWithScores 测试单个股票的基本面因子计算（包含得分）
func TestCalculateFundamentalFactorWithScores(t *testing.T) {
	// 创建模拟客户端
	mockClient := &client.MockDataSourceClient{
		DailyBasicData: &models.DailyBasic{
			TSCode:    "601669.SH",
			TradeDate: "20250831",
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(10.50)),
			Pe:        models.NewJSONDecimal(decimal.NewFromFloat(18.5)), // 良好的PE
			Pb:        models.NewJSONDecimal(decimal.NewFromFloat(2.1)),  // 良好的PB
			Ps:        models.NewJSONDecimal(decimal.NewFromFloat(1.8)),  // 优秀的PS
			DvRatio:   models.NewJSONDecimal(decimal.NewFromFloat(3.2)),
		},
		IncomeStatementData: &models.IncomeStatement{
			FinancialStatement: models.FinancialStatement{
				TSCode:  "601669.SH",
				EndDate: "20240930",
			},
			OperRevenue: models.NewJSONDecimal(decimal.NewFromFloat(1000000000)), // 10亿营收
			OperCost:    models.NewJSONDecimal(decimal.NewFromFloat(600000000)),  // 6亿成本
			OperProfit:  models.NewJSONDecimal(decimal.NewFromFloat(200000000)),  // 2亿营业利润
			NetProfit:   models.NewJSONDecimal(decimal.NewFromFloat(150000000)),  // 1.5亿净利润
		},
		BalanceSheetData: &models.BalanceSheet{
			FinancialStatement: models.FinancialStatement{
				TSCode:  "601669.SH",
				EndDate: "20240930",
			},
			TotalAssets:     models.NewJSONDecimal(decimal.NewFromFloat(5000000000)), // 50亿总资产
			TotalHldrEqy:    models.NewJSONDecimal(decimal.NewFromFloat(2000000000)), // 20亿所有者权益
			TotalLiab:       models.NewJSONDecimal(decimal.NewFromFloat(3000000000)), // 30亿负债
			TotalCurAssets:  models.NewJSONDecimal(decimal.NewFromFloat(1500000000)), // 15亿流动资产
			TotalCurLiab:    models.NewJSONDecimal(decimal.NewFromFloat(1000000000)), // 10亿流动负债
			InventoryAssets: models.NewJSONDecimal(decimal.NewFromFloat(300000000)),  // 3亿存货
			AccountsReceiv:  models.NewJSONDecimal(decimal.NewFromFloat(400000000)),  // 4亿应收账款
		},
		CashFlowData: &models.CashFlowStatement{
			FinancialStatement: models.FinancialStatement{
				TSCode:  "601669.SH",
				EndDate: "20240930",
			},
			NetCashOperAct: models.NewJSONDecimal(decimal.NewFromFloat(180000000)), // 1.8亿经营现金流
			NetCashInvAct:  models.NewJSONDecimal(decimal.NewFromFloat(-50000000)), // -0.5亿投资现金流
			NetCashFinAct:  models.NewJSONDecimal(decimal.NewFromFloat(-30000000)), // -0.3亿筹资现金流
		},
	}

	service := NewFundamentalFactorService(mockClient)

	// 测试计算基本面因子
	factor, err := service.CalculateFundamentalFactor("601669", "20250831")
	if err != nil {
		t.Fatalf("计算基本面因子失败: %v", err)
	}

	// 验证基本信息
	if factor.TSCode != "601669.SH" {
		t.Errorf("期望股票代码 601669.SH，实际得到 %s", factor.TSCode)
	}

	if factor.TradeDate != "20250831" {
		t.Errorf("期望交易日期 20250831，实际得到 %s", factor.TradeDate)
	}

	// 验证原始因子值
	t.Logf("原始因子值:")
	t.Logf("  PE: %v", factor.PE.Decimal)
	t.Logf("  PB: %v", factor.PB.Decimal)
	t.Logf("  PS: %v", factor.PS.Decimal)
	t.Logf("  ROE: %v", factor.ROE.Decimal)
	t.Logf("  ROA: %v", factor.ROA.Decimal)
	t.Logf("  毛利率: %v", factor.GrossMargin.Decimal)
	t.Logf("  净利率: %v", factor.NetMargin.Decimal)

	// 验证因子得分（这是修复的关键部分）
	t.Logf("因子得分:")
	t.Logf("  价值因子得分: %v", factor.ValueScore.Decimal)
	t.Logf("  成长因子得分: %v", factor.GrowthScore.Decimal)
	t.Logf("  质量因子得分: %v", factor.QualityScore.Decimal)
	t.Logf("  盈利因子得分: %v", factor.ProfitabilityScore.Decimal)
	t.Logf("  综合得分: %v", factor.CompositeScore.Decimal)

	// 验证因子得分不为零（修复前这些都是零）
	if factor.ValueScore.Decimal.IsZero() {
		t.Error("价值因子得分不应为零")
	}

	if factor.QualityScore.Decimal.IsZero() {
		t.Error("质量因子得分不应为零")
	}

	if factor.ProfitabilityScore.Decimal.IsZero() {
		t.Error("盈利因子得分不应为零")
	}

	if factor.CompositeScore.Decimal.IsZero() {
		t.Error("综合得分不应为零")
	}

	// 验证排名和分位数
	t.Logf("排名和分位数:")
	t.Logf("  市场排名: %d", factor.MarketRank)
	t.Logf("  行业排名: %d", factor.IndustryRank)
	t.Logf("  市场分位数: %v", factor.MarketPercentile.Decimal)
	t.Logf("  行业分位数: %v", factor.IndustryPercentile.Decimal)

	// 验证分位数为默认的50%
	expectedPercentile := decimal.NewFromFloat(50.0)
	if !factor.MarketPercentile.Decimal.Equal(expectedPercentile) {
		t.Errorf("期望市场分位数 %v，实际得到 %v", expectedPercentile, factor.MarketPercentile.Decimal)
	}

	if !factor.IndustryPercentile.Decimal.Equal(expectedPercentile) {
		t.Errorf("期望行业分位数 %v，实际得到 %v", expectedPercentile, factor.IndustryPercentile.Decimal)
	}

	// 验证更新时间
	if factor.UpdatedAt.IsZero() {
		t.Error("更新时间不应为零")
	}

	t.Log("✅ 单个股票基本面因子计算测试通过")
}

// TestFactorScoreCalculation 专门测试因子得分计算逻辑
func TestFactorScoreCalculation(t *testing.T) {
	service := NewFundamentalFactorService(nil)

	// 创建测试因子
	factor := &models.FundamentalFactor{
		TSCode:    "000001.SZ",
		TradeDate: "20250831",
		// 设置优秀的价值因子
		PE: models.NewJSONDecimal(decimal.NewFromFloat(12.0)), // 优秀
		PB: models.NewJSONDecimal(decimal.NewFromFloat(1.2)),  // 优秀
		PS: models.NewJSONDecimal(decimal.NewFromFloat(1.5)),  // 优秀
		// 设置良好的质量因子
		ROE:          models.NewJSONDecimal(decimal.NewFromFloat(18.0)), // 良好
		ROA:          models.NewJSONDecimal(decimal.NewFromFloat(8.0)),  // 良好
		DebtToAssets: models.NewJSONDecimal(decimal.NewFromFloat(45.0)), // 良好
		CurrentRatio: models.NewJSONDecimal(decimal.NewFromFloat(1.5)),  // 良好
		// 设置优秀的盈利因子
		GrossMargin:     models.NewJSONDecimal(decimal.NewFromFloat(45.0)), // 优秀
		NetMargin:       models.NewJSONDecimal(decimal.NewFromFloat(18.0)), // 优秀
		OperatingMargin: models.NewJSONDecimal(decimal.NewFromFloat(25.0)), // 优秀
	}

	// 计算因子得分
	service.calculateSingleStockScores(factor)

	// 验证价值因子得分（应该是优秀的，接近2.0）
	valueScore, _ := factor.ValueScore.Decimal.Float64()
	if valueScore < 1.5 {
		t.Errorf("价值因子得分应该较高，实际得到 %f", valueScore)
	}
	t.Logf("价值因子得分: %f", valueScore)

	// 验证质量因子得分（应该是良好的，接近1.0）
	qualityScore, _ := factor.QualityScore.Decimal.Float64()
	if qualityScore < 0.5 {
		t.Errorf("质量因子得分应该较高，实际得到 %f", qualityScore)
	}
	t.Logf("质量因子得分: %f", qualityScore)

	// 验证盈利因子得分（应该是优秀的，接近2.0）
	profitabilityScore, _ := factor.ProfitabilityScore.Decimal.Float64()
	if profitabilityScore < 1.5 {
		t.Errorf("盈利因子得分应该较高，实际得到 %f", profitabilityScore)
	}
	t.Logf("盈利因子得分: %f", profitabilityScore)

	// 验证综合得分
	compositeScore, _ := factor.CompositeScore.Decimal.Float64()
	if compositeScore <= 0 {
		t.Errorf("综合得分应该为正数，实际得到 %f", compositeScore)
	}
	t.Logf("综合得分: %f", compositeScore)

	t.Log("✅ 因子得分计算逻辑测试通过")
}

// TestEdgeCaseFactorScores 测试边界情况下的因子得分
func TestEdgeCaseFactorScores(t *testing.T) {
	service := NewFundamentalFactorService(nil)

	// 测试所有因子都为零的情况
	factor := &models.FundamentalFactor{
		TSCode:    "000002.SZ",
		TradeDate: "20250831",
		// 所有因子都为零
		PE:              models.NewJSONDecimal(decimal.Zero),
		PB:              models.NewJSONDecimal(decimal.Zero),
		PS:              models.NewJSONDecimal(decimal.Zero),
		ROE:             models.NewJSONDecimal(decimal.Zero),
		ROA:             models.NewJSONDecimal(decimal.Zero),
		GrossMargin:     models.NewJSONDecimal(decimal.Zero),
		NetMargin:       models.NewJSONDecimal(decimal.Zero),
		OperatingMargin: models.NewJSONDecimal(decimal.Zero),
	}

	// 计算因子得分
	service.calculateSingleStockScores(factor)

	// 验证所有得分都应该是零（因为没有有效数据）
	if !factor.ValueScore.Decimal.IsZero() {
		t.Errorf("价值因子得分应该为零，实际得到 %v", factor.ValueScore.Decimal)
	}

	if !factor.GrowthScore.Decimal.IsZero() {
		t.Errorf("成长因子得分应该为零，实际得到 %v", factor.GrowthScore.Decimal)
	}

	if !factor.QualityScore.Decimal.IsZero() {
		t.Errorf("质量因子得分应该为零，实际得到 %v", factor.QualityScore.Decimal)
	}

	if !factor.ProfitabilityScore.Decimal.IsZero() {
		t.Errorf("盈利因子得分应该为零，实际得到 %v", factor.ProfitabilityScore.Decimal)
	}

	if !factor.CompositeScore.Decimal.IsZero() {
		t.Errorf("综合得分应该为零，实际得到 %v", factor.CompositeScore.Decimal)
	}

	t.Log("✅ 边界情况因子得分测试通过")
}
