/**
 * 收藏功能模块
 * 负责收藏股票的UI交互和管理
 */

class FavoritesModule {
    constructor(client, favoritesService, apiService) {
        this.client = client;
        this.favoritesService = favoritesService;
        this.apiService = apiService;
        this.favorites = [];
        this.currentStockCode = null;
        this.currentStockName = null;
        
        this.init();
    }

    /**
     * 初始化收藏模块
     */
    init() {
        this.createFavoritesUI();
        this.setupEventListeners();
        this.loadFavorites();
    }

    /**
     * 创建收藏功能的UI
     */
    createFavoritesUI() {
        // 在搜索区域添加收藏按钮
        const searchSection = document.querySelector('.search-section .card');
        if (searchSection) {
            const buttonGroup = searchSection.querySelector('.button-group');
            if (buttonGroup) {
                const favoriteBtn = document.createElement('button');
                favoriteBtn.id = 'favoriteBtn';
                favoriteBtn.className = 'btn btn-favorite';
                favoriteBtn.innerHTML = '⭐ 收藏';
                favoriteBtn.title = '收藏当前股票和时间范围';
                buttonGroup.appendChild(favoriteBtn);
            }
        }

        // 创建收藏列表区域
        const mainContent = document.querySelector('.main-content');
        if (mainContent) {
            const favoritesSection = document.createElement('section');
            favoritesSection.className = 'favorites-section';
            favoritesSection.innerHTML = `
                <div class="card">
                    <h2>⭐ 收藏股票</h2>
                    <div class="favorites-container">
                        <div class="favorites-header">
                            <span class="favorites-count">共 0 支股票</span>
                            <button id="refreshFavoritesBtn" class="btn btn-outline btn-small">🔄 刷新</button>
                        </div>
                        <div class="favorites-list" id="favoritesList">
                            <div class="favorites-empty">
                                <p>还没有收藏任何股票</p>
                                <p>点击上方的"收藏"按钮来添加股票</p>
                            </div>
                        </div>
                    </div>
                </div>
            `;
            
            // 插入到搜索区域之后
            const searchSection = document.querySelector('.search-section');
            if (searchSection) {
                searchSection.insertAdjacentElement('afterend', favoritesSection);
            } else {
                mainContent.insertBefore(favoritesSection, mainContent.firstChild);
            }
        }
    }

    /**
     * 设置事件监听器
     */
    setupEventListeners() {
        // 收藏按钮点击事件
        const favoriteBtn = document.getElementById('favoriteBtn');
        if (favoriteBtn) {
            favoriteBtn.addEventListener('click', () => {
                this.handleFavoriteClick();
            });
        }

        // 刷新收藏列表按钮
        const refreshBtn = document.getElementById('refreshFavoritesBtn');
        if (refreshBtn) {
            refreshBtn.addEventListener('click', () => {
                this.loadFavorites();
            });
        }

        // 监听股票代码变化，更新收藏按钮状态
        const stockCodeInput = document.getElementById('stockCode');
        if (stockCodeInput) {
            stockCodeInput.addEventListener('input', () => {
                this.updateFavoriteButtonState();
            });
        }
    }

    /**
     * 处理收藏按钮点击
     */
    async handleFavoriteClick() {
        const stockCode = document.getElementById('stockCode')?.value?.trim();
        const startDate = document.getElementById('startDate')?.value;
        const endDate = document.getElementById('endDate')?.value;
        
        if (!stockCode) {
            this.showMessage('请先输入股票代码', 'error');
            return;
        }

        // 检查是否已收藏
        const isFavorite = await this.favoritesService.checkFavorite(stockCode);
        if (isFavorite) {
            this.showMessage('该股票已在收藏列表中', 'warning');
            return;
        }

        try {
            // 获取股票名称
            let stockName = stockCode;
            try {
                const stockBasic = await this.apiService.getStockBasic(stockCode);
                if (stockBasic && stockBasic.name) {
                    stockName = stockBasic.name;
                }
            } catch (error) {
                console.warn('获取股票名称失败，使用代码作为名称:', error);
            }

            // 添加收藏
            const favorite = await this.favoritesService.addFavorite(
                stockCode,
                stockName,
                startDate || '',
                endDate || ''
            );

            this.showMessage(`成功收藏 ${stockName}`, 'success');
            this.loadFavorites(); // 刷新收藏列表
            this.updateFavoriteButtonState();
            
        } catch (error) {
            console.error('添加收藏失败:', error);
            this.showMessage(error.message || '添加收藏失败', 'error');
        }
    }

