package service

import (
	"crypto/md5"
	"fmt"
	"log"
	"stock-a-future/internal/models"
	"sync"
	"time"
)

// CacheEntry 缓存条目
type CacheEntry struct {
	Data      []models.StockDaily `json:"data"`       // 缓存的日线数据
	ExpiresAt time.Time           `json:"expires_at"` // 过期时间
	CreatedAt time.Time           `json:"created_at"` // 创建时间
}

// IsExpired 检查缓存是否过期
func (c *CacheEntry) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// DailyCacheService 日线数据缓存服务
type DailyCacheService struct {
	cache         sync.Map      // 并发安全的缓存存储
	defaultTTL    time.Duration // 默认过期时间
	maxCacheAge   time.Duration // 最大缓存时间
	cleanupTicker *time.Ticker  // 清理定时器
	stats         CacheStats    // 缓存统计
	statsMutex    sync.RWMutex  // 统计信息的读写锁
}

// CacheStats 缓存统计信息
type CacheStats struct {
	Hits        int64     `json:"hits"`         // 命中次数
	Misses      int64     `json:"misses"`       // 未命中次数
	Entries     int64     `json:"entries"`      // 缓存条目数
	Evictions   int64     `json:"evictions"`    // 清理次数
	LastCleanup time.Time `json:"last_cleanup"` // 上次清理时间
}

// CacheConfig 缓存配置
type CacheConfig struct {
	DefaultTTL      time.Duration // 默认过期时间，建议 30分钟 到 2小时
	MaxCacheAge     time.Duration // 最大缓存时间，建议 1天
	CleanupInterval time.Duration // 清理间隔，建议 10分钟
}

// DefaultCacheConfig 默认缓存配置
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		DefaultTTL:      1 * time.Hour,    // 默认1小时过期
		MaxCacheAge:     24 * time.Hour,   // 最大缓存1天
		CleanupInterval: 10 * time.Minute, // 每10分钟清理一次
	}
}

// NewDailyCacheService 创建日线数据缓存服务
func NewDailyCacheService(config *CacheConfig) *DailyCacheService {
	if config == nil {
		config = DefaultCacheConfig()
	}

	service := &DailyCacheService{
		cache:       sync.Map{},
		defaultTTL:  config.DefaultTTL,
		maxCacheAge: config.MaxCacheAge,
		stats: CacheStats{
			LastCleanup: time.Now(),
		},
	}

	// 启动定期清理
	service.startCleanupRoutine(config.CleanupInterval)

	log.Printf("日线数据缓存服务已启动 - 默认TTL: %v, 最大缓存时间: %v, 清理间隔: %v",
		config.DefaultTTL, config.MaxCacheAge, config.CleanupInterval)

	return service
}

// generateCacheKey 生成缓存键
func (s *DailyCacheService) generateCacheKey(stockCode, startDate, endDate string) string {
	// 使用MD5哈希来生成固定长度的缓存键，避免键过长
	key := fmt.Sprintf("%s_%s_%s", stockCode, startDate, endDate)
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("daily_%x", hash)
}

// Get 从缓存获取日线数据
func (s *DailyCacheService) Get(stockCode, startDate, endDate string) ([]models.StockDaily, bool) {
	cacheKey := s.generateCacheKey(stockCode, startDate, endDate)

	value, exists := s.cache.Load(cacheKey)
	if !exists {
		s.recordMiss()
		return nil, false
	}

	entry, ok := value.(*CacheEntry)
	if !ok {
		// 缓存数据类型错误，删除并返回未命中
		s.cache.Delete(cacheKey)
		s.recordMiss()
		return nil, false
	}

	// 检查是否过期
	if entry.IsExpired() {
		s.cache.Delete(cacheKey)
		s.recordEviction()
		s.recordMiss()
		return nil, false
	}

	s.recordHit()
	log.Printf("缓存命中: %s (创建时间: %v, 过期时间: %v)",
		cacheKey, entry.CreatedAt.Format("15:04:05"), entry.ExpiresAt.Format("15:04:05"))

	return entry.Data, true
}

// Set 设置日线数据到缓存
func (s *DailyCacheService) Set(stockCode, startDate, endDate string, data []models.StockDaily) {
	s.SetWithTTL(stockCode, startDate, endDate, data, s.defaultTTL)
}

