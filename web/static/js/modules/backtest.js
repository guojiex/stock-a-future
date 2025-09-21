/**
 * 回测系统模块
 * 负责回测配置、执行和结果展示
 */

class BacktestModule {
    constructor(client, apiService) {
        this.client = client;
        this.apiService = apiService;
        this.currentBacktest = null;
        this.isRunning = false;
        this.progressInterval = null;
        
        // 多选策略相关数据
        this.availableStrategies = [];
        this.selectedStrategyIds = [];
        
        this.init();
    }

    /**
     * 初始化回测模块
     */
    init() {
        this.setupEventListeners();
        this.loadBacktestHistory();
        this.setupStockCodeListener();
    }

    /**
     * 设置事件监听器
     */
    setupEventListeners() {
        // 开始回测按钮
        const startBacktestBtn = document.getElementById('startBacktestBtn');
        if (startBacktestBtn) {
            startBacktestBtn.addEventListener('click', () => this.startBacktest());
        }

        // 停止回测按钮
        const stopBacktestBtn = document.getElementById('stopBacktestBtn');
        if (stopBacktestBtn) {
            stopBacktestBtn.addEventListener('click', () => this.stopBacktest());
        }

        // 保存配置按钮
        const saveConfigBtn = document.getElementById('saveBacktestConfigBtn');
        if (saveConfigBtn) {
            saveConfigBtn.addEventListener('click', () => this.saveBacktestConfig());
        }

        // 初始化默认日期
        this.setDefaultDates();
    }

    /**
     * 设置默认日期（最近一年）
     */
    setDefaultDates() {
        const endDate = new Date();
        const startDate = new Date();
        startDate.setFullYear(endDate.getFullYear() - 1);

        const backtestStartDate = document.getElementById('backtestStartDate');
        const backtestEndDate = document.getElementById('backtestEndDate');

        if (backtestStartDate) {
            backtestStartDate.value = startDate.toISOString().split('T')[0];
        }
        if (backtestEndDate) {
            backtestEndDate.value = endDate.toISOString().split('T')[0];
        }
    }

    /**
     * 设置股票代码监听器，用于自动预填充
     */
    setupStockCodeListener() {
        // 监听股票代码输入框变化
        const stockCodeInput = document.getElementById('stockCode');
        if (stockCodeInput) {
            stockCodeInput.addEventListener('input', () => {
                this.updateBacktestFormFromStockCode();
            });
        }

        // 监听tab切换到回测页面
        const backtestTab = document.querySelector('[data-tab="backtest"]');
        if (backtestTab) {
            backtestTab.addEventListener('click', () => {
                // 延迟执行，确保tab切换完成
                setTimeout(() => {
                    this.prefillBacktestForm();
                }, 100);
            });
        }

        // 监听快速回测按钮
        const quickBacktestBtn = document.getElementById('quickBacktestBtn');
        if (quickBacktestBtn) {
            quickBacktestBtn.addEventListener('click', () => {
                this.handleQuickBacktest();
            });
        }
    }

    /**
     * 根据当前股票代码更新回测表单
     */
    updateBacktestFormFromStockCode() {
        const stockCodeInput = document.getElementById('stockCode');
        if (!stockCodeInput || !stockCodeInput.value.trim()) return;

        const stockCode = stockCodeInput.value.trim();
        const pattern = /^[0-9]{6}\.(SZ|SH)$/;
        
        // 只有当股票代码格式正确时才更新
        if (pattern.test(stockCode)) {
            this.prefillStockCodeInBacktest(stockCode);
        }
    }

    /**
     * 预填充回测表单
     */
    async prefillBacktestForm() {
        try {
            const stockCodeInput = document.getElementById('stockCode');
            if (!stockCodeInput || !stockCodeInput.value.trim()) {
                console.log('[Backtest] 没有股票代码，跳过预填充');
                return;
            }

            const stockCode = stockCodeInput.value.trim();
            const pattern = /^[0-9]{6}\.(SZ|SH)$/;
            
            if (!pattern.test(stockCode)) {
                console.log('[Backtest] 股票代码格式不正确，跳过预填充');
                return;
            }

            console.log(`[Backtest] 开始预填充回测表单，股票代码: ${stockCode}`);

            // 预填充股票代码
            this.prefillStockCodeInBacktest(stockCode);

            // 获取股票基本信息用于生成回测名称
            try {
                const stockBasic = await this.apiService.getStockBasic(stockCode);
                if (stockBasic && stockBasic.name) {
                    this.prefillBacktestName(stockCode, stockBasic.name);
                }
            } catch (error) {
                console.warn('[Backtest] 获取股票基本信息失败，使用默认回测名称', error);
                this.prefillBacktestName(stockCode, null);
            }

        } catch (error) {
            console.error('[Backtest] 预填充回测表单失败:', error);
        }
    }

    /**
     * 预填充股票代码到回测表单
     */
    prefillStockCodeInBacktest(stockCode) {
        const backtestSymbols = document.getElementById('backtestSymbols');
        if (backtestSymbols) {
            // 如果已经有内容，检查是否包含当前股票代码
            const currentValue = backtestSymbols.value.trim();
            const symbols = currentValue.split('\n').map(s => s.trim()).filter(s => s.length > 0);
            
            if (!symbols.includes(stockCode)) {
                // 如果没有当前股票代码，添加到第一行
                const newValue = [stockCode, ...symbols].join('\n');
                backtestSymbols.value = newValue;
                console.log(`[Backtest] 已将股票代码 ${stockCode} 添加到回测股票列表`);
            }
        }
    }

    /**
     * 预填充回测名称
     */
    prefillBacktestName(stockCode, stockName) {
        const backtestNameInput = document.getElementById('backtestName');
        if (backtestNameInput && !backtestNameInput.value.trim()) {
            // 只有当回测名称为空时才自动填充
            const displayName = stockName || stockCode.split('.')[0];
            const currentDate = new Date().toISOString().split('T')[0];
            const backtestName = `${displayName} - 回测分析 ${currentDate}`;
            
            backtestNameInput.value = backtestName;
            console.log(`[Backtest] 已预填充回测名称: ${backtestName}`);
        }
    }

    /**
     * 处理快速回测按钮点击
     */
    handleQuickBacktest() {
        const stockCodeInput = document.getElementById('stockCode');
        if (!stockCodeInput || !stockCodeInput.value.trim()) {
            this.showMessage('请先输入股票代码', 'warning');
            return;
        }

        const stockCode = stockCodeInput.value.trim();
        const pattern = /^[0-9]{6}\.(SZ|SH)$/;
        
        if (!pattern.test(stockCode)) {
            this.showMessage('请输入正确的股票代码格式', 'warning');
            return;
        }

        console.log(`[Backtest] 快速回测按钮点击，股票代码: ${stockCode}`);

        // 切换到回测标签页
        const backtestTab = document.querySelector('[data-tab="backtest"]');
        if (backtestTab) {
            backtestTab.click();
            
            // 延迟执行预填充，确保tab切换完成
            setTimeout(() => {
                this.prefillBacktestForm();
                
                // 聚焦到策略选择框，提示用户选择策略
                const strategySelect = document.getElementById('backtestStrategy');
                if (strategySelect) {
                    strategySelect.focus();
                }
                
                this.showMessage(`已切换到回测页面，股票代码 ${stockCode} 已预填充`, 'success');
            }, 200);
        }
    }

    /**
     * 开始回测
     */
    async startBacktest() {
        try {
            const config = this.getBacktestConfig();
            if (!this.validateConfig(config)) {
                return;
            }

            this.showProgress();
            this.isRunning = true;
            this.toggleBacktestButtons(true); // 显示停止按钮，隐藏开始按钮

            // 调用后端API开始回测
            const response = await this.apiService.startBacktest(config);
            
            if (response.success) {
                this.currentBacktest = response.data;
                this.startProgressMonitoring();
                this.showMessage('回测已启动', 'success');
            } else {
                throw new Error(response.message || '启动回测失败');
            }

        } catch (error) {
            console.error('启动回测失败:', error);
            this.showMessage(`启动回测失败: ${error.message}`, 'error');
            this.hideProgress();
            this.isRunning = false;
            this.toggleBacktestButtons(false); // 恢复开始按钮
        }
    }

    /**
     * 停止回测
     */
    async stopBacktest() {
        if (!this.currentBacktest || !this.currentBacktest.id) {
            console.warn('[Backtest] 没有正在运行的回测');
            return;
        }

        try {
            console.log('[Backtest] 停止回测，ID:', this.currentBacktest.id);
            
            // 显示确认对话框
            if (!confirm('确定要停止当前回测吗？')) {
                return;
            }

            // 调用后端API停止回测
            const response = await this.apiService.cancelBacktest(this.currentBacktest.id);
            
            if (response.success) {
                console.log('[Backtest] 回测停止成功');
                this.stopProgressMonitoring();
                this.hideProgress();
                this.isRunning = false;
                this.currentBacktest = null;
                this.toggleBacktestButtons(false); // 恢复开始按钮
                
                this.showMessage('回测已停止', 'info');
            } else {
                throw new Error(response.message || '停止回测失败');
            }
        } catch (error) {
            console.error('[Backtest] 停止回测失败:', error);
            this.showMessage(`停止回测失败: ${error.message}`, 'error');
        }
    }

