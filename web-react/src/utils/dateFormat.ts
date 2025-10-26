/**
 * 日期格式化工具函数
 */

/**
 * 格式化日期字符串为 YYYY-MM-DD 格式
 * 支持多种输入格式：
 * - YYYYMMDD (例如: "20251024")
 * - ISO 8601 (例如: "2025-10-24T00:00:00.000Z")
 * - YYYY-MM-DD (例如: "2025-10-24")
 * - Date 对象
 * 
 * @param dateStr 日期字符串或 Date 对象
 * @returns 格式化后的日期字符串 YYYY-MM-DD，如果无法解析则返回原字符串
 */
export function formatDate(dateStr: string | Date | undefined | null): string {
  if (!dateStr) return '';

  try {
    // 如果是 Date 对象
    if (dateStr instanceof Date) {
      return dateStr.toISOString().split('T')[0];
    }

    const str = String(dateStr);

    // 如果已经是 YYYY-MM-DD 格式，直接返回
    if (/^\d{4}-\d{2}-\d{2}$/.test(str)) {
      return str;
    }

    // 如果是 YYYYMMDD 格式 (8位数字)
    if (/^\d{8}$/.test(str)) {
      return str.replace(/(\d{4})(\d{2})(\d{2})/, '$1-$2-$3');
    }

    // 如果是 ISO 8601 格式或其他可解析的日期字符串
    if (str.includes('T') || str.includes('-')) {
      const date = new Date(str);
      if (!isNaN(date.getTime())) {
        return date.toISOString().split('T')[0];
      }
    }

    // 无法识别的格式，返回原字符串
    return str;
  } catch (error) {
    console.warn('日期格式化失败:', dateStr, error);
    return String(dateStr);
  }
}

/**
 * 格式化日期时间字符串为 YYYY-MM-DD HH:mm:ss 格式
 * 
 * @param dateStr 日期时间字符串或 Date 对象
 * @returns 格式化后的日期时间字符串
 */
export function formatDateTime(dateStr: string | Date | undefined | null): string {
  if (!dateStr) return '';

  try {
    const date = dateStr instanceof Date ? dateStr : new Date(dateStr);
    
    if (isNaN(date.getTime())) {
      return String(dateStr);
    }

    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
  } catch (error) {
    console.warn('日期时间格式化失败:', dateStr, error);
    return String(dateStr);
  }
}

/**
 * 格式化为友好的相对时间
 * 例如：刚刚、5分钟前、1小时前、昨天、2天前等
 * 
 * @param dateStr 日期时间字符串或 Date 对象
 * @returns 友好的相对时间字符串
 */
export function formatRelativeTime(dateStr: string | Date | undefined | null): string {
  if (!dateStr) return '';

  try {
    const date = dateStr instanceof Date ? dateStr : new Date(dateStr);
    const now = new Date();
    const diff = now.getTime() - date.getTime();

    // 小于1分钟
    if (diff < 60 * 1000) {
      return '刚刚';
    }

    // 小于1小时
    if (diff < 60 * 60 * 1000) {
      const minutes = Math.floor(diff / (60 * 1000));
      return `${minutes}分钟前`;
    }

    // 小于1天
    if (diff < 24 * 60 * 60 * 1000) {
      const hours = Math.floor(diff / (60 * 60 * 1000));
      return `${hours}小时前`;
    }

    // 小于7天
    if (diff < 7 * 24 * 60 * 60 * 1000) {
      const days = Math.floor(diff / (24 * 60 * 60 * 1000));
      if (days === 1) return '昨天';
      return `${days}天前`;
    }

    // 超过7天，返回具体日期
    return formatDate(date);
  } catch (error) {
    console.warn('相对时间格式化失败:', dateStr, error);
    return String(dateStr);
  }
}

