/**
 * 骨架屏加载组件
 */

import React, { useEffect, useRef } from 'react';
import { View, StyleSheet, Animated } from 'react-native';
import { useTheme } from 'react-native-paper';

interface SkeletonLoaderProps {
  width?: number | string;
  height?: number;
  borderRadius?: number;
  style?: any;
}

const SkeletonLoader: React.FC<SkeletonLoaderProps> = ({
  width = '100%',
  height = 20,
  borderRadius = 4,
  style,
}) => {
  const theme = useTheme();
  const animatedValue = useRef(new Animated.Value(0)).current;

  useEffect(() => {
    const animation = Animated.loop(
      Animated.sequence([
        Animated.timing(animatedValue, {
          toValue: 1,
          duration: 1000,
          useNativeDriver: true,
        }),
        Animated.timing(animatedValue, {
          toValue: 0,
          duration: 1000,
          useNativeDriver: true,
        }),
      ])
    );

    animation.start();

    return () => animation.stop();
  }, [animatedValue]);

  const opacity = animatedValue.interpolate({
    inputRange: [0, 1],
    outputRange: [0.3, 0.7],
  });

  return (
    <Animated.View
      style={[
        styles.skeleton,
        {
          width,
          height,
          borderRadius,
          backgroundColor: theme.colors.outline,
          opacity,
        },
        style,
      ]}
    />
  );
};

// 预设的骨架屏组件
export const SkeletonCard: React.FC = () => {
  const theme = useTheme();
  
  return (
    <View style={[styles.card, { backgroundColor: theme.colors.surface }]}>
      <View style={styles.cardHeader}>
        <SkeletonLoader width={120} height={20} />
        <SkeletonLoader width={60} height={16} />
      </View>
      <View style={styles.cardContent}>
        <SkeletonLoader width="100%" height={16} style={styles.skeletonLine} />
        <SkeletonLoader width="80%" height={16} style={styles.skeletonLine} />
        <SkeletonLoader width="60%" height={16} style={styles.skeletonLine} />
      </View>
    </View>
  );
};

export const SkeletonIndicatorCard: React.FC = () => {
  const theme = useTheme();
  
  return (
    <View style={[styles.card, { backgroundColor: theme.colors.surface }]}>
      <View style={styles.indicatorHeader}>
        <SkeletonLoader width={80} height={18} />
        <SkeletonLoader width={50} height={24} borderRadius={12} />
      </View>
      <SkeletonLoader width="100%" height={14} style={styles.skeletonLine} />
      <View style={styles.indicatorValues}>
        <View style={styles.valueItem}>
          <SkeletonLoader width={40} height={12} />
          <SkeletonLoader width={60} height={16} />
        </View>
        <View style={styles.valueItem}>
          <SkeletonLoader width={40} height={12} />
          <SkeletonLoader width={60} height={16} />
        </View>
        <View style={styles.valueItem}>
          <SkeletonLoader width={40} height={12} />
          <SkeletonLoader width={60} height={16} />
        </View>
      </View>
    </View>
  );
};

export const SkeletonPredictionCard: React.FC = () => {
  const theme = useTheme();
  
  return (
    <View style={[styles.card, { backgroundColor: theme.colors.surface }]}>
      <View style={styles.predictionHeader}>
        <View style={styles.predictionLeft}>
          <SkeletonLoader width={36} height={36} borderRadius={18} />
          <View style={styles.predictionInfo}>
            <SkeletonLoader width={60} height={16} />
            <SkeletonLoader width={80} height={12} />
          </View>
        </View>
        <View style={styles.predictionRight}>
          <SkeletonLoader width={80} height={18} />
          <SkeletonLoader width={50} height={20} borderRadius={10} />
        </View>
      </View>
    </View>
  );
};

export const SkeletonChart: React.FC<{ height?: number }> = ({ height = 300 }) => {
  const theme = useTheme();
  
  return (
    <View style={[styles.chartContainer, { backgroundColor: theme.colors.surface, height }]}>
      <View style={styles.chartHeader}>
        <SkeletonLoader width={100} height={18} />
      </View>
      <View style={styles.chartArea}>
        <SkeletonLoader width="100%" height={height - 60} borderRadius={8} />
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  skeleton: {
    // Base skeleton styles handled by animated view
  },
  card: {
    borderRadius: 12,
    padding: 16,
    marginBottom: 16,
    elevation: 2,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  cardContent: {
    gap: 8,
  },
  skeletonLine: {
    marginBottom: 8,
  },
  indicatorHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  indicatorValues: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 16,
    marginTop: 12,
  },
  valueItem: {
    minWidth: '30%',
    gap: 4,
  },
  predictionHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  predictionLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  predictionInfo: {
    marginLeft: 12,
    gap: 4,
  },
  predictionRight: {
    alignItems: 'flex-end',
    gap: 4,
  },
  chartContainer: {
    borderRadius: 12,
    padding: 16,
    marginBottom: 16,
    elevation: 2,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
  },
  chartHeader: {
    marginBottom: 16,
  },
  chartArea: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
});

export default SkeletonLoader;
