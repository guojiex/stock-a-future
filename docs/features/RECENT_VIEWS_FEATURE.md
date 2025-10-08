# 最近查看功能实现文档

## 概述

实现了用户最近查看股票的持久化功能，将用户的浏览历史保存到数据库中，并设置自动过期机制（2天）。

## 实现时间

2025-10-08

## 功能特性

### 1. 数据持久化
- 最近查看记录保存到SQLite数据库
- 每个股票只保留一条记录（按ts_code唯一）
- 记录包含查看时间、过期时间等元数据

### 2. 自动过期
- 查看记录2天后自动过期
- 后台定时任务（每小时）自动清理过期记录
- 手动清理API接口

### 3. 用户体验
- 前端立即更新UI（乐观更新）
- 异步保存到后端，不阻塞用户操作
- 自动加载历史记录

## 技术架构

### 后端实现

#### 1. 数据库表设计

文件：`sql/04_create_recent_views_table.sql`

```sql
CREATE TABLE IF NOT EXISTS recent_views (
    id TEXT PRIMARY KEY,
    ts_code TEXT NOT NULL,
    name TEXT NOT NULL,
    symbol TEXT,
    market TEXT,
    viewed_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE(ts_code)
);
```

**索引：**
- `idx_recent_views_ts_code`: 股票代码索引
- `idx_recent_views_viewed_at`: 查看时间倒序索引
- `idx_recent_views_expires_at`: 过期时间索引

#### 2. Go数据模型

文件：`internal/models/stock.go`

```go
// RecentView 最近查看记录
type RecentView struct {
    ID        string    `json:"id"`
    TSCode    string    `json:"ts_code"`
    Name      string    `json:"name"`
    Symbol    string    `json:"symbol"`
    Market    string    `json:"market"`
    ViewedAt  time.Time `json:"viewed_at"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// AddRecentViewRequest 添加最近查看请求
type AddRecentViewRequest struct {
    TSCode string `json:"ts_code"`
    Name   string `json:"name"`
    Symbol string `json:"symbol"`
    Market string `json:"market"`
}
```

#### 3. 服务层

文件：`internal/service/recent_view.go`

**核心方法：**
- `AddOrUpdateRecentView()`: 添加或更新记录
- `GetRecentViews()`: 获取最近查看列表
- `DeleteRecentView()`: 删除指定记录
- `ClearExpiredViews()`: 清理过期记录
- `StartAutoCleanup()`: 启动自动清理任务

**特性：**
- 如果股票已存在，更新查看时间和过期时间
- 过期时间 = 当前时间 + 48小时
- 支持过滤已过期记录

#### 4. API Handler

文件：`internal/handler/stock.go`

**API端点：**
- `POST /api/v1/recent-views` - 添加/更新记录
- `GET /api/v1/recent-views?limit=20` - 获取列表
- `DELETE /api/v1/recent-views/{code}` - 删除记录
- `POST /api/v1/recent-views/cleanup` - 清理过期记录
- `DELETE /api/v1/recent-views` - 清空所有记录

#### 5. 路由注册

文件：`cmd/server/main.go`

- 初始化RecentViewService
- 启动自动清理任务（每小时执行一次）
- 注册API路由

### 前端实现

#### 1. API服务

文件：`web-react/src/services/api.ts`

**RTK Query端点：**
```typescript
getRecentViews: builder.query<ApiResponse, { limit?: number; includeExpired?: boolean }>
addRecentView: builder.mutation<ApiResponse, AddRecentViewRequest>
deleteRecentView: builder.mutation<ApiResponse, string>
clearExpiredRecentViews: builder.mutation<ApiResponse, void>
clearAllRecentViews: builder.mutation<ApiResponse, void>
```

**导出的Hooks：**
- `useGetRecentViewsQuery`
- `useAddRecentViewMutation`
- `useDeleteRecentViewMutation`
- `useClearExpiredRecentViewsMutation`
- `useClearAllRecentViewsMutation`

#### 2. Redux状态管理

文件：`web-react/src/store/slices/searchSlice.ts`

**新增Action：**
- `setRecentlyViewed`: 从后端API设置最近查看列表

**状态流：**
1. 组件加载时通过`useGetRecentViewsQuery`获取数据
2. 使用`setRecentlyViewed`更新Redux状态
3. 用户点击股票时：
   - 立即调用`addRecentlyViewed`更新本地状态
   - 异步调用`addRecentViewMutation`保存到后端

#### 3. UI组件集成

文件：`web-react/src/pages/MarketSearchPage.tsx`

**关键改动：**
```typescript
// 获取最近查看数据
const { data: recentViewsData } = useGetRecentViewsQuery({ limit: 20 });
const [addRecentViewMutation] = useAddRecentViewMutation();

