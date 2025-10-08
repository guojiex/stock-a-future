/**
 * 分组管理页面
 */

import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Card, useTheme } from 'react-native-paper';

import { FavoritesStackParamList } from '@/navigation/AppNavigator';

type GroupManageRouteProp = RouteProp<FavoritesStackParamList, 'GroupManage'>;

const GroupManageScreen: React.FC = () => {
  const route = useRoute<GroupManageRouteProp>();
  const theme = useTheme();
  const { groupId } = route.params;

  return (
    <View style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            {groupId ? '编辑分组' : '新建分组'}
          </Text>
          <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
            分组管理页面开发中...
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

export default GroupManageScreen;