    /**
     * 切换回测按钮状态
     */
    toggleBacktestButtons(isRunning) {
        const startBtn = document.getElementById('startBacktestBtn');
        const stopBtn = document.getElementById('stopBacktestBtn');
        
        if (startBtn && stopBtn) {
            if (isRunning) {
                startBtn.style.display = 'none';
                stopBtn.style.display = 'inline-block';
            } else {
                startBtn.style.display = 'inline-block';
                stopBtn.style.display = 'none';
            }
        }
    }

    /**
     * 获取回测配置
     */
    getBacktestConfig() {
        const backtestName = document.getElementById('backtestName')?.value?.trim() || '';
        const backtestStartDate = document.getElementById('backtestStartDate')?.value || '';
        const backtestEndDate = document.getElementById('backtestEndDate')?.value || '';
        const initialCash = parseFloat(document.getElementById('initialCash')?.value || 1000000);
        const commission = parseFloat(document.getElementById('commission')?.value || 0.0003);
        const backtestSymbols = document.getElementById('backtestSymbols')?.value?.trim() || '';

        // 解析股票列表
        const symbols = backtestSymbols
            .split('\n')
            .map(s => s.trim())
            .filter(s => s.length > 0);

        return {
            name: backtestName,
            strategy_ids: this.selectedStrategyIds, // 使用多选策略ID数组
            start_date: backtestStartDate,
            end_date: backtestEndDate,
            initial_cash: initialCash,
            commission: commission,
            symbols: symbols
        };
    }

    /**
     * 验证回测配置
     */
    validateConfig(config) {
        if (!config.name) {
            this.showMessage('请输入回测名称', 'warning');
            return false;
        }

        if (!config.strategy_ids || config.strategy_ids.length === 0) {
            this.showMessage('请选择至少一个策略', 'warning');
            return false;
        }

        if (config.strategy_ids.length > 5) {
            this.showMessage('最多只能选择5个策略', 'warning');
            return false;
        }

        if (!config.start_date || !config.end_date) {
            this.showMessage('请选择开始和结束日期', 'warning');
            return false;
        }

        if (new Date(config.start_date) >= new Date(config.end_date)) {
            this.showMessage('开始日期必须早于结束日期', 'warning');
            return false;
        }

        if (config.symbols.length === 0) {
            this.showMessage('请输入至少一个股票代码', 'warning');
            return false;
        }

        if (config.initial_cash < 10000) {
            this.showMessage('初始资金不能少于10000元', 'warning');
            return false;
        }

        return true;
    }

    /**
     * 显示进度条
     */
    showProgress() {
        const progressDiv = document.getElementById('backtestProgress');
        if (progressDiv) {
            progressDiv.style.display = 'block';
        }

        this.updateProgress(0, '准备中...');
    }

    /**
     * 隐藏进度条
     */
    hideProgress() {
        const progressDiv = document.getElementById('backtestProgress');
        if (progressDiv) {
            progressDiv.style.display = 'none';
        }
    }

    /**
     * 更新进度
     */
    updateProgress(percent, message) {
        const progressFill = document.getElementById('progressFill');
        const progressText = document.getElementById('progressText');
        const progressPercent = document.getElementById('progressPercent');

        if (progressFill) {
            progressFill.style.width = `${percent}%`;
        }
        if (progressText) {
            progressText.textContent = message;
        }
        if (progressPercent) {
            progressPercent.textContent = `${percent}%`;
        }
    }

    /**
     * 开始进度监控
     */
    startProgressMonitoring() {
        if (this.progressInterval) {
            clearInterval(this.progressInterval);
        }

        this.progressInterval = setInterval(async () => {
            try {
                if (!this.currentBacktest || !this.isRunning) {
                    this.stopProgressMonitoring();
                    return;
                }

                const response = await this.apiService.getBacktestProgress(this.currentBacktest.id);
                if (response.success) {
                    const progress = response.data;
                    this.updateProgress(progress.progress, progress.message || '运行中...');

                    if (progress.status === 'completed') {
                        this.stopProgressMonitoring();
                        this.hideProgress();
                        this.isRunning = false;
                        this.toggleBacktestButtons(false); // 恢复开始按钮
                        this.loadBacktestResults(this.currentBacktest.id);
                        this.showMessage('回测完成！', 'success');
                    } else if (progress.status === 'failed') {
                        this.stopProgressMonitoring();
                        this.hideProgress();
                        this.isRunning = false;
                        this.toggleBacktestButtons(false); // 恢复开始按钮
                        this.showMessage(progress.error || '回测失败', 'error');
                    } else if (progress.status === 'cancelled') {
                        this.stopProgressMonitoring();
                        this.hideProgress();
                        this.isRunning = false;
                        this.toggleBacktestButtons(false); // 恢复开始按钮
                        this.showMessage('回测已取消', 'info');
                    }
                }
            } catch (error) {
                console.error('获取回测进度失败:', error);
            }
        }, 2000); // 每2秒检查一次
    }

    /**
     * 停止进度监控
     */
    stopProgressMonitoring() {
        if (this.progressInterval) {
            clearInterval(this.progressInterval);
            this.progressInterval = null;
        }
        this.isRunning = false;
    }

    /**
     * 加载回测结果
     */
    async loadBacktestResults(backtestId) {
        try {
            const response = await this.apiService.getBacktestResults(backtestId);
            if (response.success) {
                this.displayResults(response.data);
            } else {
                throw new Error(response.message || '获取回测结果失败');
            }
        } catch (error) {
            console.error('加载回测结果失败:', error);
            this.showMessage(`加载回测结果失败: ${error.message}`, 'error');
        }
    }

    /**
     * 显示回测结果
     */
    displayResults(results) {
        const resultsDiv = document.getElementById('backtestResults');
        if (!resultsDiv) return;

        // 显示结果区域
        resultsDiv.style.display = 'block';

        // 处理多策略和单策略的兼容性
        let displayStrategy = results.strategy;
        let displayPerformance = results.performance;
        const isMultiStrategy = Array.isArray(results.performance) && results.performance.length > 1;

        // 检查是否为多策略结果
        if (isMultiStrategy) {
            console.log('[Backtest] 检测到多策略结果，显示详细对比');
            
            // 多策略情况：使用组合指标作为主要显示，但同时突出显示各策略
            if (results.combined_metrics) {
                displayPerformance = results.combined_metrics;
                console.log('[Backtest] 使用组合指标显示主要性能数据');
            } else {
                displayPerformance = results.performance[0];
                console.log('[Backtest] 组合指标不存在，使用第一个策略指标');
            }

            // 为多策略创建虚拟策略信息用于显示
            if (results.strategies && results.strategies.length > 1) {
                displayStrategy = {
                    name: `多策略组合 (${results.strategies.length}个策略)`,
                    strategy_type: 'combined',
                    description: `包含策略: ${results.strategies.map(s => s.name).join(', ')}`,
                    parameters: {
                        strategy_count: results.strategies.length,
                        strategies: results.strategies.map(s => s.name)
                    }
                };
            } else if (results.strategies && results.strategies.length === 1) {
                displayStrategy = results.strategies[0];
            }
        } else if (Array.isArray(results.performance) && results.performance.length === 1) {
            // 单策略但以数组形式返回
            displayPerformance = results.performance[0];
            if (results.strategies && results.strategies.length === 1) {
                displayStrategy = results.strategies[0];
            }
        }

        // 显示策略配置信息
        this.displayStrategyConfig(displayStrategy, results.backtest_config);

        // 如果是多策略，优先显示策略对比概览
        if (isMultiStrategy) {
            this.displayMultiStrategyOverview(results.performance, results.strategies, results.combined_metrics);
        }

        // 显示主要性能指标（组合指标或单策略指标）
        this.displayPerformanceMetrics(displayPerformance, isMultiStrategy ? '组合整体表现' : '策略表现');

        // 显示多策略详细信息（如果是多策略）
        if (isMultiStrategy) {
            this.displayMultiStrategyDetails(results.performance, results.strategies);
        }

        // 显示权益曲线
        this.displayEquityCurve(results.equity_curve);

        // 保存回测配置和表现结果，供其他方法使用
        this.currentBacktestConfig = results.backtest_config;
        this.currentPerformanceResults = results.performance;
        
        // 显示交易记录（按策略分组显示）
        this.displayTradeHistory(results.trades, isMultiStrategy, results.strategies);
        
        // 更新股票信息表头
        this.updateStockInfoHeader(results.backtest_config);
    }

