package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stock-a-future/config"
	"stock-a-future/internal/client"
	"stock-a-future/internal/indicators"
	"stock-a-future/internal/models"
	"stock-a-future/internal/service"
	"strconv"
	"strings"
	"time"
)

// 记录服务器启动时间，用于健康检查
var startTime = time.Now()

// StockHandler 股票数据处理器
type StockHandler struct {
	dataSourceClient  client.DataSourceClient
	calculator        *indicators.Calculator
	predictionService *service.PredictionService
	localStockService *service.LocalStockService
	dailyCacheService *service.DailyCacheService
	favoriteService   *service.FavoriteService
	config            *config.Config // 添加配置字段
	app               *service.App   // 应用上下文，用于获取其他服务
}

// NewStockHandler 创建股票处理器
func NewStockHandler(dataSourceClient client.DataSourceClient, cacheService *service.DailyCacheService, favoriteService *service.FavoriteService, app *service.App) *StockHandler {
	return &StockHandler{
		dataSourceClient:  dataSourceClient,
		calculator:        indicators.NewCalculator(),
		predictionService: service.NewPredictionService(),
		localStockService: service.NewLocalStockService("data"),
		dailyCacheService: cacheService,
		favoriteService:   favoriteService,
		config:            config.Load(), // 初始化配置
		app:               app,           // 应用上下文
	}
}

// GetDailyData 获取股票日线数据
func (h *StockHandler) GetDailyData(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 记录请求开始
	startTime := time.Now()
	log.Printf("[GetDailyData] 开始处理请求 - 路径: %s, 方法: %s",
		r.URL.Path, r.Method)

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		errorMsg := "无效的股票代码"
		log.Printf("[GetDailyData] 参数解析失败 - 路径: %s, 错误: %s", r.URL.Path, errorMsg)
		h.writeErrorResponse(w, http.StatusBadRequest, errorMsg)
		return
	}

	stockCode := pathParts[3] // /api/v1/stocks/{code}/daily

	// 解析查询参数
	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")
	adjust := query.Get("adjust") // 复权方式：qfq(前复权), hfq(后复权), none(不复权)

	// 默认获取最近30天数据
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("20060102")
	}
	if endDate == "" {
		endDate = time.Now().Format("20060102")
	}

	// 记录请求参数
	log.Printf("[GetDailyData] 请求参数 - 股票代码: %s, 开始日期: %s, 结束日期: %s, 复权方式: %s",
		stockCode, startDate, endDate, adjust)

	// 尝试从缓存获取数据
	var data []models.StockDaily
	var err error

	if h.dailyCacheService != nil {
		log.Printf("[GetDailyData] 尝试从缓存获取数据 - 股票代码: %s", stockCode)
		if cachedData, found := h.dailyCacheService.Get(stockCode, startDate, endDate); found {
			log.Printf("[GetDailyData] 缓存命中 - 股票代码: %s, 数据条数: %d", stockCode, len(cachedData))
			data = cachedData
		} else {
			log.Printf("[GetDailyData] 缓存未命中 - 股票代码: %s, 从API获取数据", stockCode)
			// 缓存未命中，从API获取数据
			data, err = h.dataSourceClient.GetDailyData(stockCode, startDate, endDate, adjust)
			if err != nil {
				// 详细记录错误信息
				log.Printf("[GetDailyData] 获取股票数据失败 - 股票代码: %s, 开始日期: %s, 结束日期: %s, 复权方式: %s, 错误: %v",
					stockCode, startDate, endDate, adjust, err)
				log.Printf("[GetDailyData] 错误详情 - 类型: %T, 消息: %s", err, err.Error())

				// 记录请求上下文信息
				log.Printf("[GetDailyData] 请求上下文 - 远程地址: %s, 引用页: %s",
					r.RemoteAddr, r.Referer())

				h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
				return
			}

			log.Printf("[GetDailyData] 从API获取数据成功 - 股票代码: %s, 数据条数: %d", stockCode, len(data))
			// 将数据存入缓存
			h.dailyCacheService.Set(stockCode, startDate, endDate, data)
			log.Printf("[GetDailyData] 数据已存入缓存 - 股票代码: %s", stockCode)
		}
	} else {
		log.Printf("[GetDailyData] 缓存服务未启用 - 股票代码: %s, 直接从API获取", stockCode)
		// 如果缓存服务未启用，直接从API获取
		data, err = h.dataSourceClient.GetDailyData(stockCode, startDate, endDate, adjust)
		if err != nil {
			// 详细记录错误信息
			log.Printf("[GetDailyData] 获取股票数据失败 - 股票代码: %s, 开始日期: %s, 结束日期: %s, 复权方式: %s, 错误: %v",
				stockCode, startDate, endDate, adjust, err)
			log.Printf("[GetDailyData] 错误详情 - 类型: %T, 消息: %s", err, err.Error())

			// 记录请求上下文信息
			log.Printf("[GetDailyData] 请求上下文 - 远程地址: %s, 引用页: %s",
				r.RemoteAddr, r.Referer())

			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
			return
		}

		log.Printf("[GetDailyData] 从API获取数据成功 - 股票代码: %s, 数据条数: %d", stockCode, len(data))
	}

	// 记录响应信息
	responseTime := time.Since(startTime)
	log.Printf("[GetDailyData] 请求处理完成 - 股票代码: %s, 响应时间: %v, 数据条数: %d",
		stockCode, responseTime, len(data))

	// 返回成功响应
	h.writeSuccessResponse(w, data)
}

