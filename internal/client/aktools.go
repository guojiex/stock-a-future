package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"stock-a-future/internal/models"

	"github.com/shopspring/decimal"
)

// AKToolsClient AKTools HTTP API客户端
type AKToolsClient struct {
	baseURL string
	client  *http.Client
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

	return &AKToolsClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
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

// GetDailyData 获取股票日线数据
func (c *AKToolsClient) GetDailyData(symbol, startDate, endDate string) ([]models.StockDaily, error) {
	// 清理股票代码，移除市场后缀
	cleanSymbol := c.CleanStockSymbol(symbol)

	// 构建查询参数
	params := url.Values{}
	params.Set("symbol", cleanSymbol)
	if startDate != "" {
		params.Set("start_date", startDate)
	}
	if endDate != "" {
		params.Set("end_date", endDate)
	}
	// 使用默认参数
	params.Set("period", "daily")
	params.Set("adjust", "hfq") // 默认后复权

	// 构建完整URL
	apiURL := fmt.Sprintf("%s/api/public/stock_zh_a_hist?%s", c.baseURL, params.Encode())

	// 发送HTTP请求
	resp, err := c.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools API失败: %w", err)
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
	apiURL := fmt.Sprintf("%s/api/public/stock_zh_a_info?%s", c.baseURL, params.Encode())

	// 发送HTTP请求
	resp, err := c.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools股票信息API失败: %w", err)
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

	// 解析JSON响应
	var aktoolsResp []AKToolsStockBasicResponse
	if err := json.Unmarshal(body, &aktoolsResp); err != nil {
		return nil, fmt.Errorf("解析AKTools股票信息响应失败: %w", err)
	}

	if len(aktoolsResp) == 0 {
		return nil, fmt.Errorf("未找到股票基本信息: %s", symbol)
	}

	// 转换为内部模型
	return c.convertToStockBasic(aktoolsResp[0], symbol), nil
}

// GetStockList 获取股票列表
func (c *AKToolsClient) GetStockList() ([]models.StockBasic, error) {
	// 构建完整URL - 使用股票列表API
	apiURL := fmt.Sprintf("%s/api/public/stock_zh_a_spot", c.baseURL)

	// 发送HTTP请求
	resp, err := c.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求AKTools股票列表API失败: %w", err)
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

		daily := models.StockDaily{
			TSCode:    tsCode,
			TradeDate: data.Date,
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

// GetBaseURL 获取AKTools API基础URL
func (c *AKToolsClient) GetBaseURL() string {
	return c.baseURL
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

	// 发送测试请求
	resp, err := testClient.Get(apiURL)
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
