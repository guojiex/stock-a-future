#!/bin/bash

# Stock-A-Future Docker 重新构建脚本 (Linux/Mac)

set -e

# 脚本配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DOCKER_DIR="$(dirname "$SCRIPT_DIR")"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 重新构建服务
rebuild_services() {
    log_info "重新构建 Docker 服务..."
    
    cd "$DOCKER_DIR"
    
    # 使用 docker-compose 或 docker compose
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    else
        COMPOSE_CMD="docker compose"
    fi
    
    # 停止现有服务
    log_info "停止现有服务..."
    $COMPOSE_CMD down
    
    # 清理旧镜像
    log_info "清理旧镜像..."
    $COMPOSE_CMD down --rmi local
    
    # 重新构建并启动
    log_info "重新构建并启动服务..."
    $COMPOSE_CMD up --build -d
    
    if [ $? -eq 0 ]; then
        log_success "服务重新构建成功"
    else
        log_error "服务重新构建失败"
        exit 1
    fi
}

# 等待服务就绪
wait_for_services() {
    log_info "等待服务就绪..."
    
    local max_attempts=30
    local attempt=0
    
    # 等待 AKTools 服务
    log_info "等待 AKTools 服务启动..."
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            log_success "AKTools 服务已就绪"
            break
        fi
        
        attempt=$((attempt + 1))
        echo -n "."
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        log_error "AKTools 服务启动超时"
        exit 1
    fi
    
    # 等待 Stock-A-Future 服务
    log_info "等待 Stock-A-Future 服务启动..."
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:8081/api/v1/health > /dev/null 2>&1; then
            log_success "Stock-A-Future 服务已就绪"
            break
        fi
        
        attempt=$((attempt + 1))
        echo -n "."
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        log_error "Stock-A-Future 服务启动超时"
        exit 1
    fi
}

# 显示服务状态
show_status() {
    log_info "服务状态:"
    
    cd "$DOCKER_DIR"
    
    if command -v docker-compose &> /dev/null; then
        docker-compose ps
    else
        docker compose ps
    fi
}

# 主函数
main() {
    echo "🔄 Stock-A-Future Docker 重新构建脚本"
    echo "====================================="
    
    rebuild_services
    wait_for_services
    show_status
    
    echo
    log_success "重新构建完成！"
    echo
    echo "📊 Stock-A-Future Web界面: http://localhost:8081"
    echo "🔗 Stock-A-Future API:    http://localhost:8081/api/v1/health"
    echo "📈 AKTools API:           http://localhost:8080"
}

# 错误处理
trap 'log_error "脚本执行失败"; exit 1' ERR

# 执行主函数
main "$@"
