'use client';

import React, { useState, useEffect } from 'react';
import ClipboardGrid from '@/components/clipboard/ClipboardGrid';
import EditModal from '@/components/clipboard/EditModal';
import PreviewModal from '@/components/clipboard/PreviewModal';
import { ClipboardItem, SaveClipboardRequest } from '@/types/clipboard';
import { clipboardService } from '@/services/api';
import { useToast } from '@/contexts/ToastContext';
import { writeClipboardRich } from '@/utils/richClipboard';
import { ClipboardGridSkeleton } from '@/components/ui/LoadingStates';

export default function FavoritesPage() {
  const [favoriteItems, setFavoriteItems] = useState<ClipboardItem[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<ClipboardItem | undefined>();
  const [isPreviewOpen, setIsPreviewOpen] = useState(false);
  const [previewItem, setPreviewItem] = useState<ClipboardItem | undefined>();
  const [isLoading, setIsLoading] = useState(true);

  const { showToast } = useToast();

  // 获取收藏的剪贴板项目
  const fetchFavorites = async () => {
    setIsLoading(true);

    try {
      const response = await clipboardService.getFavorites();
      if (response.success && response.data) {
        setFavoriteItems(response.data);
      } else {
        showToast(response.message || '获取收藏失败', 'error');
      }
    } catch (error) {
      showToast('获取收藏失败', 'error');
    } finally {
      setIsLoading(false);
    }
  };

  // 初始化加载
  useEffect(() => {
    fetchFavorites();
  }, []);

  // 处理复制操作
  const handleCopy = (item: ClipboardItem) => {
    writeClipboardRich(item)
      .then(() => showToast('已复制到剪贴板', 'success'))
      .catch(() => showToast('复制失败', 'error'));
  };

  // 打开编辑模态窗
  const handleEdit = (item?: ClipboardItem) => {
    setEditingItem(item);
    setIsModalOpen(true);
  };

  // 处理删除操作
  const handleDelete = async (item: ClipboardItem) => {
    if (!window.confirm('确定要删除这个收藏项目吗？')) {
      return;
    }
    
    try {
      const response = await clipboardService.deleteClipboard(item.id);
      if (response.success) {
        // 更新列表
        setFavoriteItems(prevItems => prevItems.filter(i => i.id !== item.id));
        showToast('删除成功', 'success');
      } else {
        showToast(response.message || '删除失败', 'error');
      }
    } catch (error) {
      showToast('删除失败', 'error');
    }
  };

  // 切换收藏状态
  const handleToggleFavorite = async (item: ClipboardItem) => {
    try {
      const response = await clipboardService.toggleFavorite(item.id, !item.isFavorite);
      if (response.success) {
        // 如果取消收藏，则从收藏列表中移除
        if (item.isFavorite) {
          setFavoriteItems(prevItems => prevItems.filter(i => i.id !== item.id));
          showToast('已取消收藏', 'success');
        } else {
          // 这里理论上不应该出现，因为收藏页面只显示已收藏的项目
          await fetchFavorites();
          showToast('已添加到收藏', 'success');
        }
      } else {
        showToast(response.message || '切换收藏状态失败', 'error');
      }
    } catch (error) {
      showToast('切换收藏状态失败', 'error');
    }
  };

  // 保存编辑
  const handleSave = async (data: SaveClipboardRequest): Promise<boolean> => {
    if (!editingItem) return false;

    const newIsFavorite = data.isFavorite;
    const currentIsFavorite = editingItem.isFavorite;
    const favoriteChanged = newIsFavorite !== undefined && newIsFavorite !== currentIsFavorite;

    const { isFavorite: _ignored, ...updateData } = data;

    try {
      const response = await clipboardService.updateClipboard(editingItem.id, updateData);
      if (response.success && response.data) {
        if (favoriteChanged) {
          const favResponse = await clipboardService.toggleFavorite(editingItem.id, newIsFavorite!);
          if (!favResponse.success) {
            showToast('内容已保存，收藏状态保存失败', 'error');
            return false;
          }
          if (newIsFavorite) {
            setFavoriteItems(prevItems =>
              prevItems.map(i => i.id === editingItem.id ? { ...favResponse.data! } : i)
            );
          } else {
            setFavoriteItems(prevItems => prevItems.filter(i => i.id !== editingItem.id));
          }
          showToast(newIsFavorite ? '已保存并添加到收藏' : '已保存并取消收藏', 'success');
        } else {
          setFavoriteItems(prevItems =>
            prevItems.map(i => i.id === editingItem.id ? response.data! : i)
          );
          showToast('保存成功', 'success');
        }
        return true;
      } else {
        showToast(response.message || '保存失败', 'error');
        return false;
      }
    } catch (error) {
      showToast('保存失败', 'error');
      return false;
    }
  };

  // 处理预览按钮点击
  const handlePreview = (item: ClipboardItem) => {
    setPreviewItem(item);
    setIsPreviewOpen(true);
  };

  return (
    <>
      <div className="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700 p-4">
        <h1 className="text-lg font-medium text-gray-900 dark:text-white">收藏夹</h1>
        <p className="text-sm text-gray-500 dark:text-gray-400">在这里查看您收藏的所有剪贴板项目</p>
      </div>
      
      <div className="flex-1 overflow-hidden bg-gray-50 dark:bg-gray-900">
        <div className="h-full overflow-y-auto custom-scrollbar p-4">
          {isLoading ? (
            <ClipboardGridSkeleton />
          ) : (
            <ClipboardGrid 
              items={favoriteItems}
              onCopy={handleCopy}
              onEdit={handleEdit}
              onDelete={handleDelete}
              onToggleFavorite={handleToggleFavorite}
              onPreview={handlePreview}
            />
          )}
        </div>
      </div>
      
      <EditModal 
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false);
          setEditingItem(undefined);
        }}
        onSave={handleSave}
        initialData={editingItem}
      />

      {/* 预览模态框 */}
      <PreviewModal 
        isOpen={isPreviewOpen}
        onClose={() => {
          setIsPreviewOpen(false);
          setPreviewItem(undefined);
        }}
        item={previewItem}
      />
    </>
  );
} 
