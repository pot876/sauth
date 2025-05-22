package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pot876/sauth/internal/config"
	"github.com/pot876/sauth/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-exit
		cancel()
	}()

	cfg := &config.Config{}
	err := envconfig.Process(ctx, cfg)
	if err != nil {
		log.Error().Err(err).Caller().Send()
		os.Exit(1)
	}

	if err := runHttpServer(ctx, cfg); err != nil {
		os.Exit(1)
	}
}

func runHttpServer(ctx context.Context, cfg *config.Config) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	api, err := NewApi(ctx, cfg)
	if err != nil {
		return err
	}
	api.RegisterEndpoints(cfg, r)
	api.RegisterMetrics(prometheus.DefaultRegisterer)

	if cfg.HTTPEndpointMetrics != "" {
		r.GET(cfg.HTTPEndpointMetrics, metrics.MetricsHandlerGin())
	}
	server := &http.Server{
		Addr:    cfg.HTTPListenAddr,
		Handler: r.Handler(),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		log.Info().Msgf("start http on %s", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				log.Info().Msgf("%v", err)
			} else {
				log.Error().Err(err).Caller().Send()
			}
		}
	}()

	<-ctx.Done()

	shutdownContext, shutdownContextCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownContextCancel()

	log.Info().Msgf("closing http")
	err = server.Shutdown(shutdownContext)
	if err != nil {
		log.Error().Err(err).Caller().Send()
	}

	return nil
}
