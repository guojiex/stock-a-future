package handler

import (
	"encoding/json"
	"net/http"
	"stock-a-future/internal/service"
	"time"
)

// CleanupHandler 数据清理处理器
type CleanupHandler struct {
	cleanupService *service.CleanupService
}

// NewCleanupHandler 创建数据清理处理器
func NewCleanupHandler(cleanupService *service.CleanupService) *CleanupHandler {
	return &CleanupHandler{
		cleanupService: cleanupService,
	}
}

// GetStatus 获取清理服务状态
func (h *CleanupHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if h.cleanupService == nil {
		http.Error(w, "数据清理服务未启用", http.StatusServiceUnavailable)
		return
	}

	status := map[string]interface{}{
		"enabled":          h.cleanupService.IsRunning(),
		"cleanup_interval": h.cleanupService.GetCleanupInterval().String(),
		"retention_days":   h.cleanupService.GetRetentionDays(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// ManualCleanup 手动触发数据清理
func (h *CleanupHandler) ManualCleanup(w http.ResponseWriter, r *http.Request) {
	if h.cleanupService == nil {
		http.Error(w, "数据清理服务未启用", http.StatusServiceUnavailable)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	// 在goroutine中执行清理，避免阻塞请求
	go func() {
		h.cleanupService.PerformManualCleanup()
	}()

	response := map[string]string{
		"message": "数据清理任务已启动，请查看服务器日志了解进度",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

// UpdateConfig 更新清理配置
func (h *CleanupHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	if h.cleanupService == nil {
		http.Error(w, "数据清理服务未启用", http.StatusServiceUnavailable)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "只支持PUT方法", http.StatusMethodNotAllowed)
		return
	}

	var config struct {
		CleanupInterval string `json:"cleanup_interval"`
		RetentionDays   int    `json:"retention_days"`
	}

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	// 更新配置
	if config.CleanupInterval != "" {
		if interval, err := time.ParseDuration(config.CleanupInterval); err == nil {
			h.cleanupService.SetCleanupInterval(interval)
		}
	}

	if config.RetentionDays > 0 {
		h.cleanupService.SetRetentionDays(config.RetentionDays)
	}

	response := map[string]string{
		"message": "清理配置已更新",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
