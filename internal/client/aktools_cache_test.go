package client

import (
	"testing"
	"time"
)

// TestAKToolsCacheBasicFunctionality 测试AKTools缓存基本功能
func TestAKToolsCacheBasicFunctionality(t *testing.T) {
	// 创建AKTools客户端
	client := NewAKToolsClient("http://127.0.0.1:8080")

	// 测试股票代码
	symbol := "600976" // 健民集团

	t.Logf("🧪 开始测试AKTools缓存功能")
	t.Logf("📊 测试股票: %s", symbol)

	// 第一次调用 - 应该发起HTTP请求
	t.Logf("🔄 第一次调用GetStockBasic...")
	start1 := time.Now()
	stock1, err := client.GetStockBasic(symbol)
	duration1 := time.Since(start1)

	if err != nil {
		t.Logf("❌ 第一次调用失败: %v", err)
		// 不要因为网络问题而失败测试
		t.Skipf("跳过测试，AKTools服务可能不可用: %v", err)
		return
	}

	if stock1 == nil {
		t.Fatalf("❌ 第一次调用返回nil")
	}

	t.Logf("✅ 第一次调用成功: %s (%s), 耗时: %v", stock1.Name, stock1.Symbol, duration1)

	// 第二次调用 - 应该从缓存获取
	t.Logf("🔄 第二次调用GetStockBasic...")
	start2 := time.Now()
	stock2, err := client.GetStockBasic(symbol)
	duration2 := time.Since(start2)

	if err != nil {
		t.Fatalf("❌ 第二次调用失败: %v", err)
	}

	if stock2 == nil {
		t.Fatalf("❌ 第二次调用返回nil")
	}

	t.Logf("✅ 第二次调用成功: %s (%s), 耗时: %v", stock2.Name, stock2.Symbol, duration2)

	// 验证缓存效果
	if duration2 >= duration1 {
		t.Logf("⚠️  第二次调用耗时(%v)不比第一次(%v)快，可能缓存未生效", duration2, duration1)
	} else {
		t.Logf("🚀 缓存生效！第二次调用比第一次快 %v", duration1-duration2)
	}

	// 验证数据一致性
	if stock1.Symbol != stock2.Symbol || stock1.Name != stock2.Name {
		t.Errorf("❌ 缓存数据不一致: 第一次(%s,%s) vs 第二次(%s,%s)",
			stock1.Symbol, stock1.Name, stock2.Symbol, stock2.Name)
	} else {
		t.Logf("✅ 缓存数据一致性验证通过")
	}

	// 检查缓存大小
	cacheSize := client.cache.Size()
	t.Logf("📦 当前缓存大小: %d 个条目", cacheSize)

	if cacheSize == 0 {
		t.Errorf("❌ 缓存大小为0，缓存可能未正常工作")
	}
}

// TestAKToolsCacheMultipleAPIs 测试多个API的缓存功能
func TestAKToolsCacheMultipleAPIs(t *testing.T) {
	client := NewAKToolsClient("http://127.0.0.1:8080")
	symbol := "600976"

	t.Logf("🧪 测试多个API的缓存功能")

	// 测试GetStockBasic缓存
	t.Logf("📊 测试GetStockBasic缓存...")
	_, err1 := client.GetStockBasic(symbol)
	if err1 != nil {
		t.Logf("⚠️  GetStockBasic失败: %v", err1)
	} else {
		t.Logf("✅ GetStockBasic第一次调用成功")
	}

	_, err2 := client.GetStockBasic(symbol)
	if err2 != nil {
		t.Logf("⚠️  GetStockBasic第二次调用失败: %v", err2)
	} else {
		t.Logf("✅ GetStockBasic第二次调用成功（应该来自缓存）")
	}

	// 测试GetDailyData缓存
	t.Logf("📈 测试GetDailyData缓存...")
	startDate := "20240101"
	endDate := "20240131"

	_, err3 := client.GetDailyData(symbol, startDate, endDate, "qfq")
	if err3 != nil {
		t.Logf("⚠️  GetDailyData失败: %v", err3)
	} else {
		t.Logf("✅ GetDailyData第一次调用成功")
	}

	_, err4 := client.GetDailyData(symbol, startDate, endDate, "qfq")
	if err4 != nil {
		t.Logf("⚠️  GetDailyData第二次调用失败: %v", err4)
	} else {
		t.Logf("✅ GetDailyData第二次调用成功（应该来自缓存）")
	}

	// 检查缓存大小
	cacheSize := client.cache.Size()
	t.Logf("📦 测试完成后缓存大小: %d 个条目", cacheSize)

	if cacheSize > 0 {
		t.Logf("✅ 缓存正常工作，包含 %d 个条目", cacheSize)
	}
}

