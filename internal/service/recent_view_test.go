package service

import (
	"database/sql"
	"os"
	"path/filepath"
	"stock-a-future/internal/models"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// setupTestDB 创建测试数据库
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// 创建临时测试数据库
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("打开测试数据库失败: %v", err)
	}

	// 创建表
	createTableSQL := `
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

	if _, err := db.Exec(createTableSQL); err != nil {
		t.Fatalf("创建表失败: %v", err)
	}

	// 创建索引
	indexSQL := []string{
		"CREATE INDEX IF NOT EXISTS idx_recent_views_ts_code ON recent_views(ts_code);",
		"CREATE INDEX IF NOT EXISTS idx_recent_views_viewed_at ON recent_views(viewed_at DESC);",
		"CREATE INDEX IF NOT EXISTS idx_recent_views_expires_at ON recent_views(expires_at);",
	}

	for _, sql := range indexSQL {
		if _, err := db.Exec(sql); err != nil {
			t.Fatalf("创建索引失败: %v", err)
		}
	}

	// 返回清理函数
	cleanup := func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}

	return db, cleanup
}

func TestAddOrUpdateRecentView(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	service := NewRecentViewService(db)

	// 测试添加新记录
	req := &models.AddRecentViewRequest{
		TSCode: "600976.SH",
		Name:   "健民集团",
		Symbol: "健民集团",
		Market: "主板",
	}

	view, err := service.AddOrUpdateRecentView(req)
	if err != nil {
		t.Fatalf("添加记录失败: %v", err)
	}

	if view.TSCode != req.TSCode {
		t.Errorf("股票代码不匹配: got %s, want %s", view.TSCode, req.TSCode)
	}

	if view.Name != req.Name {
		t.Errorf("股票名称不匹配: got %s, want %s", view.Name, req.Name)
	}

	// 验证过期时间设置（应该是2天后）
	expectedExpiry := time.Now().Add(48 * time.Hour)
	timeDiff := view.ExpiresAt.Sub(expectedExpiry)
	if timeDiff > time.Minute || timeDiff < -time.Minute {
		t.Errorf("过期时间设置不正确: got %v, expected around %v", view.ExpiresAt, expectedExpiry)
	}

	// 测试更新记录
	time.Sleep(time.Second) // 等待一秒确保时间不同
	viewUpdated, err := service.AddOrUpdateRecentView(req)
	if err != nil {
		t.Fatalf("更新记录失败: %v", err)
	}

	if viewUpdated.ID != view.ID {
		t.Errorf("更新应该保持相同的ID: got %s, want %s", viewUpdated.ID, view.ID)
	}

	if !viewUpdated.UpdatedAt.After(view.UpdatedAt) {
		t.Error("更新时间应该更新")
	}
}

func TestGetRecentViews(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	service := NewRecentViewService(db)

	// 添加多条记录
	stocks := []struct {
		code   string
		name   string
		symbol string
		market string
	}{
		{"600976.SH", "健民集团", "健民集团", "主板"},
		{"000001.SZ", "平安银行", "平安银行", "主板"},
		{"002415.SZ", "海康威视", "海康威视", "中小板"},
	}

	for _, stock := range stocks {
		_, err := service.AddOrUpdateRecentView(&models.AddRecentViewRequest{
			TSCode: stock.code,
			Name:   stock.name,
			Symbol: stock.symbol,
			Market: stock.market,
		})
		if err != nil {
			t.Fatalf("添加记录失败: %v", err)
		}
		time.Sleep(time.Millisecond * 10) // 确保时间戳不同
	}

	// 获取记录
	views, err := service.GetRecentViews(10, false)
	if err != nil {
		t.Fatalf("获取记录失败: %v", err)
	}

	if len(views) != len(stocks) {
		t.Errorf("记录数量不匹配: got %d, want %d", len(views), len(stocks))
	}

	// 验证按时间倒序排列（最新的在前）
	if len(views) >= 2 {
		if views[0].ViewedAt.Before(views[1].ViewedAt) {
			t.Error("记录应该按查看时间倒序排列")
		}
	}
}

func TestDeleteRecentView(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	service := NewRecentViewService(db)

	// 添加记录
	req := &models.AddRecentViewRequest{
		TSCode: "600976.SH",
		Name:   "健民集团",
	}

	_, err := service.AddOrUpdateRecentView(req)
	if err != nil {
		t.Fatalf("添加记录失败: %v", err)
	}

	// 删除记录
	err = service.DeleteRecentView(req.TSCode)
	if err != nil {
		t.Fatalf("删除记录失败: %v", err)
	}

	// 验证记录已删除
	views, err := service.GetRecentViews(10, false)
	if err != nil {
		t.Fatalf("获取记录失败: %v", err)
	}

	if len(views) != 0 {
		t.Errorf("记录应该已被删除: got %d records", len(views))
	}

	// 删除不存在的记录应该返回错误
	err = service.DeleteRecentView("notexist.SH")
	if err == nil {
		t.Error("删除不存在的记录应该返回错误")
	}
}

func TestClearExpiredViews(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	service := NewRecentViewService(db)

	// 手动插入一条过期记录
	pastTime := time.Now().Add(-72 * time.Hour) // 3天前
	expiredView := &models.RecentView{
		ID:        "expired-1",
		TSCode:    "000001.SZ",
		Name:      "过期股票",
		Symbol:    "过期股票",
		Market:    "主板",
		ViewedAt:  pastTime,
		ExpiresAt: pastTime.Add(48 * time.Hour), // 已过期
		CreatedAt: pastTime,
		UpdatedAt: pastTime,
	}

	_, err := db.Exec(`
		INSERT INTO recent_views 
		(id, ts_code, name, symbol, market, viewed_at, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, expiredView.ID, expiredView.TSCode, expiredView.Name, expiredView.Symbol,
		expiredView.Market, expiredView.ViewedAt, expiredView.ExpiresAt,
		expiredView.CreatedAt, expiredView.UpdatedAt)

	if err != nil {
		t.Fatalf("插入过期记录失败: %v", err)
	}

	// 添加一条未过期的记录
	_, err = service.AddOrUpdateRecentView(&models.AddRecentViewRequest{
		TSCode: "600976.SH",
		Name:   "健民集团",
	})
	if err != nil {
		t.Fatalf("添加记录失败: %v", err)
	}

	// 清理过期记录
	count, err := service.ClearExpiredViews()
	if err != nil {
		t.Fatalf("清理过期记录失败: %v", err)
	}

	if count != 1 {
		t.Errorf("应该清理1条过期记录: got %d", count)
	}

	// 验证未过期的记录还在
	views, err := service.GetRecentViews(10, false)
	if err != nil {
		t.Fatalf("获取记录失败: %v", err)
	}

	if len(views) != 1 {
		t.Errorf("应该还有1条未过期记录: got %d", len(views))
	}

	if views[0].TSCode != "600976.SH" {
		t.Errorf("剩余记录错误: got %s, want 600976.SH", views[0].TSCode)
	}
}

func TestGetRecentViewCount(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	service := NewRecentViewService(db)

	// 初始应该为0
	count, err := service.GetRecentViewCount()
	if err != nil {
		t.Fatalf("获取记录数失败: %v", err)
	}

	if count != 0 {
		t.Errorf("初始记录数应该为0: got %d", count)
	}

	// 添加记录
	for i := 0; i < 5; i++ {
		_, err := service.AddOrUpdateRecentView(&models.AddRecentViewRequest{
			TSCode: "60097" + string(rune('0'+i)) + ".SH",
			Name:   "测试股票",
		})
		if err != nil {
			t.Fatalf("添加记录失败: %v", err)
		}
	}

	// 验证记录数
	count, err = service.GetRecentViewCount()
	if err != nil {
		t.Fatalf("获取记录数失败: %v", err)
	}

	if count != 5 {
		t.Errorf("记录数不匹配: got %d, want 5", count)
	}
}
