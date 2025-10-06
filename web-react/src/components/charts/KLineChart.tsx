/**
 * K线图组件 - 使用 Recharts 绘制真正的蜡烛图
 */

import React, { useMemo } from 'react';
import {
  ComposedChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Bar,
  Cell,
  Customized,
} from 'recharts';
import { Box, Typography, useTheme } from '@mui/material';
import { StockDaily } from '../../types/stock';

interface KLineChartProps {
  data: StockDaily[];
}

// 自定义蜡烛图层 - 作为 Recharts 的自定义元素
const CandlestickLayer: React.FC<any> = ({ xAxisMap, yAxisMap, data, offset }) => {
  if (!xAxisMap || !yAxisMap || !data) return null;
  
  const xAxis = xAxisMap[0];
  const yAxis = yAxisMap[0];
  
  if (!xAxis || !yAxis) return null;
  
  const { scale: xScale } = xAxis;
  const { scale: yScale } = yAxis;
  
  // 计算每个蜡烛的宽度
  const dataKeys = data.map((item: any) => item.date);
  const bandwidth = xScale.bandwidth ? xScale.bandwidth() : 
    (dataKeys.length > 1 ? Math.abs(xScale(dataKeys[1]) - xScale(dataKeys[0])) * 0.8 : 20);
  
  return (
    <g className="recharts-candlestick-layer">
      {data.map((item: any, index: number) => {
        const { open, close, high, low, isUp, date } = item;
        
        if (!open || !close || !high || !low) return null;
        
        const color = isUp ? '#26a69a' : '#ef5350';
        
        // 使用 yScale 将价格转换为 Y 坐标
        const highY = yScale(high);
        const lowY = yScale(low);
        const openY = yScale(open);
        const closeY = yScale(close);
        
        // 使用 xScale 获取 X 坐标
        const x = xScale(date);
        if (x === undefined) return null;
        
        const centerX = x + bandwidth / 2;
        
        // 蜡烛实体
        const bodyTop = Math.min(openY, closeY);
        let bodyHeight = Math.abs(closeY - openY);
        
        // 十字星的最小高度
        if (bodyHeight < 2) bodyHeight = 2;
        
        // 计算宽度
        const wickWidth = Math.max(1, bandwidth * 0.15);
        const bodyWidth = Math.max(3, bandwidth * 0.7);
        
        return (
          <g key={`candle-${date}-${index}`}>
            {/* 上下影线 */}
            <line
              x1={centerX}
              y1={highY}
              x2={centerX}
              y2={lowY}
              stroke={color}
              strokeWidth={wickWidth}
            />
            {/* 蜡烛实体 */}
            <rect
              x={centerX - bodyWidth / 2}
              y={bodyTop}
              width={bodyWidth}
              height={bodyHeight}
              fill={color}
              stroke={color}
              strokeWidth={0.5}
            />
          </g>
        );
      })}
    </g>
  );
};

