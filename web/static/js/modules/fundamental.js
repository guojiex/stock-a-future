/**
 * 基本面数据处理模块
 */

class FundamentalModule {
    constructor(apiService) {
        this.apiService = apiService;
        this.currentStockCode = null;
        this.init();
    }

    init() {
        this.bindEvents();
        console.log('[Fundamental] 基本面模块已初始化');
    }

    bindEvents() {
        const loadBtn = document.getElementById('loadFundamentalBtn');
        const queryBtn = document.getElementById('queryFundamental');
        
        if (loadBtn) {
            loadBtn.addEventListener('click', () => this.loadFundamentalData());
        }
        if (queryBtn) {
            queryBtn.addEventListener('click', () => this.handleFundamentalQuery());
        }
    }

    async handleFundamentalQuery() {
        const stockCode = document.getElementById('stockCode').value.trim();
        if (!stockCode) {
            alert('请输入股票代码');
            return;
        }

        this.switchToFundamentalTab();
        this.currentStockCode = stockCode;
        await this.loadFundamentalData();
    }

    switchToFundamentalTab() {
        const section = document.getElementById('daily-chart-section');
        if (section) section.style.display = 'block';

        document.querySelectorAll('.tab-btn').forEach(tab => tab.classList.remove('active'));
        document.querySelectorAll('.tab-pane').forEach(pane => pane.classList.remove('active'));

        const tab = document.querySelector('[data-tab="fundamental"]');
        const pane = document.getElementById('fundamental-tab');
        
        if (tab && pane) {
            tab.classList.add('active');
            pane.classList.add('active');
        }
    }

    async loadFundamentalData() {
        if (!this.currentStockCode) {
            this.currentStockCode = document.getElementById('stockCode').value.trim();
        }
        if (!this.currentStockCode) {
            this.showError('请先选择股票代码');
            return;
        }

        try {
            this.showLoading(true);
            this.hideError();

            const period = document.getElementById('reportPeriod').value;
            const reportType = document.getElementById('reportType').value;

            const data = await this.apiService.getFundamentalData(
                this.currentStockCode, period, reportType
            );

            this.displayFundamentalData(data);
            this.showLoading(false);
        } catch (error) {
            console.error('[Fundamental] 加载失败:', error);
            this.showError(`加载失败: ${error.message}`);
            this.showLoading(false);
        }
    }

    displayFundamentalData(data) {
        const container = document.getElementById('fundamentalData');
        if (container) container.style.display = 'block';

        if (data.stock_basic) this.displayStockBasicInfo(data.stock_basic);
        if (data.daily_basic) this.displayDailyBasicInfo(data.daily_basic);
        if (data.cash_flow_statement) this.displayCashFlowStatement(data.cash_flow_statement);
        if (data.income_statement) this.displayIncomeStatement(data.income_statement);
        if (data.balance_sheet) this.displayBalanceSheet(data.balance_sheet);
        
        // 自动加载基本面因子分析
        this.loadFactorAnalysis();
    }

    displayStockBasicInfo(stockBasic) {
        const container = document.getElementById('stockBasicInfo');
        if (!container) return;

        const items = [
            { label: '股票代码', value: stockBasic.ts_code || stockBasic.symbol || '-' },
            { label: '股票名称', value: stockBasic.name || '-' },
            { label: '所属行业', value: stockBasic.industry || '-' },
            { label: '上市日期', value: stockBasic.list_date || '-' }
        ];

        container.innerHTML = items.map(item => `
            <div class="info-item">
                <span class="info-label">${item.label}:</span>
                <span class="info-value">${item.value}</span>
            </div>
        `).join('');
    }

