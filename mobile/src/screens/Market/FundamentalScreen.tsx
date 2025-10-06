/**
 * 基本面分析页面
 */

import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Card, useTheme } from 'react-native-paper';

import { MarketStackParamList } from '@/navigation/AppNavigator';

type FundamentalRouteProp = RouteProp<MarketStackParamList, 'Fundamental'>;

const FundamentalScreen: React.FC = () => {
  const route = useRoute<FundamentalRouteProp>();
  const theme = useTheme();
  const { stockCode, stockName } = route.params;

  return (
    <View style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            基本面分析
          </Text>
          <Text style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
            {stockName} ({stockCode})
          </Text>
          <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
            基本面分析页面开发中...
          </Text>
          <Text style={[styles.description, { color: theme.colors.onSurfaceVariant }]}>
            这里将显示：
            {'\n'}• 利润表数据
            {'\n'}• 资产负债表
            {'\n'}• 现金流量表
            {'\n'}• 每日基本面指标
            {'\n'}• 财务比率分析
            {'\n'}• 估值指标
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

export default FundamentalScreen;
