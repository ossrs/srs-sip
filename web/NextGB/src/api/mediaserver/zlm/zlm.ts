import type { ClientInfo, StreamInfo, VersionInfo, MediaServer, RtcPlayer } from '@/api/mediaserver/types'
import { MediaServerType } from '@/api/mediaserver/types'
import { BaseMediaServer } from '@/api/mediaserver/base'
import axios from 'axios'

// {
//   "code": 0,
//   "data": {
//       "branchName": "master",
//       "buildTime": "2023-04-19T10:34:34",
//       "commitHash": "f143898"
//    }
// }
interface ZLMVersionResponse {
  code: number
  data: {
    branchName: string
    buildTime: string
    commitHash: string
  }
}

// {
//   "code" : 0,
//   "data" : [
//   {
//      "app" : "live",  # 应用名
//      "readerCount" : 0, # 本协议观看人数
//      "totalReaderCount" : 0, # 观看总人数，包括hls/rtsp/rtmp/http-flv/ws-flv
//      "schema" : "rtsp", # 协议
//      "stream" : "obs", # 流id
//      "originSock": {  # 客户端和服务器网络信息，可能为null类型
//             "identifier": "140241931428384",
//             "local_ip": "127.0.0.1",
//             "local_port": 1935,
//             "peer_ip": "127.0.0.1",
//             "peer_port": 50097
//         },
//      "originType": 1, # 产生源类型，包括 unknown = 0,rtmp_push=1,rtsp_push=2,rtp_push=3,pull=4,ffmpeg_pull=5,mp4_vod=6,device_chn=7
//      "originTypeStr": "MediaOriginType::rtmp_push",
//      "originUrl": "rtmp://127.0.0.1:1935/live/hks2", #产生源的url
//      "createStamp": 1602205811, #GMT unix系统时间戳，单位秒
//      "aliveSecond": 100, #存活时间，单位秒
//      "bytesSpeed": 12345, #数据产生速度，单位byte/s
//      "tracks" : [    # 音视频轨道
//         {
//            "channels" : 1, # 音频通道数
//            "codec_id" : 2, # H264 = 0, H265 = 1, AAC = 2, G711A = 3, G711U = 4
//            "codec_id_name" : "CodecAAC", # 编码类型名称 
//            "codec_type" : 1, # Video = 0, Audio = 1
//            "ready" : true, # 轨道是否准备就绪
//            "frames" : 1119, #累计接收帧数
//            "sample_bit" : 16, # 音频采样位数
//            "sample_rate" : 8000 # 音频采样率
//         },
//         {
//            "codec_id" : 0, # H264 = 0, H265 = 1, AAC = 2, G711A = 3, G711U = 4
//            "codec_id_name" : "CodecH264", # 编码类型名称  
//            "codec_type" : 0, # Video = 0, Audio = 1
//            "fps" : 59,  # 视频fps
//            "frames" : 1119, #累计接收帧数，不包含sei/aud/sps/pps等不能解码的帧
//            "gop_interval_ms" : 1993, #gop间隔时间，单位毫秒
//            "gop_size" : 60, #gop大小，单位帧数
//            "key_frames" : 21, #累计接收关键帧数
//            "height" : 720, # 视频高
//            "ready" : true,  # 轨道是否准备就绪
//            "width" : 1280 # 视频宽
//         }
//      ],
//      "vhost" : "__defaultVhost__" # 虚拟主机名
//    }
//   ]
// }

interface ZLMTrackInfo {
  channels?: number        // 音频通道数
  codec_id: number        // 编码器ID
  codec_id_name: string   // 编码器名称
  codec_type: number      // 编码类型 (0: Video, 1: Audio)
  ready: boolean          // 轨道是否就绪
  frames: number          // 累计接收帧数
  sample_bit?: number     // 音频采样位数
  sample_rate?: number    // 音频采样率
  // 视频特有属性
  fps?: number           // 视频帧率
  width?: number         // 视频宽度
  height?: number        // 视频高度
  gop_interval_ms?: number // GOP间隔时间
  gop_size?: number      // GOP大小
  key_frames?: number    // 关键帧数
  bytesSpeed?: number    // 数据速率
}

interface ZLMStreamInfo {
  app: string
  readerCount: number
  totalReaderCount: number
  schema: string
  stream: string
  originSock: {
    identifier: string
    local_ip: string
    local_port: number
    peer_ip: string
    peer_port: number
  }
  originType: number
  originTypeStr: string
  originUrl: string
  createStamp: number
  aliveSecond: number
  bytesSpeed: number
  tracks: ZLMTrackInfo[]
  vhost: string
}

interface ZLMStreamResponse {
  code: number
  data: ZLMStreamInfo[]
}