// 同步到Redux
useEffect(() => {
  if (recentViewsData?.success && recentViewsData?.data?.views) {
    dispatch(setRecentlyViewed(recentViewsData.data.views));
  }
}, [recentViewsData, dispatch]);

// 点击处理（乐观更新）
const handleStockClick = async (stock: StockBasic) => {
  // 立即更新UI
  dispatch(addRecentlyViewed(stock));
  
  // 异步保存
  try {
    await addRecentViewMutation({
      ts_code: stock.ts_code,
      name: stock.name,
      symbol: stock.symbol,
      market: stock.market,
    }).unwrap();
  } catch (error) {
    console.error('保存失败:', error);
  }
  
  navigate(`/stock/${stock.ts_code}`);
};
```

## 使用方式

### 后端API使用

#### 添加最近查看记录
```bash
curl -X POST http://localhost:8080/api/v1/recent-views \
  -H "Content-Type: application/json" \
  -d '{
    "ts_code": "600976.SH",
    "name": "健民集团",
    "symbol": "健民集团",
    "market": "主板"
  }'
```

#### 获取最近查看列表
```bash
curl http://localhost:8080/api/v1/recent-views?limit=20
```

#### 删除指定记录
```bash
curl -X DELETE http://localhost:8080/api/v1/recent-views/600976.SH
```

#### 清理过期记录
```bash
curl -X POST http://localhost:8080/api/v1/recent-views/cleanup
```

### 前端使用

用户只需正常浏览股票，系统会自动记录。在首页的"最近查看"标签页中可以看到历史记录。

## 性能考虑

### 1. 数据库优化
- 使用索引加速查询
- UNIQUE约束防止重复记录
- 自动清理过期数据，控制表大小

### 2. 前端优化
- 乐观更新提升响应速度
- 异步保存不阻塞用户操作
- RTK Query自动缓存和重新验证

### 3. 后端优化
- 定时任务异步清理，不影响主线程
- 数据库连接池管理
- 批量查询减少数据库访问

## 测试建议

### 单元测试
```go
// 测试添加记录
func TestAddRecentView(t *testing.T) {
    service := NewRecentViewService(db)
    view, err := service.AddOrUpdateRecentView(&models.AddRecentViewRequest{
        TSCode: "600976.SH",
        Name:   "健民集团",
    })
    assert.NoError(t, err)
    assert.NotNil(t, view)
}

// 测试过期清理
func TestClearExpiredViews(t *testing.T) {
    service := NewRecentViewService(db)
    count, err := service.ClearExpiredViews()
    assert.NoError(t, err)
    assert.GreaterOrEqual(t, count, 0)
}
```

### 集成测试
```go
func TestRecentViewsAPI(t *testing.T) {
    // 1. 添加记录
    // 2. 获取列表
    // 3. 验证记录存在
    // 4. 删除记录
    // 5. 验证记录已删除
}
```

## 未来改进

### 1. 功能增强
- [ ] 支持按时间段过滤
- [ ] 支持搜索历史记录
- [ ] 导出/导入历史记录
- [ ] 多设备同步（需要用户系统）

### 2. 性能优化
- [ ] 使用Redis缓存热点数据
- [ ] 实现更精细的缓存策略
- [ ] 批量操作优化

### 3. 用户体验
- [ ] 添加删除确认提示
- [ ] 支持批量删除
- [ ] 显示查看次数
- [ ] 按访问频率排序

## 相关文件

### 后端
- `sql/04_create_recent_views_table.sql` - 数据库建表脚本
- `internal/models/stock.go` - 数据模型
- `internal/service/recent_view.go` - 业务逻辑
- `internal/service/database.go` - 数据库初始化
- `internal/handler/stock.go` - API处理器
- `cmd/server/main.go` - 服务启动

### 前端
- `web-react/src/services/api.ts` - API服务
- `web-react/src/store/slices/searchSlice.ts` - 状态管理
- `web-react/src/pages/MarketSearchPage.tsx` - UI组件

## 注意事项

1. **数据库迁移**: 首次运行时会自动创建表，已有数据库需要手动执行SQL脚本
2. **时区处理**: 所有时间使用服务器本地时区
3. **并发安全**: SQLite使用串行化模式，天然支持并发安全
4. **错误处理**: 前端异步保存失败不影响用户体验，只在控制台输出错误

## 参考资料

- [SQLite索引优化](https://www.sqlite.org/queryplanner.html)
- [RTK Query文档](https://redux-toolkit.js.org/rtk-query/overview)
- [Go定时任务最佳实践](https://pkg.go.dev/time#NewTicker)

