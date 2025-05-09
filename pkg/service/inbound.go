package service

import (
	"bytes"
	"encoding/xml"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/emiago/sipgo/sip"
	"github.com/ossrs/srs-sip/pkg/models"
	"github.com/ossrs/srs-sip/pkg/service/stack"
	"golang.org/x/net/html/charset"
)

const GB28181_ID_LENGTH = 20

func (s *UAS) isSameIP(addr1, addr2 string) bool {
	ip1, _, err1 := net.SplitHostPort(addr1)
	ip2, _, err2 := net.SplitHostPort(addr2)

	// 如果解析出错，回退到完整字符串比较
	if err1 != nil || err2 != nil {
		return addr1 == addr2
	}

	return ip1 == ip2
}

func (s *UAS) onRegister(req *sip.Request, tx sip.ServerTransaction) {
	id := req.From().Address.User
	if len(id) != GB28181_ID_LENGTH {
		slog.Error("invalid device ID")
		return
	}

	if s.conf.GB28181.Auth.Enable {
		// Check if Authorization header exists
		authHeader := req.GetHeaders("Authorization")

		// If no Authorization header, send 401 response to request authentication
		if len(authHeader) == 0 {
			nonce := GenerateNonce()
			resp := stack.NewUnauthorizedResponse(req, http.StatusUnauthorized, "Unauthorized", nonce, s.conf.GB28181.Realm)
			_ = tx.Respond(resp)
			return
		}

		// Validate Authorization
		authInfo := ParseAuthorization(authHeader[0].Value())
		if !ValidateAuth(authInfo, s.conf.GB28181.Auth.Password) {
			slog.Error("auth failed", "device_id", id, "source", req.Source())
			s.respondRegister(req, http.StatusForbidden, "Auth Failed", tx)
			return
		}
	}

	isUnregister := false
	if exps := req.GetHeaders("Expires"); len(exps) > 0 {
		exp := exps[0]
		expSec, err := strconv.ParseInt(exp.Value(), 10, 32)
		if err != nil {
			slog.Error("parse expires header error", "error", err.Error())
			return
		}
		if expSec == 0 {
			isUnregister = true
		}
	} else {
		slog.Error("empty expires header")
		return
	}

	if isUnregister {
		DM.RemoveDevice(id)
		slog.Warn("Device unregistered", "device_id", id)
		return
	} else {
		if d, ok := DM.GetDevice(id); !ok {
			DM.AddDevice(id, &DeviceInfo{
				DeviceID:    id,
				SourceAddr:  req.Source(),
				NetworkType: req.Transport(),
			})
			s.respondRegister(req, http.StatusOK, "OK", tx)
			slog.Info("Register success", "device_id", id, "source", req.Source(), "request", req.String())

			go s.ConfigDownload(id)
			go s.Catalog(id)
		} else {
			if d.SourceAddr != "" && !s.isSameIP(d.SourceAddr, req.Source()) {
				slog.Error("Device already registered", "device_id", id, "old_source", d.SourceAddr, "new_source", req.Source())
				// TODO: 如果ID重复，应采用虚拟ID
				s.respondRegister(req, http.StatusBadRequest, "Conflict Device ID", tx)
			} else {
				d.SourceAddr = req.Source()
				d.NetworkType = req.Transport()
				DM.UpdateDevice(id, d)
				s.respondRegister(req, http.StatusOK, "OK", tx)

				slog.Info("Re-register success", "device_id", id, "source", req.Source(), "request", req.String())
			}
		}
	}
}

func (s *UAS) respondRegister(req *sip.Request, code sip.StatusCode, reason string, tx sip.ServerTransaction) {
	res := stack.NewRegisterResponse(req, code, reason)
	_ = tx.Respond(res)

}

func (s *UAS) onMessage(req *sip.Request, tx sip.ServerTransaction) {
	id := req.From().Address.User
	if len(id) != 20 {
		slog.Error("invalid device ID", "request", req.String())
	}

	slog.Debug("Received MESSAGE", "request", req.String())

	temp := &models.XmlMessageInfo{}
	decoder := xml.NewDecoder(bytes.NewReader([]byte(req.Body())))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(temp); err != nil {
		slog.Error("decode message error", "error", err.Error(), "message", req.Body())
	}
	var body string
	switch temp.CmdType {
	case "Keepalive":
		slog.Info("Keepalive")
		if d, ok := DM.GetDevice(temp.DeviceID); ok && d.Online {
			// 更新设备心跳时间
			DM.UpdateDeviceHeartbeat(temp.DeviceID)
		} else {
			tx.Respond(sip.NewResponseFromRequest(req, http.StatusBadRequest, "", nil))
			return
		}
	case "SensorCatalog": // 兼容宇视，非国标
	case "Catalog":
		slog.Info("Catalog")
		DM.UpdateChannels(temp.DeviceID, temp.DeviceList...)
		//go s.AutoInvite(temp.DeviceID, temp.DeviceList...)
	case "ConfigDownload":
		slog.Info("ConfigDownload")
		DM.UpdateDeviceConfig(temp.DeviceID, &temp.BasicParam)
	case "Alarm":
		slog.Info("Alarm")
	case "RecordInfo":
		slog.Info("RecordInfo")
		// 从 recordQueryResults 中获取对应通道的结果通道
		if ch, ok := s.recordQueryResults.Load(temp.DeviceID); ok {
			// 发送查询结果
			resultChan := ch.(chan *models.XmlMessageInfo)
			resultChan <- temp
		}
	default:
		slog.Warn("Not supported CmdType", "cmd_type", temp.CmdType)
		response := sip.NewResponseFromRequest(req, http.StatusBadRequest, "", nil)
		tx.Respond(response)
		return
	}
	tx.Respond(sip.NewResponseFromRequest(req, http.StatusOK, "OK", []byte(body)))
}

func (s *UAS) onNotify(req *sip.Request, tx sip.ServerTransaction) {
	slog.Info("Received NOTIFY request")
	tx.Respond(sip.NewResponseFromRequest(req, http.StatusOK, "OK", nil))
}