// GetIndicators 获取技术指标
func (h *StockHandler) GetIndicators(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 记录请求开始
	startTime := time.Now()
	log.Printf("[GetIndicators] 开始处理请求 - 路径: %s, 方法: %s",
		r.URL.Path, r.Method)

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		errorMsg := "无效的股票代码"
		log.Printf("[GetIndicators] 参数解析失败 - 路径: %s, 错误: %s", r.URL.Path, errorMsg)
		h.writeErrorResponse(w, http.StatusBadRequest, errorMsg)
		return
	}

	stockCode := pathParts[3] // /api/v1/stocks/{code}/indicators

	// 解析查询参数
	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	// 默认获取最近60天数据（计算技术指标需要更多历史数据）
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -60).Format("20060102")
	}
	if endDate == "" {
		endDate = time.Now().Format("20060102")
	}

	// 记录请求参数
	log.Printf("[GetIndicators] 请求参数 - 股票代码: %s, 开始日期: %s, 结束日期: %s",
		stockCode, startDate, endDate)

	// 获取股票数据（优先使用缓存）
	var data []models.StockDaily
	var err error

	if h.dailyCacheService != nil {
		log.Printf("[GetIndicators] 尝试从缓存获取数据 - 股票代码: %s", stockCode)
		if cachedData, found := h.dailyCacheService.Get(stockCode, startDate, endDate); found {
			log.Printf("[GetIndicators] 缓存命中 - 股票代码: %s, 数据条数: %d", stockCode, len(cachedData))
			data = cachedData
		} else {
			log.Printf("[GetIndicators] 缓存未命中 - 股票代码: %s, 从API获取数据", stockCode)
			// 缓存未命中，从API获取数据
			data, err = h.dataSourceClient.GetDailyData(stockCode, startDate, endDate, "")
			if err != nil {
				// 详细记录错误信息
				log.Printf("[GetIndicators] 获取股票数据失败 - 股票代码: %s, 开始日期: %s, 结束日期: %s, 错误: %v",
					stockCode, startDate, endDate, err)
				log.Printf("[GetIndicators] 错误详情 - 类型: %T, 消息: %s", err, err.Error())

				// 记录请求上下文信息
				log.Printf("[GetIndicators] 请求上下文 - 远程地址: %s, 引用页: %s",
					r.RemoteAddr, r.Referer())

				h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
				return
			}

			log.Printf("[GetIndicators] 从API获取数据成功 - 股票代码: %s, 数据条数: %d", stockCode, len(data))
			// 将数据存入缓存
			h.dailyCacheService.Set(stockCode, startDate, endDate, data)
			log.Printf("[GetIndicators] 数据已存入缓存 - 股票代码: %s", stockCode)
		}
	} else {
		log.Printf("[GetIndicators] 缓存服务未启用 - 股票代码: %s, 直接从API获取", stockCode)
		// 如果缓存服务未启用，直接从API获取
		data, err = h.dataSourceClient.GetDailyData(stockCode, startDate, endDate, "")
		if err != nil {
			// 详细记录错误信息
			log.Printf("[GetIndicators] 获取股票数据失败 - 股票代码: %s, 开始日期: %s, 结束日期: %s, 错误: %v",
				stockCode, startDate, endDate, err)
			log.Printf("[GetIndicators] 错误详情 - 类型: %T, 消息: %s", err, err.Error())

			// 记录请求上下文信息
			log.Printf("[GetIndicators] 请求上下文 - 远程地址: %s, 引用页: %s",
				r.RemoteAddr, r.Referer())

			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
			return
		}

		log.Printf("[GetIndicators] 从API获取数据成功 - 股票代码: %s, 数据条数: %d", stockCode, len(data))
	}

	if len(data) == 0 {
		errorMsg := "未找到股票数据"
		log.Printf("[GetIndicators] 数据为空 - 股票代码: %s, 开始日期: %s, 结束日期: %s",
			stockCode, startDate, endDate)
		h.writeErrorResponse(w, http.StatusNotFound, errorMsg)
		return
	}

	log.Printf("[GetIndicators] 开始计算技术指标 - 股票代码: %s, 数据条数: %d", stockCode, len(data))

	// 计算技术指标
	latestData := data[len(data)-1]
	indicators := &models.TechnicalIndicators{
		TSCode:    latestData.TSCode,
		TradeDate: latestData.TradeDate,
	}

	// 计算MACD
	macdResults := h.calculator.CalculateMACD(data)
	if len(macdResults) > 0 {
		indicators.MACD = &macdResults[len(macdResults)-1]
	}

	// 计算RSI
	rsiResults := h.calculator.CalculateRSI(data, 14)
	if len(rsiResults) > 0 {
		indicators.RSI = &rsiResults[len(rsiResults)-1]
		log.Printf("[GetIndicators] RSI计算完成 - 股票代码: %s, 最新值: %v", stockCode, indicators.RSI)
	}

	// 计算布林带
	bollingerResults := h.calculator.CalculateBollingerBands(data, 20, 2)
	if len(bollingerResults) > 0 {
		indicators.BOLL = &bollingerResults[len(bollingerResults)-1]
		log.Printf("[GetIndicators] 布林带计算完成 - 股票代码: %s, 最新值: %v", stockCode, indicators.BOLL)
	}

	// 记录响应信息
	responseTime := time.Since(startTime)
	log.Printf("[GetIndicators] 请求处理完成 - 股票代码: %s, 响应时间: %v, 指标数量: %d",
		stockCode, responseTime, 3) // MACD, RSI, 布林带

	h.writeSuccessResponse(w, indicators)
}

