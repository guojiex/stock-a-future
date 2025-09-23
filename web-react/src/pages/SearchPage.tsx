/**
 * 搜索页面 - Web版本
 */

import React, { useState, useCallback } from 'react';
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
  InputAdornment,
  Divider,
} from '@mui/material';
import {
  Search as SearchIcon,
  TrendingUp as TrendingUpIcon,
  History as HistoryIcon,
} from '@mui/icons-material';

import { useLazySearchStocksQuery } from '../services/api';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import {
  setCurrentQuery,
  setSearchResults,
  addSearchHistory,
  addRecentlyViewed,
} from '../store/slices/searchSlice';
import { StockBasic } from '../types/stock';

const SearchPage: React.FC = () => {
  const dispatch = useAppDispatch();
  
  // Redux状态
  const currentQuery = useAppSelector(state => state.search.currentQuery);
  const searchResults = useAppSelector(state => state.search.results);
  const searchHistory = useAppSelector(state => state.search.history);
  const recentlyViewed = useAppSelector(state => state.search.recentlyViewed);
  const hotSearches = useAppSelector(state => state.search.hotSearches);
  
  // 本地状态
  const [searchInput, setSearchInput] = useState(currentQuery);
  const [showResults, setShowResults] = useState(false);
  
  // API查询
  const [searchStocks, { isLoading: isSearchLoading }] = useLazySearchStocksQuery();
  
  // 防抖搜索
  const debounceSearch = useCallback(
    debounce(async (query: string) => {
      if (query.trim().length >= 1) {
        try {
          const result = await searchStocks({
            q: query.trim(),
            limit: 20,
          }).unwrap();
          
          if (result.success && result.data) {
            dispatch(setSearchResults(result.data));
            setShowResults(true);
          }
        } catch (error) {
          console.error('搜索失败:', error);
        }
      } else {
        dispatch(setSearchResults([]));
        setShowResults(false);
      }
    }, 300),
    [searchStocks, dispatch]
  );
  
  // 处理搜索输入变化
  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const query = event.target.value;
    setSearchInput(query);
    dispatch(setCurrentQuery(query));
    
    if (query.trim()) {
      debounceSearch(query);
    } else {
      setShowResults(false);
      dispatch(setSearchResults([]));
    }
  };
  
  // 处理搜索提交
  const handleSearchSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    const query = searchInput.trim();
    if (query) {
      // 添加到搜索历史
      dispatch(addSearchHistory({ query }));
      debounceSearch(query);
    }
  };
  
  // 处理股票选择
  const handleStockSelect = (stock: StockBasic) => {
    // 添加到搜索历史
    dispatch(addSearchHistory({
      query: searchInput || stock.name,
      stockCode: stock.ts_code,
      stockName: stock.name,
    }));
    
    // 添加到最近查看
    dispatch(addRecentlyViewed(stock));
    
    // TODO: 导航到股票详情页面
    console.log('选择股票:', stock);
    
    // 清除搜索状态
    setSearchInput('');
    setShowResults(false);
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
      // 如果是股票历史，直接导航
      console.log('导航到历史股票:', historyItem);
    } else {
      // 如果是搜索历史，重新搜索
      setSearchInput(historyItem.query);
      dispatch(setCurrentQuery(historyItem.query));
      debounceSearch(historyItem.query);
    }
  };
  
  // 渲染搜索结果
  const renderSearchResults = () => {
    if (!showResults || searchResults.length === 0) return null;
    
    return (
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            搜索结果
          </Typography>
          <List>
            {searchResults.map((stock, index) => (
              <React.Fragment key={stock.ts_code}>
                <ListItem disablePadding>
                  <ListItemButton onClick={() => handleStockSelect(stock)}>
                    <Box
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                      width={32}
                      height={32}
                      borderRadius="50%"
                      bgcolor="primary.light"
                      color="primary.contrastText"
                      mr={2}
                    >
                      <TrendingUpIcon fontSize="small" />
                    </Box>
                    <ListItemText
                      primary={stock.name}
                      secondary={
                        <Box display="flex" alignItems="center" gap={1}>
                          <Typography variant="body2" color="text.secondary">
                            {stock.ts_code}
                          </Typography>
                          {stock.industry && (
                            <Chip label={stock.industry} size="small" variant="outlined" />
                          )}
                        </Box>
                      }
                    />
                  </ListItemButton>
                </ListItem>
                {index < searchResults.length - 1 && <Divider />}
              </React.Fragment>
            ))}
          </List>
        </CardContent>
      </Card>
    );
  };
  
  // 渲染热门搜索
  const renderHotSearches = () => (
    <Card sx={{ mb: 3 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          热门搜索
        </Typography>
        <Box display="flex" flexWrap="wrap" gap={1}>
          {hotSearches.map((query, index) => (
            <Chip
              key={index}
              label={query}
              onClick={() => handleHotSearchSelect(query)}
              variant="outlined"
              clickable
            />
          ))}
        </Box>
      </CardContent>
    </Card>
  );
  
  // 渲染最近查看
  const renderRecentlyViewed = () => {
    if (recentlyViewed.length === 0) return null;
    
    return (
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            最近查看
          </Typography>
          <List>
            {recentlyViewed.slice(0, 5).map((stock, index) => (
              <React.Fragment key={stock.ts_code}>
                <ListItem disablePadding>
                  <ListItemButton onClick={() => handleStockSelect(stock)}>
                    <Box
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                      width={32}
                      height={32}
                      borderRadius="50%"
                      bgcolor="primary.light"
                      color="primary.contrastText"
                      mr={2}
                    >
                      <TrendingUpIcon fontSize="small" />
                    </Box>
                    <ListItemText
                      primary={stock.name}
                      secondary={stock.ts_code}
                    />
                  </ListItemButton>
                </ListItem>
                {index < recentlyViewed.length - 1 && <Divider />}
              </React.Fragment>
            ))}
          </List>
        </CardContent>
      </Card>
    );
  };
  
  // 渲染搜索历史
  const renderSearchHistory = () => {
    if (searchHistory.length === 0) return null;
    
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            搜索历史
          </Typography>
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
    );
  };
  
  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Typography variant="h4" gutterBottom>
        搜索
      </Typography>
      
      {/* 搜索框 */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <form onSubmit={handleSearchSubmit}>
            <TextField
              fullWidth
              placeholder="输入股票名称或代码"
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
        </CardContent>
      </Card>
      
      {/* 搜索结果 */}
      {renderSearchResults()}
      
      {/* 默认内容 */}
      {!showResults && (
        <>
          {renderHotSearches()}
          {renderRecentlyViewed()}
          {renderSearchHistory()}
        </>
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

export default SearchPage;
