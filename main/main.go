package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	defer WaitTerminationSignal(cancel)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           http.FileServer(http.Dir("../web/html")),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ErrorLog:          log.New(os.Stderr, "http: ", log.LstdFlags),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	conf := utils.Parse(ctx)
	service.Run(ctx, conf)
}
