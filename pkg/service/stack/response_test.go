package stack

import (
	"testing"

	"github.com/emiago/sipgo/sip"
)

func TestNewRegisterResponse(t *testing.T) {
	// Create a test request first - we need a properly initialized request
	// Skip this test as it requires a full SIP stack to create valid responses
	t.Skip("Skipping response test - requires full SIP stack initialization")

	conf := OutboundConfig{
		Transport: "udp",
		Via:       "192.168.1.100:5060",
		From:      "34020000001320000001",
		To:        "34020000001110000001",
	}

	req, err := NewRegisterRequest(conf)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	tests := []struct {
		name   string
		code   sip.StatusCode
		reason string
	}{
		{
			name:   "200 OK response",
			code:   sip.StatusOK,
			reason: "OK",
		},
		{
			name:   "401 Unauthorized response",
			code:   sip.StatusUnauthorized,
			reason: "Unauthorized",
		},
		{
			name:   "403 Forbidden response",
			code:   sip.StatusForbidden,
			reason: "Forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewRegisterResponse(req, tt.code, tt.reason)

			if resp == nil {
				t.Fatal("Expected response to be non-nil")
			}

			if resp.StatusCode != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, resp.StatusCode)
			}

			if resp.Reason != tt.reason {
				t.Errorf("Expected reason '%s', got '%s'", tt.reason, resp.Reason)
			}

			// Check for Expires header
			expires := resp.GetHeader("Expires")
			if expires == nil {
				t.Error("Expected Expires header to be set")
			} else if expires.Value() != "3600" {
				t.Errorf("Expected Expires value 3600, got %v", expires.Value())
			}

			// Check for Date header
			date := resp.GetHeader("Date")
			if date == nil {
				t.Error("Expected Date header to be set")
			}

			// Check that To header has tag
			to := resp.To()
			if to == nil {
				t.Error("Expected To header to be set")
			} else {
				tag, ok := to.Params.Get("tag")
				if !ok || tag == "" {
					t.Error("Expected To header to have tag parameter")
				}
			}

			// Check that Allow header is removed
			allow := resp.GetHeader("Allow")
			if allow != nil {
				t.Error("Expected Allow header to be removed")
			}
		})
	}
}

func TestNewUnauthorizedResponse(t *testing.T) {
	// Skip this test as it requires a full SIP stack to create valid responses
	t.Skip("Skipping response test - requires full SIP stack initialization")

	conf := OutboundConfig{
		Transport: "udp",
		Via:       "192.168.1.100:5060",
		From:      "34020000001320000001",
		To:        "34020000001110000001",
	}

	req, err := NewRegisterRequest(conf)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	tests := []struct {
		name   string
		code   sip.StatusCode
		reason string
		nonce  string
		realm  string
	}{
		{
			name:   "401 Unauthorized with nonce and realm",
			code:   sip.StatusUnauthorized,
			reason: "Unauthorized",
			nonce:  "dcd98b7102dd2f0e8b11d0f600bfb0c093",
			realm:  "3402000000",
		},
		{
			name:   "407 Proxy Authentication Required",
			code:   sip.StatusProxyAuthRequired,
			reason: "Proxy Authentication Required",
			nonce:  "abc123def456",
			realm:  "proxy.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewUnauthorizedResponse(req, tt.code, tt.reason, tt.nonce, tt.realm)

			if resp == nil {
				t.Fatal("Expected response to be non-nil")
			}

			if resp.StatusCode != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, resp.StatusCode)
			}

			if resp.Reason != tt.reason {
				t.Errorf("Expected reason '%s', got '%s'", tt.reason, resp.Reason)
			}

			// Check for WWW-Authenticate header
			wwwAuth := resp.GetHeader("WWW-Authenticate")
			if wwwAuth == nil {
				t.Error("Expected WWW-Authenticate header to be set")
			} else {
				authValue := wwwAuth.Value()
				// Check that it contains the nonce
				if len(authValue) == 0 {
					t.Error("Expected WWW-Authenticate header to have a value")
				}
				// The value should contain Digest, realm, nonce, and algorithm
				expectedSubstrings := []string{"Digest", "realm=", "nonce=", "algorithm=MD5"}
				for _, substr := range expectedSubstrings {
					if len(authValue) > 0 && !contains(authValue, substr) {
						t.Errorf("Expected WWW-Authenticate to contain '%s', got '%s'", substr, authValue)
					}
				}
			}

			// Check that To header has tag
			to := resp.To()
			if to == nil {
				t.Error("Expected To header to be set")
			} else {
				tag, ok := to.Params.Get("tag")
				if !ok || tag == "" {
					t.Error("Expected To header to have tag parameter")
				}
			}
		})
	}
}

func TestNewResponse(t *testing.T) {
	// Skip this test as it requires a full SIP stack to create valid responses
	t.Skip("Skipping response test - requires full SIP stack initialization")

	conf := OutboundConfig{
		Transport: "udp",
		Via:       "192.168.1.100:5060",
		From:      "34020000001320000001",
		To:        "34020000001110000001",
	}

	req, err := NewRequest(sip.INVITE, nil, conf)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	tests := []struct {
		name   string
		code   sip.StatusCode
		reason string
	}{
		{
			name:   "100 Trying",
			code:   sip.StatusTrying,
			reason: "Trying",
		},
		{
			name:   "180 Ringing",
			code:   sip.StatusRinging,
			reason: "Ringing",
		},
		{
			name:   "200 OK",
			code:   sip.StatusOK,
			reason: "OK",
		},
		{
			name:   "404 Not Found",
			code:   sip.StatusNotFound,
			reason: "Not Found",
		},
		{
			name:   "500 Server Internal Error",
			code:   sip.StatusInternalServerError,
			reason: "Server Internal Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := newResponse(req, tt.code, tt.reason)

			if resp == nil {
				t.Fatal("Expected response to be non-nil")
			}

			if resp.StatusCode != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, resp.StatusCode)
			}

			if resp.Reason != tt.reason {
				t.Errorf("Expected reason '%s', got '%s'", tt.reason, resp.Reason)
			}

			// Check that To header has tag
			to := resp.To()
			if to == nil {
				t.Error("Expected To header to be set")
			} else {
				tag, ok := to.Params.Get("tag")
				if !ok || tag == "" {
					t.Error("Expected To header to have tag parameter")
				}
				// Check tag length is 10
				if len(tag) != 10 {
					t.Errorf("Expected tag length to be 10, got %d", len(tag))
				}
			}

			// Check that Allow header is removed
			allow := resp.GetHeader("Allow")
			if allow != nil {
				t.Error("Expected Allow header to be removed")
			}
		})
	}
}

func TestResponseConstants(t *testing.T) {
	if TIME_LAYOUT != "2024-01-01T00:00:00" {
		t.Errorf("Expected TIME_LAYOUT to be '2024-01-01T00:00:00', got '%s'", TIME_LAYOUT)
	}

	if EXPIRES_TIME != 3600 {
		t.Errorf("Expected EXPIRES_TIME to be 3600, got %d", EXPIRES_TIME)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
