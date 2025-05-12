package service

import (
	"context"

	"github.com/emiago/sipgo"
	"github.com/ossrs/srs-sip/pkg/config"
	"github.com/rs/zerolog"
)

type Service struct {
	ctx  context.Context
	conf *config.MainConfig
	Uac  *UAC
	Uas  *UAS
}

func NewService(ctx context.Context, r0 interface{}) (*Service, error) {
	s := &Service{
		ctx:  ctx,
		conf: r0.(*config.MainConfig),
	}
	s.Uac = NewUac()
	s.Uas = NewUas()
	return s, nil
}

func (s *Service) Start(conf *config.MainConfig) error {
	if conf.Common.LogLevel != "debug" {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	ua, err := sipgo.NewUA(
		sipgo.WithUserAgent(UserAgent),
	)
	if err != nil {
		return err
	}

	if err := s.Uas.Start(ua, s.conf); err != nil {
		return err
	}

	// if err := s.Uac.Start(ua, s.conf); err != nil {
	// 	return err
	// }

	return nil
}

func (s *Service) Stop() {
	s.Uac.Stop()
	s.Uas.Stop()
	// 停止设备心跳检查器
	GetDeviceManager().stopHeartbeatChecker()
}
