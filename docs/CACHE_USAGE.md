# 日线数据缓存系统使用指南

## 概述

Stock-A-Future API 现在支持日线数据缓存功能，可以显著提高API响应速度并减少对Tushare API的请求次数。

## 特性

### 🚀 核心功能
- **高效内存缓存**：使用并发安全的 `sync.Map` 存储日线数据
- **灵活过期时间**：支持为不同数据设置不同的TTL（生存时间）
- **自动清理机制**：定期清理过期缓存，避免内存泄漏
- **智能缓存键**：基于股票代码、开始日期和结束日期生成唯一缓存键
- **统计监控**：提供详细的缓存命中率和统计信息

### 📊 性能优化
- **首次请求**：从Tushare API获取数据并缓存
- **后续请求**：直接从内存缓存返回，响应时间从秒级降至毫秒级
- **并发安全**：支持多个请求同时访问缓存

## 配置选项

### 环境变量配置

```bash
# 是否启用缓存（默认：true）
CACHE_ENABLED=true

# 默认缓存过期时间（默认：1h）
CACHE_DEFAULT_TTL=1h

# 最大缓存时间（默认：24h）
CACHE_MAX_AGE=24h

# 清理过期缓存的间隔（默认：10m）
CACHE_CLEANUP_INTERVAL=10m
```

### 时间格式说明
- `s` - 秒
- `m` - 分钟
- `h` - 小时
- 示例：`30m`（30分钟），`2h30m`（2小时30分钟）

## API端点

### 缓存统计信息
```http
GET /api/v1/cache/stats
```

**响应示例：**
```json
{
  "success": true,
  "data": {
    "stats": {
      "hits": 45,
      "misses": 12,
      "entries": 8,
      "evictions": 3,
      "last_cleanup": "2024-01-15T10:30:00Z"
    },
    "hit_rate": "78.95%",
    "config": {
      "default_ttl": "1h0m0s",
      "max_cache_age": "24h0m0s"
    }
  }
}
```

### 清空缓存
```http
DELETE /api/v1/cache
```

**响应示例：**
```json
{
  "success": true,
  "data": {
    "message": "缓存已清空",
    "timestamp": "2024-01-15T10:35:00Z"
  }
}
```

## 缓存策略

### 缓存键生成
缓存键基于以下参数生成：
- 股票代码（如：000001.SZ）
- 开始日期（如：20240101）
- 结束日期（如：20240131）

### 缓存生命周期
1. **数据获取**：首次请求时从Tushare API获取数据
2. **缓存存储**：将数据存储到内存缓存中，设置过期时间
3. **缓存命中**：后续相同请求直接从缓存返回数据
4. **自动过期**：达到TTL时间后缓存自动失效
5. **定期清理**：后台定期清理过期的缓存条目

### 适用场景
- **相同股票相同时间段**的重复查询
- **技术指标计算**：需要获取历史数据进行计算
- **买卖点预测**：需要大量历史数据支持
- **Web界面**：用户频繁查看相同股票数据

## 使用示例

### 1. 查看缓存状态
```bash
curl http://localhost:8080/api/v1/cache/stats
```

### 2. 获取日线数据（首次请求 - 缓存未命中）
```bash
curl "http://localhost:8080/api/v1/stocks/000001.SZ/daily?start_date=20240101&end_date=20240131"
```

### 3. 再次获取相同数据（缓存命中）
```bash
curl "http://localhost:8080/api/v1/stocks/000001.SZ/daily?start_date=20240101&end_date=20240131"
```

### 4. 清空缓存
```bash
curl -X DELETE http://localhost:8080/api/v1/cache
```

## 最佳实践

### 🎯 推荐配置

**开发环境：**
```bash
CACHE_ENABLED=true
CACHE_DEFAULT_TTL=30m
CACHE_MAX_AGE=2h
CACHE_CLEANUP_INTERVAL=5m
```

**生产环境：**
```bash
CACHE_ENABLED=true
CACHE_DEFAULT_TTL=1h
CACHE_MAX_AGE=24h
CACHE_CLEANUP_INTERVAL=10m
```

**高频交易场景：**
```bash
CACHE_ENABLED=true
CACHE_DEFAULT_TTL=10m
CACHE_MAX_AGE=1h
CACHE_CLEANUP_INTERVAL=2m
```

### 💡 优化建议

1. **合理设置TTL**：
   - 日内交易：10-30分钟
   - 技术分析：1-2小时
   - 历史研究：4-24小时

2. **监控缓存性能**：
   - 定期检查命中率（建议 > 60%）
   - 关注缓存条目数量
   - 监控内存使用情况

3. **缓存清理策略**：
   - 市场开盘前清空缓存获取最新数据
   - 系统升级后清空缓存
   - 内存使用过高时手动清理

## 注意事项

### ⚠️ 重要提醒

1. **数据时效性**：缓存的数据可能不是最新的，适用于对实时性要求不高的场景
2. **内存使用**：大量缓存会占用服务器内存，请根据服务器配置合理设置
3. **缓存穿透**：不同的时间范围会产生不同的缓存键，无法复用缓存
4. **服务重启**：服务重启后缓存会清空，需要重新构建

### 🔧 故障排除

**缓存未生效**：
- 检查 `CACHE_ENABLED` 是否为 `true`
- 确认缓存服务是否正常启动
- 查看日志中的缓存相关信息

**内存使用过高**：
- 减少 `CACHE_MAX_AGE` 时间
- 缩短 `CACHE_CLEANUP_INTERVAL` 间隔
- 手动清空缓存

**命中率过低**：
- 增加 `CACHE_DEFAULT_TTL` 时间
- 分析请求模式，优化缓存策略
- 检查是否有大量不同时间范围的请求

## 技术实现

### 架构组件
- **DailyCacheService**：缓存服务核心
- **CacheEntry**：缓存条目结构
- **CacheStats**：统计信息收集
- **自动清理**：后台goroutine定期清理

### 并发安全
- 使用 `sync.Map` 保证并发读写安全
- 使用 `sync.RWMutex` 保护统计信息
- 无锁设计，高并发性能优异

### 内存管理
- 定期清理过期条目
- 限制最大缓存时间
- 支持手动清空缓存
