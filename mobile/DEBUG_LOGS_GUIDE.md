# React Native调试日志指南

## 概述

本指南说明如何使用已添加的详细调试日志来排查React Native Web端的数据显示问题。

## 快速开始

### 1. 运行后端API测试

```bash
cd mobile
node test-api.js
```

这将测试陕西煤业(601225)的后端API数据，验证数据源是否正确。

### 2. 启动React Native应用

```bash
# 在mobile目录下
npm start

# 或者根据你的环境
npx react-native run-web
npx react-native run-android
npx react-native run-ios
```

### 3. 打开浏览器开发者工具

- Chrome: F12 或 Ctrl+Shift+I (Windows/Linux) / Cmd+Option+I (Mac)
- 选择 "Console" 标签

### 4. 在应用中操作并查看日志

1. 搜索股票（如：陕西煤业、601225）
2. 进入股票详情页
3. 切换不同的时间范围
4. 观察控制台日志输出

## 日志结构

### 日志标识符

所有日志都有清晰的标识符，方便过滤和搜索：

| 标识符 | 模块 | 说明 |
|--------|------|------|
| `[ApiService]` | API服务层 | HTTP请求和响应 |
| `[StockDetail]` | 股票详情页 | 页面数据加载和处理 |
| `[KLineChart]` | K线图表组件 | 图表渲染和数据处理 |
| `[KLineChart WebView]` | WebView内部 | ECharts渲染 |

### Emoji图标含义

| Emoji | 含义 | 使用场景 |
|-------|------|---------|
| 🚀 | 开始/启动 | 请求开始 |
| 🔍 | 搜索/查询 | 开始加载数据 |
| 📡 | 网络请求 | API调用 |
| 📊 | 数据统计 | 数据量、基本信息 |
| 📈 | 价格数据 | K线、涨跌数据 |
| 📉 | 详细数据 | 详细的数据分析 |
| 📅 | 日期 | 日期相关信息 |
| 📐 | 计算 | 技术指标计算 |
| ⚙️ | 配置 | ECharts配置 |
| 🌐 | HTML/Web | HTML生成 |
| 🎨 | 渲染 | 图表渲染 |
| 📥 | 响应 | API响应 |
| 📄 | 文档/数据 | 数据样本 |
| ✅ | 成功 | 操作成功 |
| ❌ | 错误 | 操作失败 |
| ⚠️ | 警告 | 异常情况 |

## 日志流程图

```
用户操作
    ↓
📊 [StockDetail] 开始加载股票数据
    ↓
📅 [StockDetail] 请求日线数据 (start_date, end_date)
    ↓
🚀 [ApiService] API请求 (url, params)
    ↓
📥 [ApiService] API响应 (status, dataLength)
    ↓
📈 [ApiService] 日线数据响应详情 (priceRange, dateRange)
    ↓
📈 [StockDetail] 日线数据响应 (success, dataLength)
    ↓
📉 [StockDetail] 日线数据详情 (总记录数, 价格范围)
    ↓
📊 [KLineChart] 开始渲染图表
    ↓
📈 [KLineChart] K线数据处理完成 (价格统计)
    ↓
📊 [KLineChart] 成交量数据
    ↓
📅 [KLineChart] 日期数据
    ↓
📐 [KLineChart] 移动平均线计算完成
    ↓
⚙️ [KLineChart] ECharts配置生成完成
    ↓
🌐 [KLineChart] HTML生成完成
    ↓
🎨 [KLineChart WebView] 开始初始化ECharts
    ↓
📊 [KLineChart WebView] ECharts配置
    ↓
✅ [KLineChart WebView] ECharts渲染完成
```

## 如何使用日志调试

### 场景1: 数据价格异常

**问题**: 显示的价格明显不对（如陕西煤业显示17元）

**调试步骤**:

1. **先验证后端数据**:
   ```bash
   node mobile/test-api.js
   ```
   查看输出中的"价格统计"和"警告"部分

