package service

import (
	"testing"

	"stock-a-future/internal/client"
)

// TestCalculateFundamentalFactorTiming 测试基本面因子计算的时间打印功能
func TestCalculateFundamentalFactorTiming(t *testing.T) {
	// 创建AKTools客户端
	aktoolsClient := client.NewAKToolsClient("http://127.0.0.1:8080")

	// 创建基本面因子服务
	service := NewFundamentalFactorService(aktoolsClient)

	// 测试股票代码和日期
	symbol := "600976" // 健民集团
	tradeDate := "20240831"

	t.Logf("🧪 测试基本面因子计算时间打印功能")
	t.Logf("📊 测试股票: %s, 交易日期: %s", symbol, tradeDate)

	// 调用计算函数 - 这将打印总耗时
	factor, err := service.CalculateFundamentalFactor(symbol, tradeDate)

	if err != nil {
		t.Logf("⚠️  计算失败: %v", err)
		// 不要因为网络问题而失败测试
		t.Skipf("跳过测试，可能是AKTools服务不可用: %v", err)
		return
	}

	if factor == nil {
		t.Fatalf("❌ 返回的因子为nil")
	}

	t.Logf("✅ 基本面因子计算成功")
	t.Logf("📈 股票代码: %s", factor.TSCode)
	t.Logf("📅 交易日期: %s", factor.TradeDate)
	t.Logf("💰 PE: %v", factor.PE.Decimal)
	t.Logf("📊 PB: %v", factor.PB.Decimal)
	t.Logf("🏭 ROE: %v", factor.ROE.Decimal)
	t.Logf("💼 ROA: %v", factor.ROA.Decimal)

	// 注意：实际的耗时会在日志中显示，格式如：
	// [FundamentalFactorService] 基本面因子计算完成: 600976, 总耗时: 1.234567s
}

// TestBatchCalculateFundamentalFactorsTiming 测试批量计算的时间打印
func TestBatchCalculateFundamentalFactorsTiming(t *testing.T) {
	// 创建AKTools客户端
	aktoolsClient := client.NewAKToolsClient("http://127.0.0.1:8080")

	// 创建基本面因子服务
	service := NewFundamentalFactorService(aktoolsClient)

	// 测试多个股票代码
	symbols := []string{"600976", "000001", "002415"} // 健民集团、平安银行、海康威视
	tradeDate := "20240831"

	t.Logf("🧪 测试批量基本面因子计算时间打印功能")
	t.Logf("📊 测试股票数量: %d", len(symbols))

	// 调用批量计算函数 - 这将为每个股票打印单独的耗时
	factors, err := service.BatchCalculateFundamentalFactors(symbols, tradeDate)

	if err != nil {
		t.Logf("⚠️  批量计算失败: %v", err)
		t.Skipf("跳过测试，可能是AKTools服务不可用: %v", err)
		return
	}

	t.Logf("✅ 批量计算完成，成功处理: %d 个股票", len(factors))

	// 打印每个股票的基本信息
	for i, factor := range factors {
		t.Logf("📈 股票 %d: %s, PE: %v, PB: %v", i+1, factor.TSCode, factor.PE.Decimal, factor.PB.Decimal)
	}

	// 注意：每个股票的计算耗时会在日志中单独显示
}
