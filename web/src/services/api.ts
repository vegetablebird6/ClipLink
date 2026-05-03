import axios, { AxiosError } from 'axios';
import { ApiResponse, ClipboardItem, SaveClipboardRequest, ClipboardType, RawClipboardItem } from '@/types/clipboard';
import { deviceIdUtil } from '@/utils/deviceId';
import { detectDeviceType } from '@/utils/deviceDetection';
import { settingsManager } from '@/utils/settings';

// 获取通道ID
const getChannelId = (): string | null => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('clipboard_channel_id');
  }
  return null;
};

const baseUrl = process.env.NODE_ENV === 'development' ? 'http://localhost:8080/api' : '/api';
// 创建Axios实例
const api = axios.create({
  baseURL: baseUrl,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 添加请求拦截器，自动添加deviceId和channelId到每个请求
api.interceptors.request.use((config) => {
  // 获取通道ID
  const channelId = getChannelId();
  config.headers['X-Channel-ID'] = channelId;
  // 如果是GET请求，添加deviceId到查询参数
  if (config.method?.toLowerCase() === 'get') {
    config.params = {
      ...config.params,
      device_id: deviceIdUtil.getDeviceId()
    };
  }
  return config;
}, (error) => {
  return Promise.reject(error);
});

// 处理统一API响应格式
const handleApiResponse = <T>(response: unknown): ApiResponse<T> => {
 if (response && typeof response === 'object' && 'code' in response && 
      typeof response.code === 'number' && 
      'success' in response && 
      typeof response.success === 'boolean') {
    return response as ApiResponse<T>;
  }
  
  return {
    code: 500,
    message: 'API响应格式错误',
    success: false,
    error: 'API响应格式错误'
  };
};

// 处理API错误
const handleApiError = <T>(error: unknown, defaultMessage: string): ApiResponse<T> => {
  let errorMessage = defaultMessage;
  let statusCode = 500;
  
  if (axios.isAxiosError(error)) {
    // 服务器返回了错误状态码
    const axiosError = error as AxiosError;
    if (axiosError.response) {
      statusCode = axiosError.response.status;
      const responseData = axiosError.response.data as Record<string, unknown>;
      errorMessage = `服务器错误 (${statusCode}): ${
        responseData && typeof responseData === 'object' && 'message' in responseData 
          ? String(responseData.message) 
          : axiosError.message
      }`;
    } else if (axiosError.request) {
      // 请求发出但没有收到响应
      errorMessage = '服务器无响应，请检查网络连接';
      statusCode = 503;
    } else {
      // 其他错误
      errorMessage = axiosError.message || '未知错误';
    }
  } else if (error instanceof Error) {
    errorMessage = error.message;
  }
  
  return {
    code: statusCode,
    message: errorMessage,
    success: false,
    error: errorMessage
  };
};

// 将API返回的原始格式转换为前端使用的格式
type ClipboardItemResponse = RawClipboardItem & {
  isFavorite?: boolean;
};

const convertRawClipboardItem = (raw: ClipboardItemResponse): ClipboardItem => {
  return {
    id: raw.id,
    title: raw.title || '',
    content: raw.content,
    type: raw.type,
    isFavorite: raw.favorite || raw.isFavorite || false,
    created_at: raw.created_at || '',
    updated_at: raw.updated_at,
    device_id: raw.device_id,
    device_type: raw.device_type,
    content_html: raw.content_html || undefined,
    content_format: raw.content_format || undefined,
  };
};

// API服务类
export const clipboardService = {
  // 通道相关接口
  // 创建通道
  createChannel: async (channelId?: string, instanceToken?: string): Promise<ApiResponse<{id: string, created_at: string}>> => {
    try {
      const response = await api.post<unknown>(
        '/channel',
        channelId ? { channel_id: channelId } : undefined,
        instanceToken ? { headers: { 'X-Instance-Token': instanceToken } } : undefined
      );
      return handleApiResponse<{id: string, created_at: string}>(response.data);
    } catch (error) {
      return handleApiError<{id: string, created_at: string}>(error, '创建通道失败');
    }
  },

  // 验证通道（header自动带channelId）
  verifyChannel: async (channelId?: string): Promise<ApiResponse<null>> => {
    try {
      // 如果提供了channelId参数，在请求体中传递
      const data = channelId ? { channel_id: channelId } : undefined;
      const response = await api.post<unknown>('/channel/verify', data);
      return handleApiResponse<null>(response.data);
    } catch (error) {
      return handleApiError<null>(error, '验证通道失败');
    }
  },

  // 获取当前通道信息（header自动带channelId）
  getChannel: async (): Promise<ApiResponse<any>> => {
    try {
      const response = await api.get<unknown>('/channel');
      return handleApiResponse<any>(response.data);
    } catch (error) {
      return handleApiError<any>(error, '获取通道信息失败');
    }
  },

  // 获取通道统计（header自动带channelId）
  getChannelStats: async (): Promise<ApiResponse<any>> => {
    try {
      const response = await api.get<unknown>('/stats');
      return handleApiResponse<any>(response.data);
    } catch (error) {
      return handleApiError<any>(error, '获取统计数据失败');
    }
  },

  // 删除当前通道及其关联数据
  deleteChannel: async (instanceToken: string): Promise<ApiResponse<any>> => {
    try {
      const response = await api.delete<unknown>('/channel', {
        headers: { 'X-Instance-Token': instanceToken }
      });
      return handleApiResponse<any>(response.data);
    } catch (error) {
      return handleApiError<any>(error, '删除通道失败');
    }
  },

  // 设备相关接口
  // 注册设备
  registerDevice: async (deviceData: {
    device_id: string;
    device_name: string;
    device_type: string;
  }): Promise<ApiResponse<any>> => {
    try {
      const response = await api.post<unknown>('/devices', deviceData);
      return handleApiResponse<any>(response.data);
    } catch (error) {
      return handleApiError<any>(error, '设备注册失败');
    }
  },

  // 获取设备列表
  getDevices: async (): Promise<ApiResponse<any[]>> => {
    try {
      const response = await api.get<unknown>('/devices');
      return handleApiResponse<any[]>(response.data);
    } catch (error) {
      return handleApiError<any[]>(error, '获取设备列表失败');
    }
  },

  // 更新设备状态
  updateDeviceStatus: async (deviceId: string, isOnline: boolean): Promise<ApiResponse<any>> => {
    try {
      const response = await api.put<unknown>(`/devices/${deviceId}/status`, { is_online: isOnline });
      return handleApiResponse<any>(response.data);
    } catch (error) {
      return handleApiError<any>(error, '更新设备状态失败');
    }
  },

  // 删除设备
  removeDevice: async (deviceId: string): Promise<ApiResponse<null>> => {
    try {
      const response = await api.delete<unknown>(`/devices/${deviceId}`);
      return handleApiResponse<null>(response.data);
    } catch (error) {
      return handleApiError<null>(error, '删除设备失败');
    }
  },

  // 更新设备名称
  updateDeviceName: async (deviceId: string, name: string): Promise<ApiResponse<any>> => {
    try {
      const response = await api.put<unknown>(`/devices/${deviceId}/name`, { device_name: name });
      return handleApiResponse<any>(response.data);
    } catch (error) {
      return handleApiError<any>(error, '更新设备名称失败');
    }
  },

  // 获取同步历史（keyset 游标分页）
  getSyncHistory: async (limit: number = 20, after?: string, afterId?: string): Promise<ApiResponse<{items: any[], has_more: boolean}>> => {
    try {
      const params: Record<string, string | number> = { limit };
      if (after && afterId) {
        params.after = after;
        params.after_id = afterId;
      }
      const response = await api.get<unknown>(`/sync/history`, { params });
      return handleApiResponse<{items: any[], has_more: boolean}>(response.data);
    } catch (error) {
      return handleApiError<{items: any[], has_more: boolean}>(error, '获取同步历史失败');
    }
  },

  // 获取最新剪贴板内容
  getLatestClipboard: async (): Promise<ApiResponse<ClipboardItem>> => {
    try {
      const response = await api.get<unknown>('/clipboard/current');
      const apiResponse = handleApiResponse<any>(response.data);
      
      // 转换为标准格式
      if (apiResponse.success && apiResponse.data) {
        apiResponse.data = convertRawClipboardItem(apiResponse.data);
      }
      
      return apiResponse as ApiResponse<ClipboardItem>;
    } catch (error) {
      return handleApiError<ClipboardItem>(error, '获取最新剪贴板失败');
    }
  },

  // 获取当前剪贴板内容（专用接口，确保始终能获取到内容）
  getCurrentClipboard: async (): Promise<ApiResponse<ClipboardItem>> => {
    try {
      const response = await api.get<unknown>('/clipboard/current');
      const apiResponse = handleApiResponse<any>(response.data);
      
      // 处理空数组或没有数据的情况
      if (apiResponse.success) {
        if (Array.isArray(apiResponse.data)) {
          // 如果返回的是数组而不是单个对象
          if (apiResponse.data.length > 0) {
            // 取第一个项目
            apiResponse.data = convertRawClipboardItem(apiResponse.data[0]);
          } else {
            // 空数组，设置data为null
            apiResponse.data = null;
          }
        } else if (apiResponse.data) {
          // 如果是单个对象，转换为标准格式
          apiResponse.data = convertRawClipboardItem(apiResponse.data);
        } else {
          // 没有数据，设置data为null
          apiResponse.data = null;
        }
      }
      
      return apiResponse as ApiResponse<ClipboardItem>;
    } catch (error) {
      return handleApiError<ClipboardItem>(error, '获取当前剪贴板失败');
    }
  },

  // 获取剪贴板历史（keyset 游标分页）
  getClipboardHistory: async (size = 12, after?: string, afterId?: string): Promise<ApiResponse<{items: ClipboardItem[], has_more: boolean}>> => {
    try {
      const params: Record<string, string | number> = { size };
      if (after && afterId) {
        params.after = after;
        params.after_id = afterId;
      }

      const response = await api.get<unknown>('/clipboard/history', { params });
      const apiResponse = handleApiResponse<any>(response.data);

      // 转换 items
      if (apiResponse.success && apiResponse.data) {
        if (apiResponse.data.items) {
          apiResponse.data.items = apiResponse.data.items.map(convertRawClipboardItem);
        }
      }

      return apiResponse as ApiResponse<{items: ClipboardItem[], has_more: boolean}>;
    } catch (error) {
      return handleApiError<{items: ClipboardItem[], has_more: boolean}>(error, '获取剪贴板历史失败');
    }
  },

  // 按类型获取剪贴板历史（keyset 游标分页）
  getClipboardByType: async (type: ClipboardType, size = 12, after?: string, afterId?: string): Promise<ApiResponse<{items: ClipboardItem[], has_more: boolean}>> => {
    try {
      const params: Record<string, string | number> = { size };
      if (after && afterId) {
        params.after = after;
        params.after_id = afterId;
      }

      const response = await api.get<unknown>(`/clipboard/type/${type}`, { params });
      const apiResponse = handleApiResponse<any>(response.data);

      if (apiResponse.success && apiResponse.data) {
        if (apiResponse.data.items) {
          apiResponse.data.items = apiResponse.data.items.map(convertRawClipboardItem);
        }
      }

      return apiResponse as ApiResponse<{items: ClipboardItem[], has_more: boolean}>;
    } catch (error) {
      return handleApiError<{items: ClipboardItem[], has_more: boolean}>(error, '按类型获取剪贴板历史失败');
    }
  },

  // 获取收藏的剪贴板项目
  getFavorites: async (limit = 50): Promise<ApiResponse<ClipboardItem[]>> => {
    try {
      const response = await api.get<unknown>('/clipboard/favorites', {
        params: { limit }
      });
      
      const apiResponse = handleApiResponse<any>(response.data);
      
      if (apiResponse.success && apiResponse.data) {
        apiResponse.data = (apiResponse.data as any[]).map(convertRawClipboardItem);
      }
      
      return apiResponse as ApiResponse<ClipboardItem[]>;
    } catch (error) {
      return handleApiError<ClipboardItem[]>(error, '获取收藏夹失败');
    }
  },

  // 保存剪贴板内容
  saveClipboard: async (data: Omit<SaveClipboardRequest, 'device_id' | 'device_type'>): Promise<ApiResponse<ClipboardItem>> => {
    try {
      // 添加设备ID和设备类型
      const requestData: SaveClipboardRequest = {
        ...data,
        device_id: deviceIdUtil.getDeviceId(),
        device_type: detectDeviceType(),
        clean_duplicates: settingsManager.getSetting('autoCleanDuplicates')
      };
      
      const response = await api.post<unknown>('/clipboard', requestData);
      const apiResponse = handleApiResponse<any>(response.data);
      
      // 转换为标准格式
      if (apiResponse.success && apiResponse.data) {
        apiResponse.data = convertRawClipboardItem(apiResponse.data);
      }
      
      return apiResponse as ApiResponse<ClipboardItem>;
    } catch (error) {
      return handleApiError<ClipboardItem>(error, '保存剪贴板失败');
    }
  },

  // 更新剪贴板项目
  updateClipboard: async (id: string, data: SaveClipboardRequest): Promise<ApiResponse<ClipboardItem>> => {
    try {
      // 确保请求数据中包含设备ID
      const requestData: SaveClipboardRequest = {
        ...data,
        device_id: data.device_id || deviceIdUtil.getDeviceId(),
        device_type: data.device_type || detectDeviceType()
      };
      const response = await api.put<unknown>(`/clipboard/${id}`, requestData);
      const apiResponse = handleApiResponse<any>(response.data);
      
      // 转换为标准格式
      if (apiResponse.success && apiResponse.data) {
        apiResponse.data = convertRawClipboardItem(apiResponse.data);
      }
      
      return apiResponse as ApiResponse<ClipboardItem>;
    } catch (error) {
      return handleApiError<ClipboardItem>(error, '更新剪贴板失败');
    }
  },

  // 切换收藏状态
  toggleFavorite: async (id: string, favorite: boolean): Promise<ApiResponse<ClipboardItem>> => {
    try {
      const response = await api.put<unknown>(`/clipboard/${id}/favorite`, { 
        favorite,
        device_id: deviceIdUtil.getDeviceId(),
      });
      const apiResponse = handleApiResponse<any>(response.data);
      
      // 转换为标准格式
      if (apiResponse.success && apiResponse.data) {
        apiResponse.data = convertRawClipboardItem(apiResponse.data);
      }
      
      return apiResponse as ApiResponse<ClipboardItem>;
    } catch (error) {
      return handleApiError<ClipboardItem>(error, '切换收藏状态失败');
    }
  },

  // 删除剪贴板项目
  deleteClipboard: async (id: string): Promise<ApiResponse<null>> => {
    try {
      const response = await api.delete<unknown>(`/clipboard/${id}`, {
        data: { device_id: deviceIdUtil.getDeviceId() }
      });
      return handleApiResponse<null>(response.data);
    } catch (error) {
      return handleApiError<null>(error, '删除剪贴板失败');
    }
  },

  // 清理当前通道下已存在的重复剪贴板内容
  cleanupDuplicates: async (): Promise<ApiResponse<{deleted: number}>> => {
    try {
      const response = await api.post<unknown>('/clipboard/cleanup-duplicates');
      return handleApiResponse<{deleted: number}>(response.data);
    } catch (error) {
      return handleApiError<{deleted: number}>(error, '清理重复内容失败');
    }
  },

  // 获取指定类型的剪贴板项目数量
  getClipboardCount: async (type?: ClipboardType): Promise<ApiResponse<Record<string, number>>> => {
    try {
      const params: Record<string, string | undefined> = {};
      if (type) {
        params.type = type;
      }
      
      const response = await api.get<unknown>('/clipboard/count', { params });
      return handleApiResponse<Record<string, number>>(response.data);
    } catch (error) {
      return handleApiError<Record<string, number>>(error, '获取剪贴板数量失败');
    }
  },

  // 获取单个剪贴板项目
  getClipboardItem: async (id: string): Promise<ApiResponse<ClipboardItem>> => {
    try {
      const response = await api.get<unknown>(`/clipboard/${id}`);
      const apiResponse = handleApiResponse<any>(response.data);
      
      // 转换为标准格式
      if (apiResponse.success && apiResponse.data) {
        apiResponse.data = convertRawClipboardItem(apiResponse.data);
      }
      
      return apiResponse as ApiResponse<ClipboardItem>;
    } catch (error) {
      return handleApiError(error, '获取剪贴板项目失败');
    }
  },

  // 搜索剪贴板项目
  searchClipboard: async (keyword: string, page = 1, size = 12): Promise<ApiResponse<{items: ClipboardItem[], total: number, page: number, size: number, totalPages: number, keyword: string}>> => {
    try {
      const response = await api.get<unknown>('/clipboard/search', {
        params: {
          q: keyword,
          page,
          size
        }
      });
      
      const apiResponse = handleApiResponse<any>(response.data);
      
      // 转换每个项目为标准格式
      if (apiResponse.success && apiResponse.data && apiResponse.data.items) {
        apiResponse.data.items = apiResponse.data.items.map(convertRawClipboardItem);
      }
      
      return apiResponse as ApiResponse<{items: ClipboardItem[], total: number, page: number, size: number, totalPages: number, keyword: string}>;
    } catch (error) {
      return handleApiError<{items: ClipboardItem[], total: number, page: number, size: number, totalPages: number, keyword: string}>(error, '搜索失败');
    }
  }
}; 
