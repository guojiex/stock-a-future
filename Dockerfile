# 使用官方Go镜像作为构建环境
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 使用轻量级的alpine镜像作为运行环境
FROM alpine:latest

# 安装ca-certificates以支持HTTPS请求
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建环境复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件示例
COPY --from=builder /app/.env.example .

# 暴露端口
EXPOSE 8080

# 运行应用程序
CMD ["./main"]
