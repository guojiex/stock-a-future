/**
 * 策略编辑页面
 * 功能：添加或编辑收藏的股票策略
 */

import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
} from 'react-native';
import {
  Card,
  TextInput,
  Button,
  useTheme,
  Chip,
  ActivityIndicator,
  HelperText,
} from 'react-native-paper';
import { RouteProp, useRoute, useNavigation } from '@react-navigation/native';
import DateTimePicker from '@react-native-community/datetimepicker';
import Icon from 'react-native-vector-icons/MaterialIcons';

import {
  useAddFavoriteMutation,
  useUpdateFavoriteMutation,
  useGetFavoritesQuery,
  useGetGroupsQuery,
  useLazySearchStocksQuery,
  useLazyGetStockBasicQuery,
} from '@/services/api';
import { BacktestStackParamList } from '@/navigation/AppNavigator';

type StrategyEditRouteProp = RouteProp<BacktestStackParamList, 'StrategyEdit'>;

const StrategyEditScreen: React.FC = () => {
  const route = useRoute<StrategyEditRouteProp>();
  const navigation = useNavigation();
  const theme = useTheme();
  const { strategyId } = route.params;

  // API hooks
  const [searchStocks, { data: searchResults, isLoading: searchLoading }] =
    useLazySearchStocksQuery();
  const [getStockBasic] = useLazyGetStockBasicQuery();
  const { data: favoritesData } = useGetFavoritesQuery();
  const { data: groupsData } = useGetGroupsQuery();
  const [addFavorite, { isLoading: addLoading }] = useAddFavoriteMutation();
  const [updateFavorite, { isLoading: updateLoading }] = useUpdateFavoriteMutation();

  // 本地状态
  const [stockCode, setStockCode] = useState('');
  const [stockName, setStockName] = useState('');
  const [startDate, setStartDate] = useState<Date>(
    new Date(new Date().getFullYear() - 1, new Date().getMonth(), new Date().getDate())
  );
  const [endDate, setEndDate] = useState<Date>(new Date());
  const [selectedGroup, setSelectedGroup] = useState<string>('default');
  const [showStartDatePicker, setShowStartDatePicker] = useState(false);
  const [showEndDatePicker, setShowEndDatePicker] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [showSearchResults, setShowSearchResults] = useState(false);
  const [errors, setErrors] = useState<{ [key: string]: string }>({});

  // 数据
  const favorites = favoritesData?.data?.favorites || [];
  const groups = groupsData?.data?.groups || [];
  const stocks = searchResults?.data || [];

  // 如果是编辑模式，加载现有数据
  useEffect(() => {
    if (strategyId) {
      const favorite = favorites.find((f: any) => f.id === strategyId);
      if (favorite) {
        setStockCode(favorite.ts_code);
        setStockName(favorite.name);
        setSelectedGroup(favorite.group_id || 'default');
        if (favorite.start_date) {
          const start = parseDate(favorite.start_date);
          if (start) setStartDate(start);
        }
        if (favorite.end_date) {
          const end = parseDate(favorite.end_date);
          if (end) setEndDate(end);
        }
      }
    }
  }, [strategyId, favorites]);

  // 解析日期字符串 (YYYYMMDD 或 YYYY-MM-DD)
  const parseDate = (dateStr: string): Date | null => {
    if (!dateStr) return null;
    
    if (dateStr.length === 8) {
      // YYYYMMDD格式
      const year = parseInt(dateStr.substring(0, 4));
      const month = parseInt(dateStr.substring(4, 6)) - 1;
      const day = parseInt(dateStr.substring(6, 8));
      return new Date(year, month, day);
    } else if (dateStr.includes('-')) {
      // YYYY-MM-DD格式
      return new Date(dateStr);
    }
    
    return null;
  };

  // 格式化日期为 YYYYMMDD
  const formatDate = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}${month}${day}`;
  };

  // 格式化日期显示
  const formatDateDisplay = (date: Date): string => {
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  // 搜索股票
  const handleSearch = async (query: string) => {
    setSearchQuery(query);
    if (query.trim().length >= 2) {
      try {
        await searchStocks({ q: query, limit: 10 });
        setShowSearchResults(true);
      } catch (error) {
        console.error('搜索失败:', error);
      }
    } else {
      setShowSearchResults(false);
    }
  };

  // 选择股票
  const handleSelectStock = async (code: string, name: string) => {
    setStockCode(code);
    setStockName(name);
    setSearchQuery('');
    setShowSearchResults(false);
    setErrors({ ...errors, stock: '' });

    // 获取股票基本信息，设置默认时间范围
    try {
      const result = await getStockBasic(code).unwrap();
      if (result.data) {
        // 可以根据股票的上市时间等信息设置更合理的默认时间范围
        console.log('Stock basic info:', result.data);
      }
    } catch (error) {
      console.error('获取股票信息失败:', error);
    }
  };

  // 验证表单
  const validateForm = (): boolean => {
    const newErrors: { [key: string]: string } = {};

    if (!stockCode || !stockName) {
      newErrors.stock = '请选择股票';
    }

    if (startDate >= endDate) {
      newErrors.dateRange = '开始日期必须早于结束日期';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // 保存策略
  const handleSave = async () => {
    if (!validateForm()) {
      return;
    }

    try {
      const favoriteData = {
        ts_code: stockCode,
        name: stockName,
        start_date: formatDate(startDate),
        end_date: formatDate(endDate),
        group_id: selectedGroup === 'default' ? undefined : selectedGroup,
      };

      if (strategyId) {
        // 更新现有策略
        await updateFavorite({
          id: Number(strategyId),
          data: favoriteData,
        }).unwrap();
        Alert.alert('成功', '策略已更新', [
          { text: '确定', onPress: () => navigation.goBack() },
        ]);
      } else {
        // 添加新策略
        await addFavorite(favoriteData).unwrap();
        Alert.alert('成功', '策略已添加', [
          { text: '确定', onPress: () => navigation.goBack() },
        ]);
      }
    } catch (error: any) {
      Alert.alert(
        '错误',
        error.data?.message || '保存失败，请重试'
      );
    }
  };

  const isLoading = addLoading || updateLoading;

  return (
    <ScrollView style={styles.container} keyboardShouldPersistTaps="handled">
      {/* 股票选择 */}
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.sectionTitle, { color: theme.colors.onSurface }]}>
            选择股票
          </Text>
          
          {stockCode ? (
            <View style={styles.selectedStock}>
              <View style={styles.stockInfo}>
                <Text style={[styles.stockName, { color: theme.colors.onSurface }]}>
                  {stockName}
                </Text>
                <Text style={[styles.stockCode, { color: theme.colors.onSurfaceVariant }]}>
                  {stockCode}
                </Text>
              </View>
              <Button
                mode="outlined"
                onPress={() => {
                  setStockCode('');
                  setStockName('');
                }}
                compact
              >
                更换
              </Button>
            </View>
          ) : (
            <>
              <TextInput
                label="搜索股票代码或名称"
                value={searchQuery}
                onChangeText={handleSearch}
                mode="outlined"
                style={styles.input}
                left={<TextInput.Icon icon="magnify" />}
                error={!!errors.stock}
              />
              {errors.stock && (
                <HelperText type="error">{errors.stock}</HelperText>
              )}
              
              {searchLoading && (
                <View style={styles.searchLoading}>
                  <ActivityIndicator size="small" />
                </View>
              )}
              
              {showSearchResults && stocks.length > 0 && (
                <View style={styles.searchResults}>
                  {stocks.map((stock: any) => (
                    <TouchableOpacity
                      key={stock.ts_code}
                      style={styles.searchResultItem}
                      onPress={() => handleSelectStock(stock.ts_code, stock.name)}
                    >
                      <View>
                        <Text style={[styles.resultName, { color: theme.colors.onSurface }]}>
                          {stock.name}
                        </Text>
                        <Text style={[styles.resultCode, { color: theme.colors.onSurfaceVariant }]}>
                          {stock.ts_code}
                        </Text>
                      </View>
                      <Icon
                        name="chevron-right"
                        size={24}
                        color={theme.colors.onSurfaceVariant}
                      />
                    </TouchableOpacity>
                  ))}
                </View>
              )}
            </>
          )}
        </Card.Content>
      </Card>

      {/* 时间范围 */}
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.sectionTitle, { color: theme.colors.onSurface }]}>
            分析时间范围
          </Text>
          
          <TouchableOpacity
            style={styles.dateButton}
            onPress={() => setShowStartDatePicker(true)}
          >
            <Icon name="calendar-today" size={20} color={theme.colors.primary} />
            <View style={styles.dateContent}>
              <Text style={[styles.dateLabel, { color: theme.colors.onSurfaceVariant }]}>
                开始日期
              </Text>
              <Text style={[styles.dateValue, { color: theme.colors.onSurface }]}>
                {formatDateDisplay(startDate)}
              </Text>
            </View>
          </TouchableOpacity>

          <TouchableOpacity
            style={styles.dateButton}
            onPress={() => setShowEndDatePicker(true)}
          >
            <Icon name="calendar-today" size={20} color={theme.colors.primary} />
            <View style={styles.dateContent}>
              <Text style={[styles.dateLabel, { color: theme.colors.onSurfaceVariant }]}>
                结束日期
              </Text>
              <Text style={[styles.dateValue, { color: theme.colors.onSurface }]}>
                {formatDateDisplay(endDate)}
              </Text>
            </View>
          </TouchableOpacity>

          {errors.dateRange && (
            <HelperText type="error">{errors.dateRange}</HelperText>
          )}

          {showStartDatePicker && (
            <DateTimePicker
              value={startDate}
              mode="date"
              display="default"
              onChange={(event, selectedDate) => {
                setShowStartDatePicker(false);
                if (selectedDate) {
                  setStartDate(selectedDate);
                  setErrors({ ...errors, dateRange: '' });
                }
              }}
            />
          )}

          {showEndDatePicker && (
            <DateTimePicker
              value={endDate}
              mode="date"
              display="default"
              onChange={(event, selectedDate) => {
                setShowEndDatePicker(false);
                if (selectedDate) {
                  setEndDate(selectedDate);
                  setErrors({ ...errors, dateRange: '' });
                }
              }}
            />
          )}
        </Card.Content>
      </Card>

      {/* 分组选择 */}
      <Card style={styles.card}>
        <Card.Content>
          <Text style={[styles.sectionTitle, { color: theme.colors.onSurface }]}>
            选择分组
          </Text>
          
          <View style={styles.groupChips}>
            {groups.map((group: any) => (
              <Chip
                key={group.id}
                mode={selectedGroup === group.id ? 'flat' : 'outlined'}
                selected={selectedGroup === group.id}
                onPress={() => setSelectedGroup(group.id)}
                style={[
                  styles.groupChip,
                  selectedGroup === group.id && { backgroundColor: group.color },
                ]}
              >
                {group.name}
              </Chip>
            ))}
          </View>
        </Card.Content>
      </Card>

      {/* 操作按钮 */}
      <View style={styles.actions}>
        <Button
          mode="outlined"
          onPress={() => navigation.goBack()}
          style={styles.button}
          disabled={isLoading}
        >
          取消
        </Button>
        <Button
          mode="contained"
          onPress={handleSave}
          style={styles.button}
          loading={isLoading}
          disabled={isLoading}
        >
          {strategyId ? '更新' : '添加'}
        </Button>
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  card: {
    margin: 16,
    marginBottom: 8,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  input: {
    marginBottom: 8,
  },
  selectedStock: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 12,
    backgroundColor: '#f5f5f5',
    borderRadius: 8,
  },
  stockInfo: {
    flex: 1,
  },
  stockName: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  stockCode: {
    fontSize: 14,
  },
  searchLoading: {
    padding: 16,
    alignItems: 'center',
  },
  searchResults: {
    marginTop: 8,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#e0e0e0',
    maxHeight: 300,
  },
  searchResultItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  resultName: {
    fontSize: 14,
    fontWeight: '500',
    marginBottom: 4,
  },
  resultCode: {
    fontSize: 12,
  },
  dateButton: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 12,
    marginBottom: 12,
    backgroundColor: '#f5f5f5',
    borderRadius: 8,
    gap: 12,
  },
  dateContent: {
    flex: 1,
  },
  dateLabel: {
    fontSize: 12,
    marginBottom: 4,
  },
  dateValue: {
    fontSize: 16,
    fontWeight: '500',
  },
  groupChips: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  groupChip: {
    marginBottom: 8,
  },
  actions: {
    flexDirection: 'row',
    gap: 12,
    padding: 16,
    paddingBottom: 32,
  },
  button: {
    flex: 1,
  },
});

export default StrategyEditScreen;
