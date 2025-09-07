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
        
        this.init();
    }

    /**
     * 初始化回测模块
     */
    init() {
        this.setupEventListeners();
        this.loadBacktestHistory();
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
        const backtestStrategy = document.getElementById('backtestStrategy')?.value || '';
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
            strategy_id: backtestStrategy,
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

        if (!config.strategy_id) {
            this.showMessage('请选择策略', 'warning');
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

        // 显示性能指标
        this.displayPerformanceMetrics(results.performance);

        // 显示权益曲线
        this.displayEquityCurve(results.equity_curve);

        // 显示交易记录
        this.displayTradeHistory(results.trades);
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
            <div class="metric-card">
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
     * 更新策略选择框
     */
    updateStrategiesSelect(strategies) {
        const select = document.getElementById('backtestStrategy');
        if (!select) {
            console.warn('[Backtest] 找不到策略选择框元素');
            return;
        }

        // 验证strategies是数组
        if (!Array.isArray(strategies)) {
            console.error('[Backtest] strategies不是数组:', typeof strategies, strategies);
            return;
        }

        console.log(`[Backtest] 更新策略选择框，共 ${strategies.length} 个策略`);

        // 清空现有选项（保留默认选项）
        const defaultOption = select.querySelector('option[value=""]');
        select.innerHTML = '';
        if (defaultOption) {
            select.appendChild(defaultOption);
        } else {
            // 添加默认选项
            const defaultOpt = document.createElement('option');
            defaultOpt.value = '';
            defaultOpt.textContent = '请选择策略';
            select.appendChild(defaultOpt);
        }

        // 添加策略选项
        strategies.forEach((strategy, index) => {
            try {
                const option = document.createElement('option');
                option.value = strategy.id;
                option.textContent = `${strategy.name} (${strategy.strategy_type})`;
                select.appendChild(option);
                
                console.log(`[Backtest] 添加策略选项 ${index + 1}: ${strategy.name}`);
            } catch (error) {
                console.error(`[Backtest] 添加策略选项失败:`, strategy, error);
            }
        });
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
