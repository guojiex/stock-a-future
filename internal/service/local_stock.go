package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"stock-a-future/internal/models"

	"github.com/xuri/excelize/v2"
)

// LocalStockService 本地股票数据服务
type LocalStockService struct {
	stockCache map[string]*models.StockBasic
	mutex      sync.RWMutex
	dataDir    string
}

// NewLocalStockService 创建本地股票数据服务
func NewLocalStockService(dataDir string) *LocalStockService {
	if dataDir == "" {
		dataDir = "data"
	}

	service := &LocalStockService{
		stockCache: make(map[string]*models.StockBasic),
		dataDir:    dataDir,
	}

	// 启动时加载数据
	if err := service.LoadStockData(); err != nil {
		log.Printf("警告: 加载本地股票数据失败: %v", err)
	}

	return service
}

// LoadStockData 加载股票数据到缓存
func (s *LocalStockService) LoadStockData() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 清空现有缓存
	s.stockCache = make(map[string]*models.StockBasic)

	// 加载上交所数据
	if err := s.loadSSEData(); err != nil {
		log.Printf("加载上交所数据失败: %v", err)
	}

	// 加载深交所Excel数据
	if err := s.loadSZSEExcelData(); err != nil {
		log.Printf("加载深交所Excel数据失败: %v", err)
	}

	log.Printf("本地股票数据加载完成，共 %d 只股票", len(s.stockCache))
	return nil
}

// loadSSEData 加载上交所JSON数据
func (s *LocalStockService) loadSSEData() error {
	filename := fmt.Sprintf("%s/sse_stocks.json", s.dataDir)

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("读取上交所数据文件失败: %w", err)
	}

	var stocks []models.StockBasic
	if err := json.Unmarshal(data, &stocks); err != nil {
		return fmt.Errorf("解析上交所JSON数据失败: %w", err)
	}

	// 添加到缓存
	for _, stock := range stocks {
		// 使用股票代码作为key，支持多种格式查询
		s.stockCache[stock.TSCode] = &stock
		s.stockCache[stock.Symbol] = &stock

		// 如果TSCode包含.SH后缀，也添加不含后缀的版本
		if strings.Contains(stock.TSCode, ".") {
			code := strings.Split(stock.TSCode, ".")[0]
			s.stockCache[code] = &stock
		}
	}

	log.Printf("加载上交所数据: %d 只股票", len(stocks))
	return nil
}

// loadSZSEExcelData 加载深交所Excel数据
func (s *LocalStockService) loadSZSEExcelData() error {
	filename := fmt.Sprintf("%s/A股列表.xlsx", s.dataDir)

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("深交所Excel文件不存在: %s", filename)
	}

	// 打开Excel文件
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("打开Excel文件失败: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("关闭Excel文件失败: %v", err)
		}
	}()

	// 获取第一个工作表名称
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return fmt.Errorf("无法获取Excel工作表")
	}

	// 获取所有行数据
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("读取Excel行数据失败: %w", err)
	}

	if len(rows) <= 1 {
		return fmt.Errorf("Excel文件没有有效数据")
	}

	// 分析表头，查找股票代码和名称列
	header := rows[0]
	codeCol, nameCol := s.findExcelColumns(header)

	if codeCol == -1 || nameCol == -1 {
		log.Printf("Excel表头: %v", header)
		return fmt.Errorf("无法找到股票代码或名称列")
	}

	log.Printf("Excel列映射: 代码列=%d, 名称列=%d", codeCol, nameCol)

	// 处理数据行
	stockCount := 0
	for i, row := range rows[1:] { // 跳过表头
		if len(row) <= codeCol || len(row) <= nameCol {
			continue
		}

		code := strings.TrimSpace(row[codeCol])
		name := strings.TrimSpace(row[nameCol])

		if code == "" || name == "" {
			continue
		}

		// 验证是否是有效的深交所股票代码
		if !s.isValidSZSECode(code) {
			continue
		}

		// 创建股票信息
		stock := &models.StockBasic{
			TSCode:   code + ".SZ",
			Symbol:   code,
			Name:     name,
			Market:   "SZ",
			Area:     "",
			Industry: "",
			ListDate: "",
		}

		// 添加到缓存
		s.stockCache[stock.TSCode] = stock
		s.stockCache[stock.Symbol] = stock
		s.stockCache[code] = stock

		stockCount++

		// 记录前几条数据用于调试
		if i < 5 {
			log.Printf("Excel数据示例: %s -> %s (%s)", code, name, stock.TSCode)
		}
	}

	log.Printf("加载深交所Excel数据: %d 只股票", stockCount)
	return nil
}

