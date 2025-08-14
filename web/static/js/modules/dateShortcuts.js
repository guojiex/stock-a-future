/**
 * 日期快捷按钮模块
 * 负责处理时间段快捷选择功能
 */

class DateShortcutsModule {
    constructor(client) {
        this.client = client;
        this.shortcutButtons = [];
        this.activeButton = null;
        
        this.init();
    }

    /**
     * 初始化日期快捷按钮模块
     */
    init() {
        this.setupShortcutButtons();
        this.setupEventListeners();
    }

    /**
     * 设置快捷按钮
     */
    setupShortcutButtons() {
        this.shortcutButtons = document.querySelectorAll('.shortcut-btn');
        console.log(`找到 ${this.shortcutButtons.length} 个时间段快捷按钮`);
    }

    /**
     * 设置事件监听器
     */
    setupEventListeners() {
        // 为每个快捷按钮添加点击事件
        this.shortcutButtons.forEach(button => {
            button.addEventListener('click', (e) => {
                e.preventDefault();
                this.handleShortcutClick(button);
            });
        });

        // 监听日期输入框的手动更改
        const startDateInput = document.getElementById('startDate');
        const endDateInput = document.getElementById('endDate');
        
        if (startDateInput) {
            startDateInput.addEventListener('change', () => {
                this.clearActiveButton();
            });
        }
        
        if (endDateInput) {
            endDateInput.addEventListener('change', () => {
                this.clearActiveButton();
            });
        }
    }

    /**
     * 处理快捷按钮点击
     */
    handleShortcutClick(button) {
        try {
            const days = parseInt(button.getAttribute('data-days'));
            if (isNaN(days) || days <= 0) {
                console.error('无效的天数值:', button.getAttribute('data-days'));
                return;
            }

            // 计算日期范围
            const endDate = new Date();
            const startDate = new Date();
            startDate.setDate(endDate.getDate() - days);

            // 设置日期输入框
            this.setDateInputs(startDate, endDate);

            // 更新活跃按钮状态
            this.setActiveButton(button);

            // 显示反馈消息
            this.showMessage(`已选择${button.textContent}的数据范围`, 'success');

            console.log(`设置时间范围: ${this.client.formatDate(startDate)} 到 ${this.client.formatDate(endDate)}`);

        } catch (error) {
            console.error('处理快捷按钮点击失败:', error);
            this.showMessage('设置时间范围失败', 'error');
        }
    }

    /**
     * 设置日期输入框的值
     */
    setDateInputs(startDate, endDate) {
        const startDateInput = document.getElementById('startDate');
        const endDateInput = document.getElementById('endDate');

        if (startDateInput) {
            startDateInput.value = this.client.formatDate(startDate);
        }
        
        if (endDateInput) {
            endDateInput.value = this.client.formatDate(endDate);
        }
    }

    /**
     * 设置活跃按钮
     */
    setActiveButton(button) {
        // 清除所有按钮的活跃状态
        this.clearActiveButton();
        
        // 设置当前按钮为活跃状态
        button.classList.add('active');
        this.activeButton = button;
    }

    /**
     * 清除活跃按钮状态
     */
    clearActiveButton() {
        this.shortcutButtons.forEach(button => {
            button.classList.remove('active');
        });
        this.activeButton = null;
    }

    /**
     * 显示消息提示
     */
    showMessage(message, type = 'info') {
        // 创建消息元素
        const messageEl = document.createElement('div');
        messageEl.className = `message message-${type}`;
        messageEl.textContent = message;

        // 添加到页面
        document.body.appendChild(messageEl);

        // 3秒后自动移除
        setTimeout(() => {
            if (messageEl.parentNode) {
                messageEl.parentNode.removeChild(messageEl);
            }
        }, 3000);
    }

    /**
     * 获取当前选中的时间范围天数
     */
    getCurrentRangeDays() {
        if (!this.activeButton) {
            return null;
        }
        return parseInt(this.activeButton.getAttribute('data-days'));
    }

    /**
     * 根据当前日期范围设置对应的活跃按钮
     */
    updateActiveButtonByDateRange() {
        const startDateInput = document.getElementById('startDate');
        const endDateInput = document.getElementById('endDate');

        if (!startDateInput || !endDateInput || !startDateInput.value || !endDateInput.value) {
            this.clearActiveButton();
            return;
        }

        try {
            const startDate = new Date(startDateInput.value);
            const endDate = new Date(endDateInput.value);
            const daysDiff = Math.round((endDate - startDate) / (1000 * 60 * 60 * 24));

            // 查找匹配的快捷按钮
            const matchingButton = Array.from(this.shortcutButtons).find(button => {
                const buttonDays = parseInt(button.getAttribute('data-days'));
                return Math.abs(buttonDays - daysDiff) <= 1; // 允许1天的误差
            });

            if (matchingButton) {
                this.setActiveButton(matchingButton);
            } else {
                this.clearActiveButton();
            }

        } catch (error) {
            console.error('更新活跃按钮状态失败:', error);
            this.clearActiveButton();
        }
    }

    /**
     * 重新初始化（用于动态添加的按钮）
     */
    reinit() {
        this.setupShortcutButtons();
        this.setupEventListeners();
    }
}

// 导出日期快捷按钮模块类
window.DateShortcutsModule = DateShortcutsModule;
