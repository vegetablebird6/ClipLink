import type { Metadata } from 'next';
import { Inter } from "next/font/google";
import "./globals.css";
import AppProviders from './providers';

const inter = Inter({ subsets: ["latin"], variable: '--font-inter' });

const siteUrl = process.env.NEXT_PUBLIC_SITE_URL || 'https://cliplink.mmmss.com';

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  applicationName: 'ClipLink',
  title: {
    default: 'ClipLink - 跨平台剪贴板共享工具',
    template: '%s | ClipLink',
  },
  description: 'ClipLink 是一个轻量、安全的跨平台剪贴板同步工具，可在电脑、手机和平板之间通过网页共享文本、链接、代码和密码等内容。',
  keywords: ['ClipLink', '剪贴板同步', '跨设备剪贴板', '跨平台剪贴板', 'clipboard sync'],
  authors: [{ name: 'ClipLink' }],
  creator: 'ClipLink',
  publisher: 'ClipLink',
  icons: {
    icon: '/favicon.ico',
  },
  openGraph: {
    type: 'website',
    url: '/',
    siteName: 'ClipLink',
    title: 'ClipLink - 跨平台剪贴板共享工具',
    description: '通过网页在多台设备之间快速、安全地共享剪贴板内容。',
    locale: 'zh_CN',
  },
  twitter: {
    card: 'summary',
    title: 'ClipLink - 跨平台剪贴板共享工具',
    description: '通过网页在多台设备之间快速、安全地共享剪贴板内容。',
  },
  robots: {
    index: true,
    follow: true,
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body className={`${inter.className} bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100`}>
        <AppProviders>{children}</AppProviders>
      </body>
    </html>
  );
}
