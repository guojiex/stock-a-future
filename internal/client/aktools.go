package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"stock-a-future/config"
	"stock-a-future/internal/models"

	"github.com/shopspring/decimal"
)

// AKToolsClient AKTools HTTP API客户端
type AKToolsClient struct {
	baseURL string
	client  *http.Client
	config  *config.Config
}

// AKToolsDailyResponse AKTools日线数据响应结构
type AKToolsDailyResponse struct {
	Date      string  `json:"日期"`
	Open      float64 `json:"开盘"`
	Close     float64 `json:"收盘"`
	High      float64 `json:"最高"`
	Low       float64 `json:"最低"`
	Volume    float64 `json:"成交量"`
	Amount    float64 `json:"成交额"`
	Amplitude float64 `json:"振幅"`
	ChangePct float64 `json:"涨跌幅"`
	Change    float64 `json:"涨跌额"`
	Turnover  float64 `json:"换手率"`
}

// AKToolsStockBasicResponse AKTools股票基本信息响应结构
type AKToolsStockBasicResponse struct {
	Code     string `json:"代码"`
	Name     string `json:"名称"`
	Area     string `json:"地区"`
	Industry string `json:"行业"`
	Market   string `json:"市场"`
	ListDate string `json:"上市日期"`
}

// NewAKToolsClient 创建AKTools客户端
func NewAKToolsClient(baseURL string) *AKToolsClient {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8080"
	}

	// 创建一个默认配置（用于测试）
	cfg := &config.Config{
		DataSourceType: "aktools", // 默认使用aktools数据源
		AKToolsBaseURL: baseURL,
		Debug:          false,
	}

	// 尝试加载配置，但不要在失败时panic
	// 这样可以在测试环境中使用，即使没有完整的环境变量配置
	if os.Getenv("TUSHARE_TOKEN") != "" || os.Getenv("DATA_SOURCE_TYPE") == "aktools" {
		// 只有在必要的环境变量存在时才尝试加载配置
		loadedCfg := config.Load()
		if loadedCfg != nil {
			cfg = loadedCfg
		}
	}

	return &AKToolsClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: cfg,
	}
}

// CleanStockSymbol 清理股票代码，移除市场后缀
func (c *AKToolsClient) CleanStockSymbol(symbol string) string {
	// 移除常见的市场后缀
	suffixes := []string{".SH", ".SZ", ".BJ", ".sh", ".sz", ".bj"}
	for _, suffix := range suffixes {
		if len(symbol) > len(suffix) && strings.HasSuffix(symbol, suffix) {
			return symbol[:len(symbol)-len(suffix)]
		}
	}
	return symbol
}

// DetermineTSCode 智能判断股票代码的市场后缀
func (c *AKToolsClient) DetermineTSCode(symbol string) string {
	// 如果symbol已经包含市场后缀，直接返回
	if strings.Contains(symbol, ".") {
		return symbol
	}

	// 根据股票代码规则判断市场
	// 600xxx, 601xxx, 603xxx, 688xxx -> 上海
	// 000xxx, 002xxx, 300xxx -> 深圳
	// 430xxx, 830xxx, 870xxx -> 北京
	if strings.HasPrefix(symbol, "600") || strings.HasPrefix(symbol, "601") ||
		strings.HasPrefix(symbol, "603") || strings.HasPrefix(symbol, "688") {
		return symbol + ".SH"
	} else if strings.HasPrefix(symbol, "000") || strings.HasPrefix(symbol, "002") ||
		strings.HasPrefix(symbol, "300") {
		return symbol + ".SZ"
	} else if strings.HasPrefix(symbol, "430") || strings.HasPrefix(symbol, "830") ||
		strings.HasPrefix(symbol, "870") {
		return symbol + ".BJ"
	}

	// 默认返回上海市场
	return symbol + ".SH"
}

// DetermineAKShareSymbol 转换为AKShare财务报表API需要的股票代码格式
func (c *AKToolsClient) DetermineAKShareSymbol(symbol string) string {
	// 清理股票代码
	cleanSymbol := c.CleanStockSymbol(symbol)

	// 根据股票代码规则判断市场并返回AKShare格式
	// 600xxx, 601xxx, 603xxx, 688xxx -> SH600xxx
	// 000xxx, 002xxx, 300xxx -> SZ000xxx
	// 430xxx, 830xxx, 870xxx -> BJ430xxx
	if strings.HasPrefix(cleanSymbol, "600") || strings.HasPrefix(cleanSymbol, "601") ||
		strings.HasPrefix(cleanSymbol, "603") || strings.HasPrefix(cleanSymbol, "688") {
		return "SH" + cleanSymbol
	} else if strings.HasPrefix(cleanSymbol, "000") || strings.HasPrefix(cleanSymbol, "002") ||
		strings.HasPrefix(cleanSymbol, "300") {
		return "SZ" + cleanSymbol
	} else if strings.HasPrefix(cleanSymbol, "430") || strings.HasPrefix(cleanSymbol, "830") ||
		strings.HasPrefix(cleanSymbol, "870") {
		return "BJ" + cleanSymbol
	}

	// 默认返回上海市场格式
	return "SH" + cleanSymbol
}

// saveResponseToFile 保存HTTP响应到JSON文件用于调试
func (c *AKToolsClient) saveResponseToFile(responseBody []byte, apiName, symbol string, debug bool) error {
	// 如果未启用调试模式，直接返回
	if !debug {
		return nil
	}

	// 创建debug目录（如果不存在）
	debugDir := "debug"
	if _, err := os.Stat(debugDir); os.IsNotExist(err) {
		err := os.Mkdir(debugDir, 0755)
		if err != nil {
			return fmt.Errorf("创建debug目录失败: %v", err)
		}
	}

	// 生成文件名: api名称_股票代码_时间戳.json
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s/%s_%s_%s.json", debugDir, apiName, symbol, timestamp)

	// 写入文件
	err := os.WriteFile(filename, responseBody, 0644)
	if err != nil {
		return fmt.Errorf("写入响应文件失败: %v", err)
	}

	log.Printf("HTTP响应已保存到文件: %s", filename)
	return nil
}