2. **检查API响应**:
   搜索日志: `[ApiService] 日线数据响应详情`
   ```
   📈 [ApiService] 日线数据响应详情: {
     stockCode: "601225",
     priceRange: {
       highest: XX.XX,  // ← 检查这个值
       lowest: XX.XX,
       ...
     }
   }
   ```

3. **检查数据处理**:
   搜索日志: `[StockDetail] 日线数据详情`
   ```
   📉 [StockDetail] 日线数据详情: {
     价格范围: {
       最高: XX.XX,  // ← 检查这个值
       最低: XX.XX,
       ...
     }
   }
   ```

4. **检查图表渲染**:
   搜索日志: `[KLineChart] K线数据处理完成`
   ```
   📈 [KLineChart] K线数据处理完成: {
     价格统计: {
       最高价: XX.XX,  // ← 检查这个值
       ...
     }
   }
   ```

5. **定位问题**:
   - 如果步骤1显示后端数据就有问题 → 后端或数据源问题
   - 如果步骤2显示API响应有问题 → API服务问题
   - 如果步骤3或4显示处理后有问题 → 前端数据处理问题
   - 如果所有数据都正常但图表显示错误 → ECharts配置或渲染问题

### 场景2: 日期范围问题

**问题**: 选择3个月但显示了1年的数据

**调试步骤**:

1. **检查请求参数**:
   搜索日志: `[StockDetail] 请求日线数据`
   ```
   📅 [StockDetail] 请求日线数据: {
     stockCode: "601225",
     start_date: "20240710",  // ← 检查日期计算
     end_date: "20241008",
     timeRange: 90
   }
   ```

2. **检查API响应的日期范围**:
   搜索日志: `[ApiService] 日线数据响应详情`
   ```
   📈 [ApiService] 日线数据响应详情: {
     dateRange: {
       start: "2024-07-10",  // ← 与请求对比
       end: "2024-10-08"
     }
   }
   ```

3. **检查图表的日期数据**:
   搜索日志: `[KLineChart] 日期数据`
   ```
   📅 [KLineChart] 日期数据: {
     第一个日期: "2024-07-10",
     最后一个日期: "2024-10-08",
     ...
   }
   ```

### 场景3: 图表不显示或渲染失败

**问题**: K线图不显示或显示空白

**调试步骤**:

1. **检查是否有数据**:
   搜索: `[StockDetail] 日线数据响应`
   查看 `dataLength` 是否 > 0

2. **检查渲染是否开始**:
   搜索: `[KLineChart] 开始渲染图表`
   如果没有这条日志，说明渲染没有触发

3. **检查数据处理**:
   按顺序查看:
   - `[KLineChart] K线数据处理完成`
   - `[KLineChart] 成交量数据`
   - `[KLineChart] 日期数据`
   - `[KLineChart] 移动平均线计算完成`
   
   如果某一步缺失，说明在该步骤出错

4. **检查ECharts渲染**:
   搜索: `[KLineChart WebView]`
   应该看到:
   - `开始初始化ECharts`
   - `ECharts配置`
   - `ECharts渲染完成`

5. **检查浏览器控制台错误**:
   查看是否有红色的错误信息

## 控制台过滤技巧

### Chrome DevTools过滤

在控制台的过滤框中输入：

- 只看API层: `[ApiService]`
- 只看详情页: `[StockDetail]`
- 只看图表: `[KLineChart]`
- 只看错误: `❌`
- 只看警告: `⚠️`
- 只看价格相关: `价格` 或 `price`
- 看特定股票: `601225`

### 组合过滤

- API错误: `[ApiService] ❌`
- 图表价格数据: `[KLineChart] 价格`
- 详情页日期: `[StockDetail] 📅`

## 性能影响

这些详细日志会产生大量输出，可能影响性能。

### 临时禁用日志

**方法1**: 注释掉不需要的 `console.log`

**方法2**: 添加调试开关

