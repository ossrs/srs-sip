package service

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/emiago/sipgo/sip"
	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/srs-sip/pkg/media"
	"github.com/ossrs/srs-sip/pkg/models"
	"github.com/ossrs/srs-sip/pkg/service/stack"
	"github.com/ossrs/srs-sip/pkg/utils"
)

type Session struct {
	ID        string
	ParentID  string
	MediaHost string
	MediaPort int
	Ssrc      string
	Status    string
	URL       string
	RefCount  int
	CSeq      int
	InviteReq *sip.Request
	InviteRes *sip.Response
}

func (s *Session) NewRequest(method sip.RequestMethod, body []byte) *sip.Request {
	request := sip.NewRequest(method, s.InviteReq.Recipient)

	request.SipVersion = s.InviteRes.SipVersion

	maxForwardsHeader := sip.MaxForwardsHeader(70)
	request.AppendHeader(&maxForwardsHeader)

	if h := s.InviteReq.From(); h != nil {
		request.AppendHeader(h)
	}

	if h := s.InviteRes.To(); h != nil {
		request.AppendHeader(h)
	}

	if h := s.InviteReq.CallID(); h != nil {
		request.AppendHeader(h)
	}

	if h := s.InviteReq.CSeq(); h != nil {
		h.SeqNo++
		request.AppendHeader(h)
	}

	request.SetSource(s.InviteReq.Source())
	request.SetDestination(s.InviteReq.Destination())
	request.SetTransport(s.InviteReq.Transport())
	request.SetBody(body)

	s.CSeq++

	return request
}

func (s *Session) NewByeRequest() *sip.Request {
	return s.NewRequest(sip.BYE, nil)
}

// PAUSE RTSP/1.0
// CSeq:1
// PauseTime:now
func (s *Session) NewPauseRequest() *sip.Request {
	body := []byte(fmt.Sprintf(`PAUSE RTSP/1.0
CSeq: %d
PauseTime: now
`, s.CSeq))
	s.CSeq++
	pauseRequest := s.NewRequest(sip.INFO, body)
	pauseRequest.AppendHeader(sip.NewHeader("Content-Type", "Application/MANSRTSP"))
	return pauseRequest
}

// PLAY RTSP/1.0
// CSeq:2
// Range:npt=now
func (s *Session) NewResumeRequest() *sip.Request {
	body := []byte(fmt.Sprintf(`PLAY RTSP/1.0
CSeq: %d
Range: npt=now
`, s.CSeq))
	s.CSeq++
	resumeRequest := s.NewRequest(sip.INFO, body)
	resumeRequest.AppendHeader(sip.NewHeader("Content-Type", "Application/MANSRTSP"))
	return resumeRequest
}

// PLAY RTSP/1.0
// CSeq:3
// Scale:2.0
func (s *Session) NewSpeedRequest(speed float32) *sip.Request {
	body := []byte(fmt.Sprintf(`PLAY RTSP/1.0
CSeq: %d
Scale: %.1f
`, s.CSeq, speed))
	s.CSeq++
	speedRequest := s.NewRequest(sip.INFO, body)
	speedRequest.AppendHeader(sip.NewHeader("Content-Type", "Application/MANSRTSP"))
	return speedRequest
}

func (s *UAS) AddSession(key string, status Session) {
	slog.Info("AddSession", "key", key, "status", status)
	s.Streams.Store(key, status)
}

func (s *UAS) GetSession(key string) (Session, bool) {
	v, ok := s.Streams.Load(key)
	if !ok {
		return Session{}, false
	}
	return v.(Session), true
}

func (s *UAS) GetSessionByURL(url string) (string, Session) {
	var k string
	var result Session
	s.Streams.Range(func(key, value interface{}) bool {
		stream := value.(Session)
		if stream.URL == url {
			k = key.(string)
			result = stream
			return false // break
		}
		return true // continue
	})
	return k, result
}

func (s *UAS) RemoveSession(key string) {
	s.Streams.Delete(key)
}

