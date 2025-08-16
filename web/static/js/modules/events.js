/**
 * 事件处理模块
 * 负责所有用户交互事件的处理
 */

class EventsModule {
    constructor(client, apiService, displayModule, favoritesModule = null, dateShortcutsModule = null) {
        this.client = client;
        this.apiService = apiService;
        this.displayModule = displayModule;
        this.favoritesModule = favoritesModule;
        this.dateShortcutsModule = dateShortcutsModule;
        
        // 添加数据缓存，避免重复加载
        this.dataCache = new Map();
        this.lastStockCode = null;
        
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

        // Tab切换事件
        this.setupTabNavigation();

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
     * 设置Tab导航事件
     */
    setupTabNavigation() {
        const tabButtons = document.querySelectorAll('.tab-btn');
        tabButtons.forEach(btn => {
            btn.addEventListener('click', () => {
                this.switchTab(btn.dataset.tab);
            });
        });
    }

    /**
     * 切换Tab
     */
    switchTab(tabName, loadData = true) {
        console.log(`[Events] 切换到tab: ${tabName}, 是否加载数据: ${loadData}`);
        
        // 移除所有tab按钮的active状态
        const tabButtons = document.querySelectorAll('.tab-btn');
        tabButtons.forEach(btn => {
            btn.classList.remove('active');
        });

        // 隐藏所有tab内容
        const tabPanes = document.querySelectorAll('.tab-pane');
        tabPanes.forEach(pane => {
            pane.classList.remove('active');
        });

        // 激活选中的tab按钮
        const activeTabBtn = document.querySelector(`[data-tab="${tabName}"]`);
        if (activeTabBtn) {
            activeTabBtn.classList.add('active');
        }

        // 显示选中的tab内容
        const activeTabPane = document.getElementById(`${tabName}-tab`);
        if (activeTabPane) {
            activeTabPane.classList.add('active');
        }
        
        // 根据tab类型获取对应数据（只有在需要时才加载）
        if (loadData) {
            this.loadDataForTab(tabName);
        }
        
        console.log(`[Events] Tab切换完成: ${tabName}`);
    }

    /**
     * 根据tab类型加载对应数据
     */
    async loadDataForTab(tabName) {
        try {
            const stockCode = this.client.getStockCode();
            console.log(`[Events] 为tab ${tabName} 加载数据，股票代码: ${stockCode}`);
            
            // 显示加载状态
            this.showTabLoadingState(tabName, true);
            
            switch (tabName) {
                case 'daily-data':
                    await this.handleDailyQuery();
                    break;
                case 'indicators':
                    await this.handleIndicatorsQuery();
                    break;
                case 'predictions':
                    await this.handlePredictionsQuery();
                    break;
                default:
                    console.warn(`[Events] 未知的tab类型: ${tabName}`);
            }
        } catch (error) {
            console.error(`[Events] 为tab ${tabName} 加载数据失败:`, error);
            // 显示错误信息但不阻止tab切换
            this.client.showError(`加载${this.getTabDisplayName(tabName)}数据失败: ${error.message}`);
        } finally {
            // 隐藏加载状态
            this.showTabLoadingState(tabName, false);
        }
    }

    /**
     * 显示/隐藏tab的加载状态
     */
    showTabLoadingState(tabName, isLoading) {
        const tabPane = document.getElementById(`${tabName}-tab`);
        if (!tabPane) return;
        
        if (isLoading) {
            // 显示加载状态
            if (!tabPane.querySelector('.tab-loading')) {
                const loadingDiv = document.createElement('div');
                loadingDiv.className = 'tab-loading';
                loadingDiv.innerHTML = `
                    <div class="loading-spinner">
                        <div class="spinner"></div>
                        <p>正在加载${this.getTabDisplayName(tabName)}数据...</p>
                    </div>
                `;
                tabPane.appendChild(loadingDiv);
            }
        } else {
            // 隐藏加载状态
            const loadingDiv = tabPane.querySelector('.tab-loading');
            if (loadingDiv) {
                loadingDiv.remove();
            }
        }
    }

    /**
     * 获取tab的显示名称
     */
    getTabDisplayName(tabName) {
        const nameMap = {
            'daily-data': '日线数据',
            'indicators': '技术指标',
            'predictions': '买卖预测'
        };
        return nameMap[tabName] || tabName;
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

        // 更新日期快捷按钮状态
        if (this.dateShortcutsModule) {
            // 延迟更新，确保DOM已更新
            setTimeout(() => {
                this.dateShortcutsModule.updateActiveButtonByDateRange();
            }, 100);
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
        const requestId = Math.random().toString(36).substr(2, 9);
        
        console.log(`[Events] 开始处理日线数据查询 - ID: ${requestId}`, {
            timestamp: new Date().toISOString()
        });
        
        try {
            this.client.setLoading(true);
            this.client.hideAllResultSections();
            
            const stockCode = this.client.getStockCode();
            const { startDate, endDate } = this.client.getDateRange();
            
            console.log(`[Events] 日线数据查询参数 - ID: ${requestId}`, {
                stockCode,
                startDate,
                endDate,
                timestamp: new Date().toISOString()
            });
            
            // 并行获取股票基本信息和日线数据
            const [stockBasic, dailyData] = await Promise.all([
                this.apiService.getStockBasic(stockCode),
                this.apiService.getDailyData(stockCode, startDate, endDate)
            ]);
            
            console.log(`[Events] 日线数据查询成功 - ID: ${requestId}`, {
                stockCode,
                hasStockBasic: !!stockBasic,
                hasDailyData: !!dailyData,
                dailyDataLength: dailyData ? dailyData.length : 0,
                timestamp: new Date().toISOString()
            });
            
            this.displayModule.displayDailyData(dailyData, stockCode, stockBasic);
            
        } catch (error) {
            // 详细记录错误信息
            const errorDetails = {
                requestId,
                message: error.message,
                errorType: error.constructor.name,
                stack: error.stack,
                timestamp: new Date().toISOString()
            };
            
            console.error(`[Events] 日线数据查询失败 - ID: ${requestId}`, errorDetails);
            this.client.showError(`获取日线数据失败: ${error.message}`);
        } finally {
            this.client.setLoading(false);
            console.log(`[Events] 日线数据查询完成 - ID: ${requestId}`, {
                timestamp: new Date().toISOString()
            });
        }
    }

    /**
     * 处理技术指标查询
     */
    async handleIndicatorsQuery() {
        try {
            this.client.setLoading(true);
            this.client.hideAllResultSections();
            
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
            this.client.hideAllResultSections();
            
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