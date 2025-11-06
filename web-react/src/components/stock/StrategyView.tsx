/**
 * 策略视图组件 - Web版本
 * 在股票详情页显示该股票的收藏策略信息和信号
 */

import React, { useState, useMemo } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Button,
  Chip,
  CircularProgress,
  Alert,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Tab,
  Tabs,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Divider,
} from '@mui/material';
import {
  BookmarkBorder as BookmarkBorderIcon,
  Bookmark as BookmarkIcon,
  Delete as DeleteIcon,
  Refresh as RefreshIcon,
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
} from '@mui/icons-material';

import {
  useGetFavoritesQuery,
  useGetGroupsQuery,
  useAddFavoriteMutation,
  useDeleteFavoriteMutation,
  useGetPredictionsQuery,
  useGetPatternSummaryQuery,
} from '../../services/api';
import { PatternSummary } from '../../types/stock';

interface StrategyViewProps {
  stockCode: string;
  stockName: string;
}

interface TabPanelProps {
  children?: React.ReactNode;
  value: number;
  index: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`strategy-tabpanel-${index}`}
      aria-labelledby={`strategy-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
    </div>
  );
}

const StrategyView: React.FC<StrategyViewProps> = ({ stockCode, stockName }) => {
  // 状态
  const [selectedTab, setSelectedTab] = useState(0);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  // API查询
  const {
    data: favoritesData,
    isLoading: favoritesLoading,
    refetch: refetchFavorites,
  } = useGetFavoritesQuery();

  const {
    data: groupsData,
    isLoading: groupsLoading,
  } = useGetGroupsQuery();

  const {
    data: predictionsData,
    isLoading: predictionsLoading,
    refetch: refetchPredictions,
  } = useGetPredictionsQuery(stockCode, { skip: !stockCode });

  const {
    data: patternSummaryData,
    isLoading: patternSummaryLoading,
    refetch: refetchPatternSummary,
  } = useGetPatternSummaryQuery(stockCode, { skip: !stockCode });

  // API变更
  const [addFavorite, { isLoading: addLoading }] = useAddFavoriteMutation();
  const [deleteFavorite, { isLoading: deleteLoading }] = useDeleteFavoriteMutation();

  // 数据处理
  const favorites = favoritesData?.data?.favorites || [];
  const groups = groupsData?.data?.groups || [];
  const predictionResult = predictionsData?.data;
  const predictions = predictionResult?.predictions || [];
  const patternSummary: PatternSummary | undefined = patternSummaryData?.data;

  // 检查当前股票是否已收藏
  const currentFavorite = useMemo(() => {
    return favorites.find((fav: any) => fav.ts_code === stockCode);
  }, [favorites, stockCode]);

  const isFavorited = !!currentFavorite;

  // 获取分组信息
  const getGroupInfo = (groupId: string) => {
    return groups.find((g: any) => g.id === groupId);
  };

  // 处理添加收藏
  const handleAddFavorite = async () => {
    try {
      const defaultGroup = groups.find((g: any) => g.name === '默认分组') || groups[0];

      // 默认使用最近一年的数据
      const endDate = new Date();
      const startDate = new Date();
      startDate.setFullYear(startDate.getFullYear() - 1);

      const formatDate = (date: Date) => {
        return date.toISOString().split('T')[0];
      };

      await addFavorite({
        ts_code: stockCode,
        name: stockName,
        start_date: formatDate(startDate),
        end_date: formatDate(endDate),
        group_id: defaultGroup?.id || 'default',
      }).unwrap();

      refetchFavorites();
    } catch (error: any) {
      console.error('添加收藏失败:', error);
      alert(error?.data?.error || '添加收藏失败');
    }
  };

  // 处理删除收藏
  const handleDeleteFavorite = async () => {
    if (!currentFavorite) return;

    try {
      await deleteFavorite(currentFavorite.id).unwrap();
      refetchFavorites();
      setDeleteDialogOpen(false);
    } catch (error: any) {
      console.error('删除失败:', error);
      alert(error?.data?.error || '删除失败');
    }
  };

  // 处理刷新
  const handleRefresh = () => {
    refetchFavorites();
    refetchPredictions();
    refetchPatternSummary();
  };

  // 获取信号类型文本
  const getSignalText = (signalType: string) => {
    switch (signalType) {
      case 'BUY':
        return '买入';
      case 'SELL':
        return '卖出';
      case 'HOLD':
        return '持有';
      default:
        return signalType;
    }
  };

  // 获取信号颜色
  const getSignalColor = (signalType: string) => {
    switch (signalType) {
      case 'BUY':
        return 'error'; // 红色
      case 'SELL':
        return 'success'; // 绿色
      case 'HOLD':
        return 'default'; // 灰色
      default:
        return 'default';
    }
  };

  // 加载状态
  if (favoritesLoading || groupsLoading) {
    return (
      <Box display="flex" justifyContent="center" py={4}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      {/* 顶部工具栏 */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
        <Typography variant="h6">
          策略管理
        </Typography>
        <Box>
          <IconButton onClick={handleRefresh} size="small">
            <RefreshIcon />
          </IconButton>
        </Box>
      </Box>

      {/* Tab导航 */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        <Tabs value={selectedTab} onChange={(_, newValue) => setSelectedTab(newValue)}>
          <Tab label="策略信息" />
          <Tab label="信号汇总" />
        </Tabs>
      </Box>

      {/* 策略信息Tab */}
      <TabPanel value={selectedTab} index={0}>
        {!isFavorited ? (
          <Card>
            <CardContent>
              <Box textAlign="center" py={4}>
                <BookmarkBorderIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
                <Typography variant="h6" gutterBottom>
                  未收藏此股票
                </Typography>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  添加到收藏后，可以跟踪该股票的策略信号
                </Typography>
                <Button
                  variant="contained"
                  startIcon={<BookmarkIcon />}
                  onClick={handleAddFavorite}
                  disabled={addLoading}
                  sx={{ mt: 2 }}
                >
                  添加到收藏
                </Button>
              </Box>
            </CardContent>
          </Card>
        ) : (
          <Card>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                <Typography variant="h6">收藏信息</Typography>
                <IconButton
                  color="error"
                  size="small"
                  onClick={() => setDeleteDialogOpen(true)}
                  disabled={deleteLoading}
                >
                  <DeleteIcon />
                </IconButton>
              </Box>

              <TableContainer>
                <Table>
                  <TableBody>
                    <TableRow>
                      <TableCell component="th" scope="row" sx={{ width: '30%' }}>
                        股票代码
                      </TableCell>
                      <TableCell>{currentFavorite.ts_code}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell component="th" scope="row">
                        股票名称
                      </TableCell>
                      <TableCell>{currentFavorite.name}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell component="th" scope="row">
                        分析周期
                      </TableCell>
                      <TableCell>
                        {currentFavorite.start_date} 至 {currentFavorite.end_date}
                      </TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell component="th" scope="row">
                        所属分组
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={getGroupInfo(currentFavorite.group_id)?.name || '默认分组'}
                          size="small"
                          sx={{
                            bgcolor: getGroupInfo(currentFavorite.group_id)?.color || '#1976d2',
                            color: '#fff',
                          }}
                        />
                      </TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell component="th" scope="row">
                        添加时间
                      </TableCell>
                      <TableCell>
                        {new Date(currentFavorite.created_at).toLocaleString('zh-CN')}
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        )}
      </TabPanel>

      {/* 信号汇总Tab */}
      <TabPanel value={selectedTab} index={1}>
        {predictionsLoading || patternSummaryLoading ? (
          <Box display="flex" justifyContent="center" py={4}>
            <CircularProgress />
          </Box>
        ) : (
          <>
            {/* 预测信号 */}
            {predictions.length > 0 && (
              <Card sx={{ mb: 3 }}>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    预测信号
                  </Typography>
                  <TableContainer component={Paper} variant="outlined">
                    <Table>
                      <TableHead>
                        <TableRow>
                          <TableCell>信号类型</TableCell>
                          <TableCell>置信度</TableCell>
                          <TableCell>理由</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {predictions.map((pred: any, index: number) => (
                          <TableRow key={index}>
                            <TableCell>
                              <Chip
                                label={getSignalText(pred.type)}
                                color={getSignalColor(pred.type) as any}
                                size="small"
                              />
                            </TableCell>
                            <TableCell>
                              {(pred.probability * 100).toFixed(1)}%
                            </TableCell>
                            <TableCell>{pred.reason || '无详细理由'}</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </TableContainer>
                </CardContent>
              </Card>
            )}

            {/* 形态识别 */}
            {(patternSummary?.patterns && Object.keys(patternSummary.patterns).length > 0) && (
              <Card sx={{ mb: 3 }}>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    形态识别统计
                  </Typography>
                  
                  {/* 统计周期信息 */}
                  {patternSummary.start_date && patternSummary.end_date && (
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      统计周期: {patternSummary.start_date} - {patternSummary.end_date} ({patternSummary.period}天)
                    </Typography>
                  )}
                  
                  <Divider sx={{ my: 2 }} />
                  
                  {/* 形态模式统计 */}
                  <Typography variant="subtitle1" gutterBottom sx={{ mt: 2 }}>
                    形态模式
                  </Typography>
                  <TableContainer component={Paper} variant="outlined" sx={{ mb: 2 }}>
                    <Table size="small">
                      <TableHead>
                        <TableRow>
                          <TableCell>形态名称</TableCell>
                          <TableCell align="right">出现次数</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {Object.entries(patternSummary.patterns || {}).map(([key, value]: [string, any]) => (
                          <TableRow key={key}>
                            <TableCell>{key}</TableCell>
                            <TableCell align="right">
                              <Chip 
                                label={value} 
                                size="small" 
                                color="primary" 
                                variant="outlined"
                              />
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </TableContainer>

                  {/* 信号统计 */}
                  {patternSummary.signals && Object.keys(patternSummary.signals).length > 0 && (
                    <>
                      <Typography variant="subtitle1" gutterBottom>
                        信号统计
                      </Typography>
                      <TableContainer component={Paper} variant="outlined">
                        <Table size="small">
                          <TableHead>
                            <TableRow>
                              <TableCell>信号类型</TableCell>
                              <TableCell align="right">出现次数</TableCell>
                            </TableRow>
                          </TableHead>
                          <TableBody>
                            {Object.entries(patternSummary.signals || {}).map(([key, value]: [string, any]) => (
                              <TableRow key={key}>
                                <TableCell>
                                  <Chip
                                    label={key}
                                    color={
                                      key === 'BUY' ? 'success' : 
                                      key === 'SELL' ? 'error' : 
                                      'default'
                                    }
                                    size="small"
                                  />
                                </TableCell>
                                <TableCell align="right">
                                  <Typography variant="body2" fontWeight="bold">
                                    {value}
                                  </Typography>
                                </TableCell>
                              </TableRow>
                            ))}
                          </TableBody>
                        </Table>
                      </TableContainer>
                    </>
                  )}
                </CardContent>
              </Card>
            )}

            {/* 空状态 */}
            {predictions.length === 0 && 
             (!patternSummary || !patternSummary.patterns || Object.keys(patternSummary.patterns).length === 0) && (
              <Alert severity="info">
                暂无信号数据，请稍后刷新查看最新的预测信号
              </Alert>
            )}
          </>
        )}
      </TabPanel>

      {/* 删除确认对话框 */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>确认删除</DialogTitle>
        <DialogContent>
          <Typography>
            确定要将 <strong>{stockName}</strong> 从收藏中移除吗？
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>取消</Button>
          <Button
            onClick={handleDeleteFavorite}
            color="error"
            variant="contained"
            disabled={deleteLoading}
          >
            删除
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default StrategyView;

