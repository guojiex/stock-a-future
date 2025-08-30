package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// TestAKToolsStockBasicAPI 测试AKTools获取股票基本信息API
func TestAKToolsStockBasicAPI(t *testing.T) {
	// 创建AKTools客户端
	client := NewAKToolsClient("http://127.0.0.1:8080")

	// 测试用的股票代码
	testSymbol := "600976.SH"

	t.Logf("测试获取股票基本信息: %s", testSymbol)

	// 使用修复后的GetStockBasic方法
	stockBasic, err := client.GetStockBasic(testSymbol)
	if err != nil {
		t.Fatalf("获取股票基本信息失败: %v", err)
	}

	// 验证返回的数据
	if stockBasic == nil {
		t.Fatal("返回的股票基本信息为nil")
	}

	t.Logf("✅ 成功获取股票基本信息:")
	t.Logf("  股票代码: %s", stockBasic.Symbol)
	t.Logf("  股票名称: %s", stockBasic.Name)
	t.Logf("  TSCode: %s", stockBasic.TSCode)
	t.Logf("  地区: %s", stockBasic.Area)
	t.Logf("  行业: %s", stockBasic.Industry)
	t.Logf("  市场: %s", stockBasic.Market)
	t.Logf("  上市日期: %s", stockBasic.ListDate)

	// 验证必要字段不为空
	if stockBasic.Symbol == "" {
		t.Error("股票代码为空")
	}
	if stockBasic.Name == "" {
		t.Error("股票名称为空")
	}
	if stockBasic.TSCode == "" {
		t.Error("TSCode为空")
	}
}

// TestAKToolsStockCodeFormats 测试不同股票代码格式对AKTools API的影响
func TestAKToolsStockCodeFormats(t *testing.T) {
	// 测试用的股票代码
	testSymbol := "600976"

	// 创建AKTools客户端
	client := NewAKToolsClient("http://127.0.0.1:8080")

	// 测试不同的股票代码格式
	testCases := []struct {
		name        string
		symbol      string
		description string
	}{
		{"原始代码", testSymbol, "直接使用6位股票代码"},
		{"上海后缀", testSymbol + ".SH", "带上海市场后缀"},
		{"清理后代码", client.CleanStockSymbol(testSymbol + ".SH"), "使用CleanStockSymbol清理后的代码"},
		{"AKShare格式", client.DetermineAKShareSymbol(testSymbol), "使用AKShare财务API格式"},
		{"TSCode格式", client.DetermineTSCode(testSymbol), "使用TSCode格式"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("测试股票代码格式: %s (%s) - %s", tc.symbol, tc.name, tc.description)

			// 使用Go客户端方法测试股票基本信息
			stockBasic, err := client.GetStockBasic(tc.symbol)
			if err != nil {
				t.Logf("❌ 股票基本信息API失败: %v", err)
			} else {
				t.Logf("✅ 股票基本信息API成功: %s - %s", stockBasic.Symbol, stockBasic.Name)
			}

			// 使用Go客户端方法测试日线数据
			dailyData, err := client.GetDailyData(tc.symbol, "20241201", "20241231", "qfq")
			if err != nil {
				t.Logf("❌ 日线数据API失败: %v", err)
			} else {
				t.Logf("✅ 日线数据API成功: 获取到 %d 条数据", len(dailyData))
			}

			t.Logf("---")
		})
	}
}

// TestAKToolsConnectionAndAPI 测试AKTools连接和API可用性
func TestAKToolsConnectionAndAPI(t *testing.T) {
	client := NewAKToolsClient("http://127.0.0.1:8080")

	// 测试连接
	t.Run("连接测试", func(t *testing.T) {
		err := client.TestConnection()
		if err != nil {
			t.Fatalf("AKTools连接失败: %v", err)
		}
		t.Logf("✅ AKTools连接成功")
	})

	// 测试股票列表API
	t.Run("股票列表API", func(t *testing.T) {
		apiURL := "http://127.0.0.1:8080/api/public/stock_zh_a_spot"

		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Get(apiURL)
		if err != nil {
			t.Fatalf("请求股票列表API失败: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("股票列表API返回非200状态码: %d, 响应: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("读取响应体失败: %v", err)
		}

		var result []map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			t.Fatalf("解析JSON失败: %v, 响应: %s", err, string(body))
		}

		if len(result) == 0 {
			t.Fatalf("股票列表API返回空数据")
		}

		t.Logf("✅ 股票列表API成功，获取到 %d 只股票", len(result))

		// 打印前几只股票的信息
		for i, stock := range result[:min(5, len(result))] {
			t.Logf("  股票 %d: 代码=%v, 名称=%v", i+1, stock["代码"], stock["名称"])
		}
	})
}

// TestSpecificStockCode 测试特定股票代码
func TestSpecificStockCode(t *testing.T) {
	testCodes := []string{
		"600976", // 健民集团
		"000001", // 平安银行
		"002415", // 海康威视
		"688001", // 华兴源创
	}

	client := NewAKToolsClient("http://127.0.0.1:8080")

	for _, code := range testCodes {
		t.Run(fmt.Sprintf("股票_%s", code), func(t *testing.T) {
			t.Logf("测试股票代码: %s", code)

			// 使用修复后的Go客户端方法测试股票基本信息
			stockBasic, err := client.GetStockBasic(code)
			if err != nil {
				t.Logf("❌ 获取股票基本信息失败: %v", err)
			} else {
				t.Logf("✅ 获取股票基本信息成功: %s - %s", stockBasic.Symbol, stockBasic.Name)
			}

			// 测试日线数据获取
			dailyData, err := client.GetDailyData(code, "20241201", "20241231", "qfq")
			if err != nil {
				t.Logf("❌ 获取日线数据失败: %v", err)
			} else {
				t.Logf("✅ 获取日线数据成功: 获取到 %d 条数据", len(dailyData))
			}
		})
	}
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
