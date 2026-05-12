package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"izanr.com/chat/config"
	sql1 "izanr.com/chat/sql"
)

func main() {
	start := time.Now()
	cfg := config.Get()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	app := fiber.New()

	go func() {
		sig := <-ch
		cancel()

		slog.Warn("Cancellation received", "signal", sig.String())

		start := time.Now()
		err := app.ShutdownWithTimeout(3 * time.Second)
		if err != nil {
			slog.Error("Failed to shutdown app", "error", err)
		} else {
			slog.Info(
				"App shut down",
				"took", time.Since(start).Round(time.Millisecond),
			)
		}
	}()

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		slog.Info(
			"HTTP: Listening",
			"addr", cfg.GetAddr(),
			"prefork", cfg.App.Prefork,
			"took", time.Since(start).Round(time.Microsecond),
		)
		return nil
	})

	pool, err := pgxpool.New(ctx, cfg.DB.URL())
	if err != nil {
		log.Fatalf("Failed to open postgres DB: %s\n", err)
	}

	if cfg.DB.Migrate {
		err := sql1.Migrate(ctx, pool)
		if err != nil {
			log.Fatalf("Failed to migrate DB: %s\n", err)
		}
	}

	app.Listen(cfg.GetAddr(), fiber.ListenConfig{
		EnablePrefork:         cfg.App.Prefork,
		DisableStartupMessage: true,
	})
}
