import type { VersionInfo } from './types';
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
} 