import type { ClientInfo, StreamInfo, VersionInfo, RtcPlayer } from '@/api/mediaserver/types'
import { MediaServerType } from '@/api/mediaserver/types'
import { BaseMediaServer } from '@/api/mediaserver/base'
import axios from 'axios'

interface SRSVersionResponse {
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

interface SRSClientsResponse {
  code: number
  server: string
  service: string
  pid: string
  clients: {
    id: string
    vhost: string
    stream: string
    ip: string
    pageUrl: string
    swfUrl: string
    tcUrl: string
    url: string
    name: string
    type: string
    publish: boolean
    alive: number
    send_bytes: number
    recv_bytes: number
    kbps: {
      recv_30s: number
      send_30s: number
    }
  }[]
}

interface SRSStreamResponse {
  code: number
  server: string
  service: string
  pid: string
  streams: {
    id: string
    name: string
    vhost: string
    app: string
    tcUrl: string
    url: string
    live_ms: number
    clients: number
    frames: number
    send_bytes: number
    recv_bytes: number
    kbps: {
      recv_30s: number
      send_30s: number
    }
    publish: {
      active: boolean
      cid: string
    }
    video?: {
      codec: string
      profile: string
      level: string
      width: number
      height: number
    }
    audio?: {
      codec: string
      sample_rate: number
      channel: number
      profile: string
    }
  }[]
}

interface UserQuery {
  [key: string]: string | undefined
  schema?: string
  play?: string
}

interface ParsedUrl {
  url: string
  schema: string
  server: string
  port: number
  vhost: string
  app: string
  stream: string
  user_query: UserQuery
}



export class SRSServer extends BaseMediaServer {
  private baseUrl: string

  constructor(host: string, port: number) {
    super(MediaServerType.SRS)
    this.baseUrl = `http://${host}:${port}`
  }

  async getVersion(): Promise<VersionInfo> {
    try {
      const response = await axios.get<SRSVersionResponse>(`${this.baseUrl}/api/v1/versions`)

      return {
        version: response.data.data.version,
        buildDate: undefined, // SRS API 没有提供构建日期
      }
    } catch (error) {
      throw new Error(`Failed to get SRS version: ${error}`)
    }
  }

  async getStreamInfo(): Promise<StreamInfo[]> {
    try {
      const response = await axios.get<SRSStreamResponse>(`${this.baseUrl}/api/v1/streams/`)

      return response.data.streams.map((stream) => ({
        id: stream.id,
        name: stream.name,
        vhost: stream.vhost,
        url: stream.tcUrl,
        clients: stream.clients - 1,
        active: stream.publish.active,
        send_bytes: stream.send_bytes,
        recv_bytes: stream.recv_bytes,
        video: stream.video
          ? {
              codec: stream.video.codec,
              width: stream.video.width,
              height: stream.video.height,
              fps: 0, // SRS API 没有直接提供 fps 信息
            }
          : undefined,
        audio: stream.audio
          ? {
              codec: stream.audio.codec,
              sampleRate: stream.audio.sample_rate,
              channels: stream.audio.channel,
            }
          : undefined,
      }))
    } catch (error) {
      throw new Error(`Failed to get SRS streams info: ${error}`)
    }
  }

  async getClientInfo(params?: { stream_id?: string }): Promise<ClientInfo[]> {
    try {
      const response = await axios.get<SRSClientsResponse>(`${this.baseUrl}/api/v1/clients/`)
      let clients = response.data.clients.filter((client) => !client.publish)
      
      // 如果指定了 stream_id，则过滤出对应的流
      if (params?.stream_id) {
        clients = clients.filter(client => client.stream === params.stream_id)
      }

      return clients.map((client) => {
        console.log('Client alive value:', client.alive, typeof client.alive)
        return {
          id: client.id,
          vhost: client.vhost,
          stream: client.stream,
          ip: client.ip,
          url: client.url,
          alive: Math.round(client.alive * 1000), // 转换为毫秒并四舍五入
          type: client.type,
        }
      })
    } catch (error) {
      throw new Error(`Failed to get SRS clients info: ${error}`)
    }
  }

  async kickClient(clientId: string) {
    const response = await axios.post(`${this.baseUrl}/api/v1/clients/${clientId}/kick`)
    return response.data
  }

