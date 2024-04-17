package signaling

import (
	"context"

	"github.com/ossrs/go-oryx-lib/errors"
)

type Srs struct {
	Ctx  context.Context
	Addr string // The address of SRS, eg: http://localhost:1985
}

func (s *Srs) Publish(id, ssrc string) (int, error) {
	req := struct {
		Id   string `json:"id"`
		SSRC string `json:"ssrc"`
	}{
		id, ssrc,
	}

	res := struct {
		Code int `json:"code"`
		Port int `json:"port"`
	}{}

	if err := apiRequest(s.Ctx, s.Addr+"/gb/v1/publish/", req, &res); err != nil {
		return 0, errors.Wrapf(err, "gb/v1/publish")
	}

	return res.Port, nil
}

func (s *Srs) Unpublish(id string) error {
	return nil
}

//	{
//	    "code": 0,
//	    "server": "vid-y19n6nm",
//	    "service": "382k456r",
//	    "pid": "9495",
//	    "streams": [{
//	        "id": "vid-9y0ozy0",
//	        "name": "0551954854",
//	        "vhost": "vid-v2ws53u",
//	        "app": "live",
//	        "tcUrl": "webrtc://127.0.0.1:1985/live",
//	        "url": "/live/0551954854",
//	        "live_ms": 1720428680003,
//	        "clients": 1,
//	        "frames": 8431,
//	        "send_bytes": 66463941,
//	        "recv_bytes": 89323998,
//	        "kbps": {
//	            "recv_30s": 0,
//	            "send_30s": 0
//	        },
//	        "publish": {
//	            "active": false,
//	            "cid": "b3op069g"
//	        },
//	        "video": null,
//	        "audio": null
//	    }]
//	}
func (s *Srs) GetStreamStatus(id string) (bool, error) {
	type Stream struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Publish struct {
			Active bool   `json:"active"`
			Cid    string `json:"cid"`
		} `json:"publish"`
	}
	res := struct {
		Code    int      `json:"code"`
		Streams []Stream `json:"streams"`
	}{}

	if err := apiRequest(s.Ctx, s.Addr+"/api/v1/streams?count=99", nil, &res); err != nil {
		return false, errors.Wrapf(err, "api/v1/stream")
	}

	if len(res.Streams) == 0 {
		return false, nil
	} else {
		for _, v := range res.Streams {
			if v.Name == id {
				return v.Publish.Active, nil
			}
		}
	}
	return false, nil
}
