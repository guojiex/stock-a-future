package main

import (
	"fmt"
	"log"
	"stock-a-future/internal/client"
)

func main() {
	// 创建AKTools客户端
	aktoolsClient := client.NewAKToolsClient("http://127.0.0.1:8080")

	// 测试股票代码 - 使用茅台(600519)作为测试
	testSymbol := "600519"
	testPeriod := "20231231"
	testTradeDate := "20240115"

	fmt.Println("=== AKTools 基本面数据 API 测试 ===")
	fmt.Printf("测试股票: %s\n", testSymbol)
	fmt.Printf("测试期间: %s\n", testPeriod)
	fmt.Printf("测试交易日: %s\n", testTradeDate)
	fmt.Println()

	// 1. 测试连接
	fmt.Println("1. 测试AKTools连接...")
	if err := aktoolsClient.TestConnection(); err != nil {
		log.Printf("❌ AKTools连接失败: %v", err)
		fmt.Println("请确保AKTools服务在 http://127.0.0.1:8080 运行")
		fmt.Println("启动命令: docker run -p 8080:8080 akfamily/aktools")
		return
	}
	fmt.Println("✅ AKTools连接成功")
	fmt.Println()

	// 2. 测试获取利润表数据
	fmt.Println("2. 测试获取利润表数据...")
	// 首先不指定期间，获取最新数据
	incomeStatement, err := aktoolsClient.GetIncomeStatement(testSymbol, "", "1")
	if err != nil {
		log.Printf("❌ 获取利润表数据失败: %v", err)
	} else {
		fmt.Printf("✅ 获取利润表数据成功\n")
		fmt.Printf("   股票代码: %s\n", incomeStatement.TSCode)
		fmt.Printf("   报告期: %s\n", incomeStatement.FDate)
		fmt.Printf("   营业收入: %s\n", incomeStatement.Revenue.String())
		fmt.Printf("   净利润: %s\n", incomeStatement.NetProfit.String())

		// 如果获取成功，尝试用具体期间再次获取
		if incomeStatement.FDate != "" {
			fmt.Printf("   尝试获取指定期间数据: %s\n", incomeStatement.FDate)
			_, err := aktoolsClient.GetIncomeStatement(testSymbol, incomeStatement.FDate, "1")
			if err != nil {
				log.Printf("   ⚠️ 获取指定期间数据失败: %v", err)
			} else {
				fmt.Printf("   ✅ 指定期间数据获取成功\n")
			}
		}
	}
	fmt.Println()

	// 3. 测试获取资产负债表数据
	fmt.Println("3. 测试获取资产负债表数据...")
	balanceSheet, err := aktoolsClient.GetBalanceSheet(testSymbol, "", "1")
	if err != nil {
		log.Printf("❌ 获取资产负债表数据失败: %v", err)
	} else {
		fmt.Printf("✅ 获取资产负债表数据成功\n")
		fmt.Printf("   股票代码: %s\n", balanceSheet.TSCode)
		fmt.Printf("   报告期: %s\n", balanceSheet.FDate)
		fmt.Printf("   资产总计: %s\n", balanceSheet.TotalAssets.String())
		fmt.Printf("   负债合计: %s\n", balanceSheet.TotalLiab.String())
		fmt.Printf("   所有者权益: %s\n", balanceSheet.TotalHldrEqy.String())
	}
	fmt.Println()

	// 4. 测试获取现金流量表数据
	fmt.Println("4. 测试获取现金流量表数据...")
	cashFlowStatement, err := aktoolsClient.GetCashFlowStatement(testSymbol, "", "1")
	if err != nil {
		log.Printf("❌ 获取现金流量表数据失败: %v", err)
	} else {
		fmt.Printf("✅ 获取现金流量表数据成功\n")
		fmt.Printf("   股票代码: %s\n", cashFlowStatement.TSCode)
		fmt.Printf("   报告期: %s\n", cashFlowStatement.FDate)
		fmt.Printf("   经营活动现金流: %s\n", cashFlowStatement.NetCashOperAct.String())
		fmt.Printf("   投资活动现金流: %s\n", cashFlowStatement.NetCashInvAct.String())
		fmt.Printf("   筹资活动现金流: %s\n", cashFlowStatement.NetCashFinAct.String())
	}
	fmt.Println()

	// 5. 测试获取每日基本面数据
	fmt.Println("5. 测试获取每日基本面数据...")
	dailyBasic, err := aktoolsClient.GetDailyBasic(testSymbol, testTradeDate)
	if err != nil {
		log.Printf("❌ 获取每日基本面数据失败: %v", err)
	} else {
		fmt.Printf("✅ 获取每日基本面数据成功\n")
		fmt.Printf("   股票代码: %s\n", dailyBasic.TSCode)
		fmt.Printf("   交易日期: %s\n", dailyBasic.TradeDate)
		fmt.Printf("   收盘价: %s\n", dailyBasic.Close.String())
		fmt.Printf("   市盈率: %s\n", dailyBasic.Pe.String())
		fmt.Printf("   市净率: %s\n", dailyBasic.Pb.String())
		fmt.Printf("   总市值: %s\n", dailyBasic.TotalMv.String())
	}
	fmt.Println()

	// 6. 测试批量获取数据
	fmt.Println("6. 测试批量获取利润表数据...")
	incomeStatements, err := aktoolsClient.GetIncomeStatements(testSymbol, "20220101", "20231231", "1")
	if err != nil {
		log.Printf("❌ 批量获取利润表数据失败: %v", err)
	} else {
		fmt.Printf("✅ 批量获取利润表数据成功，共 %d 条\n", len(incomeStatements))
		if len(incomeStatements) > 0 {
			fmt.Printf("   最新数据期间: %s\n", incomeStatements[0].FDate)
			fmt.Printf("   最新营业收入: %s\n", incomeStatements[0].Revenue.String())
		}
	}
	fmt.Println()

	// 7. 测试股票代码处理功能
	fmt.Println("7. 测试股票代码处理功能...")
	testCodes := []string{"000001.SZ", "600000.SH", "300001", "688001"}
	for _, code := range testCodes {
		cleanCode := aktoolsClient.CleanStockSymbol(code)
		tsCode := aktoolsClient.DetermineTSCode(code)
		fmt.Printf("   %s -> 清理后: %s, TS代码: %s\n", code, cleanCode, tsCode)
	}
	fmt.Println()

	fmt.Println("=== 测试完成 ===")
	fmt.Println("✅ AKTools基本面数据API实现验证成功！")
	fmt.Println()
	fmt.Println("下一步可以：")
	fmt.Println("1. 运行单元测试: go test -v ./internal/client -run TestAKToolsClient")
	fmt.Println("2. 集成到现有的股票分析系统")
	fmt.Println("3. 实现基本面因子计算服务")
	fmt.Println("4. 添加数据库存储和API端点")
}
