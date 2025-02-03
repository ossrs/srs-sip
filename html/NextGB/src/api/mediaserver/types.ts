/**
 * 流媒体服务器类型枚举
 */
export enum MediaServerType {
  ZLM = 'zlm', // ZLMediaKit
  SRS = 'srs', // SRS
  CUSTOM = 'custom', // 自定义服务器
}

export interface RtcPlayer {
  pc: RTCPeerConnection
  play(url: string): Promise<void>
  close(): Promise<void>
  ontrack: ((event: RTCTrackEvent) => void) | null
}

export interface MediaServer {
  id: number
  name: string
  ip: string
  port: number
  type: string
  username?: string
  password?: string
  secret?: string
  status: number
  created_at: string
  isDefault?: number
}

/**
 * 版本信息接口
 */
export interface VersionInfo {
  version: string
  buildDate?: string
}

/**
 * 视频编码信息
 */
export interface VideoCodecInfo {
  codec: string // 视频编码格式，如 H264, H265
  width: number // 视频宽度
  height: number // 视频高度
  fps: number // 帧率
  bitrate?: number // 比特率 (kbps)
}

/**
 * 音频编码信息
 */
export interface AudioCodecInfo {
  codec: string // 音频编码格式，如 AAC, G711
  sampleRate: number // 采样率
  channels: number // 声道数
  bitrate?: number // 比特率 (kbps)
}

export interface StreamInfo {
  id: string
  name: string // 流名称
  vhost: string // 虚拟主机
  url: string // 流地址
  clients: number // 客户端连接数
  active: boolean // 是否活跃
  video?: VideoCodecInfo // 视频编码信息
  audio?: AudioCodecInfo // 音频编码信息
  send_bytes?: number // 已传输字节数
  recv_bytes?: number // 已接收字节数
}

export interface ClientInfo {
  id: string
  vhost: string
  stream: string
  ip: string
  url: string
  alive: number
  type: string
}
