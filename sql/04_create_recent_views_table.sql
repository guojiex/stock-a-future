-- 最近查看表创建脚本
-- 文件名: 04_create_recent_views_table.sql
-- 描述: 创建最近查看股票表，记录用户最近查看的股票历史
-- 创建时间: 2025-10-08

-- 最近查看股票表
CREATE TABLE IF NOT EXISTS recent_views (
    id TEXT PRIMARY KEY,                -- 唯一标识
    ts_code TEXT NOT NULL,              -- 股票代码
    name TEXT NOT NULL,                 -- 股票名称
    symbol TEXT,                        -- 股票简称
    market TEXT,                        -- 市场类型 (主板/创业板等)
    viewed_at DATETIME NOT NULL,        -- 查看时间
    expires_at DATETIME NOT NULL,       -- 过期时间
    created_at DATETIME NOT NULL,       -- 创建时间
    updated_at DATETIME NOT NULL,       -- 更新时间
    UNIQUE(ts_code)                     -- 每个股票只保留一条最近查看记录
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_recent_views_ts_code ON recent_views(ts_code);
CREATE INDEX IF NOT EXISTS idx_recent_views_viewed_at ON recent_views(viewed_at DESC);
CREATE INDEX IF NOT EXISTS idx_recent_views_expires_at ON recent_views(expires_at);

-- 说明：
-- 1. 每个股票代码只保留一条记录（通过UNIQUE约束）
-- 2. 当用户再次查看同一股票时，更新viewed_at和expires_at
-- 3. expires_at = viewed_at + 2天，用于自动清理过期记录
-- 4. 按viewed_at降序排列，显示最近查看的股票

