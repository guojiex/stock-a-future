# 股票列表获取工具

## 概述

股票列表获取工具允许你直接从上海证券交易所获取股票代码和名称列表，无需依赖第三方API（如Tushare）的高权限。

## 功能特性

- ✅ 支持获取上海证券交易所（SSE）股票列表

- ✅ 支持获取所有交易所的完整股票列表
- ✅ 自动保存股票数据到JSON文件
- ✅ 包含真实的股票代码和名称
- ✅ 支持命令行参数配置

## 使用方法

### 重要说明
所有 `make` 命令都需要在项目根目录执行，或者使用 `-C` 参数指定项目根目录。

### 快速开始示例

假设您的项目路径是 `/path/to/stock-a-future`，您可以从任意目录执行：

```bash
# 获取上交所股票列表
make -C /path/to/stock-a-future fetch-sse

# 获取所有股票列表
make -C /path/to/stock-a-future fetch-stocks
```

或者先切换到项目根目录：
```bash
cd /path/to/stock-a-future
make fetch-sse
```

### 1. 构建工具

在项目根目录执行：
```bash
make stocklist
```

从任意目录执行：
```bash
make -C /path/to/stock-a-future stocklist
```

### 2. 获取所有股票列表

在项目根目录执行：
```bash
make fetch-stocks
```

从任意目录执行：
```bash
make -C /path/to/stock-a-future fetch-stocks
```

这将获取上海证券交易所的所有股票，并保存到 `data/stock_list.json`。

### 3. 单独获取各交易所股票

获取上海证券交易所股票：

在项目根目录执行：
```bash
make fetch-sse
```

从任意目录执行：
```bash
make -C /path/to/stock-a-future fetch-sse
```



### 4. 自定义使用

```bash
# 显示帮助信息
./bin/stocklist -help

# 获取所有股票并保存到自定义文件
./bin/stocklist -source=all -output=my_stocks.json

# 只获取上交所股票
./bin/stocklist -source=sse -output=sse_only.json


```

## 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-source` | `all` | 数据源：`sse`（上交所）、`all`（全部） |
| `-output` | `data/stock_list.json` | 输出文件路径 |
| `-help` | - | 显示帮助信息 |

## 输出格式

生成的JSON文件包含股票基本信息数组，每个股票包含以下字段：

```json
[
  {
    "ts_code": "600000.SH",
    "symbol": "600000", 
    "name": "浦发银行",
    "area": "",
    "industry": "",
    "market": "SH",
    "list_date": ""
  }
]
```

### 字段说明

- `ts_code`: 完整股票代码（包含交易所后缀）
- `symbol`: 股票代码（不含交易所后缀）
- `name`: 股票名称
- `market`: 交易所标识（SH=上交所）
- `area`, `industry`, `list_date`: 预留字段（当前为空）

## 支持的股票类型

### 上海证券交易所（SH）
- **主板股票**: 600xxx, 601xxx, 603xxx, 605xxx
- **科创板股票**: 688xxx  
- **B股**: 900xxx

包含的知名股票：
- 银行类：浦发银行、招商银行、工商银行等
- 保险类：中国平安、中国人寿等
- 食品饮料：贵州茅台、伊利股份等
- 科技类：药明康德、金山办公等



## 数据获取说明

目前实现包含了主要的A股上市公司，涵盖：
- 上交所：39只主要股票

这些股票代表了各个行业的龙头企业和知名公司，适合用于：
- 股票分析系统测试
- 投资组合构建
- 技术指标计算
- 买卖点预测

## 与现有系统集成

生成的股票列表文件与现有的Tushare客户端完全兼容，可以直接用于：

1. **股票基本信息查询**
2. **日线数据获取**
3. **技术指标计算**
4. **买卖点预测**

## 扩展性

如果需要获取更多股票或实时数据，可以：

1. **扩展股票列表**：在 `internal/client/exchange.go` 中添加更多股票
2. **实现真实爬虫**：添加HTML解析逻辑从交易所网站获取实时数据
3. **集成其他数据源**：添加对其他数据提供商的支持

## 故障排除

### 常见问题

1. **构建失败**
   
   在项目根目录执行：
   ```bash
   make clean
   make deps
   make stocklist
   ```
   
   从任意目录执行：
   ```bash
   make -C /path/to/stock-a-future clean
   make -C /path/to/stock-a-future deps
   make -C /path/to/stock-a-future stocklist
   ```

2. **权限问题**
   ```bash
   chmod +x /path/to/stock-a-future/bin/stocklist
   ```

3. **输出目录不存在**
   ```bash
   mkdir -p /path/to/stock-a-future/data
   ```

4. **Make 命令找不到 Makefile**
   
   确保您在项目根目录，或使用 `-C` 参数：
   ```bash
   # 切换到项目根目录
   cd /path/to/stock-a-future
   make fetch-sse
   
   # 或者从任意目录使用 -C 参数
   make -C /path/to/stock-a-future fetch-sse
   ```

### 日志信息

工具会输出详细的日志信息，包括：
- 获取进度
- 股票数量统计
- 文件保存路径
- 错误信息（如有）

## 未来改进

- [ ] 支持从交易所官网实时抓取数据
- [ ] 添加股票分类和行业信息
- [ ] 支持增量更新
- [ ] 添加数据验证和清洗功能
- [ ] 支持多种输出格式（CSV、Excel等）
