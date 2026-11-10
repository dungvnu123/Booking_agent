package main

import (
	"abclinic.local/backend/internal/app"
	"abclinic.local/backend/internal/config"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		slog.Error("ABClinic worker stopped", "error", err)
		os.Exit(1)
	}
}
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	application, err := app.New(ctx, cfg)
	if err != nil {
		return err
	}
	defer application.Close()
	return application.RunWorker(ctx)
}