    displayDailyBasicInfo(dailyBasic) {
        const container = document.getElementById('dailyBasicInfo');
        if (!container) return;

        const items = [
            { label: '最新价', value: dailyBasic.close ? `${dailyBasic.close}元` : '-' },
            { label: '市盈率', value: dailyBasic.pe || dailyBasic.pe_ttm || '-' },
            { label: '市净率', value: dailyBasic.pb || '-' },
            { label: '总市值', value: this.formatNumber(dailyBasic.total_mv) }
        ];

        container.innerHTML = items.map(item => `
            <div class="info-item">
                <span class="info-label">${item.label}:</span>
                <span class="info-value">${item.value}</span>
            </div>
        `).join('');
    }

    formatNumber(value) {
        if (!value || isNaN(value)) return '-';
        const num = parseFloat(value);
        if (num >= 100000000) return (num / 100000000).toFixed(2) + '亿';
        if (num >= 10000) return (num / 10000).toFixed(2) + '万';
        return num.toFixed(2);
    }

    showLoading(show) {
        const loading = document.getElementById('fundamentalLoading');
        if (loading) loading.style.display = show ? 'block' : 'none';
    }

    showError(message) {
        const errorDiv = document.getElementById('fundamentalError');
        if (errorDiv) {
            errorDiv.textContent = message;
            errorDiv.style.display = 'block';
        }
    }

    hideError() {
        const errorDiv = document.getElementById('fundamentalError');
        if (errorDiv) errorDiv.style.display = 'none';
    }

    // ===== 基本面因子分析功能 =====

    async loadFactorAnalysis() {
        if (!this.currentStockCode) return;

        try {
            console.log('[Fundamental] 开始加载因子分析:', this.currentStockCode);
            
            // 获取基本面因子
            const factor = await this.apiService.getFundamentalFactor(this.currentStockCode);
            
            // 显示因子分析
            this.displayFactorAnalysis(factor);
            
        } catch (error) {
            console.error('[Fundamental] 因子分析加载失败:', error);
            // 不显示错误，因为这是额外功能
        }
    }

    displayFactorAnalysis(factor) {
        // 创建因子分析容器（如果不存在）
        this.createFactorAnalysisContainer();
        
        // 显示因子得分
        this.displayFactorScores(factor);
        
        // 显示因子雷达图
        this.displayFactorRadarChart(factor);
        
        // 显示详细因子数据
        this.displayDetailedFactors(factor);
    }

    createFactorAnalysisContainer() {
        const fundamentalData = document.getElementById('fundamentalData');
        if (!fundamentalData) return;

        // 检查是否已存在因子分析容器
        let factorContainer = document.getElementById('factorAnalysisContainer');
        if (factorContainer) return;

        // 创建因子分析容器
        const factorHTML = `
            <div class="fundamental-section" id="factorAnalysisContainer">
                <h4 class="section-title">📊 基本面因子分析</h4>
                
                <!-- 因子得分概览 -->
                <div class="factor-scores-overview" id="factorScoresOverview">
                    <div class="score-cards">
                        <div class="score-card value-score">
                            <div class="score-label">价值因子</div>
                            <div class="score-value" id="valueScoreValue">-</div>
                            <div class="score-rank" id="valueScoreRank">-</div>
                        </div>
                        <div class="score-card growth-score">
                            <div class="score-label">成长因子</div>
                            <div class="score-value" id="growthScoreValue">-</div>
                            <div class="score-rank" id="growthScoreRank">-</div>
                        </div>
                        <div class="score-card quality-score">
                            <div class="score-label">质量因子</div>
                            <div class="score-value" id="qualityScoreValue">-</div>
                            <div class="score-rank" id="qualityScoreRank">-</div>
                        </div>
                        <div class="score-card profitability-score">
                            <div class="score-label">盈利因子</div>
                            <div class="score-value" id="profitabilityScoreValue">-</div>
                            <div class="score-rank" id="profitabilityScoreRank">-</div>
                        </div>
                        <div class="score-card composite-score">
                            <div class="score-label">综合得分</div>
                            <div class="score-value" id="compositeScoreValue">-</div>
                            <div class="score-rank" id="compositeScoreRank">-</div>
                        </div>
                    </div>
                </div>

                <!-- 因子雷达图 -->
                <div class="factor-radar-container">
                    <div id="factorRadarChart" style="width: 100%; height: 400px;"></div>
                </div>

                <!-- 详细因子数据 -->
                <div class="factor-details">
                    <div class="factor-tabs">
                        <button class="factor-tab-btn active" data-factor="value">价值因子</button>
                        <button class="factor-tab-btn" data-factor="growth">成长因子</button>
                        <button class="factor-tab-btn" data-factor="quality">质量因子</button>
                        <button class="factor-tab-btn" data-factor="profitability">盈利因子</button>
                    </div>
                    <div class="factor-content" id="factorContent">
                        <!-- 动态内容 -->
                    </div>
                </div>

                <!-- 因子排名按钮 -->
                <div class="factor-actions">
                    <button id="viewFactorRankingBtn" class="btn btn-secondary">📈 查看市场因子排名</button>
                </div>
            </div>
        `;

        fundamentalData.insertAdjacentHTML('beforeend', factorHTML);
        
        // 绑定因子标签页事件
        this.bindFactorTabEvents();
        
        // 绑定排名按钮事件
        this.bindFactorRankingEvents();
    }

