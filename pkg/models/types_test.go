package models

import (
	"encoding/json"
	"testing"
)

func TestBaseRequest(t *testing.T) {
	req := BaseRequest{
		DeviceID:  "34020000001320000001",
		ChannelID: "34020000001320000002",
	}

	if req.DeviceID != "34020000001320000001" {
		t.Errorf("Expected DeviceID '34020000001320000001', got '%s'", req.DeviceID)
	}
	if req.ChannelID != "34020000001320000002" {
		t.Errorf("Expected ChannelID '34020000001320000002', got '%s'", req.ChannelID)
	}
}

func TestInviteRequest(t *testing.T) {
	req := InviteRequest{
		BaseRequest: BaseRequest{
			DeviceID:  "device123",
			ChannelID: "channel123",
		},
		MediaServerId: 1,
		PlayType:      0,
		SubStream:     0,
		StartTime:     1234567890,
		EndTime:       1234567900,
	}

	if req.DeviceID != "device123" {
		t.Errorf("Expected DeviceID 'device123', got '%s'", req.DeviceID)
	}
	if req.MediaServerId != 1 {
		t.Errorf("Expected MediaServerId 1, got %d", req.MediaServerId)
	}
	if req.PlayType != 0 {
		t.Errorf("Expected PlayType 0, got %d", req.PlayType)
	}
}

func TestInviteRequestJSON(t *testing.T) {
	jsonStr := `{
		"device_id": "device123",
		"channel_id": "channel123",
		"media_server_id": 1,
		"play_type": 1,
		"sub_stream": 0,
		"start_time": 1234567890,
		"end_time": 1234567900
	}`

	var req InviteRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if req.DeviceID != "device123" {
		t.Errorf("Expected DeviceID 'device123', got '%s'", req.DeviceID)
	}
	if req.PlayType != 1 {
		t.Errorf("Expected PlayType 1, got %d", req.PlayType)
	}
}

func TestInviteResponse(t *testing.T) {
	resp := InviteResponse{
		ChannelID: "channel123",
		URL:       "webrtc://example.com/live/stream",
	}

	if resp.ChannelID != "channel123" {
		t.Errorf("Expected ChannelID 'channel123', got '%s'", resp.ChannelID)
	}
	if resp.URL != "webrtc://example.com/live/stream" {
		t.Errorf("Expected URL 'webrtc://example.com/live/stream', got '%s'", resp.URL)
	}
}

func TestPTZControlRequest(t *testing.T) {
	req := PTZControlRequest{
		BaseRequest: BaseRequest{
			DeviceID:  "device123",
			ChannelID: "channel123",
		},
		PTZ:   "up",
		Speed: "5",
	}

	if req.PTZ != "up" {
		t.Errorf("Expected PTZ 'up', got '%s'", req.PTZ)
	}
	if req.Speed != "5" {
		t.Errorf("Expected Speed '5', got '%s'", req.Speed)
	}
}

func TestQueryRecordRequest(t *testing.T) {
	req := QueryRecordRequest{
		BaseRequest: BaseRequest{
			DeviceID:  "device123",
			ChannelID: "channel123",
		},
		StartTime: 1234567890,
		EndTime:   1234567900,
	}

	if req.StartTime != 1234567890 {
		t.Errorf("Expected StartTime 1234567890, got %d", req.StartTime)
	}
	if req.EndTime != 1234567900 {
		t.Errorf("Expected EndTime 1234567900, got %d", req.EndTime)
	}
}

func TestMediaServer(t *testing.T) {
	ms := MediaServer{
		Name:      "SRS Server",
		Type:      "SRS",
		IP:        "192.168.1.100",
		Port:      1985,
		Username:  "admin",
		Password:  "password",
		Secret:    "secret",
		IsDefault: 1,
	}

	if ms.Name != "SRS Server" {
		t.Errorf("Expected Name 'SRS Server', got '%s'", ms.Name)
	}
	if ms.Type != "SRS" {
		t.Errorf("Expected Type 'SRS', got '%s'", ms.Type)
	}
	if ms.Port != 1985 {
		t.Errorf("Expected Port 1985, got %d", ms.Port)
	}
	if ms.IsDefault != 1 {
		t.Errorf("Expected IsDefault 1, got %d", ms.IsDefault)
	}
}