// GetPredictions 获取买卖预测
func (h *StockHandler) GetPredictions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 记录请求开始
	startTime := time.Now()
	log.Printf("[GetPredictions] 开始处理请求 - 路径: %s, 方法: %s",
		r.URL.Path, r.Method)

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		errorMsg := "无效的股票代码"
		log.Printf("[GetPredictions] 参数解析失败 - 路径: %s, 错误: %s", r.URL.Path, errorMsg)
		h.writeErrorResponse(w, http.StatusBadRequest, errorMsg)
		return
	}

	stockCode := pathParts[3] // /api/v1/stocks/{code}/predictions

	// 获取足够的历史数据用于预测（使用配置的时间窗口）
	startDate := time.Now().AddDate(0, 0, -h.config.PatternPredictionDays).Format("20060102")
	endDate := time.Now().Format("20060102")

	// 记录请求参数
	log.Printf("[GetPredictions] 请求参数 - 股票代码: %s, 开始日期: %s, 结束日期: %s",
		stockCode, startDate, endDate)

	// 获取股票数据（优先使用缓存）
	var data []models.StockDaily
	var err error

	if h.dailyCacheService != nil {
		log.Printf("[GetPredictions] 尝试从缓存获取数据 - 股票代码: %s", stockCode)
		if cachedData, found := h.dailyCacheService.Get(stockCode, startDate, endDate); found {
			log.Printf("[GetPredictions] 缓存命中 - 股票代码: %s, 数据条数: %d", stockCode, len(cachedData))
			data = cachedData
		} else {
			log.Printf("[GetPredictions] 缓存未命中 - 股票代码: %s, 从API获取数据", stockCode)
			// 缓存未命中，从API获取数据
			data, err = h.dataSourceClient.GetDailyData(stockCode, startDate, endDate, "")
			if err != nil {
				// 详细记录错误信息
				log.Printf("[GetPredictions] 获取股票数据失败 - 股票代码: %s, 开始日期: %s, 结束日期: %s, 错误: %v",
					stockCode, startDate, endDate, err)
				log.Printf("[GetPredictions] 错误详情 - 类型: %T, 消息: %s", err, err.Error())

				// 记录请求上下文信息
				log.Printf("[GetPredictions] 请求上下文 - 远程地址: %s, 引用页: %s",
					r.RemoteAddr, r.Referer())

				h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
				return
			}

			log.Printf("[GetPredictions] 从API获取数据成功 - 股票代码: %s, 数据条数: %d", stockCode, len(data))
			// 将数据存入缓存
			h.dailyCacheService.Set(stockCode, startDate, endDate, data)
			log.Printf("[GetPredictions] 数据已存入缓存 - 股票代码: %s", stockCode)
		}
	} else {
		log.Printf("[GetPredictions] 缓存服务未启用 - 股票代码: %s, 直接从API获取", stockCode)
		// 如果缓存服务未启用，直接从API获取
		data, err = h.dataSourceClient.GetDailyData(stockCode, startDate, endDate, "")
		if err != nil {
			// 详细记录错误信息
			log.Printf("[GetPredictions] 获取股票数据失败 - 股票代码: %s, 开始日期: %s, 结束日期: %s, 错误: %v",
				stockCode, startDate, endDate, err)
			log.Printf("[GetPredictions] 错误详情 - 类型: %T, 消息: %s", err, err.Error())

			// 记录请求上下文信息
			log.Printf("[GetPredictions] 请求上下文 - 远程地址: %s, 引用页: %s",
				r.RemoteAddr, r.Referer())

			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
			return
		}

		log.Printf("[GetPredictions] 从API获取数据成功 - 股票代码: %s, 数据条数: %d", stockCode, len(data))
	}

	if len(data) == 0 {
		errorMsg := "未找到股票数据"
		log.Printf("[GetPredictions] 数据为空 - 股票代码: %s, 开始日期: %s, 结束日期: %s",
			stockCode, startDate, endDate)
		h.writeErrorResponse(w, http.StatusNotFound, errorMsg)
		return
	}

	log.Printf("[GetPredictions] 开始生成预测 - 股票代码: %s, 数据条数: %d", stockCode, len(data))

	// 生成预测
	prediction, err := h.predictionService.PredictTradingPoints(data)
	if err != nil {
		// 详细记录错误信息
		log.Printf("[GetPredictions] 预测失败 - 股票代码: %s, 错误: %v", stockCode, err)
		log.Printf("[GetPredictions] 错误详情 - 类型: %T, 消息: %s", err, err.Error())

		// 记录请求上下文信息
		log.Printf("[GetPredictions] 请求上下文 - 远程地址: %s, 引用页: %s",
			r.RemoteAddr, r.Referer())

		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("预测失败: %v", err))
		return
	}

	// 记录响应信息
	responseTime := time.Since(startTime)
	log.Printf("[GetPredictions] 请求处理完成 - 股票代码: %s, 响应时间: %v, 预测数量: %d",
		stockCode, responseTime, len(prediction.Predictions))

	h.writeSuccessResponse(w, prediction)
}

