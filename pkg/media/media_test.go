package media

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApiRequest_Success(t *testing.T) {
	// Create a test server that returns a successful response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"code": 0,
			"data": map[string]string{
				"message": "success",
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	ctx := context.Background()
	req := map[string]string{"test": "data"}
	var res map[string]interface{}

	err := apiRequest(ctx, server.URL, req, &res)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if res["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", res["code"])
	}
}

func TestApiRequest_GetMethod(t *testing.T) {
	// Create a test server that checks the HTTP method
	methodReceived := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methodReceived = r.Method
		response := map[string]interface{}{
			"code": 0,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	ctx := context.Background()
	var res map[string]interface{}

	// When req is nil, should use GET method
	err := apiRequest(ctx, server.URL, nil, &res)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if methodReceived != "GET" {
		t.Errorf("Expected GET method, got %s", methodReceived)
	}
}

func TestApiRequest_PostMethod(t *testing.T) {
	// Create a test server that checks the HTTP method
	methodReceived := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methodReceived = r.Method
		response := map[string]interface{}{
			"code": 0,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	ctx := context.Background()
	req := map[string]string{"test": "data"}
	var res map[string]interface{}

	// When req is not nil, should use POST method
	err := apiRequest(ctx, server.URL, req, &res)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if methodReceived != "POST" {
		t.Errorf("Expected POST method, got %s", methodReceived)
	}
}

func TestApiRequest_ServerError(t *testing.T) {
	// Create a test server that returns an error status code
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	ctx := context.Background()
	var res map[string]interface{}

	err := apiRequest(ctx, server.URL, nil, &res)
	if err == nil {
		t.Error("Expected error for server error status code")
	}
}

func TestApiRequest_NonZeroCode(t *testing.T) {
	// Create a test server that returns a non-zero error code
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"code":    100,
			"message": "error message",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	ctx := context.Background()
	var res map[string]interface{}

	err := apiRequest(ctx, server.URL, nil, &res)
	if err == nil {
		t.Error("Expected error for non-zero code")
	}
}

func TestApiRequest_InvalidJSON(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	ctx := context.Background()
	var res map[string]interface{}

	err := apiRequest(ctx, server.URL, nil, &res)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestApiRequest_ContextCancellation(t *testing.T) {
	// Create a test server that delays the response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		response := map[string]interface{}{
			"code": 0,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create a context that is already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var res map[string]interface{}

	err := apiRequest(ctx, server.URL, nil, &res)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}
}

func TestApiRequest_InvalidURL(t *testing.T) {
	ctx := context.Background()
	var res map[string]interface{}

	// Test with invalid URL
	err := apiRequest(ctx, "://invalid-url", nil, &res)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestApiRequest_Timeout(t *testing.T) {
	// Create a test server that never responds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second) // Longer than the 10 second timeout
	}))
	defer server.Close()

	ctx := context.Background()
	var res map[string]interface{}

	err := apiRequest(ctx, server.URL, nil, &res)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestApiRequest_ComplexResponse(t *testing.T) {
	// Create a test server that returns a complex response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"code": 0,
			"data": map[string]interface{}{
				"id":     "12345",
				"status": "active",
				"items": []string{
					"item1",
					"item2",
					"item3",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	ctx := context.Background()
	var res map[string]interface{}

	err := apiRequest(ctx, server.URL, nil, &res)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the response was properly unmarshaled
	if res["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", res["code"])
	}

	data, ok := res["data"].(map[string]interface{})
	if !ok {
		t.Error("Expected data to be a map")
	} else {
		if data["id"].(string) != "12345" {
			t.Errorf("Expected id '12345', got %v", data["id"])
		}
	}
}

func TestApiRequest_EmptyResponse(t *testing.T) {
	// Create a test server that returns minimal response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"code": 0,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	ctx := context.Background()
	var res map[string]interface{}

	err := apiRequest(ctx, server.URL, nil, &res)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if res["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", res["code"])
	}
}

func TestApiRequest_WithRequestBody(t *testing.T) {
	// Create a test server that echoes back the request
	var receivedBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		response := map[string]interface{}{
			"code": 0,
			"echo": receivedBody,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	ctx := context.Background()
	req := map[string]interface{}{
		"id":   "test-id",
		"ssrc": "test-ssrc",
	}
	var res map[string]interface{}

	err := apiRequest(ctx, server.URL, req, &res)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the server received the request body
	if receivedBody["id"].(string) != "test-id" {
		t.Errorf("Expected id 'test-id', got %v", receivedBody["id"])
	}
}

