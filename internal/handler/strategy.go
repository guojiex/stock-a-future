package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"
	"stock-a-future/internal/service"

	"github.com/google/uuid"
)

// StrategyHandler 策略处理器
type StrategyHandler struct {
	strategyService *service.StrategyService
	logger          logger.Logger
}

// NewStrategyHandler 创建策略处理器
func NewStrategyHandler(strategyService *service.StrategyService, log logger.Logger) *StrategyHandler {
	return &StrategyHandler{
		strategyService: strategyService,
		logger:          log,
	}
}

// RegisterRoutes 注册路由
func (h *StrategyHandler) RegisterRoutes(mux *http.ServeMux) {
	// 策略管理路由
	mux.HandleFunc("GET /api/v1/strategies", h.handleCORS(h.getStrategiesList))
	mux.HandleFunc("POST /api/v1/strategies", h.handleCORS(h.createStrategy))
	mux.HandleFunc("GET /api/v1/strategies/{id}", h.handleCORS(h.getStrategy))
	mux.HandleFunc("PUT /api/v1/strategies/{id}", h.handleCORS(h.updateStrategy))
	mux.HandleFunc("DELETE /api/v1/strategies/{id}", h.handleCORS(h.deleteStrategy))
	mux.HandleFunc("GET /api/v1/strategies/{id}/performance", h.handleCORS(h.getStrategyPerformance))

	// 策略操作路由
	mux.HandleFunc("POST /api/v1/strategies/{id}/activate", h.handleCORS(h.activateStrategy))
	mux.HandleFunc("POST /api/v1/strategies/{id}/deactivate", h.handleCORS(h.deactivateStrategy))
	mux.HandleFunc("POST /api/v1/strategies/{id}/test", h.handleCORS(h.testStrategy))
}

// getStrategiesList 获取策略列表
func (h *StrategyHandler) getStrategiesList(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("获取策略列表请求")

	// 解析查询参数
	query := r.URL.Query()
	req := &models.StrategyListRequest{
		Page: 1,
		Size: 20,
	}

	if page := query.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Page = p
		}
	}

	if size := query.Get("size"); size != "" {
		if s, err := strconv.Atoi(size); err == nil && s > 0 && s <= 100 {
			req.Size = s
		}
	}

	if status := query.Get("status"); status != "" {
		req.Status = models.StrategyStatus(status)
	}

	if strategyType := query.Get("type"); strategyType != "" {
		req.Type = models.StrategyType(strategyType)
	}

	if keyword := query.Get("keyword"); keyword != "" {
		req.Keyword = strings.TrimSpace(keyword)
	}

	// 调用服务层
	strategies, total, err := h.strategyService.GetStrategiesList(r.Context(), req)
	if err != nil {
		h.logger.Error("获取策略列表失败: " + err.Error())
		h.writeErrorResponse(w, "获取策略列表失败", http.StatusInternalServerError)
		return
	}

	// 详细日志：记录返回的策略顺序
	strategyIDs := make([]string, len(strategies))
	for i, s := range strategies {
		strategyIDs[i] = s.ID
	}
	h.logger.Info("返回策略列表",
		logger.Int("total", total),
		logger.Int("count", len(strategies)),
		logger.Any("strategy_ids", strategyIDs),
	)

	// 打印每个策略的详细信息
	for i, s := range strategies {
		h.logger.Info("策略详情",
			logger.Int("position", i),
			logger.String("id", s.ID),
			logger.String("name", s.Name),
			logger.String("status", string(s.Status)),
			logger.String("type", string(s.Type)),
			logger.Time("created_at", s.CreatedAt),
		)
	}

	// 返回结果
	response := &models.StrategyListResponse{
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
		Items: strategies,
	}

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"data":    response,
		"message": "获取策略列表成功",
	})
}

// createStrategy 创建策略
func (h *StrategyHandler) createStrategy(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("创建策略请求")

	var req models.CreateStrategyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("解析创建策略请求失败: %v", err)
		h.writeErrorResponse(w, "请求格式错误", http.StatusBadRequest)
		return
	}

	// 验证请求参数
	if err := h.validateCreateStrategyRequest(&req); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 创建策略对象
	strategy := &models.Strategy{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Status:      models.StrategyStatusInactive,
		Parameters:  req.Parameters,
		Code:        req.Code,
		CreatedBy:   "user", // TODO: 从认证信息获取
	}

	// 调用服务层
	if err := h.strategyService.CreateStrategy(r.Context(), strategy); err != nil {
		h.logger.Errorf("创建策略失败: %v, strategy_name: %s", err, req.Name)
		h.writeErrorResponse(w, "创建策略失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("策略创建成功", logger.String("strategy_id", strategy.ID), logger.String("strategy_name", strategy.Name))

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"data":    strategy,
		"message": "策略创建成功",
	})
}

