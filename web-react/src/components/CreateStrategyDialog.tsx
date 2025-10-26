/**
 * 创建策略对话框组件
 * 提供分步表单创建量化交易策略
 */

import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Stepper,
  Step,
  StepLabel,
  Box,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  CircularProgress,
  IconButton,
  Typography,
  Paper,
  Divider,
  Chip,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';

interface CreateStrategyDialogProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

interface StrategyFormData {
  name: string;
  description: string;
  strategy_type: 'technical' | 'fundamental' | 'ml' | 'composite';
  parameters: Record<string, any>;
  code?: string;
}

interface FormErrors {
  name?: string;
  description?: string;
  parameters?: Record<string, string>;
}

const STRATEGY_TYPES = {
  technical: '技术指标策略',
  fundamental: '基本面策略',
  ml: '机器学习策略',
  composite: '复合策略',
};

const STRATEGY_TEMPLATES = [
  { value: '', label: '不使用模板' },
  { value: 'macd_strategy', label: 'MACD金叉策略' },
  { value: 'ma_crossover', label: '双均线策略' },
  { value: 'rsi_strategy', label: 'RSI超买超卖' },
  { value: 'bollinger_strategy', label: '布林带策略' },
];

const steps = ['基本信息', '策略参数', '确认创建'];