// GetDailyData 获取股票日线数据
func (c *AKToolsClient) GetDailyData(symbol, startDate, endDate, adjust string) ([]models.StockDaily, error) {
	// 清理股票代码，移除市场后缀
	cleanSymbol := c.CleanStockSymbol(symbol)

	// 构建查询参数
	params := url.Values{}
	params.Set("symbol", cleanSymbol)
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)
	params.Set("adjust", adjust)

	// 构建完整URL
	apiURL := fmt.Sprintf("%s/api/public/stock_zh_a_hist?%s", c.baseURL, params.Encode())

	// 创建带context的请求
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送HTTP请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools API失败: %w, URL: %s, 股票代码: %s", err, apiURL, symbol)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AKTools API返回非200状态码: %d, URL: %s, 股票代码: %s, 开始日期: %s, 结束日期: %s",
			resp.StatusCode, apiURL, symbol, startDate, endDate)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	// 保存响应到文件用于调试
	if err := c.saveResponseToFile(body, "daily_data", cleanSymbol, c.config.Debug); err != nil {
		log.Printf("保存响应文件失败: %v", err)
	}

	// 解析JSON响应
	var aktoolsResp []AKToolsDailyResponse
	if err := json.Unmarshal(body, &aktoolsResp); err != nil {
		return nil, fmt.Errorf("解析AKTools响应失败: %w", err)
	}

	// 转换为内部模型
	return c.convertToStockDaily(aktoolsResp, symbol), nil
}

// GetDailyDataByDate 根据交易日期获取所有股票数据
func (c *AKToolsClient) GetDailyDataByDate(tradeDate string) ([]models.StockDaily, error) {
	// AKTools暂不支持按日期批量获取，这里返回空结果
	// 可以通过其他方式实现，比如先获取股票列表再逐个查询
	return []models.StockDaily{}, nil
}

// GetStockBasic 获取股票基本信息
func (c *AKToolsClient) GetStockBasic(symbol string) (*models.StockBasic, error) {
	// 清理股票代码，移除市场后缀
	cleanSymbol := c.CleanStockSymbol(symbol)

	// 构建查询参数
	params := url.Values{}
	params.Set("symbol", cleanSymbol)

	// 构建完整URL - 使用股票基本信息API
	apiURL := fmt.Sprintf("%s/api/public/stock_individual_info_em?%s", c.baseURL, params.Encode())

	// 创建带context的请求
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送HTTP请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools股票信息API失败: %w, URL: %s, 股票代码: %s", err, apiURL, symbol)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AKTools API返回非200状态码: %d, URL: %s, 股票代码: %s", resp.StatusCode, apiURL, symbol)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	// 保存响应到文件用于调试
	if err := c.saveResponseToFile(body, "stock_basic", cleanSymbol, c.config.Debug); err != nil {
		log.Printf("保存响应文件失败: %v", err)
	}

	// 解析JSON响应 - stock_individual_info_em返回的是key-value对数组格式
	var rawResp []map[string]interface{}
	if err := json.Unmarshal(body, &rawResp); err != nil {
		return nil, fmt.Errorf("解析AKTools股票信息响应失败: %w", err)
	}

	if len(rawResp) == 0 {
		return nil, fmt.Errorf("未找到股票基本信息: %s", symbol)
	}

	// 将key-value对数组转换为map
	stockData := make(map[string]interface{})
	for _, item := range rawResp {
		if itemKey, ok := item["item"].(string); ok {
			if itemValue, exists := item["value"]; exists {
				stockData[itemKey] = itemValue
			}
		}
	}

	// 转换为内部模型
	return c.convertStockIndividualInfoToStockBasic(stockData, symbol), nil
}

// GetStockList 获取股票列表
func (c *AKToolsClient) GetStockList() ([]models.StockBasic, error) {
	// 构建完整URL - 使用股票列表API
	apiURL := fmt.Sprintf("%s/api/public/stock_zh_a_spot", c.baseURL)

	// 创建带context的请求
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送HTTP请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools股票列表API失败: %w, URL: %s", err, apiURL)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AKTools API返回非200状态码: %d, URL: %s", resp.StatusCode, apiURL)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	// 保存响应到文件用于调试
	if err := c.saveResponseToFile(body, "stock_list", "all", c.config.Debug); err != nil {
		log.Printf("保存响应文件失败: %v", err)
	}

	// 解析JSON响应 - 这里需要根据实际API响应结构调整
	var stockList []map[string]interface{}
	if err := json.Unmarshal(body, &stockList); err != nil {
		return nil, fmt.Errorf("解析AKTools股票列表响应失败: %w", err)
	}

	// 转换为内部模型
	var result []models.StockBasic
	for _, stock := range stockList {
		basic := &models.StockBasic{}

		if code, ok := stock["代码"].(string); ok {
			basic.Symbol = code
			basic.TSCode = code + ".SH" // 默认上海，实际应该根据代码判断
		}
		if name, ok := stock["名称"].(string); ok {
			basic.Name = name
		}
		if area, ok := stock["地区"].(string); ok {
			basic.Area = area
		}
		if industry, ok := stock["行业"].(string); ok {
			basic.Industry = industry
		}
		if market, ok := stock["市场"].(string); ok {
			basic.Market = market
		}

		result = append(result, *basic)
	}

	return result, nil
}

// convertToStockDaily 将AKTools日线数据转换为内部模型
func (c *AKToolsClient) convertToStockDaily(aktoolsData []AKToolsDailyResponse, symbol string) []models.StockDaily {
	var result []models.StockDaily

	for _, data := range aktoolsData {
		// 智能判断市场后缀
		tsCode := c.DetermineTSCode(symbol)

		// 转换日期格式：将AKTools的日期格式转换为前端期望的YYYYMMDD格式
		formattedDate := c.formatDateForFrontend(data.Date)

		daily := models.StockDaily{
			TSCode:    tsCode,
			TradeDate: formattedDate,
			Open:      models.NewJSONDecimal(decimal.NewFromFloat(data.Open)),
			High:      models.NewJSONDecimal(decimal.NewFromFloat(data.High)),
			Low:       models.NewJSONDecimal(decimal.NewFromFloat(data.Low)),
			Close:     models.NewJSONDecimal(decimal.NewFromFloat(data.Close)),
			PreClose:  models.NewJSONDecimal(decimal.Zero), // AKTools不提供前收盘价
			Change:    models.NewJSONDecimal(decimal.NewFromFloat(data.Change)),
			PctChg:    models.NewJSONDecimal(decimal.NewFromFloat(data.ChangePct)),
			Vol:       models.NewJSONDecimal(decimal.NewFromFloat(data.Volume)),
			Amount:    models.NewJSONDecimal(decimal.NewFromFloat(data.Amount)),
		}
		result = append(result, daily)
	}

	// 按交易日期升序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].TradeDate < result[j].TradeDate
	})

	return result
}

// convertToStockBasic 将AKTools股票基本信息转换为内部模型
func (c *AKToolsClient) convertToStockBasic(aktoolsData AKToolsStockBasicResponse, symbol string) *models.StockBasic {
	// 智能判断市场后缀
	tsCode := c.DetermineTSCode(symbol)

	return &models.StockBasic{
		TSCode:   tsCode,
		Symbol:   aktoolsData.Code,
		Name:     aktoolsData.Name,
		Area:     aktoolsData.Area,
		Industry: aktoolsData.Industry,
		Market:   aktoolsData.Market,
		ListDate: aktoolsData.ListDate,
	}
}

