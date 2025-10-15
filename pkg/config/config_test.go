package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// 测试 Common 配置
	if cfg.Common.LogLevel != "info" {
		t.Errorf("Expected log level 'info', got '%s'", cfg.Common.LogLevel)
	}
	if cfg.Common.LogFile != "app.log" {
		t.Errorf("Expected log file 'app.log', got '%s'", cfg.Common.LogFile)
	}

	// 测试 GB28181 配置
	if cfg.GB28181.Serial != "34020000002000000001" {
		t.Errorf("Expected serial '34020000002000000001', got '%s'", cfg.GB28181.Serial)
	}
	if cfg.GB28181.Realm != "3402000000" {
		t.Errorf("Expected realm '3402000000', got '%s'", cfg.GB28181.Realm)
	}
	if cfg.GB28181.Host != "0.0.0.0" {
		t.Errorf("Expected host '0.0.0.0', got '%s'", cfg.GB28181.Host)
	}
	if cfg.GB28181.Port != 5060 {
		t.Errorf("Expected port 5060, got %d", cfg.GB28181.Port)
	}
	if cfg.GB28181.Auth.Enable != false {
		t.Errorf("Expected auth enable false, got %v", cfg.GB28181.Auth.Enable)
	}
	if cfg.GB28181.Auth.Password != "123456" {
		t.Errorf("Expected auth password '123456', got '%s'", cfg.GB28181.Auth.Password)
	}

	// 测试 HTTP 配置
	if cfg.Http.Port != 8025 {
		t.Errorf("Expected http port 8025, got %d", cfg.Http.Port)
	}
	if cfg.Http.Dir != "./html" {
		t.Errorf("Expected http dir './html', got '%s'", cfg.Http.Dir)
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	// 测试加载不存在的配置文件，应该返回默认配置
	cfg, err := LoadConfig("non_existent_config.yaml")
	if err != nil {
		t.Fatalf("Expected no error for non-existent config, got: %v", err)
	}

	// 应该返回默认配置
	defaultCfg := DefaultConfig()
	if cfg.Common.LogLevel != defaultCfg.Common.LogLevel {
		t.Errorf("Expected default log level, got '%s'", cfg.Common.LogLevel)
	}
}

func TestLoadConfigValid(t *testing.T) {
	// 创建临时配置文件
	tempFile := "test_config.yaml"
	defer os.Remove(tempFile)

	configContent := `common:
  log-level: debug
  log-file: test.log
gb28181:
  serial: "12345678901234567890"
  realm: "1234567890"
  host: "127.0.0.1"
  port: 5061
  auth:
    enable: true
    password: "test123"
http:
  listen: 9000
  dir: "./test_html"
`

	err := os.WriteFile(tempFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// 加载配置
	cfg, err := LoadConfig(tempFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 验证配置
	if cfg.Common.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", cfg.Common.LogLevel)
	}
	if cfg.Common.LogFile != "test.log" {
		t.Errorf("Expected log file 'test.log', got '%s'", cfg.Common.LogFile)
	}
	if cfg.GB28181.Serial != "12345678901234567890" {
		t.Errorf("Expected serial '12345678901234567890', got '%s'", cfg.GB28181.Serial)
	}
	if cfg.GB28181.Port != 5061 {
		t.Errorf("Expected port 5061, got %d", cfg.GB28181.Port)
	}
	if cfg.GB28181.Auth.Enable != true {
		t.Errorf("Expected auth enable true, got %v", cfg.GB28181.Auth.Enable)
	}
	if cfg.Http.Port != 9000 {
		t.Errorf("Expected http port 9000, got %d", cfg.Http.Port)
	}
}

func TestLoadConfigInvalid(t *testing.T) {
	// 创建无效的配置文件
	tempFile := "test_invalid_config.yaml"
	defer os.Remove(tempFile)

	invalidContent := `invalid yaml content: [[[`

	err := os.WriteFile(tempFile, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// 加载配置应该失败
	_, err = LoadConfig(tempFile)
	if err == nil {
		t.Error("Expected error for invalid config file, got nil")
	}
}

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	
	// 在某些环境下可能没有网络接口，所以允许返回错误
	if err != nil {
		t.Logf("GetLocalIP returned error (may be expected in some environments): %v", err)
		return
	}

	// 如果成功，验证返回的是有效的 IP 地址
	if ip == "" {
		t.Error("Expected non-empty IP address")
	}

	// 简单验证 IP 格式（应该包含点号）
	if len(ip) < 7 { // 最短的 IP 是 0.0.0.0
		t.Errorf("IP address seems invalid: %s", ip)
	}

	t.Logf("Local IP: %s", ip)
}