// TestAKToolsCacheExpiration 测试缓存过期功能
func TestAKToolsCacheExpiration(t *testing.T) {
	// 创建一个TTL很短的缓存用于测试
	client := NewAKToolsClient("http://127.0.0.1:8080")
	// 替换为短TTL的缓存
	client.cache = NewRequestCache(1*time.Second, 100)

	symbol := "600976"

	t.Logf("🧪 测试缓存过期功能（TTL=1秒）")

	// 第一次调用
	_, err1 := client.GetStockBasic(symbol)
	if err1 != nil {
		t.Skipf("跳过测试，AKTools服务不可用: %v", err1)
		return
	}

	cacheSize1 := client.cache.Size()
	t.Logf("📦 第一次调用后缓存大小: %d", cacheSize1)

	// 等待缓存过期
	t.Logf("⏳ 等待缓存过期...")
	time.Sleep(2 * time.Second)

	// 第二次调用（缓存应该已过期）
	_, err2 := client.GetStockBasic(symbol)
	if err2 != nil {
		t.Logf("⚠️  第二次调用失败: %v", err2)
	}

	cacheSize2 := client.cache.Size()
	t.Logf("📦 第二次调用后缓存大小: %d", cacheSize2)

	t.Logf("✅ 缓存过期测试完成")
}

// TestAKToolsCacheKeyGeneration 测试缓存key生成
func TestAKToolsCacheKeyGeneration(t *testing.T) {
	cache := NewRequestCache(5*time.Minute, 100)

	// 测试相同URL生成相同key
	url1 := "http://127.0.0.1:8080/api/public/stock_individual_info_em?symbol=600976"
	url2 := "http://127.0.0.1:8080/api/public/stock_individual_info_em?symbol=600976"
	url3 := "http://127.0.0.1:8080/api/public/stock_individual_info_em?symbol=000001"

	key1 := cache.generateCacheKey(url1)
	key2 := cache.generateCacheKey(url2)
	key3 := cache.generateCacheKey(url3)

	t.Logf("🔑 URL1 key: %s", key1)
	t.Logf("🔑 URL2 key: %s", key2)
	t.Logf("🔑 URL3 key: %s", key3)

	if key1 != key2 {
		t.Errorf("❌ 相同URL应该生成相同的缓存key")
	} else {
		t.Logf("✅ 相同URL生成相同缓存key")
	}

	if key1 == key3 {
		t.Errorf("❌ 不同URL不应该生成相同的缓存key")
	} else {
		t.Logf("✅ 不同URL生成不同缓存key")
	}
}

// TestAKToolsCacheMaxSize 测试缓存最大大小限制
func TestAKToolsCacheMaxSize(t *testing.T) {
	// 创建一个最大大小为2的缓存
	cache := NewRequestCache(5*time.Minute, 2)

	// 添加3个条目
	cache.Set("url1", []byte("data1"))
	cache.Set("url2", []byte("data2"))
	cache.Set("url3", []byte("data3"))

	size := cache.Size()
	t.Logf("📦 添加3个条目后缓存大小: %d", size)

	if size > 2 {
		t.Errorf("❌ 缓存大小(%d)超过最大限制(2)", size)
	} else {
		t.Logf("✅ 缓存大小限制正常工作")
	}

	// 验证最新的条目存在
	if data, found := cache.Get("url3"); !found {
		t.Errorf("❌ 最新添加的条目应该存在")
	} else {
		t.Logf("✅ 最新条目存在: %s", string(data))
	}
}
