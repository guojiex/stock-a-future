# Stock-A-Future Makefile

# Go相关变量
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 项目变量
BINARY_NAME=stock-a-future
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/server
STOCKLIST_BINARY=stocklist
STOCKLIST_PATH=./bin/$(STOCKLIST_BINARY)
STOCKLIST_MAIN_PATH=./cmd/stocklist

# 默认目标
.PHONY: all build clean test deps run help stocklist fetch-stocks fetch-sse dev fmt vet tools lint env stop kill status restart test-tushare test-aktools aktools-test

all: clean deps build

# 构建二进制文件
build:
	@echo "构建应用程序..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "构建完成: $(BINARY_PATH)"

# 构建股票列表工具
stocklist:
	@echo "构建股票列表工具..."
	@mkdir -p bin
	$(GOBUILD) -o $(STOCKLIST_PATH) $(STOCKLIST_MAIN_PATH)
	@echo "构建完成: $(STOCKLIST_PATH)"

# 测试Tushare连接
test-tushare:
	@echo "测试Tushare API连接..."
	@$(GOCMD) run test_tushare_connection.go
	@echo "Tushare连接测试完成"

# 构建AKTools测试工具
aktools-test:
	@echo "构建AKTools测试工具..."
	@mkdir -p bin
	$(GOBUILD) -o ./bin/aktools-test ./cmd/aktools-test
	@echo "AKTools测试工具构建完成: ./bin/aktools-test"

# 测试AKTools连接
test-aktools: aktools-test
	@echo "测试AKTools API连接..."
	@./bin/aktools-test
	@echo "AKTools连接测试完成"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	$(GOCLEAN)
	@rm -rf bin/
	@echo "清理完成"

# 运行测试
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...

# 下载依赖
deps:
	@echo "下载依赖..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "依赖下载完成"

# 运行应用程序
run: build
	@echo "启动应用程序..."
	$(BINARY_PATH)

# 直接运行（开发模式）
dev:
	@echo "开发模式启动..."
	$(GOCMD) run $(MAIN_PATH)

# 格式化代码
fmt:
	@echo "格式化代码..."
	$(GOCMD) fmt ./...

# 代码检查
vet:
	@echo "代码检查..."
	$(GOCMD) vet ./...

# 安装工具
tools:
	@echo "安装开发工具..."
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 代码质量检查
lint:
	@echo "运行代码质量检查..."
	golangci-lint run

# 创建.env文件
env:
	@if [ ! -f .env ]; then \
		echo "创建.env文件..."; \
		cp .env.example .env; \
		echo "请编辑.env文件并填入您的Tushare Token"; \
	else \
		echo ".env文件已存在"; \
	fi

# 获取股票列表
fetch-stocks: stocklist
	@echo "从交易所官网获取股票列表..."
	@mkdir -p data
	$(STOCKLIST_PATH) -source=all -output=data/stock_list.json
	@echo "股票列表获取完成"

# 获取上交所股票列表
fetch-sse: stocklist
	@echo "从上海证券交易所获取股票列表..."
	@mkdir -p data
	$(STOCKLIST_PATH) -source=sse -output=data/sse_stocks.json

# 停止服务器
stop:
	@echo "停止Stock-A-Future服务器..."
	@pkill -f "go run.*cmd/server" || pkill -f "stock-a-future" || echo "没有找到运行中的服务器进程"
	@echo "服务器已停止"

# 强制杀死服务器进程
kill:
	@echo "强制杀死占用端口的服务器进程..."
	@# 首先尝试通过进程名杀死
	@pkill -9 -f "go run.*cmd/server" 2>/dev/null || true
	@pkill -9 -f "stock-a-future" 2>/dev/null || true
	@# 然后检查并杀死占用8080和8081端口的进程
	@for port in 8080 8081; do \
		pid=$$(lsof -t -i:$$port 2>/dev/null); \
		if [ -n "$$pid" ]; then \
			echo "杀死占用端口$$port的进程: $$pid"; \
			kill -9 $$pid 2>/dev/null || true; \
		fi; \
	done
	@echo "服务器进程已强制终止"

# 检查服务器状态
status:
	@echo "检查Stock-A-Future服务器状态..."
	@server_running=false; \
	if pgrep -f "go run.*cmd/server" > /dev/null || pgrep -f "stock-a-future" > /dev/null; then \
		echo "✅ 服务器正在运行"; \
		echo "运行中的进程:"; \
		pgrep -f "go run.*cmd/server" -l 2>/dev/null || pgrep -f "stock-a-future" -l 2>/dev/null; \
		server_running=true; \
	else \
		echo "❌ 服务器进程未找到"; \
	fi; \
	echo "端口占用情况:"; \
	for port in 8080 8081; do \
		pid=$$(lsof -t -i:$$port 2>/dev/null); \
		if [ -n "$$pid" ]; then \
			process_name=$$(ps -p $$pid -o comm= 2>/dev/null || echo "unknown"); \
			echo "  端口$$port: 被进程$$pid ($$process_name) 占用"; \
		else \
			echo "  端口$$port: 空闲"; \
		fi; \
	done

# 重启服务器
restart: stop dev
	@echo "服务器已重启"



# 显示帮助信息
help:
	@echo "Stock-A-Future 构建命令:"
	@echo "  make build       - 构建应用程序"
	@echo "  make clean       - 清理构建文件"  
	@echo "  make test        - 运行测试"
	@echo "  make deps        - 下载依赖"
	@echo "  make run         - 构建并运行"
	@echo "  make dev         - 开发模式运行"
	@echo "  make fmt         - 格式化代码"
	@echo "  make vet         - 代码检查"
	@echo "  make lint        - 代码质量检查"
	@echo "  make tools       - 安装开发工具"
	@echo "  make env         - 创建.env配置文件"
	@echo "  make stocklist   - 构建股票列表工具"
	@echo "  make fetch-stocks - 获取所有股票列表"
	@echo "  make fetch-sse   - 获取上交所股票列表"
	@echo "  make test-tushare - 测试Tushare API连接"
	@echo "  make test-aktools - 测试AKTools API连接"
	@echo "  make aktools-test - 构建AKTools测试工具"
	@echo ""
	@echo "服务器管理:"
	@echo "  make stop        - 停止运行中的服务器"
	@echo "  make kill        - 强制杀死服务器进程"
	@echo "  make status      - 检查服务器运行状态"
	@echo "  make restart     - 重启服务器"
	@echo ""
	@echo "  make help        - 显示帮助信息"
