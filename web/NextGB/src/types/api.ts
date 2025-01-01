// 设备相关类型
export interface Device {
  device_id: string
  source_addr: string
  network_type: string
  status: 'online' | 'offline'
  name: string
}

// API 响应类型
export interface ApiResponse<T = any> {
  code: number
  data: T
}

// 通道状态类型
export type ChannelStatus = 'ON' | 'OFF'

// 通道信息类型
export interface ChannelInfo {
  device_id: string
  parent_id: string
  name: string
  manufacturer: string
  model: string
  owner: string
  civil_code: string
  address: string
  port: number
  parental: number
  safety_way: number
  register_way: number
  secrecy: number
  ip_address: string
  status: ChannelStatus
  longitude: number
  latitude: number
  info: {
    ptz_type: number
    resolution: string
    download_speed: string
  }
  ssrc: string
}

// 媒体服务器类型
export interface MediaServer {
  id: number
  name: string
  ip: string
  port: number
  type: string
  username?: string
  password?: string
  status: number
  created_at: string
  isDefault?: number
}

export interface PTZControlRequest {
  device_id: string
  channel_id: string
  ptz: string
  speed: string
}
