# Stock-A-Future Mobile App

基于React Native的A股股票分析移动应用

## 功能特性

- 📊 实时股票数据展示
- 🔍 智能股票搜索
- ⭐ 收藏股票管理
- 📈 技术指标分析
- 📋 基本面数据
- 🧪 策略回测
- 🎨 深色/浅色主题

## 技术栈

- **React Native** 0.73+
- **TypeScript** 5.0+
- **Redux Toolkit** + RTK Query
- **React Navigation** 6.x
- **React Native Paper** (Material Design)
- **Victory Native** (图表库)
- **React Native SVG** (矢量图形)

## 项目结构

```
src/
├── components/          # 通用组件
│   ├── charts/         # 图表组件
│   ├── common/         # 通用UI组件
│   └── forms/          # 表单组件
├── screens/            # 页面组件
│   ├── Market/         # 市场相关页面
│   ├── Search/         # 搜索页面
│   ├── Favorites/      # 收藏页面
│   ├── Backtest/       # 回测页面
│   └── Settings/       # 设置页面
├── services/           # API服务
├── store/              # Redux store
├── navigation/         # 导航配置
├── types/              # TypeScript类型定义
├── constants/          # 常量和主题
├── hooks/              # 自定义hooks
└── utils/              # 工具函数
```

## 开发环境设置

### 前置要求

- Node.js 18+
- React Native CLI
- Android Studio (Android开发)
- Xcode (iOS开发，仅macOS)

### 安装依赖

```bash
cd mobile
npm install
```

### iOS设置 (仅macOS)

```bash
cd ios
pod install
```

### 运行应用

```bash
# Android
npm run android

# iOS
npm run ios

# 启动Metro
npm start
```

## API配置

应用默认连接到本地Go后端服务：`http://localhost:8080/api/v1/`

可以在 `src/services/api.ts` 中修改API基础URL。

## 开发状态

### ✅ 已完成
- [x] 项目基础架构
- [x] Redux状态管理
- [x] API服务层
- [x] 导航结构
- [x] 基础页面组件
- [x] 主题系统

### 🚧 开发中
- [ ] 股票详情页面
- [ ] K线图表组件
- [ ] 技术指标展示
- [ ] 基本面数据展示
- [ ] 收藏功能
- [ ] 搜索功能完善

### 📋 计划中
- [ ] 回测功能
- [ ] 推送通知
- [ ] 离线数据缓存
- [ ] 性能优化
- [ ] 单元测试

## 构建发布

### Android

```bash
cd android
./gradlew assembleRelease
```

### iOS

使用Xcode构建和发布到App Store。

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License
