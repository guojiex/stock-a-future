package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"stock-a-future/config"
	"stock-a-future/internal/client"
	"stock-a-future/internal/handler"
	"stock-a-future/internal/logger"
	"stock-a-future/internal/service"
	"strings"
	"syscall"
	"time"
)

// 日志相关配置
var (
	// 不记录日志的路径列表
	skipLogPaths = []string{
		"/api/v1/health",
		"/api/v1/groups",
		"/api/v1/favorites",
		"/api/v1/favorites/signals",
	}
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化日志系统
	logConfig := &logger.Config{
		Level:         cfg.LogLevel,
		Format:        cfg.LogFormat,
		Output:        cfg.LogOutput,
		Filename:      cfg.LogFilename,
		MaxSize:       cfg.LogMaxSize,
		MaxBackups:    cfg.LogMaxBackups,
		MaxAge:        cfg.LogMaxAge,
		Compress:      cfg.LogCompress,
		ConsoleFormat: cfg.LogConsoleFormat,
		ShowCaller:    cfg.LogShowCaller,
		ShowTime:      cfg.LogShowTime,
	}

	if err := logger.InitGlobalLogger(logConfig); err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}

	logger.Info("启动Stock-A-Future API服务器...")
	logger.Info("服务器地址", logger.String("host", cfg.ServerHost), logger.String("port", cfg.ServerPort))

	// 打印数据源配置信息
	logger.Info("=== 数据源配置信息 ===")
	logger.Info("数据源类型", logger.String("type", cfg.DataSourceType))
	logger.Info("Tushare基础URL", logger.String("url", cfg.TushareBaseURL))
	logger.Info("AKTools基础URL", logger.String("url", cfg.AKToolsBaseURL))
	if cfg.DataSourceType == "tushare" {
		logger.Info("Tushare Token", logger.String("token", config.MaskToken(cfg.TushareToken)))
	}
	logger.Info("模式预测时间窗口", logger.Int("days", cfg.PatternPredictionDays))
	logger.Info("数据清理配置:")
	logger.Info("  启用", logger.Bool("enabled", cfg.CleanupEnabled))
	logger.Info("  间隔", logger.Duration("interval", cfg.CleanupInterval))
	logger.Info("  股票信号保留", logger.Int("retention_days", cfg.CleanupRetentionDays))
	logger.Info("=====================")

	// 使用数据源工厂创建客户端
	factory := client.NewDataSourceFactory()

	var dataSourceClient client.DataSourceClient
	var err error

	switch cfg.DataSourceType {
	case "tushare":
		logger.Info("配置Tushare数据源", logger.String("base_url", cfg.TushareBaseURL))

		config := map[string]string{
			"token":    cfg.TushareToken,
			"base_url": cfg.TushareBaseURL,
		}
		dataSourceClient, err = factory.CreateClient(client.DataSourceTushare, config)
		if err != nil {
			logger.Fatal("创建Tushare客户端失败", logger.ErrorField(err))
		}

		// 启动时测试Tushare连接
		logger.Info("正在测试Tushare API连接...")
		if err := dataSourceClient.TestConnection(); err != nil {
			logger.Fatal("Tushare API连接测试失败", logger.ErrorField(err))
		}
		logger.Info("✓ Tushare API连接测试成功")
	case "aktools":
		logger.Info("配置AKTools数据源", logger.String("base_url", cfg.AKToolsBaseURL))

		config := map[string]string{
			"base_url": cfg.AKToolsBaseURL,
		}
		dataSourceClient, err = factory.CreateClient(client.DataSourceAKTools, config)
		if err != nil {
			logger.Fatal("创建AKTools客户端失败", logger.ErrorField(err))
		}

		// 启动时测试AKTools连接
		logger.Info("正在测试AKTools API连接...")
		if err := dataSourceClient.TestConnection(); err != nil {
			logger.Fatal("AKTools API连接测试失败", logger.ErrorField(err))
		}
		logger.Info("✓ AKTools API连接测试成功")
	default:
		logger.Fatal("不支持的数据源类型", logger.String("type", cfg.DataSourceType))
	}

	// 创建缓存服务
	var cacheService *service.DailyCacheService
	if cfg.CacheEnabled {
		cacheConfig := &service.CacheConfig{
			DefaultTTL:      cfg.CacheDefaultTTL,
			MaxCacheAge:     cfg.CacheMaxAge,
			CleanupInterval: cfg.CacheCleanupInterval,
		}
		cacheService = service.NewDailyCacheService(cacheConfig)
		logger.Info("日线数据缓存已启用")
	} else {
		logger.Info("日线数据缓存已禁用")
	}

	// 创建收藏服务
	favoriteService, err := service.NewFavoriteService("data")
	if err != nil {
		logger.Fatal("创建收藏服务失败", logger.ErrorField(err))
	}

	// 创建数据库服务
	databaseService, err := service.NewDatabaseService("data")
	if err != nil {
		logger.Fatal("创建数据库服务失败", logger.ErrorField(err))
	}

	// 创建最近查看服务
	recentViewService := service.NewRecentViewService(databaseService.GetDB())
	logger.Info("✓ 最近查看服务已创建")

	// 启动自动清理过期记录任务（每小时执行一次）
	recentViewService.StartAutoCleanup(1 * time.Hour)
	logger.Info("✓ 最近查看自动清理任务已启动")

	// 创建数据清理服务
	var cleanupService *service.CleanupService
	if cfg.CleanupEnabled {
		cleanupService = service.NewCleanupService(
			databaseService,
			cfg.CleanupInterval,
			cfg.CleanupRetentionDays,
		)
		logger.Info("数据清理服务已创建", logger.Duration("interval", cfg.CleanupInterval))
	} else {
		logger.Info("数据清理服务已禁用")
	}

	// 创建图形识别服务
	patternService := service.NewPatternService(dataSourceClient)

	// 创建应用上下文
	app := service.NewApp()

	// 创建信号计算服务
	signalService, err := service.NewSignalService("data", patternService, dataSourceClient, favoriteService)
	if err != nil {
		logger.Fatal("创建信号计算服务失败", logger.ErrorField(err))
	}

	// 注册服务到应用上下文
	app.RegisterSignalService(signalService)
	if cleanupService != nil {
		app.RegisterCleanupService(cleanupService)
	}

	// 启动异步信号计算服务
	signalService.Start()
	logger.Info("✓ 异步信号计算服务已启动")

	// 启动数据清理服务
	if cleanupService != nil {
		cleanupService.Start()
		logger.Info("✓ 数据清理服务已启动")
	}

	// 创建数据源服务
	dataSourceService := service.NewDataSourceService(cfg)
	logger.Info("✓ 数据源服务已创建")

	// 创建策略服务
	strategyService := service.NewStrategyService(logger.GetGlobalLogger())
	logger.Info("✓ 策略服务已创建")

	// 创建回测服务
	backtestService := service.NewBacktestService(strategyService, dataSourceService, cacheService, logger.GetGlobalLogger())
	logger.Info("✓ 回测服务已创建")

	// 创建参数优化服务
	parameterOptimizer := service.NewParameterOptimizer(backtestService, strategyService, logger.GetGlobalLogger())
	logger.Info("✓ 参数优化服务已创建")

	// 创建处理器
	stockHandler := handler.NewStockHandler(dataSourceClient, cacheService, favoriteService, recentViewService, app)
	patternHandler := handler.NewPatternHandler(patternService)
	signalHandler := handler.NewSignalHandler(signalService)
	strategyHandler := handler.NewStrategyHandler(strategyService, logger.GetGlobalLogger())
	backtestHandler := handler.NewBacktestHandler(backtestService, strategyService, logger.GetGlobalLogger())
	parameterOptimizerHandler := handler.NewParameterOptimizerHandler(parameterOptimizer, logger.GetGlobalLogger())

	// 创建清理处理器
	var cleanupHandler *handler.CleanupHandler
	if cleanupService != nil {
		cleanupHandler = handler.NewCleanupHandler(cleanupService)
	}

	// 创建路由器
	mux := http.NewServeMux()

	// 注册路由
	registerRoutes(mux, stockHandler, patternHandler, signalHandler, strategyHandler, backtestHandler, parameterOptimizerHandler, cleanupHandler)

	// 添加静态文件服务
	registerStaticRoutes(mux)

	// 添加中间件
	handler := withCORS(withLogging(mux))

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	logger.Info("服务器正在监听", logger.String("address", addr))
	logger.Info("API文档:")
	logger.Info("  健康检查: GET http://" + addr + "/api/v1/health")
	logger.Info("  股票基本信息: GET http://" + addr + "/api/v1/stocks/{code}/basic")
	logger.Info("  股票日线: GET http://" + addr + "/api/v1/stocks/{code}/daily?start_date=20240101&end_date=20240131")
	logger.Info("  技术指标: GET http://" + addr + "/api/v1/stocks/{code}/indicators")
	logger.Info("  买卖预测: GET http://" + addr + "/api/v1/stocks/{code}/predictions")

	logger.Info("  综合基本面: GET http://" + addr + "/api/v1/stocks/{code}/fundamental?period=2023-12-31")
	logger.Info("  利润表: GET http://" + addr + "/api/v1/stocks/{code}/income?period=2023-12-31")
	logger.Info("  资产负债表: GET http://" + addr + "/api/v1/stocks/{code}/balance?period=2023-12-31")
	logger.Info("  现金流量表: GET http://" + addr + "/api/v1/stocks/{code}/cashflow?period=2023-12-31")
	logger.Info("  每日基本面: GET http://" + addr + "/api/v1/stocks/{code}/dailybasic?trade_date=20240101")
	logger.Info("  本地股票列表: GET http://" + addr + "/api/v1/stocks")
	logger.Info("  股票搜索: GET http://" + addr + "/api/v1/stocks/search?q=keyword")
	logger.Info("  刷新本地数据: POST http://" + addr + "/api/v1/stocks/refresh")
	logger.Info("  收藏列表: GET http://" + addr + "/api/v1/favorites")
	logger.Info("  添加收藏: POST http://" + addr + "/api/v1/favorites")
	logger.Info("  删除收藏: DELETE http://" + addr + "/api/v1/favorites/{id}")
	logger.Info("  检查收藏: GET http://" + addr + "/api/v1/favorites/check/{code}")
	logger.Info("  最近查看: GET http://" + addr + "/api/v1/recent-views")
	logger.Info("  添加查看: POST http://" + addr + "/api/v1/recent-views")
	logger.Info("  图形识别: GET http://" + addr + "/api/v1/patterns/recognize?ts_code=000001.SZ")
	logger.Info("  图形搜索: POST http://" + addr + "/api/v1/patterns/search")
	logger.Info("  图形摘要: GET http://" + addr + "/api/v1/patterns/summary?ts_code=000001.SZ")
	logger.Info("  最近信号: GET http://" + addr + "/api/v1/patterns/recent?ts_code=000001.SZ")
	logger.Info("  可用图形: GET http://" + addr + "/api/v1/patterns/available")
	logger.Info("  图形统计: GET http://" + addr + "/api/v1/patterns/statistics?ts_code=000001.SZ")
	logger.Info("  计算信号: POST http://" + addr + "/api/v1/signals/calculate")
	logger.Info("  批量计算: POST http://" + addr + "/api/v1/signals/batch")
	logger.Info("  获取信号: GET http://" + addr + "/api/v1/signals/{code}?signal_date=20240101")
	logger.Info("  最新信号: GET http://" + addr + "/api/v1/signals?limit=20")
	if cfg.CacheEnabled {
		logger.Info("  缓存统计: GET http://" + addr + "/api/v1/cache/stats")
		logger.Info("  清空缓存: DELETE http://" + addr + "/api/v1/cache")
	}
	if cfg.CleanupEnabled {
		logger.Info("  清理状态: GET http://" + addr + "/api/v1/cleanup/status")
		logger.Info("  手动清理: POST http://" + addr + "/api/v1/cleanup/manual")
		logger.Info("  更新配置: PUT http://" + addr + "/api/v1/cleanup/config")
	}
	logger.Info("  Web客户端: http://" + addr + "/")
	logger.Info("示例: curl http://" + addr + "/api/v1/stocks/000001.SZ/daily")

	// 设置信号处理，用于优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 在单独的goroutine中启动服务器
	go func() {
		if err := http.ListenAndServe(addr, handler); err != nil {
			logger.Fatal("服务器启动失败", logger.ErrorField(err))
		}
	}()

	// 等待退出信号
	<-sigChan
	logger.Info("收到退出信号，开始优雅关闭...")

	// 停止信号计算服务
	signalService.Stop()

	// 停止数据清理服务
	if cleanupService != nil {
		cleanupService.Stop()
	}

	// 关闭收藏服务
	if err := favoriteService.Close(); err != nil {
		logger.Warn("关闭收藏服务失败", logger.ErrorField(err))
	}

	// 关闭信号服务
	if err := signalService.Close(); err != nil {
		logger.Warn("关闭信号服务失败", logger.ErrorField(err))
	}

	// 关闭数据库服务
	if err := databaseService.Close(); err != nil {
		logger.Warn("关闭数据库服务失败", logger.ErrorField(err))
	}

	logger.Info("服务器已优雅关闭")
}

