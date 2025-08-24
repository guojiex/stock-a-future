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
    createIndicatorItem(title, values, signal = null, tooltip = null) {
        const signalHTML = signal ? `<span class="signal ${signal.toLowerCase()}">${this.getSignalText(signal)}</span>` : '';
        const tooltipHTML = tooltip ? `<span class="indicator-tooltip" data-tooltip="${tooltip}">ℹ️</span>` : '';
        
        const valuesHTML = values.map(item => 
            `<div class="indicator-value">
                <span class="name">${item.name}:</span>
                <span class="value">${item.value}</span>
            </div>`
        ).join('');
        
        return `
            <div class="indicator-item slide-in">
                <h4>${title} ${signalHTML} ${tooltipHTML}</h4>
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
            
            // 设置日期点击事件
            this.setupDateClickHandlers(stockCode);
        } else {
            console.log(`[Display] 无预测数据，显示提示信息`);
            predictionsContainer.innerHTML = '<div class="no-data">暂无预测数据</div>';
        }
        
        // 滚动到结果
        section.scrollIntoView({ behavior: 'smooth' });
        
        console.log(`[Display] 买卖预测显示完成`);
    }
    
    /**
     * 设置日期点击事件处理
     */
    setupDateClickHandlers(stockCode) {
        console.log(`[Display] 设置日期点击事件处理`);
        
        // 清除之前的事件处理函数
        if (this.dateClickHandlers && this.dateClickHandlers.length > 0) {
            this.dateClickHandlers.forEach(handler => {
                if (handler.element && handler.element.removeEventListener) {
                    handler.element.removeEventListener('click', handler.callback);
                }
            });
            this.dateClickHandlers = [];
        }
        
        // 查找所有日期链接
        const dateLinks = document.querySelectorAll('.date-link');
        console.log(`[Display] 找到 ${dateLinks.length} 个日期链接`);
        
        dateLinks.forEach(link => {
            const date = link.getAttribute('data-date');
            if (!date) return;
            
            const clickHandler = (e) => {
                e.preventDefault();
                this.handleDateLinkClick(date, stockCode);
            };
            
            link.addEventListener('click', clickHandler);
            
            // 存储事件处理函数，以便后续清理
            this.dateClickHandlers.push({
                element: link,
                callback: clickHandler
            });
            
            // 添加视觉提示，表明可点击
            link.classList.add('clickable');
            link.title = '点击跳转到日K线对应日期';
        });
    }
    
    /**
     * 处理日期链接点击
     */
    handleDateLinkClick(date, stockCode) {
        console.log(`[Display] 日期链接点击: ${date}, 股票代码: ${stockCode}`);
        
        if (!date || !stockCode) {
            console.warn(`[Display] 日期或股票代码无效: ${date}, ${stockCode}`);
            return;
        }
        
        // 标准化日期格式（移除时间部分）
        let normalizedDate = date;
        if (normalizedDate && normalizedDate.includes('T')) {
            normalizedDate = normalizedDate.split('T')[0];
        }
        
        // 切换到日线数据tab
        this.switchToTab('daily-data');
        
        // 等待图表加载完成
        setTimeout(() => {
            // 使用chartsModule导航到指定日期
            if (this.chartsModule && this.chartsModule.navigateToDate) {
                // 使用标准化后的日期
                const success = this.chartsModule.navigateToDate(normalizedDate);
                if (success) {
                    console.log(`[Display] 成功导航到日期: ${normalizedDate}`);
                } else {
                    console.warn(`[Display] 导航到日期失败: ${normalizedDate}`);
                    
                    // 检查当前图表中是否有数据
                    if (this.chartsModule.currentDates && this.chartsModule.currentDates.length > 0) {
                        // 如果有数据，但找不到特定日期，提示用户调整查询范围
                        const dateStr = this.formatDateForDisplay(normalizedDate);
                        const earliestDate = this.formatDateForDisplay(this.chartsModule.currentDates[0]);
                        const latestDate = this.formatDateForDisplay(this.chartsModule.currentDates[this.chartsModule.currentDates.length - 1]);
                        
                        this.client.showMessage(
                            `提示：未能在K线图中找到 ${dateStr} 的数据。当前K线图显示的日期范围是 ${earliestDate} 到 ${latestDate}，请调整查询范围。`, 
                            'warning',
                            5000 // 显示时间延长到5秒
                        );
                    } else {
                        // 如果没有数据，提示用户先获取数据
                        this.client.showMessage(`提示：请先获取日线数据。`, 'warning');
                    }
                }
            } else {
                console.error(`[Display] 图表模块不可用或缺少navigateToDate方法`);
            }
        }, 300); // 给图表加载留出时间
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
            ], data.macd.signal, 'MACD是趋势跟踪指标，通过快慢均线的差值判断买卖时机。DIF上穿DEA为金叉买入信号，下穿为死叉卖出信号。');
        }
        
        // RSI指标
        if (data.rsi) {
            indicatorsHTML += this.createIndicatorItem('RSI', [
                { name: 'RSI14', value: data.rsi.rsi14?.toFixed(2) || 'N/A' }
            ], data.rsi.signal, 'RSI相对强弱指数衡量价格变动的速度和幅度。RSI>70为超买区域，RSI<30为超卖区域，可作为反转信号参考。');
        }
        
        // 布林带
        if (data.boll) {
            indicatorsHTML += this.createIndicatorItem('布林带', [
                { name: '上轨', value: data.boll.upper?.toFixed(2) || 'N/A' },
                { name: '中轨', value: data.boll.middle?.toFixed(2) || 'N/A' },
                { name: '下轨', value: data.boll.lower?.toFixed(2) || 'N/A' }
            ], data.boll.signal, '布林带由移动平均线和标准差构成，价格触及上轨可能回调，触及下轨可能反弹。带宽收窄预示突破，扩张表示趋势延续。');
        }
        
        // 移动平均线
        if (data.ma) {
            indicatorsHTML += this.createIndicatorItem('移动平均线', [
                { name: 'MA5', value: data.ma.ma5?.toFixed(2) || 'N/A' },
                { name: 'MA10', value: data.ma.ma10?.toFixed(2) || 'N/A' },
                { name: 'MA20', value: data.ma.ma20?.toFixed(2) || 'N/A' },
                { name: 'MA60', value: data.ma.ma60?.toFixed(2) || 'N/A' }
            ], null, '移动平均线平滑价格波动，识别趋势方向。短期均线上穿长期均线为金叉买入信号，下穿为死叉卖出信号。');
        }
        
        // KDJ指标
        if (data.kdj) {
            indicatorsHTML += this.createIndicatorItem('KDJ', [
                { name: 'K值', value: data.kdj.k?.toFixed(2) || 'N/A' },
                { name: 'D值', value: data.kdj.d?.toFixed(2) || 'N/A' },
                { name: 'J值', value: data.kdj.j?.toFixed(2) || 'N/A' }
            ], data.kdj.signal, 'KDJ随机指标反映价格在一定周期内的相对位置。K>80且D>80为超买，K<20且D<20为超卖。J值更敏感，可提前预警。');
        }
        
        // === 动量因子指标 ===
        
        // 威廉指标 (%R)
        if (data.wr) {
            indicatorsHTML += this.createIndicatorItem('威廉指标 (%R)', [
                { name: 'WR14', value: data.wr.wr14?.toFixed(2) || 'N/A' }
            ], data.wr.signal, '威廉指标衡量股价在一定周期内的相对位置。WR>-20为超买区域，WR<-80为超卖区域。数值越接近0越超买，越接近-100越超卖。');
        }
        
        // 动量指标
        if (data.momentum) {
            indicatorsHTML += this.createIndicatorItem('动量指标', [
                { name: 'Momentum10', value: data.momentum.momentum10?.toFixed(4) || 'N/A' },
                { name: 'Momentum20', value: data.momentum.momentum20?.toFixed(4) || 'N/A' }
            ], data.momentum.signal, '动量指标衡量价格变化的速度。正值表示上涨动量，负值表示下跌动量。数值越大表示趋势越强劲。');
        }
        
        // 变化率指标 (ROC)
        if (data.roc) {
            indicatorsHTML += this.createIndicatorItem('变化率指标 (ROC)', [
                { name: 'ROC10', value: data.roc.roc10?.toFixed(2) || 'N/A' },
                { name: 'ROC20', value: data.roc.roc20?.toFixed(2) || 'N/A' }
            ], data.roc.signal, 'ROC变化率指标衡量价格在一定周期内的百分比变化。正值表示上涨，负值表示下跌。可用于判断趋势强度和转折点。');
        }
        
        // === 趋势因子指标 ===
        
        // 平均方向指数 (ADX)
        if (data.adx) {
            indicatorsHTML += this.createIndicatorItem('平均方向指数 (ADX)', [
                { name: 'ADX', value: data.adx.adx?.toFixed(2) || 'N/A' },
                { name: 'PDI', value: data.adx.pdi?.toFixed(2) || 'N/A' },
                { name: 'MDI', value: data.adx.mdi?.toFixed(2) || 'N/A' }
            ], data.adx.signal, 'ADX衡量趋势强度，不判断方向。ADX>25表示强趋势，<20表示弱趋势。PDI>MDI表示上升趋势，反之为下降趋势。');
        }
        
        // 抛物线转向 (SAR)
        if (data.sar) {
            indicatorsHTML += this.createIndicatorItem('抛物线转向 (SAR)', [
                { name: 'SAR', value: data.sar.sar?.toFixed(2) || 'N/A' }
            ], data.sar.signal, 'SAR抛物线转向指标用于确定止损点和趋势转换。价格在SAR之上为上升趋势，之下为下降趋势。SAR点位可作为止损参考。');
        }
        
        // 一目均衡表
        if (data.ichimoku) {
            indicatorsHTML += this.createIndicatorItem('一目均衡表', [
                { name: '转换线', value: data.ichimoku.tenkan_sen?.toFixed(2) || 'N/A' },
                { name: '基准线', value: data.ichimoku.kijun_sen?.toFixed(2) || 'N/A' },
                { name: '先行带A', value: data.ichimoku.senkou_span_a?.toFixed(2) || 'N/A' },
                { name: '先行带B', value: data.ichimoku.senkou_span_b?.toFixed(2) || 'N/A' }
            ], data.ichimoku.signal, '一目均衡表是综合性技术指标。转换线上穿基准线为买入信号，价格突破云带(先行带)确认趋势。云带厚度反映支撑阻力强度。');
        }
        
        // === 波动率因子指标 ===
        
        // 平均真实范围 (ATR)
        if (data.atr) {
            indicatorsHTML += this.createIndicatorItem('平均真实范围 (ATR)', [
                { name: 'ATR14', value: data.atr.atr14?.toFixed(4) || 'N/A' }
            ], data.atr.signal, 'ATR衡量价格波动幅度，数值越大表示波动越剧烈。可用于设置止损位和判断市场活跃度。高ATR适合趋势交易，低ATR适合区间交易。');
        }
        
        // 标准差
        if (data.stddev) {
            indicatorsHTML += this.createIndicatorItem('标准差', [
                { name: 'StdDev20', value: data.stddev.stddev20?.toFixed(4) || 'N/A' }
            ], data.stddev.signal, '标准差衡量价格偏离平均值的程度。数值越大表示价格波动越不稳定。可用于评估投资风险和市场不确定性。');
        }
        
        // 历史波动率
        if (data.hv) {
            indicatorsHTML += this.createIndicatorItem('历史波动率', [
                { name: 'HV20', value: data.hv.hv20?.toFixed(2) || 'N/A' },
                { name: 'HV60', value: data.hv.hv60?.toFixed(2) || 'N/A' }
            ], data.hv.signal, '历史波动率反映股价在过去一段时间的波动程度。高波动率意味着高风险高收益，低波动率表示价格相对稳定。');
        }
        
        // === 成交量因子指标 ===
        
        // 成交量加权平均价 (VWAP)
        if (data.vwap) {
            indicatorsHTML += this.createIndicatorItem('成交量加权平均价 (VWAP)', [
                { name: 'VWAP', value: data.vwap.vwap?.toFixed(2) || 'N/A' }
            ], data.vwap.signal, 'VWAP是成交量加权的平均价格，反映真实的平均成交成本。价格高于VWAP表示买盘强劲，低于VWAP表示卖盘压力大。');
        }
        
        // 累积/派发线 (A/D Line)
        if (data.ad_line) {
            indicatorsHTML += this.createIndicatorItem('累积/派发线 (A/D Line)', [
                { name: 'A/D Line', value: data.ad_line.ad_line?.toFixed(0) || 'N/A' }
            ], data.ad_line.signal, 'A/D线结合价格和成交量，衡量资金流向。上升表示资金流入(累积)，下降表示资金流出(派发)。可用于确认价格趋势。');
        }
        
        // 简易波动指标 (EMV)
        if (data.emv) {
            indicatorsHTML += this.createIndicatorItem('简易波动指标 (EMV)', [
                { name: 'EMV14', value: data.emv.emv14?.toFixed(4) || 'N/A' }
            ], data.emv.signal, 'EMV衡量价格变动的难易程度。正值表示价格上涨容易，负值表示下跌容易。结合成交量分析，判断价格变动的可持续性。');
        }
        
        // 量价确认指标 (VPT)
        if (data.vpt) {
            indicatorsHTML += this.createIndicatorItem('量价确认指标 (VPT)', [
                { name: 'VPT', value: data.vpt.vpt?.toFixed(2) || 'N/A' }
            ], data.vpt.signal, 'VPT将成交量与价格变化相结合，确认价格趋势。VPT与价格同向运动确认趋势，背离时可能预示反转。');
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
     * @param {string|Date} dateInput - 日期字符串或Date对象
     * @returns {string} - 格式化后的日期字符串
     */
    formatDate(dateInput) {
        if (!dateInput) return 'N/A';
        
        try {
            // 如果是Date对象，直接格式化
            if (dateInput instanceof Date) {
                return dateInput.toLocaleDateString('zh-CN');
            }
            
            // 处理字符串
            const dateString = String(dateInput);
            
            // 处理 YYYYMMDD 格式
            if (dateString.length === 8 && !dateString.includes('-')) {
                const year = dateString.substring(0, 4);
                const month = dateString.substring(4, 6);
                const day = dateString.substring(6, 8);
                return `${year}-${month}-${day}`;
            }
            
            // 处理带T的ISO格式日期，如 "2025-08-14T00:00:00.000"
            if (dateString.includes('T')) {
                return dateString.split('T')[0];
            }
            
            // 处理其他格式，尝试直接解析
            const date = new Date(dateString);
            if (isNaN(date.getTime())) return dateString;
            
            return date.toLocaleDateString('zh-CN');
        } catch (error) {
            console.warn('[Display] 日期格式化失败:', dateInput, error);
            return String(dateInput);
        }
    }
    
    /**
     * 格式化日期为更友好的显示格式（年/月/日）
     * @param {string|Date} dateInput - 日期字符串或Date对象
     * @returns {string} - 格式化后的日期字符串，如 2025/8/18
     */
    formatDateForDisplay(dateInput) {
        if (!dateInput) return 'N/A';
        
        try {
            let dateObj;
            
            // 如果是Date对象
            if (dateInput instanceof Date) {
                dateObj = dateInput;
            } else {
                // 处理字符串
                const dateString = String(dateInput);
                
                // 处理 YYYYMMDD 格式
                if (dateString.length === 8 && !dateString.includes('-')) {
                    const year = dateString.substring(0, 4);
                    const month = dateString.substring(4, 6);
                    const day = dateString.substring(6, 8);
                    dateObj = new Date(`${year}-${month}-${day}`);
                }
                // 处理带T的ISO格式日期，如 "2025-08-14T00:00:00.000"
                else if (dateString.includes('T')) {
                    dateObj = new Date(dateString);
                }
                // 处理 YYYY-MM-DD 格式
                else if (dateString.includes('-')) {
                    dateObj = new Date(dateString);
                }
                // 其他格式
                else {
                    dateObj = new Date(dateString);
                }
            }
            
            // 检查是否有效日期
            if (isNaN(dateObj.getTime())) {
                return String(dateInput);
            }
            
            // 格式化为 YYYY/M/D 格式
            const year = dateObj.getFullYear();
            const month = dateObj.getMonth() + 1; // 月份从0开始
            const day = dateObj.getDate();
            
            return `${year}/${month}/${day}`;
            
        } catch (error) {
            console.warn('[Display] 日期显示格式化失败:', dateInput, error);
            return String(dateInput);
        }
    }

    /**
     * 创建预测数据显示
     */
    createPredictionsDisplay(data) {
        let predictionsHTML = '';
        
        // 用于存储日期点击事件的处理函数
        this.dateClickHandlers = [];
        
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
                                <div class="prediction-main-row">
                                    <div class="prediction-type">
                                        ${typeText}信号
                                        <span class="info-icon" title="买卖信号类型：BUY=买入，SELL=卖出">ℹ️</span>
                                    </div>
                                    <div class="prediction-price">
                                        ¥${prediction.price?.toFixed(2) || 'N/A'}
                                        <span class="info-icon" title="预测的目标价格">ℹ️</span>
                                    </div>
                                    <div class="prediction-signal-date">
                                        <a href="javascript:void(0);" class="date-link" data-date="${prediction.signal_date || ''}">
                                            📅 ${this.formatDateForDisplay(prediction.signal_date) || 'N/A'}
                                        </a>
                                        <span class="info-icon" title="信号产生的日期 (点击可跳转到日K线对应日期)">ℹ️</span>
                                    </div>
                                    <div style="width: 10px;"></div> <!-- 添加一个空的div作为间隔 -->
                                    <div class="prediction-probability">
                                        概率: ${(prediction.probability * 100).toFixed(1)}%
                                        <span class="info-icon" title="预测成功的概率，基于技术指标置信度和历史表现">ℹ️</span>
                                    </div>
                                </div>
                            </div>
                            <div class="collapse-toggle">
                                <span class="collapse-icon">${collapseIcon}</span>
                                <span class="strength-badge ${strength.toLowerCase()}">${strength}</span>
                            </div>
                        </div>
                        <div class="prediction-details">
                            <div class="prediction-reason">
                                ${prediction.reason || '基于技术指标分析'}
                                <span class="info-icon" title="预测依据：包含识别的技术模式、置信度和强度等级">ℹ️</span>
                            </div>
                            ${prediction.backtested ? `
                            <div class="prediction-backtest">
                                <div class="backtest-summary">
                                    <div class="backtest-result ${prediction.is_correct ? 'correct' : 'incorrect'}">
                                        回测结果: ${prediction.is_correct ? '✅ 正确' : '❌ 错误'}
                                    </div>
                                    <div class="next-day-price">
                                        次日价格: ¥${prediction.next_day_price?.toFixed(2) || 'N/A'}
                                    </div>
                                    <div class="price-diff ${prediction.price_diff >= 0 ? 'positive' : 'negative'}">
                                        价差: ${prediction.price_diff >= 0 ? '+' : ''}${prediction.price_diff?.toFixed(2) || '0.00'} 
                                        (${prediction.price_diff_ratio >= 0 ? '+' : ''}${prediction.price_diff_ratio?.toFixed(2) || '0.00'}%)
                                    </div>
                                </div>
                            </div>
                            ` : ''}
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
