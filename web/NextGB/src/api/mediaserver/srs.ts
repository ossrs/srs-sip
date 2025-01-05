import type { ClientInfo, StreamInfo, VersionInfo } from './types';
import { MediaServerType } from './types';
import { BaseMediaServer } from './base';
import axios from 'axios';

interface SRSVersionResponse {
    code: number;
    server: string;
    service: string;
    pid: string;
    data: {
        major: number;
        minor: number;
        revision: number;
        version: string;
    };
}

interface SRSClientsResponse {
    code: number;
    server: string;
    service: string;
    pid: string;
    clients: {
        id: string;
        vhost: string;
        stream: string;
        ip: string;
        pageUrl: string;
        swfUrl: string;
        tcUrl: string;
        url: string;
        name: string;
        type: string;
        publish: boolean;
        alive: number;
        send_bytes: number;
        recv_bytes: number;
        kbps: {
            recv_30s: number;
            send_30s: number;
        };
    }[];
}

interface SRSStreamResponse {
    code: number;
    server: string;
    service: string;
    pid: string;
    streams: {
        id: string;
        name: string;
        vhost: string;
        app: string;
        tcUrl: string;
        url: string;
        live_ms: number;
        clients: number;
        frames: number;
        send_bytes: number;
        recv_bytes: number;
        kbps: {
            recv_30s: number;
            send_30s: number;
        };
        publish: {
            active: boolean;
            cid: string;
        };
        video?: {
            codec: string;
            profile: string;
            level: string;
            width: number;
            height: number;
        };
        audio?: {
            codec: string;
            sample_rate: number;
            channel: number;
            profile: string;
        };
    }[];
}

export class SRSServer extends BaseMediaServer {
    private baseUrl: string;

    constructor(host: string, port: number) {
        super(MediaServerType.SRS);
        this.baseUrl = `http://${host}:${port}`;
    }

    async getVersion(): Promise<VersionInfo> {
        try {
            const response = await axios.get<SRSVersionResponse>(`${this.baseUrl}/api/v1/versions`);
            
            return {
                version: response.data.data.version,
                buildDate: undefined, // SRS API 没有提供构建日期
                platform: `SRS Server: ${response.data.server}` // 使用 server 标识作为平台信息
            };
        } catch (error) {
            throw new Error(`Failed to get SRS version: ${error}`);
        }
    }

    async getStreamInfo(): Promise<StreamInfo[]> {
        try {
            const response = await axios.get<SRSStreamResponse>(`${this.baseUrl}/api/v1/streams/`);
            
            return response.data.streams.map(stream => ({
                id: stream.id,
                name: stream.name,
                vhost: stream.vhost,
                url: stream.tcUrl,
                clients: stream.clients - 1,
                active: stream.publish.active,
                send_bytes: stream.send_bytes,
                recv_bytes: stream.recv_bytes,
                video: stream.video ? {
                    codec: stream.video.codec,
                    width: stream.video.width,
                    height: stream.video.height,
                    fps: 0,  // SRS API 没有直接提供 fps 信息
                } : undefined,
                audio: stream.audio ? {
                    codec: stream.audio.codec,
                    sampleRate: stream.audio.sample_rate,
                    channels: stream.audio.channel,
                } : undefined
            }));
        } catch (error) {
            throw new Error(`Failed to get SRS streams info: ${error}`);
        }
    }

    async getClientInfo(): Promise<ClientInfo[]> {
        try {
            const response = await axios.get<SRSClientsResponse>(`${this.baseUrl}/api/v1/clients/`);
            return response.data.clients
                .filter(client => !client.publish)
                .map(client => {
                    console.log('Client alive value:', client.alive, typeof client.alive);
                    return {
                        id: client.id,
                        vhost: client.vhost,
                        stream: client.stream,
                        ip: client.ip,
                        url: client.url,
                        alive: Math.round(client.alive * 1000), // 转换为毫秒并四舍五入
                        type: client.type
                    };
                });
        } catch (error) {
            throw new Error(`Failed to get SRS clients info: ${error}`);
        }
    }
} 