export const CreateStrategyDialog: React.FC<CreateStrategyDialogProps> = ({
  open,
  onClose,
  onSuccess,
}) => {
  const [activeStep, setActiveStep] = useState(0);
  const [loading, setLoading] = useState(false);
  const [selectedTemplate, setSelectedTemplate] = useState('');
  const [formData, setFormData] = useState<StrategyFormData>({
    name: '',
    description: '',
    strategy_type: 'technical',
    parameters: {},
  });
  const [errors, setErrors] = useState<FormErrors>({});

  // 重置表单
  const resetForm = () => {
    setActiveStep(0);
    setSelectedTemplate('');
    setFormData({
      name: '',
      description: '',
      strategy_type: 'technical',
      parameters: {},
    });
    setErrors({});
  };

  // 应用模板
  const handleTemplateSelect = async (templateId: string) => {
    setSelectedTemplate(templateId);
    
    if (!templateId) {
      return;
    }

    try {
      const response = await fetch('/api/v1/strategies/templates');
      const result = await response.json();
      
      if (result.success) {
        const template = result.data.find((t: any) => t.id === templateId);
        if (template) {
          setFormData(prev => ({
            ...prev,
            name: template.name,
            description: template.description,
            strategy_type: template.strategy_type,
            parameters: template.parameters || {},
          }));
        }
      }
    } catch (error) {
      console.error('加载模板失败:', error);
    }
  };

  // 验证当前步骤
  const validateStep = (step: number): boolean => {
    const newErrors: FormErrors = {};

    if (step === 0) {
      if (!formData.name.trim()) {
        newErrors.name = '策略名称不能为空';
      } else if (formData.name.length > 100) {
        newErrors.name = '策略名称不能超过100个字符';
      }

      if (formData.description.length > 1000) {
        newErrors.description = '策略描述不能超过1000个字符';
      }
    }

    if (step === 1) {
      // 验证参数
      const paramErrors: Record<string, string> = {};
      
      if (formData.strategy_type === 'technical') {
        const params = formData.parameters;
        
        // MACD参数验证
        if ('fast_period' in params) {
          const fastPeriod = Number(params.fast_period);
          if (isNaN(fastPeriod) || fastPeriod < 1 || fastPeriod > 50) {
            paramErrors.fast_period = '快线周期必须在1-50之间';
          }
        }
        
        if ('slow_period' in params) {
          const slowPeriod = Number(params.slow_period);
          if (isNaN(slowPeriod) || slowPeriod < 1 || slowPeriod > 100) {
            paramErrors.slow_period = '慢线周期必须在1-100之间';
          }
        }
        
        if ('fast_period' in params && 'slow_period' in params) {
          const fastPeriod = Number(params.fast_period);
          const slowPeriod = Number(params.slow_period);
          if (fastPeriod >= slowPeriod) {
            paramErrors.slow_period = '慢线周期必须大于快线周期';
          }
        }
      }
      
      if (Object.keys(paramErrors).length > 0) {
        newErrors.parameters = paramErrors;
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleNext = () => {
    if (validateStep(activeStep)) {
      setActiveStep((prev) => prev + 1);
    }
  };

  const handleBack = () => {
    setActiveStep((prev) => prev - 1);
  };

  const handleSubmit = async () => {
    if (!validateStep(activeStep)) {
      return;
    }

    setLoading(true);
    try {
      const response = await fetch('/api/v1/strategies', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });

      const result = await response.json();

      if (result.success) {
        onSuccess();
        onClose();
        resetForm();
      } else {
        setErrors({ name: result.error || '创建失败' });
      }
    } catch (error) {
      console.error('创建策略失败:', error);
      setErrors({ name: '网络错误，请稍后重试' });
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    onClose();
    // 延迟重置表单，避免关闭动画时看到表单重置
    setTimeout(resetForm, 300);
  };

  const handleParameterChange = (key: string, value: any) => {
    setFormData(prev => ({
      ...prev,
      parameters: {
        ...prev.parameters,
        [key]: value,
      },
    }));
    
    // 清除该参数的错误
    if (errors.parameters?.[key]) {
      const newParamErrors = { ...errors.parameters };
      delete newParamErrors[key];
      setErrors(prev => ({
        ...prev,
        parameters: Object.keys(newParamErrors).length > 0 ? newParamErrors : undefined,
      }));
    }
  };

  // 渲染基本信息步骤
  const renderBasicInfoStep = () => (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3, py: 2 }}>
      <TextField
        label="策略名称"
        placeholder="例如: 我的MACD策略"
        value={formData.name}
        onChange={(e) => {
          setFormData({ ...formData, name: e.target.value });
          if (errors.name) {
            setErrors(prev => ({ ...prev, name: undefined }));
          }
        }}
        fullWidth
        required
        error={!!errors.name}
        helperText={errors.name || '策略的唯一标识名称，最多100个字符'}
      />

      <TextField
        label="策略描述"
        placeholder="简要描述策略的交易逻辑和适用场景"
        value={formData.description}
        onChange={(e) => {
          setFormData({ ...formData, description: e.target.value });
          if (errors.description) {
            setErrors(prev => ({ ...prev, description: undefined }));
          }
        }}
        fullWidth
        multiline
        rows={4}
        error={!!errors.description}
        helperText={errors.description || '详细描述策略的目标和特点，最多1000个字符'}
      />

      <FormControl fullWidth>
        <InputLabel>使用模板(可选)</InputLabel>
        <Select
          value={selectedTemplate}
          onChange={(e) => handleTemplateSelect(e.target.value)}
          label="使用模板(可选)"
        >
          {STRATEGY_TEMPLATES.map(template => (
            <MenuItem key={template.value} value={template.value}>
              {template.label}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Box>
  );

  // 渲染策略参数步骤
  const renderParametersStep = () => (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3, py: 2 }}>
      <FormControl fullWidth required>
        <InputLabel>策略类型</InputLabel>
        <Select
          value={formData.strategy_type}
          onChange={(e) => setFormData({ ...formData, strategy_type: e.target.value as any })}
          label="策略类型"
        >
          {Object.entries(STRATEGY_TYPES).map(([key, label]) => (
            <MenuItem key={key} value={key}>
              {label}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      {formData.strategy_type === 'technical' && renderTechnicalParameters()}
      
      {formData.strategy_type !== 'technical' && (
        <Alert severity="info">
          {formData.strategy_type === 'fundamental' && '基本面策略参数配置开发中...'}
          {formData.strategy_type === 'ml' && '机器学习策略参数配置开发中...'}
          {formData.strategy_type === 'composite' && '复合策略参数配置开发中...'}
        </Alert>
      )}
    </Box>
  );

  // 渲染技术指标参数
  const renderTechnicalParameters = () => {
    // 根据模板判断显示哪些参数
    const isMACDStrategy = selectedTemplate === 'macd_strategy' || 'fast_period' in formData.parameters;
    const isMAStrategy = selectedTemplate === 'ma_crossover' || 'short_period' in formData.parameters;
    const isRSIStrategy = selectedTemplate === 'rsi_strategy' || ('period' in formData.parameters && 'overbought' in formData.parameters);
    const isBollingerStrategy = selectedTemplate === 'bollinger_strategy' || ('period' in formData.parameters && 'std_dev' in formData.parameters);

    if (isMACDStrategy) {
      return (
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <Typography variant="subtitle2" color="primary">MACD策略参数</Typography>
          <TextField
            label="快线周期"
            type="number"
            value={formData.parameters.fast_period || 12}
            onChange={(e) => handleParameterChange('fast_period', parseInt(e.target.value))}
            inputProps={{ min: 1, max: 50 }}
            fullWidth
            error={!!errors.parameters?.fast_period}
            helperText={errors.parameters?.fast_period || 'MACD快线的计算周期，通常为12天'}
          />
          <TextField
            label="慢线周期"
            type="number"
            value={formData.parameters.slow_period || 26}
            onChange={(e) => handleParameterChange('slow_period', parseInt(e.target.value))}
            inputProps={{ min: 1, max: 100 }}
            fullWidth
            error={!!errors.parameters?.slow_period}
            helperText={errors.parameters?.slow_period || 'MACD慢线的计算周期，通常为26天'}
          />
          <TextField
            label="信号线周期"
            type="number"
            value={formData.parameters.signal_period || 9}
            onChange={(e) => handleParameterChange('signal_period', parseInt(e.target.value))}
            inputProps={{ min: 1, max: 50 }}
            fullWidth
            helperText="信号线的计算周期，通常为9天"
          />
          <TextField
            label="买入阈值"
            type="number"
            value={formData.parameters.buy_threshold || 0}
            onChange={(e) => handleParameterChange('buy_threshold', parseFloat(e.target.value))}
            inputProps={{ min: -1, max: 1, step: 0.1 }}
            fullWidth
            helperText="MACD线超过此阈值时产生买入信号，默认为0"
          />
          <TextField
            label="卖出阈值"
            type="number"
            value={formData.parameters.sell_threshold || 0}
            onChange={(e) => handleParameterChange('sell_threshold', parseFloat(e.target.value))}
            inputProps={{ min: -1, max: 1, step: 0.1 }}
            fullWidth
            helperText="MACD线低于此阈值时产生卖出信号，默认为0"
          />
        </Box>
      );
    }

    if (isMAStrategy) {
      return (
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <Typography variant="subtitle2" color="primary">双均线策略参数</Typography>
          <TextField
            label="短期均线周期"
            type="number"
            value={formData.parameters.short_period || 5}
            onChange={(e) => handleParameterChange('short_period', parseInt(e.target.value))}
            inputProps={{ min: 1, max: 50 }}
            fullWidth
            helperText="短期移动平均线的计算周期，通常为5天"
          />
          <TextField
            label="长期均线周期"
            type="number"
            value={formData.parameters.long_period || 20}
            onChange={(e) => handleParameterChange('long_period', parseInt(e.target.value))}
            inputProps={{ min: 1, max: 200 }}
            fullWidth
            helperText="长期移动平均线的计算周期，通常为20天"
          />
          <FormControl fullWidth>
            <InputLabel>均线类型</InputLabel>
            <Select
              value={formData.parameters.ma_type || 'sma'}
              onChange={(e) => handleParameterChange('ma_type', e.target.value)}
              label="均线类型"
            >
              <MenuItem value="sma">简单移动平均(SMA)</MenuItem>
              <MenuItem value="ema">指数移动平均(EMA)</MenuItem>
              <MenuItem value="wma">加权移动平均(WMA)</MenuItem>
            </Select>
          </FormControl>
          <TextField
            label="突破阈值"
            type="number"
            value={formData.parameters.threshold || 0.01}
            onChange={(e) => handleParameterChange('threshold', parseFloat(e.target.value))}
            inputProps={{ min: 0, max: 0.1, step: 0.001 }}
            fullWidth
            helperText="均线突破的确认阈值，默认1%"
          />
        </Box>
      );
    }

    if (isRSIStrategy) {
      return (
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <Typography variant="subtitle2" color="primary">RSI策略参数</Typography>
          <TextField
            label="RSI周期"
            type="number"
            value={formData.parameters.period || 14}
            onChange={(e) => handleParameterChange('period', parseInt(e.target.value))}
            inputProps={{ min: 1, max: 50 }}
            fullWidth
            helperText="RSI指标的计算周期，通常为14天"
          />
          <TextField
            label="超买阈值"
            type="number"
            value={formData.parameters.overbought || 70}
            onChange={(e) => handleParameterChange('overbought', parseFloat(e.target.value))}
            inputProps={{ min: 50, max: 100 }}
            fullWidth
            helperText="RSI超过此值时认为超买，产生卖出信号，通常为70"
          />
          <TextField
            label="超卖阈值"
            type="number"
            value={formData.parameters.oversold || 30}
            onChange={(e) => handleParameterChange('oversold', parseFloat(e.target.value))}
            inputProps={{ min: 0, max: 50 }}
            fullWidth
            helperText="RSI低于此值时认为超卖，产生买入信号，通常为30"
          />
        </Box>
      );
    }

    if (isBollingerStrategy) {
      return (
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <Typography variant="subtitle2" color="primary">布林带策略参数</Typography>
          <TextField
            label="周期"
            type="number"
            value={formData.parameters.period || 20}
            onChange={(e) => handleParameterChange('period', parseInt(e.target.value))}
            inputProps={{ min: 1, max: 50 }}
            fullWidth
            helperText="布林带的计算周期，通常为20天"
          />
          <TextField
            label="标准差倍数"
            type="number"
            value={formData.parameters.std_dev || 2}
            onChange={(e) => handleParameterChange('std_dev', parseFloat(e.target.value))}
            inputProps={{ min: 0.5, max: 5, step: 0.1 }}
            fullWidth
            helperText="布林带宽度的标准差倍数，通常为2倍"
          />
        </Box>
      );
    }

    return (
      <Alert severity="info">
        请先选择一个策略模板，或手动配置参数
      </Alert>
    );
  };

  // 渲染预览步骤
  const renderPreviewStep = () => (
    <Box sx={{ py: 2 }}>
      <Typography variant="h6" gutterBottom>
        策略预览
      </Typography>
      
      <Paper sx={{ p: 3, mb: 2 }}>
        <Typography variant="subtitle1" gutterBottom fontWeight="bold">
          基本信息
        </Typography>
        <Divider sx={{ mb: 2 }} />
        <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 2 }}>
          <Box>
            <Typography variant="body2" color="text.secondary">策略名称</Typography>
            <Typography variant="body1">{formData.name}</Typography>
          </Box>
          <Box>
            <Typography variant="body2" color="text.secondary">策略类型</Typography>
            <Typography variant="body1">{STRATEGY_TYPES[formData.strategy_type]}</Typography>
          </Box>
          <Box sx={{ gridColumn: '1 / -1' }}>
            <Typography variant="body2" color="text.secondary">策略描述</Typography>
            <Typography variant="body1">{formData.description || '无'}</Typography>
          </Box>
        </Box>
      </Paper>

      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" gutterBottom fontWeight="bold">
          策略参数
        </Typography>
        <Divider sx={{ mb: 2 }} />
        {Object.keys(formData.parameters).length > 0 ? (
          <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 2 }}>
            {Object.entries(formData.parameters).map(([key, value]) => (
              <Box key={key}>
                <Typography variant="body2" color="text.secondary">{key}</Typography>
                <Typography variant="body1">{String(value)}</Typography>
              </Box>
            ))}
          </Box>
        ) : (
          <Typography variant="body2" color="text.secondary">暂无参数配置</Typography>
        )}
      </Paper>

      <Alert severity="info" sx={{ mt: 2 }}>
        创建后，策略将处于"非活跃"状态。您可以在策略列表中激活它，或直接运行回测验证。
      </Alert>
    </Box>
  );

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="md"
      fullWidth
    >
      <DialogTitle>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6">创建新策略</Typography>
          <IconButton onClick={handleClose} size="small">
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>

      <DialogContent dividers>
        <Stepper activeStep={activeStep} sx={{ mb: 3 }}>
          {steps.map((label) => (
            <Step key={label}>
              <StepLabel>{label}</StepLabel>
            </Step>
          ))}
        </Stepper>

        {activeStep === 0 && renderBasicInfoStep()}
        {activeStep === 1 && renderParametersStep()}
        {activeStep === 2 && renderPreviewStep()}
      </DialogContent>

      <DialogActions sx={{ px: 3, py: 2 }}>
        <Button onClick={handleClose} disabled={loading}>
          取消
        </Button>
        {activeStep > 0 && (
          <Button onClick={handleBack} disabled={loading}>
            上一步
          </Button>
        )}
        {activeStep < steps.length - 1 ? (
          <Button variant="contained" onClick={handleNext}>
            下一步
          </Button>
        ) : (
          <Button
            variant="contained"
            onClick={handleSubmit}
            disabled={loading}
          >
            {loading ? <CircularProgress size={24} /> : '创建策略'}
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
};

