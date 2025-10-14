/**
 * RTK Query API服务配置 - Web版本
 * 对应Go后端的所有API端点
 */

import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import {
  StockBasic,
  StockDaily,
  TechnicalIndicators,
  FundamentalData,
  PredictionResult,
  PatternSummary,
  Favorite,
  FavoriteGroup,
  AddFavoriteRequest,
  CreateGroupRequest,
  UpdateGroupRequest,
  UpdateFavoriteRequest,
  ApiResponse,
  SearchParams,
} from '../types/stock';
import { appConfig } from '../constants/config';

// API基础配置 - 从配置管理中获取
const API_BASE_URL = appConfig.apiBaseUrl;

// 创建API服务
export const stockApi = createApi({
  reducerPath: 'stockApi',
  baseQuery: fetchBaseQuery({
    baseUrl: API_BASE_URL,
    prepareHeaders: (headers) => {
      headers.set('Content-Type', 'application/json');
      headers.set('Accept', 'application/json');
      return headers;
    },
  }),
  tagTypes: [
    'Stock',
    'StockBasic', 
    'StockDaily',
    'TechnicalIndicators',
    'FundamentalData',
    'Predictions',
    'Favorites',
    'FavoriteGroups',
    'Signals',
    'RecentViews',
    'Strategies'
  ],
  endpoints: (builder) => ({
    // ===== 健康检查 =====
    getHealthStatus: builder.query<ApiResponse, void>({
      query: () => 'health',
    }),

    // ===== 股票基本信息 =====
    getStockBasic: builder.query<ApiResponse<StockBasic>, string>({
      query: (stockCode) => `stocks/${stockCode}/basic`,
      providesTags: (result, error, stockCode) => [
        { type: 'StockBasic', id: stockCode },
      ],
    }),

    // ===== 股票日线数据 =====
    getDailyData: builder.query<ApiResponse<StockDaily[]>, {
      stockCode: string;
      startDate?: string;
      endDate?: string;
      adjust?: 'qfq' | 'hfq' | 'none';
    }>({
      query: ({ stockCode, startDate, endDate, adjust = 'qfq' }) => {
        const params = new URLSearchParams();
        if (startDate) params.append('start_date', startDate);
        if (endDate) params.append('end_date', endDate);
        if (adjust) params.append('adjust', adjust);
        
        return `stocks/${stockCode}/daily?${params.toString()}`;
      },
      providesTags: (result, error, { stockCode }) => [
        { type: 'StockDaily', id: stockCode },
      ],
    }),

    // ===== 技术指标 =====
    getIndicators: builder.query<ApiResponse<TechnicalIndicators>, string>({
      query: (stockCode) => `stocks/${stockCode}/indicators`,
      providesTags: (result, error, stockCode) => [
        { type: 'TechnicalIndicators', id: stockCode },
      ],
    }),

    // ===== 买卖预测 =====
    getPredictions: builder.query<ApiResponse<PredictionResult>, string>({
      query: (stockCode) => `stocks/${stockCode}/predictions`,
      providesTags: (result, error, stockCode) => [
        { type: 'Predictions', id: stockCode },
      ],
    }),

    // ===== 形态识别汇总 =====
    getPatternSummary: builder.query<ApiResponse<PatternSummary>, string>({
      query: (stockCode) => `patterns/summary?ts_code=${stockCode}`,
      providesTags: (result, error, stockCode) => [
        { type: 'Signals', id: `${stockCode}-summary` },
      ],
    }),

    // ===== 基本面数据 =====
    getFundamentalData: builder.query<ApiResponse<FundamentalData>, {
      stockCode: string;
      period?: string;
    }>({
      query: ({ stockCode, period }) => {
        const params = period ? `?period=${period}` : '';
        return `stocks/${stockCode}/fundamental${params}`;
      },
      providesTags: (result, error, { stockCode }) => [
        { type: 'FundamentalData', id: stockCode },
      ],
    }),

    // ===== 股票列表和搜索 =====
    getStockList: builder.query<ApiResponse<{total: number, stocks: StockBasic[]}>, void>({
      query: () => 'stocks',
      providesTags: ['Stock'],
    }),

    searchStocks: builder.query<ApiResponse<{keyword: string, total: number, stocks: StockBasic[]}>, SearchParams>({
      query: ({ q, limit = 20, offset = 0 }) => {
        const params = new URLSearchParams({
          q,
          limit: limit.toString(),
        });
        return `stocks/search?${params.toString()}`;
      },
      providesTags: ['Stock'],
    }),

    // ===== 收藏管理 =====
    getFavorites: builder.query<ApiResponse<{total: number, favorites: Favorite[]}>, void>({
      query: () => 'favorites',
      providesTags: ['Favorites'],
    }),

    addFavorite: builder.mutation<ApiResponse, AddFavoriteRequest>({
      query: (favorite) => ({
        url: 'favorites',
        method: 'POST',
        body: favorite,
      }),
      invalidatesTags: ['Favorites'],
    }),

    deleteFavorite: builder.mutation<ApiResponse, string>({
      query: (id) => ({
        url: `favorites/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Favorites'],
    }),

    updateFavorite: builder.mutation<ApiResponse, { id: string; data: UpdateFavoriteRequest }>({
      query: ({ id, data }) => ({
        url: `favorites/${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['Favorites'],
    }),

    checkFavorite: builder.query<ApiResponse<{ is_favorite: boolean }>, string>({
      query: (stockCode) => `favorites/check/${stockCode}`,
      providesTags: (result, error, stockCode) => [
        { type: 'Favorites', id: stockCode },
      ],
    }),

    // ===== 分组管理 =====
    getGroups: builder.query<ApiResponse<{total: number, groups: FavoriteGroup[]}>, void>({
      query: () => 'groups',
      providesTags: ['FavoriteGroups'],
    }),

    createGroup: builder.mutation<ApiResponse<FavoriteGroup>, CreateGroupRequest>({
      query: (group) => ({
        url: 'groups',
        method: 'POST',
        body: group,
      }),
      invalidatesTags: ['FavoriteGroups'],
    }),

    updateGroup: builder.mutation<ApiResponse<FavoriteGroup>, { id: string; data: UpdateGroupRequest }>({
      query: ({ id, data }) => ({
        url: `groups/${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['FavoriteGroups'],
    }),

    deleteGroup: builder.mutation<ApiResponse, string>({
      query: (id) => ({
        url: `groups/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['FavoriteGroups', 'Favorites'],
    }),

    // ===== 数据刷新 =====
    refreshLocalData: builder.mutation<ApiResponse, void>({
      query: () => ({
        url: 'stocks/refresh',
        method: 'POST',
      }),
      invalidatesTags: ['Stock', 'StockBasic'],
    }),

    // ===== 最近查看 =====
    getRecentViews: builder.query<ApiResponse<{total: number, views: StockBasic[]}>, {
      limit?: number;
      includeExpired?: boolean;
    }>({
      query: ({ limit = 20, includeExpired = false }) => {
        const params = new URLSearchParams({
          limit: limit.toString(),
          include_expired: includeExpired.toString(),
        });
        return `recent-views?${params.toString()}`;
      },
      providesTags: ['RecentViews'],
    }),

    addRecentView: builder.mutation<ApiResponse, {
      ts_code: string;
      name: string;
      symbol?: string;
      market?: string;
    }>({
      query: (view) => ({
        url: 'recent-views',
        method: 'POST',
        body: view,
      }),
      invalidatesTags: ['RecentViews'],
    }),

    deleteRecentView: builder.mutation<ApiResponse, string>({
      query: (tsCode) => ({
        url: `recent-views/${tsCode}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['RecentViews'],
    }),

    clearExpiredRecentViews: builder.mutation<ApiResponse, void>({
      query: () => ({
        url: 'recent-views/cleanup',
        method: 'POST',
      }),
      invalidatesTags: ['RecentViews'],
    }),

    clearAllRecentViews: builder.mutation<ApiResponse, void>({
      query: () => ({
        url: 'recent-views',
        method: 'DELETE',
      }),
      invalidatesTags: ['RecentViews'],
    }),

    // ===== 策略管理 =====
    getStrategies: builder.query<ApiResponse<{total?: number, items?: any[], data?: any[]}>, void>({
      query: () => 'strategies',
      providesTags: ['Strategies'],
    }),

    getStrategy: builder.query<ApiResponse<any>, string>({
      query: (id) => `strategies/${id}`,
      providesTags: (result, error, id) => [
        { type: 'Strategies', id },
      ],
    }),

    createStrategy: builder.mutation<ApiResponse<any>, any>({
      query: (strategy) => ({
        url: 'strategies',
        method: 'POST',
        body: strategy,
      }),
      invalidatesTags: ['Strategies'],
    }),

    updateStrategy: builder.mutation<ApiResponse<any>, { id: string; [key: string]: any }>({
      query: ({ id, ...data }) => ({
        url: `strategies/${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['Strategies'],
    }),

    deleteStrategy: builder.mutation<ApiResponse, string>({
      query: (id) => ({
        url: `strategies/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Strategies'],
    }),

    toggleStrategy: builder.mutation<ApiResponse, string>({
      query: (id) => ({
        url: `strategies/${id}/toggle`,
        method: 'POST',
      }),
      invalidatesTags: ['Strategies'],
    }),

    getStrategyPerformance: builder.query<ApiResponse<any>, string>({
      query: (id) => `strategies/${id}/performance`,
      providesTags: (result, error, id) => [
        { type: 'Strategies', id: `${id}-performance` },
      ],
    }),

    // ===== 回测管理 =====
    startBacktest: builder.mutation<ApiResponse<any>, any>({
      query: (config) => ({
        url: 'backtest',
        method: 'POST',
        body: config,
      }),
    }),

    getBacktestProgress: builder.query<ApiResponse<any>, string>({
      query: (id) => `backtest/${id}/progress`,
    }),

    cancelBacktest: builder.mutation<ApiResponse, string>({
      query: (id) => ({
        url: `backtest/${id}/cancel`,
        method: 'POST',
      }),
    }),

    getBacktestResults: builder.query<ApiResponse<any>, string>({
      query: (id) => `backtest/${id}/results`,
    }),
  }),
});

// 导出hooks
export const {
  // 健康检查
  useGetHealthStatusQuery,
  
  // 股票基本信息
  useGetStockBasicQuery,
  useLazyGetStockBasicQuery,
  
  // 股票日线数据
  useGetDailyDataQuery,
  useLazyGetDailyDataQuery,
  
  // 技术指标
  useGetIndicatorsQuery,
  useLazyGetIndicatorsQuery,
  
  // 买卖预测
  useGetPredictionsQuery,
  useLazyGetPredictionsQuery,
  
  // 形态识别
  useGetPatternSummaryQuery,
  
  // 基本面数据
  useGetFundamentalDataQuery,
  useLazyGetFundamentalDataQuery,
  
  // 股票列表和搜索
  useGetStockListQuery,
  useSearchStocksQuery,
  useLazySearchStocksQuery,
  
  // 收藏管理
  useGetFavoritesQuery,
  useAddFavoriteMutation,
  useUpdateFavoriteMutation,
  useDeleteFavoriteMutation,
  useCheckFavoriteQuery,
  
  // 分组管理
  useGetGroupsQuery,
  useCreateGroupMutation,
  useUpdateGroupMutation,
  useDeleteGroupMutation,
  
  // 数据刷新
  useRefreshLocalDataMutation,
  
  // 最近查看
  useGetRecentViewsQuery,
  useAddRecentViewMutation,
  useDeleteRecentViewMutation,
  useClearExpiredRecentViewsMutation,
  useClearAllRecentViewsMutation,
  
  // 策略管理
  useGetStrategiesQuery,
  useGetStrategyQuery,
  useCreateStrategyMutation,
  useUpdateStrategyMutation,
  useDeleteStrategyMutation,
  useToggleStrategyMutation,
  useGetStrategyPerformanceQuery,
  
  // 回测管理
  useStartBacktestMutation,
  useGetBacktestProgressQuery,
  useLazyGetBacktestProgressQuery,
  useCancelBacktestMutation,
  useGetBacktestResultsQuery,
  useLazyGetBacktestResultsQuery,
} = stockApi;

// 导出API实例
export default stockApi;
