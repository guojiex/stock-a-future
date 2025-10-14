/**
 * 权益曲线图表组件
 * 使用 Recharts 绘制回测的权益曲线
 */

import React from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  TooltipProps,
} from 'recharts';
import { Box, Typography, useTheme } from '@mui/material';
import { NameType, ValueType } from 'recharts/types/component/DefaultTooltipContent';

interface EquityPoint {
  date: string;
  portfolio_value: number;
  benchmark_value?: number;
  cash?: number;
  holdings?: number;
}

interface EquityCurveChartProps {
  data: EquityPoint[];
  initialCash: number;
}

// 自定义 Tooltip
const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <Box
        sx={{
          backgroundColor: 'background.paper',
          border: 1,
          borderColor: 'divider',
          borderRadius: 1,
          p: 1.5,
          boxShadow: 2,
        }}
      >
        <Typography variant="body2" fontWeight="bold" gutterBottom>
          {label}
        </Typography>
        {payload.map((entry, index) => {
          const value = entry.value as number;
          const name = entry.name as string;
          const color = entry.color;
          
          let displayName = name;
          if (name === 'portfolio_value') displayName = '组合价值';
          else if (name === 'benchmark_value') displayName = '基准价值';
          else if (name === 'cash') displayName = '现金';
          else if (name === 'holdings') displayName = '持仓';
          
          return (
            <Typography key={index} variant="body2" sx={{ color }}>
              {displayName}: ¥{value.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
            </Typography>
          );
        })}
      </Box>
    );
  }
  return null;
};

const EquityCurveChart: React.FC<EquityCurveChartProps> = ({ data, initialCash }) => {
  const theme = useTheme();
  
  if (!data || data.length === 0) {
    return (
      <Box sx={{ p: 4, textAlign: 'center', bgcolor: 'grey.100', borderRadius: 1 }}>
        <Typography variant="body2" color="text.secondary">
          暂无权益曲线数据
        </Typography>
      </Box>
    );
  }
  
  // 计算收益率
  const finalValue = data[data.length - 1]?.portfolio_value || initialCash;
  const totalReturn = ((finalValue - initialCash) / initialCash) * 100;
  const isProfit = totalReturn >= 0;
  
  return (
    <Box>
      {/* 统计信息 */}
      <Box sx={{ display: 'flex', gap: 3, mb: 2, flexWrap: 'wrap' }}>
        <Box>
          <Typography variant="body2" color="text.secondary">
            初始资金
          </Typography>
          <Typography variant="h6">
            ¥{initialCash.toLocaleString('zh-CN')}
          </Typography>
        </Box>
        <Box>
          <Typography variant="body2" color="text.secondary">
            最终价值
          </Typography>
          <Typography variant="h6">
            ¥{finalValue.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
          </Typography>
        </Box>
        <Box>
          <Typography variant="body2" color="text.secondary">
            总收益率
          </Typography>
          <Typography
            variant="h6"
            color={isProfit ? 'success.main' : 'error.main'}
          >
            {isProfit ? '+' : ''}{totalReturn.toFixed(2)}%
          </Typography>
        </Box>
      </Box>
      
      {/* 图表 */}
      <ResponsiveContainer width="100%" height={400}>
        <LineChart
          data={data}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
        >
          <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.divider} />
          <XAxis
            dataKey="date"
            stroke={theme.palette.text.secondary}
            style={{ fontSize: '12px' }}
            tick={{ fill: theme.palette.text.secondary }}
          />
          <YAxis
            stroke={theme.palette.text.secondary}
            style={{ fontSize: '12px' }}
            tick={{ fill: theme.palette.text.secondary }}
            tickFormatter={(value) => `¥${(value / 10000).toFixed(0)}万`}
          />
          <Tooltip content={<CustomTooltip />} />
          <Legend
            wrapperStyle={{ paddingTop: '20px' }}
            formatter={(value) => {
              if (value === 'portfolio_value') return '组合价值';
              if (value === 'benchmark_value') return '基准价值';
              if (value === 'cash') return '现金';
              if (value === 'holdings') return '持仓';
              return value;
            }}
          />
          
          {/* 组合价值曲线 */}
          <Line
            type="monotone"
            dataKey="portfolio_value"
            stroke={theme.palette.primary.main}
            strokeWidth={2}
            dot={false}
            activeDot={{ r: 6 }}
          />
          
          {/* 基准价值曲线（如果有） */}
          {data.some(d => d.benchmark_value !== undefined) && (
            <Line
              type="monotone"
              dataKey="benchmark_value"
              stroke={theme.palette.secondary.main}
              strokeWidth={2}
              strokeDasharray="5 5"
              dot={false}
              activeDot={{ r: 6 }}
            />
          )}
          
          {/* 现金曲线 */}
          {data.some(d => d.cash !== undefined) && (
            <Line
              type="monotone"
              dataKey="cash"
              stroke={theme.palette.success.main}
              strokeWidth={1}
              dot={false}
              opacity={0.5}
            />
          )}
          
          {/* 持仓曲线 */}
          {data.some(d => d.holdings !== undefined) && (
            <Line
              type="monotone"
              dataKey="holdings"
              stroke={theme.palette.warning.main}
              strokeWidth={1}
              dot={false}
              opacity={0.5}
            />
          )}
        </LineChart>
      </ResponsiveContainer>
    </Box>
  );
};

export default EquityCurveChart;

