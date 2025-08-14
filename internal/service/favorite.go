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
	dataDir   string
	filePath  string
	favorites map[string]*models.FavoriteStock
	mutex     sync.RWMutex
}

// NewFavoriteService 创建收藏股票服务
func NewFavoriteService(dataDir string) *FavoriteService {
	service := &FavoriteService{
		dataDir:   dataDir,
		filePath:  filepath.Join(dataDir, "favorites.json"),
		favorites: make(map[string]*models.FavoriteStock),
	}

	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("创建数据目录失败: %v\n", err)
	}

	// 加载现有数据
	service.loadFavorites()

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

	// 创建收藏记录
	favorite := &models.FavoriteStock{
		ID:        s.generateID(),
		TSCode:    request.TSCode,
		Name:      request.Name,
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
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

	// 加载到内存中
	s.favorites = make(map[string]*models.FavoriteStock)
	for _, favorite := range favoritesList {
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