// SetWithTTL 设置日线数据到缓存，指定TTL
func (s *DailyCacheService) SetWithTTL(stockCode, startDate, endDate string, data []models.StockDaily, ttl time.Duration) {
	// 限制TTL不能超过最大缓存时间
	if ttl > s.maxCacheAge {
		ttl = s.maxCacheAge
	}

	cacheKey := s.generateCacheKey(stockCode, startDate, endDate)
	now := time.Now()

	entry := &CacheEntry{
		Data:      data,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}

	s.cache.Store(cacheKey, entry)
	s.recordEntry()

	log.Printf("缓存已设置: %s (TTL: %v, 过期时间: %v, 数据条数: %d)",
		cacheKey, ttl, entry.ExpiresAt.Format("15:04:05"), len(data))
}

// Delete 删除缓存
func (s *DailyCacheService) Delete(stockCode, startDate, endDate string) {
	cacheKey := s.generateCacheKey(stockCode, startDate, endDate)
	s.cache.Delete(cacheKey)
	log.Printf("缓存已删除: %s", cacheKey)
}

// Clear 清空所有缓存
func (s *DailyCacheService) Clear() {
	s.cache.Range(func(key, value interface{}) bool {
		s.cache.Delete(key)
		return true
	})

	s.statsMutex.Lock()
	s.stats.Entries = 0
	s.statsMutex.Unlock()

	log.Printf("所有缓存已清空")
}

// GetStats 获取缓存统计信息
func (s *DailyCacheService) GetStats() CacheStats {
	s.statsMutex.RLock()
	defer s.statsMutex.RUnlock()

	// 实时计算缓存条目数
	entryCount := int64(0)
	s.cache.Range(func(key, value interface{}) bool {
		entryCount++
		return true
	})

	stats := s.stats
	stats.Entries = entryCount
	return stats
}

// GetHitRate 获取缓存命中率
func (s *DailyCacheService) GetHitRate() float64 {
	s.statsMutex.RLock()
	defer s.statsMutex.RUnlock()

	total := s.stats.Hits + s.stats.Misses
	if total == 0 {
		return 0.0
	}
	return float64(s.stats.Hits) / float64(total) * 100
}

// startCleanupRoutine 启动定期清理过期缓存的例程
func (s *DailyCacheService) startCleanupRoutine(interval time.Duration) {
	s.cleanupTicker = time.NewTicker(interval)

	go func() {
		for range s.cleanupTicker.C {
			s.cleanupExpiredEntries()
		}
	}()
}

// cleanupExpiredEntries 清理过期的缓存条目
func (s *DailyCacheService) cleanupExpiredEntries() {
	now := time.Now()
	expiredKeys := make([]interface{}, 0)

	// 找出所有过期的键
	s.cache.Range(func(key, value interface{}) bool {
		if entry, ok := value.(*CacheEntry); ok {
			if entry.IsExpired() {
				expiredKeys = append(expiredKeys, key)
			}
		}
		return true
	})

	// 删除过期的缓存条目
	evictedCount := 0
	for _, key := range expiredKeys {
		s.cache.Delete(key)
		evictedCount++
	}

	if evictedCount > 0 {
		s.statsMutex.Lock()
		s.stats.Evictions += int64(evictedCount)
		s.stats.LastCleanup = now
		s.statsMutex.Unlock()

		log.Printf("清理过期缓存: 删除了 %d 个条目", evictedCount)
	}
}

// Stop 停止缓存服务
func (s *DailyCacheService) Stop() {
	if s.cleanupTicker != nil {
		s.cleanupTicker.Stop()
		log.Printf("日线数据缓存服务已停止")
	}
}

// 统计记录方法
func (s *DailyCacheService) recordHit() {
	s.statsMutex.Lock()
	s.stats.Hits++
	s.statsMutex.Unlock()
}

func (s *DailyCacheService) recordMiss() {
	s.statsMutex.Lock()
	s.stats.Misses++
	s.statsMutex.Unlock()
}

func (s *DailyCacheService) recordEntry() {
	s.statsMutex.Lock()
	s.stats.Entries++
	s.statsMutex.Unlock()
}

func (s *DailyCacheService) recordEviction() {
	s.statsMutex.Lock()
	s.stats.Evictions++
	s.statsMutex.Unlock()
}

// GetCacheInfo 获取缓存详细信息（用于调试）
func (s *DailyCacheService) GetCacheInfo() map[string]interface{} {
	stats := s.GetStats()
	hitRate := s.GetHitRate()

	return map[string]interface{}{
		"stats":    stats,
		"hit_rate": fmt.Sprintf("%.2f%%", hitRate),
		"config": map[string]interface{}{
			"default_ttl":   s.defaultTTL.String(),
			"max_cache_age": s.maxCacheAge.String(),
		},
	}
}