// convertStockIndividualInfoToStockBasic 将stock_individual_info_em的响应转换为内部模型
func (c *AKToolsClient) convertStockIndividualInfoToStockBasic(stockData map[string]interface{}, symbol string) *models.StockBasic {
	// 智能判断市场后缀
	tsCode := c.DetermineTSCode(symbol)

	stockBasic := &models.StockBasic{
		TSCode: tsCode,
	}

	// 从stockData中提取字段
	if code, ok := stockData["股票代码"].(string); ok {
		stockBasic.Symbol = code
	}
	if name, ok := stockData["股票简称"].(string); ok {
		stockBasic.Name = name
	}
	// stock_individual_info_em 不提供地区信息，设为空
	stockBasic.Area = ""

	// stock_individual_info_em 不提供行业信息，设为空
	stockBasic.Industry = ""

	// stock_individual_info_em 不提供市场信息，根据股票代码判断
	if strings.HasPrefix(stockBasic.Symbol, "60") || strings.HasPrefix(stockBasic.Symbol, "68") {
		stockBasic.Market = "上海主板"
	} else if strings.HasPrefix(stockBasic.Symbol, "00") {
		stockBasic.Market = "深圳主板"
	} else if strings.HasPrefix(stockBasic.Symbol, "30") {
		stockBasic.Market = "创业板"
	} else {
		stockBasic.Market = "未知"
	}

	// stock_individual_info_em 不提供上市日期，设为空
	stockBasic.ListDate = ""

	return stockBasic
}

// GetBaseURL 获取AKTools API基础URL
func (c *AKToolsClient) GetBaseURL() string {
	return c.baseURL
}

// formatDateForFrontend 将AKTools的日期格式转换为前端期望的YYYYMMDD格式
func (c *AKToolsClient) formatDateForFrontend(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// 移除可能的空格和换行符
	dateStr = strings.TrimSpace(dateStr)

	// 处理不同的日期格式
	if strings.Contains(dateStr, "-") {
		// 格式可能是 "2025-08-" 或 "2025-08-15"
		parts := strings.Split(dateStr, "-")
		if len(parts) >= 2 {
			year := parts[0]
			month := parts[1]

			// 检查年份和月份是否为纯数字
			if !isNumeric(year) || !isNumeric(month) {
				return dateStr // 如果不是纯数字，返回原始字符串
			}

			// 确保月份是两位数
			if len(month) == 1 {
				month = "0" + month
			}

			// 如果有日期部分，使用它；否则使用"01"
			day := "01"
			if len(parts) >= 3 && parts[2] != "" {
				day = parts[2]
				// 检查日期是否为纯数字
				if !isNumeric(day) {
					return dateStr // 如果不是纯数字，返回原始字符串
				}
				// 确保日期是两位数
				if len(day) == 1 {
					day = "0" + day
				}
			}

			// 返回YYYYMMDD格式
			return year + month + day
		}
	}

	// 如果已经是8位数字格式，直接返回
	if len(dateStr) == 8 && isNumeric(dateStr) {
		return dateStr
	}

	// 如果无法解析，返回原始字符串
	return dateStr
}

// isNumeric 检查字符串是否为纯数字
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// TestConnection 测试AKTools连接
func (c *AKToolsClient) TestConnection() error {
	// 设置较短的超时时间用于连接测试
	testClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 直接使用股票日线数据API作为测试端点
	apiURL := fmt.Sprintf("%s/api/public/stock_zh_a_hist", c.baseURL)

	// 打印测试请求的URL
	log.Printf("正在测试AKTools API连接，请求URL: %s", apiURL)

	// 创建带context的请求
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送测试请求
	resp, err := testClient.Do(req)
	if err != nil {
		return fmt.Errorf("AKTools连接测试失败: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭测试响应体失败: %v", err)
		}
	}()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		// 读取响应体用于调试
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("读取错误响应体失败: %v", err)
		} else {
			log.Printf("错误响应体内容: %s", string(body))
		}
		return fmt.Errorf("AKTools API返回非200状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取测试响应体失败: %w", err)
	}

	// 尝试解析响应 - 股票日线数据API返回的是数组格式
	var dailyData []AKToolsDailyResponse
	if err := json.Unmarshal(body, &dailyData); err != nil {
		return fmt.Errorf("解析测试响应失败: %w", err)
	}

	// 验证返回的数据
	if len(dailyData) == 0 {
		return fmt.Errorf("AKTools API返回空数据")
	}

	log.Printf("AKTools连接测试成功 - 获取到股票日线数据，共%d条", len(dailyData))
	return nil
}

// ===== 基本面数据接口实现 =====

// GetIncomeStatement 获取利润表数据
func (c *AKToolsClient) GetIncomeStatement(symbol, period, reportType string) (*models.IncomeStatement, error) {
	// 转换为AKShare财务报表API需要的股票代码格式
	akshareSymbol := c.DetermineAKShareSymbol(symbol)

	// 构建查询参数 - 不传递period参数，因为AKShare API不支持
	params := url.Values{}
	params.Set("symbol", akshareSymbol)

	// 构建完整URL - 使用利润表API
	apiURL := fmt.Sprintf("%s/api/public/stock_profit_sheet_by_report_em?%s", c.baseURL, params.Encode())

	// 创建带context的请求
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送HTTP请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools利润表API失败: %w, URL: %s", err, apiURL)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AKTools API返回非200状态码: %d, URL: %s", resp.StatusCode, apiURL)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	// 保存响应到文件用于调试 - 已完成调试
	// cleanSymbol := c.CleanStockSymbol(symbol)
	// if err := c.saveResponseToFile(body, "income_statement", cleanSymbol); err != nil {
	// 	log.Printf("保存响应文件失败: %v", err)
	// }

	// 保存响应到文件用于调试
	cleanSymbol := c.CleanStockSymbol(symbol)
	if err := c.saveResponseToFile(body, "income_statement", cleanSymbol, c.config.Debug); err != nil {
		log.Printf("保存响应文件失败: %v", err)
	}

	// 解析JSON响应
	var rawData []map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("解析AKTools利润表响应失败: %w", err)
	}

	if len(rawData) == 0 {
		return nil, fmt.Errorf("未找到利润表数据: %s, 期间: %s", symbol, period)
	}

	// 如果指定了period，尝试找到匹配的记录
	if period != "" {
		for _, data := range rawData {
			if reportDate, ok := data["报告期"].(string); ok {
				formattedDate := c.formatDateForFrontend(reportDate)
				if formattedDate == period {
					return c.convertToIncomeStatement(data, symbol, period, reportType)
				}
			}
		}
		// 如果没有找到匹配的period，返回错误
		return nil, fmt.Errorf("未找到指定期间的利润表数据: %s, 期间: %s", symbol, period)
	}

	// 如果没有指定period，返回最新的一条数据
	return c.convertToIncomeStatement(rawData[0], symbol, period, reportType)
}

