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
    async addFavorite(stockCode, stockName, startDate, endDate) {
        try {
            const requestData = {
                ts_code: stockCode,
                name: stockName,
                start_date: startDate,
                end_date: endDate
            };

            const response = await fetch(`${this.client.baseURL}/api/v1/favorites`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            });

            if (!response.ok) {
                throw new Error(`添加收藏失败: ${response.status}`);
            }

            const data = await response.json();
            
            if (!data.success) {
                throw new Error(data.error || '添加收藏失败');
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

            if (!response.ok) {
                throw new Error(`删除收藏失败: ${response.status}`);
            }

            const data = await response.json();
            
            if (!data.success) {
                throw new Error(data.error || '删除收藏失败');
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

            if (!response.ok) {
                throw new Error(`更新收藏失败: ${response.status}`);
            }

            const data = await response.json();
            
            if (!data.success) {
                throw new Error(data.error || '更新收藏失败');
            }

            return data.data;
        } catch (error) {
            console.error('更新收藏失败:', error);
            throw error;
        }
    }

    /**
     * 检查股票是否已收藏
     */
    async checkFavorite(stockCode) {
        try {
            const response = await fetch(`${this.client.baseURL}/api/v1/favorites/check/${stockCode}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            if (!response.ok) {
                throw new Error(`检查收藏状态失败: ${response.status}`);
            }

            const data = await response.json();
            
            if (!data.success) {
                throw new Error(data.error || '检查收藏状态失败');
            }

            return data.data.is_favorite || false;
        } catch (error) {
            console.error('检查收藏状态失败:', error);
            return false; // 出错时默认返回未收藏状态
        }
    }
}

// 导出收藏服务类
window.FavoritesService = FavoritesService;
