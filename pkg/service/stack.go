package service

import (
	"net"
	"strconv"

	"github.com/emiago/sipgo/sip"
)

type sipOutboundConfig struct {
	transport string
	from      string
	to        string
}

func newRequest(method sip.RequestMethod, body []byte, conf sipOutboundConfig) (*sip.Request, error) {
	dest := conf.to
	to := sip.Uri{Host: conf.to}
	if addr, sport, err := net.SplitHostPort(conf.to); err == nil {
		if port, err := strconv.Atoi(sport); err == nil {
			to.Host = addr
			to.Port = port
			dest = conf.to
		}
	}
	from := &sip.Uri{Host: conf.from}

	fromHeader := &sip.FromHeader{Address: *from, DisplayName: conf.from, Params: sip.NewParams()}
	fromHeader.Params.Add("tag", sip.GenerateTagN(16))

	req := sip.NewRequest(method, to)
	req.AppendHeader(fromHeader)
	req.AppendHeader(&sip.ToHeader{Address: to})
	req.AppendHeader(&sip.ContactHeader{Address: *from})
	req.SetBody(body)
	req.SetDestination(dest)
	req.SetTransport(conf.transport)

	return req, nil
}

func newInviteRequest(body []byte, conf sipOutboundConfig) (*sip.Request, error) {
	req, err := newRequest(sip.INVITE, body, conf)
	if err != nil {
		return nil, err
	}
	req.AppendHeader(sip.NewHeader("Content-Type", "application/sdp"))

	return req, nil
}

func newCatelogRequest(body []byte, conf sipOutboundConfig) (*sip.Request, error) {
	req, err := newRequest(sip.MESSAGE, body, conf)
	if err != nil {
		return nil, err
	}
	req.AppendHeader(sip.NewHeader("Content-Type", "Application/MANSCDP+xml"))

	return req, nil
}
