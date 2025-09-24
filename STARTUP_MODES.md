# Stock-A-Future 启动模式说明

## 🚀 新增Production模式

为了解决React Native开发服务器启动慢的问题，我们添加了Production模式选项。

## 📋 启动选项

### Windows (`start-react-web-en.bat`)
### Linux/macOS (`start-react-web-en.sh`)

```
1. React Web App (Development Mode - Hot Reload)
2. React Web App (Production Mode - Faster Startup)  ← 新增
3. React Native Mobile App (Development Mode)
4. React Native Mobile App (Production Mode)          ← 新增
5. Both Web and Mobile (Development Mode)
6. Both Web and Mobile (Production Mode)              ← 新增
```

## 🔄 模式对比

### Development Mode (开发模式)
- **优势**：热重载，实时代码更新，完整开发工具
- **劣势**：启动慢，资源占用高，运行速度慢
- **适用场景**：代码开发、调试、测试

### Production Mode (生产模式)
- **优势**：启动快，运行速度快，资源占用低，优化构建
- **劣势**：无热重载，代码更改需重新构建
- **适用场景**：演示、性能测试、接近真实环境的测试

## 📱 React Native Production 脚本

新增的npm脚本：

```json
{
  "android:release": "react-native run-android --variant=release",
  "ios:release": "react-native run-ios --configuration Release",
  "start:reset": "react-native start --reset-cache",
  "build:android": "cd android && ./gradlew assembleRelease",
  "build:ios": "cd ios && xcodebuild -workspace ... -configuration Release ...",
  "bundle:android": "react-native bundle --platform android --dev false ...",
  "bundle:ios": "react-native bundle --platform ios --dev false ..."
}
```

## 🌐 Web React Production 脚本

新增的npm脚本：

```json
{
  "start:prod": "npm run build && npx serve -s build -l 3000",
  "serve": "npx serve -s build -l 3000"
}
```

## ⚡ 性能对比

| 模式 | 启动时间 | 运行速度 | 资源占用 | 热重载 |
|------|----------|----------|----------|--------|
| Development | 慢 (30-60s) | 慢 | 高 | ✅ |
| Production | 快 (5-15s) | 快 | 低 | ❌ |

## 🛠️ 使用建议

1. **日常开发**：使用Development模式
2. **快速测试**：使用Production模式
3. **演示展示**：使用Production模式
4. **性能调试**：使用Production模式

## 🔧 技术实现

### Web Production模式
- 使用 `react-scripts build` 构建优化版本
- 使用 `serve` 包提供静态文件服务
- 启用代码分割、压缩、优化

### Mobile Production模式
- 使用 `react-native bundle` 预构建JS包
- 使用 `--dev false` 禁用开发模式
- 启用代码压缩和优化

## 📝 注意事项

1. **首次运行Production模式需要构建时间**
2. **Production模式下代码更改不会自动更新**
3. **需要重新构建才能看到代码更改**
4. **Production模式更接近真实应用性能**
