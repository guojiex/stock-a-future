# 异步信号计算服务

## 概述

异步信号计算服务是一个后台运行的服务，用于自动计算和存储股票的买卖信号。该服务具有以下特点：

- **T+1数据特性**：每天只计算一次，避免重复计算
- **异步执行**：不阻塞主服务，在后台并发处理
- **通用设计**：支持收藏股票自动计算，也支持手动触发计算
- **持久化存储**：信号数据存储在数据库中，可供查询和分析

## 功能特性

### 1. 自动计算
- 服务启动时自动计算收藏列表中所有股票的信号
- 每天定时执行一次（避免T+1数据重复计算）
- 并发处理，提高计算效率

### 2. 手动触发
- 支持单个股票信号计算
- 支持批量股票信号计算
- 支持强制重新计算（忽略今日已计算标记）

### 3. 信号查询
- 根据股票代码和日期查询信号
- 获取最新信号列表
- 支持分页查询

## 数据库结构

### stock_signals 表
```sql
CREATE TABLE IF NOT EXISTS stock_signals (
    id TEXT PRIMARY KEY,
    ts_code TEXT NOT NULL,              -- 股票代码
    name TEXT NOT NULL,                 -- 股票名称
    trade_date TEXT NOT NULL,           -- 信号基于的交易日期
    signal_date TEXT NOT NULL,          -- 信号计算日期
    signal_type TEXT NOT NULL,          -- 信号类型: BUY, SELL, HOLD
    signal_strength TEXT NOT NULL,      -- 信号强度: STRONG, MEDIUM, WEAK
    confidence REAL NOT NULL,           -- 置信度 0-1
    patterns TEXT,                      -- 识别到的图形模式(JSON格式)
    technical_indicators TEXT,          -- 技术指标数据(JSON格式)
    predictions TEXT,                   -- 预测数据(JSON格式)
    description TEXT,                   -- 信号描述
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE(ts_code, trade_date)         -- 每个股票每天只有一个信号记录
);
```

## API 接口

### 1. 计算单个股票信号
```http
POST /api/v1/signals/calculate
Content-Type: application/json

{
    "ts_code": "000001.SZ",
    "name": "平安银行",
    "force": false
}
```

### 2. 批量计算信号
```http
POST /api/v1/signals/batch
Content-Type: application/json

{
    "ts_codes": ["000001.SZ", "000002.SZ"],
    "force": false
}
```

### 3. 获取股票信号
```http
GET /api/v1/signals/{code}?signal_date=20240101
```

### 4. 获取最新信号列表
```http
GET /api/v1/signals?limit=20
```

## 信号计算逻辑

### 1. 数据收集
- 获取最近30天的股票日线数据
- 识别图形模式（蜡烛图模式、量价关系）
- 计算技术指标（移动平均线、MACD、RSI等）

### 2. 信号生成
- **买入信号**：多个买入模式汇聚、价格突破均线等
- **卖出信号**：多个卖出模式汇聚、价格跌破支撑等
- **观望信号**：信号不明确或相互矛盾时

### 3. 置信度计算
- 基于识别到的模式数量和质量
- 技术指标的一致性
- 历史准确率（未来可扩展）

## 使用示例

### 启动服务
服务会在主服务启动时自动启动：
```bash
go run cmd/server/main.go
```

### 测试信号计算
使用提供的测试脚本：
```powershell
.\test-signal-service.ps1
```

### 手动计算信号
```bash
curl -X POST http://localhost:8080/api/v1/signals/calculate \
  -H "Content-Type: application/json" \
  -d '{"ts_code":"000001.SZ","name":"平安银行","force":true}'
```

### 查询信号
```bash
curl http://localhost:8080/api/v1/signals/000001.SZ?signal_date=20240101
```

## 配置说明

### 并发控制
- 默认并发数：5个goroutine
- 可在 `SignalService.calculateFavoriteStocksSignals()` 中调整

### 计算周期
- 默认：每天执行一次
- 可在 `SignalService.runSignalCalculation()` 中调整定时器间隔

### 数据范围
- 默认：分析最近30天数据
- 可在 `SignalService.calculateStockSignal()` 中调整日期范围

## 扩展建议

1. **机器学习集成**：引入更复杂的预测模型
2. **回测功能**：验证信号的历史准确性
3. **实时计算**：支持实时数据流处理
4. **通知机制**：重要信号的邮件/消息推送
5. **性能优化**：缓存计算结果、增量更新

## 故障排查

### 常见问题
1. **信号计算失败**：检查股票数据是否可获取
2. **数据库错误**：确保数据库表已正确创建
3. **并发问题**：检查数据库连接池配置

### 日志监控
- 服务启动日志：确认信号服务已启动
- 计算进度日志：监控批量计算进展
- 错误日志：及时发现和处理异常情况
