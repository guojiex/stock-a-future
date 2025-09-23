/**
 * 技术指标页面
 */

import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Card, useTheme } from 'react-native-paper';

import { MarketStackParamList } from '@/navigation/AppNavigator';

type TechnicalIndicatorsRouteProp = RouteProp<MarketStackParamList, 'TechnicalIndicators'>;

const TechnicalIndicatorsScreen: React.FC = () => {
  const route = useRoute<TechnicalIndicatorsRouteProp>();
  const theme = useTheme();
  const { stockCode, stockName } = route.params;

  return (
    <View style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            技术指标分析
          </Text>
          <Text style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
            {stockName} ({stockCode})
          </Text>
          <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
            技术指标页面开发中...
          </Text>
          <Text style={[styles.description, { color: theme.colors.onSurfaceVariant }]}>
            这里将显示：
            {'\n'}• 移动平均线 (MA5, MA10, MA20, MA30, MA60)
            {'\n'}• MACD指标
            {'\n'}• RSI指标
            {'\n'}• 布林带
            {'\n'}• KDJ指标
            {'\n'}• 成交量分析
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
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 14,
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

export default TechnicalIndicatorsScreen;