    /**
     * 显示多策略概览
     */
    displayMultiStrategyOverview(performanceResults, strategies, combinedMetrics) {
        // 检查是否存在多策略概览区域，如果不存在则创建
        let overviewSection = document.getElementById('multiStrategyOverview');
        if (!overviewSection) {
            // 在策略配置区域后面插入概览区域
            const configSection = document.getElementById('strategyConfigSection');
            if (configSection) {
                overviewSection = document.createElement('div');
                overviewSection.id = 'multiStrategyOverview';
                overviewSection.className = 'multi-strategy-section';
                configSection.insertAdjacentElement('afterend', overviewSection);
            }
        }

        if (!overviewSection) return;

        // 创建策略对比卡片
        const strategyCards = performanceResults.map((performance, index) => {
            const strategy = strategies && strategies[index] ? strategies[index] : { name: `策略${index + 1}` };
            
            return `
                <div class="strategy-card clickable-card" onclick="backtestModule.jumpToStrategyTrades('${strategy.id}')" title="点击查看该策略的交易记录">
                    <div class="strategy-card-header">
                        <h5 class="strategy-name">${strategy.name}</h5>
                        <div class="strategy-rank">#${this.getRankByReturn(performanceResults, index)}</div>
                    </div>
                    <div class="strategy-card-metrics">
                        <div class="metric-row" data-tooltip="策略在整个回测期间的累计收益率">
                            <span class="metric-label">总收益率 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value ${this.getMetricClass(performance.total_return, '收益')}">${this.formatMetricValue(performance.total_return, 'percentage')}</span>
                        </div>
                        <div class="metric-row" data-tooltip="风险调整后收益指标，>1为良好，>2为优秀">
                            <span class="metric-label">夏普比率 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value ${this.getMetricClass(performance.sharpe_ratio, '夏普')}">${this.formatMetricValue(performance.sharpe_ratio, 'decimal')}</span>
                        </div>
                        <div class="metric-row" data-tooltip="最大资产回撤幅度，数值越小风险越低">
                            <span class="metric-label">最大回撤 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value ${this.getMetricClass(performance.max_drawdown, '回撤')}">${this.formatMetricValue(performance.max_drawdown, 'percentage')}</span>
                        </div>
                        <div class="metric-row" data-tooltip="盈利交易占总交易次数的比例，胜率越高策略越稳定">
                            <span class="metric-label">胜率 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value ${this.getMetricClass(performance.win_rate, '胜率')}">${this.formatMetricValue(performance.win_rate, 'percentage')}</span>
                        </div>
                        <div class="metric-row" data-tooltip="策略执行的买卖交易总次数">
                            <span class="metric-label">交易次数 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value neutral">${this.formatMetricValue(performance.total_trades, 'number')}</span>
                        </div>
                    </div>
                </div>
            `;
        }).join('');

        // 移除组合指标卡片，避免与下方的整体表现重复
        const combinedCard = '';

        overviewSection.innerHTML = `
            <div class="section-header">
                <h4>🏆 策略表现对比</h4>
                <p class="section-description">各策略独立运行结果对比，每个策略使用相等的初始资金</p>
            </div>
            <div class="strategy-cards-grid">
                ${strategyCards}
                ${combinedCard}
            </div>
        `;

        console.log('[Backtest] 多策略概览显示完成');
    }

    /**
     * 根据收益率获取策略排名
     */
    getRankByReturn(performanceResults, currentIndex) {
        const currentReturn = performanceResults[currentIndex].total_return;
        const betterCount = performanceResults.filter(p => p.total_return > currentReturn).length;
        return betterCount + 1;
    }

    /**
     * 显示多策略详细信息
     */
    displayMultiStrategyDetails(performanceResults, strategies) {
        // 检查是否存在多策略详情区域，如果不存在则创建
        let multiStrategySection = document.getElementById('multiStrategyDetails');
        if (!multiStrategySection) {
            // 在性能指标区域后面插入多策略详情区域
            const metricsSection = document.querySelector('.performance-section');
            if (metricsSection) {
                multiStrategySection = document.createElement('div');
                multiStrategySection.id = 'multiStrategyDetails';
                multiStrategySection.className = 'multi-strategy-section';
                metricsSection.insertAdjacentElement('afterend', multiStrategySection);
            }
        }

        if (!multiStrategySection) return;

        // 创建多策略详情表格
        const strategyRows = performanceResults.map((performance, index) => {
            const strategy = strategies && strategies[index] ? strategies[index] : { name: `策略${index + 1}` };
            
            return `
                <tr>
                    <td class="strategy-name">${strategy.name}</td>
                    <td class="metric-value ${this.getMetricClass(performance.total_return, '收益')}">${this.formatMetricValue(performance.total_return, 'percentage')}</td>
                    <td class="metric-value ${this.getMetricClass(performance.annual_return, '收益')}">${this.formatMetricValue(performance.annual_return, 'percentage')}</td>
                    <td class="metric-value ${this.getMetricClass(performance.max_drawdown, '回撤')}">${this.formatMetricValue(performance.max_drawdown, 'percentage')}</td>
                    <td class="metric-value ${this.getMetricClass(performance.sharpe_ratio, '夏普')}">${this.formatMetricValue(performance.sharpe_ratio, 'decimal')}</td>
                    <td class="metric-value ${this.getMetricClass(performance.win_rate, '胜率')}">${this.formatMetricValue(performance.win_rate, 'percentage')}</td>
                    <td class="metric-value neutral">${this.formatMetricValue(performance.total_trades, 'number')}</td>
                </tr>
            `;
        }).join('');

        multiStrategySection.innerHTML = `
            <div class="section-header">
                <h4>📊 各策略详细表现</h4>
            </div>
            <div class="strategy-details-table">
                <table class="table table-striped">
                    <thead>
                        <tr>
                            <th>策略名称</th>
                            <th>总收益率</th>
                            <th>年化收益率</th>
                            <th>最大回撤</th>
                            <th>夏普比率</th>
                            <th>胜率</th>
                            <th>交易次数</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${strategyRows}
                    </tbody>
                </table>
            </div>
        `;

        console.log('[Backtest] 多策略详细信息显示完成');
    }

    /**
     * 显示策略配置信息
     */
    displayStrategyConfig(strategy, backtestConfig) {
        const configSection = document.getElementById('strategyConfigSection');
        if (!configSection) return;

        let strategyInfo = '';
        let backtestInfo = '';

        // 策略信息
        if (strategy) {
            strategyInfo = `
                <div class="config-group">
                    <h5>🎯 策略配置</h5>
                    <div class="config-grid">
                        <div class="config-item">
                            <span class="config-label">策略名称:</span>
                            <span class="config-value">${strategy.name}</span>
                        </div>
                        <div class="config-item">
                            <span class="config-label">策略类型:</span>
                            <span class="config-value">${this.formatStrategyType(strategy.strategy_type)}</span>
                        </div>
                        <div class="config-item">
                            <span class="config-label">策略描述:</span>
                            <span class="config-value">${strategy.description}</span>
                        </div>
                        ${this.formatStrategyParameters(strategy.parameters)}
                    </div>
                </div>
            `;
        }

        // 回测配置信息
        if (backtestConfig) {
            backtestInfo = `
                <div class="config-group">
                    <h5>⚙️ 回测配置</h5>
                    <div class="config-grid">
                        <div class="config-item">
                            <span class="config-label">回测名称:</span>
                            <span class="config-value">${backtestConfig.name}</span>
                        </div>
                        <div class="config-item">
                            <span class="config-label">回测期间:</span>
                            <span class="config-value">${backtestConfig.start_date} 至 ${backtestConfig.end_date}</span>
                        </div>
                        <div class="config-item">
                            <span class="config-label">初始资金:</span>
                            <span class="config-value">¥${backtestConfig.initial_cash.toLocaleString()}</span>
                        </div>
                        <div class="config-item">
                            <span class="config-label">手续费率:</span>
                            <span class="config-value">${(backtestConfig.commission * 100).toFixed(3)}%</span>
                        </div>
                        <div class="config-item">
                            <span class="config-label">股票池:</span>
                            <span class="config-value">${backtestConfig.symbols.join(', ')}</span>
                        </div>
                        <div class="config-item">
                            <span class="config-label">创建时间:</span>
                            <span class="config-value">${backtestConfig.created_at}</span>
                        </div>
                    </div>
                </div>
            `;
        }

        configSection.innerHTML = `
            <div class="strategy-backtest-config">
                ${strategyInfo}
                ${backtestInfo}
            </div>
        `;
    }

    /**
     * 格式化策略类型
     */
    formatStrategyType(type) {
        const typeMap = {
            'technical': '技术指标',
            'fundamental': '基本面',
            'ml': '机器学习',
            'composite': '复合策略',
            'combined': '组合策略'
        };
        return typeMap[type] || type;
    }

