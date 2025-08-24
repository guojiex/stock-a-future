# Docker 快速启动指南

## 🚀 一键启动

### Windows 用户

```cmd
cd docker
start.bat
```

### Linux/Mac 用户

```bash
cd docker
./start.sh
```

## 📋 前置要求

- **Docker** 20.10+
- **Docker Compose** 2.0+

## 🌐 访问地址

启动成功后：

- **Web界面**: http://localhost:8081
- **API接口**: http://localhost:8081/api/v1/health
- **AKTools**: http://localhost:8080

## 🛠️ 常用命令

| 操作 | Windows | Linux/Mac |
|------|---------|-----------|
| 启动服务 | `start.bat` | `./start.sh` |
| 停止服务 | `stop.bat` | `./stop.sh` |
| 重新构建 | `rebuild.bat` | `./rebuild.sh` |
| 查看日志 | `logs.bat` | `./logs.sh` |

## 🔧 故障排除

### 端口冲突
如果端口 8080 或 8081 被占用，修改 `docker-compose.yml` 中的端口映射。

### 服务启动失败
运行日志查看脚本检查错误信息：
```bash
# Windows
logs.bat

# Linux/Mac
./logs.sh
```

## 📚 详细文档

查看 [README.md](README.md) 获取完整的部署和配置说明。
