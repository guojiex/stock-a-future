/**
 * Stock-A-Future React Web App
 * 主入口文件
 */

import React, { useMemo } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Provider } from 'react-redux';
import { ThemeProvider, CssBaseline, useMediaQuery } from '@mui/material';

import store from './store';
import { getTheme } from './constants/theme';
import Layout from './components/common/Layout';
import MarketSearchPage from './pages/MarketSearchPage';
import FavoritesPage from './pages/FavoritesPage';
import SignalsPage from './pages/SignalsPage';
import StrategiesPage from './pages/StrategiesPage';
import BacktestPage from './pages/BacktestPage';
import SettingsPage from './pages/SettingsPage';
import StockDetailPage from './pages/StockDetailPage';

function AppContent() {
  // 检测系统主题偏好
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
  
  // 创建主题
  const theme = useMemo(
    () => getTheme(prefersDarkMode),
    [prefersDarkMode]
  );

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<MarketSearchPage />} />
            <Route path="stock/:stockCode" element={<StockDetailPage />} />
            <Route path="favorites" element={<FavoritesPage />} />
            <Route path="signals" element={<SignalsPage />} />
            <Route path="strategies" element={<StrategiesPage />} />
            <Route path="backtest" element={<BacktestPage />} />
            <Route path="settings" element={<SettingsPage />} />
          </Route>
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

function App() {
  return (
    <Provider store={store}>
      <AppContent />
    </Provider>
  );
}

export default App;
