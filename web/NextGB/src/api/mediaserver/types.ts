/**
 * 流媒体服务器类型枚举
 */
export enum MediaServerType {
    ZLM = 'zlm',    // ZLMediaKit
    SRS = 'srs',    // SRS
    CUSTOM = 'custom' // 自定义服务器
}

/**
 * 版本信息接口
 */
export interface VersionInfo {
    version: string;
    buildDate?: string;
    platform?: string;
}

/**
 * 媒体服务器基础接口
 */
export interface MediaServer {
    type: MediaServerType;
    getVersion(): Promise<VersionInfo>;
} 