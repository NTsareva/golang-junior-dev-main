package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"gitlab.com/llcmediatel/recruiting/golang-junior-dev/internal/config"
	"gitlab.com/llcmediatel/recruiting/golang-junior-dev/internal/server/handlers"
	"gitlab.com/llcmediatel/recruiting/golang-junior-dev/internal/server/middleware"
)

var configPath = flag.String("config", "./configs/config.toml", "Path to config file")

func main() {
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logrus.Error("parse and load config failed: ", err)
	}

	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		logrus.Error("failed to parse log level: ", err)
	}

	logrus.SetLevel(level)

	logrus.Infof(fmt.Sprintf("%s:%s", cfg.Server.Hostname, cfg.Server.Port))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.LogginMiddleware(http.HandlerFunc(handlers.ExchangeHandler)))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Hostname, cfg.Server.Port),
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logrus.Infof("starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("failed to listen to serer: %v", err)
		}
	}()

	<-quit
	logrus.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatalf("failed to shutdown server: %v", err)
	}
}
