# 数据库迁移指南

## 概述

本项目已从JSON文件存储方式迁移到SQLite数据库存储，提供更好的数据管理、查询性能和扩展性。

## 新的数据库结构

### 数据库文件
- **文件名**: `stock_future.db`
- **位置**: `data/stock_future.db`
- **用途**: 存储整个股票未来系统的所有数据

### 主要表结构

#### favorite_groups (收藏分组表)
| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | TEXT | 分组唯一标识 | PRIMARY KEY |
| name | TEXT | 分组名称 | NOT NULL |
| color | TEXT | 分组颜色 | NOT NULL |
| sort_order | INTEGER | 排序顺序 | DEFAULT 0 |
| created_at | DATETIME | 创建时间 | NOT NULL |
| updated_at | DATETIME | 更新时间 | NOT NULL |

#### favorite_stocks (收藏股票表)
| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | TEXT | 收藏记录唯一标识 | PRIMARY KEY |
| ts_code | TEXT | 股票代码 | NOT NULL |
| name | TEXT | 股票名称 | NOT NULL |
| start_date | TEXT | 开始日期 | 可选 |
| end_date | TEXT | 结束日期 | 可选 |
| group_id | TEXT | 所属分组ID | NOT NULL, FOREIGN KEY |
| sort_order | INTEGER | 排序顺序 | DEFAULT 0 |
| created_at | DATETIME | 创建时间 | NOT NULL |
| updated_at | DATETIME | 更新时间 | NOT NULL |

### 索引设计
- `idx_favorite_stocks_ts_code`: 股票代码索引
- `idx_favorite_stocks_group_id`: 分组ID索引
- `idx_favorite_stocks_sort_order`: 复合索引(分组ID+排序)
- `idx_favorite_groups_sort_order`: 分组排序索引

## 迁移过程

### 自动迁移
应用启动时会自动检测并执行数据迁移：
1. 检查数据库中是否已有数据
2. 如果数据库为空，从JSON文件读取数据并插入到数据库
3. 迁移完成后，后续操作直接使用数据库

### 手动迁移
使用专门的迁移工具：
```bash
# 构建迁移工具
go build -o migrate.exe ./cmd/migrate

# 执行迁移
./migrate.exe -data data
```

## 数据库管理工具

### db-tools 工具
提供数据库管理功能：

```bash
# 构建工具
go build -o db-tools.exe ./cmd/db-tools

# 查看数据库信息
./db-tools.exe -action info -data data

# 备份数据库
./db-tools.exe -action backup -data data

# 恢复数据库
./db-tools.exe -action restore -data data -backup path/to/backup.db
```

### 支持的操作
- `info`: 显示数据库信息和统计
- `backup`: 备份数据库到指定目录
- `restore`: 从备份文件恢复数据库

## 文件结构变化

### 迁移前
```
data/
├── favorites/
│   ├── favorites.json
│   └── groups.json
```

### 迁移后
```
data/
├── favorites/
│   ├── favorites.json (保留作为备份)
│   └── groups.json (保留作为备份)
├── stock_future.db (新的SQLite数据库)
└── backups/ (数据库备份目录)
```

## 技术实现

### SQLite驱动
使用 `modernc.org/sqlite` 纯Go实现的SQLite驱动：
- 无需CGO支持
- 跨平台兼容
- 性能良好

### 连接配置
```go
// 数据库连接
db, err := sql.Open("sqlite", dbPath)

// 连接池设置
db.SetMaxOpenConns(1) // SQLite只支持一个连接
db.SetMaxIdleConns(1)
db.SetConnMaxLifetime(time.Hour)
```

## 性能考虑

### 优势
- **查询性能**: 索引优化，支持复杂查询
- **数据完整性**: 外键约束，事务支持
- **扩展性**: 易于添加新表和功能
- **并发安全**: 支持多用户访问

### 注意事项
- SQLite是文件数据库，适合中小型应用
- 单个数据库文件，便于备份和迁移
- 支持标准SQL语法

## 故障排除

### 常见问题

#### 1. 数据库连接失败
- 检查数据目录权限
- 确认SQLite驱动已正确导入
- 验证数据库文件路径

#### 2. 迁移失败
- 检查JSON文件格式
- 确认数据目录结构
- 查看错误日志

#### 3. 性能问题
- 检查索引是否正确创建
- 优化查询语句
- 考虑数据库维护

### 日志和调试
- 启用详细日志记录
- 使用SQLite客户端工具检查数据库
- 验证表结构和数据完整性

## 未来扩展

### 计划中的功能
- 用户认证和权限管理
- 股票历史数据存储
- 技术指标计算
- 回测数据管理
- 策略配置存储

### 数据库设计原则
- 保持表结构规范化
- 合理使用索引
- 支持数据版本控制
- 便于数据迁移和升级
