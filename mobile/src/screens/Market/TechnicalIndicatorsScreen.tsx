/**
 * 技术指标页面
 */

import React from 'react';
import { View, StyleSheet } from 'react-native';
import { RouteProp, useRoute } from '@react-navigation/native';

import { MarketStackParamList } from '@/navigation/AppNavigator';
import TechnicalIndicators from '@/components/TechnicalIndicators';

type TechnicalIndicatorsRouteProp = RouteProp<MarketStackParamList, 'TechnicalIndicators'>;

const TechnicalIndicatorsScreen: React.FC = () => {
  const route = useRoute<TechnicalIndicatorsRouteProp>();
  const { stockCode, stockName } = route.params;

  return (
    <View style={styles.container}>
      <TechnicalIndicators stockCode={stockCode} stockName={stockName} />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
});

export default TechnicalIndicatorsScreen;