// GetIncomeStatements 批量获取利润表数据
func (c *AKToolsClient) GetIncomeStatements(symbol, startPeriod, endPeriod, reportType string) ([]models.IncomeStatement, error) {
	// 转换为AKShare财务报表API需要的股票代码格式
	akshareSymbol := c.DetermineAKShareSymbol(symbol)

	// 构建查询参数
	params := url.Values{}
	params.Set("symbol", akshareSymbol)

	// 构建完整URL
	apiURL := fmt.Sprintf("%s/api/public/stock_profit_sheet_by_report_em?%s", c.baseURL, params.Encode())

	// 创建带context的请求
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送HTTP请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools利润表API失败: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AKTools API返回非200状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	// 保存响应到文件用于调试
	cleanSymbol := c.CleanStockSymbol(symbol)
	if err := c.saveResponseToFile(body, "income_statements", cleanSymbol, c.config.Debug); err != nil {
		log.Printf("保存响应文件失败: %v", err)
	}

	// 解析JSON响应
	var rawData []map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("解析AKTools利润表响应失败: %w", err)
	}

	// 保存响应到文件用于调试 - 已完成调试
	// cleanSymbol := c.CleanStockSymbol(symbol)
	// if err := c.saveResponseToFile(body, "income_statements", cleanSymbol); err != nil {
	// 	log.Printf("保存响应文件失败: %v", err)
	// }

	// 转换为内部模型
	var results []models.IncomeStatement
	for _, data := range rawData {
		// 从数据中提取期间信息 - 使用实际API字段名
		period := ""
		if reportDate, ok := data["REPORT_DATE"].(string); ok {
			period = c.formatDateForFrontend(reportDate)
		}

		// 过滤期间范围（如果指定了）
		if startPeriod != "" && period < startPeriod {
			continue
		}
		if endPeriod != "" && period > endPeriod {
			continue
		}

		incomeStatement, err := c.convertToIncomeStatement(data, symbol, period, reportType)
		if err != nil {
			continue // 跳过转换失败的数据
		}
		results = append(results, *incomeStatement)
	}

	return results, nil
}

// GetBalanceSheet 获取资产负债表数据
func (c *AKToolsClient) GetBalanceSheet(symbol, period, reportType string) (*models.BalanceSheet, error) {
	// 转换为AKShare财务报表API需要的股票代码格式
	akshareSymbol := c.DetermineAKShareSymbol(symbol)

	// 构建查询参数 - 不传递period参数，因为AKShare API不支持
	params := url.Values{}
	params.Set("symbol", akshareSymbol)

	// 构建完整URL - 使用资产负债表API
	apiURL := fmt.Sprintf("%s/api/public/stock_balance_sheet_by_report_em?%s", c.baseURL, params.Encode())

	// 创建带context的请求
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送HTTP请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools资产负债表API失败: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AKTools API返回非200状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 保存响应到文件用于调试
	cleanSymbol := c.CleanStockSymbol(symbol)
	if err := c.saveResponseToFile(body, "balance_sheet", cleanSymbol, c.config.Debug); err != nil {
		log.Printf("保存响应文件失败: %v", err)
	}

	// 解析JSON响应
	var rawData []map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("解析AKTools资产负债表响应失败: %w", err)
	}
	if len(rawData) == 0 {
		return nil, fmt.Errorf("未找到资产负债表数据: %s, 期间: %s", symbol, period)
	}

	// 如果指定了period，尝试找到匹配的记录
	if period != "" {
		for _, data := range rawData {
			if reportDate, ok := data["报告期"].(string); ok {
				formattedDate := c.formatDateForFrontend(reportDate)
				if formattedDate == period {
					return c.convertToBalanceSheet(data, symbol, period, reportType)
				}
			}
		}
		// 如果没有找到匹配的period，返回错误
		return nil, fmt.Errorf("未找到指定期间的资产负债表数据: %s, 期间: %s", symbol, period)
	}

	// 如果没有指定period，返回最新的一条数据
	return c.convertToBalanceSheet(rawData[0], symbol, period, reportType)
}

// GetBalanceSheets 批量获取资产负债表数据
func (c *AKToolsClient) GetBalanceSheets(symbol, startPeriod, endPeriod, reportType string) ([]models.BalanceSheet, error) {
	// 转换为AKShare财务报表API需要的股票代码格式
	akshareSymbol := c.DetermineAKShareSymbol(symbol)
	params := url.Values{}
	params.Set("symbol", akshareSymbol)
	apiURL := fmt.Sprintf("%s/api/public/stock_balance_sheet_by_report_em?%s", c.baseURL, params.Encode())

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools资产负债表API失败: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AKTools API返回非200状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	// 保存响应到文件用于调试
	cleanSymbol := c.CleanStockSymbol(symbol)
	if err := c.saveResponseToFile(body, "balance_sheets", cleanSymbol, c.config.Debug); err != nil {
		log.Printf("保存响应文件失败: %v", err)
	}

	var rawData []map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("解析AKTools资产负债表响应失败: %w", err)
	}

	var results []models.BalanceSheet
	for _, data := range rawData {
		// 使用实际API字段名
		period := ""
		if reportDate, ok := data["REPORT_DATE"].(string); ok {
			period = c.formatDateForFrontend(reportDate)
		}

		if startPeriod != "" && period < startPeriod {
			continue
		}
		if endPeriod != "" && period > endPeriod {
			continue
		}

		balanceSheet, err := c.convertToBalanceSheet(data, symbol, period, reportType)
		if err != nil {
			continue
		}
		results = append(results, *balanceSheet)
	}

	return results, nil
}