    displayFactorScores(factor) {
        // 更新因子得分卡片
        this.updateScoreCard('valueScoreValue', 'valueScoreRank', factor.value_score, factor.market_percentile);
        this.updateScoreCard('growthScoreValue', 'growthScoreRank', factor.growth_score, factor.market_percentile);
        this.updateScoreCard('qualityScoreValue', 'qualityScoreRank', factor.quality_score, factor.market_percentile);
        this.updateScoreCard('profitabilityScoreValue', 'profitabilityScoreRank', factor.profitability_score, factor.market_percentile);
        this.updateScoreCard('compositeScoreValue', 'compositeScoreRank', factor.composite_score, factor.market_percentile);
    }

    updateScoreCard(valueId, rankId, score, percentile) {
        const valueElement = document.getElementById(valueId);
        const rankElement = document.getElementById(rankId);
        
        if (valueElement) {
            const scoreValue = score ? parseFloat(score).toFixed(2) : '-';
            valueElement.textContent = scoreValue;
            
            // 根据得分设置颜色
            const numScore = parseFloat(score) || 0;
            if (numScore > 1) {
                valueElement.className = 'score-value positive';
            } else if (numScore < -1) {
                valueElement.className = 'score-value negative';
            } else {
                valueElement.className = 'score-value neutral';
            }
        }
        
        if (rankElement) {
            const percentileValue = percentile ? parseFloat(percentile).toFixed(1) : '-';
            rankElement.textContent = percentileValue !== '-' ? `${percentileValue}%` : '-';
        }
    }

    displayFactorRadarChart(factor) {
        // 检查是否有ECharts
        if (typeof echarts === 'undefined') {
            console.warn('[Fundamental] ECharts未加载，跳过雷达图显示');
            return;
        }

        const chartContainer = document.getElementById('factorRadarChart');
        if (!chartContainer) return;

        const chart = echarts.init(chartContainer);
        
        const option = {
            title: {
                text: '基本面因子雷达图',
                left: 'center',
                textStyle: {
                    fontSize: 16,
                    fontWeight: 'bold'
                }
            },
            tooltip: {
                trigger: 'item'
            },
            radar: {
                indicator: [
                    { name: '价值因子', max: 3, min: -3 },
                    { name: '成长因子', max: 3, min: -3 },
                    { name: '质量因子', max: 3, min: -3 },
                    { name: '盈利因子', max: 3, min: -3 }
                ],
                shape: 'polygon',
                splitNumber: 6,
                axisName: {
                    color: '#666'
                },
                splitLine: {
                    lineStyle: {
                        color: '#e6e6e6'
                    }
                },
                splitArea: {
                    show: true,
                    areaStyle: {
                        color: ['rgba(114, 172, 209, 0.1)', 'rgba(255, 255, 255, 0.1)']
                    }
                }
            },
            series: [{
                name: '因子得分',
                type: 'radar',
                data: [{
                    value: [
                        factor.value_score || 0,
                        factor.growth_score || 0,
                        factor.quality_score || 0,
                        factor.profitability_score || 0
                    ],
                    name: '当前股票',
                    areaStyle: {
                        color: 'rgba(54, 162, 235, 0.3)'
                    },
                    lineStyle: {
                        color: '#36A2EB',
                        width: 2
                    },
                    symbol: 'circle',
                    symbolSize: 6
                }]
            }]
        };

        chart.setOption(option);
        
        // 响应式调整
        window.addEventListener('resize', () => {
            chart.resize();
        });
    }

