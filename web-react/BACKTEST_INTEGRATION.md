# 回测功能集成 - Web React版本

## 📋 概述

本次更新实现了从策略管理页面跳转到回测页面的完整流程，参考网页版的实现方式，使用 Redux 管理状态和 React Router 进行页面导航。

## 🎯 实现的功能

### 1. Redux 状态管理 (`backtestSlice`)

创建了专门的 Redux slice 来管理回测相关状态：

**文件**: `web-react/src/store/slices/backtestSlice.ts`

**功能**:
- ✅ 管理选中的策略ID列表
- ✅ 管理回测配置（初始资金、手续费、日期范围等）
- ✅ 提供策略添加/移除/清空操作
- ✅ 提供配置更新操作
- ✅ 提供状态重置功能

**主要 Actions**:
```typescript
- setSelectedStrategies(ids: string[])  // 设置选中的策略
- addStrategy(id: string)               // 添加单个策略
- removeStrategy(id: string)            // 移除策略
- clearStrategies()                     // 清空所有选择
- updateConfig(config: Partial<Config>) // 更新回测配置
- resetBacktest()                       // 重置所有状态
```

### 2. 回测页面 (`BacktestPage`)

创建了基础的回测页面框架：

**文件**: `web-react/src/pages/BacktestPage.tsx`

**当前功能**:
- ✅ 显示选中的策略信息
- ✅ 返回策略列表按钮
- ✅ 开发中的提示和说明
- ✅ 自动读取 Redux 中的选中策略

**计划功能**:
- 🚧 回测参数配置表单
- 🚧 策略多选管理
- 🚧 回测执行和进度条
- 🚧 权益曲线图表展示
- 🚧 交易记录详情表格
- 🚧 策略性能对比分析

### 3. 策略页面集成

更新了策略管理页面，实现完整的回测跳转：

**文件**: `web-react/src/pages/StrategiesPage.tsx`

**实现流程**:
```typescript
handleRunBacktest(strategy) {
  // 1. 设置选中的策略ID到 Redux store
  dispatch(setSelectedStrategies([strategy.id]));
  
  // 2. 导航到回测页面
  navigate('/backtest');
}
```

### 4. 路由配置

添加了回测页面的路由：

**文件**: `web-react/src/App.tsx`

```typescript
<Route path="backtest" element={<BacktestPage />} />
```

## 🔄 数据流

```
用户点击"运行回测"按钮
    ↓
StrategiesPage.handleRunBacktest()
    ↓
dispatch(setSelectedStrategies([strategyId]))
    ↓
Redux Store 更新 backtest.selectedStrategyIds
    ↓
navigate('/backtest')
    ↓
BacktestPage 组件加载
    ↓
useAppSelector 读取 selectedStrategyIds
    ↓
显示选中的策略信息
```

## 📁 新增文件

```
web-react/src/
├── store/
│   └── slices/
│       └── backtestSlice.ts          # 新增: 回测状态管理
├── pages/
│   └── BacktestPage.tsx              # 新增: 回测页面
└── BACKTEST_INTEGRATION.md           # 本文档
```

## 🔧 修改的文件

```
web-react/src/
├── store/
│   └── index.ts                      # 添加 backtestSlice
├── pages/
│   └── StrategiesPage.tsx            # 实现回测跳转逻辑
└── App.tsx                           # 添加回测路由
```

## 🎨 用户体验

### 策略管理页面
1. 用户点击策略卡片上的"运行回测"按钮
2. 系统自动选中该策略
3. 页面跳转到回测页面

### 回测页面
1. 显示"已选择 X 个策略进行回测"的提示
2. 显示选中策略的ID列表
3. 提供返回策略列表的按钮

## 🆚 与网页版对比

| 功能 | 网页版 | React版(当前) | 说明 |
|------|--------|--------------|------|
| 策略选择传递 | ✅ | ✅ | 使用全局变量 vs Redux |
| 页面跳转 | Tab切换 | 路由导航 | 符合各自技术栈 |
| 多策略选择 | ✅ | ✅ | Redux支持多选 |
| 回测配置 | ✅ | 🚧 | 待实现 |
| 回测执行 | ✅ | 🚧 | 待实现 |
| 结果展示 | ✅ | 🚧 | 待实现 |

