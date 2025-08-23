package service

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	app := NewApp()

	if app == nil {
		t.Fatal("NewApp() 应该返回非空的App实例")
	}

	if app.SignalService != nil {
		t.Error("新创建的App实例的SignalService应该为nil")
	}

	if app.CleanupService != nil {
		t.Error("新创建的App实例的CleanupService应该为nil")
	}
}

func TestApp_RegisterSignalService(t *testing.T) {
	app := NewApp()
	signalService := &SignalService{}

	app.RegisterSignalService(signalService)

	if app.SignalService != signalService {
		t.Error("RegisterSignalService应该正确设置SignalService")
	}
}

func TestApp_RegisterCleanupService(t *testing.T) {
	app := NewApp()
	cleanupService := &CleanupService{}

	app.RegisterCleanupService(cleanupService)

	if app.CleanupService != cleanupService {
		t.Error("RegisterCleanupService应该正确设置CleanupService")
	}
}

func TestApp_ServiceRegistrationOrder(t *testing.T) {
	app := NewApp()
	signalService1 := &SignalService{}
	signalService2 := &SignalService{}

	// 测试服务注册的顺序性
	app.RegisterSignalService(signalService1)
	if app.SignalService != signalService1 {
		t.Error("第一次注册SignalService失败")
	}

	app.RegisterSignalService(signalService2)
	if app.SignalService != signalService2 {
		t.Error("第二次注册SignalService应该覆盖第一次")
	}
}
