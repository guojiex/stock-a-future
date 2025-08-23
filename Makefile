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
.PHONY: all build clean test deps run help stocklist fetch-stocks fetch-sse dev fmt vet tools lint env stop kill status restart test-tushare test-aktools aktools-test curl db-tools migrate echarts

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

# 构建Curl工具
curl:
	@echo "构建Curl工具..."
	@mkdir -p bin
	$(GOBUILD) -o ./bin/curl ./cmd/curl
	@echo "Curl工具构建完成: ./bin/curl"

# 构建数据库工具
db-tools:
	@echo "构建数据库工具..."
	@mkdir -p bin
	$(GOBUILD) -o ./bin/db-tools ./cmd/db-tools
	@echo "数据库工具构建完成: ./bin/db-tools"

# 构建数据迁移工具
migrate:
	@echo "构建数据迁移工具..."
	@mkdir -p bin
	$(GOBUILD) -o ./bin/migrate ./cmd/migrate
	@echo "数据迁移工具构建完成: ./bin/migrate"

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
	golangci-lint run ./...

lint-fix:
	@echo "运行代码质量检查并自动修复..."
	@golangci-lint run --fix ./...

# 代码格式化
fmt:
	@echo "格式化代码..."
	@go fmt ./...
	@echo "代码格式化完成"

# 代码导入整理
imports:
	@echo "整理代码导入..."
	@goimports -w .
	@echo "代码导入整理完成"

# 代码质量检查（快速模式）
lint-quick:
	@echo "快速代码质量检查..."
	@golangci-lint run --fast ./...

# 生成代码质量报告
lint-report:
	@echo "生成代码质量报告..."
	@golangci-lint run --out-format=html > lint-report.html
	@echo "代码质量报告已生成: lint-report.html"

# 测试覆盖率
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "测试覆盖率报告已生成: coverage.html"

# 代码质量检查（包含测试）
quality: fmt imports lint test-coverage
	@echo "代码质量检查完成！"

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

# 下载ECharts到本地
echarts:
	@echo "下载ECharts到本地..."
	@mkdir -p web/static/js/lib/echarts
	@if command -v curl >/dev/null 2>&1; then \
		curl -L "https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js" -o "web/static/js/lib/echarts/echarts.min.js"; \
	elif command -v wget >/dev/null 2>&1; then \
		wget "https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js" -O "web/static/js/lib/echarts/echarts.min.js"; \
	else \
		echo "错误: 需要安装curl或wget来下载文件"; \
		exit 1; \
	fi
	@echo "ECharts下载完成: web/static/js/lib/echarts/echarts.min.js"
	@echo "文件大小: $$(ls -lh web/static/js/lib/echarts/echarts.min.js | awk '{print $$5}')"

# 显示帮助信息
help:
	@echo "Stock-A-Future 构建命令:"
	@echo "  make build       - 构建应用程序"
	@echo "  make clean       - 清理构建文件"  
	@echo "  make test        - 运行测试"
	@echo "  make deps        - 下载依赖"
	@echo "  make run         - 构建并运行"
	@echo "  make dev         - 开发模式运行"
	@echo ""
	@echo "代码质量检查:"
	@echo "  make fmt         - 格式化代码"
	@echo "  make imports     - 整理代码导入"
	@echo "  make vet         - 代码检查"
	@echo "  make lint        - 代码质量检查"
	@echo "  make lint-fix    - 代码质量检查并自动修复"
	@echo "  make lint-quick  - 快速代码质量检查"
	@echo "  make lint-report - 生成代码质量报告"
	@echo "  make test-coverage - 生成测试覆盖率报告"
	@echo "  make quality     - 完整代码质量检查流程"
	@echo ""
	@echo "开发工具:"
	@echo "  make tools       - 安装开发工具"
	@echo "  make env         - 创建.env配置文件"
	@echo "  make stocklist   - 构建股票列表工具"
	@echo "  make fetch-stocks - 获取所有股票列表"
	@echo "  make fetch-sse   - 获取上交所股票列表"
	@echo "  make test-tushare - 测试Tushare API连接"
	@echo "  make test-aktools - 测试AKTools API连接"
	@echo "  make aktools-test - 构建AKTools测试工具"
	@echo "  make curl         - 构建内置Curl工具"
	@echo "  make db-tools     - 构建数据库管理工具"
	@echo "  make migrate      - 构建数据迁移工具"
	@echo "  make echarts      - 下载ECharts到本地"
	@echo ""
	@echo "服务器管理:"
	@echo "  make stop        - 停止运行中的服务器"
	@echo "  make kill        - 强制杀死服务器进程"
	@echo "  make status      - 检查服务器运行状态"
	@echo "  make restart     - 重启服务器"
	@echo ""
	@echo "  make help        - 显示帮助信息"
