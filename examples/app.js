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
     * 获取股票基本信息
     */
    async getStockBasic(stockCode) {
        try {
            const endpoint = `/api/v1/stocks/${stockCode}/basic`;
            const response = await this.makeRequest(endpoint);
            
            if (response.success && response.data) {
                return response.data;
            } else {
                throw new Error(response.error || '获取股票基本信息失败');
            }
        } catch (error) {
            console.warn(`获取股票基本信息失败: ${error.message}`);
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
            '600745.SH': '闻泰科技',
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
            
            // 并行获取股票基本信息和日线数据
            const [stockBasic, dailyResponse] = await Promise.all([
                this.getStockBasic(stockCode),
                this.makeRequest(`/api/v1/stocks/${stockCode}/daily?start_date=${startDate}&end_date=${endDate}`)
            ]);
            
            if (dailyResponse.success && dailyResponse.data) {
                this.displayDailyData(dailyResponse.data, stockCode, stockBasic);
            } else {
                throw new Error(dailyResponse.error || '获取数据失败');
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
            
            // 并行获取股票基本信息和技术指标
            const [stockBasic, indicatorsResponse] = await Promise.all([
                this.getStockBasic(stockCode),
                this.makeRequest(`/api/v1/stocks/${stockCode}/indicators`)
            ]);
            
            if (indicatorsResponse.success && indicatorsResponse.data) {
                this.displayIndicators(indicatorsResponse.data, stockCode, stockBasic);
            } else {
                throw new Error(indicatorsResponse.error || '获取技术指标失败');
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
            
            // 并行获取股票基本信息和预测数据
            const [stockBasic, predictionsResponse] = await Promise.all([
                this.getStockBasic(stockCode),
                this.makeRequest(`/api/v1/stocks/${stockCode}/predictions`)
            ]);
            
            if (predictionsResponse.success && predictionsResponse.data) {
                this.displayPredictions(predictionsResponse.data, stockCode, stockBasic);
            } else {
                throw new Error(predictionsResponse.error || '获取预测数据失败');
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
    displayDailyData(data, stockCode, stockBasic) {
        const card = document.getElementById('dailyDataCard');
        const summary = document.getElementById('dailyDataSummary');
        
        // 更新卡片标题
        const cardTitle = card.querySelector('h3');
        if (cardTitle && stockBasic && stockBasic.name) {
            cardTitle.textContent = `📈 ${stockBasic.name}(${stockCode}) - 日线数据`;
        }
        
        // 显示卡片
        card.style.display = 'block';
        card.classList.add('fade-in');
        
        // 创建价格图表
        this.createPriceChart(data, stockCode, stockBasic);
        
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
        const currentClose = parseFloat(latest.close);
        const previousClose = parseFloat(latest.pre_close || previous.close);
        const change = currentClose - previousClose;
        const changePercent = ((change / previousClose) * 100).toFixed(2);
        const changeClass = change >= 0 ? 'positive' : 'negative';
        const changeSymbol = change >= 0 ? '+' : '';
        
        // 计算振幅
        const amplitude = ((parseFloat(latest.high) - parseFloat(latest.low)) / previousClose * 100).toFixed(2);
        
        return `
            <div class="summary-item">
                <div class="label">最新收盘价</div>
                <div class="value">¥${currentClose.toFixed(2)}</div>
                <div class="change ${changeClass}">
                    ${changeSymbol}${change.toFixed(2)} (${changeSymbol}${changePercent}%)
                </div>
            </div>
            <div class="summary-item">
                <div class="label">成交量</div>
                <div class="value">${this.formatVolume(parseFloat(latest.vol))}</div>
            </div>
            <div class="summary-item">
                <div class="label">最高价</div>
                <div class="value">¥${parseFloat(latest.high).toFixed(2)}</div>
            </div>
            <div class="summary-item">
                <div class="label">最低价</div>
                <div class="value">¥${parseFloat(latest.low).toFixed(2)}</div>
            </div>
            <div class="summary-item">
                <div class="label">开盘价</div>
                <div class="value">¥${parseFloat(latest.open).toFixed(2)}</div>
            </div>
            <div class="summary-item">
                <div class="label">成交额</div>
                <div class="value">${this.formatAmount(parseFloat(latest.amount))}</div>
            </div>
            <div class="summary-item">
                <div class="label">振幅</div>
                <div class="value">${amplitude}%</div>
            </div>
            <div class="summary-item">
                <div class="label">昨收价</div>
                <div class="value">¥${previousClose.toFixed(2)}</div>
            </div>
        `;
    }
    
    /**
     * 格式化成交额
     */
    formatAmount(amount) {
        // 成交额单位是千元，需要转换
        const amountInYuan = amount * 1000;
        if (amountInYuan >= 100000000) {
            return (amountInYuan / 100000000).toFixed(2) + '亿元';
        } else if (amountInYuan >= 10000) {
            return (amountInYuan / 10000).toFixed(2) + '万元';
        }
        return amountInYuan.toFixed(0) + '元';
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
     * 创建K线图表（使用ECharts）
     */
    createPriceChart(data, stockCode, stockBasic) {
        const chartContainer = document.getElementById('priceChart');
        
        // 销毁现有图表
        if (this.currentChart) {
            this.currentChart.dispose();
        }
        
        // 初始化ECharts实例
        this.currentChart = echarts.init(chartContainer);
        
        // 准备K线数据 [开盘价, 收盘价, 最低价, 最高价]
        const klineData = data.map(item => [
            parseFloat(item.open),
            parseFloat(item.close),
            parseFloat(item.low),
            parseFloat(item.high)
        ]);
        
        // 准备成交量数据
        const volumeData = data.map(item => parseFloat(item.vol));
        
        // 准备日期标签
        const dates = data.map(item => {
            const dateStr = item.trade_date;
            return `${dateStr.substring(0,4)}-${dateStr.substring(4,6)}-${dateStr.substring(6,8)}`;
        });
        
        // 计算移动平均线
        const ma5 = this.calculateMA(data.map(item => parseFloat(item.close)), 5);
        const ma10 = this.calculateMA(data.map(item => parseFloat(item.close)), 10);
        const ma20 = this.calculateMA(data.map(item => parseFloat(item.close)), 20);
        
        // 构建图表标题
        let chartTitle = `${stockCode} K线图`;
        if (stockBasic && stockBasic.name) {
            chartTitle = `${stockBasic.name}(${stockCode}) K线图`;
        }
        
        // 配置图表选项
        const option = {
            title: {
                text: chartTitle,
                left: 'center',
                textStyle: {
                    fontSize: 16,
                    fontWeight: 'bold'
                }
            },
            tooltip: {
                trigger: 'axis',
                axisPointer: {
                    type: 'cross'
                },
                formatter: function(params) {
                    let result = params[0].name + '<br/>';
                    
                    // K线数据
                    if (params[0].componentSubType === 'candlestick') {
                        const data = params[0].data;
                        result += `开盘价: ¥${data[1].toFixed(2)}<br/>`;
                        result += `收盘价: ¥${data[2].toFixed(2)}<br/>`;
                        result += `最低价: ¥${data[3].toFixed(2)}<br/>`;
                        result += `最高价: ¥${data[4].toFixed(2)}<br/>`;
                        
                        // 涨跌幅计算
                        const change = data[2] - data[1];
                        const changePercent = ((change / data[1]) * 100).toFixed(2);
                        const changeText = change >= 0 ? `+${change.toFixed(2)}` : change.toFixed(2);
                        const percentText = change >= 0 ? `+${changePercent}%` : `${changePercent}%`;
                        result += `涨跌额: ${changeText}<br/>`;
                        result += `涨跌幅: ${percentText}<br/>`;
                    }
                    
                    // 成交量
                    const volumeParam = params.find(p => p.seriesName === '成交量');
                    if (volumeParam) {
                        result += `成交量: ${this.formatVolume(volumeParam.data)}<br/>`;
                    }
                    
                    // 移动平均线
                    params.forEach(param => {
                        if (param.seriesName.includes('MA')) {
                            result += `${param.seriesName}: ¥${param.data.toFixed(2)}<br/>`;
                        }
                    });
                    
                    return result;
                }.bind(this)
            },
            legend: {
                data: ['K线', 'MA5', 'MA10', 'MA20', '成交量'],
                top: 30
            },
            grid: [
                {
                    left: '5%',
                    right: '5%',
                    top: '15%',
                    height: '60%'
                },
                {
                    left: '5%',
                    right: '5%',
                    top: '80%',
                    height: '15%'
                }
            ],
            xAxis: [
                {
                    type: 'category',
                    data: dates,
                    scale: true,
                    boundaryGap: false,
                    axisLine: { onZero: false },
                    splitLine: { show: false },
                    splitNumber: 20,
                    min: 'dataMin',
                    max: 'dataMax'
                },
                {
                    type: 'category',
                    gridIndex: 1,
                    data: dates,
                    scale: true,
                    boundaryGap: false,
                    axisLine: { onZero: false },
                    axisTick: { show: false },
                    splitLine: { show: false },
                    axisLabel: { show: false },
                    splitNumber: 20,
                    min: 'dataMin',
                    max: 'dataMax'
                }
            ],
            yAxis: [
                {
                    scale: true,
                    splitArea: {
                        show: true
                    },
                    axisLabel: {
                        formatter: '¥{value}'
                    }
                },
                {
                    scale: true,
                    gridIndex: 1,
                    splitNumber: 2,
                    axisLabel: { show: false },
                    axisLine: { show: false },
                    axisTick: { show: false },
                    splitLine: { show: false }
                }
            ],
            dataZoom: [
                {
                    type: 'inside',
                    xAxisIndex: [0, 1],
                    start: 70,
                    end: 100
                },
                {
                    show: true,
                    xAxisIndex: [0, 1],
                    type: 'slider',
                    top: '90%',
                    start: 70,
                    end: 100
                }
            ],
            series: [
                {
                    name: 'K线',
                    type: 'candlestick',
                    data: klineData,
                    itemStyle: {
                        color: '#ef4444',        // 阳线颜色（红色）
                        color0: '#10b981',       // 阴线颜色（绿色）
                        borderColor: '#ef4444',   // 阳线边框
                        borderColor0: '#10b981'   // 阴线边框
                    }
                },
                {
                    name: 'MA5',
                    type: 'line',
                    data: ma5,
                    smooth: true,
                    lineStyle: {
                        opacity: 0.8,
                        width: 1,
                        color: '#2563eb'
                    },
                    showSymbol: false
                },
                {
                    name: 'MA10',
                    type: 'line',
                    data: ma10,
                    smooth: true,
                    lineStyle: {
                        opacity: 0.8,
                        width: 1,
                        color: '#f59e0b'
                    },
                    showSymbol: false
                },
                {
                    name: 'MA20',
                    type: 'line',
                    data: ma20,
                    smooth: true,
                    lineStyle: {
                        opacity: 0.8,
                        width: 1,
                        color: '#8b5cf6'
                    },
                    showSymbol: false
                },
                {
                    name: '成交量',
                    type: 'bar',
                    xAxisIndex: 1,
                    yAxisIndex: 1,
                    data: volumeData.map((vol, index) => {
                        // 根据K线涨跌设置成交量柱子颜色
                        const kline = klineData[index];
                        const color = kline[1] >= kline[0] ? '#ef4444' : '#10b981';
                        return {
                            value: vol,
                            itemStyle: { color: color, opacity: 0.6 }
                        };
                    })
                }
            ]
        };
        
        // 设置图表配置并渲染
        this.currentChart.setOption(option);
        
        // 响应式调整
        window.addEventListener('resize', () => {
            if (this.currentChart) {
                this.currentChart.resize();
            }
        });
    }
    
    /**
     * 计算移动平均线
     */
    calculateMA(data, period) {
        const result = [];
        for (let i = 0; i < data.length; i++) {
            if (i < period - 1) {
                result.push(null);
            } else {
                let sum = 0;
                for (let j = i - period + 1; j <= i; j++) {
                    sum += data[j];
                }
                result.push(sum / period);
            }
        }
        return result;
    }

    /**
     * 显示技术指标
     */
    displayIndicators(data, stockCode, stockBasic) {
        const card = document.getElementById('indicatorsCard');
        const grid = document.getElementById('indicatorsGrid');
        
        // 更新卡片标题
        const cardTitle = card.querySelector('h3');
        if (cardTitle && stockBasic && stockBasic.name) {
            cardTitle.textContent = `📊 ${stockBasic.name}(${stockCode}) - 技术指标`;
        }
        
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
    displayPredictions(data, stockCode, stockBasic) {
        const card = document.getElementById('predictionsCard');
        const container = document.getElementById('predictionsContainer');
        
        // 更新卡片标题
        const cardTitle = card.querySelector('h3');
        if (cardTitle && stockBasic && stockBasic.name) {
            cardTitle.textContent = `🎯 ${stockBasic.name}(${stockCode}) - 买卖点预测`;
        }
        
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
