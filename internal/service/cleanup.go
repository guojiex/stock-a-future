package service

import (
	"log"
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
		log.Printf("数据清理服务已在运行中")
		return
	}

	s.isRunning = true
	log.Printf("启动数据清理服务，清理间隔: %v", s.cleanupInterval)

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

	log.Printf("停止数据清理服务...")
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
			log.Printf("数据清理服务已停止")
			return
		}
	}
}

// performCleanup 执行数据清理
func (s *CleanupService) performCleanup() {
	startTime := time.Now()
	log.Printf("开始执行数据清理任务...")

	// 清理数据库过期数据
	if err := s.databaseService.CleanupExpiredData(s.retentionDays); err != nil {
		log.Printf("数据库清理失败: %v", err)
	}

	// 获取清理后的数据库统计信息
	stats, err := s.databaseService.GetDatabaseStats()
	if err != nil {
		log.Printf("获取数据库统计信息失败: %v", err)
	} else {
		log.Printf("数据库统计信息:")
		for key, value := range stats {
			log.Printf("  %s: %v", key, value)
		}
	}

	duration := time.Since(startTime)
	log.Printf("数据清理任务完成，耗时: %v", duration)
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
	log.Printf("数据清理间隔已更新为: %v", interval)
}

// GetRetentionDays 获取股票信号数据保留天数
func (s *CleanupService) GetRetentionDays() int {
	return s.retentionDays
}

// SetRetentionDays 设置股票信号数据保留天数
func (s *CleanupService) SetRetentionDays(days int) {
	s.retentionDays = days
	log.Printf("股票信号数据保留天数已更新为: %d天", days)
}

// PerformManualCleanup 手动执行数据清理
func (s *CleanupService) PerformManualCleanup() {
	log.Printf("手动触发数据清理任务...")
	s.performCleanup()
}
