# 调试功能改进总结

## 改动概述

为了调试React Native Web端显示陕西煤业价格异常的问题，我们添加了详细的调试日志和测试工具。

## 修改的文件

### 1. API服务层 - `mobile/src/services/apiService.ts`

**改动内容**:
- ✅ 增强了API请求日志，显示完整URL和参数
- ✅ 增强了API响应日志，显示数据类型和长度
- ✅ 在 `getDailyData` 方法中添加专门的日线数据分析日志
- ✅ 添加价格范围统计（最高价、最低价、开盘价、收盘价）
- ✅ 添加日期范围统计
- ✅ 添加数据样本输出（前2条和后2条）

**日志标识**: `[ApiService]`

### 2. 股票详情页 - `mobile/src/screens/Market/StockDetailScreen.tsx`

**改动内容**:
- ✅ 在 `loadStockData` 函数开始处添加加载日志
- ✅ 添加股票基本信息响应日志
- ✅ 添加日线数据请求参数日志（日期范围、时间窗口）
- ✅ 添加日线数据响应统计日志
- ✅ 添加详细的价格范围分析日志
- ✅ 添加数据样本日志（第一条和最后一条）

**日志标识**: `[StockDetail]`

### 3. K线图表组件 - `mobile/src/components/KLineChart.tsx`

**改动内容**:
- ✅ 在 `renderChart` 函数开始处添加原始数据日志
- ✅ 添加K线数据处理完成日志（OHLC数组、价格统计）
- ✅ 添加成交量数据日志
- ✅ 添加日期数据格式化日志
- ✅ 添加移动平均线计算完成日志
- ✅ 添加ECharts配置生成日志
- ✅ 添加HTML生成日志
- ✅ 在WebView内部添加ECharts初始化和渲染日志

**日志标识**: `[KLineChart]`, `[KLineChart WebView]`

## 新增的文件

### 测试和调试工具

1. **`mobile/test-api.js`** (重写)
   - 专门测试陕西煤业(601225)的API数据
   - 测试3个月、半年、1年的数据
   - 显示详细的价格统计
   - 自动检测价格 >= 17元的异常数据
   - 显示数据样本和最高价/最低价日期

2. **`mobile/debug-test.bat`** (Windows批处理脚本)
   - 自动检查Node.js环境
   - 自动安装依赖（node-fetch）
   - 运行API测试
   - 显示下一步操作指引

3. **`mobile/debug-test.sh`** (Linux/Mac Shell脚本)
   - 功能同上（Unix版本）
   - 已添加执行权限

### 文档

4. **`mobile/DEBUG_SHANXI_COAL_ISSUE.md`**
   - 陕西煤业问题专项调试指南
   - 详细的调试步骤
   - 问题原因分析
   - 检查清单

5. **`mobile/DEBUG_LOGS_GUIDE.md`**
   - 完整的调试日志使用指南
   - 日志标识符和Emoji图标说明
   - 日志流程图
   - 各种调试场景的详细步骤
   - 控制台过滤技巧
   - 故障排除检查清单

6. **`mobile/DEBUG_CHANGES_SUMMARY.md`** (本文件)
   - 所有改动的总结
   - 使用说明

## 使用方法

### 方法1: 使用自动化脚本（推荐）

**Windows**:
```cmd
cd mobile
debug-test.bat
```

**Linux/Mac**:
```bash
cd mobile
./debug-test.sh
```

### 方法2: 手动运行测试

```bash
cd mobile

# 安装依赖（首次运行）
npm install node-fetch

# 运行测试
node test-api.js
```

### 方法3: 在React Native应用中查看日志

1. 启动后端服务（确保运行在 http://127.0.0.1:8081）

2. 启动React Native应用:
   ```bash
   cd mobile
   npm start
   ```

3. 打开浏览器开发者工具（F12）

4. 在应用中搜索并查看陕西煤业(601225)

5. 观察控制台日志输出

## 日志级别

所有日志都使用 `console.log`，并带有清晰的前缀标识：

- **信息日志**: `console.log` + 🔍📊📈📅等Emoji
- **错误日志**: `console.error` + ❌
- **警告日志**: `console.warn` + ⚠️（在测试脚本中）

## 性能影响

### 预计影响
- 额外的日志输出会增加少量CPU开销
- 对于大数据量（>100条记录），使用了数据样本而非完整数据输出
- WebView内的日志会同步到外部控制台

### 性能优化
- 使用了条件日志（数据量大时只输出样本）
- 避免了深度克隆和JSON序列化大对象
- 可以通过注释或开关轻松禁用

## 禁用调试日志

### 临时禁用（开发阶段）

注释掉相关的 `console.log` 语句：

