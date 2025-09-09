package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("🧪 测试修改后的回测流程")

	// 创建回测配置
	config := map[string]interface{}{
		"name":         "测试回测 - " + time.Now().Format("2006-01-02 15:04:05"),
		"strategy_id":  "macd_strategy",
		"symbols":      []string{"000001", "600000"},
		"start_date":   "2024-01-01",
		"end_date":     "2024-01-10",
		"initial_cash": 100000,
		"commission":   0.001,
		"slippage":     0.001,
		"benchmark":    "000001",
	}

	// 转换为JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("❌ JSON序列化失败: %v\n", err)
		return
	}

	fmt.Printf("📝 回测配置: %s\n", string(jsonData))

	// 发送创建回测请求
	url := "http://localhost:8081/api/v1/backtests"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ 读取响应失败: %v\n", err)
		return
	}

	fmt.Printf("📊 HTTP状态码: %d\n", resp.StatusCode)
	fmt.Printf("📋 响应内容: %s\n", string(body))

	// 解析响应
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("❌ 解析响应失败: %v\n", err)
		return
	}

	if success, ok := response["success"].(bool); ok && success {
		if data, ok := response["data"].(map[string]interface{}); ok {
			if backtestId, ok := data["id"].(string); ok {
				fmt.Printf("✅ 回测创建并启动成功! ID: %s\n", backtestId)

				// 等待一段时间，然后检查进度
				fmt.Println("⏳ 等待3秒后检查回测进度...")
				time.Sleep(3 * time.Second)

				// 检查回测进度
				progressURL := fmt.Sprintf("http://localhost:8081/api/v1/backtests/%s/progress", backtestId)
				progressResp, err := http.Get(progressURL)
				if err != nil {
					fmt.Printf("❌ 获取进度失败: %v\n", err)
					return
				}
				defer progressResp.Body.Close()

				progressBody, err := io.ReadAll(progressResp.Body)
				if err != nil {
					fmt.Printf("❌ 读取进度响应失败: %v\n", err)
					return
				}

				fmt.Printf("📈 回测进度: %s\n", string(progressBody))
			} else {
				fmt.Println("❌ 响应中缺少回测ID")
			}
		} else {
			fmt.Println("❌ 响应中缺少data字段")
		}
	} else {
		fmt.Printf("❌ 回测创建失败: %s\n", response["message"])
	}
}
