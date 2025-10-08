/**
 * 应用配置管理
 * 从环境变量或默认值中获取配置
 */

// 开发环境默认配置
const DEFAULT_CONFIG = {
  // 后端API配置
  API_BASE_URL: 'http://127.0.0.1:8081/api/v1/',
  API_TIMEOUT: 30000,
  
  // AKTools配置
  AKTOOLS_BASE_URL: 'http://127.0.0.1:8080',
  
  // 应用配置
  APP_NAME: 'Stock-A-Future Mobile',
  APP_VERSION: '1.0.0',
  
  // 调试配置
  DEBUG_MODE: false,
  LOG_LEVEL: 'info',
  
  // 缓存配置
  CACHE_TIMEOUT: 5 * 60 * 1000, // 5分钟
  MAX_RETRIES: 3,
  
  // 刷新间隔
  REFRESH_INTERVAL: 60, // 60秒
  
  // 默认技术指标
  DEFAULT_INDICATORS: ['MA', 'MACD', 'RSI'],
} as const;

/**
 * 获取配置值
 * 优先从环境变量获取，如果不存在则使用默认值
 */
class Config {
  private static instance: Config;
  private config: typeof DEFAULT_CONFIG;

  private constructor() {
    this.config = { ...DEFAULT_CONFIG };
    this.loadFromEnvironment();
  }

  public static getInstance(): Config {
    if (!Config.instance) {
      Config.instance = new Config();
    }
    return Config.instance;
  }

  /**
   * 从环境变量加载配置
   * React Native 中可以通过 react-native-config 或其他方式获取环境变量
   */
  private loadFromEnvironment() {
    try {
      // 使用react-native-config从.env文件加载环境变量
      // 注意：react-native-config需要原生配置，如果没有配置则使用默认值
      
      // 尝试动态导入react-native-config
      // 如果导入失败，则使用默认配置
      const Config = require('react-native-config');
      
      if (Config && Config.default) {
        const envConfig = Config.default;
        
        // 更新配置，如果环境变量存在的话
        if (envConfig.API_BASE_URL) {
          this.config.API_BASE_URL = envConfig.API_BASE_URL;
        }
        if (envConfig.API_TIMEOUT) {
          this.config.API_TIMEOUT = parseInt(envConfig.API_TIMEOUT, 10) || DEFAULT_CONFIG.API_TIMEOUT;
        }
        if (envConfig.AKTOOLS_BASE_URL) {
          this.config.AKTOOLS_BASE_URL = envConfig.AKTOOLS_BASE_URL;
        }
        if (envConfig.DEBUG_MODE) {
          this.config.DEBUG_MODE = envConfig.DEBUG_MODE === 'true';
        }
        if (envConfig.LOG_LEVEL) {
          this.config.LOG_LEVEL = envConfig.LOG_LEVEL;
        }
        if (envConfig.CACHE_TIMEOUT) {
          this.config.CACHE_TIMEOUT = parseInt(envConfig.CACHE_TIMEOUT, 10) || DEFAULT_CONFIG.CACHE_TIMEOUT;
        }
        if (envConfig.MAX_RETRIES) {
          this.config.MAX_RETRIES = parseInt(envConfig.MAX_RETRIES, 10) || DEFAULT_CONFIG.MAX_RETRIES;
        }
        if (envConfig.REFRESH_INTERVAL) {
          this.config.REFRESH_INTERVAL = parseInt(envConfig.REFRESH_INTERVAL, 10) || DEFAULT_CONFIG.REFRESH_INTERVAL;
        }
        if (envConfig.DEFAULT_INDICATORS) {
          this.config.DEFAULT_INDICATORS = envConfig.DEFAULT_INDICATORS.split(',').map(s => s.trim());
        }
      }
    } catch (error) {
      // 如果react-native-config不可用，则使用默认配置
      console.warn('react-native-config not available, using default configuration:', error);
    }
  }

  /**
   * 获取API基础URL
   */
  public get apiBaseUrl(): string {
    return this.config.API_BASE_URL;
  }

  /**
   * 获取API超时时间
   */
  public get apiTimeout(): number {
    return this.config.API_TIMEOUT;
  }

  /**
   * 获取AKTools基础URL
   */
  public get aktoolsBaseUrl(): string {
    return this.config.AKTOOLS_BASE_URL;
  }

  /**
   * 获取应用名称
   */
  public get appName(): string {
    return this.config.APP_NAME;
  }

  /**
   * 获取应用版本
   */
  public get appVersion(): string {
    return this.config.APP_VERSION;
  }

  /**
   * 获取调试模式状态
   */
  public get debugMode(): boolean {
    return this.config.DEBUG_MODE;
  }

  /**
   * 获取日志级别
   */
  public get logLevel(): string {
    return this.config.LOG_LEVEL;
  }

  /**
   * 获取缓存超时时间
   */
  public get cacheTimeout(): number {
    return this.config.CACHE_TIMEOUT;
  }

  /**
   * 获取最大重试次数
   */
  public get maxRetries(): number {
    return this.config.MAX_RETRIES;
  }

  /**
   * 获取刷新间隔
   */
  public get refreshInterval(): number {
    return this.config.REFRESH_INTERVAL;
  }

  /**
   * 获取默认技术指标
   */
  public get defaultIndicators(): string[] {
    return [...this.config.DEFAULT_INDICATORS];
  }

  /**
   * 动态更新配置
   */
  public updateConfig(newConfig: Partial<typeof DEFAULT_CONFIG>) {
    this.config = { ...this.config, ...newConfig };
  }

  /**
   * 重置为默认配置
   */
  public resetToDefaults() {
    this.config = { ...DEFAULT_CONFIG };
  }

  /**
   * 获取所有配置
   */
  public getAllConfig() {
    return { ...this.config };
  }
}

// 导出单例实例
export const appConfig = Config.getInstance();

// 导出配置类型
export type AppConfig = typeof DEFAULT_CONFIG;

// 导出默认配置（用于测试或重置）
export { DEFAULT_CONFIG };
