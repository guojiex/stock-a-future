/**
 * 交易记录表格组件
 * 显示回测过程中的所有交易记录
 */

import React, { useState } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Paper,
  Chip,
  Typography,
  Box,
} from '@mui/material';
import {
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
} from '@mui/icons-material';

interface Trade {
  id: string;
  backtest_id: string;
  strategy_id?: string;
  symbol: string;
  side: 'buy' | 'sell';
  quantity: number;
  price: number;
  commission: number;
  pnl?: number;
  signal_type?: string;
  total_assets?: number;
  holding_assets?: number;
  cash_balance?: number;
  timestamp: string;
  created_at?: string;
}

interface TradesTableProps {
  trades: Trade[];
  strategies?: Array<{ id: string; name: string }>;
}

const TradesTable: React.FC<TradesTableProps> = ({ trades, strategies }) => {
  // 创建策略ID到名称的映射
  const strategyMap = React.useMemo(() => {
    const map: Record<string, string> = {};
    if (strategies) {
      strategies.forEach(s => {
        map[s.id] = s.name;
      });
    }
    return map;
  }, [strategies]);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);

  const handleChangePage = (_event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  if (!trades || trades.length === 0) {
    return (
      <Box sx={{ p: 4, textAlign: 'center', bgcolor: 'grey.100', borderRadius: 1 }}>
        <Typography variant="body2" color="text.secondary">
          暂无交易记录
        </Typography>
      </Box>
    );
  }

  // 计算统计信息
  const buyTrades = trades.filter(t => t.side === 'buy').length;
  const sellTrades = trades.filter(t => t.side === 'sell').length;
  const totalPnL = trades.reduce((sum, t) => sum + (t.pnl || 0), 0);
  const profitTrades = trades.filter(t => t.side === 'sell' && (t.pnl || 0) > 0).length;
  const lossTrades = trades.filter(t => t.side === 'sell' && (t.pnl || 0) < 0).length;

  const paginatedTrades = trades.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);

  return (
    <Box>
      {/* 统计信息 */}
      <Box sx={{ display: 'flex', gap: 3, mb: 2, flexWrap: 'wrap' }}>
        <Box>
          <Typography variant="body2" color="text.secondary">
            总交易次数
          </Typography>
          <Typography variant="h6">
            {trades.length}
          </Typography>
        </Box>
        <Box>
          <Typography variant="body2" color="text.secondary">
            买入 / 卖出
          </Typography>
          <Typography variant="h6">
            {buyTrades} / {sellTrades}
          </Typography>
        </Box>
        <Box>
          <Typography variant="body2" color="text.secondary">
            盈利 / 亏损
          </Typography>
          <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <span style={{ color: 'green' }}>{profitTrades}</span>
            <span>/</span>
            <span style={{ color: 'red' }}>{lossTrades}</span>
          </Typography>
        </Box>
        <Box>
          <Typography variant="body2" color="text.secondary">
            累计盈亏
          </Typography>
          <Typography
            variant="h6"
            color={totalPnL >= 0 ? 'success.main' : 'error.main'}
          >
            {totalPnL >= 0 ? '+' : ''}¥{totalPnL.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
          </Typography>
        </Box>
      </Box>

      {/* 交易表格 */}
      <TableContainer component={Paper} variant="outlined">
        <Table size="small">
          <TableHead>
            <TableRow sx={{ bgcolor: 'grey.50' }}>
              <TableCell>时间</TableCell>
              <TableCell>股票代码</TableCell>
              {strategies && strategies.length > 1 && <TableCell>策略</TableCell>}
              <TableCell align="center">方向</TableCell>
              <TableCell align="right">数量</TableCell>
              <TableCell align="right">价格</TableCell>
              <TableCell align="right">金额</TableCell>
              <TableCell align="right">手续费</TableCell>
              <TableCell align="right">盈亏</TableCell>
              <TableCell align="right">总资产</TableCell>
              <TableCell>信号类型</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {paginatedTrades.map((trade) => {
              const amount = trade.price * trade.quantity;
              const isPnLPositive = (trade.pnl || 0) >= 0;
              
              return (
                <TableRow
                  key={trade.id}
                  hover
                  sx={{
                    '&:last-child td, &:last-child th': { border: 0 },
                  }}
                >
                  <TableCell>
                    <Typography variant="body2" sx={{ fontSize: '0.75rem' }}>
                      {new Date(trade.timestamp).toLocaleString('zh-CN', {
                        year: 'numeric',
                        month: '2-digit',
                        day: '2-digit',
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2" fontWeight="medium">
                      {trade.symbol}
                    </Typography>
                  </TableCell>
                  {strategies && strategies.length > 1 && (
                    <TableCell>
                      {trade.strategy_id ? (
                        <Chip
                          label={strategyMap[trade.strategy_id] || trade.strategy_id}
                          size="small"
                          variant="outlined"
                          color="primary"
                        />
                      ) : (
                        <Typography variant="body2" color="text.secondary">
                          -
                        </Typography>
                      )}
                    </TableCell>
                  )}
                  <TableCell align="center">
                    <Chip
                      icon={trade.side === 'buy' ? <TrendingUpIcon /> : <TrendingDownIcon />}
                      label={trade.side === 'buy' ? '买入' : '卖出'}
                      color={trade.side === 'buy' ? 'error' : 'success'}
                      size="small"
                      sx={{ minWidth: 70 }}
                    />
                  </TableCell>
                  <TableCell align="right">
                    <Typography variant="body2">
                      {trade.quantity}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">
                    <Typography variant="body2">
                      ¥{trade.price.toFixed(2)}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">
                    <Typography variant="body2">
                      ¥{amount.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">
                    <Typography variant="body2" color="text.secondary">
                      ¥{trade.commission.toFixed(2)}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">
                    {trade.pnl !== undefined && trade.pnl !== 0 ? (
                      <Typography
                        variant="body2"
                        fontWeight="medium"
                        color={isPnLPositive ? 'success.main' : 'error.main'}
                      >
                        {isPnLPositive ? '+' : ''}¥{trade.pnl.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                      </Typography>
                    ) : (
                      <Typography variant="body2" color="text.secondary">
                        -
                      </Typography>
                    )}
                  </TableCell>
                  <TableCell align="right">
                    {trade.total_assets !== undefined ? (
                      <Typography variant="body2">
                        ¥{trade.total_assets.toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 0 })}
                      </Typography>
                    ) : (
                      <Typography variant="body2" color="text.secondary">
                        -
                      </Typography>
                    )}
                  </TableCell>
                  <TableCell>
                    {trade.signal_type ? (
                      <Chip
                        label={trade.signal_type}
                        size="small"
                        variant="outlined"
                      />
                    ) : (
                      <Typography variant="body2" color="text.secondary">
                        -
                      </Typography>
                    )}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
        <TablePagination
          component="div"
          count={trades.length}
          page={page}
          onPageChange={handleChangePage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          rowsPerPageOptions={[5, 10, 25, 50]}
          labelRowsPerPage="每页行数:"
          labelDisplayedRows={({ from, to, count }) => `${from}-${to} 共 ${count}`}
        />
      </TableContainer>
    </Box>
  );
};

export default TradesTable;

