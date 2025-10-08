/**
 * 搜索状态管理 - Web版本
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { StockBasic } from '../../types/stock';

// 搜索历史项
interface SearchHistoryItem {
  id: string;
  query: string;
  stockCode?: string;
  stockName?: string;
  timestamp: string;
}

// 搜索状态接口
interface SearchState {
  // 当前搜索查询
  currentQuery: string;
  
  // 搜索结果
  results: StockBasic[];
  isSearching: boolean;
  
  // 搜索历史
  history: SearchHistoryItem[];
  maxHistoryItems: number;
  
  // 热门搜索
  hotSearches: string[];
  
  // 最近查看的股票
  recentlyViewed: StockBasic[];
  maxRecentItems: number;
}

// 初始状态
const initialState: SearchState = {
  currentQuery: '',
  results: [],
  isSearching: false,
  history: [],
  maxHistoryItems: 20,
  hotSearches: [
    '平安银行', '茅台', '招商银行', '工商银行', '建设银行',
    '中国平安', '五粮液', '中国移动', '腾讯控股', '阿里巴巴'
  ],
  recentlyViewed: [],
  maxRecentItems: 10,
};

// 创建切片
const searchSlice = createSlice({
  name: 'search',
  initialState,
  reducers: {
    // 设置当前查询
    setCurrentQuery: (state, action: PayloadAction<string>) => {
      state.currentQuery = action.payload;
    },
    
    // 设置搜索状态
    setSearching: (state, action: PayloadAction<boolean>) => {
      state.isSearching = action.payload;
    },
    
    // 设置搜索结果
    setSearchResults: (state, action: PayloadAction<StockBasic[]>) => {
      state.results = action.payload;
      state.isSearching = false;
    },
    
    // 清除搜索结果
    clearSearchResults: (state) => {
      state.results = [];
      state.currentQuery = '';
      state.isSearching = false;
    },
    
    // 添加搜索历史
    addSearchHistory: (state, action: PayloadAction<{
      query: string;
      stockCode?: string;
      stockName?: string;
    }>) => {
      const { query, stockCode, stockName } = action.payload;
      
      // 检查是否已存在相同的搜索记录
      const existingIndex = state.history.findIndex(
        item => item.query === query || (stockCode && item.stockCode === stockCode)
      );
      
      // 如果存在，先移除旧记录
      if (existingIndex !== -1) {
        state.history.splice(existingIndex, 1);
      }
      
      // 添加新记录到开头
      state.history.unshift({
        id: Date.now().toString(),
        query,
        stockCode,
        stockName,
        timestamp: new Date().toISOString(),
      });
      
      // 限制历史记录数量
      if (state.history.length > state.maxHistoryItems) {
        state.history = state.history.slice(0, state.maxHistoryItems);
      }
    },
    
    // 添加最近查看的股票
    addRecentlyViewed: (state, action: PayloadAction<StockBasic>) => {
      const stock = action.payload;
      
      // 检查是否已存在
      const existingIndex = state.recentlyViewed.findIndex(
        item => item.ts_code === stock.ts_code
      );
      
      // 如果存在，先移除
      if (existingIndex !== -1) {
        state.recentlyViewed.splice(existingIndex, 1);
      }
      
      // 添加到开头
      state.recentlyViewed.unshift(stock);
      
      // 限制数量
      if (state.recentlyViewed.length > state.maxRecentItems) {
        state.recentlyViewed = state.recentlyViewed.slice(0, state.maxRecentItems);
      }
    },
    
    // 清除搜索历史
    clearSearchHistory: (state) => {
      state.history = [];
    },
    
    // 设置最近查看列表（从后端获取）
    setRecentlyViewed: (state, action: PayloadAction<StockBasic[]>) => {
      state.recentlyViewed = action.payload;
    },
    
    // 清除最近查看
    clearRecentlyViewed: (state) => {
      state.recentlyViewed = [];
    },
  },
});

// 导出actions
export const {
  setCurrentQuery,
  setSearching,
  setSearchResults,
  clearSearchResults,
  addSearchHistory,
  addRecentlyViewed,
  setRecentlyViewed,
  clearSearchHistory,
  clearRecentlyViewed,
} = searchSlice.actions;

// 导出reducer
export default searchSlice.reducer;
