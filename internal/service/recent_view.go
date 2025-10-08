package service

import (
	"database/sql"
	"fmt"
	"log"
	"stock-a-future/internal/models"
	"time"

	"github.com/google/uuid"
)

// RecentViewService 最近查看服务
type RecentViewService struct {
	db *sql.DB
}

// NewRecentViewService 创建最近查看服务实例
func NewRecentViewService(db *sql.DB) *RecentViewService {
	return &RecentViewService{
		db: db,
	}
}

// AddOrUpdateRecentView 添加或更新最近查看记录
// 如果股票代码已存在，则更新查看时间和过期时间
// 过期时间设置为查看时间 + 2天
func (s *RecentViewService) AddOrUpdateRecentView(req *models.AddRecentViewRequest) (*models.RecentView, error) {
	now := time.Now()
	expiresAt := now.Add(48 * time.Hour) // 2天后过期

	// 检查是否已存在该股票的记录
	var existingID string
	err := s.db.QueryRow(`
		SELECT id FROM recent_views WHERE ts_code = ?
	`, req.TSCode).Scan(&existingID)

	if err == sql.ErrNoRows {
		// 不存在，创建新记录
		id := uuid.New().String()
		_, err := s.db.Exec(`
			INSERT INTO recent_views (
				id, ts_code, name, symbol, market,
				viewed_at, expires_at, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			id, req.TSCode, req.Name, req.Symbol, req.Market,
			now, expiresAt, now, now,
		)
		if err != nil {
			return nil, fmt.Errorf("创建最近查看记录失败: %w", err)
		}

		return &models.RecentView{
			ID:        id,
			TSCode:    req.TSCode,
			Name:      req.Name,
			Symbol:    req.Symbol,
			Market:    req.Market,
			ViewedAt:  now,
			ExpiresAt: expiresAt,
			CreatedAt: now,
			UpdatedAt: now,
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("查询最近查看记录失败: %w", err)
	}

	// 已存在，更新查看时间和过期时间
	_, err = s.db.Exec(`
		UPDATE recent_views 
		SET viewed_at = ?, expires_at = ?, updated_at = ?,
		    name = ?, symbol = ?, market = ?
		WHERE id = ?
	`, now, expiresAt, now, req.Name, req.Symbol, req.Market, existingID)
	if err != nil {
		return nil, fmt.Errorf("更新最近查看记录失败: %w", err)
	}

	return &models.RecentView{
		ID:        existingID,
		TSCode:    req.TSCode,
		Name:      req.Name,
		Symbol:    req.Symbol,
		Market:    req.Market,
		ViewedAt:  now,
		ExpiresAt: expiresAt,
		UpdatedAt: now,
	}, nil
}

// GetRecentViews 获取最近查看列表
// limit: 返回记录数量限制，默认20
// includeExpired: 是否包含已过期的记录，默认false
func (s *RecentViewService) GetRecentViews(limit int, includeExpired bool) ([]models.RecentView, error) {
	if limit <= 0 {
		limit = 20
	}

	var query string
	if includeExpired {
		query = `
			SELECT id, ts_code, name, symbol, market, 
			       viewed_at, expires_at, created_at, updated_at
			FROM recent_views
			ORDER BY viewed_at DESC
			LIMIT ?
		`
	} else {
		query = `
			SELECT id, ts_code, name, symbol, market, 
			       viewed_at, expires_at, created_at, updated_at
			FROM recent_views
			WHERE expires_at > ?
			ORDER BY viewed_at DESC
			LIMIT ?
		`
	}

	var rows *sql.Rows
	var err error

	if includeExpired {
		rows, err = s.db.Query(query, limit)
	} else {
		rows, err = s.db.Query(query, time.Now(), limit)
	}

	if err != nil {
		return nil, fmt.Errorf("查询最近查看列表失败: %w", err)
	}
	defer rows.Close()

	var views []models.RecentView
	for rows.Next() {
		var view models.RecentView
		err := rows.Scan(
			&view.ID, &view.TSCode, &view.Name, &view.Symbol, &view.Market,
			&view.ViewedAt, &view.ExpiresAt, &view.CreatedAt, &view.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描最近查看记录失败: %w", err)
		}
		views = append(views, view)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历最近查看记录失败: %w", err)
	}

	return views, nil
}

// DeleteRecentView 删除指定的最近查看记录
func (s *RecentViewService) DeleteRecentView(tsCode string) error {
	result, err := s.db.Exec(`
		DELETE FROM recent_views WHERE ts_code = ?
	`, tsCode)
	if err != nil {
		return fmt.Errorf("删除最近查看记录失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取删除结果失败: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("未找到要删除的记录: %s", tsCode)
	}

	return nil
}

// ClearExpiredViews 清理过期的最近查看记录
// 返回删除的记录数
func (s *RecentViewService) ClearExpiredViews() (int, error) {
	result, err := s.db.Exec(`
		DELETE FROM recent_views WHERE expires_at <= ?
	`, time.Now())
	if err != nil {
		return 0, fmt.Errorf("清理过期记录失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("获取清理结果失败: %w", err)
	}

	if rowsAffected > 0 {
		log.Printf("清理了 %d 条过期的最近查看记录", rowsAffected)
	}

	return int(rowsAffected), nil
}

// ClearAllViews 清空所有最近查看记录
func (s *RecentViewService) ClearAllViews() error {
	_, err := s.db.Exec(`DELETE FROM recent_views`)
	if err != nil {
		return fmt.Errorf("清空最近查看记录失败: %w", err)
	}
	return nil
}

// GetRecentViewCount 获取未过期的最近查看记录数
func (s *RecentViewService) GetRecentViewCount() (int, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM recent_views WHERE expires_at > ?
	`, time.Now()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("获取最近查看记录数失败: %w", err)
	}
	return count, nil
}

// StartAutoCleanup 启动自动清理过期记录的定时任务
// interval: 清理间隔，建议设置为每小时或每天
func (s *RecentViewService) StartAutoCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			count, err := s.ClearExpiredViews()
			if err != nil {
				log.Printf("自动清理过期最近查看记录失败: %v", err)
			} else if count > 0 {
				log.Printf("自动清理了 %d 条过期的最近查看记录", count)
			}
		}
	}()
	log.Printf("已启动最近查看记录自动清理任务，间隔: %v", interval)
}
