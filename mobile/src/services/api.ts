/**
 * RTK Query API服务配置
 * 对应Go后端的所有API端点
 */

import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import {
  StockBasic,
  StockDaily,
  TechnicalIndicators,
  FundamentalData,
  IncomeStatement,
  BalanceSheet,
  CashFlowStatement,
  DailyBasic,
  Favorite,
  AddFavoriteRequest,
  StockGroup,
  PredictionData,
  ApiResponse,
  PaginatedResponse,
  SearchParams,
  DateRange,
} from '@/types/stock';
import { appConfig } from '@/constants/config';

// API基础配置 - 从配置管理中获取
const API_BASE_URL = appConfig.apiBaseUrl;

// 创建API服务
export const stockApi = createApi({
  reducerPath: 'stockApi',
  baseQuery: fetchBaseQuery({
    baseUrl: API_BASE_URL,
    timeout: appConfig.apiTimeout,
    prepareHeaders: (headers) => {
      headers.set('Content-Type', 'application/json');
      headers.set('Accept', 'application/json');
      return headers;
    },
    // 错误处理
    responseHandler: async (response) => {
      if (response.ok) {
        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/json')) {
          return response.json();
        }
        return response.text();
      }
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    },
  }),
  tagTypes: [
    'Stock',
    'StockBasic', 
    'StockDaily',
    'TechnicalIndicators',
    'FundamentalData',
    'Favorites',
    'Groups',
    'Predictions',
    'Signals'
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
      adjust?: 'qfq' | 'hfq' | 'none'; // 复权方式
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
    getIndicators: builder.query<ApiResponse<TechnicalIndicators[]>, string>({
      query: (stockCode) => `stocks/${stockCode}/indicators`,
      providesTags: (result, error, stockCode) => [
        { type: 'TechnicalIndicators', id: stockCode },
      ],
    }),

    // ===== 预测数据 =====
    getPredictions: builder.query<ApiResponse<PredictionData[]>, string>({
      query: (stockCode) => `stocks/${stockCode}/predictions`,
      providesTags: (result, error, stockCode) => [
        { type: 'Predictions', id: stockCode },
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

    // ===== 利润表 =====
    getIncomeStatement: builder.query<ApiResponse<IncomeStatement[]>, {
      stockCode: string;
      period?: string;
    }>({
      query: ({ stockCode, period }) => {
        const params = period ? `?period=${period}` : '';
        return `stocks/${stockCode}/income${params}`;
      },
      providesTags: (result, error, { stockCode }) => [
        { type: 'FundamentalData', id: `${stockCode}-income` },
      ],
    }),

    // ===== 资产负债表 =====
    getBalanceSheet: builder.query<ApiResponse<BalanceSheet[]>, {
      stockCode: string;
      period?: string;
    }>({
      query: ({ stockCode, period }) => {
        const params = period ? `?period=${period}` : '';
        return `stocks/${stockCode}/balance${params}`;
      },
      providesTags: (result, error, { stockCode }) => [
        { type: 'FundamentalData', id: `${stockCode}-balance` },
      ],
    }),

    // ===== 现金流量表 =====
    getCashFlowStatement: builder.query<ApiResponse<CashFlowStatement[]>, {
      stockCode: string;
      period?: string;
    }>({
      query: ({ stockCode, period }) => {
        const params = period ? `?period=${period}` : '';
        return `stocks/${stockCode}/cashflow${params}`;
      },
      providesTags: (result, error, { stockCode }) => [
        { type: 'FundamentalData', id: `${stockCode}-cashflow` },
      ],
    }),

    // ===== 每日基本面 =====
    getDailyBasic: builder.query<ApiResponse<DailyBasic>, {
      stockCode: string;
      tradeDate?: string;
    }>({
      query: ({ stockCode, tradeDate }) => {
        const params = tradeDate ? `?trade_date=${tradeDate}` : '';
        return `stocks/${stockCode}/dailybasic${params}`;
      },
      providesTags: (result, error, { stockCode }) => [
        { type: 'FundamentalData', id: `${stockCode}-dailybasic` },
      ],
    }),

    // ===== 股票列表和搜索 =====
    getStockList: builder.query<ApiResponse<PaginatedResponse<StockBasic>>, {
      page?: number;
      pageSize?: number;
    }>({
      query: ({ page = 1, pageSize = 50 }) => 
        `stocks?page=${page}&page_size=${pageSize}`,
      providesTags: ['Stock'],
    }),

    searchStocks: builder.query<ApiResponse<StockBasic[]>, SearchParams>({
      query: ({ q, limit = 20, offset = 0 }) => {
        const params = new URLSearchParams({
          q,
          limit: limit.toString(),
          offset: offset.toString(),
        });
        return `stocks/search?${params.toString()}`;
      },
      providesTags: ['Stock'],
    }),

    // ===== 数据刷新 =====
    refreshLocalData: builder.mutation<ApiResponse, void>({
      query: () => ({
        url: 'stocks/refresh',
        method: 'POST',
      }),
      invalidatesTags: ['Stock', 'StockBasic'],
    }),

    // ===== 缓存管理 =====
    getCacheStats: builder.query<ApiResponse, void>({
      query: () => 'cache/stats',
    }),

    clearCache: builder.mutation<ApiResponse, void>({
      query: () => ({
        url: 'cache',
        method: 'DELETE',
      }),
      invalidatesTags: ['Stock', 'StockDaily', 'TechnicalIndicators', 'FundamentalData'],
    }),

    // ===== 收藏管理 =====
    getFavorites: builder.query<ApiResponse<Favorite[]>, void>({
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

    deleteFavorite: builder.mutation<ApiResponse, number>({
      query: (id) => ({
        url: `favorites/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Favorites'],
    }),

    updateFavorite: builder.mutation<ApiResponse, { 
      id: number; 
      data: Partial<AddFavoriteRequest> 
    }>({
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

    updateFavoritesOrder: builder.mutation<ApiResponse, { 
      favorites: Array<{ id: number; order_index: number }> 
    }>({
      query: (data) => ({
        url: 'favorites/order',
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['Favorites'],
    }),

    getFavoritesSignals: builder.query<ApiResponse, void>({
      query: () => 'favorites/signals',
      providesTags: ['Favorites', 'Signals'],
    }),

    // ===== 分组管理 =====
    getGroups: builder.query<ApiResponse<StockGroup[]>, void>({
      query: () => 'groups',
      providesTags: ['Groups'],
    }),

    createGroup: builder.mutation<ApiResponse, { 
      name: string; 
      description?: string 
    }>({
      query: (group) => ({
        url: 'groups',
        method: 'POST',
        body: group,
      }),
      invalidatesTags: ['Groups'],
    }),

    updateGroup: builder.mutation<ApiResponse, { 
      id: number; 
      data: { name?: string; description?: string } 
    }>({
      query: ({ id, data }) => ({
        url: `groups/${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['Groups'],
    }),

    deleteGroup: builder.mutation<ApiResponse, number>({
      query: (id) => ({
        url: `groups/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Groups', 'Favorites'],
    }),

    // ===== 图形识别和信号 =====
    recognizePatterns: builder.query<ApiResponse, { 
      tsCode: string;
      patterns?: string[];
    }>({
      query: ({ tsCode, patterns }) => {
        const params = new URLSearchParams({ ts_code: tsCode });
        if (patterns && patterns.length > 0) {
          params.append('patterns', patterns.join(','));
        }
        return `patterns/recognize?${params.toString()}`;
      },
      providesTags: (result, error, { tsCode }) => [
        { type: 'Signals', id: tsCode },
      ],
    }),

    searchPatterns: builder.mutation<ApiResponse, {
      patterns: string[];
      start_date?: string;
      end_date?: string;
      min_confidence?: number;
    }>({
      query: (data) => ({
        url: 'patterns/search',
        method: 'POST',
        body: data,
      }),
    }),

    getPatternSummary: builder.query<ApiResponse, string>({
      query: (tsCode) => `patterns/summary?ts_code=${tsCode}`,
      providesTags: (result, error, tsCode) => [
        { type: 'Signals', id: `${tsCode}-summary` },
      ],
    }),

    getRecentSignals: builder.query<ApiResponse, string>({
      query: (tsCode) => `patterns/recent?ts_code=${tsCode}`,
      providesTags: (result, error, tsCode) => [
        { type: 'Signals', id: `${tsCode}-recent` },
      ],
    }),

    getAvailablePatterns: builder.query<ApiResponse<string[]>, void>({
      query: () => 'patterns/available',
    }),

    getPatternStatistics: builder.query<ApiResponse, string>({
      query: (tsCode) => `patterns/statistics?ts_code=${tsCode}`,
      providesTags: (result, error, tsCode) => [
        { type: 'Signals', id: `${tsCode}-stats` },
      ],
    }),

    // ===== 信号计算 =====
    calculateSignal: builder.mutation<ApiResponse, {
      ts_code: string;
      indicators?: string[];
      start_date?: string;
      end_date?: string;
    }>({
      query: (data) => ({
        url: 'signals/calculate',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: (result, error, { ts_code }) => [
        { type: 'Signals', id: ts_code },
      ],
    }),

    batchCalculateSignals: builder.mutation<ApiResponse, {
      ts_codes: string[];
      indicators?: string[];
      start_date?: string;
      end_date?: string;
    }>({
      query: (data) => ({
        url: 'signals/batch',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Signals'],
    }),

    getSignal: builder.query<ApiResponse, string>({
      query: (stockCode) => `signals/${stockCode}`,
      providesTags: (result, error, stockCode) => [
        { type: 'Signals', id: stockCode },
      ],
    }),

    getLatestSignals: builder.query<ApiResponse, {
      limit?: number;
      offset?: number;
    }>({
      query: ({ limit = 20, offset = 0 }) => {
        const params = new URLSearchParams({
          limit: limit.toString(),
          offset: offset.toString(),
        });
        return `signals?${params.toString()}`;
      },
      providesTags: ['Signals'],
    }),

    getCalculationStatus: builder.query<ApiResponse, void>({
      query: () => 'signals/status',
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
  
  // 预测数据
  useGetPredictionsQuery,
  useLazyGetPredictionsQuery,
  
  // 基本面数据
  useGetFundamentalDataQuery,
  useLazyGetFundamentalDataQuery,
  useGetIncomeStatementQuery,
  useGetBalanceSheetQuery,
  useGetCashFlowStatementQuery,
  useGetDailyBasicQuery,
  
  // 股票列表和搜索
  useGetStockListQuery,
  useSearchStocksQuery,
  useLazySearchStocksQuery,
  
  // 数据刷新
  useRefreshLocalDataMutation,
  
  // 缓存管理
  useGetCacheStatsQuery,
  useClearCacheMutation,
  
  // 收藏管理
  useGetFavoritesQuery,
  useAddFavoriteMutation,
  useDeleteFavoriteMutation,
  useUpdateFavoriteMutation,
  useCheckFavoriteQuery,
  useUpdateFavoritesOrderMutation,
  useGetFavoritesSignalsQuery,
  
  // 分组管理
  useGetGroupsQuery,
  useCreateGroupMutation,
  useUpdateGroupMutation,
  useDeleteGroupMutation,
  
  // 图形识别和信号
  useRecognizePatternsQuery,
  useSearchPatternsMutation,
  useGetPatternSummaryQuery,
  useGetRecentSignalsQuery,
  useGetAvailablePatternsQuery,
  useGetPatternStatisticsQuery,
  
  // 信号计算
  useCalculateSignalMutation,
  useBatchCalculateSignalsMutation,
  useGetSignalQuery,
  useGetLatestSignalsQuery,
  useGetCalculationStatusQuery,
} = stockApi;

// 导出API实例
export default stockApi;