    displayDetailedFactors(factor) {
        // 默认显示价值因子
        this.showFactorDetails('value', factor);
    }

    showFactorDetails(factorType, factor) {
        const content = document.getElementById('factorContent');
        if (!content) return;

        let html = '';
        
        switch (factorType) {
            case 'value':
                html = this.generateValueFactorHTML(factor);
                break;
            case 'growth':
                html = this.generateGrowthFactorHTML(factor);
                break;
            case 'quality':
                html = this.generateQualityFactorHTML(factor);
                break;
            case 'profitability':
                html = this.generateProfitabilityFactorHTML(factor);
                break;
        }
        
        content.innerHTML = html;
    }

    generateValueFactorHTML(factor) {
        const items = [
            { label: '市盈率 (PE)', value: factor.pe, desc: '股价相对每股收益的倍数，越低越好' },
            { label: '市净率 (PB)', value: factor.pb, desc: '股价相对每股净资产的倍数，越低越好' },
            { label: '市销率 (PS)', value: factor.ps, desc: '股价相对每股销售收入的倍数，越低越好' },
            { label: '市现率 (PCF)', value: factor.pcf, desc: '股价相对每股现金流的倍数，越低越好' }
        ];

        return this.generateFactorTable(items);
    }

    generateGrowthFactorHTML(factor) {
        const items = [
            { label: '营收增长率', value: factor.revenue_growth, desc: '营业收入同比增长率，越高越好', unit: '%' },
            { label: '净利润增长率', value: factor.net_profit_growth, desc: '净利润同比增长率，越高越好', unit: '%' },
            { label: 'EPS增长率', value: factor.eps_growth, desc: '每股收益同比增长率，越高越好', unit: '%' },
            { label: 'ROE增长率', value: factor.roe_growth, desc: '净资产收益率同比增长率，越高越好', unit: '%' }
        ];

        return this.generateFactorTable(items);
    }

    generateQualityFactorHTML(factor) {
        const items = [
            { label: '净资产收益率 (ROE)', value: factor.roe, desc: '净利润与净资产的比率，越高越好', unit: '%' },
            { label: '资产收益率 (ROA)', value: factor.roa, desc: '净利润与总资产的比率，越高越好', unit: '%' },
            { label: '资产负债率', value: factor.debt_to_assets, desc: '总负债与总资产的比率，越低越好', unit: '%' },
            { label: '流动比率', value: factor.current_ratio, desc: '流动资产与流动负债的比率，适中为好' }
        ];

        return this.generateFactorTable(items);
    }

    generateProfitabilityFactorHTML(factor) {
        const items = [
            { label: '毛利率', value: factor.gross_margin, desc: '毛利润与营业收入的比率，越高越好', unit: '%' },
            { label: '净利率', value: factor.net_margin, desc: '净利润与营业收入的比率，越高越好', unit: '%' },
            { label: '营业利润率', value: factor.operating_margin, desc: '营业利润与营业收入的比率，越高越好', unit: '%' },
            { label: 'ROIC', value: factor.roic, desc: '投入资本回报率，越高越好', unit: '%' }
        ];

        return this.generateFactorTable(items);
    }

