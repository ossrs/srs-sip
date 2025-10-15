package utils

import (
	"testing"
)

func TestGenRandomNumber(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"Generate 1 digit", 1},
		{"Generate 5 digits", 5},
		{"Generate 9 digits", 9},
		{"Generate 10 digits", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenRandomNumber(tt.length)

			// 验证长度
			if len(result) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(result))
			}

			// 验证所有字符都是数字
			for i, c := range result {
				if c < '0' || c > '9' {
					t.Errorf("Character at position %d is not a digit: %c", i, c)
				}
			}
		})
	}
}

func TestGenRandomNumberUniqueness(t *testing.T) {
	// 生成多个随机数，验证它们不完全相同（虽然理论上可能相同，但概率极低）
	results := make(map[string]bool)
	iterations := 100
	length := 10

	for i := 0; i < iterations; i++ {
		result := GenRandomNumber(length)
		results[result] = true
	}

	// 至少应该有一些不同的值（不太可能100次都生成相同的10位数）
	if len(results) < 50 {
		t.Errorf("Expected at least 50 unique values out of %d iterations, got %d", iterations, len(results))
	}
}

func TestCreateSSRC(t *testing.T) {
	tests := []struct {
		name     string
		isLive   bool
		expected byte
	}{
		{"Live stream SSRC", true, '0'},
		{"Non-live stream SSRC", false, '1'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssrc := CreateSSRC(tt.isLive)

			// 验证长度为10
			if len(ssrc) != 10 {
				t.Errorf("Expected SSRC length 10, got %d", len(ssrc))
			}

			// 验证第一个字符
			if ssrc[0] != tt.expected {
				t.Errorf("Expected first character '%c', got '%c'", tt.expected, ssrc[0])
			}

			// 验证所有字符都是数字
			for i, c := range ssrc {
				if c < '0' || c > '9' {
					t.Errorf("Character at position %d is not a digit: %c", i, c)
				}
			}
		})
	}
}

func TestCreateSSRCUniqueness(t *testing.T) {
	// 测试生成的 SSRC 具有唯一性
	results := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		ssrc := CreateSSRC(true)
		results[ssrc] = true
	}

	// 应该有很多不同的值
	if len(results) < 50 {
		t.Errorf("Expected at least 50 unique SSRCs out of %d iterations, got %d", iterations, len(results))
	}
}

func TestIsVideoChannel(t *testing.T) {
	tests := []struct {
		name      string
		channelID string
		expected  bool
	}{
		{
			name:      "Video channel type 131",
			channelID: "34020000001310000001",
			expected:  true,
		},
		{
			name:      "Video channel type 132",
			channelID: "34020000001320000001",
			expected:  true,
		},
		{
			name:      "Audio channel type 137",
			channelID: "34020000001370000001",
			expected:  false,
		},
		{
			name:      "Alarm channel type 134",
			channelID: "34020000001340000001",
			expected:  false,
		},
		{
			name:      "Other device type",
			channelID: "34020000001110000001",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsVideoChannel(tt.channelID)
			if result != tt.expected {
				t.Errorf("IsVideoChannel(%s) = %v, expected %v", tt.channelID, result, tt.expected)
			}
		})
	}
}

func TestGetSessionName(t *testing.T) {
	tests := []struct {
		name     string
		playType int
		expected string
	}{
		{"Live play", 0, "Play"},
		{"Playback", 1, "Playback"},
		{"Download", 2, "Download"},
		{"Talk", 3, "Talk"},
		{"Unknown type", 99, "Play"},
		{"Negative type", -1, "Play"},
		{"Type 4", 4, "Play"},
		{"Type 5", 5, "Play"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSessionName(tt.playType)
			if result != tt.expected {
				t.Errorf("GetSessionName(%d) = %s, expected %s", tt.playType, result, tt.expected)
			}
		})
	}
}

func TestGenRandomNumberZeroLength(t *testing.T) {
	result := GenRandomNumber(0)
	if len(result) != 0 {
		t.Errorf("Expected empty string for length 0, got %s", result)
	}
}

func TestCreateSSRCBothTypes(t *testing.T) {
	// Test both live and non-live in one test
	liveSSRC := CreateSSRC(true)
	nonLiveSSRC := CreateSSRC(false)

	if liveSSRC[0] != '0' {
		t.Errorf("Live SSRC should start with '0', got '%c'", liveSSRC[0])
	}

	if nonLiveSSRC[0] != '1' {
		t.Errorf("Non-live SSRC should start with '1', got '%c'", nonLiveSSRC[0])
	}

	// They should be different (with very high probability)
	if liveSSRC == nonLiveSSRC {
		t.Error("Live and non-live SSRCs should be different")
	}
}
