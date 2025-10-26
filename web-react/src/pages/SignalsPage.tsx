/**
 * 收藏股票信号汇总页面
 * 展示所有收藏股票的买卖信号汇总
 */

import React, { useState, useMemo } from 'react';
import {
  Box,
  Typography,
  Paper,
  Alert,
  Chip,
  Button,
  Card,
  CardContent,
  CircularProgress,
  Tabs,
  Tab,
  IconButton,
  Tooltip,
  useTheme,
  alpha,
} from '@mui/material';
import RefreshIcon from '@mui/icons-material/Refresh';
import ShowChartIcon from '@mui/icons-material/ShowChart';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import TrendingDownIcon from '@mui/icons-material/TrendingDown';
import RemoveCircleOutlineIcon from '@mui/icons-material/RemoveCircleOutline';
import { useGetFavoritesSignalsQuery } from '../services/api';
import { FavoriteSignal, TradingPointPrediction } from '../types/stock';
import { useNavigate } from 'react-router-dom';

// 信号类型选项
type SignalType = 'all' | 'buy' | 'sell' | 'hold';

/**
 * 信号卡片组件
 */
const SignalCard: React.FC<{ signal: FavoriteSignal; onViewStock: (tsCode: string) => void }> = ({
  signal,
  onViewStock,
}) => {
  const theme = useTheme();
  
  // 获取预测数组（处理可能的对象格式）
  const predictions = useMemo(() => {
    if (!signal.predictions) return [];
    // 如果是数组，直接返回
    if (Array.isArray(signal.predictions)) return signal.predictions;
    // 如果是对象且有predictions属性，返回该属性
    if (typeof signal.predictions === 'object' && 'predictions' in signal.predictions) {
      const preds = (signal.predictions as any).predictions;
      return Array.isArray(preds) ? preds : [];
    }
    return [];
  }, [signal.predictions]);

  // 分析信号类型
  const signalType = useMemo(() => {
    if (predictions.length === 0) return 'hold';
    
    const hasBuy = predictions.some((p) => p.type === 'BUY');
    const hasSell = predictions.some((p) => p.type === 'SELL');
    
    if (hasBuy && !hasSell) return 'buy';
    if (hasSell && !hasBuy) return 'sell';
    if (hasBuy && hasSell) return 'mixed';
    return 'hold';
  }, [predictions]);

  // 对于混合信号，计算主要倾向
  const dominantSignal = useMemo(() => {
    if (signalType !== 'mixed') return null;
    
    const buyPreds = predictions.filter((p) => p.type === 'BUY');
    const sellPreds = predictions.filter((p) => p.type === 'SELL');
    const buyConfidence = buyPreds.length > 0 ? Math.max(...buyPreds.map((p) => p.probability)) : 0;
    const sellConfidence = sellPreds.length > 0 ? Math.max(...sellPreds.map((p) => p.probability)) : 0;
    
    if (buyConfidence > sellConfidence) return 'buy';
    if (sellConfidence > buyConfidence) return 'sell';
    return 'neutral';
  }, [signalType, predictions]);

  // 获取最高置信度
  const maxConfidence = useMemo(() => {
    if (predictions.length === 0) return 0;
    return Math.max(...predictions.map((p) => p.probability));
  }, [predictions]);

  // 信号类型的样式配置
  const signalConfig = {
    buy: {
      color: theme.palette.success.main,
      bgColor: alpha(theme.palette.success.main, 0.1),
      icon: <TrendingUpIcon />,
      label: '买入信号',
    },
    sell: {
      color: theme.palette.error.main,
      bgColor: alpha(theme.palette.error.main, 0.1),
      icon: <TrendingDownIcon />,
      label: '卖出信号',
    },
    hold: {
      color: theme.palette.warning.main,
      bgColor: alpha(theme.palette.warning.main, 0.1),
      icon: <RemoveCircleOutlineIcon />,
      label: '持有',
    },
    mixed: {
      color: theme.palette.secondary.main,
      bgColor: alpha(theme.palette.secondary.main, 0.1),
      icon: <ShowChartIcon />,
      label: '混合信号',
    },
  };

  const config = signalConfig[signalType as keyof typeof signalConfig] || signalConfig.hold;

  return (
    <Card
      sx={{
        mb: 2,
        border: `2px solid ${config.color}`,
        borderRadius: 2,
        transition: 'all 0.3s ease',
        '&:hover': {
          transform: 'translateY(-4px)',
          boxShadow: 4,
        },
      }}
    >
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Box display="flex" alignItems="center" gap={1}>
            <Box sx={{ color: config.color }}>{config.icon}</Box>
            <Typography variant="h6" component="div">
              {signal.name}
              <Typography variant="caption" color="text.secondary" ml={1}>
                {signal.ts_code}
              </Typography>
            </Typography>
          </Box>
          <Box>
            <Chip
              label={config.label}
              sx={{
                backgroundColor: config.bgColor,
                color: config.color,
                fontWeight: 'bold',
              }}
            />
            {signalType === 'mixed' && dominantSignal && (
              <Chip
                label={dominantSignal === 'buy' ? '倾向买入' : dominantSignal === 'sell' ? '倾向卖出' : '中性'}
                size="small"
                sx={{
                  ml: 1,
                  backgroundColor: 
                    dominantSignal === 'buy' 
                      ? alpha(theme.palette.success.main, 0.1) 
                      : dominantSignal === 'sell' 
                      ? alpha(theme.palette.error.main, 0.1) 
                      : alpha(theme.palette.grey[500], 0.1),
                  color: 
                    dominantSignal === 'buy' 
                      ? theme.palette.success.main 
                      : dominantSignal === 'sell' 
                      ? theme.palette.error.main 
                      : theme.palette.text.secondary,
                  fontSize: '0.75rem',
                }}
              />
            )}
          </Box>
        </Box>

        {/* 价格和日期信息 */}
        <Box display="flex" flexWrap="wrap" gap={2} mb={2}>
          <Box flex="1 1 30%">
            <Typography variant="body2" color="text.secondary">
              当前价格
            </Typography>
            <Typography variant="h6" color="primary">
              ¥{parseFloat(signal.current_price).toFixed(2)}
            </Typography>
          </Box>
          <Box flex="1 1 30%">
            <Typography variant="body2" color="text.secondary">
              交易日期
            </Typography>
            <Typography variant="body1">{signal.trade_date}</Typography>
          </Box>
          <Box flex="1 1 30%">
            <Typography variant="body2" color="text.secondary">
              置信度
            </Typography>
            <Box display="flex" alignItems="center" gap={1}>
              <CircularProgress
                variant="determinate"
                value={maxConfidence * 100}
                size={24}
                sx={{ color: config.color }}
              />
              <Typography variant="body1">{(maxConfidence * 100).toFixed(1)}%</Typography>
            </Box>
          </Box>
        </Box>

        {/* 预测信号列表 */}
        {predictions.length > 0 && (
          <Box mb={2}>
            <Typography variant="subtitle2" color="text.secondary" mb={1}>
              预测信号：
            </Typography>
            {predictions.map((pred, index) => (
              <Box
                key={index}
                sx={{
                  p: 1,
                  mb: 1,
                  bgcolor: pred.type === 'BUY' 
                    ? alpha(theme.palette.success.main, 0.08) 
                    : alpha(theme.palette.error.main, 0.08),
                  borderRadius: 1,
                  borderLeft: `4px solid ${
                    pred.type === 'BUY' 
                      ? theme.palette.success.main 
                      : theme.palette.error.main
                  }`,
                }}
              >
                <Box display="flex" flexWrap="wrap" alignItems="center" justifyContent="space-between">
                  <Box flex="1 1 60%">
                    <Typography variant="body2" fontWeight="medium">
                      {pred.type === 'BUY' ? '🟢 买入' : '🔴 卖出'}: {pred.reason}
                    </Typography>
                    {pred.indicators && pred.indicators.length > 0 && (
                      <Typography variant="caption" color="text.secondary">
                        指标: {pred.indicators.join(', ')}
                      </Typography>
                    )}
                  </Box>
                  <Box flex="1 1 35%" textAlign="right">
                    <Typography variant="caption" color="text.secondary">
                      概率: {(pred.probability * 100).toFixed(1)}%
                    </Typography>
                  </Box>
                </Box>
              </Box>
            ))}
          </Box>
        )}

        {/* 查看详情按钮 */}
        <Box display="flex" justifyContent="flex-end">
          <Button
            variant="outlined"
            size="small"
            startIcon={<ShowChartIcon />}
            onClick={() => onViewStock(signal.ts_code)}
            sx={{ borderColor: config.color, color: config.color }}
          >
            查看详情
          </Button>
        </Box>

        {/* 更新时间 */}
        <Typography variant="caption" color="text.secondary" display="block" mt={1}>
          更新时间: {signal.updated_at}
        </Typography>
      </CardContent>
    </Card>
  );
};

