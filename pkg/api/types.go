package api

// InviteRequest 定义了invite接口的请求参数
type InviteRequest struct {
	MediaServerAddr string `json:"media_server_addr"`
	DeviceID        string `json:"device_id"`
	ChannelID       string `json:"channel_id"`
	SubStream       int    `json:"sub_stream"`
	PlayType        int    `json:"play_type"` // 0: live, 1: playback, 2: download
	StartTime       int64  `json:"start_time"`
	EndTime         int64  `json:"end_time"`
}

// InviteResponse 定义了invite接口的响应参数
type InviteResponse struct {
	ChannelID string `json:"channel_id"`
	URL       string `json:"url"`
}

// PTZControlRequest 定义了PTZ控制接口的请求参数
type PTZControlRequest struct {
	DeviceID  string `json:"device_id"`
	ChannelID string `json:"channel_id"`
	PTZ       string `json:"ptz"`
	Speed     string `json:"speed"`
}

// QueryRecordRequest 定义了录像查询接口的请求参数
type QueryRecordRequest struct {
	DeviceID  string `json:"device_id"`
	ChannelID string `json:"channel_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

// MediaServerRequest 定义了媒体服务器添加接口的请求参数
type MediaServerRequest struct {
	Name      string `json:"name"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Type      string `json:"type"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	IsDefault int    `json:"is_default"`
}

// CommonResponse 定义了通用的响应格式
type CommonResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
