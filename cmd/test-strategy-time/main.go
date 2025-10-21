package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Strategy struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	Success bool `json:"success"`
	Data    struct {
		Items []Strategy `json:"items"`
	} `json:"data"`
}

func main() {
	fmt.Println("🔍 测试策略API - 检查创建时间")
	fmt.Println("================================")
	fmt.Println()

	// 等待服务器启动
	time.Sleep(2 * time.Second)

	// 调用API
	resp, err := http.Get("http://localhost:8081/api/v1/strategies")
	if err != nil {
		fmt.Printf("❌ API调用失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ 读取响应失败: %v\n", err)
		return
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("❌ 解析JSON失败: %v\n", err)
		return
	}

	if !response.Success {
		fmt.Println("❌ API返回失败")
		return
	}

	fmt.Printf("✅ 获取到 %d 个策略\n\n", len(response.Data.Items))

	// 显示每个策略的创建时间
	timeMap := make(map[string]bool)
	for i, strategy := range response.Data.Items {
		fmt.Printf("[%d] %s\n", i+1, strategy.ID)
		fmt.Printf("    名称: %s\n", strategy.Name)
		fmt.Printf("    状态: %s\n", strategy.Status)
		fmt.Printf("    创建时间: %s\n", strategy.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()

		timeStr := strategy.CreatedAt.Format("2006-01-02 15:04:05")
		timeMap[timeStr] = true
	}

	// 检查创建时间是否都不同
	fmt.Println("================================")
	fmt.Printf("创建时间唯一值数量: %d\n", len(timeMap))

	if len(timeMap) == 1 {
		fmt.Println("⚠️  问题：所有策略的创建时间都相同！")
		fmt.Println("    这会导致排序不稳定")
	} else {
		fmt.Println("✅ 策略创建时间各不相同")
		fmt.Println("    排序应该是稳定的")
	}
}
