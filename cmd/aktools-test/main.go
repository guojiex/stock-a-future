package main

import (
	"fmt"
	"log"
	"os"
	"stock-a-future/config"
	"stock-a-future/internal/client"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("=== AKTools连接测试工具 ===")
	fmt.Println()

	// 加载环境变量
	if err := loadEnvFile("config.env"); err != nil {
		log.Println("未找到config.env文件，使用系统环境变量")
	}

	// 获取配置
	cfg := config.Load()

	if cfg.DataSourceType != "aktools" {
		fmt.Printf("❌ 错误: 当前配置的数据源类型是 %s，不是 aktools\n", cfg.DataSourceType)
		fmt.Println("请修改config.env文件中的DATA_SOURCE_TYPE为aktools")
		os.Exit(1)
	}

	fmt.Printf("AKTools API地址: %s\n", cfg.AKToolsBaseURL)
	fmt.Printf("服务器端口: %s\n", cfg.ServerPort)
	fmt.Println()

	// 测试连接
	fmt.Println("正在测试AKTools API连接...")

	// 创建AKTools客户端
	aktoolsClient := client.NewAKToolsClient(cfg.AKToolsBaseURL)

	// 测试连接
	if err := aktoolsClient.TestConnection(); err != nil {
		fmt.Printf("❌ AKTools API连接测试失败: %v\n", err)
		fmt.Println()
		fmt.Println("可能的原因:")
		fmt.Println("1. AKTools服务未启动")
		fmt.Println("2. AKTools服务运行在不同的端口")
		fmt.Println("3. 网络连接问题")
		fmt.Println()
		fmt.Println("解决方案:")
		fmt.Println("1. 启动AKTools服务: python -m aktools")
		fmt.Println("2. 检查AKTools服务端口")
		fmt.Println("3. 修改config.env中的AKTOOLS_BASE_URL")
		os.Exit(1)
	}

	fmt.Println("✓ AKTools API连接测试成功!")
	fmt.Println()
	fmt.Println("连接测试完成，可以启动Stock-A-Future服务器了")
	fmt.Printf("服务器将在 http://%s:%s 启动\n", cfg.ServerHost, cfg.ServerPort)
}

// loadEnvFile 加载环境变量文件
func loadEnvFile(filename string) error {
	return godotenv.Load(filename)
}
