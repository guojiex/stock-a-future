# Section 展开/折叠功能使用说明

## 功能概述

为Stock-A-Future网页应用的所有section添加了展开和折叠功能，用户可以通过点击section标题右侧的箭头按钮来控制内容的显示和隐藏。

## 功能特性

### 1. 通用设计
- 所有section都自动支持展开/折叠功能
- 统一的UI设计和交互体验
- 支持动态生成的section（如收藏列表）

### 2. 视觉效果
- 标题右侧显示向下箭头（▼）表示可折叠
- 折叠时箭头向左旋转90度（◀）
- 平滑的动画过渡效果
- 折叠状态下显示"(已折叠)"提示文字

### 3. 状态保持
- 用户的折叠状态自动保存到localStorage
- 页面刷新后保持之前的折叠状态
- 每个section独立记录状态

### 4. 响应式设计
- 支持桌面和移动设备
- 移动设备上隐藏状态文字以节省空间
- 按钮大小适配不同屏幕尺寸

## 支持的Section

### 静态Section
1. **股票查询区域** (`#search-section`)
   - 包含搜索表单、日期选择、快捷按钮等
   
2. **收藏股票区域** (`#favorites-section`)
   - 动态生成的收藏列表和管理功能

### 动态Section
1. **日线数据** (`#daily-data-section`)
   - 股票价格图表和数据摘要
   
2. **技术指标** (`#indicators-section`)
   - MACD、RSI、布林带等技术指标
   
3. **买卖预测** (`#predictions-section`)
   - 基于技术分析的买卖点预测
   
4. **错误信息** (`#error-section`)
   - 错误提示和异常信息

## 使用方法

### 用户操作
1. 点击section标题栏的任意位置来展开/折叠
2. 也可以直接点击右侧的箭头按钮
3. 折叠状态会自动保存，下次访问时保持

### 程序接口
```javascript
// 获取全局实例
const toggleModule = window.sectionToggleModule;

// 展开指定section
toggleModule.expandSection(document.getElementById('search-section'));

// 折叠指定section
toggleModule.collapseSection(document.getElementById('search-section'));

// 切换section状态
toggleModule.toggleSection(document.getElementById('search-section'));

// 展开所有section
toggleModule.expandAllSections();

// 折叠所有section
toggleModule.collapseAllSections();

// 重置所有折叠状态
toggleModule.resetCollapsedState();

// 获取section状态
const isCollapsed = toggleModule.getSectionState('search-section');

// 设置section状态
toggleModule.setSectionState('search-section', true); // 折叠
toggleModule.setSectionState('search-section', false); // 展开
```

## 技术实现

### CSS类结构
```css
.section-header          /* 标题容器 */
.section-toggle-btn      /* 切换按钮 */
.toggle-icon            /* 箭头图标 */
.section-content        /* 内容容器 */
.section-status         /* 状态提示文字 */
```

### HTML结构变化
原始结构：
```html
<section class="search-section">
    <div class="card">
        <h2>标题</h2>
        <div>内容</div>
    </div>
</section>
```

转换后结构：
```html
<section class="search-section collapsible">
    <div class="card">
        <div class="section-header">
            <h2>标题<span class="section-status">(已折叠)</span></h2>
            <button class="section-toggle-btn">
                <span class="toggle-icon">▼</span>
            </button>
        </div>
        <div class="section-content expanded">
            <div>内容</div>
        </div>
    </div>
</section>
```

### 状态管理
- 使用localStorage保存折叠状态
- 键名：`stockafuture_collapsed_sections`
- 存储格式：JSON数组，包含已折叠的section ID列表

## 浏览器兼容性

- 支持所有现代浏览器（Chrome 60+, Firefox 55+, Safari 12+）
- 使用CSS Grid和Flexbox布局
- 使用ES6类语法和现代JavaScript API
- 优雅降级，不支持的浏览器仍可正常使用（无折叠功能）

## 注意事项

1. **动态内容**：新添加的section会自动获得折叠功能
2. **性能**：使用MutationObserver监听DOM变化，性能开销很小
3. **存储**：折叠状态保存在localStorage，清除浏览器数据会重置状态
4. **无障碍性**：按钮包含aria-label属性，支持键盘导航
5. **动画**：使用CSS transition，在低性能设备上可能需要优化

## 自定义配置

如需自定义配置，可以修改以下CSS变量：
```css
:root {
    --section-toggle-duration: 0.3s;    /* 动画持续时间 */
    --section-toggle-easing: cubic-bezier(0.4, 0, 0.2, 1); /* 动画缓动函数 */
}
```

## 故障排除

### 常见问题

1. **折叠功能不工作**
   - 检查section是否有唯一的ID
   - 确认JavaScript模块已正确加载
   - 查看浏览器控制台是否有错误

2. **状态不保存**
   - 检查localStorage是否可用
   - 确认浏览器没有禁用本地存储

3. **动画效果异常**
   - 检查CSS是否正确加载
   - 确认没有其他样式冲突

4. **移动设备显示问题**
   - 检查响应式样式是否生效
   - 确认视口设置正确
