package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/splashk1e/avito-shop/internal/app"
	"github.com/splashk1e/avito-shop/internal/config"
)

func main() {
	secret := os.Getenv("SECRET")
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	application := app.New(log, cfg.Port)
	fmt.Print(cfg.Storagepath)
	go application.Run(cfg.Storagepath, secret, cfg.TokenTTL)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sign := <-stop
	application.Stop()
	log.Info("application stopped", slog.String("signal", sign.String()))
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