// GetStockBasic 获取股票基本信息
func (h *StockHandler) GetStockBasic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		h.writeErrorResponse(w, http.StatusBadRequest, "无效的股票代码")
		return
	}

	stockCode := pathParts[3] // /api/v1/stocks/{code}/basic

	// 优先从本地数据获取股票基本信息
	stockBasic, err := h.localStockService.GetStockBasic(stockCode)
	if err != nil {
		log.Printf("从本地获取股票基本信息失败: %v，尝试从API获取", err)

		// 如果本地获取失败，尝试从数据源API获取
		stockBasic, err = h.dataSourceClient.GetStockBasic(stockCode)
		if err != nil {
			log.Printf("从API获取股票基本信息也失败: %v", err)
			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票基本信息失败: %v", err))
			return
		}
	}

	h.writeSuccessResponse(w, stockBasic)
}

// GetStockList 获取本地股票列表
func (h *StockHandler) GetStockList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 获取所有本地股票
	stocks := h.localStockService.GetAllStocks()

	// 构建响应数据
	response := map[string]interface{}{
		"total":  len(stocks),
		"stocks": stocks,
	}

	h.writeSuccessResponse(w, response)
}

// SearchStocks 搜索股票
func (h *StockHandler) SearchStocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 获取搜索参数
	query := r.URL.Query()
	keyword := query.Get("q")

	if keyword == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "搜索关键词不能为空")
		return
	}

	// 限制返回结果数量，默认10条
	limit := 10
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	// 执行搜索
	stocks := h.localStockService.SearchStocks(keyword, limit)

	// 构建响应数据
	response := map[string]interface{}{
		"keyword": keyword,
		"total":   len(stocks),
		"stocks":  stocks,
	}

	h.writeSuccessResponse(w, response)
}

