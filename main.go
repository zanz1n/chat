package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"izanr.com/chat/config"
	sql1 "izanr.com/chat/sql"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	cfg := config.Get()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	pool, err := setupDB(ctx)
	if err != nil {
		log.Fatalf("Setup DB failed: %s", err)
	}
	defer pool.Close()

	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")

	mux := http.NewServeMux()

	mux.Handle("POST /rpc", server)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error(
				"HTTP: Upgrade connection failed",
				"from", r.RemoteAddr,
				"error", err,
			)
			return
		}
		defer conn.Close()

		for {
			kind, rd, err := conn.NextReader()
			if err != nil {
				break
			}

			_, _ = kind, rd
		}
	})

	srv := &http.Server{
		Addr:    cfg.GetAddr(),
		Handler: mux,
	}

	go func() {
		sig := <-ch
		cancel()

		slog.Warn("Cancellation received", "signal", sig.String())

		ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel2()

		start := time.Now()
		err := srv.Shutdown(ctx2)
		if err != nil {
			slog.Error("Failed to shutdown app", "error", err)
		} else {
			slog.Info(
				"App shut down",
				"took", time.Since(start).Round(time.Millisecond),
			)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %s", err)
	}
}

func setupDB(ctx context.Context) (*pgxpool.Pool, error) {
	cfg := config.Get()

	pool, err := pgxpool.New(ctx, cfg.DB.URL())
	if err != nil {
		return nil, fmt.Errorf("open pool: %w", err)
	}

	if cfg.DB.Migrate {
		if err = sql1.Migrate(ctx, pool); err != nil {
			return nil, fmt.Errorf("migrate: %w", err)
		}
	}

	return pool, nil
}
