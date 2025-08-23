package service

import (
	"stock-a-future/internal/models"
	"testing"
)

// MockStockService 模拟股票服务接口
type MockStockService struct {
	shouldFail bool
}

func (m *MockStockService) GetDailyData(tsCode, startDate, endDate, adjust string) ([]models.StockDaily, error) {
	if m.shouldFail {
		return nil, &MockError{message: "模拟错误"}
	}
	return []models.StockDaily{}, nil
}

// MockError 模拟错误类型
type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

func TestNewSignalService(t *testing.T) {
	// 测试创建信号服务
	patternService := &PatternService{}
	stockService := &MockStockService{}
	favoriteService := &FavoriteService{}

	service, err := NewSignalService("test_data", patternService, stockService, favoriteService)

	if err != nil {
		t.Fatalf("创建SignalService失败: %v", err)
	}

	if service == nil {
		t.Fatal("NewSignalService应该返回非空的SignalService实例")
	}

	if service.db == nil {
		t.Error("SignalService应该包含数据库服务")
	}

	if service.patternService != patternService {
		t.Error("PatternService应该正确设置")
	}

	if service.stockService != stockService {
		t.Error("StockService应该正确设置")
	}

	if service.favoriteService != favoriteService {
		t.Error("FavoriteService应该正确设置")
	}

	if service.running {
		t.Error("新创建的SignalService不应该处于运行状态")
	}

	if service.calculationStatus == nil {
		t.Error("CalculationStatus应该被初始化")
	}
}

func TestSignalService_StartStop(t *testing.T) {
	patternService := &PatternService{}
	stockService := &MockStockService{}
	favoriteService := &FavoriteService{}

	service, err := NewSignalService("test_data", patternService, stockService, favoriteService)
	if err != nil {
		t.Fatalf("创建SignalService失败: %v", err)
	}

	// 测试启动服务
	service.Start()

	if !service.running {
		t.Error("Start()后服务应该处于运行状态")
	}

	// 测试重复启动
	service.Start()
	if !service.running {
		t.Error("重复启动后服务应该仍然处于运行状态")
	}

	// 测试停止服务
	service.Stop()

	if service.running {
		t.Error("Stop()后服务应该停止运行")
	}

	// 测试重复停止
	service.Stop()
	if service.running {
		t.Error("重复停止后服务应该仍然停止")
	}
}

func TestSignalService_GetCalculationStatus(t *testing.T) {
	patternService := &PatternService{}
	stockService := &MockStockService{}
	favoriteService := &FavoriteService{}

	service, err := NewSignalService("test_data", patternService, stockService, favoriteService)
	if err != nil {
		t.Fatalf("创建SignalService失败: %v", err)
	}

	status := service.GetCalculationStatus()

	if status == nil {
		t.Fatal("GetCalculationStatus应该返回非空的状态")
	}

	if status.IsCalculating {
		t.Error("新创建的服务状态不应该显示正在计算")
	}

	if status.Total != 0 {
		t.Error("新创建的服务总任务数应该为0")
	}

	if status.Completed != 0 {
		t.Error("新创建的服务已完成任务数应该为0")
	}

	if status.Failed != 0 {
		t.Error("新创建的服务失败任务数应该为0")
	}
}

func TestSignalService_Close(t *testing.T) {
	patternService := &PatternService{}
	stockService := &MockStockService{}
	favoriteService := &FavoriteService{}

	service, err := NewSignalService("test_data", patternService, stockService, favoriteService)
	if err != nil {
		t.Fatalf("创建SignalService失败: %v", err)
	}

	// 启动服务
	service.Start()

	// 关闭服务
	err = service.Close()
	if err != nil {
		t.Errorf("Close()应该成功执行: %v", err)
	}

	if service.running {
		t.Error("Close()后服务应该停止运行")
	}
}

func TestSignalService_ConcurrentAccess(t *testing.T) {
	patternService := &PatternService{}
	stockService := &MockStockService{}
	favoriteService := &FavoriteService{}

	service, err := NewSignalService("test_data", patternService, stockService, favoriteService)
	if err != nil {
		t.Fatalf("创建SignalService失败: %v", err)
	}

	// 测试并发访问状态查询
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			status := service.GetCalculationStatus()
			if status == nil {
				t.Error("并发访问时GetCalculationStatus应该返回非空状态")
			}
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 测试启动和停止
	service.Start()
	if !service.running {
		t.Error("Start()后服务应该处于运行状态")
	}

	service.Stop()
	if service.running {
		t.Error("Stop()后服务应该停止运行")
	}
}
