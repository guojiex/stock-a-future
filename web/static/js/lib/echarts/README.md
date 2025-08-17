# 本地ECharts资源管理

## 概述
本目录包含ECharts的本地JavaScript文件，确保应用在网络不稳定或CDN不可用时仍能正常显示图表。

## 文件说明
- `echarts.min.js` - ECharts 5.4.3 压缩版本（约1001KB）
- `README.md` - 本说明文件

## 优势
1. **网络稳定性** - 不依赖外部CDN，避免网络问题
2. **加载速度** - 本地文件加载更快，减少延迟
3. **离线支持** - 完全离线环境下仍可使用
4. **版本控制** - 可以精确控制ECharts版本

## 更新方法

### 方法一：使用脚本（推荐）
```bash
# PowerShell
.\scripts\download-echarts.ps1

# 批处理
.\scripts\download-echarts.bat
```

### 方法二：手动下载
1. 访问 [ECharts官网](https://echarts.apache.org/zh/download.html)
2. 下载需要的版本
3. 替换 `echarts.min.js` 文件

### 方法三：使用npm
```bash
npm install echarts
cp node_modules/echarts/dist/echarts.min.js web/static/js/lib/echarts/
```

## 版本管理
- 当前版本：5.4.3
- 更新频率：建议跟随项目需求，不必频繁更新
- 兼容性：确保新版本与现有代码兼容

## 注意事项
1. 更新后需要测试图表功能是否正常
2. 如果使用按需加载，需要相应调整代码
3. 建议在更新前备份当前版本
4. 检查文件大小，确保下载完整

## 故障排除
如果图表无法显示：
1. 检查文件路径是否正确
2. 确认文件是否完整下载
3. 查看浏览器控制台错误信息
4. 验证HTML中的script标签路径
