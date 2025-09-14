/**
 * 安全调用模块方法的工具函数
 * 用于处理模块可能还没有完全初始化的情况
 */

window.safeCall = function(moduleName, methodName, ...args) {
    const module = window[moduleName];
    
    if (!module) {
        console.warn(`模块 ${moduleName} 尚未初始化，正在等待...`);
        
        // 等待模块初始化，最多等待5秒
        let attempts = 0;
        const maxAttempts = 50; // 5秒，每100ms检查一次
        
        const checkModule = () => {
            attempts++;
            const currentModule = window[moduleName];
            
            if (currentModule && typeof currentModule[methodName] === 'function') {
                console.log(`模块 ${moduleName} 已初始化，执行方法 ${methodName}`);
                currentModule[methodName](...args);
            } else if (attempts < maxAttempts) {
                setTimeout(checkModule, 100);
            } else {
                console.error(`模块 ${moduleName} 初始化超时，方法 ${methodName} 无法执行`);
                // 显示用户友好的错误信息
                if (window.showMessage) {
                    window.showMessage('功能加载中，请稍后再试', 'warning');
                } else {
                    alert('功能加载中，请稍后再试');
                }
            }
        };
        
        setTimeout(checkModule, 100);
        return;
    }
    
    if (typeof module[methodName] !== 'function') {
        console.error(`模块 ${moduleName} 中不存在方法 ${methodName}`);
        return;
    }
    
    // 直接调用方法
    try {
        module[methodName](...args);
    } catch (error) {
        console.error(`调用 ${moduleName}.${methodName} 时发生错误:`, error);
    }
};

// 专门为策略模块提供的安全调用函数
window.safeCallStrategy = function(methodName, ...args) {
    return window.safeCall('strategiesModule', methodName, ...args);
};

// 专门为回测模块提供的安全调用函数
window.safeCallBacktest = function(methodName, ...args) {
    return window.safeCall('backtestModule', methodName, ...args);
};

console.log('[工具] 安全调用工具已加载');
