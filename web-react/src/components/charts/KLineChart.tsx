/**
 * K线图组件 - 使用 Recharts
 */

import React, { useMemo } from 'react';
import {
  ComposedChart,
  Line,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { Box, Typography, useTheme } from '@mui/material';
import { StockDaily } from '../../types/stock';

interface KLineChartProps {
  data: StockDaily[];
}

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
        bodyTop,
        bodyBottom,
        bodyHeight,
        isUp: close >= open,
        // 用于显示K线上下影线
        shadowTop: high,
        shadowBottom: low,
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
      {/* 价格走势图 */}
      <ResponsiveContainer width="100%" height={400}>
        <ComposedChart data={chartData} margin={{ top: 20, right: 30, left: 0, bottom: 0 }}>
          <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.divider} />
          <XAxis
            dataKey="date"
            tickFormatter={formatDate}
            stroke={theme.palette.text.secondary}
          />
          <YAxis
            yAxisId="price"
            domain={['auto', 'auto']}
            tickFormatter={formatPrice}
            stroke={theme.palette.text.secondary}
          />
          <Tooltip content={<CustomTooltip />} />
          <Legend />

          {/* 收盘价折线 */}
          <Line
            yAxisId="price"
            type="monotone"
            dataKey="close"
            stroke={theme.palette.primary.main}
            strokeWidth={2}
            dot={false}
            name="收盘价"
          />

          {/* MA5 */}
          <Line
            yAxisId="price"
            type="monotone"
            dataKey={(data) => {
              // 简单计算5日均线
              const index = chartData.indexOf(data);
              if (index < 4) return null;
              const sum = chartData.slice(index - 4, index + 1).reduce((acc, item) => acc + item.close, 0);
              return sum / 5;
            }}
            stroke={theme.palette.warning.main}
            strokeWidth={1}
            dot={false}
            name="MA5"
          />

          {/* MA10 */}
          <Line
            yAxisId="price"
            type="monotone"
            dataKey={(data) => {
              // 简单计算10日均线
              const index = chartData.indexOf(data);
              if (index < 9) return null;
              const sum = chartData.slice(index - 9, index + 1).reduce((acc, item) => acc + item.close, 0);
              return sum / 10;
            }}
            stroke={theme.palette.info.main}
            strokeWidth={1}
            dot={false}
            name="MA10"
          />
        </ComposedChart>
      </ResponsiveContainer>

      {/* 成交量图 */}
      <ResponsiveContainer width="100%" height={150}>
        <ComposedChart data={chartData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
          <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.divider} />
          <XAxis
            dataKey="date"
            tickFormatter={formatDate}
            stroke={theme.palette.text.secondary}
          />
          <YAxis
            tickFormatter={formatVolume}
            stroke={theme.palette.text.secondary}
          />
          <Tooltip
            formatter={(value: number) => formatVolume(value)}
            labelFormatter={(label) => `日期: ${label}`}
          />
          <Bar
            dataKey="vol"
            fill={theme.palette.primary.main}
            name="成交量"
          />
        </ComposedChart>
      </ResponsiveContainer>
    </Box>
  );
};

export default KLineChart;

