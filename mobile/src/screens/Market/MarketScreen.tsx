/**
 * 市场页面 - 股票查询和分析
 * 采用原始Web版设计: 搜索框 + 日期选择 + 数据展示
 */

import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Platform,
} from 'react-native';
import {
  Appbar,
  Card,
  Button,
  useTheme,
  TextInput,
  Menu,
  Chip,
  Divider,
  Portal,
  Modal,
} from 'react-native-paper';
import Icon from 'react-native-vector-icons/MaterialIcons';
// @ts-ignore - DateTimePicker types will be available after installation
import DateTimePicker from '@react-native-community/datetimepicker';

import { useGetHealthStatusQuery } from '@/services/api';
import { apiService } from '@/services/apiService';
import { useAppDispatch, useAppSelector } from '@/hooks/redux';
import { setConnectionStatus } from '@/store/slices/appSlice';
import KLineChart from '@/components/KLineChart';
import TechnicalIndicators from '@/components/TechnicalIndicators';
import PredictionSignals from '@/components/PredictionSignals';
import FundamentalData from '@/components/FundamentalData';

interface SearchSuggestion {
  ts_code: string;
  name: string;
}

type TabType = 'daily-data' | 'indicators' | 'predictions' | 'fundamental';

const MarketScreen: React.FC = () => {
  const theme = useTheme();
  const dispatch = useAppDispatch();
  
  // 搜索相关状态
  const [searchQuery, setSearchQuery] = useState('');
  const [stockCode, setStockCode] = useState('000001.SZ');
  const [stockName, setStockName] = useState('平安银行');
  const [suggestions, setSuggestions] = useState<SearchSuggestion[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  
  // 日期选择状态
  const [startDate, setStartDate] = useState(new Date(Date.now() - 180 * 24 * 60 * 60 * 1000)); // 默认6个月前
  const [endDate, setEndDate] = useState(new Date());
  const [showStartDatePicker, setShowStartDatePicker] = useState(false);
  const [showEndDatePicker, setShowEndDatePicker] = useState(false);
  
  // 结果展示状态
  const [showResults, setShowResults] = useState(false);
  const [selectedTab, setSelectedTab] = useState<TabType>('daily-data');
  const [loading, setLoading] = useState(false);
  
  // Redux状态
  const connectionStatus = useAppSelector(state => state.app.connectionStatus);
  
  // API查询
  const {
    data: healthData,
    error: healthError,
  } = useGetHealthStatusQuery(undefined, {
    pollingInterval: 30000,
  });
  
  // 监听健康状态
  useEffect(() => {
    if (healthData?.success) {
      dispatch(setConnectionStatus('connected'));
    } else if (healthError) {
      dispatch(setConnectionStatus('error'));
    }
  }, [healthData, healthError, dispatch]);
  
  // 处理搜索
  const handleSearch = async (query: string) => {
    setSearchQuery(query);
    
    if (!query.trim()) {
      setSuggestions([]);
      setShowSuggestions(false);
      return;
    }
    
    try {
      const response = await apiService.searchStocks(query, 10);
      if (response.success && response.data) {
        setSuggestions(response.data as SearchSuggestion[]);
        setShowSuggestions(true);
      }
    } catch (error) {
      console.error('搜索失败:', error);
    }
  };
  
  // 选择建议项
  const handleSelectSuggestion = (suggestion: SearchSuggestion) => {
    setStockCode(suggestion.ts_code);
    setStockName(suggestion.name);
    setSearchQuery('');
    setSuggestions([]);
    setShowSuggestions(false);
  };
  
  // 日期快捷选择
  const handleQuickDateSelect = (days: number) => {
    const end = new Date();
    const start = new Date(Date.now() - days * 24 * 60 * 60 * 1000);
    setStartDate(start);
    setEndDate(end);
  };
  
  // 格式化日期
  const formatDate = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };
  
  // 获取数据并显示结果
  const handleShowData = (tab: TabType) => {
    setSelectedTab(tab);
    setShowResults(true);
  };
  
  // 渲染连接状态
  const renderConnectionStatus = () => {
    const getStatusColor = () => {
      switch (connectionStatus) {
        case 'connected':
          return theme.colors.primary;
        case 'connecting':
          return theme.colors.secondary;
        case 'error':
          return theme.colors.error;
        default:
          return theme.colors.onSurfaceVariant;
      }
    };
    
    const getStatusText = () => {
      switch (connectionStatus) {
        case 'connected':
          return '已连接';
        case 'connecting':
          return '连接中...';
        case 'error':
          return '连接失败';
        default:
          return '未连接';
      }
    };
    
    return (
      <View style={styles.statusContainer}>
        <View style={[styles.statusDot, { backgroundColor: getStatusColor() }]} />
        <Text style={[styles.statusText, { color: getStatusColor() }]}>
          {getStatusText()}
        </Text>
      </View>
    );
  };
  
  return (
    <View style={styles.container}>
      <Appbar.Header>
        <Appbar.Content title="股票查询" />
        {renderConnectionStatus()}
      </Appbar.Header>
      
      <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
        {/* 头部卡片 */}
        <Card style={styles.headerCard} mode="elevated">
        <Card.Content>
          <View style={styles.headerContent}>
            <View>
                <Text style={[styles.headerTitle, { color: theme.colors.onSurface }]}>
                  🚀 Stock-A-Future
                </Text>
                <Text style={[styles.headerSubtitle, { color: theme.colors.onSurfaceVariant }]}>
                  A股股票买卖点预测系统
                </Text>
            </View>
          </View>
        </Card.Content>
      </Card>
      
        {/* 股票查询卡片 */}
        <Card style={styles.card} mode="elevated">
          <Card.Content>
            <Text style={[styles.cardTitle, { color: theme.colors.onSurface }]}>
              📊 股票查询
            </Text>
            
            {/* 搜索框 */}
            <View style={styles.searchContainer}>
              <Text style={[styles.label, { color: theme.colors.onSurfaceVariant }]}>
                股票搜索
              </Text>
              <TextInput
                mode="outlined"
                placeholder="输入股票名称或代码，如: 平安银行、000001"
                value={searchQuery}
                onChangeText={handleSearch}
                right={<TextInput.Icon icon="magnify" />}
                style={styles.input}
              />
              <Text style={[styles.hint, { color: theme.colors.onSurfaceVariant }]}>
                支持股票名称和代码搜索
              </Text>
              
              {/* 搜索建议 */}
              {showSuggestions && suggestions.length > 0 && (
                <Card style={styles.suggestionsCard} mode="outlined">
                  {suggestions.map((suggestion, index) => (
                    <React.Fragment key={suggestion.ts_code}>
                      <TouchableOpacity
                        style={styles.suggestionItem}
                        onPress={() => handleSelectSuggestion(suggestion)}
                      >
                        <Text style={[styles.suggestionName, { color: theme.colors.onSurface }]}>
                          {suggestion.name}
                        </Text>
                        <Text style={[styles.suggestionCode, { color: theme.colors.onSurfaceVariant }]}>
                          {suggestion.ts_code}
                        </Text>
                      </TouchableOpacity>
                      {index < suggestions.length - 1 && <Divider />}
                    </React.Fragment>
                  ))}
                </Card>
              )}
    </View>
            
            {/* 股票代码输入 */}
            <View style={styles.inputContainer}>
              <Text style={[styles.label, { color: theme.colors.onSurfaceVariant }]}>
                股票代码
              </Text>
              <TextInput
                mode="outlined"
                placeholder="例如: 000001.SZ"
                value={stockCode}
                onChangeText={setStockCode}
                style={styles.input}
              />
              <Text style={[styles.hint, { color: theme.colors.onSurfaceVariant }]}>
                支持格式: 000001.SZ (深圳), 600000.SH (上海)
              </Text>
        </View>
            
            {/* 当前选中的股票 */}
            {stockName && (
              <View style={styles.selectedStockContainer}>
                <Chip
                  icon="check-circle"
                  mode="flat"
                  style={[styles.selectedStockChip, { backgroundColor: theme.colors.primaryContainer }]}
                  textStyle={{ color: theme.colors.onPrimaryContainer }}
                >
                  {stockName} ({stockCode})
                </Chip>
      </View>
            )}
            
            {/* 日期选择 */}
            <Divider style={styles.divider} />
            <Text style={[styles.sectionTitle, { color: theme.colors.onSurface }]}>
              时间范围
            </Text>
            
            <View style={styles.dateContainer}>
              <View style={styles.dateInputGroup}>
                <Text style={[styles.label, { color: theme.colors.onSurfaceVariant }]}>
                  开始日期
                </Text>
                <TouchableOpacity onPress={() => setShowStartDatePicker(true)}>
                  <TextInput
                    mode="outlined"
                    value={formatDate(startDate)}
                    editable={false}
                    right={<TextInput.Icon icon="calendar" />}
                    style={styles.input}
                    pointerEvents="none"
                  />
                </TouchableOpacity>
              </View>
              
              <View style={styles.dateInputGroup}>
                <Text style={[styles.label, { color: theme.colors.onSurfaceVariant }]}>
                  结束日期
                </Text>
                <TouchableOpacity onPress={() => setShowEndDatePicker(true)}>
                  <TextInput
                    mode="outlined"
                    value={formatDate(endDate)}
                    editable={false}
                    right={<TextInput.Icon icon="calendar" />}
                    style={styles.input}
                    pointerEvents="none"
                  />
          </TouchableOpacity>
        </View>
      </View>
            
            {/* 日期快捷按钮 */}
            <View style={styles.quickDateContainer}>
              <Text style={[styles.label, { color: theme.colors.onSurfaceVariant }]}>
                快捷选择
              </Text>
              <View style={styles.quickDateButtons}>
                <Chip mode="outlined" onPress={() => handleQuickDateSelect(30)} style={styles.quickDateChip}>
                  近1个月
                </Chip>
                <Chip mode="outlined" onPress={() => handleQuickDateSelect(90)} style={styles.quickDateChip}>
                  近3个月
                </Chip>
                <Chip mode="outlined" onPress={() => handleQuickDateSelect(180)} style={styles.quickDateChip}>
                  近半年
                </Chip>
                <Chip mode="outlined" onPress={() => handleQuickDateSelect(365)} style={styles.quickDateChip}>
                  近1年
                </Chip>
              </View>
            </View>
            
            {/* 操作按钮 */}
            <Divider style={styles.divider} />
            <View style={styles.buttonGroup}>
              <Button
                mode="contained"
                icon="chart-line"
                onPress={() => handleShowData('daily-data')}
                style={styles.actionButton}
              >
                获取日线数据
              </Button>
              <Button
                mode="contained-tonal"
                icon="chart-bell-curve"
                onPress={() => handleShowData('indicators')}
                style={styles.actionButton}
              >
                技术指标
              </Button>
              <Button
                mode="contained-tonal"
                icon="target"
                onPress={() => handleShowData('predictions')}
                style={styles.actionButton}
              >
                买卖预测
              </Button>
              <Button
                mode="contained-tonal"
                icon="cash"
                onPress={() => handleShowData('fundamental')}
                style={styles.actionButton}
              >
                基本面数据
              </Button>
            </View>
          </Card.Content>
        </Card>
        
        {/* 结果展示区域 */}
        {showResults && (
          <Card style={styles.card} mode="elevated">
            <Card.Content>
              <View style={styles.resultHeader}>
                <Text style={[styles.cardTitle, { color: theme.colors.onSurface }]}>
                  📊 {stockName} ({stockCode})
                </Text>
                <TouchableOpacity onPress={() => setShowResults(false)}>
                  <Icon name="close" size={24} color={theme.colors.onSurfaceVariant} />
                </TouchableOpacity>
              </View>
              
              {/* Tab 导航 */}
              <ScrollView
                horizontal
                showsHorizontalScrollIndicator={false}
                style={styles.tabContainer}
              >
                <Chip
                  mode={selectedTab === 'daily-data' ? 'flat' : 'outlined'}
                  selected={selectedTab === 'daily-data'}
                  onPress={() => setSelectedTab('daily-data')}
                  style={styles.tabChip}
                >
                  📈 日线数据
                </Chip>
                <Chip
                  mode={selectedTab === 'indicators' ? 'flat' : 'outlined'}
                  selected={selectedTab === 'indicators'}
                  onPress={() => setSelectedTab('indicators')}
                  style={styles.tabChip}
                >
                  📊 技术指标
                </Chip>
                <Chip
                  mode={selectedTab === 'predictions' ? 'flat' : 'outlined'}
                  selected={selectedTab === 'predictions'}
                  onPress={() => setSelectedTab('predictions')}
                  style={styles.tabChip}
                >
                  🎯 买卖预测
                </Chip>
                <Chip
                  mode={selectedTab === 'fundamental' ? 'flat' : 'outlined'}
                  selected={selectedTab === 'fundamental'}
                  onPress={() => setSelectedTab('fundamental')}
                  style={styles.tabChip}
                >
                  💰 基本面
                </Chip>
              </ScrollView>
              
              {/* Tab 内容 */}
              <View style={styles.tabContent}>
                {selectedTab === 'daily-data' && (
                  <KLineChart
                    stockCode={stockCode}
                    stockName={stockName}
                  />
                )}
                
                {selectedTab === 'indicators' && (
                  <TechnicalIndicators
                    stockCode={stockCode}
                    stockName={stockName}
                  />
                )}
                
                {selectedTab === 'predictions' && (
                  <PredictionSignals
                    stockCode={stockCode}
                    stockName={stockName}
                  />
                )}
                
                {selectedTab === 'fundamental' && (
                  <FundamentalData
                    stockCode={stockCode}
                    stockName={stockName}
                  />
                )}
              </View>
            </Card.Content>
          </Card>
        )}
      </ScrollView>
      
      {/* 日期选择器 */}
      {showStartDatePicker && (
        <DateTimePicker
          value={startDate}
          mode="date"
          display={Platform.OS === 'ios' ? 'spinner' : 'default'}
          onChange={(_event: any, selectedDate?: Date) => {
            setShowStartDatePicker(false);
            if (selectedDate) {
              setStartDate(selectedDate);
            }
          }}
        />
      )}
      
      {showEndDatePicker && (
        <DateTimePicker
          value={endDate}
          mode="date"
          display={Platform.OS === 'ios' ? 'spinner' : 'default'}
          onChange={(_event: any, selectedDate?: Date) => {
            setShowEndDatePicker(false);
            if (selectedDate) {
              setEndDate(selectedDate);
            }
          }}
        />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  content: {
    flex: 1,
  },
  statusContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginRight: 16,
  },
  statusDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginRight: 6,
  },
  statusText: {
    fontSize: 12,
    fontWeight: '500',
  },
  headerCard: {
    margin: 16,
    marginBottom: 8,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: 'bold',
  },
  headerSubtitle: {
    fontSize: 14,
    marginTop: 4,
  },
  card: {
    margin: 16,
    marginTop: 8,
  },
  cardTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  searchContainer: {
    marginBottom: 16,
  },
  inputContainer: {
    marginBottom: 16,
  },
  label: {
    fontSize: 14,
    fontWeight: '500',
    marginBottom: 8,
  },
  input: {
    marginBottom: 4,
  },
  hint: {
    fontSize: 12,
    marginTop: 4,
  },
  suggestionsCard: {
    marginTop: 8,
    maxHeight: 250,
  },
  suggestionItem: {
    padding: 12,
  },
  suggestionName: {
    fontSize: 16,
    fontWeight: '600',
  },
  suggestionCode: {
    fontSize: 12,
    marginTop: 2,
  },
  selectedStockContainer: {
    marginBottom: 16,
  },
  selectedStockChip: {
    alignSelf: 'flex-start',
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 12,
  },
  dateContainer: {
    flexDirection: 'row',
    gap: 12,
    marginBottom: 16,
  },
  dateInputGroup: {
    flex: 1,
  },
  quickDateContainer: {
    marginBottom: 16,
  },
  quickDateButtons: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
    marginTop: 8,
  },
  quickDateChip: {
    marginRight: 4,
    marginBottom: 4,
  },
  divider: {
    marginVertical: 16,
  },
  buttonGroup: {
    gap: 8,
  },
  actionButton: {
    marginBottom: 8,
  },
  resultHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  tabContainer: {
    marginBottom: 16,
  },
  tabChip: {
    marginRight: 8,
  },
  tabContent: {
    minHeight: 300,
  },
});

export default MarketScreen;
