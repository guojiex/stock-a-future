# Docker 部署使用示例

## 🚀 基本使用流程

### 1. 首次部署

```bash
# 克隆项目
git clone <your-repo-url>
cd stock-a-future

# 进入 Docker 目录
cd docker

# 一键启动 (Windows)
start.bat

# 一键启动 (Linux/Mac)
./start.sh
```

### 2. 验证部署

```bash
# 检查服务状态
status.bat        # Windows
./status.sh       # Linux/Mac

# 查看服务日志
logs.bat          # Windows
./logs.sh         # Linux/Mac
```

### 3. 访问服务

- **Web界面**: http://localhost:8081
- **API文档**: http://localhost:8081/api
- **健康检查**: http://localhost:8081/api/v1/health

## 📋 常用操作示例

### 服务管理

```bash
# 启动服务
start.bat / ./start.sh

# 停止服务
stop.bat / ./stop.sh

# 重新构建并启动
rebuild.bat / ./rebuild.sh

# 查看服务状态
status.bat / ./status.sh
```

### 日志管理

```bash
# 查看所有服务日志 (最后50行)
logs.bat / ./logs.sh

# 实时跟踪日志
logs.bat -f / ./logs.sh -f

# 查看特定服务日志
logs.bat aktools / ./logs.sh aktools
logs.bat stock-a-future / ./logs.sh stock-a-future

# 查看更多行数的日志
logs.bat -t 100 / ./logs.sh -t 100
```

### Docker Compose 直接操作

```bash
cd docker

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs
docker-compose logs -f aktools
docker-compose logs -f stock-a-future

# 重启特定服务
docker-compose restart aktools
docker-compose restart stock-a-future

# 进入容器调试
docker-compose exec stock-a-future sh
docker-compose exec aktools bash
```

## 🔧 开发模式

### 启用开发模式

```bash
cd docker

# 使用开发配置启动
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

开发模式特性：
- 源码热重载
- 调试日志级别
- 禁用健康检查

### 代码修改后重新构建

```bash
# 停止服务
docker-compose down

# 重新构建镜像
docker-compose build --no-cache

# 启动服务
docker-compose up -d
```

## 🌐 API 使用示例

### 健康检查

```bash
curl http://localhost:8081/api/v1/health
```

### 获取股票基本信息

```bash
curl http://localhost:8081/api/v1/stocks/000001/basic
```

### 获取股票日线数据

```bash
curl "http://localhost:8081/api/v1/stocks/000001/daily?start_date=20240101&end_date=20240131"
```

### 获取技术指标

```bash
curl http://localhost:8081/api/v1/stocks/000001/indicators
```

### 获取买卖预测

```bash
curl http://localhost:8081/api/v1/stocks/000001/predictions
```

### 搜索股票

```bash
curl "http://localhost:8081/api/v1/stocks/search?q=平安银行"
```

## 🔍 故障排除示例

### 1. 端口冲突

**问题**: 启动时提示端口被占用

**解决方案**:
```bash
# 查看端口占用 (Windows)
netstat -ano | findstr :8080
netstat -ano | findstr :8081

# 查看端口占用 (Linux/Mac)
lsof -i :8080
lsof -i :8081

# 修改端口映射
# 编辑 docker-compose.yml
ports:
  - "8082:8081"  # 改为其他端口
```

### 2. 服务启动失败

**问题**: 服务无法正常启动

**解决方案**:
```bash
# 查看详细日志
logs.bat -f / ./logs.sh -f

# 检查容器状态
docker-compose ps

# 重新构建
rebuild.bat / ./rebuild.sh
```

### 3. 数据持久化问题

**问题**: 数据丢失或权限错误

**解决方案**:
```bash
# 检查卷挂载
docker-compose config

# 检查目录权限 (Linux/Mac)
ls -la docker/volumes/

# 修复权限 (Linux/Mac)
sudo chown -R $USER:$USER docker/volumes/
```

### 4. 网络连接问题

**问题**: 服务间无法通信

**解决方案**:
```bash
# 检查网络
docker network ls
docker network inspect stock-network

# 测试连接
docker-compose exec stock-a-future ping aktools
docker-compose exec aktools ping stock-a-future
```

## 📊 监控和维护

### 资源使用监控

```bash
# 查看容器资源使用
docker stats

# 查看特定容器资源使用
docker stats stock-a-future-app stock-aktools
```

### 数据备份

```bash
# 备份数据目录
tar -czf backup-$(date +%Y%m%d).tar.gz docker/volumes/data/

# 备份日志
tar -czf logs-backup-$(date +%Y%m%d).tar.gz docker/volumes/logs/
```

### 清理和维护

```bash
# 清理未使用的镜像
docker image prune

# 清理未使用的容器
docker container prune

# 清理未使用的网络
docker network prune

# 完全清理 (谨慎使用)
docker system prune -a
```

## 🔄 更新和升级

### 更新应用代码

```bash
# 拉取最新代码
git pull origin main

# 重新构建并启动
cd docker
rebuild.bat / ./rebuild.sh
```

### 更新 Docker 镜像

```bash
# 拉取最新基础镜像
docker pull golang:1.24-alpine
docker pull python:3.11-slim

# 重新构建
docker-compose build --no-cache
docker-compose up -d
```

## 📝 配置自定义

### 修改环境变量

编辑 `docker/config/docker.env`:

```env
# 修改日志级别
LOG_LEVEL=debug

# 修改缓存配置
CACHE_DEFAULT_TTL=2h
CACHE_MAX_AGE=48h

# 修改数据清理配置
CLEANUP_RETENTION_DAYS=60
```

### 修改端口映射

编辑 `docker/docker-compose.yml`:

```yaml
services:
  stock-a-future:
    ports:
      - "8082:8081"  # 修改外部端口
  
  aktools:
    ports:
      - "8083:8080"  # 修改外部端口
```

### 添加自定义配置

```yaml
services:
  stock-a-future:
    environment:
      - CUSTOM_CONFIG=value
    volumes:
      - ./custom-config:/app/config
```

这些示例涵盖了Docker部署的各种使用场景，帮助用户快速上手和解决常见问题。
