# Go语言版本Curl工具 - 项目总结

## 🎯 项目目标

为Windows环境下的Go项目开发一个内置的curl工具，解决Windows系统上缺少curl命令的问题，方便API调试和测试。

## ✨ 功能特性

### 核心功能
- **HTTP方法支持**: GET、POST、PUT、DELETE、PATCH等所有HTTP方法
- **请求头管理**: 支持自定义请求头设置
- **数据发送**: 支持JSON和表单数据POST
- **认证支持**: 基本认证（用户名:密码）
- **SSL控制**: 可选的SSL证书验证跳过
- **超时设置**: 可配置的请求超时时间
- **重定向处理**: 可选的HTTP重定向跟随
- **详细输出**: 显示响应头、状态码、响应时间等详细信息

### 响应处理
- **JSON美化**: 自动格式化JSON响应
- **HTML截断**: 长HTML内容智能截断显示
- **响应头显示**: 详细的HTTP响应头信息
- **状态码显示**: HTTP状态码和状态信息

## 🏗️ 项目结构

```
cmd/curl/
├── main.go          # 主程序文件
└── README.md        # 详细使用说明

构建脚本:
├── build-curl.bat   # Windows批处理编译脚本
├── build-curl.ps1   # PowerShell编译脚本
└── test-curl.bat    # 功能测试脚本

Makefile集成:
└── make curl        # 使用Makefile编译
```

## 🚀 使用方法

### 编译
```bash
# 方法1: 使用批处理文件（推荐Windows用户）
build-curl.bat

# 方法2: 使用PowerShell脚本
.\build-curl.ps1

# 方法3: 使用Makefile
make curl

# 方法4: 手动编译
go build -o curl.exe ./cmd/curl
```

### 基本用法
```bash
# 显示帮助
.\curl.exe

# GET请求
.\curl.exe http://localhost:8081/api/v1/health

# POST请求
.\curl.exe -X POST -d "test=data" http://localhost:8081/api/v1/stocks

# JSON POST请求
.\curl.exe -X POST -H "Content-Type: application/json" -d "{\"query\":\"平安\"}" http://localhost:8081/api/v1/stocks/search

# 详细输出
.\curl.exe -v http://localhost:8081/api/v1/stocks

# 基本认证
.\curl.exe -u "username:password" http://localhost:8081/api/v1/stocks
```

## 🔧 技术实现

### 核心技术
- **Go标准库**: 使用`net/http`包进行HTTP请求
- **命令行解析**: 使用`flag`包解析命令行参数
- **JSON处理**: 使用`encoding/json`包处理JSON数据
- **错误处理**: 完善的错误处理和用户友好的错误信息

### 架构设计
- **模块化设计**: 清晰的函数分离和职责划分
- **配置管理**: 结构化的配置管理
- **响应处理**: 智能的响应内容处理
- **扩展性**: 易于添加新功能和选项

## 📊 测试结果

### 外部API测试
- ✅ httpbin.org GET请求测试通过
- ✅ httpbin.org POST请求测试通过
- ✅ JSON数据处理测试通过
- ✅ 基本认证测试通过
- ✅ 详细输出模式测试通过

### 本地API测试
- ✅ 健康检查接口测试通过
- ✅ 股票列表接口测试通过
- ✅ 股票搜索接口测试通过
- ✅ 股票详情接口测试通过

## 🌟 优势特点

### 相比标准curl的优势
1. **Windows原生支持**: 无需安装额外工具
2. **中文友好**: 更好的中文字符处理
3. **JSON智能格式化**: 自动美化JSON输出
4. **Go语言集成**: 与Go项目完美集成
5. **轻量级**: 单个可执行文件，无依赖

### 开发体验提升
1. **快速调试**: 无需切换到其他工具
2. **一致性**: 与项目使用相同的Go环境
3. **可定制**: 可以根据项目需求扩展功能
4. **文档完善**: 详细的使用说明和示例

## 🔮 未来扩展

### 功能增强
- [ ] 支持Cookie管理
- [ ] 支持文件上传
- [ ] 支持代理设置
- [ ] 支持请求/响应保存
- [ ] 支持批量请求

### 性能优化
- [ ] 连接池优化
- [ ] 并发请求支持
- [ ] 缓存机制
- [ ] 压缩支持

## 📝 使用建议

### 开发阶段
- 使用`-v`参数查看详细请求/响应信息
- 使用`--timeout`设置合理的超时时间
- 使用`-k`参数跳过SSL验证（仅开发环境）

### 生产环境
- 避免使用`-k`参数跳过SSL验证
- 设置合理的超时时间
- 使用适当的请求头进行身份验证

## 🎉 总结

这个Go语言版本的curl工具成功解决了Windows环境下缺少curl命令的问题，为项目开发提供了便利的API调试工具。工具功能完整、使用简单、性能良好，完全满足日常API开发和测试需求。

通过这个工具，开发者可以：
1. 快速测试API接口
2. 调试请求/响应问题
3. 验证API功能
4. 进行集成测试

这是一个成功的工具开发项目，体现了Go语言在系统工具开发方面的优势。
