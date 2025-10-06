/**
 * 技术指标展示组件
 */

import React from 'react';
import {
  Box,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from '@mui/material';
import { TechnicalIndicators } from '../../types/stock';

interface TechnicalIndicatorsViewProps {
  data: TechnicalIndicators[];
}

const TechnicalIndicatorsView: React.FC<TechnicalIndicatorsViewProps> = ({ data }) => {
  if (!data || data.length === 0) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
        <Typography variant="body1" color="text.secondary">
          暂无技术指标数据
        </Typography>
      </Box>
    );
  }

  // 获取最新数据
  const latest = data[data.length - 1];

  // 指标卡片
  const IndicatorCard = ({ title, value, color }: { title: string; value?: number; color?: string }) => (
    <Paper elevation={2} sx={{ p: 2, textAlign: 'center' }}>
      <Typography variant="body2" color="text.secondary" gutterBottom>
        {title}
      </Typography>
      <Typography variant="h5" sx={{ color: color || 'text.primary' }}>
        {value !== undefined && value !== null ? value.toFixed(2) : '--'}
      </Typography>
    </Paper>
  );

  return (
    <Box>
      {/* 最新指标概览 */}
      <Box mb={3}>
        <Typography variant="h6" gutterBottom>
          当前指标 ({latest.trade_date})
        </Typography>

        <Box
          sx={{
            display: 'grid',
            gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' },
            gap: 2,
          }}
        >
          {/* 移动平均线 */}
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              移动平均线 (MA)
            </Typography>
            <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 1 }}>
              <IndicatorCard title="MA5" value={latest.ma5} />
              <IndicatorCard title="MA10" value={latest.ma10} />
              <IndicatorCard title="MA20" value={latest.ma20} />
            </Box>
          </Box>

          {/* MACD */}
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              MACD
            </Typography>
            <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 1 }}>
              <IndicatorCard title="DIF" value={latest.macd_dif} />
              <IndicatorCard title="DEA" value={latest.macd_dea} />
              <IndicatorCard
                title="MACD"
                value={latest.macd_histogram}
                color={
                  latest.macd_histogram && latest.macd_histogram > 0
                    ? 'success.main'
                    : 'error.main'
                }
              />
            </Box>
          </Box>

          {/* RSI */}
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              RSI 相对强弱指标
            </Typography>
            <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 1 }}>
              <IndicatorCard title="RSI6" value={latest.rsi6} />
              <IndicatorCard title="RSI12" value={latest.rsi12} />
              <IndicatorCard title="RSI24" value={latest.rsi24} />
            </Box>
          </Box>

          {/* 布林带 */}
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              布林带 (BOLL)
            </Typography>
            <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 1 }}>
              <IndicatorCard title="上轨" value={latest.boll_upper} />
              <IndicatorCard title="中轨" value={latest.boll_mid} />
              <IndicatorCard title="下轨" value={latest.boll_lower} />
            </Box>
          </Box>

          {/* KDJ */}
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              KDJ 随机指标
            </Typography>
            <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 1 }}>
              <IndicatorCard title="K" value={latest.kdj_k} />
              <IndicatorCard title="D" value={latest.kdj_d} />
              <IndicatorCard title="J" value={latest.kdj_j} />
            </Box>
          </Box>
        </Box>
      </Box>

      {/* 历史数据表格 */}
      <Box mt={4}>
        <Typography variant="h6" gutterBottom>
          历史技术指标
        </Typography>
        <TableContainer component={Paper}>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>日期</TableCell>
                <TableCell align="right">MA5</TableCell>
                <TableCell align="right">MA10</TableCell>
                <TableCell align="right">MA20</TableCell>
                <TableCell align="right">RSI6</TableCell>
                <TableCell align="right">MACD</TableCell>
                <TableCell align="right">KDJ-K</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {data.slice(-20).reverse().map((row) => (
                <TableRow key={row.trade_date} hover>
                  <TableCell>{row.trade_date}</TableCell>
                  <TableCell align="right">{row.ma5?.toFixed(2) || '--'}</TableCell>
                  <TableCell align="right">{row.ma10?.toFixed(2) || '--'}</TableCell>
                  <TableCell align="right">{row.ma20?.toFixed(2) || '--'}</TableCell>
                  <TableCell align="right">{row.rsi6?.toFixed(2) || '--'}</TableCell>
                  <TableCell align="right">{row.macd_histogram?.toFixed(4) || '--'}</TableCell>
                  <TableCell align="right">{row.kdj_k?.toFixed(2) || '--'}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Box>
    </Box>
  );
};

export default TechnicalIndicatorsView;

