/**
 * 股票详情页面 - Web版本
 */

import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Container,
  Box,
  Typography,
  Card,
  CardContent,
  Tabs,
  Tab,
  Chip,
  IconButton,
  CircularProgress,
  Alert,
  Paper,
  Button,
  ButtonGroup,
  Snackbar,
} from '@mui/material';
import {
  ArrowBack as BackIcon,
  Refresh as RefreshIcon,
  TrendingUp,
  TrendingDown,
  FavoriteBorder as FavoriteBorderIcon,
  Favorite as FavoriteIcon,
} from '@mui/icons-material';

import {
  useGetStockBasicQuery,
  useGetDailyDataQuery,
  useGetIndicatorsQuery,
  useGetFundamentalDataQuery,
  useCheckFavoriteQuery,
  useGetFavoritesQuery,
  useAddFavoriteMutation,
  useDeleteFavoriteMutation,
} from '../services/api';

import KLineChart from '../components/charts/KLineChart';
import TechnicalIndicatorsView from '../components/stock/TechnicalIndicatorsView';
import PredictionSignalsView from '../components/stock/PredictionSignalsView';
import FundamentalDataView from '../components/stock/FundamentalDataView';

// 标签页类型
type TabValue = 'kline' | 'indicators' | 'predictions' | 'fundamental';

// 时间范围选项
const TIME_RANGES = [
  { label: '1个月', days: 30 },
  { label: '3个月', days: 90 },
  { label: '6个月', days: 180 },
  { label: '1年', days: 365 },
  { label: '2年', days: 730 },
];

