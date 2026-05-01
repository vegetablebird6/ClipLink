import { useState, useRef, useEffect, useCallback } from 'react';
import { ClipboardItem, ClipboardType, SaveClipboardRequest } from '@/types/clipboard';
import { settingsManager } from '@/utils/settings';
import { RichClipboardContent, readClipboardRich } from '@/utils/richClipboard';

interface ClipboardFilterOptions {
  hasClipboardPermission: boolean;
  isIOSDevice: boolean;
  isChannelVerified: boolean;
  onSaveContent: (payload: RichClipboardContent) => Promise<boolean>;
  onConfirmSave?: (payload: RichClipboardContent) => void;
  debug?: boolean;
}

interface UseClipboardFilterReturn {
  syncClipboard: (force?: boolean) => Promise<void>;
  handleDeletedContent: (content: string) => void;
  handleFilteredContent: (content: string | RichClipboardContent) => Promise<boolean>;
  /** 纯过滤判断，不触发保存。用于新建路径自行控制保存时机。 */
  shouldAllowContent: (content: string) => boolean;
  trackProcessedContent: (content: string) => void;
  /** 确认保存完成后，清除待确认标记 */
  clearPendingConfirm: (content: string) => void;
  processedContents: Set<string>;
  deletedContents: Set<string>;
  resetFilter: () => void;
}

// 工具函数：去除内容首尾空格
function trimContent(content: string): string {
  return content.trim();
}

// 工具函数：内容是否为空
function isEmptyContent(content: string): boolean {
  return !content || content.trim() === '';
}

