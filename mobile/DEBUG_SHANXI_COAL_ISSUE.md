# 陕西煤业(601225)数据显示异常调试指南

## 问题描述

React Native Web端显示陕西煤业今年有到17块钱，但实际上：
- 原本的Web端没有这个问题
- 实际上今年陕西煤业并没有到过17块钱
- 后端API是一样的

## 调试日志已添加

### 1. API服务层 (`mobile/src/services/apiService.ts`)

已添加详细日志：
- ✅ 请求参数日志
- ✅ 响应数据统计（数据长度、日期范围、价格范围）
- ✅ 数据样本输出（前2条和后2条）

**关键日志标识**：
- `[ApiService]` - API服务层日志
- `📊` - 请求日志
- `📈` - 响应详情日志
- `📥` - 一般响应日志

### 2. 股票详情页 (`mobile/src/screens/Market/StockDetailScreen.tsx`)

已添加详细日志：
- ✅ 加载股票数据开始
- ✅ 股票基本信息响应
- ✅ 日线数据请求参数（日期范围）
- ✅ 日线数据响应统计
- ✅ 价格范围分析（最高价、最低价、开盘价、收盘价）

**关键日志标识**：
- `[StockDetail]` - 股票详情页日志
- `🔍` - 开始加载
- `📊` - 基本信息
- `📅` - 日期参数
- `📈` - 响应统计
- `📉` - 数据详情

### 3. K线图表组件 (`mobile/src/components/KLineChart.tsx`)

已添加详细日志：
- ✅ 渲染开始（原始数据）
- ✅ K线数据处理（OHLC数组）
- ✅ 成交量数据处理
- ✅ 日期数据格式化
- ✅ 移动平均线计算
- ✅ ECharts配置生成
- ✅ WebView内部日志

**关键日志标识**：
- `[KLineChart]` - K线图表组件日志
- `[KLineChart WebView]` - WebView内部日志
- `📊` - 开始渲染
- `📈` - K线数据处理
- `📅` - 日期数据
- `📐` - 移动平均线
- `⚙️` - ECharts配置
- `🌐` - HTML生成
- `🎨` - WebView初始化

## 调试步骤

### 步骤1: 运行测试脚本验证后端数据

```bash
cd mobile
node test-api.js
```

这个脚本会：
1. 测试3个月、半年、1年的数据
2. 显示价格统计（最高价、最低价、价格区间）
3. 找出最高价和最低价对应的日期
4. 检查是否有价格 >= 17元的数据
5. 显示数据样本（前5条和后5条）

**期望结果**：后端数据不应该有今年超过17元的价格

### 步骤2: 启动React Native应用并查看日志

```bash
# 确保后端服务运行在 http://127.0.0.1:8081

# 启动React Native（根据你的环境）
npm start
# 或
npx react-native start
```

### 步骤3: 在应用中搜索并查看陕西煤业

1. 打开应用
2. 搜索 "陕西煤业" 或 "601225"
3. 进入股票详情页
4. 选择不同的时间范围（1月、3月、半年、1年）
5. 查看K线图表

### 步骤4: 分析日志输出

打开开发者工具的控制台，查看日志输出，按照以下顺序分析：

#### 4.1 API请求阶段
```
🚀 [ApiService] 请求
📊 [StockDetail] 开始加载股票数据
📅 [StockDetail] 请求日线数据
```

**检查项**：
- [ ] `stockCode` 是否正确 (601225)
- [ ] `start_date` 和 `end_date` 格式是否正确 (YYYYMMDD)
- [ ] `timeRange` 是否符合预期

#### 4.2 API响应阶段
```
📥 [ApiService] API响应
📈 [ApiService] 日线数据响应详情
📈 [StockDetail] 日线数据响应
📉 [StockDetail] 日线数据详情
```

**检查项**：
- [ ] `dataLength` - 数据条数是否合理
- [ ] `dateRange` - 日期范围是否正确
- [ ] `priceRange.highest` - **最高价是否超过17元？**
- [ ] `priceRange.lowest` - 最低价是否合理
- [ ] `sampleData` - 前后数据样本是否正常

#### 4.3 图表渲染阶段
```
📊 [KLineChart] 开始渲染图表
📈 [KLineChart] K线数据处理完成
📊 [KLineChart] 成交量数据
📅 [KLineChart] 日期数据
📐 [KLineChart] 移动平均线计算完成
⚙️ [KLineChart] ECharts配置生成完成
🌐 [KLineChart] HTML生成完成
🎨 [KLineChart WebView] 开始初始化ECharts
📊 [KLineChart WebView] ECharts配置
✅ [KLineChart WebView] ECharts渲染完成
```

**检查项**：
- [ ] `价格统计.最高价` - **K线数据中的最高价是否超过17元？**
- [ ] `价格统计.最低价` - 最低价是否合理
- [ ] K线数据的OHLC值是否正常
- [ ] 日期格式化是否正确

## 可能的问题原因

### 原因1: 后端数据问题
如果测试脚本（步骤1）显示后端数据有 >= 17元的价格：
- 可能是数据源问题
- 可能是复权计算问题
- 需要检查后端AKTools服务

### 原因2: 数据类型转换问题
如果后端数据正常，但K线图显示异常：
- 检查 `parseFloat(String(...))` 转换
- 检查数据字段名是否正确 (`high`, `low`, `open`, `close`)
- 检查是否有NaN或Infinity值

### 原因3: ECharts配置问题
如果数据处理正常，但图表渲染异常：
- 检查ECharts的 `yAxis.scale` 配置
- 检查 `yAxis.boundaryGap` 设置
- 检查dataZoom的 `start` 和 `end` 值

### 原因4: 日期范围问题
如果不同时间范围显示不同结果：
- 检查日期计算逻辑 (`getDateRange`)
- 检查是否获取了超出预期范围的数据
- 检查日期格式转换

## 对比Web端实现

Web端的K线图实现在：
- `web/static/js/modules/charts.js`
- `web/static/js/core/client.js`

**关键差异点**：
1. Web端使用本地ECharts库，React Native使用CDN
2. Web端直接操作DOM，React Native使用WebView
3. 数据处理逻辑应该相同

建议对比：
- [ ] 数据获取API调用是否一致
- [ ] K线数据格式化逻辑是否一致
- [ ] ECharts配置是否一致

## 下一步行动

1. ✅ 运行 `node mobile/test-api.js` 验证后端数据
2. ✅ 启动React Native应用并搜索陕西煤业
3. ✅ 查看控制台日志输出
4. ✅ 根据日志分析问题出在哪个阶段
5. ✅ 对比Web端和React Native的差异

## 临时禁用日志

如果调试完成，想禁用详细日志，可以：

1. 注释掉相关的 `console.log` 语句
2. 或者添加一个调试开关：

```typescript
// 在 mobile/src/constants/config.ts 中添加
export const DEBUG_MODE = false;

// 然后在代码中使用
if (DEBUG_MODE) {
  console.log('...');
}
```

## 联系信息

如果问题仍然存在，请提供：
- 测试脚本的完整输出
- React Native控制台的完整日志
- 截图显示问题
- 使用的时间范围和股票代码

