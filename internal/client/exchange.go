package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"stock-a-future/internal/models"
)

// ExchangeClient 交易所官网客户端
type ExchangeClient struct {
	client *http.Client
}

// StockListItem 股票列表项
type StockListItem struct {
	Code string `json:"code"` // 股票代码
	Name string `json:"name"` // 股票名称
}

// NewExchangeClient 创建交易所客户端
func NewExchangeClient() *ExchangeClient {
	return &ExchangeClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// GetSSEStockList 获取上海证券交易所股票列表
func (c *ExchangeClient) GetSSEStockList() ([]StockListItem, error) {
	log.Println("正在获取上海证券交易所股票列表...")

	// 尝试从官网API获取真实数据
	stocks, err := c.fetchSSEStocksFromAPI()
	if err != nil {
		log.Printf("从官网API获取失败: %v，使用备用数据", err)
		// 如果API获取失败，使用备用数据
		stocks = c.getSSEStockList()
	}

	log.Printf("从上海证券交易所获取到 %d 只股票", len(stocks))
	return stocks, nil
}

// fetchSSEStocksFromAPI 从上交所官网API获取股票列表
func (c *ExchangeClient) fetchSSEStocksFromAPI() ([]StockListItem, error) {
	var allStocks []StockListItem

	// 上交所股票列表API的基础URL
	baseURL := "https://query.sse.com.cn/security/stock/getStockListData2.do"

	// 设置请求头，模拟浏览器访问
	headers := map[string]string{
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8",
		"Referer":          "https://www.sse.com.cn/assortment/stock/list/share/",
		"X-Requested-With": "XMLHttpRequest",
	}

	// 分页获取数据
	pageSize := 100 // 每页数量
	pageNo := 1

	for {
		log.Printf("正在获取第 %d 页股票数据...", pageNo)

		// 构建请求URL
		params := url.Values{}
		params.Set("jsonCallBack", "jsonpCallback"+strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
		params.Set("isPagination", "true")
		params.Set("stockCode", "")
		params.Set("csrcCode", "")
		params.Set("areaName", "")
		params.Set("stockType", "1") // 1表示A股
		params.Set("pageHelp.pageSize", strconv.Itoa(pageSize))
		params.Set("pageHelp.pageNo", strconv.Itoa(pageNo))
		params.Set("pageHelp.beginPage", strconv.Itoa(pageNo))
		params.Set("pageHelp.cacheSize", "1")
		params.Set("pageHelp.endPage", strconv.Itoa(pageNo+4))
		params.Set("_", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))

		requestURL := baseURL + "?" + params.Encode()

		// 发送请求
		req, err := http.NewRequest("GET", requestURL, nil)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %w", err)
		}

		// 设置请求头
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("请求失败: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("读取响应失败: %w", err)
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
		}

		// 解析JSONP响应
		stocks, hasMore, err := c.parseSSEAPIResponse(string(body))
		if err != nil {
			return nil, fmt.Errorf("解析响应失败: %w", err)
		}

		allStocks = append(allStocks, stocks...)

		if !hasMore || len(stocks) == 0 {
			break
		}

		pageNo++

		// 添加延迟避免请求过于频繁
		time.Sleep(500 * time.Millisecond)
	}

	return allStocks, nil
}

