/**
 * 收藏股票状态管理 - Web版本
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Favorite, FavoriteGroup } from '../../types/stock';

// 收藏状态接口
interface FavoritesState {
  // 收藏列表
  favorites: Favorite[];
  
  // 分组列表
  groups: FavoriteGroup[];
  
  // 当前选中的分组
  selectedGroupId?: string;  // 改为string
  
  // 排序方式
  sortBy: 'name' | 'code' | 'created_at' | 'updated_at' | 'custom';
  sortOrder: 'asc' | 'desc';
  
  // 加载状态
  isLoading: boolean;
  
  // 错误状态
  error?: string;
  
  // 同步状态
  lastSyncTime?: string;
  needsSync: boolean;
}

// 初始状态
const initialState: FavoritesState = {
  favorites: [],
  groups: [],
  sortBy: 'custom',
  sortOrder: 'asc',
  isLoading: false,
  needsSync: false,
};

// 创建切片
const favoritesSlice = createSlice({
  name: 'favorites',
  initialState,
  reducers: {
    // 设置收藏列表
    setFavorites: (state, action: PayloadAction<Favorite[]>) => {
      state.favorites = action.payload;
      state.isLoading = false;
      state.needsSync = false;
      state.lastSyncTime = new Date().toISOString();
    },
    
    // 添加收藏
    addFavorite: (state, action: PayloadAction<Favorite>) => {
      const favorite = action.payload;
      
      // 检查是否已存在
      const existingIndex = state.favorites.findIndex(
        item => item.ts_code === favorite.ts_code
      );
      
      if (existingIndex === -1) {
        state.favorites.push(favorite);
        state.needsSync = true;
      }
    },
    
    // 删除收藏
    removeFavorite: (state, action: PayloadAction<string>) => {  // 改为string
      const favoriteId = action.payload;
      state.favorites = state.favorites.filter(item => item.id !== favoriteId);
      state.needsSync = true;
    },
    
    // 设置分组列表
    setGroups: (state, action: PayloadAction<FavoriteGroup[]>) => {
      state.groups = action.payload;
    },
    
    // 选择分组
    selectGroup: (state, action: PayloadAction<string | undefined>) => {  // 改为string
      state.selectedGroupId = action.payload;
    },
    
    // 设置排序方式
    setSortBy: (state, action: PayloadAction<FavoritesState['sortBy']>) => {
      state.sortBy = action.payload;
    },
    
    // 设置加载状态
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    
    // 设置错误
    setError: (state, action: PayloadAction<string | undefined>) => {
      state.error = action.payload;
      state.isLoading = false;
    },
    
    // 清除错误
    clearError: (state) => {
      state.error = undefined;
    },
  },
});

// 导出actions
export const {
  setFavorites,
  addFavorite,
  removeFavorite,
  setGroups,
  selectGroup,
  setSortBy,
  setLoading,
  setError,
  clearError,
} = favoritesSlice.actions;

// 导出reducer
export default favoritesSlice.reducer;
