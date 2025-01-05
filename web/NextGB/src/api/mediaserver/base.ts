import type { ClientInfo, StreamInfo, VersionInfo } from './types';
import { MediaServerType } from './types';

/**
 * 媒体服务器基础接口
 */
export interface MediaServer {
    type: MediaServerType;
    getVersion(): Promise<VersionInfo>;
    getStreamInfo(): Promise<StreamInfo[]>;
    getClientInfo(): Promise<ClientInfo[]>;
}

/**
 * 媒体服务器基础实现类
 */
export abstract class BaseMediaServer implements MediaServer {
    type: MediaServerType;

    constructor(type: MediaServerType) {
        this.type = type;
    }

    abstract getVersion(): Promise<VersionInfo>;
    abstract getStreamInfo(): Promise<StreamInfo[]>;
    abstract getClientInfo(): Promise<ClientInfo[]>;
} 