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
        
        // 添加防抖机制
        this.switchTabTimeout = null;
        this.isSwitching = false;
        
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
     * 清理数据缓存
     */
    clearCache() {
        this.dataCache.clear();
        console.log('[Events] 数据缓存已清理');
    }

    /**
     * 清理指定股票代码的缓存
     */
    clearCacheForStock(stockCode) {
        const keysToDelete = [];
        for (const key of this.dataCache.keys()) {
            if (key.startsWith(stockCode + '_')) {
                keysToDelete.push(key);
            }
        }
        
        keysToDelete.forEach(key => this.dataCache.delete(key));
        console.log(`[Events] 已清理股票 ${stockCode} 的缓存数据`);
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
            
            // 添加预加载功能
            btn.addEventListener('mouseenter', () => {
                this.preloadTabData(btn.dataset.tab);
            });
        });
    }

    /**
     * 预加载tab数据
     */
    async preloadTabData(tabName) {
        const stockCode = this.client.getStockCode();
        if (!stockCode) return;
        
        const cacheKey = `${stockCode}_${tabName}`;
        if (this.dataCache.has(cacheKey)) {
            return; // 数据已缓存，无需预加载
        }
        
        console.log(`[Events] 预加载tab数据: ${tabName}`);
        
        try {
            // 静默预加载数据
            switch (tabName) {
                case 'daily-data':
                    await this.handleDailyQuery(false);
                    break;
                case 'indicators':
                    await this.handleIndicatorsQuery(false);
                    break;
                case 'predictions':
                    await this.handlePredictionsQuery(false);
                    break;
            }
            console.log(`[Events] 预加载完成: ${tabName}`);
        } catch (error) {
            console.warn(`[Events] 预加载失败: ${tabName}`, error);
            // 预加载失败不影响用户体验，只记录日志
        }
    }

    /**
     * 切换Tab
     */
    async switchTab(tabName, loadData = true) {
        console.log(`[Events] 切换到tab: ${tabName}, 是否加载数据: ${loadData}`);
        
        // 防抖处理：如果正在切换，则取消之前的操作
        if (this.isSwitching) {
            console.log(`[Events] Tab切换被取消，正在切换中: ${tabName}`);
            return;
        }
        
        // 清除之前的超时
        if (this.switchTabTimeout) {
            clearTimeout(this.switchTabTimeout);
        }
        
        // 设置切换状态
        this.isSwitching = true;
        
        try {
            // 获取当前激活的tab
            const currentActiveTab = document.querySelector('.tab-btn.active');
            const currentActivePane = document.querySelector('.tab-pane.active');
            
            if (currentActiveTab && currentActivePane) {
                // 添加退出动画
                currentActivePane.classList.add('sliding-out');
                
                // 等待退出动画完成
                await this.waitForTransition(currentActivePane);
            }
            
            // 移除所有tab按钮的active状态
            const tabButtons = document.querySelectorAll('.tab-btn');
            tabButtons.forEach(btn => {
                btn.classList.remove('active');
            });

            // 隐藏所有tab内容
            const tabPanes = document.querySelectorAll('.tab-pane');
            tabPanes.forEach(pane => {
                pane.classList.remove('active', 'sliding-out', 'sliding-in');
            });

            // 激活选中的tab按钮
            const activeTabBtn = document.querySelector(`[data-tab="${tabName}"]`);
            if (activeTabBtn) {
                activeTabBtn.classList.add('active');
            }

            // 显示选中的tab内容
            const activeTabPane = document.getElementById(`${tabName}-tab`);
            if (activeTabPane) {
                // 添加进入动画
                activeTabPane.classList.add('sliding-in');
                activeTabPane.classList.add('active');
                
                // 等待进入动画完成
                await this.waitForTransition(activeTabPane);
                activeTabPane.classList.remove('sliding-in');
                
                // 确保DOM完全准备好
                await this.waitForDOMReady(activeTabPane);
            }
            
            // 根据tab类型获取对应数据（只有在需要时才加载）
            if (loadData) {
                await this.loadDataForTab(tabName);
            }
            
            // 调试：检查tab状态
            this.debugTabState(tabName);
            
            console.log(`[Events] Tab切换完成: ${tabName}`);
            
        } catch (error) {
            console.error(`[Events] Tab切换失败: ${tabName}`, error);
        } finally {
            // 重置切换状态
            this.isSwitching = false;
        }
    }

    /**
     * 调试tab状态
     */
    debugTabState(tabName) {
        const tabPane = document.getElementById(`${tabName}-tab`);
        const tabBtn = document.querySelector(`[data-tab="${tabName}"]`);
        
        console.log(`[Events] Tab状态调试 - ${tabName}:`, {
            tabPane: tabPane ? {
                id: tabPane.id,
                classes: tabPane.className,
                isVisible: tabPane.offsetParent !== null,
                display: window.getComputedStyle(tabPane).display,
                opacity: window.getComputedStyle(tabPane).opacity,
                transform: window.getComputedStyle(tabPane).transform
            } : 'null',
            tabBtn: tabBtn ? {
                classes: tabBtn.className,
                isActive: tabBtn.classList.contains('active')
            } : 'null'
        });
    }

    /**
     * 等待CSS过渡动画完成
     */
    waitForTransition(element) {
        return new Promise(resolve => {
            const duration = 400; // 与CSS中的transition-duration匹配
            setTimeout(resolve, duration);
        });
    }

    /**
     * 等待DOM元素完全准备好
     */
    waitForDOMReady(element) {
        return new Promise(resolve => {
            // 使用 requestAnimationFrame 确保DOM更新完成
            requestAnimationFrame(() => {
                // 再等待一小段时间确保CSS过渡完全完成
                setTimeout(resolve, 50);
            });
        });
    }

    /**
     * 根据tab类型加载对应数据
     */
    async loadDataForTab(tabName) {
        try {
            const stockCode = this.client.getStockCode();
            console.log(`[Events] 为tab ${tabName} 加载数据，股票代码: ${stockCode}`);
            
            if (!stockCode) {
                console.warn(`[Events] 没有股票代码，跳过数据加载`);
                this.showNoDataMessage(tabName, '请先输入股票代码');
                return;
            }
            
            // 检查缓存
            const cacheKey = `${stockCode}_${tabName}`;
            if (this.dataCache.has(cacheKey)) {
                console.log(`[Events] 使用缓存数据: ${cacheKey}`);
                this.displayCachedData(tabName, this.dataCache.get(cacheKey));
                return;
            }
            
            // 显示加载状态
            this.showTabLoadingState(tabName, true);
            
            let data;
            switch (tabName) {
                case 'daily-data':
                    data = await this.handleDailyQuery(false); // 不显示结果，只返回数据
                    break;
                case 'indicators':
                    data = await this.handleIndicatorsQuery(false);
                    break;
                case 'predictions':
                    data = await this.handlePredictionsQuery(false);
                    break;
                default:
                    console.warn(`[Events] 未知的tab类型: ${tabName}`);
                    this.showNoDataMessage(tabName, `未知的tab类型: ${tabName}`);
                    return;
            }
            
            // 缓存数据
            if (data) {
                this.dataCache.set(cacheKey, data);
                console.log(`[Events] 数据已缓存: ${cacheKey}`);
                
                // 显示加载的数据
                this.displayCachedData(tabName, data);
            } else {
                console.warn(`[Events] 没有获取到数据: ${tabName}`);
                this.showNoDataMessage(tabName, '暂无数据');
            }
            
        } catch (error) {
            console.error(`[Events] 为tab ${tabName} 加载数据失败:`, error);
            // 显示错误信息但不阻止tab切换
            this.client.showError(`加载${this.getTabDisplayName(tabName)}数据失败: ${error.message}`);
            this.showNoDataMessage(tabName, `加载失败: ${error.message}`);
        } finally {
            // 隐藏加载状态
            this.showTabLoadingState(tabName, false);
        }
    }

    /**
     * 显示缓存的tab数据
     */
    displayCachedData(tabName, data) {
        const stockCode = this.client.getStockCode();
        console.log(`[Events] 显示缓存数据: ${tabName}`, { data, stockCode });
        
        // 验证数据完整性
        if (!data) {
            console.warn(`[Events] 数据为空: ${tabName}`);
            this.showNoDataMessage(tabName, '数据为空');
            return;
        }
        
        // 确保tab内容区域是可见的
        const tabPane = document.getElementById(`${tabName}-tab`);
        if (!tabPane) {
            console.error(`[Events] 找不到tab面板: ${tabName}-tab`);
            return;
        }
        
        // 确保tab面板是激活状态
        if (!tabPane.classList.contains('active')) {
            console.warn(`[Events] Tab面板未激活: ${tabName}-tab`);
            tabPane.classList.add('active');
        }
        
        // 清除无数据消息
        const existingNoDataMessage = tabPane.querySelector('.no-data-message');
        if (existingNoDataMessage) {
            existingNoDataMessage.remove();
        }
        
        try {
            switch (tabName) {
                case 'daily-data':
                    if (data.dailyData && data.stockBasic) {
                        console.log(`[Events] 显示日线数据`, { dailyDataLength: data.dailyData.length, stockBasic: data.stockBasic });
                        this.displayModule.displayDailyData(data.dailyData, stockCode, data.stockBasic);
                    } else {
                        console.warn(`[Events] 日线数据不完整`, { dailyData: !!data.dailyData, stockBasic: !!data.stockBasic });
                        this.showNoDataMessage(tabName, '日线数据不完整');
                        return;
                    }
                    break;
                case 'indicators':
                    if (data.indicatorsData && data.stockBasic) {
                        console.log(`[Events] 显示技术指标数据`, { indicatorsData: data.indicatorsData, stockBasic: data.stockBasic });
                        this.displayModule.displayIndicators(data.indicatorsData, stockCode, data.stockBasic);
                    } else {
                        console.warn(`[Events] 技术指标数据不完整`, { indicatorsData: !!data.indicatorsData, stockBasic: !!data.stockBasic });
                        this.showNoDataMessage(tabName, '技术指标数据不完整');
                        return;
                    }
                    break;
                case 'predictions':
                    if (data.predictionsData && data.stockBasic) {
                        console.log(`[Events] 显示预测数据`, { predictionsData: data.predictionsData, stockBasic: data.stockBasic });
                        this.displayModule.displayPredictions(data.predictionsData, stockCode, data.stockBasic);
                    } else {
                        console.warn(`[Events] 预测数据不完整`, { predictionsData: !!data.predictionsData, stockBasic: !!data.stockBasic });
                        this.showNoDataMessage(tabName, '预测数据不完整');
                        return;
                    }
                    break;
                default:
                    console.warn(`[Events] 未知的tab类型: ${tabName}`);
                    this.showNoDataMessage(tabName, `未知的tab类型: ${tabName}`);
                    return;
            }
            
            // 添加内容显示动画
            setTimeout(() => {
                this.animateTabContent(tabName);
            }, 100);
            
        } catch (error) {
            console.error(`[Events] 显示数据失败: ${tabName}`, error);
            this.showNoDataMessage(tabName, `显示数据失败: ${error.message}`);
        }
    }

    /**
     * 为tab内容添加显示动画
     */
    animateTabContent(tabName) {
        const tabPane = document.getElementById(`${tabName}-tab`);
        if (!tabPane) {
            console.error(`[Events] 找不到tab面板用于动画: ${tabName}-tab`);
            return;
        }
        
        console.log(`[Events] 开始为tab ${tabName} 添加内容动画`);
        
        // 获取所有需要动画的元素
        const animatedElements = tabPane.querySelectorAll('.chart-container, .indicators-grid, .predictions-container, .data-summary');
        console.log(`[Events] 找到 ${animatedElements.length} 个需要动画的元素`);
        
        if (animatedElements.length === 0) {
            // 如果没有找到预定义的元素，尝试查找所有可能的内容元素
            const allContentElements = tabPane.querySelectorAll('*');
            console.log(`[Events] Tab ${tabName} 中的所有元素:`, allContentElements);
        }
        
        animatedElements.forEach((element, index) => {
            console.log(`[Events] 为元素 ${index} 添加动画:`, element.className || element.tagName);
            
            element.style.opacity = '0';
            element.style.transform = 'translateY(20px)';
            element.style.transition = `all 0.4s cubic-bezier(0.4, 0, 0.2, 1) ${index * 0.1}s`;
            
            // 触发动画
            requestAnimationFrame(() => {
                element.style.opacity = '1';
                element.style.transform = 'translateY(0)';
            });
        });
        
        console.log(`[Events] Tab ${tabName} 的内容动画设置完成`);
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
                        <div class="loading-dots">
                            <span></span>
                            <span></span>
                            <span></span>
                        </div>
                    </div>
                `;
                
                // 添加淡入动画
                loadingDiv.style.opacity = '0';
                tabPane.appendChild(loadingDiv);
                
                // 触发重排后添加淡入效果
                requestAnimationFrame(() => {
                    loadingDiv.style.transition = 'opacity 0.3s ease';
                    loadingDiv.style.opacity = '1';
                });
            }
        } else {
            // 隐藏加载状态
            const loadingDiv = tabPane.querySelector('.tab-loading');
            if (loadingDiv) {
                // 添加淡出动画
                loadingDiv.style.transition = 'opacity 0.3s ease';
                loadingDiv.style.opacity = '0';
                
                // 等待动画完成后移除元素
                setTimeout(() => {
                    if (loadingDiv.parentNode) {
                        loadingDiv.remove();
                    }
                }, 300);
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
        
        // 检查股票代码是否改变
        if (this.lastStockCode && this.lastStockCode !== code) {
            console.log(`[Events] 股票代码从 ${this.lastStockCode} 改变为 ${code}，清理缓存`);
            this.clearCacheForStock(this.lastStockCode);
        }
        
        if (code && !pattern.test(code)) {
            event.target.setCustomValidity('请输入正确的股票代码格式，如：000001.SZ 或 600000.SH');
        } else {
            event.target.setCustomValidity('');
        }
        
        // 更新最后使用的股票代码
        this.lastStockCode = code;
        
        // 更新收藏按钮状态
        if (this.favoritesModule) {
            this.favoritesModule.updateFavoriteButtonState();
        }
    }

    /**
     * 处理日线数据查询
     */
    async handleDailyQuery(displayResult = true) {
        const requestId = Math.random().toString(36).substr(2, 9);
        
        console.log(`[Events] 开始处理日线数据查询 - ID: ${requestId}`, {
            timestamp: new Date().toISOString()
        });
        
        try {
            if (displayResult) {
                this.client.setLoading(true);
                this.client.hideAllResultSections();
            }
            
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
            
            if (displayResult) {
                this.displayModule.displayDailyData(dailyData, stockCode, stockBasic);
            }
            
            // 返回数据用于缓存
            return {
                dailyData,
                stockBasic
            };
            
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
            if (displayResult) {
                this.client.showError(`获取日线数据失败: ${error.message}`);
            }
            throw error;
        } finally {
            if (displayResult) {
                this.client.setLoading(false);
            }
            console.log(`[Events] 日线数据查询完成 - ID: ${requestId}`, {
                timestamp: new Date().toISOString()
            });
        }
    }

    /**
     * 处理技术指标查询
     */
    async handleIndicatorsQuery(displayResult = true) {
        try {
            if (displayResult) {
                this.client.setLoading(true);
                this.client.hideAllResultSections();
            }
            
            const stockCode = this.client.getStockCode();
            
            // 并行获取股票基本信息和技术指标
            const [stockBasic, indicatorsData] = await Promise.all([
                this.apiService.getStockBasic(stockCode),
                this.apiService.getIndicators(stockCode)
            ]);
            
            if (displayResult) {
                this.displayModule.displayIndicators(indicatorsData, stockCode, stockBasic);
            }
            
            // 返回数据用于缓存
            return {
                indicatorsData,
                stockBasic
            };
            
        } catch (error) {
            if (displayResult) {
                this.client.showError(`获取技术指标失败: ${error.message}`);
            }
            throw error;
        } finally {
            if (displayResult) {
                this.client.setLoading(false);
            }
        }
    }

    /**
     * 处理买卖预测查询
     */
    async handlePredictionsQuery(displayResult = true) {
        try {
            if (displayResult) {
                this.client.setLoading(true);
                this.client.hideAllResultSections();
            }
            
            const stockCode = this.client.getStockCode();
            
            // 并行获取股票基本信息和预测数据
            const [stockBasic, predictionsData] = await Promise.all([
                this.apiService.getStockBasic(stockCode),
                this.apiService.getPredictions(stockCode)
            ]);
            
            if (displayResult) {
                this.displayModule.displayPredictions(predictionsData, stockCode, stockBasic);
            }
            
            // 返回数据用于缓存
            return {
                predictionsData,
                stockBasic
            };
            
        } catch (error) {
            if (displayResult) {
                this.client.showError(`获取买卖预测失败: ${error.message}`);
            }
            throw error;
        } finally {
            if (displayResult) {
                this.client.setLoading(false);
            }
        }
    }

    /**
     * 显示无数据消息
     */
    showNoDataMessage(tabName, message) {
        const tabPane = document.getElementById(`${tabName}-tab`);
        if (!tabPane) return;
        
        // 清除现有内容
        const existingContent = tabPane.querySelector('.no-data-message');
        if (existingContent) {
            existingContent.remove();
        }
        
        const noDataDiv = document.createElement('div');
        noDataDiv.className = 'no-data-message';
        noDataDiv.innerHTML = `
            <div class="no-data">
                <p>${message}</p>
            </div>
        `;
        
        tabPane.appendChild(noDataDiv);
        console.log(`[Events] 显示无数据消息: ${tabName} - ${message}`);
    }
}

// 导出事件处理模块类
window.EventsModule = EventsModule;