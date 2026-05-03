import { useState, useCallback, useEffect, useRef } from 'react';
import { useToast } from '@/contexts/ToastContext';
import { clipboardService } from '@/services/api';
import { ClipboardItem, ClipboardType, SaveClipboardRequest } from '@/types/clipboard';
import { ClipboardFilterType } from '@/components/clipboard/TabBar';
import { detectClipboardType } from '@/utils/clipboardTypeDetector';
import { getClipboardTypeName } from '@/utils/clipboardTypeDetector';
import { settingsManager } from '@/utils/settings';
import { RichClipboardContent, writeClipboardRich } from '@/utils/richClipboard';

interface UseClipboardDataProps {
  pageSize?: number;
  isChannelVerified: boolean;
}

interface UseClipboardDataReturn {
  currentClipboard: ClipboardItem | undefined;
  clipboardItems: ClipboardItem[];
  isLoading: boolean;
  isLoadingMore: boolean;
  hasMore: boolean;
  fetchClipboardData: () => Promise<void>;
  fetchTabData: (tab: ClipboardFilterType) => Promise<void>;
  loadMoreData: (activeTab: ClipboardFilterType) => Promise<void>;
  handleSaveClipboardContent: (payload: RichClipboardContent) => Promise<boolean>;
  handleCopy: (item?: ClipboardItem) => void;
  handleEdit: (item?: ClipboardItem) => void;
  handleDelete: (item: ClipboardItem) => Promise<void>;
  handleToggleFavorite: (item: ClipboardItem) => Promise<void>;
  handleSave: (data: SaveClipboardRequest) => Promise<boolean>;
  handleRefresh: () => Promise<void>;
  handleSaveManualInput: (content: string, type?: ClipboardType, isManualInput?: boolean, contentHTML?: string, contentFormat?: 'plain' | 'html') => Promise<boolean>;
  setClipboardItems: React.Dispatch<React.SetStateAction<ClipboardItem[]>>;
  setCurrentClipboard: React.Dispatch<React.SetStateAction<ClipboardItem | undefined>>;
}

