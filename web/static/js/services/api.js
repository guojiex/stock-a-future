/**
 * API服务模块
 * 负责所有与后端API的交互逻辑
 */

class ApiService {
    constructor(client) {
        this.client = client;
    }

    /**
     * 获取股票基本信息
     */
    async getStockBasic(stockCode) {
        const endpoint = `/api/v1/stocks/${stockCode}/basic`;
        
        // 记录请求详情
        console.log(`[API] 获取股票基本信息请求:`, {
            stockCode,
            endpoint,
            timestamp: new Date().toISOString()
        });
        
        try {
            const response = await this.client.makeRequest(endpoint);
            
            // 记录响应详情
            console.log(`[API] 股票基本信息响应:`, {
                success: response.success,
                hasData: !!response.data,
                error: response.error || null,
                timestamp: new Date().toISOString()
            });
            
            if (response.success && response.data) {
                return response.data;
            } else {
                // 详细记录错误信息
                const errorDetails = {
                    message: response.error || '获取股票基本信息失败',
                    stockCode,
                    response: response,
                    timestamp: new Date().toISOString()
                };
                
                console.error(`[API] 获取股票基本信息失败 - 详细信息:`, errorDetails);
                throw new Error(response.error || '获取股票基本信息失败');
            }
        } catch (error) {
            // 详细记录异常信息
            const errorDetails = {
                message: error.message,
                stockCode,
                endpoint,
                errorType: error.constructor.name,
                stack: error.stack,
                timestamp: new Date().toISOString()
            };
            
            console.error(`[API] 获取股票基本信息异常 - 详细信息:`, errorDetails);
            
            // 返回默认信息，使用常见股票的本地映射
            const stockName = this.getStockNameFromLocal(stockCode);
            return {
                ts_code: stockCode,
                name: stockName,
                symbol: stockCode.split('.')[0]
            };
        }
    }

    /**
     * 获取股票日线数据
     */
    async getDailyData(stockCode, startDate, endDate) {
        const endpoint = `/api/v1/stocks/${stockCode}/daily?start_date=${startDate}&end_date=${endDate}`;
        
        // 记录请求详情
        console.log(`[API] 获取日线数据请求:`, {
            stockCode,
            startDate,
            endDate,
            endpoint,
            timestamp: new Date().toISOString()
        });
        
        try {
            const response = await this.client.makeRequest(endpoint);
            
            // 记录响应详情
            console.log(`[API] 日线数据响应:`, {
                success: response.success,
                hasData: !!response.data,
                dataLength: response.data ? response.data.length : 0,
                error: response.error || null,
                timestamp: new Date().toISOString()
            });
            
            if (response.success && response.data) {
                return response.data;
            } else {
                // 详细记录错误信息
                const errorDetails = {
                    message: response.error || '获取日线数据失败',
                    stockCode,
                    startDate,
                    endDate,
                    response: response,
                    timestamp: new Date().toISOString()
                };
                
                console.error(`[API] 获取日线数据失败 - 详细信息:`, errorDetails);
                throw new Error(response.error || '获取日线数据失败');
            }
        } catch (error) {
            // 详细记录异常信息
            const errorDetails = {
                message: error.message,
                stockCode,
                startDate,
                endDate,
                endpoint,
                errorType: error.constructor.name,
                stack: error.stack,
                timestamp: new Date().toISOString()
            };
            
            console.error(`[API] 获取日线数据异常 - 详细信息:`, errorDetails);
            throw error;
        }
    }

    /**
     * 获取技术指标
     */
    async getIndicators(stockCode) {
        const endpoint = `/api/v1/stocks/${stockCode}/indicators`;
        
        // 记录请求详情
        console.log(`[API] 获取技术指标请求:`, {
            stockCode,
            endpoint,
            timestamp: new Date().toISOString()
        });
        
        try {
            const response = await this.client.makeRequest(endpoint);
            
            // 记录响应详情
            console.log(`[API] 技术指标响应:`, {
                success: response.success,
                hasData: !!response.data,
                error: response.error || null,
                timestamp: new Date().toISOString()
            });
            
            if (response.success && response.data) {
                return response.data;
            } else {
                // 详细记录错误信息
                const errorDetails = {
                    message: response.error || '获取技术指标失败',
                    stockCode,
                    response: response,
                    timestamp: new Date().toISOString()
                };
                
                console.error(`[API] 获取技术指标失败 - 详细信息:`, errorDetails);
                throw new Error(response.error || '获取技术指标失败');
            }
        } catch (error) {
            // 详细记录异常信息
            const errorDetails = {
                message: error.message,
                stockCode,
                endpoint,
                errorType: error.constructor.name,
                stack: error.stack,
                timestamp: new Date().toISOString()
            };
            
            console.error(`[API] 获取技术指标异常 - 详细信息:`, errorDetails);
            throw error;
        }
    }

    /**
     * 获取买卖预测
     */
    async getPredictions(stockCode) {
        const endpoint = `/api/v1/stocks/${stockCode}/predictions`;
        
        // 记录请求详情
        console.log(`[API] 获取买卖预测请求:`, {
            stockCode,
            endpoint,
            timestamp: new Date().toISOString()
        });
        
        try {
            const response = await this.client.makeRequest(endpoint);
            
            // 记录响应详情
            console.log(`[API] 买卖预测响应:`, {
                success: response.success,
                hasData: !!response.data,
                error: response.error || null,
                timestamp: new Date().toISOString()
            });
            
            if (response.success && response.data) {
                return response.data;
            } else {
                // 详细记录错误信息
                const errorDetails = {
                    message: response.error || '获取预测数据失败',
                    stockCode,
                    response: response,
                    timestamp: new Date().toISOString()
                };
                
                console.error(`[API] 获取买卖预测失败 - 详细信息:`, errorDetails);
                throw new Error(response.error || '获取预测数据失败');
            }
        } catch (error) {
            // 详细记录异常信息
            const errorDetails = {
                message: error.message,
                stockCode,
                endpoint,
                errorType: error.constructor.name,
                stack: error.stack,
                timestamp: new Date().toISOString()
            };
            
            console.error(`[API] 获取买卖预测异常 - 详细信息:`, errorDetails);
            throw error;
        }
    }

