package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 应用配置结构
type Config struct {
	// 数据源配置
	DataSourceType string // 数据源类型: tushare, aktools

	// Tushare API配置
	TushareToken   string
	TushareBaseURL string

	// AKTools配置
	AKToolsBaseURL string

	// 服务器配置
	ServerPort string
	ServerHost string

	// 日志配置
	LogLevel string

	// 缓存配置
	CacheDefaultTTL      time.Duration // 默认缓存过期时间
	CacheMaxAge          time.Duration // 最大缓存时间
	CacheCleanupInterval time.Duration // 清理间隔
	CacheEnabled         bool          // 是否启用缓存

	// 模式预测配置
	PatternPredictionDays int // 模式预测的时间窗口（天数），默认14天（两周）
}

// Load 加载配置
func Load() *Config {
	// 尝试加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用环境变量")
	}

	config := &Config{
		DataSourceType: getEnv("DATA_SOURCE_TYPE", "tushare"),
		TushareToken:   getEnv("TUSHARE_TOKEN", ""),
		TushareBaseURL: getEnv("TUSHARE_BASE_URL", "http://api.tushare.pro"),
		AKToolsBaseURL: getEnv("AKTOOLS_BASE_URL", "http://127.0.0.1:8080"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		ServerHost:     getEnv("SERVER_HOST", "localhost"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),

		// 缓存配置 - 默认值
		CacheDefaultTTL:      getDurationEnv("CACHE_DEFAULT_TTL", 1*time.Hour),         // 默认1小时
		CacheMaxAge:          getDurationEnv("CACHE_MAX_AGE", 24*time.Hour),            // 最大24小时
		CacheCleanupInterval: getDurationEnv("CACHE_CLEANUP_INTERVAL", 10*time.Minute), // 每10分钟清理
		CacheEnabled:         getBoolEnv("CACHE_ENABLED", true),                        // 默认启用缓存

		// 模式预测配置 - 默认值
		PatternPredictionDays: getIntEnv("PATTERN_PREDICTION_DAYS", 14), // 默认14天
	}

	// 验证必要配置
	if config.DataSourceType == "tushare" && config.TushareToken == "" {
		log.Fatal("使用Tushare数据源时，TUSHARE_TOKEN 环境变量是必需的")
	}

	// 打印数据源配置信息
	log.Printf("=== 数据源配置信息 ===")
	log.Printf("数据源类型: %s", config.DataSourceType)
	log.Printf("Tushare基础URL: %s", config.TushareBaseURL)
	log.Printf("AKTools基础URL: %s", config.AKToolsBaseURL)
	if config.DataSourceType == "tushare" {
		log.Printf("Tushare Token: %s", maskToken(config.TushareToken))
	}
	log.Printf("模式预测时间窗口: %d天", config.PatternPredictionDays)
	log.Printf("=====================")

	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBoolEnv 获取布尔类型环境变量
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
		log.Printf("警告: 无法解析环境变量 %s=%s 为布尔值，使用默认值 %v", key, value, defaultValue)
	}
	return defaultValue
}

// getDurationEnv 获取时间间隔类型环境变量
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
		log.Printf("警告: 无法解析环境变量 %s=%s 为时间间隔，使用默认值 %v", key, value, defaultValue)
	}
	return defaultValue
}

// getIntEnv 获取整数类型环境变量
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
		log.Printf("警告: 无法解析环境变量 %s=%s 为整数，使用默认值 %d", key, value, defaultValue)
	}
	return defaultValue
}

// maskToken 掩码Token，只显示前4位和后4位
func maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
