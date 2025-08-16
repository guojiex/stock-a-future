package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"stock-a-future/internal/models"
	"testing"
	"time"
)

// MockDataSourceClient 模拟数据源客户端
type MockDataSourceClient struct {
	shouldFail bool
	baseURL    string
}

func (m *MockDataSourceClient) GetDailyData(symbol, startDate, endDate, adjust string) ([]models.StockDaily, error) {
	return nil, nil
}

func (m *MockDataSourceClient) GetDailyDataByDate(tradeDate string) ([]models.StockDaily, error) {
	return nil, nil
}

func (m *MockDataSourceClient) GetStockBasic(symbol string) (*models.StockBasic, error) {
	return nil, nil
}

func (m *MockDataSourceClient) GetStockList() ([]models.StockBasic, error) {
	return nil, nil
}

func (m *MockDataSourceClient) GetBaseURL() string {
	return m.baseURL
}

func (m *MockDataSourceClient) TestConnection() error {
	if m.shouldFail {
		return errors.New("connection failed")
	}
	return nil
}

func TestStockHandler_GetHealthStatus(t *testing.T) {
	tests := []struct {
		name            string
		checkConnection bool
		shouldFail      bool
		expectedStatus  string
		expectedCode    int
	}{
		{
			name:            "健康状态检查 - 数据源正常",
			checkConnection: false,
			shouldFail:      false,
			expectedStatus:  "healthy",
			expectedCode:    200,
		},
		{
			name:            "健康状态检查 - 数据源异常",
			checkConnection: true,
			shouldFail:      true,
			expectedStatus:  "degraded",
			expectedCode:    503,
		},
		{
			name:            "健康状态检查 - 强制检查连接且正常",
			checkConnection: true,
			shouldFail:      false,
			expectedStatus:  "healthy",
			expectedCode:    200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟数据源客户端
			mockClient := &MockDataSourceClient{
				shouldFail: tt.shouldFail,
				baseURL:    "http://127.0.0.1:8080",
			}

			// 创建处理器
			handler := &StockHandler{
				dataSourceClient: mockClient,
			}

			// 创建测试请求
			req, err := http.NewRequest("GET", "/api/v1/health", nil)
			if err != nil {
				t.Fatal(err)
			}

			// 添加查询参数
			if tt.checkConnection {
				q := req.URL.Query()
				q.Add("check_connection", "true")
				req.URL.RawQuery = q.Encode()
			}

			// 创建响应记录器
			rr := httptest.NewRecorder()

			// 调用处理器
			handler.GetHealthStatus(rr, req)

			// 检查状态码
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			// 检查响应格式
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatal("Failed to parse response JSON:", err)
			}

			// 检查success字段
			if success, exists := response["success"]; !exists {
				t.Error("Response missing 'success' field")
			} else if success != true {
				t.Errorf("Expected success=true, got %v", success)
			}

			// 检查data字段
			data, exists := response["data"]
			if !exists {
				t.Error("Response missing 'data' field")
				return
			}

			// 检查健康状态
			healthData, ok := data.(map[string]interface{})
			if !ok {
				t.Error("Data field is not a map")
				return
			}

			if status, exists := healthData["status"]; !exists {
				t.Error("Health data missing 'status' field")
			} else if status != tt.expectedStatus {
				t.Errorf("Expected status=%s, got %v", tt.expectedStatus, status)
			}

			// 检查时间戳
			if timestamp, exists := healthData["timestamp"]; !exists {
				t.Error("Health data missing 'timestamp' field")
			} else {
				// 验证时间戳格式
				if _, err := time.Parse(time.RFC3339, timestamp.(string)); err != nil {
					t.Errorf("Invalid timestamp format: %v", err)
				}
			}

			// 检查服务状态
			services, exists := healthData["services"]
			if !exists {
				t.Error("Health data missing 'services' field")
				return
			}

			servicesData, ok := services.(map[string]interface{})
			if !ok {
				t.Error("Services field is not a map")
				return
			}

			// 检查数据源状态
			dataSource, exists := servicesData["data_source"]
			if !exists {
				t.Error("Services missing 'data_source' field")
				return
			}

			dataSourceData, ok := dataSource.(map[string]interface{})
			if !ok {
				t.Error("Data source field is not a map")
				return
			}

			if url, exists := dataSourceData["url"]; !exists {
				t.Error("Data source missing 'url' field")
			} else if url != mockClient.baseURL {
				t.Errorf("Expected URL %s, got %v", mockClient.baseURL, url)
			}
		})
	}
}

func TestStockHandler_WriteSuccessResponse(t *testing.T) {
	handler := &StockHandler{}

	tests := []struct {
		name       string
		data       interface{}
		statusCode int
	}{
		{
			name:       "标准成功响应",
			data:       map[string]string{"message": "success"},
			statusCode: 200,
		},
		{
			name:       "带状态码的成功响应",
			data:       map[string]string{"message": "created"},
			statusCode: 201,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建响应记录器
			rr := httptest.NewRecorder()

			// 调用相应的方法
			if tt.statusCode == 200 {
				handler.writeSuccessResponse(rr, tt.data)
			} else {
				handler.writeSuccessResponseWithStatus(rr, tt.data, tt.statusCode)
			}

			// 检查状态码
			if status := rr.Code; status != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.statusCode)
			}

			// 检查响应格式
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatal("Failed to parse response JSON:", err)
			}

			// 检查success字段
			if success, exists := response["success"]; !exists {
				t.Error("Response missing 'success' field")
			} else if success != true {
				t.Errorf("Expected success=true, got %v", success)
			}

			// 检查data字段
			if data, exists := response["data"]; !exists {
				t.Error("Response missing 'data' field")
			} else if data == nil {
				t.Error("Data field is nil")
			}
		})
	}
}
