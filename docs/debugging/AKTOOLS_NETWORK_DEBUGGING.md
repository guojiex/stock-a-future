# AKTools 网络调试指南

## 📋 概述

本文档描述了如何使用Go单元测试进行AKTools API的网络调用调试，这种方法比使用curl等命令行工具更具跨平台兼容性。

## 🎯 调试原则

### 为什么使用Go单元测试而不是curl？

1. **跨平台兼容性**: Go测试在Windows/Linux/macOS上表现一致
2. **集成测试**: 直接使用项目中的客户端代码，更接近实际使用场景
3. **详细输出**: 可以输出结构化的调试信息
4. **自动化**: 可以集成到CI/CD流程中
5. **参数化测试**: 轻松测试多种参数组合

### 默认配置假设

- **AKTools服务端口**: `http://127.0.0.1:8080`
- **超时时间**: 30秒（某些API响应较慢）
- **测试数据**: 使用健民集团(600976)、平安银行(000001)等作为测试股票

## 🔧 调试模板使用

### 快速开始

```bash
# 运行完整的调试测试套件
go test -v ./internal/client -run TestAKToolsDebugTemplate

# 发现可用的API接口
go test -v ./internal/client -run TestAKToolsAPIDiscovery

# 测试不同股票代码格式
go test -v ./internal/client -run TestStockCodeFormatsDebug
```

### 调试模板文件

调试模板位于 `internal/client/aktools_debug_template_test.go`，包含以下测试函数：

1. **TestAKToolsDebugTemplate**: 基础连接和API文档测试
2. **TestAKToolsAPIDiscovery**: 发现可用的API接口
3. **TestStockCodeFormatsDebug**: 测试不同股票代码格式
4. **TestRawAPIResponse**: 分析原始API响应格式
5. **TestMultipleStocks**: 测试多个股票代码
6. **TestAKToolsPerformance**: API性能测试

## 📊 常见问题排查

### 404错误：接口不存在

```go
// 错误的接口名称
apiURL := "http://127.0.0.1:8080/api/public/stock_zh_a_info"

// 正确的接口名称
apiURL := "http://127.0.0.1:8080/api/public/stock_individual_info_em"
```

### 超时错误：API响应慢

```go
// 设置更长的超时时间
httpClient := &http.Client{
    Timeout: 30 * time.Second, // 增加到30秒
}
```

### 数据格式错误：解析失败

```go
// stock_individual_info_em 返回key-value对数组
var rawData []map[string]interface{}
json.Unmarshal(body, &rawData)

// 需要转换为map格式
stockData := make(map[string]interface{})
for _, item := range rawData {
    if itemKey, ok := item["item"].(string); ok {
        if itemValue, exists := item["value"]; exists {
            stockData[itemKey] = itemValue
        }
    }
}
```

## 🚀 实战示例

### 调试新API接口

当需要调试一个新的AKTools API接口时：

1. **复制调试模板**:
   ```go
   func TestNewAPIDebug(t *testing.T) {
       client := NewAKToolsClient("http://127.0.0.1:8080")
       
       // 测试连接
       err := client.TestConnection()
       if err != nil {
           t.Fatalf("连接失败: %v", err)
       }
       
       // 测试新API
       result, err := client.NewAPIMethod("test_param")
       if err != nil {
           t.Logf("❌ API调用失败: %v", err)
       } else {
           t.Logf("✅ API调用成功: %+v", result)
       }
   }
   ```

2. **运行调试测试**:
   ```bash
   go test -v ./internal/client -run TestNewAPIDebug
   ```

3. **分析输出结果**，根据错误信息调整参数或接口名称

### 验证API修复

当修复了API问题后：

1. **运行相关测试**:
   ```bash
   go test -v ./internal/client -run TestStockCodeFormatsDebug
   ```

2. **验证多个股票代码**:
   ```bash
   go test -v ./internal/client -run TestMultipleStocks
   ```

3. **检查性能**:
   ```bash
   go test -v ./internal/client -run TestAKToolsPerformance
   ```

## 📝 最佳实践

1. **总是先测试连接**: 确保AKTools服务正常运行
2. **使用详细日志**: 利用`t.Logf()`输出调试信息
3. **测试多种格式**: 验证不同的股票代码格式
4. **检查原始响应**: 分析API返回的原始数据格式
5. **性能监控**: 关注API响应时间，识别性能瓶颈

## 🔗 相关文档

- [项目开发规范](.cursor/rules/my-custom-rule.mdc)
- [AKTools集成文档](../integration/AKTOOLS_INTEGRATION.md)
- [基本面数据API](../integration/AKTOOLS_FUNDAMENTAL_API.md)

---

**提示**: 这种调试方法不仅适用于AKTools，也适用于其他HTTP API的调试。关键是使用Go的标准库构建可重复、可维护的测试用例。