// findExcelColumns 查找Excel中的股票代码和名称列
func (s *LocalStockService) findExcelColumns(header []string) (codeCol, nameCol int) {
	codeCol, nameCol = -1, -1

	// 可能的列名变体（按优先级排序）
	codeNames := []string{"A股代码", "股票代码", "代码", "证券代码", "股代码", "code", "Code", "CODE", "symbol", "Symbol"}
	// 优先选择中文名称列
	nameNames := []string{"A股简称", "股票简称", "简称", "股票名称", "名称", "证券简称", "证券名称", "中文简称", "中文名称"}
	// 英文名称作为备选
	englishNames := []string{"英文名称", "name", "Name", "NAME", "stock_name", "english_name"}

	// 记录找到的所有名称列
	var nameColumns []int
	var chineseNameCol = -1

	for i, col := range header {
		colName := strings.TrimSpace(col)

		// 查找代码列
		if codeCol == -1 {
			for _, codeName := range codeNames {
				if strings.Contains(colName, codeName) {
					codeCol = i
					break
				}
			}
		}

		// 查找中文名称列（优先）- 使用精确匹配避免误匹配
		if chineseNameCol == -1 {
			for _, nameName := range nameNames {
				if colName == nameName {
					chineseNameCol = i
					break
				}
			}
		}

		// 查找所有可能的名称列
		for _, nameName := range append(nameNames, englishNames...) {
			if strings.Contains(colName, nameName) {
				nameColumns = append(nameColumns, i)
				break
			}
		}
	}

	// 优先使用中文名称列
	if chineseNameCol != -1 {
		nameCol = chineseNameCol
	} else if len(nameColumns) > 0 {
		// 如果没有中文列，使用第一个找到的名称列
		nameCol = nameColumns[0]
	}

	return codeCol, nameCol
}

// isValidSZSECode 验证是否是有效的深交所股票代码
func (s *LocalStockService) isValidSZSECode(code string) bool {
	if len(code) != 6 {
		return false
	}

	// 深交所主板: 000xxx, 001xxx
	// 中小板: 002xxx
	// 创业板: 300xxx
	// B股: 200xxx
	return strings.HasPrefix(code, "00") || strings.HasPrefix(code, "30") || strings.HasPrefix(code, "20")
}

// GetStockBasic 获取股票基本信息
func (s *LocalStockService) GetStockBasic(stockCode string) (*models.StockBasic, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 标准化股票代码
	normalizedCode := s.normalizeStockCode(stockCode)

	// 尝试多种格式查找
	searchCodes := []string{
		normalizedCode,
		stockCode,
		strings.ToUpper(stockCode),
		strings.ToLower(stockCode),
	}

	// 如果没有后缀，尝试添加后缀
	if !strings.Contains(stockCode, ".") {
		if len(stockCode) == 6 {
			// 根据代码前缀判断交易所
			if strings.HasPrefix(stockCode, "60") || strings.HasPrefix(stockCode, "68") || strings.HasPrefix(stockCode, "90") {
				searchCodes = append(searchCodes, stockCode+".SH")
			} else if strings.HasPrefix(stockCode, "00") || strings.HasPrefix(stockCode, "30") || strings.HasPrefix(stockCode, "20") {
				searchCodes = append(searchCodes, stockCode+".SZ")
			}
		}
	}

	for _, code := range searchCodes {
		if stock, exists := s.stockCache[code]; exists {
			return stock, nil
		}
	}

	// 如果在缓存中没找到，创建一个基本的股票信息
	return s.createBasicStockInfo(stockCode), nil
}

// normalizeStockCode 标准化股票代码
func (s *LocalStockService) normalizeStockCode(code string) string {
	code = strings.TrimSpace(code)
	code = strings.ToUpper(code)
	return code
}

