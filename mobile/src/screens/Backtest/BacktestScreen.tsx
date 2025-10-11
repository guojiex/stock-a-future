/**
 * 策略管理页面
 * 功能：
 * - 显示收藏的股票列表（策略）
 * - 分组管理
 * - 信号汇总展示
 * - 快速访问和编辑
 */

import React, { useState, useMemo } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  RefreshControl,
  Alert,
  FlatList,
} from 'react-native';
import {
  Card,
  Button,
  IconButton,
  Chip,
  useTheme,
  FAB,
  Portal,
  Dialog,
  TextInput,
  ActivityIndicator,
  Divider,
  Menu,
} from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import Icon from 'react-native-vector-icons/MaterialIcons';

import {
  useGetFavoritesQuery,
  useGetGroupsQuery,
  useDeleteFavoriteMutation,
  useUpdateFavoriteMutation,
  useCreateGroupMutation,
  useUpdateGroupMutation,
  useDeleteGroupMutation,
  useGetFavoritesSignalsQuery,
} from '@/services/api';

interface FavoriteItem {
  id: string;
  ts_code: string;
  name: string;
  start_date?: string;
  end_date?: string;
  group_id?: string;
  sort_order?: number;
  created_at: string;
  updated_at: string;
}

interface Group {
  id: string;
  name: string;
  color?: string;
  sort_order?: number;
  created_at: string;
  updated_at: string;
}

