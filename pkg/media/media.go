package media

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
)

type IMedia interface {
	Publish(id, ssrc string) (int, error)
	Unpublish(id string) error
	GetStreamStatus(id string) (bool, error)
	GetAddr() string
	GetWebRTCAddr(id string) string
}

// The r is HTTP API to request, like "http://localhost:1985/gb/v1/publish".
// The req is the HTTP request body, will be marshal to JSON object. nil is no body
// The res is the HTTP response body, already unmarshal to JSON object.
func apiRequest(ctx context.Context, r string, req interface{}, res interface{}) error {
	var buf bytes.Buffer
	if req != nil {
		if err := json.NewEncoder(&buf).Encode(req); err != nil {
			return errors.Wrapf(err, "Marshal body %v", req)
		}
	}
	slog.Debug("API request", "url", r, "size", buf.Len())

	method := "POST"
	if req == nil {
		method = "GET"
	}
	reqObj, err := http.NewRequest(method, r, &buf)
	if err != nil {
		return errors.Wrapf(err, "HTTP request %v", buf.String())
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resObj, err := client.Do(reqObj.WithContext(ctx))
	if err != nil {
		return errors.Wrapf(err, "Do HTTP request %v", buf.String())
	}
	defer resObj.Body.Close()

	if resObj.StatusCode != http.StatusOK {
		return errors.Errorf("Server returned status code=%v", resObj.StatusCode)
	}

	b2, err := io.ReadAll(resObj.Body)
	if err != nil {
		return errors.Wrapf(err, "Read response for %v", buf.String())
	}
	slog.Debug("API response", "url", r, "size", len(b2))

	errorCode := struct {
		Code int `json:"code"`
	}{}
	if err := json.Unmarshal(b2, &errorCode); err != nil {
		return errors.Wrapf(err, "Unmarshal %v", string(b2))
	}
	if errorCode.Code != 0 {
		return errors.Errorf("Server fail code=%v %v", errorCode.Code, string(b2))
	}

	if err := json.Unmarshal(b2, res); err != nil {
		return errors.Wrapf(err, "Unmarshal %v", string(b2))
	}
	slog.Debug("Parse API response", "code", errorCode.Code, "response", res)

	return nil
}