// registerRoutes 注册路由
func registerRoutes(mux *http.ServeMux, stockHandler *handler.StockHandler, patternHandler *handler.PatternHandler, signalHandler *handler.SignalHandler, strategyHandler *handler.StrategyHandler, backtestHandler *handler.BacktestHandler, parameterOptimizerHandler *handler.ParameterOptimizerHandler, cleanupHandler *handler.CleanupHandler) {
	// 健康检查
	mux.HandleFunc("GET /api/v1/health", stockHandler.GetHealthStatus)

	// 股票相关API
	mux.HandleFunc("GET /api/v1/stocks/{code}/basic", stockHandler.GetStockBasic)
	mux.HandleFunc("GET /api/v1/stocks/{code}/daily", stockHandler.GetDailyData)
	mux.HandleFunc("GET /api/v1/stocks/{code}/indicators", stockHandler.GetIndicators)
	mux.HandleFunc("GET /api/v1/stocks/{code}/predictions", stockHandler.GetPredictions)

	// 基本面数据API
	mux.HandleFunc("GET /api/v1/stocks/{code}/fundamental", stockHandler.GetFundamentalData)
	mux.HandleFunc("GET /api/v1/stocks/{code}/income", stockHandler.GetIncomeStatement)
	mux.HandleFunc("GET /api/v1/stocks/{code}/balance", stockHandler.GetBalanceSheet)
	mux.HandleFunc("GET /api/v1/stocks/{code}/cashflow", stockHandler.GetCashFlowStatement)
	mux.HandleFunc("GET /api/v1/stocks/{code}/dailybasic", stockHandler.GetDailyBasic)

	// 基本面因子分析API
	mux.HandleFunc("GET /api/v1/stocks/{code}/factor", stockHandler.GetFundamentalFactor)
	mux.HandleFunc("GET /api/v1/factors/ranking", stockHandler.GetFundamentalFactorRanking)
	mux.HandleFunc("POST /api/v1/factors/batch", stockHandler.BatchCalculateFundamentalFactors)

	// 本地股票数据API
	mux.HandleFunc("GET /api/v1/stocks", stockHandler.GetStockList)
	mux.HandleFunc("GET /api/v1/stocks/search", stockHandler.SearchStocks)
	mux.HandleFunc("POST /api/v1/stocks/refresh", stockHandler.RefreshLocalData)

	// 缓存管理API
	mux.HandleFunc("GET /api/v1/cache/stats", stockHandler.GetCacheStats)
	mux.HandleFunc("DELETE /api/v1/cache", stockHandler.ClearCache)

	// 收藏股票API
	mux.HandleFunc("GET /api/v1/favorites", stockHandler.GetFavorites)
	mux.HandleFunc("POST /api/v1/favorites", stockHandler.AddFavorite)
	mux.HandleFunc("DELETE /api/v1/favorites/{id}", stockHandler.DeleteFavorite)
	mux.HandleFunc("PUT /api/v1/favorites/{id}", stockHandler.UpdateFavorite)
	mux.HandleFunc("GET /api/v1/favorites/check/{code}", stockHandler.CheckFavorite)
	mux.HandleFunc("PUT /api/v1/favorites/order", stockHandler.UpdateFavoritesOrder)
	mux.HandleFunc("GET /api/v1/favorites/signals", stockHandler.GetFavoritesSignals)

	// 分组管理API
	mux.HandleFunc("GET /api/v1/groups", stockHandler.GetGroups)
	mux.HandleFunc("POST /api/v1/groups", stockHandler.CreateGroup)
	mux.HandleFunc("PUT /api/v1/groups/{id}", stockHandler.UpdateGroup)
	mux.HandleFunc("DELETE /api/v1/groups/{id}", stockHandler.DeleteGroup)

	// 最近查看API
	mux.HandleFunc("GET /api/v1/recent-views", stockHandler.GetRecentViews)
	mux.HandleFunc("POST /api/v1/recent-views", stockHandler.AddRecentView)
	mux.HandleFunc("DELETE /api/v1/recent-views/{code}", stockHandler.DeleteRecentView)
	mux.HandleFunc("POST /api/v1/recent-views/cleanup", stockHandler.ClearExpiredRecentViews)
	mux.HandleFunc("DELETE /api/v1/recent-views", stockHandler.ClearAllRecentViews)

	// 图形识别API
	mux.HandleFunc("GET /api/v1/patterns/recognize", patternHandler.RecognizePatterns)
	mux.HandleFunc("POST /api/v1/patterns/search", patternHandler.SearchPatterns)
	mux.HandleFunc("GET /api/v1/patterns/summary", patternHandler.GetPatternSummary)
	mux.HandleFunc("GET /api/v1/patterns/recent", patternHandler.GetRecentSignals)
	mux.HandleFunc("GET /api/v1/patterns/available", patternHandler.GetAvailablePatterns)
	mux.HandleFunc("GET /api/v1/patterns/statistics", patternHandler.GetPatternStatistics)

	// 信号计算API
	mux.HandleFunc("POST /api/v1/signals/calculate", signalHandler.CalculateSignal)
	mux.HandleFunc("POST /api/v1/signals/batch", signalHandler.BatchCalculateSignals)
	mux.HandleFunc("GET /api/v1/signals/{code}", signalHandler.GetSignal)
	mux.HandleFunc("GET /api/v1/signals", signalHandler.GetLatestSignals)
	mux.HandleFunc("GET /api/v1/signals/status", signalHandler.GetCalculationStatus)

	// 策略管理API
	strategyHandler.RegisterRoutes(mux)

	// 回测管理API
	backtestHandler.RegisterRoutes(mux)

	// 参数优化API
	parameterOptimizerHandler.RegisterRoutes(mux)

	// 数据清理API
	if cleanupHandler != nil {
		mux.HandleFunc("GET /api/v1/cleanup/status", cleanupHandler.GetStatus)
		mux.HandleFunc("POST /api/v1/cleanup/manual", cleanupHandler.ManualCleanup)
		mux.HandleFunc("PUT /api/v1/cleanup/config", cleanupHandler.UpdateConfig)
	}
}