// GetCashFlowStatement 获取现金流量表数据
func (c *AKToolsClient) GetCashFlowStatement(symbol, period, reportType string) (*models.CashFlowStatement, error) {
	// 转换为AKShare财务报表API需要的股票代码格式
	akshareSymbol := c.DetermineAKShareSymbol(symbol)
	params := url.Values{}
	params.Set("symbol", akshareSymbol)
	// 不传递period参数，因为AKShare API不支持

	apiURL := fmt.Sprintf("%s/api/public/stock_cash_flow_sheet_by_report_em?%s", c.baseURL, params.Encode())

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools现金流量表API失败: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AKTools API返回非200状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	// 保存响应到文件用于调试
	cleanSymbol := c.CleanStockSymbol(symbol)
	if err := c.saveResponseToFile(body, "cash_flow_statement", cleanSymbol, c.config.Debug); err != nil {
		log.Printf("保存响应文件失败: %v", err)
	}

	var rawData []map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("解析AKTools现金流量表响应失败: %w", err)
	}
	if len(rawData) == 0 {
		return nil, fmt.Errorf("未找到现金流量表数据: %s, 期间: %s", symbol, period)
	}

	// 如果指定了period，尝试找到匹配的记录
	if period != "" {
		for _, data := range rawData {
			if reportDate, ok := data["报告期"].(string); ok {
				formattedDate := c.formatDateForFrontend(reportDate)
				if formattedDate == period {
					return c.convertToCashFlowStatement(data, symbol, period, reportType)
				}
			}
		}
		// 如果没有找到匹配的period，返回错误
		return nil, fmt.Errorf("未找到指定期间的现金流量表数据: %s, 期间: %s", symbol, period)
	}

	// 如果没有指定period，返回最新的一条数据
	return c.convertToCashFlowStatement(rawData[0], symbol, period, reportType)
}

// GetCashFlowStatements 批量获取现金流量表数据
func (c *AKToolsClient) GetCashFlowStatements(symbol, startPeriod, endPeriod, reportType string) ([]models.CashFlowStatement, error) {
	// 转换为AKShare财务报表API需要的股票代码格式
	akshareSymbol := c.DetermineAKShareSymbol(symbol)
	params := url.Values{}
	params.Set("symbol", akshareSymbol)
	apiURL := fmt.Sprintf("%s/api/public/stock_cash_flow_sheet_by_report_em?%s", c.baseURL, params.Encode())

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools现金流量表API失败: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AKTools API返回非200状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	// 保存响应到文件用于调试
	cleanSymbol := c.CleanStockSymbol(symbol)
	if err := c.saveResponseToFile(body, "cash_flow_statements", cleanSymbol, c.config.Debug); err != nil {
		log.Printf("保存响应文件失败: %v", err)
	}

	var rawData []map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("解析AKTools现金流量表响应失败: %w", err)
	}

	var results []models.CashFlowStatement
	for _, data := range rawData {
		// 使用实际API字段名
		period := ""
		if reportDate, ok := data["REPORT_DATE"].(string); ok {
			period = c.formatDateForFrontend(reportDate)
		}

		if startPeriod != "" && period < startPeriod {
			continue
		}
		if endPeriod != "" && period > endPeriod {
			continue
		}

		cashFlowStatement, err := c.convertToCashFlowStatement(data, symbol, period, reportType)
		if err != nil {
			continue
		}
		results = append(results, *cashFlowStatement)
	}

	return results, nil
}

// GetFinancialIndicator 获取财务指标数据
func (c *AKToolsClient) GetFinancialIndicator(symbol, period, reportType string) (*models.FinancialIndicator, error) {
	// AKTools暂不支持直接的财务指标API，返回空实现
	// 可以通过其他API组合计算得出
	return &models.FinancialIndicator{
		FinancialStatement: models.FinancialStatement{
			TSCode:     c.DetermineTSCode(symbol),
			FDate:      period,
			EndDate:    period,
			ReportType: reportType,
		},
	}, nil
}

// GetFinancialIndicators 批量获取财务指标数据
func (c *AKToolsClient) GetFinancialIndicators(symbol, startPeriod, endPeriod, reportType string) ([]models.FinancialIndicator, error) {
	// AKTools暂不支持直接的财务指标API，返回空切片
	return []models.FinancialIndicator{}, nil
}

// GetDailyBasic 获取每日基本面指标
func (c *AKToolsClient) GetDailyBasic(symbol, tradeDate string) (*models.DailyBasic, error) {
	cleanSymbol := c.CleanStockSymbol(symbol)

	// stock_zh_a_spot_em 不接受任何参数，返回所有A股实时数据
	// 我们需要在返回的数据中找到对应的股票
	apiURL := fmt.Sprintf("%s/api/public/stock_zh_a_spot_em", c.baseURL)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools每日基本面API失败: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AKTools API返回非200状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 保存响应到文件用于调试
	if err := c.saveResponseToFile(body, "daily_basic", cleanSymbol, c.config.Debug); err != nil {
		log.Printf("保存响应文件失败: %v", err)
	}

	var rawData []map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("解析AKTools每日基本面响应失败: %w", err)
	}

	if len(rawData) == 0 {
		return nil, fmt.Errorf("未找到每日基本面数据: %s, 日期: %s", symbol, tradeDate)
	}

	// 在返回的所有股票数据中查找指定的股票
	for _, data := range rawData {
		if code, ok := data["代码"].(string); ok {
			// 比较清理后的代码
			if code == cleanSymbol {
				return c.convertToDailyBasic(data, symbol, tradeDate)
			}
		}
		// 也尝试匹配带后缀的代码
		if code, ok := data["代码"].(string); ok {
			expectedTSCode := c.DetermineTSCode(cleanSymbol)
			if code == expectedTSCode {
				return c.convertToDailyBasic(data, symbol, tradeDate)
			}
		}
	}

	return nil, fmt.Errorf("未找到指定股票的每日基本面数据: %s, 日期: %s", symbol, tradeDate)
}

// GetDailyBasics 批量获取每日基本面指标
func (c *AKToolsClient) GetDailyBasics(symbol, startDate, endDate string) ([]models.DailyBasic, error) {
	// AKTools不支持历史每日基本面数据批量获取，返回空切片
	// 实际应用中需要逐日调用GetDailyBasic
	return []models.DailyBasic{}, nil
}

// GetDailyBasicsByDate 根据交易日期获取所有股票的每日基本面指标
func (c *AKToolsClient) GetDailyBasicsByDate(tradeDate string) ([]models.DailyBasic, error) {
	// AKTools不支持按日期获取所有股票数据，返回空切片
	return []models.DailyBasic{}, nil
}

// ===== 基本面因子接口实现 =====

// GetFundamentalFactor 获取基本面因子数据
func (c *AKToolsClient) GetFundamentalFactor(symbol, tradeDate string) (*models.FundamentalFactor, error) {
	// AKTools不直接提供基本面因子，需要通过计算服务生成
	// 这里返回基础结构，实际计算由FundamentalFactorCalculator完成
	return &models.FundamentalFactor{
		TSCode:    c.DetermineTSCode(symbol),
		TradeDate: tradeDate,
		UpdatedAt: time.Now(),
	}, nil
}

