package media

import (
	"context"

	"github.com/ossrs/go-oryx-lib/errors"
)

type Zlm struct {
	Ctx    context.Context
	Schema string // The schema of ZLM, eg: http
	Addr   string // The address of ZLM, eg: localhost:8085
	Secret string // The secret of ZLM, eg: ZLMediaKit_secret
}

// /index/api/openRtpServer
// secret={{ZLMediaKit_secret}}&port=0&enable_tcp=1&stream_id=test2
func (z *Zlm) Publish(id, ssrc string) (int, error) {

	res := struct {
		Code int `json:"code"`
		Port int `json:"port"`
	}{}

	if err := apiRequest(z.Ctx, z.Schema+"://"+z.Addr+"/index/api/openRtpServer?secret="+z.Secret+"&port=0&enable_tcp=1&stream_id="+id+"&ssrc="+ssrc, nil, &res); err != nil {
		return 0, errors.Wrapf(err, "gb/v1/publish")
	}
	return res.Port, nil
}

// /index/api/closeRtpServer
func (z *Zlm) Unpublish(id string) error {
	res := struct {
		Code int `json:"code"`
	}{}
	if err := apiRequest(z.Ctx, z.Schema+"://"+z.Addr+"/index/api/closeRtpServer?secret="+z.Secret+"&stream_id="+id, nil, &res); err != nil {
		return errors.Wrapf(err, "gb/v1/publish")
	}
	return nil
}

// /index/api/getMediaList
func (z *Zlm) GetStreamStatus(id string) (bool, error) {
	res := struct {
		Code int `json:"code"`
	}{}
	if err := apiRequest(z.Ctx, z.Schema+"://"+z.Addr+"/index/api/getMediaList?secret="+z.Secret+"&stream_id="+id, nil, &res); err != nil {
		return false, errors.Wrapf(err, "gb/v1/publish")
	}
	return res.Code == 0, nil
}

func (z *Zlm) GetAddr() string {
	return z.Addr
}

func (z *Zlm) GetWebRTCAddr(id string) string {
	return "http://" + z.Addr + "/index/api/webrtc?app=rtp&stream=" + id + "&type=play"
}
