package framework

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/booking-service/internal/app"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/booking-service/internal/router/handler/booking"
)

type Registry struct {
	app *app.Application
}

func NewHttpRegistry(app *app.Application) *Registry {
	return &Registry{
		app,
	}
}

func (r *Registry) RegisterHttpRoutes(rootCtx context.Context, router *http.ServeMux) {

	booking.RegisterRoutes(router, r.app)
	slog.Info("routes registered")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			return rootCtx
		},
	}

	go func() {
		slog.Info("starting http server")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Http Server", "Error", err)
		}
	}()
}
