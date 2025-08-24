#!/bin/bash

# Stock-A-Future Docker 停止脚本 (Linux/Mac)

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

# 停止服务
stop_services() {
    log_info "停止 Docker 服务..."
    
    cd "$DOCKER_DIR"
    
    # 使用 docker-compose 或 docker compose
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    else
        COMPOSE_CMD="docker compose"
    fi
    
    # 停止服务
    $COMPOSE_CMD down
    
    if [ $? -eq 0 ]; then
        log_success "服务停止成功"
    else
        log_error "服务停止失败"
        exit 1
    fi
}

# 主函数
main() {
    echo "⏹️ Stock-A-Future Docker 停止脚本"
    echo "================================="
    
    stop_services
    
    log_success "所有服务已停止"
}

# 错误处理
trap 'log_error "脚本执行失败"; exit 1' ERR

# 执行主函数
main "$@"
