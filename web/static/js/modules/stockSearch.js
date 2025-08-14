/**
 * 股票搜索模块
 * 负责股票搜索和建议功能
 */

class StockSearchModule {
    constructor(client, apiService) {
        this.client = client;
        this.apiService = apiService;
        this.searchTimeout = null;
        this.selectedSuggestionIndex = -1;
        
        this.init();
    }

    /**
     * 初始化搜索模块
     */
    init() {
        this.setupEventListeners();
    }

    /**
     * 设置事件监听器
     */
    setupEventListeners() {
        const searchInput = document.getElementById('stockSearch');
        const suggestionsContainer = document.getElementById('searchSuggestions');

        if (!searchInput || !suggestionsContainer) return;

        // 输入事件 - 实时搜索
        searchInput.addEventListener('input', (e) => {
            const keyword = e.target.value.trim();
            
            // 清除之前的搜索定时器
            if (this.searchTimeout) {
                clearTimeout(this.searchTimeout);
            }

            // 如果输入为空，隐藏建议
            if (!keyword) {
                this.hideSuggestions();
                return;
            }

            // 防抖搜索，300ms后执行
            this.searchTimeout = setTimeout(() => {
                this.performStockSearch(keyword);
            }, 300);
        });

        // 键盘导航
        searchInput.addEventListener('keydown', (e) => {
            const suggestions = suggestionsContainer.querySelectorAll('.search-suggestion-item');
            
            switch (e.key) {
                case 'ArrowDown':
                    e.preventDefault();
                    this.selectedSuggestionIndex = Math.min(this.selectedSuggestionIndex + 1, suggestions.length - 1);
                    this.highlightSuggestion();
                    break;
                    
                case 'ArrowUp':
                    e.preventDefault();
                    this.selectedSuggestionIndex = Math.max(this.selectedSuggestionIndex - 1, -1);
                    this.highlightSuggestion();
                    break;
                    
                case 'Enter':
                    e.preventDefault();
                    if (this.selectedSuggestionIndex >= 0 && suggestions.length > 0) {
                        this.selectSuggestion(suggestions[this.selectedSuggestionIndex]);
                    }
                    break;
                    
                case 'Escape':
                    this.hideSuggestions();
                    break;
            }
        });

        // 失去焦点时延迟隐藏建议（允许点击建议）
        searchInput.addEventListener('blur', () => {
            setTimeout(() => {
                this.hideSuggestions();
            }, 200);
        });

        // 获得焦点时显示建议（如果有内容）
        searchInput.addEventListener('focus', () => {
            if (searchInput.value.trim()) {
                this.showSuggestions();
            }
        });

        // 点击页面其他地方时隐藏建议
        document.addEventListener('click', (e) => {
            if (!searchInput.contains(e.target) && !suggestionsContainer.contains(e.target)) {
                this.hideSuggestions();
            }
        });
    }

    /**
     * 执行股票搜索
     */
    async performStockSearch(keyword) {
        const suggestionsContainer = document.getElementById('searchSuggestions');
        
        try {
            // 显示加载状态
            suggestionsContainer.innerHTML = '<div class="search-loading">搜索中...</div>';
            this.showSuggestions();

            // 调用搜索API
            const stocks = await this.apiService.searchStocks(keyword, 10);
            this.displaySearchSuggestions(stocks);
        } catch (error) {
            console.error('股票搜索失败:', error);
            this.displaySearchError();
        }
    }

    /**
     * 显示搜索建议
     */
    displaySearchSuggestions(stocks) {
        const suggestionsContainer = document.getElementById('searchSuggestions');
        
        if (!stocks || stocks.length === 0) {
            this.displayNoResults();
            return;
        }

        let suggestionsHTML = '';
        stocks.forEach((stock, index) => {
            const marketText = stock.market === 'SH' ? '上交所' : stock.market === 'SZ' ? '深交所' : stock.market;
            
            suggestionsHTML += `
                <div class="search-suggestion-item" data-stock-code="${stock.ts_code}" data-index="${index}">
                    <div class="stock-info">
                        <div class="stock-name">${stock.name}</div>
                        <div class="stock-code">${stock.ts_code}</div>
                    </div>
                    <div class="stock-market">${marketText}</div>
                </div>
            `;
        });

        suggestionsContainer.innerHTML = suggestionsHTML;
        this.showSuggestions();

        // 添加点击事件
        suggestionsContainer.querySelectorAll('.search-suggestion-item').forEach(item => {
            item.addEventListener('click', () => {
                this.selectSuggestion(item);
            });
        });

        this.selectedSuggestionIndex = -1;
    }

