/**
 * 应用主题配置
 */

import { MD3LightTheme, MD3DarkTheme } from 'react-native-paper';

// 股票相关颜色
export const stockColors = {
  // 中国股市习惯：红涨绿跌
  up: '#FF4444',        // 红色 - 上涨
  down: '#00AA44',      // 绿色 - 下跌
  upLight: '#FFE5E5',   // 浅红色背景
  downLight: '#E5F5E5', // 浅绿色背景
  
  // 技术指标颜色
  ma5: '#FF6600',       // 橙色 MA5
  ma10: '#9900FF',      // 紫色 MA10
  ma20: '#0099FF',      // 蓝色 MA20
  ma30: '#FF9900',      // 橙黄色 MA30
  ma60: '#00FFFF',      // 青色 MA60
  
  // MACD颜色
  macdDif: '#FF6600',   // DIF线颜色
  macdDea: '#0099FF',   // DEA线颜色
  macdHistogram: '#CCCCCC', // 柱状图颜色
  
  // RSI颜色
  rsi: '#9900FF',       // RSI线颜色
  rsiOverbought: '#FF4444', // 超买区颜色
  rsiOversold: '#00AA44',   // 超卖区颜色
  
  // 布林带颜色
  bollUpper: '#FF6600', // 上轨颜色
  bollMid: '#0099FF',   // 中轨颜色
  bollLower: '#00AA44', // 下轨颜色
  
  // KDJ颜色
  kdjK: '#FF6600',      // K线颜色
  kdjD: '#0099FF',      // D线颜色
  kdjJ: '#9900FF',      // J线颜色
  
  // 成交量颜色
  volume: '#CCCCCC',    // 成交量柱颜色
  volumeUp: '#FFE5E5',  // 上涨成交量
  volumeDown: '#E5F5E5', // 下跌成交量
};

// 浅色主题
export const lightTheme = {
  ...MD3LightTheme,
  colors: {
    ...MD3LightTheme.colors,
    // 自定义颜色
    primary: '#007AFF',
    secondary: '#5AC8FA',
    accent: '#FF9500',
    
    // 背景颜色
    background: '#FFFFFF',
    surface: '#F8F9FA',
    surfaceVariant: '#F1F3F4',
    
    // 文本颜色
    onBackground: '#1D1D1F',
    onSurface: '#1D1D1F',
    onSurfaceVariant: '#48484A',
    
    // 边框颜色
    outline: '#C7C7CC',
    outlineVariant: '#E5E5EA',
    
    // 错误颜色
    error: '#FF3B30',
    onError: '#FFFFFF',
    errorContainer: '#FFEBEE',
    onErrorContainer: '#B71C1C',
    
    // 成功颜色
    success: '#34C759',
    onSuccess: '#FFFFFF',
    successContainer: '#E8F5E8',
    onSuccessContainer: '#1B5E20',
    
    // 警告颜色
    warning: '#FF9500',
    onWarning: '#FFFFFF',
    warningContainer: '#FFF3E0',
    onWarningContainer: '#E65100',
    
    // 股票颜色
    ...stockColors,
    
    // 图表背景
    chartBackground: '#FFFFFF',
    chartGrid: '#E0E0E0',
    chartText: '#333333',
    
    // 卡片阴影
    shadow: 'rgba(0, 0, 0, 0.1)',
  },
};

// 深色主题
export const darkTheme = {
  ...MD3DarkTheme,
  colors: {
    ...MD3DarkTheme.colors,
    // 自定义颜色
    primary: '#0A84FF',
    secondary: '#64D2FF',
    accent: '#FF9F0A',
    
    // 背景颜色
    background: '#000000',
    surface: '#1C1C1E',
    surfaceVariant: '#2C2C2E',
    
    // 文本颜色
    onBackground: '#FFFFFF',
    onSurface: '#FFFFFF',
    onSurfaceVariant: '#AEAEB2',
    
    // 边框颜色
    outline: '#38383A',
    outlineVariant: '#48484A',
    
    // 错误颜色
    error: '#FF453A',
    onError: '#FFFFFF',
    errorContainer: '#5F2120',
    onErrorContainer: '#FFBAB1',
    
    // 成功颜色
    success: '#30D158',
    onSuccess: '#000000',
    successContainer: '#1E4620',
    onSuccessContainer: '#B7E4C7',
    
    // 警告颜色
    warning: '#FF9F0A',
    onWarning: '#000000',
    warningContainer: '#4D3319',
    onWarningContainer: '#FFCC74',
    
    // 股票颜色
    ...stockColors,
    
    // 图表背景
    chartBackground: '#1C1C1E',
    chartGrid: '#333333',
    chartText: '#FFFFFF',
    
    // 卡片阴影
    shadow: 'rgba(0, 0, 0, 0.3)',
  },
};

// 主题类型
export type AppTheme = typeof lightTheme;

// 导出主题选择函数
export const getTheme = (isDark: boolean): AppTheme => {
  return isDark ? darkTheme : lightTheme;
};

// 主题常量
export const THEME_CONSTANTS = {
  // 圆角
  borderRadius: {
    small: 4,
    medium: 8,
    large: 12,
    xlarge: 16,
  },
  
  // 间距
  spacing: {
    xs: 4,
    sm: 8,
    md: 16,
    lg: 24,
    xl: 32,
    xxl: 48,
  },
  
  // 字体大小
  fontSize: {
    xs: 10,
    sm: 12,
    md: 14,
    lg: 16,
    xl: 18,
    xxl: 20,
    xxxl: 24,
  },
  
  // 字体权重
  fontWeight: {
    light: '300' as const,
    normal: '400' as const,
    medium: '500' as const,
    semibold: '600' as const,
    bold: '700' as const,
  },
  
  // 阴影
  shadow: {
    small: {
      shadowOffset: { width: 0, height: 1 },
      shadowOpacity: 0.1,
      shadowRadius: 2,
      elevation: 2,
    },
    medium: {
      shadowOffset: { width: 0, height: 2 },
      shadowOpacity: 0.15,
      shadowRadius: 4,
      elevation: 4,
    },
    large: {
      shadowOffset: { width: 0, height: 4 },
      shadowOpacity: 0.2,
      shadowRadius: 8,
      elevation: 8,
    },
  },
  
  // 动画持续时间
  animation: {
    fast: 150,
    normal: 250,
    slow: 350,
  },
  
  // 屏幕尺寸断点
  breakpoints: {
    small: 320,
    medium: 768,
    large: 1024,
  },
};
