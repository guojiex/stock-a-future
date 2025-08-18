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

-- 股票信号存储表
CREATE TABLE IF NOT EXISTS stock_signals (
    id TEXT PRIMARY KEY,
    ts_code TEXT NOT NULL,
    name TEXT NOT NULL,
    trade_date TEXT NOT NULL,           -- 信号基于的交易日期
    signal_date TEXT NOT NULL,          -- 信号计算日期
    signal_type TEXT NOT NULL,          -- 信号类型: BUY, SELL, HOLD
    signal_strength TEXT NOT NULL,      -- 信号强度: STRONG, MEDIUM, WEAK
    confidence REAL NOT NULL,           -- 置信度 0-1
    patterns TEXT,                      -- 识别到的图形模式(JSON格式)
    technical_indicators TEXT,          -- 技术指标数据(JSON格式)
    predictions TEXT,                   -- 预测数据(JSON格式)
    description TEXT,                   -- 信号描述
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE(ts_code, trade_date)         -- 每个股票每天只有一个信号记录
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_favorite_stocks_ts_code ON favorite_stocks(ts_code);
CREATE INDEX IF NOT EXISTS idx_favorite_stocks_group_id ON favorite_stocks(group_id);
CREATE INDEX IF NOT EXISTS idx_favorite_stocks_sort_order ON favorite_stocks(group_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_favorite_groups_sort_order ON favorite_groups(sort_order);

-- 股票信号表索引
CREATE INDEX IF NOT EXISTS idx_stock_signals_ts_code ON stock_signals(ts_code);
CREATE INDEX IF NOT EXISTS idx_stock_signals_trade_date ON stock_signals(trade_date);
CREATE INDEX IF NOT EXISTS idx_stock_signals_signal_date ON stock_signals(signal_date);
CREATE INDEX IF NOT EXISTS idx_stock_signals_signal_type ON stock_signals(signal_type);
CREATE INDEX IF NOT EXISTS idx_stock_signals_ts_code_trade_date ON stock_signals(ts_code, trade_date);
