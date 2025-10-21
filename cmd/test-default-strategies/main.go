package main

import (
	"fmt"
	"stock-a-future/internal/models"
)

func main() {
	fmt.Println("🔍 检查 DefaultStrategies 的创建时间")
	fmt.Println("========================================")
	fmt.Println()

	for i, strategy := range models.DefaultStrategies {
		fmt.Printf("[%d] %s\n", i+1, strategy.ID)
		fmt.Printf("    创建时间: %s\n", strategy.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	// 检查唯一性
	timeMap := make(map[string]bool)
	for _, strategy := range models.DefaultStrategies {
		timeStr := strategy.CreatedAt.Format("2006-01-02 15:04:05")
		timeMap[timeStr] = true
	}

	fmt.Println("========================================")
	fmt.Printf("唯一创建时间数量: %d\n", len(timeMap))

	if len(timeMap) > 1 {
		fmt.Println("✅ 代码中的创建时间是不同的")
	} else {
		fmt.Println("⚠️  代码中的创建时间是相同的")
	}
}
