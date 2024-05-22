package service

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/sip"
	"github.com/gorilla/mux"
	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
	"golang.org/x/net/html/charset"

	"github.com/ossrs/srs-sip/pkg/config"
	"github.com/ossrs/srs-sip/pkg/utils"
)

const TIME_LAYOUT = "2024-01-01T00:00:00"

type GB28181Server struct {
	ctx  context.Context
	conf *config.MainConfig

	SN uint32

	srv gosip.Server
}

var dm = GetDeviceManager()

func NewGB28181Server() *GB28181Server {
	return &GB28181Server{}
}

func Run(ctx context.Context, r0 interface{}) {
	s := NewGB28181Server()

	s.startHttpServer()
	s.startGbServer(ctx, r0)
}

func (s *GB28181Server) startHttpServer() {
	router := mux.NewRouter().StrictSlash(true)
	RegisterRoutes(router)

	go func() {
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			panic(err)
		}
	}()
}

func (s *GB28181Server) startGbServer(ctx context.Context, r0 interface{}) {
	conf := r0.(*config.MainConfig)
	s.ctx = ctx
	s.conf = conf

	addr := conf.SipHost + ":" + strconv.Itoa(int(conf.SipPort))

	srvConf := gosip.ServerConfig{}
	srvConf.Host = conf.SipHost

	s.srv = gosip.NewServer(srvConf, nil, nil, log.NewDefaultLogrusLogger().WithPrefix("GB-SIP"))
	s.srv.OnRequest(sip.REGISTER, s.OnRegister)
	s.srv.OnRequest(sip.MESSAGE, s.OnMessage)
	s.srv.OnRequest(sip.NOTIFY, s.OnNotify)
	s.srv.OnRequest(sip.BYE, s.OnBye)

	if err := s.srv.Listen("udp", addr); err != nil {
		logger.Ef(s.ctx, "listen udp %d error %s", conf.SipPort, err.Error())
	}

	if err := s.srv.Listen("tcp", addr); err != nil {
		logger.Ef(s.ctx, "listen tcp %d error %s", conf.SipPort, err.Error())
	}

	logger.Tf(s.ctx, "GB28181 server started, listen on udp and tcp %d", conf.SipPort)
}

func (s *GB28181Server) OnRegister(req sip.Request, tx sip.ServerTransaction) {
	from, err := req.From()
	if !err || from.Address == nil || from.Address.User() == nil {
		logger.E(s.ctx, "empty device ID")
		return
	}
	id := from.Address.User().String()

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

	if len(id) != 20 {
		logger.E(s.ctx, "invalid device ID")
		return
	}

	if isUnregister {
		dm.RemoveDevice(id)
		logger.Wf(s.ctx, "Device %s unregistered", id)
		return
	} else {
		if _, ok := dm.GetDevice(id); !ok {
			dm.AddDevice(id, &DeviceInfo{
				DeviceID:    id,
				SourceAddr:  req.Source(),
				NetworkType: req.Transport(),
			})
			s.respondRegister(req, http.StatusOK, "OK", tx)
			logger.Tf(s.ctx, "%s Register success, source:%s, req: %s", id, req.Source(), req.String())

			go s.Catalog(id)
		} else {
			logger.Ef(s.ctx, "Device %s already registered", id)
			// TODO: 国标没有明确定义重复ID注册的处理方式，这里暂时返回冲突
			s.respondRegister(req, http.StatusConflict, "Conflict Device ID", tx)
		}
	}
}

func (s *GB28181Server) respondRegister(req sip.Request, code sip.StatusCode, reason string, tx sip.ServerTransaction) {
	resp := sip.NewResponseFromRequest("", req, code, reason, "")
	to, _ := resp.To()
	resp.ReplaceHeaders("To", []sip.Header{&sip.ToHeader{Address: to.Address, Params: sip.NewParams().Add("tag", sip.String{Str: utils.GenRandomNumber(9)})}})
	resp.RemoveHeader("Allow")
	expires := sip.Expires(3600)
	resp.AppendHeader(&expires)
	resp.AppendHeader(&sip.GenericHeader{
		HeaderName: "Date",
		Contents:   time.Now().Format(TIME_LAYOUT),
	})
	_ = tx.Respond(resp)

}

