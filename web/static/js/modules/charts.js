/**
 * 图表渲染模块
 * 负责K线图和技术指标的图表展示
 */

class ChartsModule {
    constructor(client) {
        this.client = client;
        this.currentChart = null;
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
                        show: true,
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
}

// 导出图表模块类
window.ChartsModule = ChartsModule;
