/**
 * 回测页面 - 完整实现
 * 参考网页版 web/static/js/modules/backtest.js
 */

import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import {
  Box,
  Container,
  Typography,
  Paper,
  Alert,
  Button,
  TextField,
  LinearProgress,
  Card,
  CardContent,
  Chip,
  Divider,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  CircularProgress,
  Grid,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tabs,
  Tab,
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  PlayArrow as PlayArrowIcon,
  Stop as StopIcon,
  Save as SaveIcon,
  Close as CloseIcon,
} from '@mui/icons-material';
import { useAppSelector } from '../hooks/redux';
import {
  useGetStrategiesQuery,
  useCreateBacktestMutation,
  useStartBacktestMutation,
  useCancelBacktestMutation,
  useLazyGetBacktestProgressQuery,
  useLazyGetBacktestResultsQuery,
} from '../services/api';
import {
  updateConfig,
  addStrategy,
  removeStrategy,
  clearStrategies,
  setSelectedStrategies,
} from '../store/slices/backtestSlice';
import EquityCurveChart from '../components/EquityCurveChart';
import TradesTable from '../components/TradesTable';

// 定义接口
interface BacktestConfig {
  name: string;
  strategy_ids: string[];
  start_date: string;
  end_date: string;
  initial_cash: number;
  commission: number;
  symbols: string[];
}

interface BacktestResult {
  id: string;
  status: string;
  progress: number;
  message: string;
  error?: string;
}

interface PerformanceMetrics {
  total_return: number;
  annual_return: number;
  max_drawdown: number;
  sharpe_ratio: number;
  win_rate: number;
  total_trades: number;
  avg_trade_return?: number;
  profit_factor?: number;
}

