package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"stock-a-future/internal/client"
)

func main() {
	// 定义命令行参数
	var (
		output = flag.String("output", "data/stock_list.json", "输出文件路径")
		source = flag.String("source", "all", "数据源 (sse|all)")
		help   = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		fmt.Println("Stock List Fetcher - 从交易所官网获取股票列表")
		fmt.Println("")
		fmt.Println("用法:")
		fmt.Printf("  %s [选项]\n", os.Args[0])
		fmt.Println("")
		fmt.Println("选项:")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("示例:")
		fmt.Printf("  %s -source=sse -output=sse_stocks.json\n", os.Args[0])

		fmt.Printf("  %s -source=all -output=all_stocks.json\n", os.Args[0])
		return
	}

	log.Printf("开始获取股票列表...")
	log.Printf("数据源: %s", *source)
	log.Printf("输出文件: %s", *output)

	// 创建交易所客户端
	exchangeClient := client.NewExchangeClient()

	// 确保输出目录存在
	outputDir := filepath.Dir(*output)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	var err error
	switch *source {
	case "sse":
		err = fetchSSEStocks(exchangeClient, *output)

	case "all":
		err = fetchAllStocks(exchangeClient, *output)
	default:
		log.Fatalf("不支持的数据源: %s (支持: sse, all)", *source)
	}

	if err != nil {
		log.Fatalf("获取股票列表失败: %v", err)
	}

	log.Printf("股票列表获取完成，已保存到: %s", *output)
}

// fetchSSEStocks 获取上交所股票
func fetchSSEStocks(client *client.ExchangeClient, output string) error {
	stocks, err := client.GetSSEStockList()
	if err != nil {
		return fmt.Errorf("获取上交所股票失败: %w", err)
	}

	if err := client.SaveStockListToFile(stocks, output); err != nil {
		return fmt.Errorf("保存股票列表失败: %w", err)
	}

	log.Printf("成功获取上交所股票 %d 只", len(stocks))
	return nil
}

// fetchAllStocks 获取所有交易所股票
func fetchAllStocks(client *client.ExchangeClient, output string) error {
	stocks, err := client.GetAllStockList()
	if err != nil {
		return fmt.Errorf("获取所有股票失败: %w", err)
	}

	if err := client.SaveStockListToFile(stocks, output); err != nil {
		return fmt.Errorf("保存股票列表失败: %w", err)
	}

	log.Printf("成功获取所有股票 %d 只", len(stocks))
	return nil
}
