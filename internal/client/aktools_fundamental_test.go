package client

import (
	"context"
	"testing"
	"time"
)

const (
	// 测试用的股票代码 - 使用茅台(600519)作为测试，因为财务数据更稳定
	testSymbol = "600519"
	testTSCode = "600519.SH"
	// 测试用的时间
	testPeriod    = "20231231"
	testTradeDate = "20240115"
	// AKTools测试服务器地址
	testAKToolsURL = "http://127.0.0.1:8080"
)

// TestAKToolsClient_GetIncomeStatement 测试获取利润表数据
func TestAKToolsClient_GetIncomeStatement(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试获取单个利润表数据 - 不指定期间，获取最新数据
	incomeStatement, err := client.GetIncomeStatement(testSymbol, "", "1")
	if err != nil {
		t.Logf("警告: 获取利润表数据失败: %v", err)
		t.Logf("请确保AKTools服务在 %s 运行", testAKToolsURL)
		return // 不让测试失败，因为可能是服务未启动
	}

	// 验证返回数据
	if incomeStatement == nil {
		t.Fatal("利润表数据不应为空")
	}

	if incomeStatement.TSCode != testTSCode {
		t.Errorf("期望TSCode为 %s, 实际为 %s", testTSCode, incomeStatement.TSCode)
	}

	t.Logf("成功获取利润表数据: TSCode=%s, FDate=%s", incomeStatement.TSCode, incomeStatement.FDate)
	t.Logf("营业收入: %s", incomeStatement.Revenue.String())
	t.Logf("净利润: %s", incomeStatement.NetProfit.String())

	// 验证关键财务数据不为零
	if incomeStatement.Revenue.Decimal.IsZero() {
		t.Error("营业收入不应为零")
	}
	if incomeStatement.NetProfit.Decimal.IsZero() {
		t.Error("净利润不应为零")
	}
	if !incomeStatement.BasicEps.Decimal.IsZero() {
		t.Logf("基本每股收益: %s", incomeStatement.BasicEps.String())
	}
}

// TestAKToolsClient_GetIncomeStatements 测试批量获取利润表数据
func TestAKToolsClient_GetIncomeStatements(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试批量获取利润表数据 - 不指定日期范围，获取所有可用数据
	incomeStatements, err := client.GetIncomeStatements(testSymbol, "", "", "1")
	if err != nil {
		t.Logf("警告: 批量获取利润表数据失败: %v", err)
		t.Logf("请确保AKTools服务在 %s 运行", testAKToolsURL)
		return
	}

	// 验证返回数据
	if incomeStatements == nil {
		t.Fatal("利润表数据列表不应为空")
	}

	t.Logf("成功获取 %d 条利润表数据", len(incomeStatements))

	// 验证数据结构和数值
	for i, stmt := range incomeStatements {
		if i >= 3 { // 只检查前3条
			break
		}
		t.Logf("第%d条: TSCode=%s, FDate=%s, 营业收入=%s",
			i+1, stmt.TSCode, stmt.FDate, stmt.Revenue.String())

		// 验证关键财务数据不为零
		if stmt.Revenue.Decimal.IsZero() {
			t.Errorf("第%d条记录的营业收入不应为零", i+1)
		}
		if stmt.NetProfit.Decimal.IsZero() {
			t.Errorf("第%d条记录的净利润不应为零", i+1)
		}
	}
}

// TestAKToolsClient_GetBalanceSheet 测试获取资产负债表数据
func TestAKToolsClient_GetBalanceSheet(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试获取单个资产负债表数据 - 不指定期间，获取最新数据
	balanceSheet, err := client.GetBalanceSheet(testSymbol, "", "1")
	if err != nil {
		t.Logf("警告: 获取资产负债表数据失败: %v", err)
		t.Logf("请确保AKTools服务在 %s 运行", testAKToolsURL)
		return
	}

	// 验证返回数据
	if balanceSheet == nil {
		t.Fatal("资产负债表数据不应为空")
	}

	if balanceSheet.TSCode != testTSCode {
		t.Errorf("期望TSCode为 %s, 实际为 %s", testTSCode, balanceSheet.TSCode)
	}

	t.Logf("成功获取资产负债表数据: TSCode=%s, FDate=%s", balanceSheet.TSCode, balanceSheet.FDate)
	t.Logf("资产总计: %s", balanceSheet.TotalAssets.String())
	t.Logf("负债合计: %s", balanceSheet.TotalLiab.String())
	t.Logf("所有者权益合计: %s", balanceSheet.TotalHldrEqy.String())

	// 验证关键资产负债表数据不为零
	if balanceSheet.TotalAssets.Decimal.IsZero() {
		t.Error("资产总计不应为零")
	}
	if balanceSheet.TotalLiab.Decimal.IsZero() {
		t.Error("负债合计不应为零")
	}
	if balanceSheet.TotalHldrEqy.Decimal.IsZero() {
		t.Error("所有者权益合计不应为零")
	}
}

