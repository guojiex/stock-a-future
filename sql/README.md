# SQL 目录

本目录包含股票未来系统相关的数据库SQL语句。

## 文件说明

### 01_create_tables.sql
- **用途**: 创建股票收藏功能所需的数据库表结构
- **包含内容**:
  - `favorite_groups` 表: 收藏分组表
  - `favorite_stocks` 表: 收藏股票表
  - 相关索引: 提高查询性能

## 表结构说明

### favorite_groups (收藏分组表)
| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | TEXT | 分组唯一标识 | PRIMARY KEY |
| name | TEXT | 分组名称 | NOT NULL |
| color | TEXT | 分组颜色 | NOT NULL |
| sort_order | INTEGER | 排序顺序 | DEFAULT 0 |
| created_at | DATETIME | 创建时间 | NOT NULL |
| updated_at | DATETIME | 更新时间 | NOT NULL |

### favorite_stocks (收藏股票表)
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

## 索引说明

- `idx_favorite_stocks_ts_code`: 股票代码索引，提高按股票代码查询的性能
- `idx_favorite_stocks_group_id`: 分组ID索引，提高按分组查询的性能
- `idx_favorite_stocks_sort_order`: 复合索引(分组ID+排序)，提高分组内排序查询的性能
- `idx_favorite_groups_sort_order`: 分组排序索引，提高分组排序查询的性能

## 使用方法

1. **自动执行**: 应用启动时会自动执行这些SQL语句
2. **手动执行**: 可以使用SQLite客户端工具手动执行这些语句
3. **数据库初始化**: 首次运行时会自动创建表结构和索引

## 注意事项

- 所有表都使用 `IF NOT EXISTS` 语句，避免重复创建
- 外键约束确保数据完整性
- 时间字段使用 `DATETIME` 类型，便于时间计算和排序
- 索引设计考虑了常见的查询场景，优化查询性能

## 技术实现说明

### SQLite驱动选择
本项目使用 `modernc.org/sqlite` 作为SQLite驱动，这是一个纯Go实现的SQLite驱动，具有以下优势：

- **无需CGO**: 不需要C编译器，跨平台兼容性更好
- **纯Go实现**: 编译简单，部署方便
- **功能完整**: 支持所有标准SQLite功能
- **性能良好**: 性能接近原生SQLite

### 数据库连接
```go
// 使用 modernc.org/sqlite 驱动
db, err := sql.Open("sqlite", dbPath)
```

### 数据库文件
- **文件名**: `stock_future.db`
- **用途**: 存储整个股票未来系统的所有数据
- **位置**: `data/stock_future.db`

### 编译要求
- Go 1.22+
- 无需CGO支持
- 跨平台兼容（Windows、Linux、macOS）
