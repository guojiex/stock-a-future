package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"stock-a-future/internal/models"
	"sync"
	"time"
)

// FavoriteService 收藏股票服务
type FavoriteService struct {
	dataDir        string
	filePath       string
	groupsPath     string
	favorites      map[string]*models.FavoriteStock
	groups         map[string]*models.FavoriteGroup
	mutex          sync.RWMutex
	defaultGroupID string
}

// NewFavoriteService 创建收藏股票服务
func NewFavoriteService(dataDir string) *FavoriteService {
	// 在 data 目录下创建 favorites 子目录
	favoritesDir := filepath.Join(dataDir, "favorites")

	service := &FavoriteService{
		dataDir:        favoritesDir,
		filePath:       filepath.Join(favoritesDir, "favorites.json"),
		groupsPath:     filepath.Join(favoritesDir, "groups.json"),
		favorites:      make(map[string]*models.FavoriteStock),
		groups:         make(map[string]*models.FavoriteGroup),
		defaultGroupID: "default",
	}

	// 确保收藏数据目录存在
	if err := os.MkdirAll(favoritesDir, 0755); err != nil {
		fmt.Printf("创建收藏数据目录失败: %v\n", err)
	}

	// 加载现有数据
	service.loadFavorites()
	service.loadGroups()
	service.ensureDefaultGroup()

	return service
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

	// 检查是否已经收藏
	for _, favorite := range s.favorites {
		if favorite.TSCode == request.TSCode {
			return nil, fmt.Errorf("股票 %s 已经在收藏列表中", request.TSCode)
		}
	}

	// 确定分组ID
	groupID := request.GroupID
	if groupID == "" {
		groupID = s.defaultGroupID
	}

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

	// 添加到内存中
	s.favorites[favorite.ID] = favorite

	// 保存到文件
	if err := s.saveFavorites(); err != nil {
		// 如果保存失败，从内存中移除
		delete(s.favorites, favorite.ID)
		return nil, fmt.Errorf("保存收藏失败: %v", err)
	}

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
	favorite.UpdatedAt = time.Now()

	// 保存到文件
	if err := s.saveFavorites(); err != nil {
		return nil, fmt.Errorf("保存收藏更新失败: %v", err)
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

	// 从内存中删除
	delete(s.favorites, id)

	// 保存到文件
	if err := s.saveFavorites(); err != nil {
		return fmt.Errorf("保存删除操作失败: %v", err)
	}

	return nil
}

// IsFavorite 检查股票是否已收藏
func (s *FavoriteService) IsFavorite(tsCode string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, favorite := range s.favorites {
		if favorite.TSCode == tsCode {
			return true
		}
	}

	return false
}

// GetFavoriteCount 获取收藏数量
func (s *FavoriteService) GetFavoriteCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.favorites)
}

// loadFavorites 从文件加载收藏数据
func (s *FavoriteService) loadFavorites() error {
	// 检查文件是否存在
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return nil // 文件不存在是正常的，表示还没有收藏数据
	}

	// 读取文件内容
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("读取收藏文件失败: %v", err)
	}

	// 如果文件为空，直接返回
	if len(data) == 0 {
		return nil
	}

	// 解析JSON数据
	var favoritesList []*models.FavoriteStock
	if err := json.Unmarshal(data, &favoritesList); err != nil {
		return fmt.Errorf("解析收藏数据失败: %v", err)
	}

	// 加载到内存中，并为旧数据设置默认值
	s.favorites = make(map[string]*models.FavoriteStock)
	for _, favorite := range favoritesList {
		// 为旧数据设置默认分组和排序
		if favorite.GroupID == "" {
			favorite.GroupID = s.defaultGroupID
		}
		if favorite.SortOrder == 0 {
			favorite.SortOrder = s.getNextSortOrderInGroup(favorite.GroupID)
		}
		s.favorites[favorite.ID] = favorite
	}

	fmt.Printf("成功加载 %d 个收藏股票\n", len(s.favorites))
	return nil
}

// saveFavorites 保存收藏数据到文件
func (s *FavoriteService) saveFavorites() error {
	// 转换为切片
	favoritesList := make([]*models.FavoriteStock, 0, len(s.favorites))
	for _, favorite := range s.favorites {
		favoritesList = append(favoritesList, favorite)
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(favoritesList, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化收藏数据失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("写入收藏文件失败: %v", err)
	}

	return nil
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
		s.saveGroups()
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

	// 添加到内存中
	s.groups[group.ID] = group

	// 保存到文件
	if err := s.saveGroups(); err != nil {
		// 如果保存失败，从内存中移除
		delete(s.groups, group.ID)
		return nil, fmt.Errorf("保存分组失败: %v", err)
	}

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

	// 保存到文件
	if err := s.saveGroups(); err != nil {
		return nil, fmt.Errorf("保存分组更新失败: %v", err)
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
		}
	}

	// 从内存中删除分组
	delete(s.groups, id)

	// 保存分组和收藏数据
	if err := s.saveGroups(); err != nil {
		return fmt.Errorf("保存分组删除操作失败: %v", err)
	}
	if err := s.saveFavorites(); err != nil {
		return fmt.Errorf("保存收藏数据失败: %v", err)
	}

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
		}
	}

	// 保存到文件
	if err := s.saveFavorites(); err != nil {
		return fmt.Errorf("保存排序更新失败: %v", err)
	}

	return nil
}

// loadGroups 从文件加载分组数据
func (s *FavoriteService) loadGroups() error {
	if _, err := os.Stat(s.groupsPath); os.IsNotExist(err) {
		fmt.Printf("分组文件不存在，将创建默认分组: %s\n", s.groupsPath)
		return nil
	}

	data, err := os.ReadFile(s.groupsPath)
	if err != nil {
		return fmt.Errorf("读取分组文件失败: %v", err)
	}

	var groupsList []*models.FavoriteGroup
	if err := json.Unmarshal(data, &groupsList); err != nil {
		return fmt.Errorf("解析分组数据失败: %v", err)
	}

	// 加载到内存
	for _, group := range groupsList {
		s.groups[group.ID] = group
	}

	fmt.Printf("成功加载 %d 个分组\n", len(s.groups))
	return nil
}

// saveGroups 保存分组数据到文件
func (s *FavoriteService) saveGroups() error {
	// 转换为切片
	groupsList := make([]*models.FavoriteGroup, 0, len(s.groups))
	for _, group := range s.groups {
		groupsList = append(groupsList, group)
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(groupsList, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化分组数据失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(s.groupsPath, data, 0644); err != nil {
		return fmt.Errorf("写入分组文件失败: %v", err)
	}

	return nil
}