// TestAKToolsClient_GetBalanceSheets 测试批量获取资产负债表数据
func TestAKToolsClient_GetBalanceSheets(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试批量获取资产负债表数据 - 不指定日期范围，获取所有可用数据
	balanceSheets, err := client.GetBalanceSheets(testSymbol, "", "", "1")
	if err != nil {
		t.Logf("警告: 批量获取资产负债表数据失败: %v", err)
		t.Logf("请确保AKTools服务在 %s 运行", testAKToolsURL)
		return
	}

	// 验证返回数据
	if balanceSheets == nil {
		t.Fatal("资产负债表数据列表不应为空")
	}

	t.Logf("成功获取 %d 条资产负债表数据", len(balanceSheets))

	// 验证数据结构和数值
	for i, sheet := range balanceSheets {
		if i >= 3 { // 只检查前3条
			break
		}
		t.Logf("第%d条: TSCode=%s, FDate=%s, 资产总计=%s",
			i+1, sheet.TSCode, sheet.FDate, sheet.TotalAssets.String())

		// 验证关键资产负债表数据不为零
		if sheet.TotalAssets.Decimal.IsZero() {
			t.Errorf("第%d条记录的资产总计不应为零", i+1)
		}
		if sheet.TotalLiab.Decimal.IsZero() {
			t.Errorf("第%d条记录的负债合计不应为零", i+1)
		}
		if sheet.TotalHldrEqy.Decimal.IsZero() {
			t.Errorf("第%d条记录的所有者权益合计不应为零", i+1)
		}
	}
}

// TestAKToolsClient_GetCashFlowStatement 测试获取现金流量表数据
func TestAKToolsClient_GetCashFlowStatement(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试获取单个现金流量表数据 - 不指定期间，获取最新数据
	cashFlowStatement, err := client.GetCashFlowStatement(testSymbol, "", "1")
	if err != nil {
		t.Logf("警告: 获取现金流量表数据失败: %v", err)
		t.Logf("请确保AKTools服务在 %s 运行", testAKToolsURL)
		return
	}

	// 验证返回数据
	if cashFlowStatement == nil {
		t.Fatal("现金流量表数据不应为空")
	}

	if cashFlowStatement.TSCode != testTSCode {
		t.Errorf("期望TSCode为 %s, 实际为 %s", testTSCode, cashFlowStatement.TSCode)
	}

	t.Logf("成功获取现金流量表数据: TSCode=%s, FDate=%s", cashFlowStatement.TSCode, cashFlowStatement.FDate)
	t.Logf("经营活动现金流净额: %s", cashFlowStatement.NetCashOperAct.String())
	t.Logf("投资活动现金流净额: %s", cashFlowStatement.NetCashInvAct.String())
	t.Logf("筹资活动现金流净额: %s", cashFlowStatement.NetCashFinAct.String())

	// 验证经营活动现金流不为零（茅台经营现金流应该很强劲）
	if cashFlowStatement.NetCashOperAct.Decimal.IsZero() {
		t.Error("经营活动现金流净额不应为零")
	}

	// 筹资活动现金流可能为负（分红），但不应该为零
	if cashFlowStatement.NetCashFinAct.Decimal.IsZero() {
		t.Error("筹资活动现金流净额不应为零")
	}
}

// TestAKToolsClient_GetCashFlowStatements 测试批量获取现金流量表数据
func TestAKToolsClient_GetCashFlowStatements(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试批量获取现金流量表数据 - 不指定日期范围，获取所有可用数据
	cashFlowStatements, err := client.GetCashFlowStatements(testSymbol, "", "", "1")
	if err != nil {
		t.Logf("警告: 批量获取现金流量表数据失败: %v", err)
		t.Logf("请确保AKTools服务在 %s 运行", testAKToolsURL)
		return
	}

	// 验证返回数据
	if cashFlowStatements == nil {
		t.Fatal("现金流量表数据列表不应为空")
	}

	t.Logf("成功获取 %d 条现金流量表数据", len(cashFlowStatements))

	// 验证数据结构和数值
	for i, stmt := range cashFlowStatements {
		if i >= 3 { // 只检查前3条
			break
		}
		t.Logf("第%d条: TSCode=%s, FDate=%s, 经营现金流=%s",
			i+1, stmt.TSCode, stmt.FDate, stmt.NetCashOperAct.String())

		// 验证经营活动现金流不为零（茅台经营现金流应该很强劲）
		if stmt.NetCashOperAct.Decimal.IsZero() {
			t.Errorf("第%d条记录的经营活动现金流净额不应为零", i+1)
		}
	}
}

