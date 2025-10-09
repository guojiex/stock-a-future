/**
 * 买卖预测视图组件
 * 展示基于技术指标的买卖点预测
 */

import React, { useState } from 'react';
import { PredictionResult, TradingPointPrediction } from '../../types/stock';

interface PredictionsViewProps {
  data: PredictionResult | null;
  isLoading?: boolean;
  onDateClick?: (date: string) => void;
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
  return 'WEAK'; // 默认值
};

// 格式化日期显示 (YYYYMMDD -> YYYY-MM-DD)
const formatDateForDisplay = (dateStr: string): string => {
  if (!dateStr || dateStr.length !== 8) return dateStr;
  return `${dateStr.substring(0, 4)}-${dateStr.substring(4, 6)}-${dateStr.substring(6, 8)}`;
};

// 预测项组件
const PredictionItem: React.FC<{
  prediction: TradingPointPrediction;
  index: number;
  onDateClick?: (date: string) => void;
}> = ({ prediction, index, onDateClick }) => {
  const typeClass = prediction.type.toLowerCase();
  const icon = prediction.type === 'BUY' ? '📈' : '📉';
  const typeText = prediction.type === 'BUY' ? '买入' : '卖出';
  const strength = extractStrengthFromReason(prediction.reason);
  const isWeak = strength === 'WEAK';
  
  const [isCollapsed, setIsCollapsed] = useState(isWeak);

  const handleDateClick = (e: React.MouseEvent) => {
    e.preventDefault();
    if (onDateClick && prediction.signal_date) {
      onDateClick(prediction.signal_date);
    }
  };

  return (
    <div 
      className={`card bg-base-100 shadow-md hover:shadow-lg transition-all duration-300 ${
        isCollapsed ? 'opacity-75' : ''
      }`}
      data-index={index}
    >
      {/* 预测头部 */}
      <div 
        className="card-body p-4 cursor-pointer"
        onClick={() => setIsCollapsed(!isCollapsed)}
      >
        <div className="flex items-start justify-between">
          {/* 左侧：图标和主要信息 */}
          <div className="flex items-start gap-3 flex-1">
            <div className="text-3xl">{icon}</div>
            <div className="flex-1">
              <div className="flex flex-wrap items-center gap-3 mb-2">
                {/* 信号类型 */}
                <div className={`badge ${
                  prediction.type === 'BUY' ? 'badge-success' : 'badge-error'
                } badge-lg gap-1`}>
                  {typeText}信号
                  <div className="tooltip tooltip-right" data-tip="买卖信号类型：BUY=买入，SELL=卖出">
                    <span className="text-xs opacity-70">ℹ️</span>
                  </div>
                </div>

                {/* 预测价格 */}
                <div className="flex items-center gap-1">
                  <span className="font-bold text-lg">¥{prediction.price.toFixed(2)}</span>
                  <div className="tooltip tooltip-right" data-tip="预测的目标价格">
                    <span className="text-xs opacity-70">ℹ️</span>
                  </div>
                </div>

                {/* 信号日期 */}
                <button
                  onClick={handleDateClick}
                  className="btn btn-ghost btn-xs gap-1 hover:btn-primary"
                  title="点击可跳转到日K线对应日期"
                >
                  📅 {formatDateForDisplay(prediction.signal_date)}
                  <div className="tooltip tooltip-right" data-tip="信号产生的日期 (点击可跳转到日K线对应日期)">
                    <span className="text-xs opacity-70">ℹ️</span>
                  </div>
                </button>

                {/* 概率 */}
                <div className="flex items-center gap-1">
                  <div className="badge badge-outline">
                    概率: {(prediction.probability * 100).toFixed(1)}%
                  </div>
                  <div className="tooltip tooltip-right" data-tip="预测成功的概率，基于技术指标置信度和历史表现">
                    <span className="text-xs opacity-70">ℹ️</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* 右侧：折叠按钮和强度标签 */}
          <div className="flex flex-col items-end gap-2">
            <button 
              className={`btn btn-ghost btn-xs transition-transform duration-200 ${
                isCollapsed ? '' : 'rotate-180'
              }`}
              onClick={(e) => {
                e.stopPropagation();
                setIsCollapsed(!isCollapsed);
              }}
            >
              🔽
            </button>
            <div className={`badge ${
              strength === 'STRONG' ? 'badge-success' :
              strength === 'MEDIUM' ? 'badge-warning' :
              'badge-ghost'
            }`}>
              {strength}
            </div>
          </div>
        </div>

        {/* 预测详情（可折叠） */}
        {!isCollapsed && (
          <div className="mt-4 space-y-3 border-t pt-3">
            {/* 预测理由 */}
            <div className="flex items-start gap-2">
              <span className="font-semibold text-sm">理由:</span>
              <div className="flex-1 flex items-start gap-1">
                <span className="text-sm">{prediction.reason || '基于技术指标分析'}</span>
                <div className="tooltip tooltip-right" data-tip="预测依据：包含识别的技术模式、置信度和强度等级">
                  <span className="text-xs opacity-70">ℹ️</span>
                </div>
              </div>
            </div>

            {/* 回测结果 */}
            {prediction.backtested && (
              <div className="bg-base-200 rounded-lg p-3 space-y-2">
                <div className="font-semibold text-sm mb-2">回测结果</div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-2">
                  <div className="flex items-center gap-2">
                    <span className="text-sm opacity-70">结果:</span>
                    <span className={`badge ${
                      prediction.is_correct ? 'badge-success' : 'badge-error'
                    }`}>
                      {prediction.is_correct ? '✅ 正确' : '❌ 错误'}
                    </span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-sm opacity-70">次日价格:</span>
                    <span className="font-semibold">
                      ¥{prediction.next_day_price?.toFixed(2) || 'N/A'}
                    </span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-sm opacity-70">价差:</span>
                    <span className={`font-semibold ${
                      (prediction.price_diff || 0) >= 0 ? 'text-success' : 'text-error'
                    }`}>
                      {(prediction.price_diff || 0) >= 0 ? '+' : ''}
                      {prediction.price_diff?.toFixed(2) || '0.00'} 
                      ({(prediction.price_diff_ratio || 0) >= 0 ? '+' : ''}
                      {prediction.price_diff_ratio?.toFixed(2) || '0.00'}%)
                    </span>
                  </div>
                </div>
              </div>
            )}

            {/* 相关指标 */}
            {prediction.indicators && prediction.indicators.length > 0 && (
              <div className="space-y-2">
                <div className="flex items-center gap-2 flex-wrap">
                  <span className="text-sm font-semibold">
                    {prediction.indicators.length > 1 ? '🔗 综合信号:' : '📊 相关指标:'}
                  </span>
                  {prediction.indicators.length > 1 && (
                    <div className="badge badge-warning badge-sm gap-1 animate-pulse">
                      ✨ {prediction.indicators.length}个指标共识
                    </div>
                  )}
                </div>
                <div className="flex items-center gap-2 flex-wrap">
                  {prediction.indicators.map((indicator, idx) => (
                    <div 
                      key={idx} 
                      className={`badge badge-sm ${
                        prediction.indicators.length > 1 
                          ? 'badge-primary font-medium shadow-md' 
                          : 'badge-outline'
                      }`}
                    >
                      {indicator}
                    </div>
                  ))}
                </div>
                {prediction.indicators.length > 1 && (
                  <div className="text-xs opacity-70 italic">
                    💡 多个技术指标共识，置信度已提升
                  </div>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

// 主组件
const PredictionsView: React.FC<PredictionsViewProps> = ({ 
  data, 
  isLoading = false,
  onDateClick 
}) => {
  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center py-12 space-y-4">
        <span className="loading loading-spinner loading-lg text-primary"></span>
        <p className="text-base-content/70">正在加载预测数据...</p>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="alert alert-info">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="stroke-current shrink-0 w-6 h-6">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <span>暂无预测数据</span>
      </div>
    );
  }

  const hasPredictions = data.predictions && data.predictions.length > 0;

  return (
    <div className="space-y-6">
      {/* 预测概览 */}
      <div className="stats shadow w-full">
        <div className="stat">
          <div className="stat-figure text-primary">
            🎯
          </div>
          <div className="stat-title">预测数量</div>
          <div className="stat-value text-primary">{data.predictions?.length || 0}</div>
          <div className="stat-desc">基于技术指标分析</div>
        </div>
        
        <div className="stat">
          <div className="stat-figure text-secondary">
            📊
          </div>
          <div className="stat-title">整体置信度</div>
          <div className="stat-value text-secondary">
            {(data.confidence * 100).toFixed(1)}%
          </div>
          <div className="stat-desc">预测准确度参考</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-accent">
            📅
          </div>
          <div className="stat-title">分析日期</div>
          <div className="stat-value text-accent text-2xl">
            {formatDateForDisplay(data.trade_date)}
          </div>
          <div className="stat-desc">
            更新于: {new Date(data.updated_at).toLocaleDateString('zh-CN')}
          </div>
        </div>
      </div>

      {/* 预测列表 */}
      {hasPredictions ? (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-bold">买卖信号列表</h3>
            <div className="text-sm text-base-content/70">
              点击卡片展开/收起详情
            </div>
          </div>
          
          {data.predictions.map((prediction, index) => (
            <PredictionItem
              key={index}
              prediction={prediction}
              index={index}
              onDateClick={onDateClick}
            />
          ))}
        </div>
      ) : (
        <div className="card bg-base-100 shadow-md">
          <div className="card-body items-center text-center">
            <div className="text-6xl mb-4">📊</div>
            <h3 className="card-title">暂无明确信号</h3>
            <p className="text-base-content/70">
              当前市场条件下，没有明确的买卖信号
            </p>
          </div>
        </div>
      )}

      {/* 免责声明 */}
      <div className="alert alert-warning">
        <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div>
          <h4 className="font-bold">⚠️ 免责声明</h4>
          <p className="text-sm">
            本系统提供的预测结果仅供参考，不构成投资建议。投资有风险，入市需谨慎。
          </p>
        </div>
      </div>
    </div>
  );
};

export default PredictionsView;
