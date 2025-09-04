/**
 * 图表渲染模块
 * 负责K线图和技术指标的图表展示
 */

class ChartsModule {
    constructor(client) {
        this.client = client;
        this.currentChart = null;
        this.currentDates = []; // 存储当前图表的日期数据
    }

    /**
     * 创建K线图表（使用ECharts）
     */
    createPriceChart(data, stockCode, stockBasic) {
        const chartContainer = document.getElementById('priceChart');
        
        if (!chartContainer) return;
        
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
            // 处理ISO 8601格式的日期字符串 (如: "2025-07-17T00:00:00.000")
            if (dateStr && dateStr.includes('T')) {
                // 提取日期部分，去掉时间部分
                return dateStr.split('T')[0];
            } else if (dateStr && dateStr.length === 8) {
                // 处理传统的8位数字格式 (如: "20250717")
                return `${dateStr.substring(0,4)}-${dateStr.substring(4,6)}-${dateStr.substring(6,8)}`;
            } else {
                // 如果格式不匹配，返回原始字符串
                return dateStr;
            }
        });
        
        // 存储日期数据，用于后续导航
        this.currentDates = dates;
        
        // 计算移动平均线
        const ma5 = this.calculateMA(data.map(item => parseFloat(item.close)), 5);
        const ma10 = this.calculateMA(data.map(item => parseFloat(item.close)), 10);
        const ma20 = this.calculateMA(data.map(item => parseFloat(item.close)), 20);
        
        // 获取最新的移动平均线数值
        const latestMA5 = ma5.length > 0 ? ma5[ma5.length - 1] : null;
        const latestMA10 = ma10.length > 0 ? ma10[ma10.length - 1] : null;
        const latestMA20 = ma20.length > 0 ? ma20[ma20.length - 1] : null;
        
        // 构建图表标题，包含最新的涨跌幅信息
        let chartTitle = `${stockCode} K线图`;
        let titleTextColor = '#333';
        
        if (stockBasic && stockBasic.name) {
            chartTitle = `${stockBasic.name}(${stockCode}) K线图`;
        }
        
        // 添加最新的涨跌幅信息到标题
        if (data && data.length > 0) {
            const latestData = data[data.length - 1];
            const currentClose = parseFloat(latestData.close);
            const preClose = parseFloat(latestData.pre_close) || (data.length > 1 ? parseFloat(data[data.length - 2].close) : currentClose);
            const change = parseFloat(latestData.change) || (currentClose - preClose);
            const changePercent = parseFloat(latestData.pct_chg) || ((change / preClose) * 100);
            
            const changeText = change >= 0 ? `+${change.toFixed(2)}` : change.toFixed(2);
            const percentText = change >= 0 ? `+${changePercent.toFixed(2)}%` : `${changePercent.toFixed(2)}%`;
            const changeIcon = change >= 0 ? '📈' : '📉';
            
            chartTitle += ` - ¥${currentClose.toFixed(2)} ${changeIcon} ${changeText} (${percentText})`;
            titleTextColor = change >= 0 ? '#10b981' : '#ef4444';
        }
        
        // 配置图表选项
        const option = {
            title: {
                text: chartTitle,
                left: 'center',
                textStyle: {
                    fontSize: 16,
                    fontWeight: 'bold',
                    color: titleTextColor
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
                        const klineData = params[0].data;
                        const dataIndex = params[0].dataIndex;
                        
                        result += `开盘价: ¥${klineData[1].toFixed(2)}<br/>`;
                        result += `收盘价: ¥${klineData[2].toFixed(2)}<br/>`;
                        result += `最低价: ¥${klineData[3].toFixed(2)}<br/>`;
                        result += `最高价: ¥${klineData[4].toFixed(2)}<br/>`;
                        
                        // 从原始数据中获取涨跌幅信息
                        if (data && data[dataIndex]) {
                            const currentData = data[dataIndex];
                            const currentClose = parseFloat(currentData.close);
                            const preClose = parseFloat(currentData.pre_close) || (dataIndex > 0 ? parseFloat(data[dataIndex - 1].close) : currentClose);
                            const change = parseFloat(currentData.change) || (currentClose - preClose);
                            const changePercent = parseFloat(currentData.pct_chg) || ((change / preClose) * 100);
                            
                            const changeText = change >= 0 ? `+${change.toFixed(2)}` : change.toFixed(2);
                            const percentText = change >= 0 ? `+${changePercent.toFixed(2)}%` : `${changePercent.toFixed(2)}%`;
                            const changeIcon = change >= 0 ? '📈' : '📉';
                            
                            result += `昨收价: ¥${preClose.toFixed(2)}<br/>`;
                            result += `${changeIcon} 涨跌额: ${changeText}<br/>`;
                            result += `${changeIcon} 涨跌幅: ${percentText}<br/>`;
                            
                            // 计算振幅
                            const amplitude = ((parseFloat(currentData.high) - parseFloat(currentData.low)) / preClose * 100).toFixed(2);
                            result += `振幅: ${amplitude}%<br/>`;
                        }
                    }
                    
                    // 成交量
                    const volumeParam = params.find(p => p.seriesName === '成交量');
                    if (volumeParam) {
                        // 修复成交量数据显示问题：正确处理数据格式
                        let volumeValue;
                        if (typeof volumeParam.data === 'object' && volumeParam.data.value !== undefined) {
                            volumeValue = volumeParam.data.value;
                        } else {
                            volumeValue = volumeParam.data;
                        }
                        result += `成交量: ${this.formatVolume(volumeValue)}<br/>`;
                    }
                    
                    return result;
                }.bind(this)
            },
            legend: {
                data: ['K线', 'MA5', 'MA10', 'MA20', '成交量'],
                top: 30,
                formatter: function(name) {
                    if (name === 'MA5' && latestMA5) {
                        return `MA5: ¥${latestMA5.toFixed(2)}`;
                    } else if (name === 'MA10' && latestMA10) {
                        return `MA10: ¥${latestMA10.toFixed(2)}`;
                    } else if (name === 'MA20' && latestMA20) {
                        return `MA20: ¥${latestMA20.toFixed(2)}`;
                    }
                    return name;
                }
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
                    top: '78%',
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
                    },
                    // 修复量柱遮挡价格问题：增加边界间距
                    boundaryGap: ['10%', '10%']
                },
                {
                    scale: true,
                    gridIndex: 1,
                    splitNumber: 4,
                    axisLabel: { 
                        show: true,
                        formatter: function(value) {
                            return this.formatVolume(value);
                        }.bind(this),
                        fontSize: 10,
                        color: '#666'
                    },
                    axisLine: { show: true, lineStyle: { color: '#ddd' } },
                    axisTick: { show: true, lineStyle: { color: '#ddd' } },
                    splitLine: { show: true, lineStyle: { color: '#eee' } },
                    // 修复量柱遮挡问题：增加边界间距
                    boundaryGap: ['10%', '10%']
                }
            ],
            dataZoom: [
                {
                    type: 'inside',
                    xAxisIndex: [0, 1],
                    start: 0,
                    end: 100
                },
                {
                    show: true,
                    xAxisIndex: [0, 1],
                    type: 'slider',
                    top: '95%',
                    start: 0,
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
                    }),
                    label: {
                        show: false,
                        position: 'top',
                        formatter: function(params) {
                            return this.formatVolume(params.value);
                        }.bind(this),
                        fontSize: 10,
                        color: '#666'
                    }
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
     * 销毁当前图表
     */
    destroyChart() {
        if (this.currentChart) {
            this.currentChart.dispose();
            this.currentChart = null;
        }
    }
    
    /**
     * 导航到特定日期
     * @param {string} targetDate - 目标日期，格式为 YYYY-MM-DD 或带时间的 ISO 格式
     * @returns {boolean} - 是否成功导航
     */
    navigateToDate(targetDate) {
        console.log(`[Charts] 尝试导航到日期: ${targetDate}`);
        
        if (!this.currentChart || !this.currentDates || this.currentDates.length === 0) {
            console.warn('[Charts] 无法导航：图表未初始化或无日期数据');
            return false;
        }
        
        // 标准化日期格式
        let normalizedTargetDate = targetDate;
        
        // 处理带T的ISO格式日期，如 "2025-08-14T00:00:00.000"
        if (normalizedTargetDate && normalizedTargetDate.includes('T')) {
            normalizedTargetDate = normalizedTargetDate.split('T')[0];
        }
        
        // 处理 YYYYMMDD 格式
        if (normalizedTargetDate && normalizedTargetDate.length === 8 && !normalizedTargetDate.includes('-')) {
            normalizedTargetDate = `${normalizedTargetDate.substring(0,4)}-${normalizedTargetDate.substring(4,6)}-${normalizedTargetDate.substring(6,8)}`;
        }
        
        console.log(`[Charts] 标准化后的日期: ${normalizedTargetDate}`);
        
        // 查找日期索引
        const dateIndex = this.currentDates.findIndex(date => date === normalizedTargetDate);
        if (dateIndex === -1) {
            console.warn(`[Charts] 未找到目标日期: ${normalizedTargetDate}`);
            return false;
        }
        
        console.log(`[Charts] 找到日期索引: ${dateIndex}, 总日期数: ${this.currentDates.length}`);
        
        // 计算数据窗口位置
        // 默认显示30天数据，将目标日期居中
        const totalDates = this.currentDates.length;
        const windowSize = Math.min(30, totalDates); // 显示窗口大小，最多30天
        
        let start, end;
        if (totalDates <= windowSize) {
            // 如果总数据量小于窗口大小，显示全部
            start = 0;
            end = 100;
        } else {
            // 计算百分比位置，使目标日期居中
            const halfWindow = Math.floor(windowSize / 2);
            let centerIndex = dateIndex;
            
            // 确保窗口不会超出数据范围
            if (centerIndex < halfWindow) {
                centerIndex = halfWindow;
            } else if (centerIndex > totalDates - halfWindow) {
                centerIndex = totalDates - halfWindow;
            }
            
            start = Math.max(0, (centerIndex - halfWindow) / totalDates * 100);
            end = Math.min(100, (centerIndex + halfWindow) / totalDates * 100);
        }
        
        // 设置数据区域缩放
        this.currentChart.dispatchAction({
            type: 'dataZoom',
            start: start,
            end: end
        });
        
        // 高亮目标日期点
        this.currentChart.dispatchAction({
            type: 'highlight',
            seriesIndex: 0,
            dataIndex: dateIndex
        });
        
        // 显示提示框
        this.currentChart.dispatchAction({
            type: 'showTip',
            seriesIndex: 0,
            dataIndex: dateIndex
        });
        
        console.log(`[Charts] 已导航到日期: ${normalizedTargetDate}, 缩放范围: ${start.toFixed(2)}% - ${end.toFixed(2)}%`);
        return true;
    }
}

// 导出图表模块类
window.ChartsModule = ChartsModule;
