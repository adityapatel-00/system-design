package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/seats-service/internal/adapter/pubsub/consumer"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/seats-service/internal/app"
)

func main() {

	// Init root ctx
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	defer rootCtxCancel()

	newApp := app.NewApp()
	consumer.InitPubSubConsumers(rootCtx, newApp)

	initGracefulShutDown(rootCtxCancel)
}

func initGracefulShutDown(rootCtxCancelFunc context.CancelFunc) {
	// Wait for stop signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	slog.Debug("Service waiting for signal")
	sig := <-signals
	slog.Debug("Service received signal", slog.Any("signal", sig))

	slog.Debug("Stopping service")
	rootCtxCancelFunc()

	slog.Debug("Service stopped successfully.")
}
