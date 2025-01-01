package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-sip/pkg/config"
	"github.com/ossrs/srs-sip/pkg/signaling"
)

type UAS struct {
	*Cascade

	SN             uint32
	channelsStatue sync.Map
	signal         signaling.ISignaling

	sipConnUDP *net.UDPConn
	sipConnTCP *net.TCPListener
}

var DM = GetDeviceManager()

func NewUas() *UAS {
	return &UAS{
		Cascade: &Cascade{},
	}
}

func (s *UAS) Start(agent *sipgo.UserAgent, r0 interface{}) error {
	ctx := context.Background()
	sig := &signaling.Srs{
		Ctx: ctx,
	}
	s.signal = sig
	s.startSipServer(agent, ctx, r0)
	return nil
}

func (s *UAS) Stop() {
	s.sipCli.Close()
	s.sipSvr.Close()
}

func (s *UAS) startSipServer(agent *sipgo.UserAgent, ctx context.Context, r0 interface{}) error {
	conf := r0.(*config.MainConfig)
	s.ctx = ctx
	s.conf = conf

	if agent == nil {
		ua, err := sipgo.NewUA(sipgo.WithUserAgent(UserAgent))
		if err != nil {
			return err
		}
		agent = ua
	}

	cli, err := sipgo.NewClient(agent)
	if err != nil {
		return err
	}
	s.sipCli = cli

	svr, err := sipgo.NewServer(agent)
	if err != nil {
		return err
	}
	s.sipSvr = svr

	s.sipSvr.OnRegister(s.onRegister)
	s.sipSvr.OnMessage(s.onMessage)
	s.sipSvr.OnNotify(s.onNotify)

	if err := s.startUDP(); err != nil {
		return err
	}
	if err := s.startTCP(); err != nil {
		return err
	}

	return nil
}

func (s *UAS) startUDP() error {
	lis, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: s.conf.SipPort,
	})
	if err != nil {
		return fmt.Errorf("cannot listen on the UDP signaling port %d: %w", s.conf.SipPort, err)
	}
	s.sipConnUDP = lis
	logger.Tf(s.ctx, "sip signaling listening on UDP %s:%d", lis.LocalAddr().String(), s.conf.SipPort)

	go func() {
		if err := s.sipSvr.ServeUDP(lis); err != nil {
			panic(fmt.Errorf("SIP listen UDP error: %w", err))
		}
	}()
	return nil
}

func (s *UAS) startTCP() error {
	lis, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: s.conf.SipPort,
	})
	if err != nil {
		return fmt.Errorf("cannot listen on the TCP signaling port %d: %w", s.conf.SipPort, err)
	}
	s.sipConnTCP = lis
	logger.Tf(s.ctx, "sip signaling listening on TCP %s:%d", lis.Addr().String(), s.conf.SipPort)

	go func() {
		if err := s.sipSvr.ServeTCP(lis); err != nil && !errors.Is(err, net.ErrClosed) {
			panic(fmt.Errorf("SIP listen TCP error: %w", err))
		}
	}()
	return nil
}

func sipErrorResponse(tx sip.ServerTransaction, req *sip.Request) {
	_ = tx.Respond(sip.NewResponseFromRequest(req, 400, "", nil))
}

func (s *UAS) getSN() uint32 {
	s.SN++
	return s.SN
}
