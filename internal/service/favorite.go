package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"stock-a-future/internal/models"
	"sync"
	"time"
)

// FavoriteService 收藏股票服务
type FavoriteService struct {
	db             *DatabaseService
	favorites      map[string]*models.FavoriteStock
	groups         map[string]*models.FavoriteGroup
	mutex          sync.RWMutex
	defaultGroupID string
}

// NewFavoriteService 创建收藏股票服务
func NewFavoriteService(dataDir string) (*FavoriteService, error) {
	// 创建数据库服务
	dbService, err := NewDatabaseService(dataDir)
	if err != nil {
		return nil, fmt.Errorf("创建数据库服务失败: %v", err)
	}

	service := &FavoriteService{
		db:             dbService,
		favorites:      make(map[string]*models.FavoriteStock),
		groups:         make(map[string]*models.FavoriteGroup),
		defaultGroupID: "default",
	}

	// 从JSON文件迁移数据到数据库
	favoritesPath := filepath.Join(dataDir, "favorites", "favorites.json")
	groupsPath := filepath.Join(dataDir, "favorites", "groups.json")

	if err := service.db.MigrateFromJSON(favoritesPath, groupsPath); err != nil {
		return nil, fmt.Errorf("数据迁移失败: %v", err)
	}

	// 从数据库加载数据到内存
	if err := service.loadFavoritesFromDB(); err != nil {
		return nil, fmt.Errorf("从数据库加载收藏数据失败: %v", err)
	}
	if err := service.loadGroupsFromDB(); err != nil {
		return nil, fmt.Errorf("从数据库加载分组数据失败: %v", err)
	}

	// 确保默认分组存在
	service.ensureDefaultGroup()

	return service, nil
}

// Close 关闭服务
func (s *FavoriteService) Close() error {
	return s.db.Close()
}

// generateID 生成唯一ID
func (s *FavoriteService) generateID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// AddFavorite 添加收藏股票
func (s *FavoriteService) AddFavorite(request *models.FavoriteStockRequest) (*models.FavoriteStock, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 验证请求参数
	if request.TSCode == "" {
		return nil, fmt.Errorf("股票代码不能为空")
	}
	if request.Name == "" {
		return nil, fmt.Errorf("股票名称不能为空")
	}

	// 确定分组ID
	groupID := request.GroupID
	if groupID == "" {
		groupID = s.defaultGroupID
	}

	// 添加调试日志
	fmt.Printf("添加收藏请求: 股票=%s, 名称=%s, 分组=%s\n", request.TSCode, request.Name, groupID)

	// 检查在指定分组中是否已经收藏（分组级别的重复检查）
	for _, favorite := range s.favorites {
		if favorite.TSCode == request.TSCode && favorite.GroupID == groupID {
			fmt.Printf("检查结果: 股票 %s 已在分组 %s 中收藏\n", request.TSCode, groupID)
			return nil, fmt.Errorf("股票 %s 已在当前分组中收藏", request.TSCode)
		}
	}

	fmt.Printf("检查结果: 股票 %s 在分组 %s 中未收藏，允许添加\n", request.TSCode, groupID)

	// 计算排序顺序（在分组内的最后位置）
	sortOrder := s.getNextSortOrderInGroup(groupID)

	// 创建收藏记录
	favorite := &models.FavoriteStock{
		ID:        s.generateID(),
		TSCode:    request.TSCode,
		Name:      request.Name,
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
		GroupID:   groupID,
		SortOrder: sortOrder,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库
	if err := s.saveFavoriteToDB(favorite); err != nil {
		return nil, fmt.Errorf("保存收藏到数据库失败: %v", err)
	}

	// 添加到内存中
	s.favorites[favorite.ID] = favorite

	fmt.Printf("成功添加收藏: ID=%s, 股票=%s, 分组=%s\n", favorite.ID, favorite.TSCode, favorite.GroupID)
	return favorite, nil
}

// GetFavorites 获取所有收藏股票
func (s *FavoriteService) GetFavorites() []*models.FavoriteStock {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	favorites := make([]*models.FavoriteStock, 0, len(s.favorites))
	for _, favorite := range s.favorites {
		favorites = append(favorites, favorite)
	}

	// 按照分组ID和排序顺序进行稳定排序
	sortFavorites(favorites)

	return favorites
}

// GetFavorite 根据ID获取收藏股票
func (s *FavoriteService) GetFavorite(id string) (*models.FavoriteStock, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	favorite, exists := s.favorites[id]
	if !exists {
		return nil, fmt.Errorf("收藏记录不存在")
	}

	return favorite, nil
}

// UpdateFavorite 更新收藏股票的时间范围
func (s *FavoriteService) UpdateFavorite(id string, request *models.UpdateFavoriteRequest) (*models.FavoriteStock, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	favorite, exists := s.favorites[id]
	if !exists {
		return nil, fmt.Errorf("收藏记录不存在")
	}

	// 更新时间范围
	if request.StartDate != "" {
		favorite.StartDate = request.StartDate
	}
	if request.EndDate != "" {
		favorite.EndDate = request.EndDate
	}
	if request.GroupID != "" {
		favorite.GroupID = request.GroupID
	}
	if request.SortOrder > 0 {
		favorite.SortOrder = request.SortOrder
	}
	favorite.UpdatedAt = time.Now()

	// 保存到数据库
	if err := s.updateFavoriteInDB(favorite); err != nil {
		return nil, fmt.Errorf("更新数据库失败: %v", err)
	}

	return favorite, nil
}

// DeleteFavorite 删除收藏股票
func (s *FavoriteService) DeleteFavorite(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.favorites[id]; !exists {
		return fmt.Errorf("收藏记录不存在")
	}

	// 从数据库中删除
	if err := s.deleteFavoriteFromDB(id); err != nil {
		return fmt.Errorf("从数据库删除失败: %v", err)
	}

	// 从内存中删除
	delete(s.favorites, id)

	return nil
}

// IsFavorite 检查股票是否已收藏（全局检查，不区分分组）
func (s *FavoriteService) IsFavorite(tsCode string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	fmt.Printf("全局收藏检查: 股票=%s\n", tsCode)

	for _, favorite := range s.favorites {
		if favorite.TSCode == tsCode {
			fmt.Printf("全局收藏检查结果: 股票 %s 已在分组 %s 中收藏\n", tsCode, favorite.GroupID)
			return true
		}
	}

	fmt.Printf("全局收藏检查结果: 股票 %s 未收藏\n", tsCode)
	return false
}

// GetFavoriteCount 获取收藏数量
func (s *FavoriteService) GetFavoriteCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.favorites)
}

