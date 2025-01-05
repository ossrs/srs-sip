import type { MediaServer, VersionInfo } from './types';
import { MediaServerType } from './types';

/**
 * 媒体服务器基础实现类
 */
export abstract class BaseMediaServer implements MediaServer {
    type: MediaServerType;

    constructor(type: MediaServerType) {
        this.type = type;
    }

    abstract getVersion(): Promise<VersionInfo>;
} 