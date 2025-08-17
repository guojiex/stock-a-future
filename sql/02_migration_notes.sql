-- 数据库迁移说明
-- 文件名: 02_migration_notes.sql
-- 描述: 从JSON文件存储迁移到SQLite数据库的说明
-- 创建时间: 2024年

-- 迁移概述
-- 本文件记录了从原有的JSON文件存储方式迁移到SQLite数据库的过程和注意事项

-- 迁移前的数据结构
-- 1. data/favorites/favorites.json - 存储收藏股票列表
-- 2. data/favorites/groups.json - 存储收藏分组信息

-- 迁移后的数据结构
-- 1. data/stock_future.db - SQLite数据库文件
-- 2. 包含 favorite_groups 和 favorite_stocks 两个表

-- 迁移过程说明
-- 1. 应用启动时自动检测是否需要迁移
-- 2. 如果数据库为空，则从JSON文件读取数据并插入到数据库
-- 3. 迁移完成后，后续操作直接使用数据库
-- 4. 原有的JSON文件保留作为备份

-- 迁移检查逻辑
-- SELECT COUNT(*) FROM favorite_stocks;
-- 如果返回0，说明需要迁移；如果大于0，说明已有数据，跳过迁移

-- 迁移步骤
-- 1. 读取 groups.json 文件，插入到 favorite_groups 表
-- 2. 读取 favorites.json 文件，插入到 favorite_stocks 表
-- 3. 在数据库中添加迁移完成标记（可选）

-- 注意事项
-- 1. 迁移过程中保持原有数据的完整性
-- 2. 迁移后验证数据是否正确导入
-- 3. 建议在迁移前备份原有JSON文件
-- 4. 迁移失败时应有回滚机制

-- 验证迁移结果
-- 可以使用以下SQL语句验证迁移是否成功：
-- SELECT COUNT(*) FROM favorite_groups;
-- SELECT COUNT(*) FROM favorite_stocks;
-- SELECT * FROM favorite_groups ORDER BY sort_order;
-- SELECT * FROM favorite_stocks ORDER BY group_id, sort_order;
