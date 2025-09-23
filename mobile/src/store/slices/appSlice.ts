/**
 * 应用全局状态管理
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { appConfig } from '../../constants/config';

// 应用状态接口
interface AppState {
  // 连接状态
  isConnected: boolean;
  connectionStatus: 'connecting' | 'connected' | 'disconnected' | 'error';
  
  // 加载状态
  isLoading: boolean;
  loadingMessage?: string;
  
  // 错误状态
  error?: string;
  
  // 用户偏好
  preferences: {
    theme: 'light' | 'dark' | 'auto';
    currency: 'CNY' | 'USD';
    language: 'zh-CN' | 'en-US';
    refreshInterval: number; // 秒
    chartType: 'candlestick' | 'line';
    showVolume: boolean;
    showIndicators: string[]; // 默认显示的技术指标
  };
  
  // 应用配置
  config: {
    apiBaseUrl: string;
    requestTimeout: number;
    maxRetries: number;
    cacheTimeout: number; // 毫秒
  };
  
  // 最后更新时间
  lastUpdated?: string;
}

// 初始状态 - 从配置管理中获取默认值
const initialState: AppState = {
  isConnected: false,
  connectionStatus: 'disconnected',
  isLoading: false,
  preferences: {
    theme: 'auto',
    currency: 'CNY',
    language: 'zh-CN',
    refreshInterval: appConfig.refreshInterval, // 从配置获取
    chartType: 'candlestick',
    showVolume: true,
    showIndicators: appConfig.defaultIndicators, // 从配置获取默认技术指标
  },
  config: {
    apiBaseUrl: appConfig.apiBaseUrl, // 从配置获取API基础URL
    requestTimeout: appConfig.apiTimeout, // 从配置获取超时时间
    maxRetries: appConfig.maxRetries, // 从配置获取最大重试次数
    cacheTimeout: appConfig.cacheTimeout, // 从配置获取缓存超时时间
  },
};

// 创建切片
const appSlice = createSlice({
  name: 'app',
  initialState,
  reducers: {
    // 设置连接状态
    setConnectionStatus: (state, action: PayloadAction<AppState['connectionStatus']>) => {
      state.connectionStatus = action.payload;
      state.isConnected = action.payload === 'connected';
    },
    
    // 设置加载状态
    setLoading: (state, action: PayloadAction<{ isLoading: boolean; message?: string }>) => {
      state.isLoading = action.payload.isLoading;
      state.loadingMessage = action.payload.message;
    },
    
    // 设置错误
    setError: (state, action: PayloadAction<string | undefined>) => {
      state.error = action.payload;
    },
    
    // 清除错误
    clearError: (state) => {
      state.error = undefined;
    },
    
    // 更新用户偏好
    updatePreferences: (state, action: PayloadAction<Partial<AppState['preferences']>>) => {
      state.preferences = { ...state.preferences, ...action.payload };
    },
    
    // 设置主题
    setTheme: (state, action: PayloadAction<AppState['preferences']['theme']>) => {
      state.preferences.theme = action.payload;
    },
    
    // 设置图表类型
    setChartType: (state, action: PayloadAction<AppState['preferences']['chartType']>) => {
      state.preferences.chartType = action.payload;
    },
    
    // 切换成交量显示
    toggleVolumeDisplay: (state) => {
      state.preferences.showVolume = !state.preferences.showVolume;
    },
    
    // 设置显示的技术指标
    setShowIndicators: (state, action: PayloadAction<string[]>) => {
      state.preferences.showIndicators = action.payload;
    },
    
    // 添加技术指标
    addIndicator: (state, action: PayloadAction<string>) => {
      if (!state.preferences.showIndicators.includes(action.payload)) {
        state.preferences.showIndicators.push(action.payload);
      }
    },
    
    // 移除技术指标
    removeIndicator: (state, action: PayloadAction<string>) => {
      state.preferences.showIndicators = state.preferences.showIndicators.filter(
        indicator => indicator !== action.payload
      );
    },
    
    // 更新配置
    updateConfig: (state, action: PayloadAction<Partial<AppState['config']>>) => {
      state.config = { ...state.config, ...action.payload };
    },

    // 从配置管理器同步配置
    syncConfigFromManager: (state) => {
      state.config = {
        apiBaseUrl: appConfig.apiBaseUrl,
        requestTimeout: appConfig.apiTimeout,
        maxRetries: appConfig.maxRetries,
        cacheTimeout: appConfig.cacheTimeout,
      };
      state.preferences.refreshInterval = appConfig.refreshInterval;
      state.preferences.showIndicators = appConfig.defaultIndicators;
    },
    
    // 设置最后更新时间
    setLastUpdated: (state, action: PayloadAction<string>) => {
      state.lastUpdated = action.payload;
    },
    
    // 重置应用状态
    resetApp: (state) => {
      return { ...initialState, preferences: state.preferences, config: state.config };
    },
  },
});

// 导出actions
export const {
  setConnectionStatus,
  setLoading,
  setError,
  clearError,
  updatePreferences,
  setTheme,
  setChartType,
  toggleVolumeDisplay,
  setShowIndicators,
  addIndicator,
  removeIndicator,
  updateConfig,
  syncConfigFromManager,
  setLastUpdated,
  resetApp,
} = appSlice.actions;

// 选择器
export const selectConnectionStatus = (state: { app: AppState }) => state.app.connectionStatus;
export const selectIsConnected = (state: { app: AppState }) => state.app.isConnected;
export const selectIsLoading = (state: { app: AppState }) => state.app.isLoading;
export const selectLoadingMessage = (state: { app: AppState }) => state.app.loadingMessage;
export const selectError = (state: { app: AppState }) => state.app.error;
export const selectPreferences = (state: { app: AppState }) => state.app.preferences;
export const selectTheme = (state: { app: AppState }) => state.app.preferences.theme;
export const selectChartType = (state: { app: AppState }) => state.app.preferences.chartType;
export const selectShowVolume = (state: { app: AppState }) => state.app.preferences.showVolume;
export const selectShowIndicators = (state: { app: AppState }) => state.app.preferences.showIndicators;
export const selectConfig = (state: { app: AppState }) => state.app.config;
export const selectLastUpdated = (state: { app: AppState }) => state.app.lastUpdated;

// 导出reducer
export default appSlice.reducer;
