# 回测页面 Bug 修复说明

## 问题描述

用户报告了两个问题：
1. **已选择策略但点击回测显示"请选择至少一个策略"**
2. **策略选择器没有默认全选功能**

## 问题分析

### 问题 1：验证逻辑使用了错误的状态变量

**根本原因：**
- 在 `BacktestPage.tsx` 中，`validateConfig()` 函数检查的是 `config.strategy_ids.length`
- 但用户通过策略选择对话框选择策略时，更新的是 Redux state 中的 `selectedStrategyIds`
- `config.strategy_ids` 只在初始化时设置一次，不会自动同步

**代码位置：**
```typescript:156:web-react/src/pages/BacktestPage.tsx
// 修复前
if (config.strategy_ids.length === 0) return '请选择至少一个策略';
if (config.strategy_ids.length > 5) return '最多只能选择5个策略';

// 修复后
if (selectedStrategyIds.length === 0) return '请选择至少一个策略';
if (selectedStrategyIds.length > 5) return '最多只能选择5个策略';
```

### 问题 2：缺少默认全选和全选按钮

**根本原因：**
- 策略选择对话框需要用户手动点击每个策略
- 没有提供快速全选的按钮
- 没有在页面加载时默认选择策略

## 修复方案

### 修复 1：统一使用 Redux state 进行验证

将 `validateConfig()` 函数中的策略数量检查改为使用 `selectedStrategyIds`，确保验证逻辑与实际选择的策略保持一致。

### 修复 2：添加全选功能

**2.1 添加全选按钮**
在策略选择对话框中添加"全选"按钮：

```typescript:532-546:web-react/src/pages/BacktestPage.tsx
<Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
  <Typography variant="body2" color="text.secondary">
    最多可以选择5个策略进行回测
  </Typography>
  <Button
    size="small"
    onClick={() => {
      // 全选所有策略（最多5个）
      const allStrategyIds = strategies.slice(0, 5).map(s => s.id);
      dispatch(setSelectedStrategies(allStrategyIds));
    }}
  >
    全选
  </Button>
</Box>
```

**2.2 添加默认全选逻辑**
在页面加载时，如果没有选择任何策略，自动选择所有策略（最多5个）：

```typescript:125-132:web-react/src/pages/BacktestPage.tsx
// 默认全选策略（仅在首次加载时且没有选择任何策略时）
useEffect(() => {
  if (strategies.length > 0 && selectedStrategyIds.length === 0) {
    // 默认选择所有策略（最多5个）
    const allStrategyIds = strategies.slice(0, 5).map(s => s.id);
    dispatch(setSelectedStrategies(allStrategyIds));
  }
}, [strategies, selectedStrategyIds.length, dispatch]);
```

**2.3 导入必要的 action**
添加 `setSelectedStrategies` 到导入列表：

```typescript:50-56:web-react/src/pages/BacktestPage.tsx
import {
  updateConfig,
  addStrategy,
  removeStrategy,
  clearStrategies,
  setSelectedStrategies,
} from '../store/slices/backtestSlice';
```

## 修改文件

- `web-react/src/pages/BacktestPage.tsx`
  - 修改验证逻辑使用 `selectedStrategyIds`（第 156-158 行）
  - 添加 `setSelectedStrategies` 导入（第 50-56 行）
  - 添加默认全选 useEffect（第 125-132 行）
  - 在策略选择对话框添加全选按钮（第 532-546 行）

## 测试建议

1. **测试验证逻辑修复**
   - 打开回测页面，确认策略已默认全选
   - 点击"开始回测"按钮，确认不再显示"请选择至少一个策略"错误
   - 清空所有策略选择，点击"开始回测"，确认显示正确的错误提示

2. **测试全选功能**
   - 打开回测页面，确认所有策略（最多5个）默认被选中
   - 点击"管理策略"按钮，取消选择某些策略
   - 点击"全选"按钮，确认所有策略被重新选中
   - 如果有超过5个策略，确认只选择前5个

3. **测试边界情况**
   - 当只有1个策略时，确认默认被选中
   - 当有6个或更多策略时，确认只选择前5个
   - 清空策略后重新进入页面，确认策略重新被默认选中

## 用户体验改进

- ✅ 用户不再需要手动选择策略，提高了效率
- ✅ 提供了"全选"快捷按钮，方便快速选择
- ✅ 验证逻辑更加准确，避免了误报错误
- ✅ 保持了最多选择5个策略的限制

## 注意事项

1. **Redux State 管理**：确保 Redux store 正确配置并包含 `backtestSlice`
2. **策略数据格式**：确保后端返回的策略数据包含 `id` 字段
3. **性能考虑**：useEffect 的依赖数组使用 `selectedStrategyIds.length` 而非整个数组，避免不必要的重新渲染

## 相关文件

- `web-react/src/pages/BacktestPage.tsx` - 回测页面主组件
- `web-react/src/store/slices/backtestSlice.ts` - 回测状态管理
- `web-react/src/services/api.ts` - API 服务定义

