package sql

import (
	"context"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	start := time.Now()

	m, err := newMigrate(pool)
	if err != nil {
		return err
	}
	defer m.Close()

	version, dirty, err := _migrate(ctx, m)
	if err != nil {
		slog.Error(
			"Migrate: Failed to migrate",
			"took", time.Since(start).Round(time.Millisecond),
			"error", err,
		)
		return err
	}

	slog.Info(
		"Migrate: Successfully migrated",
		"version", version,
		"drity", dirty,
		"took", time.Since(start).Round(time.Millisecond),
	)
	return nil
}

func MigrateSkipKV(ctx context.Context, pool *pgxpool.Pool) error {
	m, err := newMigrate(pool)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Force(1); err != nil {
		return err
	}

	_, _, err = _migrate(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func _migrate(ctx context.Context, m *migrate.Migrate) (uint, bool, error) {
	finish := make(chan struct{})
	go func() {
		select {
		case _, _ = <-finish:
		case <-ctx.Done():
			m.GracefulStop <- true
		}
	}()
	defer close(finish)

	if err := m.Up(); err != nil {
		finish <- struct{}{}
		return 0, false, err
	}
	finish <- struct{}{}

	return m.Version()
}

func newMigrate(pool *pgxpool.Pool) (m *migrate.Migrate, err error) {
	var source source.Driver
	db := stdlib.OpenDBFromPool(pool)

	cfg := pool.Config().ConnConfig
	driver, err := mpgx.WithInstance(db, &mpgx.Config{
		MigrationsTable: mpgx.DefaultMigrationsTable,
		DatabaseName:    cfg.Database,
	})
	if err != nil {
		goto __err
	}

	source, err = iofs.New(Migrations, "migrations")
	if err != nil {
		goto __err
	}

	m, err = migrate.NewWithInstance("iofs", source, "pgx/v5", driver)
	if err != nil {
		goto __err
	}

	return m, nil

__err:
	db.Close()
	return nil, err
}
