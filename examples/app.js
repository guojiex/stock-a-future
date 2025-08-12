/**
 * Stock-A-Future API 网页客户端
 * 提供与后端API的交互功能
 */

class StockAFutureClient {
    constructor(baseURL = null) {
        this.baseURL = baseURL || this.detectServerURL();
        this.currentChart = null;
        this.isLoading = false;
        
        // 初始化
        this.init();
    }

    /**
     * 自动检测服务器URL和端口
     */
    detectServerURL() {
        // 1. 优先从URL参数中读取
        const urlParams = new URLSearchParams(window.location.search);
        const serverPort = urlParams.get('port');
        if (serverPort) {
            return `http://localhost:${serverPort}`;
        }

        // 2. 从localStorage中读取上次保存的配置
        const savedURL = localStorage.getItem('stockapi_server_url');
        if (savedURL) {
            return savedURL;
        }

        // 3. 尝试常见端口
        return 'http://localhost:8081'; // 默认使用8081而不是8080
    }

    /**
     * 设置服务器URL并保存到localStorage
     */
    setServerURL(url) {
        this.baseURL = url;
        localStorage.setItem('stockapi_server_url', url);
        console.log(`服务器URL已更新为: ${url}`);
    }

    /**
     * 初始化客户端
     */
    async init() {
        this.setupEventListeners();
        this.setDefaultDates();
        await this.checkHealth();
    }

