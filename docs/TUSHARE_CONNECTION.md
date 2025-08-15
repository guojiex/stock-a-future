# Tushare连接测试功能说明

## 概述

Stock-A-Future API服务器现在在启动时会自动测试Tushare API连接，确保服务可用性。如果连接测试失败，服务器将无法启动，避免运行时出现API调用错误。

## 功能特性

### 1. 启动时连接测试
- 服务器启动时自动测试Tushare API连接
- 使用轻量级的`stock_basic` API进行测试
- 测试股票：000001.SZ（平安银行）
- 超时时间：10秒

### 2. 健康检查集成
- 健康检查端点`/api/v1/health`包含Tushare连接状态
- 实时监控Tushare API可用性
- 返回详细的连接状态信息

### 3. 错误处理
- 连接失败时服务器无法启动
- 详细的错误日志记录
- 清晰的错误提示信息

## 使用方法

### 启动服务器

```bash
# 使用PowerShell脚本启动
.\start-server.ps1

# 或直接使用Go命令
go run cmd/server/main.go
```

### 启动流程

1. **环境检查**：验证Go环境和依赖
2. **配置验证**：检查Tushare Token配置
3. **连接测试**：测试Tushare API连接
4. **服务启动**：启动HTTP服务器

### 启动日志示例

```
启动Stock-A-Future API服务器...
服务器地址: localhost:8080
Tushare API: http://api.tushare.pro
正在测试Tushare API连接...
Tushare连接测试成功 - 获取到股票数据: [000001.SZ 000001 平安银行]
✓ Tushare API连接测试成功
日线数据缓存已启用
服务器正在监听 localhost:8080
```

## 健康检查

### 健康状态端点

```
GET /api/v1/health
```

### 响应格式

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "services": {
    "tushare_api": {
      "status": "healthy",
      "url": "http://api.tushare.pro"
    },
    "cache": {
      "enabled": true
    },
    "favorites": {
      "status": "healthy"
    }
  },
  "uptime": "2h30m15s"
}
```

### 状态说明

- `healthy`: 所有服务正常
- `degraded`: 部分服务异常（如Tushare连接失败）
- `unhealthy`: 服务不可用

## 故障排除

### 常见问题

#### 1. Tushare Token无效

```
Tushare API连接测试失败: tushare API错误: token无效 (代码: 10001)
```

**解决方案**：
- 检查`.env`文件中的`TUSHARE_TOKEN`配置
- 验证Token是否有效且未过期
- 确认Token有足够的API调用权限

#### 2. 网络连接问题

```
Tushare连接测试失败: 发送HTTP请求失败: dial tcp: lookup api.tushare.pro: no such host
```

**解决方案**：
- 检查网络连接
- 验证DNS解析
- 检查防火墙设置

#### 3. API限流

```
Tushare API错误: 请求过于频繁，请稍后再试 (代码: 10002)
```

**解决方案**：
- 等待一段时间后重试
- 检查API调用频率限制
- 考虑升级Tushare账户等级

### 调试模式

启用详细日志记录：

```bash
# 设置日志级别
export LOG_LEVEL=debug

# 启动服务器
go run cmd/server/main.go
```

## 配置选项

### 环境变量

```bash
# Tushare配置
TUSHARE_TOKEN=your_token_here
TUSHARE_BASE_URL=http://api.tushare.pro

# 服务器配置
SERVER_HOST=localhost
SERVER_PORT=8080

# 日志配置
LOG_LEVEL=info
```

### .env文件示例

```env
TUSHARE_TOKEN=your_tushare_token_here
TUSHARE_BASE_URL=http://api.tushare.pro
SERVER_HOST=localhost
SERVER_PORT=8080
LOG_LEVEL=info
CACHE_ENABLED=true
CACHE_DEFAULT_TTL=1h
CACHE_MAX_AGE=24h
CACHE_CLEANUP_INTERVAL=10m
```

## 性能考虑

### 连接测试开销

- 测试请求使用轻量级API
- 超时时间设置为10秒
- 仅在启动时执行一次
- 不影响运行时性能

### 健康检查开销

- 健康检查包含实时连接测试
- 建议设置合理的检查频率
- 可考虑缓存连接状态（短期）

## 最佳实践

### 1. 监控和告警

- 定期检查健康状态端点
- 设置Tushare连接失败的告警
- 监控API调用成功率

### 2. 容错处理

- 实现重试机制
- 考虑备用数据源
- 优雅降级策略

### 3. 安全考虑

- 保护Tushare Token
- 限制API访问频率
- 监控异常API调用

## 更新日志

- **v1.0.0**: 初始版本，基础连接测试
- **v1.1.0**: 集成健康检查，改进错误处理
- **v1.2.0**: 优化测试性能，添加详细日志

## 技术支持

如果遇到问题，请：

1. 检查本文档的故障排除部分
2. 查看服务器启动日志
3. 验证Tushare API配置
4. 提交Issue到项目仓库