func (s *GB28181Server) OnMessage(req sip.Request, tx sip.ServerTransaction) {
	from, err := req.From()
	if !err || from.Address == nil || from.Address.User() == nil {
		logger.E(s.ctx, "empty device ID")
		return
	}

	//logger.Tf(s.ctx, "Received MESSAGE: %s", req.String())

	temp := &struct {
		XMLName      xml.Name
		CmdType      string
		SN           int // 请求序列号，一般用于对应 request 和 response
		DeviceID     string
		DeviceName   string
		Manufacturer string
		Model        string
		Channel      string
		DeviceList   []ChannelInfo `xml:"DeviceList>Item"`
		// RecordList   []*Record     `xml:"RecordList>Item"`
		// SumNum int // 录像结果的总数 SumNum，录像结果会按照多条消息返回，可用于判断是否全部返回
	}{}
	decoder := xml.NewDecoder(bytes.NewReader([]byte(req.Body())))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(temp); err != nil {
		logger.Ef(s.ctx, "decode message error: %s", err.Error())
	}
	var body string
	switch temp.CmdType {
	case "Keepalive":
		logger.T(s.ctx, "Keepalive")
		if _, ok := dm.GetDevice(temp.DeviceID); !ok {
			// unregister device
			tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, "", body))
			return
		}
	case "Catalog":
		logger.T(s.ctx, "Catalog")
		dm.UpdateChannels(temp.DeviceID, temp.DeviceList...)
		go s.AutoInvite(temp.DeviceList...)
	case "Alarm":
		logger.T(s.ctx, "Alarm")
	default:
		logger.Wf(s.ctx, "Not supported CmdType: %s", temp.CmdType)
		response := sip.NewResponseFromRequest("", req, http.StatusBadRequest, "", "")
		tx.Respond(response)
		return
	}
	tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", body))
}

func (s *GB28181Server) OnNotify(req sip.Request, tx sip.ServerTransaction) {
	logger.T(s.ctx, "Received NOTIFY request")
	tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", ""))
}

func (s *GB28181Server) OnBye(req sip.Request, tx sip.ServerTransaction) {
	logger.T(s.ctx, "Received BYE request")
	tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", ""))
}

func (s *GB28181Server) AutoInvite(list ...ChannelInfo) {
	for _, c := range list {
		if c.Status == "ON" && s.isVideoChannel(c.DeviceID) {
			if err := s.Invite(c.DeviceID); err != nil {
				logger.Ef(s.ctx, "invite error: %s", err.Error())
			}
		}
	}
}

// @see GB/T28181—2016 附录D 统一编码规则
func (s *GB28181Server) isVideoChannel(channelID string) bool {
	deviceType := channelID[10:13]
	return deviceType == "131" || deviceType == "132"
}

