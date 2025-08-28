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
}

window.FundamentalModule = FundamentalModule;