const BacktestPage: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  
  // Redux state
  const selectedStrategyIds = useAppSelector((state) => state.backtest.selectedStrategyIds);
  const storedConfig = useAppSelector((state) => state.backtest.config);
  
  // Local state
  const [config, setConfig] = useState<BacktestConfig>({
    name: storedConfig.name || '',
    strategy_ids: selectedStrategyIds,
    start_date: storedConfig.startDate || getDefaultStartDate(),
    end_date: storedConfig.endDate || getDefaultEndDate(),
    initial_cash: storedConfig.initialCash || 1000000,
    commission: storedConfig.commission || 0.0003,
    symbols: storedConfig.symbols || [],
  });
  
  const [symbolsText, setSymbolsText] = useState(storedConfig.symbols.join('\n'));
  const [isRunning, setIsRunning] = useState(false);
  const [progress, setProgress] = useState(0);
  const [progressMessage, setProgressMessage] = useState('');
  const [currentBacktestId, setCurrentBacktestId] = useState<string | null>(null);
  const [showResults, setShowResults] = useState(false);
  const [results, setResults] = useState<any>(null);
  const [selectedStrategiesDialog, setSelectedStrategiesDialog] = useState(false);
  const [selectedStrategyTab, setSelectedStrategyTab] = useState(0);
  
  // API hooks
  const { data: strategiesData } = useGetStrategiesQuery();
  const [createBacktest] = useCreateBacktestMutation();
  const [startBacktest] = useStartBacktestMutation();
  const [cancelBacktest] = useCancelBacktestMutation();
  const [getProgress] = useLazyGetBacktestProgressQuery();
  const [getResults] = useLazyGetBacktestResultsQuery();
  
  const strategies = (strategiesData?.data?.items || strategiesData?.data?.data || []) as any[];
  
  // 默认全选策略（仅在首次加载时且没有选择任何策略时）
  useEffect(() => {
    if (strategies.length > 0 && selectedStrategyIds.length === 0) {
      // 默认选择所有策略（最多5个）
      const allStrategyIds = strategies.slice(0, 5).map(s => s.id);
      dispatch(setSelectedStrategies(allStrategyIds));
    }
  }, [strategies, selectedStrategyIds.length, dispatch]);
  
  // 默认日期
  function getDefaultStartDate() {
    const date = new Date();
    date.setFullYear(date.getFullYear() - 1);
    return date.toISOString().split('T')[0];
  }
  
  function getDefaultEndDate() {
    return new Date().toISOString().split('T')[0];
  }
  
  // 更新配置
  const handleConfigChange = (field: keyof BacktestConfig, value: any) => {
    setConfig((prev) => ({ ...prev, [field]: value }));
  };
  
  // 解析股票代码
  const parseSymbols = useCallback(() => {
    const symbols = symbolsText
      .split('\n')
      .map((s) => s.trim())
      .filter((s) => s.length > 0);
    handleConfigChange('symbols', symbols);
  }, [symbolsText]);
  
  useEffect(() => {
    parseSymbols();
  }, [symbolsText, parseSymbols]);
  
  // 验证配置
  const validateConfig = (): string | null => {
    if (!config.name) return '请输入回测名称';
    // 使用 Redux state 中的 selectedStrategyIds 进行验证
    if (selectedStrategyIds.length === 0) return '请选择至少一个策略';
    if (selectedStrategyIds.length > 5) return '最多只能选择5个策略';
    if (!config.start_date || !config.end_date) return '请选择开始和结束日期';
    if (new Date(config.start_date) >= new Date(config.end_date)) {
      return '开始日期必须早于结束日期';
    }
    if (config.symbols.length === 0) return '请输入至少一个股票代码';
    if (config.initial_cash < 10000) return '初始资金不能少于10000元';
    return null;
  };
  
  // 开始回测
  const handleStartBacktest = async () => {
    const error = validateConfig();
    if (error) {
      alert(error);
      return;
    }
    
    try {
      setIsRunning(true);
      setProgress(0);
      setProgressMessage('准备中...');
      setShowResults(false);
      
      // 更新Redux配置
      dispatch(updateConfig({
        name: config.name,
        startDate: config.start_date,
        endDate: config.end_date,
        initialCash: config.initial_cash,
        commission: config.commission,
        symbols: config.symbols,
      }));
      
      // 创建回测（后端会自动启动）
      setProgressMessage('创建并启动回测...');
      const createResponse = await createBacktest({
        name: config.name,
        strategy_ids: selectedStrategyIds,
        start_date: config.start_date,
        end_date: config.end_date,
        initial_cash: config.initial_cash,
        commission: config.commission,
        symbols: config.symbols,
      }).unwrap();
      
      if (!createResponse.success || !createResponse.data) {
        throw new Error(createResponse.message || '创建回测失败');
      }
      
      const backtestId = createResponse.data.id;
      setCurrentBacktestId(backtestId);
      
      // 后端已自动启动回测，直接开始监控进度
      setProgressMessage('回测运行中...');
      startProgressMonitoring(backtestId);
    } catch (error: any) {
      console.error('启动回测失败:', error);
      const errorMessage = error.data?.message || error.message || '未知错误';
      alert(`启动回测失败: ${errorMessage}`);
      setIsRunning(false);
    }
  };
  
  // 监控进度
  const startProgressMonitoring = (backtestId: string) => {
    const interval = setInterval(async () => {
      try {
        const response = await getProgress(backtestId).unwrap();
        
        if (response.success && response.data) {
          const data = response.data;
          setProgress(data.progress || 0);
          setProgressMessage(data.message || '运行中...');
          
          if (data.status === 'completed') {
            clearInterval(interval);
            setIsRunning(false);
            loadResults(backtestId);
          } else if (data.status === 'failed') {
            clearInterval(interval);
            setIsRunning(false);
            alert(data.error || '回测失败');
          } else if (data.status === 'cancelled') {
            clearInterval(interval);
            setIsRunning(false);
            alert('回测已取消');
          }
        }
      } catch (error) {
        console.error('获取进度失败:', error);
      }
    }, 2000);
    
    // 存储interval ID以便取消
    (window as any).backtestInterval = interval;
  };
  
  // 停止回测
  const handleStopBacktest = async () => {
    if (!currentBacktestId) return;
    
    if (!window.confirm('确定要停止当前回测吗？')) return;
    
    try {
      await cancelBacktest(currentBacktestId).unwrap();
      
      // 清除进度监控
      if ((window as any).backtestInterval) {
        clearInterval((window as any).backtestInterval);
      }
      
      setIsRunning(false);
      setCurrentBacktestId(null);
      alert('回测已停止');
    } catch (error: any) {
      console.error('停止回测失败:', error);
      alert(`停止回测失败: ${error.message || '未知错误'}`);
    }
  };
  
  // 加载结果
  const loadResults = async (backtestId: string) => {
    try {
      const response = await getResults(backtestId).unwrap();
      
      if (response.success && response.data) {
        setResults(response.data);
        setShowResults(true);
      } else {
        throw new Error(response.message || '获取回测结果失败');
      }
    } catch (error: any) {
      console.error('加载回测结果失败:', error);
      alert(`加载回测结果失败: ${error.message || '未知错误'}`);
    }
  };
  
  // 获取选中策略的信息
  const getSelectedStrategiesInfo = () => {
    return selectedStrategyIds.map((id) => {
      const strategy = strategies.find((s) => s.id === id);
      return strategy || { id, name: id };
    });
  };
  
  // 渲染紧凑的性能指标（横向）
  const renderCompactMetrics = (metrics: PerformanceMetrics) => {
    const metricsList = [
      { label: '总收益率', value: `${(metrics.total_return * 100).toFixed(2)}%`, positive: metrics.total_return >= 0 },
      { label: '年化收益', value: `${(metrics.annual_return * 100).toFixed(2)}%`, positive: metrics.annual_return >= 0 },
      { label: '最大回撤', value: `${(metrics.max_drawdown * 100).toFixed(2)}%`, positive: metrics.max_drawdown >= -0.05 },
      { label: '夏普比率', value: metrics.sharpe_ratio.toFixed(2), positive: metrics.sharpe_ratio >= 1 },
      { label: '胜率', value: `${(metrics.win_rate * 100).toFixed(2)}%`, positive: metrics.win_rate >= 0.5 },
      { label: '交易次数', value: metrics.total_trades.toString(), positive: true },
    ];
    
    return (
      <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap', alignItems: 'center' }}>
        {metricsList.map((metric) => (
          <Box key={metric.label} sx={{ minWidth: 140 }}>
            <Typography variant="caption" color="text.secondary" display="block">
              {metric.label}
            </Typography>
            <Typography
              variant="h6"
              color={metric.positive ? 'success.main' : 'error.main'}
            >
              {metric.value}
            </Typography>
          </Box>
        ))}
      </Box>
    );
  };
  
  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      {/* 页面标题 */}
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 4 }}>
        <Button
          variant="outlined"
          startIcon={<ArrowBackIcon />}
          onClick={() => navigate('/strategies')}
        >
          返回策略列表
        </Button>
        <Typography variant="h4" component="h1" fontWeight="bold">
          🚀 回测系统
        </Typography>
      </Box>
      
      {/* 选中的策略信息 */}
      {selectedStrategyIds.length > 0 ? (
        <Alert severity="success" sx={{ mb: 3 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Typography>已选择 {selectedStrategyIds.length} 个策略:</Typography>
            {getSelectedStrategiesInfo().map((s) => (
              <Chip
                key={s.id}
                label={s.name}
                size="small"
                onDelete={() => dispatch(removeStrategy(s.id))}
              />
            ))}
            <Button
              size="small"
              onClick={() => setSelectedStrategiesDialog(true)}
            >
              管理策略
            </Button>
          </Box>
        </Alert>
      ) : (
        <Alert severity="info" sx={{ mb: 3 }}>
          请先选择要回测的策略
          <Button
            size="small"
            onClick={() => setSelectedStrategiesDialog(true)}
            sx={{ ml: 2 }}
          >
            选择策略
          </Button>
        </Alert>
      )}
      
      {/* 回测配置表单 */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom>
          ⚙️ 回测配置
        </Typography>
        <Divider sx={{ mb: 3 }} />
        
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          {/* 回测名称 */}
          <TextField
            label="回测名称"
            value={config.name}
            onChange={(e) => handleConfigChange('name', e.target.value)}
            fullWidth
            required
            placeholder="例如: MACD策略-2024年测试"
          />
          
          {/* 时间范围 */}
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 2 }}>
            <TextField
              label="开始日期"
              type="date"
              value={config.start_date}
              onChange={(e) => handleConfigChange('start_date', e.target.value)}
              InputLabelProps={{ shrink: true }}
              required
            />
            <TextField
              label="结束日期"
              type="date"
              value={config.end_date}
              onChange={(e) => handleConfigChange('end_date', e.target.value)}
              InputLabelProps={{ shrink: true }}
              required
            />
          </Box>
          
          {/* 资金和手续费 */}
          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 2 }}>
            <TextField
              label="初始资金"
              type="number"
              value={config.initial_cash}
              onChange={(e) => handleConfigChange('initial_cash', parseFloat(e.target.value))}
              required
              InputProps={{ inputProps: { min: 10000, step: 10000 } }}
              helperText="单位: 元"
            />
            <TextField
              label="手续费率"
              type="number"
              value={config.commission}
              onChange={(e) => handleConfigChange('commission', parseFloat(e.target.value))}
              required
              InputProps={{ inputProps: { min: 0, step: 0.0001 } }}
              helperText="默认: 0.03%"
            />
          </Box>
          
          {/* 股票列表 */}
          <TextField
            label="股票列表"
            value={symbolsText}
            onChange={(e) => setSymbolsText(e.target.value)}
            multiline
            rows={4}
            placeholder="输入股票代码，每行一个，例如:&#10;000001.SZ&#10;600000.SH"
            helperText="每行输入一个股票代码"
          />
          
          {/* 操作按钮 */}
          <Box sx={{ display: 'flex', gap: 2 }}>
            <Button
              variant="contained"
              size="large"
              startIcon={isRunning ? <StopIcon /> : <PlayArrowIcon />}
              onClick={isRunning ? handleStopBacktest : handleStartBacktest}
              disabled={selectedStrategyIds.length === 0}
              color={isRunning ? 'error' : 'primary'}
            >
              {isRunning ? '停止回测' : '开始回测'}
            </Button>
            <Button
              variant="outlined"
              startIcon={<SaveIcon />}
              onClick={() => {
                dispatch(updateConfig({
                  name: config.name,
                  startDate: config.start_date,
                  endDate: config.end_date,
                  initialCash: config.initial_cash,
                  commission: config.commission,
                  symbols: config.symbols,
                }));
                alert('配置已保存');
              }}
            >
              保存配置
            </Button>
          </Box>
        </Box>
      </Paper>
      
      {/* 回测进度 */}
      {isRunning && (
        <Paper sx={{ p: 3, mb: 3 }}>
          <Typography variant="h6" gutterBottom>
            ⏳ 回测进度
          </Typography>
          <Box sx={{ mt: 2 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
              <Typography variant="body2">{progressMessage}</Typography>
              <Typography variant="body2">{progress}%</Typography>
            </Box>
            <LinearProgress variant="determinate" value={progress} />
          </Box>
        </Paper>
      )}
      
      {/* 回测结果 */}
      {showResults && results && (
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            📊 回测结果
          </Typography>
          <Divider sx={{ mb: 3 }} />
          
          {/* 组合整体性能（可选） */}
          {results.combined_metrics && (
            <Box sx={{ mb: 3, p: 2, bgcolor: 'primary.50', borderRadius: 1 }}>
              <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                组合整体表现
              </Typography>
              {renderCompactMetrics(results.combined_metrics)}
            </Box>
          )}
          
          {/* 多策略 Tab 展示 */}
          {results.performance && Array.isArray(results.performance) && results.performance.length > 1 ? (
            <Box>
              <Tabs 
                value={selectedStrategyTab} 
                onChange={(_, newValue) => setSelectedStrategyTab(newValue)}
                variant="scrollable"
                scrollButtons="auto"
                sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}
              >
                {results.performance.map((_, index) => {
                  const strategy = results.strategies?.[index];
                  return (
                    <Tab 
                      key={index} 
                      label={strategy?.name || `策略 ${index + 1}`}
                      icon={<Chip label={index + 1} size="small" color="primary" />}
                      iconPosition="start"
                    />
                  );
                })}
              </Tabs>
              
              {results.performance.map((perfMetrics, index) => {
                const strategy = results.strategies?.[index];
                const strategyTrades = results.trades?.filter((t: any) => t.strategy_id === strategy?.id) || [];
                
                return (
                  <Box 
                    key={index} 
                    role="tabpanel"
                    hidden={selectedStrategyTab !== index}
                  >
                    {selectedStrategyTab === index && (
                      <Box>
                        {/* 策略性能指标 */}
                        <Box sx={{ mb: 3, p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
                          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                            {strategy?.name || `策略 ${index + 1}`}
                          </Typography>
                          {strategy?.description && (
                            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                              {strategy.description}
                            </Typography>
                          )}
                          {renderCompactMetrics(perfMetrics)}
                        </Box>
                        
                        {/* 策略权益曲线 */}
                        {results.equity_curve && results.equity_curve.length > 0 && (
                          <Box sx={{ mb: 3 }}>
                            <Typography variant="subtitle2" gutterBottom>
                              📈 权益曲线
                            </Typography>
                            <EquityCurveChart
                              data={results.equity_curve}
                              initialCash={config.initial_cash}
                            />
                          </Box>
                        )}
                        
                        {/* 策略交易记录 */}
                        {strategyTrades.length > 0 && (
                          <Box>
                            <Typography variant="subtitle2" gutterBottom>
                              📋 交易记录 ({strategyTrades.length} 笔)
                            </Typography>
                            <TradesTable 
                              trades={strategyTrades}
                              strategies={[strategy]}
                            />
                          </Box>
                        )}
                      </Box>
                    )}
                  </Box>
                );
              })}
            </Box>
          ) : results.performance && Array.isArray(results.performance) && results.performance.length === 1 ? (
            // 单策略展示
            <Box>
              <Box sx={{ mb: 3, p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
                {renderCompactMetrics(results.performance[0])}
              </Box>
              
              {results.equity_curve && results.equity_curve.length > 0 && (
                <Box sx={{ mb: 3 }}>
                  <Typography variant="subtitle2" gutterBottom>
                    📈 权益曲线
                  </Typography>
                  <EquityCurveChart
                    data={results.equity_curve}
                    initialCash={config.initial_cash}
                  />
                </Box>
              )}
              
              {results.trades && results.trades.length > 0 && (
                <Box>
                  <Typography variant="subtitle2" gutterBottom>
                    📋 交易记录
                  </Typography>
                  <TradesTable 
                    trades={results.trades}
                    strategies={results.strategies}
                  />
                </Box>
              )}
            </Box>
          ) : null}
        </Paper>
      )}
      
      {/* 策略选择对话框 */}
      <Dialog
        open={selectedStrategiesDialog}
        onClose={() => setSelectedStrategiesDialog(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>选择策略</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
            <Typography variant="body2" color="text.secondary">
              最多可以选择5个策略进行回测
            </Typography>
            <Button
              size="small"
              onClick={() => {
                // 全选所有策略（最多5个）
                const allStrategyIds = strategies.slice(0, 5).map(s => s.id);
                dispatch(setSelectedStrategies(allStrategyIds));
              }}
            >
              全选
            </Button>
          </Box>
          <Box sx={{ mt: 2 }}>
            {strategies.map((strategy) => {
              const isSelected = selectedStrategyIds.includes(strategy.id);
              return (
                <Box
                  key={strategy.id}
                  sx={{
                    p: 2,
                    mb: 1,
                    border: 1,
                    borderColor: isSelected ? 'primary.main' : 'grey.300',
                    borderRadius: 1,
                    cursor: 'pointer',
                    bgcolor: isSelected ? 'primary.50' : 'transparent',
                  }}
                  onClick={() => {
                    if (isSelected) {
                      dispatch(removeStrategy(strategy.id));
                    } else if (selectedStrategyIds.length < 5) {
                      dispatch(addStrategy(strategy.id));
                    } else {
                      alert('最多只能选择5个策略');
                    }
                  }}
                >
                  <Typography variant="subtitle1">{strategy.name}</Typography>
                  <Typography variant="body2" color="text.secondary">
                    {strategy.description || '暂无描述'}
                  </Typography>
                </Box>
              );
            })}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => dispatch(clearStrategies())}>清空选择</Button>
          <Button onClick={() => setSelectedStrategiesDialog(false)}>关闭</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default BacktestPage;
