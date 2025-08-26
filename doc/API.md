# SRS-SIP API 接口文档

## 概述

SRS-SIP 是一个基于 GB28181 协议的视频监控系统，提供设备管理、视频流控制、PTZ控制等功能。本文档描述了系统提供的 HTTP API 接口。

**基础信息：**
- 基础URL: `/srs-sip/v1`
- 内容类型: `application/json`
- 响应格式: JSON

## 通用响应格式

所有API接口都使用统一的响应格式：

```json
{
  "code": 0,
  "data": {}
}
```

- `code`: 状态码，0表示成功，其他值表示错误
- `data`: 响应数据，具体内容根据接口而定

## 1. 系统信息接口

### 1.1 获取API版本信息

**接口地址：** `GET /srs-sip`

**描述：** 获取API版本信息

**响应示例：**
```json
{
  "version": "v1"
}
```

### 1.2 获取所有API路由

**接口地址：** `GET /srs-sip/v1`

**描述：** 获取系统所有可用的API路由列表

**响应示例：**
```json
{
  "code": 0,
  "data": [
    {
      "method": "GET",
      "path": "/srs-sip/v1/devices"
    },
    {
      "method": "POST",
      "path": "/srs-sip/v1/invite"
    }
  ]
}
```

## 2. 设备管理接口

### 2.1 获取设备列表

**接口地址：** `GET /srs-sip/v1/devices`

**描述：** 获取所有已注册的设备列表

**响应示例：**
```json
{
  "code": 0,
  "data": [
    {
      "device_id": "34020000001320000001",
      "source_addr": "192.168.1.100:5060",
      "network_type": "UDP",
      "online": true,
      "heart_beat_interval": 60,
      "heart_beat_count": 3
    }
  ]
}
```

**响应字段说明：**
- `device_id`: 设备ID
- `source_addr`: 设备源地址
- `network_type`: 网络类型（UDP/TCP）
- `online`: 设备在线状态
- `heart_beat_interval`: 心跳间隔时间（秒）
- `heart_beat_count`: 心跳超时次数

### 2.2 获取设备通道列表

**接口地址：** `GET /srs-sip/v1/devices/{id}/channels`

**描述：** 获取指定设备的所有通道信息

**路径参数：**
- `id`: 设备ID

**响应示例：**
```json
{
  "code": 0,
  "data": [
    {
      "device_id": "34020000001320000002",
      "parent_id": "34020000001320000001",
      "name": "摄像头01",
      "manufacturer": "HIKVISION",
      "model": "DS-2CD2T47G1-L",
      "owner": "Camera Owner",
      "civil_code": "320100",
      "address": "测试地址",
      "port": 5060,
      "parental": 1,
      "safety_way": 0,
      "register_way": 1,
      "secrecy": 0,
      "ip_address": "192.168.1.100",
      "status": "ON",
      "longitude": 120.123456,
      "latitude": 30.123456,
      "info": {
        "ptz_type": 1,
        "resolution": "1920*1080",
        "download_speed": "4"
      },
      "ssrc": "1234567890"
    }
  ]
}
```

### 2.3 获取所有通道列表

**接口地址：** `GET /srs-sip/v1/channels`

**描述：** 获取系统中所有设备的通道信息

**响应格式：** 与2.2接口相同，返回所有设备的通道列表

## 3. 视频流控制接口

### 3.1 发起视频邀请

**接口地址：** `POST /srs-sip/v1/invite`

**描述：** 向指定通道发起视频流邀请

**请求参数：**
```json
{
  "device_id": "34020000001320000001",
  "channel_id": "34020000001320000002",
  "media_server_id": 1,
  "play_type": 0,
  "sub_stream": 0,
  "start_time": 1640995200,
  "end_time": 1640998800
}
```

