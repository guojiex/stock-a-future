# Stock-A-Future 快速启动指南

## 🚀 一键启动方式

### 方式1: 双击批处理文件 (推荐新手)
- **`start-aktools-server.bat`** - 一键启动AKTools服务和Stock-A-Future服务器（推荐）
- **`start-server.bat`** - 完整启动脚本，包含环境检查和依赖管理
- **`quick-start.bat`** - 简单快速启动，适合有经验的用户

### 方式2: 在Cursor中运行任务（推荐）
1. 按 `Ctrl+Shift+P` 打开命令面板
2. 输入 "Tasks: Run Task"
3. 选择相应的任务：
   - **"一键启动 AKTools + Stock-A-Future"** - 推荐新手使用
   - **"启动 AKTools 服务"** - 仅启动数据服务
   - **"启动 Stock-A-Future 服务器"** - 仅启动API服务器

### 方式3: PowerShell脚本
- **`start-aktools-server.ps1`** - AKTools专用启动脚本
- **`start-server.ps1`** - 功能强大的PowerShell启动脚本
- 右键选择 "使用PowerShell运行"

## ⚠️ 重要：启动顺序说明

**Stock-A-Future服务器依赖于数据源服务，必须按以下顺序启动：**

### 1. 首先启动数据源服务
- **使用AKTools数据源**（推荐免费使用）：
  ```bash
  # 启动AKTools服务（端口8080）
  python -m aktools
  ```
- **或使用Tushare数据源**：
  - 需要配置TUSHARE_TOKEN
  - 服务会自动连接Tushare API

### 2. 然后启动Stock-A-Future服务器
```bash
# 启动Stock-A-Future API服务器（端口8081）
go run cmd/server/main.go
```

## 📋 启动前准备

### 1. 确保Go环境已安装
```bash
go version
# 需要Go 1.22或更高版本
```

### 2. 确保Python环境已安装（使用AKTools时）
```bash
python --version
# 需要Python 3.7或更高版本
```

### 3. 配置环境变量
- **使用AKTools数据源**：
  - 复制 `cache.env.example` 为 `config.env`
  - 设置 `DATA_SOURCE_TYPE=aktools`
  - 设置 `AKTOOLS_BASE_URL=http://127.0.0.1:8080`

- **使用Tushare数据源**：
  - 复制 `cache.env.example` 为 `.env`
  - 编辑 `.env` 文件，设置你的 `TUSHARE_TOKEN`

### 4. 安装依赖
```bash
# Go依赖
go mod tidy

# Python依赖（使用AKTools时）
pip install aktools
```

## 🔧 启动选项

### 推荐：一键启动（AKTools数据源）
```bash
# 使用批处理文件（推荐新手）
start-aktools-server.bat

# 或使用PowerShell
.\start-aktools-server.ps1

# 或在Cursor中运行任务
# 选择 "一键启动 AKTools + Stock-A-Future"
```

### 分步启动（AKTools数据源）
```bash
# 步骤1：启动AKTools服务
python -m aktools

# 步骤2：在新终端启动Stock-A-Future服务器
go run cmd/server/main.go
```

### 普通启动（Tushare数据源）
```bash
# 使用批处理文件
start-server.bat

# 或使用PowerShell
.\start-server.ps1
```

### 开发模式 (支持热重载)
```bash
# 安装Air工具
go install github.com/cosmtrek/air@latest

# 启动开发模式
.\start-server.ps1 -Dev
```

### 构建后启动
```bash
# 构建并启动
.\start-server.ps1 -Build
```

## 🌐 访问地址

启动成功后，可以通过以下地址访问：

### 使用AKTools数据源
- **AKTools服务**: http://localhost:8080
- **Stock-A-Future API**: http://localhost:8081
- **Web客户端**: http://localhost:8081/
- **健康检查**: http://localhost:8081/api/v1/health

### 使用Tushare数据源
- **API服务器**: http://localhost:8081
- **Web客户端**: http://localhost:8081/
- **健康检查**: http://localhost:8081/api/v1/health

## 📱 在Cursor中运行

### 使用任务面板
1. 按 `Ctrl+Shift+P`
2. 输入 "Tasks: Run Task"
3. 选择相应的任务

### 可用的任务

#### 🚀 一键启动任务（推荐）
- **一键启动 AKTools + Stock-A-Future** - 推荐新手使用，自动处理启动顺序

#### 🔧 分步启动任务
- **启动 AKTools 服务** - 启动数据源服务（端口8080）
- **启动 Stock-A-Future 服务器** - 启动API服务器（端口8081）
- **使用 AKTools 启动 Stock-A-Future** - 依赖启动，自动先启动AKTools

#### 🛠️ 环境检查和维护任务
- **检查Go环境** - 验证Go安装
- **检查 Python 环境** - 验证Python安装（AKTools需要）
- **安装 AKTools** - 安装Python包
- **更新依赖** - 更新Go模块依赖

#### 🧪 测试任务
- **测试 AKTools 连接** - 测试数据源服务
- **测试 Stock-A-Future API (AKTools)** - 测试API服务器（端口8081）
- **测试API** - 测试健康检查端点

#### 🔄 开发任务
- **启动服务器 (开发模式)** - 使用Air热重载
- **构建并启动服务器** - 先构建再启动

## 🛠️ 故障排除

### 常见问题

1. **Go未安装**
   - 下载安装Go: https://golang.org/dl/
   - 确保版本 >= 1.22

2. **Python未安装（使用AKTools时）**
   - 下载安装Python: https://www.python.org/downloads/
   - 确保版本 >= 3.7

3. **AKTools服务未启动**
   - 确保AKTools服务在端口8080运行
   - 检查：`curl http://127.0.0.1:8080/api/public/stock_zh_a_info?symbol=000001`

4. **环境变量未配置**
   - 检查配置文件是否存在
   - 确保数据源类型和配置正确

5. **依赖问题**
   - 运行 `go mod tidy`
   - 运行 `pip install aktools`
   - 检查 `go.mod` 文件

6. **端口被占用**
   - AKTools默认使用端口8080
   - Stock-A-Future默认使用端口8081
   - 修改配置文件中的端口设置

### 获取帮助
```bash
# PowerShell脚本帮助
.\start-server.ps1 -Help

# 检查Go环境
go version
go env

# 检查Python环境
python --version

# 测试AKTools连接
curl http://127.0.0.1:8080/api/public/stock_zh_a_info?symbol=000001
```

## 📝 开发建议

- **新手用户**: 使用 **"一键启动 AKTools + Stock-A-Future"** 任务
- **开发时**: 使用 **开发模式** (`-Dev` 参数)，支持代码热重载
- **生产环境**: 使用 **构建模式** (`-Build` 参数)，性能更好
- **定期维护**: 运行 `go mod tidy` 和 `pip install --upgrade aktools`
- **使用Cursor任务**: 可以更方便地管理开发流程和启动顺序

---

🎯 **推荐新手使用**: `start-aktools-server.bat` 或 Cursor任务面板中的 **"一键启动 AKTools + Stock-A-Future"**
⚡ **推荐开发使用**: `start-aktools-server.ps1` 或 Air开发模式
🚀 **推荐生产使用**: `start-aktools-server.ps1` 或 构建后运行