    generateFactorTable(items) {
        return `
            <table class="fundamental-table">
                <thead>
                    <tr>
                        <th>指标名称</th>
                        <th>数值</th>
                        <th>说明</th>
                    </tr>
                </thead>
                <tbody>
                    ${items.map(item => `
                        <tr>
                            <td class="item-name">${item.label}</td>
                            <td class="item-value">${this.formatFactorValue(item.value, item.unit)}</td>
                            <td class="item-desc">${item.desc}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    }

    formatFactorValue(value, unit = '') {
        if (!value || isNaN(value)) return '-';
        const num = parseFloat(value);
        return `${num.toFixed(2)}${unit}`;
    }

    bindFactorTabEvents() {
        const tabBtns = document.querySelectorAll('.factor-tab-btn');
        tabBtns.forEach(btn => {
            btn.addEventListener('click', (e) => {
                // 更新标签页状态
                tabBtns.forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                
                // 显示对应内容
                const factorType = e.target.dataset.factor;
                // 这里需要重新获取factor数据，简化处理
                console.log(`[Fundamental] 切换到${factorType}因子`);
            });
        });
    }

    bindFactorRankingEvents() {
        const rankingBtn = document.getElementById('viewFactorRankingBtn');
        if (rankingBtn) {
            rankingBtn.addEventListener('click', () => this.showFactorRanking());
        }
    }

    async showFactorRanking() {
        try {
            console.log('[Fundamental] 获取因子排名...');
            
            const ranking = await this.apiService.getFundamentalFactorRanking('composite', '', 20);
            
            this.displayFactorRanking(ranking);
            
        } catch (error) {
            console.error('[Fundamental] 获取因子排名失败:', error);
            alert('获取因子排名失败: ' + error.message);
        }
    }

