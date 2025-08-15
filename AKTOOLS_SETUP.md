# AKTools 集成设置指南

## 概述

本指南将帮助您设置 Stock-A-Future 项目使用 AKTools 作为数据源。AKTools 是一个开源的 A 股数据工具，提供免费的股票数据 API。

## 前置要求

1. **Go 1.22+** - 确保已安装最新版本的 Go
2. **Python 3.7+** - 用于运行 AKTools
3. **AKTools** - 需要安装并运行 AKTools 服务

## 安装 AKTools

### 方法1: 使用 pip 安装
```bash
pip install aktools
```

### 方法2: 从源码安装
```bash
git clone https://github.com/akfamily/aktools.git
cd aktools
pip install -e .
```

## 启动 AKTools 服务

```bash
# 启动 AKTools 服务（默认端口 8080）
python -m aktools

# 或者指定端口
python -m aktools --port 8080
```

服务启动后，您应该看到类似输出：
```
INFO:     Started server process [xxxxx]
INFO:     Waiting for application startup.
INFO:     Application startup complete.
INFO:     Uvicorn running on http://127.0.0.1:8080
```

## 配置 Stock-A-Future

### 1. 创建配置文件

项目根目录已包含 `config.env` 文件，配置如下：

```env
# 数据源配置 - 使用AKTools
DATA_SOURCE_TYPE=aktools

# AKTools配置 - 运行在本地8080端口
AKTOOLS_BASE_URL=http://127.0.0.1:8080

# 服务器配置 - 使用不同的端口避免冲突
SERVER_PORT=8081
SERVER_HOST=localhost

# 缓存配置
CACHE_ENABLED=true
CACHE_DEFAULT_TTL=1h
CACHE_MAX_AGE=24h
CACHE_CLEANUP_INTERVAL=10m
```

### 2. 验证配置

运行连接测试工具：
```bash
# Windows
test-aktools-connection.bat

# 或者直接运行
go run cmd/aktools-test/main.go
```

如果测试成功，您将看到：
```
✓ AKTools API连接测试成功!
连接测试完成，可以启动Stock-A-Future服务器了
服务器将在 http://localhost:8081 启动
```

## 启动服务器

### 方法1: 使用启动脚本

**Windows:**
```bash
start-aktools-server.bat
```

**PowerShell:**
```powershell
.\start-aktools-server.ps1
```

### 方法2: 直接运行

```bash
go run cmd/server/main.go
```

## 验证服务

服务器启动后，您可以通过以下方式验证：

### 1. 健康检查
```bash
curl http://localhost:8081/api/v1/health
```

### 2. 获取股票基本信息
```bash
curl http://localhost:8081/api/v1/stocks/000001/basic
```

### 3. 获取股票日线数据
```bash
curl http://localhost:8081/api/v1/stocks/000001/daily
```

### 4. 访问 Web 客户端
在浏览器中打开：http://localhost:8081

## 故障排除

### 常见问题

#### 1. 端口冲突
**错误信息：** `bind: address already in use`
**解决方案：** 
- 确保 AKTools 运行在 8080 端口
- Stock-A-Future 运行在 8081 端口
- 检查 `config.env` 中的端口配置

#### 2. AKTools 连接失败
**错误信息：** `AKTools API连接测试失败: 404`
**解决方案：**
- 确保 AKTools 服务正在运行
- 检查 AKTools 服务端口
- 验证 `AKTOOLS_BASE_URL` 配置

#### 3. 数据获取失败
**错误信息：** `获取股票数据失败`
**解决方案：**
- 检查 AKTools 服务状态
- 验证股票代码格式
- 检查网络连接

### 调试步骤

1. **检查 AKTools 服务状态**
   ```bash
   curl http://127.0.0.1:8080/health
   ```

2. **检查 Stock-A-Future 服务状态**
   ```bash
   curl http://localhost:8081/api/v1/health
   ```

3. **查看服务器日志**
   启动服务器时会显示详细的连接和错误信息

4. **验证环境变量**
   ```bash
   # Windows
   echo %DATA_SOURCE_TYPE%
   echo %AKTOOLS_BASE_URL%
   
   # PowerShell
   $env:DATA_SOURCE_TYPE
   $env:AKTOOLS_BASE_URL
   ```

## API 端点

使用 AKTools 数据源时，以下 API 端点可用：

- `GET /api/v1/health` - 健康检查
- `GET /api/v1/stocks/{code}/basic` - 股票基本信息
- `GET /api/v1/stocks/{code}/daily` - 股票日线数据
- `GET /api/v1/stocks/{code}/indicators` - 技术指标
- `GET /api/v1/stocks/{code}/predictions` - 买卖点预测
- `GET /api/v1/stocks` - 本地股票列表
- `GET /api/v1/stocks/search` - 股票搜索
- `GET /api/v1/favorites` - 收藏股票管理

## 性能优化

1. **启用缓存** - 默认已启用，减少重复 API 调用
2. **调整缓存时间** - 根据数据更新频率调整 `CACHE_DEFAULT_TTL`
3. **批量请求** - 避免频繁的单次请求

## 安全注意事项

1. **本地部署** - AKTools 服务运行在本地，确保网络安全
2. **端口限制** - 避免在公网暴露服务端口
3. **数据验证** - 所有输入数据都会进行验证

## 支持

如果遇到问题，请：

1. 检查本文档的故障排除部分
2. 查看服务器启动日志
3. 运行连接测试工具
4. 检查 AKTools 官方文档

## 更新日志

- **v1.0.0** - 初始 AKTools 集成
- 支持股票基本信息、日线数据、技术指标等
- 完整的错误处理和日志记录
- 缓存机制优化