// GetFundamentalFactors 批量获取基本面因子数据
func (c *AKToolsClient) GetFundamentalFactors(symbol, startDate, endDate string) ([]models.FundamentalFactor, error) {
	// AKTools不直接提供基本面因子，返回空切片
	return []models.FundamentalFactor{}, nil
}

// GetFundamentalFactorsByDate 根据交易日期获取所有股票的基本面因子
func (c *AKToolsClient) GetFundamentalFactorsByDate(tradeDate string) ([]models.FundamentalFactor, error) {
	// AKTools不直接提供基本面因子，返回空切片
	return []models.FundamentalFactor{}, nil
}

// ===== 数据转换辅助函数 =====

// convertToIncomeStatement 将AKTools利润表数据转换为内部模型
func (c *AKToolsClient) convertToIncomeStatement(data map[string]interface{}, symbol, period, reportType string) (*models.IncomeStatement, error) {
	incomeStatement := &models.IncomeStatement{}

	// 设置基础字段 - 根据实际API响应
	if secucode, ok := data["SECUCODE"].(string); ok {
		incomeStatement.TSCode = secucode
	} else {
		incomeStatement.TSCode = c.DetermineTSCode(symbol)
	}

	// 报告期 - 从REPORT_DATE字段提取
	if reportDate, ok := data["REPORT_DATE"].(string); ok {
		// 格式化日期: "2025-06-30 00:00:00" -> "20250630"
		incomeStatement.FDate = c.formatDateForFrontend(reportDate)
		incomeStatement.EndDate = c.formatDateForFrontend(reportDate)
	} else {
		incomeStatement.FDate = period
		incomeStatement.EndDate = period
	}

	// 从REPORT_TYPE字段获取报告类型
	if rptType, ok := data["REPORT_TYPE"].(string); ok {
		incomeStatement.ReportType = rptType
	} else {
		incomeStatement.ReportType = reportType
	}

	// 提取公告日期 - 从NOTICE_DATE字段
	if noticeDate, ok := data["NOTICE_DATE"].(string); ok {
		incomeStatement.AnnDate = c.formatDateForFrontend(noticeDate)
	}

	// 营业收入相关 - 使用实际API字段名
	if totalOperIncome, ok := data["TOTAL_OPERATE_INCOME"]; ok {
		incomeStatement.Revenue = models.NewJSONDecimal(c.parseDecimalFromInterface(totalOperIncome))
		incomeStatement.OperRevenue = models.NewJSONDecimal(c.parseDecimalFromInterface(totalOperIncome))
	}
	if operIncome, ok := data["OPERATE_INCOME"]; ok {
		incomeStatement.OperRevenue = models.NewJSONDecimal(c.parseDecimalFromInterface(operIncome))
	}

	// 成本费用 - 使用实际API字段名
	if totalOperCost, ok := data["TOTAL_OPERATE_COST"]; ok {
		incomeStatement.OperCost = models.NewJSONDecimal(c.parseDecimalFromInterface(totalOperCost))
	}
	if operCost, ok := data["OPERATE_COST"]; ok {
		incomeStatement.OperCost = models.NewJSONDecimal(c.parseDecimalFromInterface(operCost))
	}
	if manageExp, ok := data["MANAGE_EXPENSE"]; ok {
		incomeStatement.AdminExp = models.NewJSONDecimal(c.parseDecimalFromInterface(manageExp))
	}
	if financeExp, ok := data["FINANCE_EXPENSE"]; ok {
		incomeStatement.FinExp = models.NewJSONDecimal(c.parseDecimalFromInterface(financeExp))
	}
	if researchExp, ok := data["RESEARCH_EXPENSE"]; ok {
		incomeStatement.RdExp = models.NewJSONDecimal(c.parseDecimalFromInterface(researchExp))
	}
	if saleExp, ok := data["SALE_EXPENSE"]; ok {
		incomeStatement.OperExp = models.NewJSONDecimal(c.parseDecimalFromInterface(saleExp))
	}

	// 利润相关 - 使用实际API字段名
	if operProfit, ok := data["OPERATE_PROFIT"]; ok {
		incomeStatement.OperProfit = models.NewJSONDecimal(c.parseDecimalFromInterface(operProfit))
	}
	if totalProfit, ok := data["TOTAL_PROFIT"]; ok {
		incomeStatement.TotalProfit = models.NewJSONDecimal(c.parseDecimalFromInterface(totalProfit))
	}
	if netProfit, ok := data["NETPROFIT"]; ok {
		incomeStatement.NetProfit = models.NewJSONDecimal(c.parseDecimalFromInterface(netProfit))
	}
	if deductParentNetProfit, ok := data["DEDUCT_PARENT_NETPROFIT"]; ok {
		incomeStatement.NetProfitDedt = models.NewJSONDecimal(c.parseDecimalFromInterface(deductParentNetProfit))
	}

	// 每股收益 - 使用实际API字段名
	if basicEps, ok := data["BASIC_EPS"]; ok {
		incomeStatement.BasicEps = models.NewJSONDecimal(c.parseDecimalFromInterface(basicEps))
	}
	if dilutedEps, ok := data["DILUTED_EPS"]; ok {
		incomeStatement.DilutedEps = models.NewJSONDecimal(c.parseDecimalFromInterface(dilutedEps))
	}

	return incomeStatement, nil
}

