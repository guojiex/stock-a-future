package main

import (
	"fmt"
	"log"
	"net/http"
	"stock-a-future/config"
	"stock-a-future/internal/client"
	"stock-a-future/internal/handler"
)

func main() {
	// 加载配置
	cfg := config.Load()

	log.Printf("启动Stock-A-Future API服务器...")
	log.Printf("服务器地址: %s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Tushare API: %s", cfg.TushareBaseURL)

	// 创建Tushare客户端
	tushareClient := client.NewTushareClient(cfg.TushareToken, cfg.TushareBaseURL)

	// 创建处理器
	stockHandler := handler.NewStockHandler(tushareClient)

	// 创建路由器
	mux := http.NewServeMux()

	// 注册路由
	registerRoutes(mux, stockHandler)

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
	mux.HandleFunc("POST /api/v1/stocks/refresh", stockHandler.RefreshLocalData)

	// 根路径
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
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
			"example": "curl http://localhost:8080/api/v1/stocks/000001.SZ/daily"
		}`)); err != nil {
			log.Printf("写入根路径响应失败: %v", err)
		}
	})
}

// withLogging 日志中间件
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
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

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
