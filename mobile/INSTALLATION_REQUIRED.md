# 📦 需要安装的依赖

## React Native DateTimePicker

市场页面重构后需要日期选择器组件。请按照以下步骤安装:

### 1. 安装npm包

```bash
cd mobile
npm install @react-native-community/datetimepicker
```

### 2. iOS额外步骤

```bash
cd ios
pod install
cd ..
```

### 3. Android配置

无需额外配置，npm install后即可使用。

### 4. 重启Metro Server

```bash
# 清除缓存并重启
npm run start:reset
```

### 5. 重新运行应用

**Android:**
```bash
npm run android
```

**iOS:**
```bash
npm run ios
```

## 验证安装

安装完成后，市场页面应该可以:
- ✅ 显示股票搜索框
- ✅ 显示日期选择器（点击日期输入框）
- ✅ 使用快捷日期按钮
- ✅ 显示K线图和其他分析数据

如果遇到问题，请查看[官方文档](https://github.com/react-native-datetimepicker/datetimepicker)

