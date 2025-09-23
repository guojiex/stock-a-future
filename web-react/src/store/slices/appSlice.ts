/**
 * 应用全局状态管理 - Web版本
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';

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

// 初始状态
const initialState: AppState = {
  isConnected: false,
  connectionStatus: 'disconnected',
  isLoading: false,
  preferences: {
    theme: 'auto',
    language: 'zh-CN',
    refreshInterval: 60, // 60秒
    chartType: 'candlestick',
    showVolume: true,
    showIndicators: ['MA', 'MACD', 'RSI'], // 默认显示移动平均线、MACD、RSI
  },
  config: {
    apiBaseUrl: 'http://localhost:8080/api/v1/',
    requestTimeout: 30000, // 30秒
    maxRetries: 3,
    cacheTimeout: 5 * 60 * 1000, // 5分钟
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
    
    // 更新配置
    updateConfig: (state, action: PayloadAction<Partial<AppState['config']>>) => {
      state.config = { ...state.config, ...action.payload };
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
  updateConfig,
  setLastUpdated,
  resetApp,
} = appSlice.actions;

// 导出reducer
export default appSlice.reducer;