## 🚀 使用示例

### 从策略页面跳转到回测

```typescript
// 在策略卡片中点击"运行回测"按钮
<Button
  variant="contained"
  startIcon={<AssessmentIcon />}
  onClick={() => handleRunBacktest(strategy)}
>
  运行回测
</Button>

// 自动跳转到 /backtest 并传递策略ID
```

### 在回测页面读取选中的策略

```typescript
import { useAppSelector } from '../hooks/redux';

const BacktestPage = () => {
  const selectedStrategyIds = useAppSelector(
    (state) => state.backtest.selectedStrategyIds
  );
  
  // selectedStrategyIds 包含选中的策略ID
  console.log('选中的策略:', selectedStrategyIds);
};
```

### 手动添加策略到回测

```typescript
import { useDispatch } from 'react-redux';
import { addStrategy } from '../store/slices/backtestSlice';

const dispatch = useDispatch();

// 添加单个策略
dispatch(addStrategy('macd_strategy'));

// 设置多个策略
dispatch(setSelectedStrategies(['macd_strategy', 'ma_crossover']));
```

## 📊 状态结构

```typescript
// Redux State 结构
{
  backtest: {
    selectedStrategyIds: ['strategy-id-1', 'strategy-id-2'],
    config: {
      name: '回测名称',
      startDate: '2024-01-01',
      endDate: '2024-12-31',
      initialCash: 1000000,
      commission: 0.0003,
      symbols: ['000001.SZ', '600000.SH']
    }
  }
}
```

## 🔮 后续开发计划

### Phase 1: 回测配置 (P0)
- [ ] 实现回测参数表单
  - [ ] 回测名称输入
  - [ ] 时间范围选择器
  - [ ] 初始资金输入
  - [ ] 手续费率配置
  - [ ] 股票代码列表输入
- [ ] 表单验证和错误提示
- [ ] 配置保存到 Redux

### Phase 2: 策略管理 (P0)
- [ ] 显示选中策略的详细信息
- [ ] 支持添加更多策略（多选）
- [ ] 策略移除功能
- [ ] 策略参数展示

### Phase 3: 回测执行 (P1)
- [ ] 调用后端回测API
- [ ] 显示回测进度
- [ ] 实时更新进度条
- [ ] 支持取消回测
- [ ] 错误处理和重试

### Phase 4: 结果展示 (P1)
- [ ] 权益曲线图表 (ECharts)
- [ ] 性能指标卡片
- [ ] 交易记录表格
- [ ] 策略对比分析
- [ ] 数据导出功能

### Phase 5: 高级功能 (P2)
- [ ] 回测历史记录
- [ ] 配置模板保存
- [ ] 批量回测
- [ ] 参数优化建议
- [ ] 风险分析报告

## 🐛 已知限制

1. **回测页面功能不完整**: 当前仅为基础框架，核心功能待实现
2. **策略信息展示简单**: 仅显示ID，需要关联策略详细信息
3. **无历史记录**: 暂不支持查看历史回测结果
4. **无参数验证**: 配置表单验证待实现

## 💡 技术亮点

1. **Redux 集中式状态管理**: 清晰的状态流转，易于调试
2. **类型安全**: TypeScript 提供完整的类型检查
3. **模块化设计**: 每个功能独立的 slice 和组件
4. **可扩展架构**: 易于添加新功能和优化

## 📚 参考资料

- 网页版回测模块: `web/static/js/modules/backtest.js`
- 网页版策略模块: `web/static/js/modules/strategies.js`
- Redux Toolkit文档: https://redux-toolkit.js.org/
- React Router文档: https://reactrouter.com/

---

**更新日期**: 2025-10-14
**版本**: v0.1.0 (基础框架)
**状态**: 🚧 开发中

