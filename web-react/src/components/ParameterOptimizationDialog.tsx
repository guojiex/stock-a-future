/**
 * 参数优化对话框组件
 * 支持网格搜索和遗传算法两种优化方式
 */

import React, { useState, useEffect, useMemo } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Box,
  Typography,
  LinearProgress,
  Alert,
  Chip,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tabs,
  Tab,
  CircularProgress,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  Close as CloseIcon,
  Settings as SettingsIcon,
  TrendingUp as TrendingUpIcon,
  Speed as SpeedIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
  Bookmark as BookmarkIcon,
  Search as SearchIcon,
} from '@mui/icons-material';
import {
  useStartParameterOptimizationMutation,
  useLazyGetOptimizationProgressQuery,
  useLazyGetOptimizationResultsQuery,
  useCancelOptimizationMutation,
  useGetFavoritesQuery,
  useSearchStocksQuery,
} from '../services/api';

interface Strategy {
  id: string;
  name: string;
  strategy_type: string;
  parameters: Record<string, any>;
}

interface ParameterOptimizationDialogProps {
  open: boolean;
  onClose: () => void;
  strategy: Strategy | null;
  onOptimizationComplete?: (bestParams: Record<string, any>) => void;
}

interface ParameterRange {
  min: number;
  max: number;
  step: number;
}