// === 数据库操作方法 ===

// loadFavoritesFromDB 从数据库加载收藏数据
func (s *FavoriteService) loadFavoritesFromDB() error {
	query := `
		SELECT id, ts_code, name, start_date, end_date, group_id, sort_order, created_at, updated_at
		FROM favorite_stocks
		ORDER BY group_id, created_at DESC
	`

	rows, err := s.db.GetDB().Query(query)
	if err != nil {
		return fmt.Errorf("查询收藏数据失败: %v", err)
	}
	defer rows.Close()

	s.favorites = make(map[string]*models.FavoriteStock)
	for rows.Next() {
		var favorite models.FavoriteStock
		err := rows.Scan(
			&favorite.ID,
			&favorite.TSCode,
			&favorite.Name,
			&favorite.StartDate,
			&favorite.EndDate,
			&favorite.GroupID,
			&favorite.SortOrder,
			&favorite.CreatedAt,
			&favorite.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("扫描收藏数据失败: %v", err)
		}
		s.favorites[favorite.ID] = &favorite
	}

	fmt.Printf("成功从数据库加载 %d 个收藏股票\n", len(s.favorites))
	return nil
}

// saveFavoriteToDB 保存收藏到数据库
func (s *FavoriteService) saveFavoriteToDB(favorite *models.FavoriteStock) error {
	stmt := `
		INSERT INTO favorite_stocks 
		(id, ts_code, name, start_date, end_date, group_id, sort_order, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.GetDB().Exec(stmt,
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

// sortFavorites 对收藏股票进行排序
func sortFavorites(favorites []*models.FavoriteStock) {
	// 按照分组ID和排序顺序进行排序
	// 首先按分组ID排序，然后按sort_order排序
	for i := 0; i < len(favorites)-1; i++ {
		for j := i + 1; j < len(favorites); j++ {
			if shouldSwap(favorites[i], favorites[j]) {
				favorites[i], favorites[j] = favorites[j], favorites[i]
			}
		}
	}
}

// shouldSwap 判断是否需要交换两个收藏股票的位置
func shouldSwap(a, b *models.FavoriteStock) bool {
	// 首先按分组ID排序
	if a.GroupID != b.GroupID {
		return a.GroupID > b.GroupID
	}

	// 分组ID相同时，按排序顺序排序
	if a.SortOrder != b.SortOrder {
		return a.SortOrder > b.SortOrder
	}

	// 排序顺序相同时，按创建时间排序
	return a.CreatedAt.After(b.CreatedAt)
}

// updateFavoriteInDB 更新数据库中的收藏
func (s *FavoriteService) updateFavoriteInDB(favorite *models.FavoriteStock) error {
	stmt := `
		UPDATE favorite_stocks 
		SET start_date = ?, end_date = ?, group_id = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := s.db.GetDB().Exec(stmt,
		favorite.StartDate,
		favorite.EndDate,
		favorite.GroupID,
		favorite.SortOrder,
		favorite.UpdatedAt,
		favorite.ID,
	)

	return err
}

