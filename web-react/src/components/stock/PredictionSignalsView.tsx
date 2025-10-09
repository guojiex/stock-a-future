/**
 * 买卖预测信号展示组件
 * 使用 Material-UI 展示预测数据
 */

import React, { useState } from 'react';
import {
  Box,
  Typography,
  Alert,
  AlertTitle,
  Paper,
  Chip,
  CircularProgress,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Grid,
  Divider,
} from '@mui/material';
import {
  TrendingUp as BuyIcon,
  TrendingDown as SellIcon,
  ExpandMore as ExpandMoreIcon,
  Info as InfoIcon,
} from '@mui/icons-material';
import { useGetPredictionsQuery } from '../../services/api';
import { TradingPointPrediction } from '../../types/stock';

interface PredictionSignalsViewProps {
  stockCode: string;
}

// 从预测理由中提取强度等级
const extractStrengthFromReason = (reason: string): 'STRONG' | 'MEDIUM' | 'WEAK' => {
  if (reason.includes('强度：STRONG')) {
    return 'STRONG';
  } else if (reason.includes('强度：MEDIUM')) {
    return 'MEDIUM';
  } else if (reason.includes('强度：WEAK')) {
    return 'WEAK';
  }
  return 'WEAK';
};

// 格式化日期显示 (YYYYMMDD 或 ISO 格式 -> YYYY-MM-DD)
const formatDateForDisplay = (dateStr: string): string => {
  if (!dateStr) return dateStr;
  
  // 如果是 ISO 格式 (包含 T 或已经有 -)
  if (dateStr.includes('T') || (dateStr.includes('-') && dateStr.length > 10)) {
    const date = new Date(dateStr);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  }
  
  // 如果是 YYYYMMDD 格式
  if (dateStr.length === 8 && !dateStr.includes('-')) {
    return `${dateStr.substring(0, 4)}-${dateStr.substring(4, 6)}-${dateStr.substring(6, 8)}`;
  }
  
  return dateStr;
};

// 预测项组件
const PredictionItem: React.FC<{
  prediction: TradingPointPrediction;
  index: number;
}> = ({ prediction, index }) => {
  const isBuy = prediction.type === 'BUY';
  const icon = isBuy ? <BuyIcon sx={{ fontSize: 40 }} /> : <SellIcon sx={{ fontSize: 40 }} />;
  const typeText = isBuy ? '买入' : '卖出';
  const strength = extractStrengthFromReason(prediction.reason);
  const isWeak = strength === 'WEAK';

  const [expanded, setExpanded] = useState(!isWeak);

  const getStrengthColor = () => {
    if (strength === 'STRONG') return 'success';
    if (strength === 'MEDIUM') return 'warning';
    return 'default';
  };

  return (
    <Accordion 
      expanded={expanded} 
      onChange={() => setExpanded(!expanded)}
      sx={{ mb: 2 }}
    >
      <AccordionSummary expandIcon={<ExpandMoreIcon />}>
        <Box display="flex" alignItems="center" width="100%" gap={2}>
          <Box color={isBuy ? 'success.main' : 'error.main'}>
            {icon}
          </Box>
          <Box flexGrow={1}>
            <Typography variant="h6" color={isBuy ? 'success.main' : 'error.main'}>
              {typeText}信号
            </Typography>
            <Box display="flex" gap={1} mt={0.5} flexWrap="wrap">
              <Chip 
                label={`¥${prediction.price.toFixed(2)}`} 
                size="small" 
                variant="outlined"
              />
              <Chip 
                label={`📅 ${formatDateForDisplay(prediction.signal_date)}`}
                size="small"
                variant="outlined"
              />
              <Chip 
                label={`概率: ${(prediction.probability * 100).toFixed(1)}%`}
                size="small"
                color={isBuy ? 'success' : 'error'}
              />
            </Box>
          </Box>
          <Chip 
            label={strength} 
            color={getStrengthColor() as any}
            size="small"
          />
        </Box>
      </AccordionSummary>
      <AccordionDetails>
        <Box>
          {/* 预测理由 */}
          <Box mb={2}>
            <Typography variant="subtitle2" color="text.secondary" gutterBottom>
              预测理由
            </Typography>
            <Typography variant="body2">
              {prediction.reason || '基于技术指标分析'}
            </Typography>
          </Box>

          {/* 回测结果 */}
          {prediction.backtested && (
            <Paper variant="outlined" sx={{ p: 2, bgcolor: 'background.default' }}>
              <Typography variant="subtitle2" gutterBottom>
                回测结果
              </Typography>
              <Grid container spacing={2}>
                <Grid size={{ xs: 12, sm: 4 }}>
                  <Typography variant="body2" color="text.secondary">
                    结果
                  </Typography>
                  <Chip
                    label={prediction.is_correct ? '✅ 正确' : '❌ 错误'}
                    color={prediction.is_correct ? 'success' : 'error'}
                    size="small"
                  />
                </Grid>
                <Grid size={{ xs: 12, sm: 4 }}>
                  <Typography variant="body2" color="text.secondary">
                    次日价格
                  </Typography>
                  <Typography variant="body1" fontWeight="bold">
                    ¥{prediction.next_day_price?.toFixed(2) || 'N/A'}
                  </Typography>
                </Grid>
                <Grid size={{ xs: 12, sm: 4 }}>
                  <Typography variant="body2" color="text.secondary">
                    价差
                  </Typography>
                  <Typography 
                    variant="body1" 
                    fontWeight="bold"
                    color={(prediction.price_diff || 0) >= 0 ? 'success.main' : 'error.main'}
                  >
                    {(prediction.price_diff || 0) >= 0 ? '+' : ''}
                    {prediction.price_diff?.toFixed(2) || '0.00'} 
                    ({(prediction.price_diff_ratio || 0) >= 0 ? '+' : ''}
                    {prediction.price_diff_ratio?.toFixed(2) || '0.00'}%)
                  </Typography>
                </Grid>
              </Grid>
            </Paper>
          )}

          {/* 相关指标 */}
          {prediction.indicators && prediction.indicators.length > 0 && (
            <Box mt={2}>
              <Box display="flex" alignItems="center" gap={1} mb={1}>
                <Typography variant="subtitle2" color="text.secondary">
                  {prediction.indicators.length > 1 ? '🔗 综合信号' : '📊 相关指标'}
                </Typography>
                {prediction.indicators.length > 1 && (
                  <Chip 
                    label={`✨ ${prediction.indicators.length}个指标共识`}
                    size="small"
                    sx={{
                      background: 'linear-gradient(135deg, #f59e0b, #d97706)',
                      color: 'white',
                      fontWeight: 600,
                      animation: 'pulse 2s ease-in-out infinite',
                      '@keyframes pulse': {
                        '0%, 100%': { opacity: 1 },
                        '50%': { opacity: 0.85 }
                      }
                    }}
                  />
                )}
              </Box>
              <Box display="flex" gap={1} flexWrap="wrap">
                {prediction.indicators.map((indicator, idx) => (
                  <Chip 
                    key={idx} 
                    label={indicator} 
                    size="small" 
                    variant={prediction.indicators.length > 1 ? "filled" : "outlined"}
                    color={prediction.indicators.length > 1 ? "primary" : "default"}
                    sx={prediction.indicators.length > 1 ? {
                      fontWeight: 500,
                      boxShadow: 1
                    } : undefined}
                  />
                ))}
              </Box>
              {prediction.indicators.length > 1 && (
                <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block', fontStyle: 'italic' }}>
                  💡 多个技术指标共识，置信度已提升
                </Typography>
              )}
            </Box>
          )}
        </Box>
      </AccordionDetails>
    </Accordion>
  );
};

