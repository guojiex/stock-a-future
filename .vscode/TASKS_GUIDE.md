# VSCode Tasks 使用指南

## 📋 可用任务列表

### 单独服务任务

#### 1. 启动 Go API 服务器 (端口8081) ✅ 跨平台兼容
- **描述**: 启动Go后端API服务器
- **端口**: 8081
- **快捷键**: `Cmd+Shift+B` (macOS) 或 `Ctrl+Shift+B` (Windows/Linux)
- **✨ 新特性**: 已修复Windows PowerShell兼容性问题，使用VSCode环境变量配置

#### 2. 启动 AKTools 服务 (端口8080) ✅ 跨平台兼容
- **描述**: 启动AKTools数据服务
- **端口**: 8080
- **依赖**: 需要安装Python和aktools包
- **✨ 新特性**: 跨平台兼容，支持Windows/macOS/Linux
- **备选方案**: 如果跨平台版本失败，可使用"启动 AKTools 服务 (Windows专用)"

#### 3. 启动 React Web (开发模式) ✅ 跨平台兼容
- **描述**: 启动React Web应用开发服务器
- **端口**: 3000
- **特性**: 热重载、开发工具
- **✨ 新特性**: 使用VSCode包管理器检测，支持npm/yarn/pnpm

#### 4. 启动 React Web (生产模式) ✅ 跨平台兼容
- **描述**: 构建并启动React Web应用生产版本
- **端口**: 3000
- **特性**: 优化构建、无源码映射
- **✨ 新特性**: 使用VSCode包管理器检测，支持npm/yarn/pnpm

#### 5. 启动 React Native Metro ✅ 跨平台兼容
- **描述**: 启动React Native Metro bundler
- **端口**: 8081 (Metro默认端口)
- **注意**: 与Go API端口冲突，建议单独使用
- **✨ 新特性**: 使用VSCode包管理器检测，修复Windows PowerShell兼容性

#### 6. 启动 React Native Android ✅ 跨平台兼容
- **描述**: 启动React Native Android应用
- **依赖**: 需要先启动Metro bundler
- **要求**: Android SDK和模拟器/真机
- **✨ 新特性**: 使用VSCode包管理器检测，支持Windows/macOS/Linux

#### 7. 启动 React Native iOS ✅ 跨平台兼容
- **描述**: 启动React Native iOS应用
- **依赖**: 需要先启动Metro bundler
- **要求**: macOS、Xcode、模拟器/真机
- **✨ 新特性**: 使用VSCode包管理器检测，但iOS开发仍需要macOS

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

## 🌍 跨平台兼容性改进 ✨

### 主要改进
1. **VSCode命令变量**: 使用 `${command:nodejs.packageManager}` 和 `${command:python.interpreterPath}` 替代硬编码命令
2. **包管理器检测**: 自动检测npm/yarn/pnpm，无需手动指定
3. **Python解释器**: 自动使用VSCode配置的Python解释器路径
4. **Shell配置**: 添加shell配置确保跨平台兼容性

### 支持的平台
- ✅ **Windows**: PowerShell, Command Prompt, Git Bash
- ✅ **macOS**: Terminal, iTerm2
- ✅ **Linux**: Bash, Zsh

### 环境要求
- **Node.js**: 18+ (自动检测包管理器)
- **Python**: 3.8+ (自动检测解释器路径)
- **Go**: 1.22+ (标准安装)

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

### Q2: React Native任务失败 ✅ 已修复
**问题**: `'react-scripts' 不是内部或外部命令` 错误
**原因**: React Native项目使用 `react-native` 命令，不是 `react-scripts`
**解决方案**: 
- 已更新所有React Native任务使用 `${command:nodejs.packageManager}`
- 自动检测npm/yarn/pnpm包管理器
- 修复Windows PowerShell兼容性问题

### Q3: Python命令找不到 ✅ 已修复
**问题**: `python3` 在Windows上找不到，或 `command 'python.interpreterPath' not found` 错误
**解决方案**: 
- 跨平台版本：使用简单的 `python` 命令
- Windows专用版本：如果跨平台版本失败，使用"启动 AKTools 服务 (Windows专用)"任务
- 确保Python已安装并在PATH中：`python --version`

### Q4: 任务无法启动
**A**: 检查依赖是否安装
- Go: `go version`
- Python: `python --version` 或 `python3 --version`
- Node.js: `node --version`
- npm: `npm --version`

### Q5: 端口被占用
**A**: 检查并关闭占用端口的进程
```bash
# macOS/Linux
lsof -i :8080
lsof -i :8081
lsof -i :3000

# Windows
netstat -ano | findstr :8080
netstat -ano | findstr :8081
netstat -ano | findstr :3000

# 关闭进程 (Windows)
taskkill /PID <PID> /F
```

### Q6: Metro和Go API端口冲突
**A**: 两种解决方案：
1. **推荐**: 修改Go API端口（已在tasks中配置为8081）
2. 修改Metro端口：
   ```bash
   # 在mobile目录下
   npm start -- --port 8082
   ```

### Q7: React Native任务失败
**A**: 检查依赖
- Android: 确保Android SDK已安装，`ANDROID_HOME`环境变量已设置
- iOS: 确保在macOS上运行，Xcode已安装
- Metro: 先手动运行 `npm start` 确保没有错误

### Q8: 组合任务部分失败
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