const StockDetailPage: React.FC = () => {
  const { stockCode } = useParams<{ stockCode: string }>();
  const navigate = useNavigate();

  // 状态
  const [selectedTab, setSelectedTab] = useState<TabValue>('kline');
  const [timeRange, setTimeRange] = useState(90); // 默认3个月
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState('');

  // 获取日期范围
  const getDateRange = () => {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(startDate.getDate() - timeRange);

    const formatDate = (date: Date) => {
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}${month}${day}`;
    };

    return {
      startDate: formatDate(startDate),
      endDate: formatDate(endDate),
    };
  };

  const { startDate, endDate } = getDateRange();

  // API查询
  const {
    data: basicData,
    isLoading: isBasicLoading,
    error: basicError,
    refetch: refetchBasic,
  } = useGetStockBasicQuery(stockCode || '', { skip: !stockCode });

  const {
    data: dailyData,
    isLoading: isDailyLoading,
    error: dailyError,
    refetch: refetchDaily,
  } = useGetDailyDataQuery(
    { stockCode: stockCode || '', startDate, endDate },
    { skip: !stockCode }
  );

  const {
    data: indicatorsData,
    isLoading: isIndicatorsLoading,
    refetch: refetchIndicators,
  } = useGetIndicatorsQuery(stockCode || '', { skip: !stockCode || selectedTab !== 'indicators' });

  const {
    data: fundamentalData,
    isLoading: isFundamentalLoading,
    refetch: refetchFundamental,
  } = useGetFundamentalDataQuery(
    { stockCode: stockCode || '' },
    { skip: !stockCode || selectedTab !== 'fundamental' }
  );

  const {
    data: favoriteCheck,
    refetch: refetchFavorite,
  } = useCheckFavoriteQuery(stockCode || '', { skip: !stockCode });

  const {
    data: favoritesData,
  } = useGetFavoritesQuery();

  const [addFavorite, { isLoading: isAddingFavorite }] = useAddFavoriteMutation();
  const [deleteFavorite, { isLoading: isDeletingFavorite }] = useDeleteFavoriteMutation();

  // 处理收藏
  const handleToggleFavorite = async () => {
    if (!stockCode || !basicData?.data) return;

    try {
      if (favoriteCheck?.data?.is_favorite) {
        // 取消收藏 - 从收藏列表中找到对应的ID
        const favoriteItem = favoritesData?.data?.find(
          (fav) => fav.ts_code === stockCode
        );
        
        if (favoriteItem) {
          await deleteFavorite(favoriteItem.id).unwrap();
          refetchFavorite();
          setSnackbarMessage('已取消收藏');
          setSnackbarOpen(true);
        } else {
          console.error('未找到收藏记录');
          setSnackbarMessage('取消收藏失败：未找到收藏记录');
          setSnackbarOpen(true);
        }
      } else {
        // 添加收藏
        await addFavorite({
          ts_code: stockCode,
          name: basicData.data.name,
        }).unwrap();
        refetchFavorite();
        setSnackbarMessage('已添加到收藏');
        setSnackbarOpen(true);
      }
    } catch (error) {
      console.error('收藏操作失败:', error);
      setSnackbarMessage('收藏操作失败，请重试');
      setSnackbarOpen(true);
    }
  };

  // 处理Snackbar关闭
  const handleSnackbarClose = () => {
    setSnackbarOpen(false);
  };

  // 处理刷新
  const handleRefresh = () => {
    refetchBasic();
    refetchDaily();
    if (selectedTab === 'indicators') refetchIndicators();
    if (selectedTab === 'fundamental') refetchFundamental();
  };

  // 计算价格变化
  const getPriceChange = () => {
    if (!dailyData?.data || dailyData.data.length === 0) {
      return { change: 0, changePercent: 0, latestPrice: 0 };
    }

    const data = dailyData.data;
    const latest = data[data.length - 1];
    const previous = data.length > 1 ? data[data.length - 2] : null;

    const latestPrice = parseFloat(latest.close);

    if (!previous) {
      return { change: 0, changePercent: 0, latestPrice };
    }

    const previousPrice = parseFloat(previous.close);
    const change = latestPrice - previousPrice;
    const changePercent = (change / previousPrice) * 100;

    return { change, changePercent, latestPrice };
  };

  const { change, changePercent, latestPrice } = getPriceChange();
  const isPositive = change >= 0;
  const isFavorite = favoriteCheck?.data?.is_favorite || false;

  // 根据涨跌幅度获取颜色（红涨绿跌，涨跌幅越大颜色越深）
  const getPriceColor = () => {
    const absPercent = Math.abs(changePercent);
    
    if (isPositive) {
      // 上涨 - 红色系，涨幅越大越深
      if (absPercent >= 9) return '#d32f2f'; // 深红（接近涨停）
      if (absPercent >= 5) return '#e53935'; // 中深红
      if (absPercent >= 3) return '#f44336'; // 标准红
      if (absPercent >= 1) return '#ef5350'; // 浅红
      return '#e57373'; // 微涨红
    } else {
      // 下跌 - 绿色系，跌幅越大越深
      if (absPercent >= 9) return '#1b5e20'; // 深绿（接近跌停）
      if (absPercent >= 5) return '#2e7d32'; // 中深绿
      if (absPercent >= 3) return '#388e3c'; // 标准绿
      if (absPercent >= 1) return '#43a047'; // 浅绿
      return '#66bb6a'; // 微跌绿
    }
  };

  const priceColor = getPriceColor();

  // 加载状态
  if (isBasicLoading || isDailyLoading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
          <Box textAlign="center">
            <CircularProgress size={60} />
            <Typography variant="body1" sx={{ mt: 2 }}>
              正在加载股票数据...
            </Typography>
          </Box>
        </Box>
      </Container>
    );
  }

  // 错误状态
  if (basicError || dailyError) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error" action={
          <Button color="inherit" onClick={handleRefresh}>
            重试
          </Button>
        }>
          加载失败: {basicError?.toString() || dailyError?.toString()}
        </Alert>
      </Container>
    );
  }

  const stockBasic = basicData?.data;

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      {/* 顶部导航栏 */}
      <Box display="flex" alignItems="center" mb={3}>
        <IconButton onClick={() => navigate(-1)} sx={{ mr: 2 }}>
          <BackIcon />
        </IconButton>
        <Typography variant="h5" sx={{ flexGrow: 1 }}>
          股票详情
        </Typography>
        <IconButton onClick={handleRefresh}>
          <RefreshIcon />
        </IconButton>
      </Box>

      {/* 股票头部信息卡片 */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: { xs: '1fr', md: '2fr 1fr' },
              gap: 2,
            }}
          >
            <Box display="flex" alignItems="flex-start" justifyContent="space-between">
              <Box>
                <Typography variant="h4" component="h1" gutterBottom>
                  {stockBasic?.name || '未知'}
                </Typography>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  {stockCode}
                </Typography>
                <Box display="flex" gap={1} mt={1}>
                  {stockBasic?.industry && (
                    <Chip label={stockBasic.industry} size="small" />
                  )}
                  {stockBasic?.market && (
                    <Chip label={stockBasic.market} size="small" variant="outlined" />
                  )}
                </Box>
              </Box>
              <IconButton
                onClick={handleToggleFavorite}
                disabled={isAddingFavorite || isDeletingFavorite}
                color={isFavorite ? 'error' : 'default'}
              >
                {isFavorite ? <FavoriteIcon /> : <FavoriteBorderIcon />}
              </IconButton>
            </Box>
            <Paper elevation={0} sx={{ bgcolor: isPositive ? 'rgba(239, 83, 80, 0.08)' : 'rgba(38, 166, 154, 0.08)', p: 2 }}>
              <Typography variant="h3" align="right" sx={{ color: priceColor, fontWeight: 'bold' }}>
                ¥{latestPrice.toFixed(2)}
              </Typography>
              <Box display="flex" justifyContent="flex-end" alignItems="center" mt={1}>
                {isPositive ? <TrendingUp sx={{ color: priceColor }} /> : <TrendingDown sx={{ color: priceColor }} />}
                <Typography variant="h6" sx={{ ml: 1, color: priceColor, fontWeight: 'bold' }}>
                  {isPositive ? '+' : ''}{change.toFixed(2)} ({isPositive ? '+' : ''}{changePercent.toFixed(2)}%)
                </Typography>
              </Box>
            </Paper>
          </Box>
        </CardContent>
      </Card>

      {/* 时间范围选择 */}
      <Box display="flex" justifyContent="center" mb={2}>
        <ButtonGroup size="small" variant="outlined">
          {TIME_RANGES.map((range) => (
            <Button
              key={range.days}
              onClick={() => setTimeRange(range.days)}
              variant={timeRange === range.days ? 'contained' : 'outlined'}
            >
              {range.label}
            </Button>
          ))}
        </ButtonGroup>
      </Box>

      {/* 标签页 */}
      <Card>
        <Tabs
          value={selectedTab}
          onChange={(_, newValue) => setSelectedTab(newValue)}
          variant="fullWidth"
          sx={{ borderBottom: 1, borderColor: 'divider' }}
        >
          <Tab label="K线图" value="kline" />
          <Tab label="技术指标" value="indicators" />
          <Tab label="买卖预测" value="predictions" />
          <Tab label="基本面" value="fundamental" />
        </Tabs>

        <CardContent>
          {/* K线图标签页 */}
          {selectedTab === 'kline' && (
            <Box>
              {dailyData?.data && dailyData.data.length > 0 ? (
                <KLineChart data={dailyData.data} />
              ) : (
                <Typography variant="body1" color="text.secondary" align="center" py={4}>
                  暂无K线数据
                </Typography>
              )}
            </Box>
          )}

          {/* 技术指标标签页 */}
          {selectedTab === 'indicators' && (
            <Box>
              {isIndicatorsLoading ? (
                <Box display="flex" justifyContent="center" py={4}>
                  <CircularProgress />
                </Box>
              ) : indicatorsData?.data && indicatorsData.data.length > 0 ? (
                <TechnicalIndicatorsView data={indicatorsData.data} />
              ) : (
                <Typography variant="body1" color="text.secondary" align="center" py={4}>
                  暂无技术指标数据
                </Typography>
              )}
            </Box>
          )}

          {/* 买卖预测标签页 */}
          {selectedTab === 'predictions' && (
            <Box>
              <PredictionSignalsView stockCode={stockCode || ''} />
            </Box>
          )}

          {/* 基本面标签页 */}
          {selectedTab === 'fundamental' && (
            <Box>
              {isFundamentalLoading ? (
                <Box display="flex" justifyContent="center" py={4}>
                  <CircularProgress />
                </Box>
              ) : fundamentalData?.data ? (
                <FundamentalDataView data={fundamentalData.data} />
              ) : (
                <Typography variant="body1" color="text.secondary" align="center" py={4}>
                  暂无基本面数据
                </Typography>
              )}
            </Box>
          )}
        </CardContent>
      </Card>

      {/* 消息提示 */}
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={3000}
        onClose={handleSnackbarClose}
        message={snackbarMessage}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      />
    </Container>
  );
};

export default StockDetailPage;