    /**
     * 获取收藏股票信号汇总
     */
    async getFavoritesSignals() {
        const endpoint = `/api/v1/favorites/signals`;
        
        // 记录请求详情
        console.log(`[API] 获取收藏股票信号汇总请求:`, {
            endpoint,
            timestamp: new Date().toISOString()
        });
        
        try {
            const response = await this.client.makeRequest(endpoint);
            
            // 记录响应详情
            console.log(`[API] 收藏股票信号汇总响应:`, {
                success: response.success,
                hasData: !!response.data,
                hasSignals: !!(response.data && response.data.signals),
                signalsCount: response.data && response.data.signals ? response.data.signals.length : 0,
                error: response.error || null,
                timestamp: new Date().toISOString()
            });
            
            if (response.success && response.data) {
                return response.data;
            } else {
                // 详细记录错误信息
                const errorDetails = {
                    message: response.error || '获取信号汇总失败',
                    response: response,
                    timestamp: new Date().toISOString()
                };
                
                console.error(`[API] 获取信号汇总失败 - 详细信息:`, errorDetails);
                throw new Error(response.error || '获取信号汇总失败');
            }
        } catch (error) {
            // 详细记录异常信息
            const errorDetails = {
                message: error.message,
                endpoint,
                errorType: error.constructor.name,
                stack: error.stack,
                timestamp: new Date().toISOString()
            };
            
            console.error(`[API] 获取信号汇总异常 - 详细信息:`, errorDetails);
            throw error;
        }
    }

    /**
     * 搜索股票
     */
    async searchStocks(keyword, limit = 10) {
        const endpoint = `/api/v1/stocks/search?q=${encodeURIComponent(keyword)}&limit=${limit}`;
        
        // 记录请求详情
        console.log(`[API] 搜索股票请求:`, {
            keyword,
            limit,
            endpoint,
            timestamp: new Date().toISOString()
        });
        
        try {
            const response = await this.client.makeRequest(endpoint);
            
            // 记录响应详情
            console.log(`[API] 搜索股票响应:`, {
                success: response.success,
                hasData: !!response.data,
                hasStocks: !!(response.data && response.data.stocks),
                stocksCount: response.data && response.data.stocks ? response.data.stocks.length : 0,
                error: response.error || null,
                timestamp: new Date().toISOString()
            });
            
            if (response.success && response.data.stocks) {
                return response.data.stocks;
            } else {
                // 详细记录错误信息
                const errorDetails = {
                    message: response.error || '搜索股票失败',
                    keyword,
                    limit,
                    response: response,
                    timestamp: new Date().toISOString()
                };
                
                console.error(`[API] 搜索股票失败 - 详细信息:`, errorDetails);
                throw new Error(response.error || '搜索股票失败');
            }
        } catch (error) {
            // 详细记录异常信息
            const errorDetails = {
                message: error.message,
                keyword,
                limit,
                endpoint,
                errorType: error.constructor.name,
                stack: error.stack,
                timestamp: new Date().toISOString()
            };
            
            console.error(`[API] 搜索股票异常 - 详细信息:`, errorDetails);
            throw error;
        }
    }

    /**
     * 从本地映射获取股票名称（备选方案）
     */
    getStockNameFromLocal(stockCode) {
        // 常见股票代码到名称的映射（更全面的列表）
        const stockMap = {
            // 深圳主板
            '000001.SZ': '平安银行',
            '000002.SZ': '万科A',
            '000858.SZ': '五粮液',
            '000876.SZ': '新希望',
            '000895.SZ': '双汇发展',
            '000938.SZ': '紫光股份',
            
            // 深圳中小板
            '002415.SZ': '海康威视',
            '002594.SZ': '比亚迪',
            '002714.SZ': '牧原股份',
            '002304.SZ': '洋河股份',
            
            // 深圳创业板
            '300059.SZ': '东方财富',
            '300750.SZ': '宁德时代',
            '300015.SZ': '爱尔眼科',
            '300142.SZ': '沃森生物',
            
            // 上海主板
            '600000.SH': '浦发银行',
            '600036.SH': '招商银行',
            '600519.SH': '贵州茅台',
            '600887.SH': '伊利股份',
            '600276.SH': '恒瑞医药',
            '600031.SH': '三一重工',
            '600703.SH': '三安光电',
            '601318.SH': '中国平安',
            '601166.SH': '兴业银行',
            '601012.SH': '隆基绿能',
            '601888.SH': '中国中免',
            '600009.SH': '上海机场',
            '600104.SH': '上汽集团',
            '600196.SH': '复星医药',
            '600309.SH': '万华化学',
            '600436.SH': '片仔癀',
            '600690.SH': '海尔智家',
            '600745.SZ': '闻泰科技',
            '600809.SH': '山西汾酒',
            '600893.SH': '航发动力',
            
            // 科创板
            '688111.SH': '金山办公',
            '688981.SH': '中芯国际',
            '688036.SH': '传音控股',
            '688599.SH': '天合光能'
        };
        
        return stockMap[stockCode] || `股票${stockCode.split('.')[0]}`; // 如果找不到映射，则返回格式化的代码
    }
}

// 导出API服务类
window.ApiService = ApiService;
