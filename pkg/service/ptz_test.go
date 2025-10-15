package service

import (
	"testing"
)

func TestGetPTZSpeed(t *testing.T) {
	tests := []struct {
		name     string
		speed    string
		expected uint8
	}{
		{"Speed 1", "1", 25},
		{"Speed 2", "2", 50},
		{"Speed 3", "3", 75},
		{"Speed 4", "4", 100},
		{"Speed 5", "5", 125},
		{"Speed 6", "6", 150},
		{"Speed 7", "7", 175},
		{"Speed 8", "8", 200},
		{"Speed 9", "9", 225},
		{"Speed 10", "10", 255},
		{"Invalid speed", "invalid", 125}, // 默认速度
		{"Empty speed", "", 125},          // 默认速度
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPTZSpeed(tt.speed)
			if result != tt.expected {
				t.Errorf("getPTZSpeed(%s) = %d, expected %d", tt.speed, result, tt.expected)
			}
		})
	}
}

func TestToPTZCmd(t *testing.T) {
	tests := []struct {
		name        string
		cmdName     string
		speed       string
		expectError bool
		checkPrefix bool
	}{
		{"Stop command", "stop", "5", false, true},
		{"Right command", "right", "5", false, true},
		{"Left command", "left", "5", false, true},
		{"Up command", "up", "5", false, true},
		{"Down command", "down", "5", false, true},
		{"Up-right command", "upright", "5", false, true},
		{"Up-left command", "upleft", "5", false, true},
		{"Down-right command", "downright", "5", false, true},
		{"Down-left command", "downleft", "5", false, true},
		{"Zoom in command", "zoomin", "5", false, true},
		{"Zoom out command", "zoomout", "5", false, true},
		{"Invalid command", "invalid", "5", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toPTZCmd(tt.cmdName, tt.speed)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for command %s, got nil", tt.cmdName)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for command %s: %v", tt.cmdName, err)
				return
			}

			// 验证结果格式
			if len(result) != 16 { // A50F01 + 5对字节 = 16个字符
				t.Errorf("Expected result length 16, got %d for command %s", len(result), tt.cmdName)
			}

			// 验证前缀
			if tt.checkPrefix && result[:6] != "A50F01" {
				t.Errorf("Expected prefix 'A50F01', got '%s' for command %s", result[:6], tt.cmdName)
			}
		})
	}
}

func TestToPTZCmdSpecificCases(t *testing.T) {
	// 测试停止命令
	t.Run("Stop command details", func(t *testing.T) {
		result, err := toPTZCmd("stop", "5")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Stop 命令码是 0，速度应该都是 0
		// A50F01 00 00 00 00 checksum
		if result[:8] != "A50F0100" {
			t.Errorf("Stop command should start with A50F0100, got %s", result[:8])
		}
	})

	// 测试右移命令
	t.Run("Right command details", func(t *testing.T) {
		result, err := toPTZCmd("right", "5")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Right 命令码是 1，水平速度应该是 125 (0x7D)
		// A50F01 01 7D 00 00 checksum
		if result[:8] != "A50F0101" {
			t.Errorf("Right command should start with A50F0101, got %s", result[:8])
		}
	})

	// 测试上移命令
	t.Run("Up command details", func(t *testing.T) {
		result, err := toPTZCmd("up", "5")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Up 命令码是 8，垂直速度应该是 125 (0x7D)
		// A50F01 08 00 7D 00 checksum
		if result[:8] != "A50F0108" {
			t.Errorf("Up command should start with A50F0108, got %s", result[:8])
		}
	})

	// 测试缩放命令
	t.Run("Zoom in command details", func(t *testing.T) {
		result, err := toPTZCmd("zoomin", "5")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Zoom in 命令码是 16 (0x10)
		// A50F01 10 00 00 XX checksum (XX 是速度左移4位)
		if result[:8] != "A50F0110" {
			t.Errorf("Zoom in command should start with A50F0110, got %s", result[:8])
		}
	})
}

func TestToPTZCmdWithDifferentSpeeds(t *testing.T) {
	speeds := []string{"1", "5", "10"}

	for _, speed := range speeds {
		t.Run("Right with speed "+speed, func(t *testing.T) {
			result, err := toPTZCmd("right", speed)
			if err != nil {
				t.Errorf("Unexpected error with speed %s: %v", speed, err)
			}
			if len(result) != 16 {
				t.Errorf("Expected length 16, got %d", len(result))
			}
		})
	}
}

func TestPTZCmdMap(t *testing.T) {
	// 验证所有预定义的命令都存在
	expectedCommands := []string{
		"stop", "right", "left", "down", "downright", "downleft",
		"up", "upright", "upleft", "zoomin", "zoomout",
	}

	for _, cmd := range expectedCommands {
		t.Run("Command exists: "+cmd, func(t *testing.T) {
			if _, ok := ptzCmdMap[cmd]; !ok {
				t.Errorf("Command %s not found in ptzCmdMap", cmd)
			}
		})
	}
}

func TestPTZSpeedMap(t *testing.T) {
	// 验证速度映射的正确性
	expectedSpeeds := map[string]uint8{
		"1":  25,
		"2":  50,
		"3":  75,
		"4":  100,
		"5":  125,
		"6":  150,
		"7":  175,
		"8":  200,
		"9":  225,
		"10": 255,
	}

	for speed, expectedValue := range expectedSpeeds {
		t.Run("Speed mapping: "+speed, func(t *testing.T) {
			if value, ok := ptzSpeedMap[speed]; !ok {
				t.Errorf("Speed %s not found in ptzSpeedMap", speed)
			} else if value != expectedValue {
				t.Errorf("Speed %s expected value %d, got %d", speed, expectedValue, value)
			}
		})
	}
}
