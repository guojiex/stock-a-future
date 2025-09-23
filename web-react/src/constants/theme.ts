/**
 * 主题配置 - Web版本
 */

import { createTheme, ThemeOptions } from '@mui/material/styles';

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
  
  // 成交量颜色
  volume: '#CCCCCC',    // 成交量柱颜色
};

// 浅色主题配置
const lightThemeOptions: ThemeOptions = {
  palette: {
    mode: 'light',
    primary: {
      main: '#007AFF',
      light: '#5AC8FA',
      dark: '#0056B3',
    },
    secondary: {
      main: '#FF9500',
      light: '#FFCC74',
      dark: '#E65100',
    },
    background: {
      default: '#FFFFFF',
      paper: '#F8F9FA',
    },
    text: {
      primary: '#1D1D1F',
      secondary: '#48484A',
    },
    error: {
      main: '#FF3B30',
    },
    success: {
      main: '#34C759',
    },
    warning: {
      main: '#FF9500',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", "PingFang SC", "Microsoft YaHei", sans-serif',
    h1: {
      fontSize: '2.5rem',
      fontWeight: 700,
    },
    h2: {
      fontSize: '2rem',
      fontWeight: 600,
    },
    h3: {
      fontSize: '1.5rem',
      fontWeight: 600,
    },
    h4: {
      fontSize: '1.25rem',
      fontWeight: 600,
    },
    h5: {
      fontSize: '1.125rem',
      fontWeight: 600,
    },
    h6: {
      fontSize: '1rem',
      fontWeight: 600,
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiCard: {
      styleOverrides: {
        root: {
          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
          borderRadius: 12,
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 8,
          fontWeight: 500,
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 16,
        },
      },
    },
  },
};

// 深色主题配置
const darkThemeOptions: ThemeOptions = {
  palette: {
    mode: 'dark',
    primary: {
      main: '#0A84FF',
      light: '#64D2FF',
      dark: '#0056B3',
    },
    secondary: {
      main: '#FF9F0A',
      light: '#FFCC74',
      dark: '#E65100',
    },
    background: {
      default: '#000000',
      paper: '#1C1C1E',
    },
    text: {
      primary: '#FFFFFF',
      secondary: '#AEAEB2',
    },
    error: {
      main: '#FF453A',
    },
    success: {
      main: '#30D158',
    },
    warning: {
      main: '#FF9F0A',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", "PingFang SC", "Microsoft YaHei", sans-serif',
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiCard: {
      styleOverrides: {
        root: {
          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.3)',
          borderRadius: 12,
          backgroundColor: '#1C1C1E',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 8,
          fontWeight: 500,
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 16,
        },
      },
    },
  },
};

// 创建主题实例
export const lightTheme = createTheme(lightThemeOptions);
export const darkTheme = createTheme(darkThemeOptions);

// 主题选择函数
export const getTheme = (isDark: boolean) => {
  return isDark ? darkTheme : lightTheme;
};

// 扩展主题类型以包含股票颜色
declare module '@mui/material/styles' {
  interface Theme {
    stockColors: typeof stockColors;
  }
  
  interface ThemeOptions {
    stockColors?: typeof stockColors;
  }
}

// 将股票颜色添加到主题中
lightTheme.stockColors = stockColors;
darkTheme.stockColors = stockColors;
