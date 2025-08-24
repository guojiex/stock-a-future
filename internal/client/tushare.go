package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
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
func (c *TushareClient) GetDailyData(tsCode, startDate, endDate, adjust string) ([]models.StockDaily, error) {
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

	// 设置复权方式，如果没有指定则使用前复权
	if adjust == "" {
		adjust = "qfq" // 默认前复权，更符合用户习惯
	}
	params["adjust"] = adjust

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

	// 按交易日期升序排序，确保时间轴从左到右（从早到晚）显示
	// Tushare API默认返回降序数据，我们需要反转为升序
	c.sortDailyDataByDate(dailyData)

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

// GetStockBasic 获取股票基本信息
func (c *TushareClient) GetStockBasic(tsCode string) (*models.StockBasic, error) {
	params := make(map[string]interface{})
	if tsCode != "" {
		params["ts_code"] = tsCode
	}

	request := TushareRequest{
		APIName: "stock_basic",
		Token:   c.token,
		Params:  params,
		Fields:  "ts_code,symbol,name,area,industry,market,list_date",
	}

	response, err := c.makeRequest(request)
	if err != nil {
		return nil, fmt.Errorf("请求Tushare API失败: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("tushare API错误: %s (代码: %d)", response.Msg, response.Code)
	}

	return c.parseStockBasic(response.Data, tsCode)
}

// parseStockBasic 解析股票基本信息
func (c *TushareClient) parseStockBasic(data *TushareData, tsCode string) (*models.StockBasic, error) {
	if data == nil || len(data.Items) == 0 {
		return nil, fmt.Errorf("未找到股票基本信息: %s", tsCode)
	}

	// 构建字段索引映射
	fieldMap := make(map[string]int)
	for i, field := range data.Fields {
		fieldMap[field] = i
	}

	// 取第一条记录（通常只有一条）
	item := data.Items[0]
	stockBasic := &models.StockBasic{}

	// 解析各个字段
	if idx, ok := fieldMap["ts_code"]; ok && idx < len(item) {
		if val, ok := item[idx].(string); ok {
			stockBasic.TSCode = val
		}
	}

	if idx, ok := fieldMap["symbol"]; ok && idx < len(item) {
		if val, ok := item[idx].(string); ok {
			stockBasic.Symbol = val
		}
	}

	if idx, ok := fieldMap["name"]; ok && idx < len(item) {
		if val, ok := item[idx].(string); ok {
			stockBasic.Name = val
		}
	}

	if idx, ok := fieldMap["area"]; ok && idx < len(item) {
		if val, ok := item[idx].(string); ok {
			stockBasic.Area = val
		}
	}

	if idx, ok := fieldMap["industry"]; ok && idx < len(item) {
		if val, ok := item[idx].(string); ok {
			stockBasic.Industry = val
		}
	}

	if idx, ok := fieldMap["market"]; ok && idx < len(item) {
		if val, ok := item[idx].(string); ok {
			stockBasic.Market = val
		}
	}

	if idx, ok := fieldMap["list_date"]; ok && idx < len(item) {
		if val, ok := item[idx].(string); ok {
			stockBasic.ListDate = val
		}
	}

	return stockBasic, nil
}

// sortDailyDataByDate 按交易日期升序排序日线数据
func (c *TushareClient) sortDailyDataByDate(data []models.StockDaily) {
	sort.Slice(data, func(i, j int) bool {
		// 按交易日期升序排序（从早到晚）
		return data[i].TradeDate < data[j].TradeDate
	})
}

// GetBaseURL 获取Tushare API基础URL
func (c *TushareClient) GetBaseURL() string {
	return c.baseURL
}

// GetStockList 获取股票列表
func (c *TushareClient) GetStockList() ([]models.StockBasic, error) {
	params := map[string]interface{}{
		"status": "L", // 上市状态
	}

	request := TushareRequest{
		APIName: "stock_basic",
		Token:   c.token,
		Params:  params,
		Fields:  "ts_code,symbol,name,area,industry,market,list_date",
	}

	response, err := c.makeRequest(request)
	if err != nil {
		return nil, fmt.Errorf("请求Tushare API失败: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("tushare API错误: %s (代码: %d)", response.Msg, response.Code)
	}

	return c.parseStockList(response.Data)
}

// parseStockList 解析股票列表
func (c *TushareClient) parseStockList(data *TushareData) ([]models.StockBasic, error) {
	if data == nil || len(data.Items) == 0 {
		return []models.StockBasic{}, nil
	}

	// 构建字段索引映射
	fieldMap := make(map[string]int)
	for i, field := range data.Fields {
		fieldMap[field] = i
	}

	var stockList []models.StockBasic
	for _, item := range data.Items {
		stockBasic := &models.StockBasic{}

		// 解析各个字段
		if idx, ok := fieldMap["ts_code"]; ok && idx < len(item) {
			if val, ok := item[idx].(string); ok {
				stockBasic.TSCode = val
			}
		}

		if idx, ok := fieldMap["symbol"]; ok && idx < len(item) {
			if val, ok := item[idx].(string); ok {
				stockBasic.Symbol = val
			}
		}

		if idx, ok := fieldMap["name"]; ok && idx < len(item) {
			if val, ok := item[idx].(string); ok {
				stockBasic.Name = val
			}
		}

		if idx, ok := fieldMap["area"]; ok && idx < len(item) {
			if val, ok := item[idx].(string); ok {
				stockBasic.Area = val
			}
		}

		if idx, ok := fieldMap["industry"]; ok && idx < len(item) {
			if val, ok := item[idx].(string); ok {
				stockBasic.Industry = val
			}
		}

		if idx, ok := fieldMap["market"]; ok && idx < len(item) {
			if val, ok := item[idx].(string); ok {
				stockBasic.Market = val
			}
		}

		if idx, ok := fieldMap["list_date"]; ok && idx < len(item) {
			if val, ok := item[idx].(string); ok {
				stockBasic.ListDate = val
			}
		}

		stockList = append(stockList, *stockBasic)
	}

	return stockList, nil
}

// TestConnection 测试Tushare连接
func (c *TushareClient) TestConnection() error {
	// 使用轻量级的API调用测试连接 - 获取股票基本信息
	// 选择000001.SZ（平安银行）作为测试股票，这是A股中比较稳定的股票
	testParams := map[string]interface{}{
		"ts_code": "000001.SZ",
	}

	request := TushareRequest{
		APIName: "stock_basic",
		Token:   c.token,
		Params:  testParams,
		Fields:  "ts_code,symbol,name",
	}

	// 设置较短的超时时间用于连接测试
	testClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 创建测试请求
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("序列化测试请求失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建测试HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "stock-a-future/1.0")

	// 发送测试请求
	resp, err := testClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("tushare连接测试失败: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭测试响应体失败: %v", err)
		}
	}()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tushare API返回非200状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取测试响应体失败: %w", err)
	}

	// 解析响应
	var tushareResp TushareResponse
	if err := json.Unmarshal(body, &tushareResp); err != nil {
		return fmt.Errorf("解析测试响应失败: %w", err)
	}

	// 检查API响应码
	if tushareResp.Code != 0 {
		return fmt.Errorf("tushare API错误: %s (代码: %d)", tushareResp.Msg, tushareResp.Code)
	}

	// 验证返回的数据
	if tushareResp.Data == nil || len(tushareResp.Data.Items) == 0 {
		return fmt.Errorf("tushare API返回空数据")
	}

	log.Printf("Tushare连接测试成功 - 获取到股票数据: %v", tushareResp.Data.Items[0])
	return nil
}

// ===== 基本面数据接口实现 =====

// GetIncomeStatement 获取利润表数据
// TODO: 实现Tushare利润表数据获取
func (c *TushareClient) GetIncomeStatement(symbol, period, reportType string) (*models.IncomeStatement, error) {
	panic("TODO: 实现Tushare利润表数据获取 - GetIncomeStatement")
}

// GetIncomeStatements 批量获取利润表数据
// TODO: 实现Tushare批量利润表数据获取
func (c *TushareClient) GetIncomeStatements(symbol, startPeriod, endPeriod, reportType string) ([]models.IncomeStatement, error) {
	panic("TODO: 实现Tushare批量利润表数据获取 - GetIncomeStatements")
}

// GetBalanceSheet 获取资产负债表数据
// TODO: 实现Tushare资产负债表数据获取
func (c *TushareClient) GetBalanceSheet(symbol, period, reportType string) (*models.BalanceSheet, error) {
	panic("TODO: 实现Tushare资产负债表数据获取 - GetBalanceSheet")
}

// GetBalanceSheets 批量获取资产负债表数据
// TODO: 实现Tushare批量资产负债表数据获取
func (c *TushareClient) GetBalanceSheets(symbol, startPeriod, endPeriod, reportType string) ([]models.BalanceSheet, error) {
	panic("TODO: 实现Tushare批量资产负债表数据获取 - GetBalanceSheets")
}

// GetCashFlowStatement 获取现金流量表数据
// TODO: 实现Tushare现金流量表数据获取
func (c *TushareClient) GetCashFlowStatement(symbol, period, reportType string) (*models.CashFlowStatement, error) {
	panic("TODO: 实现Tushare现金流量表数据获取 - GetCashFlowStatement")
}

// GetCashFlowStatements 批量获取现金流量表数据
// TODO: 实现Tushare批量现金流量表数据获取
func (c *TushareClient) GetCashFlowStatements(symbol, startPeriod, endPeriod, reportType string) ([]models.CashFlowStatement, error) {
	panic("TODO: 实现Tushare批量现金流量表数据获取 - GetCashFlowStatements")
}

// GetFinancialIndicator 获取财务指标数据
// TODO: 实现Tushare财务指标数据获取
func (c *TushareClient) GetFinancialIndicator(symbol, period, reportType string) (*models.FinancialIndicator, error) {
	panic("TODO: 实现Tushare财务指标数据获取 - GetFinancialIndicator")
}

// GetFinancialIndicators 批量获取财务指标数据
// TODO: 实现Tushare批量财务指标数据获取
func (c *TushareClient) GetFinancialIndicators(symbol, startPeriod, endPeriod, reportType string) ([]models.FinancialIndicator, error) {
	panic("TODO: 实现Tushare批量财务指标数据获取 - GetFinancialIndicators")
}

// GetDailyBasic 获取每日基本面指标
// TODO: 实现Tushare每日基本面指标获取
func (c *TushareClient) GetDailyBasic(symbol, tradeDate string) (*models.DailyBasic, error) {
	panic("TODO: 实现Tushare每日基本面指标获取 - GetDailyBasic")
}

// GetDailyBasics 批量获取每日基本面指标
// TODO: 实现Tushare批量每日基本面指标获取
func (c *TushareClient) GetDailyBasics(symbol, startDate, endDate string) ([]models.DailyBasic, error) {
	panic("TODO: 实现Tushare批量每日基本面指标获取 - GetDailyBasics")
}

// GetDailyBasicsByDate 根据交易日期获取所有股票的每日基本面指标
// TODO: 实现Tushare按日期获取所有股票每日基本面指标
func (c *TushareClient) GetDailyBasicsByDate(tradeDate string) ([]models.DailyBasic, error) {
	panic("TODO: 实现Tushare按日期获取所有股票每日基本面指标 - GetDailyBasicsByDate")
}

// ===== 基本面因子接口实现 =====

// GetFundamentalFactor 获取基本面因子数据
// TODO: 实现Tushare基本面因子数据获取
func (c *TushareClient) GetFundamentalFactor(symbol, tradeDate string) (*models.FundamentalFactor, error) {
	panic("TODO: 实现Tushare基本面因子数据获取 - GetFundamentalFactor")
}

// GetFundamentalFactors 批量获取基本面因子数据
// TODO: 实现Tushare批量基本面因子数据获取
func (c *TushareClient) GetFundamentalFactors(symbol, startDate, endDate string) ([]models.FundamentalFactor, error) {
	panic("TODO: 实现Tushare批量基本面因子数据获取 - GetFundamentalFactors")
}

// GetFundamentalFactorsByDate 根据交易日期获取所有股票的基本面因子
// TODO: 实现Tushare按日期获取所有股票基本面因子
func (c *TushareClient) GetFundamentalFactorsByDate(tradeDate string) ([]models.FundamentalFactor, error) {
	panic("TODO: 实现Tushare按日期获取所有股票基本面因子 - GetFundamentalFactorsByDate")
}
