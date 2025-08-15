package client

import (
	"fmt"
	"stock-a-future/internal/models"
)

// DataSourceClient 数据源客户端接口
type DataSourceClient interface {
	// 获取股票日线数据
	GetDailyData(symbol, startDate, endDate string) ([]models.StockDaily, error)

	// 根据交易日期获取所有股票数据
	GetDailyDataByDate(tradeDate string) ([]models.StockDaily, error)

	// 获取股票基本信息
	GetStockBasic(symbol string) (*models.StockBasic, error)

	// 获取股票列表
	GetStockList() ([]models.StockBasic, error)

	// 获取基础URL
	GetBaseURL() string

	// 测试连接
	TestConnection() error
}

// DataSourceType 数据源类型
type DataSourceType string

const (
	DataSourceTushare DataSourceType = "tushare"
	DataSourceAKTools DataSourceType = "aktools"
)

// DataSourceFactory 数据源工厂
type DataSourceFactory struct{}

// NewDataSourceFactory 创建数据源工厂
func NewDataSourceFactory() *DataSourceFactory {
	return &DataSourceFactory{}
}

// CreateClient 根据类型创建数据源客户端
func (f *DataSourceFactory) CreateClient(sourceType DataSourceType, config map[string]string) (DataSourceClient, error) {
	switch sourceType {
	case DataSourceTushare:
		token := config["token"]
		baseURL := config["base_url"]
		if baseURL == "" {
			baseURL = "http://api.tushare.pro"
		}
		return NewTushareClient(token, baseURL), nil

	case DataSourceAKTools:
		baseURL := config["base_url"]
		if baseURL == "" {
			baseURL = "http://127.0.0.1:8080"
		}
		return NewAKToolsClient(baseURL), nil

	default:
		return nil, fmt.Errorf("不支持的数据源类型: %s", sourceType)
	}
}
