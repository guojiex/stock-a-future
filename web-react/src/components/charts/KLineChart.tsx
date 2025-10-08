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

// 自定义蜡烛图形状 - 使用 Bar 的 shape 属性
const Candlestick = (props: any) => {
  const { x, y, width, height, payload } = props;
  
  if (!payload || !payload.open || !payload.close || !payload.high || !payload.low) {
    return null;
  }

  const { open, close, high, low, isUp } = payload;
  // 中国股市：红色代表涨，绿色代表跌
  const color = isUp ? '#ef5350' : '#26a69a';
  
  // y 是 Bar 的顶部位置（对应 high 值）
  // height 是 Bar 的高度（对应 high - low）
  
  // 计算 open 和 close 在 bar 中的相对位置
  const range = high - low;
  const openRatio = range > 0 ? (high - open) / range : 0.5;
  const closeRatio = range > 0 ? (high - close) / range : 0.5;
  
  const openY = y + height * openRatio;
  const closeY = y + height * closeRatio;
  
  // 蜡烛实体的top和height
  const bodyTop = Math.min(openY, closeY);
  let bodyHeight = Math.abs(closeY - openY);
  
  // 十字星的最小高度
  if (bodyHeight < 1) bodyHeight = 1;
  
  // 计算宽度
  const wickWidth = Math.max(1, width * 0.15);
  const bodyWidth = Math.max(2, width * 0.6);
  const centerX = x + width / 2;
  
  return (
    <g>
      {/* 上下影线 */}
      <line
        x1={centerX}
        y1={y}
        x2={centerX}
        y2={y + height}
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
};

const KLineChart: React.FC<KLineChartProps> = ({ data }) => {
  const theme = useTheme();

  // 转换数据格式
  const chartData = useMemo(() => {
    if (!data || data.length === 0) return [];

    return data.map((item, index) => {
      const open = parseFloat(item.open);
      const close = parseFloat(item.close);
      const high = parseFloat(item.high);
      const low = parseFloat(item.low);
      const vol = parseFloat(item.vol);

      // 计算涨跌幅
      let change = 0;
      let changePercent = 0;
      
      if (index > 0) {
        const prevClose = parseFloat(data[index - 1].close);
        change = close - prevClose;
        changePercent = (change / prevClose) * 100;
      }

      return {
        date: item.trade_date,
        open,
        close,
        high,
        low,
        vol,
        isUp: close >= open,
        change,
        changePercent,
        // 用于绘制蜡烛图：使用一个从 low 开始、高度为 high-low 的数组
        candle: [low, high],
      };
    });
  }, [data]);

  // 格式化日期 - 将日期转换为 MM/DD 格式
  const formatDate = (dateStr: string) => {
    if (!dateStr) return dateStr;
    
    // 处理 YYYYMMDD 格式 (如: "20250717")
    if (dateStr.length === 8 && /^\d{8}$/.test(dateStr)) {
      return `${dateStr.substring(4, 6)}/${dateStr.substring(6, 8)}`;
    }
    
    // 处理 ISO 日期格式 (如: "2025-08-01T00:00:00.000" 或 "2025-08-01")
    if (dateStr.includes('T')) {
      // 提取日期部分，去掉时间部分
      const datePart = dateStr.split('T')[0];
      const [, month, day] = datePart.split('-');
      return `${month}/${day}`;
    }
    
    // 处理 YYYY-MM-DD 格式
    if (dateStr.includes('-')) {
      const [, month, day] = dateStr.split('-');
      return `${month}/${day}`;
    }
    
    return dateStr;
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

    // 格式化完整日期显示（YYYY-MM-DD 格式）
    const formatFullDate = (dateStr: string) => {
      if (!dateStr) return dateStr;
      
      // 处理 YYYYMMDD 格式
      if (dateStr.length === 8 && /^\d{8}$/.test(dateStr)) {
        return `${dateStr.substring(0, 4)}-${dateStr.substring(4, 6)}-${dateStr.substring(6, 8)}`;
      }
      
      // 处理 ISO 日期格式，提取日期部分
      if (dateStr.includes('T')) {
        return dateStr.split('T')[0];
      }
      
      // 已经是 YYYY-MM-DD 格式
      if (dateStr.includes('-')) {
        return dateStr;
      }
      
      return dateStr;
    };

    // 判断涨跌（基于收盘价相对于前一日）
    // 中国股市：红色代表涨，绿色代表跌
    const isPriceUp = data.change >= 0;
    const priceColor = isPriceUp ? 'error.main' : 'success.main'; // 红涨绿跌

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
        <Typography variant="body2" gutterBottom fontWeight="bold">
          日期: {formatFullDate(data.date)}
        </Typography>
        <Typography variant="body2" color={data.isUp ? 'error.main' : 'success.main'}>
          开: {formatPrice(data.open)}
        </Typography>
        <Typography variant="body2" color={data.isUp ? 'error.main' : 'success.main'}>
          高: {formatPrice(data.high)}
        </Typography>
        <Typography variant="body2" color={data.isUp ? 'error.main' : 'success.main'}>
          低: {formatPrice(data.low)}
        </Typography>
        <Typography variant="body2" color={data.isUp ? 'error.main' : 'success.main'}>
          收: {formatPrice(data.close)}
        </Typography>
        {data.change !== 0 && (
          <Typography variant="body2" color={priceColor} fontWeight="bold">
            涨跌: {isPriceUp ? '+' : ''}{data.change.toFixed(2)} ({isPriceUp ? '+' : ''}{data.changePercent.toFixed(2)}%)
          </Typography>
        )}
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
            dataKey="candle"
            shape={Candlestick}
            name="K线"
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
              <Cell key={`cell-${index}`} fill={entry.isUp ? '#ef5350' : '#26a69a'} />
            ))}
          </Bar>
        </ComposedChart>
      </ResponsiveContainer>
    </Box>
  );
};

export default KLineChart;

