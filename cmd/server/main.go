package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"stock-a-future/config"
	"stock-a-future/internal/client"
	"stock-a-future/internal/handler"
	"stock-a-future/internal/service"
	"sync/atomic"
)

// 健康检查请求计数器，用于抽样记录日志
var healthCheckCounter int64

func main() {
	// 加载配置
	cfg := config.Load()

	log.Printf("启动Stock-A-Future API服务器...")
	log.Printf("服务器地址: %s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Tushare API: %s", cfg.TushareBaseURL)

	// 创建Tushare客户端
	tushareClient := client.NewTushareClient(cfg.TushareToken, cfg.TushareBaseURL)

	// 创建缓存服务
	var cacheService *service.DailyCacheService
	if cfg.CacheEnabled {
		cacheConfig := &service.CacheConfig{
			DefaultTTL:      cfg.CacheDefaultTTL,
			MaxCacheAge:     cfg.CacheMaxAge,
			CleanupInterval: cfg.CacheCleanupInterval,
		}
		cacheService = service.NewDailyCacheService(cacheConfig)
		log.Printf("日线数据缓存已启用")
	} else {
		log.Printf("日线数据缓存已禁用")
	}

	// 创建处理器
	stockHandler := handler.NewStockHandler(tushareClient, cacheService)

	// 创建路由器
	mux := http.NewServeMux()

	// 注册路由
	registerRoutes(mux, stockHandler)

	// 添加静态文件服务
	registerStaticRoutes(mux)

	// 添加中间件
	handler := withCORS(withLogging(mux))

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("服务器正在监听 %s", addr)
	log.Printf("API文档:")
	log.Printf("  健康检查: GET http://%s/api/v1/health", addr)
	log.Printf("  股票基本信息: GET http://%s/api/v1/stocks/{code}/basic", addr)
	log.Printf("  股票日线: GET http://%s/api/v1/stocks/{code}/daily?start_date=20240101&end_date=20240131", addr)
	log.Printf("  技术指标: GET http://%s/api/v1/stocks/{code}/indicators", addr)
	log.Printf("  买卖预测: GET http://%s/api/v1/stocks/{code}/predictions", addr)
	log.Printf("  本地股票列表: GET http://%s/api/v1/stocks", addr)
	log.Printf("  刷新本地数据: POST http://%s/api/v1/stocks/refresh", addr)
	if cfg.CacheEnabled {
		log.Printf("  缓存统计: GET http://%s/api/v1/cache/stats", addr)
		log.Printf("  清空缓存: DELETE http://%s/api/v1/cache", addr)
	}
	log.Printf("  Web客户端: http://%s/", addr)
	log.Printf("示例: curl http://%s/api/v1/stocks/000001.SZ/daily", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// registerRoutes 注册路由
func registerRoutes(mux *http.ServeMux, stockHandler *handler.StockHandler) {
	// 健康检查
	mux.HandleFunc("GET /api/v1/health", stockHandler.GetHealthStatus)

	// 股票相关API
	mux.HandleFunc("GET /api/v1/stocks/{code}/basic", stockHandler.GetStockBasic)
	mux.HandleFunc("GET /api/v1/stocks/{code}/daily", stockHandler.GetDailyData)
	mux.HandleFunc("GET /api/v1/stocks/{code}/indicators", stockHandler.GetIndicators)
	mux.HandleFunc("GET /api/v1/stocks/{code}/predictions", stockHandler.GetPredictions)

	// 本地股票数据API
	mux.HandleFunc("GET /api/v1/stocks", stockHandler.GetStockList)
	mux.HandleFunc("GET /api/v1/stocks/search", stockHandler.SearchStocks)
	mux.HandleFunc("POST /api/v1/stocks/refresh", stockHandler.RefreshLocalData)

	// 缓存管理API
	mux.HandleFunc("GET /api/v1/cache/stats", stockHandler.GetCacheStats)
	mux.HandleFunc("DELETE /api/v1/cache", stockHandler.ClearCache)

}

// registerStaticRoutes 注册静态文件路由
func registerStaticRoutes(mux *http.ServeMux) {
	// 获取静态文件目录的绝对路径
	webClientDir := filepath.Join("web", "static")

	// 服务静态文件
	fileServer := http.FileServer(http.Dir(webClientDir))

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
					"refresh": "POST /api/v1/stocks/refresh"
				},
				"web_client": "GET /",
				"example": "curl http://localhost:8080/api/v1/stocks/000001.SZ/daily"
			}`)); err != nil {
				log.Printf("写入API路径响应失败: %v", err)
			}
			return
		}

		// 其他路径，尝试服务静态文件
		fileServer.ServeHTTP(w, r)
	})

	// 明确处理静态资源路径
	mux.Handle("GET /js/", http.StripPrefix("/", fileServer))
	mux.Handle("GET /styles.css", http.StripPrefix("/", fileServer))
	mux.Handle("GET /test_modules.html", http.StripPrefix("/", fileServer))
}

// withLogging 日志中间件
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 对health接口进行抽样记录，每100次请求只记录1次
		if r.URL.Path == "/api/v1/health" {
			count := atomic.AddInt64(&healthCheckCounter, 1)
			if count%100 == 1 { // 第1次、第101次、第201次...记录日志
				log.Printf("%s %s %s (健康检查 #%d)", r.Method, r.URL.Path, r.RemoteAddr, count)
			}
		} else {
			// 非health接口正常记录日志
			log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		}
		next.ServeHTTP(w, r)
	})
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

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
