/**
 * 策略Tab组件
 * 在股票详情页中显示该股票的收藏策略信息和信号
 */

import React, { useState, useMemo, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  RefreshControl,
  Alert,
} from 'react-native';
import {
  Card,
  Button,
  Chip,
  useTheme,
  ActivityIndicator,
  Divider,
  IconButton,
} from 'react-native-paper';
import Icon from 'react-native-vector-icons/MaterialIcons';

import {
  useGetFavoritesQuery,
  useGetGroupsQuery,
  useAddFavoriteMutation,
  useDeleteFavoriteMutation,
  useUpdateFavoriteMutation,
  useGetPredictionsQuery,
  useGetPatternSummaryQuery,
} from '@/services/api';

interface StrategyTabProps {
  stockCode: string;
  stockName: string;
}

interface FavoriteItem {
  id: string;
  ts_code: string;
  name: string;
  start_date: string;
  end_date: string;
  group_id: string;
  group_name?: string;
  created_at: string;
}

interface Group {
  id: string;
  name: string;
  color: string;
  order_num: number;
}

const StrategyTab: React.FC<StrategyTabProps> = ({ stockCode, stockName }) => {
  const theme = useTheme();
  
  // API查询
  const { 
    data: favoritesData, 
    isLoading: favoritesLoading,
    refetch: refetchFavorites 
  } = useGetFavoritesQuery();
  
  const { 
    data: groupsData, 
    isLoading: groupsLoading 
  } = useGetGroupsQuery();
  
  const {
    data: predictionsData,
    isLoading: predictionsLoading,
    refetch: refetchPredictions
  } = useGetPredictionsQuery(stockCode, { skip: !stockCode });
  
  const {
    data: patternSummaryData,
    isLoading: patternSummaryLoading,
    refetch: refetchPatternSummary
  } = useGetPatternSummaryQuery(stockCode, { skip: !stockCode });

  // API变更
  const [addFavorite, { isLoading: addLoading }] = useAddFavoriteMutation();
  const [deleteFavorite, { isLoading: deleteLoading }] = useDeleteFavoriteMutation();
  const [updateFavorite, { isLoading: updateLoading }] = useUpdateFavoriteMutation();

  // 本地状态
  const [refreshing, setRefreshing] = useState(false);
  const [selectedTab, setSelectedTab] = useState<'info' | 'signals'>('info');

  // 数据处理
  const favorites = favoritesData?.data?.favorites || [];
  const groups = groupsData?.data?.groups || [];
  const predictions = predictionsData?.data || [];
  const patternSummary = patternSummaryData?.data || {};

  // 检查当前股票是否已收藏
  const currentFavorite = useMemo(() => {
    return favorites.find((fav: FavoriteItem) => fav.ts_code === stockCode);
  }, [favorites, stockCode]);

  const isFavorited = !!currentFavorite;

  // 获取分组信息
  const getGroupInfo = (groupId: string) => {
    return groups.find((g: Group) => g.id === groupId);
  };

  // 下拉刷新
  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await Promise.all([
        refetchFavorites(),
        refetchPredictions(),
        refetchPatternSummary(),
      ]);
    } finally {
      setRefreshing(false);
    }
  };

  // 添加到收藏
  const handleAddFavorite = async () => {
    try {
      const defaultGroup = groups.find((g: Group) => g.name === '默认分组') || groups[0];
      
      // 默认使用最近一年的数据
      const endDate = new Date();
      const startDate = new Date();
      startDate.setFullYear(startDate.getFullYear() - 1);
      
      const formatDate = (date: Date) => {
        return date.toISOString().split('T')[0];
      };

      await addFavorite({
        ts_code: stockCode,
        name: stockName,
        start_date: formatDate(startDate),
        end_date: formatDate(endDate),
        group_id: defaultGroup?.id || 'default',
      }).unwrap();

      Alert.alert('成功', '已添加到收藏策略');
      refetchFavorites();
    } catch (error: any) {
      Alert.alert('错误', error?.data?.error || '添加收藏失败');
    }
  };

  // 从收藏中删除
  const handleRemoveFavorite = async () => {
    if (!currentFavorite) return;

    Alert.alert(
      '确认删除',
      `确定要将 ${stockName} 从收藏中移除吗？`,
      [
        { text: '取消', style: 'cancel' },
        {
          text: '删除',
          style: 'destructive',
          onPress: async () => {
            try {
              await deleteFavorite(currentFavorite.id).unwrap();
              Alert.alert('成功', '已从收藏中移除');
              refetchFavorites();
            } catch (error: any) {
              Alert.alert('错误', error?.data?.error || '删除失败');
            }
          },
        },
      ]
    );
  };

  // 渲染收藏信息
  const renderFavoriteInfo = () => {
    if (!isFavorited) {
      return (
        <Card style={styles.emptyCard}>
          <Card.Content style={styles.emptyContent}>
            <Icon name="bookmark-border" size={64} color={theme.colors.onSurfaceVariant} />
            <Text style={[styles.emptyTitle, { color: theme.colors.onSurface }]}>
              未收藏此股票
            </Text>
            <Text style={[styles.emptyText, { color: theme.colors.onSurfaceVariant }]}>
              添加到收藏后，可以跟踪该股票的策略信号
            </Text>
            <Button
              mode="contained"
              onPress={handleAddFavorite}
              loading={addLoading}
              disabled={addLoading}
              style={styles.addButton}
              icon="bookmark"
            >
              添加到收藏
            </Button>
          </Card.Content>
        </Card>
      );
    }

    const groupInfo = getGroupInfo(currentFavorite.group_id);

    return (
      <View style={styles.infoContainer}>
        <Card style={styles.card}>
          <Card.Content>
            <View style={styles.cardHeader}>
              <Text style={[styles.cardTitle, { color: theme.colors.onSurface }]}>
                策略信息
              </Text>
              <IconButton
                icon="delete"
                size={20}
                iconColor={theme.colors.error}
                onPress={handleRemoveFavorite}
              />
            </View>
            
            <View style={styles.infoRow}>
              <Text style={[styles.infoLabel, { color: theme.colors.onSurfaceVariant }]}>
                股票代码
              </Text>
              <Text style={[styles.infoValue, { color: theme.colors.onSurface }]}>
                {currentFavorite.ts_code}
              </Text>
            </View>

            <View style={styles.infoRow}>
              <Text style={[styles.infoLabel, { color: theme.colors.onSurfaceVariant }]}>
                股票名称
              </Text>
              <Text style={[styles.infoValue, { color: theme.colors.onSurface }]}>
                {currentFavorite.name}
              </Text>
            </View>

            <View style={styles.infoRow}>
              <Text style={[styles.infoLabel, { color: theme.colors.onSurfaceVariant }]}>
                分析周期
              </Text>
              <Text style={[styles.infoValue, { color: theme.colors.onSurface }]}>
                {currentFavorite.start_date} 至 {currentFavorite.end_date}
              </Text>
            </View>

            <View style={styles.infoRow}>
              <Text style={[styles.infoLabel, { color: theme.colors.onSurfaceVariant }]}>
                所属分组
              </Text>
              <Chip
                mode="flat"
                style={[styles.groupChip, { backgroundColor: groupInfo?.color || '#1976d2' }]}
                textStyle={styles.groupChipText}
              >
                {groupInfo?.name || '默认分组'}
              </Chip>
            </View>

            <View style={styles.infoRow}>
              <Text style={[styles.infoLabel, { color: theme.colors.onSurfaceVariant }]}>
                添加时间
              </Text>
              <Text style={[styles.infoValue, { color: theme.colors.onSurface }]}>
                {new Date(currentFavorite.created_at).toLocaleDateString('zh-CN')}
              </Text>
            </View>
          </Card.Content>
        </Card>
      </View>
    );
  };

  // 渲染信号汇总
  const renderSignals = () => {
    if (predictionsLoading || patternSummaryLoading) {
      return (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" />
          <Text style={{ color: theme.colors.onSurfaceVariant, marginTop: 16 }}>
            加载信号数据...
          </Text>
        </View>
      );
    }

    const hasSignals = predictions.length > 0 || Object.keys(patternSummary).length > 0;

    if (!hasSignals) {
      return (
        <Card style={styles.emptyCard}>
          <Card.Content style={styles.emptyContent}>
            <Icon name="show-chart" size={64} color={theme.colors.onSurfaceVariant} />
            <Text style={[styles.emptyTitle, { color: theme.colors.onSurface }]}>
              暂无信号数据
            </Text>
            <Text style={[styles.emptyText, { color: theme.colors.onSurfaceVariant }]}>
              请稍后刷新查看最新的预测信号
            </Text>
          </Card.Content>
        </Card>
      );
    }

    return (
      <View style={styles.signalsContainer}>
        {/* 预测信号 */}
        {predictions.length > 0 && (
          <Card style={styles.card}>
            <Card.Content>
              <Text style={[styles.cardTitle, { color: theme.colors.onSurface }]}>
                预测信号
              </Text>
              {predictions.map((pred: any, index: number) => (
                <View key={index} style={styles.signalItem}>
                  <View style={styles.signalHeader}>
                    <Chip
                      mode="flat"
                      style={[
                        styles.signalChip,
                        {
                          backgroundColor:
                            pred.signal_type === 'BUY'
                              ? '#ef5350'
                              : pred.signal_type === 'SELL'
                              ? '#26a69a'
                              : '#9e9e9e',
                        },
                      ]}
                      textStyle={styles.signalChipText}
                    >
                      {pred.signal_type === 'BUY' ? '买入' : pred.signal_type === 'SELL' ? '卖出' : '持有'}
                    </Chip>
                    <Text style={[styles.confidence, { color: theme.colors.onSurfaceVariant }]}>
                      置信度: {(pred.confidence * 100).toFixed(1)}%
                    </Text>
                  </View>
                  <Text style={[styles.signalReason, { color: theme.colors.onSurface }]}>
                    {pred.reason || '无详细理由'}
                  </Text>
                </View>
              ))}
            </Card.Content>
          </Card>
        )}

        {/* 形态识别 */}
        {Object.keys(patternSummary).length > 0 && (
          <Card style={styles.card}>
            <Card.Content>
              <Text style={[styles.cardTitle, { color: theme.colors.onSurface }]}>
                形态识别
              </Text>
              {Object.entries(patternSummary).map(([key, value]: [string, any]) => (
                <View key={key} style={styles.patternItem}>
                  <Text style={[styles.patternName, { color: theme.colors.onSurface }]}>
                    {key}
                  </Text>
                  <Text style={[styles.patternValue, { color: theme.colors.primary }]}>
                    {typeof value === 'number' ? value.toFixed(2) : value}
                  </Text>
                </View>
              ))}
            </Card.Content>
          </Card>
        )}
      </View>
    );
  };

  const isLoading = favoritesLoading || groupsLoading;

  return (
    <ScrollView
      style={styles.container}
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />
      }
    >
      {/* Tab切换 */}
      <View style={styles.tabBar}>
        <TouchableOpacity
          style={[
            styles.tab,
            selectedTab === 'info' && styles.tabActive,
            { borderBottomColor: theme.colors.primary },
          ]}
          onPress={() => setSelectedTab('info')}
        >
          <Icon
            name="info"
            size={20}
            color={selectedTab === 'info' ? theme.colors.primary : theme.colors.onSurfaceVariant}
          />
          <Text
            style={[
              styles.tabText,
              {
                color:
                  selectedTab === 'info'
                    ? theme.colors.primary
                    : theme.colors.onSurfaceVariant,
              },
            ]}
          >
            策略信息
          </Text>
        </TouchableOpacity>
        
        <TouchableOpacity
          style={[
            styles.tab,
            selectedTab === 'signals' && styles.tabActive,
            { borderBottomColor: theme.colors.primary },
          ]}
          onPress={() => setSelectedTab('signals')}
        >
          <Icon
            name="show-chart"
            size={20}
            color={selectedTab === 'signals' ? theme.colors.primary : theme.colors.onSurfaceVariant}
          />
          <Text
            style={[
              styles.tabText,
              {
                color:
                  selectedTab === 'signals'
                    ? theme.colors.primary
                    : theme.colors.onSurfaceVariant,
              },
            ]}
          >
            信号汇总
          </Text>
        </TouchableOpacity>
      </View>

      {/* 内容区域 */}
      <View style={styles.content}>
        {isLoading ? (
          <View style={styles.loadingContainer}>
            <ActivityIndicator size="large" />
          </View>
        ) : selectedTab === 'info' ? (
          renderFavoriteInfo()
        ) : (
          renderSignals()
        )}
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  tabBar: {
    flexDirection: 'row',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
    marginBottom: 16,
  },
  tab: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 12,
    gap: 8,
    borderBottomWidth: 2,
    borderBottomColor: 'transparent',
  },
  tabActive: {
    borderBottomWidth: 2,
  },
  tabText: {
    fontSize: 14,
    fontWeight: '500',
  },
  content: {
    flex: 1,
    paddingBottom: 24,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: 48,
  },
  emptyCard: {
    marginHorizontal: 16,
  },
  emptyContent: {
    alignItems: 'center',
    paddingVertical: 32,
  },
  emptyTitle: {
    fontSize: 18,
    fontWeight: '600',
    marginTop: 16,
    marginBottom: 8,
  },
  emptyText: {
    fontSize: 14,
    textAlign: 'center',
    marginBottom: 24,
  },
  addButton: {
    marginTop: 8,
  },
  infoContainer: {
    paddingHorizontal: 16,
  },
  card: {
    marginBottom: 16,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  cardTitle: {
    fontSize: 16,
    fontWeight: '600',
  },
  infoRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  infoLabel: {
    fontSize: 14,
  },
  infoValue: {
    fontSize: 14,
    fontWeight: '500',
  },
  groupChip: {
    height: 28,
  },
  groupChipText: {
    fontSize: 12,
    color: '#fff',
  },
  signalsContainer: {
    paddingHorizontal: 16,
  },
  signalItem: {
    marginBottom: 16,
    paddingBottom: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  signalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  signalChip: {
    height: 28,
  },
  signalChipText: {
    fontSize: 12,
    fontWeight: '600',
    color: '#fff',
  },
  confidence: {
    fontSize: 12,
  },
  signalReason: {
    fontSize: 14,
    lineHeight: 20,
  },
  patternItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  patternName: {
    fontSize: 14,
  },
  patternValue: {
    fontSize: 14,
    fontWeight: '600',
  },
});

export default StrategyTab;

