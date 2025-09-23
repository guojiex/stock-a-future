/**
 * 收藏股票状态管理
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Favorite, StockGroup } from '@/types/stock';

// 收藏状态接口
interface FavoritesState {
  // 收藏列表
  favorites: Favorite[];
  
  // 分组列表
  groups: StockGroup[];
  
  // 当前选中的分组
  selectedGroupId?: number;
  
  // 排序方式
  sortBy: 'name' | 'code' | 'created_at' | 'updated_at' | 'custom';
  sortOrder: 'asc' | 'desc';
  
  // 编辑状态
  isEditing: boolean;
  
  // 加载状态
  isLoading: boolean;
  
  // 错误状态
  error?: string;
  
  // 同步状态
  lastSyncTime?: string;
  needsSync: boolean;
  
  // 显示设置
  displaySettings: {
    showGroups: boolean;
    showNotes: boolean;
    compactView: boolean;
    showPriceChange: boolean;
    autoRefresh: boolean;
    refreshInterval: number; // 秒
  };
}

// 初始状态
const initialState: FavoritesState = {
  favorites: [],
  groups: [],
  sortBy: 'custom',
  sortOrder: 'asc',
  isEditing: false,
  isLoading: false,
  needsSync: false,
  displaySettings: {
    showGroups: true,
    showNotes: true,
    compactView: false,
    showPriceChange: true,
    autoRefresh: true,
    refreshInterval: 60, // 1分钟
  },
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
    removeFavorite: (state, action: PayloadAction<number>) => {
      const favoriteId = action.payload;
      state.favorites = state.favorites.filter(item => item.id !== favoriteId);
      state.needsSync = true;
    },
    
    // 更新收藏
    updateFavorite: (state, action: PayloadAction<{
      id: number;
      updates: Partial<Favorite>;
    }>) => {
      const { id, updates } = action.payload;
      const index = state.favorites.findIndex(item => item.id === id);
      
      if (index !== -1) {
        state.favorites[index] = { ...state.favorites[index], ...updates };
        state.needsSync = true;
      }
    },
    
    // 批量更新收藏顺序
    updateFavoritesOrder: (state, action: PayloadAction<Array<{
      id: number;
      order_index: number;
    }>>) => {
      const orderUpdates = action.payload;
      
      orderUpdates.forEach(({ id, order_index }) => {
        const index = state.favorites.findIndex(item => item.id === id);
        if (index !== -1) {
          state.favorites[index].order_index = order_index;
        }
      });
      
      // 重新排序
      state.favorites.sort((a, b) => (a.order_index || 0) - (b.order_index || 0));
      state.needsSync = true;
    },
    
    // 设置分组列表
    setGroups: (state, action: PayloadAction<StockGroup[]>) => {
      state.groups = action.payload;
    },
    
    // 添加分组
    addGroup: (state, action: PayloadAction<StockGroup>) => {
      state.groups.push(action.payload);
    },
    
    // 删除分组
    removeGroup: (state, action: PayloadAction<number>) => {
      const groupId = action.payload;
      
      // 删除分组
      state.groups = state.groups.filter(group => group.id !== groupId);
      
      // 将该分组下的收藏移到默认分组（group_id = undefined）
      state.favorites.forEach(favorite => {
        if (favorite.group_id === groupId) {
          favorite.group_id = undefined;
        }
      });
      
      // 如果当前选中的是被删除的分组，清除选择
      if (state.selectedGroupId === groupId) {
        state.selectedGroupId = undefined;
      }
      
      state.needsSync = true;
    },
    
    // 更新分组
    updateGroup: (state, action: PayloadAction<{
      id: number;
      updates: Partial<StockGroup>;
    }>) => {
      const { id, updates } = action.payload;
      const index = state.groups.findIndex(group => group.id === id);
      
      if (index !== -1) {
        state.groups[index] = { ...state.groups[index], ...updates };
      }
    },
    
    // 选择分组
    selectGroup: (state, action: PayloadAction<number | undefined>) => {
      state.selectedGroupId = action.payload;
    },
    
    // 设置排序方式
    setSortBy: (state, action: PayloadAction<FavoritesState['sortBy']>) => {
      state.sortBy = action.payload;
    },
    
    // 设置排序顺序
    setSortOrder: (state, action: PayloadAction<FavoritesState['sortOrder']>) => {
      state.sortOrder = action.payload;
    },
    
    // 切换排序顺序
    toggleSortOrder: (state) => {
      state.sortOrder = state.sortOrder === 'asc' ? 'desc' : 'asc';
    },
    
    // 设置编辑状态
    setEditing: (state, action: PayloadAction<boolean>) => {
      state.isEditing = action.payload;
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
    
    // 标记需要同步
    markNeedsSync: (state) => {
      state.needsSync = true;
    },
    
    // 标记同步完成
    markSyncComplete: (state) => {
      state.needsSync = false;
      state.lastSyncTime = new Date().toISOString();
    },
    
    // 更新显示设置
    updateDisplaySettings: (state, action: PayloadAction<Partial<FavoritesState['displaySettings']>>) => {
      state.displaySettings = { ...state.displaySettings, ...action.payload };
    },
    
    // 切换显示设置
    toggleDisplaySetting: (state, action: PayloadAction<keyof FavoritesState['displaySettings']>) => {
      const key = action.payload;
      if (typeof state.displaySettings[key] === 'boolean') {
        (state.displaySettings as any)[key] = !state.displaySettings[key];
      }
    },
    
    // 重置收藏状态
    resetFavorites: (state) => {
      return { ...initialState, displaySettings: state.displaySettings };
    },
  },
});

// 导出actions
export const {
  setFavorites,
  addFavorite,
  removeFavorite,
  updateFavorite,
  updateFavoritesOrder,
  setGroups,
  addGroup,
  removeGroup,
  updateGroup,
  selectGroup,
  setSortBy,
  setSortOrder,
  toggleSortOrder,
  setEditing,
  setLoading,
  setError,
  clearError,
  markNeedsSync,
  markSyncComplete,
  updateDisplaySettings,
  toggleDisplaySetting,
  resetFavorites,
} = favoritesSlice.actions;

// 选择器
export const selectFavorites = (state: { favorites: FavoritesState }) => state.favorites.favorites;
export const selectGroups = (state: { favorites: FavoritesState }) => state.favorites.groups;
export const selectSelectedGroupId = (state: { favorites: FavoritesState }) => state.favorites.selectedGroupId;

// 根据分组筛选收藏
export const selectFavoritesByGroup = (state: { favorites: FavoritesState }) => {
  const { favorites, selectedGroupId } = state.favorites;
  
  if (selectedGroupId === undefined) {
    return favorites.filter(fav => !fav.group_id);
  }
  
  return favorites.filter(fav => fav.group_id === selectedGroupId);
};

// 根据排序方式排序的收藏
export const selectSortedFavorites = (state: { favorites: FavoritesState }) => {
  const { favorites, sortBy, sortOrder } = state.favorites;
  const sorted = [...favorites];
  
  sorted.sort((a, b) => {
    let comparison = 0;
    
    switch (sortBy) {
      case 'name':
        comparison = a.name.localeCompare(b.name);
        break;
      case 'code':
        comparison = a.ts_code.localeCompare(b.ts_code);
        break;
      case 'created_at':
        comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
        break;
      case 'updated_at':
        comparison = new Date(a.updated_at).getTime() - new Date(b.updated_at).getTime();
        break;
      case 'custom':
        comparison = (a.order_index || 0) - (b.order_index || 0);
        break;
    }
    
    return sortOrder === 'asc' ? comparison : -comparison;
  });
  
  return sorted;
};

export const selectSortBy = (state: { favorites: FavoritesState }) => state.favorites.sortBy;
export const selectSortOrder = (state: { favorites: FavoritesState }) => state.favorites.sortOrder;
export const selectIsEditing = (state: { favorites: FavoritesState }) => state.favorites.isEditing;
export const selectIsLoading = (state: { favorites: FavoritesState }) => state.favorites.isLoading;
export const selectError = (state: { favorites: FavoritesState }) => state.favorites.error;
export const selectNeedsSync = (state: { favorites: FavoritesState }) => state.favorites.needsSync;
export const selectLastSyncTime = (state: { favorites: FavoritesState }) => state.favorites.lastSyncTime;
export const selectDisplaySettings = (state: { favorites: FavoritesState }) => state.favorites.displaySettings;

// 检查股票是否已收藏
export const selectIsFavorite = (stockCode: string) => (state: { favorites: FavoritesState }) => {
  return state.favorites.favorites.some(fav => fav.ts_code === stockCode);
};

// 获取收藏的股票
export const selectFavoriteByCode = (stockCode: string) => (state: { favorites: FavoritesState }) => {
  return state.favorites.favorites.find(fav => fav.ts_code === stockCode);
};

// 导出reducer
export default favoritesSlice.reducer;