func (s *GB28181Server) Invite(channelID string) error {
	ssrc := utils.CreateSSRC(true)

	mediaAddr := "http://" + s.conf.MediaAddr
	mediaPort, err := utils.ApiGbPublishRequest(s.ctx, mediaAddr, ssrc, ssrc)
	if err != nil {
		return errors.Wrapf(err, "api gb publish request error")
	}

	mediaHost := strings.Split(s.conf.MediaAddr, ":")[0]
	if mediaHost == "" {
		return errors.Errorf("media host is empty")
	}

	sdpInfo := []string{
		"v=0",
		fmt.Sprintf("o=%s 0 0 IN IP4 %s", channelID, mediaHost),
		"s=" + "Play",
		"u=" + channelID + ":0",
		"c=IN IP4 " + mediaHost,
		"t=0 0", // start time and end time
		fmt.Sprintf("m=video %d TCP/RTP/AVP 96", mediaPort),
		"a=recvonly",
		"a=rtpmap:96 PS/90000",
		"y=" + ssrc,
	}
	if true { // support tcp only
		sdpInfo = append(sdpInfo, "a=setup:passive", "a=connection:new")
	}

	d, ok := dm.GetDeviceByChannel(channelID)
	if !ok {
		return errors.Errorf("device not found by %s", channelID)
	}
	invite, err := s.CreateRequest(sip.INVITE, channelID, d.SourceAddr, d.NetworkType)
	if err != nil {
		return errors.Wrapf(err, "create invite request error")
	}
	contentType := sip.ContentType("application/sdp")
	invite.AppendHeader(&contentType)
	invite.SetBody(strings.Join(sdpInfo, "\r\n")+"\r\n", true)
	subject := sip.GenericHeader{
		HeaderName: "Subject", Contents: fmt.Sprintf("%s:%s,%s:0", channelID, ssrc, s.conf.Serial),
	}
	invite.AppendHeader(&subject)

	logger.If(s.ctx, "Invite request: %s", invite.String())
	inviteRes, err := s.srv.RequestWithContext(s.ctx, invite)
	if err != nil {
		return errors.Wrapf(err, "invite request error")
	}
	logger.If(s.ctx, "Invite response: %s", inviteRes.String())

	s.srv.Send(sip.NewAckRequest("", invite, inviteRes, "", nil))

	return nil
}

func (s *GB28181Server) Catalog(deviceID string) error {
	var CatalogXML = `<?xml version="1.0"?><Query>
	<CmdType>Catalog</CmdType>
	<SN>%d</SN>
	<DeviceID>%s</DeviceID>
	</Query>
	`

	d, ok := dm.GetDevice(deviceID)
	if !ok {
		return errors.Errorf("device %s not found", deviceID)
	}

	request, err := s.CreateRequest(sip.MESSAGE, deviceID, d.SourceAddr, d.NetworkType)
	if err != nil {
		return errors.Wrapf(err, "create catalog request error")
	}

	expires := sip.Expires(3600)
	contentType := sip.ContentType("Application/MANSCDP+xml")

	request.AppendHeader(&contentType)
	request.AppendHeader(&expires)

	body := fmt.Sprintf(CatalogXML, s.SN, deviceID)
	s.SN++
	request.SetBody(body, true)
	resp, err := s.srv.RequestWithContext(context.Background(), request)
	if err != nil && resp != nil {
		return errors.Wrapf(err, "catalog request error")
	}
	return nil

}

func (s *GB28181Server) CreateRequest(Method sip.RequestMethod, deviceID, sourceAddr, networkType string) (req sip.Request, err error) {
	callId := sip.CallID(utils.GenRandomNumber(10))
	userAgent := sip.UserAgentHeader("SRS")
	maxForwards := sip.MaxForwards(70) //增加max-forwards为默认值 70
	s.SN++
	cseq := sip.CSeq{
		SeqNo:      s.SN,
		MethodName: Method,
	}
	port := sip.Port(s.conf.SipPort)
	serverAddr := sip.Address{
		Uri: &sip.SipUri{
			FUser: sip.String{Str: s.conf.Serial},
			FHost: s.conf.SipHost,
			FPort: &port,
		},
		Params: sip.NewParams().Add("tag", sip.String{Str: utils.GenRandomNumber(9)}),
	}

	//非同一域的目标地址需要使用@host
	host := s.conf.Realm
	h := deviceID[0:10]
	if h != host {
		host = sourceAddr
	}

	channelAddr := sip.Address{
		Uri: &sip.SipUri{FUser: sip.String{Str: deviceID}, FHost: host},
	}
	req = sip.NewRequest(
		"",
		Method,
		channelAddr.Uri,
		"SIP/2.0",
		[]sip.Header{
			serverAddr.AsFromHeader(),
			channelAddr.AsToHeader(),
			&callId,
			&userAgent,
			&cseq,
			&maxForwards,
			serverAddr.AsContactHeader(),
		},
		"",
		nil,
	)

	req.SetTransport(networkType)
	req.SetDestination(sourceAddr)
	return req, nil
}
