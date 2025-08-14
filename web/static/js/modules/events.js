/**
 * 事件处理模块
 * 负责所有用户交互事件的处理
 */

class EventsModule {
    constructor(client, apiService, displayModule, favoritesModule = null) {
        this.client = client;
        this.apiService = apiService;
        this.displayModule = displayModule;
        this.favoritesModule = favoritesModule;
        
        this.init();
    }

    /**
     * 初始化事件模块
     */
    init() {
        this.setupEventListeners();
        this.setDefaultDates();
    }

    /**
     * 设置事件监听器
     */
    setupEventListeners() {
        // 查询按钮事件
        const queryDailyBtn = document.getElementById('queryDaily');
        if (queryDailyBtn) {
            queryDailyBtn.addEventListener('click', () => this.handleDailyQuery());
        }

        const queryIndicatorsBtn = document.getElementById('queryIndicators');
        if (queryIndicatorsBtn) {
            queryIndicatorsBtn.addEventListener('click', () => this.handleIndicatorsQuery());
        }

        const queryPredictionsBtn = document.getElementById('queryPredictions');
        if (queryPredictionsBtn) {
            queryPredictionsBtn.addEventListener('click', () => this.handlePredictionsQuery());
        }

        // 回车键支持
        const stockCodeInput = document.getElementById('stockCode');
        if (stockCodeInput) {
            stockCodeInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    this.handleDailyQuery();
                }
            });

            // 实时验证股票代码格式和更新收藏按钮状态
            stockCodeInput.addEventListener('input', (e) => {
                this.validateStockCode(e);
                // 更新收藏按钮状态
                if (this.favoritesModule) {
                    this.favoritesModule.updateFavoriteButtonState();
                }
            });
        }
    }

    /**
     * 设置默认日期
     */
    setDefaultDates() {
        const today = new Date();
        const sixtyDaysAgo = new Date(today.getTime() - 60 * 24 * 60 * 60 * 1000);
        
        const endDateInput = document.getElementById('endDate');
        const startDateInput = document.getElementById('startDate');
        
        if (endDateInput) {
            endDateInput.value = this.client.formatDate(today);
        }
        if (startDateInput) {
            startDateInput.value = this.client.formatDate(sixtyDaysAgo);
        }
    }

    /**
     * 验证股票代码格式
     */
    validateStockCode(event) {
        const code = event.target.value.toUpperCase();
        const pattern = /^[0-9]{6}\.(SZ|SH)$/;
        
        if (code && !pattern.test(code)) {
            event.target.setCustomValidity('请输入正确的股票代码格式，如：000001.SZ 或 600000.SH');
        } else {
            event.target.setCustomValidity('');
        }
    }

    /**
     * 处理日线数据查询
     */
    async handleDailyQuery() {
        try {
            this.client.setLoading(true);
            this.client.hideAllResultCards();
            
            const stockCode = this.client.getStockCode();
            const { startDate, endDate } = this.client.getDateRange();
            
            // 并行获取股票基本信息和日线数据
            const [stockBasic, dailyData] = await Promise.all([
                this.apiService.getStockBasic(stockCode),
                this.apiService.getDailyData(stockCode, startDate, endDate)
            ]);
            
            this.displayModule.displayDailyData(dailyData, stockCode, stockBasic);
            
        } catch (error) {
            this.client.showError(`获取日线数据失败: ${error.message}`);
        } finally {
            this.client.setLoading(false);
        }
    }

    /**
     * 处理技术指标查询
     */
    async handleIndicatorsQuery() {
        try {
            this.client.setLoading(true);
            this.client.hideAllResultCards();
            
            const stockCode = this.client.getStockCode();
            
            // 并行获取股票基本信息和技术指标
            const [stockBasic, indicatorsData] = await Promise.all([
                this.apiService.getStockBasic(stockCode),
                this.apiService.getIndicators(stockCode)
            ]);
            
            this.displayModule.displayIndicators(indicatorsData, stockCode, stockBasic);
            
        } catch (error) {
            this.client.showError(`获取技术指标失败: ${error.message}`);
        } finally {
            this.client.setLoading(false);
        }
    }

    /**
     * 处理买卖预测查询
     */
    async handlePredictionsQuery() {
        try {
            this.client.setLoading(true);
            this.client.hideAllResultCards();
            
            const stockCode = this.client.getStockCode();
            
            // 并行获取股票基本信息和预测数据
            const [stockBasic, predictionsData] = await Promise.all([
                this.apiService.getStockBasic(stockCode),
                this.apiService.getPredictions(stockCode)
            ]);
            
            this.displayModule.displayPredictions(predictionsData, stockCode, stockBasic);
            
        } catch (error) {
            this.client.showError(`获取买卖预测失败: ${error.message}`);
        } finally {
            this.client.setLoading(false);
        }
    }
}

// 导出事件处理模块类
window.EventsModule = EventsModule;
