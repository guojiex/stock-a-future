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
        this.groups = [];
        this.currentStockCode = null;
        this.currentStockName = null;
        this.currentGroupId = 'default'; // 跟踪当前选择的分组
        
        this.init();
    }

    /**
     * 初始化收藏模块
     */
    init() {
        this.createFavoritesUI();
        this.setupEventListeners();
        this.loadGroups();
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
            favoritesSection.id = 'favorites-section'; // 添加ID以便折叠功能识别
            favoritesSection.innerHTML = `
                <div class="card">
                    <h2>⭐ 收藏股票</h2>
                    
                    <!-- Tab导航 -->
                    <div class="favorites-tab-navigation">
                        <button class="favorites-tab-btn active" data-tab="favorites-list">
                            <span class="tab-icon">📋</span>
                            <span class="tab-text">收藏列表</span>
                        </button>
                        <button class="favorites-tab-btn" data-tab="signals-summary">
                            <span class="tab-icon">📊</span>
                            <span class="tab-text">信号汇总</span>
                        </button>
                    </div>
                    
                    <!-- Tab内容区域 -->
                    <div class="favorites-tab-content">
                        <!-- 收藏列表tab -->
                        <div class="favorites-tab-pane active" id="favorites-list-tab">
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
                        
                        <!-- 信号汇总tab -->
                        <div class="favorites-tab-pane" id="signals-summary-tab">
                            <div class="signals-container">
                                <div class="signals-header">
                                    <span class="signals-count">共 0 支股票</span>
                                    <button id="refreshSignalsBtn" class="btn btn-outline btn-small">🔄 刷新信号</button>
                                </div>
                                <div class="signals-list" id="signalsList">
                                    <div class="signals-empty">
                                        <p>暂无信号数据</p>
                                        <p>点击"刷新信号"按钮获取最新信号</p>
                                    </div>
                                </div>
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

        // 刷新信号按钮
        const refreshSignalsBtn = document.getElementById('refreshSignalsBtn');
        if (refreshSignalsBtn) {
            refreshSignalsBtn.addEventListener('click', () => {
                this.loadSignals();
            });
        }

        // Tab切换事件
        const tabBtns = document.querySelectorAll('.favorites-tab-btn');
        tabBtns.forEach(btn => {
            btn.addEventListener('click', () => {
                this.switchTab(btn.dataset.tab);
            });
        });

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

        try {
            // 检查在当前分组中是否已收藏
            const isFavoriteInCurrentGroup = await this.favoritesService.checkFavoriteInGroup(stockCode, this.currentGroupId);
            
            if (isFavoriteInCurrentGroup) {
                // 如果在当前分组中已收藏，从当前分组中移除
                await this.handleRemoveFromGroupClick(stockCode, this.currentGroupId);
            } else {
                // 如果在当前分组中未收藏，添加到当前分组
                await this.handleAddFavoriteClick(stockCode, startDate, endDate);
            }
        } catch (error) {
            console.error('收藏操作失败:', error);
            this.showMessage(error.message || '操作失败', 'error');
        }
    }

    /**
     * 处理添加收藏
     */
    async handleAddFavoriteClick(stockCode, startDate, endDate) {
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

            // 直接添加收藏到当前分组
            // 后端会处理分组级别的重复检查：同一分组内不能重复，不同分组间可以重复
            const groupId = this.currentGroupId || 'default';
            
            const favorite = await this.favoritesService.addFavorite(
                stockCode,
                stockName,
                startDate || '',
                endDate || '',
                groupId // 分组ID是必传参数
            );

            // 显示成功消息，包含分组信息
            const currentGroup = this.groups.find(g => g.id === this.currentGroupId);
            const groupName = currentGroup ? currentGroup.name : '默认分组';
            const message = this.currentGroupId === 'default' 
                ? `成功收藏 ${stockName}` 
                : `成功收藏 ${stockName} 到"${groupName}"分组`;
            this.showMessage(message, 'success');
            this.loadFavorites(); // 刷新收藏列表
            this.updateFavoriteButtonState();
            
        } catch (error) {
            console.error('添加收藏失败:', error);
            throw error;
        }
    }



    /**
     * 处理从特定分组中移除收藏
     */
    async handleRemoveFromGroupClick(stockCode, groupId) {
        try {
            // 通过股票代码和分组ID查找收藏记录
            const favorite = await this.favoritesService.findFavoriteByCodeAndGroup(stockCode, groupId);
            if (!favorite) {
                throw new Error('找不到对应的收藏记录');
            }

            // 删除该分组中的收藏
            await this.favoritesService.deleteFavorite(favorite.id);
            
            const currentGroup = this.groups.find(g => g.id === groupId);
            const groupName = currentGroup ? currentGroup.name : '默认分组';
            const message = groupId === 'default' 
                ? `已取消收藏 ${favorite.name}` 
                : `已从"${groupName}"分组中移除 ${favorite.name}`;
            
            this.showMessage(message, 'success');
            this.loadFavorites(); // 刷新收藏列表
            this.updateFavoriteButtonState();
            
        } catch (error) {
            console.error('移除收藏失败:', error);
            throw error;
        }
    }

    /**
     * 处理取消收藏（完全删除，保留用于兼容性）
     */
    async handleUnfavoriteClick(stockCode) {
        try {
            // 通过股票代码查找收藏记录
            const favorite = await this.favoritesService.findFavoriteByCode(stockCode);
            if (!favorite) {
                throw new Error('找不到对应的收藏记录');
            }

            // 删除收藏
            await this.favoritesService.deleteFavorite(favorite.id);
            
            this.showMessage(`已取消收藏 ${favorite.name}`, 'success');
            this.loadFavorites(); // 刷新收藏列表
            this.updateFavoriteButtonState();
            
        } catch (error) {
            console.error('取消收藏失败:', error);
            throw error;
        }
    }

    /**
     * 加载分组列表
     */
    async loadGroups() {
        try {
            this.groups = await this.favoritesService.getGroups();
            console.log('加载分组列表:', this.groups);
        } catch (error) {
            console.error('加载分组列表失败:', error);
            this.groups = [];
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

        // 按分组整理收藏
        const favoritesByGroup = this.groupFavorites();
        
        // 渲染分组tab和收藏列表
        let listHTML = '';
        
        // 分组tab导航
        listHTML += `
            <div class="group-tabs-container">
                <div class="group-tabs">
        `;
        
        // 按分组排序渲染所有分组tab
        const sortedGroups = [...this.groups].sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0));
        
        // 如果当前分组不存在，设置为第一个分组
        if (sortedGroups.length > 0 && !sortedGroups.find(g => g.id === this.currentGroupId)) {
            this.currentGroupId = sortedGroups[0].id;
        }
        
        sortedGroups.forEach((group, index) => {
            const groupFavorites = favoritesByGroup[group.id] || [];
            const isActive = group.id === this.currentGroupId ? 'active' : '';
            
            listHTML += `
                <div class="group-tab ${isActive}" data-group-id="${group.id}">
                    <div class="group-color" style="background-color: ${group.color}"></div>
                    <span class="group-name">${group.name}</span>
                    <span class="group-count">(${groupFavorites.length})</span>
                    ${group.id !== 'default' ? `
                        <div class="group-actions">
                            <button class="btn btn-icon edit-group-btn" title="编辑分组">✏️</button>
                            <button class="btn btn-icon delete-group-btn" title="删除分组">🗑️</button>
                        </div>
                    ` : ''}
                </div>
            `;
        });
        
        listHTML += `
                    <div class="group-tab-add">
                        <button class="btn btn-outline btn-small" id="createGroupBtn">
                            ➕ 新建分组
                        </button>
                    </div>
                </div>
            </div>
        `;
        
        // 分组内容区域
        listHTML += `<div class="group-content-container">`;
        
        sortedGroups.forEach((group, index) => {
            const groupFavorites = favoritesByGroup[group.id] || [];
            const isActive = group.id === this.currentGroupId ? 'active' : '';
            
            listHTML += `
                <div class="group-content ${isActive}" data-group-id="${group.id}">
            `;
            
            // 渲染该分组下的收藏
            if (groupFavorites.length > 0) {
                groupFavorites.forEach(favorite => {
                    const createdDate = new Date(favorite.created_at).toLocaleDateString();
                    const dateRange = this.formatDateRange(favorite.start_date, favorite.end_date);
                    
                    listHTML += `
                        <div class="favorite-item" data-favorite-id="${favorite.id}" data-stock-code="${favorite.ts_code}" data-group-id="${favorite.group_id || 'default'}" draggable="true">
                            <div class="drag-handle" title="拖拽排序">⋮⋮</div>
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
            } else {
                // 空分组提示
                listHTML += `
                    <div class="empty-group-hint">
                        <p>该分组暂无收藏股票</p>
                        <p>可以拖拽其他收藏到这里进行分组</p>
                    </div>
                `;
            }
            
            listHTML += `</div>`;
        });
        
        listHTML += `</div>`;

        favoritesList.innerHTML = listHTML;

        // 添加事件监听器
        this.setupFavoriteItemEvents();
        this.setupGroupEvents();
        
        // 更新收藏按钮状态
        this.updateFavoriteButtonState();
        this.setupTabEvents();
        this.setupDragAndDrop();
    }

    /**
     * 设置分组事件监听器
     */
    setupGroupEvents() {
        // 新建分组按钮
        const createGroupBtn = document.getElementById('createGroupBtn');
        if (createGroupBtn) {
            createGroupBtn.addEventListener('click', () => {
                this.showCreateGroupDialog();
            });
        }

        // 编辑分组按钮
        document.querySelectorAll('.edit-group-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                const groupTab = e.target.closest('.group-tab');
                const groupId = groupTab.dataset.groupId;
                this.showEditGroupDialog(groupId);
            });
        });

        // 删除分组按钮
        document.querySelectorAll('.delete-group-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                const groupTab = e.target.closest('.group-tab');
                const groupId = groupTab.dataset.groupId;
                this.deleteGroup(groupId);
            });
        });
    }

    /**
     * 设置tab切换事件监听器
     */
    setupTabEvents() {
        // tab点击切换
        document.querySelectorAll('.group-tab').forEach(tab => {
            tab.addEventListener('click', (e) => {
                // 如果点击的是按钮，不处理tab切换
                if (e.target.closest('.group-actions')) {
                    return;
                }
                
                const groupId = tab.dataset.groupId;
                this.switchToTab(groupId);
            });
        });
    }

    /**
     * 切换到指定的tab
     */
    switchToTab(groupId) {
        // 更新当前分组ID
        this.currentGroupId = groupId;
        
        // 移除所有active状态
        document.querySelectorAll('.group-tab').forEach(tab => {
            tab.classList.remove('active');
        });
        document.querySelectorAll('.group-content').forEach(content => {
            content.classList.remove('active');
        });
        
        // 添加active状态到指定的tab和内容
        const targetTab = document.querySelector(`.group-tab[data-group-id="${groupId}"]`);
        const targetContent = document.querySelector(`.group-content[data-group-id="${groupId}"]`);
        
        if (targetTab) {
            targetTab.classList.add('active');
        }
        if (targetContent) {
            targetContent.classList.add('active');
        }
        
        // 更新收藏按钮状态
        this.updateFavoriteButtonState();
    }

    /**
     * 显示创建分组对话框
     */
    showCreateGroupDialog() {
        const name = prompt('请输入分组名称:');
        if (name && name.trim()) {
            this.createGroup(name.trim());
        }
    }

    /**
     * 显示编辑分组对话框
     */
    showEditGroupDialog(groupId) {
        const group = this.groups.find(g => g.id === groupId);
        if (!group) return;

        const name = prompt('请输入新的分组名称:', group.name);
        if (name && name.trim() && name.trim() !== group.name) {
            this.updateGroup(groupId, name.trim(), group.color);
        }
    }

    /**
     * 创建分组
     */
    async createGroup(name) {
        try {
            await this.favoritesService.createGroup(name);
            this.showMessage(`成功创建分组 "${name}"`, 'success');
            this.loadGroups();
            this.loadFavorites();
        } catch (error) {
            console.error('创建分组失败:', error);
            this.showMessage(error.message || '创建分组失败', 'error');
        }
    }

    /**
     * 更新分组
     */
    async updateGroup(groupId, name, color) {
        try {
            await this.favoritesService.updateGroup(groupId, name, color);
            this.showMessage(`成功更新分组`, 'success');
            this.loadGroups();
            this.loadFavorites();
        } catch (error) {
            console.error('更新分组失败:', error);
            this.showMessage(error.message || '更新分组失败', 'error');
        }
    }

    /**
     * 删除分组
     */
    async deleteGroup(groupId) {
        const group = this.groups.find(g => g.id === groupId);
        if (!group) return;

        if (!confirm(`确定要删除分组 "${group.name}" 吗？\n该分组下的收藏将移动到默认分组。`)) {
            return;
        }

        try {
            await this.favoritesService.deleteGroup(groupId);
            this.showMessage(`成功删除分组 "${group.name}"`, 'success');
            this.loadGroups();
            this.loadFavorites();
        } catch (error) {
            console.error('删除分组失败:', error);
            this.showMessage(error.message || '删除分组失败', 'error');
        }
    }

    /**
     * 设置拖拽排序功能
     */
    setupDragAndDrop() {
        let draggedElement = null;
        let draggedData = null;

        // 为所有收藏项添加拖拽事件
        document.querySelectorAll('.favorite-item').forEach(item => {
            item.addEventListener('dragstart', (e) => {
                draggedElement = e.target;
                draggedData = {
                    favoriteId: item.dataset.favoriteId,
                    groupId: item.dataset.groupId,
                    stockCode: item.dataset.stockCode
                };
                
                e.target.style.opacity = '0.5';
                e.dataTransfer.effectAllowed = 'move';
            });

            item.addEventListener('dragend', (e) => {
                e.target.style.opacity = '1';
                draggedElement = null;
                draggedData = null;
            });

            item.addEventListener('dragover', (e) => {
                e.preventDefault();
                e.dataTransfer.dropEffect = 'move';
            });

            item.addEventListener('drop', (e) => {
                e.preventDefault();
                
                if (!draggedData || e.target === draggedElement) return;

                const targetItem = e.target.closest('.favorite-item');
                if (!targetItem) return;

                const targetData = {
                    favoriteId: targetItem.dataset.favoriteId,
                    groupId: targetItem.dataset.groupId
                };

                this.handleDrop(draggedData, targetData);
            });
        });

        // 为分组tab添加拖拽支持（可以拖拽到tab上切换分组）
        document.querySelectorAll('.group-tab').forEach(tab => {
            tab.addEventListener('dragover', (e) => {
                e.preventDefault();
                e.dataTransfer.dropEffect = 'move';
                // 高亮显示可放置的tab
                tab.classList.add('drag-over');
            });

            tab.addEventListener('dragleave', (e) => {
                // 移除高亮
                tab.classList.remove('drag-over');
            });

            tab.addEventListener('drop', (e) => {
                e.preventDefault();
                tab.classList.remove('drag-over');
                
                if (!draggedData) return;

                const targetGroupId = tab.dataset.groupId;
                
                // 如果拖拽到不同的分组
                if (draggedData.groupId !== targetGroupId) {
                    this.moveFavoriteToGroup(draggedData.favoriteId, targetGroupId);
                }
            });
        });

        // 为分组内容区域添加拖拽支持
        document.querySelectorAll('.group-content').forEach(groupContent => {
            groupContent.addEventListener('dragover', (e) => {
                e.preventDefault();
                e.dataTransfer.dropEffect = 'move';
            });

            groupContent.addEventListener('drop', (e) => {
                e.preventDefault();
                
                if (!draggedData) return;

                // 检查是否拖拽到了空分组提示区域
                const emptyHint = e.target.closest('.empty-group-hint');
                if (emptyHint) {
                    const targetGroupId = groupContent.dataset.groupId;
                    
                    // 如果拖拽到不同的分组
                    if (draggedData.groupId !== targetGroupId) {
                        this.moveFavoriteToGroup(draggedData.favoriteId, targetGroupId);
                    }
                }
            });
        });
    }

    /**
     * 处理拖拽放置
     */
    async handleDrop(draggedData, targetData) {
        try {
            // 获取当前的收藏列表
            const favorites = [...this.favorites];
            
            // 找到拖拽的收藏和目标收藏
            const draggedFavorite = favorites.find(f => f.id === draggedData.favoriteId);
            const targetFavorite = favorites.find(f => f.id === targetData.favoriteId);
            
            if (!draggedFavorite || !targetFavorite) return;

            // 如果是同一个分组内的排序
            if (draggedData.groupId === targetData.groupId) {
                await this.reorderWithinGroup(draggedData.favoriteId, targetData.favoriteId, targetData.groupId);
            } else {
                // 跨分组移动
                await this.moveFavoriteToGroup(draggedData.favoriteId, targetData.groupId);
            }
        } catch (error) {
            console.error('拖拽操作失败:', error);
            this.showMessage('排序失败', 'error');
        }
    }

    /**
     * 在分组内重新排序
     */
    async reorderWithinGroup(draggedId, targetId, groupId) {
        try {
            // 获取该分组内的所有收藏
            const groupFavorites = this.favorites.filter(f => (f.group_id || 'default') === groupId);
            groupFavorites.sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0));

            // 找到拖拽项和目标项的位置
            const draggedIndex = groupFavorites.findIndex(f => f.id === draggedId);
            const targetIndex = groupFavorites.findIndex(f => f.id === targetId);

            if (draggedIndex === -1 || targetIndex === -1) return;

            // 重新排序
            const [draggedItem] = groupFavorites.splice(draggedIndex, 1);
            groupFavorites.splice(targetIndex, 0, draggedItem);

            // 生成新的排序数据
            const orderUpdates = groupFavorites.map((favorite, index) => ({
                id: favorite.id,
                group_id: groupId,
                sort_order: index + 1
            }));

            // 更新排序
            await this.favoritesService.updateFavoritesOrder(orderUpdates);
            this.showMessage('排序已更新', 'success');
            this.loadFavorites();
        } catch (error) {
            console.error('分组内排序失败:', error);
            throw error;
        }
    }

    /**
     * 移动收藏到其他分组
     */
    async moveFavoriteToGroup(favoriteId, targetGroupId) {
        try {
            const favorite = this.favorites.find(f => f.id === favoriteId);
            if (!favorite) return;

            // 获取目标分组的最大排序号
            const targetGroupFavorites = this.favorites.filter(f => (f.group_id || 'default') === targetGroupId);
            const maxSortOrder = Math.max(...targetGroupFavorites.map(f => f.sort_order || 0), 0);

            // 更新收藏的分组和排序
            const orderUpdate = [{
                id: favoriteId,
                group_id: targetGroupId,
                sort_order: maxSortOrder + 1
            }];

            await this.favoritesService.updateFavoritesOrder(orderUpdate);
            
            const targetGroup = this.groups.find(g => g.id === targetGroupId) || { name: '默认分组' };
            this.showMessage(`已移动到 "${targetGroup.name}"`, 'success');
            this.loadFavorites();
        } catch (error) {
            console.error('移动收藏失败:', error);
            throw error;
        }
    }

    /**
     * 按分组整理收藏
     */
    groupFavorites() {
        const favoritesByGroup = {};
        
        this.favorites.forEach(favorite => {
            const groupId = favorite.group_id || 'default';
            if (!favoritesByGroup[groupId]) {
                favoritesByGroup[groupId] = [];
            }
            favoritesByGroup[groupId].push(favorite);
        });

        // 按排序字段排序每个分组内的收藏
        Object.keys(favoritesByGroup).forEach(groupId => {
            favoritesByGroup[groupId].sort((a, b) => {
                const aOrder = a.sort_order || 0;
                const bOrder = b.sort_order || 0;
                return aOrder - bOrder;
            });
        });

        return favoritesByGroup;
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

            // 更新收藏按钮状态（稍微延迟以确保输入框值已更新）
            setTimeout(() => {
                this.updateFavoriteButtonState();
            }, 10);

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
        
        if (!favoriteBtn || !stockCode) {
            // 如果没有股票代码，重置按钮状态
            if (favoriteBtn) {
                favoriteBtn.innerHTML = '⭐ 收藏';
                favoriteBtn.className = 'btn btn-favorite';
                favoriteBtn.disabled = false;
                favoriteBtn.title = '收藏当前股票和时间范围';
            }
            return;
        }

        try {
            // 检查在当前分组中是否已收藏
            const isFavoriteInCurrentGroup = await this.favoritesService.checkFavoriteInGroup(stockCode, this.currentGroupId);
            
            const currentGroup = this.groups.find(g => g.id === this.currentGroupId);
            const groupName = currentGroup ? currentGroup.name : '默认分组';
            
            if (isFavoriteInCurrentGroup) {
                // 在当前分组中已收藏
                if (this.currentGroupId === 'default') {
                    favoriteBtn.innerHTML = '⭐ 已收藏';
                    favoriteBtn.title = '点击取消收藏该股票';
                } else {
                    favoriteBtn.innerHTML = `⭐ 已在${groupName}中`;
                    favoriteBtn.title = `点击从"${groupName}"分组中移除该股票`;
                }
                favoriteBtn.className = 'btn btn-favorite favorited';
                favoriteBtn.disabled = false;
            } else {
                // 在当前分组中未收藏
                if (this.currentGroupId === 'default') {
                    favoriteBtn.innerHTML = '⭐ 收藏';
                    favoriteBtn.title = '收藏当前股票和时间范围';
                } else {
                    favoriteBtn.innerHTML = `⭐ 收藏到${groupName}`;
                    favoriteBtn.title = `收藏当前股票和时间范围到"${groupName}"分组`;
                }
                favoriteBtn.className = 'btn btn-favorite';
                favoriteBtn.disabled = false;
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

    /**
     * 切换Tab
     */
    switchTab(tabName) {
        // 更新tab按钮状态
        const tabBtns = document.querySelectorAll('.favorites-tab-btn');
        tabBtns.forEach(btn => {
            if (btn.dataset.tab === tabName) {
                btn.classList.add('active');
            } else {
                btn.classList.remove('active');
            }
        });

        // 更新tab内容显示
        const tabPanes = document.querySelectorAll('.favorites-tab-pane');
        tabPanes.forEach(pane => {
            if (pane.id === `${tabName}-tab`) {
                pane.classList.add('active');
            } else {
                pane.classList.remove('active');
            }
        });

        // 如果切换到信号汇总tab，自动加载信号
        if (tabName === 'signals-summary') {
            this.loadSignals();
        }
    }

    /**
     * 加载信号汇总数据
     */
    async loadSignals() {
        try {
            const signalsList = document.getElementById('signalsList');
            const signalsCount = document.querySelector('.signals-count');
            
            // 显示加载状态
            signalsList.innerHTML = '<div class="loading">正在加载信号数据...</div>';
            
            // 调用API获取信号数据
            const response = await this.apiService.getFavoritesSignals();
            
            if (response && response.signals) {
                this.renderSignals(response.signals);
                signalsCount.textContent = `共 ${response.total} 支股票`;
            } else {
                signalsList.innerHTML = '<div class="signals-empty"><p>暂无信号数据</p></div>';
                signalsCount.textContent = '共 0 支股票';
            }
        } catch (error) {
            console.error('加载信号数据失败:', error);
            const signalsList = document.getElementById('signalsList');
            signalsList.innerHTML = '<div class="signals-error"><p>加载信号数据失败</p><p>错误信息: ' + error.message + '</p></div>';
        }
    }

    /**
     * 渲染信号汇总列表
     */
    renderSignals(signals) {
        const signalsList = document.getElementById('signalsList');
        
        if (!signals || signals.length === 0) {
            signalsList.innerHTML = '<div class="signals-empty"><p>暂无信号数据</p></div>';
            return;
        }

        let signalsHTML = '';
        
        // 按信号类型分组
        const buySignals = [];
        const sellSignals = [];
        const holdSignals = [];

        signals.forEach(signal => {
            // 分析预测信号
            let hasBuySignal = false;
            let hasSellSignal = false;
            
            if (signal.predictions && signal.predictions.predictions) {
                signal.predictions.predictions.forEach(prediction => {
                    if (prediction.type === 'BUY') {
                        hasBuySignal = true;
                    } else if (prediction.type === 'SELL') {
                        hasSellSignal = true;
                    }
                });
            }

            if (hasBuySignal) {
                buySignals.push(signal);
            } else if (hasSellSignal) {
                sellSignals.push(signal);
            } else {
                holdSignals.push(signal);
            }
        });

        // 渲染买入信号
        if (buySignals.length > 0) {
            signalsHTML += `
                <div class="signal-group buy-signals">
                    <h3 class="signal-group-title buy">🟢 买入信号 (${buySignals.length})</h3>
                    ${this.renderSignalGroup(buySignals, 'buy')}
                </div>
            `;
        }

        // 渲染卖出信号
        if (sellSignals.length > 0) {
            signalsHTML += `
                <div class="signal-group sell-signals">
                    <h3 class="signal-group-title sell">🔴 卖出信号 (${sellSignals.length})</h3>
                    ${this.renderSignalGroup(sellSignals, 'sell')}
                </div>
            `;
        }

        // 渲染持有信号
        if (holdSignals.length > 0) {
            signalsHTML += `
                <div class="signal-group hold-signals">
                    <h3 class="signal-group-title hold">🟡 持有信号 (${holdSignals.length})</h3>
                    ${this.renderSignalGroup(holdSignals, 'hold')}
                </div>
            `;
        }

        signalsList.innerHTML = signalsHTML;
    }

    /**
     * 渲染信号组
     */
    renderSignalGroup(signals, type) {
        return signals.map(signal => {
            const currentPrice = signal.current_price || 'N/A';
            const tradeDate = signal.trade_date || 'N/A';
            const updatedAt = signal.updated_at || 'N/A';
            
            // 获取主要信号
            let mainSignal = 'HOLD';
            let signalReason = '';
            let signalProbability = '';
            
            if (signal.predictions && signal.predictions.predictions) {
                const prediction = signal.predictions.predictions[0]; // 取第一个预测
                if (prediction) {
                    mainSignal = prediction.type;
                    signalReason = prediction.reason;
                    signalProbability = prediction.probability;
                }
            }

            return `
                <div class="signal-item ${type}-signal" data-stock-code="${signal.ts_code}">
                    <div class="signal-header">
                        <div class="signal-stock-info">
                            <span class="signal-stock-name">${signal.name}</span>
                            <span class="signal-stock-code">${signal.ts_code}</span>
                        </div>
                        <div class="signal-price">
                            <span class="current-price">¥${currentPrice}</span>
                            <span class="trade-date">${tradeDate}</span>
                        </div>
                    </div>
                    <div class="signal-details">
                        <div class="signal-main">
                            <span class="signal-type ${mainSignal.toLowerCase()}">${this.getSignalText(mainSignal)}</span>
                            ${signalProbability ? `<span class="signal-probability">置信度: ${signalProbability}</span>` : ''}
                        </div>
                        ${signalReason ? `<div class="signal-reason">${signalReason}</div>` : ''}
                    </div>
                    <div class="signal-actions">
                        <button class="btn btn-outline btn-small view-chart-btn" title="查看K线图">
                            📈 查看
                        </button>
                        <button class="btn btn-outline btn-small view-details-btn" title="查看详细分析">
                            📊 详情
                        </button>
                    </div>
                </div>
            `;
        }).join('');
    }

    /**
     * 获取信号文本
     */
    getSignalText(signal) {
        const signalMap = {
            'BUY': '买入',
            'SELL': '卖出',
            'HOLD': '持有'
        };
        return signalMap[signal] || signal;
    }
}

// 导出收藏模块类
window.FavoritesModule = FavoritesModule;