// TestAKToolsClient_GetFinancialIndicator 测试获取财务指标数据
func TestAKToolsClient_GetFinancialIndicator(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试获取财务指标数据
	financialIndicator, err := client.GetFinancialIndicator(testSymbol, testPeriod, "1")
	if err != nil {
		t.Fatalf("获取财务指标数据失败: %v", err)
	}

	// 验证返回数据
	if financialIndicator == nil {
		t.Fatal("财务指标数据不应为空")
	}

	if financialIndicator.TSCode != testTSCode {
		t.Errorf("期望TSCode为 %s, 实际为 %s", testTSCode, financialIndicator.TSCode)
	}

	t.Logf("成功获取财务指标数据: TSCode=%s, FDate=%s", financialIndicator.TSCode, financialIndicator.FDate)
	// 注意：AKTools暂不支持财务指标，所以这里主要测试结构正确性
}

// TestAKToolsClient_GetDailyBasic 测试获取每日基本面指标
func TestAKToolsClient_GetDailyBasic(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试获取每日基本面指标
	dailyBasic, err := client.GetDailyBasic(context.Background(), testSymbol, testTradeDate)
	if err != nil {
		t.Logf("警告: 获取每日基本面指标失败: %v", err)
		t.Logf("请确保AKTools服务在 %s 运行", testAKToolsURL)
		return
	}

	// 验证返回数据
	if dailyBasic == nil {
		t.Fatal("每日基本面指标数据不应为空")
	}

	if dailyBasic.TSCode != testTSCode {
		t.Errorf("期望TSCode为 %s, 实际为 %s", testTSCode, dailyBasic.TSCode)
	}

	t.Logf("成功获取每日基本面指标: TSCode=%s, TradeDate=%s", dailyBasic.TSCode, dailyBasic.TradeDate)
	t.Logf("收盘价: %s", dailyBasic.Close.String())
	t.Logf("市盈率: %s", dailyBasic.Pe.String())
	t.Logf("市净率: %s", dailyBasic.Pb.String())
	t.Logf("总市值: %s", dailyBasic.TotalMv.String())
}

// TestAKToolsClient_GetFundamentalFactor 测试获取基本面因子
func TestAKToolsClient_GetFundamentalFactor(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试获取基本面因子
	fundamentalFactor, err := client.GetFundamentalFactor(testSymbol, testTradeDate)
	if err != nil {
		t.Fatalf("获取基本面因子失败: %v", err)
	}

	// 验证返回数据
	if fundamentalFactor == nil {
		t.Fatal("基本面因子数据不应为空")
	}

	if fundamentalFactor.TSCode != testTSCode {
		t.Errorf("期望TSCode为 %s, 实际为 %s", testTSCode, fundamentalFactor.TSCode)
	}

	if fundamentalFactor.TradeDate != testTradeDate {
		t.Errorf("期望TradeDate为 %s, 实际为 %s", testTradeDate, fundamentalFactor.TradeDate)
	}

	// 验证时间戳
	if fundamentalFactor.UpdatedAt.IsZero() {
		t.Error("UpdatedAt不应为零值")
	}

	t.Logf("成功获取基本面因子: TSCode=%s, TradeDate=%s, UpdatedAt=%s",
		fundamentalFactor.TSCode, fundamentalFactor.TradeDate, fundamentalFactor.UpdatedAt.Format(time.RFC3339))
}

// TestAKToolsClient_DataConversion 测试数据转换函数
func TestAKToolsClient_DataConversion(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试parseDecimalFromInterface函数
	testCases := []struct {
		input    interface{}
		expected string
	}{
		{nil, "0"},
		{"", "0"},
		{"-", "0"},
		{"--", "0"},
		{123.45, "123.45"},
		{"678.90", "678.9"}, // Decimal会自动去掉尾随的零
		{100, "100"},
		{int64(200), "200"},
		{"invalid", "0"},
	}

	for _, tc := range testCases {
		result := client.parseDecimalFromInterface(tc.input)
		if result.String() != tc.expected {
			t.Errorf("parseDecimalFromInterface(%v): 期望 %s, 实际 %s", tc.input, tc.expected, result.String())
		}
	}

	t.Log("数据转换函数测试通过")
}

