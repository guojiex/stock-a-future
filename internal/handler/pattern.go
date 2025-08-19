package handler

import (
	"encoding/json"
	"net/http"
	"stock-a-future/internal/models"
	"stock-a-future/internal/service"
	"strconv"
)

// PatternHandler 图形识别处理器
type PatternHandler struct {
	patternService *service.PatternService
}

// NewPatternHandler 创建新的图形识别处理器
func NewPatternHandler(patternService *service.PatternService) *PatternHandler {
	return &PatternHandler{
		patternService: patternService,
	}
}

// RecognizePatterns 识别图形模式
func (h *PatternHandler) RecognizePatterns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	tsCode := r.URL.Query().Get("ts_code")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if tsCode == "" {
		http.Error(w, "ts_code is required", http.StatusBadRequest)
		return
	}

	// 设置默认日期范围
	if startDate == "" {
		startDate = "20240101"
	}
	if endDate == "" {
		endDate = "20241231"
	}

	// 识别图形模式
	patterns, err := h.patternService.RecognizePatterns(tsCode, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	response := models.APIResponse{
		Success: true,
		Data:    patterns,
	}

	// 响应头已在中间件中设置
	json.NewEncoder(w).Encode(response)
}

// SearchPatterns 搜索指定图形模式
func (h *PatternHandler) SearchPatterns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var request models.PatternSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.TSCode == "" {
		http.Error(w, "ts_code is required", http.StatusBadRequest)
		return
	}

	// 设置默认日期范围
	if request.StartDate == "" {
		request.StartDate = "20240101"
	}
	if request.EndDate == "" {
		request.EndDate = "20241231"
	}

	// 搜索图形模式
	result, err := h.patternService.SearchPatterns(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	response := models.APIResponse{
		Success: true,
		Data:    result,
	}

	// 响应头已在中间件中设置
	json.NewEncoder(w).Encode(response)
}

// GetPatternSummary 获取图形模式摘要
func (h *PatternHandler) GetPatternSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	tsCode := r.URL.Query().Get("ts_code")
	daysStr := r.URL.Query().Get("days")

	if tsCode == "" {
		http.Error(w, "ts_code is required", http.StatusBadRequest)
		return
	}

	// 设置默认天数
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	// 获取图形模式摘要
	summary, err := h.patternService.GetPatternSummary(tsCode, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	response := models.APIResponse{
		Success: true,
		Data:    summary,
	}

	// 响应头已在中间件中设置
	json.NewEncoder(w).Encode(response)
}

// GetRecentSignals 获取最近的图形信号
func (h *PatternHandler) GetRecentSignals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	tsCode := r.URL.Query().Get("ts_code")
	limitStr := r.URL.Query().Get("limit")

	if tsCode == "" {
		http.Error(w, "ts_code is required", http.StatusBadRequest)
		return
	}

	// 设置默认限制
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// 获取最近的图形信号
	signals, err := h.patternService.GetRecentSignals(tsCode, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	response := models.APIResponse{
		Success: true,
		Data:    signals,
	}

	// 响应头已在中间件中设置
	json.NewEncoder(w).Encode(response)
}

// GetAvailablePatterns 获取可用的图形模式列表
func (h *PatternHandler) GetAvailablePatterns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 定义可用的图形模式
	patterns := map[string]interface{}{
		"candlestick": map[string]string{
			"双响炮":  "连续两根大阳线，成交量放大，强势上涨信号",
			"红三兵":  "连续三根上涨K线，稳步上涨信号",
			"乌云盖顶": "大阳线后跟大阴线，可能见顶信号",
			"锤子线":  "下影线很长，可能见底信号",
			"启明星":  "下跌趋势中的反转信号",
			"黄昏星":  "上涨趋势中的反转信号",
		},
		"volume_price": map[string]string{
			"量价齐升": "价格和成交量同时上涨，强势信号",
			"量价背离": "价格与成交量变化不一致，可能见顶或见底",
			"放量突破": "价格突破重要阻力位，成交量放大",
		},
	}

	// 返回结果
	response := models.APIResponse{
		Success: true,
		Data:    patterns,
	}

	// 响应头已在中间件中设置
	json.NewEncoder(w).Encode(response)
}

// GetPatternStatistics 获取图形模式统计信息
func (h *PatternHandler) GetPatternStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	tsCode := r.URL.Query().Get("ts_code")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if tsCode == "" {
		http.Error(w, "ts_code is required", http.StatusBadRequest)
		return
	}

	// 设置默认日期范围
	if startDate == "" {
		startDate = "20240101"
	}
	if endDate == "" {
		endDate = "20241231"
	}

	// 识别图形模式
	patterns, err := h.patternService.RecognizePatterns(tsCode, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 统计各种图形模式
	statistics := map[string]interface{}{
		"total_patterns": 0,
		"pattern_types":  make(map[string]int),
		"signal_distribution": map[string]int{
			"BUY":  0,
			"SELL": 0,
			"HOLD": 0,
		},
		"strength_distribution": map[string]int{
			"STRONG": 0,
			"MEDIUM": 0,
			"WEAK":   0,
		},
		"confidence_ranges": map[string]int{
			"80-100": 0,
			"60-79":  0,
			"40-59":  0,
			"0-39":   0,
		},
	}

	// 统计图形模式
	for _, pattern := range patterns {
		statistics["total_patterns"] = statistics["total_patterns"].(int) + 1

		// 统计蜡烛图模式
		for _, candlestick := range pattern.Candlestick {
			patternTypes := statistics["pattern_types"].(map[string]int)
			patternTypes[candlestick.Pattern]++

			signals := statistics["signal_distribution"].(map[string]int)
			signals[candlestick.Signal]++

			strengths := statistics["strength_distribution"].(map[string]int)
			strengths[candlestick.Strength]++

			// 统计置信度范围
			confidence, _ := candlestick.Confidence.Decimal.Float64()
			confidenceRanges := statistics["confidence_ranges"].(map[string]int)
			if confidence >= 80 {
				confidenceRanges["80-100"]++
			} else if confidence >= 60 {
				confidenceRanges["60-79"]++
			} else if confidence >= 40 {
				confidenceRanges["40-59"]++
			} else {
				confidenceRanges["0-39"]++
			}
		}

		// 统计量价图形
		for _, volumePrice := range pattern.VolumePrice {
			patternTypes := statistics["pattern_types"].(map[string]int)
			patternTypes[volumePrice.Pattern]++

			signals := statistics["signal_distribution"].(map[string]int)
			signals[volumePrice.Signal]++

			strengths := statistics["strength_distribution"].(map[string]int)
			strengths[volumePrice.Strength]++

			// 统计置信度范围
			confidence, _ := volumePrice.Confidence.Decimal.Float64()
			confidenceRanges := statistics["confidence_ranges"].(map[string]int)
			if confidence >= 80 {
				confidenceRanges["80-100"]++
			} else if confidence >= 60 {
				confidenceRanges["60-79"]++
			} else if confidence >= 40 {
				confidenceRanges["40-59"]++
			} else {
				confidenceRanges["0-39"]++
			}
		}
	}

	// 返回结果
	response := models.APIResponse{
		Success: true,
		Data:    statistics,
	}

	// 响应头已在中间件中设置
	json.NewEncoder(w).Encode(response)
}
