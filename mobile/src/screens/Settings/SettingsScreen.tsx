/**
 * 设置页面
 */

import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Card, useTheme } from 'react-native-paper';

const SettingsScreen: React.FC = () => {
  const theme = useTheme();

  return (
    <View style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            设置
          </Text>
          <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
            设置页面开发中...
          </Text>
          <Text style={[styles.description, { color: theme.colors.onSurfaceVariant }]}>
            这里将显示：
            {'\n'}• 主题设置
            {'\n'}• 数据刷新间隔
            {'\n'}• 通知设置
            {'\n'}• 关于应用
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

export default SettingsScreen;
