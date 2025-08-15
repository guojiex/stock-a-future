/**
 * 配置管理模块
 * 负责服务器配置和连接测试
 */

class ConfigModule {
    constructor(client) {
        this.client = client;
        this.init();
    }

    /**
     * 初始化配置模块
     */
    init() {
        this.setupEventListeners();
    }

    /**
     * 设置配置相关事件监听器
     */
    setupEventListeners() {
        // 配置按钮
        const configBtn = document.getElementById('configBtn');
        if (configBtn) {
            configBtn.addEventListener('click', () => this.showConfigModal());
        }
        
        // 关闭配置模态框
        const closeConfigBtn = document.getElementById('closeConfigBtn');
        if (closeConfigBtn) {
            closeConfigBtn.addEventListener('click', () => this.hideConfigModal());
        }
        
        // 点击模态框背景关闭
        const configModal = document.getElementById('configModal');
        if (configModal) {
            configModal.addEventListener('click', (e) => {
                if (e.target.id === 'configModal') {
                    this.hideConfigModal();
                }
            });
        }

        // 测试连接按钮
        const testConnectionBtn = document.getElementById('testConnectionBtn');
        if (testConnectionBtn) {
            testConnectionBtn.addEventListener('click', () => this.testConnection());
        }
        
        // 保存配置按钮
        const saveConfigBtn = document.getElementById('saveConfigBtn');
        if (saveConfigBtn) {
            saveConfigBtn.addEventListener('click', () => this.saveConfig());
        }

        // 端口快捷按钮
        const portBtns = document.querySelectorAll('.port-btn');
        portBtns.forEach(btn => {
            btn.addEventListener('click', (e) => {
                const port = e.target.dataset.port;
                const serverURLInput = document.getElementById('serverURL');
                if (serverURLInput) {
                    serverURLInput.value = `http://localhost:${port}`;
                }
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
        
        if (!modal || !serverURLInput) return;
        
        // 设置当前服务器URL
        serverURLInput.value = this.client.baseURL;
        
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
        
        if (modal) modal.style.display = 'none';
        if (testResult) testResult.style.display = 'none';
    }

    /**
     * 测试连接
     */
    async testConnection() {
        const serverURL = document.getElementById('serverURL')?.value.trim();
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

        if (testBtn) {
            testBtn.disabled = true;
            testBtn.textContent = '🔍 测试中...';
        }
        
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
            if (testBtn) {
                testBtn.disabled = false;
                testBtn.textContent = '🔗 测试连接';
            }
        }
    }

    /**
     * 显示测试结果
     */
    showTestResult(message, type) {
        const testResult = document.getElementById('connectionTestResult');
        if (testResult) {
            testResult.className = `test-result ${type}`;
            testResult.textContent = message;
            testResult.style.display = 'block';
        }
    }

    /**
     * 保存配置
     */
    saveConfig() {
        const serverURL = document.getElementById('serverURL')?.value.trim();
        
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
        this.client.setServerURL(serverURL);
        
        // 保存到本地存储
        localStorage.setItem('serverURL', serverURL);
        
        // 显示保存成功消息
        this.showTestResult('✅ 配置已保存', 'success');
        
        // 隐藏配置模态框
        setTimeout(() => {
            this.hideConfigModal();
        }, 1500);
    }
}

// 导出配置模块类
window.ConfigModule = ConfigModule;
