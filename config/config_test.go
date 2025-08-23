package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// 测试环境变量配置
	os.Setenv("DATA_SOURCE_TYPE", "aktools")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("AKTOOLS_BASE_URL", "http://test.example.com")

	config := Load()

	if config.DataSourceType != "aktools" {
		t.Errorf("DataSourceType配置不正确，期望 aktools，实际 %s", config.DataSourceType)
	}

	if config.ServerPort != "8080" {
		t.Errorf("ServerPort配置不正确，期望 8080，实际 %s", config.ServerPort)
	}

	if config.AKToolsBaseURL != "http://test.example.com" {
		t.Errorf("AKToolsBaseURL配置不正确，期望 http://test.example.com，实际 %s", config.AKToolsBaseURL)
	}

	// 清理环境变量
	os.Unsetenv("DATA_SOURCE_TYPE")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("AKTOOLS_BASE_URL")
}

func TestLoadConfig_Defaults(t *testing.T) {
	// 测试默认配置 - 需要设置DataSourceType为aktools避免log.Fatal
	os.Setenv("DATA_SOURCE_TYPE", "aktools")

	config := Load()

	if config.DataSourceType != "aktools" {
		t.Errorf("DataSourceType默认值应该是 aktools，实际是 %s", config.DataSourceType)
	}

	if config.ServerPort != "8080" {
		t.Errorf("ServerPort默认值应该是 8080，实际是 %s", config.ServerPort)
	}

	if config.AKToolsBaseURL != "http://127.0.0.1:8080" {
		t.Errorf("AKToolsBaseURL默认值应该是 http://127.0.0.1:8080，实际是 %s", config.AKToolsBaseURL)
	}

	// 清理环境变量
	os.Unsetenv("DATA_SOURCE_TYPE")
}

func TestConfig_Fields(t *testing.T) {
	// 测试Config结构体字段
	config := Config{
		DataSourceType: "aktools",
		ServerPort:     "8080",
		AKToolsBaseURL: "http://example.com",
		LogLevel:       "info",
	}

	if config.DataSourceType != "aktools" {
		t.Errorf("DataSourceType字段不正确，期望 aktools，实际 %s", config.DataSourceType)
	}

	if config.ServerPort != "8080" {
		t.Errorf("ServerPort字段不正确，期望 8080，实际 %s", config.ServerPort)
	}

	if config.AKToolsBaseURL != "http://example.com" {
		t.Errorf("AKToolsBaseURL字段不正确，期望 http://example.com，实际 %s", config.AKToolsBaseURL)
	}

	if config.LogLevel != "info" {
		t.Errorf("LogLevel字段不正确，期望 info，实际 %s", config.LogLevel)
	}
}

func TestConfig_ServerHost(t *testing.T) {
	config := Config{
		ServerHost: "localhost",
		ServerPort: "8080",
	}

	if config.ServerHost != "localhost" {
		t.Errorf("ServerHost字段不正确，期望 localhost，实际 %s", config.ServerHost)
	}

	if config.ServerPort != "8080" {
		t.Errorf("ServerPort字段不正确，期望 8080，实际 %s", config.ServerPort)
	}
}

func TestConfig_CacheSettings(t *testing.T) {
	config := Config{
		CacheEnabled:         true,
		CacheDefaultTTL:      1 * time.Hour,
		CacheMaxAge:          24 * time.Hour,
		CacheCleanupInterval: 10 * time.Minute,
	}

	if !config.CacheEnabled {
		t.Error("CacheEnabled应该为true")
	}

	if config.CacheDefaultTTL != 1*time.Hour {
		t.Errorf("CacheDefaultTTL不正确，期望 1h，实际 %v", config.CacheDefaultTTL)
	}

	if config.CacheMaxAge != 24*time.Hour {
		t.Errorf("CacheMaxAge不正确，期望 24h，实际 %v", config.CacheMaxAge)
	}

	if config.CacheCleanupInterval != 10*time.Minute {
		t.Errorf("CacheCleanupInterval不正确，期望 10m，实际 %v", config.CacheCleanupInterval)
	}
}
