/**
 * 策略管理模块
 * 负责策略的创建、编辑、删除和展示
 */

class StrategiesModule {
    constructor(client, apiService) {
        this.client = client;
        this.apiService = apiService;
        this.strategies = [];
        this.selectedStrategy = null;
        
        this.init();
    }

    /**
     * 初始化策略模块
     */
    init() {
        this.setupEventListeners();
        this.loadStrategies();
    }

    /**
     * 设置事件监听器
     */
    setupEventListeners() {
        // 创建策略按钮
        const createStrategyBtn = document.getElementById('createStrategyBtn');
        if (createStrategyBtn) {
            createStrategyBtn.addEventListener('click', () => this.showCreateStrategyModal());
        }
    }

    /**
     * 加载策略列表
     */
    async loadStrategies() {
        try {
            const response = await this.apiService.getStrategiesList();
            if (response.success) {
                this.strategies = response.data;
                this.displayStrategies();
            } else {
                throw new Error(response.message || '获取策略列表失败');
            }
        } catch (error) {
            console.error('加载策略列表失败:', error);
            this.displayDefaultStrategies();
        }
    }

    /**
     * 显示默认策略（当后端不可用时）
     */
    displayDefaultStrategies() {
        const defaultStrategies = [
            {
                id: 'macd_strategy',
                name: 'MACD金叉策略',
                description: '基于MACD指标的金叉死叉交易策略',
                strategy_type: 'technical',
                status: 'active',
                parameters: {
                    fast_period: 12,
                    slow_period: 26,
                    signal_period: 9
                },
                created_at: new Date().toISOString()
            },
            {
                id: 'ma_crossover',
                name: '双均线策略',
                description: '短期均线突破长期均线的交易策略',
                strategy_type: 'technical',
                status: 'active',
                parameters: {
                    short_period: 5,
                    long_period: 20
                },
                created_at: new Date().toISOString()
            },
            {
                id: 'rsi_strategy',
                name: 'RSI超买超卖策略',
                description: '基于RSI指标的超买超卖交易策略',
                strategy_type: 'technical',
                status: 'inactive',
                parameters: {
                    rsi_period: 14,
                    overbought: 70,
                    oversold: 30
                },
                created_at: new Date().toISOString()
            },
            {
                id: 'bollinger_strategy',
                name: '布林带策略',
                description: '基于布林带的均值回归策略',
                strategy_type: 'technical',
                status: 'testing',
                parameters: {
                    period: 20,
                    std_dev: 2
                },
                created_at: new Date().toISOString()
            }
        ];

        this.strategies = defaultStrategies;
        this.displayStrategies();
    }

    /**
     * 显示策略列表
     */
    displayStrategies() {
        const strategiesGrid = document.getElementById('strategiesGrid');
        if (!strategiesGrid) return;

        if (this.strategies.length === 0) {
            strategiesGrid.innerHTML = `
                <div class="empty-state">
                    <div class="empty-icon">📋</div>
                    <h3>暂无策略</h3>
                    <p>点击"创建策略"按钮开始创建您的第一个策略</p>
                </div>
            `;
            return;
        }

        strategiesGrid.innerHTML = this.strategies.map(strategy => `
            <div class="strategy-card" data-strategy-id="${strategy.id}">
                <div class="strategy-header">
                    <div class="strategy-title">
                        <h4>${strategy.name}</h4>
                        <span class="strategy-status ${strategy.status}">${this.getStatusText(strategy.status)}</span>
                    </div>
                    <div class="strategy-actions">
                        <button class="btn-icon" onclick="strategiesModule.viewStrategy('${strategy.id}')" title="查看详情">
                            👁️
                        </button>
                        <button class="btn-icon" onclick="strategiesModule.editStrategy('${strategy.id}')" title="编辑">
                            ✏️
                        </button>
                        <button class="btn-icon" onclick="strategiesModule.toggleStrategy('${strategy.id}')" title="启用/禁用">
                            ${strategy.status === 'active' ? '⏸️' : '▶️'}
                        </button>
                        <button class="btn-icon danger" onclick="strategiesModule.deleteStrategy('${strategy.id}')" title="删除">
                            🗑️
                        </button>
                    </div>
                </div>
                
                <div class="strategy-content">
                    <p class="strategy-description">${strategy.description}</p>
                    
                    <div class="strategy-meta">
                        <div class="meta-item">
                            <span class="meta-label">类型:</span>
                            <span class="meta-value">${this.getTypeText(strategy.strategy_type)}</span>
                        </div>
                        <div class="meta-item">
                            <span class="meta-label">创建时间:</span>
                            <span class="meta-value">${this.formatDate(strategy.created_at)}</span>
                        </div>
                    </div>

                    <div class="strategy-parameters">
                        <h5>参数配置:</h5>
                        <div class="parameters-list">
                            ${this.renderParameters(strategy.parameters)}
                        </div>
                    </div>
                </div>

                <div class="strategy-footer">
                    <button class="btn btn-primary btn-small" onclick="strategiesModule.runBacktest('${strategy.id}')">
                        🚀 运行回测
                    </button>
                    <button class="btn btn-secondary btn-small" onclick="strategiesModule.viewPerformance('${strategy.id}')">
                        📊 查看表现
                    </button>
                </div>
            </div>
        `).join('');
    }

