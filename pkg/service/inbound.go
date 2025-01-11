package service

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/emiago/sipgo/sip"
	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-sip/pkg/models"
	"github.com/ossrs/srs-sip/pkg/service/stack"
	"golang.org/x/net/html/charset"
)

const GB28181_ID_LENGTH = 20

func (s *UAS) onRegister(req *sip.Request, tx sip.ServerTransaction) {
	id := req.From().Address.User
	if len(id) != GB28181_ID_LENGTH {
		logger.E(s.ctx, "invalid device ID")
		return
	}

	if s.conf.EnableAuth {
		// Check if Authorization header exists
		authHeader := req.GetHeaders("Authorization")

		// If no Authorization header, send 401 response to request authentication
		if len(authHeader) == 0 {
			nonce := GenerateNonce()
			resp := stack.NewUnauthorizedResponse(req, http.StatusUnauthorized, "Unauthorized", nonce, s.conf.Realm)
			_ = tx.Respond(resp)
			return
		}

		// Validate Authorization
		authInfo := ParseAuthorization(authHeader[0].Value())
		if !ValidateAuth(authInfo, s.conf.Password) {
			logger.Ef(s.ctx, "%s auth failed, source: %s", id, req.Source())
			s.respondRegister(req, http.StatusForbidden, "Auth Failed", tx)
			return
		}
	}

	isUnregister := false
	if exps := req.GetHeaders("Expires"); len(exps) > 0 {
		exp := exps[0]
		expSec, err := strconv.ParseInt(exp.Value(), 10, 32)
		if err != nil {
			logger.Ef(s.ctx, "parse expires header error: %s", err.Error())
			return
		}
		if expSec == 0 {
			isUnregister = true
		}
	} else {
		logger.E(s.ctx, "empty expires header")
		return
	}

	if isUnregister {
		DM.RemoveDevice(id)
		logger.Wf(s.ctx, "Device %s unregistered", id)
		return
	} else {
		if d, ok := DM.GetDevice(id); !ok {
			DM.AddDevice(id, &DeviceInfo{
				DeviceID:    id,
				SourceAddr:  req.Source(),
				NetworkType: req.Transport(),
			})
			s.respondRegister(req, http.StatusOK, "OK", tx)
			logger.Tf(s.ctx, "%s Register success, source:%s, req: %s", id, req.Source(), req.String())

			go s.Catalog(id)
		} else {
			if d.SourceAddr != req.Source() {
				logger.Ef(s.ctx, "Device %s[%s] already registered, %s is NOT allowed.", id, d.SourceAddr, req.Source())
				// TODO: 国标没有明确定义重复ID注册的处理方式，这里暂时返回冲突
				s.respondRegister(req, http.StatusConflict, "Conflict Device ID", tx)
			} else {
				// TODO: 刷新DM里面的设备信息
				s.respondRegister(req, http.StatusOK, "OK", tx)
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
		logger.Ef(s.ctx, "invalid device ID %s", req.String())
	}

	//logger.Tf(s.ctx, "Received MESSAGE: %s", req.String())

	temp := &models.XmlMessageInfo{}
	decoder := xml.NewDecoder(bytes.NewReader([]byte(req.Body())))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(temp); err != nil {
		logger.Ef(s.ctx, "decode message error: %s\n message:%s", err.Error(), req.Body())
	}
	var body string
	switch temp.CmdType {
	case "Keepalive":
		logger.T(s.ctx, "Keepalive")
		if _, ok := DM.GetDevice(temp.DeviceID); !ok {
			// unregister device
			tx.Respond(sip.NewResponseFromRequest(req, http.StatusBadRequest, "", nil))
			return
		}
	case "SensorCatalog": // 兼容宇视，非国标
	case "Catalog":
		logger.T(s.ctx, "Catalog")
		DM.UpdateChannels(temp.DeviceID, temp.DeviceList...)
		//go s.AutoInvite(temp.DeviceID, temp.DeviceList...)
	case "Alarm":
		logger.T(s.ctx, "Alarm")
	case "RecordInfo":
		logger.T(s.ctx, "RecordInfo")
		// 从 recordQueryResults 中获取对应通道的结果通道
		if ch, ok := s.recordQueryResults.Load(temp.DeviceID); ok {
			// 发送查询结果
			resultChan := ch.(chan *models.XmlMessageInfo)
			resultChan <- temp
		}
	default:
		logger.Wf(s.ctx, "Not supported CmdType: %s", temp.CmdType)
		response := sip.NewResponseFromRequest(req, http.StatusBadRequest, "", nil)
		tx.Respond(response)
		return
	}
	tx.Respond(sip.NewResponseFromRequest(req, http.StatusOK, "OK", []byte(body)))
}

func (s *UAS) onNotify(req *sip.Request, tx sip.ServerTransaction) {
	logger.T(s.ctx, "Received NOTIFY request")
	tx.Respond(sip.NewResponseFromRequest(req, http.StatusOK, "OK", nil))
}
