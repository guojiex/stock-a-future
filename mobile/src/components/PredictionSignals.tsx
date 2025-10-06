/**
 * 买卖预测信号展示组件
 */

import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, ScrollView, ActivityIndicator, TouchableOpacity } from 'react-native';
import { Card, useTheme, Chip, Button, Divider } from 'react-native-paper';
import { apiService } from '@/services/apiService';
import { SkeletonPredictionCard, SkeletonCard } from '@/components/common/SkeletonLoader';

interface PredictionSignalsProps {
  stockCode: string;
  stockName: string;
}

interface PredictionItem {
  type: 'BUY' | 'SELL' | 'HOLD';
  price: number;
  signal_date: string;
  probability: number;
  confidence?: number;
  reasons?: string[];
  description?: string;
  stop_loss?: number;
  target_price?: number;
}

interface PredictionData {
  predictions: PredictionItem[];
  summary?: {
    total_signals: number;
    buy_signals: number;
    sell_signals: number;
    avg_confidence: number;
  };
}

const PredictionSignals: React.FC<PredictionSignalsProps> = ({ stockCode, stockName }) => {
  const theme = useTheme();
  const [loading, setLoading] = useState(false);
  const [predictions, setPredictions] = useState<PredictionData | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [expandedItems, setExpandedItems] = useState<Set<number>>(new Set());

  // 加载预测数据
  const loadPredictions = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await apiService.getPredictions(stockCode);
      
      if (response.success && response.data) {
        setPredictions(response.data as PredictionData);
      } else {
        setError(response.error || '加载预测数据失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载预测数据失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPredictions();
  }, [stockCode]);

  // 切换展开状态
  const toggleExpanded = (index: number) => {
    const newExpanded = new Set(expandedItems);
    if (newExpanded.has(index)) {
      newExpanded.delete(index);
    } else {
      newExpanded.add(index);
    }
    setExpandedItems(newExpanded);
  };

  // 获取信号类型颜色
  const getSignalColor = (type: string) => {
    switch (type) {
      case 'BUY':
        return theme.colors.primary;
      case 'SELL':
        return theme.colors.error;
      default:
        return theme.colors.outline;
    }
  };

  // 获取信号类型图标
  const getSignalIcon = (type: string) => {
    switch (type) {
      case 'BUY':
        return '📈';
      case 'SELL':
        return '📉';
      default:
        return '⏸️';
    }
  };

  // 获取信号类型文本
  const getSignalText = (type: string) => {
    switch (type) {
      case 'BUY':
        return '买入';
      case 'SELL':
        return '卖出';
      default:
        return '持有';
    }
  };

  // 格式化日期
  const formatDate = (dateString: string) => {
    if (!dateString) return 'N/A';
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
      });
    } catch {
      return dateString;
    }
  };

  // 渲染预测项
  const renderPredictionItem = (prediction: PredictionItem, index: number) => {
    const isExpanded = expandedItems.has(index);
    const signalColor = getSignalColor(prediction.type);
    
    return (
      <Card key={index} style={[styles.predictionCard, { borderLeftColor: signalColor }]}>
        <TouchableOpacity onPress={() => toggleExpanded(index)}>
          <View style={styles.predictionHeader}>
            <View style={styles.headerLeft}>
              <View style={[styles.signalIcon, { backgroundColor: signalColor }]}>
                <Text style={styles.iconText}>{getSignalIcon(prediction.type)}</Text>
              </View>
              <View style={styles.signalInfo}>
                <Text style={[styles.signalType, { color: theme.colors.onSurface }]}>
                  {getSignalText(prediction.type)}信号
                </Text>
                <Text style={[styles.signalDate, { color: theme.colors.onSurfaceVariant }]}>
                  {formatDate(prediction.signal_date)}
                </Text>
              </View>
            </View>
            
            <View style={styles.headerRight}>
              <Text style={[styles.signalPrice, { color: theme.colors.onSurface }]}>
                ¥{prediction.price?.toFixed(2) || 'N/A'}
              </Text>
              <View style={styles.probabilityContainer}>
                <Chip
                  mode="flat"
                  style={[styles.probabilityChip, { backgroundColor: signalColor + '20' }]}
                  textStyle={[styles.probabilityText, { color: signalColor }]}
                  compact
                >
                  {(prediction.probability * 100).toFixed(1)}%
                </Chip>
              </View>
            </View>
            
            <Text style={[styles.expandIcon, { color: theme.colors.onSurfaceVariant }]}>
              {isExpanded ? '▲' : '▼'}
            </Text>
          </View>
        </TouchableOpacity>
        
        {isExpanded && (
          <View style={styles.predictionDetails}>
            <Divider style={styles.detailsDivider} />
            
            {/* 价格信息 */}
            <View style={styles.detailSection}>
              <Text style={[styles.detailTitle, { color: theme.colors.onSurface }]}>
                价格信息
              </Text>
              <View style={styles.priceGrid}>
                <View style={styles.priceItem}>
                  <Text style={[styles.priceLabel, { color: theme.colors.onSurfaceVariant }]}>
                    信号价格
                  </Text>
                  <Text style={[styles.priceValue, { color: theme.colors.onSurface }]}>
                    ¥{prediction.price?.toFixed(2) || 'N/A'}
                  </Text>
                </View>
                {prediction.target_price && (
                  <View style={styles.priceItem}>
                    <Text style={[styles.priceLabel, { color: theme.colors.onSurfaceVariant }]}>
                      目标价格
                    </Text>
                    <Text style={[styles.priceValue, { color: getSignalColor('BUY') }]}>
                      ¥{prediction.target_price.toFixed(2)}
                    </Text>
                  </View>
                )}
                {prediction.stop_loss && (
                  <View style={styles.priceItem}>
                    <Text style={[styles.priceLabel, { color: theme.colors.onSurfaceVariant }]}>
                      止损价格
                    </Text>
                    <Text style={[styles.priceValue, { color: getSignalColor('SELL') }]}>
                      ¥{prediction.stop_loss.toFixed(2)}
                    </Text>
                  </View>
                )}
              </View>
            </View>
            
            {/* 置信度信息 */}
            {prediction.confidence && (
              <View style={styles.detailSection}>
                <Text style={[styles.detailTitle, { color: theme.colors.onSurface }]}>
                  置信度
                </Text>
                <View style={styles.confidenceContainer}>
                  <View style={styles.confidenceBar}>
                    <View 
                      style={[
                        styles.confidenceFill, 
                        { 
                          width: `${prediction.confidence * 100}%`,
                          backgroundColor: signalColor,
                        }
                      ]} 
                    />
                  </View>
                  <Text style={[styles.confidenceText, { color: theme.colors.onSurface }]}>
                    {(prediction.confidence * 100).toFixed(1)}%
                  </Text>
                </View>
              </View>
            )}
            
            {/* 预测理由 */}
            {prediction.reasons && prediction.reasons.length > 0 && (
              <View style={styles.detailSection}>
                <Text style={[styles.detailTitle, { color: theme.colors.onSurface }]}>
                  预测理由
                </Text>
                {prediction.reasons.map((reason, idx) => (
                  <View key={idx} style={styles.reasonItem}>
                    <Text style={[styles.reasonBullet, { color: theme.colors.onSurfaceVariant }]}>
                      •
                    </Text>
                    <Text style={[styles.reasonText, { color: theme.colors.onSurfaceVariant }]}>
                      {reason}
                    </Text>
                  </View>
                ))}
              </View>
            )}
            
            {/* 描述信息 */}
            {prediction.description && (
              <View style={styles.detailSection}>
                <Text style={[styles.detailTitle, { color: theme.colors.onSurface }]}>
                  详细说明
                </Text>
                <Text style={[styles.descriptionText, { color: theme.colors.onSurfaceVariant }]}>
                  {prediction.description}
                </Text>
              </View>
            )}
          </View>
        )}
      </Card>
    );
  };

  if (loading) {
    return (
      <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
        <View style={styles.header}>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            买卖预测信号
          </Text>
          <Text style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
            {stockName} ({stockCode})
          </Text>
        </View>
        <SkeletonCard />
        <View style={styles.predictionsContainer}>
          <SkeletonPredictionCard />
          <SkeletonPredictionCard />
          <SkeletonPredictionCard />
          <SkeletonPredictionCard />
        </View>
      </ScrollView>
    );
  }

  if (error) {
    return (
      <View style={styles.centerContainer}>
        <Text style={[styles.errorText, { color: theme.colors.error }]}>
          {error}
        </Text>
        <Button mode="outlined" onPress={loadPredictions} style={styles.retryButton}>
          重试
        </Button>
      </View>
    );
  }

  if (!predictions || !predictions.predictions || predictions.predictions.length === 0) {
    return (
      <View style={styles.centerContainer}>
        <Text style={[styles.emptyText, { color: theme.colors.onSurfaceVariant }]}>
          暂无预测信号
        </Text>
        <Button mode="outlined" onPress={loadPredictions} style={styles.retryButton}>
          刷新
        </Button>
      </View>
    );
  }

  return (
    <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
      <View style={styles.header}>
        <Text style={[styles.title, { color: theme.colors.onSurface }]}>
          买卖预测信号
        </Text>
        <Text style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
          {stockName} ({stockCode})
        </Text>
      </View>

      {/* 摘要信息 */}
      {predictions.summary && (
        <Card style={styles.summaryCard}>
          <Card.Content>
            <Text style={[styles.summaryTitle, { color: theme.colors.onSurface }]}>
              信号摘要
            </Text>
            <View style={styles.summaryGrid}>
              <View style={styles.summaryItem}>
                <Text style={[styles.summaryLabel, { color: theme.colors.onSurfaceVariant }]}>
                  总信号数
                </Text>
                <Text style={[styles.summaryValue, { color: theme.colors.onSurface }]}>
                  {predictions.summary.total_signals}
                </Text>
              </View>
              <View style={styles.summaryItem}>
                <Text style={[styles.summaryLabel, { color: theme.colors.onSurfaceVariant }]}>
                  买入信号
                </Text>
                <Text style={[styles.summaryValue, { color: getSignalColor('BUY') }]}>
                  {predictions.summary.buy_signals}
                </Text>
              </View>
              <View style={styles.summaryItem}>
                <Text style={[styles.summaryLabel, { color: theme.colors.onSurfaceVariant }]}>
                  卖出信号
                </Text>
                <Text style={[styles.summaryValue, { color: getSignalColor('SELL') }]}>
                  {predictions.summary.sell_signals}
                </Text>
              </View>
              <View style={styles.summaryItem}>
                <Text style={[styles.summaryLabel, { color: theme.colors.onSurfaceVariant }]}>
                  平均置信度
                </Text>
                <Text style={[styles.summaryValue, { color: theme.colors.onSurface }]}>
                  {(predictions.summary.avg_confidence * 100).toFixed(1)}%
                </Text>
              </View>
            </View>
          </Card.Content>
        </Card>
      )}

      {/* 预测信号列表 */}
      <View style={styles.predictionsContainer}>
        {predictions.predictions.map((prediction, index) => 
          renderPredictionItem(prediction, index)
        )}
      </View>

      <View style={styles.refreshContainer}>
        <Button mode="outlined" onPress={loadPredictions} disabled={loading}>
          刷新数据
        </Button>
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  header: {
    padding: 16,
    paddingBottom: 8,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  subtitle: {
    fontSize: 14,
  },
  summaryCard: {
    margin: 16,
    marginTop: 8,
  },
  summaryTitle: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 12,
  },
  summaryGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 16,
  },
  summaryItem: {
    minWidth: '45%',
    alignItems: 'center',
  },
  summaryLabel: {
    fontSize: 12,
    marginBottom: 4,
  },
  summaryValue: {
    fontSize: 18,
    fontWeight: '600',
  },
  predictionsContainer: {
    padding: 16,
    paddingTop: 8,
    gap: 12,
  },
  predictionCard: {
    borderLeftWidth: 4,
    marginBottom: 8,
  },
  predictionHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 16,
  },
  headerLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  signalIcon: {
    width: 36,
    height: 36,
    borderRadius: 18,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  iconText: {
    fontSize: 16,
  },
  signalInfo: {
    flex: 1,
  },
  signalType: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 2,
  },
  signalDate: {
    fontSize: 12,
  },
  headerRight: {
    alignItems: 'flex-end',
    marginRight: 8,
  },
  signalPrice: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  probabilityContainer: {
    alignItems: 'flex-end',
  },
  probabilityChip: {
    height: 24,
  },
  probabilityText: {
    fontSize: 11,
    fontWeight: '600',
  },
  expandIcon: {
    fontSize: 12,
    marginLeft: 8,
  },
  predictionDetails: {
    paddingHorizontal: 16,
    paddingBottom: 16,
  },
  detailsDivider: {
    marginBottom: 16,
  },
  detailSection: {
    marginBottom: 16,
  },
  detailTitle: {
    fontSize: 14,
    fontWeight: '600',
    marginBottom: 8,
  },
  priceGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 16,
  },
  priceItem: {
    minWidth: '30%',
  },
  priceLabel: {
    fontSize: 12,
    marginBottom: 4,
  },
  priceValue: {
    fontSize: 16,
    fontWeight: '600',
  },
  confidenceContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  confidenceBar: {
    flex: 1,
    height: 6,
    backgroundColor: '#e0e0e0',
    borderRadius: 3,
    overflow: 'hidden',
  },
  confidenceFill: {
    height: '100%',
    borderRadius: 3,
  },
  confidenceText: {
    fontSize: 12,
    fontWeight: '600',
    minWidth: 40,
  },
  reasonItem: {
    flexDirection: 'row',
    marginBottom: 4,
  },
  reasonBullet: {
    fontSize: 12,
    marginRight: 8,
    marginTop: 2,
  },
  reasonText: {
    fontSize: 12,
    flex: 1,
    lineHeight: 16,
  },
  descriptionText: {
    fontSize: 13,
    lineHeight: 18,
  },
  loadingText: {
    marginTop: 16,
    fontSize: 14,
  },
  errorText: {
    fontSize: 14,
    textAlign: 'center',
    marginBottom: 16,
  },
  emptyText: {
    fontSize: 14,
    textAlign: 'center',
    marginBottom: 16,
  },
  retryButton: {
    marginTop: 8,
  },
  refreshContainer: {
    padding: 16,
    alignItems: 'center',
  },
});

export default PredictionSignals;
