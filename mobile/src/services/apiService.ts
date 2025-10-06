/**
 * API服务 - 简化的API调用包装器
 * 为了与现有代码兼容，提供promise风格的API调用
 */

import { appConfig } from '@/constants/config';

interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

class ApiService {
  private baseURL: string;

  constructor() {
    this.baseURL = appConfig.apiBaseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      const url = `${this.baseURL}${endpoint}`;
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          ...options.headers,
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('API请求失败:', error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  // 获取股票基本信息
  async getStockBasic(stockCode: string) {
    return this.request(`/api/v1/stocks/${stockCode}/basic`);
  }

  // 获取股票日线数据
  async getDailyData(
    stockCode: string,
    startDate?: string,
    endDate?: string,
    adjust: string = 'qfq'
  ) {
    const params = new URLSearchParams();
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);
    params.append('adjust', adjust);

    return this.request(`/api/v1/stocks/${stockCode}/daily?${params.toString()}`);
  }

  // 获取技术指标
  async getIndicators(stockCode: string) {
    return this.request(`/api/v1/stocks/${stockCode}/indicators`);
  }

  // 获取预测数据
  async getPredictions(stockCode: string) {
    return this.request(`/api/v1/stocks/${stockCode}/predictions`);
  }

  // 获取基本面数据
  async getFundamentalData(stockCode: string, period?: string) {
    const params = period ? `?period=${period}` : '';
    return this.request(`/api/v1/stocks/${stockCode}/fundamental${params}`);
  }

  // 搜索股票
  async searchStocks(query: string, limit: number = 20) {
    const params = new URLSearchParams({
      q: query,
      limit: limit.toString(),
    });
    return this.request(`/api/v1/stocks/search?${params.toString()}`);
  }

  // 获取股票列表
  async getStockList(page: number = 1, pageSize: number = 50) {
    return this.request(`/api/v1/stocks?page=${page}&page_size=${pageSize}`);
  }

  // 获取收藏列表
  async getFavorites() {
    return this.request('/api/v1/favorites');
  }

  // 添加收藏
  async addFavorite(data: {
    stock_code: string;
    stock_name: string;
    group_id?: string;
  }) {
    return this.request('/api/v1/favorites', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // 删除收藏
  async deleteFavorite(id: number) {
    return this.request(`/api/v1/favorites/${id}`, {
      method: 'DELETE',
    });
  }

  // 获取分组列表
  async getGroups() {
    return this.request('/api/v1/groups');
  }

  // 创建分组
  async createGroup(data: { name: string; description?: string }) {
    return this.request('/api/v1/groups', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // 删除分组
  async deleteGroup(id: number) {
    return this.request(`/api/v1/groups/${id}`, {
      method: 'DELETE',
    });
  }

  // 健康检查
  async healthCheck() {
    return this.request('/api/v1/health');
  }
}

export const apiService = new ApiService();
export default apiService;
