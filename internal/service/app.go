package service

// App 应用上下文，存储全局服务实例
type App struct {
	SignalService *SignalService
	// 可以添加其他全局服务
}

// NewApp 创建应用上下文
func NewApp() *App {
	return &App{}
}

// RegisterSignalService 注册信号服务
func (a *App) RegisterSignalService(service *SignalService) {
	a.SignalService = service
}