// convertToBalanceSheet 将AKTools资产负债表数据转换为内部模型
func (c *AKToolsClient) convertToBalanceSheet(data map[string]interface{}, symbol, period, reportType string) (*models.BalanceSheet, error) {
	balanceSheet := &models.BalanceSheet{}

	// 设置基础字段 - 根据实际API响应
	if secucode, ok := data["SECUCODE"].(string); ok {
		balanceSheet.TSCode = secucode
	} else {
		balanceSheet.TSCode = c.DetermineTSCode(symbol)
	}

	// 报告期 - 从REPORT_DATE字段提取
	if reportDate, ok := data["REPORT_DATE"].(string); ok {
		balanceSheet.FDate = c.formatDateForFrontend(reportDate)
		balanceSheet.EndDate = c.formatDateForFrontend(reportDate)
	} else {
		balanceSheet.FDate = period
		balanceSheet.EndDate = period
	}

	// 从REPORT_TYPE字段获取报告类型
	if rptType, ok := data["REPORT_TYPE"].(string); ok {
		balanceSheet.ReportType = rptType
	} else {
		balanceSheet.ReportType = reportType
	}

	// 提取公告日期 - 从NOTICE_DATE字段
	if noticeDate, ok := data["NOTICE_DATE"].(string); ok {
		balanceSheet.AnnDate = c.formatDateForFrontend(noticeDate)
	}

	// 资产 - 使用实际API字段名
	if totalAssets, ok := data["TOTAL_ASSETS"]; ok {
		balanceSheet.TotalAssets = models.NewJSONDecimal(c.parseDecimalFromInterface(totalAssets))
	}
	if totalCurAssets, ok := data["TOTAL_CURRENT_ASSETS"]; ok {
		balanceSheet.TotalCurAssets = models.NewJSONDecimal(c.parseDecimalFromInterface(totalCurAssets))
	}
	if money, ok := data["MONEY_CAP"]; ok {
		balanceSheet.Money = models.NewJSONDecimal(c.parseDecimalFromInterface(money))
	}
	if tradAssets, ok := data["TRADE_FINASSET"]; ok {
		balanceSheet.TradAssets = models.NewJSONDecimal(c.parseDecimalFromInterface(tradAssets))
	}
	if notesReceiv, ok := data["NOTE_ACCOUNTS_RECE"]; ok {
		balanceSheet.NotesReceiv = models.NewJSONDecimal(c.parseDecimalFromInterface(notesReceiv))
	}
	if accountsReceiv, ok := data["ACCOUNTS_RECE"]; ok {
		balanceSheet.AccountsReceiv = models.NewJSONDecimal(c.parseDecimalFromInterface(accountsReceiv))
	}
	if inventoryAssets, ok := data["INVENTORY"]; ok {
		balanceSheet.InventoryAssets = models.NewJSONDecimal(c.parseDecimalFromInterface(inventoryAssets))
	}
	if totalNcaAssets, ok := data["TOTAL_NONCURRENT_ASSETS"]; ok {
		balanceSheet.TotalNcaAssets = models.NewJSONDecimal(c.parseDecimalFromInterface(totalNcaAssets))
	}
	if fixAssets, ok := data["FIXED_ASSET"]; ok {
		balanceSheet.FixAssets = models.NewJSONDecimal(c.parseDecimalFromInterface(fixAssets))
	}
	if cipAssets, ok := data["CIP"]; ok {
		balanceSheet.CipAssets = models.NewJSONDecimal(c.parseDecimalFromInterface(cipAssets))
	}
	if intangAssets, ok := data["INTANGIBLE_ASSET"]; ok {
		balanceSheet.IntangAssets = models.NewJSONDecimal(c.parseDecimalFromInterface(intangAssets))
	}

	// 负债 - 使用实际API字段名
	if totalLiab, ok := data["TOTAL_LIABILITIES"]; ok {
		balanceSheet.TotalLiab = models.NewJSONDecimal(c.parseDecimalFromInterface(totalLiab))
	}
	if totalCurLiab, ok := data["TOTAL_CURRENT_LIAB"]; ok {
		balanceSheet.TotalCurLiab = models.NewJSONDecimal(c.parseDecimalFromInterface(totalCurLiab))
	}
	if shortLoan, ok := data["SHORT_LOAN"]; ok {
		balanceSheet.ShortLoan = models.NewJSONDecimal(c.parseDecimalFromInterface(shortLoan))
	}
	if notesPayable, ok := data["NOTE_ACCOUNTS_PAYABLE"]; ok {
		balanceSheet.NotesPayable = models.NewJSONDecimal(c.parseDecimalFromInterface(notesPayable))
	}
	if accountsPayable, ok := data["ACCOUNTS_PAYABLE"]; ok {
		balanceSheet.AccountsPayable = models.NewJSONDecimal(c.parseDecimalFromInterface(accountsPayable))
	}
	if totalNcaLiab, ok := data["TOTAL_NONCURRENT_LIAB"]; ok {
		balanceSheet.TotalNcaLiab = models.NewJSONDecimal(c.parseDecimalFromInterface(totalNcaLiab))
	}
	if longLoan, ok := data["LONG_LOAN"]; ok {
		balanceSheet.LongLoan = models.NewJSONDecimal(c.parseDecimalFromInterface(longLoan))
	}

	// 所有者权益 - 使用实际API字段名
	if totalHldrEqy, ok := data["TOTAL_EQUITY"]; ok {
		balanceSheet.TotalHldrEqy = models.NewJSONDecimal(c.parseDecimalFromInterface(totalHldrEqy))
	}
	if capRese, ok := data["CAP_RESE"]; ok {
		balanceSheet.CapRese = models.NewJSONDecimal(c.parseDecimalFromInterface(capRese))
	}
	if undistrProfit, ok := data["UNDISTR_PORFIT"]; ok {
		balanceSheet.UndistrProfit = models.NewJSONDecimal(c.parseDecimalFromInterface(undistrProfit))
	}
	if totalShare, ok := data["SHARE_CAP"]; ok {
		balanceSheet.TotalShare = models.NewJSONDecimal(c.parseDecimalFromInterface(totalShare))
	}

	return balanceSheet, nil
}

