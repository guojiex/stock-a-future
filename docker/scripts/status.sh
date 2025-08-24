#!/bin/bash

# Stock-A-Future Docker 状态检查脚本 (Linux/Mac)

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

# 检查服务状态
check_services() {
    log_info "检查 Docker 服务状态..."
    
    cd "$DOCKER_DIR"
    
    # 使用 docker-compose 或 docker compose
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    else
        COMPOSE_CMD="docker compose"
    fi
    
    # 显示服务状态
    echo
    $COMPOSE_CMD ps
    echo
}

# 检查服务健康状态
check_health() {
    log_info "检查服务健康状态..."
    
    # 检查 AKTools
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        log_success "AKTools 服务正常"
    else
        log_error "AKTools 服务不可用"
    fi
    
    # 检查 Stock-A-Future
    if curl -s http://localhost:8081/api/v1/health > /dev/null 2>&1; then
        log_success "Stock-A-Future 服务正常"
    else
        log_error "Stock-A-Future 服务不可用"
    fi
}

# 显示访问信息
show_access_info() {
    echo
    log_info "服务访问地址:"
    echo "📊 Stock-A-Future Web界面: http://localhost:8081"
    echo "🔗 Stock-A-Future API:    http://localhost:8081/api/v1/health"
    echo "📈 AKTools API:           http://localhost:8080"
    echo
}

# 主函数
main() {
    echo "📊 Stock-A-Future Docker 状态检查"
    echo "================================="
    
    check_services
    check_health
    show_access_info
}

# 执行主函数
main "$@"
