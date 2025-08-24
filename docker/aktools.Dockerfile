# AKTools Dockerfile

FROM python:3.11-slim

# 设置工作目录
WORKDIR /app

# 安装系统依赖
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    wget \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 升级 pip
RUN pip install --upgrade pip

# 安装 AKTools
RUN pip install aktools

# 创建非 root 用户
RUN groupadd -r aktools && useradd -r -g aktools aktools

# 创建必要的目录
RUN mkdir -p /app/data /app/logs && \
    chown -R aktools:aktools /app

# 切换到非 root 用户
USER aktools

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 启动 AKTools 服务
CMD ["python", "-m", "aktools", "--host", "0.0.0.0", "--port", "8080"]