/**
 * 收藏股票信号汇总页面
 */
const SignalsPage: React.FC = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const [currentTab, setCurrentTab] = useState<SignalType>('all');

  // 获取信号数据
  const { data, isLoading, isError, error, refetch } = useGetFavoritesSignalsQuery();

  // 过滤信号
  const filteredSignals = useMemo(() => {
    if (!data?.data?.signals) return [];

    const signals = data.data.signals;

    if (currentTab === 'all') return signals;

    return signals.filter((signal) => {
      // 获取预测数组
      let predictions: TradingPointPrediction[] = [];
      if (Array.isArray(signal.predictions)) {
        predictions = signal.predictions;
      } else if (signal.predictions && typeof signal.predictions === 'object' && 'predictions' in signal.predictions) {
        const preds = (signal.predictions as any).predictions;
        predictions = Array.isArray(preds) ? preds : [];
      }

      if (predictions.length === 0) {
        return currentTab === 'hold';
      }

      const hasBuy = predictions.some((p) => p.type === 'BUY');
      const hasSell = predictions.some((p) => p.type === 'SELL');

      switch (currentTab) {
        case 'buy':
          // 纯买入信号
          if (hasBuy && !hasSell) return true;
          // 混合信号，根据置信度判断
          if (hasBuy && hasSell) {
            const buyPreds = predictions.filter((p) => p.type === 'BUY');
            const sellPreds = predictions.filter((p) => p.type === 'SELL');
            const buyConfidence = buyPreds.length > 0 ? Math.max(...buyPreds.map((p) => p.probability)) : 0;
            const sellConfidence = sellPreds.length > 0 ? Math.max(...sellPreds.map((p) => p.probability)) : 0;
            return buyConfidence > sellConfidence;
          }
          return false;
        case 'sell':
          // 纯卖出信号
          if (hasSell && !hasBuy) return true;
          // 混合信号，根据置信度判断
          if (hasBuy && hasSell) {
            const buyPreds = predictions.filter((p) => p.type === 'BUY');
            const sellPreds = predictions.filter((p) => p.type === 'SELL');
            const buyConfidence = buyPreds.length > 0 ? Math.max(...buyPreds.map((p) => p.probability)) : 0;
            const sellConfidence = sellPreds.length > 0 ? Math.max(...sellPreds.map((p) => p.probability)) : 0;
            return sellConfidence > buyConfidence;
          }
          return false;
        case 'hold':
          return !hasBuy && !hasSell;
        default:
          return true;
      }
    });
  }, [data, currentTab]);

  // 统计信息
  const statistics = useMemo(() => {
    if (!data?.data?.signals) {
      return { total: 0, buy: 0, sell: 0, hold: 0 };
    }

    const signals = data.data.signals;
    let buy = 0;
    let sell = 0;
    let hold = 0;

    signals.forEach((signal) => {
      // 获取预测数组
      let predictions: TradingPointPrediction[] = [];
      if (Array.isArray(signal.predictions)) {
        predictions = signal.predictions;
      } else if (signal.predictions && typeof signal.predictions === 'object' && 'predictions' in signal.predictions) {
        const preds = (signal.predictions as any).predictions;
        predictions = Array.isArray(preds) ? preds : [];
      }

      if (predictions.length === 0) {
        hold++;
        return;
      }

      const hasBuy = predictions.some((p) => p.type === 'BUY');
      const hasSell = predictions.some((p) => p.type === 'SELL');

      if (hasBuy && !hasSell) buy++;
      else if (hasSell && !hasBuy) sell++;
      else if (!hasBuy && !hasSell) hold++;
      else {
        // 混合信号，根据置信度判断
        const buyPreds = predictions.filter((p) => p.type === 'BUY');
        const sellPreds = predictions.filter((p) => p.type === 'SELL');
        const buyConfidence = buyPreds.length > 0 ? Math.max(...buyPreds.map((p) => p.probability)) : 0;
        const sellConfidence = sellPreds.length > 0 ? Math.max(...sellPreds.map((p) => p.probability)) : 0;
        if (buyConfidence > sellConfidence) buy++;
        else if (sellConfidence > buyConfidence) sell++;
        else hold++;
      }
    });

    return { total: signals.length, buy, sell, hold };
  }, [data]);

  // 处理查看股票详情
  const handleViewStock = (tsCode: string) => {
    navigate(`/market?stock=${tsCode}`);
  };

  // 处理刷新
  const handleRefresh = () => {
    refetch();
  };

  return (
    <Box sx={{ p: 3 }}>
      {/* 标题和刷新按钮 */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1">
          收藏股票信号汇总
        </Typography>
        <Tooltip title="刷新信号数据">
          <span>
            <IconButton onClick={handleRefresh} color="primary" disabled={isLoading}>
              <RefreshIcon />
            </IconButton>
          </span>
        </Tooltip>
      </Box>

      {/* 统计卡片 */}
      <Box display="flex" flexWrap="wrap" gap={2} mb={3}>
        <Box flex="1 1 calc(50% - 8px)" minWidth="150px">
          <Paper sx={{ p: 2, textAlign: 'center' }}>
            <Typography variant="h4" color="primary">
              {statistics.total}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              总数
            </Typography>
          </Paper>
        </Box>
        <Box flex="1 1 calc(50% - 8px)" minWidth="150px">
          <Paper sx={{ p: 2, textAlign: 'center', bgcolor: alpha(theme.palette.success.main, 0.08) }}>
            <Typography variant="h4" sx={{ color: theme.palette.success.main }}>
              {statistics.buy}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              买入信号
            </Typography>
          </Paper>
        </Box>
        <Box flex="1 1 calc(50% - 8px)" minWidth="150px">
          <Paper sx={{ p: 2, textAlign: 'center', bgcolor: alpha(theme.palette.error.main, 0.08) }}>
            <Typography variant="h4" sx={{ color: theme.palette.error.main }}>
              {statistics.sell}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              卖出信号
            </Typography>
          </Paper>
        </Box>
        <Box flex="1 1 calc(50% - 8px)" minWidth="150px">
          <Paper sx={{ p: 2, textAlign: 'center', bgcolor: alpha(theme.palette.warning.main, 0.08) }}>
            <Typography variant="h4" sx={{ color: theme.palette.warning.main }}>
              {statistics.hold}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              持有
            </Typography>
          </Paper>
        </Box>
      </Box>

      {/* 计算状态提示 */}
      {data?.data?.calculating && (
        <Alert severity="info" sx={{ mb: 2 }}>
          信号正在计算中... (
          {data.data.calculation_status?.detail?.completed || 0}/
          {data.data.calculation_status?.detail?.total || 0})
        </Alert>
      )}

      {/* 加载状态 */}
      {isLoading && (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="300px">
          <CircularProgress />
        </Box>
      )}

      {/* 错误状态 */}
      {isError && (
        <Alert severity="error" sx={{ mb: 2 }}>
          加载信号数据失败: {error?.toString() || '未知错误'}
        </Alert>
      )}

      {/* 信号列表 */}
      {!isLoading && !isError && (
        <>
          {/* 标签页 */}
          <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
            <Tabs
              value={currentTab}
              onChange={(_, newValue) => setCurrentTab(newValue)}
              aria-label="信号类型筛选"
            >
              <Tab
                label={`全部 (${statistics.total})`}
                value="all"
                icon={<ShowChartIcon />}
                iconPosition="start"
              />
              <Tab
                label={`买入 (${statistics.buy})`}
                value="buy"
                icon={<TrendingUpIcon />}
                iconPosition="start"
                sx={{ 
                  color: currentTab === 'buy' ? theme.palette.success.main : 'inherit',
                  '&.Mui-selected': { color: theme.palette.success.main }
                }}
              />
              <Tab
                label={`卖出 (${statistics.sell})`}
                value="sell"
                icon={<TrendingDownIcon />}
                iconPosition="start"
                sx={{ 
                  color: currentTab === 'sell' ? theme.palette.error.main : 'inherit',
                  '&.Mui-selected': { color: theme.palette.error.main }
                }}
              />
              <Tab
                label={`持有 (${statistics.hold})`}
                value="hold"
                icon={<RemoveCircleOutlineIcon />}
                iconPosition="start"
                sx={{ 
                  color: currentTab === 'hold' ? theme.palette.warning.main : 'inherit',
                  '&.Mui-selected': { color: theme.palette.warning.main }
                }}
              />
            </Tabs>
          </Box>

          {/* 信号卡片列表 */}
          {filteredSignals.length > 0 ? (
            <Box>
              {filteredSignals.map((signal) => (
                <SignalCard key={signal.id} signal={signal} onViewStock={handleViewStock} />
              ))}
            </Box>
          ) : (
            <Paper sx={{ p: 4, textAlign: 'center' }}>
              <Typography variant="h6" color="text.secondary">
                暂无{currentTab === 'all' ? '' : currentTab === 'buy' ? '买入' : currentTab === 'sell' ? '卖出' : '持有'}信号数据
              </Typography>
              <Typography variant="body2" color="text.secondary" mt={1}>
                请先添加收藏股票，系统会自动计算买卖信号
              </Typography>
            </Paper>
          )}
        </>
      )}
    </Box>
  );
};

export default SignalsPage;

