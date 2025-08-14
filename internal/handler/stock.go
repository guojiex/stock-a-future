package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stock-a-future/internal/client"
	"stock-a-future/internal/indicators"
	"stock-a-future/internal/models"
	"stock-a-future/internal/service"
	"strconv"
	"strings"
	"time"
)

// StockHandler 股票数据处理器
type StockHandler struct {
	tushareClient     *client.TushareClient
	calculator        *indicators.Calculator
	predictionService *service.PredictionService
	localStockService *service.LocalStockService
	dailyCacheService *service.DailyCacheService
	favoriteService   *service.FavoriteService
}

// NewStockHandler 创建股票处理器
func NewStockHandler(tushareClient *client.TushareClient, cacheService *service.DailyCacheService, favoriteService *service.FavoriteService) *StockHandler {
	return &StockHandler{
		tushareClient:     tushareClient,
		calculator:        indicators.NewCalculator(),
		predictionService: service.NewPredictionService(),
		localStockService: service.NewLocalStockService("data"),
		dailyCacheService: cacheService,
		favoriteService:   favoriteService,
	}
}

// GetDailyData 获取股票日线数据
func (h *StockHandler) GetDailyData(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		h.writeErrorResponse(w, http.StatusBadRequest, "无效的股票代码")
		return
	}

	stockCode := pathParts[3] // /api/v1/stocks/{code}/daily

	// 解析查询参数
	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	// 默认获取最近30天数据
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("20060102")
	}
	if endDate == "" {
		endDate = time.Now().Format("20060102")
	}

	// 尝试从缓存获取数据
	var data []models.StockDaily
	var err error

	if h.dailyCacheService != nil {
		if cachedData, found := h.dailyCacheService.Get(stockCode, startDate, endDate); found {
			data = cachedData
		} else {
			// 缓存未命中，从API获取数据
			data, err = h.tushareClient.GetDailyData(stockCode, startDate, endDate)
			if err != nil {
				log.Printf("获取股票数据失败: %v", err)
				h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
				return
			}

			// 将数据存入缓存
			h.dailyCacheService.Set(stockCode, startDate, endDate, data)
		}
	} else {
		// 如果缓存服务未启用，直接从API获取
		data, err = h.tushareClient.GetDailyData(stockCode, startDate, endDate)
		if err != nil {
			log.Printf("获取股票数据失败: %v", err)
			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
			return
		}
	}

	// 返回成功响应
	h.writeSuccessResponse(w, data)
}

// GetIndicators 获取技术指标
func (h *StockHandler) GetIndicators(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		h.writeErrorResponse(w, http.StatusBadRequest, "无效的股票代码")
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

	// 获取股票数据（优先使用缓存）
	var data []models.StockDaily
	var err error

	if h.dailyCacheService != nil {
		if cachedData, found := h.dailyCacheService.Get(stockCode, startDate, endDate); found {
			data = cachedData
		} else {
			// 缓存未命中，从API获取数据
			data, err = h.tushareClient.GetDailyData(stockCode, startDate, endDate)
			if err != nil {
				log.Printf("获取股票数据失败: %v", err)
				h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
				return
			}

			// 将数据存入缓存
			h.dailyCacheService.Set(stockCode, startDate, endDate, data)
		}
	} else {
		// 如果缓存服务未启用，直接从API获取
		data, err = h.tushareClient.GetDailyData(stockCode, startDate, endDate)
		if err != nil {
			log.Printf("获取股票数据失败: %v", err)
			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
			return
		}
	}

	if len(data) == 0 {
		h.writeErrorResponse(w, http.StatusNotFound, "未找到股票数据")
		return
	}

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
	}

	// 计算布林带
	bollResults := h.calculator.CalculateBollingerBands(data, 20, 2.0)
	if len(bollResults) > 0 {
		indicators.BOLL = &bollResults[len(bollResults)-1]
	}

	// 计算移动平均线
	ma5 := h.calculator.CalculateMA(data, 5)
	ma10 := h.calculator.CalculateMA(data, 10)
	ma20 := h.calculator.CalculateMA(data, 20)
	ma60 := h.calculator.CalculateMA(data, 60)

	if len(ma5) > 0 && len(ma10) > 0 && len(ma20) > 0 {
		indicators.MA = &models.MovingAverageIndicator{
			MA5:  models.NewJSONDecimal(ma5[len(ma5)-1]),
			MA10: models.NewJSONDecimal(ma10[len(ma10)-1]),
			MA20: models.NewJSONDecimal(ma20[len(ma20)-1]),
		}
		if len(ma60) > 0 {
			indicators.MA.MA60 = models.NewJSONDecimal(ma60[len(ma60)-1])
		}
	}

	// 计算KDJ
	kdjResults := h.calculator.CalculateKDJ(data, 9)
	if len(kdjResults) > 0 {
		indicators.KDJ = &kdjResults[len(kdjResults)-1]
	}

	h.writeSuccessResponse(w, indicators)
}

// GetPredictions 获取买卖点预测
func (h *StockHandler) GetPredictions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析路径参数
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		h.writeErrorResponse(w, http.StatusBadRequest, "无效的股票代码")
		return
	}

	stockCode := pathParts[3] // /api/v1/stocks/{code}/predictions

	// 获取足够的历史数据用于预测（120天）
	startDate := time.Now().AddDate(0, 0, -120).Format("20060102")
	endDate := time.Now().Format("20060102")

	// 获取股票数据（优先使用缓存）
	var data []models.StockDaily
	var err error

	if h.dailyCacheService != nil {
		if cachedData, found := h.dailyCacheService.Get(stockCode, startDate, endDate); found {
			data = cachedData
		} else {
			// 缓存未命中，从API获取数据
			data, err = h.tushareClient.GetDailyData(stockCode, startDate, endDate)
			if err != nil {
				log.Printf("获取股票数据失败: %v", err)
				h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
				return
			}

			// 将数据存入缓存
			h.dailyCacheService.Set(stockCode, startDate, endDate, data)
		}
	} else {
		// 如果缓存服务未启用，直接从API获取
		data, err = h.tushareClient.GetDailyData(stockCode, startDate, endDate)
		if err != nil {
			log.Printf("获取股票数据失败: %v", err)
			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
			return
		}
	}

	if len(data) == 0 {
		h.writeErrorResponse(w, http.StatusNotFound, "未找到股票数据")
		return
	}

	// 生成预测
	prediction, err := h.predictionService.PredictTradingPoints(data)
	if err != nil {
		log.Printf("预测失败: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("预测失败: %v", err))
		return
	}

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

		// 如果本地获取失败，尝试从Tushare API获取
		stockBasic, err = h.tushareClient.GetStockBasic(stockCode)
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

// GetHealthStatus 健康检查
func (h *StockHandler) GetHealthStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	status := &models.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Services:  make(map[string]string),
	}

	// 测试Tushare连接
	if err := h.tushareClient.TestConnection(); err != nil {
		status.Status = "unhealthy"
		status.Services["tushare"] = "error: " + err.Error()
	} else {
		status.Services["tushare"] = "healthy"
	}

	// 根据整体状态设置HTTP状态码
	if status.Status == "healthy" {
		h.writeSuccessResponse(w, status)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		h.writeSuccessResponse(w, status)
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

	// 添加收藏
	favorite, err := h.favoriteService.AddFavorite(&request)
	if err != nil {
		log.Printf("添加收藏失败: %v", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

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

// writeSuccessResponse 写入成功响应
func (h *StockHandler) writeSuccessResponse(w http.ResponseWriter, data interface{}) {
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

	w.WriteHeader(http.StatusOK)
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
