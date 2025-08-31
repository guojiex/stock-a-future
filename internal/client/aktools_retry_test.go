package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestAKToolsRetryMechanism 测试AKTools重试机制
func TestAKToolsRetryMechanism(t *testing.T) {
	// 创建一个模拟服务器，前两次返回500错误，第三次返回成功
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		t.Logf("📡 模拟服务器收到第%d次请求", attemptCount)

		if attemptCount <= 2 {
			// 前两次返回500错误
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal Server Error"}`))
			t.Logf("❌ 返回500错误（第%d次）", attemptCount)
		} else {
			// 第三次返回成功
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"日期": "2024-08-31", "开盘": 10.0, "收盘": 11.0}]`))
			t.Logf("✅ 返回成功响应（第%d次）", attemptCount)
		}
	}))
	defer server.Close()

	// 创建AKTools客户端，使用模拟服务器
	client := NewAKToolsClient(server.URL)

	t.Logf("🧪 测试AKTools重试机制")
	t.Logf("🌐 模拟服务器地址: %s", server.URL)

	// 构建测试URL
	testURL := server.URL + "/api/public/stock_zh_a_hist?symbol=600976"

	// 调用带重试的请求方法
	ctx := context.Background()
	startTime := time.Now()

	body, err := client.doRequestWithRetry(ctx, testURL)

	duration := time.Since(startTime)
	t.Logf("⏱️  总耗时: %v", duration)

	if err != nil {
		t.Fatalf("❌ 重试后仍然失败: %v", err)
	}

	if body == nil {
		t.Fatalf("❌ 返回的响应体为nil")
	}

	// 验证响应内容
	responseStr := string(body)
	if !strings.Contains(responseStr, "日期") {
		t.Errorf("❌ 响应内容不正确: %s", responseStr)
	}

	// 验证重试次数
	if attemptCount != 3 {
		t.Errorf("❌ 期望重试3次，实际重试%d次", attemptCount)
	} else {
		t.Logf("✅ 重试机制正常工作，共尝试%d次", attemptCount)
	}

	// 验证总耗时（应该包含重试延迟）
	expectedMinDuration := 3 * time.Second // 1s + 2s 的重试延迟
	if duration < expectedMinDuration {
		t.Logf("⚠️  总耗时(%v)小于预期最小值(%v)，可能重试延迟未生效", duration, expectedMinDuration)
	} else {
		t.Logf("✅ 重试延迟正常工作，总耗时: %v", duration)
	}
}

// TestAKToolsRetryWith4xxError 测试4xx错误不重试
func TestAKToolsRetryWith4xxError(t *testing.T) {
	// 创建一个模拟服务器，始终返回404错误
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		t.Logf("📡 模拟服务器收到第%d次请求", attemptCount)

		// 始终返回404错误
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Not Found"}`))
		t.Logf("❌ 返回404错误（第%d次）", attemptCount)
	}))
	defer server.Close()

	// 创建AKTools客户端
	client := NewAKToolsClient(server.URL)

	t.Logf("🧪 测试4xx错误不重试机制")

	// 构建测试URL
	testURL := server.URL + "/api/public/nonexistent"

	// 调用带重试的请求方法
	ctx := context.Background()
	startTime := time.Now()

	_, err := client.doRequestWithRetry(ctx, testURL)

	duration := time.Since(startTime)
	t.Logf("⏱️  总耗时: %v", duration)

	if err == nil {
		t.Fatalf("❌ 期望返回错误，但成功了")
	}

	// 验证只尝试了1次（不重试4xx错误）
	if attemptCount != 1 {
		t.Errorf("❌ 期望只尝试1次，实际尝试%d次", attemptCount)
	} else {
		t.Logf("✅ 4xx错误不重试机制正常工作，只尝试%d次", attemptCount)
	}

	// 验证错误信息
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("❌ 错误信息不包含404: %v", err)
	} else {
		t.Logf("✅ 错误信息正确: %v", err)
	}
}

