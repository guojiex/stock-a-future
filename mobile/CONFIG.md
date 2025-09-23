# React Native 应用配置说明

## 环境变量配置

React Native 应用现在支持通过环境变量配置后端端口和其他设置。

### 配置文件

在 `mobile/` 目录下创建 `.env` 文件：

```bash
# Stock-A-Future Mobile 环境变量配置

# 后端API配置
# 注意：确保后端服务器运行在此端口上
API_BASE_URL=http://127.0.0.1:8081/api/v1/
API_TIMEOUT=30000

# AKTools配置 (用于直接访问AKTools服务，如果需要的话)
AKTOOLS_BASE_URL=http://127.0.0.1:8080

# 应用配置
APP_NAME=Stock-A-Future Mobile
APP_VERSION=1.0.0

# 调试配置
DEBUG_MODE=false
LOG_LEVEL=info

# 缓存配置 (毫秒)
CACHE_TIMEOUT=300000
MAX_RETRIES=3

# 刷新间隔 (秒)
REFRESH_INTERVAL=60

# 默认显示的技术指标 (逗号分隔)
DEFAULT_INDICATORS=MA,MACD,RSI
```

### 使用说明

1. **创建配置文件**：
   ```bash
   cd mobile
   cp .env.example .env
   # 或者手动创建 .env 文件
   ```

2. **修改配置**：
   - 根据你的后端服务器配置修改 `API_BASE_URL`
   - 如果后端运行在不同端口，更新端口号
   - 其他配置可根据需要调整

3. **安装依赖**：
   ```bash
   cd mobile
   npm install
   # 或者
   yarn install
   ```

4. **配置react-native-config**（如果需要）：
   
   对于iOS：
   ```bash
   cd ios && pod install
   ```
   
   对于Android：
   - 配置已自动处理

### 配置项说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `API_BASE_URL` | 后端API基础URL | `http://127.0.0.1:8081/api/v1/` |
| `API_TIMEOUT` | API请求超时时间（毫秒） | `30000` |
| `AKTOOLS_BASE_URL` | AKTools服务URL | `http://127.0.0.1:8080` |
| `DEBUG_MODE` | 调试模式开关 | `false` |
| `LOG_LEVEL` | 日志级别 | `info` |
| `CACHE_TIMEOUT` | 缓存超时时间（毫秒） | `300000` |
| `MAX_RETRIES` | 最大重试次数 | `3` |
| `REFRESH_INTERVAL` | 数据刷新间隔（秒） | `60` |
| `DEFAULT_INDICATORS` | 默认技术指标 | `MA,MACD,RSI` |

### 动态配置更新

应用支持运行时动态更新配置。可以通过以下方式：

1. **Redux Action**：
   ```typescript
   import { useDispatch } from 'react-redux';
   import { updateConfig, syncConfigFromManager } from '../store/slices/appSlice';
   
   const dispatch = useDispatch();
   
   // 更新特定配置
   dispatch(updateConfig({ 
     apiBaseUrl: 'http://new-server:8082/api/v1/' 
   }));
   
   // 从配置管理器同步所有配置
   dispatch(syncConfigFromManager());
   ```

2. **配置管理器**：
   ```typescript
   import { appConfig } from '../constants/config';
   
   // 动态更新配置
   appConfig.updateConfig({
     API_BASE_URL: 'http://new-server:8082/api/v1/',
     API_TIMEOUT: 60000
   });
   ```

### 故障排除

1. **环境变量不生效**：
   - 确保 `.env` 文件在 `mobile/` 目录下
   - 重新构建应用（不是热重载）
   - 检查react-native-config是否正确安装

2. **API连接失败**：
   - 检查 `API_BASE_URL` 是否正确
   - 确保后端服务器正在运行
   - 检查网络连接和防火墙设置

3. **配置不更新**：
   - 使用 `syncConfigFromManager()` action 强制同步
   - 检查配置管理器的实现

### 开发环境 vs 生产环境

可以为不同环境创建不同的配置文件：

- `.env` - 开发环境
- `.env.staging` - 测试环境
- `.env.production` - 生产环境

React Native Config 会自动选择相应的配置文件。
