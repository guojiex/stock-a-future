/**
 * 策略管理页面
 * 参考网页版 web/static/js/modules/strategies.js 实现
 * 提供策略的完整CRUD功能和性能展示
 */

import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import {
  Box,
  Container,
  Typography,
  Card,
  CardContent,
  CardActions,
  Button,
  Chip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Alert,
  CircularProgress,
  Tooltip,
  Paper,
  Divider,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  PlayArrow as PlayArrowIcon,
  Pause as PauseIcon,
  Visibility as VisibilityIcon,
  ContentCopy as ContentCopyIcon,
  Assessment as AssessmentIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material';
import {
  useGetStrategiesQuery,
  useGetStrategyQuery,
  useCreateStrategyMutation,
  useUpdateStrategyMutation,
  useDeleteStrategyMutation,
  useToggleStrategyMutation,
} from '../services/api';
import { setSelectedStrategies } from '../store/slices/backtestSlice';

// 策略类型映射
const STRATEGY_TYPES = {
  technical: '技术指标',
  fundamental: '基本面',
  ml: '机器学习',
  composite: '复合策略',
};

// 策略状态映射
const STRATEGY_STATUS = {
  active: { label: '活跃', color: 'success' as const },
  inactive: { label: '非活跃', color: 'default' as const },
  testing: { label: '测试中', color: 'warning' as const },
};

interface Strategy {
  id: string;
  name: string;
  description: string;
  strategy_type: keyof typeof STRATEGY_TYPES;
  status: keyof typeof STRATEGY_STATUS;
  parameters: Record<string, any>;
  created_at: string;
  updated_at?: string;
}

interface StrategyPerformance {
  total_return: number;
  annual_return: number;
  max_drawdown: number;
  sharpe_ratio: number;
  win_rate: number;
  total_trades: number;
  last_updated: string;
}

const StrategiesPage: React.FC = () => {
  // Hooks
  const navigate = useNavigate();
  const dispatch = useDispatch();
  
  // 状态
  const [selectedStrategy, setSelectedStrategy] = useState<Strategy | null>(null);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);
  const [editFormData, setEditFormData] = useState<Partial<Strategy>>({});

  // API查询
  const {
    data: strategiesData,
    isLoading: strategiesLoading,
    error: strategiesError,
    refetch: refetchStrategies,
  } = useGetStrategiesQuery();

  // API突变
  const [updateStrategy, { isLoading: updating }] = useUpdateStrategyMutation();
  const [deleteStrategy, { isLoading: deleting }] = useDeleteStrategyMutation();
  const [toggleStrategy, { isLoading: toggling }] = useToggleStrategyMutation();

  // 处理策略数据并按ID排序（确保显示顺序稳定）
  const strategies: Strategy[] = React.useMemo(() => {
    const data = (strategiesData?.data?.items || strategiesData?.data?.data || []) as Strategy[];
    // 按ID字母顺序排序，确保每次显示顺序一致
    return [...data].sort((a, b) => a.id.localeCompare(b.id));
  }, [strategiesData]);

  // 打开编辑对话框
  const handleEditClick = (strategy: Strategy) => {
    setSelectedStrategy(strategy);
    setEditFormData(strategy);
    setEditDialogOpen(true);
  };

  // 关闭编辑对话框
  const handleEditClose = () => {
    setEditDialogOpen(false);
    setSelectedStrategy(null);
    setEditFormData({});
  };

  // 保存策略修改
  const handleEditSave = async () => {
    if (!selectedStrategy) return;

    try {
      await updateStrategy({
        id: selectedStrategy.id,
        ...editFormData,
      }).unwrap();

      handleEditClose();
      refetchStrategies();
    } catch (error) {
      console.error('更新策略失败:', error);
    }
  };

  // 打开删除确认对话框
  const handleDeleteClick = (strategy: Strategy) => {
    setSelectedStrategy(strategy);
    setDeleteDialogOpen(true);
  };

  // 确认删除策略
  const handleDeleteConfirm = async () => {
    if (!selectedStrategy) return;

    try {
      await deleteStrategy(selectedStrategy.id).unwrap();
      setDeleteDialogOpen(false);
      setSelectedStrategy(null);
      refetchStrategies();
    } catch (error) {
      console.error('删除策略失败:', error);
    }
  };

  // 切换策略状态
  const handleToggleStrategy = async (strategy: Strategy) => {
    try {
      await toggleStrategy(strategy.id).unwrap();
      refetchStrategies();
    } catch (error) {
      console.error('切换策略状态失败:', error);
    }
  };

  // 查看策略详情
  const handleViewDetails = (strategy: Strategy) => {
    setSelectedStrategy(strategy);
    setDetailsDialogOpen(true);
  };

  // 运行回测
  const handleRunBacktest = (strategy: Strategy) => {
    console.log('运行回测:', strategy.name);
    
    // 参考网页版的实现：
    // 1. 设置选中的策略ID
    dispatch(setSelectedStrategies([strategy.id]));
    
    // 2. 导航到回测页面
    navigate('/backtest');
  };

  // 渲染策略卡片
  const renderStrategyCard = (strategy: Strategy) => {
    const statusInfo = STRATEGY_STATUS[strategy.status];

    return (
      <Box key={strategy.id}>
        <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
          <CardContent sx={{ flexGrow: 1 }}>
            {/* 策略标题和状态 */}
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
              <Typography variant="h6" component="h3" sx={{ flexGrow: 1 }}>
                {strategy.name}
              </Typography>
              <Chip
                label={statusInfo.label}
                color={statusInfo.color}
                size="small"
              />
            </Box>

            {/* 策略类型 */}
            <Chip
              label={STRATEGY_TYPES[strategy.strategy_type]}
              size="small"
              variant="outlined"
              sx={{ mb: 2 }}
            />

            {/* 策略描述 */}
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              {strategy.description || '暂无描述'}
            </Typography>

            {/* 策略参数 */}
            {strategy.parameters && Object.keys(strategy.parameters).length > 0 && (
              <Box sx={{ mt: 2 }}>
                <Typography variant="subtitle2" gutterBottom>
                  参数配置:
                </Typography>
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                  {Object.entries(strategy.parameters).slice(0, 3).map(([key, value]) => (
                    <Chip
                      key={key}
                      label={`${key}: ${value}`}
                      size="small"
                      variant="outlined"
                    />
                  ))}
                  {Object.keys(strategy.parameters).length > 3 && (
                    <Chip
                      label={`+${Object.keys(strategy.parameters).length - 3}更多`}
                      size="small"
                      variant="outlined"
                    />
                  )}
                </Box>
              </Box>
            )}

            {/* 创建时间 */}
            <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 2 }}>
              创建于: {new Date(strategy.created_at).toLocaleDateString('zh-CN')}
            </Typography>
          </CardContent>

          {/* 操作按钮 */}
          <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
            <Box>
              <Tooltip title="查看详情">
                <IconButton size="small" onClick={() => handleViewDetails(strategy)}>
                  <VisibilityIcon fontSize="small" />
                </IconButton>
              </Tooltip>
              <Tooltip title="编辑">
                <IconButton size="small" onClick={() => handleEditClick(strategy)}>
                  <EditIcon fontSize="small" />
                </IconButton>
              </Tooltip>
              <Tooltip title={strategy.status === 'active' ? '暂停' : '启用'}>
                <IconButton
                  size="small"
                  onClick={() => handleToggleStrategy(strategy)}
                  disabled={toggling}
                >
                  {strategy.status === 'active' ? (
                    <PauseIcon fontSize="small" />
                  ) : (
                    <PlayArrowIcon fontSize="small" />
                  )}
                </IconButton>
              </Tooltip>
              <Tooltip title="删除">
                <IconButton
                  size="small"
                  color="error"
                  onClick={() => handleDeleteClick(strategy)}
                >
                  <DeleteIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </Box>

            <Button
              size="small"
              variant="contained"
              startIcon={<AssessmentIcon />}
              onClick={() => handleRunBacktest(strategy)}
            >
              运行回测
            </Button>
          </CardActions>
        </Card>
      </Box>
    );
  };

  // 渲染编辑对话框
  const renderEditDialog = () => (
    <Dialog
      open={editDialogOpen}
      onClose={handleEditClose}
      maxWidth="md"
      fullWidth
    >
      <DialogTitle>编辑策略</DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
          {/* 策略名称 */}
          <TextField
            label="策略名称"
            value={editFormData.name || ''}
            onChange={(e) => setEditFormData({ ...editFormData, name: e.target.value })}
            fullWidth
            required
          />

          {/* 策略类型 */}
          <FormControl fullWidth disabled>
            <InputLabel>策略类型</InputLabel>
            <Select
              value={editFormData.strategy_type || ''}
              label="策略类型"
            >
              {Object.entries(STRATEGY_TYPES).map(([key, label]) => (
                <MenuItem key={key} value={key}>
                  {label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          {/* 策略状态 */}
          <FormControl fullWidth>
            <InputLabel>策略状态</InputLabel>
            <Select
              value={editFormData.status || ''}
              label="策略状态"
              onChange={(e) => setEditFormData({ ...editFormData, status: e.target.value as any })}
            >
              {Object.entries(STRATEGY_STATUS).map(([key, { label }]) => (
                <MenuItem key={key} value={key}>
                  {label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          {/* 策略描述 */}
          <TextField
            label="策略描述"
            value={editFormData.description || ''}
            onChange={(e) => setEditFormData({ ...editFormData, description: e.target.value })}
            fullWidth
            multiline
            rows={3}
          />

          {/* 参数配置 */}
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              参数配置
            </Typography>
            {renderParameterFields()}
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleEditClose}>取消</Button>
        <Button onClick={handleEditSave} variant="contained" disabled={updating}>
          {updating ? <CircularProgress size={24} /> : '保存'}
        </Button>
      </DialogActions>
    </Dialog>
  );

  // 渲染参数字段(根据策略类型动态生成)
  const renderParameterFields = () => {
    if (!selectedStrategy) return null;

    const parameters = editFormData.parameters || {};

    // 根据策略ID或类型渲染不同的参数字段
    // 这里简化处理,实际应该根据策略类型动态生成
    return (
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        {Object.entries(parameters).map(([key, value]) => (
          <TextField
            key={key}
            label={key}
            value={value}
            onChange={(e) => {
              const newValue = typeof value === 'number'
                ? parseFloat(e.target.value) || 0
                : e.target.value;
              setEditFormData({
                ...editFormData,
                parameters: {
                  ...parameters,
                  [key]: newValue,
                },
              });
            }}
            type={typeof value === 'number' ? 'number' : 'text'}
            fullWidth
            size="small"
          />
        ))}
      </Box>
    );
  };

  // 渲染策略详情对话框
  const renderDetailsDialog = () => {
    if (!selectedStrategy) return null;

    return (
      <Dialog
        open={detailsDialogOpen}
        onClose={() => setDetailsDialogOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h6">{selectedStrategy.name}</Typography>
            <Chip
              label={STRATEGY_STATUS[selectedStrategy.status].label}
              color={STRATEGY_STATUS[selectedStrategy.status].color}
              size="small"
            />
          </Box>
        </DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
            {/* 基本信息 */}
            <Paper sx={{ p: 2 }}>
              <Typography variant="subtitle1" gutterBottom fontWeight="bold">
                基本信息
              </Typography>
              <Divider sx={{ mb: 2 }} />
              <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 2 }}>
                <Box>
                  <Typography variant="body2" color="text.secondary">
                    策略类型
                  </Typography>
                  <Typography variant="body1">
                    {STRATEGY_TYPES[selectedStrategy.strategy_type]}
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="body2" color="text.secondary">
                    创建时间
                  </Typography>
                  <Typography variant="body1">
                    {new Date(selectedStrategy.created_at).toLocaleString('zh-CN')}
                  </Typography>
                </Box>
                <Box sx={{ gridColumn: '1 / -1' }}>
                  <Typography variant="body2" color="text.secondary">
                    描述
                  </Typography>
                  <Typography variant="body1">
                    {selectedStrategy.description || '暂无描述'}
                  </Typography>
                </Box>
              </Box>
            </Paper>

            {/* 参数配置 */}
            {selectedStrategy.parameters && Object.keys(selectedStrategy.parameters).length > 0 && (
              <Paper sx={{ p: 2 }}>
                <Typography variant="subtitle1" gutterBottom fontWeight="bold">
                  参数配置
                </Typography>
                <Divider sx={{ mb: 2 }} />
                <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 2 }}>
                  {Object.entries(selectedStrategy.parameters).map(([key, value]) => (
                    <Box key={key}>
                      <Typography variant="body2" color="text.secondary">
                        {key}
                      </Typography>
                      <Typography variant="body1">{String(value)}</Typography>
                    </Box>
                  ))}
                </Box>
              </Paper>
            )}

            {/* TODO: 添加策略表现指标 */}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDetailsDialogOpen(false)}>关闭</Button>
          <Button
            variant="contained"
            startIcon={<EditIcon />}
            onClick={() => {
              setDetailsDialogOpen(false);
              handleEditClick(selectedStrategy);
            }}
          >
            编辑策略
          </Button>
          <Button
            variant="contained"
            color="primary"
            startIcon={<AssessmentIcon />}
            onClick={() => {
              setDetailsDialogOpen(false);
              handleRunBacktest(selectedStrategy);
            }}
          >
            运行回测
          </Button>
        </DialogActions>
      </Dialog>
    );
  };

  // 渲染删除确认对话框
  const renderDeleteDialog = () => (
    <Dialog
      open={deleteDialogOpen}
      onClose={() => setDeleteDialogOpen(false)}
    >
      <DialogTitle>确认删除</DialogTitle>
      <DialogContent>
        <Typography>
          确定要删除策略 "{selectedStrategy?.name}" 吗？此操作不可撤销。
        </Typography>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => setDeleteDialogOpen(false)}>取消</Button>
        <Button
          onClick={handleDeleteConfirm}
          color="error"
          variant="contained"
          disabled={deleting}
        >
          {deleting ? <CircularProgress size={24} /> : '删除'}
        </Button>
      </DialogActions>
    </Dialog>
  );

  // 加载状态
  if (strategiesLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '60vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  // 错误状态
  if (strategiesError) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Alert severity="error">
          加载策略列表失败，请稍后重试
        </Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      {/* 页面标题和操作栏 */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Typography variant="h4" component="h1" fontWeight="bold">
          📋 策略管理
        </Typography>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<RefreshIcon />}
            onClick={() => refetchStrategies()}
          >
            刷新
          </Button>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {
              // TODO: 打开创建策略对话框
              console.log('创建新策略');
            }}
          >
            创建策略
          </Button>
        </Box>
      </Box>

      {/* 策略列表 */}
      {strategies.length === 0 ? (
        <Paper sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h6" color="text.secondary" gutterBottom>
            暂无策略
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            点击"创建策略"按钮开始创建您的第一个交易策略
          </Typography>
          <Button variant="contained" startIcon={<AddIcon />}>
            创建策略
          </Button>
        </Paper>
      ) : (
        <Box
          sx={{
            display: 'grid',
            gridTemplateColumns: {
              xs: '1fr',
              sm: 'repeat(2, 1fr)',
              md: 'repeat(3, 1fr)',
            },
            gap: 3,
          }}
        >
          {strategies.map((strategy) => renderStrategyCard(strategy))}
        </Box>
      )}

      {/* 对话框 */}
      {renderEditDialog()}
      {renderDetailsDialog()}
      {renderDeleteDialog()}
    </Container>
  );
};

export default StrategiesPage;

