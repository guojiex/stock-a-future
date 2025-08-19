package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"stock-a-future/internal/models"
	"stock-a-future/internal/service"
	"strconv"
	"time"
)

// SignalHandler 信号处理器
type SignalHandler struct {
	signalService *service.SignalService
}

// NewSignalHandler 创建信号处理器
func NewSignalHandler(signalService *service.SignalService) *SignalHandler {
	return &SignalHandler{
		signalService: signalService,
	}
}

// CalculateSignal 计算单个股票信号
func (h *SignalHandler) CalculateSignal(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SignalHandler] 收到计算信号请求: %s %s", r.Method, r.URL.Path)

	var request models.SignalCalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("解析请求失败: %v", err)
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	// 验证请求参数
	if request.TSCode == "" {
		http.Error(w, "股票代码不能为空", http.StatusBadRequest)
		return
	}
	if request.Name == "" {
		http.Error(w, "股票名称不能为空", http.StatusBadRequest)
		return
	}

	// 计算信号
	err := h.signalService.CalculateStockSignal(request.TSCode, request.Name, request.Force)
	if err != nil {
		log.Printf("计算股票信号失败: %v", err)
		response := models.SignalCalculationResponse{
			Success: false,
			Message: "计算失败",
			Error:   err.Error(),
		}
		// 响应头已在中间件中设置
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取计算结果
	today := time.Now().Format("20060102")
	signal, err := h.signalService.GetSignal(request.TSCode, today)
	if err != nil {
		log.Printf("获取信号结果失败: %v", err)
		response := models.SignalCalculationResponse{
			Success: false,
			Message: "获取结果失败",
			Error:   err.Error(),
		}
		// 响应头已在中间件中设置
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.SignalCalculationResponse{
		Success: true,
		Message: "计算成功",
		Signal:  signal,
	}

	// 响应头已在中间件中设置
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码响应失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}

	log.Printf("[SignalHandler] 计算信号成功: %s", request.TSCode)
}

// BatchCalculateSignals 批量计算信号
func (h *SignalHandler) BatchCalculateSignals(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SignalHandler] 收到批量计算信号请求: %s %s", r.Method, r.URL.Path)

	var request models.BatchSignalRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("解析请求失败: %v", err)
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	if len(request.TSCodes) == 0 {
		http.Error(w, "股票代码列表不能为空", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	total := len(request.TSCodes)
	success := 0
	failed := 0
	var results []models.SignalCalculationResponse

	// 批量计算信号
	for _, tsCode := range request.TSCodes {
		err := h.signalService.CalculateStockSignal(tsCode, "", request.Force)

		result := models.SignalCalculationResponse{
			Success: err == nil,
		}

		if err != nil {
			result.Message = "计算失败"
			result.Error = err.Error()
			failed++
		} else {
			result.Message = "计算成功"
			success++

			// 获取计算结果
			today := time.Now().Format("20060102")
			if signal, err := h.signalService.GetSignal(tsCode, today); err == nil {
				result.Signal = signal
			}
		}

		results = append(results, result)
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	response := models.BatchSignalResponse{
		Total:     total,
		Success:   success,
		Failed:    failed,
		Results:   results,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration.String(),
	}

	// 响应头已在中间件中设置
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码响应失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}

	log.Printf("[SignalHandler] 批量计算信号完成: 总数=%d, 成功=%d, 失败=%d, 耗时=%v",
		total, success, failed, duration)
}

// GetSignal 获取股票信号
func (h *SignalHandler) GetSignal(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SignalHandler] 收到获取信号请求: %s %s", r.Method, r.URL.Path)

	// 从路径中获取股票代码
	tsCode := r.PathValue("code")
	if tsCode == "" {
		http.Error(w, "股票代码不能为空", http.StatusBadRequest)
		return
	}

	// 获取信号日期参数，默认为今天
	signalDate := r.URL.Query().Get("signal_date")
	if signalDate == "" {
		signalDate = time.Now().Format("20060102")
	}

	signal, err := h.signalService.GetSignal(tsCode, signalDate)
	if err != nil {
		log.Printf("获取信号失败: %v", err)
		http.Error(w, "获取信号失败: "+err.Error(), http.StatusNotFound)
		return
	}

	response := models.APIResponse{
		Success: true,
		Data:    signal,
	}

	// 响应头已在中间件中设置
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码响应失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}

	log.Printf("[SignalHandler] 获取信号成功: %s", tsCode)
}

// GetLatestSignals 获取最新信号列表
func (h *SignalHandler) GetLatestSignals(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SignalHandler] 收到获取最新信号请求: %s %s", r.Method, r.URL.Path)

	// 获取限制数量参数，默认为20
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	signals, err := h.signalService.GetLatestSignals(limit)
	if err != nil {
		log.Printf("获取最新信号失败: %v", err)
		http.Error(w, "获取最新信号失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"total":   len(signals),
			"signals": signals,
		},
	}

	// 响应头已在中间件中设置
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码响应失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}

	log.Printf("[SignalHandler] 获取最新信号成功: 返回 %d 条记录", len(signals))
}

// GetCalculationStatus 获取信号计算状态
func (h *SignalHandler) GetCalculationStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SignalHandler] 收到获取信号计算状态请求: %s %s", r.Method, r.URL.Path)

	// 获取计算状态
	status := h.signalService.GetCalculationStatus()

	response := models.APIResponse{
		Success: true,
		Data:    status,
	}

	// 响应头已在中间件中设置
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("编码响应失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}

	log.Printf("[SignalHandler] 获取信号计算状态成功")
}
