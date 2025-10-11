# 策略管理功能实现文档

## 📋 概述

基于Web版本的收藏管理功能,在React Native移动端实现了完整的策略管理系统。用户可以将关注的股票添加为策略,进行分组管理,并查看信号汇总。

## ✅ 已完成功能

### 1. 策略管理主页面 (`BacktestScreen.tsx`)

**位置**: `mobile/src/screens/Backtest/BacktestScreen.tsx`

**功能特性**:
- ✅ 策略列表展示 - 显示所有收藏的股票策略
- ✅ 信号汇总展示 - 查看所有策略的买卖信号
- ✅ Tab切换 - 在列表和信号之间切换
- ✅ 分组管理 - 创建、编辑、删除分组
- ✅ 分组筛选 - 按分组查看策略
- ✅ 下拉刷新 - 支持数据刷新
- ✅ 空状态提示 - 友好的空数据展示
- ✅ 删除策略 - 从分组中移除策略
- ✅ 快速操作 - 查看图表、编辑策略

**UI组件**:
- 顶部Tab导航 (策略列表 / 信号汇总)
- 水平滚动的分组标签
- 卡片式列表展示
- 浮动操作按钮 (FAB) - 添加新策略
- 分组编辑对话框

### 2. 策略编辑页面 (`StrategyEditScreen.tsx`)

**位置**: `mobile/src/screens/Backtest/StrategyEditScreen.tsx`

**功能特性**:
- ✅ 股票搜索 - 智能搜索股票代码和名称
- ✅ 日期选择 - 设置分析时间范围
- ✅ 分组选择 - 将策略添加到指定分组
- ✅ 表单验证 - 确保数据完整性
- ✅ 编辑模式 - 支持编辑现有策略
- ✅ 新建模式 - 添加新的策略

**UI组件**:
- 股票搜索输入框
- 搜索结果下拉列表
- 日期选择器 (DateTimePicker)
- 分组选择chips
- 操作按钮 (取消 / 保存)

### 3. 策略信号详情页面 (`BacktestResultScreen.tsx`)

**位置**: `mobile/src/screens/Backtest/BacktestResultScreen.tsx`

**功能特性**:
- ✅ 预测信号展示 - 显示买卖预测信号
- ✅ 形态识别 - 显示技术形态分析
- ✅ 信号汇总 - 综合评估和统计
- ✅ Tab切换 - 在不同视图间切换
- ✅ 下拉刷新 - 更新信号数据
- ✅ 详细信息 - 置信度、理由、目标价位等

**UI组件**:
- Tab导航 (预测信号 / 形态识别 / 汇总)
- 信号卡片 - 显示买卖信号和置信度
- 形态卡片 - 显示看涨/看跌形态
- 统计卡片 - 显示信号统计数据

## 🔌 后端接口集成

所有功能完全基于现有的后端API,无需额外开发:

### 收藏管理API
```typescript
// 获取收藏列表
useGetFavoritesQuery()

// 添加收藏
useAddFavoriteMutation({
  ts_code: string,
  name: string,
  start_date: string,
  end_date: string,
  group_id?: string
})

// 更新收藏
useUpdateFavoriteMutation({ id, data })

// 删除收藏
useDeleteFavoriteMutation(id)

// 获取收藏信号
useGetFavoritesSignalsQuery()
```

### 分组管理API
```typescript
// 获取分组列表
useGetGroupsQuery()

// 创建分组
useCreateGroupMutation({
  name: string,
  color?: string
})

// 更新分组
useUpdateGroupMutation({ id, data })

// 删除分组
useDeleteGroupMutation(id)
```

### 信号分析API
```typescript
// 获取预测信号
useGetPredictionsQuery(stockCode)

// 形态识别
useRecognizePatternsQuery({ tsCode })

// 获取信号摘要
useGetPatternSummaryQuery(stockCode)
```

### 股票搜索API
```typescript
// 搜索股票
useLazySearchStocksQuery()

// 获取股票基本信息
useLazyGetStockBasicQuery()
```

## 📱 用户界面设计

