# 收藏股票信号汇总页面

## 概述

React Web版现在已经添加了与原生Web版相同的**收藏股票信号汇总**功能，可以在一个页面中查看所有收藏股票的买卖信号。

## 功能特性

### 1. 信号汇总展示
- **实时数据**: 自动获取所有收藏股票的最新买卖信号
- **智能分类**: 将股票按买入/卖出/持有信号自动分类
- **置信度排序**: 在每个分类中按预测置信度降序排列

### 2. 统计面板
页面顶部显示四个统计卡片：
- **总数**: 所有收藏股票数量
- **买入信号**: 有买入信号的股票数量（绿色）
- **卖出信号**: 有卖出信号的股票数量（红色）
- **持有**: 建议持有的股票数量（橙色）

### 3. 标签页筛选
通过标签页快速筛选不同类型的信号：
- **全部**: 显示所有股票信号
- **买入**: 只显示有买入信号的股票
- **卖出**: 只显示有卖出信号的股票
- **持有**: 只显示建议持有的股票

### 4. 信号卡片
每个股票显示为一张卡片，包含：
- **股票信息**: 名称、代码、当前价格
- **信号类型**: 买入/卖出/持有/混合，带颜色标识
- **置信度**: 圆形进度条显示预测置信度
- **预测详情**: 
  - 买入/卖出理由
  - 相关技术指标
  - 预测概率
- **交易日期**: 数据更新时间
- **操作按钮**: 点击查看股票详情页

### 5. 计算状态提示
当后台正在计算信号时，会显示进度提示：
```
信号正在计算中... (已完成/总数)
```

### 6. 响应式设计
- **移动端**: 通过底部导航访问
- **桌面端**: 通过顶部标签栏访问
- 卡片布局自动适应不同屏幕尺寸

## 技术实现

### API端点
```typescript
// 获取收藏股票信号
GET /api/v1/favorites/signals

// 响应格式
{
  "success": true,
  "data": {
    "total": 10,
    "signals": [
      {
        "id": "uuid",
        "ts_code": "600976.SH",
        "name": "健民集团",
        "current_price": "32.50",
        "trade_date": "20240118",
        "indicators": { /* 技术指标 */ },
        "predictions": [
          {
            "type": "BUY",
            "probability": 0.85,
            "reason": "MACD金叉",
            "indicators": ["MACD", "RSI"]
          }
        ],
        "updated_at": "2024-01-18 15:30:00"
      }
    ],
    "calculating": false,
    "calculation_status": { /* 计算状态 */ }
  }
}
```

### 核心组件

#### SignalsPage
主页面组件，负责：
- 调用API获取数据
- 实现标签页筛选逻辑
- 统计数据计算
- 刷新数据功能

#### SignalCard
信号卡片组件，负责：
- 展示单个股票的信号详情
- 信号类型颜色标识
- 置信度可视化
- 跳转到股票详情页

### 数据类型
```typescript
// 收藏股票信号
interface FavoriteSignal {
  id: string;
  ts_code: string;
  name: string;
  group_id?: string;
  current_price: string;
  trade_date: string;
  indicators: TechnicalIndicators;
  predictions: TradingPointPrediction[];
  updated_at: string;
}

// 信号汇总响应
interface FavoritesSignalsResponse {
  total: number;
  signals: FavoriteSignal[];
  calculating?: boolean;
  calculation_status?: {
    status: string;
    message: { message: string };
    detail: {
      is_calculating: boolean;
      completed: number;
      total: number;
      last_updated: string;
    };
  };
}
```

## 使用方式

### 桌面端
1. 点击顶部导航栏的 **"信号"** 标签
2. 查看信号统计和列表
3. 使用标签页筛选不同类型的信号
4. 点击 **"查看详情"** 按钮查看股票完整信息

### 移动端
1. 点击底部导航栏的 **"信号"** 图标（铃铛图标）
2. 在移动设备上浏览信号列表
3. 卡片自动适应小屏幕显示

### 刷新数据
点击页面右上角的刷新按钮（🔄）可以手动刷新信号数据

## 信号分类规则

### 买入信号
- 至少有一个 `type: "BUY"` 的预测
- 没有 `type: "SELL"` 的预测
- 卡片边框为绿色

### 卖出信号
- 至少有一个 `type: "SELL"` 的预测
- 没有 `type: "BUY"` 的预测
- 卡片边框为红色

### 持有信号
- 没有买入或卖出预测
- 卡片边框为橙色

### 混合信号
- 同时有买入和卖出预测
- 根据置信度判断主要信号类型
- 卡片边框为紫色

## 置信度排序

在每个信号类型中，股票按以下规则排序：
1. 计算每个股票所有预测的最高置信度
2. 置信度高的排在前面（降序）
3. 帮助用户快速找到最可靠的信号

## 性能优化

### 后端优化
- 使用预计算的信号数据，避免实时计算
- 对缺失数据进行实时补充（最多20个）
- 异步信号计算服务持续更新数据

### 前端优化
- 使用 RTK Query 自动缓存数据
- 避免不必要的重新渲染
- 懒加载股票详情数据

## 与原生Web版的对比

| 特性 | 原生Web版 | React Web版 |
|------|-----------|------------|
| 信号汇总 | ✅ | ✅ |
| 实时刷新 | ✅ | ✅ |
| 信号分类 | ✅ | ✅ |
| 置信度排序 | ✅ | ✅ |
| 响应式设计 | ✅ | ✅ |
| 动画效果 | 淡入淡出 | Material-UI动画 |
| 框架 | 原生JS | React + TypeScript |
| 状态管理 | 本地状态 | Redux Toolkit |

## 文件清单

### 新增文件
- `web-react/src/pages/SignalsPage.tsx` - 信号汇总页面
- `web-react/SIGNALS_PAGE.md` - 功能文档

### 修改文件
- `web-react/src/types/stock.ts` - 添加信号类型定义
- `web-react/src/services/api.ts` - 添加API端点
- `web-react/src/App.tsx` - 添加路由
- `web-react/src/components/common/Layout.tsx` - 添加导航入口

## 后续改进建议

1. **添加筛选功能**
   - 按分组筛选
   - 按置信度区间筛选
   - 搜索功能

2. **增强交互**
   - 信号详情弹窗
   - 批量操作功能
   - 导出信号列表

3. **性能优化**
   - 虚拟滚动（处理大量数据）
   - 分页加载
   - WebSocket实时更新

4. **数据可视化**
   - 信号趋势图表
   - 成功率统计
   - 收益率分析

## 故障排查

### 无法加载信号数据
1. 检查后端服务是否正常运行
2. 确认有收藏的股票
3. 查看浏览器控制台错误信息

### 信号数据不更新
1. 点击刷新按钮手动更新
2. 检查后端信号计算服务状态
3. 查看 `/api/v1/favorites/signals` API响应

### 页面显示异常
1. 清除浏览器缓存
2. 检查网络连接
3. 查看控制台错误日志

## 联系支持

如有问题，请查看：
- 后端API文档: `docs/integration/AKTOOLS_INTEGRATION.md`
- 信号服务文档: `docs/features/SIGNAL_SERVICE.md`
- 项目主README: `README.md`

