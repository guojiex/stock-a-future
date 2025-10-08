/**
 * 技术指标展示组件 - 矩阵网格布局
 * 以卡片矩阵形式展示所有技术指标
 */

import React from 'react';
import {
  Box,
  Typography,
  Paper,
  Chip,
  Tooltip,
} from '@mui/material';
import {
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
const getSignalColor = (signal?: string): string => {
  if (!signal) return '#9e9e9e';
  const lower = signal.toLowerCase();
  if (lower === 'buy' || lower === 'oversold') return '#10b981';
  if (lower === 'sell' || lower === 'overbought') return '#ef4444';
  if (lower === 'hold' || lower === 'neutral') return '#f59e0b';
  return '#9e9e9e';
};

// 获取信号图标
const getSignalIcon = (signal?: string) => {
  if (!signal) return <RemoveIcon fontSize="small" />;
  const lower = signal.toLowerCase();
  if (lower === 'buy' || lower === 'oversold') return <TrendingUp fontSize="small" />;
  if (lower === 'sell' || lower === 'overbought') return <TrendingDown fontSize="small" />;
  return <RemoveIcon fontSize="small" />;
};

// 格式化日期显示 (YYYYMMDD 或 ISO 格式 -> YYYY-MM-DD)
const formatDateForDisplay = (dateStr: string): string => {
  if (!dateStr) return dateStr;
  
  // 如果是 ISO 格式 (包含 T 或 -)
  if (dateStr.includes('T') || dateStr.includes('-')) {
    const date = new Date(dateStr);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  }
  
  // 如果是 YYYYMMDD 格式
  if (dateStr.length === 8) {
    return `${dateStr.substring(0, 4)}-${dateStr.substring(4, 6)}-${dateStr.substring(6, 8)}`;
  }
  
  return dateStr;
};

// 指标卡片组件
interface IndicatorCardProps {
  title: string;
  tooltip: string;
  signal?: string;
  values: Array<{ label: string; value?: number; color?: string }>;
}

const IndicatorCard: React.FC<IndicatorCardProps> = ({ title, tooltip, signal, values }) => (
  <Paper 
    elevation={2} 
    sx={{ 
      p: 1.5,
      height: '100%',
      display: 'flex',
      flexDirection: 'column',
      transition: 'all 0.2s',
      '&:hover': {
        elevation: 4,
        transform: 'translateY(-2px)',
      }
    }}
  >
    {/* 标题和信号 */}
    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
      <Tooltip title={tooltip} arrow placement="top">
        <Typography 
          variant="subtitle2" 
          sx={{ 
            fontWeight: 'bold',
            cursor: 'help',
            fontSize: '0.85rem',
            flex: 1,
          }}
        >
          {title}
        </Typography>
      </Tooltip>
      {signal && (
        <Chip
          icon={getSignalIcon(signal)}
          label={getSignalText(signal)}
          size="small"
          sx={{
            height: 20,
            fontSize: '0.7rem',
            bgcolor: getSignalColor(signal),
            color: 'white',
            '& .MuiChip-icon': { fontSize: '0.9rem' },
          }}
        />
      )}
    </Box>

    {/* 指标数值 */}
    <Box sx={{ 
      display: 'grid', 
      gridTemplateColumns: values.length <= 2 ? `repeat(${values.length}, 1fr)` : 'repeat(3, 1fr)',
      gap: 0.5,
      flex: 1,
    }}>
      {values.map((item, idx) => (
        <Box 
          key={idx} 
          sx={{ 
            textAlign: 'center',
            p: 0.5,
            borderLeft: item.color ? `2px solid ${item.color}` : 'none',
            bgcolor: 'background.default',
            borderRadius: 0.5,
          }}
        >
          <Typography variant="caption" color="text.secondary" sx={{ fontSize: '0.65rem', display: 'block' }}>
            {item.label}
          </Typography>
          <Typography 
            variant="body2" 
            sx={{ 
              fontWeight: 'bold',
              fontSize: '0.8rem',
              color: item.color || 'text.primary',
            }}
          >
            {item.value !== undefined && item.value !== null ? item.value.toFixed(4) : 'N/A'}
          </Typography>
        </Box>
      ))}
    </Box>
  </Paper>
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
    <Box sx={{ p: 2 }}>
      {/* 标题 */}
      <Typography variant="h6" gutterBottom sx={{ mb: 2, fontWeight: 'bold' }}>
        📊 技术指标分析 ({formatDateForDisplay(data.trade_date)})
      </Typography>

      {/* 传统技术指标 */}
      <Typography variant="subtitle1" sx={{ mb: 1.5, color: 'primary.main', fontWeight: 'bold' }}>
        📈 传统技术指标
      </Typography>
      <Box sx={{ 
        display: 'grid', 
        gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
        gap: 1.5,
        mb: 3,
      }}>
        {/* 移动平均线 */}
        {data.ma && (
          <IndicatorCard
            title="移动平均线 (MA)"
            tooltip="移动平均线平滑价格波动，识别趋势方向。短期均线上穿长期均线为金叉买入信号，下穿为死叉卖出信号。"
            values={[
              { label: 'MA5', value: data.ma.ma5, color: '#2563eb' },
              { label: 'MA10', value: data.ma.ma10, color: '#f59e0b' },
              { label: 'MA20', value: data.ma.ma20, color: '#8b5cf6' },
              { label: 'MA60', value: data.ma.ma60, color: '#10b981' },
              { label: 'MA120', value: data.ma.ma120, color: '#ef4444' },
            ]}
          />
        )}

        {/* MACD */}
        {data.macd && (
          <IndicatorCard
            title="MACD指标"
            tooltip="MACD是趋势跟踪指标，通过快慢均线的差值判断买卖时机。DIF上穿DEA为金叉买入信号，下穿为死叉卖出信号。"
            signal={data.macd.signal}
            values={[
              { label: 'DIF', value: data.macd.dif },
              { label: 'DEA', value: data.macd.dea },
              { label: '柱状图', value: data.macd.histogram, color: data.macd.histogram > 0 ? '#10b981' : '#ef4444' },
            ]}
          />
        )}

        {/* RSI */}
        {data.rsi && (
          <IndicatorCard
            title="RSI 相对强弱指标"
            tooltip="RSI衡量价格变动的速度和幅度。RSI>70为超买区域，RSI<30为超卖区域，可作为反转信号参考。"
            signal={data.rsi.signal}
            values={[
              { label: 'RSI14', value: data.rsi.rsi14, color: data.rsi.rsi14 > 70 ? '#ef4444' : data.rsi.rsi14 < 30 ? '#10b981' : undefined },
            ]}
          />
        )}

        {/* 布林带 */}
        {data.boll && (
          <IndicatorCard
            title="布林带 (BOLL)"
            tooltip="布林带由移动平均线和标准差构成，价格触及上轨可能回调，触及下轨可能反弹。带宽收窄预示突破，扩张表示趋势延续。"
            signal={data.boll.signal}
            values={[
              { label: '上轨', value: data.boll.upper, color: '#ef4444' },
              { label: '中轨', value: data.boll.middle },
              { label: '下轨', value: data.boll.lower, color: '#10b981' },
            ]}
          />
        )}

        {/* KDJ */}
        {data.kdj && (
          <IndicatorCard
            title="KDJ 随机指标"
            tooltip="KDJ随机指标反映价格在一定周期内的相对位置。K>80且D>80为超买，K<20且D<20为超卖。J值更敏感，可提前预警。"
            signal={data.kdj.signal}
            values={[
              { label: 'K值', value: data.kdj.k },
              { label: 'D值', value: data.kdj.d },
              { label: 'J值', value: data.kdj.j },
            ]}
          />
        )}
      </Box>

      {/* 动量因子指标 */}
      <Typography variant="subtitle1" sx={{ mb: 1.5, color: 'primary.main', fontWeight: 'bold' }}>
        🔄 动量因子指标
      </Typography>
      <Box sx={{ 
        display: 'grid', 
        gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
        gap: 1.5,
        mb: 3,
      }}>
        {/* 威廉指标 */}
        {data.wr && (
          <IndicatorCard
            title="威廉指标 (%R)"
            tooltip="威廉指标衡量股价在一定周期内的相对位置。WR>-20为超买区域，WR<-80为超卖区域。数值越接近0越超买，越接近-100越超卖。"
            signal={data.wr.signal}
            values={[
              { label: 'WR14', value: data.wr.wr14 },
            ]}
          />
        )}

        {/* 动量指标 */}
        {data.momentum && (
          <IndicatorCard
            title="动量指标 (Momentum)"
            tooltip="动量指标衡量价格变化的速度。正值表示上涨动量，负值表示下跌动量。数值越大表示趋势越强劲。"
            signal={data.momentum.signal}
            values={[
              { label: 'Momentum10', value: data.momentum.momentum10 },
              { label: 'Momentum20', value: data.momentum.momentum20 },
            ]}
          />
        )}

        {/* 变化率指标 */}
        {data.roc && (
          <IndicatorCard
            title="变化率指标 (ROC)"
            tooltip="ROC变化率指标衡量价格在一定周期内的百分比变化。正值表示上涨，负值表示下跌。可用于判断趋势强度和转折点。"
            signal={data.roc.signal}
            values={[
              { label: 'ROC10', value: data.roc.roc10 },
              { label: 'ROC20', value: data.roc.roc20 },
            ]}
          />
        )}
      </Box>

      {/* 趋势因子指标 */}
      <Typography variant="subtitle1" sx={{ mb: 1.5, color: 'primary.main', fontWeight: 'bold' }}>
        📊 趋势因子指标
      </Typography>
      <Box sx={{ 
        display: 'grid', 
        gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
        gap: 1.5,
        mb: 3,
      }}>
        {/* ADX */}
        {data.adx && (
          <IndicatorCard
            title="平均方向指数 (ADX)"
            tooltip="ADX衡量趋势强度，不判断方向。ADX>25表示强趋势，<20表示弱趋势。PDI>MDI表示上升趋势，反之为下降趋势。"
            signal={data.adx.signal}
            values={[
              { label: 'ADX', value: data.adx.adx },
              { label: 'PDI', value: data.adx.pdi, color: '#10b981' },
              { label: 'MDI', value: data.adx.mdi, color: '#ef4444' },
            ]}
          />
        )}

        {/* SAR */}
        {data.sar && (
          <IndicatorCard
            title="抛物线转向 (SAR)"
            tooltip="SAR抛物线转向指标用于确定止损点和趋势转换。价格在SAR之上为上升趋势，之下为下降趋势。SAR点位可作为止损参考。"
            signal={data.sar.signal}
            values={[
              { label: 'SAR', value: data.sar.sar },
            ]}
          />
        )}

        {/* 一目均衡表 */}
        {data.ichimoku && (
          <IndicatorCard
            title="一目均衡表 (Ichimoku)"
            tooltip="一目均衡表是综合性技术指标。转换线上穿基准线为买入信号，价格突破云带(先行带)确认趋势。云带厚度反映支撑阻力强度。"
            signal={data.ichimoku.signal}
            values={[
              { label: '转换线', value: data.ichimoku.tenkan_sen },
              { label: '基准线', value: data.ichimoku.kijun_sen },
              { label: '先行带A', value: data.ichimoku.senkou_span_a },
              { label: '先行带B', value: data.ichimoku.senkou_span_b },
            ]}
          />
        )}
      </Box>

      {/* 波动率因子指标 */}
      <Typography variant="subtitle1" sx={{ mb: 1.5, color: 'primary.main', fontWeight: 'bold' }}>
        📉 波动率因子指标
      </Typography>
      <Box sx={{ 
        display: 'grid', 
        gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
        gap: 1.5,
        mb: 3,
      }}>
        {/* ATR */}
        {data.atr && (
          <IndicatorCard
            title="平均真实范围 (ATR)"
            tooltip="ATR衡量价格波动幅度，数值越大表示波动越剧烈。可用于设置止损位和判断市场活跃度。高ATR适合趋势交易，低ATR适合区间交易。"
            signal={data.atr.signal}
            values={[
              { label: 'ATR14', value: data.atr.atr14 },
            ]}
          />
        )}

        {/* 标准差 */}
        {data.stddev && (
          <IndicatorCard
            title="标准差 (StdDev)"
            tooltip="标准差衡量价格偏离平均值的程度。数值越大表示价格波动越不稳定。可用于评估投资风险和市场不确定性。"
            signal={data.stddev.signal}
            values={[
              { label: 'StdDev20', value: data.stddev.stddev20 },
            ]}
          />
        )}

        {/* 历史波动率 */}
        {data.hv && (
          <IndicatorCard
            title="历史波动率 (HV)"
            tooltip="历史波动率反映股价在过去一段时间的波动程度。高波动率意味着高风险高收益，低波动率表示价格相对稳定。"
            signal={data.hv.signal}
            values={[
              { label: 'HV20', value: data.hv.hv20 },
              { label: 'HV60', value: data.hv.hv60 },
            ]}
          />
        )}
      </Box>

      {/* 成交量因子指标 */}
      <Typography variant="subtitle1" sx={{ mb: 1.5, color: 'primary.main', fontWeight: 'bold' }}>
        📊 成交量因子指标
      </Typography>
      <Box sx={{ 
        display: 'grid', 
        gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
        gap: 1.5,
        mb: 3,
      }}>
        {/* VWAP */}
        {data.vwap && (
          <IndicatorCard
            title="成交量加权平均价 (VWAP)"
            tooltip="VWAP是成交量加权的平均价格，反映真实的平均成交成本。价格高于VWAP表示买盘强劲，低于VWAP表示卖盘压力大。"
            signal={data.vwap.signal}
            values={[
              { label: 'VWAP', value: data.vwap.vwap },
            ]}
          />
        )}

        {/* A/D线 */}
        {data.ad_line && (
          <IndicatorCard
            title="累积/派发线 (A/D Line)"
            tooltip="A/D线结合价格和成交量，衡量资金流向。上升表示资金流入(累积)，下降表示资金流出(派发)。可用于确认价格趋势。"
            signal={data.ad_line.signal}
            values={[
              { label: 'A/D Line', value: data.ad_line.ad_line },
            ]}
          />
        )}

        {/* EMV */}
        {data.emv && (
          <IndicatorCard
            title="简易波动指标 (EMV)"
            tooltip="EMV衡量价格变动的难易程度。正值表示价格上涨容易，负值表示下跌容易。结合成交量分析，判断价格变动的可持续性。"
            signal={data.emv.signal}
            values={[
              { label: 'EMV14', value: data.emv.emv14 },
            ]}
          />
        )}

        {/* VPT */}
        {data.vpt && (
          <IndicatorCard
            title="量价确认指标 (VPT)"
            tooltip="VPT将成交量与价格变化相结合，确认价格趋势。VPT与价格同向运动确认趋势，背离时可能预示反转。"
            signal={data.vpt.signal}
            values={[
              { label: 'VPT', value: data.vpt.vpt },
            ]}
          />
        )}
      </Box>

      {/* 免责声明 */}
      <Paper sx={{ p: 1.5, mt: 2, bgcolor: 'warning.light', color: 'warning.contrastText' }}>
        <Typography variant="caption" display="block">
          <strong>📢 免责声明：</strong>
          本系统提供的技术指标分析和预测结果仅供参考，不构成投资建议。投资有风险，入市需谨慎。请在专业人士指导下进行投资决策。
        </Typography>
      </Paper>
    </Box>
  );
};

export default TechnicalIndicatorsView;