    /**
     * 加载收藏列表
     */
    async loadFavorites() {
        try {
            this.favorites = await this.favoritesService.getFavorites();
            this.renderFavoritesList();
        } catch (error) {
            console.error('加载收藏列表失败:', error);
            this.showMessage('加载收藏列表失败', 'error');
        }
    }

    /**
     * 渲染收藏列表
     */
    renderFavoritesList() {
        const favoritesList = document.getElementById('favoritesList');
        const favoritesCount = document.querySelector('.favorites-count');
        
        if (!favoritesList) return;

        // 更新计数
        if (favoritesCount) {
            favoritesCount.textContent = `共 ${this.favorites.length} 支股票`;
        }

        // 如果没有收藏，显示空状态
        if (this.favorites.length === 0) {
            favoritesList.innerHTML = `
                <div class="favorites-empty">
                    <p>还没有收藏任何股票</p>
                    <p>点击上方的"收藏"按钮来添加股票</p>
                </div>
            `;
            return;
        }

        // 渲染收藏列表
        let listHTML = '';
        this.favorites.forEach(favorite => {
            const createdDate = new Date(favorite.created_at).toLocaleDateString();
            const dateRange = this.formatDateRange(favorite.start_date, favorite.end_date);
            
            listHTML += `
                <div class="favorite-item" data-favorite-id="${favorite.id}" data-stock-code="${favorite.ts_code}">
                    <div class="favorite-info">
                        <div class="favorite-stock">
                            <span class="stock-name">${favorite.name}</span>
                            <span class="stock-code">${favorite.ts_code}</span>
                        </div>
                        <div class="favorite-details">
                            <span class="date-range">${dateRange}</span>
                            <span class="created-date">收藏于 ${createdDate}</span>
                        </div>
                    </div>
                    <div class="favorite-actions">
                        <button class="btn btn-outline btn-small load-chart-btn" title="查看K线图">
                            📈 查看
                        </button>
                        <button class="btn btn-outline btn-small delete-favorite-btn" title="删除收藏">
                            🗑️ 删除
                        </button>
                    </div>
                </div>
            `;
        });

        favoritesList.innerHTML = listHTML;

        // 添加事件监听器
        this.setupFavoriteItemEvents();
    }

