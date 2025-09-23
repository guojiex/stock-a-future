/**
 * 市场页面 - 显示股票列表和市场概况
 */

import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  RefreshControl,
  TouchableOpacity,
  Alert,
} from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import {
  Appbar,
  Card,
  List,
  ActivityIndicator,
  Chip,
  Divider,
  useTheme,
} from 'react-native-paper';
import Icon from 'react-native-vector-icons/MaterialIcons';

import { MarketStackParamList } from '@/navigation/AppNavigator';
import { useGetStockListQuery, useGetHealthStatusQuery } from '@/services/api';
import { StockBasic } from '@/types/stock';
import { useAppDispatch, useAppSelector } from '@/hooks/redux';
import { setConnectionStatus } from '@/store/slices/appSlice';

type NavigationProp = NativeStackNavigationProp<MarketStackParamList, 'MarketHome'>;

const MarketScreen: React.FC = () => {
  const navigation = useNavigation<NavigationProp>();
  const theme = useTheme();
  const dispatch = useAppDispatch();
  
  // 本地状态
  const [refreshing, setRefreshing] = useState(false);
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
  } = useGetStockListQuery({
    page: 1,
    pageSize: 100,
  });
  
  // 监听健康状态
  useEffect(() => {
    if (healthData?.success) {
      dispatch(setConnectionStatus('connected'));
    } else if (healthError) {
      dispatch(setConnectionStatus('error'));
    }
  }, [healthData, healthError, dispatch]);
  
  // 处理下拉刷新
  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await Promise.all([
        refetchHealth(),
        refetchStockList(),
      ]);
    } catch (error) {
      console.error('刷新失败:', error);
    } finally {
      setRefreshing(false);
    }
  };
  
  // 筛选股票列表
  const getFilteredStocks = (): StockBasic[] => {
    if (!stockListData?.data?.data) return [];
    
    const stocks = stockListData.data.data;
    
    switch (selectedMarket) {
      case 'sh':
        return stocks.filter(stock => stock.ts_code.endsWith('.SH'));
      case 'sz':
        return stocks.filter(stock => stock.ts_code.endsWith('.SZ'));
      default:
        return stocks;
    }
  };
  
  // 导航到股票详情
  const navigateToStockDetail = (stock: StockBasic) => {
    navigation.navigate('StockDetail', {
      stockCode: stock.ts_code,
      stockName: stock.name,
    });
  };
  
  // 渲染股票列表项
  const renderStockItem = ({ item }: { item: StockBasic }) => (
    <Card style={styles.stockCard} mode="outlined">
      <TouchableOpacity
        onPress={() => navigateToStockDetail(item)}
        activeOpacity={0.7}
      >
        <List.Item
          title={item.name}
          description={item.ts_code}
          left={(props) => (
            <View style={styles.stockIconContainer}>
              <Icon
                name="trending-up"
                size={24}
                color={theme.colors.primary}
                {...props}
              />
            </View>
          )}
          right={(props) => (
            <View style={styles.stockInfoContainer}>
              <Text style={[styles.stockPrice, { color: theme.colors.onSurface }]}>
                --
              </Text>
              <Text style={[styles.stockChange, { color: theme.colors.onSurfaceVariant }]}>
                --
              </Text>
            </View>
          )}
          titleStyle={styles.stockTitle}
          descriptionStyle={styles.stockCode}
        />
      </TouchableOpacity>
    </Card>
  );
  
  // 渲染市场筛选器
  const renderMarketFilter = () => (
    <View style={styles.filterContainer}>
      <Chip
        selected={selectedMarket === 'all'}
        onPress={() => setSelectedMarket('all')}
        style={styles.filterChip}
      >
        全部
      </Chip>
      <Chip
        selected={selectedMarket === 'sh'}
        onPress={() => setSelectedMarket('sh')}
        style={styles.filterChip}
      >
        上海
      </Chip>
      <Chip
        selected={selectedMarket === 'sz'}
        onPress={() => setSelectedMarket('sz')}
        style={styles.filterChip}
      >
        深圳
      </Chip>
    </View>
  );
  
  // 渲染连接状态
  const renderConnectionStatus = () => {
    const getStatusColor = () => {
      switch (connectionStatus) {
        case 'connected':
          return theme.colors.success;
        case 'connecting':
          return theme.colors.warning;
        case 'error':
          return theme.colors.error;
        default:
          return theme.colors.onSurfaceVariant;
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
      <View style={styles.statusContainer}>
        <View style={[styles.statusDot, { backgroundColor: getStatusColor() }]} />
        <Text style={[styles.statusText, { color: getStatusColor() }]}>
          {getStatusText()}
        </Text>
      </View>
    );
  };
  
  // 渲染头部
  const renderHeader = () => (
    <View>
      <Card style={styles.headerCard} mode="contained">
        <Card.Content>
          <View style={styles.headerContent}>
            <View>
              <Text style={styles.headerTitle}>Stock-A-Future</Text>
              <Text style={styles.headerSubtitle}>A股股票分析系统</Text>
            </View>
            {renderConnectionStatus()}
          </View>
        </Card.Content>
      </Card>
      
      {renderMarketFilter()}
      <Divider style={styles.divider} />
    </View>
  );
  
  // 渲染加载状态
  if (isStockListLoading && !stockListData) {
    return (
      <View style={styles.container}>
        <Appbar.Header>
          <Appbar.Content title="市场" />
        </Appbar.Header>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={theme.colors.primary} />
          <Text style={styles.loadingText}>正在加载股票列表...</Text>
        </View>
      </View>
    );
  }
  
  // 渲染错误状态
  if (stockListError && !stockListData) {
    return (
      <View style={styles.container}>
        <Appbar.Header>
          <Appbar.Content title="市场" />
        </Appbar.Header>
        <View style={styles.errorContainer}>
          <Icon name="error-outline" size={48} color={theme.colors.error} />
          <Text style={styles.errorText}>加载失败</Text>
          <Text style={styles.errorSubtext}>
            {stockListError.toString()}
          </Text>
          <TouchableOpacity
            style={[styles.retryButton, { backgroundColor: theme.colors.primary }]}
            onPress={() => refetchStockList()}
          >
            <Text style={[styles.retryButtonText, { color: theme.colors.onPrimary }]}>
              重试
            </Text>
          </TouchableOpacity>
        </View>
      </View>
    );
  }
  
  const filteredStocks = getFilteredStocks();
  
  return (
    <View style={styles.container}>
      <Appbar.Header>
        <Appbar.Content title="市场" />
        <Appbar.Action
          icon="refresh"
          onPress={handleRefresh}
          disabled={refreshing}
        />
      </Appbar.Header>
      
      <FlatList
        data={filteredStocks}
        renderItem={renderStockItem}
        keyExtractor={(item) => item.ts_code}
        ListHeaderComponent={renderHeader}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={handleRefresh}
            colors={[theme.colors.primary]}
          />
        }
        contentContainerStyle={styles.listContainer}
        showsVerticalScrollIndicator={false}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  listContainer: {
    paddingBottom: 16,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  loadingText: {
    marginTop: 16,
    fontSize: 16,
  },
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  errorText: {
    marginTop: 16,
    fontSize: 18,
    fontWeight: 'bold',
  },
  errorSubtext: {
    marginTop: 8,
    fontSize: 14,
    textAlign: 'center',
    opacity: 0.7,
  },
  retryButton: {
    marginTop: 24,
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 8,
  },
  retryButtonText: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  headerCard: {
    margin: 16,
    marginBottom: 8,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: 'bold',
  },
  headerSubtitle: {
    fontSize: 14,
    opacity: 0.7,
    marginTop: 4,
  },
  statusContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  statusDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginRight: 8,
  },
  statusText: {
    fontSize: 12,
    fontWeight: '500',
  },
  filterContainer: {
    flexDirection: 'row',
    paddingHorizontal: 16,
    paddingVertical: 8,
  },
  filterChip: {
    marginRight: 8,
  },
  divider: {
    marginVertical: 8,
  },
  stockCard: {
    marginHorizontal: 16,
    marginVertical: 4,
  },
  stockIconContainer: {
    justifyContent: 'center',
    alignItems: 'center',
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: 'rgba(0, 122, 255, 0.1)',
  },
  stockInfoContainer: {
    alignItems: 'flex-end',
    justifyContent: 'center',
  },
  stockTitle: {
    fontSize: 16,
    fontWeight: '600',
  },
  stockCode: {
    fontSize: 12,
    opacity: 0.7,
  },
  stockPrice: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  stockChange: {
    fontSize: 12,
    marginTop: 2,
  },
});

export default MarketScreen;
