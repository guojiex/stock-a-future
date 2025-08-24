# Docker 部署指南

## 概述

本目录包含了 Stock-A-Future 项目的 Docker 部署配置，可以一键启动 AKTools 和 Golang 程序，简化部署和依赖管理。

## 目录结构

```
docker/
├── README.md                 # 本文档
├── docker-compose.yml        # Docker Compose 配置文件
├── Dockerfile               # Golang 应用的 Dockerfile
├── aktools.Dockerfile       # AKTools 的 Dockerfile
├── scripts/
│   ├── start.sh            # 启动脚本 (Linux/Mac)
│   ├── start.bat           # 启动脚本 (Windows)
│   ├── stop.sh             # 停止脚本 (Linux/Mac)
│   ├── stop.bat            # 停止脚本 (Windows)
│   ├── rebuild.sh          # 重新构建脚本 (Linux/Mac)
│   └── rebuild.bat         # 重新构建脚本 (Windows)
├── config/
│   └── docker.env          # Docker 环境变量配置
└── volumes/
    ├── data/               # 数据持久化目录
    └── logs/               # 日志持久化目录

```

## 快速开始

### 前置要求

- Docker 20.10+
- Docker Compose 2.0+

### 一键启动

**Windows:**
```cmd
cd docker
start.bat
```

**Linux/Mac:**
```bash
cd docker
./start.sh
```

### 访问服务

启动成功后，可以通过以下地址访问：

- **Stock-A-Future Web界面**: http://localhost:8081
- **Stock-A-Future API**: http://localhost:8081/api/v1/health
- **AKTools API**: http://localhost:8080

## 详细说明

### 服务架构

本 Docker 部署包含两个主要服务：

1. **aktools**: Python 数据服务，提供股票数据 API
   - 端口: 8080
   - 基于官方 Python 镜像
   - 自动安装和启动 AKTools

2. **stock-a-future**: Go 应用服务，提供股票分析功能
   - 端口: 8081
   - 基于官方 Go 镜像
   - 多阶段构建，优化镜像大小
   - 依赖 aktools 服务

### 环境变量配置

配置文件位于 `config/docker.env`，主要配置项：

```env
# 数据源配置
DATA_SOURCE_TYPE=aktools
AKTOOLS_BASE_URL=http://aktools:8080

# 服务器配置
SERVER_PORT=8081
SERVER_HOST=0.0.0.0

# 日志配置
LOG_LEVEL=info

# 缓存配置
CACHE_ENABLED=true
CACHE_DEFAULT_TTL=1h
CACHE_MAX_AGE=24h
CACHE_CLEANUP_INTERVAL=10m
```

### 数据持久化

以下目录会被持久化到宿主机：

- `volumes/data/`: 数据库文件和缓存数据
- `volumes/logs/`: 应用日志文件

## 常用操作

### 查看服务状态

```bash
docker-compose ps
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs stock-a-future
docker-compose logs aktools

# 实时跟踪日志
docker-compose logs -f
```

### 重新构建和启动

**Windows:**
```cmd
rebuild.bat
```

**Linux/Mac:**
```bash
./rebuild.sh
```

### 停止服务

**Windows:**
```cmd
stop.bat
```

**Linux/Mac:**
```bash
./stop.sh
```

### 清理资源

```bash
# 停止并删除容器
docker-compose down

# 删除镜像
docker-compose down --rmi all

# 删除所有数据（谨慎使用）
docker-compose down -v
```

## 开发模式

如果需要在开发过程中实时更新代码，可以使用开发模式：

```bash
# 启动开发模式（挂载源码目录）
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

## 故障排除

### 常见问题

#### 1. 端口冲突
**错误**: `port is already allocated`
**解决**: 修改 `docker-compose.yml` 中的端口映射，或停止占用端口的服务

#### 2. AKTools 启动失败
**错误**: `aktools service unhealthy`
**解决**: 
- 检查 Python 依赖是否正确安装
- 查看 aktools 服务日志: `docker-compose logs aktools`

#### 3. Go 应用连接 AKTools 失败
**错误**: `AKTools API连接测试失败`
**解决**:
- 确保 aktools 服务已启动
- 检查网络连接: `docker-compose exec stock-a-future ping aktools`

#### 4. 数据持久化问题
**错误**: 数据丢失或权限问题
**解决**:
- 检查 volumes 目录权限
- 确保 Docker 有权限访问挂载目录

### 调试步骤

1. **检查服务状态**
   ```bash
   docker-compose ps
   ```

2. **进入容器调试**
   ```bash
   # 进入 Go 应用容器
   docker-compose exec stock-a-future sh
   
   # 进入 AKTools 容器
   docker-compose exec aktools bash
   ```

3. **检查网络连通性**
   ```bash
   # 从 Go 应用测试 AKTools 连接
   docker-compose exec stock-a-future wget -qO- http://aktools:8080/health
   ```

4. **查看详细日志**
   ```bash
   docker-compose logs --details
   ```

## 性能优化

### 资源限制

在 `docker-compose.yml` 中已配置了合理的资源限制：

- **aktools**: 内存限制 512MB，CPU 限制 0.5 核
- **stock-a-future**: 内存限制 1GB，CPU 限制 1 核

### 缓存优化

- 启用了应用层缓存，减少重复数据请求
- Docker 镜像使用多阶段构建，减少镜像大小
- 数据持久化避免重复初始化

## 安全注意事项

1. **网络隔离**: 服务间通过 Docker 内部网络通信
2. **端口暴露**: 只暴露必要的端口到宿主机
3. **数据保护**: 敏感数据通过环境变量管理
4. **权限控制**: 容器以非 root 用户运行

## 更新和维护

### 更新应用

1. 拉取最新代码
2. 运行重新构建脚本
3. 重启服务

```bash
git pull
./rebuild.sh
```

### 备份数据

```bash
# 备份数据目录
tar -czf backup-$(date +%Y%m%d).tar.gz volumes/data/

# 备份日志
tar -czf logs-backup-$(date +%Y%m%d).tar.gz volumes/logs/
```

## 支持

如果遇到问题：

1. 查看本文档的故障排除部分
2. 检查服务日志
3. 验证配置文件
4. 检查 Docker 和 Docker Compose 版本

## 版本历史

- **v1.0.0**: 初始 Docker 部署配置
  - 支持一键启动 AKTools 和 Go 应用
  - 数据持久化和日志管理
  - 完整的脚本工具集
