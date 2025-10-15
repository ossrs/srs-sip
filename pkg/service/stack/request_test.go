package stack

import (
	"testing"

	"github.com/emiago/sipgo/sip"
)

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name      string
		method    sip.RequestMethod
		body      []byte
		conf      OutboundConfig
		wantErr   bool
		errMsg    string
		checkFunc func(*testing.T, *sip.Request)
	}{
		{
			name:   "Valid REGISTER request",
			method: sip.REGISTER,
			body:   nil,
			conf: OutboundConfig{
				Transport: "udp",
				Via:       "192.168.1.100:5060",
				From:      "34020000001320000001",
				To:        "34020000001110000001",
			},
			wantErr: false,
			checkFunc: func(t *testing.T, req *sip.Request) {
				if req.Method != sip.REGISTER {
					t.Errorf("Expected method REGISTER, got %v", req.Method)
				}
				if req.Destination() != "192.168.1.100:5060" {
					t.Errorf("Expected destination 192.168.1.100:5060, got %v", req.Destination())
				}
				if req.Transport() != "udp" {
					t.Errorf("Expected transport udp, got %v", req.Transport())
				}
			},
		},
		{
			name:   "Valid INVITE request with body",
			method: sip.INVITE,
			body:   []byte("v=0\r\no=- 0 0 IN IP4 127.0.0.1\r\n"),
			conf: OutboundConfig{
				Transport: "tcp",
				Via:       "192.168.1.100:5060",
				From:      "34020000001320000001",
				To:        "34020000001110000001",
			},
			wantErr: false,
			checkFunc: func(t *testing.T, req *sip.Request) {
				if req.Method != sip.INVITE {
					t.Errorf("Expected method INVITE, got %v", req.Method)
				}
				if req.Body() == nil {
					t.Error("Expected body to be set")
				}
			},
		},
		{
			name:   "Invalid From length - too short",
			method: sip.REGISTER,
			body:   nil,
			conf: OutboundConfig{
				Transport: "udp",
				Via:       "192.168.1.100:5060",
				From:      "123456789",
				To:        "34020000001110000001",
			},
			wantErr: true,
			errMsg:  "From or To length is not 20",
		},
		{
			name:   "Invalid To length - too long",
			method: sip.REGISTER,
			body:   nil,
			conf: OutboundConfig{
				Transport: "udp",
				Via:       "192.168.1.100:5060",
				From:      "34020000001320000001",
				To:        "340200000011100000012345",
			},
			wantErr: true,
			errMsg:  "From or To length is not 20",
		},
		{
			name:   "Valid MESSAGE request",
			method: sip.MESSAGE,
			body:   []byte("<?xml version=\"1.0\"?><Query></Query>"),
			conf: OutboundConfig{
				Transport: "udp",
				Via:       "192.168.1.100:5060",
				From:      "34020000001320000001",
				To:        "34020000001110000001",
			},
			wantErr: false,
			checkFunc: func(t *testing.T, req *sip.Request) {
				if req.Method != sip.MESSAGE {
					t.Errorf("Expected method MESSAGE, got %v", req.Method)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := NewRequest(tt.method, tt.body, tt.conf)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if req == nil {
				t.Error("Expected request to be non-nil")
				return
			}

			// Run custom checks
			if tt.checkFunc != nil {
				tt.checkFunc(t, req)
			}

			// Common checks
			if req.From() == nil {
				t.Error("Expected From header to be set")
			}
			if req.To() == nil {
				t.Error("Expected To header to be set")
			}
			if req.Contact() == nil {
				t.Error("Expected Contact header to be set")
			}
		})
	}
}

func TestNewRegisterRequest(t *testing.T) {
	conf := OutboundConfig{
		Transport: "udp",
		Via:       "192.168.1.100:5060",
		From:      "34020000001320000001",
		To:        "34020000001110000001",
	}

	req, err := NewRegisterRequest(conf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if req.Method != sip.REGISTER {
		t.Errorf("Expected method REGISTER, got %v", req.Method)
	}

	// Check for Expires header
	expires := req.GetHeader("Expires")
	if expires == nil {
		t.Error("Expected Expires header to be set")
	} else if expires.Value() != "3600" {
		t.Errorf("Expected Expires value 3600, got %v", expires.Value())
	}
}

func TestNewInviteRequest(t *testing.T) {
	conf := OutboundConfig{
		Transport: "udp",
		Via:       "192.168.1.100:5060",
		From:      "34020000001320000001",
		To:        "34020000001110000001",
	}

	body := []byte("v=0\r\no=- 0 0 IN IP4 127.0.0.1\r\ns=Play\r\n")
	subject := "34020000001320000001:0,34020000001110000001:0"

	req, err := NewInviteRequest(body, subject, conf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if req.Method != sip.INVITE {
		t.Errorf("Expected method INVITE, got %v", req.Method)
	}

	// Check for Content-Type header
	contentType := req.GetHeader("Content-Type")
	if contentType == nil {
		t.Error("Expected Content-Type header to be set")
	} else if contentType.Value() != "application/sdp" {
		t.Errorf("Expected Content-Type value application/sdp, got %v", contentType.Value())
	}

	// Check for Subject header
	subjectHeader := req.GetHeader("Subject")
	if subjectHeader == nil {
		t.Error("Expected Subject header to be set")
	} else if subjectHeader.Value() != subject {
		t.Errorf("Expected Subject value %s, got %v", subject, subjectHeader.Value())
	}

	// Check body
	if req.Body() == nil {
		t.Error("Expected body to be set")
	}
}

func TestNewMessageRequest(t *testing.T) {
	conf := OutboundConfig{
		Transport: "udp",
		Via:       "192.168.1.100:5060",
		From:      "34020000001320000001",
		To:        "34020000001110000001",
	}

	body := []byte("<?xml version=\"1.0\"?><Query><CmdType>Catalog</CmdType></Query>")

	req, err := NewMessageRequest(body, conf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if req.Method != sip.MESSAGE {
		t.Errorf("Expected method MESSAGE, got %v", req.Method)
	}

	// Check for Content-Type header
	contentType := req.GetHeader("Content-Type")
	if contentType == nil {
		t.Error("Expected Content-Type header to be set")
	} else if contentType.Value() != "Application/MANSCDP+xml" {
		t.Errorf("Expected Content-Type value Application/MANSCDP+xml, got %v", contentType.Value())
	}

	// Check body
	if req.Body() == nil {
		t.Error("Expected body to be set")
	}
}

func TestNewRequestWithInvalidConfig(t *testing.T) {
	tests := []struct {
		name string
		conf OutboundConfig
	}{
		{
			name: "Empty From",
			conf: OutboundConfig{
				Transport: "udp",
				Via:       "192.168.1.100:5060",
				From:      "",
				To:        "34020000001110000001",
			},
		},
		{
			name: "Empty To",
			conf: OutboundConfig{
				Transport: "udp",
				Via:       "192.168.1.100:5060",
				From:      "34020000001320000001",
				To:        "",
			},
		},
		{
			name: "From length 19",
			conf: OutboundConfig{
				Transport: "udp",
				Via:       "192.168.1.100:5060",
				From:      "3402000000132000001",
				To:        "34020000001110000001",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRequest(sip.REGISTER, nil, tt.conf)
			if err == nil {
				t.Error("Expected error but got nil")
			}
		})
	}
}

