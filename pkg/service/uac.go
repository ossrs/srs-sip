package service

import (
	"context"
	"fmt"

	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-sip/pkg/config"
)

const (
	UserAgent = "SRS-SIP/1.0"
)

type UAC struct {
	*Cascade

	SN uint32
}

func NewUac() *UAC {
	c := &UAC{
		Cascade: &Cascade{},
	}
	return c
}

func (c *UAC) Start(agent *sipgo.UserAgent, r0 interface{}) error {
	var err error

	c.ctx = context.Background()
	c.conf = r0.(*config.MainConfig)

	if agent == nil {
		ua, err := sipgo.NewUA(sipgo.WithUserAgent(UserAgent))
		if err != nil {
			return err
		}
		agent = ua
	}

	c.sipCli, err = sipgo.NewClient(agent, sipgo.WithClientHostname("172.18.5.44"))
	if err != nil {
		return err
	}

	c.sipSvr, err = sipgo.NewServer(agent)
	if err != nil {
		return err
	}

	c.sipSvr.OnInvite(c.onInvite)
	c.sipSvr.OnBye(c.onBye)
	c.sipSvr.OnMessage(c.onMessage)

	go c.doRegister()

	return nil
}

func (c *UAC) Stop() {
	// TODO: 断开所有当前连接
	c.sipCli.Close()
	c.sipSvr.Close()
}

func (c *UAC) doRegister() error {
	// Create basic REGISTER request structure
	username := "34020000002000000001"
	src := "172.18.5.44:5080"
	recipient := sip.Uri{
		User: username,
		Host: c.conf.SipHost,
		Port: c.conf.SipPort,
	}
	req := sip.NewRequest(sip.REGISTER, recipient)
	sipgo.ClientRequestBuild(c.sipCli, req)
	req.SetSource(src)
	req.AppendHeader(sip.NewHeader("Expires", "3600"))
	req.From().Address.User = "34020000001110000001"
	// req.AppendHeader(sip.NewHeader("From", fmt.Sprintf("<34020000001110000001@%s>", src)))

	// Send REGISTER request
	ctx := context.Background()
	tx, err := c.sipCli.TransactionRequest(ctx, req)
	if err != nil {
		return err
	}

	// Wait for response
	res, err := c.getResponse(tx)
	if err != nil {
		return err
	}
	logger.Tf("response: %s", res.String())
	return nil
}

func (c *UAC) OnRequest(req *sip.Request, tx sip.ServerTransaction) {
	switch req.Method {
	case "INVITE":
		c.onInvite(req, tx)
	}
}

func (c *UAC) onInvite(req *sip.Request, tx sip.ServerTransaction) {
	logger.T("onInvite")
}

func (c *UAC) onBye(req *sip.Request, tx sip.ServerTransaction) {
	logger.T("onBye")
}

func (c *UAC) onMessage(req *sip.Request, tx sip.ServerTransaction) {
	logger.T("onMessage")
}

func (c *UAC) getResponse(tx sip.ClientTransaction) (*sip.Response, error) {
	select {
	case <-tx.Done():
		return nil, fmt.Errorf("transaction died")
	case res := <-tx.Responses():
		return res, nil
	}
}
