# 策略排序调试指南

## 概述

已在前端和后端添加了详细的调试日志，用于追踪策略列表排序变化的根本原因。

## 如何使用

### 1. 启动后端服务器

```bash
# Windows
.\bin\server.exe

# Linux/Mac
./bin/server
```

后端会在控制台输出详细的日志信息。

### 2. 启动前端（React Web）

```bash
cd web-react
npm start
```

### 3. 打开浏览器开发者工具

1. 访问策略管理页面
2. 按 `F12` 打开开发者工具
3. 切换到 **Console（控制台）** 标签

### 4. 测试排序稳定性

在策略管理页面多次点击"刷新"按钮（建议至少点击 5-10 次），观察日志输出。

## 日志说明

### 前端日志（浏览器控制台）

每次获取策略列表时，会输出以下信息：

```
========== 第 1 次获取策略列表 ==========
策略顺序: ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]

┌─────────┬──────────────────────┬────────────────┬──────────┬───────────┬─────────────────────────┬─────────────────────────┐
│ (index) │          ID          │      名称      │   状态   │   类型    │        创建时间         │        更新时间         │
├─────────┼──────────────────────┼────────────────┼──────────┼───────────┼─────────────────────────┼─────────────────────────┤
│    0    │ 'bollinger_strategy' │   '布林带策略'  │ 'active' │ 'technical' │ '2025-10-21T17:36:48...' │ '2025-10-21T20:36:48...' │
│    1    │   'ma_crossover'     │  '双均线策略'   │ 'active' │ 'technical' │ '2025-10-21T18:36:48...' │ '2025-10-21T20:36:48...' │
│    2    │   'macd_strategy'    │ 'MACD金叉策略' │ 'active' │ 'technical' │ '2025-10-21T17:36:48...' │ '2025-10-21T20:36:48...' │
│    3    │   'rsi_strategy'     │'RSI超买超卖策略'│'inactive'│ 'technical' │ '2025-10-21T19:36:48...' │ '2025-10-21T20:36:48...' │
└─────────┴──────────────────────┴────────────────┴──────────┴───────────┴─────────────────────────┴─────────────────────────┘

✅ 策略顺序保持一致
```

#### 如果顺序发生变化

会看到警告日志：

```
⚠️ 策略顺序发生变化！
上次顺序: ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]
本次顺序: ["ma_crossover", "bollinger_strategy", "macd_strategy", "rsi_strategy"]
变化详情: ["位置 0: bollinger_strategy -> ma_crossover", "位置 1: ma_crossover -> bollinger_strategy"]

原始API响应: {
  "success": true,
  "data": {
    "items": [...],
    "total": 4,
    "page": 1,
    "size": 20
  }
}
```

### 后端日志（服务器控制台）

#### 服务层日志 (internal/service/strategy.go)

```
INFO    开始获取策略列表    {"total_strategies": 4, "page": 1, "size": 20}
INFO    过滤后策略数量      {"count": 4}
INFO    排序前顺序          {"ids": ["macd_strategy", "bollinger_strategy", "rsi_strategy", "ma_crossover"]}
INFO    排序后顺序          {"ids": ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]}
INFO    返回策略列表        {"total": 4, "returned": 4, "ids": ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]}
```

**重点关注**：
- `排序前顺序`：从 map 中提取的原始顺序（每次可能不同）
- `排序后顺序`：应用排序算法后的顺序（应该稳定）

#### 处理器层日志 (internal/handler/strategy.go)

```
INFO    返回策略列表    {"total": 4, "count": 4, "strategy_ids": ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]}
INFO    策略详情        {"position": 0, "id": "bollinger_strategy", "name": "布林带策略", "status": "active", "type": "technical", "created_at": "2025-10-21T20:36:48..."}
INFO    策略详情        {"position": 1, "id": "ma_crossover", "name": "双均线策略", "status": "active", "type": "technical", "created_at": "2025-10-21T18:36:48..."}
INFO    策略详情        {"position": 2, "id": "macd_strategy", "name": "MACD金叉策略", "status": "active", "type": "technical", "created_at": "2025-10-21T17:36:48..."}
INFO    策略详情        {"position": 3, "id": "rsi_strategy", "name": "RSI超买超卖策略", "status": "inactive", "type": "technical", "created_at": "2025-10-21T19:36:48..."}
```

