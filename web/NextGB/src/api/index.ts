import axios from 'axios'
import { ElMessage } from 'element-plus'
import type { Device, ApiResponse, ChannelInfo } from '@/types/api'

const api = axios.create({
  baseURL: import.meta.env.VITE_APP_API_BASE_URL,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 定义 summaries 响应类型
interface SRSSummaries {
  code: number
  server: {
    id: string
    pid: number
    argv: string
    cwd: string
    srs_version: string
    version: number
    peers: number
    now_ms: number
  }
  system: {
    cpu_percent: number
    disk_read_KBps: number
    disk_write_KBps: number
    mem_percent: number
    mem_ram_KBytes: number
    mem_swap_KBytes: number
    net_recv_bytes: number
    net_send_bytes: number
    srs_cpu_percent: number
    srs_mem_percent: number
    uptime_seconds: number
  }
}

// 定义 versions 响应类型
interface SRSVersions {
  code: number
  server: string
  service: string
  pid: string
  data: {
    major: number
    minor: number
    revision: number
    version: string
  }
}

// 服务器相关 API
export const serverApi = {
  // 获取媒体服务器信息
  getMediaServer: () => api.get<ApiResponse<{ type: string; address: string }>>('/srs-sip/v1/media-server'),
  
  // 检查服务器状态
  checkStatus: (address: string) => 
    api.get<SRSVersions>(`http://${address}/api/v1/versions`),
    
  // 获取服务器详细信息
  getSummaries: (address: string) =>
    api.get<SRSSummaries>(`http://${address}/api/v1/summaries`),
}

// 设备相关 API
export const deviceApi = {
  // 获取设备列表
  getDevices: () => api.get<ApiResponse<Device[]>>('/srs-sip/v1/devices'),

  // 添加设备
  addDevice: (device: Omit<Device, 'device_id'>) =>
    api.post<ApiResponse<Device>>('/srs-sip/v1/devices', device),

  // 更新设备
  updateDevice: (deviceId: string, device: Partial<Device>) =>
    api.put<ApiResponse<Device>>(`/srs-sip/v1/devices/${deviceId}`, device),

  // 删除设备
  deleteDevice: (deviceId: string) =>
    api.delete<ApiResponse<null>>(`/srs-sip/v1/devices/${deviceId}`),

  // 获取设备通道
  getDeviceChannels: (deviceId: string) =>
    api.get<ApiResponse<ChannelInfo[]>>(`/srs-sip/v1/devices/${deviceId}/channels`),

  // 添加 invite API
  inviteStream: (params: { device_id: string; channel_id: string; sub_stream: string }) =>
    api.post<ApiResponse<any>>('/srs-sip/v1/invite', params),

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
