package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	"github.com/ossrs/go-oryx-lib/errors"
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
		slog.Error("get local ip failed", "error", err)
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
		Via:       fmt.Sprintf("%s:%d", c.LocalIP, c.conf.GB28181.Port),
	})
	tx, err := c.sipCli.TransactionRequest(c.ctx, r)
	if err != nil {
		return errors.Wrapf(err, "transaction request error")
	}

	rs, _ := c.getResponse(tx)
	slog.Info("register response", "response", rs.String())
	return nil
}

func (c *UAC) OnRequest(req *sip.Request, tx sip.ServerTransaction) {
	switch req.Method {
	case "INVITE":
		c.onInvite(req, tx)
	}
}

func (c *UAC) onInvite(req *sip.Request, tx sip.ServerTransaction) {
	slog.Debug("onInvite")
}

func (c *UAC) onBye(req *sip.Request, tx sip.ServerTransaction) {
	slog.Debug("onBye")
}

func (c *UAC) onMessage(req *sip.Request, tx sip.ServerTransaction) {
	slog.Debug("onMessage", "request", req.String())
}

func (c *UAC) getResponse(tx sip.ClientTransaction) (*sip.Response, error) {
	select {
	case <-tx.Done():
		return nil, fmt.Errorf("transaction died")
	case res := <-tx.Responses():
		return res, nil
	}
}