    /**
     * 设置事件监听器
     */
    setupEventListeners() {
        // 查询按钮事件
        document.getElementById('queryDaily').addEventListener('click', () => this.handleDailyQuery());
        document.getElementById('queryIndicators').addEventListener('click', () => this.handleIndicatorsQuery());
        document.getElementById('queryPredictions').addEventListener('click', () => this.handlePredictionsQuery());

        // 回车键支持
        document.getElementById('stockCode').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.handleDailyQuery();
            }
        });

        // 实时验证股票代码格式
        document.getElementById('stockCode').addEventListener('input', this.validateStockCode);

        // 配置相关事件
        this.setupConfigEventListeners();
    }

    /**
     * 设置配置相关事件监听器
     */
    setupConfigEventListeners() {
        // 配置按钮
        document.getElementById('configBtn').addEventListener('click', () => this.showConfigModal());
        
        // 关闭配置模态框
        document.getElementById('closeConfigBtn').addEventListener('click', () => this.hideConfigModal());
        
        // 点击模态框背景关闭
        document.getElementById('configModal').addEventListener('click', (e) => {
            if (e.target.id === 'configModal') {
                this.hideConfigModal();
            }
        });

        // 测试连接按钮
        document.getElementById('testConnectionBtn').addEventListener('click', () => this.testConnection());
        
        // 保存配置按钮
        document.getElementById('saveConfigBtn').addEventListener('click', () => this.saveConfig());

        // 端口快捷按钮
        document.querySelectorAll('.port-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const port = e.target.dataset.port;
                document.getElementById('serverURL').value = `http://localhost:${port}`;
            });
        });

        // ESC键关闭模态框
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                this.hideConfigModal();
            }
        });
    }

    /**
     * 显示配置模态框
     */
    showConfigModal() {
        const modal = document.getElementById('configModal');
        const serverURLInput = document.getElementById('serverURL');
        
        // 设置当前服务器URL
        serverURLInput.value = this.baseURL;
        
        // 显示模态框
        modal.style.display = 'flex';
        modal.classList.add('fade-in');
        
        // 聚焦到输入框
        setTimeout(() => serverURLInput.focus(), 100);
    }

    /**
     * 隐藏配置模态框
     */
    hideConfigModal() {
        const modal = document.getElementById('configModal');
        const testResult = document.getElementById('connectionTestResult');
        
        modal.style.display = 'none';
        testResult.style.display = 'none';
    }

    /**
     * 测试连接
     */
    async testConnection() {
        const serverURL = document.getElementById('serverURL').value.trim();
        const testResult = document.getElementById('connectionTestResult');
        const testBtn = document.getElementById('testConnectionBtn');
        
        if (!serverURL) {
            this.showTestResult('请输入服务器地址', 'error');
            return;
        }

        // 验证URL格式
        try {
            new URL(serverURL);
        } catch (error) {
            this.showTestResult('无效的URL格式', 'error');
            return;
        }

        testBtn.disabled = true;
        testBtn.textContent = '🔍 测试中...';
        
        try {
            const response = await fetch(`${serverURL}/api/v1/health`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
                signal: AbortSignal.timeout(5000) // 5秒超时
            });

            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    this.showTestResult('✅ 连接成功！服务器运行正常', 'success');
                } else {
                    this.showTestResult('⚠️ 服务器响应异常', 'warning');
                }
            } else {
                this.showTestResult(`❌ 连接失败: HTTP ${response.status}`, 'error');
            }
        } catch (error) {
            if (error.name === 'AbortError') {
                this.showTestResult('❌ 连接超时，请检查服务器地址和端口', 'error');
            } else {
                this.showTestResult(`❌ 连接失败: ${error.message}`, 'error');
            }
        } finally {
            testBtn.disabled = false;
            testBtn.textContent = '🔗 测试连接';
        }
    }

    /**
     * 显示测试结果
     */
    showTestResult(message, type) {
        const testResult = document.getElementById('connectionTestResult');
        testResult.className = `test-result ${type}`;
        testResult.textContent = message;
        testResult.style.display = 'block';
    }

    /**
     * 保存配置
     */
    saveConfig() {
        const serverURL = document.getElementById('serverURL').value.trim();
        
        if (!serverURL) {
            this.showTestResult('请输入服务器地址', 'error');
            return;
        }

        // 验证URL格式
        try {
            new URL(serverURL);
        } catch (error) {
            this.showTestResult('无效的URL格式', 'error');
            return;
        }

        // 更新服务器URL
        this.setServerURL(serverURL);
        
        // 立即进行健康检查
        this.checkHealth();
        
        // 关闭模态框
        this.hideConfigModal();
        
        console.log(`配置已保存: ${serverURL}`);
    }

    /**
     * 设置默认日期
     */
    setDefaultDates() {
        const today = new Date();
        const thirtyDaysAgo = new Date(today.getTime() - 30 * 24 * 60 * 60 * 1000);
        
        document.getElementById('endDate').value = this.formatDate(today);
        document.getElementById('startDate').value = this.formatDate(thirtyDaysAgo);
    }

    /**
     * 格式化日期为 YYYY-MM-DD
     */
    formatDate(date) {
        return date.toISOString().split('T')[0];
    }

    /**
     * 将日期转换为API需要的格式 YYYYMMDD
     */
    formatDateForAPI(dateString) {
        return dateString.replace(/-/g, '');
    }

    /**
     * 解析YYYYMMDD格式的日期字符串为Date对象
     */
    parseTradeDate(tradeDateStr) {
        if (!tradeDateStr || tradeDateStr.length !== 8) {
            return new Date(); // 返回当前日期作为fallback
        }
        
        const year = tradeDateStr.substring(0, 4);
        const month = tradeDateStr.substring(4, 6);
        const day = tradeDateStr.substring(6, 8);
        return new Date(`${year}-${month}-${day}`);
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
     * 健康检查
     */
    async checkHealth() {
        try {
            const response = await this.makeRequest('/api/v1/health');
            
            if (response.success) {
                this.updateConnectionStatus('online', `服务连接正常 (${this.baseURL})`);
                // 保存成功的URL配置
                this.setServerURL(this.baseURL);
            } else {
                this.updateConnectionStatus('offline', '服务异常');
            }
        } catch (error) {
            console.error('健康检查失败:', error);
            
            // 如果当前URL连接失败，尝试其他常见端口
            if (!await this.tryAlternativePorts()) {
                this.updateConnectionStatus('offline', `连接失败 (${this.baseURL})`);
            }
        }
    }

    /**
     * 尝试连接其他常见端口
     */
    async tryAlternativePorts() {
        const commonPorts = ['8081', '8080', '3000', '8000', '9000'];
        
        for (const port of commonPorts) {
            const testURL = `http://localhost:${port}`;
            
            // 跳过当前已经测试过的URL
            if (testURL === this.baseURL) continue;
            
            try {
                console.log(`尝试连接端口 ${port}...`);
                
                const response = await fetch(`${testURL}/api/v1/health`, {
                    method: 'GET',
                    headers: { 'Content-Type': 'application/json' },
                    signal: AbortSignal.timeout(3000) // 3秒超时
                });
                
                if (response.ok) {
                    const data = await response.json();
                    if (data.success) {
                        console.log(`成功连接到端口 ${port}`);
                        this.setServerURL(testURL);
                        this.updateConnectionStatus('online', `服务连接正常 (${testURL})`);
                        return true;
                    }
                }
            } catch (error) {
                // 继续尝试下一个端口
                console.log(`端口 ${port} 连接失败:`, error.message);
            }
        }
        
        return false;
    }

    /**
     * 更新连接状态显示
     */
    updateConnectionStatus(status, message) {
        const statusDot = document.getElementById('connectionStatus');
        const statusText = document.getElementById('statusText');
        
        statusDot.className = `status-dot ${status}`;
        statusText.textContent = message;
    }

    /**
     * 发起API请求
     */
    async makeRequest(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        
        const defaultOptions = {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            ...options
        };

        try {
            const response = await fetch(url, defaultOptions);
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('API请求失败:', error);
            throw error;
        }
    }

    /**
     * 显示/隐藏加载状态
     */
    setLoading(loading) {
        this.isLoading = loading;
        const overlay = document.getElementById('loadingOverlay');
        const buttons = document.querySelectorAll('.btn');
        
        if (loading) {
            overlay.style.display = 'flex';
            buttons.forEach(btn => btn.disabled = true);
        } else {
            overlay.style.display = 'none';
            buttons.forEach(btn => btn.disabled = false);
        }
    }

    /**
     * 显示错误信息
     */
    showError(message) {
        const errorCard = document.getElementById('errorCard');
        const errorMessage = document.getElementById('errorMessage');
        
        errorMessage.textContent = message;
        errorCard.style.display = 'block';
        errorCard.classList.add('fade-in');
        
        // 隐藏其他卡片
        this.hideAllResultCards();
        errorCard.style.display = 'block';
        
        // 滚动到错误卡片
        errorCard.scrollIntoView({ behavior: 'smooth' });
    }

    /**
     * 隐藏所有结果卡片
     */
    hideAllResultCards() {
        const cards = ['dailyDataCard', 'indicatorsCard', 'predictionsCard', 'errorCard'];
        cards.forEach(cardId => {
            document.getElementById(cardId).style.display = 'none';
        });
    }

    /**
     * 获取股票代码
     */
    getStockCode() {
        const code = document.getElementById('stockCode').value.trim().toUpperCase();
        if (!code) {
            throw new Error('请输入股票代码');
        }
        
        const pattern = /^[0-9]{6}\.(SZ|SH)$/;
        if (!pattern.test(code)) {
            throw new Error('股票代码格式不正确，请使用如：000001.SZ 或 600000.SH');
        }
        
        return code;
    }

    /**
     * 获取日期范围
     */
    getDateRange() {
        const startDate = document.getElementById('startDate').value;
        const endDate = document.getElementById('endDate').value;
        
        if (!startDate || !endDate) {
            throw new Error('请选择开始和结束日期');
        }
        
        if (new Date(startDate) > new Date(endDate)) {
            throw new Error('开始日期不能晚于结束日期');
        }
        
        return {
            startDate: this.formatDateForAPI(startDate),
            endDate: this.formatDateForAPI(endDate)
        };
    }

    /**
     * 处理日线数据查询
     */
    async handleDailyQuery() {
        try {
            this.setLoading(true);
            this.hideAllResultCards();
            
            const stockCode = this.getStockCode();
            const { startDate, endDate } = this.getDateRange();
            
            const endpoint = `/api/v1/stocks/${stockCode}/daily?start_date=${startDate}&end_date=${endDate}`;
            const response = await this.makeRequest(endpoint);
            
            if (response.success && response.data) {
                this.displayDailyData(response.data, stockCode);
            } else {
                throw new Error(response.error || '获取数据失败');
            }
            
        } catch (error) {
            this.showError(`获取日线数据失败: ${error.message}`);
        } finally {
            this.setLoading(false);
        }
    }

    /**
     * 处理技术指标查询
     */
    async handleIndicatorsQuery() {
        try {
            this.setLoading(true);
            this.hideAllResultCards();
            
            const stockCode = this.getStockCode();
            const endpoint = `/api/v1/stocks/${stockCode}/indicators`;
            const response = await this.makeRequest(endpoint);
            
            if (response.success && response.data) {
                this.displayIndicators(response.data, stockCode);
            } else {
                throw new Error(response.error || '获取技术指标失败');
            }
            
        } catch (error) {
            this.showError(`获取技术指标失败: ${error.message}`);
        } finally {
            this.setLoading(false);
        }
    }

    /**
     * 处理买卖预测查询
     */
    async handlePredictionsQuery() {
        try {
            this.setLoading(true);
            this.hideAllResultCards();
            
            const stockCode = this.getStockCode();
            const endpoint = `/api/v1/stocks/${stockCode}/predictions`;
            const response = await this.makeRequest(endpoint);
            
            if (response.success && response.data) {
                this.displayPredictions(response.data, stockCode);
            } else {
                throw new Error(response.error || '获取预测数据失败');
            }
            
        } catch (error) {
            this.showError(`获取买卖预测失败: ${error.message}`);
        } finally {
            this.setLoading(false);
        }
    }

    /**
     * 显示日线数据
     */
    displayDailyData(data, stockCode) {
        const card = document.getElementById('dailyDataCard');
        const summary = document.getElementById('dailyDataSummary');
        
        // 显示卡片
        card.style.display = 'block';
        card.classList.add('fade-in');
        
        // 创建价格图表
        this.createPriceChart(data, stockCode);
        
        // 显示数据摘要
        if (data.length > 0) {
            const latest = data[data.length - 1];
            const previous = data.length > 1 ? data[data.length - 2] : latest;
            
            summary.innerHTML = this.createDataSummary(latest, previous);
        }
        
        // 滚动到结果
        card.scrollIntoView({ behavior: 'smooth' });
    }

    /**
     * 创建数据摘要HTML
     */
    createDataSummary(latest, previous) {
        const change = latest.close - previous.close;
        const changePercent = ((change / previous.close) * 100).toFixed(2);
        const changeClass = change >= 0 ? 'positive' : 'negative';
        const changeSymbol = change >= 0 ? '+' : '';
        
        return `
            <div class="summary-item">
                <div class="label">最新收盘价</div>
                <div class="value">¥${latest.close.toFixed(2)}</div>
                <div class="change ${changeClass}">
                    ${changeSymbol}${change.toFixed(2)} (${changeSymbol}${changePercent}%)
                </div>
            </div>
            <div class="summary-item">
                <div class="label">成交量</div>
                <div class="value">${this.formatVolume(latest.vol)}</div>
            </div>
            <div class="summary-item">
                <div class="label">最高价</div>
                <div class="value">¥${latest.high.toFixed(2)}</div>
            </div>
            <div class="summary-item">
                <div class="label">最低价</div>
                <div class="value">¥${latest.low.toFixed(2)}</div>
            </div>
            <div class="summary-item">
                <div class="label">开盘价</div>
                <div class="value">¥${latest.open.toFixed(2)}</div>
            </div>
            <div class="summary-item">
                <div class="label">换手率</div>
                <div class="value">${latest.turnover_rate ? latest.turnover_rate.toFixed(2) + '%' : 'N/A'}</div>
            </div>
        `;
    }

    /**
     * 格式化成交量
     */
    formatVolume(volume) {
        if (volume >= 100000000) {
            return (volume / 100000000).toFixed(2) + '亿';
        } else if (volume >= 10000) {
            return (volume / 10000).toFixed(2) + '万';
        }
        return volume.toString();
    }

    /**
     * 创建价格图表
     */
    createPriceChart(data, stockCode) {
        const canvas = document.getElementById('priceChart');
        const ctx = canvas.getContext('2d');
        
        // 销毁现有图表
        if (this.currentChart) {
            this.currentChart.destroy();
        }
        
        // 准备数据
        const labels = data.map(item => {
            const date = this.parseTradeDate(item.trade_date);
            return date.toLocaleDateString('zh-CN');
        });
        
        const prices = data.map(item => item.close);
        const volumes = data.map(item => item.vol);
        
        // 创建新图表
        this.currentChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [{
                    label: '收盘价',
                    data: prices,
                    borderColor: '#2563eb',
                    backgroundColor: 'rgba(37, 99, 235, 0.1)',
                    borderWidth: 2,
                    fill: true,
                    tension: 0.1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: `${stockCode} 价格走势`,
                        font: {
                            size: 16,
                            weight: 'bold'
                        }
                    },
                    legend: {
                        position: 'top'
                    }
                },
                scales: {
                    x: {
                        title: {
                            display: true,
                            text: '日期'
                        }
                    },
                    y: {
                        title: {
                            display: true,
                            text: '价格 (¥)'
                        },
                        beginAtZero: false
                    }
                },
                interaction: {
                    intersect: false,
                    mode: 'index'
                }
            }
        });
    }

    /**
     * 显示技术指标
     */
    displayIndicators(data, stockCode) {
        const card = document.getElementById('indicatorsCard');
        const grid = document.getElementById('indicatorsGrid');
        
        // 显示卡片
        card.style.display = 'block';
        card.classList.add('fade-in');
        
        // 构建指标HTML
        let indicatorsHTML = '';
        
        // MACD指标
        if (data.macd) {
            indicatorsHTML += this.createIndicatorItem('MACD', [
                { name: 'DIF', value: data.macd.dif?.toFixed(4) || 'N/A' },
                { name: 'DEA', value: data.macd.dea?.toFixed(4) || 'N/A' },
                { name: '柱状图', value: data.macd.histogram?.toFixed(4) || 'N/A' }
            ], data.macd.signal);
        }
        
        // RSI指标
        if (data.rsi) {
            indicatorsHTML += this.createIndicatorItem('RSI', [
                { name: 'RSI14', value: data.rsi.rsi14?.toFixed(2) || 'N/A' }
            ], data.rsi.signal);
        }
        
        // 布林带
        if (data.boll) {
            indicatorsHTML += this.createIndicatorItem('布林带', [
                { name: '上轨', value: data.boll.upper?.toFixed(2) || 'N/A' },
                { name: '中轨', value: data.boll.middle?.toFixed(2) || 'N/A' },
                { name: '下轨', value: data.boll.lower?.toFixed(2) || 'N/A' }
            ], data.boll.signal);
        }
        
        // 移动平均线
        if (data.ma) {
            indicatorsHTML += this.createIndicatorItem('移动平均线', [
                { name: 'MA5', value: data.ma.ma5?.toFixed(2) || 'N/A' },
                { name: 'MA10', value: data.ma.ma10?.toFixed(2) || 'N/A' },
                { name: 'MA20', value: data.ma.ma20?.toFixed(2) || 'N/A' },
                { name: 'MA60', value: data.ma.ma60?.toFixed(2) || 'N/A' }
            ]);
        }
        
        // KDJ指标
        if (data.kdj) {
            indicatorsHTML += this.createIndicatorItem('KDJ', [
                { name: 'K值', value: data.kdj.k?.toFixed(2) || 'N/A' },
                { name: 'D值', value: data.kdj.d?.toFixed(2) || 'N/A' },
                { name: 'J值', value: data.kdj.j?.toFixed(2) || 'N/A' }
            ], data.kdj.signal);
        }
        
        grid.innerHTML = indicatorsHTML;
        
        // 滚动到结果
        card.scrollIntoView({ behavior: 'smooth' });
    }

    /**
     * 创建指标项HTML
     */
    createIndicatorItem(title, values, signal = null) {
        const signalHTML = signal ? `<span class="signal ${signal.toLowerCase()}">${this.getSignalText(signal)}</span>` : '';
        
        const valuesHTML = values.map(item => 
            `<div class="indicator-value">
                <span class="name">${item.name}:</span>
                <span class="value">${item.value}</span>
            </div>`
        ).join('');
        
        return `
            <div class="indicator-item slide-in">
                <h4>${title} ${signalHTML}</h4>
                <div class="indicator-values">
                    ${valuesHTML}
                </div>
            </div>
        `;
    }

    /**
     * 获取信号文本
     */
    getSignalText(signal) {
        const signalMap = {
            'buy': '买入',
            'sell': '卖出',
            'hold': '持有',
            'neutral': '中性'
        };
        return signalMap[signal.toLowerCase()] || signal;
    }

    /**
     * 显示预测结果
     */
    displayPredictions(data, stockCode) {
        const card = document.getElementById('predictionsCard');
        const container = document.getElementById('predictionsContainer');
        
        // 显示卡片
        card.style.display = 'block';
        card.classList.add('fade-in');
        
        // 构建预测HTML
        let predictionsHTML = '';
        
        // 置信度摘要
        const confidence = data.confidence || 0;
        predictionsHTML += `
            <div class="prediction-summary">
                <div class="confidence-score">${(confidence * 100).toFixed(1)}%</div>
                <p>预测置信度</p>
            </div>
        `;
        
        // 预测列表
        if (data.predictions && data.predictions.length > 0) {
            predictionsHTML += '<div class="predictions-list">';
            
            data.predictions.forEach(prediction => {
                const typeClass = prediction.type.toLowerCase();
                const icon = prediction.type.toLowerCase() === 'buy' ? '📈' : '📉';
                const typeText = prediction.type.toLowerCase() === 'buy' ? '买入' : '卖出';
                
                predictionsHTML += `
                    <div class="prediction-item ${typeClass} slide-in">
                        <div class="prediction-icon">${icon}</div>
                        <div class="prediction-content">
                            <div class="prediction-type">${typeText}信号</div>
                            <div class="prediction-price">¥${prediction.price?.toFixed(2) || 'N/A'}</div>
                            <div class="prediction-probability">概率: ${(prediction.probability * 100).toFixed(1)}%</div>
                            <div class="prediction-reason">${prediction.reason || '基于技术指标分析'}</div>
                        </div>
                    </div>
                `;
            });
            
            predictionsHTML += '</div>';
        } else {
            predictionsHTML += `
                <div class="prediction-item">
                    <div class="prediction-content">
                        <div class="prediction-type">暂无明确信号</div>
                        <div class="prediction-reason">当前市场条件下，没有明确的买卖信号</div>
                    </div>
                </div>
            `;
        }
        
        container.innerHTML = predictionsHTML;
        
        // 滚动到结果
        card.scrollIntoView({ behavior: 'smooth' });
    }
}

// 页面加载完成后初始化客户端
document.addEventListener('DOMContentLoaded', () => {
    // 创建客户端实例
    window.stockClient = new StockAFutureClient();
    
    // 定期检查健康状态
    setInterval(() => {
        if (!window.stockClient.isLoading) {
            window.stockClient.checkHealth();
        }
    }, 30000); // 每30秒检查一次
    
    console.log('Stock-A-Future 网页客户端已初始化');
});
