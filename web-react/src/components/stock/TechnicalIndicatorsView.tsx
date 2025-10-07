/**
 * 技术指标展示组件 - 全面重构版
 * 参考原有web版本的实现，支持18种核心技术指标
 */

import React from 'react';
import {
  Box,
  Typography,
  Paper,
  Chip,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Tooltip,
  IconButton,
} from '@mui/material';
import {
  ExpandMore as ExpandMoreIcon,
  Info as InfoIcon,
  TrendingUp,
  TrendingDown,
  Remove as RemoveIcon,
} from '@mui/icons-material';
import { TechnicalIndicators } from '../../types/stock';

interface TechnicalIndicatorsViewProps {
  data: TechnicalIndicators;
}

// 信号类型到中文的映射
const getSignalText = (signal?: string): string => {
  if (!signal) return '中性';
  const signalMap: Record<string, string> = {
    'buy': '买入',
    'sell': '卖出',
    'hold': '持有',
    'neutral': '中性',
    'overbought': '超买',
    'oversold': '超卖',
  };
  return signalMap[signal.toLowerCase()] || signal;
};

// 获取信号颜色
const getSignalColor = (signal?: string): 'success' | 'error' | 'warning' | 'default' => {
  if (!signal) return 'default';
  const lower = signal.toLowerCase();
  if (lower === 'buy' || lower === 'oversold') return 'success';
  if (lower === 'sell' || lower === 'overbought') return 'error';
  if (lower === 'hold' || lower === 'neutral') return 'warning';
  return 'default';
};

// 获取信号图标
const getSignalIcon = (signal?: string) => {
  if (!signal) return <RemoveIcon fontSize="small" />;
  const lower = signal.toLowerCase();
  if (lower === 'buy' || lower === 'oversold') return <TrendingUp fontSize="small" />;
  if (lower === 'sell' || lower === 'overbought') return <TrendingDown fontSize="small" />;
  return <RemoveIcon fontSize="small" />;
};

// 指标卡片组件
const IndicatorCard = ({ 
  label, 
  value, 
  color 
}: { 
  label: string; 
  value?: number; 
  color?: string;
}) => (
  <Box sx={{ textAlign: 'center', p: 1 }}>
    <Typography variant="caption" color="text.secondary" display="block">
      {label}
    </Typography>
    <Typography 
      variant="h6" 
      sx={{ 
        color: color || 'text.primary',
        fontWeight: 'bold',
        fontSize: '1.1rem'
      }}
    >
      {value !== undefined && value !== null ? value.toFixed(2) : 'N/A'}
    </Typography>
  </Box>
);

// 指标分组组件
const IndicatorSection = ({
  title,
  tooltip,
  signal,
  children,
  defaultExpanded = false,
}: {
  title: string;
  tooltip: string;
  signal?: string;
  children: React.ReactNode;
  defaultExpanded?: boolean;
}) => (
  <Accordion defaultExpanded={defaultExpanded}>
    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
      <Box display="flex" alignItems="center" width="100%">
        <Typography variant="subtitle1" sx={{ flexGrow: 1, fontWeight: 'bold' }}>
          {title}
        </Typography>
        {signal && (
          <Chip
            icon={getSignalIcon(signal)}
            label={getSignalText(signal)}
            color={getSignalColor(signal)}
            size="small"
            sx={{ mr: 2 }}
          />
        )}
        <Tooltip title={tooltip}>
          <IconButton size="small">
            <InfoIcon fontSize="small" />
          </IconButton>
        </Tooltip>
      </Box>
    </AccordionSummary>
    <AccordionDetails>
      {children}
    </AccordionDetails>
  </Accordion>
);

