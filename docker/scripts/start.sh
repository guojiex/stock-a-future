#!/bin/bash

# Stock-A-Future Docker 启动脚本 (Linux/Mac)
# 用于一键启动 AKTools 和 Golang 程序

set -e

# 脚本配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DOCKER_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$DOCKER_DIR")"

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

# 检查 Docker 和 Docker Compose
check_docker() {
    log_info "检查 Docker 环境..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装或不在 PATH 中"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose 未安装或不在 PATH 中"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        log_error "Docker 守护进程未运行"
        exit 1
    fi
    
    log_success "Docker 环境检查通过"
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."
    
    mkdir -p "$DOCKER_DIR/volumes/data"
    mkdir -p "$DOCKER_DIR/volumes/logs"
    mkdir -p "$DOCKER_DIR/volumes/aktools-data"
    mkdir -p "$DOCKER_DIR/volumes/aktools-logs"
    
    log_success "目录创建完成"
}

# 检查端口占用
check_ports() {
    log_info "检查端口占用情况..."
    
    local ports=(8080 8081)
    local occupied_ports=()
    
    for port in "${ports[@]}"; do
        if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            occupied_ports+=($port)
        fi
    done
    
    if [ ${#occupied_ports[@]} -gt 0 ]; then
        log_warning "以下端口已被占用: ${occupied_ports[*]}"
        log_warning "请确保这些端口可用，或修改 docker-compose.yml 中的端口映射"
        
        read -p "是否继续启动? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "启动已取消"
            exit 0
        fi
    else
        log_success "端口检查通过"
    fi
}

# 启动服务
start_services() {
    log_info "启动 Docker 服务..."
    
    cd "$DOCKER_DIR"
    
    # 使用 docker-compose 或 docker compose
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    else
        COMPOSE_CMD="docker compose"
    fi
    
    # 构建并启动服务
    $COMPOSE_CMD up --build -d
    
    if [ $? -eq 0 ]; then
        log_success "服务启动成功"
    else
        log_error "服务启动失败"
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
        show_logs
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
        show_logs
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

# 显示访问信息
show_access_info() {
    echo
    log_success "=== 服务启动完成 ==="
    echo
    echo "📊 Stock-A-Future Web界面: http://localhost:8081"
    echo "🔗 Stock-A-Future API:    http://localhost:8081/api/v1/health"
    echo "📈 AKTools API:           http://localhost:8080"
    echo
    echo "📋 常用 API 端点:"
    echo "   健康检查: curl http://localhost:8081/api/v1/health"
    echo "   股票信息: curl http://localhost:8081/api/v1/stocks/000001/basic"
    echo "   日线数据: curl http://localhost:8081/api/v1/stocks/000001/daily"
    echo
    echo "📝 查看日志: $SCRIPT_DIR/logs.sh"
    echo "⏹️  停止服务: $SCRIPT_DIR/stop.sh"
    echo "🔄 重新构建: $SCRIPT_DIR/rebuild.sh"
    echo
}

# 显示日志
show_logs() {
    log_info "显示最近的服务日志:"
    
    cd "$DOCKER_DIR"
    
    if command -v docker-compose &> /dev/null; then
        docker-compose logs --tail=20
    else
        docker compose logs --tail=20
    fi
}

# 主函数
main() {
    echo "🚀 Stock-A-Future Docker 启动脚本"
    echo "=================================="
    
    check_docker
    create_directories
    check_ports
    start_services
    wait_for_services
    show_status
    show_access_info
    
    log_success "启动完成！"
}

# 错误处理
trap 'log_error "脚本执行失败"; exit 1' ERR

# 执行主函数
main "$@"