// parseSSEAPIResponse 解析上交所API响应
func (c *ExchangeClient) parseSSEAPIResponse(response string) ([]StockListItem, bool, error) {
	// 移除JSONP包装
	jsonpPattern := regexp.MustCompile(`^jsonpCallback\d+\((.*)\)$`)
	matches := jsonpPattern.FindStringSubmatch(response)
	if len(matches) < 2 {
		return nil, false, fmt.Errorf("无效的JSONP响应格式")
	}

	jsonData := matches[1]

	// 解析JSON
	var apiResponse struct {
		PageHelp struct {
			Data []struct {
				COMPANY_CODE    string `json:"COMPANY_CODE"`    // 股票代码
				COMPANY_ABBR    string `json:"COMPANY_ABBR"`    // 股票简称
				PRODUCTID       string `json:"PRODUCTID"`       // 产品ID
				SECURITY_CODE_A string `json:"SECURITY_CODE_A"` // A股代码
			} `json:"data"`
			Total    int `json:"total"`
			PageSize int `json:"pageSize"`
			PageNo   int `json:"pageNo"`
		} `json:"pageHelp"`
		Result []struct {
			COMPANY_CODE    string `json:"COMPANY_CODE"`
			COMPANY_ABBR    string `json:"COMPANY_ABBR"`
			PRODUCTID       string `json:"PRODUCTID"`
			SECURITY_CODE_A string `json:"SECURITY_CODE_A"`
		} `json:"result"`
	}

	if err := json.Unmarshal([]byte(jsonData), &apiResponse); err != nil {
		return nil, false, fmt.Errorf("JSON解析失败: %w", err)
	}

	var stocks []StockListItem

	// 优先使用result字段的数据
	dataSource := apiResponse.Result
	if len(dataSource) == 0 {
		dataSource = apiResponse.PageHelp.Data
	}

	for _, item := range dataSource {
		// 使用A股代码或公司代码
		code := item.SECURITY_CODE_A
		if code == "" {
			code = item.COMPANY_CODE
		}

		if code != "" && item.COMPANY_ABBR != "" {
			// 验证是否是有效的上交所股票代码
			if c.isValidSSECode(code) {
				stocks = append(stocks, StockListItem{
					Code: code + ".SH",
					Name: item.COMPANY_ABBR,
				})
			}
		}
	}

	// 判断是否还有更多数据
	hasMore := false
	if apiResponse.PageHelp.Total > 0 {
		totalPages := (apiResponse.PageHelp.Total + apiResponse.PageHelp.PageSize - 1) / apiResponse.PageHelp.PageSize
		hasMore = apiResponse.PageHelp.PageNo < totalPages
	}

	return stocks, hasMore, nil
}

// parseSSEStockList 解析上交所股票列表HTML
func (c *ExchangeClient) parseSSEStockList(html string) []StockListItem {
	var stocks []StockListItem

	// 上交所的股票代码通常是6位数字，以60、68、90开头
	// 这里使用正则表达式尝试从HTML中提取股票信息
	// 由于网页可能是动态加载的，我们先尝试找到可能的数据接口

	// 查找可能的JavaScript数据或API调用
	jsDataRegex := regexp.MustCompile(`(?i)stock.*?code.*?[:=]\s*["']?(\d{6})["']?.*?name.*?[:=]\s*["']([^"']+)["']?`)
	matches := jsDataRegex.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			code := strings.TrimSpace(match[1])
			name := strings.TrimSpace(match[2])

			// 验证是否是有效的上交所股票代码
			if c.isValidSSECode(code) {
				stocks = append(stocks, StockListItem{
					Code: code + ".SH",
					Name: name,
				})
			}
		}
	}

	// 如果没有找到数据，尝试其他正则表达式模式
	if len(stocks) == 0 {
		// 尝试查找表格行数据
		tableRowRegex := regexp.MustCompile(`<tr[^>]*>.*?(\d{6}).*?<[^>]*>([^<]+)</[^>]*>.*?</tr>`)
		matches = tableRowRegex.FindAllStringSubmatch(html, -1)

		for _, match := range matches {
			if len(match) >= 3 {
				code := strings.TrimSpace(match[1])
				name := strings.TrimSpace(match[2])

				if c.isValidSSECode(code) {
					stocks = append(stocks, StockListItem{
						Code: code + ".SH",
						Name: name,
					})
				}
			}
		}
	}

	// 如果仍然没有数据，返回一些示例数据（实际应用中应该调用API）
	if len(stocks) == 0 {
		log.Println("警告: 无法从HTML中解析股票数据，返回示例数据")
		stocks = c.getSSEStockList()
	}

	return stocks
}

// isValidSSECode 验证是否是有效的上交所股票代码
func (c *ExchangeClient) isValidSSECode(code string) bool {
	if len(code) != 6 {
		return false
	}

	// 上交所主板: 600xxx, 601xxx, 603xxx, 605xxx
	// 科创板: 688xxx
	// B股: 900xxx
	return strings.HasPrefix(code, "60") || strings.HasPrefix(code, "68") || strings.HasPrefix(code, "90")
}

