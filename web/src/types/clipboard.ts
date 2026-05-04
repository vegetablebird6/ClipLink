// 剪贴板项目类型
export enum ClipboardType {
  TEXT = 'text',
  LINK = 'link',
  CODE = 'code',
  PASSWORD = 'password'
}

// 设备类型
export enum DeviceType {
  PHONE = 'phone',
  TABLET = 'tablet',
  DESKTOP = 'desktop',
  OTHER = 'other'
}

// API返回的原始剪贴板项目接口
export interface RawClipboardItem {
  id: string;
  title: string;
  content: string;
  type: ClipboardType;
  favorite: boolean;  // API返回使用favorite而不是isFavorite
  created_at: string; // API返回使用created_at而不是createdAt
  updated_at?: string;
  device_id: string;
  device_type?: DeviceType;
  content_html?: string;
  content_format?: 'plain' | 'html';
}

// 前端使用的剪贴板项目接口
export interface ClipboardItem {
  id: string;
  title: string;
  content: string;
  type: ClipboardType;
  isFavorite: boolean;   // 前端使用isFavorite
  created_at: string;    // 保持和API一致使用created_at
  updated_at?: string;
  device_id?: string;    // 可选字段
  device_type?: DeviceType; // 可选的设备类型字段
  content_html?: string;
  content_format?: 'plain' | 'html';
}

// 创建或更新剪贴板的请求
export interface SaveClipboardRequest {
  id?: string; // 添加id字段，用于编辑时传递
  title?: string;
  content: string;
  isFavorite?: boolean;
  type?: ClipboardType;
  device_id?: string;
  device_type?: DeviceType; // 添加设备类型字段
  clean_duplicates?: boolean;
  content_html?: string;
  content_format?: 'plain' | 'html';
}

// 通道统计
export interface ChannelStats {
  total_devices: number;
  online_devices: number;
  sync_count: number;
  clipboard_item_count: number;
  created_at: string;
}

// API响应类型 - 新的统一格式
export interface ApiResponse<T> {
  code: number;
  message: string;
  success: boolean;
  data?: T;
  error?: string;
}

// --- 后端 DTO 对齐类型 ---

// 通道
export interface ChannelResponse {
  id: string;
  created_at: string;
}

export interface ChannelDeleteResponse {
  channel_id: string;
  clipboard_items_deleted: number;
  sync_events_deleted: number;
  device_links_deleted: number;
  orphan_devices_deleted: number;
}

// 设备
export interface RegisterDeviceRequest {
  device_id: string;
  device_name: string;
  device_type: string;
}

export interface UpdateDeviceStatusRequest {
  is_online: boolean;
}

export interface UpdateDeviceNameRequest {
  device_name: string;
}

export interface DeviceResponse {
  id: string;
  name: string;
  type: string;
  channel_id: string;
  last_seen: string;
  is_online: boolean;
  created_at: string;
  joined_at: string;
}

// 同步事件
export interface SyncEventResponse {
  id: number;
  channel_id: string;
  action: string;
  target_type: string;
  target_id: string;
  content: string;
  summary: string;
  actor_device_id: string;
  actor_device_name: string;
  actor_device_type: string;
  created_at: string;
}

// 通用分页
export interface KeysetPageResponse<T> {
  items: T[];
  has_more: boolean;
  next_after?: string;
  next_after_id?: string;
}

export interface OffsetPageResponse<T> {
  items: T[];
  total: number;
  page: number;
  size: number;
  totalPages: number;
} 
