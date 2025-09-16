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
                        this.loadBacktestResults(this.currentBacktest.id);
                        this.showMessage('回测完成！', 'success');
                    } else if (progress.status === 'failed') {
                        this.stopProgressMonitoring();
                        this.hideProgress();
                        this.showMessage(progress.error || '回测失败', 'error');
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

        // 显示策略配置信息
        this.displayStrategyConfig(results.strategy, results.backtest_config);

        // 显示性能指标
        this.displayPerformanceMetrics(results.performance);

        // 显示权益曲线
        this.displayEquityCurve(results.equity_curve);

        // 显示交易记录
        this.displayTradeHistory(results.trades);
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
            'composite': '复合策略'
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
            'std_dev': '标准差倍数'
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
    displayPerformanceMetrics(performance) {
        const metricsGrid = document.getElementById('metricsGrid');
        if (!metricsGrid) return;

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
     * 显示交易记录
     */
    displayTradeHistory(trades) {
        const tableBody = document.querySelector('#tradesTable tbody');
        if (!tableBody || !trades || trades.length === 0) {
            if (tableBody) {
                tableBody.innerHTML = '<tr><td colspan="7">暂无交易记录</td></tr>';
            }
            return;
        }

        tableBody.innerHTML = trades.map(trade => `
            <tr>
                <td>${trade.timestamp}</td>
                <td>${trade.symbol}</td>
                <td class="${trade.side === 'buy' ? 'buy' : 'sell'}">${trade.side === 'buy' ? '买入' : '卖出'}</td>
                <td>${trade.quantity}</td>
                <td>¥${trade.price.toFixed(2)}</td>
                <td>¥${trade.commission.toFixed(2)}</td>
                <td class="${trade.pnl >= 0 ? 'profit' : 'loss'}">${trade.pnl ? '¥' + trade.pnl.toFixed(2) : '-'}</td>
            </tr>
        `).join('');
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
