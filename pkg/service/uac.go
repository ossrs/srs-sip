package service

import (
	"context"
	"fmt"

	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-sip/pkg/config"
	"github.com/ossrs/srs-sip/pkg/service/stack"
)

const (
	UserAgent = "SRS-SIP/1.0"
)

type UAC struct {
	*Cascade

	SN      uint32
	LocalIP string
}

func NewUac() *UAC {
	ip, err := config.GetLocalIP()
	if err != nil {
		logger.E("get local ip failed")
		return nil
	}

	c := &UAC{
		Cascade: &Cascade{},
		LocalIP: ip,
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

	c.sipCli, err = sipgo.NewClient(agent, sipgo.WithClientHostname(c.LocalIP))
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
	r, _ := stack.NewRegisterRequest(stack.OutboundConfig{
		From:      "34020000001110000001",
		To:        "34020000002000000001",
		Transport: "UDP",
		Via:       fmt.Sprintf("%s:%d", c.LocalIP, c.conf.SipPort),
	})
	tx, err := c.sipCli.TransactionRequest(c.ctx, r)
	if err != nil {
		return errors.Wrapf(err, "transaction request error")
	}

	rs, _ := c.getResponse(tx)
	logger.Tf(c.ctx, "register response: %s", rs.String())
	return nil
}

func (c *UAC) OnRequest(req *sip.Request, tx sip.ServerTransaction) {
	switch req.Method {
	case "INVITE":
		c.onInvite(req, tx)
	}
}

func (c *UAC) onInvite(req *sip.Request, tx sip.ServerTransaction) {
	logger.T(c.ctx, "onInvite")
}

func (c *UAC) onBye(req *sip.Request, tx sip.ServerTransaction) {
	logger.T(c.ctx, "onBye")
}

func (c *UAC) onMessage(req *sip.Request, tx sip.ServerTransaction) {
	logger.Tf(c.ctx, "onMessage %s", req.String())
}

func (c *UAC) getResponse(tx sip.ClientTransaction) (*sip.Response, error) {
	select {
	case <-tx.Done():
		return nil, fmt.Errorf("transaction died")
	case res := <-tx.Responses():
		return res, nil
	}
}
