import axios from 'axios'
import { ElMessage } from 'element-plus'
import type * as Types from './types'
import type { MediaServer } from '@/api/mediaserver/types'

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
  getMediaServers: () =>
    api.get<Types.ApiResponse<MediaServer[]>>('/srs-sip/v1/media-servers'),

  // 添加媒体服务器
  addMediaServer: (data: Omit<MediaServer, 'id' | 'status' | 'created_at'>) => 
    api.post<Types.ApiResponse<{ msg: string }>>('/srs-sip/v1/media-servers', data),

  // 删除媒体服务器
  deleteMediaServer: (id: number) =>
    api.delete<Types.ApiResponse<{ msg: string }>>(`/srs-sip/v1/media-servers/${id}`),

  // 设置默认媒体服务器
  setDefaultMediaServer: (id: number) =>
    api.post<Types.ApiResponse<{ msg: string }>>(`/srs-sip/v1/media-servers/default/${id}`),
}

// 设备相关 API
export const deviceApi = {
  // 获取设备列表
  getDevices: () => api.get<Types.ApiResponse<Types.Device[]>>('/srs-sip/v1/devices'),

  // 获取设备通道
  getDeviceChannels: (deviceId: string) =>
    api.get<Types.ApiResponse<Types.ChannelInfo[]>>(`/srs-sip/v1/devices/${deviceId}/channels`),

  // 添加 invite API
  invite: (params: Types.InviteRequest) =>
    api.post<Types.ApiResponse<Types.InviteResponse>>('/srs-sip/v1/invite', params),

  // 停止播放
  bye: (params: Types.ByeRequest) => api.post<Types.ApiResponse<any>>('/srs-sip/v1/bye', params),

  // 暂停播放
  pause: (params: Types.PauseRequest) => api.post<Types.ApiResponse<any>>('/srs-sip/v1/pause', params),

  // 恢复播放
  resume: (params: Types.ResumeRequest) => api.post<Types.ApiResponse<any>>('/srs-sip/v1/resume', params),

  // 设置播放速度
  speed: (params: Types.SpeedRequest) => api.post<Types.ApiResponse<any>>('/srs-sip/v1/speed', params),

  // 云台控制
  controlPTZ: (params: Types.PTZControlRequest) =>
    api.post<Types.ApiResponse<any>>('/srs-sip/v1/ptz', params),

  // 查询录像
  queryRecord: (params: Types.RecordInfoRequest) =>
    api.post<Types.ApiResponse<Types.RecordInfoResponse[]>>('/srs-sip/v1/query-record', params),
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
    const res = response.data as Types.ApiResponse<any>
    if (res.code !== 0) {
      ElMessage.error('请求失败')
      return Promise.reject(new Error('请求失败'))
    }
    response.data = res.data
    return response
  },
  (error) => {
    ElMessage.error(error.response?.data?.message || '网络错误')
    return Promise.reject(error)
  },
)

export default api