## 分析指南

### 场景 1：后端排序稳定，前端显示不稳定

**症状**：
- 后端日志显示 `排序后顺序` 每次都相同
- 前端控制台显示 `⚠️ 策略顺序发生变化`

**可能原因**：
- 前端数据处理问题
- Redux/RTK Query 缓存问题
- React 渲染问题

**排查步骤**：
1. 查看前端日志中的 `原始API响应`
2. 检查 `strategiesData.data.items` 的顺序
3. 检查是否有前端排序或过滤逻辑

### 场景 2：后端排序不稳定

**症状**：
- 后端日志显示 `排序后顺序` 每次都不同
- `排序前顺序` 每次都不同

**可能原因**：
- Map 遍历顺序影响
- 排序算法问题
- 策略数据本身有变化

**排查步骤**：
1. 检查 `排序前顺序` 是否每次都不同（正常现象）
2. 检查每个策略的 `status` 和 `created_at` 字段
3. 确认排序逻辑是否正确执行

### 场景 3：策略数据发生变化

**症状**：
- 策略的 `created_at` 或 `updated_at` 字段在变化
- 策略的 `status` 字段在变化

**可能原因**：
- 默认策略使用了 `time.Now()`（已修复）
- 其他服务在修改策略数据
- 并发问题

**排查步骤**：
1. 对比前后两次日志中相同策略的字段值
2. 检查是否有其他操作在修改策略
3. 检查 `DefaultStrategies` 的初始化代码

## 预期的正常日志

### 后端日志应该显示：

```
# 第一次请求
INFO    排序前顺序    {"ids": ["rsi_strategy", "macd_strategy", "ma_crossover", "bollinger_strategy"]}  ← 可能不同
INFO    排序后顺序    {"ids": ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]}  ← 应该稳定

# 第二次请求  
INFO    排序前顺序    {"ids": ["ma_crossover", "bollinger_strategy", "rsi_strategy", "macd_strategy"]}  ← 可能不同
INFO    排序后顺序    {"ids": ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]}  ← 应该相同

# 第三次请求
INFO    排序前顺序    {"ids": ["macd_strategy", "rsi_strategy", "bollinger_strategy", "ma_crossover"]}  ← 可能不同
INFO    排序后顺序    {"ids": ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]}  ← 应该相同
```

**关键点**：
- ✅ `排序前顺序` 可以不同（map 遍历顺序不确定）
- ✅ `排序后顺序` 必须每次都相同（排序算法稳定）

### 前端日志应该显示：

```
========== 第 1 次获取策略列表 ==========
策略顺序: ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]
✅ 策略顺序保持一致

========== 第 2 次获取策略列表 ==========
策略顺序: ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]
✅ 策略顺序保持一致

========== 第 3 次获取策略列表 ==========
策略顺序: ["bollinger_strategy", "ma_crossover", "macd_strategy", "rsi_strategy"]
✅ 策略顺序保持一致
```

## 如何报告问题

如果发现排序仍然不稳定，请提供以下信息：

1. **前端控制台日志**：
   - 完整的日志输出（包括变化前后的顺序）
   - `原始API响应` 的内容

2. **后端服务器日志**：
   - 对应时间段的完整日志
   - 特别是 `排序前顺序` 和 `排序后顺序`

3. **复现步骤**：
   - 点击刷新的次数
   - 是否有其他操作（编辑、删除等）

4. **环境信息**：
   - 浏览器类型和版本
   - 操作系统
   - Go 版本

## 清理日志

测试完成后，如果确认问题已解决，可以：

1. **前端**：移除 `web-react/src/pages/StrategiesPage.tsx` 中的 `useEffect` 日志
2. **后端**：调整日志级别或移除详细日志

但建议在生产环境保留关键日志（不包括 Debug 级别），方便问题追踪。

