package models

type BaseRequest struct {
	DeviceID  string `json:"device_id"`
	ChannelID string `json:"channel_id"`
}

type InviteRequest struct {
	BaseRequest
	MediaServerId int   `json:"media_server_id"`
	PlayType      int   `json:"play_type"` // 0: live, 1: playback, 2: download
	SubStream     int   `json:"sub_stream"`
	StartTime     int64 `json:"start_time"`
	EndTime       int64 `json:"end_time"`
}

type InviteResponse struct {
	ChannelID string `json:"channel_id"`
	URL       string `json:"url"`
}

type SessionRequest struct {
	BaseRequest
	URL string `json:"url"`
}

type ByeRequest struct {
	SessionRequest
}

type PauseRequest struct {
	SessionRequest
}

type ResumeRequest struct {
	SessionRequest
}

type SpeedRequest struct {
	SessionRequest
	Speed float32 `json:"speed"`
}

type PTZControlRequest struct {
	BaseRequest
	PTZ   string `json:"ptz"`
	Speed string `json:"speed"`
}

type QueryRecordRequest struct {
	BaseRequest
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type MediaServer struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Secret    string `json:"secret"`
	IsDefault int    `json:"is_default"`
}

type MediaServerRequest struct {
	MediaServer
}

type MediaServerResponse struct {
	MediaServer
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
}

type CommonResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
