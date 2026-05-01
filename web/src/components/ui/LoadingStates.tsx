import React from 'react';

export function ClipboardGridSkeleton({ count = 12 }: { count?: number }) {
  return (
    <div className="w-full pb-3" aria-busy="true" aria-label="剪贴板内容加载中">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 3xl:grid-cols-6 4xl:grid-cols-7 gap-2.5">
        {Array.from({ length: count }).map((_, index) => (
          <div
            key={index}
            className="h-40 overflow-hidden rounded-xl border border-white/30 bg-white/75 shadow-md backdrop-blur-md dark:border-dark-border-primary/30 dark:bg-dark-surface-primary/75 dark:shadow-dark-md"
          >
            <div className="flex h-full flex-col p-3">
              <div className="mb-3 flex items-center justify-between">
                <div className="flex items-center">
                  <span className="h-6 w-6 rounded-lg bg-neutral-200 dark:bg-dark-surface-tertiary skeleton-pulse" />
                  <span className="ml-2 h-3 w-24 rounded bg-neutral-200 dark:bg-dark-surface-tertiary skeleton-pulse" />
                </div>
                <span className="h-4 w-4 rounded-full bg-neutral-200 dark:bg-dark-surface-tertiary skeleton-pulse" />
              </div>
              <div className="space-y-2">
                <span className="block h-3 rounded bg-neutral-200 dark:bg-dark-surface-tertiary skeleton-pulse" />
                <span className="block h-3 w-5/6 rounded bg-neutral-200 dark:bg-dark-surface-tertiary skeleton-pulse" />
                <span className="block h-3 w-2/3 rounded bg-neutral-200 dark:bg-dark-surface-tertiary skeleton-pulse" />
              </div>
              <div className="mt-auto flex items-center justify-between border-t border-neutral-100 pt-2 dark:border-dark-border-primary/50">
                <span className="h-3 w-20 rounded bg-neutral-200 dark:bg-dark-surface-tertiary skeleton-pulse" />
                <div className="flex gap-1">
                  <span className="h-6 w-6 rounded-lg bg-neutral-200 dark:bg-dark-surface-tertiary skeleton-pulse" />
                  <span className="h-6 w-6 rounded-lg bg-neutral-200 dark:bg-dark-surface-tertiary skeleton-pulse" />
                  <span className="h-6 w-6 rounded-lg bg-neutral-200 dark:bg-dark-surface-tertiary skeleton-pulse" />
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export function PageLoadingFallback({ label = '加载中' }: { label?: string }) {
  return (
    <div className="flex h-full min-h-64 flex-col items-center justify-center gap-3 bg-gray-50 text-neutral-500 dark:bg-gradient-dark dark:text-dark-text-tertiary">
      <span className="relative h-8 w-8">
        <span className="absolute inset-0 rounded-full border-2 border-neutral-200 dark:border-dark-border-secondary" />
        <span className="absolute inset-0 rounded-full border-2 border-brand-600 border-t-transparent animate-spin dark:border-brand-dark-400 dark:border-t-transparent" />
      </span>
      <span className="text-sm font-medium">{label}</span>
    </div>
  );
}
