/**
 * 工具函数模块
 * 包含一些通用的辅助函数
 */

class Helpers {
    /**
     * 防抖函数
     * @param {Function} func 要防抖的函数
     * @param {number} wait 等待时间（毫秒）
     * @returns {Function} 防抖后的函数
     */
    static debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    /**
     * 节流函数
     * @param {Function} func 要节流的函数
     * @param {number} limit 限制时间（毫秒）
     * @returns {Function} 节流后的函数
     */
    static throttle(func, limit) {
        let inThrottle;
        return function() {
            const args = arguments;
            const context = this;
            if (!inThrottle) {
                func.apply(context, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    }

    /**
     * 格式化数字，添加千分位分隔符
     * @param {number} num 要格式化的数字
     * @returns {string} 格式化后的字符串
     */
    static formatNumber(num) {
        return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
    }

    /**
     * 格式化货币
     * @param {number} amount 金额
     * @param {string} currency 货币符号
     * @returns {string} 格式化后的货币字符串
     */
    static formatCurrency(amount, currency = '¥') {
        return `${currency}${this.formatNumber(amount.toFixed(2))}`;
    }

    /**
     * 格式化百分比
     * @param {number} value 值
     * @param {number} decimals 小数位数
     * @returns {string} 格式化后的百分比字符串
     */
    static formatPercentage(value, decimals = 2) {
        return `${(value * 100).toFixed(decimals)}%`;
    }

    /**
     * 深拷贝对象
     * @param {*} obj 要拷贝的对象
     * @returns {*} 拷贝后的对象
     */
    static deepClone(obj) {
        if (obj === null || typeof obj !== 'object') return obj;
        if (obj instanceof Date) return new Date(obj.getTime());
        if (obj instanceof Array) return obj.map(item => this.deepClone(item));
        if (typeof obj === 'object') {
            const clonedObj = {};
            for (const key in obj) {
                if (obj.hasOwnProperty(key)) {
                    clonedObj[key] = this.deepClone(obj[key]);
                }
            }
            return clonedObj;
        }
    }

    /**
     * 生成唯一ID
     * @returns {string} 唯一ID
     */
    static generateId() {
        return Date.now().toString(36) + Math.random().toString(36).substr(2);
    }

    /**
     * 检查元素是否在视口中
     * @param {Element} element 要检查的元素
     * @returns {boolean} 是否在视口中
     */
    static isInViewport(element) {
        const rect = element.getBoundingClientRect();
        return (
            rect.top >= 0 &&
            rect.left >= 0 &&
            rect.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
            rect.right <= (window.innerWidth || document.documentElement.clientWidth)
        );
    }

    /**
     * 平滑滚动到元素
     * @param {Element|string} target 目标元素或选择器
     * @param {Object} options 滚动选项
     */
    static smoothScrollTo(target, options = {}) {
        const element = typeof target === 'string' ? document.querySelector(target) : target;
        if (!element) return;

        const defaultOptions = {
            behavior: 'smooth',
            block: 'start',
            inline: 'nearest'
        };

        element.scrollIntoView({ ...defaultOptions, ...options });
    }

    /**
     * 添加CSS类（如果不存在）
     * @param {Element} element 目标元素
     * @param {string} className 类名
     */
    static addClassIfNotExists(element, className) {
        if (element && !element.classList.contains(className)) {
            element.classList.add(className);
        }
    }

    /**
     * 移除CSS类（如果存在）
     * @param {Element} element 目标元素
     * @param {string} className 类名
     */
    static removeClassIfExists(element, className) {
        if (element && element.classList.contains(className)) {
            element.classList.remove(className);
        }
    }

    /**
     * 切换CSS类
     * @param {Element} element 目标元素
     * @param {string} className 类名
     */
    static toggleClass(element, className) {
        if (element) {
            element.classList.toggle(className);
        }
    }

    /**
     * 获取元素的计算样式
     * @param {Element} element 目标元素
     * @param {string} property 样式属性
     * @returns {string} 计算后的样式值
     */
    static getComputedStyle(element, property) {
        if (!element) return '';
        return window.getComputedStyle(element).getPropertyValue(property);
    }

    /**
     * 设置元素样式
     * @param {Element} element 目标元素
     * @param {Object} styles 样式对象
     */
    static setStyles(element, styles) {
        if (!element) return;
        Object.assign(element.style, styles);
    }

    /**
     * 获取URL参数
     * @param {string} name 参数名
     * @returns {string|null} 参数值
     */
    static getUrlParam(name) {
        const urlParams = new URLSearchParams(window.location.search);
        return urlParams.get(name);
    }

    /**
     * 设置URL参数
     * @param {string} name 参数名
     * @param {string} value 参数值
     */
    static setUrlParam(name, value) {
        const url = new URL(window.location);
        url.searchParams.set(name, value);
        window.history.replaceState({}, '', url);
    }

    /**
     * 移除URL参数
     * @param {string} name 参数名
     */
    static removeUrlParam(name) {
        const url = new URL(window.location);
        url.searchParams.delete(name);
        window.history.replaceState({}, '', url);
    }

    /**
     * 检查是否为移动设备
     * @returns {boolean} 是否为移动设备
     */
    static isMobile() {
        return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    }

    /**
     * 检查是否为触摸设备
     * @returns {boolean} 是否为触摸设备
     */
    static isTouchDevice() {
        return 'ontouchstart' in window || navigator.maxTouchPoints > 0;
    }

    /**
     * 延迟执行函数
     * @param {Function} func 要执行的函数
     * @param {number} delay 延迟时间（毫秒）
     * @returns {Promise} Promise对象
     */
    static delay(func, delay) {
        return new Promise(resolve => {
            setTimeout(() => {
                func();
                resolve();
            }, delay);
        });
    }

    /**
     * 异步等待
     * @param {number} ms 等待时间（毫秒）
     * @returns {Promise} Promise对象
     */
    static sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}

// 导出工具函数类
window.Helpers = Helpers;
