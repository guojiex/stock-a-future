# Go语言版本的Curl工具

这是一个用Go语言编写的curl工具，提供了与Linux/Unix curl命令类似的功能，特别适合在Windows环境下使用。

## 功能特性

- 支持所有HTTP方法（GET、POST、PUT、DELETE、PATCH等）
- 自定义请求头
- 请求体数据支持（JSON、表单数据）
- 基本认证
- SSL证书验证控制
- 超时设置
- 重定向处理
- 详细输出模式
- 响应格式化（JSON美化、HTML截断等）

## 编译

在项目根目录执行：

```bash
go build -o curl.exe ./cmd/curl
```

或者在Windows PowerShell中：

```powershell
go build -o curl.exe .\cmd\curl\
```

## 使用方法

### 基本语法

```bash
curl [选项] <URL>
```

### 常用选项

#### HTTP方法
```bash
# GET请求（默认）
curl http://localhost:8080/api/stocks

# POST请求
curl -X POST http://localhost:8080/api/stocks

# PUT请求
curl -X PUT http://localhost:8080/api/stocks/123

# DELETE请求
curl -X DELETE http://localhost:8080/api/stocks/123
```

#### 请求头
```bash
# 添加单个请求头
curl -H "Authorization: Bearer token123" http://localhost:8080/api/stocks

# 添加多个请求头
curl -H "Content-Type: application/json" -H "Accept: application/json" http://localhost:8080/api/stocks
```

#### 请求数据
```bash
# JSON数据
curl -X POST -d '{"name":"AAPL","price":150.50}' http://localhost:8080/api/stocks

# 表单数据
curl -X POST -d "name=AAPL&price=150.50" http://localhost:8080/api/stocks
```

#### 认证
```bash
# 基本认证
curl -u "username:password" http://localhost:8080/api/stocks
```

#### 其他选项
```bash
# 详细输出（显示响应头、状态码等）
curl -v http://localhost:8080/api/stocks

# 跳过SSL证书验证
curl -k https://localhost:8080/api/stocks

# 设置超时时间（秒）
curl --timeout 10 http://localhost:8080/api/stocks

# 跟随重定向
curl -L http://localhost:8080/api/stocks

# 限制重定向次数
curl -L --max-redirects 5 http://localhost:8080/api/stocks
```

## 实际使用示例

### 测试本地API服务器
```bash
# 获取股票列表
curl http://localhost:8080/api/stocks

# 搜索股票
curl -X POST -d '{"query":"AAPL"}' http://localhost:8080/api/stocks/search

# 添加收藏股票
curl -X POST -H "Content-Type: application/json" -d '{"symbol":"AAPL","name":"Apple Inc."}' http://localhost:8080/api/favorites
```

### 测试外部API
```bash
# 获取天气信息
curl "https://api.openweathermap.org/data/2.5/weather?q=Beijing&appid=YOUR_API_KEY"

# 测试GitHub API
curl -H "Accept: application/vnd.github.v3+json" https://api.github.com/users/octocat
```

## 输出格式

### 普通模式
只显示响应体内容

### 详细模式（-v）
- 响应时间
- HTTP状态码
- 响应头信息
- 响应体内容

### 响应体处理
- JSON响应：自动美化格式化
- HTML响应：超过500字符时截断显示
- 其他类型：完整显示

## 注意事项

1. 默认超时时间为30秒
2. 默认不跟随重定向
3. 默认进行SSL证书验证
4. 请求头格式为 `Key:Value`
5. 支持JSON和表单数据自动识别

## 与标准curl的区别

- 更简洁的命令行参数
- 自动JSON格式化
- 更好的中文支持
- Windows原生支持
- 无需额外安装依赖

这个工具特别适合在Windows环境下调试和测试API接口，提供了与标准curl相当的功能，同时保持了Go语言的简洁和高效。
