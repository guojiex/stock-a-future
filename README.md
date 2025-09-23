# Stock-A-Future React Web App

基于React的现代化A股股票分析Web应用

## 🎉 功能特性

- **📊 实时股票数据** - 与Go后端API完美集成
- **🔍 智能搜索** - 股票名称和代码搜索，搜索历史
- **⭐ 收藏管理** - 股票收藏和分组功能
- **🎨 现代化UI** - Material-UI设计，响应式布局
- **🌓 主题切换** - 自动跟随系统深色/浅色主题
- **📱 移动端优化** - 完美适配手机和平板设备

## 🚀 技术栈

- **React 18** + **TypeScript**
- **Redux Toolkit** + **RTK Query** (状态管理和数据获取)
- **Material-UI (MUI)** (UI组件库)
- **React Router** (路由管理)
- **Recharts** (图表库)

## 📦 快速开始

### 🚀 方式一：一键全栈启动（推荐）

使用全栈启动脚本，自动启动Go后端 + 选择前端类型：

**Windows:**
```bash
.\start-react-web-en.bat
```

**Mac/Linux:**
```bash
chmod +x start-react-web-en.sh  # 首次运行添加权限
./start-react-web-en.sh
```

**功能特性：**
- ✅ 自动检查并启动Go后端服务器 (http://localhost:8081)
- ✅ 选择启动React Web App或React Native App
- ✅ 自动安装npm依赖
- ✅ 完整的环境检查和错误提示

### 🔧 方式二：手动启动

#### 1. 启动后端服务

```bash
go run cmd/server/main.go
```

#### 2. 启动前端（选择一个）

**React Web应用:**
```bash
cd web-react
npm install
npm start
```

**React Native应用:**
```bash
cd mobile
npm install
npm start
# 在另一个终端运行:
npm run android  # Android
npm run ios      # iOS
```

### 📱 服务地址

- **Go后端API**: http://localhost:8081
- **React Web应用**: http://localhost:3000  
- **React Native Metro**: http://localhost:8081

## 📂 项目结构

```
src/
├── components/         # React组件
│   ├── common/        # 通用组件
│   └── charts/        # 图表组件
├── pages/             # 页面组件
│   ├── MarketPage.tsx    # 市场页面
│   ├── SearchPage.tsx    # 搜索页面
│   ├── FavoritesPage.tsx # 收藏页面
│   └── SettingsPage.tsx  # 设置页面
├── services/          # API服务
│   └── api.ts         # RTK Query API定义
├── store/             # Redux状态管理
│   ├── index.ts       # Store配置
│   └── slices/        # 状态切片
├── types/             # TypeScript类型定义
├── constants/         # 常量和主题配置
├── hooks/             # 自定义React hooks
└── utils/             # 工具函数
```

## 🔗 API集成

应用完全集成了Go后端的所有API功能：

- ✅ 健康检查和连接状态监控
- ✅ 股票基本信息和日线数据
- ✅ 技术指标计算
- ✅ 基本面数据获取
- ✅ 股票搜索和列表
- ✅ 收藏管理功能

## 🎨 主题和样式

- **自动主题切换** - 跟随系统设置
- **中国股市颜色** - 红涨绿跌的习惯配色
- **响应式设计** - 适配桌面和移动设备
- **Material Design** - 遵循Google设计规范

## 📱 移动端体验

- 底部导航栏（移动设备）
- 触摸友好的交互
- 优化的列表和卡片布局
- 流畅的页面切换动画

## 🔧 开发命令

```bash
npm start          # 启动开发服务器
npm run build      # 构建生产版本
npm test           # 运行测试
npm run lint       # 代码检查
```

## 🌟 开发状态

### ✅ 已完成
- [x] 基础项目架构和配置
- [x] Redux状态管理集成
- [x] Material-UI主题系统
- [x] 路由和导航结构
- [x] API服务层完整集成
- [x] 市场页面（股票列表展示）
- [x] 搜索页面（智能搜索功能）
- [x] 收藏页面（基础结构）
- [x] 设置页面（基础结构）
- [x] 响应式布局和移动端优化

### 🚧 开发中
- [ ] 股票详情页面和K线图表
- [ ] 收藏功能完整实现
- [ ] 技术指标图表集成
- [ ] 基本面数据展示优化
- [ ] 实时数据更新
- [ ] 性能优化和缓存策略

### 📋 计划中
- [ ] 回测功能页面
- [ ] 用户偏好设置
- [ ] 数据导出功能
- [ ] PWA支持（离线使用）
- [ ] 更多图表类型
- [ ] 单元测试覆盖

## 🔍 与移动端的关系

这个React Web应用与 `../mobile/` 目录中的React Native应用共享：

- **相同的API服务层** - 完全一致的数据获取逻辑
- **相似的状态管理** - 使用相同的Redux Toolkit模式
- **一致的类型定义** - 共享TypeScript接口
- **统一的设计理念** - 保持用户体验一致性

## 🎯 使用场景

- **桌面端** - 完整的股票分析工作站
- **移动端** - 随时随地查看股票信息
- **平板端** - 平衡的阅读和操作体验
- **开发调试** - 快速测试API功能

## 🔧 自定义配置

### API基础URL

在 `src/services/api.ts` 中修改：

```typescript
const API_BASE_URL = 'http://localhost:8080/api/v1/';
```

### 主题自定义

在 `src/constants/theme.ts` 中修改颜色和样式配置。

## 🎉 立即体验

1. 确保Go后端服务运行在 `http://localhost:8080`
2. 运行 `npm start`
3. 在浏览器中打开 `http://localhost:3000`
4. 享受现代化的股票分析体验！

---

这个React Web应用为您提供了一个现代化、响应式的股票分析界面，完美补充了您的移动端应用。立即开始使用，体验Web端的强大功能！