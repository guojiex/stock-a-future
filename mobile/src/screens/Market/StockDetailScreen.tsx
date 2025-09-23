/**
 * 股票详情页面
 */

import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Appbar, Card, useTheme } from 'react-native-paper';

import { MarketStackParamList } from '@/navigation/AppNavigator';

type StockDetailRouteProp = RouteProp<MarketStackParamList, 'StockDetail'>;

const StockDetailScreen: React.FC = () => {
  const route = useRoute<StockDetailRouteProp>();
  const theme = useTheme();
  const { stockCode, stockName } = route.params;

  return (
    <View style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            {stockName}
          </Text>
          <Text style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
            {stockCode}
          </Text>
          <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
            股票详情页面开发中...
          </Text>
          <Text style={[styles.description, { color: theme.colors.onSurfaceVariant }]}>
            这里将显示：
            {'\n'}• K线图表
            {'\n'}• 实时价格
            {'\n'}• 技术指标
            {'\n'}• 基本面数据
            {'\n'}• 买卖信号
          </Text>
        </Card.Content>
      </Card>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
  },
  card: {
    padding: 16,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 16,
    marginBottom: 24,
  },
  placeholder: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 16,
  },
  description: {
    fontSize: 14,
    lineHeight: 20,
  },
});

export default StockDetailScreen;
