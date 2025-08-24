# Docker 部署方案总结

## 🎯 部署目标达成

✅ **一键启动 AKTools 和 Golang 程序**
✅ **简化依赖管理**
✅ **暴露 HTTP 端口供宿主机访问**
✅ **支持重新编译和热部署**

## 📁 完整目录结构

```
docker/
├── README.md                 # 详细部署文档
├── QUICK_START.md           # 快速启动指南
├── USAGE_EXAMPLES.md        # 使用示例文档
├── DEPLOYMENT_SUMMARY.md    # 本总结文档
├── Dockerfile               # Go 应用容器镜像
├── aktools.Dockerfile       # AKTools 容器镜像
├── docker-compose.yml       # 主要编排配置
├── docker-compose.dev.yml   # 开发模式配置
├── .dockerignore           # Docker 忽略文件
├── config/
│   └── docker.env          # 环境变量配置
├── scripts/
│   ├── start.sh/.bat       # 启动脚本
│   ├── stop.sh/.bat        # 停止脚本
│   ├── rebuild.sh/.bat     # 重新构建脚本
│   ├── logs.sh/.bat        # 日志查看脚本
│   └── status.sh/.bat      # 状态检查脚本
└── volumes/                # 数据持久化目录 (运行时创建)
    ├── data/               # 应用数据
    ├── logs/               # 应用日志
    ├── aktools-data/       # AKTools 数据
    └── aktools-logs/       # AKTools 日志
```

## 🚀 核心功能

### 1. 一键启动
- **Windows**: `docker\start.bat`
- **Linux/Mac**: `docker/start.sh`

### 2. 服务架构
- **AKTools 服务**: Python 数据源，端口 8080
- **Stock-A-Future 服务**: Go 应用，端口 8081
- **自动依赖管理**: AKTools 启动后再启动 Go 应用

### 3. 端口暴露
- **Web界面**: http://localhost:8081
- **API接口**: http://localhost:8081/api/v1/*
- **AKTools API**: http://localhost:8080

### 4. 数据持久化
- 数据库文件持久化
- 日志文件持久化
- 配置文件持久化

## 🛠️ 技术特性

### Docker 镜像优化
- **多阶段构建**: 减少最终镜像大小
- **非 root 用户**: 提高安全性
- **健康检查**: 自动监控服务状态
- **资源限制**: 防止资源过度使用

### 网络配置
- **内部网络**: 服务间安全通信
- **端口映射**: 宿主机访问服务
- **服务发现**: 通过服务名通信

### 环境管理
- **环境变量**: 统一配置管理
- **开发模式**: 支持代码热重载
- **日志管理**: 结构化日志输出

## 📋 使用流程

### 首次部署
1. 进入 docker 目录
2. 运行启动脚本
3. 等待服务就绪
4. 访问 Web 界面

### 日常操作
- **查看状态**: `status.sh/.bat`
- **查看日志**: `logs.sh/.bat`
- **重新构建**: `rebuild.sh/.bat`
- **停止服务**: `stop.sh/.bat`

### 开发调试
- **开发模式**: `docker-compose -f docker-compose.yml -f docker-compose.dev.yml up`
- **进入容器**: `docker-compose exec stock-a-future sh`
- **查看日志**: `docker-compose logs -f`

## 🔧 配置说明

### 环境变量 (config/docker.env)
```env
DATA_SOURCE_TYPE=aktools          # 数据源类型
AKTOOLS_BASE_URL=http://aktools:8080  # AKTools 内部地址
SERVER_PORT=8081                  # Go 应用端口
LOG_LEVEL=info                    # 日志级别
CACHE_ENABLED=true               # 缓存开关
```

### 端口映射
- `8080:8080` - AKTools 服务
- `8081:8081` - Stock-A-Future 服务

### 卷挂载
- `./volumes/data:/app/data` - 数据持久化
- `./volumes/logs:/app/logs` - 日志持久化

## 🔍 故障排除

### 常见问题
1. **端口冲突**: 修改 docker-compose.yml 端口映射
2. **服务启动失败**: 查看日志 `logs.sh/.bat`
3. **权限问题**: 检查 volumes 目录权限
4. **网络问题**: 检查 Docker 网络配置

### 调试命令
```bash
# 查看服务状态
docker-compose ps

# 查看详细日志
docker-compose logs --details

# 进入容器调试
docker-compose exec stock-a-future sh

# 测试网络连通性
docker-compose exec stock-a-future ping aktools
```

## 📈 性能优化

### 资源配置
- **AKTools**: 512MB 内存，0.5 CPU 核心
- **Stock-A-Future**: 1GB 内存，1 CPU 核心

### 缓存策略
- 应用层缓存减少 API 调用
- Docker 镜像层缓存加速构建
- 数据持久化避免重复初始化

## 🔒 安全考虑

- 容器以非 root 用户运行
- 内部网络隔离
- 只暴露必要端口
- 环境变量管理敏感配置

## 🎉 部署完成

现在您可以：

1. **一键启动**: 运行 `start.sh/.bat` 启动所有服务
2. **访问应用**: 打开 http://localhost:8081 使用 Web 界面
3. **调用 API**: 使用 http://localhost:8081/api/v1/* 访问 API
4. **查看日志**: 运行 `logs.sh/.bat` 查看服务日志
5. **管理服务**: 使用提供的脚本工具管理服务生命周期

Docker 部署方案已完成，实现了您要求的所有功能！
