package stack

import (
	"github.com/emiago/sipgo/sip"
	"github.com/ossrs/go-oryx-lib/errors"
)

type OutboundConfig struct {
	Transport string
	Via       string
	From      string
	To        string
}

func newRequest(method sip.RequestMethod, body []byte, conf OutboundConfig) (*sip.Request, error) {
	if len(conf.From) != 20 || len(conf.To) != 20 {
		return nil, errors.Errorf("From or To length is not 20")
	}

	dest := conf.Via
	to := sip.Uri{User: conf.To, Host: conf.To[:10]}
	from := &sip.Uri{User: conf.From, Host: conf.From[:10]}

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

func NewRegisterRequest(conf OutboundConfig) (*sip.Request, error) {
	req, err := newRequest(sip.REGISTER, nil, conf)
	if err != nil {
		return nil, err
	}
	req.AppendHeader(sip.NewHeader("Expires", "3600"))

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

func NewMessageRequest(body []byte, conf OutboundConfig) (*sip.Request, error) {
	req, err := newRequest(sip.MESSAGE, body, conf)
	if err != nil {
		return nil, err
	}
	req.AppendHeader(sip.NewHeader("Content-Type", "Application/MANSCDP+xml"))

	return req, nil
}
