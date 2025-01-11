// API 响应类型
export interface ApiResponse<T = any> {
  code: number
  data: T
}

export interface InviteRequest {
  media_server_id: number
  device_id: string
  channel_id: string
  sub_stream: number
  play_type: number
  start_time: number
  end_time: number
}

export interface InviteResponse {
  url: string
}

export interface ByeRequest {
  device_id: string
  channel_id: string
  url: string
}

export interface PauseRequest {
  device_id: string
  channel_id: string
  url: string
}

export interface ResumeRequest {
  device_id: string
  channel_id: string
  url: string
}

export interface SpeedRequest {
  device_id: string
  channel_id: string
  url: string
  speed: number
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

export interface RecordInfoRequest {
  device_id: string
  channel_id: string
  start_time: number
  end_time: number
}

export interface RecordInfoResponse {
  device_id: string
  name: string
  file_path: string
  address: string
  start_time: string
  end_time: string
  secrecy: number
  type: string
}

export interface Device {
  device_id: string
  source_addr: string
  network_type: string
  status: number
  name: string
}

export interface PTZControlRequest {
  device_id: string
  channel_id: string
  ptz: string
  speed: string
}

// 媒体服务器类型