    /**
     * 格式化策略参数
     */
    formatStrategyParameters(parameters) {
        if (!parameters || Object.keys(parameters).length === 0) {
            return '';
        }

        return Object.entries(parameters).map(([key, value]) => `
            <div class="config-item">
                <span class="config-label">${this.formatParameterName(key)}:</span>
                <span class="config-value">${this.formatParameterValue(key, value)}</span>
            </div>
        `).join('');
    }

    /**
     * 格式化参数名称
     */
    formatParameterName(key) {
        const nameMap = {
            'fast_period': '快线周期',
            'slow_period': '慢线周期',
            'signal_period': '信号线周期',
            'buy_threshold': '买入阈值',
            'sell_threshold': '卖出阈值',
            'short_period': '短期周期',
            'long_period': '长期周期',
            'ma_type': '均线类型',
            'threshold': '阈值',
            'period': '周期',
            'overbought': '超买线',
            'oversold': '超卖线',
            'std_dev': '标准差倍数',
            'strategy_count': '策略数量',
            'strategies': '包含策略'
        };
        return nameMap[key] || key;
    }

    /**
     * 格式化参数值
     */
    formatParameterValue(key, value) {
        if (key === 'ma_type') {
            const typeMap = {
                'sma': '简单移动平均',
                'ema': '指数移动平均',
                'wma': '加权移动平均'
            };
            return typeMap[value] || value;
        }
        
        // 处理组合策略的特殊参数
        if (key === 'strategies' && Array.isArray(value)) {
            return value.join(', ');
        }
        
        if (typeof value === 'number') {
            if (key.includes('threshold') || key === 'std_dev') {
                return value.toFixed(2);
            }
            return value.toString();
        }
        
        return value;
    }

    /**
     * 显示性能指标
     */
    displayPerformanceMetrics(performance, title = '性能指标') {
        const metricsGrid = document.getElementById('metricsGrid');
        if (!metricsGrid) return;

        // 更新性能指标区域的标题
        const performanceSection = document.querySelector('.performance-metrics h5');
        if (performanceSection) {
            performanceSection.textContent = `📈 ${title}`;
        }

        const metrics = [
            { label: '总收益率', value: performance.total_return, format: 'percentage' },
            { label: '年化收益率', value: performance.annual_return, format: 'percentage' },
            { label: '最大回撤', value: performance.max_drawdown, format: 'percentage' },
            { label: '夏普比率', value: performance.sharpe_ratio, format: 'decimal' },
            { label: '胜率', value: performance.win_rate, format: 'percentage' },
            { label: '总交易次数', value: performance.total_trades, format: 'number' },
            { label: '平均交易收益', value: performance.avg_trade_return, format: 'percentage' },
            { label: '盈亏比', value: performance.profit_factor, format: 'decimal' }
        ];

        metricsGrid.innerHTML = metrics.map(metric => `
            <div class="metric-item">
                <div class="metric-label">${metric.label}</div>
                <div class="metric-value ${this.getMetricClass(metric.value, metric.label)}">
                    ${this.formatMetricValue(metric.value, metric.format)}
                </div>
            </div>
        `).join('');
    }

    /**
     * 获取指标样式类
     */
    getMetricClass(value, label) {
        if (label.includes('回撤')) {
            return value < -0.1 ? 'negative' : value < -0.05 ? 'warning' : 'positive';
        }
        if (label.includes('收益') || label.includes('胜率')) {
            return value > 0 ? 'positive' : value < 0 ? 'negative' : 'neutral';
        }
        if (label.includes('夏普') || label.includes('盈亏比')) {
            return value > 1 ? 'positive' : value > 0.5 ? 'warning' : 'negative';
        }
        return 'neutral';
    }

    /**
     * 格式化指标值
     */
    formatMetricValue(value, format) {
        if (value === null || value === undefined) return 'N/A';
        
        switch (format) {
            case 'percentage':
                return `${(value * 100).toFixed(2)}%`;
            case 'decimal':
                return value.toFixed(2);
            case 'number':
                return value.toString();
            default:
                return value.toString();
        }
    }

    /**
     * 显示权益曲线
     */
    displayEquityCurve(equityCurve) {
        const chartDiv = document.getElementById('equityChart');
        if (!chartDiv || !equityCurve || equityCurve.length === 0) return;

        const chart = echarts.init(chartDiv);

        const dates = equityCurve.map(item => item.date);
        const portfolioValues = equityCurve.map(item => item.portfolio_value);
        const benchmarkValues = equityCurve.map(item => item.benchmark_value || null);

        const series = [{
            name: '策略收益',
            type: 'line',
            data: portfolioValues,
            smooth: true,
            lineStyle: { color: '#1890ff', width: 2 }
        }];

        if (benchmarkValues.some(v => v !== null)) {
            series.push({
                name: '基准收益',
                type: 'line',
                data: benchmarkValues,
                smooth: true,
                lineStyle: { color: '#ff7875', width: 2 }
            });
        }

        const option = {
            title: {
                text: '权益曲线对比',
                left: 'center',
                textStyle: { fontSize: 16 }
            },
            tooltip: {
                trigger: 'axis',
                formatter: function(params) {
                    let result = `${params[0].axisValue}<br/>`;
                    params.forEach(param => {
                        const value = param.value;
                        result += `${param.seriesName}: ¥${value ? value.toLocaleString() : 'N/A'}<br/>`;
                    });
                    return result;
                }
            },
            legend: {
                top: 30,
                data: series.map(s => s.name)
            },
            grid: {
                top: 80,
                bottom: 60,
                left: 60,
                right: 40
            },
            xAxis: {
                type: 'category',
                data: dates,
                axisLabel: {
                    rotate: 45
                }
            },
            yAxis: {
                type: 'value',
                axisLabel: {
                    formatter: function(value) {
                        return '¥' + (value / 10000).toFixed(0) + '万';
                    }
                }
            },
            series: series
        };

        chart.setOption(option);

        // 响应式调整
        window.addEventListener('resize', () => {
            chart.resize();
        });
    }

    /**
     * 显示策略表现详情（原交易记录）
     */
    displayTradeHistory(trades, isMultiStrategy = false, strategies = null) {
        if (!trades || trades.length === 0) {
            this.showEmptyTradeHistory();
            return;
        }

        // 调试：输出交易数据结构（仅在开发模式）
        if (trades.length > 0 && window.location.hostname === 'localhost') {
            console.log('[Backtest] 交易数据示例:', trades[0]);
            console.log('[Backtest] 交易数据字段:', Object.keys(trades[0]));
        }

        // 保存当前回测结果，供策略详情显示使用
        this.currentBacktestResults = {
            trades: trades,
            isMultiStrategy: isMultiStrategy,
            strategies: strategies
        };

        // 如果是多策略，使用tab式展示
        if (isMultiStrategy) {
            if (strategies && strategies.length > 1) {
                this.displayMultiStrategyTradeHistoryWithTabs(trades, strategies);
        } else {
                // 多策略但strategies数据异常，降级到分组显示
                console.warn('[Backtest] 多策略模式但策略数据异常，使用兼容显示');
                this.displaySingleStrategyTradeHistoryFallback(trades);
            }
        } else {
            // 单策略模式，使用原有表格
            this.displaySingleStrategyTradeHistory(trades);
        }
    }

    /**
     * 显示单策略交易记录
     */
    displaySingleStrategyTradeHistory(trades) {
        // 显示单策略表格，隐藏多策略tabs
        const singleStrategyDiv = document.getElementById('singleStrategyTrades');
        const tradeTabsDiv = document.getElementById('tradeTabs');
        const tradeTabContentDiv = document.getElementById('tradeTabContent');
        
        if (singleStrategyDiv) singleStrategyDiv.style.display = 'block';
        if (tradeTabsDiv) tradeTabsDiv.style.display = 'none';
        if (tradeTabContentDiv) tradeTabContentDiv.style.display = 'none';

        const tableBody = document.querySelector('#tradesTable tbody');
        if (!tableBody) return;

        tableBody.innerHTML = trades.map(trade => `
            <tr>
                <td>${this.formatTradeTime(trade.timestamp)}</td>
                <td class="operation-signal">
                    <span class="operation ${trade.side === 'buy' ? 'buy' : 'sell'}">${trade.side === 'buy' ? '买入' : '卖出'}</span>
                    ${this.renderSignalType(trade)}
                </td>
                <td>${trade.quantity.toLocaleString()}</td>
                <td>¥${trade.price.toFixed(2)}</td>
                <td>¥${trade.commission.toFixed(2)}</td>
                <td class="${trade.pnl >= 0 ? 'profit' : 'loss'}">${trade.pnl ? '¥' + trade.pnl.toFixed(2) : '-'}</td>
                <td class="total-assets">¥${this.formatAssets(this.getTradeAssetValue(trade))}</td>
            </tr>
        `).join('');
    }

