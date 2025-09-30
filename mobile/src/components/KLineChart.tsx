/**
 * K线图表组件
 * 使用React Native ECharts渲染K线图（蜡烛图）
 */

import React, { useRef, useEffect } from 'react';
import { View, StyleSheet, Dimensions } from 'react-native';
import { Text, ActivityIndicator, useTheme } from 'react-native-paper';
import { WebView } from 'react-native-webview';
import * as echarts from 'echarts/core';
import { CandlestickChart, LineChart, BarChart } from 'echarts/charts';
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  DataZoomComponent,
} from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';

// 注册ECharts组件
echarts.use([
  CandlestickChart,
  LineChart,
  BarChart,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  DataZoomComponent,
  CanvasRenderer,
]);

const { width: SCREEN_WIDTH } = Dimensions.get('window');

interface DailyData {
  trade_date: string;
  open: number | string;
  high: number | string;
  low: number | string;
  close: number | string;
  vol: number | string;
}

interface KLineChartProps {
  data: DailyData[];
  stockCode?: string;
  stockName?: string;
  loading?: boolean;
}

/**
 * 计算移动平均线
 */
const calculateMA = (data: number[], dayCount: number): (number | null)[] => {
  const result: (number | null)[] = [];
  
  for (let i = 0; i < data.length; i++) {
    if (i < dayCount - 1) {
      result.push(null);
      continue;
    }
    
    let sum = 0;
    for (let j = 0; j < dayCount; j++) {
      sum += data[i - j];
    }
    result.push(sum / dayCount);
  }
  
  return result;
};

/**
 * 格式化日期
 */
const formatDate = (dateStr: string): string => {
  if (dateStr && dateStr.includes('T')) {
    return dateStr.split('T')[0];
  } else if (dateStr && dateStr.length === 8) {
    return `${dateStr.substring(0, 4)}-${dateStr.substring(4, 6)}-${dateStr.substring(6, 8)}`;
  }
  return dateStr;
};

