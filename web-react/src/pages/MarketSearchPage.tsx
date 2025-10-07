/**
 * 市场与搜索合并页面 - Web版本
 * 整合了市场列表和搜索功能
 */

import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Typography,
  TextField,
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
  InputAdornment,
  Tabs,
  Tab,
} from '@mui/material';
import {
  Search as SearchIcon,
  Refresh as RefreshIcon,
  TrendingUp as TrendingUpIcon,
  FavoriteBorder as FavoriteBorderIcon,
  History as HistoryIcon,
} from '@mui/icons-material';

import { useGetStockListQuery, useGetHealthStatusQuery, useLazySearchStocksQuery } from '../services/api';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { setConnectionStatus } from '../store/slices/appSlice';
import {
  setCurrentQuery,
  setSearchResults,
  addSearchHistory,
  addRecentlyViewed,
} from '../store/slices/searchSlice';
import { StockBasic } from '../types/stock';

const MarketSearchPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  
  // 本地状态
  const [selectedMarket, setSelectedMarket] = useState<'all' | 'sh' | 'sz'>('all');
  const [searchInput, setSearchInput] = useState('');
  const [showSearchResults, setShowSearchResults] = useState(false);
  const [activeTab, setActiveTab] = useState<'market' | 'recent'>('market');
  
  // Redux状态
  const connectionStatus = useAppSelector(state => state.app.connectionStatus);
  const searchResults = useAppSelector(state => state.search.results);
  const searchHistory = useAppSelector(state => state.search.history);
  const recentlyViewed = useAppSelector(state => state.search.recentlyViewed);
  const hotSearches = useAppSelector(state => state.search.hotSearches);
  
  // API查询
  const {
    data: healthData,
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
  
  const [searchStocks, { isLoading: isSearchLoading }] = useLazySearchStocksQuery();
  
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
  
  // 防抖搜索
  const debounceSearch = useCallback((query: string) => {
    const debounced = debounce(async (q: string) => {
      if (q.trim().length >= 1) {
        try {
          const result = await searchStocks({
            q: q.trim(),
            limit: 50,
          }).unwrap();
          
          if (result.success && result.data && result.data.stocks) {
            dispatch(setSearchResults(result.data.stocks));
            setShowSearchResults(true);
          }
        } catch (error) {
          console.error('搜索失败:', error);
        }
      } else {
        dispatch(setSearchResults([]));
        setShowSearchResults(false);
      }
    }, 300);
    
    debounced(query);
  }, [searchStocks, dispatch]);
  
  // 处理搜索输入变化
  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const query = event.target.value;
    setSearchInput(query);
    dispatch(setCurrentQuery(query));
    
    if (query.trim()) {
      debounceSearch(query);
    } else {
      setShowSearchResults(false);
      dispatch(setSearchResults([]));
    }
  };
  
  // 处理搜索提交
  const handleSearchSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    const query = searchInput.trim();
    if (query) {
      dispatch(addSearchHistory({ query }));
      debounceSearch(query);
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
    // 添加到搜索历史
    if (searchInput.trim()) {
      dispatch(addSearchHistory({
        query: searchInput,
        stockCode: stock.ts_code,
        stockName: stock.name,
      }));
    }
    
    // 添加到最近查看
    dispatch(addRecentlyViewed(stock));
    
    // 导航到股票详情页面
    navigate(`/stock/${stock.ts_code}`);
    
    // 清除搜索状态
    setSearchInput('');
    setShowSearchResults(false);
    dispatch(setCurrentQuery(''));
    dispatch(setSearchResults([]));
  };
  
  // 处理热门搜索选择
  const handleHotSearchSelect = (query: string) => {
    setSearchInput(query);
    dispatch(setCurrentQuery(query));
    debounceSearch(query);
  };
  
  // 处理历史搜索选择
  const handleHistorySelect = (historyItem: any) => {
    if (historyItem.stockCode) {
      navigate(`/stock/${historyItem.stockCode}`);
    } else {
      setSearchInput(historyItem.query);
      dispatch(setCurrentQuery(historyItem.query));
      debounceSearch(historyItem.query);
    }
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
  const displayStocks = showSearchResults ? searchResults : filteredStocks.slice(0, 100);
  
  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" gutterBottom>
        市场
      </Typography>
      
      {/* 头部卡片 */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
            <Box>
              <Typography variant="h5" component="h1" gutterBottom>
                Stock-A-Future
              </Typography>
              <Typography variant="subtitle2" color="text.secondary">
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
          
          {/* 搜索框 */}
          <form onSubmit={handleSearchSubmit}>
            <TextField
              fullWidth
              placeholder="搜索股票名称或代码..."
              value={searchInput}
              onChange={handleSearchChange}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
                endAdornment: isSearchLoading && (
                  <InputAdornment position="end">
                    <CircularProgress size={20} />
                  </InputAdornment>
                ),
              }}
              sx={{ mb: 2 }}
            />
          </form>
          
          {/* 热门搜索 */}
          {!showSearchResults && hotSearches.length > 0 && (
            <Box display="flex" flexWrap="wrap" gap={1} mb={2}>
              <Typography variant="caption" color="text.secondary" sx={{ mr: 1, lineHeight: '32px' }}>
                热门:
              </Typography>
              {hotSearches.slice(0, 6).map((query, index) => (
                <Chip
                  key={index}
                  label={query}
                  onClick={() => handleHotSearchSelect(query)}
                  variant="outlined"
                  size="small"
                  clickable
                />
              ))}
            </Box>
          )}
          
          {/* 市场筛选器 - 只在非搜索状态显示 */}
          {!showSearchResults && (
            <Box display="flex" gap={1}>
              <Chip
                label="全部"
                onClick={() => setSelectedMarket('all')}
                color={selectedMarket === 'all' ? 'primary' : 'default'}
                variant={selectedMarket === 'all' ? 'filled' : 'outlined'}
                clickable
              />
              <Chip
                label="上海"
                onClick={() => setSelectedMarket('sh')}
                color={selectedMarket === 'sh' ? 'primary' : 'default'}
                variant={selectedMarket === 'sh' ? 'filled' : 'outlined'}
                clickable
              />
              <Chip
                label="深圳"
                onClick={() => setSelectedMarket('sz')}
                color={selectedMarket === 'sz' ? 'primary' : 'default'}
                variant={selectedMarket === 'sz' ? 'filled' : 'outlined'}
                clickable
              />
            </Box>
          )}
        </CardContent>
      </Card>
      
      {/* 标签切换 - 非搜索状态 */}
      {!showSearchResults && (
        <Box sx={{ mb: 2 }}>
          <Tabs value={activeTab} onChange={(e, val) => setActiveTab(val)}>
            <Tab 
              label={`股票列表 (${filteredStocks.length})`} 
              value="market" 
            />
            <Tab 
              label={`最近查看 (${recentlyViewed.length})`} 
              value="recent" 
            />
          </Tabs>
        </Box>
      )}
      
      {/* 主内容区域 */}
      <Card>
        <CardContent>
          {showSearchResults ? (
            <>
              <Typography variant="h6" gutterBottom>
                搜索结果
                <Chip 
                  label={searchResults.length} 
                  size="small" 
                  color="primary" 
                  sx={{ ml: 2 }}
                />
              </Typography>
              <Divider sx={{ my: 2 }} />
            </>
          ) : (
            <>
              <Typography variant="h6" gutterBottom>
                {activeTab === 'market' ? '股票列表' : '最近查看'}
                {activeTab === 'market' && (
                  <Chip 
                    label={displayStocks.length} 
                    size="small" 
                    color="primary" 
                    sx={{ ml: 2 }}
                  />
                )}
              </Typography>
              <Divider sx={{ my: 2 }} />
            </>
          )}
          
          {/* 股票列表 */}
          {activeTab === 'market' || showSearchResults ? (
            displayStocks.length === 0 ? (
              <Box textAlign="center" py={4}>
                <Typography variant="body1" color="text.secondary">
                  {showSearchResults ? '未找到匹配的股票' : '暂无股票数据'}
                </Typography>
              </Box>
            ) : (
              <List>
                {displayStocks.map((stock, index) => (
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
                              {stock.area && ` · ${stock.area}`}
                            </Typography>
                          }
                        />
                        <IconButton edge="end" sx={{ ml: 1 }}>
                          <FavoriteBorderIcon />
                        </IconButton>
                      </ListItemButton>
                    </ListItem>
                    {index < displayStocks.length - 1 && <Divider />}
                  </React.Fragment>
                ))}
              </List>
            )
          ) : (
            // 最近查看标签
            recentlyViewed.length === 0 ? (
              <Box textAlign="center" py={4}>
                <Typography variant="body1" color="text.secondary">
                  暂无最近查看记录
                </Typography>
              </Box>
            ) : (
              <List>
                {recentlyViewed.slice(0, 20).map((stock, index) => (
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
                          bgcolor="secondary.light"
                          color="secondary.contrastText"
                          mr={2}
                        >
                          <HistoryIcon />
                        </Box>
                        <ListItemText
                          primary={
                            <Typography variant="subtitle1" fontWeight={600}>
                              {stock.name}
                            </Typography>
                          }
                          secondary={
                            <Typography variant="body2" color="text.secondary">
                              {stock.ts_code}
                            </Typography>
                          }
                        />
                      </ListItemButton>
                    </ListItem>
                    {index < recentlyViewed.length - 1 && <Divider />}
                  </React.Fragment>
                ))}
              </List>
            )
          )}
        </CardContent>
      </Card>
      
      {/* 搜索历史卡片 - 非搜索状态时显示 */}
      {!showSearchResults && searchHistory.length > 0 && (
        <Card sx={{ mt: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              搜索历史
            </Typography>
            <Divider sx={{ my: 2 }} />
            <List>
              {searchHistory.slice(0, 5).map((item) => (
                <ListItem key={item.id} disablePadding>
                  <ListItemButton onClick={() => handleHistorySelect(item)}>
                    <HistoryIcon sx={{ mr: 2, color: 'text.secondary' }} />
                    <ListItemText
                      primary={item.stockName || item.query}
                      secondary={item.stockCode}
                    />
                  </ListItemButton>
                </ListItem>
              ))}
            </List>
          </CardContent>
        </Card>
      )}
    </Container>
  );
};

// 防抖函数
function debounce<T extends (...args: any[]) => any>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => func(...args), delay);
  };
}

export default MarketSearchPage;

