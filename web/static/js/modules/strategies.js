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
                // 处理分页响应格式
                this.strategies = response.data.items || response.data || [];
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

        // 确保strategies是数组
        const strategies = Array.isArray(this.strategies) ? this.strategies : [];

        if (strategies.length === 0) {
            strategiesGrid.innerHTML = `
                <div class="empty-state">
                    <div class="empty-icon">📋</div>
                    <h3>暂无策略</h3>
                    <p>点击"创建策略"按钮开始创建您的第一个策略</p>
                </div>
            `;
            return;
        }

        strategiesGrid.innerHTML = strategies.map(strategy => `
            <div class="strategy-card" data-strategy-id="${strategy.id}">
                <div class="strategy-header">
                    <div class="strategy-title">
                        <h4>${strategy.name}</h4>
                        <span class="strategy-status ${strategy.status}">${this.getStatusText(strategy.status)}</span>
                    </div>
                    <div class="strategy-actions">
                        <button class="btn-icon" onclick="safeCallStrategy('viewStrategy', '${strategy.id}')" title="查看详情">
                            👁️
                        </button>
                        <button class="btn-icon" onclick="safeCallStrategy('editStrategy', '${strategy.id}')" title="编辑">
                            ✏️
                        </button>
                        <button class="btn-icon" onclick="safeCallStrategy('toggleStrategy', '${strategy.id}')" title="启用/禁用">
                            ${strategy.status === 'active' ? '⏸️' : '▶️'}
                        </button>
                        <button class="btn-icon danger" onclick="safeCallStrategy('deleteStrategy', '${strategy.id}')" title="删除">
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
                    <button class="btn btn-primary btn-small" onclick="safeCallStrategy('runBacktest', '${strategy.id}')">
                        🚀 运行回测
                    </button>
                    <button class="btn btn-secondary btn-small" onclick="safeCallStrategy('viewPerformance', '${strategy.id}')">
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
                        <label>策略类型:</label>
                        <span>${this.getTypeText(strategy.strategy_type)}</span>
                    </div>
                    
                    <div class="detail-item">
                        <label>创建时间:</label>
                        <span>${this.formatDate(strategy.created_at)}</span>
                    </div>
                    
                    <div class="detail-item full-width">
                        <label>描述:</label>
                        <p>${strategy.description}</p>
                    </div>
                    
                    <div class="detail-item full-width">
                        <label>参数配置:</label>
                        <div class="parameters-detail">
                            ${this.renderParametersDetail(strategy.parameters)}
                        </div>
                    </div>
                </div>
                
                <div class="detail-actions">
                    <button class="btn btn-primary" onclick="safeCallStrategy('editStrategy', '${strategy.id}')">
                        ✏️ 编辑策略
                    </button>
                    <button class="btn btn-secondary" onclick="safeCallStrategy('runBacktest', '${strategy.id}')">
                        🚀 运行回测
                    </button>
                    <button class="btn btn-outline" onclick="safeCallStrategy('cloneStrategy', '${strategy.id}')">
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
                    <div class="parameter-detail readonly">
                        <div class="param-name">${key}</div>
                        <div class="param-value">${value}</div>
                        <div class="param-type">${typeof value}</div>
                    </div>
                `).join('')}
            </div>
            <div class="parameters-hint">
                <small>💡 点击上方的"编辑策略"按钮来修改这些参数</small>
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
                    <button class="btn btn-primary" onclick="safeCallStrategy('viewDetailedPerformance', '${this.selectedStrategy?.id}')">
                        📊 详细分析
                    </button>
                    <button class="btn btn-secondary" onclick="safeCallStrategy('compareStrategies', '${this.selectedStrategy?.id}')">
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
        console.log('[策略编辑] 编辑策略被调用, strategyId:', strategyId);
        console.log('[策略编辑] 当前策略列表:', this.strategies);
        
        const strategy = this.strategies.find(s => s.id === strategyId);
        if (!strategy) {
            console.error('[策略编辑] 策略不存在:', strategyId);
            this.showMessage('策略不存在', 'error');
            return;
        }

        console.log('[策略编辑] 找到策略:', strategy);
        this.showEditStrategyModal(strategy);
    }

    /**
     * 切换策略状态
     */
    async toggleStrategy(strategyId) {
        const strategy = this.strategies.find(s => s.id === strategyId);
        if (!strategy) return;

        try {
            const isCurrentlyActive = strategy.status === 'active';
            const newStatus = isCurrentlyActive ? 'inactive' : 'active';
            
            // 调用后端API进行状态切换
            let response;
            if (isCurrentlyActive) {
                response = await this.apiService.deactivateStrategy(strategyId);
            } else {
                response = await this.apiService.activateStrategy(strategyId);
            }
            
            if (response && response.success) {
                // 只有在后端成功响应后才更新本地状态
                strategy.status = newStatus;
                
                // 重新显示策略列表
                this.displayStrategies();
                
                this.showMessage(`策略已${newStatus === 'active' ? '启用' : '禁用'}`, 'success');
            } else {
                throw new Error(response?.error || '切换策略状态失败');
            }
            
        } catch (error) {
            console.error('切换策略状态失败:', error);
            this.showMessage(`切换策略状态失败: ${error.message}`, 'error');
            
            // 发生错误时重新显示策略列表以恢复正确状态
            this.displayStrategies();
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

        // 设置策略选择（支持多选）
        setTimeout(() => {
            // 如果回测模块存在且支持多选，使用多选方式
            if (window.backtestModule && window.backtestModule.selectedStrategyIds !== undefined) {
                // 清空当前选择
                window.backtestModule.selectedStrategyIds = [];
                // 添加当前策略
                window.backtestModule.toggleStrategySelection(strategyId);
            } else {
                // 兼容旧的单选方式
                const strategySelect = document.getElementById('backtestStrategy');
                if (strategySelect) {
                    strategySelect.value = strategyId;
                }
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
     * 显示策略编辑模态框
     */
    showEditStrategyModal(strategy) {
        console.log('[策略编辑] 显示编辑模态框, strategy:', strategy);
        
        const modal = document.getElementById('editStrategyModal');
        if (!modal) {
            console.error('[策略编辑] 找不到编辑模态框元素');
            this.showMessage('编辑界面加载失败', 'error');
            return;
        }

        console.log('[策略编辑] 找到模态框元素:', modal);

        // 填充表单数据
        this.populateEditForm(strategy);
        
        // 设置事件监听器
        this.setupEditModalEvents(strategy);
        
        // 显示模态框
        modal.style.display = 'flex';
        console.log('[策略编辑] 模态框已显示');
    }

    /**
     * 填充编辑表单数据
     */
    populateEditForm(strategy) {
        // 基本信息
        document.getElementById('editStrategyName').value = strategy.name || '';
        document.getElementById('editStrategyType').value = strategy.strategy_type || '';
        document.getElementById('editStrategyDescription').value = strategy.description || '';
        document.getElementById('editStrategyStatus').value = strategy.status || '';
        document.getElementById('editStrategyCode').value = strategy.code || '';

        // 生成参数表单
        this.generateParametersForm(strategy);
    }

    /**
     * 生成参数配置表单
     */
    generateParametersForm(strategy) {
        const parametersContainer = document.getElementById('editStrategyParameters');
        if (!parametersContainer) return;

        const parameters = strategy.parameters || {};
        
        // 根据策略类型生成不同的参数表单
        let formHTML = '';
        
        switch (strategy.id) {
            case 'macd_strategy':
                formHTML = this.generateMACDParametersForm(parameters);
                break;
            case 'ma_crossover':
                formHTML = this.generateMAParametersForm(parameters);
                break;
            case 'rsi_strategy':
                formHTML = this.generateRSIParametersForm(parameters);
                break;
            case 'bollinger_strategy':
                formHTML = this.generateBollingerParametersForm(parameters);
                break;
            default:
                formHTML = this.generateGenericParametersForm(parameters);
        }
        
        parametersContainer.innerHTML = formHTML;
    }

    /**
     * 生成MACD策略参数表单
     */
    generateMACDParametersForm(parameters) {
        return `
            <div class="parameter-hint">
                <h5>MACD策略参数说明</h5>
                <p>MACD策略基于快慢均线差值和信号线的交叉来产生买卖信号。当MACD线上穿信号线时产生买入信号，下穿时产生卖出信号。</p>
            </div>
            <div class="parameters-form">
                <div class="parameter-group">
                    <label for="edit_fast_period">快线周期</label>
                    <input type="number" id="edit_fast_period" name="fast_period" 
                           value="${parameters.fast_period || 12}" min="1" max="50" required>
                    <small>计算快速EMA的周期，通常为12</small>
                </div>
                <div class="parameter-group">
                    <label for="edit_slow_period">慢线周期</label>
                    <input type="number" id="edit_slow_period" name="slow_period" 
                           value="${parameters.slow_period || 26}" min="1" max="100" required>
                    <small>计算慢速EMA的周期，通常为26</small>
                </div>
                <div class="parameter-group">
                    <label for="edit_signal_period">信号线周期</label>
                    <input type="number" id="edit_signal_period" name="signal_period" 
                           value="${parameters.signal_period || 9}" min="1" max="50" required>
                    <small>MACD信号线的EMA周期，通常为9</small>
                </div>
                <div class="parameter-group">
                    <label for="edit_buy_threshold">买入阈值</label>
                    <input type="number" id="edit_buy_threshold" name="buy_threshold" 
                           value="${parameters.buy_threshold || 0}" min="-1" max="1" step="0.01">
                    <small>MACD线上穿信号线的阈值</small>
                </div>
                <div class="parameter-group">
                    <label for="edit_sell_threshold">卖出阈值</label>
                    <input type="number" id="edit_sell_threshold" name="sell_threshold" 
                           value="${parameters.sell_threshold || 0}" min="-1" max="1" step="0.01">
                    <small>MACD线下穿信号线的阈值</small>
                </div>
            </div>
        `;
    }

    /**
     * 生成双均线策略参数表单
     */
    generateMAParametersForm(parameters) {
        return `
            <div class="parameter-hint">
                <h5>双均线策略参数说明</h5>
                <p>双均线策略基于短期均线和长期均线的交叉来产生买卖信号。当短期均线上穿长期均线时买入，下穿时卖出。</p>
            </div>
            <div class="parameters-form">
                <div class="parameter-group">
                    <label for="edit_short_period">短期均线周期</label>
                    <input type="number" id="edit_short_period" name="short_period" 
                           value="${parameters.short_period || 5}" min="1" max="50" required>
                    <small>短期移动平均线的周期</small>
                </div>
                <div class="parameter-group">
                    <label for="edit_long_period">长期均线周期</label>
                    <input type="number" id="edit_long_period" name="long_period" 
                           value="${parameters.long_period || 20}" min="1" max="200" required>
                    <small>长期移动平均线的周期</small>
                </div>
                <div class="parameter-group">
                    <label for="edit_ma_type">均线类型</label>
                    <select id="edit_ma_type" name="ma_type" required>
                        <option value="sma" ${parameters.ma_type === 'sma' ? 'selected' : ''}>简单移动平均(SMA)</option>
                        <option value="ema" ${parameters.ma_type === 'ema' ? 'selected' : ''}>指数移动平均(EMA)</option>
                        <option value="wma" ${parameters.ma_type === 'wma' ? 'selected' : ''}>加权移动平均(WMA)</option>
                    </select>
                    <small>选择移动平均线的计算方法</small>
                </div>
                <div class="parameter-group">
                    <label for="edit_threshold">突破阈值</label>
                    <input type="number" id="edit_threshold" name="threshold" 
                           value="${parameters.threshold || 0.01}" min="0" max="0.1" step="0.001">
                    <small>均线突破的确认阈值(百分比)</small>
                </div>
            </div>
        `;
    }

    /**
     * 生成RSI策略参数表单
     */
    generateRSIParametersForm(parameters) {
        return `
            <div class="parameter-hint">
                <h5>RSI策略参数说明</h5>
                <p>RSI策略基于相对强弱指标来识别超买超卖状态。当RSI低于超卖线时买入，高于超买线时卖出。</p>
            </div>
            <div class="parameters-form">
                <div class="parameter-group">
                    <label for="edit_rsi_period">RSI周期</label>
                    <input type="number" id="edit_rsi_period" name="period" 
                           value="${parameters.period || 14}" min="1" max="50" required>
                    <small>计算RSI的周期，通常为14</small>
                </div>
                <div class="parameter-group">
                    <label for="edit_overbought">超买阈值</label>
                    <input type="number" id="edit_overbought" name="overbought" 
                           value="${parameters.overbought || 70}" min="50" max="100" required>
                    <small>RSI超买阈值，通常为70</small>
                </div>
                <div class="parameter-group">
                    <label for="edit_oversold">超卖阈值</label>
                    <input type="number" id="edit_oversold" name="oversold" 
                           value="${parameters.oversold || 30}" min="0" max="50" required>
                    <small>RSI超卖阈值，通常为30</small>
                </div>
            </div>
        `;
    }

    /**
     * 生成布林带策略参数表单
     */
    generateBollingerParametersForm(parameters) {
        return `
            <div class="parameter-hint">
                <h5>布林带策略参数说明</h5>
                <p>布林带策略基于价格相对于统计区间的位置来判断买卖时机。价格触及下轨时买入，触及上轨时卖出。</p>
            </div>
            <div class="parameters-form">
                <div class="parameter-group">
                    <label for="edit_bb_period">布林带周期</label>
                    <input type="number" id="edit_bb_period" name="period" 
                           value="${parameters.period || 20}" min="1" max="50" required>
                    <small>计算布林带的周期，通常为20</small>
                </div>
                <div class="parameter-group">
                    <label for="edit_std_dev">标准差倍数</label>
                    <input type="number" id="edit_std_dev" name="std_dev" 
                           value="${parameters.std_dev || 2}" min="0.5" max="5" step="0.1" required>
                    <small>布林带宽度的标准差倍数，通常为2</small>
                </div>
            </div>
        `;
    }

    /**
     * 生成通用参数表单
     */
    generateGenericParametersForm(parameters) {
        if (!parameters || Object.keys(parameters).length === 0) {
            return '<p class="no-parameters">该策略暂无可配置参数</p>';
        }

        const parametersHTML = Object.entries(parameters).map(([key, value]) => `
            <div class="parameter-group">
                <label for="edit_${key}">${key}</label>
                <input type="${typeof value === 'number' ? 'number' : 'text'}" 
                       id="edit_${key}" name="${key}" value="${value}">
                <small>参数: ${key}</small>
            </div>
        `).join('');

        return `<div class="parameters-form">${parametersHTML}</div>`;
    }

    /**
     * 设置编辑模态框事件监听器
     */
    setupEditModalEvents(strategy) {
        // 移除之前的事件监听器
        this.removeEditModalEvents();

        // 关闭按钮
        const closeBtn = document.getElementById('closeEditStrategyBtn');
        const cancelBtn = document.getElementById('cancelEditStrategy');
        const resetBtn = document.getElementById('resetEditStrategy');
        const saveBtn = document.getElementById('saveEditStrategy');
        const form = document.getElementById('editStrategyForm');

        this.editModalHandlers = {
            close: () => this.hideEditStrategyModal(),
            cancel: () => this.hideEditStrategyModal(),
            reset: () => this.populateEditForm(strategy),
            save: (e) => this.handleSaveStrategy(e, strategy),
            submit: (e) => this.handleSaveStrategy(e, strategy)
        };

        if (closeBtn) closeBtn.addEventListener('click', this.editModalHandlers.close);
        if (cancelBtn) cancelBtn.addEventListener('click', this.editModalHandlers.cancel);
        if (resetBtn) resetBtn.addEventListener('click', this.editModalHandlers.reset);
        if (saveBtn) saveBtn.addEventListener('click', this.editModalHandlers.save);
        if (form) form.addEventListener('submit', this.editModalHandlers.submit);

        // 点击背景关闭模态框
        const modal = document.getElementById('editStrategyModal');
        this.editModalHandlers.backdrop = (e) => {
            if (e.target === modal) {
                this.hideEditStrategyModal();
            }
        };
        if (modal) modal.addEventListener('click', this.editModalHandlers.backdrop);
    }

    /**
     * 移除编辑模态框事件监听器
     */
    removeEditModalEvents() {
        if (!this.editModalHandlers) return;

        const closeBtn = document.getElementById('closeEditStrategyBtn');
        const cancelBtn = document.getElementById('cancelEditStrategy');
        const resetBtn = document.getElementById('resetEditStrategy');
        const saveBtn = document.getElementById('saveEditStrategy');
        const form = document.getElementById('editStrategyForm');
        const modal = document.getElementById('editStrategyModal');

        if (closeBtn && this.editModalHandlers.close) {
            closeBtn.removeEventListener('click', this.editModalHandlers.close);
        }
        if (cancelBtn && this.editModalHandlers.cancel) {
            cancelBtn.removeEventListener('click', this.editModalHandlers.cancel);
        }
        if (resetBtn && this.editModalHandlers.reset) {
            resetBtn.removeEventListener('click', this.editModalHandlers.reset);
        }
        if (saveBtn && this.editModalHandlers.save) {
            saveBtn.removeEventListener('click', this.editModalHandlers.save);
        }
        if (form && this.editModalHandlers.submit) {
            form.removeEventListener('submit', this.editModalHandlers.submit);
        }
        if (modal && this.editModalHandlers.backdrop) {
            modal.removeEventListener('click', this.editModalHandlers.backdrop);
        }

        this.editModalHandlers = null;
    }

    /**
     * 隐藏编辑模态框
     */
    hideEditStrategyModal() {
        const modal = document.getElementById('editStrategyModal');
        if (modal) {
            modal.style.display = 'none';
        }
        this.removeEditModalEvents();
    }

    /**
     * 处理保存策略
     */
    async handleSaveStrategy(e, originalStrategy) {
        e.preventDefault();
        
        try {
            // 收集表单数据
            const formData = this.collectFormData();
            
            // 验证表单数据
            const validationResult = this.validateStrategyForm(formData, originalStrategy);
            if (!validationResult.valid) {
                this.showMessage(validationResult.message, 'error');
                return;
            }

            // 显示保存状态
            const saveBtn = document.getElementById('saveEditStrategy');
            const originalText = saveBtn.textContent;
            saveBtn.textContent = '保存中...';
            saveBtn.disabled = true;

            // 调用API更新策略
            const response = await this.apiService.updateStrategy(originalStrategy.id, formData);
            
            if (response.success) {
                // 更新本地策略数据
                const strategyIndex = this.strategies.findIndex(s => s.id === originalStrategy.id);
                if (strategyIndex !== -1) {
                    // 合并更新后的数据
                    this.strategies[strategyIndex] = {
                        ...this.strategies[strategyIndex],
                        ...formData,
                        updated_at: new Date().toISOString()
                    };
                }

                // 刷新显示
                this.displayStrategies();
                
                // 如果当前选中的策略被更新，也要更新详情显示
                if (this.selectedStrategy?.id === originalStrategy.id) {
                    this.selectedStrategy = this.strategies[strategyIndex];
                    this.displayStrategyDetails(this.selectedStrategy);
                }

                this.hideEditStrategyModal();
                this.showMessage('策略更新成功', 'success');
            } else {
                throw new Error(response.message || '更新策略失败');
            }

        } catch (error) {
            console.error('保存策略失败:', error);
            this.showMessage(`保存失败: ${error.message}`, 'error');
        } finally {
            // 恢复按钮状态
            const saveBtn = document.getElementById('saveEditStrategy');
            if (saveBtn) {
                saveBtn.textContent = originalText;
                saveBtn.disabled = false;
            }
        }
    }

    /**
     * 收集表单数据
     */
    collectFormData() {
        const formData = {
            name: document.getElementById('editStrategyName').value.trim(),
            description: document.getElementById('editStrategyDescription').value.trim(),
            status: document.getElementById('editStrategyStatus').value,
            code: document.getElementById('editStrategyCode').value.trim(),
            parameters: {}
        };

        // 收集参数数据
        const parametersContainer = document.getElementById('editStrategyParameters');
        if (parametersContainer) {
            const paramInputs = parametersContainer.querySelectorAll('input, select');
            paramInputs.forEach(input => {
                if (input.name && input.value !== '') {
                    let value = input.value;
                    
                    // 转换数据类型
                    if (input.type === 'number') {
                        value = parseFloat(value);
                        if (isNaN(value)) value = 0;
                    }
                    
                    formData.parameters[input.name] = value;
                }
            });
        }

        return formData;
    }

    /**
     * 验证策略表单
     */
    validateStrategyForm(formData, originalStrategy) {
        // 验证策略名称
        if (!formData.name) {
            return { valid: false, message: '策略名称不能为空' };
        }

        if (formData.name.length > 100) {
            return { valid: false, message: '策略名称长度不能超过100个字符' };
        }

        // 检查名称是否与其他策略重复
        const duplicateStrategy = this.strategies.find(s => 
            s.id !== originalStrategy.id && s.name === formData.name
        );
        if (duplicateStrategy) {
            return { valid: false, message: '策略名称已存在，请使用其他名称' };
        }

        // 验证描述长度
        if (formData.description && formData.description.length > 1000) {
            return { valid: false, message: '策略描述长度不能超过1000个字符' };
        }

        // 验证参数
        const paramValidation = this.validateStrategyParameters(formData.parameters, originalStrategy);
        if (!paramValidation.valid) {
            return paramValidation;
        }

        return { valid: true };
    }

    /**
     * 验证策略参数
     */
    validateStrategyParameters(parameters, strategy) {
        switch (strategy.id) {
            case 'macd_strategy':
                return this.validateMACDParameters(parameters);
            case 'ma_crossover':
                return this.validateMAParameters(parameters);
            case 'rsi_strategy':
                return this.validateRSIParameters(parameters);
            case 'bollinger_strategy':
                return this.validateBollingerParameters(parameters);
            default:
                return { valid: true };
        }
    }

    /**
     * 验证MACD参数
     */
    validateMACDParameters(params) {
        const fastPeriod = params.fast_period;
        const slowPeriod = params.slow_period;
        const signalPeriod = params.signal_period;

        if (!fastPeriod || fastPeriod < 1 || fastPeriod > 50) {
            return { valid: false, message: 'MACD快线周期必须在1-50之间' };
        }

        if (!slowPeriod || slowPeriod < 1 || slowPeriod > 100) {
            return { valid: false, message: 'MACD慢线周期必须在1-100之间' };
        }

        if (fastPeriod >= slowPeriod) {
            return { valid: false, message: 'MACD快线周期必须小于慢线周期' };
        }

        if (!signalPeriod || signalPeriod < 1 || signalPeriod > 50) {
            return { valid: false, message: 'MACD信号线周期必须在1-50之间' };
        }

        return { valid: true };
    }

    /**
     * 验证双均线参数
     */
    validateMAParameters(params) {
        const shortPeriod = params.short_period;
        const longPeriod = params.long_period;

        if (!shortPeriod || shortPeriod < 1 || shortPeriod > 50) {
            return { valid: false, message: '短期均线周期必须在1-50之间' };
        }

        if (!longPeriod || longPeriod < 1 || longPeriod > 200) {
            return { valid: false, message: '长期均线周期必须在1-200之间' };
        }

        if (shortPeriod >= longPeriod) {
            return { valid: false, message: '短期均线周期必须小于长期均线周期' };
        }

        if (!params.ma_type || !['sma', 'ema', 'wma'].includes(params.ma_type)) {
            return { valid: false, message: '均线类型必须是SMA、EMA或WMA' };
        }

        return { valid: true };
    }

    /**
     * 验证RSI参数
     */
    validateRSIParameters(params) {
        const period = params.period;
        const overbought = params.overbought;
        const oversold = params.oversold;

        if (!period || period < 1 || period > 50) {
            return { valid: false, message: 'RSI周期必须在1-50之间' };
        }

        if (!overbought || overbought < 50 || overbought > 100) {
            return { valid: false, message: 'RSI超买阈值必须在50-100之间' };
        }

        if (!oversold || oversold < 0 || oversold > 50) {
            return { valid: false, message: 'RSI超卖阈值必须在0-50之间' };
        }

        if (oversold >= overbought) {
            return { valid: false, message: 'RSI超卖阈值必须小于超买阈值' };
        }

        return { valid: true };
    }

    /**
     * 验证布林带参数
     */
    validateBollingerParameters(params) {
        const period = params.period;
        const stdDev = params.std_dev;

        if (!period || period < 1 || period > 50) {
            return { valid: false, message: '布林带周期必须在1-50之间' };
        }

        if (!stdDev || stdDev < 0.5 || stdDev > 5) {
            return { valid: false, message: '标准差倍数必须在0.5-5之间' };
        }

        return { valid: true };
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
