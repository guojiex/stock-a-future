/**
 * 收藏管理页面
 */

import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Card, useTheme } from 'react-native-paper';

const FavoritesManageScreen: React.FC = () => {
  const theme = useTheme();

  return (
    <View style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            管理收藏
          </Text>
          <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
            收藏管理页面开发中...
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

export default FavoritesManageScreen;
