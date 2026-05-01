'use client';

import React, {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState
} from 'react';
import { usePathname } from 'next/navigation';

interface NavigationProgressContextValue {
  isNavigating: boolean;
  startNavigation: (href: string) => void;
}

const NavigationProgressContext = createContext<NavigationProgressContextValue | undefined>(undefined);

export function NavigationProgressProvider({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  const [isNavigating, setIsNavigating] = useState(false);
  const pendingHrefRef = useRef<string | null>(null);
  const fallbackTimerRef = useRef<number | null>(null);

  const clearFallbackTimer = useCallback(() => {
    if (fallbackTimerRef.current !== null) {
      window.clearTimeout(fallbackTimerRef.current);
      fallbackTimerRef.current = null;
    }
  }, []);

  const startNavigation = useCallback((href: string) => {
    if (href === pathname) {
      return;
    }

    pendingHrefRef.current = href;
    setIsNavigating(true);
    clearFallbackTimer();

    fallbackTimerRef.current = window.setTimeout(() => {
      pendingHrefRef.current = null;
      setIsNavigating(false);
      fallbackTimerRef.current = null;
    }, 8000);
  }, [clearFallbackTimer, pathname]);

  useEffect(() => {
    if (!isNavigating || pendingHrefRef.current !== pathname) {
      return;
    }

    const settleTimer = window.setTimeout(() => {
      pendingHrefRef.current = null;
      setIsNavigating(false);
      clearFallbackTimer();
    }, 180);

    return () => window.clearTimeout(settleTimer);
  }, [clearFallbackTimer, isNavigating, pathname]);

  useEffect(() => clearFallbackTimer, [clearFallbackTimer]);

  return (
    <NavigationProgressContext.Provider value={{ isNavigating, startNavigation }}>
      <RouteProgress visible={isNavigating} />
      {children}
    </NavigationProgressContext.Provider>
  );
}

export function useNavigationProgress() {
  const context = useContext(NavigationProgressContext);

  if (!context) {
    throw new Error('useNavigationProgress must be used within NavigationProgressProvider');
  }

  return context;
}

function RouteProgress({ visible }: { visible: boolean }) {
  if (!visible) {
    return null;
  }

  return (
    <div className="fixed inset-x-0 top-0 z-[80] pointer-events-none">
      <div className="h-0.5 w-full overflow-hidden bg-brand-100/80 dark:bg-brand-900/30">
        <div className="h-full w-1/3 animate-route-progress rounded-full bg-brand-600 shadow-glow-brand dark:bg-brand-dark-400" />
      </div>
      <div className="absolute right-4 top-3 hidden items-center gap-2 rounded-full border border-white/40 bg-white/85 px-3 py-1.5 text-xs font-medium text-neutral-700 shadow-lg backdrop-blur-md dark:border-dark-border-primary/60 dark:bg-dark-surface-primary/85 dark:text-dark-text-secondary sm:flex">
        <span className="relative h-3.5 w-3.5">
          <span className="absolute inset-0 rounded-full border-2 border-brand-200 dark:border-brand-900" />
          <span className="absolute inset-0 rounded-full border-2 border-brand-600 border-t-transparent animate-spin dark:border-brand-dark-400 dark:border-t-transparent" />
        </span>
        切换中
      </div>
    </div>
  );
}
