# Stock Components 使用说明

本目录包含股票相关的展示组件。

## 组件列表

### 1. PredictionSignalsView (Material-UI 版本)

**文件**: `PredictionSignalsView.tsx`

**用途**: 在股票详情页面展示买卖预测信号

**特点**:
- 使用 Material-UI 组件库
- 自动从 API 获取数据
- Accordion 折叠面板展示

**使用示例**:
```tsx
import PredictionSignalsView from './components/stock/PredictionSignalsView';

<PredictionSignalsView stockCode="000001.SZ" />
```

**Props**:
- `stockCode: string` - 股票代码（必需）

---

### 2. PredictionsView (DaisyUI 版本)

**文件**: `PredictionsView.tsx`

**用途**: 独立的买卖预测展示组件，可在任何页面使用

**特点**:
- 使用 DaisyUI (TailwindCSS) 组件库
- 需要手动传入数据
- 支持日期点击回调
- 卡片式现代化设计

**使用示例**:
```tsx
import PredictionsView from './components/stock/PredictionsView';
import { useGetPredictionsQuery } from './services/api';

function MyPage() {
  const { data, isLoading } = useGetPredictionsQuery('000001.SZ');
  
  return (
    <PredictionsView 
      data={data?.data || null}
      isLoading={isLoading}
      onDateClick={(date) => {
        // 处理日期点击，例如跳转到 K 线图
        console.log('点击日期:', date);
      }}
    />
  );
}
```

**Props**:
- `data: PredictionResult | null` - 预测数据（必需）
- `isLoading?: boolean` - 加载状态（可选）
- `onDateClick?: (date: string) => void` - 日期点击回调（可选）

---

### 3. TechnicalIndicatorsView

**文件**: `TechnicalIndicatorsView.tsx`

**用途**: 展示技术指标数据

**使用示例**:
```tsx
import TechnicalIndicatorsView from './components/stock/TechnicalIndicatorsView';

<TechnicalIndicatorsView data={indicatorsData} />
```

---

### 4. FundamentalDataView

**文件**: `FundamentalDataView.tsx`

**用途**: 展示基本面数据

**使用示例**:
```tsx
import FundamentalDataView from './components/stock/FundamentalDataView';

<FundamentalDataView data={fundamentalData} />
```

---

## 完整示例：在自定义页面中使用预测组件

### Material-UI 版本（推荐用于已有 MUI 的页面）

```tsx
import React from 'react';
import { Container, Box, Typography } from '@mui/material';
import PredictionSignalsView from '../components/stock/PredictionSignalsView';

const MyPredictionsPage: React.FC = () => {
  const stockCode = '000001.SZ'; // 可以从路由参数获取

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" gutterBottom>
        买卖预测分析
      </Typography>
      
      <PredictionSignalsView stockCode={stockCode} />
    </Container>
  );
};

export default MyPredictionsPage;
```

### DaisyUI 版本（推荐用于需要更多控制的场景）

```tsx
import React, { useState } from 'react';
import PredictionsView from '../components/stock/PredictionsView';
import { useGetPredictionsQuery } from '../services/api';

const MyPredictionsPage: React.FC = () => {
  const [stockCode, setStockCode] = useState('000001.SZ');
  const { data, isLoading, refetch } = useGetPredictionsQuery(stockCode);

  const handleDateClick = (date: string) => {
    // 实现跳转到 K 线图对应日期的逻辑
    console.log('跳转到日期:', date);
    // 例如：navigate(`/stock/${stockCode}/kline?date=${date}`);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">买卖预测分析</h1>
      
      {/* 股票选择器 */}
      <div className="mb-6">
        <input
          type="text"
          value={stockCode}
          onChange={(e) => setStockCode(e.target.value)}
          className="input input-bordered w-full max-w-xs"
          placeholder="输入股票代码"
        />
        <button onClick={() => refetch()} className="btn btn-primary ml-2">
          刷新
        </button>
      </div>

      {/* 预测展示 */}
      <PredictionsView 
        data={data?.data || null}
        isLoading={isLoading}
        onDateClick={handleDateClick}
      />
    </div>
  );
};

export default MyPredictionsPage;
```

## 数据类型

### PredictionResult

```typescript
interface PredictionResult {
  ts_code: string;              // 股票代码
  trade_date: string;           // 交易日期 YYYYMMDD
  predictions: TradingPointPrediction[];  // 预测列表
  confidence: number;           // 整体置信度 (0-1)
  updated_at: string;           // 更新时间
}
```

### TradingPointPrediction

```typescript
interface TradingPointPrediction {
  type: 'BUY' | 'SELL';        // 买入或卖出
  price: number;                // 预测价格
  date: string;                 // 预测日期 YYYYMMDD
  probability: number;          // 概率 (0-1)
  reason: string;               // 预测理由
  indicators: string[];         // 相关指标
  signal_date: string;          // 信号产生日期
  backtested: boolean;          // 是否已回测
  is_correct?: boolean;         // 预测是否正确
  next_day_price?: number;      // 次日价格
  price_diff?: number;          // 价格差值
  price_diff_ratio?: number;    // 价格差值百分比
}
```

## 样式定制

### Material-UI 版本

可以通过 MUI 的 `sx` prop 或主题定制样式：

```tsx
<PredictionSignalsView 
  stockCode="000001.SZ"
  // 组件内部使用 MUI 主题，可通过 ThemeProvider 全局定制
/>
```

### DaisyUI 版本

可以通过 TailwindCSS 类名或 DaisyUI 主题定制：

```tsx
// 在 tailwind.config.js 中定制主题
module.exports = {
  daisyui: {
    themes: [
      {
        mytheme: {
          "primary": "#your-color",
          "secondary": "#your-color",
          // ...
        },
      },
    ],
  },
}
```

## 注意事项

1. **API 依赖**: 两个组件都依赖后端 API `/api/v1/stocks/{code}/predictions`
2. **错误处理**: 组件内部已处理加载和错误状态
3. **性能**: 使用 RTK Query 自动缓存，避免重复请求
4. **响应式**: 两个版本都支持移动端适配

## 选择建议

- **使用 PredictionSignalsView (MUI)** 如果：
  - 项目已使用 Material-UI
  - 需要快速集成，不需要额外配置
  - 想要统一的 MUI 设计风格

- **使用 PredictionsView (DaisyUI)** 如果：
  - 项目使用 TailwindCSS/DaisyUI
  - 需要更多的自定义控制
  - 需要日期点击等额外功能
  - 偏好卡片式现代化设计