const PredictionSignalsView: React.FC<PredictionSignalsViewProps> = ({ stockCode }) => {
  const { data, isLoading, error } = useGetPredictionsQuery(stockCode, {
    skip: !stockCode,
  });

  // 加载状态
  if (isLoading) {
    return (
      <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" py={8}>
        <CircularProgress size={60} />
        <Typography variant="body1" sx={{ mt: 2 }}>
          正在加载预测数据...
        </Typography>
      </Box>
    );
  }

  // 错误状态
  if (error) {
    return (
      <Alert severity="error">
        <AlertTitle>加载失败</AlertTitle>
        {error.toString()}
      </Alert>
    );
  }

  // 无数据
  if (!data?.data) {
    return (
      <Alert severity="info">
        <AlertTitle>暂无数据</AlertTitle>
        暂无预测数据
      </Alert>
    );
  }

  const predictionResult = data.data;
  const hasPredictions = predictionResult.predictions && predictionResult.predictions.length > 0;

  return (
    <Box>
      {/* 预测概览 */}
      <Paper elevation={2} sx={{ p: 3, mb: 3 }}>
        <Grid container spacing={3}>
          <Grid size={{ xs: 12, sm: 4 }}>
            <Box textAlign="center">
              <Typography variant="h4" color="primary.main" fontWeight="bold">
                {predictionResult.predictions?.length || 0}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                预测数量
              </Typography>
            </Box>
          </Grid>
          <Grid size={{ xs: 12, sm: 4 }}>
            <Box textAlign="center">
              <Typography variant="h4" color="secondary.main" fontWeight="bold">
                {(predictionResult.confidence * 100).toFixed(1)}%
              </Typography>
              <Typography variant="body2" color="text.secondary">
                整体置信度
              </Typography>
            </Box>
          </Grid>
          <Grid size={{ xs: 12, sm: 4 }}>
            <Box textAlign="center">
              <Typography variant="h6" fontWeight="bold">
                {formatDateForDisplay(predictionResult.trade_date)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                分析日期
              </Typography>
            </Box>
          </Grid>
        </Grid>
        <Divider sx={{ my: 2 }} />
        <Typography variant="caption" color="text.secondary" display="block">
          更新于: {new Date(predictionResult.updated_at).toLocaleDateString('zh-CN')}
        </Typography>
      </Paper>

      {/* 预测列表 */}
      {hasPredictions ? (
        <Box>
          <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
            <Typography variant="h6">买卖信号列表</Typography>
            <Typography variant="caption" color="text.secondary">
              点击展开/收起详情
            </Typography>
          </Box>
          
          {predictionResult.predictions.map((prediction, index) => (
            <PredictionItem
              key={index}
              prediction={prediction}
              index={index}
            />
          ))}
        </Box>
      ) : (
        <Paper elevation={2} sx={{ p: 4, textAlign: 'center' }}>
          <Box fontSize={60} mb={2}>📊</Box>
          <Typography variant="h6" gutterBottom>
            暂无明确信号
          </Typography>
          <Typography variant="body2" color="text.secondary">
            当前市场条件下，没有明确的买卖信号
          </Typography>
        </Paper>
      )}

      {/* 免责声明 */}
      <Alert severity="warning" sx={{ mt: 3 }} icon={<InfoIcon />}>
        <AlertTitle>⚠️ 免责声明</AlertTitle>
        本系统提供的预测结果仅供参考，不构成投资建议。投资有风险，入市需谨慎。
      </Alert>
    </Box>
  );
};

export default PredictionSignalsView;

