package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stock-a-future/config"
	"stock-a-future/internal/client"
	"stock-a-future/internal/handler"
	"time"
)

func main() {
	log.Println("最终API测试...")

	// 加载配置
	cfg := config.Load()

	// 创建Tushare客户端
	tushareClient := client.NewTushareClient(cfg.TushareToken, cfg.TushareBaseURL)

	// 创建处理器
	stockHandler := handler.NewStockHandler(tushareClient)

	// 创建路由器
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/stocks", stockHandler.GetStockList)
	mux.HandleFunc("GET /api/v1/stocks/{code}/basic", stockHandler.GetStockBasic)

	// 启动服务器
	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	go func() {
		log.Println("服务器启动在端口 8082")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("服务器错误: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(3 * time.Second)

	// 测试API
	testFinalAPIs()

	// 关闭服务器
	server.Close()
	log.Println("测试完成！")
}

func testFinalAPIs() {
	fmt.Println("=== 最终API测试 ===")

	// 1. 测试获取股票列表统计
	fmt.Println("\n1. 股票列表统计")
	resp, err := http.Get("http://localhost:8082/api/v1/stocks")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var stockListResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stockListResp); err != nil {
		fmt.Printf("解析错误: %v\n", err)
		return
	}

	if data, ok := stockListResp["data"].(map[string]interface{}); ok {
		if total, ok := data["total"].(float64); ok {
			fmt.Printf("  总股票数: %.0f\n", total)
		}
		if stocks, ok := data["stocks"].([]interface{}); ok {
			sseCount, szseCount := 0, 0
			for _, stock := range stocks {
				if s, ok := stock.(map[string]interface{}); ok {
					if market, ok := s["market"].(string); ok {
						if market == "SH" {
							sseCount++
						} else if market == "SZ" {
							szseCount++
						}
					}
				}
			}
			fmt.Printf("  上交所: %d 只\n", sseCount)
			fmt.Printf("  深交所: %d 只\n", szseCount)
		}
	}

	// 2. 测试混合股票查询（上交所+深交所）
	testCodes := []struct {
		code     string
		expected string
		market   string
	}{
		{"600000", "浦发银行", "SH"},
		{"000001", "平安银行", "SZ"},
		{"300750", "宁德时代", "SZ"},
		{"688111", "金山办公", "SH"},
		{"002415", "海康威视", "SZ"},
	}

	fmt.Println("\n2. 混合股票基本信息查询")
	for _, test := range testCodes {
		fmt.Printf("\n测试股票: %s\n", test.code)
		resp, err := http.Get(fmt.Sprintf("http://localhost:8082/api/v1/stocks/%s/basic", test.code))
		if err != nil {
			fmt.Printf("  错误: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		var basicResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&basicResp); err != nil {
			fmt.Printf("  解析错误: %v\n", err)
			continue
		}

		if data, ok := basicResp["data"].(map[string]interface{}); ok {
			name := data["name"].(string)
			market := data["market"].(string)
			tsCode := data["ts_code"].(string)

			fmt.Printf("  代码: %s\n", tsCode)
			fmt.Printf("  名称: %s\n", name)
			fmt.Printf("  市场: %s\n", market)

			if market == test.market {
				fmt.Printf("  ✅ 市场匹配\n")
			} else {
				fmt.Printf("  ❌ 市场不匹配 (期望: %s, 实际: %s)\n", test.market, market)
			}
		}
	}
}