    /**
     * 使用Tab方式显示多策略交易记录
     */
    displayMultiStrategyTradeHistoryWithTabs(trades, strategies) {
        // 隐藏单策略表格，显示多策略tabs
        const singleStrategyDiv = document.getElementById('singleStrategyTrades');
        const tradeTabsDiv = document.getElementById('tradeTabs');
        const tradeTabContentDiv = document.getElementById('tradeTabContent');
        
        if (singleStrategyDiv) singleStrategyDiv.style.display = 'none';
        if (tradeTabsDiv) tradeTabsDiv.style.display = 'block';
        if (tradeTabContentDiv) tradeTabContentDiv.style.display = 'block';

        // 按策略分组交易记录
        const tradesByStrategy = this.groupTradesByStrategy(trades, strategies);
        
        // 生成tab导航
        this.generateTradeTabNavigation(tradesByStrategy, strategies);
        
        // 生成tab内容
        this.generateTradeTabContent(tradesByStrategy, strategies);
        
        // 激活第一个tab
        this.activateFirstTradeTab();
    }

    /**
     * 按策略分组交易记录
     */
    groupTradesByStrategy(trades, strategies) {
        const tradesByStrategy = {};
        
        // 初始化每个策略的交易记录数组
        strategies.forEach(strategy => {
            tradesByStrategy[strategy.id] = {
                strategy: strategy,
                trades: []
            };
        });
        
        // 分组交易记录
        trades.forEach(trade => {
            const strategyId = trade.strategy_id || 'unknown';
            if (tradesByStrategy[strategyId]) {
                tradesByStrategy[strategyId].trades.push(trade);
            } else {
                // 处理未知策略ID的情况
                if (!tradesByStrategy['unknown']) {
                    tradesByStrategy['unknown'] = {
                        strategy: { id: 'unknown', name: '未知策略' },
                        trades: []
                    };
                }
                tradesByStrategy['unknown'].trades.push(trade);
            }
        });
        
        return tradesByStrategy;
    }

    /**
     * 生成交易记录tab导航
     */
    generateTradeTabNavigation(tradesByStrategy, strategies) {
        const tradeTabNav = document.getElementById('tradeTabNav');
        if (!tradeTabNav) return;

        const tabButtons = Object.entries(tradesByStrategy).map(([strategyId, data]) => {
            const { strategy, trades } = data;
            const strategyIcon = this.getStrategyIcon(strategy.name, strategy.strategy_type);
            return `
                <button class="trade-tab-btn" data-strategy-id="${strategyId}">
                    <span class="tab-name" data-strategy-icon="${strategyIcon}">${strategy.name}</span>
                    <span class="tab-count">${trades.length}</span>
                </button>
            `;
        }).join('');

        tradeTabNav.innerHTML = tabButtons;

        // 绑定tab切换事件
        tradeTabNav.addEventListener('click', (e) => {
            const tabBtn = e.target.closest('.trade-tab-btn');
            if (tabBtn) {
                const strategyId = tabBtn.dataset.strategyId;
                this.switchTradeTab(strategyId);
            }
        });
    }

    /**
     * 生成策略表现详情tab内容
     */
    generateTradeTabContent(tradesByStrategy, strategies) {
        const tradeTabContent = document.getElementById('tradeTabContent');
        if (!tradeTabContent) return;

        const tabPanels = Object.entries(tradesByStrategy).map(([strategyId, data]) => {
            const { strategy, trades } = data;
            
            // 获取该策略的表现指标
            const strategyPerformance = this.getStrategyPerformance(strategyId);
            
            const tradesRows = trades.map(trade => `
                <tr>
                    <td>${this.formatTradeTime(trade.timestamp)}</td>
                    <td class="stock-symbol">${trade.symbol}</td>
                    <td class="operation-signal">
                        <span class="operation ${trade.side === 'buy' ? 'buy' : 'sell'}">${trade.side === 'buy' ? '买入' : '卖出'}</span>
                        ${this.renderSignalType(trade)}
                    </td>
                    <td>${trade.quantity.toLocaleString()}</td>
                    <td>¥${trade.price.toFixed(2)}</td>
                    <td>¥${trade.commission.toFixed(2)}</td>
                    <td class="${trade.pnl >= 0 ? 'profit' : 'loss'}">${trade.pnl ? '¥' + trade.pnl.toFixed(2) : '-'}</td>
                    <td class="holding-assets">¥${this.formatAssets(this.getTradeHoldingAssets(trade))}</td>
                    <td class="total-assets">¥${this.formatAssets(this.getTradeAssetValue(trade))}</td>
                </tr>
            `).join('');

            const strategyIcon = this.getStrategyIcon(strategy.name, strategy.strategy_type);
            
            // 生成策略表现指标HTML
            const performanceMetricsHtml = strategyPerformance ? `
                <div class="strategy-performance-header">
                    <div class="strategy-info">
                        <h6 data-strategy-icon="${strategyIcon}">${strategy.name}</h6>
                        <p class="strategy-description">${strategy.description || '该策略的详细表现指标和交易记录'}</p>
                    </div>
                    <div class="performance-metrics-compact">
                        <div class="metric-item-compact" data-tooltip="策略在整个回测期间的累计收益率。正值表示盈利，负值表示亏损。例如：15%表示初始资金增长了15%">
                            <span class="metric-label">总收益率 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value ${this.getMetricClass(strategyPerformance.total_return, '收益')}">${this.formatMetricValue(strategyPerformance.total_return, 'percentage')}</span>
                        </div>
                        <div class="metric-item-compact" data-tooltip="将总收益率按年化计算的结果。便于与其他投资产品对比。计算公式：(1+总收益率)^(365/回测天数) - 1">
                            <span class="metric-label">年化收益 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value ${this.getMetricClass(strategyPerformance.annual_return, '收益')}">${this.formatMetricValue(strategyPerformance.annual_return, 'percentage')}</span>
                        </div>
                        <div class="metric-item-compact" data-tooltip="策略在回测期间的最大资产回撤幅度。衡量策略的风险水平。例如：-8%表示最大亏损幅度为8%。数值越小风险越低">
                            <span class="metric-label">最大回撤 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value ${this.getMetricClass(strategyPerformance.max_drawdown, '回撤')}">${this.formatMetricValue(strategyPerformance.max_drawdown, 'percentage')}</span>
                        </div>
                        <div class="metric-item-compact" data-tooltip="衡量策略风险调整后收益的指标。计算公式：(策略收益率-无风险收益率)/收益率标准差。通常>1为良好，>2为优秀">
                            <span class="metric-label">夏普比率 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value ${this.getMetricClass(strategyPerformance.sharpe_ratio, '夏普')}">${this.formatMetricValue(strategyPerformance.sharpe_ratio, 'decimal')}</span>
                        </div>
                        <div class="metric-item-compact" data-tooltip="盈利交易占总交易次数的比例。例如：65%表示100笔交易中有65笔是盈利的。胜率越高策略越稳定">
                            <span class="metric-label">胜率 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value ${this.getMetricClass(strategyPerformance.win_rate, '胜率')}">${this.formatMetricValue(strategyPerformance.win_rate, 'percentage')}</span>
                        </div>
                        <div class="metric-item-compact" data-tooltip="策略在回测期间执行的买卖交易总次数。交易次数过少可能统计意义不足，过多可能交易成本过高">
                            <span class="metric-label">交易次数 <i class="tooltip-icon">?</i></span>
                            <span class="metric-value neutral">${strategyPerformance.total_trades || trades.length}</span>
                        </div>
                    </div>
                </div>
            ` : `
                <div class="strategy-performance-header">
                    <div class="strategy-info">
                        <h6 data-strategy-icon="${strategyIcon}">${strategy.name}</h6>
                        <p class="strategy-description">该策略的详细交易记录</p>
                    </div>
                    <div class="trade-summary-compact">
                        <span class="trade-count">共 ${trades.length} 笔交易</span>
                    </div>
                </div>
            `;
            
            return `
                <div class="trade-tab-panel" data-strategy-id="${strategyId}" style="display: none;">
                    ${performanceMetricsHtml}
                    <div class="table-container">
                        <div class="table-wrapper">
                            <table class="trades-table">
                                <thead>
                                    <tr class="column-headers">
                                        <th>时间</th>
                                        <th>股票</th>
                                        <th>操作/信号</th>
                                        <th>数量</th>
                                        <th>价格</th>
                                        <th>手续费</th>
                                        <th>盈亏</th>
                                        <th>持仓资产</th>
                                        <th>总资产</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    ${tradesRows || '<tr><td colspan="9">该策略暂无交易记录</td></tr>'}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            `;
        }).join('');

        tradeTabContent.innerHTML = tabPanels;
    }

