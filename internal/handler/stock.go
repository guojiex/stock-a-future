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
	"strings"
	"time"
)

// StockHandler 股票数据处理器
type StockHandler struct {
	tushareClient     *client.TushareClient
	calculator        *indicators.Calculator
	predictionService *service.PredictionService
}

// NewStockHandler 创建股票处理器
func NewStockHandler(tushareClient *client.TushareClient) *StockHandler {
	return &StockHandler{
		tushareClient:     tushareClient,
		calculator:        indicators.NewCalculator(),
		predictionService: service.NewPredictionService(),
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

	// 调用Tushare API
	data, err := h.tushareClient.GetDailyData(stockCode, startDate, endDate)
	if err != nil {
		log.Printf("获取股票数据失败: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
		return
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

	// 获取股票数据
	data, err := h.tushareClient.GetDailyData(stockCode, startDate, endDate)
	if err != nil {
		log.Printf("获取股票数据失败: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
		return
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
			MA5:  ma5[len(ma5)-1],
			MA10: ma10[len(ma10)-1],
			MA20: ma20[len(ma20)-1],
		}
		if len(ma60) > 0 {
			indicators.MA.MA60 = ma60[len(ma60)-1]
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

	// 获取股票数据
	data, err := h.tushareClient.GetDailyData(stockCode, startDate, endDate)
	if err != nil {
		log.Printf("获取股票数据失败: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取股票数据失败: %v", err))
		return
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
	w.Write(jsonData)
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
		w.Write([]byte(`{"success":false,"error":"内部服务器错误"}`))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonData)
}
