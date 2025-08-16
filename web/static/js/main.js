/**
 * Stock-A-Future API 网页客户端主入口文件
 * 负责初始化所有模块和协调它们之间的交互
 */

class StockAFutureApp {
    constructor() {
        this.client = null;
        this.apiService = null;
        this.favoritesService = null;
        this.stockSearchModule = null;
        this.chartsModule = null;
        this.favoritesModule = null;
        this.configModule = null;
        this.displayModule = null;
        this.eventsModule = null;
        this.dateShortcutsModule = null;
        
        this.init();
    }

    /**
     * 初始化应用
     */
    async init() {
        try {
            console.log('正在初始化 Stock-A-Future 应用...');
            
            // 等待DOM加载完成
            if (document.readyState === 'loading') {
                await new Promise(resolve => {
                    document.addEventListener('DOMContentLoaded', resolve);
                });
            }
            
            // 初始化各个模块
            await this.initializeModules();
            
            // 设置默认值
            this.setDefaultValues();
            
            // 启动定期健康检查
            this.startHealthCheck();
            
            console.log('Stock-A-Future 应用初始化完成');
            
        } catch (error) {
            console.error('应用初始化失败:', error);
        }
    }

    /**
     * 初始化各个模块
     */
    async initializeModules() {
        // 1. 初始化核心客户端
        this.client = new StockAFutureClient();
        
        // 2. 初始化API服务
        this.apiService = new ApiService(this.client);
        
        // 3. 初始化收藏服务
        this.favoritesService = new FavoritesService(this.client);
        
        // 4. 初始化图表模块
        this.chartsModule = new ChartsModule(this.client);
        
        // 5. 初始化数据展示模块
        this.displayModule = new DisplayModule(this.client, this.chartsModule);
        
        // 6. 初始化股票搜索模块
        this.stockSearchModule = new StockSearchModule(this.client, this.apiService);
        
        // 7. 初始化收藏功能模块
        this.favoritesModule = new FavoritesModule(this.client, this.favoritesService, this.apiService);
        
        // 8. 初始化日期快捷按钮模块
        this.dateShortcutsModule = new DateShortcutsModule(this.client);
        
        // 9. 初始化事件处理模块（需要收藏模块和日期快捷按钮模块引用）
        this.eventsModule = new EventsModule(this.client, this.apiService, this.displayModule, this.favoritesModule, this.dateShortcutsModule);
        
        // 10. 初始化配置管理模块
        this.configModule = new ConfigModule(this.client);
        
        // 等待客户端初始化完成
        await this.waitForClientReady();
    }

    /**
     * 等待客户端准备就绪
     */
    async waitForClientReady() {
        let attempts = 0;
        const maxAttempts = 50; // 最多等待5秒
        
        while (!this.client.baseURL && attempts < maxAttempts) {
            await new Promise(resolve => setTimeout(resolve, 100));
            attempts++;
        }
        
        if (!this.client.baseURL) {
            console.warn('客户端初始化超时，使用默认配置');
        }
    }

    /**
     * 设置默认值
     */
    setDefaultValues() {
        // 设置默认股票代码（从缓存读取）
        if (this.stockSearchModule) {
            this.stockSearchModule.setDefaultStockCode();
        }
    }

    /**
     * 启动定期健康检查
     */
    startHealthCheck() {
        // 每5分钟检查一次健康状态，减少不必要的连接测试
        setInterval(() => {
            if (this.client && !this.client.isLoading) {
                this.client.checkHealth();
            }
        }, 300000); // 5分钟 = 300秒 = 300000毫秒
    }

    /**
     * 获取客户端实例
     */
    getClient() {
        return this.client;
    }

    /**
     * 获取API服务实例
     */
    getApiService() {
        return this.apiService;
    }

    /**
     * 获取图表模块实例
     */
    getChartsModule() {
        return this.chartsModule;
    }

    /**
     * 获取数据展示模块实例
     */
    getDisplayModule() {
        return this.displayModule;
    }

    /**
     * 获取事件处理模块实例
     */
    getEventsModule() {
        return this.eventsModule;
    }

    /**
     * 获取股票搜索模块实例
     */
    getStockSearchModule() {
        return this.stockSearchModule;
    }

    /**
     * 获取配置管理模块实例
     */
    getConfigModule() {
        return this.configModule;
    }

    /**
     * 获取收藏服务实例
     */
    getFavoritesService() {
        return this.favoritesService;
    }

    /**
     * 获取收藏模块实例
     */
    getFavoritesModule() {
        return this.favoritesModule;
    }

    /**
     * 获取日期快捷按钮模块实例
     */
    getDateShortcutsModule() {
        return this.dateShortcutsModule;
    }
}

// 全局应用实例
let stockApp = null;

// 页面加载完成后初始化应用
document.addEventListener('DOMContentLoaded', async () => {
    try {
        // 创建应用实例
        stockApp = new StockAFutureApp();
        
        // 将应用实例挂载到全局对象，方便调试
        window.stockApp = stockApp;
        
        // 将关键模块挂载到全局对象，方便其他模块调用
        if (stockApp.eventsModule) {
            window.eventsModule = stockApp.eventsModule;
        }
        
        console.log('Stock-A-Future 网页客户端已初始化');
        
    } catch (error) {
        console.error('Stock-A-Future 网页客户端初始化失败:', error);
    }
});

// 导出应用类（如果需要）
window.StockAFutureApp = StockAFutureApp;
