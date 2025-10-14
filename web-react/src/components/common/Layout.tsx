/**
 * 应用布局组件 - Web版本
 */

import React from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import {
  Box,
  AppBar,
  Toolbar,
  Typography,
  BottomNavigation,
  BottomNavigationAction,
  Paper,
  useTheme,
  useMediaQuery,
  Tabs,
  Tab,
} from '@mui/material';
import {
  TrendingUp as MarketIcon,
  Star as FavoritesIcon,
  AccountTree as StrategiesIcon,
  Settings as SettingsIcon,
} from '@mui/icons-material';

const Layout: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  
  // 导航项配置
  const navItems = [
    { label: '市场', value: '/', icon: <MarketIcon /> },
    { label: '收藏', value: '/favorites', icon: <FavoritesIcon /> },
    { label: '策略', value: '/strategies', icon: <StrategiesIcon /> },
    { label: '设置', value: '/settings', icon: <SettingsIcon /> },
  ];
  
  // 获取当前路径对应的导航值
  const getCurrentNavValue = () => {
    return navItems.find(item => item.value === location.pathname)?.value || '/';
  };
  
  // 处理导航变化
  const handleNavChange = (event: React.SyntheticEvent, newValue: string) => {
    navigate(newValue);
  };
  
  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      {/* 顶部应用栏 */}
      <AppBar position="static" elevation={1}>
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Stock-A-Future
          </Typography>
          
          {/* 桌面端：顶部标签导航 */}
          {!isMobile && (
            <Tabs
              value={getCurrentNavValue()}
              onChange={handleNavChange}
              textColor="inherit"
              indicatorColor="secondary"
              sx={{ 
                ml: 'auto',
                '& .MuiTab-root': { 
                  color: 'rgba(255, 255, 255, 0.7)',
                  minHeight: 64,
                },
                '& .Mui-selected': {
                  color: 'white',
                }
              }}
            >
              {navItems.map((item) => (
                <Tab
                  key={item.value}
                  label={item.label}
                  value={item.value}
                  icon={item.icon}
                  iconPosition="start"
                />
              ))}
            </Tabs>
          )}
          
          {/* 移动端：显示版本信息 */}
          {isMobile && (
            <Typography variant="body2" color="inherit" sx={{ opacity: 0.8 }}>
              Web版
            </Typography>
          )}
        </Toolbar>
      </AppBar>
      
      {/* 主内容区域 */}
      <Box sx={{ flex: 1, pb: isMobile ? 7 : 0 }}>
        <Outlet />
      </Box>
      
      {/* 移动端：底部导航栏 */}
      {isMobile && (
        <Paper 
          sx={{ position: 'fixed', bottom: 0, left: 0, right: 0 }} 
          elevation={3}
        >
          <BottomNavigation
            value={getCurrentNavValue()}
            onChange={handleNavChange}
            showLabels
          >
            {navItems.map((item) => (
              <BottomNavigationAction
                key={item.value}
                label={item.label}
                value={item.value}
                icon={item.icon}
              />
            ))}
          </BottomNavigation>
        </Paper>
      )}
    </Box>
  );
};

export default Layout;
