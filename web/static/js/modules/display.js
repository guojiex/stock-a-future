/**
 * 数据展示模块
 * 负责各种数据的展示和格式化
 */

class DisplayModule {
    constructor(client, chartsModule) {
        this.client = client;
        this.chartsModule = chartsModule;
    }

    /**
     * 显示日线数据
     */
    displayDailyData(data, stockCode, stockBasic) {
        console.log(`[Display] 显示日线数据 - 数据:`, data);
        
        const section = document.getElementById('daily-chart-section');
        const card = document.getElementById('dailyChartCard');
        const summary = document.getElementById('dailyDataSummary');
        
        if (!section || !card || !summary) {
            console.error(`[Display] 找不到必要的DOM元素:`, { section: !!section, card: !!card, summary: !!summary });
            return;
        }
        
        // 显示section
        section.style.display = 'block';
        section.classList.add('fade-in');
        
        // 更新股票分析标题
        this.updateStockAnalysisTitle(stockCode, stockBasic);
        
        // 切换到日线数据tab（不触发数据加载）
        this.switchToTabWithoutDataLoad('daily-data');
        
        // 创建价格图表
        this.chartsModule.createPriceChart(data, stockCode, stockBasic);
        
        // 显示数据摘要
        if (data && data.length > 0) {
            const latest = data[data.length - 1];
            const previous = data.length > 1 ? data[data.length - 2] : latest;
            
            summary.innerHTML = this.createDataSummary(latest, previous);
        }
        
        // 滚动到结果
        section.scrollIntoView({ behavior: 'smooth' });
        
        console.log(`[Display] 日线数据显示完成`);
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
     * 显示技术指标
     */
    displayIndicators(data, stockCode, stockBasic) {
        console.log(`[Display] 显示技术指标 - 数据:`, data);
        
        const section = document.getElementById('daily-chart-section');
        const indicatorsGrid = document.getElementById('indicatorsGrid');
        
        if (!section || !indicatorsGrid) {
            console.error(`[Display] 找不到必要的DOM元素:`, { section: !!section, indicatorsGrid: !!indicatorsGrid });
            return;
        }
        
        // 显示section
        section.style.display = 'block';
        section.classList.add('fade-in');
        
        // 更新股票分析标题
        this.updateStockAnalysisTitle(stockCode, stockBasic);
        
        // 切换到技术指标tab（不触发数据加载）
        this.switchToTabWithoutDataLoad('indicators');
        
        // 显示技术指标数据 - 修复数据检查逻辑
        if (data && (Array.isArray(data) ? data.length > 0 : Object.keys(data).length > 0)) {
            console.log(`[Display] 有技术指标数据，创建显示内容`);
            indicatorsGrid.innerHTML = this.createIndicatorsGrid(data);
        } else {
            console.log(`[Display] 无技术指标数据，显示提示信息`);
            indicatorsGrid.innerHTML = '<div class="no-data">暂无技术指标数据</div>';
        }
        
        // 滚动到结果
        section.scrollIntoView({ behavior: 'smooth' });
        
        console.log(`[Display] 技术指标显示完成`);
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
     * 显示买卖预测
     */
    displayPredictions(data, stockCode, stockBasic) {
        console.log(`[Display] 显示买卖预测 - 数据:`, data);
        
        const section = document.getElementById('daily-chart-section');
        const predictionsContainer = document.getElementById('predictionsContainer');
        
        if (!section || !predictionsContainer) {
            console.error(`[Display] 找不到必要的DOM元素:`, { section: !!section, predictionsContainer: !!predictionsContainer });
            return;
        }
        
        // 显示section
        section.style.display = 'block';
        section.classList.add('fade-in');
        
        // 更新股票分析标题
        this.updateStockAnalysisTitle(stockCode, stockBasic);
        
        // 切换到买卖预测tab（不触发数据加载）
        console.log(`[Display] 准备切换到predictions tab`);
        this.switchToTabWithoutDataLoad('predictions');
        
        // 显示预测数据 - 修复数据检查逻辑
        if (data && (Array.isArray(data) ? data.length > 0 : Object.keys(data).length > 0)) {
            console.log(`[Display] 有预测数据，创建显示内容`);
            predictionsContainer.innerHTML = this.createPredictionsDisplay(data);
        } else {
            console.log(`[Display] 无预测数据，显示提示信息`);
            predictionsContainer.innerHTML = '<div class="no-data">暂无预测数据</div>';
        }
        
        // 滚动到结果
        section.scrollIntoView({ behavior: 'smooth' });
        
        console.log(`[Display] 买卖预测显示完成`);
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
     * 切换到指定tab
     */
    switchToTab(tabName) {
        // 调用events模块中的switchTab方法
        if (window.eventsModule && window.eventsModule.switchTab) {
            window.eventsModule.switchTab(tabName);
        } else {
            // 备用方案：直接操作DOM
            // 移除所有tab按钮的active状态
            const tabButtons = document.querySelectorAll('.tab-btn');
            tabButtons.forEach(btn => {
                btn.classList.remove('active');
            });

            // 隐藏所有tab内容
            const tabPanes = document.querySelectorAll('.tab-pane');
            tabPanes.forEach(pane => {
                pane.classList.remove('active');
            });

            // 激活选中的tab按钮
            const activeTabBtn = document.querySelector(`[data-tab="${tabName}"]`);
            if (activeTabBtn) {
                activeTabBtn.classList.add('active');
            }

            // 显示选中的tab内容
            const activeTabPane = document.getElementById(`${tabName}-tab`);
            if (activeTabPane) {
                activeTabPane.classList.add('active');
            }
        }
    }

    /**
     * 切换到指定tab（不触发数据加载）
     */
    switchToTabWithoutDataLoad(tabName) {
        // 调用events模块中的switchTab方法，但不加载数据
        if (window.eventsModule && window.eventsModule.switchTab) {
            window.eventsModule.switchTab(tabName, false);
        } else {
            // 备用方案：直接操作DOM
            // 移除所有tab按钮的active状态
            const tabButtons = document.querySelectorAll('.tab-btn');
            tabButtons.forEach(btn => {
                btn.classList.remove('active');
            });

            // 隐藏所有tab内容
            const tabPanes = document.querySelectorAll('.tab-pane');
            tabPanes.forEach(pane => {
                pane.classList.remove('active');
            });

            // 激活选中的tab按钮
            const activeTabBtn = document.querySelector(`[data-tab="${tabName}"]`);
            if (activeTabBtn) {
                activeTabBtn.classList.add('active');
            }

            // 显示选中的tab内容
            const activeTabPane = document.getElementById(`${tabName}-tab`);
            if (activeTabPane) {
                activeTabPane.classList.add('active');
            }
        }
    }

    /**
     * 创建技术指标网格
     */
    createIndicatorsGrid(data) {
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
        
        return indicatorsHTML;
    }

    /**
     * 更新股票分析标题
     */
    updateStockAnalysisTitle(stockCode, stockBasic) {
        const titleElement = document.getElementById('stockAnalysisTitle');
        if (!titleElement) {
            console.warn('[Display] 找不到股票分析标题元素');
            return;
        }
        
        let stockName = stockCode;
        if (stockBasic && stockBasic.name) {
            stockName = `${stockBasic.name} (${stockCode})`;
        }
        
        titleElement.textContent = `📊 ${stockName} - 股票数据分析`;
        console.log(`[Display] 更新股票分析标题: ${stockName}`);
    }

    /**
     * 格式化日期
     */
    formatDate(dateString) {
        if (!dateString) return 'N/A';
        
        try {
            // 处理 YYYYMMDD 格式
            if (dateString.length === 8) {
                const year = dateString.substring(0, 4);
                const month = dateString.substring(4, 6);
                const day = dateString.substring(6, 8);
                return `${year}-${month}-${day}`;
            }
            
            // 处理其他格式，尝试直接解析
            const date = new Date(dateString);
            if (isNaN(date.getTime())) return dateString;
            
            return date.toLocaleDateString('zh-CN');
        } catch (error) {
            console.warn('[Display] 日期格式化失败:', dateString, error);
            return dateString;
        }
    }

    /**
     * 创建预测数据显示
     */
    createPredictionsDisplay(data) {
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
            
            data.predictions.forEach((prediction, index) => {
                const typeClass = prediction.type.toLowerCase();
                const icon = prediction.type.toLowerCase() === 'buy' ? '📈' : '📉';
                const typeText = prediction.type.toLowerCase() === 'buy' ? '买入' : '卖出';
                
                // 提取强度等级
                const strength = this.extractStrengthFromReason(prediction.reason);
                const isWeak = strength === 'WEAK';
                const isCollapsed = isWeak ? 'collapsed' : '';
                // 使用统一的下拉图标，通过CSS旋转控制方向
                const collapseIcon = '🔽';
                
                predictionsHTML += `
                    <div class="prediction-item ${typeClass} slide-in ${isCollapsed}" data-index="${index}">
                        <div class="prediction-header" onclick="this.parentElement.classList.toggle('collapsed')">
                            <div class="prediction-icon">${icon}</div>
                            <div class="prediction-content">
                                <div class="prediction-type">
                                    ${typeText}信号
                                    <span class="info-icon" title="买卖信号类型：BUY=买入，SELL=卖出">ℹ️</span>
                                </div>
                                <div class="prediction-price">
                                    ¥${prediction.price?.toFixed(2) || 'N/A'}
                                    <span class="info-icon" title="预测的目标价格">ℹ️</span>
                                </div>
                                <div class="prediction-signal-date">
                                    📅 ${this.formatDate(prediction.signal_date) || 'N/A'}
                                    <span class="info-icon" title="信号产生的日期">ℹ️</span>
                                </div>
                            </div>
                            <div class="collapse-toggle">
                                <span class="collapse-icon">${collapseIcon}</span>
                                <span class="strength-badge ${strength.toLowerCase()}">${strength}</span>
                            </div>
                        </div>
                        <div class="prediction-details">
                            <div class="prediction-probability">
                                概率: ${(prediction.probability * 100).toFixed(1)}%
                                <span class="info-icon" title="预测成功的概率，基于技术指标置信度和历史表现">ℹ️</span>
                            </div>
                            <div class="prediction-reason">
                                ${prediction.reason || '基于技术指标分析'}
                                <span class="info-icon" title="预测依据：包含识别的技术模式、置信度和强度等级">ℹ️</span>
                            </div>
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
        
        return predictionsHTML;
    }

    /**
     * 从预测理由中提取强度等级
     */
    extractStrengthFromReason(reason) {
        if (!reason) return 'WEAK';
        
        if (reason.includes('强度：STRONG')) {
            return 'STRONG';
        } else if (reason.includes('强度：MEDIUM')) {
            return 'MEDIUM';
        } else if (reason.includes('强度：WEAK')) {
            return 'WEAK';
        }
        
        return 'WEAK'; // 默认值
    }
}

// 导出数据展示模块类
window.DisplayModule = DisplayModule;