// deleteFavoriteFromDB 从数据库删除收藏
func (s *FavoriteService) deleteFavoriteFromDB(id string) error {
	stmt := `DELETE FROM favorite_stocks WHERE id = ?`
	_, err := s.db.GetDB().Exec(stmt, id)
	return err
}

// === 分组相关方法 ===

// ensureDefaultGroup 确保默认分组存在
func (s *FavoriteService) ensureDefaultGroup() {
	if _, exists := s.groups[s.defaultGroupID]; !exists {
		defaultGroup := &models.FavoriteGroup{
			ID:        s.defaultGroupID,
			Name:      "默认分组",
			Color:     "#3b82f6",
			SortOrder: 0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		s.groups[s.defaultGroupID] = defaultGroup
		s.saveGroupToDB(defaultGroup)
	}
}

// getNextSortOrderInGroup 获取分组内的下一个排序序号
func (s *FavoriteService) getNextSortOrderInGroup(groupID string) int {
	maxOrder := 0
	for _, favorite := range s.favorites {
		if favorite.GroupID == groupID && favorite.SortOrder > maxOrder {
			maxOrder = favorite.SortOrder
		}
	}
	return maxOrder + 1
}

// CreateGroup 创建分组
func (s *FavoriteService) CreateGroup(request *models.CreateGroupRequest) (*models.FavoriteGroup, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 验证请求参数
	if request.Name == "" {
		return nil, fmt.Errorf("分组名称不能为空")
	}

	// 检查分组名是否已存在
	for _, group := range s.groups {
		if group.Name == request.Name {
			return nil, fmt.Errorf("分组名称 %s 已存在", request.Name)
		}
	}

	// 计算排序顺序
	maxOrder := 0
	for _, group := range s.groups {
		if group.SortOrder > maxOrder {
			maxOrder = group.SortOrder
		}
	}

	// 创建分组
	group := &models.FavoriteGroup{
		ID:        s.generateID(),
		Name:      request.Name,
		Color:     request.Color,
		SortOrder: maxOrder + 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库
	if err := s.saveGroupToDB(group); err != nil {
		return nil, fmt.Errorf("保存分组到数据库失败: %v", err)
	}

	// 添加到内存中
	s.groups[group.ID] = group

	return group, nil
}

// GetGroups 获取所有分组
func (s *FavoriteService) GetGroups() []*models.FavoriteGroup {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	groups := make([]*models.FavoriteGroup, 0, len(s.groups))
	for _, group := range s.groups {
		groups = append(groups, group)
	}

	return groups
}

// UpdateGroup 更新分组
func (s *FavoriteService) UpdateGroup(id string, request *models.UpdateGroupRequest) (*models.FavoriteGroup, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	group, exists := s.groups[id]
	if !exists {
		return nil, fmt.Errorf("分组不存在")
	}

	// 如果是默认分组，不允许修改名称
	if id == s.defaultGroupID && request.Name != group.Name {
		return nil, fmt.Errorf("不能修改默认分组的名称")
	}

	// 检查分组名是否已存在（排除自己）
	if request.Name != "" && request.Name != group.Name {
		for _, g := range s.groups {
			if g.ID != id && g.Name == request.Name {
				return nil, fmt.Errorf("分组名称 %s 已存在", request.Name)
			}
		}
	}

	// 更新分组信息
	if request.Name != "" {
		group.Name = request.Name
	}
	if request.Color != "" {
		group.Color = request.Color
	}
	if request.SortOrder > 0 {
		group.SortOrder = request.SortOrder
	}
	group.UpdatedAt = time.Now()

	// 保存到数据库
	if err := s.updateGroupInDB(group); err != nil {
		return nil, fmt.Errorf("更新数据库失败: %v", err)
	}

	return group, nil
}

// DeleteGroup 删除分组
func (s *FavoriteService) DeleteGroup(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 不允许删除默认分组
	if id == s.defaultGroupID {
		return fmt.Errorf("不能删除默认分组")
	}

	if _, exists := s.groups[id]; !exists {
		return fmt.Errorf("分组不存在")
	}

	// 将该分组下的所有收藏移动到默认分组
	for _, favorite := range s.favorites {
		if favorite.GroupID == id {
			favorite.GroupID = s.defaultGroupID
			favorite.UpdatedAt = time.Now()

			// 更新数据库
			if err := s.updateFavoriteInDB(favorite); err != nil {
				return fmt.Errorf("更新收藏分组失败: %v", err)
			}
		}
	}

	// 从数据库中删除分组
	if err := s.deleteGroupFromDB(id); err != nil {
		return fmt.Errorf("从数据库删除分组失败: %v", err)
	}

	// 从内存中删除分组
	delete(s.groups, id)

	return nil
}

// UpdateFavoritesOrder 批量更新收藏排序
func (s *FavoriteService) UpdateFavoritesOrder(request *models.UpdateFavoritesOrderRequest) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 更新每个收藏的分组和排序
	for _, orderItem := range request.FavoriteOrders {
		if favorite, exists := s.favorites[orderItem.ID]; exists {
			favorite.GroupID = orderItem.GroupID
			favorite.SortOrder = orderItem.SortOrder
			favorite.UpdatedAt = time.Now()

			// 更新数据库
			if err := s.updateFavoriteInDB(favorite); err != nil {
				return fmt.Errorf("更新收藏排序失败: %v", err)
			}
		}
	}

	return nil
}

// loadGroupsFromDB 从数据库加载分组数据
func (s *FavoriteService) loadGroupsFromDB() error {
	query := `
		SELECT id, name, color, sort_order, created_at, updated_at
		FROM favorite_groups
		ORDER BY sort_order
	`

	rows, err := s.db.GetDB().Query(query)
	if err != nil {
		return fmt.Errorf("查询分组数据失败: %v", err)
	}
	defer rows.Close()

	s.groups = make(map[string]*models.FavoriteGroup)
	for rows.Next() {
		var group models.FavoriteGroup
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.Color,
			&group.SortOrder,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("扫描分组数据失败: %v", err)
		}
		s.groups[group.ID] = &group
	}

	fmt.Printf("成功从数据库加载 %d 个分组\n", len(s.groups))
	return nil
}

// saveGroupToDB 保存分组到数据库
func (s *FavoriteService) saveGroupToDB(group *models.FavoriteGroup) error {
	stmt := `
		INSERT OR REPLACE INTO favorite_groups 
		(id, name, color, sort_order, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.GetDB().Exec(stmt,
		group.ID,
		group.Name,
		group.Color,
		group.SortOrder,
		group.CreatedAt,
		group.UpdatedAt,
	)

	return err
}

// updateGroupInDB 更新数据库中的分组
func (s *FavoriteService) updateGroupInDB(group *models.FavoriteGroup) error {
	stmt := `
		UPDATE favorite_groups 
		SET name = ?, color = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := s.db.GetDB().Exec(stmt,
		group.Name,
		group.Color,
		group.SortOrder,
		group.UpdatedAt,
		group.ID,
	)

	return err
}

// deleteGroupFromDB 从数据库删除分组
func (s *FavoriteService) deleteGroupFromDB(id string) error {
	stmt := `DELETE FROM favorite_groups WHERE id = ?`
	_, err := s.db.GetDB().Exec(stmt, id)
	return err
}
