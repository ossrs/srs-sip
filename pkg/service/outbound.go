package service

import (
	"fmt"
	"strings"

	"github.com/emiago/sipgo/sip"
	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-sip/pkg/service/stack"
	"github.com/ossrs/srs-sip/pkg/utils"
)

func (s *UAS) Invite(mediaServerAddr, deviceID, channelID string) error {
	if s.isPublishing(channelID) {
		return nil
	}

	ssrc := utils.CreateSSRC(true)

	// 更新s.signal.Addr
	s.signal.SetAddr("http://" + mediaServerAddr)

	mediaPort, err := s.signal.Publish(ssrc, ssrc)
	if err != nil {
		return errors.Wrapf(err, "api gb publish request error")
	}

	mediaHost := strings.Split(mediaServerAddr, ":")[0]
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
		"\r\n",
	}
	if true { // support tcp only
		sdpInfo = append(sdpInfo, "a=setup:passive", "a=connection:new")
	}

	// TODO: 需要考虑不同设备，通道ID相同的情况
	d, ok := DM.GetDeviceInfoByChannel(channelID)
	if !ok {
		return errors.Errorf("device not found by %s", channelID)
	}

	subject := fmt.Sprintf("%s:%s,%s:0", channelID, ssrc, s.conf.Serial)

	req, err := stack.NewInviteRequest([]byte(strings.Join(sdpInfo, "\r\n")), subject, stack.OutboundConfig{
		Via:       d.SourceAddr,
		To:        d.DeviceID,
		From:      s.conf.Serial,
		Transport: d.NetworkType,
	})
	if err != nil {
		return errors.Wrapf(err, "build invite request error")
	}
	tx, err := s.sipCli.TransactionRequest(s.ctx, req)
	if err != nil {
		return errors.Wrapf(err, "transaction request error")
	}

	res, err := s.waitAnswer(tx)
	if err != nil {
		return errors.Wrapf(err, "wait answer error")
	}
	if res.StatusCode != 200 {
		return errors.Errorf("invite response error: %s", res.String())
	}

	ack := sip.NewAckRequest(req, res, nil)
	s.sipCli.WriteRequest(ack)

	s.AddVideoChannelStatue(channelID, VideoChannelStatus{
		ID:        channelID,
		ParentID:  deviceID,
		MediaHost: mediaHost,
		MediaPort: mediaPort,
		Ssrc:      ssrc,
		Status:    "ON",
	})

	return nil
}

func (s *UAS) isPublishing(channelID string) bool {
	c, err := s.GetVideoChannelStatue(channelID)
	if !err {
		return false
	}

	if p, err := s.signal.GetStreamStatus(c.Ssrc); err != nil || !p {
		return false
	}
	return true
}

func (s *UAS) Bye() error {
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
		From:      s.conf.Serial,
		Transport: d.NetworkType,
	})
	if err != nil {
		return errors.Wrapf(err, "build catalog request error")
	}

	tx, err := s.sipCli.TransactionRequest(s.ctx, req)
	if err != nil {
		return errors.Wrapf(err, "transaction request error")
	}

	res, err := s.waitAnswer(tx)
	if err != nil {
		return errors.Wrapf(err, "wait answer error")
	}
	logger.Tf(s.ctx, "catalog response: %s", res.String())

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
		From:      s.conf.Serial,
		Transport: d.NetworkType,
	})
	if err != nil {
		return errors.Wrapf(err, "build ptz request error")
	}

	tx, err := s.sipCli.TransactionRequest(s.ctx, req)
	if err != nil {
		return errors.Wrapf(err, "transaction request error")
	}

	res, err := s.waitAnswer(tx)
	if err != nil {
		return errors.Wrapf(err, "wait answer error")
	}

	if res.StatusCode != 200 {
		return errors.Errorf("ptz response error: %s", res.String())
	}

	return nil
}
