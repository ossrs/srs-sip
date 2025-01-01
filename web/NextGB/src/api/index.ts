import axios from 'axios'
import { ElMessage } from 'element-plus'
import type { Device, ApiResponse, ChannelInfo, MediaServer } from '@/types/api'

const api = axios.create({
  baseURL: import.meta.env.VITE_APP_API_BASE_URL,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 媒体服务器相关 API
export const mediaServerApi = {
  // 获取媒体服务器列表
  getMediaServers: () => api.get<ApiResponse<MediaServer[]>>('/srs-sip/v1/media-servers'),

  // 添加媒体服务器
  addMediaServer: (data: {
    name: string
    ip: string
    port: number
    type: string
    username?: string
    password?: string
  }) => api.post<ApiResponse<{ msg: string }>>('/srs-sip/v1/media-servers', data),

  // 删除媒体服务器
  deleteMediaServer: (id: number) =>
    api.delete<ApiResponse<{ msg: string }>>(`/srs-sip/v1/media-servers/${id}`),

  // 设置默认媒体服务器
  setDefaultMediaServer: (id: number) =>
    api.post<ApiResponse<{ msg: string }>>(`/srs-sip/v1/media-servers/default/${id}`),

  // 检查服务器状态
  checkStatus: (server: MediaServer) => {
    const url = `http://${server.ip}:${server.port}/api/v1/versions`
    return api.get<ApiResponse<any>>(url, { timeout: 2000 })
  },
}

// 设备相关 API
export const deviceApi = {
  // 获取设备列表
  getDevices: () => api.get<ApiResponse<Device[]>>('/srs-sip/v1/devices'),

  // 获取设备通道
  getDeviceChannels: (deviceId: string) =>
    api.get<ApiResponse<ChannelInfo[]>>(`/srs-sip/v1/devices/${deviceId}/channels`),

  // 添加 invite API
  inviteStream: (params: {
    media_server_addr: string
    device_id: string
    channel_id: string
    sub_stream: string
  }) => api.post<ApiResponse<any>>('/srs-sip/v1/invite', params),

  // 云台控制
  controlPTZ: (params: { device_id: string; channel_id: string; ptz: string; speed: string }) =>
    api.post<ApiResponse<any>>('/srs-sip/v1/ptz', params),
}

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 配置处理逻辑
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res
  },
  (error) => {
    ElMessage.error(error.response?.data?.message || '网络错误')
    return Promise.reject(error)
  },
)

export default api
