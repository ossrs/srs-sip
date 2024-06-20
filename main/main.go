package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-sip/pkg/api"
	"github.com/ossrs/srs-sip/pkg/service"
	"github.com/ossrs/srs-sip/pkg/utils"
)

func WaitTerminationSignal(cancel context.CancelFunc) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)
	<-sigc
	cancel()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	conf := utils.Parse(ctx)
	sipSvr, err := service.NewService(ctx, conf)
	if err != nil {
		logger.Ef("create service failed. err is %v", err.Error())
		return
	}

	if err := sipSvr.Start(); err != nil {
		logger.Ef("start sip service failed. err is %v", err.Error())
		return
	}

	httpSvr, err := api.NewHttpServer(conf, sipSvr)
	if err != nil {
		logger.Ef("create http service failed. err is %v", err.Error())
		return
	}
	httpSvr.Start()

	WaitTerminationSignal(cancel)

	sipSvr.Stop()
}