const TechnicalIndicatorsView: React.FC<TechnicalIndicatorsViewProps> = ({ data }) => {
  if (!data) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
        <Typography variant="body1" color="text.secondary">
          暂无技术指标数据
        </Typography>
      </Box>
    );
  }

  return (
    <Box>
      {/* 标题 */}
      <Typography variant="h6" gutterBottom sx={{ mb: 3 }}>
        技术指标分析 ({data.trade_date})
      </Typography>

      {/* 传统技术指标 */}
      <Typography variant="h6" gutterBottom sx={{ mt: 3, mb: 2, color: 'primary.main' }}>
        📈 传统技术指标
        </Typography>

      {/* 移动平均线 */}
      {data.ma && (
        <IndicatorSection
          title="移动平均线 (MA)"
          tooltip="移动平均线平滑价格波动，识别趋势方向。短期均线上穿长期均线为金叉买入信号，下穿为死叉卖出信号。"
          defaultExpanded
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: 'repeat(2, 1fr)', sm: 'repeat(5, 1fr)' }, gap: 2 }}>
            <IndicatorCard label="MA5" value={data.ma.ma5} color="#2563eb" />
            <IndicatorCard label="MA10" value={data.ma.ma10} color="#f59e0b" />
            <IndicatorCard label="MA20" value={data.ma.ma20} color="#8b5cf6" />
            <IndicatorCard label="MA60" value={data.ma.ma60} color="#10b981" />
            <IndicatorCard label="MA120" value={data.ma.ma120} color="#ef4444" />
          </Box>
        </IndicatorSection>
      )}

          {/* MACD */}
      {data.macd && (
        <IndicatorSection
          title="MACD指标"
          tooltip="MACD是趋势跟踪指标，通过快慢均线的差值判断买卖时机。DIF上穿DEA为金叉买入信号，下穿为死叉卖出信号。"
          signal={data.macd.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 2 }}>
            <IndicatorCard label="DIF" value={data.macd.dif} />
            <IndicatorCard label="DEA" value={data.macd.dea} />
              <IndicatorCard
              label="柱状图" 
              value={data.macd.histogram}
              color={data.macd.histogram > 0 ? '#10b981' : '#ef4444'}
            />
          </Box>
        </IndicatorSection>
      )}

          {/* RSI */}
      {data.rsi && (
        <IndicatorSection
          title="RSI 相对强弱指标"
          tooltip="RSI衡量价格变动的速度和幅度。RSI>70为超买区域，RSI<30为超卖区域，可作为反转信号参考。"
          signal={data.rsi.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr', gap: 2 }}>
            <IndicatorCard 
              label="RSI14" 
              value={data.rsi.rsi14}
              color={data.rsi.rsi14 > 70 ? '#ef4444' : data.rsi.rsi14 < 30 ? '#10b981' : undefined}
            />
          </Box>
        </IndicatorSection>
      )}

          {/* 布林带 */}
      {data.boll && (
        <IndicatorSection
          title="布林带 (BOLL)"
          tooltip="布林带由移动平均线和标准差构成，价格触及上轨可能回调，触及下轨可能反弹。带宽收窄预示突破，扩张表示趋势延续。"
          signal={data.boll.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 2 }}>
            <IndicatorCard label="上轨" value={data.boll.upper} color="#ef4444" />
            <IndicatorCard label="中轨" value={data.boll.middle} />
            <IndicatorCard label="下轨" value={data.boll.lower} color="#10b981" />
          </Box>
        </IndicatorSection>
      )}

          {/* KDJ */}
      {data.kdj && (
        <IndicatorSection
          title="KDJ 随机指标"
          tooltip="KDJ随机指标反映价格在一定周期内的相对位置。K>80且D>80为超买，K<20且D<20为超卖。J值更敏感，可提前预警。"
          signal={data.kdj.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 2 }}>
            <IndicatorCard label="K值" value={data.kdj.k} />
            <IndicatorCard label="D值" value={data.kdj.d} />
            <IndicatorCard label="J值" value={data.kdj.j} />
          </Box>
        </IndicatorSection>
      )}

      {/* 动量因子指标 */}
      <Typography variant="h6" gutterBottom sx={{ mt: 4, mb: 2, color: 'primary.main' }}>
        ⚡ 动量因子指标
            </Typography>

      {/* 威廉指标 */}
      {data.wr && (
        <IndicatorSection
          title="威廉指标 (%R)"
          tooltip="威廉指标衡量股价在一定周期内的相对位置。WR>-20为超买区域，WR<-80为超卖区域。数值越接近0越超买，越接近-100越超卖。"
          signal={data.wr.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr', gap: 2 }}>
            <IndicatorCard label="WR14" value={data.wr.wr14} />
          </Box>
        </IndicatorSection>
      )}

      {/* 动量指标 */}
      {data.momentum && (
        <IndicatorSection
          title="动量指标 (Momentum)"
          tooltip="动量指标衡量价格变化的速度。正值表示上涨动量，负值表示下跌动量。数值越大表示趋势越强劲。"
          signal={data.momentum.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 2 }}>
            <IndicatorCard label="Momentum10" value={data.momentum.momentum10} />
            <IndicatorCard label="Momentum20" value={data.momentum.momentum20} />
          </Box>
        </IndicatorSection>
      )}

      {/* 变化率指标 */}
      {data.roc && (
        <IndicatorSection
          title="变化率指标 (ROC)"
          tooltip="ROC变化率指标衡量价格在一定周期内的百分比变化。正值表示上涨，负值表示下跌。可用于判断趋势强度和转折点。"
          signal={data.roc.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 2 }}>
            <IndicatorCard label="ROC10" value={data.roc.roc10} />
            <IndicatorCard label="ROC20" value={data.roc.roc20} />
          </Box>
        </IndicatorSection>
      )}

      {/* 趋势因子指标 */}
      <Typography variant="h6" gutterBottom sx={{ mt: 4, mb: 2, color: 'primary.main' }}>
        📊 趋势因子指标
      </Typography>

      {/* ADX */}
      {data.adx && (
        <IndicatorSection
          title="平均方向指数 (ADX)"
          tooltip="ADX衡量趋势强度，不判断方向。ADX>25表示强趋势，<20表示弱趋势。PDI>MDI表示上升趋势，反之为下降趋势。"
          signal={data.adx.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 2 }}>
            <IndicatorCard label="ADX" value={data.adx.adx} />
            <IndicatorCard label="PDI" value={data.adx.pdi} color="#10b981" />
            <IndicatorCard label="MDI" value={data.adx.mdi} color="#ef4444" />
          </Box>
        </IndicatorSection>
      )}

      {/* SAR */}
      {data.sar && (
        <IndicatorSection
          title="抛物线转向 (SAR)"
          tooltip="SAR抛物线转向指标用于确定止损点和趋势转换。价格在SAR之上为上升趋势，之下为下降趋势。SAR点位可作为止损参考。"
          signal={data.sar.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr', gap: 2 }}>
            <IndicatorCard label="SAR" value={data.sar.sar} />
          </Box>
        </IndicatorSection>
      )}

      {/* 一目均衡表 */}
      {data.ichimoku && (
        <IndicatorSection
          title="一目均衡表 (Ichimoku)"
          tooltip="一目均衡表是综合性技术指标。转换线上穿基准线为买入信号，价格突破云带(先行带)确认趋势。云带厚度反映支撑阻力强度。"
          signal={data.ichimoku.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: 'repeat(2, 1fr)', sm: 'repeat(4, 1fr)' }, gap: 2 }}>
            <IndicatorCard label="转换线" value={data.ichimoku.tenkan_sen} />
            <IndicatorCard label="基准线" value={data.ichimoku.kijun_sen} />
            <IndicatorCard label="先行带A" value={data.ichimoku.senkou_span_a} />
            <IndicatorCard label="先行带B" value={data.ichimoku.senkou_span_b} />
          </Box>
        </IndicatorSection>
      )}

      {/* 波动率因子指标 */}
      <Typography variant="h6" gutterBottom sx={{ mt: 4, mb: 2, color: 'primary.main' }}>
        📊 波动率因子指标
      </Typography>

      {/* ATR */}
      {data.atr && (
        <IndicatorSection
          title="平均真实范围 (ATR)"
          tooltip="ATR衡量价格波动幅度，数值越大表示波动越剧烈。可用于设置止损位和判断市场活跃度。高ATR适合趋势交易，低ATR适合区间交易。"
          signal={data.atr.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr', gap: 2 }}>
            <IndicatorCard label="ATR14" value={data.atr.atr14} />
          </Box>
        </IndicatorSection>
      )}

      {/* 标准差 */}
      {data.stddev && (
        <IndicatorSection
          title="标准差 (StdDev)"
          tooltip="标准差衡量价格偏离平均值的程度。数值越大表示价格波动越不稳定。可用于评估投资风险和市场不确定性。"
          signal={data.stddev.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr', gap: 2 }}>
            <IndicatorCard label="StdDev20" value={data.stddev.stddev20} />
          </Box>
        </IndicatorSection>
      )}

      {/* 历史波动率 */}
      {data.hv && (
        <IndicatorSection
          title="历史波动率 (HV)"
          tooltip="历史波动率反映股价在过去一段时间的波动程度。高波动率意味着高风险高收益，低波动率表示价格相对稳定。"
          signal={data.hv.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 2 }}>
            <IndicatorCard label="HV20" value={data.hv.hv20} />
            <IndicatorCard label="HV60" value={data.hv.hv60} />
          </Box>
        </IndicatorSection>
      )}

      {/* 成交量因子指标 */}
      <Typography variant="h6" gutterBottom sx={{ mt: 4, mb: 2, color: 'primary.main' }}>
        📊 成交量因子指标
            </Typography>

      {/* VWAP */}
      {data.vwap && (
        <IndicatorSection
          title="成交量加权平均价 (VWAP)"
          tooltip="VWAP是成交量加权的平均价格，反映真实的平均成交成本。价格高于VWAP表示买盘强劲，低于VWAP表示卖盘压力大。"
          signal={data.vwap.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr', gap: 2 }}>
            <IndicatorCard label="VWAP" value={data.vwap.vwap} />
          </Box>
        </IndicatorSection>
      )}

      {/* A/D线 */}
      {data.ad_line && (
        <IndicatorSection
          title="累积/派发线 (A/D Line)"
          tooltip="A/D线结合价格和成交量，衡量资金流向。上升表示资金流入(累积)，下降表示资金流出(派发)。可用于确认价格趋势。"
          signal={data.ad_line.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr', gap: 2 }}>
            <IndicatorCard label="A/D Line" value={data.ad_line.ad_line} />
          </Box>
        </IndicatorSection>
      )}

      {/* EMV */}
      {data.emv && (
        <IndicatorSection
          title="简易波动指标 (EMV)"
          tooltip="EMV衡量价格变动的难易程度。正值表示价格上涨容易，负值表示下跌容易。结合成交量分析，判断价格变动的可持续性。"
          signal={data.emv.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr', gap: 2 }}>
            <IndicatorCard label="EMV14" value={data.emv.emv14} />
          </Box>
        </IndicatorSection>
      )}

      {/* VPT */}
      {data.vpt && (
        <IndicatorSection
          title="量价确认指标 (VPT)"
          tooltip="VPT将成交量与价格变化相结合，确认价格趋势。VPT与价格同向运动确认趋势，背离时可能预示反转。"
          signal={data.vpt.signal}
        >
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr', gap: 2 }}>
            <IndicatorCard label="VPT" value={data.vpt.vpt} />
          </Box>
        </IndicatorSection>
      )}

      {/* 说明文字 */}
      <Paper sx={{ p: 2, mt: 3, bgcolor: 'warning.light', color: 'warning.contrastText' }}>
        <Typography variant="body2">
          <strong>📢 免责声明：</strong>
          本系统提供的技术指标分析和预测结果仅供参考，不构成投资建议。投资有风险，入市需谨慎。请在专业人士指导下进行投资决策。
        </Typography>
      </Paper>
    </Box>
  );
};

export default TechnicalIndicatorsView;

