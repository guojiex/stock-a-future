# 跨平台启动指南

本项目支持在Windows、macOS和Linux上运行。由于不同操作系统的环境变量设置方式不同，我们提供了多种启动方式。

## 🚀 VSCode Tasks（推荐）

已修复跨平台兼容性问题，现在可以在所有平台上使用：

1. 按 `Ctrl+Shift+P` (Windows/Linux) 或 `Cmd+Shift+P` (macOS)
2. 输入 "Tasks: Run Task"
3. 选择 "启动 Go API 服务器 (端口8081)"

## 🔧 命令行方式

### Windows

#### 方式1：使用批处理脚本
```cmd
scripts\start-go-server.bat
```

#### 方式2：使用PowerShell脚本
```powershell
.\scripts\start-go-server.ps1
```

#### 方式3：使用Makefile (需要安装make)
```cmd
make dev-win
```

#### 方式4：手动设置环境变量
```cmd
set SERVER_PORT=8081
go run cmd/server/main.go
```

#### 方式5：PowerShell环境变量
```powershell
$env:SERVER_PORT = "8081"
go run cmd/server/main.go
```

### macOS/Linux

#### 方式1：使用Shell脚本
```bash
./scripts/start-go-server.sh
```

#### 方式2：使用Makefile
```bash
make dev-8081
```

#### 方式3：直接命令行
```bash
SERVER_PORT=8081 go run cmd/server/main.go
```

## 🔍 验证服务器启动

服务器启动后，访问以下URL验证：
- API健康检查: http://localhost:8081/health
- API文档: http://localhost:8081/api/docs

## ⚠️ 常见问题

### Windows PowerShell执行策略
如果PowerShell脚本无法执行，运行：
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 端口占用
如果端口被占用，可以：
1. 更改端口号（修改SERVER_PORT环境变量）
2. 停止占用端口的进程：
   ```bash
   # Unix/Linux/macOS
   make kill
   
   # Windows
   netstat -ano | findstr :8081
   taskkill /PID <PID> /F
   ```

## 📝 开发建议

- **推荐使用VSCode Tasks**：最简单，跨平台兼容
- **Windows用户**：优先使用PowerShell脚本或批处理脚本
- **Unix用户**：可以使用Shell脚本或Makefile
- **CI/CD环境**：使用Makefile命令确保一致性