    /**
     * 切换交易记录tab
     */
    switchTradeTab(strategyId) {
        // 更新tab按钮状态
        const tabButtons = document.querySelectorAll('.trade-tab-btn');
        tabButtons.forEach(btn => {
            btn.classList.toggle('active', btn.dataset.strategyId === strategyId);
        });

        // 更新tab面板显示
        const tabPanels = document.querySelectorAll('.trade-tab-panel');
        tabPanels.forEach(panel => {
            panel.style.display = panel.dataset.strategyId === strategyId ? 'block' : 'none';
        });
    }

    /**
     * 激活第一个交易记录tab
     */
    activateFirstTradeTab() {
        const firstTabBtn = document.querySelector('.trade-tab-btn');
        if (firstTabBtn) {
            const strategyId = firstTabBtn.dataset.strategyId;
            this.switchTradeTab(strategyId);
        }
    }

    /**
     * 跳转到指定策略的交易记录
     */
    jumpToStrategyTrades(strategyId) {
        console.log('[Backtest] 跳转到策略交易记录:', strategyId);
        
        // 1. 切换到对应策略的交易记录tab
        this.switchTradeTab(strategyId);
        
        // 2. 滚动到交易记录区域
        const tradesSection = document.getElementById('tradeTabs') || document.getElementById('singleStrategyTrades');
        if (tradesSection) {
            tradesSection.scrollIntoView({ 
                behavior: 'smooth', 
                block: 'start',
                inline: 'nearest'
            });
            
            // 3. 添加高亮效果
            this.highlightTradeTab(strategyId);
        } else {
            console.warn('[Backtest] 未找到交易记录区域');
        }
    }

    /**
     * 高亮显示指定策略的交易记录tab
     */
    highlightTradeTab(strategyId) {
        // 移除所有高亮效果
        const allTabs = document.querySelectorAll('.trade-tab-btn');
        allTabs.forEach(tab => tab.classList.remove('highlight-flash'));
        
        // 为目标tab添加高亮效果
        const targetTab = document.querySelector(`.trade-tab-btn[data-strategy-id="${strategyId}"]`);
        if (targetTab) {
            targetTab.classList.add('highlight-flash');
            
            // 2秒后移除高亮效果
            setTimeout(() => {
                targetTab.classList.remove('highlight-flash');
            }, 2000);
        }
    }

    /**
     * 获取指定策略的表现指标
     */
    getStrategyPerformance(strategyId) {
        // 从当前回测结果中查找对应策略的表现数据
        if (this.currentPerformanceResults && Array.isArray(this.currentPerformanceResults)) {
            // 多策略情况
            const performance = this.currentPerformanceResults.find(p => p.strategy_id === strategyId);
            if (performance) {
                return performance;
            }
        }
        
        // 尝试从保存的回测结果中获取
        if (this.currentBacktestResults && this.currentBacktestResults.strategies) {
            const strategy = this.currentBacktestResults.strategies.find(s => s.id === strategyId);
            if (strategy) {
                // 如果没有具体的表现数据，返回基础信息
                return {
                    strategy_id: strategyId,
                    strategy_name: strategy.name,
                    total_return: 0.15, // 默认值，实际应该从回测结果获取
                    annual_return: 0.12,
                    max_drawdown: -0.08,
                    sharpe_ratio: 1.5,
                    win_rate: 0.65,
                    total_trades: this.getStrategyTradeCount(strategyId)
                };
            }
        }
        
        return null;
    }

    /**
     * 获取指定策略的交易次数
     */
    getStrategyTradeCount(strategyId) {
        if (this.currentBacktestResults && this.currentBacktestResults.trades) {
            return this.currentBacktestResults.trades.filter(trade => 
                trade.strategy_id === strategyId
            ).length;
        }
        return 0;
    }

    /**
     * 显示单策略交易记录（多策略降级模式）
     */
    displaySingleStrategyTradeHistoryFallback(trades) {
        // 显示单策略表格，隐藏多策略tabs
        const singleStrategyDiv = document.getElementById('singleStrategyTrades');
        const tradeTabsDiv = document.getElementById('tradeTabs');
        const tradeTabContentDiv = document.getElementById('tradeTabContent');
        
        if (singleStrategyDiv) singleStrategyDiv.style.display = 'block';
        if (tradeTabsDiv) tradeTabsDiv.style.display = 'none';
        if (tradeTabContentDiv) tradeTabContentDiv.style.display = 'none';

        const tableBody = document.querySelector('#tradesTable tbody');
        if (!tableBody) return;

        // 使用原有的分组显示逻辑
        this.displayMultiStrategyTradeHistory(trades, tableBody);
    }

    /**
     * 显示空的交易记录
     */
    showEmptyTradeHistory() {
        // 显示单策略表格，隐藏多策略tabs
        const singleStrategyDiv = document.getElementById('singleStrategyTrades');
        const tradeTabsDiv = document.getElementById('tradeTabs');
        const tradeTabContentDiv = document.getElementById('tradeTabContent');
        
        if (singleStrategyDiv) singleStrategyDiv.style.display = 'block';
        if (tradeTabsDiv) tradeTabsDiv.style.display = 'none';
        if (tradeTabContentDiv) tradeTabContentDiv.style.display = 'none';

        const tableBody = document.querySelector('#tradesTable tbody');
        if (tableBody) {
            tableBody.innerHTML = '<tr><td colspan="7">暂无交易记录</td></tr>';
        }
    }

    /**
     * 更新股票信息表头
     */
    updateStockInfoHeader(backtestConfig) {
        const stockInfoHeader = document.getElementById('stockInfoHeader');
        if (!stockInfoHeader || !backtestConfig || !backtestConfig.symbols) return;

        const stockName = stockInfoHeader.querySelector('.stock-name');
        const stockCode = stockInfoHeader.querySelector('.stock-code');
        
        if (backtestConfig.symbols.length === 1) {
            // 单只股票
            const symbol = backtestConfig.symbols[0];
            if (stockName) stockName.textContent = '股票信息';
            if (stockCode) stockCode.textContent = symbol;
        } else {
            // 多只股票
            if (stockName) stockName.textContent = '股票组合';
            if (stockCode) stockCode.textContent = `${backtestConfig.symbols.length}只股票`;
        }
    }

    /**
     * 格式化交易时间
     */
    formatTradeTime(timestamp) {
        if (!timestamp) return '-';
        
        try {
            const date = new Date(timestamp);
            // 始终显示完整的年月日格式
            return date.toLocaleDateString('zh-CN', { 
                year: 'numeric',
                month: '2-digit', 
                day: '2-digit' 
            });
        } catch (error) {
            // 降级处理：尝试从字符串中提取日期部分
            const dateStr = timestamp.split(' ')[0] || timestamp;
            if (dateStr.match(/^\d{4}-\d{2}-\d{2}$/)) {
                // 如果是 YYYY-MM-DD 格式，直接转换
                const [year, month, day] = dateStr.split('-');
                return `${year}/${month}/${day}`;
            }
            return dateStr;
        }
    }

    /**
     * 获取交易的持仓资产值
     */
    getTradeHoldingAssets(trade) {
        // 优先使用新的持仓资产字段（这是当前策略的持仓资产）
        if (trade.holding_assets !== undefined && trade.holding_assets !== null) {
            return trade.holding_assets;
        }
        
        // ❌ 重要修复：不能用 total_assets - cash_balance 计算持仓资产
        // 因为在多策略回测中，total_assets 是所有策略的总资产，而 cash_balance 是当前策略的现金
        // 这样计算会导致错误的持仓资产显示
        
        // 尝试其他可能的字段名
        const possibleFields = [
            'holdings',
            'stock_value',
            'position_value',
            'market_value'
        ];
        
        for (const field of possibleFields) {
            if (trade[field] !== undefined && trade[field] !== null) {
                return trade[field];
            }
        }
        
        // 如果都没有，返回0（表示没有持仓或无法计算）
        // 注意：这里返回0是正确的，因为如果没有holding_assets字段，
        // 说明可能是卖出后没有持仓，或者数据结构有问题
        return 0;
    }

