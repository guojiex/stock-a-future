# Stock-A-Future 快速启动指南

## 🚀 一键启动方式

### 方式1: 双击批处理文件 (推荐新手)
- **`start-server.bat`** - 完整启动脚本，包含环境检查和依赖管理
- **`quick-start.bat`** - 简单快速启动，适合有经验的用户

### 方式2: 在Cursor中运行任务
1. 按 `Ctrl+Shift+P` 打开命令面板
2. 输入 "Tasks: Run Task"
3. 选择 "启动 Stock-A-Future 服务器"

### 方式3: PowerShell脚本
- **`start-server.ps1`** - 功能强大的PowerShell启动脚本
- 右键选择 "使用PowerShell运行"

## 📋 启动前准备

### 1. 确保Go环境已安装
```bash
go version
# 需要Go 1.22或更高版本
```

### 2. 配置环境变量
- 复制 `cache.env.example` 为 `.env`
- 编辑 `.env` 文件，设置你的 `TUSHARE_TOKEN`

### 3. 安装依赖 (可选)
```bash
go mod tidy
```

## 🔧 启动选项

### 普通启动
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

- **API服务器**: http://localhost:8080
- **Web客户端**: http://localhost:8080/
- **健康检查**: http://localhost:8080/api/v1/health

## 📱 在Cursor中运行

### 使用任务面板
1. 按 `Ctrl+Shift+P`
2. 输入 "Tasks: Run Task"
3. 选择相应的任务

### 可用的任务
- **启动 Stock-A-Future 服务器** - 默认启动方式
- **启动服务器 (开发模式)** - 使用Air热重载
- **构建并启动服务器** - 先构建再启动
- **检查Go环境** - 验证Go安装
- **更新依赖** - 更新Go模块
- **测试API** - 测试健康检查端点

## 🛠️ 故障排除

### 常见问题

1. **Go未安装**
   - 下载安装Go: https://golang.org/dl/
   - 确保版本 >= 1.22

2. **环境变量未配置**
   - 检查 `.env` 文件是否存在
   - 确保 `TUSHARE_TOKEN` 已设置

3. **依赖问题**
   - 运行 `go mod tidy`
   - 检查 `go.mod` 文件

4. **端口被占用**
   - 修改 `.env` 中的 `SERVER_PORT`
   - 或关闭占用端口的程序

### 获取帮助
```bash
# PowerShell脚本帮助
.\start-server.ps1 -Help

# 检查Go环境
go version
go env
```

## 📝 开发建议

- 开发时使用 **开发模式** (`-Dev` 参数)，支持代码热重载
- 生产环境使用 **构建模式** (`-Build` 参数)，性能更好
- 定期运行 `go mod tidy` 保持依赖更新
- 使用Cursor的任务功能可以更方便地管理开发流程

---

🎯 **推荐新手使用**: `start-server.bat` 或 Cursor任务面板
⚡ **推荐开发使用**: `start-server.ps1 -Dev` 或 Air开发模式
🚀 **推荐生产使用**: `start-server.ps1 -Build` 或 构建后运行
