package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置结构
type Config struct {
	// Tushare API配置
	TushareToken   string
	TushareBaseURL string

	// 服务器配置
	ServerPort string
	ServerHost string

	// 日志配置
	LogLevel string
}

// Load 加载配置
func Load() *Config {
	// 尝试加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用环境变量")
	}

	config := &Config{
		TushareToken:   getEnv("TUSHARE_TOKEN", ""),
		TushareBaseURL: getEnv("TUSHARE_BASE_URL", "http://api.tushare.pro"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		ServerHost:     getEnv("SERVER_HOST", "localhost"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}

	// 验证必要配置
	if config.TushareToken == "" {
		log.Fatal("TUSHARE_TOKEN 环境变量是必需的")
	}

	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
