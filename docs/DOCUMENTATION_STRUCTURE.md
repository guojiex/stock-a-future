# 文档结构说明

本文档说明了 Stock-A-Future 项目文档的组织结构和分类规则。

## 📁 目录结构

```
docs/
├── architecture/           # 系统架构和设计
│   ├── BACKTEST_SYSTEM_DESIGN.md
│   ├── DATABASE_MIGRATION.md
│   ├── LOG_SYSTEM_REFACTOR.md
│   └── QUANTITATIVE_SYSTEM_ROADMAP.md
├── deployment/             # 部署和服务器管理
│   └── SERVER_MANAGEMENT.md
├── features/               # 功能特性文档
│   ├── ENHANCED_PATTERN_RECOGNITION.md
│   ├── FAVORITES_FEATURE.md
│   ├── PATTERN_RECOGNITION.md
│   ├── PREDICTION_FIELDS_EXPLANATION.md
│   ├── SIGNAL_SERVICE.md
│   ├── TECHNICAL_INDICATORS_SUMMARY.md
│   └── TECHNICAL_INDICATORS.md
├── guides/                 # 用户指南和教程
│   ├── DATA_CLEANUP.md
│   └── README_VENV.md
├── imgs/                   # 图片资源
│   ├── 技术指标.png
│   ├── 搜索功能视图.png
│   ├── 日k线信息.png
│   ├── 股票收藏.png
│   └── 股票查询.png
├── integration/            # 第三方集成
│   ├── AKTOOLS_API_FIXES.md
│   ├── AKTOOLS_FUNDAMENTAL_API.md
│   ├── AKTOOLS_INTEGRATION.md
│   ├── stock-list-fetcher.md
│   └── TUSHARE_CONNECTION.md
├── performance/            # 性能优化
│   ├── CACHE_USAGE.md
│   └── LOADING_UI_OPTIMIZATION.md
└── ui/                     # 用户界面
    └── PREDICTION_UI_ENHANCEMENT.md
```

## 📋 分类规则

### 🏗️ architecture/ - 系统架构和设计
包含系统设计、架构决策、技术方案等文档：
- 数据库设计和迁移
- 日志系统架构
- 回测系统设计
- 量化系统路线图

### 🚀 deployment/ - 部署和服务器管理
包含部署、运维、服务器管理相关文档：
- 服务器启动停止管理
- 部署配置和脚本
- 运维指南

### ⭐ features/ - 功能特性文档
包含具体功能模块的设计和说明：
- 模式识别功能
- 收藏夹功能
- 信号服务
- 技术指标
- 预测字段说明

### 📖 guides/ - 用户指南和教程
包含用户使用指南、配置教程等：
- 数据清理使用指南
- Python环境配置教程
- 功能使用教程

### 🖼️ imgs/ - 图片资源
包含文档中使用的截图、示意图等：
- 功能界面截图
- 架构图表
- 流程示意图

### 🔌 integration/ - 第三方集成
包含与外部系统集成相关的文档：
- AKTools API集成
- Tushare连接配置
- 数据源集成

### ⚡ performance/ - 性能优化
包含性能优化、缓存策略等文档：
- 缓存系统使用指南
- UI加载优化
- 性能调优建议

### 🎨 ui/ - 用户界面
包含前端UI设计、交互优化等文档：
- UI组件设计
- 交互体验优化
- 界面功能增强

## 📝 文档命名规范

1. **文件名格式**：使用大写字母和下划线，如 `FEATURE_NAME.md`
2. **中英文混合**：支持中文文件名，但推荐使用英文
3. **描述性命名**：文件名应清楚描述文档内容
4. **版本控制**：重要更新可在文件名中包含版本信息

## 🔄 维护指南

### 添加新文档
1. 确定文档类型和分类
2. 选择合适的目录
3. 遵循命名规范
4. 更新本文档的目录结构

### 移动现有文档
1. 评估文档内容和用途
2. 确定更合适的分类
3. 更新相关引用链接
4. 更新本文档记录

### 定期维护
- 检查文档分类是否合理
- 清理过时或重复的文档
- 更新目录结构说明
- 保持分类规则的一致性

---

*最后更新：2025年1月*
*整理完成：已将所有根目录文档按类型分类到对应子目录*