func TestMediaServerResponse(t *testing.T) {
	resp := MediaServerResponse{
		MediaServer: MediaServer{
			Name: "Test Server",
			Type: "ZLM",
			IP:   "10.0.0.1",
			Port: 8080,
		},
		ID:        1,
		CreatedAt: "2024-01-01 12:00:00",
	}

	if resp.ID != 1 {
		t.Errorf("Expected ID 1, got %d", resp.ID)
	}
	if resp.CreatedAt != "2024-01-01 12:00:00" {
		t.Errorf("Expected CreatedAt '2024-01-01 12:00:00', got '%s'", resp.CreatedAt)
	}
}

func TestCommonResponse(t *testing.T) {
	resp := CommonResponse{
		Code: 0,
		Data: map[string]string{"key": "value"},
	}

	if resp.Code != 0 {
		t.Errorf("Expected Code 0, got %d", resp.Code)
	}

	// 测试 JSON 序列化
	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var decoded CommonResponse
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.Code != 0 {
		t.Errorf("Expected decoded Code 0, got %d", decoded.Code)
	}
}

func TestSessionRequest(t *testing.T) {
	req := SessionRequest{
		BaseRequest: BaseRequest{
			DeviceID:  "device123",
			ChannelID: "channel123",
		},
		URL: "webrtc://example.com/live/stream",
	}

	if req.URL != "webrtc://example.com/live/stream" {
		t.Errorf("Expected URL 'webrtc://example.com/live/stream', got '%s'", req.URL)
	}
}

func TestByeRequest(t *testing.T) {
	req := ByeRequest{
		SessionRequest: SessionRequest{
			BaseRequest: BaseRequest{
				DeviceID:  "device123",
				ChannelID: "channel123",
			},
			URL: "webrtc://example.com/live/stream",
		},
	}

	if req.DeviceID != "device123" {
		t.Errorf("Expected DeviceID 'device123', got '%s'", req.DeviceID)
	}
}

func TestPauseRequest(t *testing.T) {
	req := PauseRequest{
		SessionRequest: SessionRequest{
			BaseRequest: BaseRequest{
				DeviceID:  "device123",
				ChannelID: "channel123",
			},
			URL: "webrtc://example.com/live/stream",
		},
	}

	if req.URL != "webrtc://example.com/live/stream" {
		t.Errorf("Expected URL 'webrtc://example.com/live/stream', got '%s'", req.URL)
	}
}

func TestResumeRequest(t *testing.T) {
	req := ResumeRequest{
		SessionRequest: SessionRequest{
			BaseRequest: BaseRequest{
				DeviceID:  "device123",
				ChannelID: "channel123",
			},
			URL: "webrtc://example.com/live/stream",
		},
	}

	if req.ChannelID != "channel123" {
		t.Errorf("Expected ChannelID 'channel123', got '%s'", req.ChannelID)
	}
}

func TestSpeedRequest(t *testing.T) {
	req := SpeedRequest{
		SessionRequest: SessionRequest{
			BaseRequest: BaseRequest{
				DeviceID:  "device123",
				ChannelID: "channel123",
			},
			URL: "webrtc://example.com/live/stream",
		},
		Speed: 2.0,
	}

	if req.Speed != 2.0 {
		t.Errorf("Expected Speed 2.0, got %f", req.Speed)
	}
}

func TestMediaServerRequestJSON(t *testing.T) {
	jsonStr := `{
		"name": "Test Server",
		"type": "SRS",
		"ip": "192.168.1.100",
		"port": 1985,
		"username": "admin",
		"password": "pass123",
		"secret": "secret123",
		"is_default": 1
	}`

	var req MediaServerRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if req.Name != "Test Server" {
		t.Errorf("Expected Name 'Test Server', got '%s'", req.Name)
	}
	if req.Type != "SRS" {
		t.Errorf("Expected Type 'SRS', got '%s'", req.Type)
	}
	if req.Port != 1985 {
		t.Errorf("Expected Port 1985, got %d", req.Port)
	}
}

func TestCommonResponseWithDifferentDataTypes(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{"String data", "test string"},
		{"Integer data", 123},
		{"Map data", map[string]interface{}{"key": "value"}},
		{"Array data", []string{"item1", "item2"}},
		{"Nil data", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := CommonResponse{
				Code: 0,
				Data: tt.data,
			}

			jsonData, err := json.Marshal(resp)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			var decoded CommonResponse
			err = json.Unmarshal(jsonData, &decoded)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if decoded.Code != 0 {
				t.Errorf("Expected Code 0, got %d", decoded.Code)
			}
		})
	}
}
