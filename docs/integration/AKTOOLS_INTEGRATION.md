# AKTools 集成使用说明

## 概述

Stock-A-Future 现已支持 AKTools 数据源，这是一个基于 AKShare 的免费开源财经数据接口。通过 AKTools，您可以获取 A 股、港股、美股等多种市场数据，而无需 Python 环境。

## 优势

- **免费开源**: 基于 AKShare，完全免费使用
- **数据丰富**: 支持多种市场和数据类型
- **易于部署**: 一行命令即可启动 HTTP 服务
- **语言无关**: 通过 HTTP API 调用，支持任何编程语言

## 安装和启动 AKTools

### 1. 安装 AKTools

```bash
pip install aktools
```

### 2. 启动 AKTools 服务

```bash
python -m aktools
```

默认情况下，AKTools 将在 `http://127.0.0.1:8080` 启动。

## 配置 Stock-A-Future

### 1. 环境变量配置

在 `.env` 文件中设置以下配置：

```bash
# 数据源类型
DATA_SOURCE_TYPE=aktools

# AKTools 服务地址
AKTOOLS_BASE_URL=http://127.0.0.1:8080

# 其他配置保持不变
SERVER_PORT=8080
SERVER_HOST=localhost
LOG_LEVEL=info
```

### 2. 配置说明

- `DATA_SOURCE_TYPE`: 设置为 `aktools` 使用 AKTools 数据源
- `AKTOOLS_BASE_URL`: AKTools 服务的地址和端口
- 当使用 AKTools 时，`TUSHARE_TOKEN` 不是必需的

## 测试连接

### 1. 构建测试工具

```bash
make aktools-test
```

### 2. 测试连接

```bash
make test-aktools
```

或者直接运行：

```bash
./bin/aktools-test
```

### 3. 自定义服务地址测试

```bash
./bin/aktools-test -url http://your-server:8080
```

## 支持的数据接口

### 1. 股票日线数据

- **接口**: `/api/public/stock_zh_a_hist`
- **参数**: 
  - `symbol`: 股票代码（如：000001）
  - `start_date`: 开始日期（可选）
  - `end_date`: 结束日期（可选）
  - `period`: 周期（默认：daily）
  - `adjust`: 复权类型（默认：hfq）

### 2. 股票基本信息

- **接口**: `/api/public/stock_zh_a_info`
- **参数**: 
  - `symbol`: 股票代码

### 3. 股票列表

- **接口**: `/api/public/stock_zh_a_spot`
- **参数**: 无

## 数据格式

### 日线数据字段

```json
{
  "日期": "2024-12-01",
  "开盘": 10.50,
  "收盘": 10.80,
  "最高": 11.00,
  "最低": 10.40,
  "成交量": 1000000,
  "成交额": 10800000,
  "振幅": 5.71,
  "涨跌幅": 2.86,
  "涨跌额": 0.30,
  "换手率": 1.23
}
```

### 股票基本信息字段

```json
{
  "代码": "000001",
  "名称": "平安银行",
  "地区": "深圳",
  "行业": "银行",
  "市场": "主板",
  "上市日期": "1991-04-03"
}
```

## 切换数据源

### 1. 运行时切换

通过环境变量动态切换：

```bash
export DATA_SOURCE_TYPE=aktools
export AKTOOLS_BASE_URL=http://127.0.0.1:8080
./bin/stock-a-future
```

### 2. 配置文件切换

修改 `.env` 文件中的 `DATA_SOURCE_TYPE` 值：

```bash
# 使用 AKTools
DATA_SOURCE_TYPE=aktools

# 使用 Tushare
DATA_SOURCE_TYPE=tushare
```

## 故障排除

### 1. 连接失败

**问题**: 无法连接到 AKTools 服务

**解决方案**:
- 确保 AKTools 服务正在运行
- 检查服务地址和端口是否正确
- 验证防火墙设置

### 2. 数据获取失败

**问题**: 能够连接但无法获取数据

**解决方案**:
- 检查股票代码格式是否正确
- 验证日期格式（YYYYMMDD）
- 查看 AKTools 服务日志

### 3. 性能问题

**问题**: 数据获取速度较慢

**解决方案**:
- 使用本地部署的 AKTools 服务
- 优化网络连接
- 考虑使用缓存机制

## 部署建议

### 1. 生产环境

- 将 AKTools 部署在专用服务器上
- 配置反向代理（如 Nginx）
- 设置适当的访问控制

### 2. 开发环境

- 使用本地部署的 AKTools 服务
- 配置开发环境的端口映射
- 启用详细的日志记录

## 与 Tushare 的对比

| 特性 | AKTools | Tushare |
|------|---------|---------|
| 费用 | 免费 | 付费 |
| 部署 | 自部署 | 云端服务 |
| 数据量 | 中等 | 丰富 |
| 稳定性 | 依赖自部署 | 高 |
| 支持 | 社区支持 | 官方支持 |

## 总结

AKTools 为 Stock-A-Future 提供了一个优秀的免费数据源选择。它特别适合：

- 个人投资者和小型团队
- 对数据成本敏感的用户
- 需要自定义部署的用户
- 学习和研究用途

通过合理的配置和部署，AKTools 可以为您的股票分析应用提供稳定可靠的数据支持。
