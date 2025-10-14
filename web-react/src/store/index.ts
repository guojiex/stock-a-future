/**
 * Redux Store配置 - Web版本
 * 使用Redux Toolkit和RTK Query
 */

import { configureStore } from '@reduxjs/toolkit';
import { setupListeners } from '@reduxjs/toolkit/query';
import { stockApi } from '../services/api';
import appSlice from './slices/appSlice';
import searchSlice from './slices/searchSlice';
import favoritesSlice from './slices/favoritesSlice';
import backtestSlice from './slices/backtestSlice';

// 配置store
export const store = configureStore({
  reducer: {
    // RTK Query API
    [stockApi.reducerPath]: stockApi.reducer,
    
    // 应用状态
    app: appSlice,
    search: searchSlice,
    favorites: favoritesSlice,
    backtest: backtestSlice,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // 忽略RTK Query的action类型
        ignoredActions: [
          'persist/PERSIST',
          'persist/REHYDRATE',
        ],
      },
    }).concat(stockApi.middleware),
  devTools: process.env.NODE_ENV !== 'production', // 开发环境启用DevTools
});

// 设置监听器以启用缓存、无效化、轮询等功能
setupListeners(store.dispatch);

// 导出类型
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

// 导出store
export default store;
