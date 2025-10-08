# VSCode Tasks 使用指南

## 📋 可用任务列表

### 单独服务任务

#### 1. 启动 Go API 服务器 (端口8081) ✅ 跨平台兼容
- **描述**: 启动Go后端API服务器
- **端口**: 8081
- **快捷键**: `Cmd+Shift+B` (macOS) 或 `Ctrl+Shift+B` (Windows/Linux)
- **✨ 新特性**: 已修复Windows PowerShell兼容性问题

#### 2. 启动 AKTools 服务 (端口8080)
- **描述**: 启动AKTools数据服务
- **端口**: 8080
- **依赖**: 需要安装Python和aktools包

#### 3. 启动 React Web (开发模式)
- **描述**: 启动React Web应用开发服务器
- **端口**: 3000
- **特性**: 热重载、开发工具

#### 4. 启动 React Web (生产模式)
- **描述**: 构建并启动React Web应用生产版本
- **端口**: 3000
- **特性**: 优化构建、无源码映射

#### 5. 启动 React Native Metro ✅ 跨平台兼容
- **描述**: 启动React Native Metro bundler
- **端口**: 8081 (Metro默认端口)
- **注意**: 与Go API端口冲突，建议单独使用
- **✨ 新特性**: 支持Windows/macOS/Linux跨平台

#### 6. 启动 React Native Android ✅ 跨平台兼容
- **描述**: 启动React Native Android应用
- **依赖**: 需要先启动Metro bundler
- **要求**: Android SDK和模拟器/真机
- **✨ 新特性**: 支持Windows/macOS/Linux跨平台

#### 7. 启动 React Native iOS ✅ 跨平台兼容
- **描述**: 启动React Native iOS应用
- **依赖**: 需要先启动Metro bundler
- **要求**: macOS、Xcode、模拟器/真机
- **✨ 新特性**: 支持跨平台配置，但iOS只能在macOS上运行

#### 8. 启动 React Native Metro (Windows CMD) 🆕 Windows专用
- **描述**: 使用Windows CMD启动Metro bundler
- **适用**: Windows系统，当默认任务失败时使用
- **命令**: `cmd /c npm start`

#### 9. 启动 React Native Metro (PowerShell) 🆕 Windows专用
- **描述**: 使用PowerShell启动Metro bundler
- **适用**: Windows系统，PowerShell环境
- **命令**: `powershell -Command npm start`

#### 10. 检查 React Native 环境 🆕 诊断工具
- **描述**: 检查React Native开发环境配置
- **命令**: `npx react-native doctor`
- **用途**: 诊断环境问题

#### 11. 清理 React Native 缓存 🆕 维护工具
- **描述**: 清理Metro缓存并重启
- **命令**: `npm run start:reset`
- **用途**: 解决缓存问题

---

### 组合任务

#### 🚀 启动完整开发环境 (推荐)
**按顺序启动**:
1. AKTools服务 (8080)
2. Go API服务器 (8081)
3. React Web (3000)

**适用场景**: Web开发

#### 🚀 启动完整开发环境 (含Mobile)
**按顺序启动**:
1. AKTools服务 (8080)
2. Go API服务器 (8081)
3. React Web (3000)
4. React Native Metro

**适用场景**: 同时开发Web和Mobile

---

## 🎯 如何使用

### 方法1: 命令面板 (推荐)
1. 按 `Cmd+Shift+P` (macOS) 或 `Ctrl+Shift+P` (Windows/Linux)
2. 输入 `Tasks: Run Task`
3. 选择要运行的任务

### 方法2: 快捷键
- 按 `Cmd+Shift+B` (macOS) 或 `Ctrl+Shift+B` (Windows/Linux)
- 会运行默认任务：`🚀 启动完整开发环境`

### 方法3: 终端菜单
1. 点击菜单栏 `Terminal` → `Run Task...`
2. 选择要运行的任务

---

## 🌐 跨平台兼容性

### Windows系统
- **默认Shell**: PowerShell
- **React Native**: 如果默认任务失败，使用Windows专用任务
- **推荐任务**:
  - `启动 React Native Metro (Windows CMD)` - 使用CMD
  - `启动 React Native Metro (PowerShell)` - 使用PowerShell
- **环境变量**: 确保`ANDROID_HOME`等环境变量正确设置

### macOS系统
- **默认Shell**: Bash/Zsh
- **React Native**: 支持所有任务
- **iOS开发**: 需要Xcode和iOS模拟器
- **推荐任务**: 使用默认的跨平台任务

### Linux系统
- **默认Shell**: Bash
- **React Native**: 支持所有任务
- **Android开发**: 需要Android SDK
- **推荐任务**: 使用默认的跨平台任务

### 平台特定问题解决
1. **Windows PowerShell问题**: 使用CMD版本的任务
2. **权限问题**: 确保有足够的文件系统权限
3. **环境变量**: 检查`PATH`、`ANDROID_HOME`等环境变量
4. **Node.js版本**: 确保使用Node.js 18+版本

