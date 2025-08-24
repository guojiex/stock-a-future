#!/bin/bash

# Stock-A-Future Docker 日志查看脚本 (Linux/Mac)

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

# 显示使用说明
show_usage() {
    echo "📝 Stock-A-Future Docker 日志查看脚本"
    echo "===================================="
    echo
    echo "用法: $0 [选项] [服务名]"
    echo
    echo "选项:"
    echo "  -f, --follow     实时跟踪日志"
    echo "  -t, --tail N     显示最后 N 行日志 (默认: 50)"
    echo "  -h, --help       显示此帮助信息"
    echo
    echo "服务名:"
    echo "  stock-a-future   Go 应用服务"
    echo "  aktools          Python 数据服务"
    echo "  (不指定则显示所有服务日志)"
    echo
    echo "示例:"
    echo "  $0                           # 显示所有服务的最后 50 行日志"
    echo "  $0 -f                        # 实时跟踪所有服务日志"
    echo "  $0 -t 100 stock-a-future     # 显示 Go 应用的最后 100 行日志"
    echo "  $0 -f aktools                # 实时跟踪 AKTools 服务日志"
}

# 主函数
main() {
    local follow=false
    local tail_lines=50
    local service=""
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--follow)
                follow=true
                shift
                ;;
            -t|--tail)
                tail_lines="$2"
                shift 2
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            stock-a-future|aktools)
                service="$1"
                shift
                ;;
            *)
                echo "未知选项: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    cd "$DOCKER_DIR"
    
    # 使用 docker-compose 或 docker compose
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    else
        COMPOSE_CMD="docker compose"
    fi
    
    # 构建日志命令
    local log_cmd="$COMPOSE_CMD logs"
    
    if [ "$follow" = true ]; then
        log_cmd="$log_cmd -f"
    fi
    
    log_cmd="$log_cmd --tail=$tail_lines"
    
    if [ -n "$service" ]; then
        log_cmd="$log_cmd $service"
        log_info "显示 $service 服务日志 (最后 $tail_lines 行)"
    else
        log_info "显示所有服务日志 (最后 $tail_lines 行)"
    fi
    
    if [ "$follow" = true ]; then
        log_info "实时跟踪模式 (按 Ctrl+C 退出)"
    fi
    
    echo
    
    # 执行日志命令
    eval $log_cmd
}

# 执行主函数
main "$@"