// registerStaticRoutes 注册静态文件路由
func registerStaticRoutes(mux *http.ServeMux) {
	// 获取静态文件目录的绝对路径
	webClientDir := filepath.Join("web", "static")

	// 服务静态文件
	fileServer := http.FileServer(http.Dir(webClientDir))

	// 专门处理favicon.ico请求，重定向到favicon.png
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// 设置缓存头，避免重复请求
		w.Header().Set("Cache-Control", "public, max-age=86400") // 24小时缓存
		w.Header().Set("Content-Type", "image/png")

		// 重定向到favicon.png
		http.Redirect(w, r, "/favicon.png", http.StatusMovedPermanently)
	})

	// 处理根路径，返回index.html
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// 如果是根路径，服务index.html
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(webClientDir, "index.html"))
			return
		}

		// 检查是否是API路径，如果是则返回API信息
		if r.URL.Path == "/api" || r.URL.Path == "/api/" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{
				"service": "Stock-A-Future API",
				"version": "1.0.0",
				"description": "A股股票买卖点预测API服务",
				"endpoints": {
					"health": "GET /api/v1/health",
					"basic": "GET /api/v1/stocks/{code}/basic",
					"daily": "GET /api/v1/stocks/{code}/daily",
					"indicators": "GET /api/v1/stocks/{code}/indicators", 
					"predictions": "GET /api/v1/stocks/{code}/predictions",
					"stocks": "GET /api/v1/stocks",
					"refresh": "POST /api/v1/stocks/refresh",
					"favorites": "GET /api/v1/favorites",
					"favorites_signals": "GET /api/v1/favorites/signals"
				},
				"web_client": "GET /",
				"example": "curl http://localhost:8080/api/v1/stocks/{code}/daily"
			}`)); err != nil {
				logger.Warn("写入API路径响应失败", logger.ErrorField(err))
			}
			return
		}

		// 其他路径，尝试服务静态文件
		fileServer.ServeHTTP(w, r)
	})

	// 明确处理静态资源路径
	mux.Handle("GET /js/", http.StripPrefix("/", fileServer))
	mux.Handle("GET /styles.css", http.StripPrefix("/", fileServer))
	mux.Handle("GET /favicon.png", http.StripPrefix("/", fileServer))
}

// withLogging 日志中间件
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// startTime := time.Now()

		// 生成请求ID
		requestID := logger.GenerateRequestID()
		ctx := logger.WithRequestIDContext(r.Context(), requestID)
		r = r.WithContext(ctx)

		// 检查是否跳过日志记录
		for _, path := range skipLogPaths {
			if r.URL.Path == path {
				// 包装ResponseWriter以捕获状态码
				wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
				next.ServeHTTP(wrappedWriter, r)
				return
			}
		}

		// 跳过静态资源请求的日志记录
		if isStaticResource(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// 记录请求开始
		// logger.GetGlobalLogger().InfoCtx(ctx, "HTTP请求开始",
		// 	logger.String("method", r.Method),
		// 	logger.String("path", r.URL.Path),
		// 	logger.String("remote_addr", r.RemoteAddr),
		// 	logger.String("user_agent", r.UserAgent()),
		// )

		// 包装ResponseWriter以捕获状态码
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		// 计算响应时间
		// responseTime := time.Since(startTime)

		// 对于304状态码不记录日志
		if wrappedWriter.statusCode == http.StatusNotModified {
			return
		}

		// 记录响应信息
		// logger.GetGlobalLogger().InfoCtx(ctx, "HTTP请求完成",
		// 	logger.String("method", r.Method),
		// 	logger.String("path", r.URL.Path),
		// 	logger.String("remote_addr", r.RemoteAddr),
		// 	logger.Int("status_code", wrappedWriter.statusCode),
		// 	logger.Duration("response_time", responseTime),
		// )
	})
}

// isStaticResource 判断是否为静态资源请求
func isStaticResource(path string) bool {
	// 静态资源文件扩展名
	staticExtensions := []string{
		".js", ".css", ".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg",
		".woff", ".woff2", ".ttf", ".eot", ".pdf", ".txt", ".xml",
	}

	// 检查路径是否包含静态资源扩展名
	for _, ext := range staticExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	// 检查是否为常见的静态资源路径
	staticPaths := []string{
		"/js/", "/css/", "/images/", "/fonts/", "/static/", "/assets/",
		"/favicon.ico", "/robots.txt", "/sitemap.xml",
	}

	for _, staticPath := range staticPaths {
		if strings.HasPrefix(path, staticPath) {
			return true
		}
	}

	return false
}

// responseWriter 包装ResponseWriter以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// withCORS CORS中间件
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 设置内容安全策略(CSP)头，解决CSP错误
		// 允许外部CDN资源加载，允许内联脚本执行
		cspHeader := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' https:; " +
			"connect-src 'self' https:; " +
			"frame-src 'none'; " +
			"object-src 'none';"
		w.Header().Set("Content-Security-Policy", cspHeader)

		// 为API请求设置Content-Type
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Content-Type", "application/json")
		}

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
