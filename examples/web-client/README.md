# Stock-A-Future Web Client

这是一个重构后的模块化Web客户端，用于与Stock-A-Future API进行交互。

## 项目结构

```
js/
├── core/           # 核心模块
│   └── client.js   # 核心客户端类
├── services/       # 服务层
│   └── api.js      # API服务
├── modules/        # 功能模块
│   ├── charts.js           # 图表渲染
│   ├── stockSearch.js      # 股票搜索
│   ├── config.js           # 配置管理
│   ├── display.js          # 数据展示
│   └── events.js           # 事件处理
├── utils/          # 工具函数
│   └── helpers.js # 通用辅助函数
└── main.js         # 主入口文件
```

## 模块说明

### 核心模块 (core/)

- **client.js**: 核心客户端类，负责初始化、健康检查、API请求等基础功能

### 服务层 (services/)

- **api.js**: API服务类，封装所有与后端API的交互逻辑

### 功能模块 (modules/)

- **charts.js**: 图表渲染模块，负责K线图和技术指标的图表展示
- **stockSearch.js**: 股票搜索模块，处理股票搜索和建议功能
- **config.js**: 配置管理模块，负责服务器配置和连接测试
- **display.js**: 数据展示模块，负责各种数据的展示和格式化
- **events.js**: 事件处理模块，处理所有用户交互事件

### 工具函数 (utils/)

- **helpers.js**: 通用辅助函数，包含防抖、节流、格式化等实用工具

### 主入口 (main.js)

- **main.js**: 应用主入口，负责初始化所有模块和协调它们之间的交互

## 使用方法

1. 确保所有JavaScript文件都正确加载
2. 应用会在页面加载完成后自动初始化
3. 可以通过 `window.stockApp` 访问应用实例进行调试

## 重构优势

1. **模块化**: 代码按功能分离，便于维护和扩展
2. **可读性**: 每个模块职责单一，代码结构清晰
3. **可维护性**: 修改某个功能时只需要关注对应模块
4. **可扩展性**: 新增功能时可以轻松添加新模块
5. **可测试性**: 模块化结构便于单元测试

## 依赖关系

```
main.js
├── client.js (核心)
├── api.js (依赖 client.js)
├── charts.js (依赖 client.js)
├── stockSearch.js (依赖 client.js, api.js)
├── config.js (依赖 client.js)
├── display.js (依赖 client.js, charts.js)
├── events.js (依赖 client.js, api.js, display.js)
└── helpers.js (独立工具函数)
```

## 注意事项

- 所有模块都通过 `window` 对象导出，确保加载顺序正确
- 模块间通过依赖注入的方式协作，降低耦合度
- 保持了原有功能的完整性，重构不影响用户体验