// RefreshLocalData 刷新本地数据
func (h *StockHandler) RefreshLocalData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 刷新本地数据
	if err := h.localStockService.RefreshData(); err != nil {
		log.Printf("刷新本地数据失败: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("刷新本地数据失败: %v", err))
		return
	}

	// 获取刷新后的股票数量
	count := h.localStockService.GetStockCount()

	response := map[string]interface{}{
		"message": "本地数据刷新成功",
		"count":   count,
	}

	h.writeSuccessResponse(w, response)
}

// GetCacheStats 获取缓存统计信息
func (h *StockHandler) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if h.dailyCacheService == nil {
		h.writeErrorResponse(w, http.StatusServiceUnavailable, "缓存服务未启用")
		return
	}

	cacheInfo := h.dailyCacheService.GetCacheInfo()
	h.writeSuccessResponse(w, cacheInfo)
}

// ClearCache 清空缓存
func (h *StockHandler) ClearCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if h.dailyCacheService == nil {
		h.writeErrorResponse(w, http.StatusServiceUnavailable, "缓存服务未启用")
		return
	}

	h.dailyCacheService.Clear()

	response := map[string]interface{}{
		"message":   "缓存已清空",
		"timestamp": time.Now(),
	}

	h.writeSuccessResponse(w, response)
}

