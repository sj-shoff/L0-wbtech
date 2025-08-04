package app

import (
	"L0-wbtech/internal/config"
	"L0-wbtech/internal/service"
	"net/http"
)

type App struct {
	cfg          *config.Config
	orderService *service.OrderService
	httpServer   *http.Server
}

func New(cfg *config.Config, orderService *service.OrderService) *App {
	return &App{
		cfg:          cfg,
		orderService: orderService,
	}
}
