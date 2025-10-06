/**
 * 收藏页面 - Web版本
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
import { Star as StarIcon } from '@mui/icons-material';

const FavoritesPage: React.FC = () => {
  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" gutterBottom>
        收藏
      </Typography>
      
      <Card>
        <CardContent>
          <Box textAlign="center" py={8}>
            <StarIcon sx={{ fontSize: 64, color: 'primary.main', mb: 2 }} />
            <Typography variant="h5" gutterBottom>
              我的收藏
            </Typography>
            <Alert severity="info" sx={{ mt: 3, maxWidth: 400, mx: 'auto' }}>
              收藏功能开发中...
              <br />
              这里将显示：
              <br />• 收藏的股票列表
              <br />• 分组管理
              <br />• 快速访问
              <br />• 价格提醒
            </Alert>
          </Box>
        </CardContent>
      </Card>
    </Container>
  );
};

export default FavoritesPage;