在 `mobile/src/constants/config.ts` 中添加:
```typescript
export const appConfig = {
  // ... 其他配置
  debug: {
    apiService: true,    // API服务日志
    stockDetail: true,   // 股票详情页日志
    klineChart: true,    // K线图表日志
  }
};
```

然后修改日志代码:
```typescript
if (appConfig.debug.apiService) {
  console.log('[ApiService] ...');
}
```

**方法3**: 使用环境变量

```typescript
const DEBUG = __DEV__; // React Native的开发模式标志

if (DEBUG) {
  console.log('...');
}
```

## 日志输出示例

### 正常流程示例

```
🔍 [StockDetail] 开始加载股票数据 { stockCode: "601225", stockName: "陕西煤业", timeRange: "90" }
📅 [StockDetail] 请求日线数据 { stockCode: "601225", start_date: "20240710", end_date: "20241008", timeRange: 90 }
🚀 [ApiService] API请求: { url: "http://...", method: "GET", ... }
📥 [ApiService] API响应: { url: "...", status: 200, success: true, dataLength: 63 }
📈 [ApiService] 日线数据响应详情: { stockCode: "601225", dataLength: 63, priceRange: { highest: 16.85, lowest: 13.21, ... } }
📈 [StockDetail] 日线数据响应: { success: true, hasData: true, dataLength: 63 }
📉 [StockDetail] 日线数据详情: { 总记录数: 63, 价格范围: { 最高: 16.85, 最低: 13.21, ... } }
📊 [KLineChart] 开始渲染图表 { stockCode: "601225", dataLength: 63, ... }
📈 [KLineChart] K线数据处理完成 { klineDataLength: 63, 价格统计: { 最高价: 16.85, 最低价: 13.21, ... } }
📊 [KLineChart] 成交量数据 { volumeDataLength: 63, ... }
📅 [KLineChart] 日期数据 { datesLength: 63, 第一个日期: "2024-07-10", 最后一个日期: "2024-10-08" }
📐 [KLineChart] 移动平均线计算完成 { closeDataLength: 63, ... }
⚙️ [KLineChart] ECharts配置生成完成 { series数量: 5, K线数据点数: 63, ... }
🌐 [KLineChart] HTML生成完成，准备注入WebView
🎨 [KLineChart WebView] 开始初始化ECharts
📊 [KLineChart WebView] ECharts配置: { seriesCount: 5, klineDataLength: 63, ... }
✅ [KLineChart WebView] ECharts渲染完成
```

### 异常流程示例

```
🔍 [StockDetail] 开始加载股票数据 { ... }
📅 [StockDetail] 请求日线数据 { ... }
🚀 [ApiService] API请求: { ... }
❌ [ApiService] API请求失败: { url: "...", error: "Network request failed" }
```

## 故障排除检查清单

- [ ] 后端服务是否运行 (http://127.0.0.1:8081)
- [ ] API响应是否成功 (success: true)
- [ ] 数据是否返回 (dataLength > 0)
- [ ] 价格数据是否合理 (检查最高价、最低价)
- [ ] 日期范围是否正确
- [ ] K线数据是否处理成功
- [ ] ECharts是否渲染完成
- [ ] 是否有JavaScript错误

## 相关文件

- `mobile/src/services/apiService.ts` - API服务层
- `mobile/src/screens/Market/StockDetailScreen.tsx` - 股票详情页
- `mobile/src/components/KLineChart.tsx` - K线图表组件
- `mobile/test-api.js` - 后端API测试脚本
- `mobile/DEBUG_SHANXI_COAL_ISSUE.md` - 陕西煤业问题专项调试指南

## 需要帮助？

如果通过日志仍然无法定位问题，请提供：

1. 完整的控制台日志输出（从"开始加载"到"渲染完成"）
2. `node test-api.js` 的输出
3. 问题截图
4. 操作步骤（股票代码、时间范围等）
5. 环境信息（浏览器、操作系统、React Native版本）

