/**
 * 买卖预测信号展示组件
 */

import React from 'react';
import {
  Box,
  Typography,
  Alert,
  AlertTitle,
  Paper,
  Chip,
} from '@mui/material';
import {
  TrendingUp as BuyIcon,
  TrendingDown as SellIcon,
  RemoveCircleOutline as HoldIcon,
} from '@mui/icons-material';

interface PredictionSignalsViewProps {
  stockCode: string;
}

const PredictionSignalsView: React.FC<PredictionSignalsViewProps> = ({ stockCode }) => {
  // TODO: 这里需要从API获取预测信号数据
  // 目前暂时显示占位内容

  return (
    <Box>
      <Alert severity="info" sx={{ mb: 3 }}>
        <AlertTitle>功能开发中</AlertTitle>
        买卖预测信号功能正在开发中，敬请期待！
      </Alert>

      {/* 示例界面布局 */}
      <Paper elevation={2} sx={{ p: 3, mb: 2 }}>
        <Box display="flex" alignItems="center" mb={2}>
          <BuyIcon color="success" sx={{ mr: 1, fontSize: 40 }} />
          <Box>
            <Typography variant="h6" color="success.main">
              买入信号
            </Typography>
            <Typography variant="body2" color="text.secondary">
              综合技术指标显示买入机会
            </Typography>
          </Box>
          <Chip label="强度: 75%" color="success" sx={{ ml: 'auto' }} />
        </Box>
        <Typography variant="body2" color="text.secondary">
          • MACD金叉形成<br />
          • RSI指标超卖反弹<br />
          • 成交量放大
        </Typography>
      </Paper>

      <Paper elevation={2} sx={{ p: 3, mb: 2 }}>
        <Box display="flex" alignItems="center" mb={2}>
          <SellIcon color="error" sx={{ mr: 1, fontSize: 40 }} />
          <Box>
            <Typography variant="h6" color="error.main">
              卖出信号
            </Typography>
            <Typography variant="body2" color="text.secondary">
              技术指标提示风险
            </Typography>
          </Box>
          <Chip label="强度: 60%" color="error" sx={{ ml: 'auto' }} />
        </Box>
        <Typography variant="body2" color="text.secondary">
          • KDJ指标超买<br />
          • 价格突破布林带上轨<br />
          • 成交量萎缩
        </Typography>
      </Paper>

      <Paper elevation={2} sx={{ p: 3 }}>
        <Box display="flex" alignItems="center" mb={2}>
          <HoldIcon color="action" sx={{ mr: 1, fontSize: 40 }} />
          <Box>
            <Typography variant="h6">
              观望信号
            </Typography>
            <Typography variant="body2" color="text.secondary">
              当前无明确买卖信号
            </Typography>
          </Box>
          <Chip label="等待中" sx={{ ml: 'auto' }} />
        </Box>
        <Typography variant="body2" color="text.secondary">
          • 指标处于中性区间<br />
          • 建议等待更明确信号<br />
          • 注意市场整体趋势
        </Typography>
      </Paper>

      <Box mt={3}>
        <Typography variant="caption" color="text.secondary">
          * 预测信号仅供参考，投资有风险，入市需谨慎
        </Typography>
      </Box>
    </Box>
  );
};

export default PredictionSignalsView;

