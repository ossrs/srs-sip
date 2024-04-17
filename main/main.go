package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-sip/pkg/api"
	"github.com/ossrs/srs-sip/pkg/config"
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

	apiSvr, err := api.NewHttpApiServer(conf, sipSvr)
	if err != nil {
		logger.Ef("create http service failed. err is %v", err.Error())
		return
	}
	apiSvr.Start()

	var targetDir string
	targetDirs := []string{"./web/html", "../web/html"}
	for _, dir := range targetDirs {
		if _, err := os.Stat(path.Join(dir, "index.html")); err == nil {
			targetDir = dir
			break
		}
	}
	if targetDir == "" {
		logger.Ef(ctx, "index.html not found in %v", targetDirs)
		return
	}

	go func() {
		c := conf.(*config.MainConfig)
		httpPort := strconv.Itoa(c.HttpServerPort)
		server := &http.Server{
			Addr:              ":" + httpPort,
			Handler:           http.FileServer(http.Dir(targetDir)),
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		}
		logger.Tf(ctx, "http server listen on %s, home is %v", httpPort, targetDir)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Ef(ctx, "listen on %s failed", httpPort)
		}
	}()

	logger.Tf(ctx, "media server address is %v", conf.(*config.MainConfig).MediaAddr)

	WaitTerminationSignal(cancel)

	sipSvr.Stop()
}
