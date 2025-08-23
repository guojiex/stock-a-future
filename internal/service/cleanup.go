package service

import (
	"stock-a-future/internal/logger"
	"time"
)

// CleanupService 数据清理服务
type CleanupService struct {
	databaseService *DatabaseService
	cleanupInterval time.Duration
	retentionDays   int
	stopChan        chan struct{}
	isRunning       bool
}

// NewCleanupService 创建数据清理服务
func NewCleanupService(databaseService *DatabaseService, cleanupInterval time.Duration, retentionDays int) *CleanupService {
	return &CleanupService{
		databaseService: databaseService,
		cleanupInterval: cleanupInterval,
		retentionDays:   retentionDays,
		stopChan:        make(chan struct{}),
		isRunning:       false,
	}
}

// Start 启动数据清理服务
func (s *CleanupService) Start() {
	if s.isRunning {
		logger.Warn("数据清理服务已在运行中")
		return
	}

	s.isRunning = true
	logger.Info("启动数据清理服务", logger.Duration("interval", s.cleanupInterval))

	// 启动时立即执行一次清理
	go s.performCleanup()

	// 启动定期清理任务
	go s.runCleanupLoop()
}

// Stop 停止数据清理服务
func (s *CleanupService) Stop() {
	if !s.isRunning {
		return
	}

	logger.Info("停止数据清理服务...")
	close(s.stopChan)
	s.isRunning = false
}

// runCleanupLoop 运行清理循环
func (s *CleanupService) runCleanupLoop() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go s.performCleanup()
		case <-s.stopChan:
			logger.Info("数据清理服务已停止")
			return
		}
	}
}

// performCleanup 执行数据清理
func (s *CleanupService) performCleanup() {
	startTime := time.Now()
	logger.Info("开始执行数据清理任务...")

	// 清理数据库过期数据
	if err := s.databaseService.CleanupExpiredData(s.retentionDays); err != nil {
		logger.ErrorLog("数据库清理失败", logger.ErrorField(err))
	}

	// 获取清理后的数据库统计信息
	stats, err := s.databaseService.GetDatabaseStats()
	if err != nil {
		logger.Warn("获取数据库统计信息失败", logger.ErrorField(err))
	} else {
		logger.Info("数据库统计信息:")
		for key, value := range stats {
			logger.Info("  统计项", logger.String("key", key), logger.Any("value", value))
		}
	}

	duration := time.Since(startTime)
	logger.Info("数据清理任务完成", logger.Duration("duration", duration))
}

// IsRunning 检查服务是否正在运行
func (s *CleanupService) IsRunning() bool {
	return s.isRunning
}

// GetCleanupInterval 获取清理间隔
func (s *CleanupService) GetCleanupInterval() time.Duration {
	return s.cleanupInterval
}

// SetCleanupInterval 设置清理间隔
func (s *CleanupService) SetCleanupInterval(interval time.Duration) {
	s.cleanupInterval = interval
	logger.Info("数据清理间隔已更新", logger.Duration("interval", interval))
}

// GetRetentionDays 获取股票信号数据保留天数
func (s *CleanupService) GetRetentionDays() int {
	return s.retentionDays
}

// SetRetentionDays 设置股票信号数据保留天数
func (s *CleanupService) SetRetentionDays(days int) {
	s.retentionDays = days
	logger.Info("股票信号数据保留天数已更新", logger.Int("retention_days", days))
}

// PerformManualCleanup 手动执行数据清理
func (s *CleanupService) PerformManualCleanup() {
	logger.Info("手动触发数据清理任务...")
	s.performCleanup()
}
