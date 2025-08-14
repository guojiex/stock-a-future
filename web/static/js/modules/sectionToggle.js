/**
 * Section 展开/折叠功能模块
 * 为页面中的所有section添加展开和折叠功能
 */

class SectionToggleModule {
    constructor() {
        this.collapsedSections = new Set();
        this.storageKey = 'stockafuture_collapsed_sections';
        
        this.init();
    }

    /**
     * 初始化模块
     */
    init() {
        // 加载保存的折叠状态
        this.loadCollapsedState();
        
        // 等待DOM加载完成后初始化
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => {
                this.initializeSections();
            });
        } else {
            this.initializeSections();
        }
    }

    /**
     * 初始化所有section的折叠功能
     */
    initializeSections() {
        // 为现有的section添加折叠功能
        this.addToggleToExistingSections();
        
        // 监听动态添加的section（如收藏列表）
        this.observeNewSections();
    }

    /**
     * 为现有的section添加折叠功能
     */
    addToggleToExistingSections() {
        const sections = document.querySelectorAll('section:not(.collapsible)');
        
        sections.forEach(section => {
            this.makeSectionCollapsible(section);
        });
    }

    /**
     * 监听动态添加的section
     */
    observeNewSections() {
        const observer = new MutationObserver((mutations) => {
            mutations.forEach((mutation) => {
                mutation.addedNodes.forEach((node) => {
                    if (node.nodeType === Node.ELEMENT_NODE) {
                        // 检查新添加的节点是否是section
                        if (node.tagName === 'SECTION' && !node.classList.contains('collapsible')) {
                            this.makeSectionCollapsible(node);
                        }
                        
                        // 检查新添加节点内部的section
                        const sections = node.querySelectorAll?.('section:not(.collapsible)');
                        sections?.forEach(section => {
                            this.makeSectionCollapsible(section);
                        });
                    }
                });
            });
        });

        observer.observe(document.body, {
            childList: true,
            subtree: true
        });
    }

    /**
     * 让section变为可折叠的
     */
    makeSectionCollapsible(section) {
        if (section.classList.contains('collapsible')) {
            return; // 已经处理过了
        }



        // 添加可折叠标记
        section.classList.add('collapsible');

        // 获取section ID，如果没有则生成一个
        let sectionId = section.id;
        if (!sectionId) {
            sectionId = this.generateSectionId(section);
            section.id = sectionId;
        }

        // 查找标题元素
        const titleElement = this.findTitleElement(section);
        if (!titleElement) {
            console.warn('未找到section标题元素:', section);
            return;
        }

        // 创建section内容容器
        this.createContentContainer(section, titleElement);

        // 创建标题容器和切换按钮
        this.createToggleHeader(section, titleElement, sectionId);

        // 应用保存的折叠状态
        if (this.collapsedSections.has(sectionId)) {
            this.collapseSection(section, false);
        }
    }

    /**
     * 生成section ID
     */
    generateSectionId(section) {
        const className = section.className.split(' ')[0] || 'section';
        const timestamp = Date.now();
        return `${className}-${timestamp}`;
    }

    /**
     * 查找标题元素
     */
    findTitleElement(section) {
        // 优先查找直接子元素中的标题
        let titleElement = section.querySelector(':scope > .card > h2, :scope > .card > h3');
        
        if (!titleElement) {
            titleElement = section.querySelector(':scope > h2, :scope > h3');
        }

        return titleElement;
    }

    /**
     * 创建内容容器
     */
    createContentContainer(section, titleElement) {
        const card = section.querySelector('.card');
        if (!card) return;

        // 检查是否已经有内容容器
        if (card.querySelector('.section-content')) {
            return;
        }

        // 获取除标题外的所有内容
        const contentElements = Array.from(card.children).filter(child => 
            child !== titleElement && !child.classList.contains('section-header')
        );

        if (contentElements.length === 0) return;

        // 创建内容容器
        const contentContainer = document.createElement('div');
        contentContainer.className = 'section-content expanded';

        // 将内容移动到容器中
        contentElements.forEach(element => {
            contentContainer.appendChild(element);
        });

        // 将内容容器添加到card中
        card.appendChild(contentContainer);
    }

    /**
     * 创建切换标题和按钮
     */
    createToggleHeader(section, titleElement, sectionId) {
        const card = section.querySelector('.card');
        if (!card) return;

        // 创建标题容器
        const headerContainer = document.createElement('div');
        headerContainer.className = 'section-header';

        // 创建切换按钮
        const toggleBtn = document.createElement('button');
        toggleBtn.className = 'section-toggle-btn';
        toggleBtn.setAttribute('aria-label', '展开/折叠section');
        toggleBtn.innerHTML = '<span class="toggle-icon">▼</span>';

        // 创建状态指示器
        const statusSpan = document.createElement('span');
        statusSpan.className = 'section-status';
        statusSpan.textContent = '(已折叠)';

        // 将标题移动到标题容器中
        const titleClone = titleElement.cloneNode(true);
        titleClone.appendChild(statusSpan);
        headerContainer.appendChild(titleClone);
        headerContainer.appendChild(toggleBtn);

        // 替换原标题
        titleElement.replaceWith(headerContainer);

        // 添加点击事件到标题容器
        headerContainer.addEventListener('click', (e) => {
            e.preventDefault();
            this.toggleSection(section);
        });

        // 按钮也应该触发切换，而不是阻止事件
        toggleBtn.addEventListener('click', (e) => {
            e.preventDefault();
            e.stopPropagation();
            this.toggleSection(section);
        });
    }

    /**
     * 切换section状态
     */
    toggleSection(section) {
        const contentContainer = section.querySelector('.section-content');
        const isCollapsed = contentContainer?.classList.contains('collapsed');

        if (isCollapsed) {
            this.expandSection(section);
        } else {
            this.collapseSection(section);
        }
    }

    /**
     * 展开section
     */
    expandSection(section, saveState = true) {
        const contentContainer = section.querySelector('.section-content');
        const headerContainer = section.querySelector('.section-header');
        const toggleBtn = section.querySelector('.section-toggle-btn');

        if (!contentContainer) return;

        // 移除折叠状态
        contentContainer.classList.remove('collapsed');
        contentContainer.classList.add('expanded');
        headerContainer?.classList.remove('collapsed');
        toggleBtn?.classList.remove('collapsed');

        // 设置最大高度为自动，允许内容完全展开
        requestAnimationFrame(() => {
            const height = contentContainer.scrollHeight;
            contentContainer.style.maxHeight = height + 'px';
            
            // 动画完成后移除内联样式
            setTimeout(() => {
                if (!contentContainer.classList.contains('collapsed')) {
                    contentContainer.style.maxHeight = '';
                }
            }, 300);
        });

        // 保存状态
        if (saveState && section.id) {
            this.collapsedSections.delete(section.id);
            this.saveCollapsedState();
        }
    }

    /**
     * 折叠section
     */
    collapseSection(section, saveState = true) {
        const contentContainer = section.querySelector('.section-content');
        const headerContainer = section.querySelector('.section-header');
        const toggleBtn = section.querySelector('.section-toggle-btn');

        if (!contentContainer) return;

        // 设置当前高度
        const height = contentContainer.scrollHeight;
        contentContainer.style.maxHeight = height + 'px';

        // 强制重绘
        requestAnimationFrame(() => {
            // 添加折叠状态
            contentContainer.classList.add('collapsed');
            contentContainer.classList.remove('expanded');
            headerContainer?.classList.add('collapsed');
            toggleBtn?.classList.add('collapsed');
        });

        // 保存状态
        if (saveState && section.id) {
            this.collapsedSections.add(section.id);
            this.saveCollapsedState();
        }
    }

    /**
     * 展开所有section
     */
    expandAllSections() {
        const sections = document.querySelectorAll('section.collapsible');
        sections.forEach(section => {
            this.expandSection(section);
        });
    }

    /**
     * 折叠所有section
     */
    collapseAllSections() {
        const sections = document.querySelectorAll('section.collapsible');
        sections.forEach(section => {
            this.collapseSection(section);
        });
    }

    /**
     * 保存折叠状态到localStorage
     */
    saveCollapsedState() {
        try {
            const stateArray = Array.from(this.collapsedSections);
            localStorage.setItem(this.storageKey, JSON.stringify(stateArray));
        } catch (error) {
            console.warn('无法保存section折叠状态:', error);
        }
    }

    /**
     * 从localStorage加载折叠状态
     */
    loadCollapsedState() {
        try {
            const saved = localStorage.getItem(this.storageKey);
            if (saved) {
                const stateArray = JSON.parse(saved);
                this.collapsedSections = new Set(stateArray);
            }
        } catch (error) {
            console.warn('无法加载section折叠状态:', error);
            this.collapsedSections = new Set();
        }
    }

    /**
     * 重置所有折叠状态
     */
    resetCollapsedState() {
        this.collapsedSections.clear();
        this.saveCollapsedState();
        this.expandAllSections();
    }

    /**
     * 获取section的折叠状态
     */
    getSectionState(sectionId) {
        return this.collapsedSections.has(sectionId);
    }

    /**
     * 设置section的折叠状态
     */
    setSectionState(sectionId, collapsed) {
        const section = document.getElementById(sectionId);
        if (!section) return;

        if (collapsed) {
            this.collapseSection(section);
        } else {
            this.expandSection(section);
        }
    }
}

// 导出类，但不自动初始化
window.SectionToggleModule = SectionToggleModule;

// 延迟初始化，确保DOM和其他模块都已加载
window.addEventListener('load', () => {
    // 等待一小段时间，确保所有动态内容都已生成
    setTimeout(() => {
        window.sectionToggleModule = new SectionToggleModule();
        console.log('SectionToggleModule initialized after page load');
    }, 500);
});
