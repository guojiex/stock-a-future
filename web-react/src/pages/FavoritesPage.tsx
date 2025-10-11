/**
 * 收藏页面 - Web版本
 * 功能：
 * - 收藏的股票列表展示
 * - 分组管理（创建、编辑、删除分组）
 * - 快速访问股票详情
 * - 股票移动到不同分组
 */

import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Typography,
  Card,
  CardContent,
  Box,
  Grid,
  Button,
  IconButton,
  Chip,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  ListItemSecondaryAction,
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
  Divider,
  Paper,
  Tabs,
  Tab,
  Menu,
  ListItemIcon,
} from '@mui/material';
import {
  Star as StarIcon,
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  TrendingUp as TrendingUpIcon,
  Folder as FolderIcon,
  FolderOpen as FolderOpenIcon,
  MoreVert as MoreVertIcon,
} from '@mui/icons-material';

import {
  useGetFavoritesQuery,
  useGetGroupsQuery,
  useDeleteFavoriteMutation,
  useUpdateFavoriteMutation,
  useCreateGroupMutation,
  useUpdateGroupMutation,
  useDeleteGroupMutation,
} from '../services/api';
import { Favorite, FavoriteGroup } from '../types/stock';

const FavoritesPage: React.FC = () => {
  const navigate = useNavigate();

  // API查询
  const { data: favoritesData, isLoading: favoritesLoading } = useGetFavoritesQuery();
  const { data: groupsData, isLoading: groupsLoading } = useGetGroupsQuery();

  // API变更
  const [deleteFavorite] = useDeleteFavoriteMutation();
  const [updateFavorite] = useUpdateFavoriteMutation();
  const [createGroup] = useCreateGroupMutation();
  const [updateGroup] = useUpdateGroupMutation();
  const [deleteGroup] = useDeleteGroupMutation();

  // 本地状态
  const [selectedGroup, setSelectedGroup] = useState<string>('all');
  const [groupDialogOpen, setGroupDialogOpen] = useState(false);
  const [editingGroup, setEditingGroup] = useState<FavoriteGroup | null>(null);
  const [groupName, setGroupName] = useState('');
  const [groupColor, setGroupColor] = useState('#1976d2');
  const [moveDialogOpen, setMoveDialogOpen] = useState(false);
  const [movingFavorite, setMovingFavorite] = useState<Favorite | null>(null);
  const [targetGroupId, setTargetGroupId] = useState<string>('');
  const [groupMenuAnchor, setGroupMenuAnchor] = useState<null | HTMLElement>(null);
  const [menuGroupId, setMenuGroupId] = useState<string>('');

  // 获取数据
  const favorites = favoritesData?.data?.favorites || [];
  const groups = groupsData?.data?.groups || [];

  // 按分组过滤收藏
  const filteredFavorites = useMemo(() => {
    if (selectedGroup === 'all') {
      return favorites;
    }
    return favorites.filter(fav => fav.group_id === selectedGroup);
  }, [favorites, selectedGroup]);

  // 按分组统计数量
  const groupCounts = useMemo(() => {
    const counts: Record<string, number> = {
      all: favorites.length,
    };
    groups.forEach(group => {
      counts[group.id] = favorites.filter(fav => fav.group_id === group.id).length;
    });
    return counts;
  }, [favorites, groups]);

  // 处理点击股票
  const handleStockClick = (favorite: Favorite) => {
    // 导航到股票详情页
    navigate(`/stock/${favorite.ts_code}`);
  };

  // 处理删除收藏
  const handleDeleteFavorite = async (favoriteId: string, event: React.MouseEvent) => {
    event.stopPropagation();
    if (window.confirm('确定要删除这个收藏吗？')) {
      try {
        await deleteFavorite(favoriteId).unwrap();
      } catch (error) {
        console.error('删除收藏失败:', error);
        alert('删除失败，请重试');
      }
    }
  };

  // 处理移动收藏到分组
  const handleMoveToGroup = (favorite: Favorite, event: React.MouseEvent) => {
    event.stopPropagation();
    setMovingFavorite(favorite);
    setTargetGroupId(favorite.group_id || '');
    setMoveDialogOpen(true);
  };

  // 确认移动收藏
  const handleConfirmMove = async () => {
    if (!movingFavorite) return;

    try {
      await updateFavorite({
        id: movingFavorite.id,
        data: {
          group_id: targetGroupId || undefined,
        },
      }).unwrap();
      setMoveDialogOpen(false);
      setMovingFavorite(null);
    } catch (error) {
      console.error('移动收藏失败:', error);
      alert('移动失败，请重试');
    }
  };

  // 打开创建/编辑分组对话框
  const handleOpenGroupDialog = (group?: FavoriteGroup) => {
    if (group) {
      setEditingGroup(group);
      setGroupName(group.name);
      setGroupColor(group.color || '#1976d2');
    } else {
      setEditingGroup(null);
      setGroupName('');
      setGroupColor('#1976d2');
    }
    setGroupDialogOpen(true);
  };

  // 保存分组
  const handleSaveGroup = async () => {
    if (!groupName.trim()) {
      alert('请输入分组名称');
      return;
    }

    try {
      if (editingGroup) {
        // 更新分组
        await updateGroup({
          id: editingGroup.id,
          data: {
            name: groupName,
            color: groupColor,
          },
        }).unwrap();
      } else {
        // 创建分组
        await createGroup({
          name: groupName,
          color: groupColor,
        }).unwrap();
      }
      setGroupDialogOpen(false);
    } catch (error) {
      console.error('保存分组失败:', error);
      alert('保存失败，请重试');
    }
  };

  // 删除分组
  const handleDeleteGroup = async (groupId: string, event: React.MouseEvent) => {
    event.stopPropagation();
    if (window.confirm('确定要删除这个分组吗？分组内的股票将移至未分组。')) {
      try {
        await deleteGroup(groupId).unwrap();
        if (selectedGroup === groupId) {
          setSelectedGroup('all');
        }
      } catch (error) {
        console.error('删除分组失败:', error);
        alert('删除失败，请重试');
      }
    }
  };

  // 打开分组菜单
  const handleOpenGroupMenu = (event: React.MouseEvent<HTMLElement>, groupId: string) => {
    event.stopPropagation();
    setGroupMenuAnchor(event.currentTarget);
    setMenuGroupId(groupId);
  };

  // 关闭分组菜单
  const handleCloseGroupMenu = () => {
    setGroupMenuAnchor(null);
    setMenuGroupId('');
  };

  // 从菜单中编辑分组
  const handleEditFromMenu = () => {
    const group = groups.find(g => g.id === menuGroupId);
    if (group) {
      handleOpenGroupDialog(group);
    }
    handleCloseGroupMenu();
  };

  // 从菜单中删除分组
  const handleDeleteFromMenu = () => {
    handleDeleteGroup(menuGroupId, {} as React.MouseEvent);
    handleCloseGroupMenu();
  };

  // 渲染分组标签
  const renderGroupTabs = () => (
    <Paper sx={{ mb: 3 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', px: 2, py: 1, gap: 1 }}>
        <Tabs
          value={selectedGroup}
          onChange={(_, value) => setSelectedGroup(value)}
          variant="scrollable"
          scrollButtons="auto"
          sx={{ flex: 1 }}
        >
          <Tab
            label={`全部 (${groupCounts.all || 0})`}
            value="all"
            icon={<FolderOpenIcon />}
            iconPosition="start"
          />
          {groups.map(group => (
            <Tab
              key={group.id}
              label={
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                  <Box
                    sx={{
                      width: 12,
                      height: 12,
                      borderRadius: '50%',
                      bgcolor: group.color || '#1976d2',
                    }}
                  />
                  {`${group.name} (${groupCounts[group.id] || 0})`}
                  <Box
                    component="span"
                    onClick={(e) => handleOpenGroupMenu(e, group.id)}
                    sx={{
                      ml: 0.5,
                      display: 'inline-flex',
                      alignItems: 'center',
                      opacity: 0.6,
                      '&:hover': { opacity: 1 },
                    }}
                  >
                    <MoreVertIcon fontSize="small" />
                  </Box>
                </Box>
              }
              value={group.id}
              icon={<FolderIcon />}
              iconPosition="start"
            />
          ))}
        </Tabs>
        
        <Button
          startIcon={<AddIcon />}
          onClick={() => handleOpenGroupDialog()}
          size="small"
          variant="outlined"
        >
          新建分组
        </Button>
      </Box>

      {/* 分组操作菜单 */}
      <Menu
        anchorEl={groupMenuAnchor}
        open={Boolean(groupMenuAnchor)}
        onClose={handleCloseGroupMenu}
      >
        <MenuItem onClick={handleEditFromMenu}>
          <ListItemIcon>
            <EditIcon fontSize="small" />
          </ListItemIcon>
          <Typography variant="body2">编辑分组</Typography>
        </MenuItem>
        <MenuItem onClick={handleDeleteFromMenu}>
          <ListItemIcon>
            <DeleteIcon fontSize="small" color="error" />
          </ListItemIcon>
          <Typography variant="body2" color="error">
            删除分组
          </Typography>
        </MenuItem>
      </Menu>
    </Paper>
  );

  // 渲染收藏列表
  const renderFavoritesList = () => {
    if (favoritesLoading) {
      return (
        <Box display="flex" justifyContent="center" py={8}>
          <CircularProgress />
        </Box>
      );
    }

    if (filteredFavorites.length === 0) {
      return (
        <Card>
          <CardContent>
            <Box textAlign="center" py={8}>
              <StarIcon sx={{ fontSize: 64, color: 'text.disabled', mb: 2 }} />
              <Typography variant="h6" color="text.secondary" gutterBottom>
                {selectedGroup === 'all' ? '还没有收藏任何股票' : '该分组暂无收藏'}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {selectedGroup === 'all' 
                  ? '在股票详情页点击星标按钮添加收藏' 
                  : '点击股票右侧的文件夹图标可将股票移动到该分组'}
              </Typography>
            </Box>
          </CardContent>
        </Card>
      );
    }

    return (
      <Card>
        <List>
          {filteredFavorites.map((favorite, index) => (
            <React.Fragment key={favorite.id}>
              <ListItem disablePadding>
                <ListItemButton onClick={() => handleStockClick(favorite)}>
                  <Box
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    width={48}
                    height={48}
                    borderRadius="50%"
                    bgcolor="primary.light"
                    color="primary.contrastText"
                    mr={2}
                  >
                    <TrendingUpIcon />
                  </Box>
                  <ListItemText
                    primary={
                      <Box display="flex" alignItems="center" gap={1}>
                        <Typography variant="subtitle1" fontWeight="medium">
                          {favorite.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {favorite.ts_code}
                        </Typography>
                        {favorite.group_id && (
                          <Chip
                            size="small"
                            label={groups.find(g => g.id === favorite.group_id)?.name || '未知分组'}
                            sx={{
                              bgcolor: groups.find(g => g.id === favorite.group_id)?.color || '#1976d2',
                              color: 'white',
                            }}
                          />
                        )}
                      </Box>
                    }
                    secondary={
                      <Typography variant="body2" color="text.secondary">
                        收藏于 {new Date(favorite.created_at).toLocaleDateString()}
                      </Typography>
                    }
                  />
                  <ListItemSecondaryAction>
                    <IconButton
                      edge="end"
                      onClick={(e) => handleMoveToGroup(favorite, e)}
                      sx={{ mr: 1 }}
                    >
                      <FolderIcon />
                    </IconButton>
                    <IconButton
                      edge="end"
                      onClick={(e) => handleDeleteFavorite(favorite.id, e)}
                      color="error"
                    >
                      <DeleteIcon />
                    </IconButton>
                  </ListItemSecondaryAction>
                </ListItemButton>
              </ListItem>
              {index < filteredFavorites.length - 1 && <Divider />}
            </React.Fragment>
          ))}
        </List>
      </Card>
    );
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">
          我的收藏
        </Typography>
      </Box>

      {/* 分组标签 */}
      {renderGroupTabs()}

      {/* 收藏列表 */}
      {renderFavoritesList()}

      {/* 创建/编辑分组对话框 */}
      <Dialog open={groupDialogOpen} onClose={() => setGroupDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>{editingGroup ? '编辑分组' : '创建分组'}</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="分组名称"
            type="text"
            fullWidth
            value={groupName}
            onChange={(e) => setGroupName(e.target.value)}
            sx={{ mb: 2, mt: 1 }}
          />
          <FormControl fullWidth>
            <InputLabel>分组颜色</InputLabel>
            <Select
              value={groupColor}
              label="分组颜色"
              onChange={(e) => setGroupColor(e.target.value)}
            >
              <MenuItem value="#1976d2">
                <Box display="flex" alignItems="center" gap={1}>
                  <Box width={20} height={20} borderRadius="50%" bgcolor="#1976d2" />
                  蓝色
                </Box>
              </MenuItem>
              <MenuItem value="#2e7d32">
                <Box display="flex" alignItems="center" gap={1}>
                  <Box width={20} height={20} borderRadius="50%" bgcolor="#2e7d32" />
                  绿色
                </Box>
              </MenuItem>
              <MenuItem value="#ed6c02">
                <Box display="flex" alignItems="center" gap={1}>
                  <Box width={20} height={20} borderRadius="50%" bgcolor="#ed6c02" />
                  橙色
                </Box>
              </MenuItem>
              <MenuItem value="#d32f2f">
                <Box display="flex" alignItems="center" gap={1}>
                  <Box width={20} height={20} borderRadius="50%" bgcolor="#d32f2f" />
                  红色
                </Box>
              </MenuItem>
              <MenuItem value="#9c27b0">
                <Box display="flex" alignItems="center" gap={1}>
                  <Box width={20} height={20} borderRadius="50%" bgcolor="#9c27b0" />
                  紫色
                </Box>
              </MenuItem>
              <MenuItem value="#0288d1">
                <Box display="flex" alignItems="center" gap={1}>
                  <Box width={20} height={20} borderRadius="50%" bgcolor="#0288d1" />
                  青色
                </Box>
              </MenuItem>
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setGroupDialogOpen(false)}>取消</Button>
          <Button onClick={handleSaveGroup} variant="contained">保存</Button>
        </DialogActions>
      </Dialog>

      {/* 移动到分组对话框 */}
      <Dialog open={moveDialogOpen} onClose={() => setMoveDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>移动到分组</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" gutterBottom>
            将 {movingFavorite?.name} 移动到：
          </Typography>
          <FormControl fullWidth sx={{ mt: 2 }}>
            <InputLabel>目标分组</InputLabel>
            <Select
              value={targetGroupId}
              label="目标分组"
              onChange={(e) => setTargetGroupId(e.target.value)}
            >
              <MenuItem value="">未分组</MenuItem>
              {groups.map(group => (
                <MenuItem key={group.id} value={group.id}>
                  <Box display="flex" alignItems="center" gap={1}>
                    <Box
                      width={12}
                      height={12}
                      borderRadius="50%"
                      bgcolor={group.color || '#1976d2'}
                    />
                    {group.name}
                  </Box>
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setMoveDialogOpen(false)}>取消</Button>
          <Button onClick={handleConfirmMove} variant="contained">移动</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default FavoritesPage;