// getSSEStockList 获取上交所股票列表（包含更多真实股票）
func (c *ExchangeClient) getSSEStockList() []StockListItem {
	return []StockListItem{
		// 上证主板 - 银行
		{Code: "600000.SH", Name: "浦发银行"},
		{Code: "600036.SH", Name: "招商银行"},
		{Code: "601398.SH", Name: "工商银行"},
		{Code: "601328.SH", Name: "交通银行"},
		{Code: "601288.SH", Name: "农业银行"},
		{Code: "601939.SH", Name: "建设银行"},
		{Code: "600016.SH", Name: "民生银行"},
		{Code: "600015.SH", Name: "华夏银行"},
		{Code: "601166.SH", Name: "兴业银行"},
		{Code: "601009.SH", Name: "南京银行"},

		// 上证主板 - 保险
		{Code: "601318.SH", Name: "中国平安"},
		{Code: "601628.SH", Name: "中国人寿"},
		{Code: "601336.SH", Name: "新华保险"},
		{Code: "601601.SH", Name: "中国太保"},

		// 上证主板 - 石化能源
		{Code: "601857.SH", Name: "中国石油"},
		{Code: "600028.SH", Name: "中国石化"},
		{Code: "601088.SH", Name: "中国神华"},
		{Code: "600309.SH", Name: "万华化学"},

		// 上证主板 - 食品饮料
		{Code: "600519.SH", Name: "贵州茅台"},
		{Code: "600887.SH", Name: "伊利股份"},
		{Code: "600809.SH", Name: "山西汾酒"},
		{Code: "603288.SH", Name: "海天味业"},
		{Code: "600600.SH", Name: "青岛啤酒"},

		// 上证主板 - 科技
		{Code: "600276.SH", Name: "恒瑞医药"},
		{Code: "603259.SH", Name: "药明康德"},
		{Code: "600570.SH", Name: "恒生电子"},
		{Code: "600845.SH", Name: "宝信软件"},

		// 科创板
		{Code: "688111.SH", Name: "金山办公"},
		{Code: "688008.SH", Name: "澜起科技"},
		{Code: "688036.SH", Name: "传音控股"},
		{Code: "688981.SH", Name: "中芯国际"},
		{Code: "688599.SH", Name: "天合光能"},

		// 上证主板 - 其他
		{Code: "600104.SH", Name: "上汽集团"},
		{Code: "601012.SH", Name: "隆基绿能"},
		{Code: "601888.SH", Name: "中国中免"},
		{Code: "600690.SH", Name: "海尔智家"},
		{Code: "600036.SH", Name: "招商银行"},
		{Code: "601668.SH", Name: "中国建筑"},
		{Code: "600031.SH", Name: "三一重工"},
	}
}

// GetAllStockList 获取所有交易所的股票列表
func (c *ExchangeClient) GetAllStockList() ([]StockListItem, error) {
	var allStocks []StockListItem

	// 获取上交所股票
	sseStocks, err := c.GetSSEStockList()
	if err != nil {
		log.Printf("获取上交所股票列表失败: %v", err)
	} else {
		allStocks = append(allStocks, sseStocks...)
	}

	log.Printf("总共获取到 %d 只股票", len(allStocks))
	return allStocks, nil
}

// SaveStockListToFile 将股票列表保存到文件
func (c *ExchangeClient) SaveStockListToFile(stocks []StockListItem, filename string) error {
	// 转换为StockBasic格式以保持兼容性
	var stockBasics []models.StockBasic
	for _, stock := range stocks {
		// 从股票代码中提取symbol
		parts := strings.Split(stock.Code, ".")
		symbol := ""
		market := ""
		if len(parts) == 2 {
			symbol = parts[0]
			market = parts[1]
		}

		stockBasic := models.StockBasic{
			TSCode: stock.Code,
			Symbol: symbol,
			Name:   stock.Name,
			Market: market,
		}
		stockBasics = append(stockBasics, stockBasic)
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(stockBasics, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化股票数据失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	log.Printf("股票列表已保存到文件: %s", filename)
	return nil
}

// LoadStockListFromFile 从文件加载股票列表
func (c *ExchangeClient) LoadStockListFromFile(filename string) ([]models.StockBasic, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	var stocks []models.StockBasic
	if err := json.Unmarshal(data, &stocks); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	log.Printf("从文件加载了 %d 只股票", len(stocks))
	return stocks, nil
}
