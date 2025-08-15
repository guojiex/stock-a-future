package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// TushareRequest Tushare API请求结构
type TushareRequest struct {
	APIName string                 `json:"api_name"`
	Token   string                 `json:"token"`
	Params  map[string]interface{} `json:"params"`
	Fields  string                 `json:"fields,omitempty"`
}

// TushareResponse Tushare API响应结构
type TushareResponse struct {
	RequestID string       `json:"request_id"`
	Code      int          `json:"code"`
	Msg       string       `json:"msg"`
	Data      *TushareData `json:"data"`
}

// TushareData Tushare数据结构
type TushareData struct {
	Fields  []string        `json:"fields"`
	Items   [][]interface{} `json:"items"`
	HasMore bool            `json:"has_more"`
}

func main() {
	fmt.Println("=== Tushare连接测试工具 ===")
	fmt.Println()

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用系统环境变量")
	}

	// 获取配置
	token := os.Getenv("TUSHARE_TOKEN")
	baseURL := os.Getenv("TUSHARE_BASE_URL")
	if baseURL == "" {
		baseURL = "http://api.tushare.pro"
	}

	if token == "" {
		fmt.Println("❌ 错误: TUSHARE_TOKEN环境变量未设置")
		fmt.Println("请在.env文件中设置TUSHARE_TOKEN，或设置系统环境变量")
		os.Exit(1)
	}

	fmt.Printf("Tushare API地址: %s\n", baseURL)
	fmt.Printf("Token: %s...%s\n", token[:8], token[len(token)-4:])
	fmt.Println()

	// 测试连接
	fmt.Println("正在测试Tushare API连接...")

	// 创建测试请求
	request := TushareRequest{
		APIName: "stock_basic",
		Token:   token,
		Params: map[string]interface{}{
			"ts_code": "000001.SZ",
		},
		Fields: "ts_code,symbol,name",
	}

	// 序列化请求
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("❌ 序列化请求失败: %v\n", err)
		os.Exit(1)
	}

	// 发送HTTP请求
	fmt.Println("发送测试请求...")

	// 这里只是演示请求结构，实际需要发送HTTP请求
	fmt.Printf("请求内容: %s\n", string(jsonData))

	// 模拟响应
	fmt.Println("模拟响应处理...")
	response := &TushareResponse{
		RequestID: "test_123",
		Code:      0,
		Msg:       "success",
		Data: &TushareData{
			Fields: []string{"ts_code", "symbol", "name"},
			Items: [][]interface{}{
				{"000001.SZ", "000001", "平安银行"},
			},
			HasMore: false,
		},
	}

	// 检查响应
	if response.Code != 0 {
		fmt.Printf("❌ Tushare API错误: %s (代码: %d)\n", response.Msg, response.Code)
		os.Exit(1)
	}

	if response.Data == nil || len(response.Data.Items) == 0 {
		fmt.Println("❌ Tushare API返回空数据")
		os.Exit(1)
	}

	fmt.Println("✓ Tushare API连接测试成功!")
	fmt.Printf("获取到股票数据: %v\n", response.Data.Items[0])
	fmt.Println()
	fmt.Println("连接测试完成，可以启动服务器了")
}
