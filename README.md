# Stock-A-Future - A股股票买卖点预测API

[![Go](https://github.com/guojiex/stock-a-future/actions/workflows/go.yml/badge.svg)](https://github.com/guojiex/stock-a-future/actions/workflows/go.yml)
![SQLite](https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

![股票查询](docs/imgs/股票查询.png)

基于Go语言开发的A股股票买卖点预测系统，支持多种数据源（Tushare、AKTools），提供技术指标计算和买卖点预测功能。

## 功能特性

### 📊 数据获取
- **多数据源支持**：集成Tushare API和AKTools（基于AKShare的免费开源数据源）
- 支持股票日线数据查询
- 自动数据预处理和清洗
- **股票列表工具**：支持从上交所在线获取最新股票列表
- **深交所数据**：使用本地Excel文件（`data/A股列表.xlsx`）提供完整深交所股票数据

### 📈 技术指标计算

![技术指标](docs/imgs/技术指标.png)

- **MACD** - 指数平滑异同平均线，识别趋势转折
- **RSI** - 相对强弱指数，判断超买超卖
- **布林带** - 价格波动区间分析
- **移动平均线** - MA5/MA10/MA20/MA60/MA120多周期均线
- **KDJ** - 随机指标，短期买卖信号

### 🎯 智能预测
- 基于多指标综合分析的买卖点预测
- 预测概率和置信度计算
- 详细的预测理由说明
- 支持多时间周期预测
- 预测信号回测功能，显示历史预测准确性

### 🚀 RESTful API
- 完整的REST API接口
- JSON格式数据交换
- CORS支持，便于前端集成
- 详细的错误处理和日志记录

### 🌐 Web界面
- **专业K线图**: 使用ECharts显示完整的OHLC数据
- **技术指标叠加**: MA5/MA10/MA20移动平均线
- **成交量副图**: 底部显示成交量柱状图，颜色与K线联动
- **智能搜索**: 支持股票名称和代码实时搜索
- **股票代码缓存**: 自动保存用户选择的股票，下次访问时自动恢复
- **响应式设计**: 自适应桌面端和移动端
- **交互体验**: 缩放、平移、十字光标等专业功能
- **收藏管理**: 支持股票收藏、分组管理和拖拽排序

### 🛠️ 开发工具
- **内置Curl工具**: Go语言版本的curl工具，支持Windows环境下的API调试
- **多平台支持**: 提供批处理和PowerShell脚本，简化编译和测试流程
- **数据库管理**: 支持SQLite数据库，提供数据迁移、备份和恢复工具

## 快速开始

### 环境要求
- **Go**: 1.24.0 或更高版本
- **Python**: 3.7+ (用于AKTools服务)
- **内存**: 至少2GB可用内存
- **磁盘**: 至少1GB可用空间

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/yourusername/stock-a-future.git
   cd stock-a-future
   ```

2. **安装Go依赖**
   ```bash
   go mod download
   go mod tidy
   ```

3. **配置环境变量**
   ```bash
   # 复制配置文件模板
   cp config.env.example config.env
   
   # 编辑配置文件，设置数据源类型
   # 推荐使用AKTools（免费，无需API密钥）
   ```

4. **启动AKTools服务**
   ```bash
   # 安装AKTools
   pip install aktools
   
   # 启动服务
   python -m aktools --port 8080
   
   # 后台运行（Linux/Mac）
   nohup python -m aktools > aktools.log 2>&1 &
   
   # Windows后台运行
   start /B python -m aktools
   
   # 使用项目提供的启动脚本（推荐）
   # Windows批处理
   start-aktools-server.bat
   
   # PowerShell脚本
   .\start-aktools-server.ps1
   ```

5. **启动Stock-A-Future服务**
   ```bash
   # 开发模式运行
   go run cmd/server/main.go
   
   # 或使用Makefile
   make dev
   
   # 或构建后运行
   make build
   make run
   ```

6. **使用Web界面**
   - 启动服务器后，直接在浏览器访问 `http://localhost:8081/`
   - Web界面已集成到服务器中，无需单独打开HTML文件
   - 使用智能搜索框输入股票名称或代码（如：平安银行、000001）
   - 选择股票后查看专业K线图和技术指标

7. **代码质量检查**
   ```bash
   # 快速代码质量检查
   golangci-lint run ./...
   
   # 格式化代码
   go fmt ./...
   
   # 运行测试
   go test -v ./...
   
   # 生成测试覆盖率报告
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```

### 快速命令参考

| 命令 | 说明 | 用途 |
|------|------|------|
| `make dev` | 开发模式运行 | 快速启动开发服务器 |
| `make build` | 构建应用 | 生成可执行文件 |
| `make test` | 运行测试 | 验证代码功能 |
| `make lint` | 代码质量检查 | 检查代码规范 |
| `make quality` | 完整质量检查 | 格式化+检查+测试 |
| `make fmt` | 格式化代码 | 统一代码风格 |
| `make imports` | 整理导入 | 优化导入顺序 |
| `make stop` | 停止服务 | 停止运行中的服务器 |
| `make status` | 查看状态 | 检查服务运行状态 |

### 验证安装

1. **检查AKTools服务**
   ```bash
   curl http://127.0.0.1:8080/api/public/stock_zh_a_info?symbol=000001
   ```

2. **检查Stock-A-Future服务**
   ```bash
   curl http://localhost:8081/api/v1/health
   ```

3. **访问Web界面**
   - 浏览器访问: `http://localhost:8081/`
   - 搜索股票: 输入"平安银行"或"000001"
   - 查看图表: 选择股票后查看K线图

## API文档

### 使用内置Curl工具测试API

项目内置了Go语言版本的curl工具，特别适合在Windows环境下调试API接口：

```bash
# 编译curl工具
build-curl.bat

# 测试基本API
curl.exe http://localhost:8081/api/v1/stocks

# 搜索股票
curl.exe -X POST -d '{"query":"平安"}' http://localhost:8081/api/v1/stocks/search

# 获取技术指标（详细输出模式）
curl.exe -v http://localhost:8081/api/v1/stocks/000001.SZ/indicators

# 添加请求头
curl.exe -H "Authorization: Bearer token123" http://localhost:8081/api/v1/stocks
```

更多curl工具使用方法，请参考 [cmd/curl/README.md](cmd/curl/README.md)

### 基础信息
- **Base URL**: `http://localhost:8081`
- **Content-Type**: `application/json`

### 接口列表

#### 1. 健康检查
```http
GET /api/v1/health
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-15T10:30:00Z",
    "version": "1.0.0",
    "services": {
      "tushare": "healthy"
    }
  }
}
```

#### 2. 获取股票日线数据
```http
GET /api/v1/stocks/{code}/daily?start_date=20240101&end_date=20240131&adjust=qfq
```

**参数说明**:
- `code`: 股票代码 (如: 000001.SZ, 600000.SH)
- `start_date`: 开始日期 (YYYYMMDD格式，可选，默认30天前)
- `end_date`: 结束日期 (YYYYMMDD格式，可选，默认今天)
- `adjust`: 复权方式 (可选，默认qfq)
  - `qfq`: 前复权 - 以当前价格为基准，向前调整历史价格
  - `hfq`: 后复权 - 以历史价格为基准，向后调整当前价格  
  - `none`: 不复权 - 原始价格数据

**复权说明**:
- **前复权(qfq)**: 推荐使用，价格连续性更好，便于技术分析
- **后复权(hfq)**: 历史价格真实，但当前价格可能偏离实际交易价格
- **不复权(none)**: 原始数据，包含除权除息缺口

**响应示例**:
```json
{
  "success": true,
  "data": [
    {
      "ts_code": "000001.SZ",
      "trade_date": "2024-01-15T00:00:00.000",
      "open": 8.75,
      "high": 8.85,
      "low": 8.69,
      "close": 8.70,
      "pre_close": 8.72,
      "change": -0.02,
      "pct_chg": -0.23,
      "vol": 525152.77,
      "amount": 460697.377
    }
  ]
}
```

#### 3. 获取技术指标
```http
GET /api/v1/stocks/{code}/indicators
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "ts_code": "000001.SZ",
    "trade_date": "20240115",
    "macd": {
      "dif": 0.05,
      "dea": 0.03,
      "macd": 0.04,
      "signal": "BUY"
    },
    "rsi": {
      "rsi6": 45.2,
      "rsi12": 48.5,
      "rsi24": 52.1,
      "signal": "HOLD"
    },
    "boll": {
      "upper": 9.20,
      "middle": 8.70,
      "lower": 8.20,
      "signal": "HOLD"
    }
  }
}
```

#### 4. 获取股票基本信息
```http
GET /api/v1/stocks/{code}/basic
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "ts_code": "000001.SZ",
    "symbol": "000001",
    "name": "平安银行",
    "area": "深圳",
    "industry": "银行",
    "market": "SZ",
    "list_date": "19910403"
  }
}
```

#### 5. 搜索股票
```http
GET /api/v1/stocks/search?q={keyword}&limit={limit}
```

**参数说明**:
- `q`: 搜索关键词 (股票名称或代码)
- `limit`: 返回结果数量限制 (可选，默认10，最大50)

**响应示例**:
```json
{
  "success": true,
  "data": {
    "keyword": "平安",
    "total": 3,
    "stocks": [
      {
        "ts_code": "000001.SZ",
        "symbol": "000001",
        "name": "平安银行",
        "market": "SZ"
      },
      {
        "ts_code": "601318.SH",
        "symbol": "601318",
        "name": "中国平安",
        "market": "SH"
      }
    ]
  }
}
```

#### 6. 获取股票列表
```http
GET /api/v1/stocks
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "total": 8950,
    "stocks": [
      {
        "ts_code": "000001.SZ",
        "symbol": "000001",
        "name": "平安银行",
        "market": "SZ"
      }
    ]
  }
}
```

#### 7. 获取买卖点预测
```http
GET /api/v1/stocks/{code}/predictions
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "ts_code": "000001.SZ",
    "trade_date": "2024-01-15",
    "predictions": [
      {
        "type": "BUY",
        "price": 8.70,
        "date": "2024-01-16",
        "probability": 0.65,
        "reason": "MACD金叉信号，DIF线上穿DEA线",
        "indicators": ["MACD"],
        "signal_date": "2024-01-15",
        "backtested": true,
        "is_correct": true,
        "next_day_price": 8.85,
        "price_diff": 0.15,
        "price_diff_ratio": 1.72
      }
    ],
    "confidence": 0.68,
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 复权选择建议

#### 不同场景的复权选择
- **技术分析**: 推荐使用前复权(`adjust=qfq`)，价格连续性更好
- **历史研究**: 推荐使用后复权(`adjust=hfq`)，历史价格真实
- **实时交易**: 推荐使用不复权(`adjust=none`)，当前价格准确

#### 复权参数使用示例
```bash
# 获取前复权数据（推荐用于技术分析）
curl "http://localhost:8081/api/v1/stocks/000001.SZ/daily?adjust=qfq"

# 获取后复权数据（历史价格真实）
curl "http://localhost:8081/api/v1/stocks/000001.SZ/daily?adjust=hfq"

# 获取不复权数据（原始价格）
curl "http://localhost:8081/api/v1/stocks/000001.SZ/daily?adjust=none"

# 不指定复权参数时，默认使用前复权
curl "http://localhost:8081/api/v1/stocks/000001.SZ/daily"
```

#### 常见问题
**Q: 为什么同一只股票在不同复权方式下价格差异很大？**
A: 这是因为股票在除权除息后，不同复权方式对历史价格的调整方式不同：
- 前复权：以当前价格为基准，向前调整历史价格
- 后复权：以历史价格为基准，向后调整当前价格
- 不复权：保持原始价格，但存在价格跳跃

**Q: 哪种复权方式最适合K线图分析？**
A: 推荐使用前复权(`qfq`)，因为：
1. 价格连续性最好，便于识别趋势
2. 技术指标计算更准确
3. 符合大多数交易软件的习惯

## 技术架构

### 项目结构
```
stock-a-future/
├── cmd/
│   ├── server/          # 主应用程序入口
│   └── stocklist/       # 股票列表获取工具
├── config/              # 配置管理
├── data/                # 数据文件
│   └── A股列表.xlsx      # 深交所股票列表（Excel格式）
├── internal/            # 内部包
│   ├── client/          # API客户端（Tushare + 交易所）
│   ├── handler/         # HTTP处理器
│   ├── indicators/      # 技术指标计算
│   ├── models/          # 数据模型
│   └── service/         # 业务逻辑服务
├── web/                 # Web资源
│   └── static/          # 静态文件（HTML、CSS、JS）
├── docs/                # 项目文档
├── examples/            # 使用示例
├── Makefile            # 构建脚本
└── README.md           # 项目文档
```

### 技术栈
- **语言**: Go 1.22
- **HTTP框架**: 标准库 net/http + ServeMux
- **数据源**: Tushare Pro API + 交易所官网 + 本地Excel文件
- **数值计算**: shopspring/decimal
- **配置管理**: godotenv

### 核心算法

#### 买卖点预测逻辑
1. **多指标综合分析**: 结合MACD、RSI、布林带、KDJ、移动平均线
2. **概率计算**: 基于指标强度和历史表现计算预测概率
3. **置信度评估**: 根据信号一致性评估整体置信度
4. **风险控制**: 设置概率阈值，过滤低质量信号

#### 技术指标实现
- **精确计算**: 使用decimal库确保金融计算精度
- **标准算法**: 严格按照技术分析标准公式实现
- **性能优化**: 高效的滑动窗口算法

## 开发指南

### 本地开发

#### 开发工具和代码质量
```bash
# 安装开发工具
make tools

# 代码格式化
make fmt

# 代码检查
make vet

# 运行测试
make test

# 代码质量检查
make lint
```

#### 服务器管理
```bash
# 检查服务器状态
make status

# 开发模式启动
make dev

# 停止服务器
make stop

# 强制停止服务器
make kill

# 重启服务器
make restart
```

#### 构建和部署
```bash
# 下载依赖
make deps

# 构建应用程序
make build

# 构建并运行
make run

# 清理构建文件
make clean
```

#### 数据管理
```bash
# 获取上交所股票列表
make fetch-sse

# 获取所有股票列表
make fetch-stocks

# 构建股票列表工具
make stocklist
```

### 环境配置
创建`.env`文件：
```bash
TUSHARE_TOKEN=your_tushare_token_here
TUSHARE_BASE_URL=http://api.tushare.pro
SERVER_PORT=8081
SERVER_HOST=localhost
LOG_LEVEL=info
```

### 部署
```bash
# 构建生产版本
make build

# 运行
./bin/stock-a-future
```

## 功能展示

### 🖥️ Web界面特性

#### K线图升级

![日K线信息](docs/imgs/日k线信息.png)

- **从简单折线图到专业K线图**: 显示完整的开盘、最高、最低、收盘价
- **红绿涨跌色彩**: 红色阳线表示上涨，绿色阴线表示下跌
- **成交量联动**: 底部成交量柱状图，颜色与K线保持一致
- **技术指标叠加**: 自动计算并显示MA5、MA10、MA20移动平均线

#### 智能搜索功能

![搜索功能视图](docs/imgs/搜索功能视图.png)

- **实时搜索**: 输入股票名称或代码，300ms防抖搜索
- **模糊匹配**: 支持部分匹配，如输入"平安"可找到"平安银行"、"中国平安"
- **键盘导航**: 支持上下箭头键选择，回车确认
- **自动填充**: 选择搜索结果后自动填入股票代码框
- **股票代码缓存**: 自动保存用户选择的股票，下次访问时自动恢复，提升用户体验

#### 股票收藏功能

![股票收藏](docs/imgs/股票收藏.png)

- **收藏管理**: 支持添加和删除股票收藏，便于快速访问常关注的股票
- **分组管理**: 支持将收藏的股票按不同分组进行管理
- **快速切换**: 一键切换至收藏的股票进行查看和分析
- **数据库存储**: 使用SQLite数据库存储，支持数据迁移、备份和恢复
- **拖拽排序**: 支持拖拽操作调整收藏顺序和分组

#### 数据清理功能

- **自动清理过期数据**: 定期清理超过指定天数的历史股票信号数据
- **智能清理策略**: 只清理股票信号数据，保护用户收藏信息
- **可配置保留期**: 支持自定义数据保留天数（默认90天）
- **手动清理支持**: 提供API接口手动触发清理任务
- **实时监控**: 查看清理服务状态和数据库统计信息
- **优雅关闭**: 服务器关闭时自动停止清理服务

#### 启用数据清理
在 `config.env` 文件中添加：
```bash
CLEANUP_ENABLED=true           # 启用数据清理服务
CLEANUP_INTERVAL=24h           # 清理间隔（默认24小时）
CLEANUP_RETENTION_DAYS=90      # 股票信号数据保留天数
```

#### 数据清理API
- `GET /api/v1/cleanup/status` - 查看清理服务状态
- `POST /api/v1/cleanup/manual` - 手动触发清理
- `PUT /api/v1/cleanup/config` - 更新清理配置

#### 交互体验
- **图表缩放**: 鼠标滚轮缩放，拖拽平移
- **数据提示**: 鼠标悬停显示详细的OHLC数据、成交量、涨跌幅
- **响应式设计**: 自适应不同屏幕尺寸
- **数据摘要**: 显示8个关键指标（收盘价、成交量、振幅等）

### 🔧 服务器管理

新增的Make命令让服务器管理更加便捷：

```bash
# 检查服务器状态（显示进程和端口占用）
make status

# 优雅停止服务器
make stop

# 强制停止（包括端口清理）
make kill

# 一键重启
make restart
```

## 使用示例

### cURL示例

#### 基础API测试
```bash
# 健康检查
curl http://localhost:8081/api/v1/health

# 获取平安银行日线数据
curl "http://localhost:8081/api/v1/stocks/000001.SZ/daily?start_date=20240101&end_date=20240131"

# 获取股票基本信息
curl http://localhost:8081/api/v1/stocks/000001.SZ/basic

# 搜索股票
curl "http://localhost:8081/api/v1/stocks/search?q=平安&limit=5"

# 获取股票列表
curl http://localhost:8081/api/v1/stocks

# 获取技术指标
curl http://localhost:8081/api/v1/stocks/000001.SZ/indicators

# 获取买卖点预测
curl http://localhost:8081/api/v1/stocks/000001.SZ/predictions

# 数据清理管理（需要启用清理服务）
curl http://localhost:8081/api/v1/cleanup/status                    # 查看清理服务状态
curl -X POST http://localhost:8081/api/v1/cleanup/manual           # 手动触发清理
curl -X PUT http://localhost:8081/api/v1/cleanup/config \          # 更新清理配置
  -H "Content-Type: application/json" \
  -d '{"retention_days": 60}'
```

#### AKTools快速测试
```bash
# 测试AKTools服务连接
curl http://127.0.0.1:8080/api/public/stock_zh_a_info?symbol=000001

# 测试Stock-A-Future API（使用AKTools数据源）
curl http://localhost:8081/api/v1/stocks/000001/basic

# 获取股票日线数据
curl http://localhost:8081/api/v1/stocks/000001/daily
```

### Python示例
```python
import requests

# 基础配置
base_url = "http://localhost:8081"
stock_code = "000001.SZ"

# 1. 搜索股票
def search_stocks(keyword):
    response = requests.get(f"{base_url}/api/v1/stocks/search", 
                          params={"q": keyword, "limit": 5})
    data = response.json()
    if data["success"]:
        print(f"搜索 '{keyword}' 的结果:")
        for stock in data["data"]["stocks"]:
            print(f"  {stock['name']} ({stock['ts_code']}) - {stock['market']}")
        return data["data"]["stocks"]
    return []

# 2. 获取股票基本信息
def get_stock_basic(stock_code):
    response = requests.get(f"{base_url}/api/v1/stocks/{stock_code}/basic")
    data = response.json()
    if data["success"]:
        stock = data["data"]
        print(f"股票信息: {stock['name']} ({stock['ts_code']})")
        print(f"所属市场: {stock['market']}, 行业: {stock.get('industry', 'N/A')}")
        return stock
    return None

# 3. 获取日线数据
def get_daily_data(stock_code, start_date="20250101", end_date="20250131"):
    response = requests.get(f"{base_url}/api/v1/stocks/{stock_code}/daily",
                          params={"start_date": start_date, "end_date": end_date})
    data = response.json()
    if data["success"]:
        daily_data = data["data"]
        print(f"获取到 {len(daily_data)} 条日线数据")
        if daily_data:
            latest = daily_data[-1]
            print(f"最新数据 ({latest['trade_date']}): 收盘价 {latest['close']}")
        return daily_data
    return []

# 4. 获取预测结果
def get_predictions(stock_code):
    response = requests.get(f"{base_url}/api/v1/stocks/{stock_code}/predictions")
    data = response.json()
    if data["success"]:
        predictions = data["data"]["predictions"]
        confidence = data["data"]["confidence"]
        print(f"预测置信度: {confidence:.2%}")
        for pred in predictions:
            print(f"预测类型: {pred['type']}")
            print(f"预测价格: {pred['price']}")
            print(f"预测概率: {pred['probability']:.2%}")
            print(f"预测理由: {pred['reason']}")
            
            # 显示回测结果
            if pred.get("backtested"):
                result = "✅ 正确" if pred.get("is_correct") else "❌ 错误"
                print(f"回测结果: {result}")
                print(f"次日价格: {pred.get('next_day_price')}")
                
                # 格式化价格差异
                price_diff = pred.get("price_diff", 0)
                price_diff_ratio = pred.get("price_diff_ratio", 0)
                sign = "+" if price_diff >= 0 else ""
                print(f"价格差异: {sign}{price_diff:.2f} ({sign}{price_diff_ratio:.2f}%)")
            
            print("---")
        return predictions
    return []

# 使用示例
if __name__ == "__main__":
    # 搜索包含"平安"的股票
    stocks = search_stocks("平安")
    
    if stocks:
        # 使用第一个搜索结果
        stock_code = stocks[0]["ts_code"]
        
        # 获取基本信息
        get_stock_basic(stock_code)
        
        # 获取日线数据
        get_daily_data(stock_code)
        
        # 获取预测结果
        get_predictions(stock_code)
```

## 注意事项

### Tushare使用限制
- 需要注册Tushare Pro账号获取Token
- 免费账号有调用频率限制
- 部分高级数据需要积分

### 风险提示
- 本系统仅供学习和研究使用
- 预测结果不构成投资建议
- 股市有风险，投资需谨慎
- 请根据自身情况做出投资决策

### 性能考虑
- 技术指标计算需要足够的历史数据
- 建议为计算密集型操作添加缓存
- 生产环境建议使用数据库存储历史数据
- 收藏功能已升级为SQLite数据库存储，提供更好的性能和可靠性

## 重要配置说明

### 开始使用前的必要步骤

**⚠️ 需要注意的是，您需要：**

1. **获取Tushare API Token**
   - 访问 [Tushare官网](https://tushare.pro) 注册账号
   - 在个人中心获取您的API Token
   - 免费账号有一定的调用限制，请合理使用

2. **配置环境变量**
   ```bash
   # 复制配置文件模板
   cp .env.example .env
   
   # 编辑.env文件，将your_tushare_token_here替换为您的真实Token
   vim .env
   ```

3. **验证配置**
   ```bash
   # 启动服务
   make dev
   
   # 在另一个终端测试健康检查
   curl http://localhost:8081/api/v1/health
   ```

如果健康检查显示Tushare服务状态为"healthy"，说明配置成功。

## 📈 更新日志

### 2025-08-20

#### 🆕 新增功能
- **预测信号回测功能**: 新增对历史预测信号的回测，显示预测准确性和价格差异
  - 自动判断预测是否正确（买入信号后涨/卖出信号后跌）
  - 计算预测价格与实际价格的差值和百分比
  - 直观显示回测结果（正确/错误）和价格变化

### 2025-08-13

#### 🆕 新增功能
- **多数据源支持**: 新增AKTools数据源集成，支持免费开源财经数据
- **专业K线图**: 升级前端图表为ECharts，支持完整OHLC显示
- **智能股票搜索**: 新增股票名称和代码搜索API和前端界面
- **服务器管理命令**: 新增 `make stop/kill/status/restart` 命令
- **成交量副图**: K线图下方显示成交量柱状图
- **技术指标叠加**: 自动显示MA5/MA10/MA20移动平均线
- **复权参数支持**: 新增`adjust`参数，支持前复权(qfq)、后复权(hfq)、不复权(none)选择
- **数据库存储**: 收藏功能升级为SQLite数据库存储，支持数据迁移、备份和恢复
- **数据清理功能**: 新增自动数据清理服务，定期清理过期股票信号数据，保护用户收藏信息

#### 🔧 改进优化
- **端口更新**: 默认端口从8080改为8081，避免常见冲突
- **数据摘要增强**: 显示8个关键指标（成交额、振幅等）
- **交互体验**: 支持图表缩放、平移、键盘导航
- **响应式设计**: 优化移动端显示效果
- **错误处理**: 改进API错误处理和用户反馈
- **默认复权方式**: 默认使用前复权(qfq)，更适合技术分析

#### 🐛 修复问题
- 修复图表数据格式化问题
- 修复K线图日期显示格式错误（从"2025--0-7-"修复为"2025-07-17"）
- 优化搜索性能和防抖处理
- 改进服务器进程管理和端口检测
- 修复AKTools数据源日期格式兼容性问题

#### 📚 文档更新
- 更新所有API示例和端口号
- 新增服务器管理指南
- 完善Python使用示例
- 添加功能展示说明
- 新增复权选择指南和常见问题解答
- 完善AKTools启动和配置说明

### 2025-08-12

#### 🆕 新增功能
- **本地股票数据服务**: 支持从Excel文件读取深交所股票数据
- **混合数据源架构**: 上交所在线获取 + 深交所本地Excel文件
- **股票列表获取工具**: 支持多种数据源选择（sse/all）
- **智能列映射**: 自动识别Excel文件中的股票代码和名称列

#### 🔧 改进优化
- **智能端口检测**: 自动尝试常见端口（8081, 8080, 3000, 8000, 9000）
- **配置管理界面**: 前端配置按钮，支持手动设置服务器地址
- **自动故障转移**: 连接失败时自动尝试其他端口
- **数据排序优化**: 修复日线数据时间轴显示问题
- **进程管理增强**: 改进make命令的进程识别和端口检测

#### 🐛 修复问题
- 修复前端日期解析"Invalid Date"问题
- 修复JSON序列化中decimal类型问题
- 修复日线数据时间排序问题（从左到右正确增长）
- 优化端口占用检测和进程清理

#### 📚 文档更新
- 新增股票列表获取工具文档
- 更新项目结构和技术栈说明
- 完善服务器管理命令说明

### 2025-08-12 (早期)

#### 🆕 新增功能
- **交易所客户端**: 支持获取上交所和深交所股票列表
- **命令行工具**: 支持多种数据源选择和配置
- **Makefile集成**: 新增股票列表相关构建命令

#### 🔧 改进优化
- **智能配置读取**: 自动从.env文件读取配置
- **命令行参数支持**: 支持端口和主机地址配置
- **测试脚本优化**: 自动适配配置文件中的端口设置

### 初始版本
- 初始版本发布
- 基础API功能
- 技术指标计算
- 简单前端界面

## 贡献指南

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 代码质量检查

本项目使用 `golangci-lint` 进行代码质量检查，确保代码符合Go语言最佳实践。

### 安装golangci-lint

```bash
# 安装最新版本
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 验证安装
golangci-lint --version
```

### 运行代码检查

```bash
# 检查整个项目
golangci-lint run ./...

# 检查特定包
golangci-lint run ./internal/service

# 检查特定文件
golangci-lint run internal/service/signal.go

# 生成详细报告
golangci-lint run --out-format=html > lint-report.html
```

### 自动修复

```bash
# 自动修复可修复的问题
golangci-lint run --fix ./...

# 格式化代码
go fmt ./...

# 整理导入
goimports -w .
```

### 配置说明

项目根目录包含 `.golangci.yml` 配置文件，主要设置：

- **函数长度限制**: 最大80行，40条语句
- **函数复杂度限制**: 最大圈复杂度10
- **启用的检查器**: 包括代码质量、性能、安全等40+个检查器
- **排除规则**: 测试文件和vendor目录

### 常见问题解决

**Q: 如何忽略特定的lint警告？**
A: 在代码中添加注释：
```go
//nolint:gocyclo
func complexFunction() {
    // 复杂逻辑
}
```

**Q: 如何临时禁用某个检查器？**
A: 修改 `.golangci.yml` 配置文件，在 `disable` 部分添加检查器名称。

**Q: 如何查看检查器的详细说明？**
A: 运行 `golangci-lint help linters` 查看所有检查器的说明。

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

- 项目地址: [GitHub Repository]
- 问题反馈: [Issues]
- 文档: [Wiki]

---

**免责声明**: 本项目仅用于技术学习和研究目的，不构成任何投资建议。使用者应当自行承担投资风险。

## 常见问题解决

### 复权相关问题

**Q: 为什么K线图显示的股票价格与实际交易价格不符？**
A: 这是因为系统默认使用前复权数据。不同复权方式下的价格差异是正常的：
- **前复权(qfq)**: 以当前价格为基准调整历史价格，适合技术分析
- **后复权(hfq)**: 以历史价格为基准调整当前价格，历史价格真实
- **不复权(none)**: 原始价格，包含除权除息缺口

**解决方案**: 在API调用时指定复权参数：
```bash
# 获取后复权数据（历史价格真实）
curl "http://localhost:8081/api/v1/stocks/601225.SH/daily?adjust=hfq"

# 获取不复权数据（当前价格准确）
curl "http://localhost:8081/api/v1/stocks/601225.SH/daily?adjust=none"
```

### 日期显示问题

**Q: K线图X轴日期显示为"2025--0-7-"这样的错误格式？**
A: 这个问题已经修复。现在系统正确处理ISO 8601格式的日期字符串，自动转换为"2025-07-17"格式。

### AKTools服务问题

**Q: 启动Stock-A-Future后无法获取股票数据？**
A: 请确保AKTools服务正在运行：
```bash
# 检查AKTools服务状态
curl http://127.0.0.1:8080/api/public/stock_zh_a_info?symbol=000001

# 如果服务未启动，请启动AKTools
python -m aktools
```

**Q: AKTools服务启动失败？**
A: 检查Python环境和依赖：
```bash
# 检查Python版本（需要3.7+）
python --version

# 重新安装AKTools
pip uninstall aktools
pip install aktools

# 或者从源码安装
pip install git+https://github.com/akfamily/aktools.git
```

### 端口冲突问题

**Q: 启动服务时提示端口被占用？**
A: 使用项目提供的端口管理命令：
```bash
# 查看端口占用情况
make status

# 停止占用端口的进程
make stop

# 或者指定其他端口启动
go run cmd/server/main.go --port 8082
```

## 技术支持

如果您遇到其他问题，请：
1. 查看项目文档和更新日志
2. 搜索GitHub Issues中的类似问题
3. 创建新的Issue并提供详细的错误信息和环境描述