    displayFactorRanking(ranking) {
        // 创建模态框显示排名
        const modalHTML = `
            <div class="modal-overlay" id="factorRankingModal">
                <div class="modal-content">
                    <div class="modal-header">
                        <h3>📈 基本面因子排名 (${ranking.factor_type})</h3>
                        <button class="modal-close" onclick="this.closest('.modal-overlay').remove()">×</button>
                    </div>
                    <div class="modal-body">
                        <p>数据日期: ${ranking.trade_date} | 总计: ${ranking.total}个股票</p>
                        <table class="ranking-table">
                            <thead>
                                <tr>
                                    <th>排名</th>
                                    <th>股票代码</th>
                                    <th>综合得分</th>
                                    <th>市场分位数</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${ranking.factors.map((factor, index) => `
                                    <tr>
                                        <td>${index + 1}</td>
                                        <td>${factor.ts_code}</td>
                                        <td>${this.formatFactorValue(factor.composite_score)}</td>
                                        <td>${this.formatFactorValue(factor.market_percentile)}%</td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        `;
        
        document.body.insertAdjacentHTML('beforeend', modalHTML);
    }

    // 通用的金额格式化函数
    formatAmount(value) {
        if (!value || value === '0' || value === 0) return '0';
        const num = parseFloat(value);
        if (isNaN(num)) return '0';
        
        // 如果金额大于10亿，显示为亿元
        if (Math.abs(num) >= 1000000000) {
            const yi = (num / 100000000).toFixed(2);
            return `${yi}亿元`;
        }
        // 如果金额大于1万，显示为万元
        else if (Math.abs(num) >= 10000) {
            const wan = (num / 10000).toFixed(0);
            return `${wan}万元`;
        }
        // 小于1万，显示为元
        else {
            return `${num.toFixed(0)}元`;
        }
    }

    displayCashFlowStatement(cashFlow) {
        console.log('[Fundamental] 显示现金流量表数据:', cashFlow);
        
        const container = document.getElementById('cashFlowInfo');
        if (!container) {
            console.warn('[Fundamental] 未找到现金流量表容器 #cashFlowInfo');
            return;
        }

        const items = [
            { 
                label: '经营现金流', 
                value: this.formatAmount(cashFlow.net_cash_oper_act),
                desc: '经营活动产生的现金流量净额'
            },
            { 
                label: '投资现金流', 
                value: this.formatAmount(cashFlow.net_cash_inv_act),
                desc: '投资活动产生的现金流量净额'
            },
            { 
                label: '筹资现金流', 
                value: this.formatAmount(cashFlow.net_cash_fin_act),
                desc: '筹资活动产生的现金流量净额'
            },
            { 
                label: '期间', 
                value: cashFlow.end_date || cashFlow.ann_date || '-',
                desc: '报告期间'
            }
        ];

        container.innerHTML = `
            <div class="financial-statement">
                <h4 class="statement-title">💰 现金流量表</h4>
                <div class="statement-items">
                    ${items.map(item => `
                        <div class="info-item" title="${item.desc}">
                            <span class="info-label">${item.label}:</span>
                            <span class="info-value ${item.label.includes('现金流') ? 'cash-flow-value' : ''}">${item.value}</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
        
        container.style.display = 'block';
    }

    displayIncomeStatement(income) {
        console.log('[Fundamental] 显示利润表数据:', income);
        
        const container = document.getElementById('incomeStatementInfo');
        if (!container) {
            console.warn('[Fundamental] 未找到利润表容器 #incomeStatementInfo');
            return;
        }

        const items = [
            { 
                label: '营业收入', 
                value: this.formatAmount(income.total_revenue),
                desc: '营业收入总额'
            },
            { 
                label: '净利润', 
                value: this.formatAmount(income.n_income),
                desc: '净利润'
            },
            { 
                label: '营业利润', 
                value: this.formatAmount(income.operate_profit),
                desc: '营业利润'
            },
            { 
                label: '期间', 
                value: income.end_date || income.ann_date || '-',
                desc: '报告期间'
            }
        ];

        container.innerHTML = `
            <div class="financial-statement">
                <h4 class="statement-title">📊 利润表</h4>
                <div class="statement-items">
                    ${items.map(item => `
                        <div class="info-item" title="${item.desc}">
                            <span class="info-label">${item.label}:</span>
                            <span class="info-value">${item.value}</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
        
        container.style.display = 'block';
    }

    displayBalanceSheet(balance) {
        console.log('[Fundamental] 显示资产负债表数据:', balance);
        console.log('[Fundamental] 净资产字段值:', balance.total_hldr_eqy);
        console.log('[Fundamental] 总资产字段值:', balance.total_assets);
        console.log('[Fundamental] 总负债字段值:', balance.total_liab);
        
        const container = document.getElementById('balanceSheetInfo');
        if (!container) {
            console.warn('[Fundamental] 未找到资产负债表容器 #balanceSheetInfo');
            return;
        }

        const items = [
            { 
                label: '总资产', 
                value: this.formatAmount(balance.total_assets),
                desc: '资产总计'
            },
            { 
                label: '净资产', 
                value: this.formatAmount(balance.total_hldr_eqy),
                desc: '所有者权益合计'
            },
            { 
                label: '总负债', 
                value: this.formatAmount(balance.total_liab),
                desc: '负债合计'
            },
            { 
                label: '期间', 
                value: balance.end_date || balance.ann_date || '-',
                desc: '报告期间'
            }
        ];

        container.innerHTML = `
            <div class="financial-statement">
                <h4 class="statement-title">🏦 资产负债表</h4>
                <div class="statement-items">
                    ${items.map(item => `
                        <div class="info-item" title="${item.desc}">
                            <span class="info-label">${item.label}:</span>
                            <span class="info-value">${item.value}</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
        
        container.style.display = 'block';
    }
}

window.FundamentalModule = FundamentalModule;