// createBasicStockInfo 创建基本股票信息
func (s *LocalStockService) createBasicStockInfo(stockCode string) *models.StockBasic {
	// 标准化代码
	normalizedCode := s.normalizeStockCode(stockCode)

	var tsCode, symbol, market, name string

	if strings.Contains(normalizedCode, ".") {
		parts := strings.Split(normalizedCode, ".")
		symbol = parts[0]
		market = parts[1]
		tsCode = normalizedCode
	} else {
		symbol = normalizedCode
		// 根据代码前缀判断市场
		if len(symbol) == 6 {
			if strings.HasPrefix(symbol, "60") || strings.HasPrefix(symbol, "68") || strings.HasPrefix(symbol, "90") {
				market = "SH"
				tsCode = symbol + ".SH"
			} else if strings.HasPrefix(symbol, "00") || strings.HasPrefix(symbol, "30") || strings.HasPrefix(symbol, "20") {
				market = "SZ"
				tsCode = symbol + ".SZ"
			} else {
				market = "UNKNOWN"
				tsCode = symbol
			}
		} else {
			market = "UNKNOWN"
			tsCode = symbol
		}
	}

	// 生成基本名称
	name = s.generateStockName(symbol, market)

	return &models.StockBasic{
		TSCode:   tsCode,
		Symbol:   symbol,
		Name:     name,
		Market:   market,
		Area:     "",
		Industry: "",
		ListDate: "",
	}
}

// generateStockName 生成股票名称
func (s *LocalStockService) generateStockName(symbol, market string) string {
	// 常见股票的映射表（作为后备方案）
	stockNameMap := map[string]string{
		// 上交所主要股票
		"600000": "浦发银行",
		"600036": "招商银行",
		"600519": "贵州茅台",
		"600887": "伊利股份",
		"600276": "恒瑞医药",
		"601318": "中国平安",
		"601166": "兴业银行",
		"601398": "工商银行",
		"601288": "农业银行",
		"601939": "建设银行",
		"601857": "中国石油",
		"600028": "中国石化",
		"688111": "金山办公",
		"688981": "中芯国际",

		// 深交所主要股票
		"000001": "平安银行",
		"000002": "万科A",
		"000858": "五粮液",
		"002415": "海康威视",
		"002594": "比亚迪",
		"300059": "东方财富",
		"300750": "宁德时代",
		"300015": "爱尔眼科",
		"300760": "迈瑞医疗",
	}

	if name, exists := stockNameMap[symbol]; exists {
		return name
	}

	// 如果没有找到映射，生成默认名称
	marketName := ""
	switch market {
	case "SH":
		marketName = "上交所"
	case "SZ":
		marketName = "深交所"
	default:
		marketName = "股票"
	}

	return fmt.Sprintf("%s%s", marketName, symbol)
}

// GetAllStocks 获取所有股票列表
func (s *LocalStockService) GetAllStocks() []*models.StockBasic {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 去重并返回所有股票
	uniqueStocks := make(map[string]*models.StockBasic)

	for _, stock := range s.stockCache {
		uniqueStocks[stock.TSCode] = stock
	}

	result := make([]*models.StockBasic, 0, len(uniqueStocks))
	for _, stock := range uniqueStocks {
		result = append(result, stock)
	}

	return result
}

// RefreshData 刷新数据
func (s *LocalStockService) RefreshData() error {
	return s.LoadStockData()
}

// GetStockCount 获取股票数量
func (s *LocalStockService) GetStockCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 计算唯一股票数量
	uniqueStocks := make(map[string]bool)
	for _, stock := range s.stockCache {
		uniqueStocks[stock.TSCode] = true
	}

	return len(uniqueStocks)
}

// SearchStocks 根据关键词搜索股票
func (s *LocalStockService) SearchStocks(keyword string, limit int) []*models.StockBasic {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if keyword == "" {
		return []*models.StockBasic{}
	}

	keyword = strings.ToUpper(strings.TrimSpace(keyword))
	var results []*models.StockBasic
	uniqueStocks := make(map[string]*models.StockBasic)

	// 收集所有唯一的股票
	for _, stock := range s.stockCache {
		uniqueStocks[stock.TSCode] = stock
	}

	// 搜索匹配的股票
	for _, stock := range uniqueStocks {
		// 搜索股票代码
		if strings.Contains(strings.ToUpper(stock.TSCode), keyword) ||
			strings.Contains(strings.ToUpper(stock.Symbol), keyword) {
			results = append(results, stock)
			continue
		}

		// 搜索股票名称
		if strings.Contains(strings.ToUpper(stock.Name), keyword) {
			results = append(results, stock)
			continue
		}
	}

	// 限制返回结果数量
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results
}