**参数说明：**
- `device_id`: 设备ID
- `channel_id`: 通道ID
- `media_server_id`: 媒体服务器ID
- `play_type`: 播放类型（0:实时播放, 1:回放, 2:下载）
- `sub_stream`: 子码流（0:主码流, 1:子码流）
- `start_time`: 开始时间（Unix时间戳，回放时使用）
- `end_time`: 结束时间（Unix时间戳，回放时使用）

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "channel_id": "34020000001320000002",
    "url": "rtmp://192.168.1.200:1935/live/34020000001320000002"
  }
}
```

### 3.2 结束视频会话

**接口地址：** `POST /srs-sip/v1/bye`

**描述：** 结束指定的视频会话

**请求参数：**
```json
{
  "device_id": "34020000001320000001",
  "channel_id": "34020000001320000002",
  "url": "rtmp://192.168.1.200:1935/live/34020000001320000002"
}
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "msg": "success"
  }
}
```

### 3.3 暂停视频流

**接口地址：** `POST /srs-sip/v1/pause`

**描述：** 暂停指定的视频流

**请求参数：**
```json
{
  "device_id": "34020000001320000001",
  "channel_id": "34020000001320000002",
  "url": "rtmp://192.168.1.200:1935/live/34020000001320000002"
}
```

### 3.4 恢复视频流

**接口地址：** `POST /srs-sip/v1/resume`

**描述：** 恢复已暂停的视频流

**请求参数：** 与暂停接口相同

### 3.5 调整播放速度

**接口地址：** `POST /srs-sip/v1/speed`

**描述：** 调整视频播放速度（用于回放）

**请求参数：**
```json
{
  "device_id": "34020000001320000001",
  "channel_id": "34020000001320000002",
  "url": "rtmp://192.168.1.200:1935/live/34020000001320000002",
  "speed": 2.0
}
```

**参数说明：**
- `speed`: 播放速度倍数（0.5表示0.5倍速，2.0表示2倍速）

## 4. PTZ控制接口

### 4.1 PTZ控制

**接口地址：** `POST /srs-sip/v1/ptz`

**描述：** 控制摄像头的PTZ（云台）功能

**请求参数：**
```json
{
  "device_id": "34020000001320000001",
  "channel_id": "34020000001320000002",
  "ptz": "up",
  "speed": "5"
}
```

**参数说明：**
- `ptz`: PTZ控制命令
  - `up`: 向上
  - `down`: 向下
  - `left`: 向左
  - `right`: 向右
  - `zoom_in`: 放大
  - `zoom_out`: 缩小
  - `stop`: 停止
- `speed`: 控制速度（1-9）

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "msg": "success"
  }
}
```

## 5. 录像查询接口

### 5.1 查询录像记录

**接口地址：** `POST /srs-sip/v1/query-record`

**描述：** 查询指定设备通道的录像记录

**请求参数：**
```json
{
  "device_id": "34020000001320000001",
  "channel_id": "34020000001320000002",
  "start_time": 1640995200,
  "end_time": 1640998800
}
```

**参数说明：**
- `device_id`: 设备ID
- `channel_id`: 通道ID
- `start_time`: 查询开始时间（Unix时间戳）
- `end_time`: 查询结束时间（Unix时间戳）

**响应示例：**
```json
{
  "code": 0,
  "data": [
    {
      "device_id": "34020000001320000002",
      "name": "录像文件01",
      "file_path": "/record/20220101/001.mp4",
      "address": "192.168.1.100",
      "start_time": "2022-01-01T10:00:00Z",
      "end_time": "2022-01-01T11:00:00Z",
      "secrecy": 0,
      "type": "time"
    }
  ]
}
```

**响应字段说明：**
- `device_id`: 设备ID
- `name`: 录像文件名称
- `file_path`: 录像文件路径
- `address`: 录像文件地址
- `start_time`: 录像开始时间
- `end_time`: 录像结束时间
- `secrecy`: 保密属性（0:不保密, 1:保密）
- `type`: 录像类型

## 6. 媒体服务器管理接口

### 6.1 获取媒体服务器列表

**接口地址：** `GET /srs-sip/v1/media-servers`

**描述：** 获取所有已配置的媒体服务器列表

**响应示例：**
```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "name": "SRS服务器01",
      "type": "SRS",
      "ip": "192.168.1.200",
      "port": 1935,
      "username": "admin",
      "password": "123456",
      "secret": "secret_key",
      "is_default": 1,
      "created_at": "2022-01-01T10:00:00Z"
    }
  ]
}
```

**响应字段说明：**
- `id`: 媒体服务器ID
- `name`: 服务器名称
- `type`: 服务器类型（SRS、ZLM等）
- `ip`: 服务器IP地址
- `port`: 服务器端口
- `username`: 用户名
- `password`: 密码
- `secret`: 密钥
- `is_default`: 是否为默认服务器（1:是, 0:否）
- `created_at`: 创建时间

### 6.2 添加媒体服务器

**接口地址：** `POST /srs-sip/v1/media-servers`

