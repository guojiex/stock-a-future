/**
 * 市场页面 - Web版本
 */

import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Card,
  CardContent,
  List,
  ListItem,
  ListItemText,
  ListItemButton,
  Chip,
  Box,
  CircularProgress,
  Alert,
  IconButton,
  Divider,
  Badge,
} from '@mui/material';
import {
  Refresh as RefreshIcon,
  TrendingUp as TrendingUpIcon,
  Favorite as FavoriteIcon,
  FavoriteBorder as FavoriteBorderIcon,
} from '@mui/icons-material';

import { useGetStockListQuery, useGetHealthStatusQuery } from '../services/api';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { setConnectionStatus } from '../store/slices/appSlice';
import { StockBasic } from '../types/stock';

const MarketPage: React.FC = () => {
  const dispatch = useAppDispatch();
  
  // 本地状态
  const [selectedMarket, setSelectedMarket] = useState<'all' | 'sh' | 'sz'>('all');
  
  // Redux状态
  const connectionStatus = useAppSelector(state => state.app.connectionStatus);
  
  // API查询
  const {
    data: healthData,
    isLoading: isHealthLoading,
    error: healthError,
    refetch: refetchHealth,
  } = useGetHealthStatusQuery(undefined, {
    pollingInterval: 30000, // 30秒轮询
  });
  
  const {
    data: stockListData,
    isLoading: isStockListLoading,
    error: stockListError,
    refetch: refetchStockList,
  } = useGetStockListQuery();
  
  // 监听健康状态
  useEffect(() => {
    if (healthData?.success) {
      dispatch(setConnectionStatus('connected'));
    } else if (healthError) {
      dispatch(setConnectionStatus('error'));
    }
  }, [healthData, healthError, dispatch]);
  
  // 处理刷新
  const handleRefresh = async () => {
    try {
      await Promise.all([
        refetchHealth(),
        refetchStockList(),
      ]);
    } catch (error) {
      console.error('刷新失败:', error);
    }
  };
  
  // 筛选股票列表
  const getFilteredStocks = (): StockBasic[] => {
    if (!stockListData?.data?.stocks) return [];
    
    const stocks = stockListData.data.stocks;
    
    switch (selectedMarket) {
      case 'sh':
        return stocks.filter(stock => stock.ts_code.endsWith('.SH'));
      case 'sz':
        return stocks.filter(stock => stock.ts_code.endsWith('.SZ'));
      default:
        return stocks;
    }
  };
  
  // 处理股票点击
  const handleStockClick = (stock: StockBasic) => {
    // TODO: 导航到股票详情页面
    console.log('点击股票:', stock);
  };
  
  // 渲染连接状态
  const renderConnectionStatus = () => {
    const getStatusColor = () => {
      switch (connectionStatus) {
        case 'connected':
          return 'success';
        case 'connecting':
          return 'warning';
        case 'error':
          return 'error';
        default:
          return 'default';
      }
    };
    
    const getStatusText = () => {
      switch (connectionStatus) {
        case 'connected':
          return '已连接';
        case 'connecting':
          return '连接中...';
        case 'error':
          return '连接失败';
        default:
          return '未连接';
      }
    };
    
    return (
      <Chip
        label={getStatusText()}
        color={getStatusColor() as any}
        size="small"
        variant="outlined"
      />
    );
  };
  
  // 渲染头部
  const renderHeader = () => (
    <Card sx={{ mb: 3 }}>
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Box>
            <Typography variant="h4" component="h1" gutterBottom>
              Stock-A-Future
            </Typography>
            <Typography variant="subtitle1" color="text.secondary">
              A股股票分析系统
            </Typography>
          </Box>
          <Box display="flex" alignItems="center" gap={2}>
            {renderConnectionStatus()}
            <IconButton onClick={handleRefresh} disabled={isStockListLoading}>
              <RefreshIcon />
            </IconButton>
          </Box>
        </Box>
        
        {/* 市场筛选器 */}
        <Box display="flex" gap={1}>
          <Chip
            label="全部"
            onClick={() => setSelectedMarket('all')}
            color={selectedMarket === 'all' ? 'primary' : 'default'}
            variant={selectedMarket === 'all' ? 'filled' : 'outlined'}
          />
          <Chip
            label="上海"
            onClick={() => setSelectedMarket('sh')}
            color={selectedMarket === 'sh' ? 'primary' : 'default'}
            variant={selectedMarket === 'sh' ? 'filled' : 'outlined'}
          />
          <Chip
            label="深圳"
            onClick={() => setSelectedMarket('sz')}
            color={selectedMarket === 'sz' ? 'primary' : 'default'}
            variant={selectedMarket === 'sz' ? 'filled' : 'outlined'}
          />
        </Box>
      </CardContent>
    </Card>
  );
  
  // 渲染加载状态
  if (isStockListLoading && !stockListData) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography variant="h4" gutterBottom>
          市场
        </Typography>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
          <Box textAlign="center">
            <CircularProgress size={60} />
            <Typography variant="body1" sx={{ mt: 2 }}>
              正在加载股票列表...
            </Typography>
          </Box>
        </Box>
      </Container>
    );
  }
  
  // 渲染错误状态
  if (stockListError && !stockListData) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography variant="h4" gutterBottom>
          市场
        </Typography>
        <Alert 
          severity="error" 
          action={
            <IconButton color="inherit" onClick={() => refetchStockList()}>
              <RefreshIcon />
            </IconButton>
          }
        >
          加载失败: {stockListError.toString()}
        </Alert>
      </Container>
    );
  }
  
  const filteredStocks = getFilteredStocks();
  
  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" gutterBottom>
        市场
      </Typography>
      
      {renderHeader()}
      
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            股票列表
            <Chip 
              label={filteredStocks.length} 
              size="small" 
              color="primary" 
              sx={{ ml: 2 }}
            />
          </Typography>
          
          <Divider sx={{ my: 2 }} />
          
          {filteredStocks.length === 0 ? (
            <Box textAlign="center" py={4}>
              <Typography variant="body1" color="text.secondary">
                暂无股票数据
              </Typography>
            </Box>
          ) : (
            <List>
              {filteredStocks.slice(0, 50).map((stock, index) => (
                <React.Fragment key={stock.ts_code}>
                  <ListItem disablePadding>
                    <ListItemButton onClick={() => handleStockClick(stock)}>
                      <Box
                        display="flex"
                        alignItems="center"
                        justifyContent="center"
                        width={40}
                        height={40}
                        borderRadius="50%"
                        bgcolor="primary.light"
                        color="primary.contrastText"
                        mr={2}
                      >
                        <TrendingUpIcon />
                      </Box>
                      <ListItemText
                        primary={
                          <Box display="flex" alignItems="center" gap={1}>
                            <Typography variant="subtitle1" fontWeight={600}>
                              {stock.name}
                            </Typography>
                            {stock.industry && (
                              <Chip label={stock.industry} size="small" variant="outlined" />
                            )}
                          </Box>
                        }
                        secondary={
                          <Typography variant="body2" color="text.secondary">
                            {stock.ts_code}
                          </Typography>
                        }
                      />
                      <Box textAlign="right">
                        <Typography variant="subtitle1" fontWeight={600}>
                          --
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          --
                        </Typography>
                      </Box>
                      <IconButton edge="end" sx={{ ml: 1 }}>
                        <FavoriteBorderIcon />
                      </IconButton>
                    </ListItemButton>
                  </ListItem>
                  {index < filteredStocks.length - 1 && <Divider />}
                </React.Fragment>
              ))}
            </List>
          )}
        </CardContent>
      </Card>
    </Container>
  );
};

export default MarketPage;
