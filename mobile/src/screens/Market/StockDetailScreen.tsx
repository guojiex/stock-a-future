/**
 * 股票详情页面
 */

import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, ScrollView, RefreshControl } from 'react-native';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Card, useTheme, SegmentedButtons, Chip } from 'react-native-paper';

import { MarketStackParamList } from '@/navigation/AppNavigator';
import KLineChart from '@/components/KLineChart';
import { apiService } from '@/services/apiService';

type StockDetailRouteProp = RouteProp<MarketStackParamList, 'StockDetail'>;

interface StockBasic {
  name: string;
  industry?: string;
  market?: string;
  list_date?: string;
}

interface DailyData {
  trade_date: string;
  open: number | string;
  high: number | string;
  low: number | string;
  close: number | string;
  vol: number | string;
  amount?: number | string;
  pct_chg?: number | string;
}

const StockDetailScreen: React.FC = () => {
  const route = useRoute<StockDetailRouteProp>();
  const theme = useTheme();
  const { stockCode, stockName } = route.params;

  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [stockBasic, setStockBasic] = useState<StockBasic | null>(null);
  const [dailyData, setDailyData] = useState<DailyData[]>([]);
  const [selectedTab, setSelectedTab] = useState('kline');
  const [timeRange, setTimeRange] = useState('90'); // 默认3个月

  // 获取日期范围
  const getDateRange = (days: number) => {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(startDate.getDate() - days);

    const formatDate = (date: Date) => {
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}${month}${day}`;
    };

    return {
      start_date: formatDate(startDate),
      end_date: formatDate(endDate),
    };
  };

  // 加载股票数据
  const loadStockData = async (showLoading = true) => {
    try {
      if (showLoading) {
        setLoading(true);
      }

      // 获取股票基本信息
      const basicResponse = await apiService.getStockBasic(stockCode);
      if (basicResponse.success && basicResponse.data) {
        setStockBasic(basicResponse.data);
      }

      // 获取日线数据
      const { start_date, end_date } = getDateRange(parseInt(timeRange));
      const dailyResponse = await apiService.getDailyData(stockCode, start_date, end_date);
      
      if (dailyResponse.success && dailyResponse.data) {
        setDailyData(dailyResponse.data);
      }

    } catch (error) {
      console.error('加载股票数据失败:', error);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  // 下拉刷新
  const onRefresh = () => {
    setRefreshing(true);
    loadStockData(false);
  };

  // 初始加载
  useEffect(() => {
    loadStockData();
  }, [stockCode, timeRange]);

  // 计算涨跌幅和价格变化
  const getPriceChange = () => {
    if (!dailyData || dailyData.length === 0) return { change: 0, changePercent: 0 };
    
    const latest = dailyData[dailyData.length - 1];
    const previous = dailyData.length > 1 ? dailyData[dailyData.length - 2] : null;
    
    if (!previous) return { change: 0, changePercent: 0 };
    
    const currentClose = parseFloat(String(latest.close));
    const previousClose = parseFloat(String(previous.close));
    const change = currentClose - previousClose;
    const changePercent = (change / previousClose) * 100;
    
    return { change, changePercent };
  };

  const { change, changePercent } = getPriceChange();
  const isPositive = change >= 0;
  const latestPrice = dailyData.length > 0 ? parseFloat(String(dailyData[dailyData.length - 1].close)) : 0;

  return (
    <ScrollView 
      style={styles.container}
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
      }
    >
      {/* 股票头部信息 */}
      <Card style={styles.headerCard}>
        <Card.Content>
          <View style={styles.headerRow}>
            <View style={styles.headerLeft}>
              <Text style={[styles.stockName, { color: theme.colors.onSurface }]}>
                {stockName}
              </Text>
              <Text style={[styles.stockCode, { color: theme.colors.onSurfaceVariant }]}>
                {stockCode}
              </Text>
            </View>
            <View style={styles.headerRight}>
              <Text style={[styles.price, { color: theme.colors.onSurface }]}>
                ¥{latestPrice.toFixed(2)}
              </Text>
              <View style={styles.changeRow}>
                <Chip
                  mode="flat"
                  style={[
                    styles.changeChip,
                    { backgroundColor: isPositive ? '#ef5350' : '#26a69a' }
                  ]}
                  textStyle={styles.changeText}
                >
                  {isPositive ? '+' : ''}{change.toFixed(2)} ({isPositive ? '+' : ''}{changePercent.toFixed(2)}%)
                </Chip>
              </View>
            </View>
          </View>

          {stockBasic && (
            <View style={styles.basicInfo}>
              {stockBasic.industry && (
                <Chip mode="outlined" style={styles.infoChip} compact>
                  {stockBasic.industry}
                </Chip>
              )}
              {stockBasic.market && (
                <Chip mode="outlined" style={styles.infoChip} compact>
                  {stockBasic.market === 'SH' ? '上交所' : '深交所'}
                </Chip>
              )}
            </View>
          )}
        </Card.Content>
      </Card>

      {/* 时间范围选择 */}
      <View style={styles.timeRangeContainer}>
        <SegmentedButtons
          value={timeRange}
          onValueChange={setTimeRange}
          buttons={[
            { value: '30', label: '1月' },
            { value: '90', label: '3月' },
            { value: '180', label: '半年' },
            { value: '365', label: '1年' },
          ]}
          style={styles.segmentedButtons}
        />
      </View>

      {/* Tab导航 */}
      <View style={styles.tabContainer}>
        <SegmentedButtons
          value={selectedTab}
          onValueChange={setSelectedTab}
          buttons={[
            { value: 'kline', label: 'K线', icon: 'chart-line' },
            { value: 'indicators', label: '指标', icon: 'chart-bar' },
            { value: 'predictions', label: '预测', icon: 'target' },
            { value: 'fundamental', label: '基本面', icon: 'cash' },
          ]}
          style={styles.segmentedButtons}
        />
      </View>

      {/* Tab内容 */}
      <View style={styles.contentContainer}>
        {selectedTab === 'kline' && (
          <Card style={styles.chartCard}>
            <Card.Content>
              <Text style={[styles.cardTitle, { color: theme.colors.onSurface }]}>
                K线图表
              </Text>
              <KLineChart 
                data={dailyData} 
                stockCode={stockCode}
                stockName={stockName}
                loading={loading}
              />
            </Card.Content>
          </Card>
        )}

        {selectedTab === 'indicators' && (
          <Card style={styles.card}>
            <Card.Content>
              <Text style={[styles.cardTitle, { color: theme.colors.onSurface }]}>
                技术指标
              </Text>
              <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
                技术指标功能开发中...
              </Text>
            </Card.Content>
          </Card>
        )}

        {selectedTab === 'predictions' && (
          <Card style={styles.card}>
            <Card.Content>
              <Text style={[styles.cardTitle, { color: theme.colors.onSurface }]}>
                买卖预测
              </Text>
              <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
                买卖预测功能开发中...
              </Text>
            </Card.Content>
          </Card>
        )}

        {selectedTab === 'fundamental' && (
          <Card style={styles.card}>
            <Card.Content>
              <Text style={[styles.cardTitle, { color: theme.colors.onSurface }]}>
                基本面数据
              </Text>
              <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
                基本面数据功能开发中...
              </Text>
            </Card.Content>
          </Card>
        )}
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  headerCard: {
    margin: 16,
    marginBottom: 8,
  },
  headerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 12,
  },
  headerLeft: {
    flex: 1,
  },
  headerRight: {
    alignItems: 'flex-end',
  },
  stockName: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  stockCode: {
    fontSize: 14,
  },
  price: {
    fontSize: 28,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  changeRow: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  changeChip: {
    height: 28,
  },
  changeText: {
    fontSize: 12,
    fontWeight: '600',
    color: '#ffffff',
  },
  basicInfo: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
    marginTop: 8,
  },
  infoChip: {
    marginRight: 8,
  },
  timeRangeContainer: {
    paddingHorizontal: 16,
    marginBottom: 8,
  },
  tabContainer: {
    paddingHorizontal: 16,
    marginBottom: 16,
  },
  segmentedButtons: {
    width: '100%',
  },
  contentContainer: {
    paddingHorizontal: 16,
    paddingBottom: 32,
  },
  chartCard: {
    marginBottom: 16,
  },
  card: {
    marginBottom: 16,
    padding: 16,
  },
  cardTitle: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 16,
  },
  placeholder: {
    fontSize: 14,
    textAlign: 'center',
    paddingVertical: 32,
  },
});

export default StockDetailScreen;
