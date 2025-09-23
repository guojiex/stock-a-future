/**
 * 搜索结果页面
 */

import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { RouteProp, useRoute } from '@react-navigation/native';
import { Card, useTheme } from 'react-native-paper';

import { SearchStackParamList } from '@/navigation/AppNavigator';

type SearchResultRouteProp = RouteProp<SearchStackParamList, 'SearchResult'>;

const SearchResultScreen: React.FC = () => {
  const route = useRoute<SearchResultRouteProp>();
  const theme = useTheme();
  const { query } = route.params;

  return (
    <View style={styles.container}>
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            搜索结果
          </Text>
          <Text style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
            关键词: {query}
          </Text>
          <Text style={[styles.placeholder, { color: theme.colors.onSurfaceVariant }]}>
            搜索结果页面开发中...
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
});

export default SearchResultScreen;
