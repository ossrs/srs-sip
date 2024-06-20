package stack

import (
	"github.com/emiago/sipgo/sip"
)

type OutboundConfig struct {
	Transport string
	Via       string
	From      string
	To        string
}

func newRequest(method sip.RequestMethod, body []byte, conf OutboundConfig) (*sip.Request, error) {
	dest := conf.Via
	to := sip.Uri{Host: conf.To}
	from := &sip.Uri{Host: conf.From}

	fromHeader := &sip.FromHeader{Address: *from, Params: sip.NewParams()}
	fromHeader.Params.Add("tag", sip.GenerateTagN(16))

	req := sip.NewRequest(method, to)
	req.AppendHeader(fromHeader)
	req.AppendHeader(&sip.ToHeader{Address: to})
	req.AppendHeader(&sip.ContactHeader{Address: *from})
	req.AppendHeader(sip.NewHeader("Max-Forwards", "70"))
	req.SetBody(body)
	req.SetDestination(dest)
	req.SetTransport(conf.Transport)

	return req, nil
}

func NewInviteRequest(body []byte, subject string, conf OutboundConfig) (*sip.Request, error) {
	req, err := newRequest(sip.INVITE, body, conf)
	if err != nil {
		return nil, err
	}
	req.AppendHeader(sip.NewHeader("Content-Type", "application/sdp"))
	req.AppendHeader(sip.NewHeader("Subject", subject))

	return req, nil
}

func NewCatelogRequest(body []byte, conf OutboundConfig) (*sip.Request, error) {
	req, err := newRequest(sip.MESSAGE, body, conf)
	if err != nil {
		return nil, err
	}
	req.AppendHeader(sip.NewHeader("Content-Type", "Application/MANSCDP+xml"))

	return req, nil
}
