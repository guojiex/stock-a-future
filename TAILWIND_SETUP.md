# Tailwind CSS + DaisyUI 安装和使用指南

## 🎉 安装完成！

您的项目已经成功安装了 Tailwind CSS 和 DaisyUI，可以开始使用了！

## 📁 文件结构

```
web/static/
├── css/
│   ├── tailwind.css      # Tailwind CSS 入口文件
│   └── output.css        # 构建后的 CSS 文件（包含所有样式）
├── tailwind-test.html    # 测试页面
└── ...其他文件
```

## 🚀 使用方法

### 1. 在HTML中引入CSS

```html
<link href="./css/output.css" rel="stylesheet">
```

### 2. 使用Tailwind CSS类

```html
<!-- 响应式容器 -->
<div class="container mx-auto p-8">
    <!-- 网格布局 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- 卡片组件 -->
        <div class="card bg-base-200 shadow-xl">
            <div class="card-body">
                <h2 class="card-title text-success">标题</h2>
                <p>内容描述</p>
                <div class="card-actions justify-end">
                    <button class="btn btn-primary">按钮</button>
                </div>
            </div>
        </div>
    </div>
</div>
```

### 3. 使用DaisyUI组件

DaisyUI 提供了丰富的组件，如：
- `btn` - 按钮
- `card` - 卡片
- `modal` - 模态框
- `navbar` - 导航栏
- `table` - 表格
- `form` - 表单

## 🔧 构建命令

### 开发模式（监听文件变化）
```bash
npm run build:css
```

### 生产模式（压缩CSS）
```bash
npm run build:css:prod
```

## 🎨 主题系统

DaisyUI 支持多种主题，可以通过设置 `data-theme` 属性来切换：

```html
<html data-theme="light">     <!-- 浅色主题 -->
<html data-theme="dark">      <!-- 深色主题 -->
<html data-theme="corporate"> <!-- 企业主题 -->
<html data-theme="business">  <!-- 商务主题 -->
```

## 📱 响应式设计

Tailwind CSS 的响应式前缀：
- `sm:` - 640px+
- `md:` - 768px+
- `lg:` - 1024px+
- `xl:` - 1280px+
- `2xl:` - 1536px+

## 🧪 测试页面

打开 `web/static/tailwind-test.html` 来查看所有功能是否正常工作。

## 📚 学习资源

- [Tailwind CSS 官方文档](https://tailwindcss.com/docs)
- [DaisyUI 官方文档](https://daisyui.com/components/)
- [Tailwind CSS 中文文档](https://tailwindcss.cn/)

## 💡 下一步

1. 查看测试页面，熟悉各种组件
2. 开始重构现有的HTML，使用Tailwind CSS类
3. 利用DaisyUI组件快速构建界面
4. 根据需要自定义主题和样式

祝您开发愉快！ 🎊