```typescript
// console.log('📊 [KLineChart] 开始渲染图表', { ... });
```

### 永久禁用（生产环境）

**方法1**: 使用环境变量

```typescript
if (__DEV__) {
  console.log('...');
}
```

**方法2**: 添加全局调试开关

在 `mobile/src/constants/config.ts`:
```typescript
export const DEBUG_CONFIG = {
  enableApiLogs: __DEV__,
  enableComponentLogs: __DEV__,
  enableChartLogs: __DEV__,
};
```

然后在代码中:
```typescript
import { DEBUG_CONFIG } from '@/constants/config';

if (DEBUG_CONFIG.enableApiLogs) {
  console.log('[ApiService] ...');
}
```

**方法3**: 使用日志库

使用 `react-native-logs` 或类似库，支持日志级别控制。

## 调试流程

```
1. 运行 debug-test.bat/sh
   ↓
2. 查看测试输出，确认后端数据是否正常
   ↓
3. 如果后端数据正常 → 前端问题
   ↓
4. 启动React Native应用
   ↓
5. 打开浏览器开发者工具
   ↓
6. 搜索并查看陕西煤业
   ↓
7. 按照日志流程分析问题出在哪个环节
   ↓
8. 参考 DEBUG_SHANXI_COAL_ISSUE.md 定位问题
```

## 问题排查路径

### 路径1: 后端数据问题
```
test-api.js 显示价格 >= 17元
    ↓
检查后端服务
    ↓
检查AKTools数据源
    ↓
检查复权计算
```

### 路径2: 前端数据处理问题
```
后端数据正常，但前端显示异常
    ↓
检查 [ApiService] 日志
    ↓
检查 [StockDetail] 日志
    ↓
定位数据转换或解析问题
```

### 路径3: 图表渲染问题
```
数据处理正常，但图表显示异常
    ↓
检查 [KLineChart] 日志
    ↓
检查 [KLineChart WebView] 日志
    ↓
检查 ECharts 配置
```

## 关键检查点

在调试过程中，重点关注以下数据点：

### API响应阶段
- ✅ `priceRange.highest` - 最高价是否异常？
- ✅ `priceRange.lowest` - 最低价是否合理？
- ✅ `dateRange` - 日期范围是否正确？
- ✅ `dataLength` - 数据量是否合理？

### 数据处理阶段
- ✅ `价格范围.最高` - 处理后的最高价
- ✅ `价格范围.最低` - 处理后的最低价
- ✅ `第一条/最后一条` - 数据样本是否正常

### 图表渲染阶段
- ✅ `价格统计.最高价` - K线数据的最高价
- ✅ `价格统计.最低价` - K线数据的最低价
- ✅ `klineData[0]` - 第一条K线 [开, 收, 低, 高]
- ✅ `dates[0]` - 第一个日期

## 对比分析

### 与Web端对比

原始Web端实现文件：
- `web/static/js/core/client.js`
- `web/static/js/modules/charts.js`

对比要点：
1. API调用参数是否一致
2. 数据格式化逻辑是否一致
3. K线数据数组顺序是否一致（[开, 收, 低, 高]）
4. ECharts配置是否一致

### 数据流对比

**Web端**:
```
client.makeRequest() 
  → charts.createPriceChart() 
  → ECharts渲染
```

**React Native Web**:
```
apiService.getDailyData() 
  → StockDetailScreen.loadStockData() 
  → KLineChart.renderChart() 
  → WebView中ECharts渲染
```

## 后续优化建议

1. **日志分级**: 实现 DEBUG/INFO/WARN/ERROR 级别
2. **日志开关**: 添加全局配置开关
3. **日志过滤**: 在生产环境自动禁用
4. **性能监控**: 添加性能计时日志
5. **错误追踪**: 集成错误追踪服务（如Sentry）

## 相关资源

- [陕西煤业问题调试指南](DEBUG_SHANXI_COAL_ISSUE.md)
- [调试日志使用指南](DEBUG_LOGS_GUIDE.md)
- [React Native调试文档](https://reactnative.dev/docs/debugging)
- [ECharts配置文档](https://echarts.apache.org/zh/option.html)

## 版本信息

- 创建日期: 2025-10-08
- React Native版本: (根据项目实际版本)
- 调试目标: 陕西煤业(601225)价格显示异常

## 联系与反馈

如果您在使用过程中遇到问题或有改进建议，请：

1. 查看现有文档是否有解决方案
2. 检查控制台是否有错误信息
3. 提供完整的日志输出和问题描述
4. 包含环境信息（操作系统、浏览器、Node版本等）

---

**重要提示**: 这些调试日志仅用于开发和调试阶段，请勿在生产环境中保留过多的详细日志。