const KLineChart: React.FC<KLineChartProps> = ({ data }) => {
  const theme = useTheme();

  // 转换数据格式
  const chartData = useMemo(() => {
    if (!data || data.length === 0) return [];

    return data.map((item) => {
      const open = parseFloat(item.open);
      const close = parseFloat(item.close);
      const high = parseFloat(item.high);
      const low = parseFloat(item.low);
      const vol = parseFloat(item.vol);

      // 计算K线实体的起点和高度
      const bodyTop = Math.max(open, close);
      const bodyBottom = Math.min(open, close);
      const bodyHeight = bodyTop - bodyBottom;

      return {
        date: item.trade_date,
        open,
        close,
        high,
        low,
        vol,
        isUp: close >= open,
        // 用于占位的 Bar（不可见）
        placeholder: 0,
      };
    });
  }, [data]);

  // 格式化日期
  const formatDate = (dateStr: string) => {
    if (!dateStr || dateStr.length !== 8) return dateStr;
    return `${dateStr.substring(4, 6)}/${dateStr.substring(6, 8)}`;
  };

  // 格式化价格
  const formatPrice = (value: number) => {
    return `¥${value.toFixed(2)}`;
  };

  // 格式化成交量
  const formatVolume = (value: number) => {
    if (value >= 100000000) {
      return `${(value / 100000000).toFixed(2)}亿`;
    } else if (value >= 10000) {
      return `${(value / 10000).toFixed(2)}万`;
    }
    return value.toString();
  };

  // 自定义Tooltip
  const CustomTooltip = ({ active, payload }: any) => {
    if (!active || !payload || payload.length === 0) return null;

    const data = payload[0].payload;

    return (
      <Box
        sx={{
          bgcolor: 'background.paper',
          p: 1.5,
          border: 1,
          borderColor: 'divider',
          borderRadius: 1,
          boxShadow: 2,
        }}
      >
        <Typography variant="body2" gutterBottom>
          日期: {data.date}
        </Typography>
        <Typography variant="body2" color={data.isUp ? 'success.main' : 'error.main'}>
          开: {formatPrice(data.open)}
        </Typography>
        <Typography variant="body2" color={data.isUp ? 'success.main' : 'error.main'}>
          高: {formatPrice(data.high)}
        </Typography>
        <Typography variant="body2" color={data.isUp ? 'success.main' : 'error.main'}>
          低: {formatPrice(data.low)}
        </Typography>
        <Typography variant="body2" color={data.isUp ? 'success.main' : 'error.main'}>
          收: {formatPrice(data.close)}
        </Typography>
        <Typography variant="body2">
          量: {formatVolume(data.vol)}
        </Typography>
      </Box>
    );
  };

  if (!chartData || chartData.length === 0) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight={400}>
        <Typography variant="body1" color="text.secondary">
          暂无数据
        </Typography>
      </Box>
    );
  }

  return (
    <Box>
      {/* K线蜡烛图 */}
      <ResponsiveContainer width="100%" height={400}>
        <ComposedChart data={chartData} margin={{ top: 20, right: 30, left: 20, bottom: 0 }}>
          <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.divider} />
          <XAxis
            dataKey="date"
            tickFormatter={formatDate}
            stroke={theme.palette.text.secondary}
          />
          <YAxis
            domain={['auto', 'auto']}
            tickFormatter={formatPrice}
            stroke={theme.palette.text.secondary}
            width={80}
          />
          <Tooltip content={<CustomTooltip />} />
          <Legend />

          {/* 自定义蜡烛图层 - 使用 Customized 组件 */}
          <Customized component={<CandlestickLayer data={chartData} />} />

          {/* MA5 */}
          <Line
            type="monotone"
            dataKey={(data) => {
              const index = chartData.indexOf(data);
              if (index < 4) return null;
              const sum = chartData.slice(index - 4, index + 1).reduce((acc, item) => acc + item.close, 0);
              return sum / 5;
            }}
            stroke="#FFA726"
            strokeWidth={1.5}
            dot={false}
            name="MA5"
            connectNulls
            isAnimationActive={false}
          />

          {/* MA10 */}
          <Line
            type="monotone"
            dataKey={(data) => {
              const index = chartData.indexOf(data);
              if (index < 9) return null;
              const sum = chartData.slice(index - 9, index + 1).reduce((acc, item) => acc + item.close, 0);
              return sum / 10;
            }}
            stroke="#42A5F5"
            strokeWidth={1.5}
            dot={false}
            name="MA10"
            connectNulls
            isAnimationActive={false}
          />

          {/* MA20 */}
          <Line
            type="monotone"
            dataKey={(data) => {
              const index = chartData.indexOf(data);
              if (index < 19) return null;
              const sum = chartData.slice(index - 19, index + 1).reduce((acc, item) => acc + item.close, 0);
              return sum / 20;
            }}
            stroke="#AB47BC"
            strokeWidth={1.5}
            dot={false}
            name="MA20"
            connectNulls
            isAnimationActive={false}
          />
        </ComposedChart>
      </ResponsiveContainer>

      {/* 成交量图 */}
      <ResponsiveContainer width="100%" height={150}>
        <ComposedChart data={chartData} margin={{ top: 10, right: 30, left: 20, bottom: 0 }}>
          <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.divider} />
          <XAxis
            dataKey="date"
            tickFormatter={formatDate}
            stroke={theme.palette.text.secondary}
          />
          <YAxis
            tickFormatter={formatVolume}
            stroke={theme.palette.text.secondary}
            width={80}
          />
          <Tooltip
            formatter={(value: number) => formatVolume(value)}
            labelFormatter={(label) => `日期: ${label}`}
          />
          <Bar
            dataKey="vol"
            name="成交量"
          >
            {chartData.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={entry.isUp ? '#26a69a' : '#ef5350'} />
            ))}
          </Bar>
        </ComposedChart>
      </ResponsiveContainer>
    </Box>
  );
};

export default KLineChart;