  createRtcPlayer(): RtcPlayer {
    const self = {
      pc: new RTCPeerConnection({
        iceServers: [],
      }),

      async play(url: string) {
        const conf = this.__internal.prepareUrl(url)
        this.pc.addTransceiver('audio', { direction: 'recvonly' })
        this.pc.addTransceiver('video', { direction: 'recvonly' })

        const offer = await this.pc.createOffer()
        await this.pc.setLocalDescription(offer)

        const session = await fetch(conf.apiUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            api: conf.apiUrl,
            streamurl: conf.streamUrl,
            clientip: null,
            sdp: offer.sdp,
          }),
        }).then((res) => res.json())

        if (session.code) {
          throw session
        }

        await this.pc.setRemoteDescription(
          new RTCSessionDescription({ type: 'answer', sdp: session.sdp }),
        )
        return session
      },

      async close() {
        this.pc.close()
      },

      ontrack: null as ((event: RTCTrackEvent) => void) | null,

      __internal: {
        defaultPath: '/rtc/v1/play/',

        prepareUrl(webrtcUrl: string) {
          const urlObject = this.parse(webrtcUrl) as ParsedUrl
          const schema = urlObject.user_query.schema
            ? urlObject.user_query.schema + ':'
            : window.location.protocol

          let port = urlObject.port || 1985
          if (schema === 'https:') {
            port = urlObject.port || 443
          }

          let api = urlObject.user_query.play || this.defaultPath
          if (api.lastIndexOf('/') !== api.length - 1) {
            api += '/'
          }

          let apiUrl = schema + '//' + urlObject.server + ':' + port + api
          for (const key in urlObject.user_query) {
            if (key !== 'api' && key !== 'play') {
              apiUrl += '&' + key + '=' + urlObject.user_query[key]
            }
          }
          apiUrl = apiUrl.replace(api + '&', api + '?')

          return {
            apiUrl,
            streamUrl: urlObject.url,
            schema,
            urlObject,
            port,
          }
        },

        parse(url: string): ParsedUrl {
          const a = document.createElement('a')
          a.href = url
            .replace('rtmp://', 'http://')
            .replace('webrtc://', 'http://')
            .replace('rtc://', 'http://')

          let vhost = a.hostname
          let app = a.pathname.substring(1, a.pathname.lastIndexOf('/'))
          const stream = a.pathname.slice(a.pathname.lastIndexOf('/') + 1)

          app = app.replace('...vhost...', '?vhost=')
          if (app.indexOf('?') >= 0) {
            const params = app.slice(app.indexOf('?'))
            app = app.slice(0, app.indexOf('?'))

            if (params.indexOf('vhost=') > 0) {
              vhost = params.slice(params.indexOf('vhost=') + 'vhost='.length)
              if (vhost.indexOf('&') > 0) {
                vhost = vhost.slice(0, vhost.indexOf('&'))
              }
            }
          }

          if (a.hostname === vhost) {
            const re = /^(\d+)\.(\d+)\.(\d+)\.(\d+)$/
            if (re.test(a.hostname)) {
              vhost = '__defaultVhost__'
            }
          }

          let schema = 'rtmp'
          if (url.indexOf('://') > 0) {
            schema = url.slice(0, url.indexOf('://'))
          }

          let port = parseInt(a.port)
          if (!port) {
            if (schema === 'http') {
              port = 80
            } else if (schema === 'https') {
              port = 443
            } else if (schema === 'rtmp') {
              port = 1935
            }
          }

          const ret: ParsedUrl = {
            url,
            schema,
            server: a.hostname,
            port,
            vhost,
            app,
            stream,
            user_query: {},
          }

          this.fill_query(a.search, ret)
          return ret
        },

        fill_query(query_string: string, obj: ParsedUrl) {
          if (query_string.length === 0) {
            return
          }

          if (query_string.indexOf('?') >= 0) {
            query_string = query_string.split('?')[1]
          }

          const queries = query_string.split('&')
          for (const elem of queries) {
            const query = elem.split('=')
            obj.user_query[query[0]] = query[1]
          }

          if (obj.user_query.domain) {
            obj.vhost = obj.user_query.domain
          }
        },
      },
    }

    self.pc.ontrack = (event: RTCTrackEvent) => {
      if (self.ontrack) {
        self.ontrack(event)
      }
    }

    return self
  }
}