// {
//   "code": 0,
//   "data": [
//       {
//           "identifier": "3-309",
//           "local_ip": "::",
//           "local_port": 8000,
//           "peer_ip": "172.18.190.159",
//           "peer_port": 52996,
//           "typeid": "mediakit::WebRtcSession"
//       }
//   ]
// }

interface ZLMClientInfo {
  identifier: string
  local_ip: string
  local_port: number
  peer_ip: string
  peer_port: number
  typeid: string
}

interface ZLMClientInfoResponse {
  code: number
  data: ZLMClientInfo[]
}

interface ZLMRtcResponse {
  code: number
  id: string
  sdp: string
  type: string
}

export class ZLMServer extends BaseMediaServer {
  private baseUrl: string
  private secret: string

  constructor(host: string, port: number, secret: string = '') {
    super(MediaServerType.ZLM)
    this.baseUrl = `http://${host}:${port}`
    this.secret = secret
  }

  async getVersion(): Promise<VersionInfo> {
    try {
      const response = await axios.get<ZLMVersionResponse>(`${this.baseUrl}/index/api/version${this.secret ? '?secret=' + this.secret : ''}`)
      return {
        version: response.data.data.buildTime,
        buildDate: response.data.data.buildTime,
      }
    } catch (error) {
      throw new Error(`Failed to get ZLM version: ${error}`)
    }
  }

  // /index/api/getMediaList?schema=rtsp&secret=
  async getStreamInfo(): Promise<StreamInfo[]> {
    try {
      const response = await axios.get<ZLMStreamResponse>(`${this.baseUrl}/index/api/getMediaList?schema=rtsp&secret=${this.secret}`)
      return response.data.data.map((stream) => {
        const videoTrack = stream.tracks.find((track) => track.codec_type === 0)
        const audioTrack = stream.tracks.find((track) => track.codec_type === 1)

        return {
          id: stream.stream,
          name: stream.stream,
          vhost: stream.vhost,
          url: stream.originUrl,
          clients: stream.readerCount,
          active: stream.aliveSecond > 0,
          video: videoTrack?.codec_id_name === 'CodecH264' ? {
            codec: videoTrack.codec_id_name,
            width: videoTrack.width ?? 0,
            height: videoTrack.height ?? 0,
            fps: videoTrack.fps ?? 0,
            bitrate: videoTrack.bytesSpeed ?? 0,
          } : undefined,
          audio: audioTrack?.codec_id_name === 'CodecAAC' ? {
            codec: audioTrack.codec_id_name,
            sampleRate: audioTrack.sample_rate ?? 0,
            channels: audioTrack.channels ?? 0,
            bitrate: audioTrack.bytesSpeed ?? 0,
          } : undefined,
        }
      })
    } catch (error) {
      throw new Error(`Failed to get ZLM streams info: ${error}`)
    }
  }

  // /index/api/getMediaPlayerList?secret=035c73f7-bb6b-4889-a715-d9eb2d1925cc&schema=rtsp&vhost=defaultVhost&app=live&stream=test
  async getClientInfo(): Promise<ClientInfo[]> {
    try {
      const response = await axios.get<ZLMClientInfoResponse>(
        `${this.baseUrl}/index/api/getMediaPlayerList?secret=${this.secret}`
      )
      
      return response.data.data.map((client) => ({
        id: client.identifier,
        vhost: '__defaultVhost__', // ZLM默认虚拟主机
        stream: client.identifier.split('-')[1] || '', // 假设identifier格式为 "type-streamid"
        ip: client.peer_ip,
        url: `${client.local_ip}:${client.local_port}`,
        alive: 1, // 在线状态
        type: client.typeid.replace('mediakit::', '') // 移除 mediakit:: 前缀
      }))
    } catch (error) {
      throw new Error(`Failed to get ZLM clients info: ${error}`)
    }
  }

  createRtcPlayer(): RtcPlayer {
    const self = {
      pc: new RTCPeerConnection({
        iceServers: [],
      }),

      async play(url: string): Promise<void> {
        this.pc.addTransceiver('audio', { direction: 'recvonly' })
        this.pc.addTransceiver('video', { direction: 'recvonly' })

        const offer = await this.pc.createOffer()
        await this.pc.setLocalDescription(offer)

        const response = await axios.post<ZLMRtcResponse>(
          url,
          offer.sdp,
          {
            headers: {
              'Content-Type': 'text/plain;charset=utf-8'
            }
          }
        )

        if (response.data.code !== 0) {
          throw new Error('创建WebRTC播放器失败')
        }

        await this.pc.setRemoteDescription(
          new RTCSessionDescription({ type: 'answer', sdp: response.data.sdp })
        )
      },

      async close(): Promise<void> {
        this.pc.close()
      },

      ontrack: null as ((event: RTCTrackEvent) => void) | null,
    }

    self.pc.ontrack = (event: RTCTrackEvent) => {
      if (self.ontrack) {
        self.ontrack(event)
      }
    }

    return self
  }
}
