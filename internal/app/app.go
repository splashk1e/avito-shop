package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/splashk1e/avito-shop/internal/handlers"
	"github.com/splashk1e/avito-shop/internal/server"
	"github.com/splashk1e/avito-shop/internal/services"
	"github.com/splashk1e/avito-shop/internal/storage/postgres"
)

type App struct {
	log    *slog.Logger
	port   string
	server *server.Server
}

func New(log *slog.Logger, port string) *App {
	return &App{
		log:    log,
		port:   port,
		server: new(server.Server),
	}
}
func (a *App) Run(storagePath, secret string, tokenTTL time.Duration) error {
	storage := postgres.New(context.Background(), storagePath)
	authService := services.NewAuthService(a.log, storage, storage, tokenTTL, secret)
	transacService := services.NewTransactionsService(storage, storage, a.log)
	handlers := handlers.NewHandler(authService, transacService)

	return a.server.Run(a.port, handlers.InitRoutes())
}

func (a *App) Stop() {
	a.server.Shutdown(context.Background())
}
