# 回测功能快速测试指南

## 🧪 如何测试回测跳转功能

### 步骤 1: 启动应用

```bash
cd web-react
npm start
```

### 步骤 2: 访问策略管理页面

在浏览器中访问:
```
http://localhost:3000/strategies
```

### 步骤 3: 点击"运行回测"按钮

1. 在任意策略卡片上找到"运行回测"按钮
2. 点击该按钮

### 步骤 4: 验证跳转

应该看到:
- ✅ 页面自动跳转到 `/backtest`
- ✅ 显示"已选择 1 个策略进行回测"提示
- ✅ 显示选中策略的ID
- ✅ 控制台输出: `[BacktestPage] 选中的策略IDs: ['strategy-id']`

### 步骤 5: 测试返回

点击"返回策略列表"按钮，应该回到策略管理页面

## 🔍 调试技巧

### 查看 Redux 状态

打开浏览器的 Redux DevTools:

```javascript
// 应该看到 backtest 状态
{
  backtest: {
    selectedStrategyIds: ['布林带策略'],
    config: {
      name: '',
      startDate: '',
      endDate: '',
      initialCash: 1000000,
      commission: 0.0003,
      symbols: []
    }
  }
}
```

### 查看控制台日志

点击运行回测时，控制台应该输出:
```
运行回测: 布林带策略
[BacktestPage] 选中的策略IDs: ['bollinger_strategy']
```

## 🎯 测试场景

### 场景 1: 单策略回测
1. 点击任一策略的"运行回测"
2. 验证跳转到回测页面
3. 验证显示1个策略

### 场景 2: 返回继续选择
1. 从回测页面返回策略列表
2. 点击另一个策略的"运行回测"
3. 验证之前的选择被替换

### 场景 3: 刷新页面
1. 在回测页面刷新浏览器
2. ⚠️ 注意: Redux 状态会丢失（未持久化）
3. 需要重新从策略页面选择

## 📊 预期行为

| 操作 | 预期结果 |
|------|---------|
| 点击"运行回测" | 跳转到 `/backtest` 页面 |
| Redux 状态更新 | `selectedStrategyIds` 包含策略ID |
| 回测页面显示 | 显示选中策略数量和ID |
| 点击"返回" | 回到 `/strategies` 页面 |
| 刷新回测页面 | 状态丢失，提示选择策略 |

## 🐛 常见问题

### Q: 点击按钮没反应？
**A**: 检查:
1. 浏览器控制台是否有错误
2. Redux DevTools 是否安装
3. `backtestSlice` 是否正确导入

### Q: 跳转后没有显示策略？
**A**: 检查:
1. Redux store 是否正确配置
2. `useAppSelector` 是否正确使用
3. 控制台日志是否输出策略ID

### Q: 状态为什么会丢失？
**A**: 
- Redux 状态默认存储在内存中
- 刷新页面会清空状态
- 后续可以使用 `redux-persist` 实现持久化

## 🔧 手动测试 Redux Action

在浏览器控制台中:

```javascript
// 获取 store
const store = window.store; // 需要在开发模式暴露

// 手动设置策略
store.dispatch({
  type: 'backtest/setSelectedStrategies',
  payload: ['macd_strategy', 'ma_crossover']
});

// 查看状态
store.getState().backtest;
```

## 📝 测试清单

- [ ] 点击"运行回测"按钮
- [ ] 页面跳转到 `/backtest`
- [ ] 显示选中策略信息
- [ ] 控制台输出正确日志
- [ ] Redux 状态正确更新
- [ ] "返回策略列表"按钮工作
- [ ] 可以选择不同的策略
- [ ] 策略ID正确传递

## 🎉 成功标志

如果所有测试通过，你应该看到:

```
✅ 策略管理页面正常显示
✅ 点击"运行回测"成功跳转
✅ 回测页面显示: "已选择 1 个策略进行回测"
✅ 策略ID列表正确显示
✅ 返回按钮正常工作
✅ 控制台无错误
```

## 🚀 下一步

测试通过后，可以开始实现:
1. 回测配置表单
2. 策略详情展示
3. 回测API调用
4. 结果可视化

---

**提示**: 这是功能的基础框架测试，完整的回测功能需要后续开发。