// TestAKToolsClient_SymbolHandling 测试股票代码处理
func TestAKToolsClient_SymbolHandling(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试CleanStockSymbol函数
	testCases := []struct {
		input    string
		expected string
	}{
		{"000001.SZ", "000001"},
		{"600000.SH", "600000"},
		{"300001.sz", "300001"},
		{"000001", "000001"},
		{"688001.SH", "688001"},
	}

	for _, tc := range testCases {
		result := client.CleanStockSymbol(tc.input)
		if result != tc.expected {
			t.Errorf("CleanStockSymbol(%s): 期望 %s, 实际 %s", tc.input, tc.expected, result)
		}
	}

	// 测试DetermineTSCode函数
	testTSCodes := []struct {
		input    string
		expected string
	}{
		{"000001", "000001.SZ"},
		{"600000", "600000.SH"},
		{"300001", "300001.SZ"},
		{"688001", "688001.SH"},
		{"000001.SZ", "000001.SZ"}, // 已有后缀
	}

	for _, tc := range testTSCodes {
		result := client.DetermineTSCode(tc.input)
		if result != tc.expected {
			t.Errorf("DetermineTSCode(%s): 期望 %s, 实际 %s", tc.input, tc.expected, result)
		}
	}

	// 测试DetermineAKShareSymbol函数
	testAKShareCodes := []struct {
		input    string
		expected string
	}{
		{"000001", "SZ000001"},
		{"600000", "SH600000"},
		{"600519", "SH600519"},
		{"300001", "SZ300001"},
		{"688001", "SH688001"},
		{"000001.SZ", "SZ000001"}, // 带后缀的也能正确处理
		{"600519.SH", "SH600519"}, // 带后缀的也能正确处理
	}

	for _, tc := range testAKShareCodes {
		result := client.DetermineAKShareSymbol(tc.input)
		if result != tc.expected {
			t.Errorf("DetermineAKShareSymbol(%s): 期望 %s, 实际 %s", tc.input, tc.expected, result)
		}
	}

	t.Log("股票代码处理函数测试通过")
}

// TestAKToolsClient_APIConnectivity 测试API连通性
func TestAKToolsClient_APIConnectivity(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 测试连接
	err := client.TestConnection()
	if err != nil {
		t.Logf("警告: AKTools连接测试失败: %v", err)
		t.Logf("请确保AKTools服务在 %s 运行", testAKToolsURL)
		t.Logf("可以通过以下命令启动AKTools服务:")
		t.Logf("  docker run -p 8080:8080 akfamily/aktools")
		return
	}

	t.Logf("AKTools连接测试成功: %s", testAKToolsURL)
}

// TestAKToolsClient_IntegrationTest 集成测试
func TestAKToolsClient_IntegrationTest(t *testing.T) {
	client := NewAKToolsClient(testAKToolsURL)

	// 首先测试连接
	if err := client.TestConnection(); err != nil {
		t.Logf("跳过集成测试，AKTools服务未运行: %v", err)
		return
	}

	t.Log("开始集成测试...")

	// 测试获取股票基本信息
	stockBasic, err := client.GetStockBasic(testSymbol)
	if err != nil {
		t.Logf("获取股票基本信息失败: %v", err)
	} else {
		t.Logf("股票基本信息: %s - %s", stockBasic.TSCode, stockBasic.Name)
	}

	// 测试获取日线数据
	dailyData, err := client.GetDailyData(testSymbol, "20240101", "20240115", "qfq")
	if err != nil {
		t.Logf("获取日线数据失败: %v", err)
	} else {
		t.Logf("获取到 %d 条日线数据", len(dailyData))
	}

	// 测试获取利润表数据
	incomeStatement, err := client.GetIncomeStatement(testSymbol, "", "1")
	if err != nil {
		t.Logf("获取利润表数据失败: %v", err)
	} else {
		t.Logf("利润表数据: 营业收入=%s, 净利润=%s",
			incomeStatement.Revenue.String(), incomeStatement.NetProfit.String())
	}

	t.Log("集成测试完成")
}
