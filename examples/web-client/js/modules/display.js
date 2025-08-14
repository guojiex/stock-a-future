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
        const card = document.getElementById('dailyDataCard');
        const summary = document.getElementById('dailyDataSummary');
        
        if (!card || !summary) return;
        
        // 更新卡片标题
        const cardTitle = card.querySelector('h3');
        if (cardTitle && stockBasic && stockBasic.name) {
            cardTitle.textContent = `📈 ${stockBasic.name}(${stockCode}) - 日线数据`;
        }
        
        // 显示卡片
        card.style.display = 'block';
        card.classList.add('fade-in');
        
        // 创建价格图表
        this.chartsModule.createPriceChart(data, stockCode, stockBasic);
        
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
     * 显示技术指标
     */
    displayIndicators(data, stockCode, stockBasic) {
        const card = document.getElementById('indicatorsCard');
        const grid = document.getElementById('indicatorsGrid');
        
        if (!card || !grid) return;
        
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
        
        if (!card || !container) return;
        
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
}

// 导出数据展示模块类
window.DisplayModule = DisplayModule;
