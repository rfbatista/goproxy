package server

import (
	"context"
	"fmt"
	"goproxy/internal/infrastructure/config"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHTTPServer(
	lc fx.Lifecycle,
	mux *http.ServeMux,
	log *zap.Logger,
	cfg config.AppConfig,
) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%s", "8080"), Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("starting http server at", zap.String("port", srv.Addr))
			go srv.Serve(ln)
			log.Info("proxying requests to", zap.String("url", cfg.BackendURL))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
