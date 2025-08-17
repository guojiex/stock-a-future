-- 股票未来系统数据库建表语句
-- 文件名: 01_create_tables.sql
-- 描述: 创建收藏分组表和收藏股票表
-- 创建时间: 2024年

-- 收藏分组表
CREATE TABLE IF NOT EXISTS favorite_groups (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    color TEXT NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- 收藏股票表
CREATE TABLE IF NOT EXISTS favorite_stocks (
    id TEXT PRIMARY KEY,
    ts_code TEXT NOT NULL,
    name TEXT NOT NULL,
    start_date TEXT,
    end_date TEXT,
    group_id TEXT NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (group_id) REFERENCES favorite_groups(id)
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_favorite_stocks_ts_code ON favorite_stocks(ts_code);
CREATE INDEX IF NOT EXISTS idx_favorite_stocks_group_id ON favorite_stocks(group_id);
CREATE INDEX IF NOT EXISTS idx_favorite_stocks_sort_order ON favorite_stocks(group_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_favorite_groups_sort_order ON favorite_groups(sort_order);