const KLineChart: React.FC<KLineChartProps> = ({ 
  data, 
  stockCode = '', 
  stockName = '',
  loading = false 
}) => {
  const theme = useTheme();
  const webViewRef = useRef<WebView>(null);

  useEffect(() => {
    if (data && data.length > 0 && webViewRef.current) {
      renderChart();
    }
  }, [data, theme]);

  const renderChart = () => {
    if (!data || data.length === 0) return;

    // 准备K线数据 [开盘, 收盘, 最低, 最高]
    const klineData = data.map(item => [
      parseFloat(String(item.open)),
      parseFloat(String(item.close)),
      parseFloat(String(item.low)),
      parseFloat(String(item.high)),
    ]);

    // 准备成交量数据
    const volumeData = data.map(item => parseFloat(String(item.vol)));

    // 准备日期标签
    const dates = data.map(item => formatDate(item.trade_date));

    // 计算移动平均线
    const closeData = data.map(item => parseFloat(String(item.close)));
    const ma5 = calculateMA(closeData, 5);
    const ma10 = calculateMA(closeData, 10);
    const ma20 = calculateMA(closeData, 20);

    const option = {
      backgroundColor: theme.dark ? '#1a1a1a' : '#ffffff',
      animation: true,
      legend: {
        data: ['日K', 'MA5', 'MA10', 'MA20'],
        top: 10,
        textStyle: {
          color: theme.dark ? '#ffffff' : '#333333',
        },
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
        },
        backgroundColor: theme.dark ? 'rgba(50, 50, 50, 0.9)' : 'rgba(255, 255, 255, 0.9)',
        borderColor: theme.colors.outline,
        textStyle: {
          color: theme.dark ? '#ffffff' : '#333333',
        },
        formatter: (params: any) => {
          const param = params[0];
          const klineParam = params.find((p: any) => p.seriesName === '日K');
          
          if (!klineParam) return '';
          
          const value = klineParam.data;
          const date = param.name;
          
          return `
            <div style="padding: 8px;">
              <div style="font-weight: bold; margin-bottom: 4px;">${date}</div>
              <div>开盘: ${value[1].toFixed(2)}</div>
              <div>收盘: ${value[2].toFixed(2)}</div>
              <div>最低: ${value[3].toFixed(2)}</div>
              <div>最高: ${value[4].toFixed(2)}</div>
              ${params.find((p: any) => p.seriesName === 'MA5') 
                ? `<div>MA5: ${params.find((p: any) => p.seriesName === 'MA5').data.toFixed(2)}</div>` 
                : ''}
              ${params.find((p: any) => p.seriesName === 'MA10') 
                ? `<div>MA10: ${params.find((p: any) => p.seriesName === 'MA10').data.toFixed(2)}</div>` 
                : ''}
              ${params.find((p: any) => p.seriesName === 'MA20') 
                ? `<div>MA20: ${params.find((p: any) => p.seriesName === 'MA20').data.toFixed(2)}</div>` 
                : ''}
            </div>
          `;
        },
      },
      grid: [
        {
          left: '10%',
          right: '8%',
          top: '15%',
          height: '50%',
        },
        {
          left: '10%',
          right: '8%',
          top: '70%',
          height: '15%',
        },
      ],
      xAxis: [
        {
          type: 'category',
          data: dates,
          scale: true,
          boundaryGap: false,
          axisLine: { 
            onZero: false,
            lineStyle: { color: theme.colors.outline }
          },
          splitLine: { show: false },
          axisLabel: {
            color: theme.dark ? '#999999' : '#666666',
          },
        },
        {
          type: 'category',
          gridIndex: 1,
          data: dates,
          scale: true,
          boundaryGap: false,
          axisLine: { 
            onZero: false,
            lineStyle: { color: theme.colors.outline }
          },
          axisTick: { show: false },
          splitLine: { show: false },
          axisLabel: { show: false },
        },
      ],
      yAxis: [
        {
          scale: true,
          splitArea: {
            show: true,
          },
          axisLabel: {
            color: theme.dark ? '#999999' : '#666666',
          },
          splitLine: {
            lineStyle: { color: theme.colors.outlineVariant }
          },
        },
        {
          scale: true,
          gridIndex: 1,
          splitNumber: 2,
          axisLabel: { show: false },
          axisLine: { show: false },
          axisTick: { show: false },
          splitLine: { show: false },
        },
      ],
      dataZoom: [
        {
          type: 'inside',
          xAxisIndex: [0, 1],
          start: 70,
          end: 100,
        },
        {
          show: true,
          xAxisIndex: [0, 1],
          type: 'slider',
          top: '90%',
          start: 70,
          end: 100,
          backgroundColor: theme.dark ? '#2a2a2a' : '#f0f0f0',
          borderColor: theme.colors.outline,
          fillerColor: 'rgba(47, 69, 84, 0.25)',
          handleStyle: {
            color: theme.colors.primary,
          },
          textStyle: {
            color: theme.dark ? '#999999' : '#666666',
          },
        },
      ],
      series: [
        {
          name: '日K',
          type: 'candlestick',
          data: klineData,
          itemStyle: {
            color: '#ef5350',
            color0: '#26a69a',
            borderColor: '#ef5350',
            borderColor0: '#26a69a',
          },
        },
        {
          name: 'MA5',
          type: 'line',
          data: ma5,
          smooth: true,
          lineStyle: {
            opacity: 0.8,
            width: 1.5,
            color: '#FF6B6B',
          },
          showSymbol: false,
        },
        {
          name: 'MA10',
          type: 'line',
          data: ma10,
          smooth: true,
          lineStyle: {
            opacity: 0.8,
            width: 1.5,
            color: '#4ECDC4',
          },
          showSymbol: false,
        },
        {
          name: 'MA20',
          type: 'line',
          data: ma20,
          smooth: true,
          lineStyle: {
            opacity: 0.8,
            width: 1.5,
            color: '#95E1D3',
          },
          showSymbol: false,
        },
        {
          name: '成交量',
          type: 'bar',
          xAxisIndex: 1,
          yAxisIndex: 1,
          data: volumeData,
          itemStyle: {
            color: (params: any) => {
              const index = params.dataIndex;
              const current = klineData[index];
              // 收盘价 > 开盘价，涨，红色；否则跌，绿色
              return current[2] >= current[1] ? '#ef5350' : '#26a69a';
            },
          },
        },
      ],
    };

    // 生成HTML内容，嵌入ECharts
    const html = `
      <!DOCTYPE html>
      <html>
        <head>
          <meta charset="utf-8">
          <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
          <style>
            * { margin: 0; padding: 0; }
            body { 
              background-color: ${theme.dark ? '#1a1a1a' : '#ffffff'}; 
            }
            #chart { 
              width: 100vw; 
              height: 100vh; 
            }
          </style>
        </head>
        <body>
          <div id="chart"></div>
          <script src="https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js"></script>
          <script>
            var chart = echarts.init(document.getElementById('chart'));
            var option = ${JSON.stringify(option)};
            chart.setOption(option);
            
            window.addEventListener('resize', function() {
              chart.resize();
            });
          </script>
        </body>
      </html>
    `;

    // 注入HTML到WebView
    if (webViewRef.current) {
      webViewRef.current.injectJavaScript(`
        document.open();
        document.write(\`${html.replace(/`/g, '\\`')}\`);
        document.close();
      `);
    }
  };

  if (loading) {
    return (
      <View style={[styles.container, styles.centerContent]}>
        <ActivityIndicator size="large" color={theme.colors.primary} />
        <Text style={styles.loadingText}>加载图表数据中...</Text>
      </View>
    );
  }

  if (!data || data.length === 0) {
    return (
      <View style={[styles.container, styles.centerContent]}>
        <Text style={styles.emptyText}>暂无K线数据</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <WebView
        ref={webViewRef}
        originWhitelist={['*']}
        source={{ html: '<html><body></body></html>' }}
        style={styles.webview}
        javaScriptEnabled={true}
        domStorageEnabled={true}
        onLoad={renderChart}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    width: '100%',
    height: 400,
  },
  centerContent: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  webview: {
    flex: 1,
    backgroundColor: 'transparent',
  },
  loadingText: {
    marginTop: 12,
    fontSize: 14,
  },
  emptyText: {
    fontSize: 16,
    opacity: 0.6,
  },
});

export default KLineChart;