// GetHealthStatus 获取服务器健康状态
func (h *StockHandler) GetHealthStatus(w http.ResponseWriter, r *http.Request) {
	// 使用缓存的数据源状态，避免每次健康检查都进行连接测试
	// 只有在明确需要检查连接时才进行实际测试
	dataSourceStatus := "healthy"

	// 检查是否需要强制刷新连接状态（通过查询参数控制）
	if r.URL.Query().Get("check_connection") == "true" {
		if err := h.dataSourceClient.TestConnection(); err != nil {
			dataSourceStatus = "unhealthy"
			log.Printf("健康检查: 数据源连接失败: %v", err)
		}
	}

	// 构建健康状态响应
	healthStatus := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"services": map[string]interface{}{
			"data_source": map[string]interface{}{
				"status": dataSourceStatus,
				"url":    h.dataSourceClient.GetBaseURL(),
				"note":   "使用 ?check_connection=true 参数可强制检查连接状态",
			},
			"cache": map[string]interface{}{
				"enabled": h.dailyCacheService != nil,
			},
			"favorites": map[string]interface{}{
				"status": "healthy",
			},
		},
		"uptime": time.Since(startTime).String(),
	}

	// 如果数据源不健康，设置整体状态为不健康
	if dataSourceStatus == "unhealthy" {
		healthStatus["status"] = "degraded"
		// 使用503状态码表示服务降级
		h.writeSuccessResponseWithStatus(w, healthStatus, http.StatusServiceUnavailable)
	} else {
		// 使用标准响应格式返回健康状态
		h.writeSuccessResponse(w, healthStatus)
	}
}

// GetFavorites 获取收藏股票列表
func (h *StockHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	favorites := h.favoriteService.GetFavorites()

	response := map[string]interface{}{
		"total":     len(favorites),
		"favorites": favorites,
	}

	h.writeSuccessResponse(w, response)
}

// AddFavorite 添加收藏股票
func (h *StockHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "只支持POST方法")
		return
	}

	// 解析请求体
	var request models.FavoriteStockRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "请求数据格式错误")
		return
	}

	// 添加调试日志
	log.Printf("收到添加收藏请求: %+v", request)

	// 添加收藏
	favorite, err := h.favoriteService.AddFavorite(&request)
	if err != nil {
		log.Printf("添加收藏失败: %v", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("添加收藏成功: %+v", favorite)
	h.writeSuccessResponse(w, favorite)
}

// DeleteFavorite 删除收藏股票
func (h *StockHandler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodDelete {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "只支持DELETE方法")
		return
	}

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		h.writeErrorResponse(w, http.StatusBadRequest, "无效的收藏ID")
		return
	}

	favoriteID := pathParts[3] // /api/v1/favorites/{id}

	// 删除收藏
	if err := h.favoriteService.DeleteFavorite(favoriteID); err != nil {
		log.Printf("删除收藏失败: %v", err)
		h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	response := map[string]interface{}{
		"message": "收藏删除成功",
		"id":      favoriteID,
	}

	h.writeSuccessResponse(w, response)
}

// UpdateFavorite 更新收藏股票的时间范围
func (h *StockHandler) UpdateFavorite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPut {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "只支持PUT方法")
		return
	}

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		h.writeErrorResponse(w, http.StatusBadRequest, "无效的收藏ID")
		return
	}

	favoriteID := pathParts[3] // /api/v1/favorites/{id}

	// 解析请求体
	var request models.UpdateFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "请求数据格式错误")
		return
	}

	// 更新收藏
	favorite, err := h.favoriteService.UpdateFavorite(favoriteID, &request)
	if err != nil {
		log.Printf("更新收藏失败: %v", err)
		h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeSuccessResponse(w, favorite)
}

// CheckFavorite 检查股票是否已收藏
func (h *StockHandler) CheckFavorite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 5 {
		h.writeErrorResponse(w, http.StatusBadRequest, "无效的股票代码")
		return
	}

	stockCode := pathParts[4] // /api/v1/favorites/check/{code}

	isFavorite := h.favoriteService.IsFavorite(stockCode)

	response := map[string]interface{}{
		"ts_code":     stockCode,
		"is_favorite": isFavorite,
	}

	h.writeSuccessResponse(w, response)
}