### 设计原则
- **Material Design** - 遵循Material Design 3规范
- **响应式布局** - 适配不同屏幕尺寸
- **中国股市配色** - 红涨绿跌,符合用户习惯
- **触摸友好** - 按钮大小适合触摸操作
- **视觉反馈** - 清晰的状态提示和加载动画

### 颜色方案
```typescript
// 信号颜色
买入信号: #4caf50 (绿色)
卖出信号: #f44336 (红色)
持有信号: #ff9800 (橙色)

// 分组颜色
蓝色: #1976d2
绿色: #2e7d32
橙色: #ed6c02
红色: #d32f2f
紫色: #9c27b0
青色: #0288d1
```

## 🔄 数据流

```
用户操作
    ↓
React Native组件
    ↓
RTK Query Hook
    ↓
API服务层 (api.ts)
    ↓
Go后端服务
    ↓
SQLite数据库
```

## 📦 依赖项

### 已使用的依赖
- `react-native-paper` - Material Design组件库
- `@react-navigation/native` - 导航管理
- `@reduxjs/toolkit` - 状态管理
- `react-native-vector-icons` - 图标库
- `@react-native-community/datetimepicker` - 日期选择器

### 需要安装的依赖
```bash
# 如果还未安装,需要执行:
npm install @react-native-community/datetimepicker
```

## 🚀 使用方法

### 启动应用
```bash
# 进入mobile目录
cd mobile

# 安装依赖
npm install

# 启动Metro bundler
npm start

# 运行Android
npm run android

# 运行iOS (仅macOS)
npm run ios
```

### 确保后端服务运行
```bash
# 在项目根目录
go run cmd/server/main.go
# 或
./server.exe
```

后端服务应运行在 `http://localhost:8080`

## 📝 代码结构

```
mobile/src/screens/Backtest/
├── BacktestScreen.tsx          # 策略管理主页面
├── StrategyEditScreen.tsx      # 策略编辑页面
└── BacktestResultScreen.tsx    # 策略信号详情页面
```

## 🎯 导航流程

```
BacktestScreen (策略管理)
    ├─→ StrategyEditScreen (添加/编辑策略)
    │       └─→ 返回 BacktestScreen
    │
    └─→ BacktestResultScreen (查看信号详情)
            └─→ 返回 BacktestScreen
```

## ⚠️ 注意事项

1. **数据同步**: 确保后端服务正常运行,API端点可访问
2. **日期格式**: 使用YYYYMMDD格式与后端通信
3. **分组ID**: 默认分组ID为'default',不可删除
4. **导航参数**: BacktestResultScreen的backtestId实际是股票代码
5. **类型安全**: 已定义TypeScript接口确保类型安全

## 🔮 未来优化方向

### 建议添加的功能
1. **拖拽排序** - 使用`react-native-draggable-flatlist`实现策略排序
2. **批量操作** - 支持批量删除、移动策略
3. **导出功能** - 导出策略数据为Excel或CSV
4. **价格预警** - 设置价格提醒通知
5. **策略模板** - 预设常用策略模板
6. **数据可视化** - 添加更多图表展示
7. **离线缓存** - 实现本地缓存提升性能

### 性能优化
- 使用`React.memo`优化列表渲染
- 实现虚拟列表提升长列表性能
- 添加分页加载减少初始加载时间
- 使用`useMemo`和`useCallback`优化计算

## 📚 参考文档

- [React Native Paper文档](https://callstack.github.io/react-native-paper/)
- [React Navigation文档](https://reactnavigation.org/)
- [RTK Query文档](https://redux-toolkit.js.org/rtk-query/overview)
- [Material Design指南](https://m3.material.io/)

## 🤝 贡献

如需改进或添加新功能,请遵循以下规范:
1. 保持代码风格一致
2. 添加适当的TypeScript类型定义
3. 编写清晰的注释
4. 确保向后兼容
5. 更新相关文档

## 📄 许可证

遵循项目主许可证

---

**最后更新**: 2025-10-11
**作者**: AI助手
**版本**: 1.0.0