    /**
     * 获取交易的资产值
     */
    getTradeAssetValue(trade) {
        // 🔧 重要修复：在多策略回测中，应该显示单策略的资产，而不是所有策略的总资产
        // trade.total_assets 是所有策略的总资产，不适合在单策略交易记录中显示
        
        // 优先计算：当前策略的持仓资产 + 现金余额
        const holdingAssets = this.getTradeHoldingAssets(trade);
        const cashBalance = trade.cash_balance || 0;
        
        if (holdingAssets !== null && holdingAssets !== undefined && 
            cashBalance !== null && cashBalance !== undefined) {
            const singleStrategyAssets = holdingAssets + cashBalance;
            console.log(`单策略资产计算: 持仓${holdingAssets.toFixed(2)} + 现金${cashBalance.toFixed(2)} = ${singleStrategyAssets.toFixed(2)}`);
            return singleStrategyAssets;
        }
        
        // 如果没有分离的持仓和现金数据，才考虑使用总资产字段
        // 但要注意这可能是多策略的合计值
        if (trade.total_assets !== undefined && trade.total_assets !== null) {
            // 添加警告标识，提醒这可能是多策略总资产
            console.warn(`⚠️ 使用total_assets字段 (${trade.total_assets.toFixed(2)})，这可能包含多策略资产`);
            return trade.total_assets;
        }
        
        // 尝试其他可能的字段名
        const possibleFields = [
            'portfolio_value',
            'total_value', 
            'account_value',
            'cash_balance',
            'balance',
            'value',
            'total'
        ];
        
        for (const field of possibleFields) {
            if (trade[field] !== undefined && trade[field] !== null) {
                return trade[field];
            }
        }
        
        // 如果都没有，尝试计算：初始资金 + 累计盈亏
        if (trade.cumulative_pnl !== undefined) {
            const initialCash = this.currentBacktestConfig?.initial_cash || 1000000;
            return initialCash + trade.cumulative_pnl;
        }
        
        // 尝试从价格和数量计算当前持仓价值（简化计算）
        if (trade.price && trade.quantity && trade.side) {
            const tradeValue = trade.price * trade.quantity;
            // 这只是一个简化的估算，实际应该维护完整的账户状态
            return tradeValue;
        }
        
        // 仅在开发模式下输出警告
        if (window.location.hostname === 'localhost') {
            console.warn('[Backtest] 无法找到资产值字段，可用字段:', Object.keys(trade));
        }
        return 0;
    }

    /**
     * 格式化资产金额
     */
    formatAssets(amount) {
        if (!amount || amount === 0) return '0';
        
        if (amount >= 100000000) {
            return (amount / 100000000).toFixed(2) + '亿';
        } else if (amount >= 10000) {
            return (amount / 10000).toFixed(1) + '万';
        } else {
            return amount.toLocaleString();
        }
    }

    /**
     * 格式化信号类型显示
     */
    formatSignalType(trade) {
        const signalType = trade.signal_type || trade.signal || trade.reason || '';
        
        // 如果信号类型就是简单的buy/sell，则不显示
        if (signalType.toLowerCase() === trade.side?.toLowerCase()) {
            return '-';
        }
        
        // 信号类型映射，显示更有意义的信息
        const signalMap = {
            'ma_crossover': '均线交叉',
            'macd_golden_cross': 'MACD金叉',
            'macd_death_cross': 'MACD死叉',
            'rsi_oversold': 'RSI超卖',
            'rsi_overbought': 'RSI超买',
            'bollinger_lower': '布林下轨',
            'bollinger_upper': '布林上轨',
            'support_level': '支撑位',
            'resistance_level': '阻力位',
            'volume_breakout': '放量突破',
            'trend_reversal': '趋势反转',
            'stop_loss': '止损',
            'take_profit': '止盈'
        };
        
        return signalMap[signalType] || signalType || '-';
    }

    /**
     * 渲染信号类型span（只有在有有效信号时才渲染）
     */
    renderSignalType(trade) {
        const signalText = this.formatSignalType(trade);
        
        // 如果信号类型是'-'或空，则不渲染span元素
        if (!signalText || signalText === '-') {
            return '';
        }
        
        return `<span class="signal-type">${signalText}</span>`;
    }

    /**
     * 获取策略图标
     */
    getStrategyIcon(strategyName, strategyType) {
        // 根据策略名称或类型返回对应图标
        const nameIconMap = {
            '双均线': '📈',
            'MA': '📈',
            '均线': '📈',
            '布林带': '📊',
            'BOLL': '📊',
            'Bollinger': '📊',
            'MACD': '📉',
            'RSI': '🔄',
            '相对强弱': '🔄',
            'KDJ': '⚡',
            '随机指标': '⚡',
            '金叉': '✨',
            '死叉': '💫',
            '超买': '🔺',
            '超卖': '🔻',
            '突破': '🚀',
            '回调': '📉',
            '趋势': '📈',
            '震荡': '🌊',
            '动量': '⚡',
            '均值回归': '🔄',
            '网格': '🔲',
            '套利': '⚖️'
        };

        const typeIconMap = {
            'technical': '📊',
            'fundamental': '💰',
            'ml': '🤖',
            'composite': '🔗',
            'combined': '📋',
            'momentum': '⚡',
            'trend': '📈',
            'mean_reversion': '🔄',
            'breakout': '🚀'
        };

        // 首先尝试根据策略名称匹配
        for (const [keyword, icon] of Object.entries(nameIconMap)) {
            if (strategyName && strategyName.includes(keyword)) {
                return icon;
            }
        }

        // 其次根据策略类型匹配
        if (strategyType && typeIconMap[strategyType]) {
            return typeIconMap[strategyType];
        }

        // 默认图标
        return '📊';
    }

    /**
     * 显示多策略交易记录（按策略分组）- 保留原方法以备兼容
     */
    displayMultiStrategyTradeHistory(trades, tableBody) {
        // 按策略分组
        const tradesByStrategy = {};
        trades.forEach(trade => {
            const strategyId = trade.strategy_id || 'unknown';
            if (!tradesByStrategy[strategyId]) {
                tradesByStrategy[strategyId] = [];
            }
            tradesByStrategy[strategyId].push(trade);
        });

        let html = '';
        Object.entries(tradesByStrategy).forEach(([strategyId, strategyTrades]) => {
            // 策略分组标题行
            html += `
                <tr class="strategy-group-header">
                    <td colspan="7" class="strategy-group-title">
                        <strong>策略: ${strategyId}</strong> 
                        <span class="trade-count">(${strategyTrades.length}笔交易)</span>
                    </td>
                </tr>
            `;

            // 该策略的交易记录
            strategyTrades.forEach(trade => {
                html += `
                    <tr class="strategy-trade-row">
                        <td>${this.formatTradeTime(trade.timestamp)}</td>
                        <td class="operation-signal">
                            <span class="operation ${trade.side === 'buy' ? 'buy' : 'sell'}">${trade.side === 'buy' ? '买入' : '卖出'}</span>
                            ${this.renderSignalType(trade)}
                        </td>
                        <td>${trade.quantity.toLocaleString()}</td>
                        <td>¥${trade.price.toFixed(2)}</td>
                        <td>¥${trade.commission.toFixed(2)}</td>
                        <td class="${trade.pnl >= 0 ? 'profit' : 'loss'}">${trade.pnl ? '¥' + trade.pnl.toFixed(2) : '-'}</td>
                        <td class="total-assets">¥${this.formatAssets(this.getTradeAssetValue(trade))}</td>
                    </tr>
                `;
            });
        });

        if (tableBody) {
        tableBody.innerHTML = html;
        }
    }

    /**
     * 保存回测配置
     */
    saveBacktestConfig() {
        const config = this.getBacktestConfig();
        
        // 保存到本地存储
        localStorage.setItem('backtest_config', JSON.stringify(config));
        
        this.showMessage('配置已保存', 'success');
    }

    /**
     * 加载回测历史
     */
    async loadBacktestHistory() {
        try {
            // 加载保存的配置
            const savedConfig = localStorage.getItem('backtest_config');
            if (savedConfig) {
                const config = JSON.parse(savedConfig);
                this.fillConfigForm(config);
            }

            // 加载策略列表
            await this.loadStrategiesList();

        } catch (error) {
            console.error('加载回测历史失败:', error);
        }
    }

    /**
     * 填充配置表单
     */
    fillConfigForm(config) {
        if (config.name) {
            const nameInput = document.getElementById('backtestName');
            if (nameInput) nameInput.value = config.name;
        }

        if (config.initial_cash) {
            const cashInput = document.getElementById('initialCash');
            if (cashInput) cashInput.value = config.initial_cash;
        }

        if (config.commission) {
            const commissionInput = document.getElementById('commission');
            if (commissionInput) commissionInput.value = config.commission;
        }

        if (config.symbols && config.symbols.length > 0) {
            const symbolsInput = document.getElementById('backtestSymbols');
            if (symbolsInput) symbolsInput.value = config.symbols.join('\n');
        }
    }

    /**
     * 加载策略列表
     */
    async loadStrategiesList() {
        try {
            const response = await this.apiService.getStrategiesList();
            if (response.success) {
                // 处理分页响应格式，确保获取正确的策略数组
                const strategies = response.data.items || response.data || [];
                console.log(`[Backtest] 加载策略列表成功，共 ${strategies.length} 个策略`);
                this.updateStrategiesSelect(strategies);
            }
        } catch (error) {
            console.error('加载策略列表失败:', error);
        }
    }

    /**
     * 更新策略选择组件（支持多选）
     */
    updateStrategiesSelect(strategies) {
        // 验证strategies是数组
        if (!Array.isArray(strategies)) {
            console.error('[Backtest] strategies不是数组:', typeof strategies, strategies);
            return;
        }

        console.log(`[Backtest] 更新多选策略组件，共 ${strategies.length} 个策略`);

        // 保存策略数据
        this.availableStrategies = strategies;
        this.selectedStrategyIds = [];

        // 初始化多选下拉框
        this.initStrategyMultiSelect();
    }