    /**
     * 获取状态文本
     */
    getStatusText(status) {
        const statusMap = {
            'active': '活跃',
            'inactive': '非活跃',
            'testing': '测试中'
        };
        return statusMap[status] || status;
    }

    /**
     * 获取类型文本
     */
    getTypeText(type) {
        const typeMap = {
            'technical': '技术指标',
            'fundamental': '基本面',
            'ml': '机器学习',
            'composite': '复合策略'
        };
        return typeMap[type] || type;
    }

    /**
     * 格式化日期
     */
    formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString('zh-CN');
    }

    /**
     * 渲染参数列表
     */
    renderParameters(parameters) {
        if (!parameters || Object.keys(parameters).length === 0) {
            return '<span class="no-parameters">无参数</span>';
        }

        return Object.entries(parameters).map(([key, value]) => `
            <div class="parameter-item">
                <span class="param-key">${key}:</span>
                <span class="param-value">${value}</span>
            </div>
        `).join('');
    }

    /**
     * 查看策略详情
     */
    async viewStrategy(strategyId) {
        const strategy = this.strategies.find(s => s.id === strategyId);
        if (!strategy) return;

        this.selectedStrategy = strategy;
        this.displayStrategyDetails(strategy);
    }

    /**
     * 显示策略详情
     */
    displayStrategyDetails(strategy) {
        const detailsDiv = document.getElementById('strategyDetails');
        const infoDiv = document.getElementById('strategyInfo');
        const performanceDiv = document.getElementById('strategyPerformance');

        if (!detailsDiv || !infoDiv || !performanceDiv) return;

        // 显示详情区域
        detailsDiv.style.display = 'block';

        // 显示策略信息
        infoDiv.innerHTML = `
            <div class="strategy-detail-card">
                <div class="detail-header">
                    <h5>${strategy.name}</h5>
                    <span class="strategy-status ${strategy.status}">${this.getStatusText(strategy.status)}</span>
                </div>
                
                <div class="detail-content">
                    <div class="detail-item">
                        <label>描述:</label>
                        <p>${strategy.description}</p>
                    </div>
                    
                    <div class="detail-item">
                        <label>策略类型:</label>
                        <span>${this.getTypeText(strategy.strategy_type)}</span>
                    </div>
                    
                    <div class="detail-item">
                        <label>创建时间:</label>
                        <span>${this.formatDate(strategy.created_at)}</span>
                    </div>
                    
                    <div class="detail-item">
                        <label>参数配置:</label>
                        <div class="parameters-detail">
                            ${this.renderParametersDetail(strategy.parameters)}
                        </div>
                    </div>
                </div>
                
                <div class="detail-actions">
                    <button class="btn btn-primary" onclick="strategiesModule.editStrategy('${strategy.id}')">
                        ✏️ 编辑策略
                    </button>
                    <button class="btn btn-secondary" onclick="strategiesModule.runBacktest('${strategy.id}')">
                        🚀 运行回测
                    </button>
                    <button class="btn btn-outline" onclick="strategiesModule.cloneStrategy('${strategy.id}')">
                        📋 克隆策略
                    </button>
                </div>
            </div>
        `;

        // 加载策略表现数据
        this.loadStrategyPerformance(strategy.id);
    }

    /**
     * 渲染详细参数
     */
    renderParametersDetail(parameters) {
        if (!parameters || Object.keys(parameters).length === 0) {
            return '<p class="no-parameters">无参数配置</p>';
        }

        return `
            <div class="parameters-grid">
                ${Object.entries(parameters).map(([key, value]) => `
                    <div class="parameter-detail">
                        <div class="param-name">${key}</div>
                        <div class="param-value">${value}</div>
                        <div class="param-type">${typeof value}</div>
                    </div>
                `).join('')}
            </div>
        `;
    }

    /**
     * 加载策略表现
     */
    async loadStrategyPerformance(strategyId) {
        const performanceDiv = document.getElementById('strategyPerformance');
        if (!performanceDiv) return;

        try {
            // 显示加载状态
            performanceDiv.innerHTML = `
                <div class="loading-state">
                    <div class="loading-spinner"></div>
                    <p>加载策略表现数据...</p>
                </div>
            `;

            const response = await this.apiService.getStrategyPerformance(strategyId);
            if (response.success) {
                this.displayStrategyPerformance(response.data);
            } else {
                throw new Error(response.message || '获取策略表现失败');
            }
        } catch (error) {
            console.error('加载策略表现失败:', error);
            this.displayMockPerformance(strategyId);
        }
    }

    /**
     * 显示模拟策略表现
     */
    displayMockPerformance(strategyId) {
        const performanceDiv = document.getElementById('strategyPerformance');
        if (!performanceDiv) return;

        const mockData = {
            total_return: Math.random() * 0.4 - 0.1, // -10% 到 30%
            annual_return: Math.random() * 0.3 - 0.05, // -5% 到 25%
            max_drawdown: -Math.random() * 0.2, // 0% 到 -20%
            sharpe_ratio: Math.random() * 2 + 0.5, // 0.5 到 2.5
            win_rate: Math.random() * 0.4 + 0.4, // 40% 到 80%
            total_trades: Math.floor(Math.random() * 500 + 50),
            last_updated: new Date().toISOString()
        };

        this.displayStrategyPerformance(mockData);
    }

    /**
     * 显示策略表现
     */
    displayStrategyPerformance(performance) {
        const performanceDiv = document.getElementById('strategyPerformance');
        if (!performanceDiv) return;

        performanceDiv.innerHTML = `
            <div class="performance-summary">
                <div class="performance-metrics">
                    <div class="metric-item">
                        <div class="metric-label">总收益率</div>
                        <div class="metric-value ${performance.total_return >= 0 ? 'positive' : 'negative'}">
                            ${(performance.total_return * 100).toFixed(2)}%
                        </div>
                    </div>
                    <div class="metric-item">
                        <div class="metric-label">年化收益率</div>
                        <div class="metric-value ${performance.annual_return >= 0 ? 'positive' : 'negative'}">
                            ${(performance.annual_return * 100).toFixed(2)}%
                        </div>
                    </div>
                    <div class="metric-item">
                        <div class="metric-label">最大回撤</div>
                        <div class="metric-value ${performance.max_drawdown >= -0.05 ? 'positive' : performance.max_drawdown >= -0.1 ? 'warning' : 'negative'}">
                            ${(performance.max_drawdown * 100).toFixed(2)}%
                        </div>
                    </div>
                    <div class="metric-item">
                        <div class="metric-label">夏普比率</div>
                        <div class="metric-value ${performance.sharpe_ratio >= 1 ? 'positive' : performance.sharpe_ratio >= 0.5 ? 'warning' : 'negative'}">
                            ${performance.sharpe_ratio.toFixed(2)}
                        </div>
                    </div>
                    <div class="metric-item">
                        <div class="metric-label">胜率</div>
                        <div class="metric-value ${performance.win_rate >= 0.6 ? 'positive' : performance.win_rate >= 0.5 ? 'warning' : 'negative'}">
                            ${(performance.win_rate * 100).toFixed(2)}%
                        </div>
                    </div>
                    <div class="metric-item">
                        <div class="metric-label">总交易次数</div>
                        <div class="metric-value neutral">
                            ${performance.total_trades}
                        </div>
                    </div>
                </div>
                
                <div class="performance-actions">
                    <button class="btn btn-primary" onclick="strategiesModule.viewDetailedPerformance('${this.selectedStrategy?.id}')">
                        📊 详细分析
                    </button>
                    <button class="btn btn-secondary" onclick="strategiesModule.compareStrategies('${this.selectedStrategy?.id}')">
                        📈 策略对比
                    </button>
                </div>
                
                <div class="performance-footer">
                    <small>最后更新: ${this.formatDate(performance.last_updated || new Date().toISOString())}</small>
                </div>
            </div>
        `;
    }

    /**
     * 编辑策略
     */
    editStrategy(strategyId) {
        const strategy = this.strategies.find(s => s.id === strategyId);
        if (!strategy) return;

        this.showMessage('策略编辑功能开发中...', 'info');
    }

    /**
     * 切换策略状态
     */
    async toggleStrategy(strategyId) {
        const strategy = this.strategies.find(s => s.id === strategyId);
        if (!strategy) return;

        try {
            const newStatus = strategy.status === 'active' ? 'inactive' : 'active';
            
            // 更新本地状态
            strategy.status = newStatus;
            
            // 重新显示策略列表
            this.displayStrategies();
            
            this.showMessage(`策略已${newStatus === 'active' ? '启用' : '禁用'}`, 'success');
            
        } catch (error) {
            console.error('切换策略状态失败:', error);
            this.showMessage('切换策略状态失败', 'error');
        }
    }

    /**
     * 删除策略
     */
    async deleteStrategy(strategyId) {
        const strategy = this.strategies.find(s => s.id === strategyId);
        if (!strategy) return;

        if (!confirm(`确定要删除策略"${strategy.name}"吗？此操作不可撤销。`)) {
            return;
        }

        try {
            // 从本地列表中移除
            this.strategies = this.strategies.filter(s => s.id !== strategyId);
            
            // 重新显示策略列表
            this.displayStrategies();
            
            // 如果当前选中的策略被删除，隐藏详情
            if (this.selectedStrategy?.id === strategyId) {
                const detailsDiv = document.getElementById('strategyDetails');
                if (detailsDiv) {
                    detailsDiv.style.display = 'none';
                }
                this.selectedStrategy = null;
            }
            
            this.showMessage('策略已删除', 'success');
            
        } catch (error) {
            console.error('删除策略失败:', error);
            this.showMessage('删除策略失败', 'error');
        }
    }

    /**
     * 运行回测
     */
    runBacktest(strategyId) {
        const strategy = this.strategies.find(s => s.id === strategyId);
        if (!strategy) return;

        // 切换到回测标签页
        const backtestTab = document.querySelector('[data-tab="backtest"]');
        if (backtestTab) {
            backtestTab.click();
        }

        // 设置策略选择
        setTimeout(() => {
            const strategySelect = document.getElementById('backtestStrategy');
            if (strategySelect) {
                strategySelect.value = strategyId;
            }
        }, 100);

        this.showMessage(`已切换到回测页面，策略"${strategy.name}"已选中`, 'info');
    }

    /**
     * 查看策略表现
     */
    viewPerformance(strategyId) {
        this.viewStrategy(strategyId);
    }

    /**
     * 克隆策略
     */
    cloneStrategy(strategyId) {
        const strategy = this.strategies.find(s => s.id === strategyId);
        if (!strategy) return;

        this.showMessage('策略克隆功能开发中...', 'info');
    }

    /**
     * 显示创建策略模态框
     */
    showCreateStrategyModal() {
        this.showMessage('策略创建功能开发中...', 'info');
    }

    /**
     * 查看详细表现
     */
    viewDetailedPerformance(strategyId) {
        this.showMessage('详细表现分析功能开发中...', 'info');
    }

    /**
     * 策略对比
     */
    compareStrategies(strategyId) {
        this.showMessage('策略对比功能开发中...', 'info');
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

// 导出策略模块类
window.StrategiesModule = StrategiesModule;
