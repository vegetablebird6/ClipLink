'use client';

import { useState, useEffect, useRef } from 'react';
import { useChannel } from '@/contexts/ChannelContext';
import { clipboardService } from '@/services/api';
import { deviceIdUtil } from '@/utils/deviceId';
import { detectDeviceType, generateDeviceName } from '@/utils/deviceDetection';

// 本地存储的设备名称键名
const DEVICE_NAME_KEY = 'clipboard_device_name';

// 获取或生成设备名称
function getOrGenerateDeviceName(): string {
  // 尝试从本地存储获取设备名称
  const storedName = typeof window !== 'undefined' ? localStorage.getItem(DEVICE_NAME_KEY) : null;
  if (storedName) return storedName;

  // 生成默认设备名称（OS + Browser 格式）
  const defaultName = generateDeviceName();

  // 存储生成的名称
  if (typeof window !== 'undefined') localStorage.setItem(DEVICE_NAME_KEY, defaultName);
  return defaultName;
}

export function useDeviceRegistration() {
  const { channelId, isChannelVerified } = useChannel();
  const [deviceId] = useState(() => deviceIdUtil.getDeviceId());
  const [deviceName] = useState(getOrGenerateDeviceName);
  const [deviceType] = useState(detectDeviceType);
  const [registeredChannelId, setRegisteredChannelId] = useState<string | null>(null);
  const [isRegistering, setIsRegistering] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [retryCount, setRetryCount] = useState(0);
  // 用于在异步回调中检测 channel 是否已切换
  const currentChannelRef = useRef<string | null>(null);

  // channelId 变化时重置注册状态和错误状态，以便在新通道重新注册
  useEffect(() => {
    currentChannelRef.current = channelId;
    setRegisteredChannelId(null);
    setError(null);
    setRetryCount(0);
  }, [channelId]);

  // 当通道连接时，注册设备
  useEffect(() => {
    // 只在客户端执行
    if (typeof window === 'undefined') return;

    // 已在当前通道注册过，跳过
    if (registeredChannelId === channelId) return;

    // 如果达到最大重试次数，或没有有效的通道，则不继续
    if (retryCount >= 5 || !channelId || !isChannelVerified) return;

    // 如果正在注册，不重复发起
    if (isRegistering) return;

    const registerDevice = async () => {
      setIsRegistering(true);
      setError(null);

      try {
        const response = await clipboardService.registerDevice({
          device_id: deviceId,
          device_name: deviceName,
          device_type: deviceType
        });

        // 只在 channel 未切换过时接受结果（通过 ref 检测）
        if (response.success && currentChannelRef.current === channelId) {
          setRegisteredChannelId(channelId);
        } else if (currentChannelRef.current === channelId) {
          // 当前通道注册失败，计入重试
          const msg = response.message || '设备注册失败';
          setError(msg);
          setRetryCount(count => count + 1);
        }
        // 通道已切换则静默丢弃结果
      } catch (err) {
        if (currentChannelRef.current !== channelId) return;
        console.error('设备注册错误:', err);
        setError(err instanceof Error ? err.message : '设备注册失败');
        setRetryCount(count => count + 1);
      } finally {
        setIsRegistering(false);
      }
    };

    registerDevice();
  }, [channelId, isChannelVerified, deviceId, deviceName, deviceType, isRegistering, retryCount, registeredChannelId]);

  const isRegistered = registeredChannelId === channelId;

  // 定期发送心跳更新设备状态
  useEffect(() => {
    const currentRegistered = registeredChannelId === channelId;
    if (!channelId || !isChannelVerified || !currentRegistered) return;

    // 发送心跳更新
    const sendHeartbeat = async () => {
      try {
        // 使用clipboardService发送心跳
        const response = await clipboardService.updateDeviceStatus(deviceId, true);
        if (!response.success) {
          console.warn('心跳发送返回非成功状态:', response.message);
        }
      } catch (err) {
        console.error('发送心跳错误:', err);
      }
    };

    // 立即发送第一次心跳
    sendHeartbeat();

    // 设置定期发送心跳的间隔 - 每5分钟一次
    const heartbeatInterval = setInterval(sendHeartbeat, 300000);

    return () => {
      clearInterval(heartbeatInterval);
    };
  }, [channelId, isChannelVerified, deviceId, registeredChannelId]);

  // 在组件卸载或页面关闭时更新设备为离线状态
  useEffect(() => {
    const currentRegistered = registeredChannelId === channelId;
    if (!channelId || !isChannelVerified || !currentRegistered) return;

    const updateOfflineStatus = async () => {
      try {
        // 使用clipboardService更新离线状态
        await clipboardService.updateDeviceStatus(deviceId, false);
      } catch (err) {
        console.error('更新设备离线状态错误:', err);
      }
    };

    // 在页面卸载时更新设备状态为离线
    const handleBeforeUnload = () => {
      updateOfflineStatus();
    };

    window.addEventListener('beforeunload', handleBeforeUnload);

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
      updateOfflineStatus();
    };
  }, [channelId, isChannelVerified, deviceId, registeredChannelId]);

  return {
    deviceId,
    deviceName,
    deviceType,
    isRegistered,
    isRegistering,
    error
  };
}