    /**
     * 设置收藏项的事件监听器
     */
    setupFavoriteItemEvents() {
        const favoritesList = document.getElementById('favoritesList');
        if (!favoritesList) return;

        // 查看图表按钮
        favoritesList.querySelectorAll('.load-chart-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                const favoriteItem = e.target.closest('.favorite-item');
                this.loadFavoriteChart(favoriteItem);
            });
        });

        // 删除收藏按钮
        favoritesList.querySelectorAll('.delete-favorite-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                const favoriteItem = e.target.closest('.favorite-item');
                this.deleteFavorite(favoriteItem);
            });
        });

        // 点击收藏项也加载图表
        favoritesList.querySelectorAll('.favorite-item').forEach(item => {
            item.addEventListener('click', () => {
                this.loadFavoriteChart(item);
            });
        });
    }

    /**
     * 加载收藏股票的图表
     */
    async loadFavoriteChart(favoriteItem) {
        const favoriteId = favoriteItem.dataset.favoriteId;
        const stockCode = favoriteItem.dataset.stockCode;
        
        // 找到对应的收藏记录
        const favorite = this.favorites.find(f => f.id === favoriteId);
        if (!favorite) {
            this.showMessage('收藏记录不存在', 'error');
            return;
        }

        try {
            // 高亮当前选中的收藏项
            document.querySelectorAll('.favorite-item').forEach(item => {
                item.classList.remove('active');
            });
            favoriteItem.classList.add('active');

            // 设置股票代码和名称到搜索框
            const stockCodeInput = document.getElementById('stockCode');
            const stockSearchInput = document.getElementById('stockSearch');
            
            if (stockCodeInput) {
                stockCodeInput.value = favorite.ts_code;
            }
            if (stockSearchInput) {
                stockSearchInput.value = `${favorite.name} (${favorite.ts_code})`;
            }

            // 恢复收藏时的时间范围
            const startDateInput = document.getElementById('startDate');
            const endDateInput = document.getElementById('endDate');
            
            if (startDateInput && favorite.start_date) {
                const formattedStartDate = this.formatDateForInput(favorite.start_date);
                if (formattedStartDate) {
                    startDateInput.value = formattedStartDate;
                }
            }
            if (endDateInput && favorite.end_date) {
                const formattedEndDate = this.formatDateForInput(favorite.end_date);
                if (formattedEndDate) {
                    endDateInput.value = formattedEndDate;
                }
            }

            // 触发查询日线数据
            const queryDailyBtn = document.getElementById('queryDaily');
            if (queryDailyBtn) {
                queryDailyBtn.click();
            }

            this.showMessage(`已加载 ${favorite.name} 的K线图`, 'success');
            
        } catch (error) {
            console.error('加载收藏图表失败:', error);
            this.showMessage('加载图表失败', 'error');
        }
    }

    /**
     * 删除收藏
     */
    async deleteFavorite(favoriteItem) {
        const favoriteId = favoriteItem.dataset.favoriteId;
        const favorite = this.favorites.find(f => f.id === favoriteId);
        
        if (!favorite) return;

        if (!confirm(`确定要删除收藏的股票 "${favorite.name}" 吗？`)) {
            return;
        }

        try {
            await this.favoritesService.deleteFavorite(favoriteId);
            this.showMessage(`已删除收藏 ${favorite.name}`, 'success');
            this.loadFavorites(); // 刷新列表
            this.updateFavoriteButtonState();
        } catch (error) {
            console.error('删除收藏失败:', error);
            this.showMessage(error.message || '删除收藏失败', 'error');
        }
    }

    /**
     * 更新收藏按钮状态
     */
    async updateFavoriteButtonState() {
        const favoriteBtn = document.getElementById('favoriteBtn');
        const stockCode = document.getElementById('stockCode')?.value?.trim();
        
        if (!favoriteBtn || !stockCode) return;

        try {
            const isFavorite = await this.favoritesService.checkFavorite(stockCode);
            
            if (isFavorite) {
                favoriteBtn.innerHTML = '⭐ 已收藏';
                favoriteBtn.className = 'btn btn-favorite favorited';
                favoriteBtn.disabled = true;
                favoriteBtn.title = '该股票已在收藏列表中';
            } else {
                favoriteBtn.innerHTML = '⭐ 收藏';
                favoriteBtn.className = 'btn btn-favorite';
                favoriteBtn.disabled = false;
                favoriteBtn.title = '收藏当前股票和时间范围';
            }
        } catch (error) {
            console.error('更新收藏按钮状态失败:', error);
        }
    }

    /**
     * 格式化日期范围显示
     */
    formatDateRange(startDate, endDate) {
        if (!startDate && !endDate) {
            return '全部时间';
        }
        
        const formatDate = (dateStr) => {
            if (!dateStr) return '';
            if (dateStr.length === 8) {
                // YYYYMMDD格式
                return `${dateStr.substring(0,4)}-${dateStr.substring(4,6)}-${dateStr.substring(6,8)}`;
            }
            return dateStr;
        };

        const start = formatDate(startDate);
        const end = formatDate(endDate);
        
        if (start && end) {
            return `${start} 至 ${end}`;
        } else if (start) {
            return `从 ${start}`;
        } else if (end) {
            return `到 ${end}`;
        }
        
        return '全部时间';
    }

    /**
     * 格式化日期为input[type="date"]的格式
     */
    formatDateForInput(dateStr) {
        if (!dateStr) return '';
        
        if (dateStr.length === 8) {
            // YYYYMMDD格式转换为YYYY-MM-DD
            return `${dateStr.substring(0,4)}-${dateStr.substring(4,6)}-${dateStr.substring(6,8)}`;
        }
        
        return dateStr;
    }

    /**
     * 显示消息提示
     */
    showMessage(message, type = 'info') {
        // 创建消息元素
        const messageDiv = document.createElement('div');
        messageDiv.className = `message message-${type}`;
        messageDiv.textContent = message;
        
        // 添加到页面
        document.body.appendChild(messageDiv);
        
        // 自动消失
        setTimeout(() => {
            if (messageDiv.parentNode) {
                messageDiv.parentNode.removeChild(messageDiv);
            }
        }, 3000);
    }

    /**
     * 获取收藏列表
     */
    getFavorites() {
        return this.favorites;
    }

    /**
     * 获取收藏数量
     */
    getFavoritesCount() {
        return this.favorites.length;
    }
}

// 导出收藏模块类
window.FavoritesModule = FavoritesModule;
