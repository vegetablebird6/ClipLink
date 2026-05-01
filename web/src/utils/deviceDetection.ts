import { DeviceType } from '@/types/clipboard';

/**
 * 检测当前设备类型
 * 简化版：只区分电脑和手机
 * @returns {DeviceType} 设备类型枚举值
 */
export function detectDeviceType(): DeviceType {
  // 只在客户端执行
  if (typeof window === 'undefined' || typeof navigator === 'undefined') {
    return DeviceType.OTHER;
  }
  
  const userAgent = navigator.userAgent.toLowerCase();
  
  // 简化检测：只检查是否是移动设备
  const isMobile = /android|webos|iphone|ipod|blackberry|iemobile|opera mini|mobile/i.test(userAgent) || 
                   (window.innerWidth < 768); // 额外使用屏幕宽度判断
  
  if (isMobile) {
    return DeviceType.PHONE;
  } else {
    return DeviceType.DESKTOP;
  }
}

/**
 * 获取设备类型的本地化描述
 * @param {DeviceType} deviceType 设备类型
 * @returns {string} 设备类型的中文描述
 */
export function getDeviceTypeLabel(deviceType: DeviceType): string {
  switch (deviceType) {
    case DeviceType.DESKTOP:
      return '电脑';
    case DeviceType.PHONE:
      return '手机';
    case DeviceType.TABLET:
      return '平板';
    default:
      return '其他设备';
  }
}

/**
 * 生成友好的默认设备名称，格式：OS Browser
 * 无法识别时 fallback 到设备类型的中文描述
 */
export function generateDeviceName(): string {
  if (typeof navigator === 'undefined') return '未知设备';

  const ua = navigator.userAgent;

  // OS
  let os = '未知';
  if (/iPhone/.test(ua))        os = 'iPhone';
  else if (/iPad/.test(ua))     os = 'iPad';
  else if (/Android/.test(ua))  os = 'Android';
  else if (/Windows/.test(ua))  os = 'Windows';
  else if (/Mac OS X/.test(ua)) os = 'Mac';
  else if (/Linux/.test(ua))    os = 'Linux';

  // Browser（顺序重要：Edge/Opera 在 Chrome 之前匹配）
  let browser = '';
  if (/Edg\//.test(ua))          browser = 'Edge';
  else if (/OPR\//.test(ua))     browser = 'Opera';
  else if (/Chrome\//.test(ua))  browser = 'Chrome';
  else if (/Firefox\//.test(ua)) browser = 'Firefox';
  else if (/Safari\//.test(ua))  browser = 'Safari';

  return browser ? `${os} ${browser}` : os;
}