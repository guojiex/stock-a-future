package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"stock-a-future/internal/logger"
	"stock-a-future/internal/models"
	"stock-a-future/internal/service"

	"github.com/google/uuid"
)

// BacktestHandler 回测处理器
type BacktestHandler struct {
	backtestService *service.BacktestService
	strategyService *service.StrategyService
	logger          logger.Logger
}

// NewBacktestHandler 创建回测处理器
func NewBacktestHandler(backtestService *service.BacktestService, strategyService *service.StrategyService, log logger.Logger) *BacktestHandler {
	return &BacktestHandler{
		backtestService: backtestService,
		strategyService: strategyService,
		logger:          log,
	}
}

// RegisterRoutes 注册路由
func (h *BacktestHandler) RegisterRoutes(mux *http.ServeMux) {
	// 回测管理路由
	mux.HandleFunc("GET /api/v1/backtests", h.handleCORS(h.getBacktestsList))
	mux.HandleFunc("POST /api/v1/backtests", h.handleCORS(h.createBacktest))
	mux.HandleFunc("GET /api/v1/backtests/{id}", h.handleCORS(h.getBacktest))
	mux.HandleFunc("PUT /api/v1/backtests/{id}", h.handleCORS(h.updateBacktest))
	mux.HandleFunc("DELETE /api/v1/backtests/{id}", h.handleCORS(h.deleteBacktest))

	// 回测操作路由
	mux.HandleFunc("POST /api/v1/backtests/{id}/start", h.handleCORS(h.startBacktest))
	mux.HandleFunc("POST /api/v1/backtests/{id}/cancel", h.handleCORS(h.cancelBacktest))
	mux.HandleFunc("GET /api/v1/backtests/{id}/progress", h.handleCORS(h.getBacktestProgress))
	mux.HandleFunc("GET /api/v1/backtests/{id}/results", h.handleCORS(h.getBacktestResults))
}

// getBacktestsList 获取回测列表
func (h *BacktestHandler) getBacktestsList(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("获取回测列表请求")

	// 解析查询参数
	query := r.URL.Query()
	req := &models.BacktestListRequest{
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
		req.Status = models.BacktestStatus(status)
	}

	if strategyID := query.Get("strategy_id"); strategyID != "" {
		req.StrategyID = strategyID
	}

	if keyword := query.Get("keyword"); keyword != "" {
		req.Keyword = strings.TrimSpace(keyword)
	}

	if startDate := query.Get("start_date"); startDate != "" {
		req.StartDate = startDate
	}

	if endDate := query.Get("end_date"); endDate != "" {
		req.EndDate = endDate
	}

	// 调用服务层
	backtests, total, err := h.backtestService.GetBacktestsList(r.Context(), req)
	if err != nil {
		h.logger.Error("获取回测列表失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "获取回测列表失败", http.StatusInternalServerError)
		return
	}

	// 返回结果
	response := &models.BacktestListResponse{
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
		Items: backtests,
	}

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"data":    response,
		"message": "获取回测列表成功",
	})
}

