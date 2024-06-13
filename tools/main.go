package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-bench/gb28181"
)

func main() {
	ctx := context.Background()

	var conf interface{}
	conf = gb28181.Parse(ctx)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
		for sig := range sigs {
			logger.Wf(ctx, "Quit for signal %v", sig)
			cancel()
		}
	}()

	gb28181.Run(ctx, conf)
}
