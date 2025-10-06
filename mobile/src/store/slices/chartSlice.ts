/**
 * 图表状态管理
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { StockDaily, TechnicalIndicators, ChartDataPoint, CandlestickData } from '@/types/stock';

// 时间范围类型
type TimeRange = '1D' | '5D' | '1M' | '3M' | '6M' | '1Y' | '2Y' | '5Y' | 'ALL';

// 图表类型
type ChartType = 'candlestick' | 'line' | 'area';

// 技术指标类型
type IndicatorType = 'MA' | 'MACD' | 'RSI' | 'BOLL' | 'KDJ' | 'VOL';

// 图表配置
interface ChartConfig {
  type: ChartType;
  timeRange: TimeRange;
  showVolume: boolean;
  showGrid: boolean;
  showCrosshair: boolean;
  enableZoom: boolean;
  enablePan: boolean;
  theme: 'light' | 'dark';
  colors: {
    up: string;      // 上涨颜色
    down: string;    // 下跌颜色
    ma5: string;     // MA5线颜色
    ma10: string;    // MA10线颜色
    ma20: string;    // MA20线颜色
    ma30: string;    // MA30线颜色
    ma60: string;    // MA60线颜色
    volume: string;  // 成交量颜色
    grid: string;    // 网格颜色
    text: string;    // 文字颜色
    background: string; // 背景颜色
  };
}

// 图表数据状态
interface ChartDataState {
  stockCode?: string;
  stockName?: string;
  dailyData: StockDaily[];
  indicators: TechnicalIndicators[];
  processedData: {
    candlestickData: CandlestickData[];
    volumeData: ChartDataPoint[];
    maData: {
      ma5: ChartDataPoint[];
      ma10: ChartDataPoint[];
      ma20: ChartDataPoint[];
      ma30: ChartDataPoint[];
      ma60: ChartDataPoint[];
    };
    macdData: {
      dif: ChartDataPoint[];
      dea: ChartDataPoint[];
      histogram: ChartDataPoint[];
    };
    rsiData: {
      rsi6: ChartDataPoint[];
      rsi12: ChartDataPoint[];
      rsi24: ChartDataPoint[];
    };
    bollData: {
      upper: ChartDataPoint[];
      mid: ChartDataPoint[];
      lower: ChartDataPoint[];
    };
    kdjData: {
      k: ChartDataPoint[];
      d: ChartDataPoint[];
      j: ChartDataPoint[];
    };
  };
  dateRange: {
    start: string;
    end: string;
  };
  lastUpdated?: string;
}

// 图表状态接口
interface ChartState {
  // 图表配置
  config: ChartConfig;
  
  // 当前显示的技术指标
  activeIndicators: IndicatorType[];
  
  // 图表数据
  data: ChartDataState;
  
  // 加载状态
  isLoading: boolean;
  isLoadingIndicators: boolean;
  
  // 错误状态
  error?: string;
  
  // 缩放和平移状态
  zoom: {
    enabled: boolean;
    level: number;
    centerIndex: number;
  };
  
  // 十字线状态
  crosshair: {
    enabled: boolean;
    visible: boolean;
    position: {
      x: number;
      y: number;
    };
    data?: {
      date: string;
      open: number;
      high: number;
      low: number;
      close: number;
      volume: number;
      change: number;
      changePercent: number;
    };
  };
  
  // 图表尺寸
  dimensions: {
    width: number;
    height: number;
    mainChartHeight: number;
    volumeChartHeight: number;
    indicatorChartHeight: number;
  };
}

// 默认图表配置
const defaultChartConfig: ChartConfig = {
  type: 'candlestick',
  timeRange: '3M',
  showVolume: true,
  showGrid: true,
  showCrosshair: true,
  enableZoom: true,
  enablePan: true,
  theme: 'light',
  colors: {
    up: '#FF4444',      // 红色上涨
    down: '#00AA44',    // 绿色下跌
    ma5: '#FF6600',     // 橙色MA5
    ma10: '#9900FF',    // 紫色MA10
    ma20: '#0099FF',    // 蓝色MA20
    ma30: '#FF9900',    // 橙黄色MA30
    ma60: '#00FFFF',    // 青色MA60
    volume: '#CCCCCC',  // 灰色成交量
    grid: '#E0E0E0',    // 浅灰色网格
    text: '#333333',    // 深灰色文字
    background: '#FFFFFF', // 白色背景
  },
};

// 初始状态
const initialState: ChartState = {
  config: defaultChartConfig,
  activeIndicators: ['MA', 'VOL'],
  data: {
    dailyData: [],
    indicators: [],
    processedData: {
      candlestickData: [],
      volumeData: [],
      maData: {
        ma5: [],
        ma10: [],
        ma20: [],
        ma30: [],
        ma60: [],
      },
      macdData: {
        dif: [],
        dea: [],
        histogram: [],
      },
      rsiData: {
        rsi6: [],
        rsi12: [],
        rsi24: [],
      },
      bollData: {
        upper: [],
        mid: [],
        lower: [],
      },
      kdjData: {
        k: [],
        d: [],
        j: [],
      },
    },
    dateRange: {
      start: '',
      end: '',
    },
  },
  isLoading: false,
  isLoadingIndicators: false,
  zoom: {
    enabled: true,
    level: 1,
    centerIndex: 0,
  },
  crosshair: {
    enabled: true,
    visible: false,
    position: {
      x: 0,
      y: 0,
    },
  },
  dimensions: {
    width: 350,
    height: 400,
    mainChartHeight: 250,
    volumeChartHeight: 80,
    indicatorChartHeight: 70,
  },
};

// 创建切片
const chartSlice = createSlice({
  name: 'chart',
  initialState,
  reducers: {
    // 设置图表配置
    setChartConfig: (state, action: PayloadAction<Partial<ChartConfig>>) => {
      state.config = { ...state.config, ...action.payload };
    },
    
    // 设置图表类型
    setChartType: (state, action: PayloadAction<ChartType>) => {
      state.config.type = action.payload;
    },
    
    // 设置时间范围
    setTimeRange: (state, action: PayloadAction<TimeRange>) => {
      state.config.timeRange = action.payload;
    },
    
    // 切换成交量显示
    toggleVolume: (state) => {
      state.config.showVolume = !state.config.showVolume;
    },
    
    // 切换网格显示
    toggleGrid: (state) => {
      state.config.showGrid = !state.config.showGrid;
    },
    
    // 切换十字线
    toggleCrosshair: (state) => {
      state.config.showCrosshair = !state.config.showCrosshair;
      state.crosshair.enabled = state.config.showCrosshair;
    },
    
    // 设置主题
    setTheme: (state, action: PayloadAction<'light' | 'dark'>) => {
      state.config.theme = action.payload;
      
      // 更新主题相关颜色
      if (action.payload === 'dark') {
        state.config.colors = {
          ...state.config.colors,
          grid: '#333333',
          text: '#FFFFFF',
          background: '#1A1A1A',
        };
      } else {
        state.config.colors = {
          ...state.config.colors,
          grid: '#E0E0E0',
          text: '#333333',
          background: '#FFFFFF',
        };
      }
    },
    
    // 设置活跃的技术指标
    setActiveIndicators: (state, action: PayloadAction<IndicatorType[]>) => {
      state.activeIndicators = action.payload;
    },
    
    // 添加技术指标
    addIndicator: (state, action: PayloadAction<IndicatorType>) => {
      if (!state.activeIndicators.includes(action.payload)) {
        state.activeIndicators.push(action.payload);
      }
    },
    
    // 移除技术指标
    removeIndicator: (state, action: PayloadAction<IndicatorType>) => {
      state.activeIndicators = state.activeIndicators.filter(
        indicator => indicator !== action.payload
      );
    },
    
    // 设置股票数据
    setStockData: (state, action: PayloadAction<{
      stockCode: string;
      stockName: string;
      dailyData: StockDaily[];
      dateRange: { start: string; end: string };
    }>) => {
      const { stockCode, stockName, dailyData, dateRange } = action.payload;
      
      state.data.stockCode = stockCode;
      state.data.stockName = stockName;
      state.data.dailyData = dailyData;
      state.data.dateRange = dateRange;
      state.data.lastUpdated = new Date().toISOString();
      
      // 处理K线数据
      state.data.processedData.candlestickData = dailyData.map((item, index) => ({
        x: index,
        date: item.trade_date,
        open: parseFloat(item.open),
        high: parseFloat(item.high),
        low: parseFloat(item.low),
        close: parseFloat(item.close),
        volume: parseFloat(item.vol),
        y: parseFloat(item.close), // 用于线图
      }));
      
      // 处理成交量数据
      state.data.processedData.volumeData = dailyData.map((item, index) => ({
        x: index,
        y: parseFloat(item.vol),
        date: item.trade_date,
      }));
      
      state.isLoading = false;
    },
    
    // 设置技术指标数据
    setIndicatorsData: (state, action: PayloadAction<TechnicalIndicators[]>) => {
      const indicators = action.payload;
      state.data.indicators = indicators;
      
      // 处理移动平均线数据
      state.data.processedData.maData = {
        ma5: indicators.map((item, index) => ({
          x: index,
          y: item.ma5 || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
        ma10: indicators.map((item, index) => ({
          x: index,
          y: item.ma10 || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
        ma20: indicators.map((item, index) => ({
          x: index,
          y: item.ma20 || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
        ma30: indicators.map((item, index) => ({
          x: index,
          y: item.ma30 || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
        ma60: indicators.map((item, index) => ({
          x: index,
          y: item.ma60 || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
      };
      
      // 处理MACD数据
      state.data.processedData.macdData = {
        dif: indicators.map((item, index) => ({
          x: index,
          y: item.macd_dif || 0,
          date: item.trade_date,
        })),
        dea: indicators.map((item, index) => ({
          x: index,
          y: item.macd_dea || 0,
          date: item.trade_date,
        })),
        histogram: indicators.map((item, index) => ({
          x: index,
          y: item.macd_histogram || 0,
          date: item.trade_date,
        })),
      };
      
      // 处理RSI数据
      state.data.processedData.rsiData = {
        rsi6: indicators.map((item, index) => ({
          x: index,
          y: item.rsi6 || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
        rsi12: indicators.map((item, index) => ({
          x: index,
          y: item.rsi12 || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
        rsi24: indicators.map((item, index) => ({
          x: index,
          y: item.rsi24 || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
      };
      
      // 处理布林带数据
      state.data.processedData.bollData = {
        upper: indicators.map((item, index) => ({
          x: index,
          y: item.boll_upper || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
        mid: indicators.map((item, index) => ({
          x: index,
          y: item.boll_mid || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
        lower: indicators.map((item, index) => ({
          x: index,
          y: item.boll_lower || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
      };
      
      // 处理KDJ数据
      state.data.processedData.kdjData = {
        k: indicators.map((item, index) => ({
          x: index,
          y: item.kdj_k || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
        d: indicators.map((item, index) => ({
          x: index,
          y: item.kdj_d || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
        j: indicators.map((item, index) => ({
          x: index,
          y: item.kdj_j || 0,
          date: item.trade_date,
        })).filter(item => item.y > 0),
      };
      
      state.isLoadingIndicators = false;
    },
    
    // 设置加载状态
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    
    // 设置指标加载状态
    setLoadingIndicators: (state, action: PayloadAction<boolean>) => {
      state.isLoadingIndicators = action.payload;
    },
    
    // 设置错误
    setError: (state, action: PayloadAction<string | undefined>) => {
      state.error = action.payload;
      state.isLoading = false;
      state.isLoadingIndicators = false;
    },
    
    // 清除错误
    clearError: (state) => {
      state.error = undefined;
    },
    
    // 设置缩放
    setZoom: (state, action: PayloadAction<{ level: number; centerIndex: number }>) => {
      state.zoom.level = action.payload.level;
      state.zoom.centerIndex = action.payload.centerIndex;
    },
    
    // 设置十字线位置
    setCrosshairPosition: (state, action: PayloadAction<{ x: number; y: number }>) => {
      state.crosshair.position = action.payload;
    },
    
    // 设置十字线可见性
    setCrosshairVisible: (state, action: PayloadAction<boolean>) => {
      state.crosshair.visible = action.payload;
    },
    
    // 设置十字线数据
    setCrosshairData: (state, action: PayloadAction<ChartState['crosshair']['data']>) => {
      state.crosshair.data = action.payload;
    },
    
    // 设置图表尺寸
    setDimensions: (state, action: PayloadAction<Partial<ChartState['dimensions']>>) => {
      state.dimensions = { ...state.dimensions, ...action.payload };
    },
    
    // 清除图表数据
    clearChartData: (state) => {
      state.data = initialState.data;
      state.error = undefined;
      state.isLoading = false;
      state.isLoadingIndicators = false;
    },
    
    // 重置图表状态
    resetChart: (state) => {
      return { ...initialState, config: state.config };
    },
  },
});

// 导出actions
export const {
  setChartConfig,
  setChartType,
  setTimeRange,
  toggleVolume,
  toggleGrid,
  toggleCrosshair,
  setTheme,
  setActiveIndicators,
  addIndicator,
  removeIndicator,
  setStockData,
  setIndicatorsData,
  setLoading,
  setLoadingIndicators,
  setError,
  clearError,
  setZoom,
  setCrosshairPosition,
  setCrosshairVisible,
  setCrosshairData,
  setDimensions,
  clearChartData,
  resetChart,
} = chartSlice.actions;

// 选择器
export const selectChartConfig = (state: { chart: ChartState }) => state.chart.config;
export const selectChartType = (state: { chart: ChartState }) => state.chart.config.type;
export const selectTimeRange = (state: { chart: ChartState }) => state.chart.config.timeRange;
export const selectShowVolume = (state: { chart: ChartState }) => state.chart.config.showVolume;
export const selectActiveIndicators = (state: { chart: ChartState }) => state.chart.activeIndicators;
export const selectChartData = (state: { chart: ChartState }) => state.chart.data;
export const selectCandlestickData = (state: { chart: ChartState }) => state.chart.data.processedData.candlestickData;
export const selectVolumeData = (state: { chart: ChartState }) => state.chart.data.processedData.volumeData;
export const selectMAData = (state: { chart: ChartState }) => state.chart.data.processedData.maData;
export const selectMACDData = (state: { chart: ChartState }) => state.chart.data.processedData.macdData;
export const selectRSIData = (state: { chart: ChartState }) => state.chart.data.processedData.rsiData;
export const selectBOLLData = (state: { chart: ChartState }) => state.chart.data.processedData.bollData;
export const selectKDJData = (state: { chart: ChartState }) => state.chart.data.processedData.kdjData;
export const selectIsLoading = (state: { chart: ChartState }) => state.chart.isLoading;
export const selectIsLoadingIndicators = (state: { chart: ChartState }) => state.chart.isLoadingIndicators;
export const selectChartError = (state: { chart: ChartState }) => state.chart.error;
export const selectZoom = (state: { chart: ChartState }) => state.chart.zoom;
export const selectCrosshair = (state: { chart: ChartState }) => state.chart.crosshair;
export const selectDimensions = (state: { chart: ChartState }) => state.chart.dimensions;

// 导出reducer
export default chartSlice.reducer;
