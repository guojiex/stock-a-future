package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"stock-a-future/internal/models"
	"time"

	_ "modernc.org/sqlite"
)

// DatabaseService 数据库服务
type DatabaseService struct {
	db     *sql.DB
	dbPath string
}

// NewDatabaseService 创建数据库服务
func NewDatabaseService(dataDir string) (*DatabaseService, error) {
	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %v", err)
	}

	// 数据库文件路径
	dbPath := filepath.Join(dataDir, "stock_future.db")

	// 连接数据库
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(1) // SQLite只支持一个连接
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	service := &DatabaseService{
		db:     db,
		dbPath: dbPath,
	}

	// 初始化数据库表
	if err := service.initTables(); err != nil {
		return nil, fmt.Errorf("初始化数据库表失败: %v", err)
	}

	return service, nil
}

// Close 关闭数据库连接
func (s *DatabaseService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// GetDB 获取数据库连接
func (s *DatabaseService) GetDB() *sql.DB {
	return s.db
}

// initTables 初始化数据库表
func (s *DatabaseService) initTables() error {
	// 创建收藏分组表
	createGroupsTable := `
	CREATE TABLE IF NOT EXISTS favorite_groups (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		color TEXT NOT NULL,
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	// 创建收藏股票表
	createStocksTable := `
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
	);`

	// 创建股票信号表
	createSignalsTable := `
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
	);`

	// 创建最近查看表
	createRecentViewsTable := `
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
	);`

	// 创建索引
	createIndexes := []string{
		// 收藏股票索引
		"CREATE INDEX IF NOT EXISTS idx_favorite_stocks_ts_code ON favorite_stocks(ts_code);",
		"CREATE INDEX IF NOT EXISTS idx_favorite_stocks_group_id ON favorite_stocks(group_id);",
		"CREATE INDEX IF NOT EXISTS idx_favorite_stocks_sort_order ON favorite_stocks(group_id, sort_order);",
		"CREATE INDEX IF NOT EXISTS idx_favorite_groups_sort_order ON favorite_groups(sort_order);",

		// 股票信号索引
		"CREATE INDEX IF NOT EXISTS idx_stock_signals_ts_code ON stock_signals(ts_code);",
		"CREATE INDEX IF NOT EXISTS idx_stock_signals_trade_date ON stock_signals(trade_date);",
		"CREATE INDEX IF NOT EXISTS idx_stock_signals_signal_date ON stock_signals(signal_date);",
		"CREATE INDEX IF NOT EXISTS idx_stock_signals_signal_type ON stock_signals(signal_type);",
		"CREATE INDEX IF NOT EXISTS idx_stock_signals_ts_code_trade_date ON stock_signals(ts_code, trade_date);",

		// 最近查看索引
		"CREATE INDEX IF NOT EXISTS idx_recent_views_ts_code ON recent_views(ts_code);",
		"CREATE INDEX IF NOT EXISTS idx_recent_views_viewed_at ON recent_views(viewed_at DESC);",
		"CREATE INDEX IF NOT EXISTS idx_recent_views_expires_at ON recent_views(expires_at);",
	}

	// 执行建表语句
	statements := append([]string{createGroupsTable, createStocksTable, createSignalsTable, createRecentViewsTable}, createIndexes...)

	for _, stmt := range statements {
		if _, err := s.db.Exec(stmt); err != nil {
			return fmt.Errorf("执行SQL语句失败: %v\nSQL: %s", err, stmt)
		}
	}

	return nil
}

// MigrateFromJSON 从JSON文件迁移数据到数据库
func (s *DatabaseService) MigrateFromJSON(favoritesPath, groupsPath string) error {
	// 检查是否需要迁移
	if err := s.checkMigrationNeeded(); err != nil {
		return err
	}

	// 迁移分组数据
	if err := s.migrateGroups(groupsPath); err != nil {
		return fmt.Errorf("迁移分组数据失败: %v", err)
	}

	// 迁移收藏数据
	if err := s.migrateFavorites(favoritesPath); err != nil {
		return fmt.Errorf("迁移收藏数据失败: %v", err)
	}

	// 标记迁移完成
	if err := s.markMigrationComplete(); err != nil {
		return fmt.Errorf("标记迁移完成失败: %v", err)
	}

	return nil
}

// checkMigrationNeeded 检查是否需要迁移
func (s *DatabaseService) checkMigrationNeeded() error {
	// 检查数据库中是否已有数据
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM favorite_stocks").Scan(&count)
	if err != nil {
		return fmt.Errorf("检查数据库状态失败: %v", err)
	}

	if count > 0 {
		// 数据库中已有数据，不需要迁移
		return nil
	}

	return nil
}

// migrateGroups 迁移分组数据
func (s *DatabaseService) migrateGroups(groupsPath string) error {
	// 检查JSON文件是否存在
	if _, err := os.Stat(groupsPath); os.IsNotExist(err) {
		// 文件不存在，创建默认分组
		return s.createDefaultGroup()
	}

	// 读取JSON文件
	data, err := os.ReadFile(groupsPath)
	if err != nil {
		return fmt.Errorf("读取分组文件失败: %v", err)
	}

	// 解析JSON数据
	var groups []*models.FavoriteGroup
	if err := json.Unmarshal(data, &groups); err != nil {
		return fmt.Errorf("解析分组数据失败: %v", err)
	}

	// 插入到数据库
	for _, group := range groups {
		if err := s.insertGroup(group); err != nil {
			return fmt.Errorf("插入分组失败: %v", err)
		}
	}

	return nil
}

// migrateFavorites 迁移收藏数据
func (s *DatabaseService) migrateFavorites(favoritesPath string) error {
	// 检查JSON文件是否存在
	if _, err := os.Stat(favoritesPath); os.IsNotExist(err) {
		// 文件不存在，无需迁移
		return nil
	}

	// 读取JSON文件
	data, err := os.ReadFile(favoritesPath)
	if err != nil {
		return fmt.Errorf("读取收藏文件失败: %v", err)
	}

	// 解析JSON数据
	var favorites []*models.FavoriteStock
	if err := json.Unmarshal(data, &favorites); err != nil {
		return fmt.Errorf("解析收藏数据失败: %v", err)
	}

	// 插入到数据库
	for _, favorite := range favorites {
		if err := s.insertFavorite(favorite); err != nil {
			return fmt.Errorf("插入收藏失败: %v", err)
		}
	}

	return nil
}

// createDefaultGroup 创建默认分组
func (s *DatabaseService) createDefaultGroup() error {
	defaultGroup := &models.FavoriteGroup{
		ID:        "default",
		Name:      "默认分组",
		Color:     "#3b82f6",
		SortOrder: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.insertGroup(defaultGroup)
}

// insertGroup 插入分组到数据库
func (s *DatabaseService) insertGroup(group *models.FavoriteGroup) error {
	stmt := `
		INSERT OR REPLACE INTO favorite_groups 
		(id, name, color, sort_order, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(stmt,
		group.ID,
		group.Name,
		group.Color,
		group.SortOrder,
		group.CreatedAt,
		group.UpdatedAt,
	)

	return err
}

