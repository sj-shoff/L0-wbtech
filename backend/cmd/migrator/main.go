package main

import (
	"L0-wbtech/internal/config"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()

	migrationsPath := flag.String("migrations-path", "", "Path to migrations")
	flag.Parse()

	if *migrationsPath == "" {
		*migrationsPath = cfg.Migrations
	}

	if *migrationsPath == "" {
		slog.Error("Migrations path is required")
		os.Exit(1)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	m, err := migrate.New(
		"file://"+*migrationsPath,
		dsn,
	)
	if err != nil {
		slog.Error("Migration initialization failed", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("No new migrations to apply")
			return
		}
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Migrations applied successfully")
}