// getStrategy 获取策略详情
func (h *StrategyHandler) getStrategy(w http.ResponseWriter, r *http.Request) {
	strategyID := r.PathValue("id")
	if strategyID == "" {
		h.writeErrorResponse(w, "策略ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("获取策略详情请求", logger.String("strategy_id", strategyID))

	strategy, err := h.strategyService.GetStrategy(r.Context(), strategyID)
	if err != nil {
		if err == service.ErrStrategyNotFound {
			h.writeErrorResponse(w, "策略不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("获取策略详情失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "获取策略详情失败", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"data":    strategy,
		"message": "获取策略详情成功",
	})
}

// updateStrategy 更新策略
func (h *StrategyHandler) updateStrategy(w http.ResponseWriter, r *http.Request) {
	strategyID := r.PathValue("id")
	if strategyID == "" {
		h.writeErrorResponse(w, "策略ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("更新策略请求", logger.String("strategy_id", strategyID))

	var req models.UpdateStrategyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("解析更新策略请求失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "请求格式错误", http.StatusBadRequest)
		return
	}

	// 调用服务层
	if err := h.strategyService.UpdateStrategy(r.Context(), strategyID, &req); err != nil {
		if err == service.ErrStrategyNotFound {
			h.writeErrorResponse(w, "策略不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("更新策略失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "更新策略失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("策略更新成功", logger.String("strategy_id", strategyID))

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "策略更新成功",
	})
}

// deleteStrategy 删除策略
func (h *StrategyHandler) deleteStrategy(w http.ResponseWriter, r *http.Request) {
	strategyID := r.PathValue("id")
	if strategyID == "" {
		h.writeErrorResponse(w, "策略ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("删除策略请求", logger.String("strategy_id", strategyID))

	// 调用服务层
	if err := h.strategyService.DeleteStrategy(r.Context(), strategyID); err != nil {
		if err == service.ErrStrategyNotFound {
			h.writeErrorResponse(w, "策略不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("删除策略失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "删除策略失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("策略删除成功", logger.String("strategy_id", strategyID))

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "策略删除成功",
	})
}

// getStrategyPerformance 获取策略表现
func (h *StrategyHandler) getStrategyPerformance(w http.ResponseWriter, r *http.Request) {
	strategyID := r.PathValue("id")
	if strategyID == "" {
		h.writeErrorResponse(w, "策略ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("获取策略表现请求", logger.String("strategy_id", strategyID))

	performance, err := h.strategyService.GetStrategyPerformance(r.Context(), strategyID)
	if err != nil {
		if err == service.ErrStrategyNotFound {
			h.writeErrorResponse(w, "策略不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("获取策略表现失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "获取策略表现失败", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"data":    performance,
		"message": "获取策略表现成功",
	})
}

// activateStrategy 激活策略
func (h *StrategyHandler) activateStrategy(w http.ResponseWriter, r *http.Request) {
	strategyID := r.PathValue("id")
	if strategyID == "" {
		h.writeErrorResponse(w, "策略ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("激活策略请求", logger.String("strategy_id", strategyID))

	if err := h.strategyService.UpdateStrategyStatus(r.Context(), strategyID, models.StrategyStatusActive); err != nil {
		if err == service.ErrStrategyNotFound {
			h.writeErrorResponse(w, "策略不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("激活策略失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "激活策略失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("策略激活成功", logger.String("strategy_id", strategyID))

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "策略激活成功",
	})
}

// deactivateStrategy 停用策略
func (h *StrategyHandler) deactivateStrategy(w http.ResponseWriter, r *http.Request) {
	strategyID := r.PathValue("id")
	if strategyID == "" {
		h.writeErrorResponse(w, "策略ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("停用策略请求", logger.String("strategy_id", strategyID))

	if err := h.strategyService.UpdateStrategyStatus(r.Context(), strategyID, models.StrategyStatusInactive); err != nil {
		if err == service.ErrStrategyNotFound {
			h.writeErrorResponse(w, "策略不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("停用策略失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "停用策略失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("策略停用成功", logger.String("strategy_id", strategyID))

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "策略停用成功",
	})
}

// testStrategy 测试策略
func (h *StrategyHandler) testStrategy(w http.ResponseWriter, r *http.Request) {
	strategyID := r.PathValue("id")
	if strategyID == "" {
		h.writeErrorResponse(w, "策略ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("测试策略请求", logger.String("strategy_id", strategyID))

	if err := h.strategyService.UpdateStrategyStatus(r.Context(), strategyID, models.StrategyStatusTesting); err != nil {
		if err == service.ErrStrategyNotFound {
			h.writeErrorResponse(w, "策略不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("设置策略测试状态失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "设置策略测试状态失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("策略测试状态设置成功", logger.String("strategy_id", strategyID))

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "策略测试状态设置成功",
	})
}

// validateCreateStrategyRequest 验证创建策略请求
func (h *StrategyHandler) validateCreateStrategyRequest(req *models.CreateStrategyRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("策略名称不能为空")
	}

	if len(req.Name) > 100 {
		return fmt.Errorf("策略名称长度不能超过100个字符")
	}

	if len(req.Description) > 1000 {
		return fmt.Errorf("策略描述长度不能超过1000个字符")
	}

	// 验证策略类型
	validTypes := []models.StrategyType{
		models.StrategyTypeTechnical,
		models.StrategyTypeFundamental,
		models.StrategyTypeML,
		models.StrategyTypeComposite,
	}

	isValidType := false
	for _, validType := range validTypes {
		if req.Type == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("无效的策略类型: %s", req.Type)
	}

	if strings.TrimSpace(req.Code) == "" {
		return fmt.Errorf("策略代码不能为空")
	}

	return nil
}

// handleCORS 处理CORS
func (h *StrategyHandler) handleCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// writeJSONResponse 写入JSON响应
func (h *StrategyHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse 写入错误响应
func (h *StrategyHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
