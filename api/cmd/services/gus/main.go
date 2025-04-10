package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gradientsearch/gus/foundation/logger"
)

var build = "develop"

func main() {
	l := logger.New(os.Stdout, logger.LevelInfo, "GUS", nil)
	if err := run(context.Background(), l); err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	log.Info(ctx, "starting service", "build", build)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGINT)

	sig := <-shutdown

	log.Info(ctx, "shutdown", "status", "shutdown complete", "sig", sig)

	return nil
}
