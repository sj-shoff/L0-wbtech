package main

import (
	"L0-wbtech/internal/app"
	"L0-wbtech/internal/cache"
	"L0-wbtech/internal/config"
	"L0-wbtech/internal/service"
	"L0-wbtech/internal/storage/postgres"
	"L0-wbtech/pkg/logger/sl"
	"L0-wbtech/pkg/logger/slogsetup"
	"context"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := slogsetup.SetupLogger(cfg.Env)
	log.Info("Starting server", "env", cfg.Env)

	storage, err := postgres.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Error("Failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	orderCache := cache.NewCache()

	orderService := service.NewOrderService(storage, orderCache, log)

	if err := orderService.RestoreCache(context.Background()); err != nil {
		log.Error("Failed to restore cache", sl.Err(err))
	}

	consumer := app.NewKafkaConsumer(cfg, orderService, log)

	application := app.New(cfg, orderService, consumer, log)
	application.Run()
}