export const useClipboardData = ({
  pageSize = 12,
  isChannelVerified
}: UseClipboardDataProps): UseClipboardDataReturn => {
  const [currentClipboard, setCurrentClipboard] = useState<ClipboardItem | undefined>();
  const [clipboardItems, setClipboardItems] = useState<ClipboardItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [hasMore, setHasMore] = useState(false);
  const tabSeqRef = useRef(0);

  const { showToast } = useToast();

  const removeDuplicateContentFromList = useCallback((items: ClipboardItem[], savedItem: ClipboardItem) => {
    if (!settingsManager.getSetting('autoCleanDuplicates')) {
      return items;
    }

    const normalizedContent = savedItem.content.trim();
    const savedHTML = savedItem.content_html?.trim();
    return items.filter(item => {
      if (item.id === savedItem.id) return true;
      if (item.content.trim() === normalizedContent) return false;
      if (savedHTML && item.content_html?.trim() === savedHTML) return false;
      return true;
    });
  }, []);

  const dedupeItemsForDisplay = useCallback((items: ClipboardItem[]) => {
    if (!settingsManager.getSetting('autoCleanDuplicates')) {
      return items;
    }

    const seenText = new Set<string>();
    const seenHTML = new Set<string>();
    return items.filter(item => {
      const html = item.content_html?.trim();
      if (html) {
        if (seenHTML.has(html)) return false;
        seenHTML.add(html);
      }
      const text = item.content.trim();
      if (text) {
        if (seenText.has(text)) return false;
        seenText.add(text);
      }
      return true;
    });
  }, []);
  
  const fetchClipboardData = useCallback(async () => {
    if (!isChannelVerified) {
      showToast('请先验证通道', 'warning');
      return;
    }

    const seq = tabSeqRef.current;

    try {
      setIsLoading(true);
      const [latestRes, historyRes] = await Promise.all([
        clipboardService.getLatestClipboard(),
        clipboardService.getClipboardHistory(pageSize)
      ]);

      // tab 已切换过，只更新 currentClipboard，不覆盖当前 tab 的列表
      const tabChanged = seq !== tabSeqRef.current;

      if (latestRes.success && latestRes.data) {
        setCurrentClipboard(latestRes.data);
      }

      if (!tabChanged && historyRes.success && historyRes.data) {
        const items = historyRes.data.items || [];
        setClipboardItems(dedupeItemsForDisplay(items));
        setHasMore(historyRes.data.has_more || false);
      }
    } catch (error) {
      showToast('获取数据失败', 'error');
    } finally {
      if (seq === tabSeqRef.current) setIsLoading(false);
    }
  }, [showToast, pageSize, isChannelVerified, dedupeItemsForDisplay]);
  
  const loadMoreData = useCallback(async (activeTab: ClipboardFilterType) => {
    if (!isChannelVerified) {
      showToast('请先验证通道', 'warning');
      return;
    }

    if (!hasMore || isLoadingMore) return;

    try {
      setIsLoadingMore(true);

      // 取最后一条作为 cursor
      const lastItem = clipboardItems[clipboardItems.length - 1];
      const after = lastItem?.created_at;
      const afterId = lastItem?.id;

      let response;
      if (activeTab === 'favorite') {
        response = await clipboardService.getFavorites(1, pageSize);
      } else if (activeTab === 'all') {
        response = await clipboardService.getClipboardHistory(pageSize, after, afterId);
      } else {
        response = await clipboardService.getClipboardByType(activeTab as ClipboardType, pageSize, after, afterId);
      }

      if (response.success && response.data) {
        const existingIds = new Set(clipboardItems.map(item => item.id));

        if ('items' in response.data) {
          const items = response.data.items || [];
          const uniqueNewItems = dedupeItemsForDisplay(items).filter(item => !existingIds.has(item.id));
          setClipboardItems(prevItems => [...prevItems, ...uniqueNewItems]);
          setHasMore('has_more' in response.data ? (response.data as {has_more?: boolean}).has_more || false : items.length === pageSize);
        } else if (Array.isArray(response.data)) {
          const uniqueFilteredItems = dedupeItemsForDisplay(response.data).filter(item => !existingIds.has(item.id));
          setClipboardItems(prevItems => [...prevItems, ...uniqueFilteredItems]);
          setHasMore(false);
        }
      }
    } catch (error) {
      showToast('加载更多数据失败', 'error');
    } finally {
      setIsLoadingMore(false);
    }
  }, [hasMore, isLoadingMore, showToast, clipboardItems, pageSize, isChannelVerified, dedupeItemsForDisplay]);
  
  const handleSaveClipboardContent = useCallback(async (payload: RichClipboardContent) => {
    if (!isChannelVerified) {
      showToast('请先验证通道', 'warning');
      return false;
    }

    try {
      if (!payload.text || payload.text.trim() === '') {
        return false;
      }

      const detectedType = detectClipboardType(payload.text);

      const response = await clipboardService.saveClipboard({
        content: payload.text,
        type: detectedType,
        content_html: payload.html,
        content_format: payload.format,
      });
      
      if (!response || !response.success) {
        const errorMsg = response?.error || response?.message || '未知错误';
        showToast(`保存失败: ${errorMsg}`, 'error');
        return false;
      }
      
      if (response.data) {
        const rawData = response.data as any;
        const clipboardItem: ClipboardItem = {
          id: rawData.id || 'temp-' + Date.now(),
          content: rawData.content || payload.text,
          type: rawData.type || ClipboardType.TEXT,
          title: rawData.title || '',
          isFavorite: rawData.favorite || rawData.isFavorite || false,
          created_at: rawData.created_at || new Date().toISOString(),
          updated_at: rawData.updated_at || new Date().toISOString(),
          content_html: rawData.content_html || payload.html,
          content_format: rawData.content_format || payload.format,
        };
        
        setCurrentClipboard(clipboardItem);
        setClipboardItems(prev => [clipboardItem, ...removeDuplicateContentFromList(prev, clipboardItem)]);
      } else {
        const tempItem: ClipboardItem = {
          id: 'temp-' + Date.now(),
          content: payload.text,
          type: ClipboardType.TEXT,
          title: '',
          isFavorite: false,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          content_html: payload.html,
          content_format: payload.format,
        };
        setCurrentClipboard(tempItem);
        setClipboardItems(prev => [tempItem, ...removeDuplicateContentFromList(prev, tempItem)]);
      }
      
      showToast('新内容已同步', 'success');
      return true;
    } catch (error) {
      showToast('保存失败，请重试', 'error');
      return false;
    }
  }, [showToast, isChannelVerified, removeDuplicateContentFromList]);
  
  const handleCopy = useCallback((item?: ClipboardItem) => {
    if (!item) return;

    writeClipboardRich(item)
      .then(() => showToast('已复制到剪贴板', 'success'))
      .catch(() => {
        showToast('复制失败', 'error');
      });
  }, [showToast]);
  
  const handleEdit = useCallback((item?: ClipboardItem) => {
    return item;
  }, []);
  
  const handleDelete = useCallback(async (item: ClipboardItem) => {
    if (!isChannelVerified) {
      showToast('请先验证通道', 'warning');
      return;
    }
    
    try {
      const response = await clipboardService.deleteClipboard(item.id);
      if (response.success) {
        setClipboardItems(prevItems => prevItems.filter(i => i.id !== item.id));
        
        if (currentClipboard && currentClipboard.id === item.id) {
          const latestRes = await clipboardService.getLatestClipboard();
          if (latestRes.success && latestRes.data) {
            setCurrentClipboard(latestRes.data);
          } else {
            setCurrentClipboard(undefined);
          }
        }
        
        showToast('删除成功', 'success');
      } else {
        showToast(response.message || '删除失败', 'error');
      }
    } catch (error) {
      showToast('删除失败', 'error');
    }
  }, [currentClipboard, showToast, isChannelVerified]);
  
  const handleToggleFavorite = useCallback(async (item: ClipboardItem) => {
    if (!isChannelVerified) {
      showToast('请先验证通道', 'warning');
      return;
    }
    
    try {
      const response = await clipboardService.toggleFavorite(item.id, !item.isFavorite);
      if (response.success && response.data) {
        setClipboardItems(prevItems => 
          prevItems.map(i => i.id === item.id ? { ...i, isFavorite: !i.isFavorite } : i)
        );
        
        if (currentClipboard && currentClipboard.id === item.id) {
          setCurrentClipboard({ ...currentClipboard, isFavorite: !currentClipboard.isFavorite });
        }
        
        showToast(item.isFavorite ? '已取消收藏' : '已添加到收藏', 'success');
      } else {
        showToast(response.message || '操作失败', 'error');
      }
    } catch (error) {
      showToast('操作失败', 'error');
    }
  }, [currentClipboard, showToast, isChannelVerified]);
  
  const handleSave = useCallback(async (data: SaveClipboardRequest): Promise<boolean> => {
    if (!isChannelVerified) {
      showToast('请先验证通道', 'warning');
      return false;
    }

    if (!data.id) {
      showToast('缺少项目ID', 'error');
      return false;
    }

    // 提取 isFavorite 变更，收藏只走专用接口
    const newIsFavorite = data.isFavorite;
    const currentItem = clipboardItems.find(i => i.id === data.id);
    const currentIsFavorite = currentItem?.isFavorite;
    const favoriteChanged = newIsFavorite !== undefined && newIsFavorite !== currentIsFavorite;

    // 构造不含 isFavorite 的请求体
    const { isFavorite: _ignored, ...updateData } = data;

    try {
      const response = await clipboardService.updateClipboard(data.id, updateData);

      if (response.success && response.data) {
        // 如果收藏状态变了，调用专用接口
        if (favoriteChanged) {
          const favResponse = await clipboardService.toggleFavorite(data.id, newIsFavorite!);
          if (!favResponse.success) {
            showToast('内容已保存，收藏状态保存失败', 'error');
            return false;
          }
        }

        setClipboardItems(prevItems =>
          prevItems.map(item => item.id === data.id ? response.data! : item)
        );

        if (currentClipboard && currentClipboard.id === data.id) {
          setCurrentClipboard(response.data);
        }

        showToast('保存成功', 'success');
        return true;
      } else {
        showToast(response.message || '保存失败', 'error');
        return false;
      }
    } catch (error) {
      showToast('保存失败', 'error');
      return false;
    }
  }, [clipboardItems, currentClipboard, showToast, isChannelVerified]);
  
  const handleRefresh = useCallback(async () => {
    await fetchClipboardData();
  }, [fetchClipboardData]);
  
  const handleSaveManualInput = useCallback(async (
    content: string,
    type?: ClipboardType,
    _isManualInput: boolean = true,
    contentHTML?: string,
    contentFormat?: 'plain' | 'html'
  ): Promise<boolean> => {
    if (!isChannelVerified) {
      showToast('请先验证通道', 'warning');
      return false;
    }

      const contentType = type || detectClipboardType(content);

    if (contentType !== ClipboardType.TEXT && !type) {
      showToast(`检测到${getClipboardTypeName(contentType)}内容`, 'info');
      }

    try {
      const response = await clipboardService.saveClipboard({
        content,
        type: contentType,
        content_html: contentHTML,
        content_format: contentFormat,
      });
      
      if (response.success && response.data) {
        const clipboardItem: ClipboardItem = response.data;
        
        setCurrentClipboard(clipboardItem);
        setClipboardItems(prev => [clipboardItem, ...removeDuplicateContentFromList(prev, clipboardItem)]);
        
        return true;
      } else {
        showToast(response.message || '保存失败', 'error');
        return false;
      }
    } catch (error) {
      showToast('保存失败，请重试', 'error');
      return false;
    }
  }, [showToast, isChannelVerified, removeDuplicateContentFromList]);
  
  const fetchTabData = useCallback(async (tab: ClipboardFilterType) => {
    if (!isChannelVerified) {
      showToast('请先验证通道', 'warning');
      return;
    }

    const seq = ++tabSeqRef.current;

    try {
      setIsLoading(true);
      let response;

      if (tab === 'favorite') {
        response = await clipboardService.getFavorites(1, pageSize);
      } else if (tab === 'all') {
        response = await clipboardService.getClipboardHistory(pageSize);
      } else {
        response = await clipboardService.getClipboardByType(tab as ClipboardType, pageSize);
      }

      // 快速切换时丢弃过期响应，只允许最后一次请求落状态
      if (seq !== tabSeqRef.current) return;

      if (response.success && response.data) {
        let items: ClipboardItem[] = [];

        if (Array.isArray(response.data)) {
          items = response.data;
        } else if ('items' in response.data) {
          items = response.data.items || [];
        }

        setClipboardItems(dedupeItemsForDisplay(items));
        setHasMore('has_more' in response.data ? (response.data as {has_more?: boolean}).has_more || false : false);
      }
    } catch (error) {
      if (seq !== tabSeqRef.current) return;
      showToast('获取数据失败', 'error');
    } finally {
      if (seq === tabSeqRef.current) setIsLoading(false);
    }
  }, [showToast, pageSize, isChannelVerified, dedupeItemsForDisplay]);
  
  useEffect(() => {
    if (isChannelVerified) {
      fetchClipboardData();
    }
  }, [fetchClipboardData, isChannelVerified]);
  
  return {
    currentClipboard,
    clipboardItems,
    isLoading,
    isLoadingMore,
    hasMore,
    fetchClipboardData,
    fetchTabData,
    loadMoreData,
    handleSaveClipboardContent,
    handleCopy,
    handleEdit,
    handleDelete,
    handleToggleFavorite,
    handleSave,
    handleRefresh,
    handleSaveManualInput,
    setClipboardItems,
    setCurrentClipboard
  };
}; 
