/**
 * 搜索页面 - 股票搜索功能
 */

import React, { useState, useEffect, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  Keyboard,
} from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import {
  Appbar,
  Searchbar,
  Card,
  List,
  Chip,
  Divider,
  ActivityIndicator,
  useTheme,
} from 'react-native-paper';
import Icon from 'react-native-vector-icons/MaterialIcons';

import { SearchStackParamList } from '@/navigation/AppNavigator';
import { useLazySearchStocksQuery } from '@/services/api';
import { StockBasic } from '@/types/stock';
import { useAppDispatch, useAppSelector } from '@/hooks/redux';
import {
  setCurrentQuery,
  setSearchResults,
  addSearchHistory,
  addRecentlyViewed,
  selectCurrentQuery,
  selectSearchResults,
  selectIsSearching,
  selectSearchHistory,
  selectRecentlyViewed,
  selectHotSearches,
} from '@/store/slices/searchSlice';

type NavigationProp = NativeStackNavigationProp<SearchStackParamList, 'SearchHome'>;

const SearchScreen: React.FC = () => {
  const navigation = useNavigation<NavigationProp>();
  const theme = useTheme();
  const dispatch = useAppDispatch();
  
  // Redux状态
  const currentQuery = useAppSelector(selectCurrentQuery);
  const searchResults = useAppSelector(selectSearchResults);
  const isSearching = useAppSelector(selectIsSearching);
  const searchHistory = useAppSelector(selectSearchHistory);
  const recentlyViewed = useAppSelector(selectRecentlyViewed);
  const hotSearches = useAppSelector(selectHotSearches);
  
  // 本地状态
  const [searchInput, setSearchInput] = useState(currentQuery);
  const [showSuggestions, setShowSuggestions] = useState(false);
  
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
          }
        } catch (error) {
          console.error('搜索失败:', error);
        }
      } else {
        dispatch(setSearchResults([]));
      }
    }, 300),
    [searchStocks, dispatch]
  );
  
  // 处理搜索输入变化
  const handleSearchChange = (query: string) => {
    setSearchInput(query);
    dispatch(setCurrentQuery(query));
    
    if (query.trim()) {
      setShowSuggestions(true);
      debounceSearch(query);
    } else {
      setShowSuggestions(false);
      dispatch(setSearchResults([]));
    }
  };
  
  // 处理搜索提交
  const handleSearchSubmit = () => {
    const query = searchInput.trim();
    if (query) {
      // 添加到搜索历史
      dispatch(addSearchHistory({ query }));
      
      // 隐藏建议
      setShowSuggestions(false);
      Keyboard.dismiss();
      
      // 如果有结果，导航到结果页面
      if (searchResults.length > 0) {
        navigation.navigate('SearchResult', { query });
      }
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
    
    // 导航到股票详情页面
    // 注意：这里需要导航到Market栈中的StockDetail页面
    // 在实际实现中，可能需要使用navigate('MarketTab', { screen: 'StockDetail', params: ... })
    console.log('导航到股票详情:', stock);
    
    // 清除搜索状态
    setSearchInput('');
    setShowSuggestions(false);
    dispatch(setCurrentQuery(''));
    dispatch(setSearchResults([]));
    Keyboard.dismiss();
  };
  
  // 处理热门搜索选择
  const handleHotSearchSelect = (query: string) => {
    setSearchInput(query);
    handleSearchChange(query);
  };
  
  // 处理历史搜索选择
  const handleHistorySelect = (historyItem: any) => {
    if (historyItem.stockCode) {
      // 如果是股票历史，直接导航
      console.log('导航到历史股票:', historyItem);
    } else {
      // 如果是搜索历史，重新搜索
      setSearchInput(historyItem.query);
      handleSearchChange(historyItem.query);
    }
  };
  
  // 渲染搜索结果项
  const renderSearchResultItem = ({ item }: { item: StockBasic }) => (
    <TouchableOpacity
      onPress={() => handleStockSelect(item)}
      activeOpacity={0.7}
    >
      <List.Item
        title={item.name}
        description={`${item.ts_code} ${item.industry || ''}`}
        left={(props) => (
          <View style={styles.stockIconContainer}>
            <Icon
              name="trending-up"
              size={20}
              color={theme.colors.primary}
              {...props}
            />
          </View>
        )}
        right={(props) => (
          <Icon
            name="chevron-right"
            size={20}
            color={theme.colors.onSurfaceVariant}
            {...props}
          />
        )}
        titleStyle={styles.resultTitle}
        descriptionStyle={styles.resultDescription}
      />
      <Divider />
    </TouchableOpacity>
  );
  
  // 渲染热门搜索
  const renderHotSearches = () => (
    <View style={styles.section}>
      <Text style={[styles.sectionTitle, { color: theme.colors.onSurface }]}>
        热门搜索
      </Text>
      <View style={styles.chipContainer}>
        {hotSearches.map((query, index) => (
          <Chip
            key={index}
            onPress={() => handleHotSearchSelect(query)}
            style={styles.chip}
            textStyle={styles.chipText}
          >
            {query}
          </Chip>
        ))}
      </View>
    </View>
  );
  
  // 渲染搜索历史
  const renderSearchHistory = () => {
    if (searchHistory.length === 0) return null;
    
    return (
      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={[styles.sectionTitle, { color: theme.colors.onSurface }]}>
            搜索历史
          </Text>
          <TouchableOpacity>
            <Icon
              name="clear-all"
              size={20}
              color={theme.colors.onSurfaceVariant}
            />
          </TouchableOpacity>
        </View>
        
        {searchHistory.slice(0, 5).map((item) => (
          <TouchableOpacity
            key={item.id}
            onPress={() => handleHistorySelect(item)}
            activeOpacity={0.7}
          >
            <List.Item
              title={item.stockName || item.query}
              description={item.stockCode}
              left={(props) => (
                <Icon
                  name="history"
                  size={20}
                  color={theme.colors.onSurfaceVariant}
                  {...props}
                />
              )}
              titleStyle={styles.historyTitle}
              descriptionStyle={styles.historyDescription}
            />
          </TouchableOpacity>
        ))}
      </View>
    );
  };
  
  // 渲染最近查看
  const renderRecentlyViewed = () => {
    if (recentlyViewed.length === 0) return null;
    
    return (
      <View style={styles.section}>
        <Text style={[styles.sectionTitle, { color: theme.colors.onSurface }]}>
          最近查看
        </Text>
        
        {recentlyViewed.slice(0, 5).map((stock) => (
          <TouchableOpacity
            key={stock.ts_code}
            onPress={() => handleStockSelect(stock)}
            activeOpacity={0.7}
          >
            <List.Item
              title={stock.name}
              description={stock.ts_code}
              left={(props) => (
                <View style={styles.stockIconContainer}>
                  <Icon
                    name="trending-up"
                    size={20}
                    color={theme.colors.primary}
                    {...props}
                  />
                </View>
              )}
              titleStyle={styles.recentTitle}
              descriptionStyle={styles.recentDescription}
            />
          </TouchableOpacity>
        ))}
      </View>
    );
  };
  
  return (
    <View style={styles.container}>
      <Appbar.Header>
        <Appbar.Content title="搜索" />
      </Appbar.Header>
      
      <View style={styles.searchContainer}>
        <Searchbar
          placeholder="输入股票名称或代码"
          value={searchInput}
          onChangeText={handleSearchChange}
          onSubmitEditing={handleSearchSubmit}
          style={styles.searchbar}
          inputStyle={styles.searchInput}
          loading={isSearchLoading}
        />
      </View>
      
      {showSuggestions && searchResults.length > 0 ? (
        // 显示搜索结果
        <View style={styles.resultsContainer}>
          <FlatList
            data={searchResults}
            renderItem={renderSearchResultItem}
            keyExtractor={(item) => item.ts_code}
            showsVerticalScrollIndicator={false}
          />
        </View>
      ) : (
        // 显示默认内容
        <FlatList
          data={[]}
          renderItem={() => null}
          ListHeaderComponent={
            <View>
              {renderHotSearches()}
              <Divider style={styles.divider} />
              {renderRecentlyViewed()}
              <Divider style={styles.divider} />
              {renderSearchHistory()}
            </View>
          }
          showsVerticalScrollIndicator={false}
          contentContainerStyle={styles.contentContainer}
        />
      )}
      
      {isSearchLoading && (
        <View style={styles.loadingOverlay}>
          <ActivityIndicator size="small" color={theme.colors.primary} />
        </View>
      )}
    </View>
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

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  searchContainer: {
    padding: 16,
  },
  searchbar: {
    elevation: 2,
  },
  searchInput: {
    fontSize: 16,
  },
  resultsContainer: {
    flex: 1,
  },
  contentContainer: {
    paddingBottom: 16,
  },
  section: {
    padding: 16,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 12,
  },
  chipContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
  },
  chip: {
    marginRight: 8,
    marginBottom: 8,
  },
  chipText: {
    fontSize: 12,
  },
  stockIconContainer: {
    justifyContent: 'center',
    alignItems: 'center',
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: 'rgba(0, 122, 255, 0.1)',
  },
  resultTitle: {
    fontSize: 16,
    fontWeight: '600',
  },
  resultDescription: {
    fontSize: 12,
    opacity: 0.7,
  },
  historyTitle: {
    fontSize: 14,
  },
  historyDescription: {
    fontSize: 12,
    opacity: 0.7,
  },
  recentTitle: {
    fontSize: 14,
    fontWeight: '500',
  },
  recentDescription: {
    fontSize: 12,
    opacity: 0.7,
  },
  divider: {
    marginVertical: 8,
  },
  loadingOverlay: {
    position: 'absolute',
    top: 80,
    right: 32,
  },
});

export default SearchScreen;