func (s *UAS) InitMediaServer(req models.InviteRequest) error {
	s.mediaLock.Lock()
	defer s.mediaLock.Unlock()

	mediaServer, err := MediaDB.GetMediaServer(req.MediaServerId)
	if err != nil {
		return errors.Wrapf(err, "get media server error")
	}

	if s.media != nil && s.media.GetAddr() == fmt.Sprintf("%s:%d", mediaServer.IP, mediaServer.Port) {
		return nil
	}

	switch mediaServer.Type {
	case "SRS", "srs":
		s.media = &media.Srs{
			Ctx:      s.ctx,
			Schema:   "http",
			Addr:     fmt.Sprintf("%s:%d", mediaServer.IP, mediaServer.Port),
			Username: mediaServer.Username,
			Password: mediaServer.Password,
		}
	case "ZLM", "zlm":
		s.media = &media.Zlm{
			Ctx:    s.ctx,
			Schema: "http",
			Addr:   fmt.Sprintf("%s:%d", mediaServer.IP, mediaServer.Port),
			Secret: mediaServer.Secret,
		}
	default:
		return errors.Errorf("unsupported media server type: %s", mediaServer.Type)
	}

	return nil
}

func (s *UAS) handleSipTransaction(req *sip.Request) (*sip.Response, error) {
	tx, err := s.sipCli.TransactionRequest(s.ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "transaction request error")
	}

	res, err := s.waitAnswer(tx)
	if err != nil {
		return nil, errors.Wrapf(err, "wait answer error")
	}
	if res.StatusCode != 200 {
		return nil, errors.Errorf("response error: %s", res.String())
	}

	return res, nil
}

func (s *UAS) isPublishing(key string) bool {
	c, ok := s.GetSession(key)
	if !ok {
		return false
	}

	// Check if stream already exists
	if p, err := s.media.GetStreamStatus(c.Ssrc); err != nil || !p {
		return false
	}

	return true
}

func (s *UAS) Invite(req models.InviteRequest) (*Session, error) {
	key := fmt.Sprintf("%d:%s:%s:%d:%d:%d:%d", req.MediaServerId, req.DeviceID, req.ChannelID, req.SubStream, req.PlayType, req.StartTime, req.EndTime)

	// Check if stream already exists
	if s.isPublishing(key) {
		// Stream exists, increase reference count
		c, _ := s.GetSession(key)
		c.RefCount++
		s.AddSession(key, c)
		return &c, nil
	}

	ssrc := utils.CreateSSRC(req.PlayType == 0)

	err := s.InitMediaServer(req)
	if err != nil {
		return nil, errors.Wrapf(err, "init media server error")
	}

	mediaPort, err := s.media.Publish(ssrc, ssrc)
	if err != nil {
		return nil, errors.Wrapf(err, "api gb publish request error")
	}

	mediaHost := strings.Split(s.media.GetAddr(), ":")[0]
	if mediaHost == "" {
		return nil, errors.Errorf("media host is empty")
	}

	sessionName := utils.GetSessionName(req.PlayType)

	sdpInfo := []string{
		"v=0",
		fmt.Sprintf("o=%s 0 0 IN IP4 %s", req.ChannelID, mediaHost),
		"s=" + sessionName,
		"u=" + req.ChannelID + ":0",
		"c=IN IP4 " + mediaHost,
		"t=" + fmt.Sprintf("%d %d", req.StartTime, req.EndTime),
		fmt.Sprintf("m=video %d TCP/RTP/AVP 96", mediaPort),
		"a=recvonly",
		"a=rtpmap:96 PS/90000",
		"y=" + ssrc,
		"\r\n",
	}
	if true { // support tcp only
		sdpInfo = append(sdpInfo, "a=setup:passive", "a=connection:new")
	}

	// TODO: 需要考虑不同设备，通道ID相同的情况
	d, ok := DM.GetDeviceInfoByChannel(req.ChannelID)
	if !ok {
		return nil, errors.Errorf("device not found by %s", req.ChannelID)
	}

	subject := fmt.Sprintf("%s:%s,%s:0", req.ChannelID, ssrc, s.conf.GB28181.Serial)

	reqInvite, err := stack.NewInviteRequest([]byte(strings.Join(sdpInfo, "\r\n")), subject, stack.OutboundConfig{
		Via:       d.SourceAddr,
		To:        d.DeviceID,
		From:      s.conf.GB28181.Serial,
		Transport: d.NetworkType,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "build invite request error")
	}

	res, err := s.handleSipTransaction(reqInvite)
	if err != nil {
		return nil, err
	}

	ack := sip.NewAckRequest(reqInvite, res, nil)
	s.sipCli.WriteRequest(ack)

	session := Session{
		ID:        req.ChannelID,
		ParentID:  req.DeviceID,
		MediaHost: mediaHost,
		MediaPort: mediaPort,
		Ssrc:      ssrc,
		Status:    "ON",
		URL:       s.media.GetWebRTCAddr(ssrc),
		RefCount:  1,
		InviteReq: reqInvite,
		InviteRes: res,
	}

	s.AddSession(key, session)

	return &session, nil
}

