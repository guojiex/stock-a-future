package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/service"
)

// ParameterOptimizerHandler 参数优化处理器
type ParameterOptimizerHandler struct {
	optimizer *service.ParameterOptimizer
	logger    logger.Logger
}

// NewParameterOptimizerHandler 创建参数优化处理器
func NewParameterOptimizerHandler(optimizer *service.ParameterOptimizer, log logger.Logger) *ParameterOptimizerHandler {
	return &ParameterOptimizerHandler{
		optimizer: optimizer,
		logger:    log,
	}
}

// RegisterRoutes 注册路由
func (h *ParameterOptimizerHandler) RegisterRoutes(mux *http.ServeMux) {
	// 参数优化路由
	mux.HandleFunc("POST /api/v1/strategies/{id}/optimize", h.handleCORS(h.startOptimization))
	mux.HandleFunc("GET /api/v1/optimizations/{id}/progress", h.handleCORS(h.getOptimizationProgress))
	mux.HandleFunc("GET /api/v1/optimizations/{id}/results", h.handleCORS(h.getOptimizationResults))
	mux.HandleFunc("POST /api/v1/optimizations/{id}/cancel", h.handleCORS(h.cancelOptimization))
}

// startOptimization 启动参数优化
func (h *ParameterOptimizerHandler) startOptimization(w http.ResponseWriter, r *http.Request) {
	strategyID := r.PathValue("id")

	var config service.OptimizationConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		h.logger.Error("解析优化配置失败", logger.ErrorField(err))
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "无效的请求数据",
			Error:   err.Error(),
		})
		return
	}

	// 设置策略ID
	config.StrategyID = strategyID

	// 设置默认值
	if config.Algorithm == "" {
		config.Algorithm = "grid_search"
	}
	if config.OptimizationTarget == "" {
		config.OptimizationTarget = "sharpe_ratio"
	}
	if config.InitialCash == 0 {
		config.InitialCash = 1000000
	}
	if config.Commission == 0 {
		config.Commission = 0.0003
	}
	if config.MaxCombinations == 0 {
		config.MaxCombinations = 100
	}

	// 验证配置
	if len(config.ParameterRanges) == 0 {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "参数范围不能为空",
		})
		return
	}

	if len(config.Symbols) == 0 {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "股票列表不能为空",
		})
		return
	}

	// 启动优化（使用background context，让优化任务独立于HTTP请求生命周期）
	optimizationID, err := h.optimizer.StartOptimization(context.Background(), &config)
	if err != nil {
		h.logger.Error("启动优化失败", logger.ErrorField(err))
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "启动优化失败",
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info("参数优化已启动",
		logger.String("optimization_id", optimizationID),
		logger.String("strategy_id", strategyID),
	)

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "参数优化已启动",
		Data: map[string]interface{}{
			"optimization_id": optimizationID,
			"strategy_id":     strategyID,
		},
	})
}

// getOptimizationProgress 获取优化进度
func (h *ParameterOptimizerHandler) getOptimizationProgress(w http.ResponseWriter, r *http.Request) {
	optimizationID := r.PathValue("id")

	task, err := h.optimizer.GetOptimizationProgress(optimizationID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "优化任务不存在",
			Error:   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "获取优化进度成功",
		Data: map[string]interface{}{
			"optimization_id": task.ID,
			"strategy_id":     task.StrategyID,
			"status":          task.Status,
			"progress":        task.Progress,
			"current_combo":   task.CurrentCombo,
			"total_combos":    task.TotalCombos,
			"current_params":  task.CurrentParams,
			"best_params":     task.BestParams,
			"best_score":      task.BestScore,
			"start_time":      task.StartTime,
			"estimated_end":   task.EstimatedEndTime,
		},
	})
}

// getOptimizationResults 获取优化结果
func (h *ParameterOptimizerHandler) getOptimizationResults(w http.ResponseWriter, r *http.Request) {
	optimizationID := r.PathValue("id")

	result, err := h.optimizer.GetOptimizationResult(optimizationID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "获取优化结果失败",
			Error:   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "获取优化结果成功",
		Data:    result,
	})
}

// cancelOptimization 取消优化任务
func (h *ParameterOptimizerHandler) cancelOptimization(w http.ResponseWriter, r *http.Request) {
	optimizationID := r.PathValue("id")

	err := h.optimizer.CancelOptimization(optimizationID)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "取消优化失败",
			Error:   err.Error(),
		})
		return
	}

	h.logger.Info("优化任务已取消", logger.String("optimization_id", optimizationID))

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "优化任务已取消",
	})
}

// handleCORS CORS中间件
func (h *ParameterOptimizerHandler) handleCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// respondJSON 返回JSON响应
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// APIResponse API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
