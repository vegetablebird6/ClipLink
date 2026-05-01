'use client';

import DOMPurify from 'dompurify';

interface RichTextRendererProps {
  html: string;
}

export default function RichTextRenderer({ html }: RichTextRendererProps) {
  const cleanHtml = DOMPurify.sanitize(html, {
    USE_PROFILES: { html: true },
    FORBID_TAGS: ['style', 'script'],
    FORBID_ATTR: ['onerror', 'onclick', 'onload'],
  });

  return (
    <div
      className="rich-text-content text-sm leading-6 break-words"
      dangerouslySetInnerHTML={{ __html: cleanHtml }}
    />
  );
}
