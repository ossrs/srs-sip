package service

import (
	"context"
	"net/http"
	"time"

	"github.com/emiago/sipgo"
	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-sip/pkg/config"
)

type Service struct {
	ctx  context.Context
	conf *config.MainConfig
	uac  *UAC
	uas  *UAS
}

func NewService(ctx context.Context, r0 interface{}) (*Service, error) {
	s := &Service{
		ctx:  ctx,
		conf: r0.(*config.MainConfig),
	}
	s.uac = NewUac()
	s.uas = NewUas()
	return s, nil
}

func (s *Service) Start() error {
	ua, err := sipgo.NewUA(
		sipgo.WithUserAgent(UserAgent),
	)
	if err != nil {
		return err
	}

	// if err := s.uac.Start(ua); err != nil {
	// 	return err
	// }
	if err := s.uas.Start(ua, s.conf); err != nil {
		return err
	}

	go func() {
		httpPort := "8080"
		server := &http.Server{
			Addr:              ":" + httpPort,
			Handler:           http.FileServer(http.Dir("../web/html")),
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Ef("listen on %s failed", httpPort)
		}
	}()

	return nil
}

func (s *Service) Stop() {
	s.uac.Stop()
	s.uas.Stop()
}