const ParameterOptimizationDialog: React.FC<ParameterOptimizationDialogProps> = ({
  open,
  onClose,
  strategy,
  onOptimizationComplete,
}) => {
  // API Hooks
  const [startOptimization, { isLoading: isStarting }] = useStartParameterOptimizationMutation();
  const [getProgress] = useLazyGetOptimizationProgressQuery();
  const [getResults] = useLazyGetOptimizationResultsQuery();
  const [cancelOptimization] = useCancelOptimizationMutation();
  const { data: favoritesData } = useGetFavoritesQuery();
  const [searchQuery, setSearchQuery] = useState('');
  const { data: searchData } = useSearchStocksQuery({ q: searchQuery }, {
    skip: searchQuery.length < 2,
  });

  // State
  const [algorithm, setAlgorithm] = useState<'grid_search' | 'genetic'>('grid_search');
  const [optimizationTarget, setOptimizationTarget] = useState('sharpe_ratio');
  const [selectedSymbols, setSelectedSymbols] = useState<Array<{code: string; name: string}>>([]);
  const [startDate, setStartDate] = useState(getDefaultStartDate());
  const [endDate, setEndDate] = useState(getDefaultEndDate());
  const [initialCash, setInitialCash] = useState(1000000);
  const [commission, setCommission] = useState(0.0003);
  const [maxCombinations, setMaxCombinations] = useState(100);
  const [showStockSelector, setShowStockSelector] = useState(false);
  
  // Parameter ranges
  const [parameterRanges, setParameterRanges] = useState<Record<string, ParameterRange>>({});
  
  // Genetic algorithm config
  const [populationSize, setPopulationSize] = useState(50);
  const [generations, setGenerations] = useState(20);
  const [mutationRate, setMutationRate] = useState(0.1);
  const [crossoverRate, setCrossoverRate] = useState(0.7);
  const [elitismRate, setElitismRate] = useState(0.1);
  
  // Progress and results
  const [isRunning, setIsRunning] = useState(false);
  const [progress, setProgress] = useState(0);
  const [currentCombo, setCurrentCombo] = useState(0);
  const [totalCombos, setTotalCombos] = useState(0);
  const [currentOptimizationId, setCurrentOptimizationId] = useState<string | null>(null);
  const [optimizationResults, setOptimizationResults] = useState<any>(null);
  const [resultTab, setResultTab] = useState(0);

  // 默认日期函数
  function getDefaultStartDate() {
    const date = new Date();
    date.setFullYear(date.getFullYear() - 1);
    return date.toISOString().split('T')[0];
  }

  function getDefaultEndDate() {
    return new Date().toISOString().split('T')[0];
  }

  // 初始化参数范围（基于当前策略参数）
  useEffect(() => {
    if (strategy && strategy.parameters) {
      const ranges: Record<string, ParameterRange> = {};
      
      // 为每个数值参数生成合理的范围
      Object.entries(strategy.parameters).forEach(([key, value]) => {
        if (typeof value === 'number') {
          const val = value as number;
          
          // 根据参数名称和值生成合理的范围
          if (key.includes('period') || key.includes('window')) {
            // 周期类参数
            ranges[key] = {
              min: Math.max(1, Math.floor(val * 0.5)),
              max: Math.ceil(val * 2),
              step: 1,
            };
          } else if (key.includes('threshold') || key.includes('rate')) {
            // 阈值类参数
            ranges[key] = {
              min: val * 0.5,
              max: val * 1.5,
              step: (val * 0.1) || 0.01,
            };
          } else {
            // 默认范围
            ranges[key] = {
              min: val * 0.7,
              max: val * 1.3,
              step: (val * 0.1) || 1,
            };
          }
        }
      });
      
      setParameterRanges(ranges);
    }
  }, [strategy]);

  // 估计组合数量
  const estimatedCombinations = useMemo(() => {
    if (algorithm === 'genetic') {
      return populationSize * generations;
    }
    
    let total = 1;
    Object.values(parameterRanges).forEach((range) => {
      const steps = Math.ceil((range.max - range.min) / range.step) + 1;
      total *= steps;
    });
    return Math.min(total, maxCombinations);
  }, [algorithm, parameterRanges, maxCombinations, populationSize, generations]);

  // 轮询优化进度
  useEffect(() => {
    if (!isRunning || !currentOptimizationId) return;

    const interval = setInterval(async () => {
      try {
        const response = await getProgress(currentOptimizationId).unwrap();
        
        if (response.success && response.data) {
          setProgress(response.data.progress || 0);
          setCurrentCombo(response.data.current_combo || 0);
          setTotalCombos(response.data.total_combos || 0);
          
          if (response.data.status === 'completed') {
            clearInterval(interval);
            setIsRunning(false);
            loadResults(currentOptimizationId);
          } else if (response.data.status === 'failed' || response.data.status === 'cancelled') {
            clearInterval(interval);
            setIsRunning(false);
            alert(`优化${response.data.status === 'failed' ? '失败' : '已取消'}`);
          }
        }
      } catch (error) {
        console.error('获取优化进度失败:', error);
      }
    }, 2000);

    return () => clearInterval(interval);
  }, [isRunning, currentOptimizationId, getProgress]);

  // 加载优化结果
  const loadResults = async (optimizationId: string) => {
    try {
      const response = await getResults(optimizationId).unwrap();
      
      if (response.success && response.data) {
        setOptimizationResults(response.data);
      }
    } catch (error) {
      console.error('加载优化结果失败:', error);
    }
  };

  // 获取收藏列表和搜索结果
  const favorites = favoritesData?.data?.favorites || [];
  const searchResults = searchData?.data?.stocks || [];

  // 添加股票到选择列表
  const handleAddSymbol = (code: string, name: string) => {
    if (selectedSymbols.find(s => s.code === code)) {
      return; // 已存在
    }
    setSelectedSymbols([...selectedSymbols, { code, name }]);
  };

  // 从选择列表移除股票
  const handleRemoveSymbol = (code: string) => {
    setSelectedSymbols(selectedSymbols.filter(s => s.code !== code));
  };

  // 从收藏列表添加所有股票
  const handleAddAllFavorites = () => {
    const newSymbols = favorites
      .filter((fav: any) => !selectedSymbols.find(s => s.code === fav.ts_code))
      .map((fav: any) => ({ code: fav.ts_code, name: fav.name }));
    setSelectedSymbols([...selectedSymbols, ...newSymbols]);
  };

  // 处理启动优化
  const handleStartOptimization = async () => {
    if (!strategy) return;

    // 验证输入
    if (Object.keys(parameterRanges).length === 0) {
      alert('请至少配置一个参数范围');
      return;
    }

    if (selectedSymbols.length === 0) {
      alert('请至少选择一个股票');
      return;
    }

    try {
      const config = {
        strategy_type: strategy.strategy_type,
        parameter_ranges: parameterRanges,
        optimization_target: optimizationTarget,
        symbols: selectedSymbols.map(s => s.code),
        start_date: startDate,
        end_date: endDate,
        initial_cash: initialCash,
        commission: commission,
        algorithm: algorithm,
        max_combinations: maxCombinations,
        ...(algorithm === 'genetic' && {
          genetic_config: {
            population_size: populationSize,
            generations: generations,
            mutation_rate: mutationRate,
            crossover_rate: crossoverRate,
            elitism_rate: elitismRate,
          },
        }),
      };

      const response = await startOptimization({
        strategyId: strategy.id,
        config,
      }).unwrap();

      if (response.success && response.data) {
        setCurrentOptimizationId(response.data.optimization_id);
        setIsRunning(true);
        setProgress(0);
        setOptimizationResults(null);
      }
    } catch (error: any) {
      alert(`启动优化失败: ${error.data?.message || error.message || '未知错误'}`);
    }
  };

  // 处理取消优化
  const handleCancelOptimization = async () => {
    if (!currentOptimizationId) return;

    if (!window.confirm('确定要取消当前优化吗？')) return;

    try {
      await cancelOptimization(currentOptimizationId).unwrap();
      setIsRunning(false);
    } catch (error: any) {
      alert(`取消优化失败: ${error.message || '未知错误'}`);
    }
  };

  // 处理应用最佳参数
  const handleApplyBestParams = () => {
    if (optimizationResults && optimizationResults.best_parameters) {
      onOptimizationComplete?.(optimizationResults.best_parameters);
      onClose();
    }
  };

  // 渲染参数范围配置
  const renderParameterRanges = () => {
    return (
      <Box>
        <Typography variant="subtitle2" gutterBottom fontWeight="bold">
          参数范围配置
        </Typography>
        {Object.entries(parameterRanges).map(([paramName, range]) => (
          <Box key={paramName} sx={{ mb: 2, p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
            <Typography variant="body2" fontWeight="medium" gutterBottom>
              {paramName}
            </Typography>
            <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 2 }}>
              <TextField
                label="最小值"
                type="number"
                value={range.min}
                onChange={(e) => setParameterRanges({
                  ...parameterRanges,
                  [paramName]: { ...range, min: parseFloat(e.target.value) || 0 },
                })}
                size="small"
                fullWidth
              />
              <TextField
                label="最大值"
                type="number"
                value={range.max}
                onChange={(e) => setParameterRanges({
                  ...parameterRanges,
                  [paramName]: { ...range, max: parseFloat(e.target.value) || 0 },
                })}
                size="small"
                fullWidth
              />
              <TextField
                label="步长"
                type="number"
                value={range.step}
                onChange={(e) => setParameterRanges({
                  ...parameterRanges,
                  [paramName]: { ...range, step: parseFloat(e.target.value) || 1 },
                })}
                size="small"
                fullWidth
                inputProps={{ step: 0.01 }}
              />
            </Box>
          </Box>
        ))}
      </Box>
    );
  };

  // 渲染优化结果
  const renderOptimizationResults = () => {
    if (!optimizationResults) return null;

    return (
      <Box>
        <Alert severity="success" sx={{ mb: 2 }}>
          优化完成！测试了 {optimizationResults.total_tested} 组参数组合
        </Alert>

        <Tabs value={resultTab} onChange={(_, v) => setResultTab(v)} sx={{ mb: 2 }}>
          <Tab label="最佳参数" />
          <Tab label="性能指标" />
          <Tab label="所有结果" />
        </Tabs>

        {resultTab === 0 && (
          <Paper sx={{ p: 2 }}>
            <Typography variant="subtitle2" gutterBottom fontWeight="bold">
              最佳参数组合（得分: {optimizationResults.best_score.toFixed(4)}）
            </Typography>
            <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 2, mt: 2 }}>
              {Object.entries(optimizationResults.best_parameters).map(([key, value]) => (
                <Box key={key}>
                  <Typography variant="caption" color="text.secondary">
                    {key}
                  </Typography>
                  <Typography variant="body1" fontWeight="medium">
                    {typeof value === 'number' ? (value as number).toFixed(2) : String(value)}
                  </Typography>
                </Box>
              ))}
            </Box>
          </Paper>
        )}

        {resultTab === 1 && optimizationResults.performance && (
          <Paper sx={{ p: 2 }}>
            <Typography variant="subtitle2" gutterBottom fontWeight="bold">
              回测性能指标
            </Typography>
            <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 2, mt: 2 }}>
              <Box>
                <Typography variant="caption" color="text.secondary">总收益率</Typography>
                <Typography variant="h6" color={optimizationResults.performance.total_return >= 0 ? 'success.main' : 'error.main'}>
                  {(optimizationResults.performance.total_return * 100).toFixed(2)}%
                </Typography>
              </Box>
              <Box>
                <Typography variant="caption" color="text.secondary">年化收益</Typography>
                <Typography variant="h6">
                  {(optimizationResults.performance.annual_return * 100).toFixed(2)}%
                </Typography>
              </Box>
              <Box>
                <Typography variant="caption" color="text.secondary">夏普比率</Typography>
                <Typography variant="h6">
                  {optimizationResults.performance.sharpe_ratio.toFixed(2)}
                </Typography>
              </Box>
              <Box>
                <Typography variant="caption" color="text.secondary">最大回撤</Typography>
                <Typography variant="h6" color="error">
                  {(optimizationResults.performance.max_drawdown * 100).toFixed(2)}%
                </Typography>
              </Box>
              <Box>
                <Typography variant="caption" color="text.secondary">胜率</Typography>
                <Typography variant="h6">
                  {(optimizationResults.performance.win_rate * 100).toFixed(2)}%
                </Typography>
              </Box>
              <Box>
                <Typography variant="caption" color="text.secondary">交易次数</Typography>
                <Typography variant="h6">
                  {optimizationResults.performance.total_trades}
                </Typography>
              </Box>
            </Box>
          </Paper>
        )}

        {resultTab === 2 && optimizationResults.all_results && (
          <TableContainer component={Paper} sx={{ maxHeight: 400 }}>
            <Table stickyHeader size="small">
              <TableHead>
                <TableRow>
                  <TableCell>排名</TableCell>
                  <TableCell>得分</TableCell>
                  <TableCell>参数</TableCell>
                  <TableCell>收益率</TableCell>
                  <TableCell>夏普</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {optimizationResults.all_results.slice(0, 50).map((result: any, index: number) => (
                  <TableRow key={index}>
                    <TableCell>{index + 1}</TableCell>
                    <TableCell>{result.score.toFixed(4)}</TableCell>
                    <TableCell>
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                        {Object.entries(result.parameters).map(([k, v]) => (
                          <Chip
                            key={k}
                            label={`${k}:${typeof v === 'number' ? (v as number).toFixed(1) : v}`}
                            size="small"
                            variant="outlined"
                          />
                        ))}
                      </Box>
                    </TableCell>
                    <TableCell>
                      {result.performance ? `${(result.performance.total_return * 100).toFixed(2)}%` : '-'}
                    </TableCell>
                    <TableCell>
                      {result.performance ? result.performance.sharpe_ratio.toFixed(2) : '-'}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Box>
    );
  };

  if (!strategy) return null;

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="md"
      fullWidth
      disableEscapeKeyDown={isRunning}
    >
      <DialogTitle>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Box display="flex" alignItems="center" gap={1}>
            <SettingsIcon />
            <Typography variant="h6">参数优化 - {strategy.name}</Typography>
          </Box>
          <IconButton onClick={onClose} disabled={isRunning}>
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>

      <DialogContent dividers>
        {!optimizationResults ? (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
            {/* 优化算法选择 */}
            <Box>
              <Typography variant="subtitle2" gutterBottom fontWeight="bold">
                优化算法
              </Typography>
              <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 2 }}>
                <Paper
                  sx={{
                    p: 2,
                    cursor: 'pointer',
                    border: 2,
                    borderColor: algorithm === 'grid_search' ? 'primary.main' : 'grey.300',
                  }}
                  onClick={() => setAlgorithm('grid_search')}
                >
                  <Box display="flex" alignItems="center" gap={1} mb={1}>
                    <SpeedIcon color={algorithm === 'grid_search' ? 'primary' : 'disabled'} />
                    <Typography variant="subtitle2">网格搜索</Typography>
                  </Box>
                  <Typography variant="caption" color="text.secondary">
                    穷举所有参数组合，精确但耗时较长
                  </Typography>
                </Paper>
                <Paper
                  sx={{
                    p: 2,
                    cursor: 'pointer',
                    border: 2,
                    borderColor: algorithm === 'genetic' ? 'primary.main' : 'grey.300',
                  }}
                  onClick={() => setAlgorithm('genetic')}
                >
                  <Box display="flex" alignItems="center" gap={1} mb={1}>
                    <TrendingUpIcon color={algorithm === 'genetic' ? 'primary' : 'disabled'} />
                    <Typography variant="subtitle2">遗传算法</Typography>
                  </Box>
                  <Typography variant="caption" color="text.secondary">
                    智能搜索，速度快但可能不是全局最优
                  </Typography>
                </Paper>
              </Box>
            </Box>

            {/* 优化目标 */}
            <FormControl fullWidth>
              <InputLabel>优化目标</InputLabel>
              <Select
                value={optimizationTarget}
                label="优化目标"
                onChange={(e) => setOptimizationTarget(e.target.value)}
              >
                <MenuItem value="sharpe_ratio">夏普比率（综合风险收益）</MenuItem>
                <MenuItem value="total_return">总收益率</MenuItem>
                <MenuItem value="win_rate">胜率</MenuItem>
                <MenuItem value="profit_factor">盈亏比</MenuItem>
              </Select>
            </FormControl>

            {/* 参数范围 */}
            {renderParameterRanges()}

            {/* 回测配置 */}
            <Box>
              <Typography variant="subtitle2" gutterBottom fontWeight="bold">
                回测配置
              </Typography>
              <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 2 }}>
                <TextField
                  label="开始日期"
                  type="date"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                  InputLabelProps={{ shrink: true }}
                  fullWidth
                />
                <TextField
                  label="结束日期"
                  type="date"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                  InputLabelProps={{ shrink: true }}
                  fullWidth
                />
                <TextField
                  label="初始资金"
                  type="number"
                  value={initialCash}
                  onChange={(e) => setInitialCash(parseFloat(e.target.value) || 0)}
                  fullWidth
                />
                <TextField
                  label="手续费率"
                  type="number"
                  value={commission}
                  onChange={(e) => setCommission(parseFloat(e.target.value) || 0)}
                  fullWidth
                  inputProps={{ step: 0.0001 }}
                />
              </Box>
              {/* 股票选择 */}
              <Box sx={{ mt: 2 }}>
                <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
                  <Typography variant="subtitle2" fontWeight="bold">
                    选择股票（已选 {selectedSymbols.length} 只）
                  </Typography>
                  <Box>
                    <Button
                      size="small"
                      startIcon={<BookmarkIcon />}
                      onClick={handleAddAllFavorites}
                      disabled={favorites.length === 0}
                    >
                      添加所有收藏
                    </Button>
                    <Button
                      size="small"
                      startIcon={<SearchIcon />}
                      onClick={() => setShowStockSelector(!showStockSelector)}
                      sx={{ ml: 1 }}
                    >
                      {showStockSelector ? '收起' : '搜索添加'}
                    </Button>
                  </Box>
                </Box>

                {/* 已选股票列表 */}
                {selectedSymbols.length > 0 ? (
                  <Paper variant="outlined" sx={{ p: 1, mb: 2, maxHeight: 150, overflow: 'auto' }}>
                    <Box display="flex" flexWrap="wrap" gap={1}>
                      {selectedSymbols.map((symbol) => (
                        <Chip
                          key={symbol.code}
                          label={`${symbol.name} (${symbol.code})`}
                          onDelete={() => handleRemoveSymbol(symbol.code)}
                          color="primary"
                          variant="outlined"
                          size="small"
                        />
                      ))}
                    </Box>
                  </Paper>
                ) : (
                  <Alert severity="info" sx={{ mb: 2 }}>
                    请从收藏列表或搜索结果中添加股票
                  </Alert>
                )}

                {/* 股票选择器 */}
                {showStockSelector && (
                  <Paper variant="outlined" sx={{ p: 2, mb: 2 }}>
                    <Tabs value={0} sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
                      <Tab label={`收藏列表 (${favorites.length})`} />
                    </Tabs>

                    {/* 收藏列表 */}
                    <Box>
                      <TextField
                        placeholder="搜索股票..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        size="small"
                        fullWidth
                        sx={{ mb: 2 }}
                        InputProps={{
                          startAdornment: <SearchIcon sx={{ mr: 1, color: 'text.secondary' }} />,
                        }}
                      />

                      {/* 显示搜索结果或收藏列表 */}
                      <Box sx={{ maxHeight: 200, overflow: 'auto' }}>
                        {searchQuery.length >= 2 ? (
                          // 搜索结果
                          searchResults.length > 0 ? (
                            searchResults.map((stock: any) => (
                              <Box
                                key={stock.ts_code}
                                sx={{
                                  p: 1,
                                  display: 'flex',
                                  justifyContent: 'space-between',
                                  alignItems: 'center',
                                  '&:hover': { bgcolor: 'action.hover' },
                                  cursor: 'pointer',
                                }}
                              >
                                <Box onClick={() => handleAddSymbol(stock.ts_code, stock.name)}>
                                  <Typography variant="body2">
                                    {stock.name} ({stock.ts_code})
                                  </Typography>
                                </Box>
                                <IconButton
                                  size="small"
                                  onClick={() => handleAddSymbol(stock.ts_code, stock.name)}
                                  disabled={selectedSymbols.some(s => s.code === stock.ts_code)}
                                >
                                  <AddIcon fontSize="small" />
                                </IconButton>
                              </Box>
                            ))
                          ) : (
                            <Typography variant="body2" color="text.secondary" textAlign="center">
                              未找到匹配的股票
                            </Typography>
                          )
                        ) : (
                          // 收藏列表
                          favorites.length > 0 ? (
                            favorites.map((fav: any) => (
                              <Box
                                key={fav.ts_code}
                                sx={{
                                  p: 1,
                                  display: 'flex',
                                  justifyContent: 'space-between',
                                  alignItems: 'center',
                                  '&:hover': { bgcolor: 'action.hover' },
                                  cursor: 'pointer',
                                }}
                              >
                                <Box onClick={() => handleAddSymbol(fav.ts_code, fav.name)}>
                                  <Typography variant="body2">
                                    {fav.name} ({fav.ts_code})
                                  </Typography>
                                  {fav.group_name && (
                                    <Typography variant="caption" color="text.secondary">
                                      {fav.group_name}
                                    </Typography>
                                  )}
                                </Box>
                                <IconButton
                                  size="small"
                                  onClick={() => handleAddSymbol(fav.ts_code, fav.name)}
                                  disabled={selectedSymbols.some(s => s.code === fav.ts_code)}
                                >
                                  <AddIcon fontSize="small" />
                                </IconButton>
                              </Box>
                            ))
                          ) : (
                            <Typography variant="body2" color="text.secondary" textAlign="center">
                              暂无收藏股票，请先添加收藏
                            </Typography>
                          )
                        )}
                      </Box>
                    </Box>
                  </Paper>
                )}
              </Box>
            </Box>

            {/* 算法特定配置 */}
            {algorithm === 'grid_search' && (
              <TextField
                label="最大组合数"
                type="number"
                value={maxCombinations}
                onChange={(e) => setMaxCombinations(parseInt(e.target.value) || 100)}
                fullWidth
                helperText={`估计将测试约 ${estimatedCombinations} 组参数`}
              />
            )}

            {algorithm === 'genetic' && (
              <Box>
                <Typography variant="subtitle2" gutterBottom fontWeight="bold">
                  遗传算法参数
                </Typography>
                <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 2 }}>
                  <TextField
                    label="种群大小"
                    type="number"
                    value={populationSize}
                    onChange={(e) => setPopulationSize(parseInt(e.target.value) || 50)}
                    fullWidth
                    size="small"
                  />
                  <TextField
                    label="迭代代数"
                    type="number"
                    value={generations}
                    onChange={(e) => setGenerations(parseInt(e.target.value) || 20)}
                    fullWidth
                    size="small"
                  />
                  <TextField
                    label="变异率"
                    type="number"
                    value={mutationRate}
                    onChange={(e) => setMutationRate(parseFloat(e.target.value) || 0.1)}
                    fullWidth
                    size="small"
                    inputProps={{ step: 0.01, min: 0, max: 1 }}
                  />
                  <TextField
                    label="交叉率"
                    type="number"
                    value={crossoverRate}
                    onChange={(e) => setCrossoverRate(parseFloat(e.target.value) || 0.7)}
                    fullWidth
                    size="small"
                    inputProps={{ step: 0.01, min: 0, max: 1 }}
                  />
                </Box>
                <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                  估计将测试约 {estimatedCombinations} 组参数
                </Typography>
              </Box>
            )}

            {/* 进度显示 */}
            {isRunning && (
              <Box>
                <Alert severity="info" sx={{ mb: 2 }}>
                  优化运行中... ({currentCombo}/{totalCombos})
                </Alert>
                <LinearProgress variant="determinate" value={progress} />
                <Box display="flex" justifyContent="space-between" mt={1}>
                  <Typography variant="caption">{progress}%</Typography>
                  <Typography variant="caption">
                    预计还需 {Math.ceil((100 - progress) / 10)} 分钟
                  </Typography>
                </Box>
              </Box>
            )}
          </Box>
        ) : (
          renderOptimizationResults()
        )}
      </DialogContent>

      <DialogActions>
        {!optimizationResults ? (
          <>
            <Button onClick={onClose} disabled={isRunning}>
              取消
            </Button>
            {isRunning ? (
              <Button onClick={handleCancelOptimization} color="error" variant="outlined">
                停止优化
              </Button>
            ) : (
              <Button
                onClick={handleStartOptimization}
                variant="contained"
                disabled={isStarting}
                startIcon={isStarting ? <CircularProgress size={20} /> : <SettingsIcon />}
              >
                {isStarting ? '启动中...' : '开始优化'}
              </Button>
            )}
          </>
        ) : (
          <>
            <Button onClick={() => {
              setOptimizationResults(null);
              setIsRunning(false);
              setProgress(0);
            }}>
              重新优化
            </Button>
            <Button onClick={onClose} variant="outlined">
              关闭
            </Button>
            <Button
              onClick={handleApplyBestParams}
              variant="contained"
              color="primary"
              startIcon={<TrendingUpIcon />}
            >
              应用最佳参数
            </Button>
          </>
        )}
      </DialogActions>
    </Dialog>
  );
};

export default ParameterOptimizationDialog;