// createBacktest 创建回测
func (h *BacktestHandler) createBacktest(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("创建回测请求")

	var req models.CreateBacktestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("解析创建回测请求失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "请求格式错误", http.StatusBadRequest)
		return
	}

	// 验证请求参数
	if err := h.validateCreateBacktestRequest(&req); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 兼容性处理：如果使用了旧的单策略字段，转换为多策略格式
	if len(req.StrategyIDs) == 0 && req.StrategyID != "" {
		req.StrategyIDs = []string{req.StrategyID}
	}

	// 验证策略是否存在
	var strategies []*models.Strategy
	for _, strategyID := range req.StrategyIDs {
		strategy, err := h.strategyService.GetStrategy(r.Context(), strategyID)
		if err != nil {
			if err == service.ErrStrategyNotFound {
				h.writeErrorResponse(w, fmt.Sprintf("策略不存在: %s", strategyID), http.StatusBadRequest)
				return
			}
			h.logger.Error("验证策略失败",
				logger.String("strategy_id", strategyID),
				logger.ErrorField(err))
			h.writeErrorResponse(w, "验证策略失败", http.StatusInternalServerError)
			return
		}
		strategies = append(strategies, strategy)
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		h.writeErrorResponse(w, "开始日期格式错误", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		h.writeErrorResponse(w, "结束日期格式错误", http.StatusBadRequest)
		return
	}

	// 创建回测对象
	backtest := &models.Backtest{
		ID:          uuid.New().String(),
		Name:        req.Name,
		StrategyID:  req.StrategyID,  // 保留用于兼容性
		StrategyIDs: req.StrategyIDs, // 多策略ID列表
		Symbols:     req.Symbols,
		StartDate:   startDate,
		EndDate:     endDate,
		InitialCash: req.InitialCash,
		Commission:  req.Commission,
		Slippage:    req.Slippage,
		Benchmark:   req.Benchmark,
		Status:      models.BacktestStatusPending,
		Progress:    0,
		CreatedBy:   "user", // TODO: 从认证信息获取
	}

	// 记录原始名称，用于检查是否被重命名
	originalName := req.Name

	// 调用服务层创建回测
	if err := h.backtestService.CreateBacktest(r.Context(), backtest); err != nil {
		h.logger.Error("创建回测失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "创建回测失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("回测创建成功", logger.String("backtest_id", backtest.ID))

	// 立即启动回测 - 使用独立的上下文，不依赖HTTP请求上下文
	if err := h.backtestService.StartBacktest(context.Background(), backtest, strategies); err != nil {
		h.logger.Error("启动多策略回测失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "回测创建成功但启动失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("多策略回测创建并启动成功", logger.String("backtest_id", backtest.ID))

	// 返回创建的回测信息（包含策略名称）
	var strategyNames []string
	for _, strategy := range strategies {
		strategyNames = append(strategyNames, strategy.Name)
	}
	backtest.StrategyNames = strategyNames

	// 准备响应消息
	message := "回测创建并启动成功"
	if backtest.Name != originalName {
		message = fmt.Sprintf("回测创建并启动成功。由于名称重复，已自动重命名为：%s", backtest.Name)
	}

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"data":    backtest,
		"message": message,
	})
}

// getBacktest 获取回测详情
func (h *BacktestHandler) getBacktest(w http.ResponseWriter, r *http.Request) {
	backtestID := r.PathValue("id")
	if backtestID == "" {
		h.writeErrorResponse(w, "回测ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("获取回测详情请求", logger.String("backtest_id", backtestID))

	backtest, err := h.backtestService.GetBacktest(r.Context(), backtestID)
	if err != nil {
		if err == service.ErrBacktestNotFound {
			h.writeErrorResponse(w, "回测不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("获取回测详情失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "获取回测详情失败", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"data":    backtest,
		"message": "获取回测详情成功",
	})
}

// updateBacktest 更新回测
func (h *BacktestHandler) updateBacktest(w http.ResponseWriter, r *http.Request) {
	backtestID := r.PathValue("id")
	if backtestID == "" {
		h.writeErrorResponse(w, "回测ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("更新回测请求", logger.String("backtest_id", backtestID))

	var req models.UpdateBacktestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("解析更新回测请求失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "请求格式错误", http.StatusBadRequest)
		return
	}

	// 调用服务层
	if err := h.backtestService.UpdateBacktest(r.Context(), backtestID, &req); err != nil {
		if err == service.ErrBacktestNotFound {
			h.writeErrorResponse(w, "回测不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("更新回测失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "更新回测失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("回测更新成功", logger.String("backtest_id", backtestID))

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "回测更新成功",
	})
}

// deleteBacktest 删除回测
func (h *BacktestHandler) deleteBacktest(w http.ResponseWriter, r *http.Request) {
	backtestID := r.PathValue("id")
	if backtestID == "" {
		h.writeErrorResponse(w, "回测ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("删除回测请求", logger.String("backtest_id", backtestID))

	// 调用服务层
	if err := h.backtestService.DeleteBacktest(r.Context(), backtestID); err != nil {
		if err == service.ErrBacktestNotFound {
			h.writeErrorResponse(w, "回测不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("删除回测失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "删除回测失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("回测删除成功", logger.String("backtest_id", backtestID))

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "回测删除成功",
	})
}

// startBacktest 启动回测
func (h *BacktestHandler) startBacktest(w http.ResponseWriter, r *http.Request) {
	backtestID := r.PathValue("id")
	if backtestID == "" {
		h.writeErrorResponse(w, "回测ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("启动回测请求", logger.String("backtest_id", backtestID))

	// 获取回测信息
	backtest, err := h.backtestService.GetBacktest(r.Context(), backtestID)
	if err != nil {
		if err == service.ErrBacktestNotFound {
			h.writeErrorResponse(w, "回测不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("获取回测信息失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "获取回测信息失败", http.StatusInternalServerError)
		return
	}

	// 检查回测状态
	if backtest.Status != models.BacktestStatusPending {
		h.writeErrorResponse(w, fmt.Sprintf("回测状态不允许启动: %s", backtest.Status), http.StatusBadRequest)
		return
	}

	// 获取策略信息（支持多策略）
	var strategies []*models.Strategy
	strategyIDs := backtest.StrategyIDs

	// 兼容性处理：如果没有多策略ID但有单策略ID，使用单策略ID
	if len(strategyIDs) == 0 && backtest.StrategyID != "" {
		strategyIDs = []string{backtest.StrategyID}
	}

	for _, strategyID := range strategyIDs {
		strategy, err := h.strategyService.GetStrategy(r.Context(), strategyID)
		if err != nil {
			h.logger.Error("获取策略信息失败",
				logger.String("strategy_id", strategyID),
				logger.ErrorField(err))
			h.writeErrorResponse(w, fmt.Sprintf("获取策略信息失败: %s", strategyID), http.StatusInternalServerError)
			return
		}
		strategies = append(strategies, strategy)
	}

	// 启动回测 - 使用独立的上下文，不依赖HTTP请求上下文
	if err := h.backtestService.StartBacktest(context.Background(), backtest, strategies); err != nil {
		h.logger.Error("启动多策略回测失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "启动回测失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("回测启动成功", logger.String("backtest_id", backtestID))

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":         backtestID,
			"status":     models.BacktestStatusRunning,
			"started_at": time.Now(),
		},
		"message": "回测启动成功",
	})
}

// cancelBacktest 取消回测
func (h *BacktestHandler) cancelBacktest(w http.ResponseWriter, r *http.Request) {
	backtestID := r.PathValue("id")
	if backtestID == "" {
		h.writeErrorResponse(w, "回测ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("取消回测请求", logger.String("backtest_id", backtestID))

	if err := h.backtestService.CancelBacktest(r.Context(), backtestID); err != nil {
		if err == service.ErrBacktestNotFound {
			h.writeErrorResponse(w, "回测不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("取消回测失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "取消回测失败", http.StatusInternalServerError)
		return
	}

	h.logger.Info("回测取消成功", logger.String("backtest_id", backtestID))

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "回测取消成功",
	})
}

// getBacktestProgress 获取回测进度
func (h *BacktestHandler) getBacktestProgress(w http.ResponseWriter, r *http.Request) {
	backtestID := r.PathValue("id")
	if backtestID == "" {
		h.writeErrorResponse(w, "回测ID不能为空", http.StatusBadRequest)
		return
	}

	progress, err := h.backtestService.GetBacktestProgress(r.Context(), backtestID)
	if err != nil {
		if err == service.ErrBacktestNotFound {
			h.writeErrorResponse(w, "回测不存在", http.StatusNotFound)
			return
		}
		h.logger.Error("获取回测进度失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "获取回测进度失败", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"data":    progress,
		"message": "获取回测进度成功",
	})
}

// getBacktestResults 获取回测结果
func (h *BacktestHandler) getBacktestResults(w http.ResponseWriter, r *http.Request) {
	backtestID := r.PathValue("id")
	if backtestID == "" {
		h.writeErrorResponse(w, "回测ID不能为空", http.StatusBadRequest)
		return
	}

	h.logger.Info("获取回测结果请求", logger.String("backtest_id", backtestID))

	results, err := h.backtestService.GetBacktestResults(r.Context(), backtestID)
	if err != nil {
		if err == service.ErrBacktestNotFound {
			h.writeErrorResponse(w, "回测不存在", http.StatusNotFound)
			return
		}
		if err == service.ErrBacktestNotCompleted {
			h.writeErrorResponse(w, "回测尚未完成", http.StatusBadRequest)
			return
		}
		h.logger.Error("获取回测结果失败", logger.ErrorField(err))
		h.writeErrorResponse(w, "获取回测结果失败", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, map[string]interface{}{
		"success": true,
		"data":    results,
		"message": "获取回测结果成功",
	})
}

// validateCreateBacktestRequest 验证创建回测请求
func (h *BacktestHandler) validateCreateBacktestRequest(req *models.CreateBacktestRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("回测名称不能为空")
	}

	if len(req.Name) > 100 {
		return fmt.Errorf("回测名称长度不能超过100个字符")
	}

	// 验证策略ID：支持多策略或单策略（兼容性）
	if len(req.StrategyIDs) == 0 && strings.TrimSpace(req.StrategyID) == "" {
		return fmt.Errorf("至少需要选择一个策略")
	}

	// 验证多策略数量限制
	if len(req.StrategyIDs) > 5 {
		return fmt.Errorf("最多只能选择5个策略")
	}

	// 检查策略ID重复和空值
	strategyIDSet := make(map[string]bool)
	for i, strategyID := range req.StrategyIDs {
		// 清理空白字符
		cleanID := strings.TrimSpace(strategyID)
		if cleanID == "" {
			return fmt.Errorf("策略ID不能为空（位置: %d）", i+1)
		}

		// 更新原数组中的值为清理后的值
		req.StrategyIDs[i] = cleanID

		if strategyIDSet[cleanID] {
			return fmt.Errorf("策略ID重复: %s", cleanID)
		}
		strategyIDSet[cleanID] = true
	}

	// 同时检查兼容性字段
	if req.StrategyID != "" && strings.TrimSpace(req.StrategyID) == "" {
		return fmt.Errorf("单策略ID不能为空字符串")
	}

	if len(req.Symbols) == 0 {
		return fmt.Errorf("股票列表不能为空")
	}

	if len(req.Symbols) > 100 {
		return fmt.Errorf("股票数量不能超过100个")
	}

	// 验证日期格式
	if req.StartDate == "" {
		return fmt.Errorf("开始日期不能为空")
	}

	if req.EndDate == "" {
		return fmt.Errorf("结束日期不能为空")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return fmt.Errorf("开始日期格式错误，应为YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return fmt.Errorf("结束日期格式错误，应为YYYY-MM-DD")
	}

	if !startDate.Before(endDate) {
		return fmt.Errorf("开始日期必须早于结束日期")
	}

	// 验证回测时间范围（不能超过5年）
	if endDate.Sub(startDate).Hours() > 24*365*5 {
		return fmt.Errorf("回测时间范围不能超过5年")
	}

	if req.InitialCash < 10000 {
		return fmt.Errorf("初始资金不能少于10000元")
	}

	if req.InitialCash > 100000000 { // 1亿
		return fmt.Errorf("初始资金不能超过1亿元")
	}

	if req.Commission < 0 || req.Commission > 0.01 {
		return errors.New("手续费率必须在0-1%之间")
	}

	if req.Slippage < 0 || req.Slippage > 0.01 {
		return errors.New("滑点必须在0-1%之间")
	}

	return nil
}

// handleCORS 处理CORS
func (h *BacktestHandler) handleCORS(next http.HandlerFunc) http.HandlerFunc {
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
func (h *BacktestHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse 写入错误响应
func (h *BacktestHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
