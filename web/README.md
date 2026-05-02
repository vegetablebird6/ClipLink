# ClipLink 前端项目

## 项目概述

ClipLink 是一个跨设备剪贴板同步工具，允许用户在不同设备之间共享和同步剪贴板内容。本项目是 ClipLink 的前端部分，使用 React、Next.js 和 Tailwind CSS 构建。

## 技术栈

- React 19 - 用户界面库
- Next.js 16 - React框架
- TypeScript - 类型安全的JavaScript超集
- Tailwind CSS 4 - 实用优先的CSS框架
- Axios - HTTP客户端
- Font Awesome - 图标库

## 功能特性

- 实时剪贴板同步
- 支持文本、链接和代码片段
- 剪贴板历史记录
- 收藏功能
- 多设备同步
- 响应式设计，支持桌面和移动设备

## 项目结构

```
web/
├── src/
│   ├── app/                    # Next.js App Router目录
│   │   ├── favorites/          # 收藏页面
│   │   ├── history/            # 历史记录页面
│   │   ├── layout.tsx          # 应用布局
│   │   ├── page.tsx            # 主页
│   │   └── globals.css         # 全局样式
│   ├── components/             # React组件
│   │   ├── clipboard/          # 剪贴板相关组件
│   │   └── layout/             # 布局组件
│   ├── services/               # API服务
│   └── types/                  # TypeScript类型定义
├── public/                     # 静态资源
├── tailwind.config.js          # Tailwind CSS兼容配置
├── postcss.config.js           # Tailwind CSS v4 PostCSS插件配置
└── package.json                # 项目依赖
```

## 开发指南

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

### 生产构建

```bash
npm run build
npm start
```

### Tailwind CSS 4

Tailwind v4 的入口在 `src/app/globals.css`：

```css
@import 'tailwindcss';
@config '../../tailwind.config.js';
@plugin 'tailwind-scrollbar';
```

PostCSS 使用 `@tailwindcss/postcss`，不再需要 `postcss-import` 和 `autoprefixer`。修改样式体系后建议运行：

```bash
npm exec tsc -- --noEmit
npm run build
```

## API接口

前端通过以下API与后端通信：

- `GET /api/clipboard` - 获取最新剪贴板内容
- `GET /api/clipboard/history` - 获取剪贴板历史记录
- `GET /api/clipboard/favorites` - 获取收藏的剪贴板内容
- `POST /api/clipboard` - 保存新的剪贴板内容
- `PUT /api/clipboard/:id` - 更新剪贴板内容
- `PUT /api/clipboard/:id/favorite` - 切换收藏状态
- `DELETE /api/clipboard/:id` - 删除剪贴板内容
- `GET /api/clipboard/:id` - 获取单个剪贴板项目
