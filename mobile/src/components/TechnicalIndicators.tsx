/**
 * 技术指标展示组件
 */

import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, ScrollView, ActivityIndicator } from 'react-native';
import { Card, useTheme, Chip, Divider, Button } from 'react-native-paper';
import { apiService } from '@/services/apiService';
import { SkeletonIndicatorCard } from '@/components/common/SkeletonLoader';

interface TechnicalIndicatorsProps {
  stockCode: string;
  stockName: string;
}

interface IndicatorData {
  // MACD指标
  macd?: {
    dif: number;
    dea: number;
    histogram: number;
    signal?: string;
  };
  // RSI指标
  rsi?: {
    rsi14: number;
    signal?: string;
  };
  // 布林带
  boll?: {
    upper: number;
    middle: number;
    lower: number;
    signal?: string;
  };
  // 移动平均线
  ma?: {
    ma5: number;
    ma10: number;
    ma20: number;
    ma60?: number;
    ma120?: number;
    signal?: string;
  };
  // KDJ指标
  kdj?: {
    k: number;
    d: number;
    j: number;
    signal?: string;
  };
}

const TechnicalIndicators: React.FC<TechnicalIndicatorsProps> = ({ stockCode, stockName }) => {
  const theme = useTheme();
  const [loading, setLoading] = useState(false);
  const [indicators, setIndicators] = useState<IndicatorData | null>(null);
  const [error, setError] = useState<string | null>(null);

  // 加载技术指标数据
  const loadIndicators = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await apiService.getIndicators(stockCode);
      
      if (response.success && response.data) {
        setIndicators(response.data);
      } else {
        setError(response.error || '加载技术指标失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载技术指标失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadIndicators();
  }, [stockCode]);

  // 获取信号颜色
  const getSignalColor = (signal?: string) => {
    switch (signal?.toLowerCase()) {
      case 'buy':
      case 'bullish':
        return theme.colors.primary;
      case 'sell':
      case 'bearish':
        return theme.colors.error;
      default:
        return theme.colors.outline;
    }
  };

  // 获取信号文本
  const getSignalText = (signal?: string) => {
    switch (signal?.toLowerCase()) {
      case 'buy':
        return '买入';
      case 'sell':
        return '卖出';
      case 'bullish':
        return '看涨';
      case 'bearish':
        return '看跌';
      case 'hold':
        return '持有';
      default:
        return '中性';
    }
  };

  // 渲染指标项
  const renderIndicatorItem = (
    title: string,
    description: string,
    values: Array<{ name: string; value: string | number; unit?: string }>,
    signal?: string
  ) => (
    <Card style={styles.indicatorCard} key={title}>
      <Card.Content>
        <View style={styles.indicatorHeader}>
          <Text style={[styles.indicatorTitle, { color: theme.colors.onSurface }]}>
            {title}
          </Text>
          {signal && (
            <Chip
              mode="flat"
              style={[styles.signalChip, { backgroundColor: getSignalColor(signal) + '20' }]}
              textStyle={[styles.signalText, { color: getSignalColor(signal) }]}
              compact
            >
              {getSignalText(signal)}
            </Chip>
          )}
        </View>
        
        <Text style={[styles.indicatorDescription, { color: theme.colors.onSurfaceVariant }]}>
          {description}
        </Text>
        
        <Divider style={styles.divider} />
        
        <View style={styles.valuesGrid}>
          {values.map((item, index) => (
            <View key={index} style={styles.valueItem}>
              <Text style={[styles.valueName, { color: theme.colors.onSurfaceVariant }]}>
                {item.name}
              </Text>
              <Text style={[styles.valueNumber, { color: theme.colors.onSurface }]}>
                {typeof item.value === 'number' ? item.value.toFixed(4) : item.value}
                {item.unit && <Text style={styles.valueUnit}>{item.unit}</Text>}
              </Text>
            </View>
          ))}
        </View>
      </Card.Content>
    </Card>
  );

  if (loading) {
    return (
      <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
        <View style={styles.header}>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            技术指标分析
          </Text>
          <Text style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
            {stockName} ({stockCode})
          </Text>
        </View>
        <View style={styles.indicatorsContainer}>
          <SkeletonIndicatorCard />
          <SkeletonIndicatorCard />
          <SkeletonIndicatorCard />
          <SkeletonIndicatorCard />
          <SkeletonIndicatorCard />
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
        <Button mode="outlined" onPress={loadIndicators} style={styles.retryButton}>
          重试
        </Button>
      </View>
    );
  }

  if (!indicators) {
    return (
      <View style={styles.centerContainer}>
        <Text style={[styles.emptyText, { color: theme.colors.onSurfaceVariant }]}>
          暂无技术指标数据
        </Text>
        <Button mode="outlined" onPress={loadIndicators} style={styles.retryButton}>
          刷新
        </Button>
      </View>
    );
  }

  return (
    <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
      <View style={styles.header}>
        <Text style={[styles.title, { color: theme.colors.onSurface }]}>
          技术指标分析
        </Text>
        <Text style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
          {stockName} ({stockCode})
        </Text>
      </View>

      <View style={styles.indicatorsContainer}>
        {/* MACD指标 */}
        {indicators.macd && renderIndicatorItem(
          'MACD',
          'MACD是趋势跟踪指标，通过快慢均线的差值判断买卖时机。DIF上穿DEA为金叉买入信号，下穿为死叉卖出信号。',
          [
            { name: 'DIF', value: indicators.macd.dif },
            { name: 'DEA', value: indicators.macd.dea },
            { name: '柱状图', value: indicators.macd.histogram },
          ],
          indicators.macd.signal
        )}

        {/* RSI指标 */}
        {indicators.rsi && renderIndicatorItem(
          'RSI',
          'RSI相对强弱指数衡量价格变动的速度和幅度。RSI>70为超买区域，RSI<30为超卖区域，可作为反转信号参考。',
          [
            { name: 'RSI14', value: indicators.rsi.rsi14, unit: '%' },
          ],
          indicators.rsi.signal
        )}

        {/* 布林带 */}
        {indicators.boll && renderIndicatorItem(
          '布林带',
          '布林带由移动平均线和标准差构成，价格触及上轨可能回调，触及下轨可能反弹。带宽收窄预示突破，扩张表示趋势延续。',
          [
            { name: '上轨', value: indicators.boll.upper, unit: '¥' },
            { name: '中轨', value: indicators.boll.middle, unit: '¥' },
            { name: '下轨', value: indicators.boll.lower, unit: '¥' },
          ],
          indicators.boll.signal
        )}

        {/* 移动平均线 */}
        {indicators.ma && renderIndicatorItem(
          '移动平均线',
          '移动平均线用于平滑价格波动，识别趋势方向。金叉死叉信号可用于判断买卖时机，支撑阻力位可用于风险控制。',
          [
            { name: 'MA5', value: indicators.ma.ma5, unit: '¥' },
            { name: 'MA10', value: indicators.ma.ma10, unit: '¥' },
            { name: 'MA20', value: indicators.ma.ma20, unit: '¥' },
            ...(indicators.ma.ma60 ? [{ name: 'MA60', value: indicators.ma.ma60, unit: '¥' }] : []),
            ...(indicators.ma.ma120 ? [{ name: 'MA120', value: indicators.ma.ma120, unit: '¥' }] : []),
          ],
          indicators.ma.signal
        )}

        {/* KDJ指标 */}
        {indicators.kdj && renderIndicatorItem(
          'KDJ',
          'KDJ随机指标通过收盘价在一定周期内的相对位置判断超买超卖。K值大于D值且J值上升为买入信号，反之为卖出信号。',
          [
            { name: 'K值', value: indicators.kdj.k, unit: '%' },
            { name: 'D值', value: indicators.kdj.d, unit: '%' },
            { name: 'J值', value: indicators.kdj.j, unit: '%' },
          ],
          indicators.kdj.signal
        )}
      </View>

      <View style={styles.refreshContainer}>
        <Button mode="outlined" onPress={loadIndicators} disabled={loading}>
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
  indicatorsContainer: {
    padding: 16,
    paddingTop: 8,
    gap: 16,
  },
  indicatorCard: {
    marginBottom: 8,
  },
  indicatorHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  indicatorTitle: {
    fontSize: 18,
    fontWeight: '600',
    flex: 1,
  },
  signalChip: {
    height: 28,
  },
  signalText: {
    fontSize: 12,
    fontWeight: '600',
  },
  indicatorDescription: {
    fontSize: 13,
    lineHeight: 18,
    marginBottom: 12,
  },
  divider: {
    marginVertical: 8,
  },
  valuesGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 12,
  },
  valueItem: {
    minWidth: '30%',
    marginBottom: 8,
  },
  valueName: {
    fontSize: 12,
    marginBottom: 2,
  },
  valueNumber: {
    fontSize: 16,
    fontWeight: '600',
  },
  valueUnit: {
    fontSize: 12,
    fontWeight: '400',
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

export default TechnicalIndicators;