func (s *UAS) Bye(req models.ByeRequest) error {
	key, session := s.GetSessionByURL(req.URL)
	if key == "" {
		return errors.Errorf("stream not found: %s", req.URL)
	}

	session.RefCount--
	if session.RefCount > 0 {
		s.AddSession(key, session)
		return nil
	}

	defer func() {
		if err := s.media.Unpublish(session.Ssrc); err != nil {
			slog.Error("unpublish stream error", "error", err, "stream_id", session.Ssrc)
		}
		s.RemoveSession(key)
	}()

	reqBye := session.NewByeRequest()
	_, err := s.handleSipTransaction(reqBye)
	if err != nil {
		return err
	}

	return nil
}

func (s *UAS) Pause(req models.PauseRequest) error {
	key, session := s.GetSessionByURL(req.URL)
	if key == "" {
		return errors.Errorf("stream not found: %s", req.URL)
	}

	pauseRequest := session.NewPauseRequest()
	_, err := s.handleSipTransaction(pauseRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *UAS) Resume(req models.ResumeRequest) error {
	key, session := s.GetSessionByURL(req.URL)
	if key == "" {
		return errors.Errorf("stream not found: %s", req.URL)
	}

	resumeRequest := session.NewResumeRequest()
	_, err := s.handleSipTransaction(resumeRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *UAS) Speed(req models.SpeedRequest) error {
	key, session := s.GetSessionByURL(req.URL)
	if key == "" {
		return errors.Errorf("stream not found: %s", req.URL)
	}

	speedRequest := session.NewSpeedRequest(req.Speed)
	_, err := s.handleSipTransaction(speedRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *UAS) Catalog(deviceID string) error {
	var CatalogXML = `<?xml version="1.0"?><Query>
	<CmdType>Catalog</CmdType>
	<SN>%d</SN>
	<DeviceID>%s</DeviceID>
	</Query>
	`

	d, ok := DM.GetDevice(deviceID)
	if !ok {
		return errors.Errorf("device %s not found", deviceID)
	}

	body := fmt.Sprintf(CatalogXML, s.getSN(), deviceID)

	req, err := stack.NewMessageRequest([]byte(body), stack.OutboundConfig{
		Via:       d.SourceAddr,
		To:        d.DeviceID,
		From:      s.conf.GB28181.Serial,
		Transport: d.NetworkType,
	})
	if err != nil {
		return errors.Wrapf(err, "build catalog request error")
	}

	_, err = s.handleSipTransaction(req)
	if err != nil {
		return err
	}

	return nil
}

func (s *UAS) waitAnswer(tx sip.ClientTransaction) (*sip.Response, error) {
	select {
	case <-s.ctx.Done():
		return nil, errors.Errorf("context done")
	case res := <-tx.Responses():
		if res.StatusCode == 100 || res.StatusCode == 101 || res.StatusCode == 180 || res.StatusCode == 183 {
			return s.waitAnswer(tx)
		}
		return res, nil
	}
}

// <?xml version="1.0"?>
// <Control>
// <CmdType>DeviceControl</CmdType>
// <SN>474</SN>
// <DeviceID>33010602001310019325</DeviceID>
// <PTZCmd>a50f4d0190000092</PTZCmd>
// <Info>
// <ControlPriority>150</ControlPriority>
// </Info>
// </Control>
func (s *UAS) ControlPTZ(deviceID, channelID, ptz, speed string) error {
	var ptzXML = `<?xml version="1.0"?>
	<Control>
	<CmdType>DeviceControl</CmdType>
	<SN>%d</SN>
	<DeviceID>%s</DeviceID>
	<PTZCmd>%s</PTZCmd>
	<Info>
	<ControlPriority>150</ControlPriority>
	</Info>
	</Control>
	`

	// d, ok := DM.GetDevice(deviceID)
	d, ok := DM.GetDeviceInfoByChannel(channelID)
	if !ok {
		return errors.Errorf("device %s not found", deviceID)
	}

	ptzCmd, err := toPTZCmd(ptz, speed)
	if err != nil {
		return errors.Wrapf(err, "build ptz command error")
	}

	body := fmt.Sprintf(ptzXML, s.getSN(), channelID, ptzCmd)

	req, err := stack.NewMessageRequest([]byte(body), stack.OutboundConfig{
		Via:       d.SourceAddr,
		To:        d.DeviceID,
		From:      s.conf.GB28181.Serial,
		Transport: d.NetworkType,
	})
	if err != nil {
		return errors.Wrapf(err, "build ptz request error")
	}

	_, err = s.handleSipTransaction(req)
	return err
}

// QueryRecord 查询录像记录
func (s *UAS) QueryRecord(deviceID, channelID string, startTime, endTime int64) ([]*models.Record, error) {
	var queryXML = `<?xml version="1.0"?>
	<Query>
	<CmdType>RecordInfo</CmdType>
	<SN>%d</SN>
	<DeviceID>%s</DeviceID>
	<StartTime>%s</StartTime>
	<EndTime>%s</EndTime>
	<Secrecy>0</Secrecy>
	<Type>all</Type>
	</Query>
	`

	d, ok := DM.GetDeviceInfoByChannel(channelID)
	if !ok {
		return nil, errors.Errorf("device %s not found", deviceID)
	}

	// 时间原本是unix时间戳，需要转换为YYYY-MM-DDTHH:MM:SS
	startTimeStr := time.Unix(startTime, 0).Format("2006-01-02T00:00:00")
	endTimeStr := time.Unix(endTime, 0).Format("2006-01-02T15:04:05")

	body := fmt.Sprintf(queryXML, s.getSN(), channelID, startTimeStr, endTimeStr)

	req, err := stack.NewMessageRequest([]byte(body), stack.OutboundConfig{
		Via:       d.SourceAddr,
		To:        d.DeviceID,
		From:      s.conf.GB28181.Serial,
		Transport: d.NetworkType,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "build query request error")
	}

	if _, err := s.handleSipTransaction(req); err != nil {
		return nil, err
	}

	// 创建一个通道来接收录像查询结果
	resultChan := make(chan *models.XmlMessageInfo, 1)
	s.recordQueryResults.Store(channelID, resultChan)
	defer s.recordQueryResults.Delete(channelID)

	// 等待结果或超时
	var allRecords []*models.Record
	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-timeout:
			return allRecords, errors.Errorf("query record timeout after 30s")
		case <-s.ctx.Done():
			return nil, errors.Errorf("context done")
		case records := <-resultChan:
			allRecords = append(allRecords, records.RecordList...)
			slog.Info("Record query result",
				"channel", channelID,
				"expected_count", records.SumNum,
				"actual_count", len(allRecords),
				"batch_count", len(records.RecordList))

			if len(allRecords) == records.SumNum {
				return allRecords, nil
			}
		}
	}
}

// ConfigDownload
// <?xml version="1.0"?>
// <Control>
// <CmdType>ConfigDownload</CmdType>
// <SN>474</SN>
// <DeviceID>33010602001310019325</DeviceID>
// </Control>

func (s *UAS) ConfigDownload(deviceID string) error {
	var deviceConfigXML = `<?xml version="1.0"?>
	<Control>
	<CmdType>ConfigDownload</CmdType>
	<SN>%d</SN>
	<DeviceID>%s</DeviceID>
	<ConfigType>BasicParam</ConfigType>
	</Control>
	`

	d, ok := DM.GetDevice(deviceID)
	if !ok {
		return errors.Errorf("device %s not found", deviceID)
	}

	body := fmt.Sprintf(deviceConfigXML, s.getSN(), deviceID)

	req, err := stack.NewMessageRequest([]byte(body), stack.OutboundConfig{
		Via:       d.SourceAddr,
		To:        d.DeviceID,
		From:      s.conf.GB28181.Serial,
		Transport: d.NetworkType,
	})
	if err != nil {
		return errors.Wrapf(err, "build device config request error")
	}

	_, err = s.handleSipTransaction(req)
	return err
}
