-- 数据库维护和查询SQL
-- 文件名: 03_maintenance_queries.sql
-- 描述: 股票未来系统数据库维护操作和查询语句
-- 创建时间: 2024年

-- ========================================
-- 数据查询语句
-- ========================================

-- 查看所有分组
SELECT id, name, color, sort_order, created_at, updated_at 
FROM favorite_groups 
ORDER BY sort_order;

-- 查看所有收藏股票
SELECT fs.id, fs.ts_code, fs.name, fs.start_date, fs.end_date,
       fg.name as group_name, fg.color as group_color,
       fs.sort_order, fs.created_at, fs.updated_at
FROM favorite_stocks fs
JOIN favorite_groups fg ON fs.group_id = fg.id
ORDER BY fg.sort_order, fs.sort_order;

-- 按分组查看收藏股票
SELECT fg.name as group_name, fg.color as group_color,
       fs.ts_code, fs.name, fs.sort_order
FROM favorite_stocks fs
JOIN favorite_groups fg ON fs.group_id = fg.id
WHERE fg.id = ?  -- 替换为具体的分组ID
ORDER BY fs.sort_order;

-- 统计每个分组的股票数量
SELECT fg.name, fg.color, COUNT(fs.id) as stock_count
FROM favorite_groups fg
LEFT JOIN favorite_stocks fs ON fg.id = fs.group_id
GROUP BY fg.id, fg.name, fg.color
ORDER BY fg.sort_order;

-- ========================================
-- 数据维护语句
-- ========================================

-- 清理孤立的收藏记录（没有对应分组的记录）
DELETE FROM favorite_stocks 
WHERE group_id NOT IN (SELECT id FROM favorite_groups);

-- 重置分组排序（按名称字母顺序）
UPDATE favorite_groups 
SET sort_order = (
    SELECT rowid 
    FROM (
        SELECT id, ROW_NUMBER() OVER (ORDER BY name) as rowid 
        FROM favorite_groups
    ) ranked 
    WHERE ranked.id = favorite_groups.id
);

-- 重置分组内股票排序（按股票代码字母顺序）
UPDATE favorite_stocks 
SET sort_order = (
    SELECT rowid 
    FROM (
        SELECT id, ROW_NUMBER() OVER (PARTITION BY group_id ORDER BY ts_code) as rowid 
        FROM favorite_stocks
    ) ranked 
    WHERE ranked.id = favorite_stocks.id
);

-- 更新所有记录的更新时间
UPDATE favorite_groups SET updated_at = DATETIME('now');
UPDATE favorite_stocks SET updated_at = DATETIME('now');

-- ========================================
-- 数据备份和恢复
-- ========================================

-- 导出分组数据
SELECT id, name, color, sort_order, created_at, updated_at 
FROM favorite_groups 
ORDER BY sort_order;

-- 导出收藏股票数据
SELECT id, ts_code, name, start_date, end_date, group_id, sort_order, created_at, updated_at
FROM favorite_stocks 
ORDER BY group_id, sort_order;

-- 导出完整数据（包含分组信息）
SELECT 
    fs.id, fs.ts_code, fs.name, fs.start_date, fs.end_date,
    fg.name as group_name, fg.color as group_color,
    fs.sort_order, fs.created_at, fs.updated_at
FROM favorite_stocks fs
JOIN favorite_groups fg ON fs.group_id = fg.id
ORDER BY fg.sort_order, fs.sort_order;

-- ========================================
-- 性能优化查询
-- ========================================

-- 检查索引使用情况
SELECT name, sql FROM sqlite_master WHERE type = 'index';

-- 分析表统计信息
ANALYZE;

-- 查看表大小
SELECT 
    name,
    sqlite_compileoption_used('ENABLE_DBSTAT_VTAB') as dbstat_enabled,
    CASE 
        WHEN sqlite_compileoption_used('ENABLE_DBSTAT_VTAB') 
        THEN (SELECT pgsize FROM dbstat WHERE name = favorite_stocks.name)
        ELSE 'N/A'
    END as page_size
FROM sqlite_master 
WHERE type = 'table';

-- ========================================
-- 数据完整性检查
-- ========================================

-- 检查外键约束
PRAGMA foreign_key_check;

-- 检查表结构
.schema favorite_groups
.schema favorite_stocks

-- 检查数据一致性
SELECT 
    'orphaned_stocks' as issue,
    COUNT(*) as count
FROM favorite_stocks fs
LEFT JOIN favorite_groups fg ON fs.group_id = fg.id
WHERE fg.id IS NULL

UNION ALL

SELECT 
    'empty_groups' as issue,
    COUNT(*) as count
FROM favorite_groups fg
LEFT JOIN favorite_stocks fs ON fg.id = fs.group_id
WHERE fs.id IS NULL;
