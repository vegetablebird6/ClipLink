'use client';

import { ReactNode } from 'react';
import { config } from '@fortawesome/fontawesome-svg-core';
import { ToastProvider } from '@/contexts/ToastContext';
import { ChannelProvider } from '@/contexts/ChannelContext';
import DeviceRegistration from '@/components/device/DeviceRegistration';
import MainLayout from '@/components/layout/MainLayout';

// Prevent Font Awesome icons from flashing while keeping this browser-only setup
// out of the server layout.
config.autoAddCss = true;

interface AppProvidersProps {
  children: ReactNode;
}

export default function AppProviders({ children }: AppProvidersProps) {
  return (
    <ToastProvider>
      <ChannelProvider>
        <DeviceRegistration />
        <MainLayout>{children}</MainLayout>
      </ChannelProvider>
    </ToastProvider>
  );
}
