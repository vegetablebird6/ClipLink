'use client';

import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCopy, faFile, faLink, faLock, faXmark } from '@fortawesome/free-solid-svg-icons';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { oneLight } from 'react-syntax-highlighter/dist/cjs/styles/prism';
import AnimatedModal from '@/components/ui/AnimatedModal';
import { ClipboardItem, ClipboardType } from '@/types/clipboard';
import { detectLanguage } from '@/utils/codeHelper';
import { useToast } from '@/contexts/ToastContext';
import { writeClipboardRich } from '@/utils/richClipboard';
import RichTextRenderer from './RichTextRenderer';

interface PreviewModalProps {
  isOpen: boolean;
  onClose: () => void;
  item?: ClipboardItem;
}

export default function PreviewModal({ isOpen, onClose, item }: PreviewModalProps) {
  const { showToast } = useToast();

  const handleCopy = async () => {
    if (!item) return;

    try {
      await writeClipboardRich(item);
      showToast('已复制到剪贴板', 'success');
    } catch {
      showToast('复制失败', 'error');
    }
  };

  const renderContent = () => {
    if (!item) return null;

    // 富文本内容优先渲染
    if (item.content_format === 'html' && item.content_html) {
      return <RichTextRenderer html={item.content_html} />;
    }

    if (item.type === ClipboardType.CODE) {
      return (
        <SyntaxHighlighter
          language={detectLanguage(item.content)}
          style={oneLight}
          customStyle={{
            margin: 0,
            borderRadius: '0.5rem',
            maxHeight: '60vh',
            fontSize: '0.875rem'
          }}
          showLineNumbers
          wrapLongLines
        >
          {item.content}
        </SyntaxHighlighter>
      );
    }

    if (item.type === ClipboardType.LINK) {
      return (
        <a
          href={item.content}
          target="_blank"
          rel="noreferrer"
          className="break-all text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
        >
          {item.content}
        </a>
      );
    }

    if (item.type === ClipboardType.PASSWORD) {
      return <pre className="whitespace-pre-wrap break-words font-mono text-sm">{item.content}</pre>;
    }

    return <pre className="whitespace-pre-wrap break-words text-sm leading-6">{item.content}</pre>;
  };

  const getTypeIcon = () => {
    switch (item?.type) {
      case ClipboardType.LINK:
        return faLink;
      case ClipboardType.PASSWORD:
        return faLock;
      default:
        return faFile;
    }
  };

  return (
    <AnimatedModal isOpen={isOpen} onClose={onClose} maxWidth="max-w-3xl" showCloseButton={false}>
      <div className="overflow-hidden rounded-xl bg-white shadow-xl dark:bg-gray-900">
        <div className="flex items-center justify-between border-b border-gray-200 px-5 py-4 dark:border-gray-700">
          <div className="flex min-w-0 items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-blue-50 text-blue-600 dark:bg-blue-900/30 dark:text-blue-300">
              <FontAwesomeIcon icon={getTypeIcon()} />
            </div>
            <div className="min-w-0">
              <h2 className="truncate text-base font-semibold text-gray-900 dark:text-white">
                {item?.title || '内容预览'}
              </h2>
              <p className="text-xs text-gray-500 dark:text-gray-400">{item?.type || 'text'}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={handleCopy}
              className="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
              aria-label="复制"
            >
              <FontAwesomeIcon icon={faCopy} />
            </button>
            <button
              type="button"
              onClick={onClose}
              className="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
              aria-label="关闭"
            >
              <FontAwesomeIcon icon={faXmark} />
            </button>
          </div>
        </div>
        <div className="max-h-[70vh] overflow-auto p-5 text-gray-800 dark:text-gray-100">
          {renderContent()}
        </div>
      </div>
    </AnimatedModal>
  );
}
