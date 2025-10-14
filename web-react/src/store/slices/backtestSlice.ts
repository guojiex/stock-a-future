/**
 * 回测状态管理 - Web版本
 * 管理回测配置和选中的策略
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface BacktestState {
  // 选中的策略ID列表
  selectedStrategyIds: string[];
  // 回测配置
  config: {
    name: string;
    startDate: string;
    endDate: string;
    initialCash: number;
    commission: number;
    symbols: string[];
  };
}

const initialState: BacktestState = {
  selectedStrategyIds: [],
  config: {
    name: '',
    startDate: '',
    endDate: '',
    initialCash: 1000000,
    commission: 0.0003,
    symbols: [],
  },
};

const backtestSlice = createSlice({
  name: 'backtest',
  initialState,
  reducers: {
    // 设置选中的策略
    setSelectedStrategies: (state, action: PayloadAction<string[]>) => {
      state.selectedStrategyIds = action.payload;
    },
    // 添加策略
    addStrategy: (state, action: PayloadAction<string>) => {
      if (!state.selectedStrategyIds.includes(action.payload)) {
        state.selectedStrategyIds.push(action.payload);
      }
    },
    // 移除策略
    removeStrategy: (state, action: PayloadAction<string>) => {
      state.selectedStrategyIds = state.selectedStrategyIds.filter(
        (id) => id !== action.payload
      );
    },
    // 清空策略选择
    clearStrategies: (state) => {
      state.selectedStrategyIds = [];
    },
    // 更新回测配置
    updateConfig: (state, action: PayloadAction<Partial<BacktestState['config']>>) => {
      state.config = { ...state.config, ...action.payload };
    },
    // 重置回测状态
    resetBacktest: (state) => {
      state.selectedStrategyIds = [];
      state.config = initialState.config;
    },
  },
});

export const {
  setSelectedStrategies,
  addStrategy,
  removeStrategy,
  clearStrategies,
  updateConfig,
  resetBacktest,
} = backtestSlice.actions;

export default backtestSlice.reducer;

