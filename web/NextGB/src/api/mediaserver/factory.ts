import type { MediaServer } from './types'
import { MediaServerType } from './types'
import type { BaseMediaServer } from './base'
import { SRSServer } from './srs/srs'
import { ZLMServer } from './zlm/zlm'

/**
 * 创建媒体服务器实例的工厂函数
 */
export const createMediaServer = (config: MediaServer): BaseMediaServer => {
  // 统一转换为小写进行比较
  const serverType = config.type.toLowerCase()
  
  switch (serverType) {
    case MediaServerType.SRS:
      return new SRSServer(config.ip, config.port)
    case MediaServerType.ZLM:
      return new ZLMServer(config.ip, config.port, config.secret)
    default:
      throw new Error(`Unsupported media server type: ${config.type}`)
  }
} 