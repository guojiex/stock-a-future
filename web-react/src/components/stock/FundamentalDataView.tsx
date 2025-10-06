/**
 * 基本面数据展示组件
 */

import React from 'react';
import {
  Box,
  Typography,
  Paper,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
} from '@mui/material';
import { FundamentalData } from '../../types/stock';

interface FundamentalDataViewProps {
  data: FundamentalData;
}

const FundamentalDataView: React.FC<FundamentalDataViewProps> = ({ data }) => {
  // 格式化数字（亿元）
  const formatBillion = (value?: number) => {
    if (value === undefined || value === null) return '--';
    return (value / 100000000).toFixed(2);
  };

  // 格式化百分比
  const formatPercent = (value?: number) => {
    if (value === undefined || value === null) return '--';
    return `${value.toFixed(2)}%`;
  };

  // 数据卡片
  const DataCard = ({ title, value, unit }: { title: string; value: string; unit?: string }) => (
    <Paper elevation={1} sx={{ p: 2 }}>
      <Typography variant="body2" color="text.secondary" gutterBottom>
        {title}
      </Typography>
      <Typography variant="h6">
        {value} {unit && <Typography component="span" variant="body2" color="text.secondary">{unit}</Typography>}
      </Typography>
    </Paper>
  );

  return (
    <Box>
      {/* 市场数据 */}
      {data.daily_basic && (
        <Box mb={4}>
          <Typography variant="h6" gutterBottom>
            市场数据
          </Typography>
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: { xs: 'repeat(2, 1fr)', md: 'repeat(4, 1fr)' },
              gap: 2,
            }}
          >
            <DataCard
              title="总市值"
              value={formatBillion(data.daily_basic.total_mv)}
              unit="亿元"
            />
            <DataCard
              title="流通市值"
              value={formatBillion(data.daily_basic.circ_mv)}
              unit="亿元"
            />
            <DataCard
              title="市盈率(PE)"
              value={data.daily_basic.pe?.toFixed(2) || '--'}
            />
            <DataCard
              title="市净率(PB)"
              value={data.daily_basic.pb?.toFixed(2) || '--'}
            />
            <DataCard
              title="市销率(PS)"
              value={data.daily_basic.ps?.toFixed(2) || '--'}
            />
            <DataCard
              title="换手率"
              value={formatPercent(data.daily_basic.turnover_rate)}
            />
          </Box>
        </Box>
      )}

      <Divider sx={{ my: 3 }} />

      {/* 利润表 */}
      {data.income_statement && (
        <Box mb={4}>
          <Typography variant="h6" gutterBottom>
            利润表 ({data.income_statement.end_date})
          </Typography>
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: { xs: 'repeat(2, 1fr)', md: 'repeat(4, 1fr)' },
              gap: 2,
            }}
          >
            <DataCard
              title="营业总收入"
              value={formatBillion(data.income_statement.total_revenue)}
              unit="亿元"
            />
            <DataCard
              title="营业收入"
              value={formatBillion(data.income_statement.revenue)}
              unit="亿元"
            />
            <DataCard
              title="利润总额"
              value={formatBillion(data.income_statement.total_profit)}
              unit="亿元"
            />
            <DataCard
              title="净利润"
              value={formatBillion(data.income_statement.n_income)}
              unit="亿元"
            />
            <DataCard
              title="基本每股收益"
              value={data.income_statement.basic_eps.toFixed(2)}
              unit="元"
            />
          </Box>
        </Box>
      )}

      <Divider sx={{ my: 3 }} />

      {/* 资产负债表 */}
      {data.balance_sheet && (
        <Box mb={4}>
          <Typography variant="h6" gutterBottom>
            资产负债表 ({data.balance_sheet.end_date})
          </Typography>
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: { xs: 'repeat(2, 1fr)', md: 'repeat(4, 1fr)' },
              gap: 2,
            }}
          >
            <DataCard
              title="资产总计"
              value={formatBillion(data.balance_sheet.total_assets)}
              unit="亿元"
            />
            <DataCard
              title="负债合计"
              value={formatBillion(data.balance_sheet.total_liab)}
              unit="亿元"
            />
            <DataCard
              title="股东权益"
              value={formatBillion(data.balance_sheet.total_hldr_eqy_inc_min_int)}
              unit="亿元"
            />
            <DataCard
              title="总股本"
              value={formatBillion(data.balance_sheet.total_share)}
              unit="亿股"
            />
          </Box>
        </Box>
      )}

      <Divider sx={{ my: 3 }} />

      {/* 现金流量表 */}
      {data.cash_flow_statement && (
        <Box mb={4}>
          <Typography variant="h6" gutterBottom>
            现金流量表 ({data.cash_flow_statement.end_date})
          </Typography>
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: { xs: '1fr', md: 'repeat(3, 1fr)' },
              gap: 2,
            }}
          >
            <DataCard
              title="经营活动现金流"
              value={formatBillion(data.cash_flow_statement.n_cashflow_operate_a)}
              unit="亿元"
            />
            <DataCard
              title="投资活动现金流"
              value={formatBillion(data.cash_flow_statement.n_cashflow_invest_a)}
              unit="亿元"
            />
            <DataCard
              title="筹资活动现金流"
              value={formatBillion(data.cash_flow_statement.n_cash_flows_fnc_act)}
              unit="亿元"
            />
          </Box>
        </Box>
      )}

      {/* 股票基本信息 */}
      {data.stock_basic && (
        <Box>
          <Typography variant="h6" gutterBottom>
            股票基本信息
          </Typography>
          <TableContainer component={Paper}>
            <Table size="small">
              <TableBody>
                <TableRow>
                  <TableCell component="th">股票代码</TableCell>
                  <TableCell>{data.stock_basic.ts_code}</TableCell>
                  <TableCell component="th">股票名称</TableCell>
                  <TableCell>{data.stock_basic.name}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell component="th">所属行业</TableCell>
                  <TableCell>{data.stock_basic.industry || '--'}</TableCell>
                  <TableCell component="th">所属地域</TableCell>
                  <TableCell>{data.stock_basic.area || '--'}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell component="th">上市日期</TableCell>
                  <TableCell>{data.stock_basic.list_date || '--'}</TableCell>
                  <TableCell component="th">市场类型</TableCell>
                  <TableCell>{data.stock_basic.market || '--'}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </TableContainer>
        </Box>
      )}

      <Box mt={3}>
        <Typography variant="caption" color="text.secondary">
          * 财务数据来源于最新财报，仅供参考
        </Typography>
      </Box>
    </Box>
  );
};

export default FundamentalDataView;

