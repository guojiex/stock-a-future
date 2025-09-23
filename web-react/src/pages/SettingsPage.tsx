/**
 * 设置页面 - Web版本
 */

import React from 'react';
import {
  Container,
  Typography,
  Card,
  CardContent,
  Box,
  Alert,
} from '@mui/material';
import { Settings as SettingsIcon } from '@mui/icons-material';

const SettingsPage: React.FC = () => {
  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Typography variant="h4" gutterBottom>
        设置
      </Typography>
      
      <Card>
        <CardContent>
          <Box textAlign="center" py={8}>
            <SettingsIcon sx={{ fontSize: 64, color: 'primary.main', mb: 2 }} />
            <Typography variant="h5" gutterBottom>
              应用设置
            </Typography>
            <Alert severity="info" sx={{ mt: 3, maxWidth: 400, mx: 'auto' }}>
              设置功能开发中...
              <br />
              这里将显示：
              <br />• 主题设置
              <br />• 数据刷新间隔
              <br />• 通知设置
              <br />• 关于应用
            </Alert>
          </Box>
        </CardContent>
      </Card>
    </Container>
  );
};

export default SettingsPage;