// convertToCashFlowStatement 将AKTools现金流量表数据转换为内部模型
func (c *AKToolsClient) convertToCashFlowStatement(data map[string]interface{}, symbol, period, reportType string) (*models.CashFlowStatement, error) {
	cashFlowStatement := &models.CashFlowStatement{}

	// 设置基础字段 - 根据实际API响应
	if secucode, ok := data["SECUCODE"].(string); ok {
		cashFlowStatement.TSCode = secucode
	} else {
		cashFlowStatement.TSCode = c.DetermineTSCode(symbol)
	}

	// 报告期 - 从REPORT_DATE字段提取
	if reportDate, ok := data["REPORT_DATE"].(string); ok {
		cashFlowStatement.FDate = c.formatDateForFrontend(reportDate)
		cashFlowStatement.EndDate = c.formatDateForFrontend(reportDate)
	} else {
		cashFlowStatement.FDate = period
		cashFlowStatement.EndDate = period
	}

	// 从REPORT_TYPE字段获取报告类型
	if rptType, ok := data["REPORT_TYPE"].(string); ok {
		cashFlowStatement.ReportType = rptType
	} else {
		cashFlowStatement.ReportType = reportType
	}

	// 提取公告日期 - 从NOTICE_DATE字段
	if noticeDate, ok := data["NOTICE_DATE"].(string); ok {
		cashFlowStatement.AnnDate = c.formatDateForFrontend(noticeDate)
	}

	// 经营活动现金流量 - 使用实际API字段名
	if netCashOperAct, ok := data["NETCASH_OPERATE"]; ok {
		cashFlowStatement.NetCashOperAct = models.NewJSONDecimal(c.parseDecimalFromInterface(netCashOperAct))
	}
	if salesCash, ok := data["SALES_SERVICES"]; ok {
		cashFlowStatement.CashRecrSale = models.NewJSONDecimal(c.parseDecimalFromInterface(salesCash))
	}
	if buyCash, ok := data["BUY_SERVICES"]; ok {
		cashFlowStatement.CashPayGoods = models.NewJSONDecimal(c.parseDecimalFromInterface(buyCash))
	}
	if payStaffCash, ok := data["PAY_STAFF_CASH"]; ok {
		cashFlowStatement.CashPayBehalfEmpl = models.NewJSONDecimal(c.parseDecimalFromInterface(payStaffCash))
	}
	if payTaxCash, ok := data["PAY_ALL_TAX"]; ok {
		cashFlowStatement.CashPayTax = models.NewJSONDecimal(c.parseDecimalFromInterface(payTaxCash))
	}

	// 投资活动现金流量 - 使用实际API字段名
	if netCashInvAct, ok := data["NETCASH_INVEST"]; ok {
		cashFlowStatement.NetCashInvAct = models.NewJSONDecimal(c.parseDecimalFromInterface(netCashInvAct))
	}
	if withdrawInvest, ok := data["WITHDRAW_INVEST"]; ok {
		cashFlowStatement.CashRecvDisp = models.NewJSONDecimal(c.parseDecimalFromInterface(withdrawInvest))
	}
	if investPayCash, ok := data["INVEST_PAY_CASH"]; ok {
		cashFlowStatement.CashPayAcq = models.NewJSONDecimal(c.parseDecimalFromInterface(investPayCash))
	}

	// 筹资活动现金流量 - 使用实际API字段名
	if netCashFinAct, ok := data["NETCASH_FINANCE"]; ok {
		cashFlowStatement.NetCashFinAct = models.NewJSONDecimal(c.parseDecimalFromInterface(netCashFinAct))
	}
	if acceptInvestCash, ok := data["ACCEPT_INVEST_CASH"]; ok {
		cashFlowStatement.CashRecvInvest = models.NewJSONDecimal(c.parseDecimalFromInterface(acceptInvestCash))
	}
	if assignDividend, ok := data["ASSIGN_DIVIDEND_PORFIT"]; ok {
		cashFlowStatement.CashPayDist = models.NewJSONDecimal(c.parseDecimalFromInterface(assignDividend))
	}

	// 汇率变动影响 - 使用实际API字段名
	if fxEffectCash, ok := data["RATE_CHANGE_EFFECT"]; ok {
		cashFlowStatement.FxEffectCash = models.NewJSONDecimal(c.parseDecimalFromInterface(fxEffectCash))
	}

	// 现金净增加额和期初期末余额
	if netIncrCash, ok := data["CCE_ADD"]; ok {
		cashFlowStatement.NetIncrCashCce = models.NewJSONDecimal(c.parseDecimalFromInterface(netIncrCash))
	}
	if beginCash, ok := data["BEGIN_CCE"]; ok {
		cashFlowStatement.CashBegPeriod = models.NewJSONDecimal(c.parseDecimalFromInterface(beginCash))
	}
	if endCash, ok := data["END_CCE"]; ok {
		cashFlowStatement.CashEndPeriod = models.NewJSONDecimal(c.parseDecimalFromInterface(endCash))
	}

	return cashFlowStatement, nil
}

// convertToDailyBasic 将AKTools每日基本面数据转换为内部模型
func (c *AKToolsClient) convertToDailyBasic(data map[string]interface{}, symbol, tradeDate string) (*models.DailyBasic, error) {
	dailyBasic := &models.DailyBasic{}

	// 设置基础字段
	dailyBasic.TSCode = c.DetermineTSCode(symbol)
	dailyBasic.TradeDate = tradeDate

	// 基本数据 - 使用实际API字段名
	if close, ok := data["最新价"]; ok {
		dailyBasic.Close = models.NewJSONDecimal(c.parseDecimalFromInterface(close))
	}
	if turnover, ok := data["换手率"]; ok {
		dailyBasic.Turnover = models.NewJSONDecimal(c.parseDecimalFromInterface(turnover))
	}
	if volumeRatio, ok := data["量比"]; ok {
		dailyBasic.VolumeRatio = models.NewJSONDecimal(c.parseDecimalFromInterface(volumeRatio))
	}

	// 估值指标 - 使用实际API字段名
	if pe, ok := data["市盈率-动态"]; ok {
		dailyBasic.Pe = models.NewJSONDecimal(c.parseDecimalFromInterface(pe))
	}
	if peTtm, ok := data["市盈率TTM"]; ok {
		dailyBasic.PeTtm = models.NewJSONDecimal(c.parseDecimalFromInterface(peTtm))
	}
	if pb, ok := data["市净率"]; ok {
		dailyBasic.Pb = models.NewJSONDecimal(c.parseDecimalFromInterface(pb))
	}
	if ps, ok := data["市销率"]; ok {
		dailyBasic.Ps = models.NewJSONDecimal(c.parseDecimalFromInterface(ps))
	}
	if psTtm, ok := data["市销率TTM"]; ok {
		dailyBasic.PsTtm = models.NewJSONDecimal(c.parseDecimalFromInterface(psTtm))
	}

	// 股本和市值
	if totalShare, ok := data["总股本"]; ok {
		dailyBasic.TotalShare = models.NewJSONDecimal(c.parseDecimalFromInterface(totalShare))
	}
	if floatShare, ok := data["流通股本"]; ok {
		dailyBasic.FloatShare = models.NewJSONDecimal(c.parseDecimalFromInterface(floatShare))
	}
	if totalMv, ok := data["总市值"]; ok {
		dailyBasic.TotalMv = models.NewJSONDecimal(c.parseDecimalFromInterface(totalMv))
	}
	if circMv, ok := data["流通市值"]; ok {
		dailyBasic.CircMv = models.NewJSONDecimal(c.parseDecimalFromInterface(circMv))
	}

	// 分红指标
	if dvRatio, ok := data["股息率"]; ok {
		dailyBasic.DvRatio = models.NewJSONDecimal(c.parseDecimalFromInterface(dvRatio))
	}
	if dvTtm, ok := data["股息率TTM"]; ok {
		dailyBasic.DvTtm = models.NewJSONDecimal(c.parseDecimalFromInterface(dvTtm))
	}

	return dailyBasic, nil
}

// parseDecimalFromInterface 从interface{}解析decimal值
func (c *AKToolsClient) parseDecimalFromInterface(value interface{}) decimal.Decimal {
	if value == nil {
		return decimal.Zero
	}

	switch v := value.(type) {
	case float64:
		return decimal.NewFromFloat(v)
	case float32:
		return decimal.NewFromFloat(float64(v))
	case int:
		return decimal.NewFromInt(int64(v))
	case int64:
		return decimal.NewFromInt(v)
	case string:
		if v == "" || v == "-" || v == "--" {
			return decimal.Zero
		}
		if d, err := decimal.NewFromString(v); err == nil {
			return d
		}
	}
	return decimal.Zero
}