**描述：** 添加新的媒体服务器配置

**请求参数：**
```json
{
  "name": "SRS服务器02",
  "type": "SRS",
  "ip": "192.168.1.201",
  "port": 1935,
  "username": "admin",
  "password": "123456",
  "secret": "secret_key",
  "is_default": 0
}
```

**参数说明：**
- `name`: 服务器名称（必填）
- `type`: 服务器类型（必填，支持SRS、ZLM等）
- `ip`: 服务器IP地址（必填）
- `port`: 服务器端口（必填）
- `username`: 用户名（可选）
- `password`: 密码（可选）
- `secret`: 密钥（可选）
- `is_default`: 是否设为默认服务器（可选，默认为0）

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "msg": "success"
  }
}
```

### 6.3 删除媒体服务器

**接口地址：** `DELETE /srs-sip/v1/media-servers/{id}`

**描述：** 删除指定的媒体服务器配置

**路径参数：**
- `id`: 媒体服务器ID

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "msg": "success"
  }
}
```

### 6.4 设置默认媒体服务器

**接口地址：** `POST /srs-sip/v1/media-servers/default/{id}`

**描述：** 将指定的媒体服务器设置为默认服务器

**路径参数：**
- `id`: 媒体服务器ID

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "msg": "success"
  }
}
```

## 7. 错误码说明

| 错误码 | HTTP状态码 | 说明 |
|--------|------------|------|
| 0 | 200 | 成功 |
| 400 | 400 | 请求参数错误 |
| 500 | 500 | 服务器内部错误 |

## 8. 使用示例

### 8.1 获取设备列表并播放视频

```bash
# 1. 获取设备列表
curl -X GET "http://localhost:8080/srs-sip/v1/devices"

# 2. 获取设备通道
curl -X GET "http://localhost:8080/srs-sip/v1/devices/34020000001320000001/channels"

# 3. 发起视频邀请
curl -X POST "http://localhost:8080/srs-sip/v1/invite" \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "34020000001320000001",
    "channel_id": "34020000001320000002",
    "media_server_id": 1,
    "play_type": 0,
    "sub_stream": 0
  }'
```

### 8.2 PTZ控制示例

```bash
# 控制摄像头向上移动
curl -X POST "http://localhost:8080/srs-sip/v1/ptz" \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "34020000001320000001",
    "channel_id": "34020000001320000002",
    "ptz": "up",
    "speed": "5"
  }'

# 停止PTZ控制
curl -X POST "http://localhost:8080/srs-sip/v1/ptz" \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "34020000001320000001",
    "channel_id": "34020000001320000002",
    "ptz": "stop",
    "speed": "0"
  }'
```

### 8.3 录像查询和回放示例

```bash
# 查询录像记录
curl -X POST "http://localhost:8080/srs-sip/v1/query-record" \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "34020000001320000001",
    "channel_id": "34020000001320000002",
    "start_time": 1640995200,
    "end_time": 1640998800
  }'

# 回放录像
curl -X POST "http://localhost:8080/srs-sip/v1/invite" \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "34020000001320000001",
    "channel_id": "34020000001320000002",
    "media_server_id": 1,
    "play_type": 1,
    "start_time": 1640995200,
    "end_time": 1640998800
  }'
```

## 9. 注意事项

1. **设备ID格式**：设备ID需要符合GB28181标准，通常为20位数字
2. **时间格式**：所有时间参数使用Unix时间戳格式
3. **媒体服务器**：在发起视频邀请前，需要确保已配置可用的媒体服务器
4. **网络连接**：确保SIP服务器与设备、媒体服务器之间网络连通
5. **权限控制**：部分接口可能需要认证，具体根据系统配置而定
6. **设备状态**：只有在线状态的设备才能进行视频邀请和PTZ控制
7. **并发限制**：同一通道同时只能有一个活跃的视频会话

## 10. 常见问题

### Q1: 设备离线怎么办？
A: 检查设备网络连接，确保设备能正常发送心跳包到SIP服务器。

### Q2: 视频邀请失败？
A: 检查媒体服务器状态，确保设备支持对应的编码格式。

### Q3: PTZ控制无响应？
A: 确认设备支持PTZ功能，检查PTZ控制命令格式是否正确。

## 11. 更新日志

- **v1.0.0** (2024-01-01): 初始版本，包含基础的设备管理、视频控制、PTZ控制功能