// insertFavorite 插入收藏到数据库
func (s *DatabaseService) insertFavorite(favorite *models.FavoriteStock) error {
	stmt := `
		INSERT OR REPLACE INTO favorite_stocks 
		(id, ts_code, name, start_date, end_date, group_id, sort_order, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(stmt,
		favorite.ID,
		favorite.TSCode,
		favorite.Name,
		favorite.StartDate,
		favorite.EndDate,
		favorite.GroupID,
		favorite.SortOrder,
		favorite.CreatedAt,
		favorite.UpdatedAt,
	)

	return err
}

// markMigrationComplete 标记迁移完成
func (s *DatabaseService) markMigrationComplete() error {
	// 创建迁移记录表
	createMigrationTable := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		executed_at DATETIME NOT NULL
	);`

	if _, err := s.db.Exec(createMigrationTable); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %v", err)
	}

	// 记录迁移完成
	stmt := `INSERT INTO migrations (name, executed_at) VALUES (?, ?)`
	_, err := s.db.Exec(stmt, "json_to_sqlite_migration", time.Now())

	return err
}

// CleanupExpiredData 清理过期数据
func (s *DatabaseService) CleanupExpiredData(retentionDays int) error {
	log.Printf("开始清理过期数据...")
	log.Printf("股票信号数据保留天数: %d天", retentionDays)

	// 清理过期的股票信号数据
	if err := s.cleanupExpiredSignals(retentionDays); err != nil {
		log.Printf("清理过期信号数据失败: %v", err)
		return err
	}

	log.Printf("过期数据清理完成")
	return nil
}

// cleanupExpiredSignals 清理过期的股票信号数据
func (s *DatabaseService) cleanupExpiredSignals(retentionDays int) error {
	// 计算指定天数前的日期
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays).Format("20060102")

	// 删除指定天数前的信号数据
	result, err := s.db.Exec(`
		DELETE FROM stock_signals 
		WHERE trade_date < ? AND signal_date < ?
	`, cutoffDate, cutoffDate)
	if err != nil {
		return fmt.Errorf("删除过期信号数据失败: %v", err)
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		log.Printf("无法获取删除行数: %v", err)
	} else {
		log.Printf("清理了 %d 条过期的股票信号数据", deletedRows)
	}

	return nil
}

// GetDatabaseStats 获取数据库统计信息
func (s *DatabaseService) GetDatabaseStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取各表的记录数
	tables := []string{"favorite_groups", "favorite_stocks", "stock_signals"}

	for _, table := range tables {
		var count int
		err := s.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("获取表 %s 记录数失败: %v", table, err)
		}
		stats[table+"_count"] = count
	}

	// 尝试获取数据库文件大小（使用更安全的方式）
	if dbPath := s.getDatabasePath(); dbPath != "" {
		if fileInfo, err := os.Stat(dbPath); err == nil {
			stats["database_size_mb"] = float64(fileInfo.Size()) / (1024 * 1024)
		}
	}

	// 获取最早的信号数据日期
	var oldestSignalDate string
	err := s.db.QueryRow("SELECT MIN(trade_date) FROM stock_signals").Scan(&oldestSignalDate)
	if err == nil && oldestSignalDate != "" {
		stats["oldest_signal_date"] = oldestSignalDate
	}

	// 获取最新的信号数据日期
	var newestSignalDate string
	err = s.db.QueryRow("SELECT MAX(trade_date) FROM stock_signals").Scan(&newestSignalDate)
	if err == nil && newestSignalDate != "" {
		stats["newest_signal_date"] = newestSignalDate
	}

	return stats, nil
}

// getDatabasePath 获取数据库文件路径
func (s *DatabaseService) getDatabasePath() string {
	return s.dbPath
}
