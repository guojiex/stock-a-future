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
} from 'recharts';
import { Box, Typography, useTheme } from '@mui/material';
import { StockDaily } from '../../types/stock';

interface KLineChartProps {
  data: StockDaily[];
}

// 自定义蜡烛图形状组件
const CandleShape = (props: any) => {
  const { x, y, width, height, payload } = props;
  
  if (!payload || !payload.open || !payload.close || !payload.high || !payload.low) {
    return null;
  }

  const { open, close, high, low, isUp } = payload;
  const color = isUp ? '#26a69a' : '#ef5350';
  
  // 计算Y坐标的比例（基于图表的高度和数据范围）
  const { chartHeight = 314, yMin = 22.5, yMax = 24.5 } = props;
  const priceRange = yMax - yMin;
  const yScale = chartHeight / priceRange;
  
  // 计算各个价格点的Y坐标（相对于图表顶部）
  const highY = (yMax - high) * yScale;
  const lowY = (yMax - low) * yScale;
  const openY = (yMax - open) * yScale;
  const closeY = (yMax - close) * yScale;
  
  // 蜡烛实体的top和height
  const bodyTop = Math.min(openY, closeY);
  let bodyHeight = Math.abs(closeY - openY);
  
  // 十字星的最小高度
  if (bodyHeight < 2) bodyHeight = 2;
  
  // 计算宽度
  const wickWidth = Math.max(1, width * 0.15);
  const bodyWidth = Math.max(3, width * 0.7);
  const centerX = x + width / 2;
  
  return (
    <g>
      {/* 上下影线 */}
      <line
        x1={centerX}
        y1={highY + y}
        x2={centerX}
        y2={lowY + y}
        stroke={color}
        strokeWidth={wickWidth}
      />
      {/* 蜡烛实体 */}
      <rect
        x={centerX - bodyWidth / 2}
        y={bodyTop + y}
        width={bodyWidth}
        height={bodyHeight}
        fill={color}
        stroke={color}
        strokeWidth={0.5}
      />
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

      return {
        date: item.trade_date,
        open,
        close,
        high,
        low,
        vol,
        isUp: close >= open,
        // 用于蜡烛图的占位值（使用low作为基准）
        candleBase: low,
        candleHeight: high - low,
      };
    });
  }, [data]);

  // 计算价格范围（用于蜡烛图渲染）
  const priceRange = useMemo(() => {
    if (chartData.length === 0) return { min: 0, max: 0 };
    
    const prices = chartData.flatMap(d => [d.high, d.low]);
    const min = Math.min(...prices);
    const max = Math.max(...prices);
    
    // 添加5%的padding
    const padding = (max - min) * 0.05;
    return {
      min: min - padding,
      max: max + padding,
    };
  }, [chartData]);

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

          {/* 蜡烛图 - 使用 Bar 组件配合自定义形状 */}
          <Bar
            dataKey="candleHeight"
            shape={(props: any) => (
              <CandleShape
                {...props}
                chartHeight={314}
                yMin={priceRange.min}
                yMax={priceRange.max}
              />
            )}
            isAnimationActive={false}
          />

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

