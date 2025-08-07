package main

import (
	"L0-wbtech/internal/config"
	"L0-wbtech/pkg/logger/sl"
	"L0-wbtech/pkg/logger/slogsetup"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	migrationsPath := flag.String("migrations-path", "", "Path to migrations")
	migrateDown := flag.Bool("down", false, "Run migrations down (rollback)")
	flag.Parse()

	cfg := config.MustLoad()

	log := slogsetup.SetupLogger(cfg.Env)

	if *migrationsPath == "" {
		*migrationsPath = cfg.Migrations
	}

	if *migrationsPath == "" {
		log.Error("Migrations path is required")
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

	log.Info("Connecting to database", "dsn", dsn)

	m, err := migrate.New(
		"file://"+*migrationsPath,
		dsn,
	)
	if err != nil {
		log.Error("Migration initialization failed", sl.Err(err))
		os.Exit(1)
	}

	if *migrateDown {
		log.Info("Rolling back migrations", "path", *migrationsPath)
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				log.Info("No migrations to rollback")
				return
			}
			log.Error("Migration rollback failed", sl.Err(err))
			os.Exit(1)
		}
		log.Info("Migrations rolled back successfully")
	} else {
		log.Info("Applying migrations", "path", *migrationsPath)
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Info("No new migrations to apply")
				return
			}
			log.Error("Migration failed", sl.Err(err))
			os.Exit(1)
		}
		log.Info("Migrations applied successfully")
	}
}