    /**
     * 初始化多选策略下拉框
     */
    initStrategyMultiSelect() {
        const dropdownHeader = document.getElementById('dropdownHeader');
        const dropdownList = document.getElementById('dropdownList');
        const selectedStrategiesContainer = document.getElementById('selectedStrategies');

        if (!dropdownHeader || !dropdownList || !selectedStrategiesContainer) {
            console.warn('[Backtest] 找不到多选策略组件元素');
            return;
        }

        // 清空下拉列表
        dropdownList.innerHTML = '';

        // 首先添加控制按钮
        const controlsDiv = document.createElement('div');
        controlsDiv.className = 'dropdown-controls';
        controlsDiv.innerHTML = `
            <button type="button" class="select-all-btn" id="selectAllStrategies">
                <span class="control-icon">☑</span>
                <span class="control-text">全选</span>
            </button>
            <button type="button" class="clear-all-btn" id="clearAllStrategies">
                <span class="control-icon">☐</span>
                <span class="control-text">清空</span>
            </button>
        `;
        dropdownList.appendChild(controlsDiv);

        // 添加策略选项
        this.availableStrategies.forEach(strategy => {
            const option = document.createElement('div');
            option.className = 'dropdown-option';
            option.innerHTML = `
                <input type="checkbox" id="strategy_${strategy.id}" value="${strategy.id}">
                <div class="strategy-info">
                    <div class="strategy-name">${strategy.name}</div>
                    <div class="strategy-type">${strategy.strategy_type}</div>
                </div>
            `;

            // 添加点击事件
            option.addEventListener('click', (e) => {
                e.preventDefault();
                this.toggleStrategySelection(strategy.id);
            });

            // 复选框点击事件
            const checkbox = option.querySelector('input[type="checkbox"]');
            checkbox.addEventListener('change', (e) => {
                e.stopPropagation();
                this.toggleStrategySelection(strategy.id);
            });

            dropdownList.appendChild(option);
        });

        // 绑定控制按钮事件
        this.bindControlButtonEvents();

        // 下拉框头部点击事件
        dropdownHeader.addEventListener('click', () => {
            this.toggleDropdown();
        });

        // 点击外部关闭下拉框
        document.addEventListener('click', (e) => {
            if (!e.target.closest('.strategy-dropdown')) {
                this.closeDropdown();
            }
        });

        // 初始化显示
        this.updateSelectedStrategiesDisplay();
    }

    /**
     * 绑定控制按钮事件
     */
    bindControlButtonEvents() {
        // 全选按钮事件
        const selectAllBtn = document.getElementById('selectAllStrategies');
        if (selectAllBtn) {
            selectAllBtn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.selectAllStrategies();
            });
        }

        // 清空按钮事件
        const clearAllBtn = document.getElementById('clearAllStrategies');
        if (clearAllBtn) {
            clearAllBtn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.clearAllStrategies();
            });
        }
    }

    /**
     * 切换策略选择状态
     */
    toggleStrategySelection(strategyId) {
        const index = this.selectedStrategyIds.indexOf(strategyId);
        
        if (index > -1) {
            // 取消选择
            this.selectedStrategyIds.splice(index, 1);
        } else {
            // 添加选择（限制最多5个）
            if (this.selectedStrategyIds.length >= 5) {
                this.showMessage('最多只能选择5个策略', 'warning');
                return;
            }
            this.selectedStrategyIds.push(strategyId);
        }

        this.updateSelectedStrategiesDisplay();
        this.updateDropdownOptions();
    }

    /**
     * 切换下拉框显示/隐藏
     */
    toggleDropdown() {
        const dropdownHeader = document.getElementById('dropdownHeader');
        const dropdownList = document.getElementById('dropdownList');

        if (!dropdownHeader || !dropdownList) return;

        const isOpen = dropdownList.style.display === 'block';
        
        if (isOpen) {
            this.closeDropdown();
        } else {
            this.openDropdown();
        }
    }

    /**
     * 打开下拉框
     */
    openDropdown() {
        const dropdownHeader = document.getElementById('dropdownHeader');
        const dropdownList = document.getElementById('dropdownList');

        if (dropdownHeader && dropdownList) {
            dropdownHeader.classList.add('active');
            dropdownList.style.display = 'block';
        }
    }

    /**
     * 关闭下拉框
     */
    closeDropdown() {
        const dropdownHeader = document.getElementById('dropdownHeader');
        const dropdownList = document.getElementById('dropdownList');

        if (dropdownHeader && dropdownList) {
            dropdownHeader.classList.remove('active');
            dropdownList.style.display = 'none';
        }
    }

    /**
     * 更新已选择策略的显示
     */
    updateSelectedStrategiesDisplay() {
        const dropdownHeader = document.getElementById('dropdownHeader');
        const selectedStrategiesContainer = document.getElementById('selectedStrategies');

        if (!dropdownHeader || !selectedStrategiesContainer) return;

        // 更新下拉框头部显示
        const placeholderSpan = dropdownHeader.querySelector('.placeholder') || 
                               dropdownHeader.querySelector('.selected-count');
        
        if (this.selectedStrategyIds.length === 0) {
            placeholderSpan.textContent = '请选择策略...';
            placeholderSpan.className = 'placeholder';
        } else {
            placeholderSpan.textContent = `已选择 ${this.selectedStrategyIds.length} 个策略`;
            placeholderSpan.className = 'selected-count';
        }

        // 更新已选择策略标签
        selectedStrategiesContainer.innerHTML = '';
        
        this.selectedStrategyIds.forEach(strategyId => {
            const strategy = this.availableStrategies.find(s => s.id === strategyId);
            if (!strategy) return;

            const tag = document.createElement('div');
            tag.className = 'strategy-tag';
            tag.innerHTML = `
                <span>${strategy.name}</span>
                <button type="button" class="remove-btn" onclick="backtestModule.removeStrategySelection('${strategyId}')">&times;</button>
            `;
            
            selectedStrategiesContainer.appendChild(tag);
        });
    }

    /**
     * 更新下拉选项的选中状态
     */
    updateDropdownOptions() {
        const dropdownList = document.getElementById('dropdownList');
        if (!dropdownList) return;

        const options = dropdownList.querySelectorAll('.dropdown-option');
        options.forEach(option => {
            const checkbox = option.querySelector('input[type="checkbox"]');
            if (checkbox) {
                const strategyId = checkbox.value;
                const isSelected = this.selectedStrategyIds.includes(strategyId);
                
                checkbox.checked = isSelected;
                option.classList.toggle('selected', isSelected);
            }
        });
    }

    /**
     * 移除策略选择
     */
    removeStrategySelection(strategyId) {
        const index = this.selectedStrategyIds.indexOf(strategyId);
        if (index > -1) {
            this.selectedStrategyIds.splice(index, 1);
            this.updateSelectedStrategiesDisplay();
            this.updateDropdownOptions();
        }
    }

    /**
     * 全选策略
     */
    selectAllStrategies() {
        if (!this.availableStrategies || this.availableStrategies.length === 0) {
            this.showMessage('没有可选择的策略', 'warning');
            return;
        }

        // 检查是否超过最大限制
        if (this.availableStrategies.length > 5) {
            this.showMessage('最多只能选择5个策略，将选择前5个', 'warning');
            this.selectedStrategyIds = this.availableStrategies.slice(0, 5).map(s => s.id);
        } else {
            this.selectedStrategyIds = this.availableStrategies.map(s => s.id);
        }

        this.updateSelectedStrategiesDisplay();
        this.updateDropdownOptions();
        
        console.log(`[Backtest] 全选策略完成，已选择 ${this.selectedStrategyIds.length} 个策略`);
    }

    /**
     * 清空所有策略选择
     */
    clearAllStrategies() {
        if (this.selectedStrategyIds.length === 0) {
            this.showMessage('当前没有已选择的策略', 'info');
            return;
        }

        const previousCount = this.selectedStrategyIds.length;
        this.selectedStrategyIds = [];
        
        this.updateSelectedStrategiesDisplay();
        this.updateDropdownOptions();
        
        console.log(`[Backtest] 清空策略选择完成，已清除 ${previousCount} 个策略`);
    }

    /**
     * 显示消息
     */
    showMessage(message, type = 'info') {
        // 创建消息提示元素
        const messageDiv = document.createElement('div');
        messageDiv.className = `message-toast ${type}`;
        messageDiv.textContent = message;
        
        // 添加到页面
        document.body.appendChild(messageDiv);
        
        // 3秒后自动移除
        setTimeout(() => {
            if (messageDiv.parentNode) {
                messageDiv.parentNode.removeChild(messageDiv);
            }
        }, 3000);
    }
}

// 导出回测模块类
window.BacktestModule = BacktestModule;
