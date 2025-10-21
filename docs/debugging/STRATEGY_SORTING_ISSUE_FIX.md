# 策略管理页面排序不稳定问题修复

## 问题描述

在策略管理页面点击"刷新"按钮时，策略列表的排序顺序会发生随机变化，导致用户体验不佳。

## 问题分析

### 根本原因

问题由两个因素共同导致：

1. **Go map 遍历顺序的不确定性**
   - 在 `internal/service/strategy.go` 的 `GetStrategiesList` 方法中，策略存储在 `map[string]*models.Strategy` 中
   - Go 语言的 map 遍历顺序是**非确定性的**（为了安全考虑，Go 故意随机化 map 遍历顺序）
   - 每次从 map 中提取策略时，顺序可能不同

2. **默认策略的创建时间相同**
   - 在 `internal/models/strategy.go` 中，所有默认策略的 `CreatedAt` 都设置为 `time.Now()`
   - 多个策略具有相同的状态（Active）和相同的创建时间
   - 排序函数无法区分这些策略的顺序

### 问题代码

```go
// internal/service/strategy.go (修复前)
for _, strategy := range s.strategies {  // map 遍历顺序不确定
    if s.matchesFilter(strategy, req) {
        results = append(results, *strategy)
    }
}

sort.Slice(results, func(i, j int) bool {
    // ...
    // 状态相同时，按创建时间降序排序（新的在前）
    return results[i].CreatedAt.After(results[j].CreatedAt)
    // 问题：如果创建时间也相同，排序不稳定
})
```

```go
// internal/models/strategy.go (修复前)
var DefaultStrategies = []Strategy{
    {
        ID:        "macd_strategy",
        Status:    StrategyStatusActive,
        CreatedAt: time.Now(),  // 所有策略同时创建
        // ...
    },
    {
        ID:        "ma_crossover",
        Status:    StrategyStatusActive,
        CreatedAt: time.Now(),  // 创建时间相同
        // ...
    },
    // ...
}
```

## 解决方案

### 1. 添加 ID 作为第三排序键

在 `internal/service/strategy.go` 中修改排序逻辑，添加 ID 字段作为最终的排序依据：

```go
// 稳定排序：优先按状态排序（active > testing > inactive），
// 然后按创建时间排序（新的在前），最后按ID排序
sort.Slice(results, func(i, j int) bool {
    // 首先按状态排序
    statusPriority := map[models.StrategyStatus]int{
        models.StrategyStatusActive:   1,
        models.StrategyStatusTesting:  2,
        models.StrategyStatusInactive: 3,
    }

    statusI := statusPriority[results[i].Status]
    statusJ := statusPriority[results[j].Status]

    if statusI != statusJ {
        return statusI < statusJ
    }

    // 状态相同时，按创建时间降序排序（新的在前）
    if !results[i].CreatedAt.Equal(results[j].CreatedAt) {
        return results[i].CreatedAt.After(results[j].CreatedAt)
    }

    // 创建时间相同时，按ID字典序排序（保证稳定排序）
    return results[i].ID < results[j].ID
})
```

### 2. 为默认策略设置不同的创建时间

在 `internal/models/strategy.go` 中，为每个默认策略设置不同的创建时间：

```go
// DefaultStrategies 默认策略配置
// 注意：使用固定的创建时间顺序，确保排序的稳定性
var DefaultStrategies = []Strategy{
    {
        ID:        "macd_strategy",
        Status:    StrategyStatusActive,
        CreatedAt: time.Now().Add(-3 * time.Hour), // 最早创建
        // ...
    },
    {
        ID:        "ma_crossover",
        Status:    StrategyStatusActive,
        CreatedAt: time.Now().Add(-2 * time.Hour),
        // ...
    },
    {
        ID:        "rsi_strategy",
        Status:    StrategyStatusInactive,
        CreatedAt: time.Now().Add(-1 * time.Hour),
        // ...
    },
    {
        ID:        "bollinger_strategy",
        Status:    StrategyStatusActive,
        CreatedAt: time.Now(), // 最新创建
        // ...
    },
}
```

## 排序规则

修复后的策略列表排序规则（优先级从高到低）：

1. **按状态排序**: Active > Testing > Inactive
2. **按创建时间排序**: 新的在前（降序）
3. **按ID字典序排序**: 确保稳定性

## 测试验证

创建了完整的测试套件 `internal/service/strategy_sort_test.go`：

### 1. 排序稳定性测试

```go
func TestStrategyListSorting(t *testing.T)
```

- 多次获取策略列表（10次）
- 验证每次获取的顺序是否一致
- 测试结果：✅ 所有10次获取的顺序完全一致

