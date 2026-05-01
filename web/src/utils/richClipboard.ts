import type { ClipboardItem } from '@/types/clipboard';

export interface RichClipboardContent {
  text: string;
  html?: string;
  format: 'plain' | 'html';
}

export async function readClipboardRich(): Promise<RichClipboardContent> {
  if (!navigator.clipboard?.read) {
    const text = await navigator.clipboard.readText();
    return { text, format: 'plain' };
  }

  const items = await navigator.clipboard.read();

  for (const item of items) {
    const hasHtml = item.types.includes('text/html');
    const hasText = item.types.includes('text/plain');

    if (hasHtml) {
      const htmlBlob = await item.getType('text/html');
      const html = await htmlBlob.text();

      let text = '';
      if (hasText) {
        const textBlob = await item.getType('text/plain');
        text = await textBlob.text();
      } else {
        text = html.replace(/<[^>]+>/g, '').trim();
      }

      return { text, html, format: 'html' };
    }
  }

  const text = await navigator.clipboard.readText();
  return { text, format: 'plain' };
}

export async function writeClipboardRich(item: ClipboardItem): Promise<void> {
  if (
    item.content_format === 'html' &&
    item.content_html &&
    navigator.clipboard?.write &&
    typeof ClipboardItem !== 'undefined'
  ) {
    try {
      await navigator.clipboard.write([
        new globalThis.ClipboardItem({
          'text/plain': new Blob([item.content], { type: 'text/plain' }),
          'text/html': new Blob([item.content_html], { type: 'text/html' }),
        }),
      ]);
      return;
    } catch {
      // 富文本写入失败，回退到纯文本
    }
  }

  await navigator.clipboard.writeText(item.content);
}