// GetFavoritesSignals 获取所有收藏股票的信号汇总
func (h *StockHandler) GetFavoritesSignals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "只支持GET方法")
		return
	}

	// 获取所有收藏股票
	favorites := h.favoriteService.GetFavorites()
	if len(favorites) == 0 {
		h.writeSuccessResponse(w, map[string]interface{}{
			"total":   0,
			"signals": []interface{}{},
		})
		return
	}

	log.Printf("获取 %d 支收藏股票的信号数据", len(favorites))

	// 获取信号服务
	signalService, ok := h.getSignalService()
	if !ok {
		h.writeErrorResponse(w, http.StatusInternalServerError, "信号服务不可用")
		return
	}

	// 批量获取信号数据
	var signals []models.FavoriteSignal
	today := time.Now().Format("20060102")

	for _, favorite := range favorites {
		// 从数据库获取已计算好的信号
		stockSignal, err := signalService.GetSignal(favorite.TSCode, today)

		// 如果没有找到信号，获取最新的股票数据用于显示基本信息
		if err != nil {
			stockData, err := h.dataSourceClient.GetDailyData(favorite.TSCode, "", "", "")
			if err != nil || len(stockData) == 0 {
				continue
			}

			// 构建简单信号（只有基本信息，没有预测）
			signal := models.FavoriteSignal{
				ID:           favorite.ID,
				TSCode:       favorite.TSCode,
				Name:         favorite.Name,
				GroupID:      favorite.GroupID,
				CurrentPrice: stockData[len(stockData)-1].Close.Decimal.String(),
				TradeDate:    stockData[len(stockData)-1].TradeDate,
				Indicators: models.SimpleIndicator{
					MA5: "N/A", MA10: "N/A", MA20: "N/A",
					CurrentPrice: stockData[len(stockData)-1].Close.Decimal.String(),
					Trend:        "UNKNOWN",
					LastUpdate:   "未计算",
				},
				Predictions: map[string]interface{}{"message": "信号尚未计算"},
				UpdatedAt:   "未计算",
			}
			signals = append(signals, signal)
			continue
		}

		// 从已计算的信号构建响应
		indicators := h.buildIndicatorsFromSignal(stockSignal)

		// 获取最新的股票数据以显示当前价格
		var currentPrice string
		stockData, err := h.dataSourceClient.GetDailyData(favorite.TSCode, "", "", "")
		if err == nil && len(stockData) > 0 {
			currentPrice = stockData[len(stockData)-1].Close.Decimal.String()
		} else {
			// 如果获取失败，使用指标中的价格或默认值
			currentPrice = indicators.CurrentPrice
		}

		signal := models.FavoriteSignal{
			ID:           favorite.ID,
			TSCode:       favorite.TSCode,
			Name:         favorite.Name,
			GroupID:      favorite.GroupID,
			CurrentPrice: currentPrice,
			TradeDate:    stockSignal.TradeDate,
			Indicators:   indicators,
			Predictions:  stockSignal.Predictions,
			UpdatedAt:    stockSignal.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		signals = append(signals, signal)
	}

	log.Printf("所有股票信号处理完成，共 %d 支", len(signals))

	response := models.FavoritesSignalsResponse{
		Total:   len(signals),
		Signals: signals,
	}

	h.writeSuccessResponse(w, response)
}

// getSignalService 获取信号服务实例
func (h *StockHandler) getSignalService() (*service.SignalService, bool) {
	// 直接从App获取信号服务
	if h.app != nil && h.app.SignalService != nil {
		return h.app.SignalService, true
	}

	// 如果没有注册到App中，则返回失败
	log.Printf("警告: 信号服务不可用")
	return nil, false
}