    /**
     * 显示无结果
     */
    displayNoResults() {
        const suggestionsContainer = document.getElementById('searchSuggestions');
        suggestionsContainer.innerHTML = '<div class="no-search-results">未找到相关股票</div>';
        this.showSuggestions();
    }

    /**
     * 显示搜索错误
     */
    displaySearchError() {
        const suggestionsContainer = document.getElementById('searchSuggestions');
        suggestionsContainer.innerHTML = '<div class="no-search-results">搜索出错，请稍后重试</div>';
        this.showSuggestions();
    }

    /**
     * 选择建议项
     */
    selectSuggestion(suggestionElement) {
        const stockCode = suggestionElement.dataset.stockCode;
        const stockName = suggestionElement.querySelector('.stock-name').textContent;
        
        // 填充到代码输入框
        const stockCodeInput = document.getElementById('stockCode');
        if (stockCodeInput) {
            stockCodeInput.value = stockCode;
        }
        
        // 更新搜索框显示
        const stockSearchInput = document.getElementById('stockSearch');
        if (stockSearchInput) {
            stockSearchInput.value = `${stockName} (${stockCode})`;
        }
        
        // 保存到本地缓存
        this.saveStockToCache(stockCode, stockName);
        
        // 隐藏建议
        this.hideSuggestions();
        
        // 触发验证
        if (stockCodeInput) {
            const event = new Event('input', { bubbles: true });
            stockCodeInput.dispatchEvent(event);
        }
        
        console.log(`已选择股票: ${stockName} (${stockCode})`);
    }

    /**
     * 高亮选中的建议
     */
    highlightSuggestion() {
        const suggestions = document.querySelectorAll('.search-suggestion-item');
        
        // 清除所有高亮
        suggestions.forEach(item => {
            item.classList.remove('highlighted');
        });

        // 高亮当前选中项
        if (this.selectedSuggestionIndex >= 0 && this.selectedSuggestionIndex < suggestions.length) {
            suggestions[this.selectedSuggestionIndex].classList.add('highlighted');
        }
    }

    /**
     * 显示建议框
     */
    showSuggestions() {
        const suggestionsContainer = document.getElementById('searchSuggestions');
        if (suggestionsContainer) {
            suggestionsContainer.style.display = 'block';
        }
    }

    /**
     * 隐藏建议框
     */
    hideSuggestions() {
        const suggestionsContainer = document.getElementById('searchSuggestions');
        if (suggestionsContainer) {
            suggestionsContainer.style.display = 'none';
        }
        this.selectedSuggestionIndex = -1;
    }

    /**
     * 保存股票信息到本地缓存
     */
    saveStockToCache(stockCode, stockName) {
        localStorage.setItem('stockapi_last_stock_code', stockCode);
        localStorage.setItem('stockapi_last_stock_name', stockName);
        console.log(`股票信息已缓存: ${stockName} (${stockCode})`);
    }

    /**
     * 设置默认股票代码（从缓存读取）
     */
    setDefaultStockCode() {
        const cachedStockCode = localStorage.getItem('stockapi_last_stock_code');
        const cachedStockName = localStorage.getItem('stockapi_last_stock_name');
        
        const stockCodeInput = document.getElementById('stockCode');
        const stockSearchInput = document.getElementById('stockSearch');
        
        if (cachedStockCode && stockCodeInput) {
            // 设置股票代码
            stockCodeInput.value = cachedStockCode;
            
            // 如果有股票名称，也设置到搜索框
            if (cachedStockName && stockSearchInput) {
                stockSearchInput.value = `${cachedStockName} (${cachedStockCode})`;
            }
            console.log(`从缓存恢复股票代码: ${cachedStockCode}`);
        }
    }
}

// 导出股票搜索模块类
window.StockSearchModule = StockSearchModule;
