/**
 * 策略信号详情页面
 * 显示单个股票的详细信号分析
 */

import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  RefreshControl,
  TouchableOpacity,
} from 'react-native';
import {
  Card,
  Button,
  Chip,
  useTheme,
  ActivityIndicator,
  Divider,
} from 'react-native-paper';
import { RouteProp, useRoute, useNavigation } from '@react-navigation/native';
import Icon from 'react-native-vector-icons/MaterialIcons';

import {
  useGetPredictionsQuery,
  useRecognizePatternsQuery,
  useGetPatternSummaryQuery,
} from '@/services/api';
import { BacktestStackParamList } from '@/navigation/AppNavigator';

type BacktestResultRouteProp = RouteProp<BacktestStackParamList, 'BacktestResult'>;

interface SignalData {
  type: string;
  probability: number;
  reason: string;
  indicator?: string;
}

const BacktestResultScreen: React.FC = () => {
  const route = useRoute<BacktestResultRouteProp>();
  const navigation = useNavigation();
  const theme = useTheme();
  const { backtestId } = route.params;

  // 在实际应用中，backtestId 会是股票代码
  const stockCode = backtestId;

  // API查询
  const {
    data: predictionsData,
    isLoading: predictionsLoading,
    refetch: refetchPredictions,
  } = useGetPredictionsQuery(stockCode);
  
  const {
    data: patternsData,
    isLoading: patternsLoading,
    refetch: refetchPatterns,
  } = useRecognizePatternsQuery({ tsCode: stockCode });

  const {
    data: summaryData,
    isLoading: summaryLoading,
    refetch: refetchSummary,
  } = useGetPatternSummaryQuery(stockCode);

  // 本地状态
  const [selectedTab, setSelectedTab] = useState<'predictions' | 'patterns' | 'summary'>('predictions');
  const [refreshing, setRefreshing] = useState(false);

  // 获取数据
  const predictions = predictionsData?.data || [];
  const patterns = patternsData?.data?.patterns || [];
  const summary = summaryData?.data || {};

  // 处理刷新
  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await Promise.all([
        refetchPredictions(),
        refetchPatterns(),
        refetchSummary(),
      ]);
    } finally {
      setRefreshing(false);
    }
  };

  // 获取信号颜色
  const getSignalColor = (type: string) => {
    switch (type.toUpperCase()) {
      case 'BUY':
        return '#4caf50';
      case 'SELL':
        return '#f44336';
      case 'HOLD':
        return '#ff9800';
      default:
        return theme.colors.onSurfaceVariant;
    }
  };

  // 获取信号文本
  const getSignalText = (type: string) => {
    switch (type.toUpperCase()) {
      case 'BUY':
        return '买入';
      case 'SELL':
        return '卖出';
      case 'HOLD':
        return '持有';
      default:
        return type;
    }
  };

  // 渲染预测信号
  const renderPredictions = () => {
    if (predictions.length === 0) {
      return (
        <View style={styles.emptyState}>
          <Icon name="show-chart" size={48} color={theme.colors.onSurfaceVariant} />
          <Text style={[styles.emptyText, { color: theme.colors.onSurfaceVariant }]}>
            暂无预测信号
          </Text>
        </View>
      );
    }

    return predictions.map((prediction: any, index: number) => {
      const signalColor = getSignalColor(prediction.type);
      const signalText = getSignalText(prediction.type);
      const probability = parseFloat(prediction.probability || 0);

      return (
        <Card key={index} style={styles.signalCard}>
          <Card.Content>
            <View style={styles.signalHeader}>
              <Chip
                mode="flat"
                style={[styles.signalChip, { backgroundColor: signalColor }]}
                textStyle={{ color: '#fff', fontWeight: 'bold' }}
              >
                {signalText}
              </Chip>
              <View style={styles.probabilityContainer}>
                <Text style={[styles.probabilityLabel, { color: theme.colors.onSurfaceVariant }]}>
                  置信度
                </Text>
                <Text style={[styles.probabilityValue, { color: signalColor }]}>
                  {probability.toFixed(1)}%
                </Text>
              </View>
            </View>

            {prediction.indicator && (
              <View style={styles.indicatorRow}>
                <Icon name="trending-up" size={16} color={theme.colors.onSurfaceVariant} />
                <Text style={[styles.indicatorText, { color: theme.colors.onSurfaceVariant }]}>
                  {prediction.indicator}
                </Text>
              </View>
            )}

            {prediction.reason && (
              <View style={styles.reasonContainer}>
                <Text style={[styles.reasonLabel, { color: theme.colors.onSurfaceVariant }]}>
                  分析理由：
                </Text>
                <Text style={[styles.reasonText, { color: theme.colors.onSurface }]}>
                  {prediction.reason}
                </Text>
              </View>
            )}

            {prediction.target_price && (
              <View style={styles.targetPriceRow}>
                <Text style={[styles.targetPriceLabel, { color: theme.colors.onSurfaceVariant }]}>
                  目标价位：
                </Text>
                <Text style={[styles.targetPriceValue, { color: theme.colors.primary }]}>
                  ¥{prediction.target_price}
                </Text>
              </View>
            )}
          </Card.Content>
        </Card>
      );
    });
  };

  // 渲染形态识别
  const renderPatterns = () => {
    if (patterns.length === 0) {
      return (
        <View style={styles.emptyState}>
          <Icon name="category" size={48} color={theme.colors.onSurfaceVariant} />
          <Text style={[styles.emptyText, { color: theme.colors.onSurfaceVariant }]}>
            暂无形态识别
          </Text>
        </View>
      );
    }

    return patterns.map((pattern: any, index: number) => {
      const confidence = parseFloat(pattern.confidence || 0);
      const isBullish = pattern.type === 'bullish';

      return (
        <Card key={index} style={styles.patternCard}>
          <Card.Content>
            <View style={styles.patternHeader}>
              <View style={styles.patternInfo}>
                <Text style={[styles.patternName, { color: theme.colors.onSurface }]}>
                  {pattern.name || '未知形态'}
                </Text>
                <Chip
                  mode="flat"
                  style={[
                    styles.patternTypeChip,
                    { backgroundColor: isBullish ? '#e8f5e9' : '#ffebee' },
                  ]}
                  textStyle={{
                    color: isBullish ? '#2e7d32' : '#c62828',
                    fontSize: 12,
                  }}
                >
                  {isBullish ? '看涨' : '看跌'}
                </Chip>
              </View>
              <View style={styles.confidenceContainer}>
                <Text style={[styles.confidenceLabel, { color: theme.colors.onSurfaceVariant }]}>
                  可信度
                </Text>
                <Text
                  style={[
                    styles.confidenceValue,
                    { color: confidence >= 70 ? '#4caf50' : '#ff9800' },
                  ]}
                >
                  {confidence.toFixed(0)}%
                </Text>
              </View>
            </View>

            {pattern.description && (
              <Text style={[styles.patternDescription, { color: theme.colors.onSurfaceVariant }]}>
                {pattern.description}
              </Text>
            )}

            {pattern.detected_at && (
              <View style={styles.detectedRow}>
                <Icon name="access-time" size={14} color={theme.colors.onSurfaceVariant} />
                <Text style={[styles.detectedText, { color: theme.colors.onSurfaceVariant }]}>
                  识别于 {new Date(pattern.detected_at).toLocaleDateString('zh-CN')}
                </Text>
              </View>
            )}
          </Card.Content>
        </Card>
      );
    });
  };

  // 渲染摘要
  const renderSummary = () => {
    if (!summary || Object.keys(summary).length === 0) {
      return (
        <View style={styles.emptyState}>
          <Icon name="summarize" size={48} color={theme.colors.onSurfaceVariant} />
          <Text style={[styles.emptyText, { color: theme.colors.onSurfaceVariant }]}>
            暂无汇总信息
          </Text>
        </View>
      );
    }

    return (
      <View>
        {/* 总体评估 */}
        <Card style={styles.summaryCard}>
          <Card.Content>
            <Text style={[styles.summaryTitle, { color: theme.colors.onSurface }]}>
              总体评估
            </Text>
            <View style={styles.summaryRow}>
              <Text style={[styles.summaryLabel, { color: theme.colors.onSurfaceVariant }]}>
                推荐操作：
              </Text>
              <Chip
                mode="flat"
                style={[
                  styles.summaryChip,
                  { backgroundColor: getSignalColor(summary.recommendation || 'HOLD') },
                ]}
                textStyle={{ color: '#fff', fontWeight: 'bold' }}
              >
                {getSignalText(summary.recommendation || 'HOLD')}
              </Chip>
            </View>

            {summary.overall_confidence && (
              <View style={styles.summaryRow}>
                <Text style={[styles.summaryLabel, { color: theme.colors.onSurfaceVariant }]}>
                  综合置信度：
                </Text>
                <Text style={[styles.summaryValue, { color: theme.colors.primary }]}>
                  {parseFloat(summary.overall_confidence).toFixed(1)}%
                </Text>
              </View>
            )}
          </Card.Content>
        </Card>

        {/* 信号统计 */}
        {summary.signal_counts && (
          <Card style={styles.summaryCard}>
            <Card.Content>
              <Text style={[styles.summaryTitle, { color: theme.colors.onSurface }]}>
                信号统计
              </Text>
              <View style={styles.statsGrid}>
                <View style={styles.statItem}>
                  <Text style={[styles.statValue, { color: '#4caf50' }]}>
                    {summary.signal_counts.buy || 0}
                  </Text>
                  <Text style={[styles.statLabel, { color: theme.colors.onSurfaceVariant }]}>
                    买入信号
                  </Text>
                </View>
                <View style={styles.statItem}>
                  <Text style={[styles.statValue, { color: '#f44336' }]}>
                    {summary.signal_counts.sell || 0}
                  </Text>
                  <Text style={[styles.statLabel, { color: theme.colors.onSurfaceVariant }]}>
                    卖出信号
                  </Text>
                </View>
                <View style={styles.statItem}>
                  <Text style={[styles.statValue, { color: '#ff9800' }]}>
                    {summary.signal_counts.hold || 0}
                  </Text>
                  <Text style={[styles.statLabel, { color: theme.colors.onSurfaceVariant }]}>
                    持有信号
                  </Text>
                </View>
              </View>
            </Card.Content>
          </Card>
        )}
      </View>
    );
  };

  const isLoading = predictionsLoading || patternsLoading || summaryLoading;

  if (isLoading && !refreshing) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      {/* Tab导航 */}
      <View style={styles.tabBar}>
        <TouchableOpacity
          style={[
            styles.tab,
            selectedTab === 'predictions' && styles.tabActive,
            { borderBottomColor: theme.colors.primary },
          ]}
          onPress={() => setSelectedTab('predictions')}
        >
          <Text
            style={[
              styles.tabText,
              {
                color:
                  selectedTab === 'predictions'
                    ? theme.colors.primary
                    : theme.colors.onSurfaceVariant,
              },
            ]}
          >
            预测信号
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[
            styles.tab,
            selectedTab === 'patterns' && styles.tabActive,
            { borderBottomColor: theme.colors.primary },
          ]}
          onPress={() => setSelectedTab('patterns')}
        >
          <Text
            style={[
              styles.tabText,
              {
                color:
                  selectedTab === 'patterns'
                    ? theme.colors.primary
                    : theme.colors.onSurfaceVariant,
              },
            ]}
          >
            形态识别
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[
            styles.tab,
            selectedTab === 'summary' && styles.tabActive,
            { borderBottomColor: theme.colors.primary },
          ]}
          onPress={() => setSelectedTab('summary')}
        >
          <Text
            style={[
              styles.tabText,
              {
                color:
                  selectedTab === 'summary'
                    ? theme.colors.primary
                    : theme.colors.onSurfaceVariant,
              },
            ]}
          >
            汇总
          </Text>
        </TouchableOpacity>
      </View>

      {/* 内容区域 */}
      <ScrollView
        style={styles.content}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />}
      >
        <View style={styles.contentPadding}>
          {selectedTab === 'predictions' && renderPredictions()}
          {selectedTab === 'patterns' && renderPatterns()}
          {selectedTab === 'summary' && renderSummary()}
        </View>
      </ScrollView>
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
    paddingVertical: 16,
    alignItems: 'center',
  },
  tabActive: {
    borderBottomWidth: 2,
  },
  tabText: {
    fontSize: 14,
    fontWeight: '600',
  },
  content: {
    flex: 1,
  },
  contentPadding: {
    padding: 16,
  },
  emptyState: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 64,
  },
  emptyText: {
    fontSize: 16,
    marginTop: 16,
  },
  signalCard: {
    marginBottom: 16,
  },
  signalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  signalChip: {
    height: 32,
  },
  probabilityContainer: {
    alignItems: 'flex-end',
  },
  probabilityLabel: {
    fontSize: 12,
    marginBottom: 2,
  },
  probabilityValue: {
    fontSize: 20,
    fontWeight: 'bold',
  },
  indicatorRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginBottom: 12,
  },
  indicatorText: {
    fontSize: 14,
  },
  reasonContainer: {
    marginTop: 8,
    padding: 12,
    backgroundColor: '#f5f5f5',
    borderRadius: 8,
  },
  reasonLabel: {
    fontSize: 12,
    marginBottom: 4,
  },
  reasonText: {
    fontSize: 14,
    lineHeight: 20,
  },
  targetPriceRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: 12,
    paddingTop: 12,
    borderTopWidth: 1,
    borderTopColor: '#e0e0e0',
  },
  targetPriceLabel: {
    fontSize: 14,
    marginRight: 8,
  },
  targetPriceValue: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  patternCard: {
    marginBottom: 16,
  },
  patternHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 12,
  },
  patternInfo: {
    flex: 1,
    marginRight: 12,
  },
  patternName: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  patternTypeChip: {
    alignSelf: 'flex-start',
  },
  confidenceContainer: {
    alignItems: 'flex-end',
  },
  confidenceLabel: {
    fontSize: 12,
    marginBottom: 2,
  },
  confidenceValue: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  patternDescription: {
    fontSize: 14,
    lineHeight: 20,
    marginBottom: 12,
  },
  detectedRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  detectedText: {
    fontSize: 12,
  },
  summaryCard: {
    marginBottom: 16,
  },
  summaryTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  summaryRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  summaryLabel: {
    fontSize: 14,
  },
  summaryValue: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  summaryChip: {
    height: 28,
  },
  statsGrid: {
    flexDirection: 'row',
    justifyContent: 'space-around',
  },
  statItem: {
    alignItems: 'center',
  },
  statValue: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  statLabel: {
    fontSize: 12,
  },
});

export default BacktestResultScreen;
