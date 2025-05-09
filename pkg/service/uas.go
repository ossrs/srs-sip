package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"

	"github.com/emiago/sipgo"
	"github.com/ossrs/srs-sip/pkg/config"
	"github.com/ossrs/srs-sip/pkg/db"
	"github.com/ossrs/srs-sip/pkg/media"
)

type UAS struct {
	*Cascade

	SN                 uint32
	Streams            sync.Map
	mediaLock          sync.Mutex
	media              media.IMedia
	recordQueryResults sync.Map // channelID -> chan []Record

	sipConnUDP *net.UDPConn
	sipConnTCP *net.TCPListener
}

var DM = GetDeviceManager()
var MediaDB, _ = db.GetInstance("./media_servers.db")

func NewUas() *UAS {
	return &UAS{
		Cascade: &Cascade{},
	}
}

func (s *UAS) Start(agent *sipgo.UserAgent, r0 interface{}) error {
	ctx := context.Background()
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

	candidate := os.Getenv("CANDIDATE")
	if candidate != "" {
		MediaDB.AddMediaServer("Default", "SRS", candidate, 1985, "", "", "", 1)
	}
	return nil
}

func (s *UAS) startUDP() error {
	lis, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: s.conf.GB28181.Port,
	})
	if err != nil {
		return fmt.Errorf("cannot listen on the UDP signaling port %d: %w", s.conf.GB28181.Port, err)
	}
	s.sipConnUDP = lis
	slog.Info("sip signaling listening on UDP", "address", lis.LocalAddr().String(), "port", s.conf.GB28181.Port)

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
		Port: s.conf.GB28181.Port,
	})
	if err != nil {
		return fmt.Errorf("cannot listen on the TCP signaling port %d: %w", s.conf.GB28181.Port, err)
	}
	s.sipConnTCP = lis
	slog.Info("sip signaling listening on TCP", "address", lis.Addr().String(), "port", s.conf.GB28181.Port)

	go func() {
		if err := s.sipSvr.ServeTCP(lis); err != nil && !errors.Is(err, net.ErrClosed) {
			panic(fmt.Errorf("SIP listen TCP error: %w", err))
		}
	}()
	return nil
}

func (s *UAS) getSN() uint32 {
	s.SN++
	return s.SN
}