// buildIndicatorsFromSignal 从信号数据构建技术指标
func (h *StockHandler) buildIndicatorsFromSignal(signal *models.StockSignal) models.SimpleIndicator {
	// 如果没有技术指标数据，返回默认值
	if signal == nil || signal.TechnicalIndicators == nil || signal.TechnicalIndicators.MA == nil {
		return models.SimpleIndicator{
			MA5: "N/A", MA10: "N/A", MA20: "N/A",
			CurrentPrice: "0.00",
			Trend:        "UNKNOWN",
			LastUpdate:   "未知",
		}
	}

	// 从信号中提取技术指标
	ma := signal.TechnicalIndicators.MA

	// 判断趋势
	trend := "HOLD"
	switch signal.SignalType {
	case "BUY":
		trend = "UP"
	case "SELL":
		trend = "DOWN"
	}

	return models.SimpleIndicator{
		MA5:            ma.MA5.Decimal.StringFixed(2),
		MA10:           ma.MA10.Decimal.StringFixed(2),
		MA20:           ma.MA20.Decimal.StringFixed(2),
		CurrentPrice:   "0.00", // 这个值需要从最新数据获取，信号中可能没有
		PriceChange:    "0.00", // 同上
		PriceChangePct: "0.00", // 同上
		Trend:          trend,
		LastUpdate:     signal.UpdatedAt.Format("15:04:05"),
	}
}

// writeSuccessResponse 写入成功响应
func (h *StockHandler) writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	h.writeSuccessResponseWithStatus(w, data, http.StatusOK)
}

// writeSuccessResponseWithStatus 写入成功响应（带状态码）
func (h *StockHandler) writeSuccessResponseWithStatus(w http.ResponseWriter, data interface{}, statusCode int) {
	response := models.APIResponse{
		Success: true,
		Data:    data,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("序列化响应失败: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "序列化响应失败")
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("写入响应失败: %v", err)
	}
}

// writeErrorResponse 写入错误响应
func (h *StockHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := models.APIResponse{
		Success: false,
		Error:   message,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("序列化错误响应失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(`{"success":false,"error":"内部服务器错误"}`)); writeErr != nil {
			log.Printf("写入错误响应失败: %v", writeErr)
		}
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("写入错误响应失败: %v", err)
	}
}

// === 分组相关API ===

// GetGroups 获取分组列表
func (h *StockHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	groups := h.favoriteService.GetGroups()

	response := map[string]interface{}{
		"total":  len(groups),
		"groups": groups,
	}

	h.writeSuccessResponse(w, response)
}

// CreateGroup 创建分组
func (h *StockHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "只支持POST方法")
		return
	}

	// 解析请求体
	var request models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "请求数据格式错误")
		return
	}

	// 创建分组
	group, err := h.favoriteService.CreateGroup(&request)
	if err != nil {
		log.Printf("创建分组失败: %v", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeSuccessResponse(w, group)
}

// UpdateGroup 更新分组
func (h *StockHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPut {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "只支持PUT方法")
		return
	}

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		h.writeErrorResponse(w, http.StatusBadRequest, "无效的分组ID")
		return
	}

	groupID := pathParts[3] // /api/v1/groups/{id}

	// 解析请求体
	var request models.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "请求数据格式错误")
		return
	}

	// 更新分组
	group, err := h.favoriteService.UpdateGroup(groupID, &request)
	if err != nil {
		log.Printf("更新分组失败: %v", err)
		h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeSuccessResponse(w, group)
}

// DeleteGroup 删除分组
func (h *StockHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodDelete {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "只支持DELETE方法")
		return
	}

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		h.writeErrorResponse(w, http.StatusBadRequest, "无效的分组ID")
		return
	}

	groupID := pathParts[3] // /api/v1/groups/{id}

	// 删除分组
	if err := h.favoriteService.DeleteGroup(groupID); err != nil {
		log.Printf("删除分组失败: %v", err)
		h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	response := map[string]interface{}{
		"message": "分组删除成功",
		"id":      groupID,
	}

	h.writeSuccessResponse(w, response)
}

// UpdateFavoritesOrder 更新收藏排序
func (h *StockHandler) UpdateFavoritesOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPut {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "只支持PUT方法")
		return
	}

	// 解析请求体
	var request models.UpdateFavoritesOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "请求数据格式错误")
		return
	}

	// 更新排序
	if err := h.favoriteService.UpdateFavoritesOrder(&request); err != nil {
		log.Printf("更新收藏排序失败: %v", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"message": "排序更新成功",
	}

	h.writeSuccessResponse(w, response)
}
