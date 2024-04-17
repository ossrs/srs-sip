package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ossrs/srs-sip/gbserver"
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

	conf := gbserver.Parse(ctx)
	gbserver.Run(ctx, conf)
}
