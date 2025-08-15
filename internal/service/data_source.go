package service

import (
	"fmt"
	"log"
	"sync"

	"stock-a-future/config"
	"stock-a-future/internal/client"
)

// DataSourceService 数据源服务
type DataSourceService struct {
	config        *config.Config
	factory       *client.DataSourceFactory
	currentClient client.DataSourceClient
	mutex         sync.RWMutex
}

// NewDataSourceService 创建数据源服务
func NewDataSourceService(cfg *config.Config) *DataSourceService {
	return &DataSourceService{
		config:  cfg,
		factory: client.NewDataSourceFactory(),
	}
}

// GetClient 获取当前数据源客户端
func (s *DataSourceService) GetClient() (client.DataSourceClient, error) {
	s.mutex.RLock()
	if s.currentClient != nil {
		defer s.mutex.RUnlock()
		return s.currentClient, nil
	}
	s.mutex.RUnlock()

	// 需要初始化客户端
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 双重检查
	if s.currentClient != nil {
		return s.currentClient, nil
	}

	// 创建新的客户端
	client, err := s.createClient()
	if err != nil {
		return nil, fmt.Errorf("创建数据源客户端失败: %w", err)
	}

	s.currentClient = client
	return client, nil
}

// createClient 根据配置创建数据源客户端
func (s *DataSourceService) createClient() (client.DataSourceClient, error) {
	sourceType := client.DataSourceType(s.config.DataSourceType)

	var config map[string]string

	switch sourceType {
	case client.DataSourceTushare:
		config = map[string]string{
			"token":    s.config.TushareToken,
			"base_url": s.config.TushareBaseURL,
		}

	case client.DataSourceAKTools:
		config = map[string]string{
			"base_url": s.config.AKToolsBaseURL,
		}

	default:
		return nil, fmt.Errorf("不支持的数据源类型: %s", s.config.DataSourceType)
	}

	client, err := s.factory.CreateClient(sourceType, config)
	if err != nil {
		return nil, fmt.Errorf("创建%s客户端失败: %w", sourceType, err)
	}

	log.Printf("=== 数据源客户端创建成功 ===")
	log.Printf("类型: %s", sourceType)
	log.Printf("基础URL: %s", client.GetBaseURL())
	log.Printf("========================")
	return client, nil
}

// SwitchDataSource 切换数据源
func (s *DataSourceService) SwitchDataSource(sourceType string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 更新配置
	s.config.DataSourceType = sourceType

	// 关闭当前客户端
	if s.currentClient != nil {
		log.Printf("关闭当前数据源客户端")
		s.currentClient = nil
	}

	// 创建新的客户端
	client, err := s.createClient()
	if err != nil {
		return fmt.Errorf("切换数据源失败: %w", err)
	}

	s.currentClient = client
	log.Printf("成功切换到%s数据源", sourceType)
	return nil
}

// TestConnection 测试当前数据源连接
func (s *DataSourceService) TestConnection() error {
	client, err := s.GetClient()
	if err != nil {
		return fmt.Errorf("获取数据源客户端失败: %w", err)
	}

	return client.TestConnection()
}

// GetDataSourceInfo 获取当前数据源信息
func (s *DataSourceService) GetDataSourceInfo() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	info := map[string]interface{}{
		"type": s.config.DataSourceType,
	}

	if s.currentClient != nil {
		info["base_url"] = s.currentClient.GetBaseURL()
		info["status"] = "connected"
	} else {
		info["status"] = "disconnected"
	}

	return info
}
