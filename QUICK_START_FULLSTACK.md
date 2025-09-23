# Stock-A-Future Full Stack Quick Start

## 🚀 启动脚本说明

现在只需要一个启动脚本就能同时启动Go后端和React前端（Web或Native）！

## 📋 文件说明

- **`start-react-web-en.bat`** - Windows版本启动脚本
- **`start-react-web-en.sh`** - Mac/Linux版本启动脚本

## 🖥️ Windows 使用方法

直接双击运行或在命令行中执行：

```cmd
.\start-react-web-en.bat
```

## 🍎 Mac/Linux 使用方法

首先添加执行权限，然后运行：

```bash
chmod +x start-react-web-en.sh
./start-react-web-en.sh
```

## 🎯 功能特性

脚本会自动：

1. **检查环境**：
   - 检查Node.js是否安装
   - 检查Go是否安装

2. **启动Go后端**：
   - 检查后端服务器是否已运行
   - 如果未运行，自动在后台启动 `go run cmd/server/main.go`
   - 等待服务器启动并验证健康状态

3. **选择前端**：
   - 选项1：React Web App (浏览器)
   - 选项2：React Native Mobile App (移动端)
   - 选项3：同时启动Web和Mobile

4. **自动安装依赖**：
   - 检查并安装npm依赖
   - 支持首次运行自动配置

## 📱 启动选项详解

### 选项1：React Web App
- 启动基于浏览器的React Web应用
- 地址：http://localhost:3000
- 适合：桌面浏览器使用

### 选项2：React Native Mobile App
- 启动React Native Metro bundler
- 需要在另一个终端运行：
  - `npm run android` (Android模拟器/设备)
  - `npm run ios` (iOS模拟器/设备)
- 适合：移动应用开发

### 选项3：同时启动Web和Mobile
- 同时启动Web应用和Mobile应用
- Web应用在后台运行
- Mobile应用在前台运行
- 适合：同时开发和测试多个平台

## 🔧 服务地址

启动后可访问以下服务：

- **Go后端API**: http://localhost:8081
- **React Web应用**: http://localhost:3000
- **React Native Metro**: http://localhost:8081 (Metro bundler)
- **AKTools API**: http://localhost:8080 (需要单独启动)

## 📖 使用流程

1. **首次使用**：
   ```bash
   # Windows
   .\start-react-web-en.bat
   
   # Mac/Linux
   chmod +x start-react-web-en.sh
   ./start-react-web-en.sh
   ```

2. **选择前端类型**：
   - 输入 `1` 启动Web应用
   - 输入 `2` 启动Mobile应用
   - 输入 `3` 同时启动两者

3. **等待启动完成**：
   - Go后端会自动启动
   - 前端应用会根据选择启动
   - 首次运行会自动安装依赖

4. **开始使用**：
   - Web应用会自动在浏览器中打开
   - Mobile应用需要在另一个终端运行设备命令

## 🔍 故障排除

### Go后端启动失败
- 检查Go是否正确安装
- 检查端口8081是否被占用
- 查看终端错误信息

### 前端依赖安装失败
- 检查Node.js版本（需要18+）
- 检查网络连接
- 尝试手动运行 `npm install`

### React Native设备连接问题
- 确保Android/iOS开发环境已配置
- 检查设备/模拟器是否正常运行
- 参考React Native官方文档

## 📝 注意事项

1. **首次运行时间较长**：需要下载依赖和编译Go代码
2. **端口冲突**：确保8081和3000端口未被占用
3. **移动端开发**：需要配置Android Studio或Xcode
4. **网络要求**：需要互联网连接下载依赖

## 🎉 快速测试

启动后可以访问以下地址进行测试：

- http://localhost:8081/api/v1/health - 后端健康检查
- http://localhost:3000 - Web应用首页
- http://localhost:8081/static/ - 简单Web界面

全部正常即表示启动成功！