---

## 📊 服务端口映射

| 服务 | 端口 | 说明 |
|------|------|------|
| AKTools | 8080 | 数据服务 |
| Go API | 8081 | 后端API |
| React Web | 3000 | Web前端 |
| React Native Metro | 8081 | Metro bundler (与Go API冲突) |

⚠️ **注意**: React Native Metro默认也使用8081端口，这会与Go API冲突。建议：
- **Web开发**: 使用 `🚀 启动完整开发环境`
- **Mobile开发**: 单独启动Metro和Mobile任务

---

## 🔧 任务配置详解

### presentation配置
```json
"presentation": {
    "echo": true,              // 显示命令
    "reveal": "always",        // 总是显示终端
    "focus": false,            // 不自动聚焦
    "panel": "dedicated",      // 使用独立终端面板
    "showReuseMessage": true,  // 显示重用消息
    "clear": false             // 不清空之前的输出
}
```

### dependsOn配置
```json
"dependsOn": [
    "任务1",
    "任务2"
],
"dependsOrder": "sequence"  // 按顺序执行
```

---

## 💡 使用技巧

### 1. 查看所有运行中的任务
- 点击VSCode底部状态栏的终端图标
- 会显示所有活动的终端

### 2. 停止任务
- 在终端面板中按 `Ctrl+C`
- 或点击终端右上角的垃圾桶图标

### 3. 重新运行任务
- 按 `Cmd+Shift+B` (macOS) 再次选择任务
- VSCode会提示是否重启已运行的任务

### 4. 自定义任务
编辑 `.vscode/tasks.json` 文件：
```json
{
    "label": "我的自定义任务",
    "type": "shell",
    "command": "echo 'Hello World'",
    "group": "build"
}
```

---

## 🐛 常见问题

### Q1: Windows PowerShell环境变量问题 ✅ 已修复
**问题**: `SERVER_PORT=8081 go run cmd/server/main.go` 在Windows PowerShell中报错
**解决方案**: 已更新tasks.json使用VSCode的env配置，现在支持跨平台

### Q2: 任务无法启动
**A**: 检查依赖是否安装
- Go: `go version`
- Python: `python3 --version`
- Node.js: `node --version`
- npm: `npm --version`

### Q3: 端口被占用
**A**: 检查并关闭占用端口的进程
```bash
# macOS/Linux
lsof -i :8080
lsof -i :8081
lsof -i :3000

# 关闭进程
kill -9 <PID>
```

### Q4: Metro和Go API端口冲突
**A**: 两种解决方案：
1. **推荐**: 修改Go API端口（已在tasks中配置为8081）
2. 修改Metro端口：
   ```bash
   # 在mobile目录下
   npm start -- --port 8082
   ```

### Q5: React Native任务失败 ✅ 已修复跨平台问题
**A**: 检查依赖和环境
- **Windows用户**: 如果默认任务失败，尝试使用"启动 React Native Metro (Windows CMD)"或"启动 React Native Metro (PowerShell)"
- **Android**: 确保Android SDK已安装，`ANDROID_HOME`环境变量已设置
- **iOS**: 确保在macOS上运行，Xcode已安装
- **Metro**: 先手动运行 `npm start` 确保没有错误
- **环境检查**: 使用"检查 React Native 环境"任务诊断问题
- **缓存问题**: 使用"清理 React Native 缓存"任务解决缓存相关错误

### Q6: 组合任务部分失败
**A**: 
- 组合任务会按顺序启动，如果某个任务失败，后续任务也会失败
- 可以单独运行失败的任务进行调试
- 查看终端输出了解具体错误

---

## 📝 开发工作流建议

### Web开发工作流
```
1. 运行: 🚀 启动完整开发环境
2. 等待所有服务启动（约30秒）
3. 浏览器访问: http://localhost:3000
4. 开始开发
```

### Mobile开发工作流
```
1. 单独运行: 启动 AKTools 服务
2. 单独运行: 启动 Go API 服务器
3. 单独运行: 启动 React Native Metro
4. 单独运行: 启动 React Native Android/iOS
5. 开始开发
```

### 全栈开发工作流
```
1. 运行: 🚀 启动完整开发环境 (含Mobile)
2. 等待所有服务启动
3. Web测试: http://localhost:3000
4. 手动运行Android/iOS应用
5. 同时开发Web和Mobile
```

---

## 🔗 相关文档

- [VSCode Tasks官方文档](https://code.visualstudio.com/docs/editor/tasks)
- [跨平台启动指南](../docs/guides/CROSS_PLATFORM_STARTUP.md) ✨ 新增
- [项目启动脚本](../start-react-web-en.sh)
- [React Native开发指南](../mobile/README.md)
- [功能开发清单](../mobile/FEATURE_TODO.md)

---

**提示**: 如果你有其他常用的开发任务，可以随时编辑 `.vscode/tasks.json` 添加新的任务！
