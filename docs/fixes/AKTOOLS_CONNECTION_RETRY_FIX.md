# AKTools连接重试机制改进

## 问题描述

之前的实现中，如果服务器启动时AKTools API连接测试失败，程序会立即退出（Fatal error），导致：

1. **启动失败**：即使只是暂时的网络问题，也无法启动服务
2. **无重试机制**：不给AKTools服务恢复的机会
3. **影响开发**：本地开发时如果忘记启动AKTools，整个服务无法运行

## 改进方案

### 1. 添加重试机制

在AKTools连接测试时增加3次重试，每次间隔2秒：

```go
// 启动时测试AKTools连接（带重试机制）
logger.Info("正在测试AKTools API连接...")
maxRetries := 3
retryDelay := 2 * time.Second
var lastErr error

for i := 0; i < maxRetries; i++ {
    if err := dataSourceClient.TestConnection(); err != nil {
        lastErr = err
        logger.Warn("AKTools API连接测试失败，正在重试...",
            logger.Int("attempt", i+1),
            logger.Int("max_retries", maxRetries),
            logger.ErrorField(err),
        )
        if i < maxRetries-1 {
            time.Sleep(retryDelay)
        }
    } else {
        logger.Info("✓ AKTools API连接测试成功")
        lastErr = nil
        break
    }
}
```

### 2. 降级为警告而非Fatal

如果3次重试都失败，**不再退出程序**，而是：
- 记录警告日志
- 继续启动服务
- 提示用户检查AKTools服务

```go
// 如果所有重试都失败，记录警告但不退出程序
if lastErr != nil {
    logger.Warn("⚠️ AKTools API连接失败，服务将继续启动但数据获取功能可能不可用",
        logger.ErrorField(lastErr),
        logger.String("aktools_url", cfg.AKToolsBaseURL),
    )
    logger.Warn("请检查AKTools服务是否正常运行")
}
```

## 改进效果

### 启动日志示例

#### 场景1：AKTools正常（第一次连接成功）

```log
INFO 正在测试AKTools API连接...
INFO ✓ AKTools API连接测试成功
INFO 日线数据缓存已启用
INFO 服务器启动成功，监听端口: 8080
```

#### 场景2：AKTools暂时不可用（重试后成功）

```log
INFO 正在测试AKTools API连接...
WARN AKTools API连接测试失败，正在重试... {"attempt": 1, "max_retries": 3, "error": "connection refused"}
WARN AKTools API连接测试失败，正在重试... {"attempt": 2, "max_retries": 3, "error": "connection refused"}
INFO ✓ AKTools API连接测试成功
INFO 服务器启动成功，监听端口: 8080
```

#### 场景3：AKTools持续不可用（服务仍然启动）

```log
INFO 正在测试AKTools API连接...
WARN AKTools API连接测试失败，正在重试... {"attempt": 1, "max_retries": 3, "error": "connection refused"}
WARN AKTools API连接测试失败，正在重试... {"attempt": 2, "max_retries": 3, "error": "connection refused"}
WARN AKTools API连接测试失败，正在重试... {"attempt": 3, "max_retries": 3, "error": "connection refused"}
WARN ⚠️ AKTools API连接失败，服务将继续启动但数据获取功能可能不可用 {"error": "...", "aktools_url": "http://127.0.0.1:8080"}
WARN 请检查AKTools服务是否正常运行
INFO 服务器启动成功，监听端口: 8080
```

## 优势分析

### 1. 增强容错性

- ✅ **网络波动**：临时网络问题不会导致启动失败
- ✅ **服务重启**：AKTools服务重启期间，主服务可以等待
- ✅ **开发友好**：本地开发时更容错

### 2. 更好的用户体验

- ✅ **明确的进度**：显示重试次数和进度
- ✅ **清晰的提示**：告诉用户发生了什么，需要做什么
- ✅ **服务可用**：即使AKTools不可用，其他功能仍可使用（如参数优化、回测管理等）

### 3. 渐进式降级

```
AKTools可用 → 所有功能正常
    ↓ 连接失败
重试3次 (共6秒)
    ↓ 仍然失败
记录警告 + 服务继续运行
    ↓
部分功能可用（不依赖实时数据的功能）
```

## 受影响的功能

### AKTools不可用时：

**✅ 仍然可用的功能：**
- 参数优化（如果使用缓存数据）
- 回测管理和查看历史结果
- 策略管理
- 收藏股票管理
- 查看已缓存的股票数据

**❌ 不可用的功能：**
- 获取实时股票数据
- 获取最新日线数据
- 获取基本面数据
- 新的回测（需要获取新数据）

### 运行时处理

如果AKTools在运行时恢复：
- 后续的数据请求会自动成功
- 无需重启服务
- 用户体验无缝恢复

## 配置说明

### 调整重试参数

如果需要修改重试策略，可以在`cmd/server/main.go`中调整：

```go
maxRetries := 3           // 重试次数（默认3次）
retryDelay := 2 * time.Second  // 重试间隔（默认2秒）
```

### 推荐配置

根据不同场景：

1. **生产环境**：
   - `maxRetries`: 5次
   - `retryDelay`: 5秒
   - 原因：生产环境更稳定，但给服务更多恢复时间

2. **开发环境**：
   - `maxRetries`: 2次
   - `retryDelay`: 1秒
   - 原因：快速失败，方便调试

3. **CI/CD环境**：
   - `maxRetries`: 1次
   - `retryDelay`: 1秒
   - 原因：快速验证，不需要等待

## 未来改进

### 1. 健康检查端点

添加`/api/v1/health`端点，返回各个依赖服务的状态：

```json
{
  "status": "healthy",
  "services": {
    "aktools": {
      "status": "unhealthy",
      "last_check": "2025-10-28T23:34:26Z",
      "error": "connection refused"
    },
    "database": {
      "status": "healthy"
    },
    "cache": {
      "status": "healthy"
    }
  }
}
```

### 2. 自动重连机制

在服务运行期间，定期检查AKTools连接状态：
- 每5分钟自动测试一次
- 连接恢复时自动记录日志
- 提供Prometheus metrics

### 3. 配置化重试策略

通过配置文件控制重试行为：

```yaml
aktools:
  base_url: "http://127.0.0.1:8080"
  connection_test:
    enabled: true
    max_retries: 3
    retry_delay: 2s
    fail_behavior: "warn"  # warn | fatal
```

## 测试方法

### 1. 测试重试机制

```bash
# 确保AKTools服务未启动
# 启动主服务
./bin/server.exe

# 观察日志，应该看到3次重试
# 最后服务仍然成功启动
```

### 2. 测试正常连接

```bash
# 启动AKTools服务
cd aktools && python main.py

# 启动主服务
cd .. && ./bin/server.exe

# 应该看到第一次就连接成功
```

### 3. 测试恢复场景

```bash
# 1. AKTools未启动时启动主服务（会看到重试和警告）
./bin/server.exe

# 2. 在另一个终端启动AKTools
cd aktools && python main.py

# 3. 主服务继续正常提供服务
# 4. 尝试获取股票数据，应该能成功
curl http://localhost:8080/api/v1/stocks/000001.SZ/daily
```

## 相关文件

- `cmd/server/main.go` - 服务器启动逻辑（第108-139行）
- `internal/client/aktools_client.go` - AKTools客户端实现
- `config/config.go` - 配置管理

---

**修复日期**: 2025-10-28  
**问题类型**: 启动失败 + 缺少重试机制  
**修复状态**: ✅ 已完成  
**影响范围**: 服务启动流程

