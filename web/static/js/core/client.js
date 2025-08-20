/**
 * Stock-A-Future API 核心客户端类
 * 负责客户端的初始化和核心功能管理
 */

class StockAFutureClient {
    constructor(baseURL = null) {
        this.baseURL = baseURL || this.detectServerURL();
        console.log(`StockAFutureClient 初始化，baseURL: ${this.baseURL}`);
        this.currentChart = null;
        this.isLoading = false;
        this.searchTimeout = null;
        this.selectedSuggestionIndex = -1;
        
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
            console.log(`从URL参数检测到端口: ${serverPort}`);
            return `http://localhost:${serverPort}`;
        }

        // 2. 从localStorage中读取上次保存的配置
        const savedURL = localStorage.getItem('stockapi_server_url');
        if (savedURL) {
            console.log(`从localStorage读取到保存的URL: ${savedURL}`);
            return savedURL;
        }

        // 3. 尝试常见端口
        console.log(`使用默认端口: 8081`);
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
        // 导入并初始化各个模块
        await this.initializeModules();
        await this.checkHealth();
    }

    /**
     * 初始化各个模块
     */
    async initializeModules() {
        // 这里会在main.js中调用各个模块的初始化方法
        console.log('客户端模块初始化完成');
    }

    /**
     * 健康检查
     */
    async checkHealth() {
        try {
            // 使用轻量级健康检查，不进行数据源连接测试
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
        // 只尝试8081端口，因为我们的服务器运行在这个端口
        const commonPorts = ['8081'];
        
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
        
        if (statusDot && statusText) {
            statusDot.className = `status-dot ${status}`;
            statusText.textContent = message;
        }
    }

    /**
     * 发起API请求
     */
    async makeRequest(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const requestId = Math.random().toString(36).substr(2, 9);
        
        // 记录请求详情
        console.log(`[Client] 发起API请求 - ID: ${requestId}`, {
            url,
            baseURL: this.baseURL,
            endpoint,
            options,
            timestamp: new Date().toISOString()
        });
        
        const defaultOptions = {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            ...options
        };

        try {
            const response = await fetch(url, defaultOptions);
            
            // 记录响应状态
            console.log(`[Client] API响应状态 - ID: ${requestId}`, {
                url,
                status: response.status,
                statusText: response.statusText,
                ok: response.ok,
                headers: Object.fromEntries(response.headers.entries()),
                timestamp: new Date().toISOString()
            });
            
            if (!response.ok) {
                // 详细记录HTTP错误
                const errorDetails = {
                    requestId,
                    url,
                    status: response.status,
                    statusText: response.statusText,
                    headers: Object.fromEntries(response.headers.entries()),
                    timestamp: new Date().toISOString()
                };
                
                console.error(`[Client] HTTP错误 - ID: ${requestId}`, errorDetails);
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            const data = await response.json();
            
            // 记录响应数据
            console.log(`[Client] API响应数据 - ID: ${requestId}`, {
                url,
                success: data.success,
                hasData: !!data.data,
                dataType: data.data ? typeof data.data : 'undefined',
                dataLength: data.data && Array.isArray(data.data) ? data.data.length : 'N/A',
                error: data.error || null,
                message: data.message || null,
                timestamp: new Date().toISOString()
            });
            
            return data;
        } catch (error) {
            // 详细记录异常信息
            const errorDetails = {
                requestId,
                url,
                baseURL: this.baseURL,
                endpoint,
                options: defaultOptions,
                errorType: error.constructor.name,
                message: error.message,
                stack: error.stack,
                timestamp: new Date().toISOString()
            };
            
            console.error(`[Client] API请求异常 - ID: ${requestId}`, errorDetails);
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
            if (overlay) overlay.style.display = 'flex';
            buttons.forEach(btn => btn.disabled = true);
        } else {
            if (overlay) overlay.style.display = 'none';
            buttons.forEach(btn => btn.disabled = false);
        }
    }

    /**
     * 显示错误信息
     */
    showError(message) {
        const errorSection = document.getElementById('error-section');
        const errorCard = document.getElementById('errorCard');
        const errorMessage = document.getElementById('errorMessage');
        
        if (errorSection && errorCard && errorMessage) {
            errorMessage.textContent = message;
            errorSection.style.display = 'block';
            errorSection.classList.add('fade-in');
            
            // 隐藏其他section
            this.hideAllResultSections();
            errorSection.style.display = 'block';
            
            // 滚动到错误section
            errorSection.scrollIntoView({ behavior: 'smooth' });
        }
    }

    /**
     * 隐藏所有结果section
     */
    hideAllResultSections() {
        const sections = ['daily-chart-section', 'error-section'];
        sections.forEach(sectionId => {
            const section = document.getElementById(sectionId);
            if (section) section.style.display = 'none';
        });
    }

    /**
     * 获取股票代码
     */
    getStockCode() {
        const code = document.getElementById('stockCode')?.value.trim().toUpperCase();
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
        const startDate = document.getElementById('startDate')?.value;
        const endDate = document.getElementById('endDate')?.value;
        
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
     * 将日期转换为API需要的格式 YYYYMMDD
     */
    formatDateForAPI(dateString) {
        return dateString.replace(/-/g, '');
    }

    /**
     * 格式化日期为 YYYY-MM-DD
     * @param {Date} date - 日期对象
     * @returns {string} - 格式化后的日期字符串 YYYY-MM-DD
     */
    formatDate(date) {
        if (!date) {
            console.warn('[Client] formatDate 收到无效日期');
            return '';
        }
        
        try {
            // 确保是Date对象
            const dateObj = date instanceof Date ? date : new Date(date);
            
            // 检查是否有效日期
            if (isNaN(dateObj.getTime())) {
                console.warn('[Client] formatDate 收到无效日期:', date);
                return String(date);
            }
            
            // 获取年月日并格式化
            const year = dateObj.getFullYear();
            // 月份需要+1，因为getMonth()返回0-11
            const month = (dateObj.getMonth() + 1).toString().padStart(2, '0');
            const day = dateObj.getDate().toString().padStart(2, '0');
            
            return `${year}-${month}-${day}`;
        } catch (error) {
            console.error('[Client] 日期格式化失败:', error);
            return String(date);
        }
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
     * 显示消息提示
     * @param {string} message - 消息内容
     * @param {string} type - 消息类型: 'success', 'error', 'warning', 'info'
     * @param {number} duration - 显示时长(毫秒)，默认3000
     */
    showMessage(message, type = 'info', duration = 3000) {
        // 移除已有的消息
        const existingMessages = document.querySelectorAll('.message');
        existingMessages.forEach(msg => msg.remove());
        
        // 创建新消息元素
        const messageElement = document.createElement('div');
        messageElement.className = `message message-${type}`;
        messageElement.textContent = message;
        
        // 添加到body
        document.body.appendChild(messageElement);
        
        // 设置自动消失
        setTimeout(() => {
            messageElement.style.opacity = '0';
            messageElement.style.transform = 'translateX(100%)';
            
            // 动画结束后移除元素
            setTimeout(() => {
                messageElement.remove();
            }, 500);
        }, duration);
        
        return messageElement;
    }
}

// 导出客户端类
window.StockAFutureClient = StockAFutureClient;
