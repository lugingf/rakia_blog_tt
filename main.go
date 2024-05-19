package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"rakia_blog_tt/config"
	"rakia_blog_tt/handler"
	"rakia_blog_tt/service"
	"rakia_blog_tt/storage"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.New(ctx)
	if err != nil {
		slog.Error("Config initialization failed", "error", err)
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	postRepo := storage.NewInMemoryPostRepository(logger)

	application := service.New(postRepo, logger)

	hndl := handler.New(
		application, logger,
	)

	metrics := config.InitMetrics()
	go runMetricServer(cfg.Monitoring)

	server := http.Server{
		Addr:        fmt.Sprintf(":%s", cfg.App.Port),
		Handler:     handler.NewRouter(hndl, logger, metrics),
		ReadTimeout: cfg.Http.ReadTimeout,
	}

	srvErr := make(chan error)

	go func() {
		slog.Info("Listening on port", "port", cfg.App.Port)
		srvErr <- server.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		slog.Error("ListenAndServe failed", "error", err)
		return
	case <-ctx.Done():
		slog.Info("Context done signal received")
		stop()
	}

	// Create a deadline to wait for.
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	if err := server.Shutdown(ctxShutdown); err != nil {
		slog.Error("Server Shutdown Failed", "error", err)
		return
	}

	slog.Info("Server gracefully shutdown")
}

func runMetricServer(cfg *config.Monitoring) {
	mh := chi.NewRouter()
	mh.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	srv := &http.Server{
		Addr:        cfg.Port,
		Handler:     mh,
		ReadTimeout: cfg.ReadTimeout,
	}

	slog.Info(fmt.Sprintf("starting Metric exporter server: listening on %s", cfg.Port))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to listen promhandler server")
	}
}