// TestAKToolsRetryWithCache 测试重试机制与缓存的配合
func TestAKToolsRetryWithCache(t *testing.T) {
	// 创建一个模拟服务器
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		t.Logf("📡 模拟服务器收到第%d次请求", attemptCount)

		// 第一次返回500，第二次返回成功
		if attemptCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal Server Error"}`))
			t.Logf("❌ 返回500错误（第%d次）", attemptCount)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"日期": "2024-08-31", "开盘": 10.0, "收盘": 11.0}]`))
			t.Logf("✅ 返回成功响应（第%d次）", attemptCount)
		}
	}))
	defer server.Close()

	// 创建AKTools客户端
	client := NewAKToolsClient(server.URL)

	t.Logf("🧪 测试重试机制与缓存配合")

	// 构建测试URL
	testURL := server.URL + "/api/public/stock_zh_a_hist?symbol=600976"

	ctx := context.Background()

	// 第一次调用（会重试并成功）
	t.Logf("🔄 第一次调用...")
	body1, err1 := client.doRequestWithCache(ctx, testURL)
	if err1 != nil {
		t.Fatalf("❌ 第一次调用失败: %v", err1)
	}

	t.Logf("✅ 第一次调用成功，服务器总请求次数: %d", attemptCount)

	// 重置计数器，测试缓存
	serverRequestsBefore := attemptCount

	// 第二次调用（应该从缓存获取）
	t.Logf("🔄 第二次调用...")
	body2, err2 := client.doRequestWithCache(ctx, testURL)
	if err2 != nil {
		t.Fatalf("❌ 第二次调用失败: %v", err2)
	}

	// 验证缓存生效（服务器请求次数不变）
	if attemptCount != serverRequestsBefore {
		t.Errorf("❌ 缓存未生效，服务器又收到了请求")
	} else {
		t.Logf("✅ 缓存生效，服务器请求次数未增加")
	}

	// 验证响应内容一致
	if string(body1) != string(body2) {
		t.Errorf("❌ 两次调用返回的内容不一致")
	} else {
		t.Logf("✅ 缓存数据一致性验证通过")
	}
}

// TestAKToolsRetryMaxAttempts 测试最大重试次数限制
func TestAKToolsRetryMaxAttempts(t *testing.T) {
	// 创建一个模拟服务器，始终返回500错误
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		t.Logf("📡 模拟服务器收到第%d次请求", attemptCount)

		// 始终返回500错误
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		t.Logf("❌ 返回500错误（第%d次）", attemptCount)
	}))
	defer server.Close()

	// 创建AKTools客户端
	client := NewAKToolsClient(server.URL)

	t.Logf("🧪 测试最大重试次数限制")

	// 构建测试URL
	testURL := server.URL + "/api/public/stock_zh_a_hist?symbol=600976"

	// 调用带重试的请求方法
	ctx := context.Background()
	startTime := time.Now()

	_, err := client.doRequestWithRetry(ctx, testURL)

	duration := time.Since(startTime)
	t.Logf("⏱️  总耗时: %v", duration)

	if err == nil {
		t.Fatalf("❌ 期望返回错误，但成功了")
	}

	// 验证重试了最大次数（3次）
	expectedAttempts := 3
	if attemptCount != expectedAttempts {
		t.Errorf("❌ 期望重试%d次，实际重试%d次", expectedAttempts, attemptCount)
	} else {
		t.Logf("✅ 最大重试次数限制正常工作，共尝试%d次", attemptCount)
	}

	// 验证总耗时包含所有重试延迟（1s + 2s = 3s）
	expectedMinDuration := 3 * time.Second
	if duration < expectedMinDuration {
		t.Logf("⚠️  总耗时(%v)小于预期最小值(%v)", duration, expectedMinDuration)
	} else {
		t.Logf("✅ 重试延迟累计正常，总耗时: %v", duration)
	}
}
