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

# 默认目标
.PHONY: all build clean test deps run help

all: clean deps build

# 构建二进制文件
build:
	@echo "构建应用程序..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "构建完成: $(BINARY_PATH)"

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

# 显示帮助信息
help:
	@echo "Stock-A-Future 构建命令:"
	@echo "  make build    - 构建应用程序"
	@echo "  make clean    - 清理构建文件"  
	@echo "  make test     - 运行测试"
	@echo "  make deps     - 下载依赖"
	@echo "  make run      - 构建并运行"
	@echo "  make dev      - 开发模式运行"
	@echo "  make fmt      - 格式化代码"
	@echo "  make vet      - 代码检查"
	@echo "  make lint     - 代码质量检查"
	@echo "  make tools    - 安装开发工具"
	@echo "  make env      - 创建.env配置文件"
	@echo "  make help     - 显示帮助信息"
