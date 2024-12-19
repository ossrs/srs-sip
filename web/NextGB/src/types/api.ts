// 设备相关类型
export interface Device {
  device_id: string
  source_addr: string
  network_type: string
  status: 'online' | 'offline'
  name: string
}

// API 响应类型
export interface ApiResponse<T> {
  code: number
  data: T
  message: string
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