const BacktestScreen: React.FC = () => {
  const theme = useTheme();
  const navigation = useNavigation();

  // API查询
  const { 
    data: favoritesData, 
    isLoading: favoritesLoading,
    refetch: refetchFavorites 
  } = useGetFavoritesQuery();
  const { 
    data: groupsData, 
    isLoading: groupsLoading,
    refetch: refetchGroups 
  } = useGetGroupsQuery();
  const {
    data: signalsData,
    isLoading: signalsLoading,
    refetch: refetchSignals
  } = useGetFavoritesSignalsQuery();

  // API变更
  const [deleteFavorite] = useDeleteFavoriteMutation();
  const [updateFavorite] = useUpdateFavoriteMutation();
  const [createGroup] = useCreateGroupMutation();
  const [updateGroup] = useUpdateGroupMutation();
  const [deleteGroup] = useDeleteGroupMutation();

  // 本地状态
  const [selectedTab, setSelectedTab] = useState<'list' | 'signals'>('list');
  const [selectedGroup, setSelectedGroup] = useState<string>('all');
  const [fabOpen, setFabOpen] = useState(false);
  const [groupDialogVisible, setGroupDialogVisible] = useState(false);
  const [editingGroup, setEditingGroup] = useState<Group | null>(null);
  const [groupName, setGroupName] = useState('');
  const [groupColor, setGroupColor] = useState('#1976d2');
  const [refreshing, setRefreshing] = useState(false);
  const [groupMenuVisible, setGroupMenuVisible] = useState<string | null>(null);

  // 获取数据
  const favorites = (favoritesData?.data?.favorites || []) as FavoriteItem[];
  const groups = (groupsData?.data?.groups || []) as Group[];
  const signals = signalsData?.data?.signals || [];

  // 按分组过滤收藏
  const filteredFavorites = useMemo(() => {
    if (selectedGroup === 'all') {
      return favorites;
    }
    return favorites.filter(fav => (fav.group_id || 'default') === selectedGroup);
  }, [favorites, selectedGroup]);

  // 按分组统计数量
  const groupCounts = useMemo(() => {
    const counts: Record<string, number> = {
      all: favorites.length,
    };
    groups.forEach(group => {
      counts[group.id] = favorites.filter(
        fav => (fav.group_id || 'default') === group.id
      ).length;
    });
    return counts;
  }, [favorites, groups]);

  // 处理刷新
  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await Promise.all([
        refetchFavorites(),
        refetchGroups(),
        selectedTab === 'signals' ? refetchSignals() : Promise.resolve(),
      ]);
    } finally {
      setRefreshing(false);
    }
  };

  // 处理删除收藏
  const handleDeleteFavorite = async (favoriteId: string, name: string) => {
    Alert.alert(
      '删除策略',
      `确定要删除策略 "${name}" 吗？`,
      [
        { text: '取消', style: 'cancel' },
        {
          text: '删除',
          style: 'destructive',
          onPress: async () => {
            try {
              await deleteFavorite(Number(favoriteId)).unwrap();
            } catch (error) {
              Alert.alert('错误', '删除失败，请重试');
            }
          },
        },
      ]
    );
  };

  // 打开分组对话框
  const handleOpenGroupDialog = (group?: Group) => {
    if (group) {
      setEditingGroup(group);
      setGroupName(group.name);
      setGroupColor(group.color || '#1976d2');
    } else {
      setEditingGroup(null);
      setGroupName('');
      setGroupColor('#1976d2');
    }
    setGroupDialogVisible(true);
  };

  // 保存分组
  const handleSaveGroup = async () => {
    if (!groupName.trim()) {
      Alert.alert('提示', '请输入分组名称');
      return;
    }

    try {
      if (editingGroup) {
        await updateGroup({
          id: Number(editingGroup.id),
          data: {
            name: groupName,
            color: groupColor,
          },
        }).unwrap();
      } else {
        await createGroup({
          name: groupName,
          color: groupColor,
        }).unwrap();
      }
      setGroupDialogVisible(false);
    } catch (error) {
      Alert.alert('错误', '保存失败，请重试');
    }
  };

  // 删除分组
  const handleDeleteGroup = async (groupId: string, groupName: string) => {
    Alert.alert(
      '删除分组',
      `确定要删除分组 "${groupName}" 吗？\n该分组下的策略将移至默认分组。`,
      [
        { text: '取消', style: 'cancel' },
        {
          text: '删除',
          style: 'destructive',
          onPress: async () => {
            try {
              await deleteGroup(Number(groupId)).unwrap();
              if (selectedGroup === groupId) {
                setSelectedGroup('all');
              }
            } catch (error) {
              Alert.alert('错误', '删除失败，请重试');
            }
          },
        },
      ]
    );
  };

  // 渲染收藏项
  const renderFavoriteItem = ({ item }: { item: FavoriteItem }) => {
    const group = groups.find(g => g.id === item.group_id);
    const createdDate = new Date(item.created_at).toLocaleDateString('zh-CN');

    return (
      <Card style={styles.favoriteCard}>
        <TouchableOpacity
          onPress={() => {
            // TODO: 导航到股票详情页
            console.log('Navigate to stock detail:', item.ts_code);
          }}
        >
          <Card.Content>
            <View style={styles.favoriteHeader}>
              <View style={styles.favoriteInfo}>
                <Text style={[styles.favoriteName, { color: theme.colors.onSurface }]}>
                  {item.name}
                </Text>
                <Text style={[styles.favoriteCode, { color: theme.colors.onSurfaceVariant }]}>
                  {item.ts_code}
                </Text>
              </View>
              <View style={styles.favoriteActions}>
                <IconButton
                  icon="chart-line"
                  size={20}
                  onPress={() => {
                    // TODO: 查看图表
                    console.log('View chart:', item.ts_code);
                  }}
                />
                <IconButton
                  icon="delete"
                  size={20}
                  iconColor={theme.colors.error}
                  onPress={() => handleDeleteFavorite(item.id, item.name)}
                />
              </View>
            </View>
            <View style={styles.favoriteDetails}>
              {group && (
                <Chip
                  mode="flat"
                  style={[styles.groupChip, { backgroundColor: group.color }]}
                  textStyle={{ color: '#fff', fontSize: 12 }}
                >
                  {group.name}
                </Chip>
              )}
              <Text style={[styles.favoriteDate, { color: theme.colors.onSurfaceVariant }]}>
                收藏于 {createdDate}
              </Text>
            </View>
          </Card.Content>
        </TouchableOpacity>
      </Card>
    );
  };

  // 渲染信号项
  const renderSignalItem = ({ item }: { item: any }) => {
    const signalType = item.signal_type || 'HOLD';
    const signalColor = 
      signalType === 'BUY' ? '#4caf50' : 
      signalType === 'SELL' ? '#f44336' : '#ff9800';
    
    const signalText = 
      signalType === 'BUY' ? '买入' : 
      signalType === 'SELL' ? '卖出' : '持有';

    return (
      <Card style={styles.signalCard}>
        <Card.Content>
          <View style={styles.signalHeader}>
            <View style={styles.signalInfo}>
              <Text style={[styles.signalName, { color: theme.colors.onSurface }]}>
                {item.name}
              </Text>
              <Text style={[styles.signalCode, { color: theme.colors.onSurfaceVariant }]}>
                {item.ts_code}
              </Text>
            </View>
            <Chip
              mode="flat"
              style={[styles.signalChip, { backgroundColor: signalColor }]}
              textStyle={{ color: '#fff', fontSize: 12 }}
            >
              {signalText}
            </Chip>
          </View>
          {item.current_price && (
            <View style={styles.signalDetails}>
              <Text style={[styles.signalPrice, { color: theme.colors.onSurface }]}>
                ¥{item.current_price}
              </Text>
              {item.trade_date && (
                <Text style={[styles.signalDate, { color: theme.colors.onSurfaceVariant }]}>
                  {item.trade_date}
                </Text>
              )}
            </View>
          )}
          {item.reason && (
            <Text style={[styles.signalReason, { color: theme.colors.onSurfaceVariant }]}>
              {item.reason}
            </Text>
          )}
        </Card.Content>
      </Card>
    );
  };

  // 渲染分组标签
  const renderGroupTabs = () => (
    <ScrollView
      horizontal
      showsHorizontalScrollIndicator={false}
      style={styles.groupTabs}
      contentContainerStyle={styles.groupTabsContent}
    >
      <Chip
        mode={selectedGroup === 'all' ? 'flat' : 'outlined'}
        selected={selectedGroup === 'all'}
        onPress={() => setSelectedGroup('all')}
        style={styles.groupTab}
      >
        全部 ({groupCounts.all || 0})
      </Chip>
      {groups.map(group => (
        <View key={group.id} style={styles.groupTabContainer}>
          <Chip
            mode={selectedGroup === group.id ? 'flat' : 'outlined'}
            selected={selectedGroup === group.id}
            onPress={() => setSelectedGroup(group.id)}
            style={[
              styles.groupTab,
              selectedGroup === group.id && { backgroundColor: group.color },
            ]}
          >
            {group.name} ({groupCounts[group.id] || 0})
          </Chip>
          {group.id !== 'default' && (
            <Menu
              visible={groupMenuVisible === group.id}
              onDismiss={() => setGroupMenuVisible(null)}
              anchor={
                <IconButton
                  icon="dots-vertical"
                  size={16}
                  onPress={() => setGroupMenuVisible(group.id)}
                  style={styles.groupMenuButton}
                />
              }
            >
              <Menu.Item
                onPress={() => {
                  setGroupMenuVisible(null);
                  handleOpenGroupDialog(group);
                }}
                title="编辑"
                leadingIcon="pencil"
              />
              <Menu.Item
                onPress={() => {
                  setGroupMenuVisible(null);
                  handleDeleteGroup(group.id, group.name);
                }}
                title="删除"
                leadingIcon="delete"
              />
            </Menu>
          )}
        </View>
      ))}
      <Button
        mode="outlined"
        onPress={() => handleOpenGroupDialog()}
        style={styles.addGroupButton}
        compact
      >
        + 新建分组
      </Button>
    </ScrollView>
  );

  // 渲染空状态
  const renderEmptyState = () => (
    <View style={styles.emptyState}>
      <Icon
        name={selectedTab === 'list' ? 'star-border' : 'show-chart'}
        size={64}
        color={theme.colors.onSurfaceVariant}
      />
      <Text style={[styles.emptyTitle, { color: theme.colors.onSurface }]}>
        {selectedTab === 'list' ? '还没有收藏任何策略' : '暂无信号数据'}
      </Text>
      <Text style={[styles.emptyDescription, { color: theme.colors.onSurfaceVariant }]}>
        {selectedTab === 'list'
          ? '点击右下角的 + 按钮添加策略'
          : '收藏股票后，点击刷新按钮获取信号'}
      </Text>
    </View>
  );

  if (favoritesLoading || groupsLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      {/* 顶部Tab切换 */}
      <View style={styles.tabBar}>
        <TouchableOpacity
          style={[
            styles.tab,
            selectedTab === 'list' && styles.tabActive,
            { borderBottomColor: theme.colors.primary },
          ]}
          onPress={() => setSelectedTab('list')}
        >
          <Icon name="list" size={20} color={selectedTab === 'list' ? theme.colors.primary : theme.colors.onSurfaceVariant} />
          <Text
            style={[
              styles.tabText,
              {
                color: selectedTab === 'list'
                  ? theme.colors.primary
                  : theme.colors.onSurfaceVariant,
              },
            ]}
          >
            策略列表
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
          <Icon name="show-chart" size={20} color={selectedTab === 'signals' ? theme.colors.primary : theme.colors.onSurfaceVariant} />
          <Text
            style={[
              styles.tabText,
              {
                color: selectedTab === 'signals'
                  ? theme.colors.primary
                  : theme.colors.onSurfaceVariant,
              },
            ]}
          >
            信号汇总
          </Text>
        </TouchableOpacity>
      </View>

      {/* 分组标签 */}
      {selectedTab === 'list' && renderGroupTabs()}

      {/* 内容区域 */}
      {selectedTab === 'list' ? (
        filteredFavorites.length === 0 ? (
          renderEmptyState()
        ) : (
          <FlatList
            data={filteredFavorites}
            renderItem={renderFavoriteItem}
            keyExtractor={item => item.id}
            contentContainerStyle={styles.listContent}
            refreshControl={
              <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />
            }
          />
        )
      ) : signalsLoading ? (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" />
        </View>
      ) : signals.length === 0 ? (
        renderEmptyState()
      ) : (
        <FlatList
          data={signals}
          renderItem={renderSignalItem}
          keyExtractor={(item, index) => `${item.ts_code}-${index}`}
          contentContainerStyle={styles.listContent}
          refreshControl={
            <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />
          }
        />
      )}

      {/* 浮动操作按钮 */}
      <FAB
        icon="plus"
        style={[styles.fab, { backgroundColor: theme.colors.primary }]}
        onPress={() => {
          // TODO: 导航到添加策略页面
          console.log('Add new strategy');
        }}
      />

      {/* 分组编辑对话框 */}
      <Portal>
        <Dialog visible={groupDialogVisible} onDismiss={() => setGroupDialogVisible(false)}>
          <Dialog.Title>{editingGroup ? '编辑分组' : '创建分组'}</Dialog.Title>
          <Dialog.Content>
            <TextInput
              label="分组名称"
              value={groupName}
              onChangeText={setGroupName}
              mode="outlined"
              style={styles.input}
            />
            <Text style={[styles.label, { color: theme.colors.onSurface }]}>
              分组颜色
            </Text>
            <View style={styles.colorPicker}>
              {['#1976d2', '#2e7d32', '#ed6c02', '#d32f2f', '#9c27b0', '#0288d1'].map(
                color => (
                  <TouchableOpacity
                    key={color}
                    style={[
                      styles.colorOption,
                      { backgroundColor: color },
                      groupColor === color && styles.colorOptionSelected,
                    ]}
                    onPress={() => setGroupColor(color)}
                  />
                )
              )}
            </View>
          </Dialog.Content>
          <Dialog.Actions>
            <Button onPress={() => setGroupDialogVisible(false)}>取消</Button>
            <Button onPress={handleSaveGroup}>保存</Button>
          </Dialog.Actions>
        </Dialog>
      </Portal>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  tabBar: {
    flexDirection: 'row',
    backgroundColor: '#fff',
    elevation: 2,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
  },
  tab: {
    flex: 1,
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: 12,
    gap: 8,
  },
  tabActive: {
    borderBottomWidth: 2,
  },
  tabText: {
    fontSize: 14,
    fontWeight: '600',
  },
  groupTabs: {
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  groupTabsContent: {
    paddingHorizontal: 16,
    paddingVertical: 12,
    gap: 8,
  },
  groupTabContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  groupTab: {
    marginRight: 4,
  },
  groupMenuButton: {
    margin: 0,
    padding: 0,
  },
  addGroupButton: {
    marginLeft: 8,
  },
  listContent: {
    padding: 16,
    gap: 12,
  },
  favoriteCard: {
    marginBottom: 8,
  },
  favoriteHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  favoriteInfo: {
    flex: 1,
  },
  favoriteName: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  favoriteCode: {
    fontSize: 14,
  },
  favoriteActions: {
    flexDirection: 'row',
    gap: 4,
  },
  favoriteDetails: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  groupChip: {
    height: 24,
  },
  favoriteDate: {
    fontSize: 12,
  },
  signalCard: {
    marginBottom: 8,
  },
  signalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  signalInfo: {
    flex: 1,
  },
  signalName: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  signalCode: {
    fontSize: 14,
  },
  signalChip: {
    height: 24,
  },
  signalDetails: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  signalPrice: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  signalDate: {
    fontSize: 12,
  },
  signalReason: {
    fontSize: 12,
    fontStyle: 'italic',
  },
  emptyState: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  emptyTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginTop: 16,
    marginBottom: 8,
  },
  emptyDescription: {
    fontSize: 14,
    textAlign: 'center',
  },
  fab: {
    position: 'absolute',
    right: 16,
    bottom: 16,
  },
  input: {
    marginBottom: 16,
  },
  label: {
    fontSize: 14,
    fontWeight: '600',
    marginBottom: 8,
  },
  colorPicker: {
    flexDirection: 'row',
    gap: 12,
    marginTop: 8,
  },
  colorOption: {
    width: 40,
    height: 40,
    borderRadius: 20,
    borderWidth: 2,
    borderColor: 'transparent',
  },
  colorOptionSelected: {
    borderColor: '#000',
  },
});

export default BacktestScreen;
