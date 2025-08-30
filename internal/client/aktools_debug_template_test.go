package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func init() {
	// 设置测试环境标记，避免配置加载时因缺少token而panic
	os.Setenv("GO_TEST", "1")
	// 设置数据源类型为aktools
	os.Setenv("DATA_SOURCE_TYPE", "aktools")
}

// isLocalTestEnvironment 检查是否在本地测试环境中运行
// 在CI环境中通常会设置CI=true环境变量
func isLocalTestEnvironment() bool {
	// 检查常见的CI环境变量
	if os.Getenv("CI") == "true" ||
		os.Getenv("GITHUB_ACTIONS") == "true" ||
		os.Getenv("GITLAB_CI") == "true" ||
		os.Getenv("JENKINS_URL") != "" ||
		os.Getenv("TRAVIS") == "true" {
		return false
	}

	// 检查是否明确设置了跳过本地测试的环境变量
	if os.Getenv("SKIP_LOCAL_TESTS") == "true" {
		return false
	}

	return true
}

// TestAKToolsDebugTemplate AKTools调试测试模板
// 这个文件作为调试模板，当遇到AKTools API问题时可以参考
func TestAKToolsDebugTemplate(t *testing.T) {
	// 在非本地环境中跳过测试
	if !isLocalTestEnvironment() {
		t.Skip("在非本地环境中跳过AKTools调试测试")
	}

	// 创建AKTools客户端，默认端口8080
	client := NewAKToolsClient("http://127.0.0.1:8080")

	t.Log("=== AKTools调试测试开始 ===")

	// 1. 首先测试连接
	t.Run("连接测试", func(t *testing.T) {
		err := client.TestConnection()
		if err != nil {
			t.Logf("❌ AKTools连接失败: %v", err)
			t.Skip("AKTools服务可能未运行，跳过后续测试")
		}
		t.Logf("✅ AKTools连接成功")
	})

	// 2. 测试API文档获取
	t.Run("API文档", func(t *testing.T) {
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Get("http://127.0.0.1:8080/openapi.json")
		if err != nil {
			t.Logf("❌ 获取API文档失败: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			t.Logf("✅ API文档可访问")
		} else {
			t.Logf("❌ API文档返回状态码: %d", resp.StatusCode)
		}
	})
}

// TestAKToolsAPIDiscovery 发现可用的API接口
func TestAKToolsAPIDiscovery(t *testing.T) {
	// 在非本地环境中跳过测试
	if !isLocalTestEnvironment() {
		t.Skip("在非本地环境中跳过AKTools API发现测试")
	}

	// 常见的AKTools API接口名称
	apiInterfaces := []string{
		"stock_zh_a_hist",                    // 股票日线数据
		"stock_zh_a_spot",                    // 股票列表
		"stock_individual_info_em",           // 股票基本信息
		"stock_profit_sheet_by_report_em",    // 利润表
		"stock_balance_sheet_by_report_em",   // 资产负债表
		"stock_cash_flow_sheet_by_report_em", // 现金流量表
		"stock_zh_a_spot_em",                 // A股实时数据
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}

	// 先检查AKTools服务是否可用
	_, err := httpClient.Get("http://127.0.0.1:8080")
	if err != nil {
		t.Skip("AKTools服务可能未运行，跳过API发现测试")
	}

	for _, apiName := range apiInterfaces {
		t.Run(apiName, func(t *testing.T) {
			// 构建API URL
			apiURL := fmt.Sprintf("http://127.0.0.1:8080/api/public/%s", apiName)

			// 发送GET请求
			resp, err := httpClient.Get(apiURL)
			if err != nil {
				t.Logf("❌ %s - 请求失败: %v", apiName, err)
				return
			}
			defer resp.Body.Close()

			// 检查状态码
			switch resp.StatusCode {
			case 200:
				t.Logf("✅ %s - 接口可用", apiName)
			case 404:
				t.Logf("❌ %s - 接口不存在", apiName)
			default:
				body, _ := io.ReadAll(resp.Body)
				t.Logf("⚠️  %s - 状态码: %d, 响应: %s", apiName, resp.StatusCode, string(body))
			}
		})
	}
}

// TestStockCodeFormatsDebug 测试不同股票代码格式的调试
func TestStockCodeFormatsDebug(t *testing.T) {
	// 在非本地环境中跳过测试
	if !isLocalTestEnvironment() {
		t.Skip("在非本地环境中跳过股票代码格式测试")
	}

	client := NewAKToolsClient("http://127.0.0.1:8080")

	// 先检查AKTools服务是否可用
	httpClient := &http.Client{Timeout: 5 * time.Second}
	_, err := httpClient.Get("http://127.0.0.1:8080")
	if err != nil {
		t.Skip("AKTools服务可能未运行，跳过股票代码格式测试")
	}

	// 测试用股票代码
	testStock := "600976" // 健民集团

	// 不同的股票代码格式
	formats := []struct {
		name   string
		symbol string
	}{
		{"纯数字代码", testStock},
		{"上海后缀", testStock + ".SH"},
		{"清理后代码", client.CleanStockSymbol(testStock + ".SH")},
		{"TSCode格式", client.DetermineTSCode(testStock)},
		{"AKShare格式", client.DetermineAKShareSymbol(testStock)},
	}

	for _, format := range formats {
		t.Run(format.name, func(t *testing.T) {
			t.Logf("测试股票代码格式: %s (%s)", format.symbol, format.name)

			// 测试股票基本信息API
			stockBasic, err := client.GetStockBasic(format.symbol)
			if err != nil {
				t.Logf("❌ 股票基本信息失败: %v", err)
			} else {
				t.Logf("✅ 股票基本信息成功: %s - %s", stockBasic.Symbol, stockBasic.Name)
			}

			// 测试日线数据API
			dailyData, err := client.GetDailyData(format.symbol, "20241201", "20241231", "qfq")
			if err != nil {
				t.Logf("❌ 日线数据失败: %v", err)
			} else {
				t.Logf("✅ 日线数据成功: 获取到 %d 条数据", len(dailyData))
			}
		})
	}
}

// TestRawAPIResponse 测试原始API响应格式
func TestRawAPIResponse(t *testing.T) {
	// 在非本地环境中跳过测试
	if !isLocalTestEnvironment() {
		t.Skip("在非本地环境中跳过原始API响应测试")
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}

	// 先检查AKTools服务是否可用
	_, err := httpClient.Get("http://127.0.0.1:8080")
	if err != nil {
		t.Skip("AKTools服务可能未运行，跳过原始API响应测试")
	}

	// 测试stock_individual_info_em的原始响应
	t.Run("股票基本信息原始响应", func(t *testing.T) {
		apiURL := "http://127.0.0.1:8080/api/public/stock_individual_info_em?symbol=600976"

		resp, err := httpClient.Get(apiURL)
		if err != nil {
			t.Skipf("请求失败: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Skipf("API返回状态码: %d, 响应: %s", resp.StatusCode, string(body))
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Skipf("读取响应失败: %v", err)
			return
		}

		// 解析JSON
		var rawData []map[string]interface{}
		if err := json.Unmarshal(body, &rawData); err != nil {
			t.Skipf("解析JSON失败: %v", err)
			return
		}

		t.Logf("原始响应数据条数: %d", len(rawData))

		// 打印前几条数据
		for i, item := range rawData {
			if i >= 5 { // 只打印前5条
				break
			}
			t.Logf("  数据 %d: %+v", i+1, item)
		}

		// 转换为map格式
		stockData := make(map[string]interface{})
		for _, item := range rawData {
			if itemKey, ok := item["item"].(string); ok {
				if itemValue, exists := item["value"]; exists {
					stockData[itemKey] = itemValue
				}
			}
		}

		t.Logf("转换后的数据:")
		for key, value := range stockData {
			t.Logf("  %s: %v", key, value)
		}
	})
}

// TestMultipleStocks 测试多个股票代码
func TestMultipleStocks(t *testing.T) {
	// 在非本地环境中跳过测试
	if !isLocalTestEnvironment() {
		t.Skip("在非本地环境中跳过多股票代码测试")
	}

	client := NewAKToolsClient("http://127.0.0.1:8080")

	// 先检查AKTools服务是否可用
	httpClient := &http.Client{Timeout: 5 * time.Second}
	_, err := httpClient.Get("http://127.0.0.1:8080")
	if err != nil {
		t.Skip("AKTools服务可能未运行，跳过多股票代码测试")
	}

	testStocks := []struct {
		code string
		name string
	}{
		{"600976", "健民集团"},
		{"000001", "平安银行"},
		{"002415", "海康威视"},
		{"688001", "华兴源创"},
	}

	for _, stock := range testStocks {
		t.Run(stock.name, func(t *testing.T) {
			t.Logf("测试股票: %s (%s)", stock.name, stock.code)

			// 获取股票基本信息
			stockBasic, err := client.GetStockBasic(stock.code)
			if err != nil {
				t.Logf("❌ 获取股票基本信息失败: %v", err)
				return
			}

			t.Logf("✅ 股票基本信息:")
			t.Logf("  代码: %s", stockBasic.Symbol)
			t.Logf("  名称: %s", stockBasic.Name)
			t.Logf("  TSCode: %s", stockBasic.TSCode)
			t.Logf("  市场: %s", stockBasic.Market)

			// 验证名称是否匹配
			if stockBasic.Name != stock.name {
				t.Logf("⚠️  名称不匹配: 期望 %s, 实际 %s", stock.name, stockBasic.Name)
			}
		})
	}
}

// TestAKToolsPerformance 测试AKTools API性能
func TestAKToolsPerformance(t *testing.T) {
	// 在非本地环境中跳过测试
	if !isLocalTestEnvironment() {
		t.Skip("在非本地环境中跳过AKTools性能测试")
	}

	client := NewAKToolsClient("http://127.0.0.1:8080")

	// 先检查AKTools服务是否可用
	httpClient := &http.Client{Timeout: 5 * time.Second}
	_, err := httpClient.Get("http://127.0.0.1:8080")
	if err != nil {
		t.Skip("AKTools服务可能未运行，跳过性能测试")
	}

	testStock := "600976"
	iterations := 5

	t.Run("性能测试", func(t *testing.T) {
		var totalDuration time.Duration
		var successCount int

		for i := 0; i < iterations; i++ {
			start := time.Now()

			_, err := client.GetStockBasic(testStock)
			if err != nil {
				t.Logf("第 %d 次调用失败: %v", i+1, err)
				continue
			}

			duration := time.Since(start)
			totalDuration += duration
			successCount++
			t.Logf("第 %d 次调用耗时: %v", i+1, duration)
		}

		// 如果所有调用都失败，跳过后续测试
		if successCount == 0 {
			t.Skip("所有API调用都失败，跳过性能评估")
		}

		avgDuration := totalDuration / time.Duration(successCount)
		t.Logf("平均响应时间: %v (成功 %d/%d 次)", avgDuration, successCount, iterations)

		if avgDuration > 5*time.Second {
			t.Logf("⚠️  响应时间较慢，可能需要优化")
		} else {
			t.Logf("✅ 响应时间正常")
		}
	})
}
