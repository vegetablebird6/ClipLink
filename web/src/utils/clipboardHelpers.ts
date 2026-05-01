/**
 * 格式化剪贴板内容，用于显示
 * @param content 原始内容
 * @param maxLength 最大长度
 * @returns 格式化后的内容
 */
export const formatClipboardContent = (content: string, maxLength: number = 50): string => {
  if (!content) return '';

  // 如果内容太长，截断并添加省略号
  if (content.length > maxLength) {
    return content.substring(0, maxLength) + '...';
  }

  return content;
};
