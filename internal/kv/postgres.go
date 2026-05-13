package kv

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"izanr.com/chat/internal/utils"
)

type pgstorer struct {
	q      utils.Querier
	prefix string
}

func NewPgStorer(q utils.Querier) Storer {
	return &pgstorer{
		q:      q,
		prefix: "",
	}
}

// WithDomain implements [Storer].
func (s *pgstorer) WithDomain(dom string) Storer {
	return &pgstorer{
		q:      s.q,
		prefix: s.prefix + dom + ".",
	}
}

// Get implements [Storer].
func (s *pgstorer) Get(ctx context.Context, key string, out any) (bool, error) {
	const query = "SELECT value FROM key_value WHERE key = $1 AND expiration > now();"

	key = s.prefix + key
	return s.execGet(ctx, query, out, key)
}

// GetTTL implements [Storer].
func (s *pgstorer) GetTTL(
	ctx context.Context,
	key string,
	ttl time.Duration,
	out any,
) (bool, error) {
	const query = "UPDATE key_value SET expiration = $2 " +
		"WHERE key = $1 AND expiration > now() " +
		"RETURNING value;"

	key = s.prefix + key
	return s.execGet(ctx, query, out, key, time.Now().Add(ttl))
}

// GetString implements [Storer].
func (s *pgstorer) GetString(ctx context.Context, key string) (string, error) {
	var out string
	ok, err := s.Get(ctx, key, &out)
	if err != nil || !ok {
		return "", err
	}
	return out, nil
}

// GetStringTTL implements [Storer].
func (s *pgstorer) GetStringTTL(
	ctx context.Context,
	key string,
	ttl time.Duration,
) (string, error) {
	var out string
	ok, err := s.GetTTL(ctx, key, ttl, &out)
	if err != nil || !ok {
		return "", err
	}
	return out, nil
}

// Set implements [Storer].
func (s *pgstorer) Set(ctx context.Context, key string, value any) error {
	return s.execSet(ctx, key, value, nil)
}

// SetTTL implements [Storer].
func (s *pgstorer) SetTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	return s.execSet(ctx, key, value, &ttl)
}

// Delete implements [Storer].
func (s *pgstorer) Delete(ctx context.Context, key string) error {
	const query = "DELETE FROM key_value WHERE key = $1;"

	key = s.prefix + key
	_, err := s.q.Exec(ctx, query, key)
	return err
}

// DeleteExists implements [Storer].
func (s *pgstorer) DeleteExists(ctx context.Context, key string) (bool, error) {
	const query = "DELETE FROM key_value WHERE key = $1 AND expiration > now() " +
		"RETURNING key;"

	key = s.prefix + key
	row := s.q.QueryRow(ctx, query, key)

	var out string
	if err := row.Scan(&out); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		}
		return false, err
	}

	return true, nil
}

// DeleteGet implements [Storer].
func (s *pgstorer) DeleteGet(ctx context.Context, key string, out any) (bool, error) {
	const query = "DELETE FROM key_value WHERE key = $1 AND expiration > now() " +
		"RETURNING value;"

	key = s.prefix + key
	return s.execGet(ctx, query, out, key)
}

// DeleteGetString implements [Storer].
func (s *pgstorer) DeleteGetString(ctx context.Context, key string) (string, error) {
	var out string
	ok, err := s.DeleteGet(ctx, key, &out)
	if err != nil || !ok {
		return "", err
	}
	return out, nil
}

func (s *pgstorer) execSet(ctx context.Context, key string, value any, ttl *time.Duration) error {
	const query = "INSERT INTO key_value (key, expiration, value) " +
		"VALUES ($1, $2, $3) " +
		"ON CONFLICT (key) DO UPDATE SET " +
		"key = EXCLUDED.key, " +
		"expiration = EXCLUDED.expiration, " +
		"value = EXCLUDED.value;"

	var expiration pgtype.Timestamp
	if ttl != nil {
		expiration = pgtype.Timestamp{
			Valid: true,
			Time:  time.Now().Add(*ttl),
		}
	}

	encv, err := json.Marshal(value)
	if err != nil {
		return err
	}

	key = s.prefix + key
	cmd, err := s.q.Exec(ctx, query, key, expiration, encv)
	if err != nil {
		return err
	}

	// Logically XOR
	// One (only one) of these operations must be executed
	utils.Assert(cmd.Update() != cmd.Insert())

	return nil
}

func (s *pgstorer) execGet(
	ctx context.Context,
	query string,
	out any,
	args ...any,
) (bool, error) {
	row := s.q.QueryRow(ctx, query, args...)

	var b []byte
	if err := row.Scan(&b); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	if len(b) == 0 {
		return false, nil
	}

	if err := json.Unmarshal(b, out); err != nil {
		return false, err
	}

	return true, nil
}
