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
  IncomeStatement,
  BalanceSheet,
  CashFlowStatement,
  DailyBasic,
  Favorite,
  AddFavoriteRequest,
  StockGroup,
  ApiResponse,
  PaginatedResponse,
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
    'Favorites',
    'Groups',
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
    getIndicators: builder.query<ApiResponse<TechnicalIndicators[]>, string>({
      query: (stockCode) => `stocks/${stockCode}/indicators`,
      providesTags: (result, error, stockCode) => [
        { type: 'TechnicalIndicators', id: stockCode },
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

    checkFavorite: builder.query<ApiResponse<{ is_favorite: boolean }>, string>({
      query: (stockCode) => `favorites/check/${stockCode}`,
      providesTags: (result, error, stockCode) => [
        { type: 'Favorites', id: stockCode },
      ],
    }),

    // ===== 数据刷新 =====
    refreshLocalData: builder.mutation<ApiResponse, void>({
      query: () => ({
        url: 'stocks/refresh',
        method: 'POST',
      }),
      invalidatesTags: ['Stock', 'StockBasic'],
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
  useDeleteFavoriteMutation,
  useCheckFavoriteQuery,
  
  // 数据刷新
  useRefreshLocalDataMutation,
} = stockApi;

// 导出API实例
export default stockApi;
