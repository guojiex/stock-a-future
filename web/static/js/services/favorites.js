/**
 * 收藏股票服务
 * 负责与后端收藏API的交互
 */

class FavoritesService {
    constructor(client) {
        this.client = client;
    }

    /**
     * 获取收藏列表
     */
    async getFavorites() {
        try {
            const response = await fetch(`${this.client.baseURL}/api/v1/favorites`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            if (!response.ok) {
                throw new Error(`获取收藏列表失败: ${response.status}`);
            }

            const data = await response.json();
            
            if (!data.success) {
                throw new Error(data.error || '获取收藏列表失败');
            }

            return data.data.favorites || [];
        } catch (error) {
            console.error('获取收藏列表失败:', error);
            throw error;
        }
    }

    /**
     * 添加收藏
     */
    async addFavorite(stockCode, stockName, startDate, endDate, groupId) {
        if (!groupId) {
            throw new Error('分组ID是必需的参数');
        }
        try {
            const requestData = {
                ts_code: stockCode,
                name: stockName,
                start_date: startDate,
                end_date: endDate,
                group_id: groupId
            };

            const response = await fetch(`${this.client.baseURL}/api/v1/favorites`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            });

            const data = await response.json();
            
            if (!response.ok || !data.success) {
                throw new Error(data.error || `添加收藏失败: ${response.status}`);
            }

            return data.data;
        } catch (error) {
            console.error('添加收藏失败:', error);
            throw error;
        }
    }



    /**
     * 删除收藏
     */
    async deleteFavorite(favoriteId) {
        try {
            const response = await fetch(`${this.client.baseURL}/api/v1/favorites/${favoriteId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            const data = await response.json();
            
            if (!response.ok || !data.success) {
                throw new Error(data.error || `删除收藏失败: ${response.status}`);
            }

            return data.data;
        } catch (error) {
            console.error('删除收藏失败:', error);
            throw error;
        }
    }

    /**
     * 更新收藏的时间范围
     */
    async updateFavorite(favoriteId, startDate, endDate) {
        try {
            const requestData = {
                start_date: startDate,
                end_date: endDate
            };

            const response = await fetch(`${this.client.baseURL}/api/v1/favorites/${favoriteId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            });

            const data = await response.json();
            
            if (!response.ok || !data.success) {
                throw new Error(data.error || `更新收藏失败: ${response.status}`);
            }

            return data.data;
        } catch (error) {
            console.error('更新收藏失败:', error);
            throw error;
        }
    }

    /**
     * 检查股票是否已收藏到指定分组
     */
    async checkFavoriteInGroup(stockCode, groupId = null) {
        try {
            if (groupId === null) {
                // 如果没有指定分组，检查是否在任何分组中收藏
                return this.checkFavorite(stockCode);
            }
            
            // 使用本地数据检查特定分组中是否已收藏
            const favorite = await this.findFavoriteByCodeAndGroup(stockCode, groupId);
            return favorite !== null;
        } catch (error) {
            console.error('检查分组收藏状态失败:', error);
            return false; // 出错时默认返回未收藏状态
        }
    }

    /**
     * 检查股票是否已收藏（任何分组）- 向后兼容
     */
    async checkFavorite(stockCode) {
        return this.checkFavoriteInGroup(stockCode, null);
    }

    /**
     * 通过股票代码查找收藏记录
     */
    async findFavoriteByCode(stockCode) {
        try {
            const favorites = await this.getFavorites();
            return favorites.find(favorite => favorite.ts_code === stockCode) || null;
        } catch (error) {
            console.error('查找收藏记录失败:', error);
            return null;
        }
    }

    /**
     * 通过股票代码和分组ID查找收藏记录
     */
    async findFavoriteByCodeAndGroup(stockCode, groupId) {
        try {
            const favorites = await this.getFavorites();
            
            const result = favorites.find(favorite => 
                favorite.ts_code === stockCode && 
                (favorite.group_id || 'default') === groupId
            ) || null;
            
            return result;
        } catch (error) {
            console.error('查找分组收藏记录失败:', error);
            return null;
        }
    }

    // === 分组相关API ===

    /**
     * 获取分组列表
     */
    async getGroups() {
        try {
            const response = await fetch(`${this.client.baseURL}/api/v1/groups`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            if (!response.ok) {
                throw new Error(`获取分组列表失败: ${response.status}`);
            }

            const data = await response.json();
            
            if (!data.success) {
                throw new Error(data.error || '获取分组列表失败');
            }

            return data.data.groups || [];
        } catch (error) {
            console.error('获取分组列表失败:', error);
            throw error;
        }
    }

    /**
     * 创建分组
     */
    async createGroup(name, color = '#3b82f6') {
        try {
            const requestData = {
                name: name,
                color: color
            };

            const response = await fetch(`${this.client.baseURL}/api/v1/groups`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            });

            const data = await response.json();
            
            if (!response.ok || !data.success) {
                throw new Error(data.error || `创建分组失败: ${response.status}`);
            }

            return data.data;
        } catch (error) {
            console.error('创建分组失败:', error);
            throw error;
        }
    }

    /**
     * 更新分组
     */
    async updateGroup(groupId, name, color, sortOrder) {
        try {
            const requestData = {};
            if (name) requestData.name = name;
            if (color) requestData.color = color;
            if (sortOrder) requestData.sort_order = sortOrder;

            const response = await fetch(`${this.client.baseURL}/api/v1/groups/${groupId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            });

            const data = await response.json();
            
            if (!response.ok || !data.success) {
                throw new Error(data.error || `更新分组失败: ${response.status}`);
            }

            return data.data;
        } catch (error) {
            console.error('更新分组失败:', error);
            throw error;
        }
    }

    /**
     * 删除分组
     */
    async deleteGroup(groupId) {
        try {
            const response = await fetch(`${this.client.baseURL}/api/v1/groups/${groupId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            const data = await response.json();
            
            if (!response.ok || !data.success) {
                throw new Error(data.error || `删除分组失败: ${response.status}`);
            }

            return data.data;
        } catch (error) {
            console.error('删除分组失败:', error);
            throw error;
        }
    }

    /**
     * 更新收藏排序
     */
    async updateFavoritesOrder(favoriteOrders) {
        try {
            const requestData = {
                favorite_orders: favoriteOrders
            };

            const response = await fetch(`${this.client.baseURL}/api/v1/favorites/order`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            });

            const data = await response.json();
            
            if (!response.ok || !data.success) {
                throw new Error(data.error || `更新收藏排序失败: ${response.status}`);
            }

            return data.data;
        } catch (error) {
            console.error('更新收藏排序失败:', error);
            throw error;
        }
    }


}

// 导出收藏服务类
window.FavoritesService = FavoritesService;
