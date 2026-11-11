package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"abclinic.local/backend/internal/config"
	"abclinic.local/backend/internal/platform/postgres"
	"abclinic.local/backend/migrations"
)

func main() {
	if err := run(); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("database migrations are up to date")
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := postgres.Open(ctx, cfg.DatabaseURL, cfg.DBMaxConns)
	if err != nil {
		return err
	}
	defer pool.Close()
	return migrations.Apply(ctx, pool)
}
