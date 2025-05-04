package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/booking-service/internal/adapter/pubsub/consumer"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/booking-service/internal/app"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/booking-service/internal/router/framework"
)

func main() {

	// Init root ctx
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	defer rootCtxCancel()

	newApp := app.NewApp()
	mux := http.NewServeMux()

	consumer.InitPubSubConsumers(rootCtx, newApp)
	slog.Info("consumers initialized")
	framework.NewHttpRegistry(newApp).RegisterHttpRoutes(rootCtx, mux)

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