### 2. 排序规则测试

```go
func TestStrategyListSortingRules(t *testing.T)
```

- 创建不同状态和创建时间的测试策略
- 验证排序规则是否正确应用
- 测试结果：✅ 排序规则符合预期

### 3. 相同时间戳测试

```go
func TestStrategyListSortingWithSameTimestamp(t *testing.T)
```

- 创建具有相同状态和创建时间的策略
- 验证ID字典序排序是否生效
- 测试结果：✅ 按ID字典序稳定排序

## 测试运行结果

```bash
$ go test -v ./internal/service -run TestStrategyListSorting

=== RUN   TestStrategyListSorting
    strategy_sort_test.go:41: 第 1 次获取策略列表顺序: [bollinger_strategy ma_crossover macd_strategy rsi_strategy]
    strategy_sort_test.go:60: 第 2 次获取策略列表顺序: [bollinger_strategy ma_crossover macd_strategy rsi_strategy] (✓ 与第1次一致)
    strategy_sort_test.go:60: 第 3 次获取策略列表顺序: [bollinger_strategy ma_crossover macd_strategy rsi_strategy] (✓ 与第1次一致)
    ...
--- PASS: TestStrategyListSorting (0.03s)
=== RUN   TestStrategyListSortingRules
    strategy_sort_test.go:126: 策略排序结果:
    strategy_sort_test.go:128: 1. bollinger_strategy (状态: active, 创建时间: 20:31:05)
    strategy_sort_test.go:128: 2. test_active_2 (状态: active, 创建时间: 19:31:05)
    strategy_sort_test.go:128: 3. test_active_1 (状态: active, 创建时间: 18:31:05)
    ...
--- PASS: TestStrategyListSortingRules (0.00s)
=== RUN   TestStrategyListSortingWithSameTimestamp
--- PASS: TestStrategyListSortingWithSameTimestamp (0.00s)
PASS
ok      stock-a-future/internal/service 0.585s
```

## 预期效果

修复后，用户在策略管理页面点击刷新按钮时：

1. ✅ 策略列表的顺序将**保持稳定**，不会随机变化
2. ✅ 策略按照**状态优先级**排序（活跃策略在前）
3. ✅ 相同状态的策略按**创建时间降序**排序（新策略在前）
4. ✅ 创建时间相同的策略按**ID字典序**排序（可预测的顺序）

## 技术要点

### Go Map 遍历的特性

Go 语言的 map 遍历顺序是**故意随机化**的：

```go
// 这段代码每次运行顺序可能不同
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k, v := range m {
    fmt.Println(k, v)  // 顺序不确定
}
```

**原因**：
- 防止程序依赖 map 的遍历顺序（这是未定义的行为）
- 鼓励开发者在需要确定顺序时显式排序

**解决方法**：
1. 提取到切片后进行排序
2. 使用多级排序键确保稳定性
3. 添加唯一标识符（如ID）作为最终排序键

### 稳定排序 vs 不稳定排序

- **不稳定排序**：相等元素的相对顺序可能改变
- **稳定排序**：相等元素的相对顺序保持不变

Go 的 `sort.Slice` 是**不稳定排序**，因此需要通过多级排序键来实现稳定性：

```go
sort.Slice(items, func(i, j int) bool {
    // 第一级排序键
    if items[i].Priority != items[j].Priority {
        return items[i].Priority > items[j].Priority
    }
    // 第二级排序键
    if items[i].CreatedAt != items[j].CreatedAt {
        return items[i].CreatedAt.After(items[j].CreatedAt)
    }
    // 第三级排序键（唯一标识符，确保稳定性）
    return items[i].ID < items[j].ID
})
```

## 相关文件

修改的文件：
- `internal/service/strategy.go` - 添加ID排序键
- `internal/models/strategy.go` - 设置不同的创建时间

新增的文件：
- `internal/service/strategy_sort_test.go` - 排序稳定性测试

## 影响范围

- ✅ 前端策略管理页面：排序顺序稳定
- ✅ API `/api/v1/strategies`：返回顺序稳定
- ✅ 无破坏性变更：完全向后兼容

## 总结

这个修复通过两个简单但关键的改进，解决了由 Go map 遍历不确定性和排序键不完整导致的排序不稳定问题。修复后的代码更加健壮，用户体验得到明显改善。

**关键经验**：
1. 在 Go 中使用 map 时，不要依赖遍历顺序
2. 实现排序时，确保排序键的唯一性和完整性
3. 对于列表展示，始终使用稳定的排序规则
4. 编写测试来验证排序的稳定性