export function useClipboardFilter({
  hasClipboardPermission,
  isIOSDevice,
  isChannelVerified,
  onSaveContent,
  onConfirmSave,
  debug = false
}: ClipboardFilterOptions): UseClipboardFilterReturn {
  // 状态管理
  const [processedContents, setProcessedContents] = useState<Set<string>>(new Set());
  const [deletedContents, setDeletedContents] = useState<Set<string>>(new Set());
  const lastSyncTimeRef = useRef<number>(0);
  const hasVisibilityChangedRef = useRef<boolean>(false);
  // 确认弹窗待处理内容（防重复弹窗，不阻塞实际保存）
  const pendingConfirmRef = useRef<Set<string>>(new Set());

  // 追踪已处理内容
  const trackProcessedContent = useCallback((content: string) => {
    if (isEmptyContent(content)) return;
    const trimmedContent = trimContent(content);
    setProcessedContents(prev => {
      const newSet = new Set(prev);
      newSet.add(trimmedContent);
      return newSet;
    });
  }, []);

  // 处理已删除内容
  const handleDeletedContent = useCallback((content: string) => {
    if (isEmptyContent(content)) return;
    const trimmedContent = trimContent(content);
    setDeletedContents(prev => {
      const newSet = new Set(prev);
      newSet.add(trimmedContent);
      return newSet;
    });
    if (debug) {
      console.log('[ClipboardFilter] 内容已添加到屏蔽列表:', trimmedContent);
    }
  }, [debug]);

  // 检查是否应该同步内容
  const shouldSyncContent = useCallback((content: string): boolean => {
    if (isEmptyContent(content)) return false;
    const trimmedContent = trimContent(content);
    if (deletedContents.has(trimmedContent)) {
      if (debug) {
        console.log('[ClipboardFilter] 内容在屏蔽列表中，跳过同步:', trimmedContent);
      }
      return false;
    }
    if (processedContents.has(trimmedContent)) {
      if (debug) {
        console.log('[ClipboardFilter] 内容已存在，跳过同步:', trimmedContent);
      }
      return false;
    }
    if (pendingConfirmRef.current.has(trimmedContent)) {
      if (debug) {
        console.log('[ClipboardFilter] 内容正在等待确认，跳过同步:', trimmedContent);
      }
      return false;
    }
    return true;
  }, [deletedContents, processedContents, debug]);

  // 处理过滤后的内容
  const handleFilteredContent = useCallback(async (content: string | RichClipboardContent): Promise<boolean> => {
    const payload: RichClipboardContent = typeof content === 'string'
      ? { text: content, format: 'plain' as const }
      : content;
    const text = payload.text;

    if (isEmptyContent(text)) return false;
    if (!shouldSyncContent(text)) return false;

    // 添加与上次保存的最小时间间隔检查
    const now = Date.now();
    const minTimeBetweenSaves = 500; // 最小间隔为500毫秒
    if (now - lastSyncTimeRef.current < minTimeBetweenSaves) {
      if (debug) {
        console.log('[ClipboardFilter] 距离上次同步时间太短，跳过:', trimContent(text));
      }
      return false;
    }

    // 检查是否需要确认
    const confirmBeforeSave = settingsManager.getSetting('confirmBeforeSave');
    if (confirmBeforeSave && onConfirmSave) {
      // 标记为待确认，避免重复弹窗；不加入 processedContents，否则 shouldAllowContent 会在实际保存时拒绝
      const trimmed = trimContent(text);
      pendingConfirmRef.current.add(trimmed);
      // 如果需要确认，传递完整 payload（保留 HTML）
      onConfirmSave(payload);
      return true;
    }

    const result = await onSaveContent(payload);
    if (result) {
      trackProcessedContent(text);
      lastSyncTimeRef.current = now;
      
      // 触发剪贴板更新事件，通知应用有新内容
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new Event('clipboard-updated'));
      }
      
      if (debug) {
        console.log('[ClipboardFilter] 新内容已同步:', trimContent(text));
      }
    }
    return result;
  }, [shouldSyncContent, onSaveContent, onConfirmSave, trackProcessedContent, debug]);

  // 同步剪贴板内容
  const syncClipboard = useCallback(async (force = false) => {
    if (!isChannelVerified) return;
    if (!hasClipboardPermission || isIOSDevice) return;
    
    // 检查是否启用自动读取剪切板
    const autoReadClipboard = settingsManager.getSetting('autoReadClipboard');
    if (!autoReadClipboard && !force) {
      if (debug) {
        console.log('[ClipboardFilter] 自动读取剪切板已禁用，跳过同步');
      }
      return;
    }
    
    if (!force) {
      const now = Date.now();
      const timeSinceLastSync = now - lastSyncTimeRef.current;
      if (hasVisibilityChangedRef.current && timeSinceLastSync < 3000) {
        return;
      }
      hasVisibilityChangedRef.current = false;
    }
    try {
      const payload = await readClipboardRich();
      // 传递完整富文本 payload，保留 HTML 内容
      await handleFilteredContent(payload);
    } catch (error) {
      if (debug) {
        console.warn('[ClipboardFilter] 同步错误:', error);
      }
    }
  }, [hasClipboardPermission, isIOSDevice, isChannelVerified, handleFilteredContent, debug]);

  // 标记可见性变化
  const setVisibilityChanged = useCallback(() => {
    hasVisibilityChangedRef.current = true;
  }, []);

  // 监听可见性和焦点变化
  useEffect(() => {
    if (!isChannelVerified) return;
    
    const handleVisibilityChange = () => {
      if (document.visibilityState === 'visible') {
        setVisibilityChanged();
        syncClipboard(false);
      }
    };
    
    const handleWindowFocus = () => {
      setVisibilityChanged();
      syncClipboard(false);
    };
    
    // 初始同步（只有在启用自动读取时）
    const autoReadClipboard = settingsManager.getSetting('autoReadClipboard');
    if (hasClipboardPermission && !isIOSDevice && autoReadClipboard) {
      syncClipboard(true);
    }
    
    document.addEventListener('visibilitychange', handleVisibilityChange);
    window.addEventListener('focus', handleWindowFocus);
    
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      window.removeEventListener('focus', handleWindowFocus);
    };
  }, [hasClipboardPermission, isIOSDevice, isChannelVerified, syncClipboard, setVisibilityChanged]);

  // 重置过滤器
  const resetFilter = useCallback(() => {
    setProcessedContents(new Set());
    setDeletedContents(new Set());
    lastSyncTimeRef.current = 0;
    hasVisibilityChangedRef.current = false;
    if (debug) {
      console.log('[ClipboardFilter] 过滤器已重置');
    }
  }, [debug]);

  // 纯过滤判断：内容是否应被保存（不触发保存，不更新时间戳）
  // 注意：不检查 pendingConfirmRef，pending 只用于 shouldSyncContent 防重复弹窗
  const shouldAllowContent = useCallback((content: string): boolean => {
    if (isEmptyContent(content)) return false;
    const trimmed = trimContent(content);
    if (deletedContents.has(trimmed)) return false;
    if (processedContents.has(trimmed)) return false;
    const now = Date.now();
    if (now - lastSyncTimeRef.current < 500) return false;
    return true;
  }, [deletedContents, processedContents]);

  // 确认保存完成后，清除待确认标记
  const clearPendingConfirm = useCallback((content: string) => {
    const trimmed = trimContent(content);
    pendingConfirmRef.current.delete(trimmed);
  }, []);

  // 返回 API
  return {
    syncClipboard,
    handleDeletedContent,
    handleFilteredContent,
    shouldAllowContent,
    trackProcessedContent,
    clearPendingConfirm,
    processedContents,
    deletedContents,
    resetFilter
  };
}