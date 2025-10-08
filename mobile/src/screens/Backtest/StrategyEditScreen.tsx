/**
 * 策略编辑页面
 */

import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Card, useTheme } from 'react-native-paper';

import { BacktestStackParamList } from '@/navigation/AppNavigator';

type StrategyEditRouteProp = RouteProp<BacktestStackParamList, 'StrategyEdit'>;

const StrategyEditScreen: React.FC = () => {
  const route = useRoute<StrategyEditRouteProp>();
  const theme = useTheme();
  const { strategyId } = route.params;

  return (
    <View style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            {strategyId ? '编辑策略' : '新建策略'}
          </Text>
          <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
            策略编辑页面开发中...
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
    marginBottom: 24,
  },
  placeholder: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 16,
  },
});

export default StrategyEditScreen;
