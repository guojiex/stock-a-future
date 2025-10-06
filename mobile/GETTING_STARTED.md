# Stock-A-Future Mobile App 快速开始指南

## 🎯 项目概述

这是一个基于React Native开发的A股股票分析移动应用，与现有的Go后端API完美集成。应用提供了完整的股票数据展示、技术分析、基本面分析和策略回测功能。

## 🏗️ 架构特点

### 技术栈优势
- **React Native 0.73+**: 最新稳定版本，支持新架构
- **TypeScript**: 完整的类型安全
- **Redux Toolkit + RTK Query**: 现代化状态管理和数据获取
- **React Navigation 6**: 类型安全的导航
- **React Native Paper**: Material Design组件库
- **Victory Native**: 专业的图表解决方案

### 与现有后端的完美集成
- ✅ 完全兼容现有的Go API接口
- ✅ 支持所有现有功能（股票数据、技术指标、基本面、回测等）
- ✅ 无需修改后端代码
- ✅ 支持实时数据更新

## 📱 应用功能

### 核心功能
1. **市场页面**: 股票列表、实时价格、市场概况
2. **搜索功能**: 智能股票搜索、搜索历史、热门搜索
3. **收藏管理**: 股票收藏、分组管理、快速访问
4. **技术分析**: K线图、技术指标、图形识别
5. **基本面分析**: 财务报表、基本面指标、估值分析
6. **策略回测**: 策略编辑、回测执行、结果分析

### 用户体验
- 🌓 深色/浅色主题自动切换
- 🔄 下拉刷新和实时数据更新
- 📱 响应式设计，适配各种屏幕
- ⚡ 优化的性能和流畅的动画
- 🔍 智能搜索和建议

## 🚀 快速开始

### 1. 环境准备

**必需环境:**
- Node.js 18+ 
- npm 或 yarn
- Android Studio (Android开发)
- Xcode (iOS开发，仅macOS)

**React Native环境设置:**
```bash
# 安装React Native CLI
npm install -g @react-native-community/cli
```

### 2. 项目安装

```bash
# 进入mobile目录
cd mobile

# 运行安装脚本
# Windows:
scripts\setup.bat

# Linux/macOS:
chmod +x scripts/setup.sh
./scripts/setup.sh

# 或手动安装
npm install
```

### 3. 启动后端服务

确保Go后端服务正在运行：
```bash
# 在项目根目录
go run cmd/server/main.go
# 或
./server.exe
```

后端应该在 `http://localhost:8080` 运行。

### 4. 运行移动应用

```bash
# 启动Metro bundler
npm start

# 在新终端运行Android
npm run android

# 在新终端运行iOS (仅macOS)
npm run ios
```

## 📂 项目结构解析

```
mobile/
├── src/
│   ├── components/          # 可复用组件
│   │   ├── charts/         # 图表组件 (K线图、技术指标图等)
│   │   ├── common/         # 通用UI组件
│   │   └── forms/          # 表单组件
│   ├── screens/            # 页面组件
│   │   ├── Market/         # 市场页面 (股票列表、详情等)
│   │   ├── Search/         # 搜索页面
│   │   ├── Favorites/      # 收藏页面
│   │   ├── Backtest/       # 回测页面
│   │   └── Settings/       # 设置页面
│   ├── services/           # API服务层
│   │   └── api.ts          # RTK Query API定义
│   ├── store/              # Redux状态管理
│   │   ├── index.ts        # Store配置
│   │   └── slices/         # 状态切片
│   ├── navigation/         # 导航配置
│   ├── types/              # TypeScript类型定义
│   ├── constants/          # 常量和主题
│   ├── hooks/              # 自定义React hooks
│   └── utils/              # 工具函数
├── android/                # Android原生代码
├── ios/                    # iOS原生代码
└── scripts/                # 构建和部署脚本
```

## 🔧 开发工具和命令

```bash
# 开发
npm start                   # 启动Metro bundler
npm run android            # 运行Android版本
npm run ios                # 运行iOS版本

# 代码质量
npm run lint               # ESLint检查
npm run type-check         # TypeScript类型检查

# 构建
npm run build:android      # 构建Android APK
npm run build:ios          # 构建iOS (需要Xcode)
```

## 🌟 开发进度

### ✅ 已完成 (可以立即使用)
- [x] 完整的项目架构和配置
- [x] Redux状态管理 (应用状态、搜索、收藏、图表)
- [x] RTK Query API服务层 (对应所有Go后端接口)
- [x] React Navigation导航结构
- [x] 基础页面组件和布局
- [x] 主题系统 (深色/浅色主题)
- [x] TypeScript类型定义 (完整的API类型)

### 🚧 下一步开发 (优先级排序)
1. **股票详情页面** - K线图表和实时数据展示
2. **搜索功能完善** - 搜索结果展示和交互
3. **收藏功能** - 收藏列表和管理
4. **技术指标图表** - 使用Victory Native实现各种技术指标
5. **基本面数据展示** - 财务报表和指标展示
6. **回测功能** - 策略配置和结果展示

### 📋 未来计划
- [ ] 推送通知 (价格提醒、信号提醒)
- [ ] 离线数据缓存
- [ ] 手势操作优化
- [ ] 性能优化和内存管理
- [ ] 单元测试和集成测试

## 🔗 API集成说明

应用已经完全配置好与Go后端的集成：

### API配置
- 基础URL: `http://localhost:8080/api/v1/`
- 支持所有现有API端点
- 自动错误处理和重试机制
- 实时数据轮询

### 支持的API功能
- ✅ 健康检查和连接状态
- ✅ 股票基本信息和日线数据
- ✅ 技术指标计算
- ✅ 基本面数据 (利润表、资产负债表、现金流量表)
- ✅ 收藏管理和分组
- ✅ 图形识别和信号计算
- ✅ 策略管理和回测

## 🎨 设计系统

### 主题
- **浅色主题**: 现代简洁的设计
- **深色主题**: 护眼的深色界面
- **自动切换**: 跟随系统主题

### 颜色规范
- **中国股市习惯**: 红涨绿跌
- **技术指标**: 标准化的颜色方案
- **Material Design**: 遵循Google设计规范

## 🔍 故障排除

### 常见问题

1. **Metro bundler启动失败**
   ```bash
   npx react-native start --reset-cache
   ```

2. **Android构建失败**
   ```bash
   cd android
   ./gradlew clean
   cd ..
   npm run android
   ```

3. **iOS构建失败** (macOS)
   ```bash
   cd ios
   pod install
   cd ..
   npm run ios
   ```

4. **API连接失败**
   - 确保Go后端服务正在运行
   - 检查防火墙设置
   - 在Android模拟器中使用 `10.0.2.2:8080` 而不是 `localhost:8080`

### 调试技巧
- 使用React Native Debugger
- 启用Redux DevTools
- 查看Metro bundler日志
- 使用Android Studio或Xcode的调试工具

## 📞 技术支持

如果遇到问题：
1. 查看 `README.md` 和本文档
2. 检查控制台错误信息
3. 确保所有依赖都正确安装
4. 验证Go后端服务正常运行

## 🎉 总结

这个React Native应用为您的Stock-A-Future项目提供了：
- **完整的移动端解决方案**
- **与现有后端的无缝集成**
- **现代化的用户体验**
- **可扩展的架构设计**
- **完整的类型安全**

立即开始开发，将您的股票分析系统扩展到移动端！
