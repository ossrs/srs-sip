import type { ClientInfo, StreamInfo, VersionInfo, MediaServer } from '@/api/mediaserver/types'
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

  async getStreamInfo(): Promise<StreamInfo[]> {
    try {
      // TODO: 实现获取ZLM流信息的逻辑
      throw new Error('Method not implemented.')
    } catch (error) {
      throw new Error(`Failed to get ZLM streams info: ${error}`)
    }
  }

  async getClientInfo(): Promise<ClientInfo[]> {
    try {
      // TODO: 实现获取ZLM客户端信息的逻辑
      throw new Error('Method not implemented.')
    } catch (error) {
      throw new Error(`Failed to get ZLM clients info: ${error}`)
    }
  }
}