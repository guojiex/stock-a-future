/**
 * 收藏页面
 */

import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Card, useTheme } from 'react-native-paper';

const FavoritesScreen: React.FC = () => {
  const theme = useTheme();

  return (
    <View style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            我的收藏
          </Text>
          <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
            收藏页面开发中...
          </Text>
          <Text style={[styles.description, { color: theme.colors.onSurfaceVariant }]}>
            这里将显示：
            {'\n'}• 收藏的股票列表
            {'\n'}• 分组管理
            {'\n'}• 快速访问
            {'\n'}• 价格提醒
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
  description: {
    fontSize: 14,
    lineHeight: 20,
  },
});

export default FavoritesScreen;
