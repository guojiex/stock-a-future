package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"stock-a-future/internal/models"

	"github.com/shopspring/decimal"
)

// TushareClient Tushare API客户端
type TushareClient struct {
	token   string
	baseURL string
	client  *http.Client
}

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

// NewTushareClient 创建Tushare客户端
func NewTushareClient(token, baseURL string) *TushareClient {
	return &TushareClient{
		token:   token,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetDailyData 获取股票日线数据
func (c *TushareClient) GetDailyData(tsCode, startDate, endDate string) ([]models.StockDaily, error) {
	params := make(map[string]interface{})

	if tsCode != "" {
		params["ts_code"] = tsCode
	}
	if startDate != "" {
		params["start_date"] = startDate
	}
	if endDate != "" {
		params["end_date"] = endDate
	}

	request := TushareRequest{
		APIName: "daily",
		Token:   c.token,
		Params:  params,
		Fields:  "ts_code,trade_date,open,high,low,close,pre_close,change,pct_chg,vol,amount",
	}

	response, err := c.makeRequest(request)
	if err != nil {
		return nil, fmt.Errorf("请求Tushare API失败: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("tushare API错误: %s (代码: %d)", response.Msg, response.Code)
	}

	return c.parseDailyData(response.Data)
}

// GetDailyDataByDate 根据交易日期获取所有股票数据
func (c *TushareClient) GetDailyDataByDate(tradeDate string) ([]models.StockDaily, error) {
	params := map[string]interface{}{
		"trade_date": tradeDate,
	}

	request := TushareRequest{
		APIName: "daily",
		Token:   c.token,
		Params:  params,
		Fields:  "ts_code,trade_date,open,high,low,close,pre_close,change,pct_chg,vol,amount",
	}

	response, err := c.makeRequest(request)
	if err != nil {
		return nil, fmt.Errorf("请求Tushare API失败: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("tushare API错误: %s (代码: %d)", response.Msg, response.Code)
	}

	return c.parseDailyData(response.Data)
}

// makeRequest 发送HTTP请求到Tushare API
func (c *TushareClient) makeRequest(req TushareRequest) (*TushareResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "stock-a-future/1.0")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	var tushareResp TushareResponse
	if err := json.Unmarshal(body, &tushareResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &tushareResp, nil
}

// parseDailyData 解析日线数据
func (c *TushareClient) parseDailyData(data *TushareData) ([]models.StockDaily, error) {
	if data == nil || len(data.Items) == 0 {
		return []models.StockDaily{}, nil
	}

	// 构建字段索引映射
	fieldMap := make(map[string]int)
	for i, field := range data.Fields {
		fieldMap[field] = i
	}

	var dailyData []models.StockDaily
	for _, item := range data.Items {
		daily := models.StockDaily{}

		// 解析各个字段
		if idx, ok := fieldMap["ts_code"]; ok && idx < len(item) {
			if val, ok := item[idx].(string); ok {
				daily.TSCode = val
			}
		}

		if idx, ok := fieldMap["trade_date"]; ok && idx < len(item) {
			if val, ok := item[idx].(string); ok {
				daily.TradeDate = val
			}
		}

		// 解析数值字段
		daily.Open = models.NewJSONDecimal(c.parseDecimal(item, fieldMap, "open"))
		daily.High = models.NewJSONDecimal(c.parseDecimal(item, fieldMap, "high"))
		daily.Low = models.NewJSONDecimal(c.parseDecimal(item, fieldMap, "low"))
		daily.Close = models.NewJSONDecimal(c.parseDecimal(item, fieldMap, "close"))
		daily.PreClose = models.NewJSONDecimal(c.parseDecimal(item, fieldMap, "pre_close"))
		daily.Change = models.NewJSONDecimal(c.parseDecimal(item, fieldMap, "change"))
		daily.PctChg = models.NewJSONDecimal(c.parseDecimal(item, fieldMap, "pct_chg"))
		daily.Vol = models.NewJSONDecimal(c.parseDecimal(item, fieldMap, "vol"))
		daily.Amount = models.NewJSONDecimal(c.parseDecimal(item, fieldMap, "amount"))

		dailyData = append(dailyData, daily)
	}

	return dailyData, nil
}

// parseDecimal 解析decimal字段
func (c *TushareClient) parseDecimal(item []interface{}, fieldMap map[string]int, fieldName string) decimal.Decimal {
	if idx, ok := fieldMap[fieldName]; ok && idx < len(item) {
		if item[idx] == nil {
			return decimal.Zero
		}

		switch val := item[idx].(type) {
		case float64:
			return decimal.NewFromFloat(val)
		case string:
			if d, err := decimal.NewFromString(val); err == nil {
				return d
			}
		case int:
			return decimal.NewFromInt(int64(val))
		case int64:
			return decimal.NewFromInt(val)
		}
	}
	return decimal.Zero
}

// TestConnection 测试Tushare连接
func (c *TushareClient) TestConnection() error {
	// 使用简单的API调用测试连接
	_, err := c.GetDailyDataByDate("20240101")
	if err != nil {
		return fmt.Errorf("tushare连接测试失败: %w", err)
	}
	return nil
}
